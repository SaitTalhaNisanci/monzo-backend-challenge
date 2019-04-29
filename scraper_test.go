package scraper

import (
	"runtime"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScraperMemoryLeakage(t *testing.T) {
	before := runtime.NumGoroutine()
	scraper, err := New("https://dolarrekorkirdimi.com/")
	require.NoError(t, err)
	scraper.Scrape()
	assertTrueEventually(t, func() bool {
		after := runtime.NumGoroutine()
		return after == before
	})
}

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
