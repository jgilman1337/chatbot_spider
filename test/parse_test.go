package test

import (
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

var (
	topts = func() perplexity.Options {
		opts := perplexity.DefaultOpts()
		opts.PostProcessCitations = false
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
	parserTestsBackend(t, url, false, topts)
}

func TestParse_Two(t *testing.T) {
	url := "./data/two_qs.txt"
	parserTestsBackend(t, url, false, topts)
}
