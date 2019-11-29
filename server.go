package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	cr := NewSimpleCrawler()
	si := NewSimpleIndexer()
	go func() {
		err := cr.Crawl("http://go-colly.org/docs/examples/rate_limit/", si)
		if err != nil {
			panic(err)
		}
	}()
	r := gin.Default()
	r.GET("/jobs", func(c *gin.Context) {
		c.JSON(200, cr.Jobs())
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
