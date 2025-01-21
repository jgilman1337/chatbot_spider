package perplexity

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/spider"
)

var (
	ErrorNilDoc       = errors.New("document is not currently populated; run `Crawl()` first")
	ErrorHttpResp     = errors.New("http error")
	ErrorUnbalancedQA = errors.New("unequal question and answer array sizes")
	ErrNoQA           = errors.New("no questions or answers found")

	headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"DNT":             "1",
		"Sec-GPC":         "1",
		"Sec-Fetch-Dest":  "document",
		"Sec-Fetch-Mode":  "navigate",
		"Sec-Fetch-Site":  "cross-site",
	}
)

// Asserts that this struct conforms to the `Crawler` interface.
var _ spider.Crawler[pkg.Archive] = &Crawler[pkg.Archive]{}

// A web crawler for `https://perplexity.ai` threads that implements `spider.Crawler`.
type Crawler[T pkg.Archive] struct {
	Options Options //Crawler options.

	doc *goquery.Document //Cached version of the document.
}

// Creates a new Perplexity.ai crawler.
func NewPerplexityCrawler() *Crawler[pkg.Archive] {
	obj := Crawler[pkg.Archive]{
		Options: DefaultOpts(),
	}
	return &obj
}

// Fetches the raw HTML that is to be parsed.
func (c *Crawler[T]) Crawl(url string) ([]byte, error) {
	//Prepare the cURL command; Perplexity 403s for the "Go way" of doing this, every time, even with identical headers.
	//If someone wants to take another crack at this, have at it!
	args := make([]string, len(headers)*2)
	i := 0
	for k, v := range headers {
		args[i] = "-H"
		i++
		args[i] = fmt.Sprintf("%s: %s", k, v)
		i++
	}
	args = append(args, "-w", "%{http_code}", url) //This line adds the HTTP status code to the end of the output

	//Run the command and capture the output
	cmd := exec.Command("curl", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing cURL command: %v", err)
	}
	outStr := out.String()
	outHtml := outStr[:len(outStr)-3]

	//Get the status code (end of the output)
	statusCode, err := strconv.Atoi(outStr[len(outStr)-3:])
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, errors.Join(ErrorHttpResp,
			fmt.Errorf("%d %s", statusCode, http.StatusText(statusCode)),
		)
	}

	//Load the HTML document from the bytes
	if err := c.FromBytes([]byte(outHtml)); err != nil {
		return nil, err
	}

	return []byte(c.doc.Text()), nil
}

// Reads in a document from a byte array.
func (c *Crawler[T]) FromBytes(src []byte) error {
	var err error
	c.doc, err = goquery.NewDocumentFromReader(bytes.NewReader(src))
	if err != nil {
		return err
	}
	return nil
}

// Reads in a document from a file.
func (c *Crawler[T]) FromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := c.FromBytes(content); err != nil {
		return err
	}
	return nil
}

// Parses the raw HTML into a desired type.
func (c Crawler[T]) Aggregate(_ []byte) (*T, error) {
	//Ensure the document is non-null
	if err := c.assertNonNilDoc(); err != nil {
		return nil, err
	}

	//Create arrays for questions and answers
	questions := make([]string, 0)
	answers := make([]pkg.Reply, 0)

	/*
		Perplexity uses a JS frontend to render the content on the page, so
		scraping isn't as straightforward. Luckily, the data that is inserted
		into the rendered page is available as JSON, albeit in the form of JSON
		fragments inside of `<script>` tags.
	*/
	c.doc.Find(`script`).Each(func(_ int, s *goquery.Selection) {
		//Get the content of the script
		cont := s.Text()

		//Skip non-pushing scripts
		prefix := "self.__next_f.push("
		suffix := ")"
		scriptPrefixIdx := kmp.Search([]byte(cont), []byte(prefix))
		if scriptPrefixIdx == -1 {
			return
		}

		//Remove the prefix and ending paren
		cont = cont[scriptPrefixIdx+len(prefix) : len(cont)-len(suffix)]

		//Skip empty scripts
		if len(cont) < 1 {
			return
		}

		//TODO: multi-faceted KMP might be a good idea to use once question searches are added

		//Check if the current script content has an answer using the KMP algorithm
		//Answers begin with the following: `{\"answer\":`
		ansPrefix := `{\"answer\":`
		kmpAnsIdx := kmp.Search([]byte(cont), []byte(ansPrefix))
		if kmpAnsIdx != -1 {
			if err := c.handleEncounterAnswer(cont, &answers); err != nil {
				log.Fatalf("error while processing answer: %s", err)
			}
		}

		//Check if the current script content has a question using the KMP algorithm
		//Questions begin with the following: `{\"answer\":`
		quesPrefix := `\"queries\":[`
		kmpQuesIdx := kmp.Search([]byte(cont), []byte(quesPrefix))
		if kmpQuesIdx != -1 {
			//Skip matches with empty query arrays (false positives)
			queryArrBegin := cont[kmpQuesIdx+len(quesPrefix)]
			if queryArrBegin == ']' {
				return
			}

			//Parse out the questions from the found queries block
			if err := c.handleEncounterQuestion(cont, &questions); err != nil {
				log.Fatalf("error while processing question: %s", err)
			}
		}
	})

	//Ensure both arrays have entires
	if len(questions) < 1 || len(answers) < 1 {
		return nil, fmt.Errorf("%w; %d questions, %d answers", ErrNoQA,
			len(questions), len(answers),
		)
	}

	//At this point, the question and answer arrays are fully filled
	//It is assumed that each question has a corresponding answer, so reject if this isn't the case
	if len(questions) != len(answers) {
		return nil, errors.Join(ErrorUnbalancedQA,
			fmt.Errorf("; qs[%d], as[%d]", len(questions), len(answers)),
		)
	}

	//Build up a thread of questions and answers
	thread := make([]pkg.Question, len(questions))
	for i := range thread {
		threadItem := &thread[i]
		*threadItem = pkg.Question{
			Query: questions[i],
			Reply: answers[i],
		}
	}

	//Get the metadata of the page
	meta, err := c.GetPageMetadata()
	if err != nil {
		return nil, err
	}

	//Compile and return the archive object
	out := T(pkg.Archive{
		Metadata: *meta,
		Thread:   thread,
	})
	return &out, nil
}

// Fetches the metadata of the thread.
func (c Crawler[T]) GetPageMetadata() (*pkg.Metadata, error) {
	//Ensure the document is non-null
	if err := c.assertNonNilDoc(); err != nil {
		return nil, err
	}

	//Get the title of the thread
	title := c.doc.Find(`title`).Text()
	title = truncateTitle(title)

	//Get the URL and other attributes
	url, _ := c.doc.Find(`link[rel="canonical"]`).Attr("href")
	created, _ := c.doc.Find(`meta[name="datePublished"]`).Attr("content")

	//Parse the creation date to a Go time object
	createdTime, err := dateparse.ParseAny(created)
	if err != nil {
		return nil, fmt.Errorf("error when parsing creation time `%s`: %s", created, err)
	}

	//Construct the final metadata object
	meta := pkg.Metadata{
		Title:    title,
		Service:  "Perplexity.ai",
		URL:      url,
		Created:  createdTime,
		Archived: time.Now(),
	}

	return &meta, nil
}

// Asserts that a cached document is non-null.
func (c Crawler[T]) assertNonNilDoc() error {
	if c.doc == nil {
		return ErrorNilDoc
	} else {
		return nil
	}
}

// Truncates a thread title to the first newline or 100th character, whichever comes first.
func truncateTitle(title string) string {
	//Find the index of the first newline
	firstNewline := strings.Index(title, "\n")

	//Determine the truncation index
	truncateAt := 100 // Default to 100 characters
	if firstNewline != -1 && firstNewline < truncateAt {
		truncateAt = firstNewline
	}

	//If the title is shorter than the truncation index, return it as is
	if len(title) <= truncateAt {
		return title
	}

	//Truncate the title and add ellipsis
	return title[:truncateAt] + "..."
}
