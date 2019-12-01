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

type UserError struct {
	Error string `json:"error"`
}

func allowInsecure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "*")
}

func sendUserError(c *gin.Context, err error) {
	userError := UserError{
		Error: err.Error(),
	}
	c.JSON(400, userError)
}

func main() {
	cr := NewSimpleCrawler()
	si := NewDistanceIndexer()
	r := gin.Default()

	r.GET("/jobs", func(c *gin.Context) {
		allowInsecure(c)
		jobs := cr.Jobs()
		if len(jobs) > 5 {
			jobs = jobs[:5]
		}
		c.JSON(200, jobs)
	})

	r.DELETE("/index", func(c *gin.Context) {
		allowInsecure(c)
		si = NewDistanceIndexer()
		cr = NewSimpleCrawler()
		c.Status(204)
	})

	r.GET("/search", func(c *gin.Context) {
		allowInsecure(c)
		var searchParams SearchParams
		if err := c.ShouldBindQuery(&searchParams); err != nil {
			sendUserError(c, err)
			return
		}
		pages := si.GetPagesForWord(searchParams.Query)
		c.JSON(200, pages)
	})

	r.OPTIONS("/index", func(c *gin.Context) {
		allowInsecure(c)
		c.Status(200)
	})

	r.POST("/index", func(c *gin.Context) {
		allowInsecure(c)

		var indexBody IndexPostBody
		if err := c.ShouldBindJSON(&indexBody); err != nil {
			sendUserError(c, err)
			return
		}
		go func() {
			cr.Crawl(indexBody.URL, si)
		}()
		c.Status(201)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
