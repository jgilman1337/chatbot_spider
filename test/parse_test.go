package test

import (
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
	urlA := "./data/multiple_qs.txt"
	urlE := "./data/expected_savemyphind/multiple_qs.md.txt"
	minSim := 82.5 //Some sources were missing from the scraped page that SMChatBot picked up (possibly) via duping
	parseAndCompare(t, false, topts, urlA, urlE, minSim)
}

func TestParse_Nontrivial_Single(t *testing.T) {
	urlA := "./data/nontrivial_single_q.txt"
	urlE := "./data/expected_savemyphind/nontrivial_single_q.md.txt"
	minSim := 85.0 //85% is ok; small differences more noticeable in smaller comparisons
	parseAndCompare(t, false, topts, urlA, urlE, minSim)
}

func TestParse_Single(t *testing.T) {
	urlA := "./data/single_q.txt"
	urlE := "./data/expected_savemyphind/single_q.md.txt"
	minSim := 90.0
	parseAndCompare(t, false, topts, urlA, urlE, minSim)
}

func TestParse_Two(t *testing.T) {
	urlA := "./data/two_qs.txt"
	urlE := "./data/expected_savemyphind/two_qs.md.txt"
	minSim := 90.0
	parseAndCompare(t, false, topts, urlA, urlE, minSim)
}
