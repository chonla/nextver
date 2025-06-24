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

	"nextver/debugger"
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

func (r *Repo) RecentTaggedVersion(versionIsNotPrefixedByV bool) (string, error) {
	tags, err := r.repo.Tags()
	if err != nil {
		return "", err
	}
	tagRefs := []string{}
	tags.ForEach(func(tag *plumbing.Reference) error {
		tagRefs = append(tagRefs, tag.Name().Short())
		return nil
	})
	tags.Close()
	verRefs := lo.Filter(tagRefs, func(ver string, _ int) bool { return versioning.IsValid(ver, versionIsNotPrefixedByV) })
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

func (r *Repo) EstimatedNextVersion(currentVer string, versionIsNotPrefixedByV bool, debug *debugger.Debugger) (string, error) {
	prefix := "v"
	if versionIsNotPrefixedByV {
		debug.Log("Version is not prefixed by v.")
		prefix = ""
	} else {
		debug.Log("Version is prefixed by v.")
	}

	if currentVer == "" {
		debug.Log("Current version is missing. Start a new one at v1.0.0.")
		return fmt.Sprintf("%s1.0.0", prefix), nil
	}

	// Get HEAD
	headRef, err := r.repo.Head()
	if err != nil {
		debug.Logf("Unable to retrieve HEAD. %v", err)
		return "", err
	}
	headCommitID := headRef.Hash().String()
	debug.Logf("HEAD commit ID=%s", headCommitID)

	// Get Tag
	tagRef, err := r.repo.Tag(currentVer)
	if err != nil {
		debug.Logf("Unable to retrieve tags. %v", err)
		return "", err
	}
	tagCommitID := tagRef.Hash().String()
	debug.Logf("Latest tag commit ID=%s", tagCommitID)

	// Get commit from HEAD to tag, excluding tag itself
	newCommits := []string{}
	commits, err := r.repo.Log(&git.LogOptions{From: headRef.Hash()})
	if err != nil {
		debug.Logf("Unable to retrieve commits. %v", err)
		return "", err
	}

	commits.ForEach(func(commit *object.Commit) error {
		commitID := commit.Hash.String()
		if commitID == tagCommitID {
			return errors.New("end of new commits")
		}
		newCommits = append(newCommits, commit.Message)
		return nil
	})
	commits.Close()
	debug.Logf("%d commit(s) since latest tag", len(newCommits))

	majorCount := 0
	minorCount := 0
	revisionCount := 0
	lo.ForEach(newCommits, func(commitMessage string, _ int) {
		result, err := convcommit.Parse(commitMessage)
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
	})
	debug.Log("============")
	debug.Log("Commit stats")
	debug.Log("------------")
	debug.Logf("Major change(s) = %d", majorCount)
	debug.Logf("Minor change(s) = %d", minorCount)
	debug.Logf("Revision change(s) = %d", revisionCount)
	debug.Log("============")

	currentSemVer, _ := versioning.Parse(currentVer, versionIsNotPrefixedByV)
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
