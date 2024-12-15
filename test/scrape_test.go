package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestScrape(t *testing.T) {
	//cSetup
	url := "https://www.perplexity.ai/search/any-way-to-do-this-more-concis-QGFOL_P.RzalvOSUiNYUdQ"
	//url := "https://www.perplexity.ai/search/hi-RVnJW4CaTaSOv4cS1uPmSQ"
	//url := "https://www.perplexity.ai/page/mirror-microbe-threat-warning-_ZXTBw9cTj6jbVaTlnz4vQ"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if _, err := crawler.Crawl(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", dat)
}
