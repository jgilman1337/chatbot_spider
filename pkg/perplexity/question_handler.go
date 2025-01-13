package perplexity

import (
	"encoding/json"
	"fmt"

	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
)

// Handles what is to be done when the aggregator encounters a block containing a question
func (c Crawler[T]) handleEncounterQuestion(cont string, ques *[]string) {
	//Unmarshal to an array of interfaces
	//This unescapes the target JSON data
	data := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(cont), &data); err != nil {
		fmt.Printf("questionHandler: err during 1st parse pass: %s\n", err)
		return
	}

	//Setup the target queries struct
	var queries Queries
	foundQueries := false

	//Loop over the collected array items
	for _, dat := range data {
		//Skip non-strings
		item, ok := dat.(string)
		if !ok {
			continue
		}

		//Get the index of the beginning of the question object, skipping unmatched array items
		//The question object starts with `{"state":{`
		prefix := `{"state":{`
		suffix := "]"
		idx := kmp.Search([]byte(item), []byte(prefix))
		if idx == -1 {
			continue
		}
		item = item[idx : len(item)-len(suffix)-1]

		//fmt.Printf("raw ques json: ```%s```\n\n\n", item)

		//Unmarshal the queries object to a struct
		if err := json.Unmarshal([]byte(item), &queries); err != nil {
			fmt.Printf("questionHandler: err during 2nd parse pass: %s\n", err)
			continue
		}

		//Queries object was found; no need to continue
		foundQueries = true
		break
	}

	//Parse the questions only if it's non-null
	if !foundQueries {
		return
	}

	//Add the collected queries to the output array
	questions := queries.State.Queries[0].State.Data
	for _, question := range questions {
		qtext := question.QueryStr
		*ques = append(*ques, qtext)
	}
}
