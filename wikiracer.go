package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meditativeape/wikiracer/impl"
	"github.com/meditativeape/wikiracer/util"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()

	router.POST("/race", race)

	router.Run()
}

func race(c *gin.Context) {
	startTime := time.Now()
	startUrl := c.PostForm("startUrl")
	endUrl := c.PostForm("endUrl")
	util.Logger.Printf("[Main] Request received. Start URL: %s, End URL: %s\n", startUrl, endUrl)

	if len(startUrl) == 0 || len(endUrl) == 0 {
		c.JSON(http.StatusBadRequest, nil)
	} else {
		path, _ := impl.FindPath(startUrl, endUrl)
		c.JSON(http.StatusOK, gin.H{
			"path": path,
		})
	}

	elapsed := time.Since(startTime)
	util.Logger.Printf("[Main] Request served after %s. Start URL: %s, End URL: %s\n", elapsed, startUrl, endUrl)
}
