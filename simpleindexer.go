package main

import (
	"fmt"
)

type WebPageWithCount struct {
	page  *WebPage
	count uint
}

// SimpleIndexer is the in-memory index for word --> web page mapping
type SimpleIndexer struct {
	webPages []WebPage
	urls     map[string]*WebPage
	words    map[string]map[string]*WebPageWithCount
}

func NewSimpleIndexer() *SimpleIndexer {
	return &SimpleIndexer{
		webPages: []WebPage{},
		urls:     map[string]*WebPage{},
		words:    map[string]map[string]*WebPageWithCount{},
	}
}

func (s *SimpleIndexer) PagesForWord(word string) []WebPageWithCount {
	foundWord := s.words[word]
	if foundWord == nil {
		return []WebPageWithCount{}
	}

	w := make([]WebPageWithCount, len(foundWord))
	i := 0
	for _, v := range foundWord {
		w[i] = *v
		i++
	}

	return w
}

func (s *SimpleIndexer) IndexWord(word string, page WebPage) error {
	webPage := s.urls[page.URL]
	if webPage == nil {
		s.webPages = append(s.webPages, page)
		s.urls[page.URL] = &s.webPages[len(s.webPages)-1]
		webPage = s.urls[page.URL]
	} else if webPage.Title != page.Title {
		return fmt.Errorf("mismatch in title, trying to index %v, found %v", page, webPage)
	}

	if s.words[word] == nil {
		s.words[word] = map[string]*WebPageWithCount{
			page.URL: &WebPageWithCount{
				page:  webPage,
				count: 1,
			},
		}
	} else if s.words[word][webPage.URL] == nil {
		s.words[word][webPage.URL] = &WebPageWithCount{
			page:  webPage,
			count: 1,
		}
	} else {
		s.words[word][webPage.URL].count++
	}

	return nil
}
