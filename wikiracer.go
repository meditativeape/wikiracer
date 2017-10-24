package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meditativeape/wikiracer/impl"
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
	//fmt.Printf("Request received. Start URL: %s, End URL: %s\n", startUrl, endUrl)
	if len(startUrl) == 0 || len(endUrl) == 0 {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	path, _ := impl.FindPath(startUrl, endUrl)
	elapsed := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"path": path,
		"time": elapsed.String(),
	})
}
