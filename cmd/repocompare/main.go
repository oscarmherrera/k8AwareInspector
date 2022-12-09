package main

import (
	"context"
	"errors"
	"flag"
	"github.com/google/go-github/v48/github"
	"go.uber.org/zap"
	"k8Aware/k8aast"
	"k8Aware/k8agithub"
	"os"
	"strings"
)

// with go modules enabled (GO111MODULE=on or outside GOPATH)

var (
	zLog  *zap.Logger
	sugar *zap.SugaredLogger
)

func init() {
	//logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	zLog = logger
	sugar = logger.Sugar()
}

func main() {
	repoURL := flag.String("repo", "", "repo to inspect")
	tag1 := flag.String("newTag", "", "latest tag")
	tag2 := flag.String("oldTag", "", "previous tag to compare against")
	accessToken := flag.String("token", "", "your github access token")

	flag.Parse()

	if strings.Contains(*repoURL, "https://github.com/") == true {
		newRepoURL := strings.ReplaceAll(*repoURL, "https://github.com/", "")
		repoURL = &newRepoURL
	}
	if strings.Contains(*repoURL, "http://github.com/") == true {
		newRepoURL := strings.ReplaceAll(*repoURL, "https://github.com/", "")
		repoURL = &newRepoURL
	}
	if strings.Contains(*repoURL, ".git") == true {
		newRepoURL := strings.ReplaceAll(*repoURL, ".git", "")
		repoURL = &newRepoURL
	}

	ghClient := k8agithub.GetGithubClient(accessToken)

	repoNameComponents := strings.Split(*repoURL, "/")
	username := repoNameComponents[0]
	repoName := repoNameComponents[1]

	urls, err := getTagData(ghClient, username, repoName, tag1, tag2)
	if err != nil {
		zLog.Fatal("unable to get repository", zap.Error(err))
	}

	tempDir, err := os.MkdirTemp("./", "github")
	if err != nil {
		zLog.Fatal("unable to create temp directory", zap.Error(err))
	}
	defer os.RemoveAll("./" + tempDir)

	err = k8aast.ProcessTags(ghClient, urls, accessToken, tempDir, false, zLog)
	if err != nil {
		zLog.Error("unable process tags", zap.Error(err))
		//os.Exit(1)
	}

}

func getTagData(ghClient *github.Client, username, repoName string, tag1, tag2 *string) ([]*string, error) {
	ctx := context.Background()
	repoTagList, resp0, err := ghClient.Repositories.ListTags(ctx, username, repoName, nil)
	if err != nil {
		zLog.Error("unable to get repository", zap.Any("response", resp0), zap.Error(err))
		return nil, err
	}

	tag1Found := false
	tag1URL := ""
	tag2Found := false
	tag2URL := ""

	for _, repo := range repoTagList {
		if *repo.Name == *tag1 {
			tag1Found = true
			tag1URL = *repo.ZipballURL
		}
		if *repo.Name == *tag2 {
			tag2Found = true
			tag2URL = *repo.ZipballURL
		}
		if tag1Found == true && tag2Found == true {
			break
		}
	}

	if tag1Found == false {
		return nil, errors.New("unable to find newTag")
	}
	if tag2Found == false {
		return nil, errors.New("unable to find oldTag")
	}

	urls := make([]*string, 2)
	urls[0] = &tag1URL
	urls[1] = &tag2URL
	zLog.Debug("", zap.String("tag_url_1", tag1URL))
	zLog.Debug("", zap.String("tag_url_2", tag2URL))
	return urls, nil
}
