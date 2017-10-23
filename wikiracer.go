package main

import (
	"github.com/meditativeape/wikiracer/crawler"
)

func main() {
	crawler.GetLinksInPage("https://en.wikipedia.org/wiki/Computer_programming")
}
