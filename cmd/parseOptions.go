package cmd

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type OptionsParsed struct {
	Repository  gh.Repository
	PullRequest int
}

func ParseOptions() tea.Msg {
	// Attempt to get the current directory's repository as a default
	defaultRepository, defaultRepositoryError := gh.CurrentRepository()

	repoFlag := flag.String(
		"repo",
		ToNwo(defaultRepository),
		"Select another repository using the [HOST/]OWNER/REPO format",
	)
	pullRequestNumber := flag.Int(
		"pullRequest",
		0,
		"Specifies the number of the pull request to split",
	)
	flag.Parse()

	if *repoFlag == "" {
		// The user didn't specify a repository, and gh couldn't determine it
		// from the current working directory.
		log.Fatal(defaultRepositoryError)
	}
	repository, err := ToRepository(*repoFlag)
	if err != nil {
		return ErrMsg{err}
	}

	return OptionsParsed{
		Repository:  repository,
		PullRequest: *pullRequestNumber,
	}
}
