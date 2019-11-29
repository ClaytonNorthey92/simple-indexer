package main

import (
	"sort"
	"strings"
)

type SimpleIndexer struct {
	index []WebPage
}

func NewSimpleIndexer() *SimpleIndexer {
	return &SimpleIndexer{
		index: []WebPage{},
	}
}

func (s *SimpleIndexer) IndexTextForPage(pageContent string, url string, title string) {
	s.index = append(s.index, WebPage{
		URL:     url,
		Title:   title,
		Content: pageContent,
	})
}

func (s *SimpleIndexer) GetPagesForWord(word string) SortableWebPagesWithCount {
	pages := make(SortableWebPagesWithCount, len(s.index))

	for i, v := range s.index {
		count := strings.Count(v.Content, word)
		pages[i] = WebPageWithCount{
			Count: uint(count),
			WebPage: WebPage{
				URL:   v.URL,
				Title: v.Title,
			},
		}
	}

	sort.Sort(pages)

	return pages
}
