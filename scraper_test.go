package scraper

import (
	"runtime"
	"testing"

	"time"

	"io/ioutil"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScraper(t *testing.T) {
	testCases := []struct {
		rootUrl string
		err     bool
	}{
		{"com", true},
		{"/a", true},
		{"ww.bb.cc", true},
		{"https://monzo.com", false},
	}
	for _, testCase := range testCases {
		_, err := New(testCase.rootUrl)
		if testCase.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

	}
}

func TestScraperDoesNotReturnAnyExternalLink(t *testing.T) {
	root := "http://www.monzo.com"
	scraper, err := New(root)
	require.NoError(t, err)
	scraper.Scrape()
	urls := scraper.Urls()
	assert.NotZero(t, len(urls))
	for _, url := range urls {
		assert.True(t, hasSameDomain(url, scraper.hostName))
	}
}

func TestScraperProcessUrlsInPage(t *testing.T) {

	expectedUrls := []string{
		"https://monzo.com/",
		"https://monzo.com/about",
		"https://monzo.com/blog",
		"https://monzo.com/community",
		"https://monzo.com/help",
		"https://monzo.com/download",
		"https://monzo.com/business",
		"https://monzo.com/features/apple-pay",
		"https://monzo.com/features/google-pay",
		"https://monzo.com/features/travel",
		"https://monzo.com/features/switch",
		"https://monzo.com/features/overdrafts",
		"https://monzo.com/press",
		"https://monzo.com/careers",
		"https://monzo.com/tone-of-voice",
		"https://monzo.com/blog/how-money-works",
		"https://monzo.com/transparency",
		"https://monzo.com/community/making-monzo",
		"https://monzo.com/faq",
		"https://monzo.com/legal/terms-and-conditions",
		"https://monzo.com/legal/fscs-information",
		"https://monzo.com/legal/privacy-policy",
		"https://monzo.com/legal/cookie-policy",
		"https://monzo.com/cdn-cgi/l/email-protection#751d10190535181a1b0f1a5b161a18",
		"https://monzo.com/cdn-cgi/l/email-protection#ff979a938fbf9290918590d19c9092",
	}

	body, err := ioutil.ReadFile("./testData/page.html")
	rootUrl := "https://monzo.com"
	require.NoError(t, err)
	scraper, err := New(rootUrl)
	require.NoError(t, err)
	scraper.processUrlsInPage(ioutil.NopCloser(bytes.NewReader(body)), rootUrl)
	urls := scraper.Urls()
	assert.Equal(t, len(expectedUrls), len(urls))
	for _, url := range urls {
		assert.Contains(t, expectedUrls, url)
	}
	for _, url := range expectedUrls {
		assert.Contains(t, urls, url)
	}

}

func TestScraperRoutineLeakage(t *testing.T) {
	before := runtime.NumGoroutine()
	scraper, err := New("https://dolarrekorkirdimi.com/")
	require.NoError(t, err)
	scraper.Scrape()
	assertTrueEventually(t, func() bool {
		after := runtime.NumGoroutine()
		return after == before
	})
}

func assertTrueEventually(t *testing.T, assertions func() bool) {
	startTime := time.Now()
	for time.Since(startTime) < 3*time.Minute {
		if assertions() {
			return
		}
		time.Sleep(100 * time.Microsecond)
	}
	t.Fail()
}
