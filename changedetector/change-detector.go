package changedetector

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v35/github"
	"github.com/izanagi1995/github-watch-update/utils"
	"go.uber.org/zap"
)

// ChangeDetector holds the configuration and is responsible
// to query GitHub for changes in the branch
type ChangeDetector struct {
	RepoUrl *url.URL
	Branch  string
	Watch   bool
	Log     *zap.SugaredLogger

	localRepo    *git.Repository
	githubClient *github.Client
}

// Init looks up if the repo has already been cloned.
// If yes, it retrieves the local commit hash
func (cd *ChangeDetector) Init() error {
	repoName := utils.RepoName(*cd.RepoUrl)
	if _, err := os.Stat("./" + repoName); err != nil && os.IsNotExist(err) {
		r, errClone := utils.CloneRepo(cd.RepoUrl.String(), "./"+repoName)
		if errClone != nil {
			return errClone
		}
		cd.localRepo = r
	} else if err != nil {
		return err
	}

	// If we didn't clone the repo, it means it exists, so open it
	if cd.localRepo == nil {
		r, err := git.PlainOpen("./" + repoName)
		if err != nil {
			return err
		}
		cd.localRepo = r
	}

	cd.githubClient = github.NewClient(nil)

	return nil
}

func (cd *ChangeDetector) Run() error {
	if cd.Watch {
		for range time.NewTicker(30 * time.Second).C {
			cd.Log.Debug("Run : tick")
			err := cd.CheckChange()
			if err != nil {
				return err
			}
		}
	} else {
		return cd.CheckChange()
	}
	return nil
}

func (cd *ChangeDetector) CheckChange() error {
	// TODO : Handle authentication ?

	ctx, cancel := context.WithCancel(context.Background())
	// Let's cancel the context at the end of the function
	defer cancel()
	branch, _, err := cd.githubClient.Repositories.GetBranch(ctx, utils.RepoOwner(*cd.RepoUrl), utils.RepoName(*cd.RepoUrl), cd.Branch)
	if err != nil {
		return err
	}
	remoteHeadHash := branch.Commit.SHA
	localHeadHash, err := utils.GetHeadHash(cd.localRepo)
	if err != nil {
		return err
	}
	cd.Log.Debugw("Checking hashes", "local", localHeadHash, "remote", *remoteHeadHash)
	if *remoteHeadHash != localHeadHash {
		if worktree, err := cd.localRepo.Worktree(); err != nil {
			return err
		} else {
			cd.Log.Infof("Changes detected : remote HEAD at %s", *remoteHeadHash)
			err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
			if err == nil {
				cd.Log.Info("Successfully pulled remote origin")
			}
			return err
		}
	}
	cd.Log.Info("No changes detected")
	return nil
}
