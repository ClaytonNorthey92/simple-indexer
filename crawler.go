package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	striptags "github.com/grokify/html-strip-tags-go"
)

var depth = 3

var statusInProgress = "in progress"
var statusFailed = "failed"
var statusSucceeded = "succeeded"

type CrawlJob struct {
	ID                 uint
	Status             string
	IndexedPageCount   uint
	NewWordsAddedCount uint
}

type Crawler interface {
	Jobs() []*CrawlJob
	Crawl(url string, i *Indexer)
}

func NewSimpleCrawler() *SimpleCrawler {
	return &SimpleCrawler{
		jobs: []*CrawlJob{},
	}
}

type SimpleCrawler struct {
	jobs []*CrawlJob
}

func (s *SimpleCrawler) Jobs() []*CrawlJob {
	return s.jobs
}

func (s *SimpleCrawler) Crawl(url string, i Indexer) error {
	job := &CrawlJob{
		Status:             statusInProgress,
		IndexedPageCount:   0,
		NewWordsAddedCount: 0,
		ID:                 uint(time.Now().Nanosecond()),
	}
	if s.jobs == nil {
		s.jobs = []*CrawlJob{
			job,
		}
	} else {
		s.jobs = append(s.jobs, job)
	}

	c := colly.NewCollector()
	c.MaxDepth = depth
	c.Limit(&colly.LimitRule{
		Delay: 5 * time.Second,
	})
	titles := map[string]string{}
	visited := map[string]bool{}

	c.OnXML("//title or head", func(e *colly.XMLElement) {
		title := e.Text
		url := e.Request.URL
		body := e.Response.Body
		fmt.Printf("found title for %s\n", url)
		titles[url.String()] = title
		text := striptags.StripTags(string(body))
		i.IndexTextForPage(text, url.String(), title)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if visited[link] == true {
			return
		}
		visited[link] = true
		fmt.Printf("navigating to %s\n", link)
		e.Request.Visit(link)
	})

	c.OnScraped(func(r *colly.Response) {
		job.IndexedPageCount++
	})

	err := c.Visit(url)
	if err != nil {
		job.Status = statusFailed
		return err
	}

	job.Status = statusSucceeded

	return nil
}
