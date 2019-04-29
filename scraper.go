package scraper

import (
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"sync"

	"io"

	"github.com/PuerkitoBio/goquery"
)

type scraper struct {
	hostName    string
	rootUrl     string
	urls        sync.Map
	urlChan     chan string
	countingSem chan struct{}
	wg          *sync.WaitGroup
	done        chan struct{}
	client      http.Client
}

func New(rootUrl string) (*scraper, error) {
	return NewWithConfig(rootUrl, NewConfig())
}

func NewWithConfig(rootUrl string, config *config) (*scraper, error) {
	parsedUrl, err := neturl.Parse(rootUrl)
	if err != nil {
		return nil, err
	}
	if parsedUrl.Hostname() == "" {
		return nil, errors.New("the given URL should have a domain part")
	}
	client := http.Client{
		Timeout: config.Timeout(),
	}

	s := &scraper{
		hostName:    parsedUrl.Hostname(),
		rootUrl:     rootUrl,
		urlChan:     make(chan string, 1000),
		countingSem: make(chan struct{}, config.MaxWorkerAmount()),
		wg:          new(sync.WaitGroup),
		client:      client,
		done:        make(chan struct{}),
	}
	return s, nil
}

func (s *scraper) Scrape() {
	go s.process()
	s.init()
	s.wg.Wait()
	close(s.done)
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

func (s *scraper) process() {
	for {
		select {
		case url := <-s.urlChan:
			go s.scrape(url)
		case <-s.done:
			return
		}
	}
}

func (s *scraper) init() {
	s.processUrl(s.rootUrl)
}

func (s *scraper) processUrl(url string) {
	if s.isScraped(url) {
		return
	}
	s.addToUrls(url)
	s.wg.Add(1)
	s.urlChan <- url
}

func (s *scraper) isScraped(url string) bool {
	_, found := s.urls.Load(url)
	return found
}

func (s *scraper) addToUrls(url string) {
	s.urls.Store(url, struct{}{})
}

func (s *scraper) scrape(url string) {
	s.preScrape()
	defer s.postScrape()
	response, err := s.client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	s.processUrlsInPage(response.Body, url)
}

func (s *scraper) preScrape() {
	s.countingSem <- struct{}{}
}

func (s *scraper) postScrape() {
	s.wg.Done()
	<-s.countingSem
}

func (s *scraper) processUrlsInPage(body io.ReadCloser, baseUrl string) {
	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Println(err)
		return
	}
	document.Find("a").Each(func(_ int, element *goquery.Selection) {
		if href, exists := element.Attr("href"); exists {
			absoluteUrl := absoluteUrl(href, baseUrl)
			if absoluteUrl != "" && hasSameDomain(absoluteUrl, s.hostName) {
				s.processUrl(absoluteUrl)
			}
		}
	})
}
