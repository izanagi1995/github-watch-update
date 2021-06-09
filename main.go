package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var repoUrl string
	var branch string
	var watch bool

	flag.StringVar(&repoUrl, "repo", "REQUIRED", "The github repository URL (required)")
	flag.StringVar(&branch, "branch", "REQUIRED", "The branch in the repository to check for (required)")
	flag.BoolVar(&watch, "watch", false, "Enable watch mode (disabled by default)")
	flag.Parse()

	if repoUrl == "REQUIRED" || branch == "REQUIRED" {
		fmt.Fprintln(os.Stderr, "The repo and branch arguments are required")
		os.Exit(1)
	}

}
