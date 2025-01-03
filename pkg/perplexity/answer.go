package perplexity

/*
Represents an answer to a question in Perplexity. This object is only to be used
when unmarshalling a Perplexity answer from an answer payload. The structure of
an answer is as follows (irrelevant fields are omitted for brevity):

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
						"name": {
							"type": "string",
						},
						"url": {
							"type": "string",
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
type Answer struct {
	Answer     string       `json:"answer"`
	WebResults []WebResults `json:"web_results"`
	Chunks     []string     `json:"chunks"`
}

// Represents the inner web results type in a Perplexity answer.
type WebResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
