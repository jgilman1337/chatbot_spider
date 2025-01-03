package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestParse_Multiple(t *testing.T) {
	url := "./data/multiple_qs.txt"
	parserTestsBackend(t, url)
}

func TestParse_Nontrivial_Single(t *testing.T) {
	url := "./data/nontrivial_single_q.txt"
	parserTestsBackend(t, url)
}

func TestParse_Single(t *testing.T) {
	url := "./data/single_q.txt"
	parserTestsBackend(t, url)
}

func TestParse_Two(t *testing.T) {
	url := "./data/two_qs.txt"
	parserTestsBackend(t, url)
}

func parserTestsBackend(t *testing.T, url string) *pkg.Archive {
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
	//fmt.Printf("data: %+v\n", dat)

	//Render to markdown
	md, err := dat.RenderMD()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("md::```\n" + string(md) + "\n```")

	//Return the metadata object
	return dat
}
