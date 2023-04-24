package find

import (
	"log"
	"net/url"
	"sync"
)

type Feed struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func Find(target *url.URL) ([]Feed, error) {
	log.SetPrefix("[" + target.String() + "]")

	initClient()

	// find in third-party service
	fromService, err := tryService(target)
	if err != nil {
		log.Printf("%s: %s\n", "parse service", err)
	}
	if len(fromService) != 0 {
		return fromService, nil
	}

	feedsChan := make(chan []Feed, 2)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// find in HTML
		fromPage, err := tryPageSource(target.String())
		if err != nil {
			log.Printf("%s: %s\n", "parse page", err)
		}

		feedsChan <- fromPage
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// find well-known under this url
		fromWellKnown, err := tryWellKnown(target.Scheme + "://" + target.Host + target.Path) // https://go.dev/play/p/dVt-47_XWjU
		if err != nil {
			log.Printf("%s: %s\n", "parse wellknown", err)
		}
		if len(fromWellKnown) == 0 {
			// find well-known under url root
			fromWellKnown, err = tryWellKnown(target.Scheme + "://" + target.Host)
			if err != nil {
				log.Printf("%s: %s\n", "parse wellknown under root", err)
			}
		}

		feedsChan <- fromWellKnown
	}()

	go func() {
		wg.Wait()
		close(feedsChan)
	}()

	res := make([]Feed, 0)
	for feeds := range feedsChan {
		if len(feeds) == 0 {
			continue
		}

		res = append(res, feeds...)
	}

	return res, err
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
