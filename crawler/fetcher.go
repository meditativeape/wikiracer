package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func GetLinksInPage(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := make(map[string]string)
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			for k, v := range result {
				fmt.Printf("key[%s] value[%s]\n", k, v)
			}
			return result, nil
		case html.StartTagToken:
			token := z.Token()
			tagType := token.Data
			if tagType == "a" {
				var link, title string
				for _, attribute := range token.Attr {
					switch attribute.Key {
					case "href":
						link = attribute.Val
					case "title":
						title = attribute.Val
					}
				}
				result[link] = title
			}
		}
	}

	for k, v := range result {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
	return result, nil
}
