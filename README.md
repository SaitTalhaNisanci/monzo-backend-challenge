# Monzo Backend Challenge
[![GoDoc](https://godoc.org/github.com/SaitTalhaNisanci/monzo-backend-challenge?status.svg)](https://godoc.org/github.com/SaitTalhaNisanci/monzo-backend-challenge)
[![Go Report Card](https://goreportcard.com/badge/github.com/SaitTalhaNisanci/monzo-backend-challenge)](https://goreportcard.com/report/github.com/SaitTalhaNisanci/monzo-backend-challenge)

## Installation 

Make sure you have Go 1.9+ installed, [install it here](https://golang.org/doc/install).

Make sure your **GOROOT** and **GOPATH** are set correctly.

Run the following to get the code:

```
go get -u github.com/SaitTalhaNisanci/monzo-backend-challenge
```

## Technology

- Go
- [testify](https://github.com/stretchr/testify)
- [goquery](https://github.com/PuerkitoBio/goquery) 

[CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments) is followed 
during the development.


## How to use

```go
root := "http://www.monzo.com"
s, err := scraper.New(root)
if err != nil {
    log.Fatal(err)
}
s.Scrape()
urls := s.Urls()
for _, url := range urls {
    fmt.Println(url)
}

```

If you want to change the default config:

```go
cfg := scraper.NewConfig()
cfg.SetMaxWorkerAmount(10)
cfg.SetConTimeout(10 * time.Second)
scraper.NewWithConfig("https://www.monzo.com", cfg)
```

## How it works

When a url is scraped, all the links within the url are sent to
a url channel. For each new url(not visited), a new go routine is
created for scraping. Waitgroup is used to wait until all scraping 
is done.

A buffered channel(as a counting semaphore) is used to limit the maximum
number of go routines for the library. Otherwise, with a high load
the program could crash.

## Test

To run tests:
```
go test ./...
```

To run tests with race:

```
go test -race ./...
``` 
 