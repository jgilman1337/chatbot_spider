package perplexity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
	"github.com/jgilman1337/chatbot_spider/pkg"
	postprocess "github.com/jgilman1337/chatbot_spider/pkg/post_process"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
)

// Handles what is to be done when the aggregator encounters a block containing an answer.
func (c Crawler[T]) handleEncounterAnswer(cont string, ans *[]pkg.Reply) {
	//Unmarshal to an array of interfaces
	//This unescapes the target JSON data
	data := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(cont), &data); err != nil {
		fmt.Printf("answerHandler: err during 1st parse pass: %s\n", err)
		return
	}

	//Setup the target answer struct
	var answer Answer
	foundAns := false

	//Loop over the collected array items
	for _, dat := range data {
		//Skip non-strings
		item, ok := dat.(string)
		if !ok {
			continue
		}

		//Get the index of the beginning of the answer object, skipping unmatched array items
		//The answer object starts with `{"answer":`
		prefix := `{"answer":`
		idx := kmp.Search([]byte(item), []byte(prefix))
		if idx == -1 {
			continue
		}
		item = item[idx:]

		//fmt.Printf("raw ans json: ```%s```\n\n\n", item)

		//Unmarshal the answer to a struct
		if err := json.Unmarshal([]byte(item), &answer); err != nil {
			fmt.Printf("answerHandler: err during 2nd parse pass: %s\n", err)
			continue
		}

		//Answer was found; no need to continue
		foundAns = true
		break
	}

	//Parse the answer only if it's non-null
	if !foundAns {
		return
	}

	//Collect the list of source URLs
	sources := make([]pkg.Source, len(answer.WebResults))
	for i, result := range answer.WebResults {
		source := &sources[i]
		source.ID = i + 1
		source.Name = result.Name
		source.URL = result.URL
	}

	//Post-process the answer, if requested
	if c.Options.PostProcessCitations {
		//Setup the post-processor
		urls := make([]string, len(sources))
		for i, source := range sources {
			urls[i] = source.URL
		}
		citer := postprocess.NewInlineCitationTransformer(urls...)

		//Setup Goldmark with the post-processor
		gm := goldmark.New(
			goldmark.WithExtensions(
				citer,
			),
		)
		gm.SetRenderer(markdown.NewRenderer())
		buf := bytes.Buffer{}

		//Render the Goldmark AST to GFLM
		if err := gm.Convert([]byte(answer.Answer), &buf); err != nil {
			log.Fatalf("Encountered Markdown conversion error: %v", err)
		}

		//Replace the answer text with the post-processed answer
		answer.Answer = buf.String()
	}

	//Construct a generic reply object
	reply := pkg.Reply{
		Answer:  answer.Answer,
		Sources: sources,
	}

	//Add the answer to the answer array
	*ans = append(*ans, reply)
}
