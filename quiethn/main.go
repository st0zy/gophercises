package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/st0zy/gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {

	sc := storyCache{
		numStories: numStories,
		duration:   3 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(sc.duration)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   sc.duration,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	cache      []item
	numStories int
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if time.Now().Before(sc.expiration) {
		// fmt.Println("Returning cached stories.")
		return sc.cache, nil
	}
	var err error
	sc.cache, err = getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.expiration = time.Now().Add(time.Second * 1)
	return sc.cache, err
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	var stories []result
	resultCh := make(chan result)
	done := make(chan bool)

	for i := 0; i < len(ids); i++ {
		go func(index, id int) {
			var client hn.Client
			item, err := client.GetItem(id)
			if err != nil {
				return
			}
			resultCh <- result{
				story: parseHNItem(item),
				index: index,
			}
		}(i, ids[i])
	}

	go func() {
		for item := range resultCh {
			if !isStoryLink(item.story) {
				continue
			}
			stories = append(stories, item)
			if len(stories) >= len(ids) {
				done <- true
				return
			}
		}
	}()

	<-done
	sort.Slice(stories, func(i, j int) bool { return stories[i].index < stories[j].index })
	var result []item
	for _, story := range stories {
		result = append(result, story.story)
	}

	return result
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type result struct {
	index int
	story item
}
