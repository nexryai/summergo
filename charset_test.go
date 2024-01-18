package summergo

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestIsShiftJis(t *testing.T) {
	// 滅ぼすべきサイト
	htmlString := `<html>
					  <head>
						<meta http-equiv="Content-Type" content="text/html; charset=Shift_JIS">
					  </head>
					</html>`
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatal(err)
	}

	if !isShiftJis(doc) {
		t.Errorf("Expected: true, Got: false")
	}

	// Shift JISでないコンテンツのHTML
	htmlString = `<html>
					  <head>
						<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
					  </head>
					</html>`
	doc, err = html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatal(err)
	}

	if isShiftJis(doc) {
		t.Errorf("Expected: false, Got: true")
	}
}

func TestConvertShiftJisToUtf8(t *testing.T) {
	shiftJisStr := "\x82\xa0\x82\xa2\x82\xa4" // "あいう"
	result := convertShiftJisToUtf8(shiftJisStr)
	expected := "あいう"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}
