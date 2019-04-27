package internal

import (
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type scraper struct {
	hostName    string
	scheme      string
	urlChan     chan string
	resultChan  chan string
	visitedUrls sync.Map
}

func NewScraper(rootUrl string) (*scraper, error) {
	parsedUrl, err := neturl.Parse(rootUrl)
	log.Println(parsedUrl.EscapedPath())
	if err != nil {
		return nil, err
	}
	if parsedUrl.Hostname() == "" {
		return nil, errors.New("the given URL should have a domain part")
	}
	s := &scraper{
		hostName: parsedUrl.Hostname(),
		scheme:   parsedUrl.Scheme,
		urlChan:  make(chan string, 1000),
	}
	s.urlChan <- rootUrl
	return s, nil
}

func (s *scraper) Scrape() {
	s.process()
}

func (s *scraper) processElement(_ int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	if exists && s.hasSameDomain(href) {
		if s.isRelativePath(href) {
			href = s.prependDomain(href)
		}
		if !s.isScraped(href) {
			s.addToVisited(href)
			s.urlChan <- href
		}
	}
}

func (s *scraper) hasSameDomain(href string) bool {
	parsedUrl, err := neturl.Parse(href)
	if err != nil {
		return false
	}
	if !parsedUrl.IsAbs() {
		return true
	}
	return parsedUrl.Hostname() == s.hostName
}

func (s *scraper) isRelativePath(url string) bool {
	if parsedUrl, err := neturl.Parse(url); err == nil {
		return !parsedUrl.IsAbs()
	}
	return false
}

func (s *scraper) isScraped(url string) bool {
	_, found := s.visitedUrls.Load(url)
	return found
}

func (s *scraper) prependDomain(url string) string {
	return s.scheme + "://" + s.hostName + url
}

func (s *scraper) scrape(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Println("error", err)
		return
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	document.Find("a").Each(s.processElement)
}

func (s *scraper) addToVisited(url string) {
	s.visitedUrls.Store(url, struct{}{})
}

func (s *scraper) process() {
	for {
		select {
		case url := <-s.urlChan:
			log.Println("will process", url)
			go s.scrape(url)
		}
	}
}
