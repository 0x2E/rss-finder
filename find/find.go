package find

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type ErrBadContentType struct {
	T string
}

func (e ErrBadContentType) Error() string {
	return fmt.Sprintf("bad content type %s", e.T)
}

type Feed struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func Find(target *url.URL) ([]Feed, error) {
	feeds := make([]Feed, 0)

	fromPage, err := tryPageSource(target.String())
	if err != nil {
		log.Println(err)
	}
	feeds = append(feeds, fromPage...)

	fromWellKnown, err := tryWellKnown(target)
	if err != nil {
		log.Println(err)
	}
	feeds = append(feeds, fromWellKnown...)

	return feeds, nil
}

func tryPageSource(link string) ([]Feed, error) {
	resp, err := request(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")

	feeds, err := parseHTMLResp(contentType, resp.Body)
	if err != nil {
		log.Println(err)
	}
	if len(feeds) != 0 {
		return feeds, nil
	}

	feed, err := parseRSSResp(contentType, resp.Body)
	if err != nil {
		log.Println(err)
	}
	if feed != (Feed{}) {
		return []Feed{feed}, nil
	}

	return nil, ErrBadContentType{T: contentType}
}

func tryWellKnown(target *url.URL) ([]Feed, error) {
	wellKnown := []string{
		"atom.xml",
		"feed.xml",
		"rss.xml",
		"feed/",
		"rss/",
	}
	feeds := make([]Feed, 0)

	baseURL := target.Scheme + "://" + target.Host + target.Path // https://go.dev/play/p/dVt-47_XWjU
	for _, suffix := range wellKnown {
		newTarget, err := url.JoinPath(baseURL, suffix)
		if err != nil {
			continue
		}
		resp, err := request(newTarget)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		feed, err := parseRSSResp(resp.Header.Get("Content-Type"), resp.Body)
		if err != nil {
			continue
		}
		if feed != (Feed{}) {
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil
}

func parseRSSResp(contentType string, body io.Reader) (Feed, error) {
	var (
		feed Feed
		err  error

		// https://en.wikipedia.org/wiki/Media_type
		rssType = []string{
			"text/plain",
			"application/json",
			"application/rss+xml",
			"application/atom+xml",
			"application/json",
			"application/feed+json",
		}
	)
	for _, t := range rssType {
		if contentType == t {
			feed, err = parseRSSContent(body)
			break
		}
	}
	return feed, err
}

func parseRSSContent(content io.Reader) (Feed, error) {
	parsed, err := gofeed.NewParser().Parse(content)
	if err != nil {
		return Feed{}, err
	}
	if parsed == nil {
		return Feed{}, nil
	}
	return Feed{
		// https://github.com/mmcdole/gofeed#default-mappings
		Title: parsed.Title,
		Link:  parsed.FeedLink,
	}, nil
}

func parseHTMLResp(contentType string, body io.Reader) ([]Feed, error) {
	if contentType != "text/html" && contentType != "text/plain" {
		return nil, ErrBadContentType{T: contentType}
	}

	feeds := make([]Feed, 0)

	exprs := []string{
		"link[type='application/rss+xml']",
		"link[type='application/atom+xml']",
		"link[type='application/json']",
		"link[type='application/feed+json']",
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	pageTitle := doc.FindMatcher(goquery.Single("title")).Text()

	for _, expr := range exprs {
		doc.Find(expr).Each(func(i int, s *goquery.Selection) {
			feed := Feed{}
			feed.Title, _ = s.Attr("title")
			feed.Link, _ = s.Attr("href")

			if feed.Title == "" {
				feed.Title = pageTitle
			}
			feeds = append(feeds, feed)
		})
	}
	return feeds, nil
}

func request(link string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	ua := os.Getenv("USER_AGENT")
	if strings.TrimSpace(ua) == "" {
		ua = "rss-finder/1.0"
	}
	req.Header.Add("User-Agent", ua)

	// todo add Accept header

	return client.Do(req)
}
