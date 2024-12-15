package perplexity

import (
	"encoding/json"
	"fmt"

	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
	"github.com/jgilman1337/chatbot_spider/pkg"
)

// Handles what is to be done when the aggregator encounters a block containing an answer.
func handleEncounterAnswer(cont string, ans *[]pkg.Reply) {
	//Unmarshal to an array of interfaces
	//This unescapes the target JSON data
	data := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(cont), &data); err != nil {
		fmt.Printf("err during 1st parse pass: %s\n", err)
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

		//Unmarshal the answer to a struct
		if err := json.Unmarshal([]byte(item), &answer); err != nil {
			fmt.Printf("err during 2nd parse pass: %s\n", err)
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

	//Collect the list of source URL
	sources := make([]string, len(answer.WebResults))
	for i, result := range answer.WebResults {
		sources[i] = result.URL
	}

	//Construct a generic reply object
	reply := pkg.Reply{
		Answer:  answer.Answer,
		Sources: sources,
	}

	//Add the answer to the answer array
	*ans = append(*ans, reply)
}
