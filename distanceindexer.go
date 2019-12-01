package main

import (
	"regexp"
	"sort"
	"strings"
)

type DistanceIndexer struct {
	index map[string]SortableWebPagesWithCount
	lock  chan bool
}

func NewDistanceIndexer() *DistanceIndexer {
	d := &DistanceIndexer{
		index: map[string]SortableWebPagesWithCount{},
		lock:  make(chan bool, 1),
	}
	d.Unlock()
	return d
}

func (d *DistanceIndexer) WaitForLock() {
	<-d.lock
}

func (d *DistanceIndexer) Unlock() {
	d.lock <- true
}

func (d *DistanceIndexer) GetPagesForWord(word string) SortableWebPagesWithCount {
	word = strings.ToLower(word)
	d.WaitForLock()
	foundPages := d.index[word]
	if foundPages == nil {
		d.Unlock()
		return SortableWebPagesWithCount{}
	}
	sort.Sort(foundPages)
	d.Unlock()
	return foundPages
}

func (d *DistanceIndexer) IndexTextForPage(pageContent string, url string, title string) int {
	re := regexp.MustCompile(` |\n|\t+`)
	text := re.Split(pageContent, -1)
	wordsIndexed := 0
	for _, word := range text {
		if word == "" {
			continue
		}
		partials := allPermutationsOfWord(word)
		for _, partial := range partials {
			d.WaitForLock()
			partial = strings.ToLower(partial)
			if d.index[partial] == nil {
				d.index[partial] = []WebPageWithCount{
					WebPageWithCount{
						WebPage: WebPage{
							URL:   url,
							Title: title,
						},
						Count: 1,
					},
				}
				d.Unlock()
				continue
			}

			found := false
			for i, v := range d.index[partial] {
				if v.URL == url {
					d.index[partial][i].Count++
					found = true
					break
				}
			}

			if found == true {
				d.Unlock()
				continue
			}

			d.index[partial] = append(d.index[partial], WebPageWithCount{
				WebPage: WebPage{
					URL:   url,
					Title: title,
				},
				Count: 1,
			})
			d.Unlock()
		}
		wordsIndexed++
	}
	return wordsIndexed
}

func allPermutationsOfWord(word string) []string {
	var length int
	for i := len(word); i > 0; i-- {
		length += i
	}

	partials := make([]string, length)

	partialIndex := 0
	for i := range word {
		for s := i; s < len(word); s++ {
			partialWord := word[i : s+1]
			partials[partialIndex] = partialWord
			partialIndex++
		}
	}

	return partials
}
