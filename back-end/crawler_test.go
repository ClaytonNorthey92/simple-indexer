// +build integration

package main

import (
	"fmt"
	"testing"
)

func Test_CanCrawlStackOverflow(t *testing.T) {
	si := NewPartialWordIndexer()
	sc := NewSimpleCrawler()

	sc.Crawl("https://stackoverflow.com/questions/1998681/xpath-selection-by-innertext", si)

	w := si.GetPagesForWord("mychild")
	fmt.Printf("%+v\n", w)
}
