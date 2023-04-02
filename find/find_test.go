package find

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTMLResp(t *testing.T) {
	type testItem struct {
		ct   string
		body io.Reader
		feed []Feed
	}

	table := []testItem{
		// no feed title
		{ct: "text/plain", body: strings.NewReader(`
		<html>
		<head>
			<title>html title</title>
			<link type="application/rss+xml" href="https://example.com/x/feed.xml">
		</head>
		<body></body></html>
		`), feed: []Feed{{Title: "html title", Link: "https://example.com/x/feed.xml"}}},

		// match all types
		{ct: "text/html", body: strings.NewReader(`
		<html>
		<head>
			<title>html title</title>
			<link type="application/rss+xml" title="rss+xml feed title" href="https://example.com/x/rss.xml">
			<link type="application/atom+xml" title="atom+xml feed title" href="https://example.com/x/atom.xml">
			<link type="application/json" title="json feed title" href="https://example.com/x/feed.json">
			<link type="application/feed+json" title="feed+json feed title" href="https://example.com/x/feed.json">
		</head>
		<body></body></html>
		`), feed: []Feed{
			{Title: "rss+xml feed title", Link: "https://example.com/x/rss.xml"},
			{Title: "atom+xml feed title", Link: "https://example.com/x/atom.xml"},
			{Title: "json feed title", Link: "https://example.com/x/feed.json"},
			{Title: "feed+json feed title", Link: "https://example.com/x/feed.json"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLResp(tt.ct, tt.body)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.feed, feed)
	}
}

func TestParseRSSResp(t *testing.T) {
	type testItem struct {
		ct      string
		content io.Reader
		feed    Feed
	}

	// todo match all types, e.g. https://github.com/mmcdole/gofeed/tree/master/testdata
	table := []testItem{
		{ct: "application/rss+xml", content: strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">  
		  <channel> 
			<title>test</title>  
			<link>https://example.com/</link>  
			<description>Recent content on test</description>  
			<language>en</language>
			<lastBuildDate>Fri, 24 Feb 2023 00:43:57 +0800</lastBuildDate>  
			<atom:link href="https://example.com/feed.xml" rel="self" type="application/rss+xml"/>  
			<item> 
			  <title>post1</title>  
			  <link>https://example.com/post1/</link>  
			  <pubDate>Fri, 24 Feb 2023 00:43:57 +0800</pubDate>  
			  <guid>https://example.com/post1/</guid>  
			  <description>post1 content</description> 
			</item>  
		  </channel> 
		</rss>
		`), feed: Feed{Title: "test", Link: "https://example.com/feed.xml"}},
	}

	for _, tt := range table {
		feed, err := parseRSSResp(tt.ct, tt.content)
		assert.Nil(t, err)
		assert.Equal(t, tt.feed, feed)
	}
}
