package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/0x2e/rss-finder/find"
)

type errorResp struct {
	Error string `json:"error"`
}

type successResp struct {
	Data interface{} `json:"data"`
}

// FindHandler is a entrypoint for Vercel serverless function,
// request format: [GET] ?url=xxx
func FindHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	target, err := validInput(r.FormValue("url"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp{Error: err.Error()})
		return
	}

	res, err := find.Find(target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(successResp{
		Data: res,
	})
}

var bannedDomainSuffix = []string{
	".gov.cn",
	".edu.cn",
}

func validInput(input string) (*url.URL, error) {
	// add scheme, otherwise url.Parse cannot resolve the Hostname
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	target, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	domain := target.Hostname()

	for _, ban := range bannedDomainSuffix {
		if strings.HasSuffix(domain, ban) {
			return nil, errors.New("banned domains")
		}
	}

	// todo safe check

	return target, nil
}
