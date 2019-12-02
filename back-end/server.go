package main

import (
	"github.com/gin-gonic/gin"
)

type indexPostBody struct {
	URL string `json:"url" binding:"required"`
}

type searchParams struct {
	Query string `form:"q" binding:"required"`
}

type userError struct {
	Error string `json:"error"`
}

func allowInsecure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "*")
}

func sendUserError(c *gin.Context, err error) {
	userError := userError{
		Error: err.Error(),
	}
	c.JSON(400, userError)
}

func main() {

	cr := NewSimpleCrawler()
	si := NewPartialWordIndexer()
	r := gin.Default()
	stopChannels := new([]chan bool)

	r.Use(func(c *gin.Context) {
		allowInsecure(c)
	})

	r.GET("/jobs", func(c *gin.Context) {
		jobs := cr.Jobs()
		if len(jobs) > 5 {
			jobs = jobs[:5]
		}
		c.JSON(200, jobs)
	})

	r.DELETE("/index", func(c *gin.Context) {
		for _, c := range *stopChannels {
			c <- true
		}

		stopChannels = new([]chan bool)
		si.Clear()
		cr = NewSimpleCrawler()

		c.Status(204)
	})

	r.GET("/search", func(c *gin.Context) {
		var searchParams searchParams
		if err := c.ShouldBindQuery(&searchParams); err != nil {
			sendUserError(c, err)
			return
		}
		pages := si.GetPagesForWord(searchParams.Query)
		c.JSON(200, pages)
	})

	r.OPTIONS("/index", func(c *gin.Context) {
		c.Status(200)
	})

	r.POST("/index", func(c *gin.Context) {
		var indexBody indexPostBody
		if err := c.ShouldBindJSON(&indexBody); err != nil {
			sendUserError(c, err)
			return
		}

		stop := make(chan bool, 1)
		*stopChannels = append(*stopChannels, stop)
		go cr.Crawl(indexBody.URL, si, stop)
		c.Status(201)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
