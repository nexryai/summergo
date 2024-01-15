package summergo

import (
	"golang.org/x/net/html"
)

type findParam struct {
	tagName   string
	attrKey   string
	attrValue string
	targetKey string
}

// 属性の値の検索
func getAttributeValue(node *html.Node, attrKey string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}

func analyzeNode(node *html.Node, find ...*findParam) string {
	for _, f := range find {
		if f.tagName == "title" {
			if node.Type == html.ElementNode && node.Data == f.tagName {
				return node.FirstChild.Data
			} else {
				continue
			}
		}

		if node.Type == html.ElementNode && node.Data == f.tagName {
			for _, attr := range node.Attr {
				if attr.Key == f.attrKey && attr.Val == f.attrValue {
					return getAttributeValue(node, f.targetKey)
				}
			}
		}
	}

	// 子ノードを再帰的に検索
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := analyzeNode(child, find...); result != "" {
			return result
		}
	}

	return ""
}
