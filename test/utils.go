package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hbollon/go-edlib"
	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

// Runs a fuzzy comparison between 2 strings to determine similarity.
func compare(t *testing.T, expected, actual string, minSimilarity float64) {
	percSim := runeSimilarity(expected, actual)
	if percSim < minSimilarity {
		t.Fatalf("fuzzy comparison failed between expected and actual; need: %f, got: %f", minSimilarity, percSim)
	} else {
		t.Logf("fuzzy comparison succeeded between expected and actual; need: %f, got: %f", minSimilarity, percSim)
	}
}

// Parses a scraped Perplexity page and compares it against an expected output.
func parseAndCompare(t *testing.T, quiet bool, opts perplexity.Options, actualPath, expectedPath string, minSimilarity float64) {
	//Get the archive
	dat := perplexityProvider(t, actualPath, opts)

	//Render to markdown
	actual, err := dat.RenderMD()
	if err != nil {
		t.Fatal(err)
	}
	if !quiet {
		fmt.Println("md::'''\n" + string(actual) + "'''")
	}

	//Read in the expected output
	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read 'expected' output file: %s", err)
	}

	//Compare the two files
	compare(t, string(expected), string(actual), minSimilarity)
}

// Contains common code for parser tests.
func parserTestsBackend(t *testing.T, path string, quiet bool, opts perplexity.Options) *pkg.Archive {
	//Get the archive
	dat := perplexityProvider(t, path, opts)

	//Render to markdown
	md, err := dat.RenderMD()
	if err != nil {
		t.Fatal(err)
	}
	if !quiet {
		fmt.Println("md::'''\n" + string(md) + "'''")
	}

	//Return the archive object
	return dat
}

// Contains the parser logic specific to Perplexity.ai documents.
func perplexityProvider(t *testing.T, path string, opts perplexity.Options) *pkg.Archive {
	// Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	crawler.Options = opts
	if err := crawler.FromFile(path); err != nil {
		t.Fatal(err)
	}

	// Get the metadata
	dat, err := crawler.Aggregate([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("data: %+v\n", dat)

	return dat
}

// calculateRuneSimilarity computes the percentage similarity between two strings based on rune-level comparison.
func runeSimilarity(s1, s2 string) float64 {
	//Compute Levenshtein distance
	distance := edlib.LevenshteinDistance(s1, s2)

	//Find the maximum possible length (to normalize the similarity)
	maxLen := max(len(s1), len(s2))

	//Avoid division by zero
	if maxLen == 0 {
		return 100.0 //Both strings are empty, so they are 100% similar
	}

	//Calculate percentage similarity
	similarity := 100 * (1 - float64(distance)/float64(maxLen))
	return similarity
}
