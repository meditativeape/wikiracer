package impl

import (
	// "github.com/meditativeape/wikiracer/util"
	"golang.org/x/net/html"
	"net/http"
)

type ArticleWithParent struct {
	Article       string
	ParentArticle string
}

func crawl(articleName string, ch chan ArticleWithParent) {
	// util.Logger.Printf("Crawling article: %s\n", articleName)
	resp, err := http.Get(GetUrlStringFromArticleName(articleName))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	visited := make(map[string]bool)
	visited[articleName] = true
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
						nextArticleName, err := GetArticleNameFromUrlString(attribute.Val)
						if err == nil && !visited[nextArticleName] {
							ch <- ArticleWithParent{nextArticleName, articleName}
							visited[nextArticleName] = true
						}
					}
				}
			}
		}
	}

	return
}
