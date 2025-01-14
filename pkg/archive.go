package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jgilman1337/chatbot_spider/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	ProgIdent = "Chatbot Spider"
	ProgURL   = "https://github.com/jgilman1337/chatbot_spider"

	QuestionHeader = "User"
	ReplyHeader    = "AI Answer"
	SourcesHeader  = "---\n**Sources:**\n"
	HeaderBreak    = "\n\n\n"
	QuestionBreak  = "\n\n"
	SourcesBreak   = QuestionBreak
	ThreadBreak    = "\n\n\n\n"
	EndSentinel    = "\n"
)

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

// Renders an `Archive` object to HTML using Goldmark.
func (a Archive) RenderHTML() ([]byte, error) {
	//Get the markdown render of the object
	md, err := a.RenderMD()
	if err != nil {
		return nil, err
	}

	//Setup Goldmark
	goldmark := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	//Render the document
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, err
	}

	//Create the document header
	header := strings.ReplaceAll(
		fmt.Sprintf(`
		<head>
			<title>%s :: %s</title>
		</head>
		`,
			a.Title,
			ProgIdent,
		),
		"\t",
		"",
	)

	//Add the missing `<html>` tag
	html := "<html>" + header + "<body>\n" + buf.String() + "</body></html>\n"

	return []byte(html), nil
}

// Renders an `Archive` object to Markdown. Based off SaveMyChatbot's format.
func (a Archive) RenderMD() ([]byte, error) {
	//Create the metadata header
	mheader := fmt.Sprintf("Exported on %s [from %s](%s) - with [%s](%s)",
		a.Metadata.Archived.Format(time.DateTime),
		a.Metadata.Service,
		a.Metadata.URL,
		ProgIdent,
		ProgURL,
	)

	//Create the result string and add the header
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(
		"# %s\n%s%s",
		a.Metadata.Title,
		mheader,
		HeaderBreak,
	))

	//Add the questions and answers from the thread
	for i, question := range a.Thread {
		//Add the thread break
		if i > 0 {
			builder.WriteString(ThreadBreak)
		}

		//Add the question
		builder.WriteString(fmt.Sprintf("## %s\n%s%s",
			QuestionHeader,
			question.Query,
			QuestionBreak,
		))

		//Add the reply, sans any sources
		builder.WriteString(fmt.Sprintf("## %s\n%s",
			ReplyHeader,
			question.Reply.Answer,
		))

		//Add the sources
		if len(question.Reply.Sources) > 0 {
			builder.WriteString(SourcesBreak + SourcesHeader)
			for i, source := range question.Reply.Sources {
				if i > 0 {
					builder.WriteString("\n")
				}

				builder.WriteString(fmt.Sprintf("- [(%d) %s](%s)",
					i+1,
					util.EscapeMD(source.Name),
					source.URL,
				))
			}
		}
	}

	//Build and return the Markdown
	return []byte(builder.String() + EndSentinel), nil
}
