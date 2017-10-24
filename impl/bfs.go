package impl

import (
	"net/url"
)

func FindPath(startUrl string, endUrl string) ([]string, map[string]string) {
	urlToParent := make(map[string]string)
	parsedStartUrl, err := url.Parse(startUrl)
	if err != nil {
		panic(err.Error())
	}

	urlToParent[startUrl] = "root"
	ch := make(chan UrlWithParent, 1000)
	go crawl(parsedStartUrl, ch)

	for nextUrl := range ch {
		nextUrlString := nextUrl.Url.String()
		if urlToParent[nextUrlString] == "" {
			urlToParent[nextUrlString] = nextUrl.ParentUrl.String()
			if nextUrlString == endUrl {
				break
			} else {
				go crawl(nextUrl.Url, ch)
			}
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
