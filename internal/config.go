package internal

import "time"

const defaultWorkerAmount = 1000
const defaultTimeout = 10 * time.Second

type config struct {
	workerAmount int
	timeout      time.Duration
}

func NewConfig() *config {
	return &config{
		workerAmount: defaultWorkerAmount,
		timeout:      defaultTimeout,
	}
}

func (c *config) SetWorkerAmount(amount int) {
	c.workerAmount = amount
}

func (c *config) WorkerAmount() int {
	return c.workerAmount
}

func (c *config) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *config) Timeout() time.Duration {
	return c.timeout
}
