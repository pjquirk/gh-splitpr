package cmd

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type pullRequestRaw struct {
	Number int       `json:"number"`
	Title  string    `json:"title"`
	Author authorRaw `json:"author"`
}

type authorRaw struct {
	Login string `json:"login"`
}

type PullRequestsFetched struct {
	PullRequests []PullRequest
}

func FetchPullRequests() tea.Msg {
	stdOut, _, err := gh.Exec("pr", "list", "--json", "number,title,author")
	if err != nil {
		return ErrMsg{err}
	}

	var rawPullRequests []pullRequestRaw
	err = json.Unmarshal(stdOut.Bytes(), &rawPullRequests)
	if err != nil {
		return ErrMsg{err}
	}

	return PullRequestsFetched{convertPullRequests(rawPullRequests)}
}

func convertPullRequests(src []pullRequestRaw) []PullRequest {
	dst := make([]PullRequest, len(src))
	for i, t := range src {
		dst[i] = PullRequest{
			Number: t.Number,
			Title:  t.Title,
			Author: t.Author.Login,
		}
	}
	return dst
}
