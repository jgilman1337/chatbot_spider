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
	pathA := "./data/multiple_qs.txt"
	pathE := "./data/expected_savemyphind/multiple_qs.md.txt"
	minSim := 82.5 //Some sources were missing from the scraped page that SMChatBot picked up (possibly) via duping
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}

func TestParse_Multiple2025(t *testing.T) {
	pathA := "./data/multiple_qs_2025.html.txt"
	pathE := "./data/expected_savemyphind/multiple_qs_2025.md.txt"
	minSim := 82.5
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}

func TestParse_Nontrivial_Single(t *testing.T) {
	pathA := "./data/nontrivial_single_q.txt"
	pathE := "./data/expected_savemyphind/nontrivial_single_q.md.txt"
	minSim := 85.0 //85% is ok; small differences more noticeable in smaller comparisons
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}

func TestParse_Single(t *testing.T) {
	pathA := "./data/single_q.txt"
	pathE := "./data/expected_savemyphind/single_q.md.txt"
	minSim := 90.0
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}

func TestParse_Two(t *testing.T) {
	pathA := "./data/two_qs.txt"
	pathE := "./data/expected_savemyphind/two_qs.md.txt"
	minSim := 90.0
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}

//-----------------------------------------------------
// 3p thread tests

func TestParse_FalseAssump(t *testing.T) {
	pathA := "./data/3p/false_mem.html.txt"
	pathE := "./data/expected_savemyphind/3p/false_mem.md.txt"
	minSim := 90.0
	parseAndCompare(t, false, topts, pathA, pathE, minSim)
}
