package main

import (
	"github.com/gin-gonic/gin"
)

type IndexPostBody struct {
	URL string `json:"url" binding:"required"`
}

type SearchParams struct {
	Query string `form:"q" binding:"required"`
}

func main() {
	cr := NewSimpleCrawler()
	si := NewSimpleIndexer()
	r := gin.Default()

	r.GET("/jobs", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, cr.Jobs())
	})

	r.GET("/search", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		var searchParams SearchParams
		if err := c.ShouldBindQuery(&searchParams); err != nil {
			panic(err)
		}
		pages := si.GetPagesForWord(searchParams.Query)
		c.JSON(200, pages)
	})

	r.POST("/index", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var indexBody IndexPostBody
		if err := c.ShouldBindJSON(&indexBody); err != nil {
			panic(err)
		}
		go func() {
			cr.Crawl(indexBody.URL, si)
		}()
		c.Status(201)
	})

	r.OPTIONS("/index", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Accept", "*")
		c.Header("Allow", "POST")
		c.Header("Access-Control-Request-Method", "POST")
		c.Status(200)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
