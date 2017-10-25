package impl

import (
	"net/url"
)

func FindPath(startUrl string, endUrl string) (*[]string, *map[string]string) {
	urlToParent := make(map[string]string)
	parsedStartUrl, err := url.Parse(startUrl)
	if err != nil {
		panic(err.Error())
	}

	urlToParent[startUrl] = "root"
	ch := make(chan UrlWithParent, 1000)
	initPath := make([]string, 0)
	go crawl(parsedStartUrl, ch, &initPath)

	for nextUrl := range ch {
		nextUrlString := nextUrl.Url.String()
		currentUrlString := nextUrl.ParentUrl.String()
		if urlToParent[nextUrlString] == "" {
			urlToParent[nextUrlString] = currentUrlString
			if nextUrlString == endUrl {
				break
			} else {
				go crawl(nextUrl.Url, ch, getPath(&urlToParent, currentUrlString))
			}
		}
	}

	return getPath(&urlToParent, endUrl), &urlToParent
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
