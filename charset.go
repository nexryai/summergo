package summergo

import (
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"strings"
)

func isShiftJis(doc *html.Node) bool {
	contentType := analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "http-equiv", attrValue: "Content-Type", targetKey: "content"},
		{tagName: "meta", attrKey: "http-equiv", attrValue: "content-type", targetKey: "content"},
		{tagName: "meta", attrKey: "charset", attrValue: "Shift_JIS", targetKey: "charset"},
		{tagName: "meta", attrKey: "charset", attrValue: "shift_jis", targetKey: "charset"},
	}...)

	if strings.Contains(strings.ToLower(contentType), "shift_jis") {
		return true
	} else {
		return false
	}
}

func convertShiftJisToUtf8(str string) string {
	encoder := japanese.ShiftJIS.NewDecoder()

	// 文字列をShift JISからUTF-8に変換
	utf8Bytes, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return ""
	}

	// UTF-8のバイト列を文字列に変換
	utf8Str := string(utf8Bytes)

	return utf8Str
}

func convertEucJpToUtf8(str string) string {
	encoder := japanese.EUCJP.NewDecoder()

	// 文字列をEUC-JPからUTF-8に変換
	utf8Bytes, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return ""
	}

	// UTF-8のバイト列を文字列に変換
	utf8Str := string(utf8Bytes)

	return utf8Str
}
