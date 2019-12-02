package main

// WebPage represents info about what is rendered on
// a single url on the internet
type WebPage struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// WebPageWithCount represents a single page on the internet
// with the count of some associated word, this is usually used in
// a mapping where you map word --> WebPageWithCount
type WebPageWithCount struct {
	*WebPage
	Count uint `json:"count"`
}

// SortableWebPagesWithCount is a slice of WebPageWithCounts that will
// counterintuitively sort with the highest count first
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
