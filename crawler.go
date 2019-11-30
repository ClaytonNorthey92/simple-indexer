package main

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/gocolly/colly"
	striptags "github.com/grokify/html-strip-tags-go"
)

var depth = 3

var statusInProgress = "in progress"
var statusFailed = "failed"
var statusSucceeded = "succeeded"

type CrawlJob struct {
	ID                 uint   `json:"id"`
	Status             string `json:"status"`
	IndexedPageCount   uint   `json:"indexed-page-count"`
	NewWordsAddedCount uint   `json:"new-words-added-count"`
	Details            string `json:"details"`
	StartURL           string `json:"start-url"`
}

type SortableCrawlJobs []*CrawlJob

func (s SortableCrawlJobs) Len() int {
	return len(s)
}

func (s SortableCrawlJobs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableCrawlJobs) Less(i, j int) bool {
	return s[i].ID > s[j].ID
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
	jobs SortableCrawlJobs
}

func (s *SimpleCrawler) Jobs() SortableCrawlJobs {
	sort.Sort(sort.Reverse(s.jobs))
	return s.jobs
}

func (s *SimpleCrawler) Crawl(url string, i Indexer) error {
	job := &CrawlJob{
		Status:             statusInProgress,
		IndexedPageCount:   0,
		NewWordsAddedCount: 0,
		ID:                 uint(time.Now().Nanosecond()),
		StartURL:           url,
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
		re := regexp.MustCompile(`<script>.+<\/script>`)
		text := re.ReplaceAllString(string(body), "")
		fmt.Printf("found title for %s\n", url)
		titles[url.String()] = title
		text = striptags.StripTags(string(body))
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
		job.Details = err.Error()
		return err
	}

	job.Status = statusSucceeded

	return nil
}
