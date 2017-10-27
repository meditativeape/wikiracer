package impl

import (
	"github.com/meditativeape/wikiracer/util"
	"runtime"
	"sync"
)

var numCrawlersPerLevel int = runtime.NumCPU() * 10

func FindPath(startUrl string, endUrl string) (*[]string, *map[string]string) {
	articleToParent := make(map[string]string)
	emptyPath := make([]string, 0)

	// Make sure both startUrl and endUrl point to valid Wikipedia articles
	startArticle, err := GetArticleNameFromUrlString(startUrl)
	if err != nil || !IsReachable(startUrl) {
		util.Logger.Printf("Invalid StartURL: %s\n", startUrl)
		return &emptyPath, &articleToParent
	}
	endArticle, err := GetArticleNameFromUrlString(endUrl)
	if err != nil || !IsReachable(endUrl) {
		util.Logger.Printf("Invalid EndURL: %s\n", endUrl)
		return &emptyPath, &articleToParent
	}

	articleToParent[startArticle] = "root"
	if startUrl == endUrl {
		emptyPath = append(emptyPath, startUrl)
		return &emptyPath, &articleToParent
	}

	outputCh := make(chan ArticleWithParent)
	var wg sync.WaitGroup
	wg.Add(1)
	go crawlSingleArticleAndSync(startArticle, outputCh, &wg)
	go closeChannelOnWg(outputCh, &wg)

	level := 0
	found := false

	for {
		// collect outputs from crawlers at the current level, and skip articles that are already visited
		toBeCrawled := make([]string, 0)
		for articleWithParent := range outputCh {
			nextArticle := articleWithParent.Article
			if articleToParent[nextArticle] == "" {
				articleToParent[nextArticle] = articleWithParent.ParentArticle
				if nextArticle == endArticle {
					found = true
					break
				} else {
					toBeCrawled = append(toBeCrawled, nextArticle)
				}
			}
		}
		if found {
			util.Logger.Printf(
				"Found a path while crawling level %d! Start: %s End: %s\n",
				level, startArticle, endArticle)
			break
		} else {
			util.Logger.Printf(
				"Collected %d outputs from crawlers for level %d. Start: %s End: %s\n",
				len(toBeCrawled), level, startArticle, endArticle)
			level++
		}

		// start a fixed number of goroutines to crawl all articles at the next level
		inputCh := make(chan string)
		nextOutputCh := make(chan ArticleWithParent, 1000)
		var nextWg sync.WaitGroup
		nextWg.Add(numCrawlersPerLevel)
		for i := 0; i < numCrawlersPerLevel; i++ {
			go crawlMultipleArticlesAndSync(inputCh, nextOutputCh, &nextWg)
		}
		go closeChannelOnWg(nextOutputCh, &nextWg)
		util.Logger.Printf(
			"Started %d crawlers for level %d. Start: %s End: %s\n",
			numCrawlersPerLevel, level, startArticle, endArticle)

		// start a goroutine to feed URLs into input channel
		go feedArticlesIntoChannel(toBeCrawled, inputCh)
		outputCh = nextOutputCh
	}

	return getPath(&articleToParent, endArticle), &articleToParent
}

func feedArticlesIntoChannel(articles []string, ch chan string) {
	for _, article := range articles {
		ch <- article
	}
	close(ch)
}

func crawlSingleArticleAndSync(article string, outputCh chan ArticleWithParent, wg *sync.WaitGroup) {
	crawl(article, outputCh)
	(*wg).Done()
}

func crawlMultipleArticlesAndSync(inputCh chan string, outputCh chan ArticleWithParent, wg *sync.WaitGroup) {
	for article := range inputCh {
		crawl(article, outputCh)
	}
	(*wg).Done()
}

func closeChannelOnWg(ch chan ArticleWithParent, wg *sync.WaitGroup) {
	(*wg).Wait()
	close(ch)
}

func getPath(articleToParent *map[string]string, endArticle string) *[]string {
	reversedPath := make([]string, 0)
	currentArticle := endArticle
	for currentArticle != "root" {
		reversedPath = append(reversedPath, currentArticle)
		currentArticle = (*articleToParent)[currentArticle]
	}

	pathLen := len(reversedPath)
	path := make([]string, pathLen)
	for i := 0; i < pathLen; i++ {
		path[i] = GetUrlStringFromArticleName(reversedPath[pathLen-i-1])
	}
	return &path
}
