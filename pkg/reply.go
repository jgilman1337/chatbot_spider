package pkg

/*
Represents a singular reply to a question.Replies include the answer text and
the relevant sources.
*/
type Reply struct {
	//The answer from the chatbot.
	Answer string `json:"answer"`

	//The list of source URLs.
	Sources []Source `json:"sources"`
}

// Represents a list of replies.
type Replies []Reply

// Represents a single reply source, which has an ID, name, and URL.
type Source struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
