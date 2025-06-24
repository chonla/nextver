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
	var showLatestVersionFlag bool
	var noNewLineFlag bool
	var versionHasNoPrefixed bool
	var debug *debugger.Debugger

	flag.BoolVar(&debugFlag, "d", false, "Debug mode, print considering steps.")
	flag.BoolVar(&showVersionFlag, "v", false, "Show version of nextver.")
	flag.BoolVar(&showLatestVersionFlag, "t", false, "Show detected latest version.")
	flag.BoolVar(&noNewLineFlag, "e", false, "Suppress trailing new line. Print only version out.")
	flag.BoolVar(&versionHasNoPrefixed, "n", false, "Version is not prefixed by v, for example, 1.0.0.")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of nextver:\n\n")
		fmt.Fprintf(w, "  nextver [options...] [dir]\n\n")
		fmt.Fprintf(w, "Options:\n")
		flag.PrintDefaults() // Prints the default flag descriptions
		fmt.Fprintf(w, "\nFor more information, visit https://github.com/chonla/nextver.\n")
	}

	flag.Parse()

	if showVersionFlag {
		showVersion()
		os.Exit(0)
	}

	debug = debugger.New(debugFlag)

	cmdArgs := flag.Args()

	repoPath := "."
	if len(cmdArgs) > 0 {
		repoPath = cmdArgs[0]
	}

	if !repository.IsGitRepo(repoPath) {
		dyingMessage(errors.New("target path is not git repo path"))
	}

	repo, err := repository.Open(repoPath)
	if err != nil {
		dyingMessage(err)
	}

	currentVersion, err := repo.RecentTaggedVersion(versionHasNoPrefixed)
	if err != nil {
		dyingMessage(err)
	}

	if showLatestVersionFlag {
		fmt.Println(currentVersion)
		os.Exit(0)
	}

	if currentVersion == "" {
		debug.Log("No version detected.")
	} else {
		debug.Logf("Detected current version: %s", currentVersion)
	}

	nextVersion, err := repo.EstimatedNextVersion(currentVersion, versionHasNoPrefixed, debug)
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
