package pkg

/*
Represents a single question and answer in a chatbot thread. The data is encoded as
Markdown internally.
*/
type Question struct {
	//The initial question from the user.
	Query string `json:"query"`

	//The response from the chatbot.
	Reply Reply `json:"reply"`
}

// A thread is a list of questions, in order of when they were asked.
type Thread []Question
