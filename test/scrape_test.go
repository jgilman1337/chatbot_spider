package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestScrape(t *testing.T) {
	//cSetup
	url := "https://www.perplexity.ai/search/ai-amplifies-false-memories-9iZN5JuFT5.9asR1Ntf._A"

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
