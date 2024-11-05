package summergo

import (
	"errors"
	"fmt"
	"github.com/nexryai/archer"
	"strings"
	"testing"
	"unicode/utf8"
)

type summarizeTest struct {
	Url                  string
	ExpectUrlError       bool
	TitleWillEmpty       bool
	DescriptionWillEmpty bool
	ExpectActivityPub    bool
	ExpectPlayer         bool
}

var summarizeTests = []summarizeTest{
	{Url: "https://www.google.com/", DescriptionWillEmpty: true},
	{Url: "https://log.sda1.net/blog/how-to-use-rootless-docker/"},
	// ActivityPub
	{Url: "https://misskey.io/notes/97itm23ctg", ExpectActivityPub: true},
	// プライベートIP、一般的でないポートは弾かれる
	{Url: "http://127.0.0.1", ExpectUrlError: true},
	{Url: "https://192.168.1.1", ExpectUrlError: true},
	{Url: "https://sda1.net:3000", ExpectUrlError: true},
	// Player
	{Url: "https://www.youtube.com/watch?v=zK-RUYiYLok", ExpectPlayer: true},
	{Url: "https://www.youtube.com/watch?v=KdbnaBhJs6Y", ExpectPlayer: true},
	// shift-jis 1
	{Url: "https://www.itmedia.co.jp/mobile/articles/2401/18/news172.html"},
	// shift-jis 2
	// {Url: "https://akizukidenshi.com/catalog/contents2/news.aspx"},
	// shift-jis 3
	{Url: "https://www.clas.kitasato-u.ac.jp/~ogawa/C/C01.html", DescriptionWillEmpty: true},
	{Url: "https://www.clas.kitasato-u.ac.jp/~ogawa/C/C02.html", DescriptionWillEmpty: true},
	{Url: "https://www.clas.kitasato-u.ac.jp/~ogawa/C/C03.html", DescriptionWillEmpty: true},
	// EUC-JP
	{Url: "https://map.japanpost.jp/p/search/?&cond200=1&"},
	// 中国語
	{Url: "https://hsr.hoyoverse.com/zh-cn/home"},
	{Url: "https://hsr.hoyoverse.com/zh-tw/home"},
	// 韓国語
	{Url: "https://genshin.hoyoverse.com/ko"},
	//{Url: "https://twitter.com/honkaistarrail/status/1691299712450826240"},
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

		if err != nil && !test.ExpectUrlError {
			t.Errorf("failed to summarize: %v", err)
		} else if err != nil && test.ExpectUrlError {
			if !errors.Is(err, archer.ErrUnsafeUrlDetected) {
				t.Errorf("unexpected error: %v", err)
			}
			continue
		} else if test.ExpectUrlError {
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

		// Replace youtube.com with youtube-nocookie.com
		if strings.HasPrefix(summary.Url, "https://www.youtube.com/watch?v=") && !strings.HasPrefix(summary.Player.Url, "https://www.youtube-nocookie.com/embed/") {
			t.Errorf("youtube.com in summary.Player.Url should be replaced by youtube-nocookie.com: %v", summary)
		}

		// Thumbnailが絶対URLか
		if strings.HasPrefix(summary.Thumbnail, "/") {
			t.Errorf("thumbnail should be absolute url: %v", summary)
		}

		// テキストがUTF-8か
		if summary.Title != "" && !utf8.ValidString(summary.Title) {
			t.Errorf("title should be utf-8: %v", summary)
		} else if summary.Description != "" && !utf8.ValidString(summary.Description) {
			t.Errorf("description should be utf-8: %v", summary)
		} else if summary.SiteName != "" && !utf8.ValidString(summary.SiteName) {
			t.Errorf("sitename should be utf-8: %v", summary)
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
