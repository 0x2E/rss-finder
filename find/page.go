package find

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

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
