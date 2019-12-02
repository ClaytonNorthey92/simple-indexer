package main

import (
	"sort"
	"testing"

	"github.com/go-test/deep"
)

func Test_ShouldBeAbleToSortWebPagesByCount(t *testing.T) {
	webPages := SortableWebPagesWithCount{
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://google.com",
				Title: "google home page",
			},
			Count: 0,
		},
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://yahoo.com",
				Title: "yahoo home page",
			},
			Count: 4,
		},
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://reddit.com",
				Title: "reddit home page",
			},
			Count: 3,
		},
	}

	sort.Sort(webPages)
	diff := deep.Equal(webPages, SortableWebPagesWithCount{
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://yahoo.com",
				Title: "yahoo home page",
			},
			Count: 4,
		},
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://reddit.com",
				Title: "reddit home page",
			},
			Count: 3,
		},
		WebPageWithCount{
			WebPage: &WebPage{
				URL:   "http://google.com",
				Title: "google home page",
			},
			Count: 0,
		},
	})

	if diff != nil {
		t.Error(diff)
	}
}
