package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var site = flag.String("site", "https://www.scrapethissite.com/pages/", "enter the site to scrape")

	siteMapBuilder, err := NewSiteMapBuilder(*site)
	if err != nil {
		panic(err)
	}
	sites := siteMapBuilder.BuildSiteMap()

	siteMapPrinter := &SiteMapPrinter{
		writer: os.Stdout,
		sites:  sites,
	}
	err = siteMapPrinter.Write()
	if err != nil {
		log.Fatal(err)
	}
}
