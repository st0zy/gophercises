package parser

import (
	"bytes"
	"io"
	"strings"

	"github.com/st0zy/gophercises/link/link"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type LinkParser interface {
	Parse(r io.Reader) []link.HyperLink
}

type Parser struct {
	Reader io.Reader
}

func (p Parser) Parse() []link.HyperLink {

	z, err := html.Parse(p.Reader)
	if err != nil {
		panic(err)
	}
	links := make([]link.HyperLink, 0)

	for n := range z.Descendants() {
		// fmt.Println(n.DataAtom)
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			// fmt.Printf("%+v", n)
			// fmt.Println()
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, link.HyperLink{
						Href: a.Val,
						Text: extractText(n),
					})
				}
			}
		}
	}

	return links
}

func extractText(node *html.Node) string {
	var result bytes.Buffer

	if node.Type != html.ElementNode && node.Type != html.CommentNode {
		result.WriteString(node.Data)
	}

	for c := range node.ChildNodes() {
		result.WriteString(extractText(c))
	}

	return strings.TrimSpace(result.String())
}

func NewParser(r io.Reader) Parser {
	return Parser{
		Reader: r,
	}
}
