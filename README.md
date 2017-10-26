# wikiracer

Finds a path between two Wikipedia articles, using only Wikipedia links.

## Approach

Wikiracer runs a one-way parallel BFS (Breadth First Search) from the given start URL to crawl the graph of Wikipedia articles until it reaches the target URL.

At each level of BFS, the work is shared across a number of Goroutines. These Goroutines fetch work from a common input channel, which streams links found by Goroutines for the previous level, crawl the articles, and send links found in these articles to another common output channel. The main function collects the output into an array, removes duplicates and links that have already been crawled, and starts the next batch of Goroutines to crawl the new links.

For simplicity, Wikiracer only uses English articles (URL prefix: `en.wikipedia.org/wiki/`).

## Installation

```
$ go get github.com/meditativeape/wikiracer
$ cd $GOPATH/src/github.com/meditativeape/wikiracer
$ make install
```

## Usage

Start wikiracer by running `wikiracer`. It spins up an HTTP server that listens on port `8080`.

Wikiracer offers one REST endpoint, `POST /race`, that expects two keys in the POST form: `startUrl` and `endUrl`. It returns the path found in JSON format. You could use your favorite client, such as cURL or Postman, to query against this endpoint.

Example request as a cURL command:
```
curl localhost:8080/race -F startUrl=https://en.wikipedia.org/wiki/Computer_programming -F endUrl=https://en.wikipedia.org/wiki/Blade_Runner 
```

## Logging

Wikiracer keeps a lightweighted log under `/tmp/wikiracer/service.log`.

## License

The MIT License
