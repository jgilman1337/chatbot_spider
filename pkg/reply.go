package pkg

/*
Represents a singular reply to a question.Replies include the answer text and
the relevant sources.
*/
type Reply struct {
	//The answer from the chatbot.
	Answer string

	//The list of source URLs.
	Sources []string
}

// Represents a list of replies.
type Replies []Reply
