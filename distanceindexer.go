package main

import (
	"regexp"
	"sort"
	"strings"
)

type DistanceIndexer struct {
	index map[string][]WebPageWithCount
}

func NewDistanceIndexer() *DistanceIndexer {
	return &DistanceIndexer{
		index: map[string][]WebPageWithCount{},
	}
}

func (d *DistanceIndexer) GetPagesForWord(word string) SortableWebPagesWithCount {
	word = strings.ToLower(word)
	pages := map[string]*WebPageWithCount{}
	for k, v := range d.index {
		k = strings.ToLower(k)
		if k == word || strings.Contains(k, word) {
			for _, page := range v {
				if pages[page.URL] == nil {
					pages[page.URL] = &WebPageWithCount{
						WebPage: WebPage{
							URL:   page.URL,
							Title: page.Title,
						},
						Count: 1,
					}
				} else {
					pages[page.URL].Count++
				}
			}
		}
	}

	sortablePages := make(SortableWebPagesWithCount, len(pages))
	idx := 0
	for _, v := range pages {
		sortablePages[idx] = *v
		idx++
	}

	sort.Sort(sortablePages)
	return sortablePages
}

func (d *DistanceIndexer) IndexTextForPage(pageContent string, url string, title string) int {
	re := regexp.MustCompile(` |\n|\t+`)
	text := re.Split(pageContent, -1)
	originalIndexSize := len(d.index)
	for _, word := range text {
		if d.index[word] == nil {
			d.index[word] = []WebPageWithCount{
				WebPageWithCount{
					WebPage: WebPage{
						URL:   url,
						Title: title,
					},
					Count: 1,
				},
			}
			continue
		}

		found := false
		for i, v := range d.index[word] {
			if v.URL == url {
				d.index[word][i].Count++
				break
			}
		}

		if found == true {
			continue
		}

		d.index[word] = append(d.index[word], WebPageWithCount{
			WebPage: WebPage{
				URL:   url,
				Title: title,
			},
			Count: 1,
		})
	}
	return len(d.index) - originalIndexSize
}
