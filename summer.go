package summergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nexryai/archer"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func SummarizeHtml(siteUrl url.URL, body io.Reader) (*Summary, error) {
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

	return &Summary{
		Url:         siteUrl.String(),
		Title:       getPageTitle(doc),
		Description: getPageDescription(doc),
		Thumbnail:   getPageImage(doc),
		SiteName:    getSiteName(doc, siteUrl),
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

	requester := archer.SecureRequest{
		Request: req,
		TimeOut: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, respErr := requester.Send()

	if respErr != nil {
		return nil, errors.New("failed to send request")
	} else if resp.StatusCode != 200 {
		return nil, errors.New("non-200 status code")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	return SummarizeHtml(*parsedUrl, resp.Body)
}
