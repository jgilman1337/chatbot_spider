package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestGetMetadata(t *testing.T) {
	//Setup
	url := "https://www.perplexity.ai/search/how-long-does-perplexity-store-nBUeH74UQ8KApfeFIct7sw"
	//url := "https://www.perplexity.ai/page/mirror-microbe-threat-warning-_ZXTBw9cTj6jbVaTlnz4vQ"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if _, err := crawler.Crawl(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	meta, err := crawler.GetPageMetadata()
	if err != nil {
		t.Fatal(meta)
	}
	fmt.Printf("metadata: %v\n", meta)
}
