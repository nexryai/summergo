package summergo

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestAnalyzeNode(t *testing.T) {
	// テスト用のHTMLノードを作成
	htmlString := "<html><head><title>Test Title</title></head><body><div class=\"example\">Content</div></body></html>"
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatal(err)
	}

	// タイトル取得
	result := analyzeNode(doc, &findParam{tagName: "title"})
	expected := "Test Title"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// 属性を検索して取得する
	result = analyzeNode(doc, &findParam{tagName: "div", attrKey: "class", attrValue: "example", targetKey: "class"})
	expected = "example"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// 存在しない要素
	result = analyzeNode(doc, &findParam{tagName: "nonexistent"})
	if result != "" {
		t.Errorf("Expected empty result, Got: %s", result)
	}
}
