package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/st0zy/gophercises/cyoa"
)

func main() {
	var storyFile string
	port := flag.Int("port", 8080, "The port on which the server runs.")
	flag.StringVar(&storyFile, "story", "gopher.json", "The path from where the story is read")
	flag.Parse()
	file, err := os.Open(storyFile)
	if err != nil {
		fmt.Println(err)
		panic("failed to read story.")
	}
	story, err := cyoa.NewStoryFromReader(file)
	if err != nil {
		panic("something went wrong.")
	}
	mux := http.NewServeMux()
	// var tempTemplate = template.Must(template.New("").Parse("Hello!"))
	handler := cyoa.NewHandler(story, cyoa.WithPathFn(customPathFn()))
	mux.Handle("/story/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func customPathFn() func(r *http.Request) string {
	return func(r *http.Request) string {
		path := strings.TrimSpace(r.URL.Path)
		if path == "/story" || path == "/story/" {
			path = "/story/intro"
		}
		path = path[len("/story/"):]
		return path
	}
}
