package main

import (
	"log"

	"github.com/webCrawler/internal"
)

func main() {
	scraper, err := internal.NewScraper("https://www.monzo.com")
	if err != nil {
		log.Println(err)
	}

	scraper.Scrape()
}
