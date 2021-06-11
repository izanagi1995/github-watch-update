package utils

import (
	"os"

	"github.com/go-git/go-git/v5"
)

// CloneRepo clones a repository from `url` into `dst`
func CloneRepo(url string, dst string) (*git.Repository, error) {
	repo, err := git.PlainClone(dst, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return repo, err
}

// GetHeadHash returns the hash of of the HEAD reference commit
func GetHeadHash(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}
