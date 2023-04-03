package find

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type serviceMatcher func(*url.URL) ([]Feed, error)

func tryService(link *url.URL) ([]Feed, error) {
	matcher := []serviceMatcher{
		githubMatcher,
	}
	for _, fn := range matcher {
		feed, err := fn(link)
		if err != nil {
			continue
		}
		if len(feed) != 0 {
			return feed, nil
		}
	}
	return nil, nil
}

var githubGlobalFeed = []Feed{
	{Title: "global public timeline", Link: "https://github.com/timeline"},
	{Title: "global security advisories", Link: "https://github.com/security-advisories.atom"},
}

// https://docs.github.com/en/rest/activity/feeds?apiVersion=2022-11-28#get-feeds
func githubMatcher(link *url.URL) ([]Feed, error) {
	if !strings.HasSuffix(link.Hostname(), "github.com") {
		return nil, nil
	}
	feed := githubGlobalFeed

	splited := strings.SplitN(link.Path, "/", 4) // split "/user/repo/" -> []string{"", "user", "repo", ""}
	len := len(splited)
	if len < 2 {
		return feed, nil
	}

	username, reponame := "", ""

	if len >= 2 { // username exist
		username = splited[1]
		re, err := regexp.Compile(`^[a-zA-Z0-9][-]?[a-zA-Z0-9]{0,38}$`) // todo need improve
		if err != nil {
			return feed, err
		}
		if !re.MatchString(username) {
			return feed, nil
		}
		feed = append(feed, genGitHubUserFeed(username)...)
	}

	if len >= 3 { // reponame exist
		reponame = splited[2]
		re, err := regexp.Compile(`^[a-zA-Z0-9][a-zA-Z0-9-_\.]{0,98}[a-zA-Z0-9]$`) // todo need improve
		if err != nil {
			return feed, err
		}
		if !re.MatchString(reponame) {
			return feed, nil
		}
		feed = append(feed, genGitHubUserRepoFeed(username+"/"+reponame)...)
	}

	return feed, nil
}

func genGitHubUserFeed(username string) []Feed {
	return []Feed{{Title: username, Link: fmt.Sprintf("https://github.com/%s.atom", username)}}
}

func genGitHubUserRepoFeed(userRepo string) []Feed {
	return []Feed{
		{Title: fmt.Sprintf("%s repo commits", userRepo), Link: fmt.Sprintf("https://github.com/%s/commits.atom", userRepo)},
		{Title: fmt.Sprintf("%s repo releases", userRepo), Link: fmt.Sprintf("https://github.com/%s/releases.atom", userRepo)},
		{Title: fmt.Sprintf("%s repo tags", userRepo), Link: fmt.Sprintf("https://github.com/%s/tags.atom", userRepo)},
		{Title: fmt.Sprintf("%s repo wiki", userRepo), Link: fmt.Sprintf("https://github.com/%s/wiki.atom", userRepo)},
	}
}
