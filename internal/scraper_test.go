package internal

import (
	"runtime"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScraperMemoryLeakage(t *testing.T) {
	before := runtime.NumGoroutine()
	scraper, err := NewScraper("https://www.monzo.com/")
	require.NoError(t, err)
	scraper.Scrape()
	assertTrueEventually(t, func() bool {
		after := runtime.NumGoroutine()
		return after == before
	})
}

func TestNewScraperWithoutDomain(t *testing.T) {
	_, err := NewScraper("com")
	assert.Error(t, err)
}

func assertTrueEventually(t *testing.T, assertions func() bool) {
	startTime := time.Now()
	for time.Since(startTime) < 3*time.Minute {
		if assertions() {
			return
		}
		time.Sleep(10 * time.Microsecond)
	}
	t.Fail()
}

func TestScraperDoesNotReturnAnyExternalLink(t *testing.T) {
	domain := "https://www.monzo.com/"
	scraper, err := NewScraper(domain)
	require.NoError(t, err)
	scraper.Scrape()
	urls := scraper.Urls()
	assert.NotZero(t, len(urls))
	for _, url := range urls {
		assert.Contains(t, url, domain)
	}
}
