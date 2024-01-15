package summergo

import (
	"errors"
	"fmt"
	"github.com/nexryai/archer"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

func getPageTitle(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		&findParam{tagName: "meta", attrKey: "property", attrValue: "og:title", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "name", attrValue: "twitter:title", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "property", attrValue: "twitter:title", targetKey: "content"},
		&findParam{tagName: "title"},
	}...)
}

func getPageDescription(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		&findParam{tagName: "meta", attrKey: "property", attrValue: "og:description", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "name", attrValue: "twitter:description", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "property", attrValue: "twitter:description", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "name", attrValue: "description", targetKey: "content"},
	}...)
}

func getPageImage(doc *html.Node) string {
	return analyzeNode(doc, []*findParam{
		&findParam{tagName: "meta", attrKey: "property", attrValue: "og:image", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "name", attrValue: "twitter:image", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "property", attrValue: "twitter:image", targetKey: "content"},
		&findParam{tagName: "link", attrKey: "rel", attrValue: "image_src", targetKey: "href"},
		&findParam{tagName: "link", attrKey: "rel", attrValue: "apple-touch-icon", targetKey: "href"},
		&findParam{tagName: "link", attrKey: "rel", attrValue: "apple-touch-icon image_src", targetKey: "href"},
	}...)
}

func GetSiteName(doc *html.Node, siteUrl string) string {
	res := analyzeNode(doc, []*findParam{
		&findParam{tagName: "meta", attrKey: "property", attrValue: "og:site_name", targetKey: "content"},
		&findParam{tagName: "meta", attrKey: "name", attrValue: "twitter:site", targetKey: "content"},
	}...)

	if res == "" {
		u, urlErr := url.Parse(siteUrl)
		if urlErr != nil {
			return ""
		} else {
			res = u.Host
		}
	}

	return res
}

func GetFavicon(doc *html.Node, siteUrl string) string {
	res := analyzeNode(doc, []*findParam{
		&findParam{tagName: "link", attrKey: "rel", attrValue: "shortcut icon", targetKey: "href"},
		&findParam{tagName: "link", attrKey: "rel", attrValue: "icon", targetKey: "href"},
	}...)

	if res == "" {
		u, urlErr := url.Parse(siteUrl)
		if urlErr != nil {
			return ""
		} else {
			res = fmt.Sprintf("https://%s/favicon.ico", u.Host)
		}
	}

	return res
}

func Summarize(url string) (*Summary, error) {
	req, newReqErr := http.NewRequest("GET", url, nil)

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

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse html")
	}

	return &Summary{
		Title:       getPageTitle(doc),
		Description: getPageDescription(doc),
		Thumbnail:   getPageImage(doc),
		SiteName:    GetSiteName(doc, url),
		Icon:        GetFavicon(doc, url),
	}, nil
}
