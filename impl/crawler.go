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

func crawl(article string, ch chan ArticleWithParent) {
	// Cache hit
	if ArticleCache.Contains(article) {
		// util.Logger.Printf("Cache hit for article: %s\n", article)
		value, _ := ArticleCache.Get(article)
		nextArticles := value.(*[]string)
		for _, nextArticle := range *nextArticles {
			ch <- ArticleWithParent{nextArticle, article}
		}
		return
	}

	// Cache miss
	// util.Logger.Printf("Cache miss. Crawling article: %s\n", article)
	resp, err := http.Get(GetUrlStringFromArticleName(article))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	nextArticles := make([]string, 0)
	visited := make(map[string]bool)
	visited[article] = true
	z := html.NewTokenizer(resp.Body)
	done := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			done = true
		case html.StartTagToken:
			token := z.Token()
			tagType := token.Data
			if tagType == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						nextArticle, err := GetArticleNameFromUrlString(attribute.Val)
						if err == nil && !visited[nextArticle] {
							ch <- ArticleWithParent{nextArticle, article}
							visited[nextArticle] = true
							nextArticles = append(nextArticles, nextArticle)
						}
					}
				}
			}
		}

		if done {
			break
		}
	}

	ArticleCache.Add(article, &nextArticles)
	return
}
