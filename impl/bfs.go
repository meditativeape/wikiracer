package impl

import (
	// "fmt"
	"github.com/meditativeape/wikiracer/util"
	"net/url"
	"runtime"
	"sync"
)

var numCrawlersPerLevel int = runtime.NumCPU() * 10

func FindPath(startUrl string, endUrl string) (*[]string, *map[string]string) {
	urlToParent := make(map[string]string)
	parsedStartUrl, err := url.Parse(startUrl)
	util.PanicIfError(err)

	urlToParent[startUrl] = "root"
	if startUrl == endUrl {
		path := make([]string, 1)
		path[0] = startUrl
		return &path, &urlToParent
	}

	outputCh := make(chan UrlWithParent)
	var wg sync.WaitGroup
	wg.Add(1)
	go crawlSingleUrlAndSync(parsedStartUrl, outputCh, &wg)
	go closeChannelOnWg(outputCh, &wg)

	level := 0
	found := false

	for {
		// collect outputs from crawlers at the current level, and skip URLs that are already visited once
		toBeCrawled := make([]UrlWithParent, 0)
		for urlWithParent := range outputCh {
			nextUrlString := urlWithParent.Url.String()
			if urlToParent[nextUrlString] == "" {
				urlToParent[nextUrlString] = urlWithParent.ParentUrl.String()
				if nextUrlString == endUrl {
					found = true
					break
				} else {
					toBeCrawled = append(toBeCrawled, urlWithParent)
				}
			}
		}
		if found {
			util.Logger.Printf(
				"Found path while crawling level %d! StartURL: %s, EndURL: %s\n",
				level, startUrl, endUrl)
			break
		} else {
			util.Logger.Printf(
				"Collected %d outputs from crawlers for level %d. StartURL: %s, EndURL: %s\n",
				len(toBeCrawled), level, startUrl, endUrl)
			level++
		}

		// start a fixed number of goroutines to crawl all URLs at the next level
		inputCh := make(chan UrlWithParent)
		nextOutputCh := make(chan UrlWithParent, 1000)
		var nextWg sync.WaitGroup
		nextWg.Add(numCrawlersPerLevel)
		for i := 0; i < numCrawlersPerLevel; i++ {
			go crawlMultipleUrlsAndSync(inputCh, nextOutputCh, &nextWg)
		}
		go closeChannelOnWg(nextOutputCh, &nextWg)
		util.Logger.Printf(
			"Started %d crawlers for level %d. StartURL: %s, EndURL: %s\n",
			numCrawlersPerLevel, level, startUrl, endUrl)

		// start a goroutine to feed URLs into input channel
		go feedUrlsIntoChannel(toBeCrawled, inputCh)
		outputCh = nextOutputCh
	}

	return getPath(&urlToParent, endUrl), &urlToParent
}

func feedUrlsIntoChannel(urls []UrlWithParent, ch chan UrlWithParent) {
	for _, url := range urls {
		ch <- url
	}
	close(ch)
}

func crawlSingleUrlAndSync(urlToCrawl *url.URL, outputCh chan UrlWithParent, wg *sync.WaitGroup) {
	crawl(urlToCrawl, outputCh)
	(*wg).Done()
}

func crawlMultipleUrlsAndSync(inputCh chan UrlWithParent, outputCh chan UrlWithParent, wg *sync.WaitGroup) {
	for urlToCrawl := range inputCh {
		crawl(urlToCrawl.Url, outputCh)
	}
	(*wg).Done()
}

func closeChannelOnWg(ch chan UrlWithParent, wg *sync.WaitGroup) {
	(*wg).Wait()
	// util.Logger.Printf("All crawlers have finished. Closing channel.\n")
	close(ch)
}

func getPath(urlToParent *map[string]string, endUrl string) *[]string {
	path := make([]string, 0)
	currentUrl := endUrl
	for currentUrl != "root" {
		path = append(path, currentUrl)
		currentUrl = (*urlToParent)[currentUrl]
	}
	return &path
}
