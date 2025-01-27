package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/st0zy/gophercises/link/parser"
)

type SiteMapBuilder struct {
	site    string
	visited map[string]bool
}

func NewSiteMapBuilder(site string) (*SiteMapBuilder, error) {
	_, err := url.ParseRequestURI(site)
	if err != nil {
		return nil, errors.New("failed to parse the url")
	}
	response, err := http.Get(site)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("please enter a valid site.")
	}
	return &SiteMapBuilder{
		site:    site,
		visited: make(map[string]bool),
	}, nil
}

type Site struct {
	URL string
}

type Sites []Site

func (s SiteMapBuilder) BuildSiteMap() Sites {
	queue := make([]string, 0)
	queue = append(queue, s.site)
	s.visited[s.site] = true
	sites := s.bfs(queue)
	return sites
}

func (s SiteMapBuilder) bfs(queue []string) Sites {
	var sites Sites
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		sites = append(sites, Site{current})
		base, err := url.Parse(current)
		if err != nil {
			continue
		}
		response, err := http.Get(base.String())
		if err != nil {
			fmt.Printf("Skipping %s as it isn't reachable", current)
			continue
		}
		links := parser.NewParser(response.Body).Parse()
		// fmt.Println(links)
		for _, link := range links {
			hrefs, _ := url.Parse(link.Href)
			var absoluteUrl *url.URL
			if strings.HasPrefix(hrefs.String(), "/") {
				absoluteUrl = base.ResolveReference(hrefs)
			} else if strings.HasPrefix(hrefs.String(), base.Scheme+"://"+base.Host) {
				absoluteUrl = hrefs
			} else {
				continue
			}
			if _, ok := s.visited[absoluteUrl.String()]; !ok {
				// fmt.Println(absoluteUrl)
				queue = append(queue, absoluteUrl.String())
				s.visited[absoluteUrl.String()] = true
			}
		}
	}

	return sites
}

type SiteMapPrinter struct {
	writer io.Writer
	sites  Sites
}

func (s SiteMapPrinter) Write() error {
	bufferedWriter := bufio.NewWriter(s.writer)
	tpl := template.Must(template.New("").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{range .}}
  <url>
    <loc>{{.URL}}</loc>
  </url>
  {{end}}
</urlset>`))

	defer bufferedWriter.Flush()
	tpl.Execute(bufferedWriter, s.sites)
	return nil

}
