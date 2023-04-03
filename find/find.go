package find

import (
	"bytes"
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

type Feed struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func Find(target *url.URL) ([]Feed, error) {
	log.SetPrefix("[" + target.String() + "]")

	// find in HTML
	fromPage, err := tryPageSource(target.String())
	if err != nil {
		log.Printf("%s: %s\n", "parse page", err)
	}
	if len(fromPage) != 0 {
		return fromPage, nil
	}

	// find well-known under this url
	fromWellKnown, err := tryWellKnown(target.Scheme + "://" + target.Host + target.Path) // https://go.dev/play/p/dVt-47_XWjU
	if err != nil {
		log.Printf("%s: %s\n", "parse wellknown", err)
	}
	if len(fromWellKnown) != 0 {
		return fromWellKnown, nil
	}

	// find well-known under url root
	fromWellKnown, err = tryWellKnown(target.Scheme + "://" + target.Host)
	if err != nil {
		log.Printf("%s: %s\n", "parse wellknown under root", err)
	}
	return fromWellKnown, err
}

func tryPageSource(link string) ([]Feed, error) {
	resp, err := request(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %d", resp.StatusCode)
	}

	feeds, err := parseHTMLResp(content)
	if err != nil {
		log.Printf("parse html resp: %s\n", err)
	}
	if len(feeds) != 0 {
		for i := range feeds {
			f := &feeds[i]
			f.Link = formatLinkToAbs(link, f.Link)
		}
		return feeds, nil
	}

	feed, err := parseRSSResp(content)
	if err != nil {
		log.Printf("parse rss resp: %s\n", err)
	}
	if !isEmptyFeed(feed) {
		if feed.Link == "" {
			feed.Link = link
		}
		return []Feed{feed}, nil
	}

	return nil, nil
}

func tryWellKnown(baseURL string) ([]Feed, error) {
	wellKnown := []string{
		"atom.xml",
		"feed.xml",
		"rss.xml",
		"index.xml",
		"atom.json",
		"feed.json",
		"rss.json",
		"index.json",
		"feed/",
		"rss/",
	}
	feeds := make([]Feed, 0)

	for _, suffix := range wellKnown {
		newTarget, err := url.JoinPath(baseURL, suffix)
		if err != nil {
			continue
		}
		parse := func(target string) (Feed, error) { // func for defer close resp.Body
			resp, err := request(newTarget)
			if err != nil {
				return Feed{}, err
			}
			defer resp.Body.Close()
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return Feed{}, err
			}
			return parseRSSResp(content)
		}
		feed, err := parse(newTarget)
		if err != nil {
			continue
		}
		if !isEmptyFeed(feed) {
			feed.Link = newTarget // this may be more accurate than the link parsed from the rss content
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil
}

func parseRSSResp(content []byte) (Feed, error) {
	parsed, err := gofeed.NewParser().Parse(bytes.NewReader(content))
	if err != nil || parsed == nil {
		return Feed{}, err
	}
	return Feed{
		// https://github.com/mmcdole/gofeed#default-mappings
		Title: parsed.Title,

		// set as default value, but the value parsed from rss are not always accurate.
		// it is better to use the url that gets the rss content.
		Link: parsed.FeedLink,
	}, nil
}

func parseHTMLResp(content []byte) ([]Feed, error) {
	feeds := make([]Feed, 0)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	pageTitle := doc.FindMatcher(goquery.Single("title")).Text()

	// find <link> type rss in <header>
	linkExprs := []string{
		"link[type='application/rss+xml']",
		"link[type='application/atom+xml']",
		"link[type='application/json']",
		"link[type='application/feed+json']",
	}
	for _, expr := range linkExprs {
		doc.Find("head").Find(expr).Each(func(i int, s *goquery.Selection) {
			feed := Feed{}
			feed.Title, _ = s.Attr("title")
			feed.Link, _ = s.Attr("href")

			if feed.Title == "" {
				feed.Title = pageTitle
			}
			feeds = append(feeds, feed)
		})
	}

	// find <a> type rss in <body>
	aExpr := "a:contains('rss')"
	doc.Find("body").Find(aExpr).Each(func(i int, s *goquery.Selection) {
		feed := Feed{}
		feed.Title = s.Text()
		feed.Link, _ = s.Attr("href")

		feeds = append(feeds, feed)
	})

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

	return client.Do(req)
}

func isEmptyFeed(feed Feed) bool {
	return feed == Feed{}
}

func formatLinkToAbs(base, link string) string {
	if link == "" {
		return base
	}
	linkURL, err := url.Parse(link)
	if err != nil {
		return link
	}
	if linkURL.IsAbs() {
		return link
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return link
	}
	return baseURL.ResolveReference(linkURL).String()
}
