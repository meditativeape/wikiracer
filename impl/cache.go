package impl

import (
	"github.com/hashicorp/golang-lru"
	"github.com/meditativeape/wikiracer/util"
)

const CacheSize int = 200000

var ArticleCache *lru.Cache = initCache()

func initCache() *lru.Cache {
	cache, err := lru.New(CacheSize)
	util.PanicIfError(err)
	return cache
}
