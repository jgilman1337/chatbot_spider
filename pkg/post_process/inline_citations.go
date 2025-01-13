package postprocess

import (
	"log"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Represents a Markdown citation processor that expands all citations in a Markdown file that are unattributed.
type InlineCitationTransformer struct {
	URLs []string

	Priority int

	AddBrackets bool
	AddSpaces   bool
}

// Creates a new `InlineCitationTransformer` object.
func NewInlineCitationTransformer(urls ...string) *InlineCitationTransformer {
	obj := InlineCitationTransformer{URLs: urls}
	obj.Priority = 0
	obj.AddBrackets = false
	obj.AddSpaces = true
	return &obj
}

// Adds this object to a Goldmark parser.
func (t *InlineCitationTransformer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(t, t.Priority),
	))
}

// Transform implements goldmark.parser.ASTTransformer.
func (t *InlineCitationTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	//Read the input bytes
	source := reader.Source()

	//Get the index of the first space char
	firstSpaceIdx := strings.Index(string(source), " ")

	//Walk the AST in depth-first fashion and apply transformations
	err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		/*
			Each node will be visited twice, once when it is first encountered
			(entering), and again after all the node's children have been visited
			(if any). Skip the latter.
		*/
		if !entering {
			return ast.WalkContinue, nil
		}

		//Mark code blocks and code fences and skip them
		if node.Kind() == ast.KindCodeBlock ||
			node.Kind() == ast.KindFencedCodeBlock ||
			node.Kind() == ast.KindCodeSpan {
			//inCodeBlock = entering
			return ast.WalkContinue, nil
		}

		//Skip the children of existing links to prevent double-transformation
		if node.Kind() == ast.KindLink || node.Kind() == ast.KindAutoLink {
			return ast.WalkSkipChildren, nil
		}

		//Check any text nodes that are encountered
		if node.Kind() == ast.KindText {
			textNode := node.(*ast.Text)
			t.attribute(textNode, source, firstSpaceIdx)
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		//TODO: use fatal slog here, if possible
		log.Fatal("Error encountered while transforming AST:", err)
	}
}

/*
Finds all unattributed citations in the given Text node and replaces them with Link
nodes that point to a URL given by the index specified in the inline citation.
*/
func (t InlineCitationTransformer) attribute(node *ast.Text, source []byte, firstSpaceIdx int) {
	//Get the text segment of the current node
	content := node.Segment.Value(source)

	//Get the previous sibling, if it exists
	prevSibUnconfirmed := node.PreviousSibling()
	if prevSibUnconfirmed == nil || prevSibUnconfirmed.Kind() != ast.KindText {
		return
	}
	prevSib := prevSibUnconfirmed.(*ast.Text)
	psContent := prevSib.Segment.Value(source)

	//Get the next sibling, if it exists
	nextSibUnconfirmed := node.NextSibling()
	if nextSibUnconfirmed == nil || nextSibUnconfirmed.Kind() != ast.KindText {
		return
	}
	nextSib := nextSibUnconfirmed.(*ast.Text)
	nsContent := nextSib.Segment.Value(source)

	//Check if the current node has an adjacent citation (eg: `1][`)
	hasAdjacent := false
	if strings.HasSuffix(string(content), "][") {
		//Mark the node as being adjacent
		hasAdjacent = true

		//Remove the closing and opening brackets from this node (eg: `1][` becomes `1`)
		node.Segment = node.Segment.WithStop(node.Segment.Stop - 2)
		content = node.Segment.Value(source)

		/*
			Shift the content of the succeeding node to the right by 1
			Successor goes from `2...` to `]`
			This allows the adjacent node to be processed by the same algo
			that is used for non-adjacent nodes
		*/
		shiftNodeBy(nextSib, 1)
		nsContent = nextSib.Segment.Value(source)
	}

	/*
		Check if the current node is "bookended" by brackets. Only
		dangling citations (eg: [1], [5], [99]) will pass this guard.
		The "ending bookend" is ignored if this node has an adjacent
		citation.
	*/
	if len(psContent) < 1 || psContent[len(psContent)-1] != '[' {
		return
	}
	if !hasAdjacent && (len(nsContent) < 1 || nsContent[0] != ']') {
		return
	}

	//Check if the current node is a valid integer and falls within the range of the links array
	idx, err := strconv.Atoi(string(content))
	if err != nil {
		return
	}
	if idx < 1 || idx-1 >= len(t.URLs) {
		//fmt.Printf("skipping index (%d) that's beyond the bounds of the links array (%d)\n", idx, len(t.URLs))
		return
	}

	//Remove the brackets off the siblings
	modifyTextNode(prevSib, prevSib.Segment.Start, prevSib.Segment.Stop-1)
	if !hasAdjacent {
		modifyTextNode(nextSib, nextSib.Segment.Start+1, nextSib.Segment.Stop)
	}

	//Add the link
	t.createLinkNode(node, t.URLs[idx-1])

	/*
		If the current node has an adjacent citation, unshift the successor node
		and truncate this node to only be the opening bracket for the successor
		node check to read.

		Otherwise, "delete" the current node by setting its start and stop indices
		to be the same value.
	*/
	if hasAdjacent {
		shiftNodeBy(nextSib, -1)
		node.Segment = node.Segment.WithStop(node.Segment.Stop + 2)
		node.Segment = node.Segment.WithStart(node.Segment.Stop - 1)
	} else {
		node.Segment = node.Segment.WithStart(node.Segment.Start)
		node.Segment = node.Segment.WithStop(node.Segment.Start)
	}

	//Check if the predecessor has a space just before it
	if t.AddSpaces && (len(psContent) <= 1 || psContent[len(psContent)-2] != ' ') {
		//Skip predecessors at the beginning of a line
		if prevSib.Segment.Start == 0 || source[prevSib.Segment.Start-1] == '\n' || source[prevSib.Segment.Start-1] == '\r' {
			return
		}

		createSpacerNode(node, firstSpaceIdx)
	}
}

// Creates a link node given a target node and the URL to point to.
func (t InlineCitationTransformer) createLinkNode(targ *ast.Text, url string) {
	//Get the text segment of the current node and its parent
	tSegment := targ.Segment
	parent := targ.Parent()

	//Create a text.Segment for the link text
	segStart := tSegment.Start
	segStop := tSegment.Stop
	if t.AddBrackets {
		segStart -= 1
		segStop += 1
	}
	lSegment := text.NewSegment(segStart, segStop)

	//Create a new link node
	link := ast.NewLink()
	link.AppendChild(link, ast.NewTextSegment(lSegment))
	link.Destination = []byte(url)
	parent.InsertBefore(parent, targ, link)
}

/*
Creates a spacer node by taking the first occurrence of a space in the source
Markdown. This is done because Goldmark's AST can only source from the source
file from which its constructed.
*/
func createSpacerNode(targ *ast.Text, spOff int) {
	//Create a new spacer node
	parent := targ.Parent()
	sp := ast.NewTextSegment(
		text.NewSegment(spOff, spOff+1),
	)

	/*
		Insert the node just before the current one; for some reason using the
		target itself pushes in the new node just after the target rather than
		before it like expected, so the target's predecessor is used as the `v1`
		node instead.
	*/
	parent.InsertBefore(parent, targ.PreviousSibling(), sp)
}

// Modifies the offsets of a text node.
func modifyTextNode(node *ast.Text, newBegin, newEnd int) {
	//Get the parent of the node
	parent := node.Parent()

	//Create the new node segment
	lSegment := text.NewSegment(newBegin, newEnd)

	//Replace the node data with the new segment
	parent.ReplaceChild(parent, node, ast.NewTextSegment(lSegment))
}

// Shifts the start and end indices of a node by x places.
func shiftNodeBy(node *ast.Text, amt int) {
	node.Segment = node.Segment.WithStart(node.Segment.Start + amt)
	node.Segment = node.Segment.WithStop(node.Segment.Stop + amt)
}
