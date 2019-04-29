package scraper

import "time"

const defaultMaxWorkerAmount = 1000
const defaultTimeout = 1 * time.Minute

type config struct {
	maxWorkerAmount int
	timeout         time.Duration
}

func NewConfig() *config {
	return &config{
		maxWorkerAmount: defaultMaxWorkerAmount,
		timeout:         defaultTimeout,
	}
}

func (c *config) SetMaxWorkerAmount(amount int) {
	c.maxWorkerAmount = amount
}

func (c *config) MaxWorkerAmount() int {
	return c.maxWorkerAmount
}

func (c *config) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *config) Timeout() time.Duration {
	return c.timeout
}