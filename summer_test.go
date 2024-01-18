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
	ExpectActivityPub    bool
	ExpectPlayer         bool
}

var summarizeTests = []summarizeTest{
	{Url: "https://www.google.com/", DescriptionWillEmpty: true},
	{Url: "https://social.sda1.net/"},
	{Url: "https://log.sda1.net/blog/how-to-use-rootless-docker/"},
	{Url: "https://nyan.sda1.net/notes/9oi4vq8a27", ExpectActivityPub: true},
	{Url: "https://sda1.net:3000", ExpectError: true},
	{Url: "https://www.youtube.com/watch?v=zK-RUYiYLok", ExpectPlayer: true},
	{Url: "https://docs.gofiber.io/api/middleware/cache/"},
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
		} else if summary.ActivityPub == "" && test.ExpectActivityPub {
			t.Errorf("activitypub should not be empty: %v", summary)
		} else if summary.Player.Url == "" && test.ExpectPlayer {
			t.Errorf("player should not be empty: %v", summary)
		}

		// Playerのテスト
		if summary.Player.Url != "" {
			if summary.Player.Width == 0 {
				t.Errorf("player width should not be empty: %v", summary)
			}
			if summary.Player.Height == 0 {
				t.Errorf("player height should not be empty: %v", summary)
			}

			for _, permission := range summary.Player.IframePermissions {
				if permission == "accelerometer" || permission == "gyroscope" {
					t.Errorf("unsafe permissions should not be allowed: %s", permission)
				}
			}
		}
	}
}
