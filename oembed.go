package summergo

import (
	"regexp"
	"strings"
)

type oembed struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Html   string `json:"html"`
}

func getRequiredPermissionsFromIframe(iframe string) []string {
	pattern := `allow\s*=\s*["']([^"']+?)["']`

	// 正規表現にマッチする部分を検索
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(iframe, -1)

	// 結果を取得
	var allowList []string
	for _, match := range matches {
		// セミコロンで分割し、トリムしてから追加
		attributes := strings.Split(match[1], ";")
		for _, attr := range attributes {
			trimmedAttr := strings.TrimSpace(attr)
			if trimmedAttr != "" {
				allowList = append(allowList, trimmedAttr)
			}
		}
	}

	return allowList
}
