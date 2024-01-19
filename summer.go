package summergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nexryai/archer"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

func getPageTitle(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "property", attrValue: "og:title", targetKey: "content"},
		{tagName: "meta", attrKey: "name", attrValue: "twitter:title", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:title", targetKey: "content"},
		{tagName: "title"},
	}...)
}

func getPageDescription(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "property", attrValue: "og:description", targetKey: "content"},
		{tagName: "meta", attrKey: "name", attrValue: "twitter:description", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:description", targetKey: "content"},
		{tagName: "meta", attrKey: "name", attrValue: "description", targetKey: "content"},
	}...)
}

func getPageImage(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "property", attrValue: "og:image", targetKey: "content"},
		{tagName: "meta", attrKey: "name", attrValue: "twitter:image", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:image", targetKey: "content"},
		{tagName: "link", attrKey: "rel", attrValue: "image_src", targetKey: "href"},
		{tagName: "link", attrKey: "rel", attrValue: "apple-touch-icon", targetKey: "href"},
		{tagName: "link", attrKey: "rel", attrValue: "apple-touch-icon image_src", targetKey: "href"},
	}...)
}

func getPlayerFromOEmbed(doc *html.Node) *Player {
	oembedUrl := analyzeNode(doc, []*findParam{
		{tagName: "link", attrKey: "type", attrValue: "application/json+oembed", targetKey: "href"},
	}...)

	if oembedUrl == "" {
		return nil
	}

	// OEmbedを取得する
	req, newReqErr := http.NewRequest("GET", oembedUrl, nil)
	if newReqErr != nil {
		return nil
	}

	req.Header.Set("User-Agent", "SummerGo/0.1")

	requester := archer.SecureRequest{
		Request: req,
		TimeOut: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, respErr := requester.Send()

	if respErr != nil || resp.StatusCode != 200 {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	embed := &oembed{}
	err = json.Unmarshal(body, embed)
	if err != nil {
		return nil
	}

	// OEmbedのiframeが要求する権限のうち安全なものを許可する
	var allowed []string
	safePermissions := []string{"autoplay", "clipboard-write", "picture-in-picture", "web-share", "fullscreen"}
	for _, rp := range getRequiredPermissionsFromIframe(embed.Html) {
		for _, sp := range safePermissions {
			if rp == sp {
				allowed = append(allowed, rp)
			}
		}
	}

	return &Player{
		Url:               getPlayerUrl(doc),
		Width:             embed.Width,
		Height:            embed.Height,
		IframePermissions: allowed,
	}

}

func getPlayerUrl(doc *html.Node) string {
	twitterCard := analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "name", attrValue: "twitter:card", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:card", targetKey: "content"},
	}...)

	if twitterCard != "summary_large_image" {
		playerUrlFromTwitterCard := analyzeNode(doc, []*findParam{
			{tagName: "meta", attrKey: "name", attrValue: "twitter:player", targetKey: "content"},
			{tagName: "meta", attrKey: "property", attrValue: "twitter:player", targetKey: "content"},
		}...)

		if playerUrlFromTwitterCard != "" {
			return playerUrlFromTwitterCard
		}
	}

	return analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "property", attrValue: "og:video", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "og:video:secure_url", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "og:video:url", targetKey: "content"},
	}...)
}

func getPlayerWidth(doc *html.Node) int {
	widthStr := analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "name", attrValue: "twitter:player:width", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:player:width", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "og:video:width", targetKey: "content"},
	}...)

	w, err := strconv.Atoi(widthStr)
	if err != nil {
		return 0
	}

	return w
}

func getPlayerHeight(doc *html.Node) int {
	heightStr := analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "name", attrValue: "twitter:player:height", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "twitter:player:height", targetKey: "content"},
		{tagName: "meta", attrKey: "property", attrValue: "og:video:height", targetKey: "content"},
	}...)

	h, err := strconv.Atoi(heightStr)
	if err != nil {
		return 0
	}

	return h
}

func getActivityPubLink(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		{tagName: "link", attrKey: "type", attrValue: "application/activity+json", targetKey: "href"},
	}...)
}

func isSensitive(doc *html.Node, parsedUrl url.URL) bool {
	if strings.Contains("mixi.co.jp", parsedUrl.Host) && analyzeNode(doc, []*findParam{{tagName: "meta", attrKey: "property", attrValue: "mixi:content-rating", targetKey: "content"}}...) == "1" {
		return true
	} else {
		return false
	}
}

func getSiteName(doc *html.Node, parsedUrl url.URL) string {
	res := analyzeNode(doc, []*findParam{
		{tagName: "meta", attrKey: "property", attrValue: "og:site_name", targetKey: "content"},
		{tagName: "meta", attrKey: "name", attrValue: "twitter:site", targetKey: "content"},
	}...)

	if res == "" {
		res = parsedUrl.Host
	}

	return res
}

func getFavicon(doc *html.Node, parsedUrl url.URL) string {
	res := analyzeNode(doc, []*findParam{
		{tagName: "link", attrKey: "rel", attrValue: "shortcut icon", targetKey: "href"},
		{tagName: "link", attrKey: "rel", attrValue: "icon", targetKey: "href"},
	}...)

	if res == "" {
		res = fmt.Sprintf("https://%s/favicon.ico", parsedUrl.Host)
	} else if !strings.HasPrefix(res, "https://") {
		res = fmt.Sprintf("https://%s%s", parsedUrl.Host, res)
	}

	return res
}

func SummarizeHtml(siteUrl url.URL, body io.Reader, charSet string) (*Summary, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, errors.New("failed to parse html")
	}

	player := getPlayerFromOEmbed(doc)
	if player == nil {
		player = &Player{
			Url:    getPlayerUrl(doc),
			Width:  getPlayerWidth(doc),
			Height: getPlayerHeight(doc),
		}
	}

	title := getPageTitle(doc)
	description := getPageDescription(doc)
	siteName := getSiteName(doc, siteUrl)

	// shift_jis対策
	if charSet == "" {
		if utf8.ValidString(title) {
			charSet = "utf-8"
		} else {
			charsetDetector := chardet.NewTextDetector()
			charsetResult, err := charsetDetector.DetectBest([]byte(title))
			if err != nil {
				// fallback
				charSet = "utf-8"
			} else {
				charSet = strings.ToLower(charsetResult.Charset)
			}
		}
	}

	// そのうち他の文字コードにも対応する？
	// FixMe: 漢字のみのテキストだとshift-jisがbig5 or GB-18030認識される場合がある
	if charSet == "shift_jis" || charSet == "big5" || charSet == "gb-18030" {
		title = convertShiftJisToUtf8(title)
		description = convertShiftJisToUtf8(description)
		siteName = convertShiftJisToUtf8(siteName)
	}

	return &Summary{
		Url:         siteUrl.String(),
		Title:       title,
		Description: description,
		Thumbnail:   getPageImage(doc),
		SiteName:    siteName,
		Icon:        getFavicon(doc, siteUrl),
		ActivityPub: getActivityPubLink(doc),
		Sensitive:   isSensitive(doc, siteUrl),
		Player:      *player,
	}, nil
}

func Summarize(siteUrl string) (*Summary, error) {
	parsedUrl, err := url.Parse(siteUrl)
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	req, newReqErr := http.NewRequest("GET", siteUrl, nil)
	if newReqErr != nil {
		return nil, errors.New("failed to create request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0 SummerGo/0.1")

	// :)
	if parsedUrl.Host == "twitter.com" || parsedUrl.Host == "x.com" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)")
	}

	requester := archer.SecureRequest{
		Request: req,
		TimeOut: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, respErr := requester.Send()

	if respErr != nil {
		return nil, errors.New("failed to send request")
	} else if resp.StatusCode != 200 {
		return nil, errors.New("non-200 status code: " + resp.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// サーバーからのレスポンスでcharsetを明示しているならそれを使って高速化する
	var knownCharset string
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "utf-8") {
		knownCharset = "utf-8"
	} else if strings.Contains(strings.ToLower(contentType), "shift_jis") {
		knownCharset = "shift_jis"
	}

	return SummarizeHtml(*parsedUrl, resp.Body, knownCharset)
}
