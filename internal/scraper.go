package internal

import (
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"sync"

	"io"

	"time"

	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
)

const defaultWorkerAmount = 1000
const defaultConcurrentScrapers = 10000

type scraper struct {
	hostName     string
	rootUrl      string
	urls         sync.Map
	workerAmount int
	urlChan      chan string
	wg           *sync.WaitGroup
	countingSem  chan struct{}
	done         chan struct{}
	client       http.Client
}

func NewScraper(rootUrl string) (*scraper, error) {
	parsedUrl, err := neturl.Parse(rootUrl)
	if err != nil {
		return nil, err
	}
	if parsedUrl.Hostname() == "" {
		return nil, errors.New("the given URL should have a domain part")
	}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	s := &scraper{
		hostName:     parsedUrl.Hostname(),
		rootUrl:      rootUrl,
		urlChan:      make(chan string, 1000),
		countingSem:  make(chan struct{}, defaultConcurrentScrapers),
		workerAmount: defaultWorkerAmount,
		wg:           new(sync.WaitGroup),
		client:       client,
		done:         make(chan struct{}),
	}

	return s, nil
}

func (s *scraper) Scrape() {
	s.startWorkers()
	s.init()
	s.wg.Wait()
	close(s.done)
	log.Println("gg")
}

func (s *scraper) Urls() []string {
	urls := make([]string, 0)
	rangeF := func(key, _ interface{}) bool {
		urls = append(urls, key.(string))
		return true
	}
	s.urls.Range(rangeF)
	return urls
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

var closedCount int32 = 0

func (s *scraper) worker(id int) {
	for {
		select {
		case url := <-s.urlChan:
			s.countingSem <- struct{}{}
			s.scrape(url)
			s.wg.Done()
			<-s.countingSem
		case <-s.done:
			atomic.AddInt32(&closedCount, 1)
			log.Println(closedCount)
			return
		}
	}
}

func (s *scraper) processElement(_ int, element *goquery.Selection) string {
	href, _ := element.Attr("href")
	return href
}

func (s *scraper) scrape(url string) {
	response, err := s.client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	absoluteUrls := s.findUrlsInPage(response.Body, url)
	s.processUrls(absoluteUrls)
}

func (s *scraper) findUrlsInPage(body io.ReadCloser, url string) []string {
	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Println(err)
		return nil
	}
	urls := document.Find("a").Map(s.processElement)
	return s.convertToAbsolute(urls, url)
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
	_, found := s.urls.Load(url)
	return found
}

func (s *scraper) addToVisited(url string) {
	s.urls.Store(url, struct{}{})
}
