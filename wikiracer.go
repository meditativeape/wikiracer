package main

import (
	"fmt"
	"github.com/meditativeape/wikiracer/impl"
	"time"
)

func main() {
	startTime := time.Now()

	startUrl := "https://en.wikipedia.org/wiki/Computer_programming"
	endUrl := "https://en.wikipedia.org/wiki/Blade_Runner"
	path, _ := impl.FindPath(startUrl, endUrl)

	elapsed := time.Since(startTime)

	for i := len(path) - 1; i >= 0; i-- {
		fmt.Println(path[i])
	}
	fmt.Printf("Wikiracer took %s\n", elapsed)
}
