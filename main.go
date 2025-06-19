package main

import (
	"errors"
	"flag"
	"fmt"
	"nextver/debugger"
	"nextver/repository"
	"os"
)

var Version = "development"

func main() {
	var debugFlag bool
	var showVersionFlag bool
	var noNewLineFlag bool
	var debug *debugger.Debugger

	flag.BoolVar(&debugFlag, "d", false, "Debug mode")
	flag.BoolVar(&showVersionFlag, "v", false, "Show version")
	flag.BoolVar(&noNewLineFlag, "c", false, "Suppress trailing new line")
	flag.Parse()

	if showVersionFlag {
		showVersion()
		os.Exit(0)
	}

	debug = debugger.New(debugFlag)

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
		debug.Log("No version detected.")
	} else {
		debug.Logf("Detected current version: %s", currentVersion)
	}

	nextVersion, err := repo.EstimatedNextVersion(currentVersion, debug)
	if err != nil {
		dyingMessage(err)
	}

	if currentVersion == nextVersion {
		debug.Logf("%s is the most recent version.\n", currentVersion)
	} else {
		debug.Logf("Estimated next version: %s\n", nextVersion)
		if !debugFlag {
			if noNewLineFlag {
				fmt.Print(nextVersion)
			} else {
				fmt.Println(nextVersion)
			}
		}
	}
}

func dyingMessage(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}

func showVersion() {
	fmt.Println(Version)
}
