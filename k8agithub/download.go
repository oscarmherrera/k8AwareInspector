package k8agithub

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	"go.uber.org/zap"
	"net/http"
	"os"
	"regexp"
)

var (
	zLog  *zap.Logger
	sugar *zap.SugaredLogger
)

func setLogging(zl *zap.Logger) {
	if zLog == nil {
		zLog = zl
		sugar = zl.Sugar()
	}
}

func GetGithubTagZip(ghClient *github.Client, url string, accessToken *string, tempDir string, zl *zap.Logger) (*string, error) {
	setLogging(zl)

	ctx := context.Background()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", *accessToken))
	req.Header.Add("User-Agent", "k8Aware-client")
	//req.Header.Add("Accept", "application/octet-stream")
	//req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/vnd.github+json")

	resp1, err := ghClient.BareDo(ctx, req)
	if err != nil {
		zLog.Error("unable to get repository", zap.Any("response", resp1), zap.Error(err))
		return nil, err
	}

	disp := resp1.Header.Get("Content-disposition")
	re := regexp.MustCompile(`filename=(.+)`)
	matches := re.FindAllStringSubmatch(disp, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		zLog.Debug("WTF", zap.Any("matches", matches))
		zLog.Debug("", zap.Any("response_header", resp1.Header))
		zLog.Info("", zap.Any("request", req))
		return nil, errors.New("unable to match the filename from response")
	}

	filename := tempDir + "/" + matches[0][1]

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		zLog.Error("creating zip file for source", zap.Error(err))
		return nil, err
	}
	b := make([]byte, 4096)
	var i int
	for err == nil {
		i, err = resp1.Body.Read(b)
		f.Write(b[:i])
	}
	sugar.Info("Finished: %s -> %s\n", url, disp)

	return &filename, nil
}
