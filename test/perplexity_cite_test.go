package test

import (
	"fmt"
	"testing"

	"github.com/jgilman1337/chatbot_spider/pkg/perplexity"
)

func TestAddCitationsSingle(t *testing.T) {
	opts := perplexity.DefaultOpts()
	opts.PostProcessCitations = true
	dat := perplexityProvider(t, "./data/single_q.txt", opts)

	fmt.Printf("dat: %+v\n", dat.Thread)
}
