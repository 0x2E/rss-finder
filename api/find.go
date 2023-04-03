package api

import (
	"encoding/json"
	"errors"
	"log"
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
// resp:
// 1. ok: 200 + data
// 2. user input error: 400 + err msg
// 3. server error: 500
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
		// json.NewEncoder(w).Encode(errorResp{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(successResp{
		Data: res,
	})
}

// var bannedDomainSuffix = []string{
// 	".gov.cn",
// 	".edu.cn",
// }

func validInput(input string) (*url.URL, error) {
	// add scheme, otherwise url.Parse cannot resolve the Hostname
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	target, err := url.Parse(input)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to parse the url")
	}

	// domain := target.Hostname()
	// for _, ban := range bannedDomainSuffix {
	// 	if strings.HasSuffix(domain, ban) {
	// 		return nil, errors.New("banned domain extentions")
	// 	}
	// }

	// todo safe check

	return target, nil
}
