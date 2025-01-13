package perplexity

import (
	"log"

	"github.com/creasty/defaults"
)

// Holds the options for the Perplexity parser.
type Options struct {
	//Whether to run the citation post-processor on parsed answers.
	PostProcessCitations bool `default:"false"`
}

// Gets the default options for the Perplexity crawler.
func DefaultOpts() Options {
	opts := Options{}
	err := defaults.Set(&opts)
	if err != nil {
		log.Fatal("Error setting default options:", err)
	}
	return opts
}
