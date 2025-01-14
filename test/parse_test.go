package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

var (
	topts = func() perplexity.Options {
		opts := perplexity.DefaultOpts()
		opts.PostProcessCitations = true
		return opts
	}()
)

func TestParse_Multiple(t *testing.T) {
	url := "./data/multiple_qs.txt"
	parserTestsBackend(t, url, false, topts)
}

func TestParse_Nontrivial_Single(t *testing.T) {
	url := "./data/nontrivial_single_q.txt"
	parserTestsBackend(t, url, false, topts)
}

func TestParse_Single(t *testing.T) {
	url := "./data/single_q.txt"
	dat := parserTestsBackend(t, url, true, topts)

	txt, err := dat.RenderHTML()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("txt: ```%s```\n", string(txt))
}

func TestParse_Two(t *testing.T) {
	url := "./data/two_qs.txt"
	parserTestsBackend(t, url, false, topts)
}
