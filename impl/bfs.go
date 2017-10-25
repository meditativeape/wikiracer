package impl

import (
	"github.com/meditativeape/wikiracer/util"
	"net/url"
	"sync"
)

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

	ch := make(chan UrlWithParent, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go crawl(parsedStartUrl, ch, &wg)
	go closeChannelOnWg(ch, &wg)
	level := 1
	found := false

	for {
		var nextWg sync.WaitGroup
		nextCh := make(chan UrlWithParent, 100000)
		for urlWithParent := range ch {
			nextUrlString := urlWithParent.Url.String()
			if urlToParent[nextUrlString] == "" {
				urlToParent[nextUrlString] = urlWithParent.ParentUrl.String()
				if nextUrlString == endUrl {
					found = true
					break
				} else {
					nextWg.Add(1)
					go crawl(urlWithParent.Url, nextCh, &nextWg)
				}
			}
		}
		go closeChannelOnWg(nextCh, &nextWg)

		if found {
			util.Logger.Printf(
				"Found path while crawling level %d! StartURL: %s, EndURL: %s\n",
				level, startUrl, endUrl)
			break
		} else {
			util.Logger.Printf(
				"Finished crawling level %d. Onto the next level... StartURL: %s, EndURL: %s\n",
				level, startUrl, endUrl)
			level++
			ch = nextCh
		}
	}

	return getPath(&urlToParent, endUrl), &urlToParent
}

func closeChannelOnWg(ch chan UrlWithParent, wg *sync.WaitGroup) {
	(*wg).Wait()
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
