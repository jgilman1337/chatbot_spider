package pkg

import "time"

// Represents an archived chatbot thread.
type Archive struct {
	Metadata

	//The payload of the archive.
	Thread Thread `json:"thread"`
}

// Represents the metadata of chatbot archive.
type Metadata struct {
	//The title of the thread (usually a truncated version of the opening question).
	Title string `json:"title"`

	//The name of the service from which the data was sourced.
	Service string `json:"service"`

	//The URL of the page that was archived.
	URL string `json:"url"`

	//The time at which the thread was created.
	Created time.Time `json:"created"`

	//The time at which the thread was archived.
	Archived time.Time `json:"archived"`
}
