package pkg

/*
Represents a single question and answer in a chatbot thread. The data is encoded as
Markdown internally.
*/
type Question struct {
	//The initial question from the user.
	Query []string

	//The response from the chatbot.
	Response []string

	//The list of raw URLs that the chatbot used as sources.
	Sources []string
}

// A thread is a list of questions, in order of when they were asked.
type Thread []Question
