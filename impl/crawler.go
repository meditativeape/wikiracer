package impl

import (
	// "github.com/meditativeape/wikiracer/util"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

var blacklist = map[string]bool{
	"/wiki/Special:Random": true,
}

type UrlWithParent struct {
	Url       *url.URL
	ParentUrl *url.URL
}

func crawl(urlToCrawl *url.URL, ch chan UrlWithParent) {
	// util.Logger.Printf("Crawling article: %s\n", urlToCrawl.Path)
	resp, err := http.Get(urlToCrawl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	visited := make(map[url.URL]bool)
	visited[*urlToCrawl] = true
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			token := z.Token()
			tagType := token.Data
			if tagType == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						absoluteUrl := getAbsoluteUrl(urlToCrawl, attribute.Val)
						if isEnWikiArticle(absoluteUrl) && !visited[*absoluteUrl] {
							ch <- UrlWithParent{absoluteUrl, urlToCrawl}
							visited[*absoluteUrl] = true
						}
					}
				}
			}
		}
	}

	return
}

func getAbsoluteUrl(currentUrl *url.URL, link string) *url.URL {
	newUrl, err := url.Parse(link)
	if err != nil || len(link) == 0 || link[0] == '#' { // ignore links to fragments on current page, e.g. "#cite_ref-7"
		return nil
	}
	if !newUrl.IsAbs() { // e.g. "//shop.wikimedia.org"
		newUrl.Scheme = currentUrl.Scheme
	}
	if newUrl.Host == "" { // e.g. "/wiki/System_programming"
		newUrl.Host = currentUrl.Host
	}
	return newUrl
}

func isEnWikiArticle(urlToCheck *url.URL) bool {
	if urlToCheck != nil && urlToCheck.Host == "en.wikipedia.org" &&
		strings.HasPrefix(urlToCheck.Path, "/wiki/") && !blacklist[urlToCheck.Path] {
		return true
	}
	return false
}
