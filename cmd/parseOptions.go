package cmd

import (
	"flag"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type OptionsParsed struct {
	Repository  gh.Repository
	PullRequest int
}

type UsageShown struct {
	Usage string
}

func ParseOptions() tea.Msg {
	// Attempt to get the current directory's repository as a default
	defaultRepository, defaultRepositoryError := gh.CurrentRepository()

	// Write usage information to a buffer rather than stderr
	usageOutput := strings.Builder{}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(&usageOutput)
	flag.CommandLine.Usage = func() {
		usageOutput.WriteString("Splits the commits in one pull request into one or more smaller PRs.\n")
		usageOutput.WriteString("Usage: \n\n")
		flag.PrintDefaults()
	}

	// Set CLI options
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

	usageText := usageOutput.String()
	if len(usageText) > 0 {
		return UsageShown{usageText}
	}

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
