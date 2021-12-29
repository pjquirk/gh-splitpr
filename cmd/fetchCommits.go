package cmd

import (
	"encoding/json"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type topLevelRaw struct {
	Commits []commitRaw `json:"commits"`
}

type commitRaw struct {
	Comment string `json:"messageHeadline"`
	Sha     string `json:"oid"`
}

type CommitsFetched struct {
	Commits []Commit
}

func FetchCommits(
	repository gh.Repository,
	pullRequestId int,
) tea.Msg {
	nwo := ToNwo(repository)
	stdOut, _, err := gh.Exec("pr", "view", strconv.Itoa(pullRequestId), "--repo", nwo, "--json", "commits")
	if err != nil {
		return ErrMsg{err}
	}

	var topLevel topLevelRaw
	err = json.Unmarshal(stdOut.Bytes(), &topLevel)
	if err != nil {
		return ErrMsg{err}
	}

	return CommitsFetched{convertCommits(topLevel.Commits)}
}

func convertCommits(src []commitRaw) []Commit {
	dst := make([]Commit, len(src))
	for i, t := range src {
		dst[i] = Commit{
			Sha:     t.Sha,
			Comment: t.Comment,
		}
	}
	return dst
}
