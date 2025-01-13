package test

import (
	"fmt"
	"testing"

	"github.com/hbollon/go-edlib"
	"github.com/jgilman1337/chatbot_spider/pkg"
	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

// Contains common code for parser tests.
func parserTestsBackend(t *testing.T, url string, quiet bool, opts perplexity.Options) *pkg.Archive {
	//Get the archive
	dat := perplexityProvider(t, url, opts)

	//Render to markdown
	md, err := dat.RenderMD()
	if err != nil {
		t.Fatal(err)
	}
	if !quiet {
		fmt.Println("md::```\n" + string(md) + "\n```")
	}

	//Return the metadata object
	return dat
}

// Contains the parser logic specific to Perplexity.ai documents.
func perplexityProvider(t *testing.T, url string, opts perplexity.Options) *pkg.Archive {
	// Do the initial crawl
	crawler := perplexity.NewPerplexityCrawler()
	crawler.Options = opts
	if err := crawler.FromFile(url); err != nil {
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
