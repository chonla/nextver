package repository

import (
	"errors"
	"fmt"
	"os"

	convcommit "github.com/chonla/conv-commit"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/samber/lo"
	"golang.org/x/mod/semver"

	"nextver/versioning"
)

type Repo struct {
	repo *git.Repository
}

func IsGitRepo(path string) bool {
	gitPath := fmt.Sprintf("%s/.git", path)
	fInfo, err := os.Stat(gitPath)
	if os.IsNotExist(err) {
		return false
	}
	return fInfo.IsDir()
}

func Open(path string) (*Repo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repo{
		repo,
	}, nil
}

func (r *Repo) RecentTaggedVersion() (string, error) {
	tags, err := r.repo.Tags()
	if err != nil {
		return "", err
	}
	tagRefs := []string{}
	tags.ForEach(func(tag *plumbing.Reference) error {
		tagRefs = append(tagRefs, tag.Name().Short())
		return nil
	})
	verRefs := lo.Filter(tagRefs, func(ver string, _ int) bool { return versioning.IsValid(ver) })
	if len(verRefs) == 0 {
		return "", nil
	}

	semver.Sort(verRefs)
	recentVer, found := lo.Last(verRefs)
	if found {
		return recentVer, nil
	}
	return "", errors.New("no recent tagged version")
}

func (r *Repo) EstimatedNextVersion(currentVer string) (string, error) {
	if currentVer == "" {
		return "v1.0.0", nil
	}
	tagRef, err := r.repo.Tag(currentVer)
	if err != nil {
		return "", err
	}
	commits, err := r.repo.Log(&git.LogOptions{From: tagRef.Hash()})
	if err != nil {
		return "", err
	}

	majorCount := 0
	minorCount := 0
	revisionCount := 0
	commits.ForEach(func(commit *object.Commit) error {
		result, err := convcommit.Parse(commit.Message)
		if err == nil {
			if result.IsBreakingChange {
				majorCount += 1
			} else {
				switch result.Type {
				case "feat":
					minorCount += 1
				case "fix":
					revisionCount += 1
				}
			}
		}
		return nil
	})

	currentSemVer, _ := versioning.Parse(currentVer)
	if majorCount > 0 {
		return currentSemVer.NextMajor().String(), nil
	} else {
		if minorCount > 0 {
			return currentSemVer.NextMinor().String(), nil
		} else {
			if revisionCount > 0 {
				return currentSemVer.NextPatch().String(), nil
			}
		}
	}
	return currentSemVer.String(), nil
}
