package internal

import (
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const defaultWorkerAmount = 1000

type scraper struct {
	hostName     string
	rootUrl      string
	visitedUrls  sync.Map
	workerAmount int
	urlChan      chan string
	wg           *sync.WaitGroup
	done         chan struct{}
}

func NewScraper(rootUrl string) (*scraper, error) {
	parsedUrl, err := neturl.Parse(rootUrl)
	if err != nil {
		return nil, err
	}
	if parsedUrl.Hostname() == "" {
		return nil, errors.New("the given URL should have a domain part")
	}
	s := &scraper{
		hostName:     parsedUrl.Hostname(),
		rootUrl:      rootUrl,
		urlChan:      make(chan string, 1000),
		workerAmount: defaultWorkerAmount,
		wg:           new(sync.WaitGroup),
		done:         make(chan struct{}),
	}
	return s, nil
}

func (s *scraper) Scrape() {
	s.startWorkers()
	s.init()
	s.wg.Wait()
	close(s.done)
}

func (s *scraper) init() {
	s.processUrl(s.rootUrl)
}

func (s *scraper) startWorkers() {
	for i := 0; i < s.workerAmount; i++ {
		go s.worker(i)
	}
}

func (s *scraper) processUrl(url string) {
	if s.isScraped(url) {
		return
	}
	s.addToVisited(url)
	s.wg.Add(1)
	s.urlChan <- url
}

func (s *scraper) worker(id int) {
	for {
		select {
		case url := <-s.urlChan:
			log.Println("will process", id, url)
			s.scrape(url)
			s.wg.Done()
		case <-s.done:
			return
		}
	}
}

func (s *scraper) processElement(_ int, element *goquery.Selection) string {
	href, _ := element.Attr("href")
	return href
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
		log.Println("error", err)
		return
	}
	urls := document.Find("a").Map(s.processElement)
	absoluteUrls := s.convertToAbsolute(urls, url)
	s.processUrls(absoluteUrls)
}

func (s *scraper) convertToAbsolute(urls []string, baseUrl string) []string {
	absoluteUrls := make([]string, 0)
	for _, url := range urls {
		absoluteUrls = append(absoluteUrls, absoluteUrl(url, baseUrl))
	}
	return absoluteUrls
}

func (s *scraper) processUrls(absoluteUrls []string) {
	for _, url := range absoluteUrls {
		if !s.hasSameDomain(url) {
			continue
		}
		s.processUrl(url)
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

func (s *scraper) isScraped(url string) bool {
	_, found := s.visitedUrls.Load(url)
	return found
}

func (s *scraper) addToVisited(url string) {
	s.visitedUrls.Store(url, struct{}{})
}
