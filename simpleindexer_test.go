package main

import (
	"testing"

	"github.com/go-test/deep"
)

func Test_RetrivePageFromWord(t *testing.T) {
	si := NewSimpleIndexer()

	err := si.IndexWord("google", WebPage{
		URL:   "www.google.com",
		Title: "Google Home Page",
	})

	if err != nil {
		t.Error(err)
	}
}

func Test_ShouldNotAllowMismatchedTitles(t *testing.T) {
	si := NewSimpleIndexer()

	err := si.IndexWord("google", WebPage{
		URL:   "www.google.com",
		Title: "Google Home Page",
	})

	if err != nil {
		t.Error(err)
	}

	err = si.IndexWord("search", WebPage{
		URL:   "www.google.com",
		Title: "Google Other Page",
	})

	if err == nil || err.Error() != "mismatch in title, trying to index {www.google.com Google Other Page}, found &{www.google.com Google Home Page}" {
		t.Error(err)
	}
}

func Test_ShouldBeAbleToCoundWordsOnPages(t *testing.T) {
	si := NewSimpleIndexer()

	si.IndexWord("google", WebPage{
		URL:   "www.google.com",
		Title: "Google Home Page",
	})

	si.IndexWord("google", WebPage{
		URL:   "www.google.com",
		Title: "Google Home Page",
	})

	si.IndexWord("search", WebPage{
		URL:   "www.google.com",
		Title: "Google Home Page",
	})

	si.IndexWord("google", WebPage{
		URL:   "www.yahoo.com",
		Title: "Yahoo Home Page",
	})

	webPages := si.PagesForWord("google")
	if len(webPages) != 2 {
		t.Error("incorrect number of pages")
	}

	expectedResults := []WebPageWithCount{
		WebPageWithCount{
			page: &WebPage{
				URL:   "www.google.com",
				Title: "Google Home Page",
			},
			count: 2,
		},
		WebPageWithCount{
			page: &WebPage{
				URL:   "www.yahoo.com",
				Title: "Yahoo Home Page",
			},
			count: 1,
		},
	}

	deep.CompareUnexportedFields = true
	if diff := deep.Equal(expectedResults, webPages); diff != nil {
		t.Error(diff)
	}
}
