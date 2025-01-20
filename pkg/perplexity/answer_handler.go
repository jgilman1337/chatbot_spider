package perplexity

import (
	"bytes"
	"encoding/json"
	"fmt"

	llq "github.com/emirpasic/gods/v2/queues/linkedlistqueue"
	kmp "github.com/fbonhomm/knuth-morris-pratt/source"
	"github.com/jgilman1337/chatbot_spider/pkg"
	postprocess "github.com/jgilman1337/chatbot_spider/pkg/post_process"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
)

/*
Handles what is to be done when the aggregator encounters a block containing an answer.
Due to the way that some Perplexity responses are structured (sometimes appear as a root
object, other times nested in a `step > content` object), a BFS-style parse is used to
get the proper answer object. See: https://john-ahn.medium.com/breadth-first-search-e86a4cbc6803
*/
func (c Crawler[T]) handleEncounterAnswer(cont string, ans *[]pkg.Reply) error {
	//Setup the target answer struct
	var answer Answer
	foundAns := false

	//Unmarshal to an array of interfaces
	//This unescapes the target JSON data
	data := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(cont), &data); err != nil {
		return fmt.Errorf("answerHandler: err during 1st parse pass: %w", err)
	}

	//Loop over the collected array items
	for i, dat := range data {
		if i > 0 {
			fmt.Print("\n\n")
		}
		stritem := fmt.Sprintf("%v", dat)
		fmt.Printf("child item %d: '''%s'''\n",
			i+1,
			stritem[:min(len(stritem), 100)],
		)
		fmt.Printf("type: %T\n", dat)

		//Skip non-strings
		sdat, ok := dat.(string)
		if !ok {
			fmt.Println("skipped!")
			continue
		}
		//fmt.Printf("stritem: %s\n", item)

		//Attempt to parse the current entry to an array of untyped maps
		//Most Perplexity response payloads will not match this, but some will
		nestedDatArr := make([]map[string]interface{}, 1)
		err := json.Unmarshal([]byte(sdat), &nestedDatArr)
		if err == nil {
			//Parse the answer from the nested data
			fmt.Printf("hit #%d is nested\n", i+1)

			ok, err = parseAnsNestedMap(&answer, nestedDatArr)
			if ok && err == nil {
				//No need to continue searching if an answer was found
				foundAns = true
				break
			}
		}

		//Answer is encoded in a plain string; use that method of parsing
		ok, err = parseAnsObj(&answer, sdat)
		if ok && err == nil {
			//No need to continue searching if an answer was found
			foundAns = true
			break
		}
	}
	fmt.Print("\n--------------\n\n")

	//Parse the answer only if it's non-null
	if !foundAns {
		return nil
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
			return fmt.Errorf("encountered Markdown conversion error: %w", err)
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

	return nil
}

/*
Parses answer objects when they are not nested in a map, but exist as a string (most
common encoding type).
*/
func parseAnsObj(ans *Answer, sdata string) (bool, error) {
	//Get the index of the beginning of the answer object, skipping unmatched array items
	//The answer object starts with `{"answer":`
	prefix := `{"answer":`
	idx := kmp.Search([]byte(sdata), []byte(prefix))
	if idx == -1 {
		return false, nil
	}
	sdata = sdata[idx:]

	//Unmarshal the answer to a struct
	if err := json.Unmarshal([]byte(sdata), &ans); err != nil {
		return false, fmt.Errorf("answerHandler: err during 2nd parse pass (non-nested): %w", err)
	}

	//A result was found with no errors
	return true, nil
}

/*
Parses answer objects when they are nested in a map (sometimes possible). Uses a BFS on
each nested map object to find the answer.
*/
func parseAnsNestedMap(ans *Answer, sdata []map[string]interface{}) (bool, error) {
	//Declare a function that enqueues all array items into a llq
	enqueueArrMap := func(target *llq.Queue[interface{}], items []map[string]interface{}) {
		for _, item := range items {
			target.Enqueue(item)
		}
	}

	//Create a queue for the BFS parse
	bfsq := llq.New[interface{}]()
	enqueueArrMap(bfsq, sdata)

	//Run the BFS over the collected maps
	for bfsq.Size() > 0 {
		//Get the current item from the queue
		//A type assertion is ok to do here
		utitem, _ := bfsq.Dequeue()
		item, _ := utitem.(map[string]interface{})

		//Loop over the keys and values of the current item
		for k, rv := range item {
			//Switch over the type of the item
			switch v := rv.(type) {
			//Embedded object; enqueue it into the BFS array
			case map[string]interface{}:
				bfsq.Enqueue(v)

			//String type; if the owner key is `answer` then parse the answer struct from here and quit
			case string:
				//Check if the current key is `answer`
				if k == "answer" {
					//Parse out the answer object and bail
					return parseAnsObj(ans, v)
				}

			//Other types; skip the current kay-value pair
			default:
				continue
			}
		}
	}

	//No hits and no error, so return false and nil
	return false, nil
}
