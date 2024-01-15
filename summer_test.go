package summergo

import (
	"fmt"
	"testing"
)

type summarizeTest struct {
	Url                  string
	ExpectError          bool
	TitleWillEmpty       bool
	DescriptionWillEmpty bool
}

var summarizeTests = []summarizeTest{
	summarizeTest{
		Url:                  "https://www.google.com/",
		DescriptionWillEmpty: true,
	},
	{
		Url: "https://social.sda1.net/",
	},
	{
		Url: "https://log.sda1.net/blog/how-to-use-rootless-docker/",
	},
	{
		Url:         "https://sda1.net:3000",
		ExpectError: true,
	},
}

func TestSummarize(t *testing.T) {
	for _, test := range summarizeTests {
		summary, err := Summarize(test.Url)
		if err == nil && summary == nil {
			t.Errorf("err == nil && summary == nil")
		}

		if err == nil {
			fmt.Printf("==== %s =================================\n", test.Url)
			fmt.Printf("Title: %v\n", summary.Title)
			fmt.Printf("Description: %v\n", summary.Description)
			fmt.Printf("Thumbnail: %v\n", summary.Thumbnail)
			fmt.Printf("Icon: %v\n", summary.Icon)
			fmt.Printf("SiteName: %v\n", summary.SiteName)
			fmt.Printf("Player: %v\n", summary.Player)
			fmt.Printf("Sensitive: %v\n", summary.Sensitive)
			fmt.Printf("ActivityPub: %v\n", summary.ActivityPub)
		}

		if err != nil && !test.ExpectError {
			t.Errorf("failed to summarize: %v", err)
		} else if err != nil && test.ExpectError {
			continue
		} else if test.ExpectError {
			t.Errorf("summarize should be failed: %v", summary)
		} else if summary.Title == "" && !test.TitleWillEmpty {
			t.Errorf("title should not be empty: %v", summary)
		} else if summary.Description == "" && !test.DescriptionWillEmpty {
			t.Errorf("description should not be empty: %v", summary)
		}
	}
}
