package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	mapping "github.com/st0zy/gophercises/urlshort/mapping"
	"gopkg.in/yaml.v2"
)

var ParseError = errors.New("failed to parse")

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	mappings, err := parseYaml(yml)
	if err != nil {
		return nil, ParseError
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapping := range mappings {
			if mapping.Path == r.URL.Path {
				http.Redirect(w, r, mapping.RedirectPath, http.StatusPermanentRedirect)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil

}

func parseYaml(yml []byte) ([]mapping.PathMapping, error) {
	fmt.Println("Parsing YAMl")
	var mapping []mapping.PathMapping
	err := yaml.Unmarshal(yml, &mapping)
	fmt.Println(err)
	return mapping, err
}

func parseJson(yml []byte) ([]mapping.PathMapping, error) {
	fmt.Println("Parsing YAMl")
	var mapping []mapping.PathMapping
	err := json.Unmarshal(yml, &mapping)
	fmt.Println(err)
	return mapping, err
}

func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {

	mappings, err := parseJson(json)
	if err != nil {
		return nil, ParseError
	}
	fmt.Println(mappings)

	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapping := range mappings {
			if mapping.Path == r.URL.Path {
				http.Redirect(w, r, mapping.RedirectPath, http.StatusPermanentRedirect)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil

}

func DBHandler(pathMappings []mapping.PathMapping, fallback http.Handler) (http.HandlerFunc, error) {

	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapping := range pathMappings {
			if mapping.Path == r.URL.Path {
				http.Redirect(w, r, mapping.RedirectPath, http.StatusPermanentRedirect)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil

}
