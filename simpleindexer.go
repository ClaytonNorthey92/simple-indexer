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
	pageContent = strings.ToLower(pageContent)
	for i, v := range s.index {
		if v.URL == url {
			s.index[i] = WebPage{
				URL:     url,
				Title:   title,
				Content: pageContent,
			}
			return
		}
	}

	s.index = append(s.index, WebPage{
		URL:     url,
		Title:   title,
		Content: pageContent,
	})
}
func (s *SimpleIndexer) GetPagesForWord(word string) SortableWebPagesWithCount {
	word = strings.ToLower(word)

	pages := SortableWebPagesWithCount{}

	for _, v := range s.index {
		count := strings.Count(v.Content, word)
		if count <= 0 {
			continue
		}
		pages = append(pages, WebPageWithCount{
			Count: uint(count),
			WebPage: WebPage{
				URL:   v.URL,
				Title: v.Title,
			},
		})
	}

	sort.Sort(pages)

	return pages
}
