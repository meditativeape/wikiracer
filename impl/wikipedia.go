package impl

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var blacklist = map[string]bool{
	"/wiki/Special:Random": true,
}

func IsEnWikiArticle(inputUrl *url.URL) bool {
	if inputUrl != nil &&
		(inputUrl.Host == "" || inputUrl.Host == "en.wikipedia.org") &&
		strings.HasPrefix(inputUrl.Path, "/wiki/") &&
		!blacklist[inputUrl.Path] {
		return true
	}
	return false
}

func GetArticleNameFromUrlString(inputUrl string) (string, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	return GetArticleNameFromParsedUrl(parsedUrl)
}

func GetArticleNameFromParsedUrl(inputUrl *url.URL) (string, error) {
	if !IsEnWikiArticle(inputUrl) {
		return "", errors.New("input URL does not point to an English Wikipedia article")
	}

	return inputUrl.Path[6:], nil
}

func GetUrlStringFromArticleName(articleName string) string {
	return "https://en.wikipedia.org/wiki/" + articleName
}

func IsReachable(inputUrl string) bool {
	resp, err := http.Get(inputUrl)
	if err != nil {
		return false
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
