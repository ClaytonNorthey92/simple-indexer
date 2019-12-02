package main

import (
	_ "net/http/pprof"
	"regexp"
	"sort"
	"strings"
)

var spaceRe = regexp.MustCompile(` |\n|\t+`)

// PartialWordIndexer is an Indexer that will store a mapping of all
// "partial" words to their count on a web page
// the word "blah" would map to the following partials:
// b, bl, bla, blah, l, la, lah, a, ah, h
type PartialWordIndexer struct {
	index   map[string]SortableWebPagesWithCount
	lock    chan bool
	pageMap map[string]*WebPage
}

// NewPartialWordIndexer creates a new PartialWordIndexer and returns a reference
// to it
func NewPartialWordIndexer() *PartialWordIndexer {
	d := &PartialWordIndexer{
		index:   map[string]SortableWebPagesWithCount{},
		lock:    make(chan bool, 1),
		pageMap: map[string]*WebPage{},
	}
	d.unlock()
	return d
}

func (d *PartialWordIndexer) waitForLock() {
	<-d.lock
}

func (d *PartialWordIndexer) unlock() {
	d.lock <- true
}

// Clear will clear the old index and use a new one
// WARNING: there is a current memory leak in this program with the garbage collector
// not clearing out the old index when no longer referenced
func (d *PartialWordIndexer) Clear() {
	d.waitForLock()
	d.index = map[string]SortableWebPagesWithCount{}
	d.pageMap = map[string]*WebPage{}
	d.unlock()
}

// GetPagesForWord gets a slice of sorted web pages with the count of word occurances
// in them
func (d *PartialWordIndexer) GetPagesForWord(word string) SortableWebPagesWithCount {
	word = strings.ToLower(word)
	pages := SortableWebPagesWithCount{}
	pageMap := map[string]*WebPageWithCount{}
	d.waitForLock()
	for key, webPages := range d.index {
		if strings.Contains(key, word) == false {
			continue
		}

		for _, webPage := range webPages {
			if pageMap[webPage.URL] == nil {
				pageMap[webPage.URL] = &WebPageWithCount{
					WebPage: &WebPage{
						URL:   webPage.URL,
						Title: webPage.Title,
					},
					Count: webPage.Count,
				}
			} else {
				pageMap[webPage.URL].Count += webPage.Count
			}
		}
	}
	d.unlock()
	for _, v := range pageMap {
		pages = append(pages, *v)
	}
	sort.Sort(pages)
	return pages
}

// IndexTextForPage will take some content and index each partial word found in it with the count of
// occurances
func (d *PartialWordIndexer) IndexTextForPage(pageContent string, url string, title string) int {
	d.waitForLock()
	cachedWebPage := d.pageMap[url]
	if cachedWebPage == nil {
		d.pageMap[url] = &WebPage{
			URL:   url,
			Title: title,
		}
		cachedWebPage = d.pageMap[url]
	}

	text := spaceRe.Split(pageContent, -1)
	pageContent = strings.ToLower(pageContent)

	for _, word := range text {
		word = strings.Trim(word, " \n\t")
		if word == "" {
			continue
		}
		word = strings.ToLower(word)
		if d.index[word] == nil {
			d.index[word] = []WebPageWithCount{
				WebPageWithCount{
					WebPage: cachedWebPage,
					Count:   1,
				},
			}
			continue
		}

		found := false
		for i, v := range d.index[word] {
			if v.URL == url {
				d.index[word][i].Count++
				found = true
				break
			}
		}

		if found == true {
			continue
		}

		d.index[word] = append(d.index[word], WebPageWithCount{
			WebPage: cachedWebPage,
			Count:   1,
		})
	}
	d.unlock()
	return len(text)
}
