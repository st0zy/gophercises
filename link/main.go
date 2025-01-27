package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/st0zy/gophercises/link/parser"
)

func main() {
	fileToParseForLinks := flag.String("file", "ex3.html", "select the file to parse for links")
	flag.Parse()

	f, err := os.Open(*fileToParseForLinks)
	if err != nil {
		panic(err)
	}

	parser := parser.NewFileParser(f)
	fmt.Printf("%+v", parser.Parse())

	// parser.Parse()
}
