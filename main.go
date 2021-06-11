package main

import (
	"flag"
	"net/url"
	"os"

	"github.com/izanagi1995/github-watch-update/changedetector"
	"github.com/izanagi1995/github-watch-update/utils"
	"go.uber.org/zap"
)

func CheckErr(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func main() {
	var repoUrl string
	var branch string
	var watch bool

	debug := os.Getenv("DEBUG")

	flag.StringVar(&repoUrl, "repo", "REQUIRED", "The github repository URL (required)")
	flag.StringVar(&branch, "branch", "REQUIRED", "The branch in the repository to check for (required)")
	flag.BoolVar(&watch, "watch", false, "Enable watch mode (disabled by default)")
	flag.Parse()

	var logger *zap.Logger
	var logError error

	if debug == "1" {
		logger, logError = zap.NewDevelopment()
	} else {
		logger, logError = zap.NewProduction()
	}
	if logError != nil {
		// Cannot initialize logger, panic exit
		panic(logError)
	}
	defer logger.Sync() // flushes buffer, if any

	sugaredLogger := logger.Sugar()

	if repoUrl == "REQUIRED" || branch == "REQUIRED" {
		sugaredLogger.Fatal("The repo and branch arguments are required")
	}

	parsedUrl, err := url.Parse(repoUrl)
	CheckErr(err, sugaredLogger)

	if !utils.ValidateRepoUrl(*parsedUrl) {
		sugaredLogger.Fatalf("%s is not a valid github URL (https://github.com/OWNER/REPO)", repoUrl)
	}

	changeDetector := &changedetector.ChangeDetector{
		RepoUrl: parsedUrl,
		Branch:  branch,
		Watch:   watch,
		Log:     sugaredLogger,
	}

	err = changeDetector.Init()
	CheckErr(err, sugaredLogger)
	err = changeDetector.Run()
	CheckErr(err, sugaredLogger)

}
