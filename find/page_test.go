package find

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParseHTMLRespItem struct {
	content []byte
	want    []Feed
}

func TestParseHTMLRespMatchLink(t *testing.T) {
	table := []testParseHTMLRespItem{
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
		feed, err := parseHTMLResp(tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.want, feed)
	}
}

func TestParseHTMLRespMatchA(t *testing.T) {
	table := []testParseHTMLRespItem{
		// match <a>
		{content: []byte(`
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
		`), want: []Feed{
			{Title: "RSS1", Link: "https://example.com/index.xml"},
			{Title: "rss2", Link: "https://example.com/x/index.xml"},
		}},
	}

	for _, tt := range table {
		feed, err := parseHTMLResp(tt.content)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.want, feed)
	}
}
