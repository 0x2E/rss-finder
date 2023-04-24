package find

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParseHTMLContentItem struct {
	content []byte
	want    []Feed
}

func TestParseHTMLContentMatchLink(t *testing.T) {
	table := []testParseHTMLContentItem{
		{content: []byte(`
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
		`), want: []Feed{
			{Title: "feed title", Link: "https://example.com/x/rss.xml"},
			{Title: "html title", Link: "https://example.com/x/atom.xml"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLContent(tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.want, feed)
	}
}

func TestParseHTMLContentMatchA(t *testing.T) {
	// todo refactor request() to use httptest

	initClient()
	table := []testParseHTMLContentItem{
		// match <a>
		{content: []byte(`
		<html>
		<head><title>html title</title></head>
		<body>
			<p>xxx</p>
			<main>
				<p>xxx</p>
				<a href="https://github.com/0x2E/rss-finder/releases.atom">Release notes from rss-finder</a>
			</main>
			<footer>
				<a href="https://github.com/0x2E/rss-finder">wrong rss</a>
			</footer>
		</body>
		</html>
		`), want: []Feed{
			{Title: "Release notes from rss-finder", Link: "https://github.com/0x2E/rss-finder/releases.atom"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLContent(tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.want, feed)
	}
}
