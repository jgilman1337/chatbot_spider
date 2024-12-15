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
	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/spider"
)

var (
	ErrorNilDoc   = errors.New("document is not currently populated; run `Crawl()` first")
	ErrorHttpResp = errors.New("http error")

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
	doc *goquery.Document //Cached version of the document.
}

// Creates a new Perplexity.ai crawler.
func NewPerplexityCrawler() *Crawler[pkg.Archive] {
	return &Crawler[pkg.Archive]{}
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

	//TODO: temp
	/*
		// Create a new file with a timestamp in its name
		timestamp := time.Now().Unix()
		fileName := fmt.Sprintf("data_%d.txt", timestamp)
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	*/

	//Create arrays for questions and answers
	//questions := make([]string, 0)
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
		//TODO: might want to use KMP here too
		prefix := "self.__next_f.push("
		suffix := ")"
		if !strings.HasPrefix(cont, prefix) {
			return
		}
		//Remove the prefix and ending paren
		cont = cont[len(prefix) : len(cont)-len(suffix)]

		//Skip empty scripts
		if len(cont) < 1 {
			return
		}

		//Check if the current script content has an answer using the KMP algorithm
		//Answers begin with the following: `{\"answer\":`
		//TODO: multi-faceted KMP might be a good idea to use once question searches are added
		ansPrefix := `{\"answer\":`
		kmpIdx := kmp.Search([]byte(cont), []byte(ansPrefix))
		if kmpIdx != -1 {
			//fmt.Printf("ans: %s\n", cont)
			handleEncounterAnswer(cont, &answers)
		}

		//fmt.Printf("s: %s\n", cont)
		//idx := util.If(len(cont) <= 200, len(cont), 200)
		//fmt.Printf("s: %s\n", cont[:idx])

		// Write the content to the file instead of stdout
		/*
			if _, err := file.WriteString(fmt.Sprintf("s: %s\n\n", cont)); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
		*/
	})

	fmt.Printf("answers found: %d\n", len(answers))
	for i, answer := range answers {
		fmt.Printf("answer #%d: %v\n", i+1, answer)
	}

	return nil, nil
}

// Fetches the metadata of the thread.
func (c Crawler[T]) GetPageMetadata() (*pkg.Metadata, error) {
	//Ensure the document is non-null
	if err := c.assertNonNilDoc(); err != nil {
		return nil, err
	}

	//Get the title of the thread
	title, _ := c.doc.Find(`meta[name="twitter:title"]`).Attr("content")

	//Get the time at which the thread was created

	//Construct the final metadata object
	meta := pkg.Metadata{
		Title: title,
		//Created: time.Time,
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
