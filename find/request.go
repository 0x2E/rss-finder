package find

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

var globalClient *http.Client

func initClient() {
	// avoid ssrf
	// https://www.agwa.name/blog/post/preventing_server_side_request_forgery_in_golang
	socketControl := func(network, address string, c syscall.RawConn) error {
		if !(network == "tcp4" || network == "tcp6") {
			return fmt.Errorf("banned network type: %s", network)
		}

		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("failed to split host:port: %s", err)
		}

		ipaddress := net.ParseIP(host)
		if ipaddress == nil {
			return fmt.Errorf("invalid ip: %s", host)
		}
		if ipaddress.IsLoopback() || ipaddress.IsPrivate() || ipaddress.IsUnspecified() {
			return fmt.Errorf("banned ip range: %s", ipaddress)
		}

		return nil
	}

	safeDialer := &net.Dialer{
		Timeout: 3 * time.Second,
		Control: socketControl,
	}
	safeTransport := &http.Transport{
		DialContext:       safeDialer.DialContext,
		ForceAttemptHTTP2: true,
		// DisableKeepAlives: true,
	}
	globalClient = &http.Client{
		Transport: safeTransport,
		Timeout:   3 * time.Second,
	}
}

func request(link string) (*http.Response, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	ua := os.Getenv("USER_AGENT")
	if strings.TrimSpace(ua) == "" {
		ua = "rss-finder/1.0"
	}
	req.Header.Add("User-Agent", ua)

	return globalClient.Do(req)
}
