package main

// Indexer stores an index of web pages and counts of words
// within each one
type Indexer interface {
	GetPagesForWord(word string) SortableWebPagesWithCount
	IndexTextForPage(pageContent string, url string, title string) int
	Clear()
}
