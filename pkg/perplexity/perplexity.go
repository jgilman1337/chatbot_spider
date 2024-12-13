package perplexity

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/spider"
)

var (
	ErrorNilDoc = errors.New("document is not currently populated; run `Crawl()` first")
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
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
	//Prepare the cURL command
	cmd := exec.Command("curl",
		"-w", "%{http_code}", //This line adds the HTTP status code to the end of the output
		"-H", fmt.Sprintf("User-Agent: %s", userAgent),
		"-H", "Accept-Language: en-US,en;q=0.6",
		"-H", "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		url,
	)

	//Run the command and capture the output
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
		return nil, fmt.Errorf("HTTP error %d :: %s", statusCode, http.StatusText(statusCode))
	}

	//Load the HTML document from the bytes
	c.doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(outHtml)))
	if err != nil {
		return nil, err
	}

	return []byte(c.doc.Text()), nil
}

// Parses the raw HTML into a desired type.
func (c Crawler[T]) Aggregate(_ []byte) (*T, error) {
	//Ensure the document is non-null
	if err := c.assertNonNilDoc(); err != nil {
		return nil, err
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
