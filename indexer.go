package main

// WebPage represents info about what is rendered on
// a single url on the internet
type WebPage struct {
	URL     string
	Title   string
	Content string
	Count   string
}

// Indexer stores an index of web pages and counts of words
// within each one
type Indexer interface {
	// map of word --> web page
	GetPagesForWord(word string) SortableWebPagesWithCount
	IndexTextForPage(pageContent string, url string, title string)
}
