package utils

import (
	"net/url"
	"strings"
)

func ValidateRepoUrl(url url.URL) bool {
	return url.Host == "github.com" && len(SplitUrlPath(url)) == 2
}

// SplitUrlPath splits an url path into its parts
func SplitUrlPath(url url.URL) []string {
	path := url.Path
	// to avoid empty first and last element, let's trim the url
	// remove leading slash if present
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	// remove trailing slash if present
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}

// RepoName gives the name of the repository from a parsed Github URL
func RepoName(url url.URL) string {
	return SplitUrlPath(url)[1]
}

// RepoOwner gives the owner of the repository from a parsed Github URL
func RepoOwner(url url.URL) string {
	return SplitUrlPath(url)[0]
}
