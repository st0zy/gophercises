package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var defaultHandlerTempalte = `<html>
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <title>Choose your own Adventure</title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link rel="stylesheet" href="">
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraph}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
            <li><a href="/story/{{.Next}}">{{.Text}}</a></li>
			{{end}}
        </ul>
        <script src="" async defer></script>
    </body>
</html>`

var errorHandler = `<html>
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <title>Choose your own Adventure</title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link rel="stylesheet" href="">
    </head>
    <body>
        <h1>Page Not Found</h1>
        <script src="" async defer></script>
    </body>
</html>`

var defaultPathFn = func(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	return path
}

type Story map[string]Chapter

type Chapter struct {
	Title     string    `json:"title"`
	Paragraph []string  `json:"story"`
	Options   []Options `json:"options"`
}

type Options struct {
	Text string `json:"text"`
	Next string `json:"arc"`
}

type StoryHandler struct {
	s      Story
	tpl    *template.Template
	errTpl *template.Template
	pathFn func(r *http.Request) string
}

type HandlerOpts func(h *StoryHandler)

func NewStoryFromReader(r io.Reader) (Story, error) {

	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return story, nil

}

func (handler StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := handler.pathFn(r)
	if chapter, ok := handler.s[path]; ok {
		handler.tpl.Execute(w, chapter)
	} else {
		w.WriteHeader(http.StatusNotFound)
		handler.errTpl.Execute(w, struct{}{})
	}
}

func WithTemplate(t *template.Template) HandlerOpts {
	return func(h *StoryHandler) {
		if t != nil {
			h.tpl = t
		}
	}
}

func WithErrorTemplate(t *template.Template) HandlerOpts {
	return func(h *StoryHandler) {
		if t != nil {
			h.tpl = t
		}
	}
}

func WithPathFn(pathfn func(r *http.Request) string) HandlerOpts {
	return func(h *StoryHandler) {
		h.pathFn = pathfn
	}
}

func NewHandler(s Story, opts ...HandlerOpts) http.Handler {
	storyHandler := StoryHandler{
		s:      s,
		tpl:    template.Must(template.New("").Parse(defaultHandlerTempalte)),
		errTpl: template.Must(template.New("").Parse(errorHandler)),
		pathFn: defaultPathFn,
	}
	for _, option := range opts {
		option(&storyHandler)
	}
	return storyHandler
}
