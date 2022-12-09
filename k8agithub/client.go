package k8agithub

import (
	"context"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(accessToken *string) *github.Client {
	var ghClient *github.Client
	if *accessToken == "" {
		ghClient = github.NewClient(nil)
	} else {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		ghClient = github.NewClient(tc)
	}

	return ghClient
}
