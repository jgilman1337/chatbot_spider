package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestParse_Multiple(t *testing.T) {
	//Setup
	url := "./data/multiple_qs.txt"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if err := crawler.FromFile(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", dat)
}

func TestParse_Nontrivial_Single(t *testing.T) {
	//Setup
	url := "./data/nontrivial_single_q.txt"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if err := crawler.FromFile(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", dat)
}

func TestParse_Single(t *testing.T) {
	//Setup
	url := "./data/single_q.txt"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if err := crawler.FromFile(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", dat)
}

func TestParse_Two(t *testing.T) {
	//Setup
	url := "./data/two_qs.txt"

	//Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	if err := crawler.FromFile(url); err != nil {
		t.Fatal(err)
	}

	//Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %v\n", dat)
}
