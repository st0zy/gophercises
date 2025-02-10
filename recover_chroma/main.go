package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)

	mux.HandleFunc("/debug/", fileServe)

	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, parseAndFormatStack(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

var filePathWithLineNumberRegex = regexp.MustCompile(`([/].*[.].*):(\d+)`)

func parseAndFormatStack(stackTrace string) string {
	return filePathWithLineNumberRegex.ReplaceAllString(stackTrace, `<a href=/debug?path=$1&line=$2>$1:$2</a>`)
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func fileServe(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query()["path"][0]
	lineNumbers := r.URL.Query()["line"]
	var lineNumber int
	if len(lineNumbers) == 0 {
		lineNumber = 0
	} else {
		lineNumber, _ = strconv.Atoi(lineNumbers[0])
	}
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	writeFormatted(file, w, lineNumber)
}

func writeFormatted(reader io.Reader, w http.ResponseWriter, lineToHighlight int) {
	lexer := lexers.Get("go")
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithLineNumbers(true), html.LineNumbersInTable(true), html.HighlightLines([][2]int{{lineToHighlight, lineToHighlight}}))
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, reader)
	contents, _ := io.ReadAll(buffer)
	iterator, _ := lexer.Tokenise(nil, string(contents))
	err := formatter.Format(w, style, iterator)
	w.Header().Set("Content-Type", "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
