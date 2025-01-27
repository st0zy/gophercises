package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/st0zy/gophercises/urlshort/database"
	"github.com/st0zy/gophercises/urlshort/handlers"
)

func main() {
	mux := defaultMux()
	db, err := database.OpenDB("my.db")

	if err != nil {
		log.Fatal("failed to load database, skipping db mappings")
	}

	pathMappingsFromDB := database.GetAllMappings(db)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `
	yamlHandler, err := handlers.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := `[
  {
    "path": "/test",
    "url" : "https://reddit.com/"
  },
  {
    "path": "/hahaha",
    "url" : "https://google.com/"
  }
]`

	jsonHandler, err := handlers.JsonHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}
	dbHandler, err := handlers.DBHandler(pathMappingsFromDB, jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
