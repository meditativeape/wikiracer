package impl

import (
	"fmt"
	"net/url"
)

func FindPath(startUrl string, endUrl string) ([]string, map[string]string) {
	queue := make([]*url.URL, 0)
	urlToParent := make(map[string]string)
	parsedStartUrl, err := url.Parse(startUrl)
	if err != nil {
		panic(err.Error())
	}

	queue = append(queue, parsedStartUrl)
	urlToParent[startUrl] = "root"
	for {
		urlToVisit := queue[0]
		queue = queue[1:]

		ch := make(chan *url.URL)
		go crawl(urlToVisit, ch)

		for link := range ch {
			if urlToParent[link.String()] == "" {
				urlToParent[link.String()] = urlToVisit.String()
				queue = append(queue, link)
			}
		}
		if urlToParent[endUrl] != "" {
			fmt.Println("Found a path!")
			break
		}
	}

	path := make([]string, 0)
	currentUrl := endUrl
	for currentUrl != "root" {
		path = append(path, currentUrl)
		currentUrl = urlToParent[currentUrl]
	}

	return path, urlToParent
}
