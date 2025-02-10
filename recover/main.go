package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMw(mux, true)))
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>Panic: %v </h1> <pre> %s</pre>", r, string(stack))
			}
		}()
		nw := &newResponseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

// type ResponseWriter interface {
// 	Header() Header

// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type newResponseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (nw *newResponseWriter) Write(bytes []byte) (int, error) {
	nw.writes = append(nw.writes, bytes)
	return len(bytes), nil
}

func (nw *newResponseWriter) WriteHeader(status int) {
	nw.status = status
}

func (nw *newResponseWriter) flush() {
	if nw.status != 0 {
		nw.ResponseWriter.WriteHeader(nw.status)
	}
	for _, write := range nw.writes {
		nw.ResponseWriter.Write(write)
	}
}
