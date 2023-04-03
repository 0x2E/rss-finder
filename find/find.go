package find

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Feed struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func Find(target *url.URL) ([]Feed, error) {
	log.SetPrefix("[" + target.String() + "]")

	// find in third-party service
	fromService, err := tryService(target)
	if err != nil {
		log.Printf("%s: %s\n", "parse service", err)
	}
	if len(fromService) != 0 {
		return fromService, nil
	}

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
