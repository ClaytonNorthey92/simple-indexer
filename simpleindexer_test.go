package main

import (
	"testing"

	"github.com/go-test/deep"
)

func Test_ShouldBeAbleToIndexMultipleOccurancesOfAWordSimple(t *testing.T) {
	si := NewSimpleIndexer()

	si.IndexTextForPage("hi these are words, these are very cool", "https://google.com", "Google Home Page")

	pages := si.GetPagesForWord("these")

	diff := deep.Equal(pages, SortableWebPagesWithCount{
		WebPageWithCount{
			WebPage: WebPage{
				URL:   "https://google.com",
				Title: "Google Home Page",
			},
			Count: 2,
		},
	})

	if diff != nil {
		t.Error(diff)
	}
}

func Test_ShouldBeAbleToIndexAPartialWordSimple(t *testing.T) {
	si := NewSimpleIndexer()

	si.IndexTextForPage("I like google, and googley things, while I google", "https://google.com", "Google Home Page")

	pages := si.GetPagesForWord("google")
	diff := deep.Equal(pages, SortableWebPagesWithCount{
		WebPageWithCount{
			WebPage: WebPage{
				URL:   "https://google.com",
				Title: "Google Home Page",
			},
			Count: 3,
		},
	})

	if diff != nil {
		t.Error(diff)
	}
}

func Test_ShouldBeAbleToIndexMulitplePagesSimple(t *testing.T) {
	si := NewSimpleIndexer()

	si.IndexTextForPage("I like google, and googley things, while I google things", "https://google.com", "Google Home Page")
	si.IndexTextForPage("I hate google.  Yahoo is better.  way better at everything.", "https://yahoo.com", "Yahoo Home Page")

	pages := si.GetPagesForWord("thing")
	diff := deep.Equal(pages, SortableWebPagesWithCount{
		WebPageWithCount{
			WebPage: WebPage{
				URL:   "https://google.com",
				Title: "Google Home Page",
			},
			Count: 2,
		},
		WebPageWithCount{
			WebPage: WebPage{
				URL:   "https://yahoo.com",
				Title: "Yahoo Home Page",
			},
			Count: 1,
		},
	})

	if diff != nil {
		t.Error(diff)
	}
}
