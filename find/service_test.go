package find

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHub(t *testing.T) {
	type testItem struct {
		url  *url.URL
		want []Feed
	}

	urlBase, _ := url.Parse("https://github.com")
	urlUser, _ := url.Parse("https://github.com/user")
	urlUserFeed := append(githubGlobalFeed, genGitHubUserFeed("user")...)
	urlUserRepo, _ := url.Parse("https://github.com/user/repo")
	urlUserRepoFeed := append(urlUserFeed, genGitHubUserRepoFeed("user/repo")...)

	table := []testItem{
		{url: urlBase, want: githubGlobalFeed},
		{url: urlUser, want: urlUserFeed},
		{url: urlUserRepo, want: urlUserRepoFeed},
	}

	for _, tt := range table {
		feed, err := githubMatcher(tt.url)
		assert.Nil(t, err)
		assert.ElementsMatch(t, tt.want, feed)
	}
}
