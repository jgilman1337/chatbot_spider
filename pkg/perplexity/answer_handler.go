package perplexity

import (
	"encoding/json"
	"fmt"

	"github.com/jgilman1337/chatbot_spider/pkg"
)

// Handles what is to be done when the aggregator encounters a block containing an answer.
// The structure of an answer is as follows (irrelevant fields are omitted for brevity):
/*
	{
	    "$schema": "http://json-schema.org/draft-07/schema#",
	    "type": "object",
	    "properties": {
	        "answer": {
	            "type": "string"
	        },
	        "web_results": {
	            "type": "array",
	            "items": {
	                "type": "object",
	                "properties": {
	                    "url": {
	                        "type": "string"
	                    }
	                }
	            }
	        },
	        "chunks": {
	            "type": "array",
	            "items": {
	                "type": "string"
	            }
	        },
	        "extra_web_results": {
	            "type": "array",
	            "items": {}
	        }
	    }
	}
*/
func handleEncounterAnswer(cont string, ans pkg.Replies) {
	//Unmarshal to an array of interfaces
	//This unescapes the target JSON data
	data := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(cont), &data); err != nil {
		fmt.Printf("err during 1st parse pass: %s\n", err)
		return
	}

	//Set up a map that will eventually contain the answer
	jsons := make(map[string]interface{})

	//Loop over the collected array items
	for _, dat := range data {
		//Skip non-strings
		item, ok := dat.(string)
		if !ok {
			continue
		}

		//Unmarshal to a map; answers are simply JSON
		if err := json.Unmarshal([]byte(item), &jsons); err != nil {
			fmt.Printf("err during 2nd parse pass: %s\n", err)
			continue
		}

		//Check for the existence of an answer
		ansText, ok := jsons["answer"]
		if !ok {
			continue
		}

		//Construct a reply object
		reply := pkg.Reply{
			Answer: ansText.(string), //A type assertion is safe; `answer` is always a string
		}
		fmt.Printf("answer: %v\n", reply)

		//The answer was found; no need to continue searching
		break
	}

	fmt.Println("\n\n")
}
