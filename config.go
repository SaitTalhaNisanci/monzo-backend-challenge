package scraper

import "time"

const defaultMaxWorkerAmount = 1000
const defaultConTimeout = 1 * time.Minute

type config struct {
	maxWorkerAmount int
	conTimeout      time.Duration
}

// NewConfig returns a config for scraper.
func NewConfig() *config {
	return &config{
		maxWorkerAmount: defaultMaxWorkerAmount,
		conTimeout:      defaultConTimeout,
	}
}

// SetMaxWorkerAmount sets the maximum workers that can run concurrently
// while scraping. The value of 1 equals to scraping sequentially.
func (c *config) SetMaxWorkerAmount(amount int) {
	c.maxWorkerAmount = amount
}

// MaxWorkerAmount returns the maximum worker amounts.
func (c *config) MaxWorkerAmount() int {
	return c.maxWorkerAmount
}

// SetConTimeout sets the connection timeout for each request in scraper.
// The value of 0 means no timeout, which is not recommended as the scraper
// will be stuck when server is stuck.
func (c *config) SetConTimeout(timeout time.Duration) {
	c.conTimeout = timeout
}

// ConTimeout returns the connection timeout for scraper.
func (c *config) ConTimeout() time.Duration {
	return c.conTimeout
}
