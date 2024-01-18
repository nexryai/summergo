package summergo

import (
	"errors"
	"fmt"
	"github.com/nexryai/archer"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
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
	}

	return res
}

func Summarize(siteUrl string) (*Summary, error) {
	parsedUrl, err := url.Parse(siteUrl)
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	req, newReqErr := http.NewRequest("GET", siteUrl, nil)

	// User-Agentを設定
	// ブラウザっぽくするのはお行儀的に微妙かもしれないので変えられるようにする？
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0 SummerGo/0.1")

	requester := archer.SecureRequest{
		Request: req,
		TimeOut: 10,
		MaxSize: 1024 * 1024 * 10,
	}

	resp, respErr := requester.Send()

	if errors.Join(newReqErr, respErr) != nil {
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse html")
	}

	return &Summary{
		Title:       getPageTitle(doc),
		Description: getPageDescription(doc),
		Thumbnail:   getPageImage(doc),
		SiteName:    getSiteName(doc, *parsedUrl),
		Icon:        getFavicon(doc, *parsedUrl),
		ActivityPub: getActivityPubLink(doc),
		Sensitive:   isSensitive(doc, *parsedUrl),
	}, nil
}
