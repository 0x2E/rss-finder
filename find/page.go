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

	feeds, err := parseHTMLContent(content)
	if err != nil {
		log.Printf("parse html content: %s\n", err)
	}
	if len(feeds) != 0 {
		for i := range feeds {
			f := &feeds[i]
			f.Link = formatLinkToAbs(link, f.Link)
		}
		return feeds, nil
	}

	feed, err := parseRSSContent(content)
	if err != nil {
		log.Printf("parse rss content: %s\n", err)
	}
	if !isEmptyFeed(feed) {
		if feed.Link == "" {
			feed.Link = link
		}
		return []Feed{feed}, nil
	}

	return nil, nil
}

func parseHTMLContent(content []byte) ([]Feed, error) {
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
	suspected := make(map[string]struct{})
	doc.Find("body").Find(aExpr).Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		suspected[link] = struct{}{}
	})
	for link := range suspected {
		feed, err := parseRSSUrl(link)
		if err != nil {
			continue
		}
		if !isEmptyFeed(feed) {
			feed.Link = link // this may be more accurate than the link parsed from the rss content
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil
}
