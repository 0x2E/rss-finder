package find

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParseHTMLRespItem struct {
	ct      string
	content []byte
	feed    []Feed
}

func TestParseHTMLRespMatchLink(t *testing.T) {
	table := []testParseHTMLRespItem{
		{ct: "text/html", content: []byte(`
		<html>
		<head>
			<title>html title</title>
			<link type="application/rss+xml" title="feed title" href="https://example.com/x/rss.xml">
			<link type="application/atom+xml" href="https://example.com/x/atom.xml">
		</head>
		<body>
			<link type="application/feed+json" title="link in body" href="https://example.com/x/feed.json">
		</body>
		</html>
		`), feed: []Feed{
			{Title: "feed title", Link: "https://example.com/x/rss.xml"},
			{Title: "html title", Link: "https://example.com/x/atom.xml"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLResp(tt.ct, tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.feed, feed)
	}
}

func TestParseHTMLRespMatchA(t *testing.T) {
	table := []testParseHTMLRespItem{
		// match <a>
		{ct: "text/html", content: []byte(`
		<html>
		<head><title>html title</title></head>
		<body>
			<p>xxx</p>
			<main>
				<p>xxx</p>
				<a href="https://example.com/index.xml">RSS1</a>
			</main>
			<footer>
				<a href="https://example.com/x/index.xml">rss2</a>
			</footer>
		</body>
		</html>
		`), feed: []Feed{
			{Title: "RSS1", Link: "https://example.com/index.xml"},
			{Title: "rss2", Link: "https://example.com/x/index.xml"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLResp(tt.ct, tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.feed, feed)
	}
}

func TestParseRSSResp(t *testing.T) {
	type testItem struct {
		ct      string
		content []byte
		feed    Feed
	}

	// todo match all types, e.g. https://github.com/mmcdole/gofeed/tree/master/testdata
	table := []testItem{
		{ct: "application/rss+xml", content: []byte(`
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
