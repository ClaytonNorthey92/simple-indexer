package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
	striptags "github.com/grokify/html-strip-tags-go"
)

const depth = 3

const statusInProgress = "in progress"

const statusFailed = "failed"

const statusSucceeded = "succeeded"

const statusCancelled = "cancelled"

var scriptRe = regexp.MustCompile(`<script>.+<\/script>`)

// CrawlJob represents a single computation that is
// or was responsible for gathering data from the internet
type CrawlJob struct {
	ID                 uint   `json:"id"`
	Status             string `json:"status"`
	IndexedPageCount   uint   `json:"indexed-page-count"`
	NewWordsAddedCount uint   `json:"new-words-added-count"`
	Details            string `json:"details"`
	StartURL           string `json:"start-url"`
}

// SortableCrawlJobs represent and array of crawl jobs
// sortable by id
type SortableCrawlJobs []*CrawlJob

func (s SortableCrawlJobs) Len() int {
	return len(s)
}

func (s SortableCrawlJobs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableCrawlJobs) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// Crawler represents something that can "Crawl" a url
// and its links
type Crawler interface {
	Jobs() []*CrawlJob
	Crawl(url string, i Indexer, stop chan bool)
}

// SimpleCrawler will collect content from a page and crawl links with a max
// depth of 3
type SimpleCrawler struct {
	jobs SortableCrawlJobs
}

// NewSimpleCrawler creates a SimpleCrawler and returns a reference to it
func NewSimpleCrawler() *SimpleCrawler {
	return &SimpleCrawler{
		jobs: []*CrawlJob{},
	}
}

// Jobs retuns a sorted slice of all crawl jobs the SimpleCrawler is responsible for
func (s *SimpleCrawler) Jobs() SortableCrawlJobs {
	sort.Sort(sort.Reverse(s.jobs))
	return s.jobs
}

// Crawl gathers content from a web pages, then follows links no more than 3 links deep
func (s *SimpleCrawler) Crawl(url string, i Indexer, stop chan bool) error {
	job := &CrawlJob{
		Status:             statusInProgress,
		IndexedPageCount:   0,
		NewWordsAddedCount: 0,
		ID:                 uint(time.Now().Unix()),
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

	cancelled := false

	c.OnXML("//title or head", func(e *colly.XMLElement) {
		title := e.Text
		url := strings.Split(e.Request.URL.String(), "?")[0]
		if visited[url] == true {
			return
		}
		body := e.Response.Body
		text := scriptRe.ReplaceAllString(string(body), "")
		fmt.Printf("found title for %s\n", url)
		titles[url] = title
		if len(body) == 0 {
			return
		}
		text = striptags.StripTags(string(body))
		wordsIndexed := i.IndexTextForPage(text, url, title)
		job.NewWordsAddedCount += uint(wordsIndexed)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		select {
		case <-stop:
			fmt.Println("stopped link following")
			cancelled = true
			job.Details = "user cancelled job, there is no guarantee these terms are indexed"
			fmt.Println("attempting to send to stop channel")
			stop <- true
			fmt.Println("sent to stop channel")
			return
		default:
			// no op
		}

		link := e.Attr("href")
		link = strings.Split(link, "?")[0]
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

	fmt.Println("job finished without error")

	if cancelled {
		job.Status = statusCancelled
	} else {
		job.Status = statusSucceeded
	}

	return nil
}
