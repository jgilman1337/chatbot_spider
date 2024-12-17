package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestGetMetadata(t *testing.T) {
	//Setup
	url := "./data/multiple_qs.txt"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if err := crawler.FromFile(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	meta, err := crawler.GetPageMetadata()
	if err != nil {
		t.Fatal(meta)
	}
	fmt.Printf("metadata: %+v\n", meta)
}
