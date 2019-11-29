package main

type WebPageWithCount struct {
	WebPage
	Count uint
}

type SortableWebPagesWithCount []WebPageWithCount

func (s SortableWebPagesWithCount) Len() int {
	return len(s)
}

func (s SortableWebPagesWithCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableWebPagesWithCount) Less(i, j int) bool {
	return s[i].Count > s[j].Count
}

type FastLookupIndexer struct {
	index map[string][]WebPageWithCount
}

func NewFastLookupIndexer() *FastLookupIndexer {
	return &FastLookupIndexer{
		index: map[string][]WebPageWithCount{},
	}
}

func (s *FastLookupIndexer) GetPagesForWord(word string) []WebPageWithCount {
	return s.index[word]
}

func (s *FastLookupIndexer) IndexTextForPage(pageContent string, url string, title string) {
	contentSize := len(pageContent)
	for i, _ := range pageContent {
		for step := i; step < contentSize; step++ {
			word := pageContent[i : step+1]

			if s.index[word] == nil {
				s.index[word] = []WebPageWithCount{
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

			exists := false
			for i, v := range s.index[word] {
				if v.WebPage.URL == url {
					s.index[word][i].Count++
					exists = true
					continue
				}
			}

			if exists == true {
				continue
			}

			s.index[word] = append(s.index[word], WebPageWithCount{
				WebPage: WebPage{
					URL:   url,
					Title: title,
				},
				Count: 1,
			})

		}
	}
}
