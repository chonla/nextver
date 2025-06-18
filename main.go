package main

import (
	"errors"
	"fmt"
	"nextver/repository"
	"os"
)

func main() {
	repoPath := "."

	if !repository.IsGitRepo(repoPath) {
		dyingMessage(errors.New("target path is not git repo path"))
	}

	repo, err := repository.Open(repoPath)
	if err != nil {
		dyingMessage(err)
	}

	currentVersion, err := repo.RecentTaggedVersion()
	if err != nil {
		dyingMessage(err)
	}
	if currentVersion == "" {
		fmt.Println("No version detected.")
	} else {
		fmt.Printf("Detected current version: %s\n", currentVersion)
	}

	nextVersion, err := repo.EstimatedNextVersion(currentVersion)
	if err != nil {
		dyingMessage(err)
	}
	if currentVersion == nextVersion {
		fmt.Printf("%s is the most recent version.", currentVersion)
	} else {
		fmt.Printf("Estimated next version: %s\n", nextVersion)
	}
}

func dyingMessage(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}
