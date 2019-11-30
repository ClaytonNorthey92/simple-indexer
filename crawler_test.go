package main

import (
	"fmt"
	"testing"
)

func Test_CanCrawlGithub(t *testing.T) {
	t.Skip()
	si := NewSimpleIndexer()
	sc := NewSimpleCrawler()

	sc.Crawl("https://stackoverflow.com/questions/1998681/xpath-selection-by-innertext", si)

	w := si.GetPagesForWord("mychild")
	fmt.Printf("%+v\n", w)
}
