package cmd

import "github.com/cli/go-gh"

type ErrMsg struct {
	Error error
}

type PullRequestSelected struct {
	Repository    gh.Repository
	PullRequestId int
}
