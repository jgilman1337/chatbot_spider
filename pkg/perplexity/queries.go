package perplexity

// Represents a list of all asked questions in Perplexity.
// The structure of a `Queries` is as follows (irrelevant fields are omitted for brevity):
/*
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"state": {
			"type": "object",
			"properties": {
				"queries": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"state": {
								"type": "object",
								"properties": {
									"data": {
										"type": "array",
										"items": {
											"type": "object",
											"properties": {
												"query_str": {
													"type": "string"
												},
												"related_queries": {
													"type": "array",
													"items": {
														"type": "string"
													}
												},
												"updated_datetime": {
													"type": "string"
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
*/
type Queries struct {
	State State1 `json:"state"`
}

// Lvl 1 (mutations, queries)
type State1 struct {
	Queries []QueriesInner `json:"queries"`
}

// Lvl 2 (state)
type QueriesInner struct {
	State State2 `json:"state"`
}

// Lvl 3 (data)
type State2 struct {
	Data []Data `json:"data"`
}

// Lvl 4 (query_str, related_queries, updated_datetime)
type Data struct {
	QueryStr       string   `json:"query_str"`
	RelatedQueries []string `json:"related_queries"`
	UpdatedAt      string   `json:"updated_datetime"`
}
