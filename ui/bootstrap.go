package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type BootstrapModel struct {
	Repository    gh.Repository
	PullRequestId int
	PullRequests  []cmd.PullRequest
}

func (m BootstrapModel) IsComplete() bool {
	return m.Repository != nil && m.PullRequestId > 0
}

func (m BootstrapModel) Init() tea.Cmd {
	return nil
}

func (m BootstrapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case cmd.OptionsParsed:
		options := cmd.OptionsParsed(msg)
		m.Repository = options.Repository
		m.PullRequestId = options.PullRequest
		if m.PullRequestId > 0 {
			// Skip getting all PRs
			return m, nil
		} else {
			fetchPullRequests := func() tea.Msg {
				return cmd.FetchPullRequests(m.Repository)
			}
			return m, fetchPullRequests
		}

	case cmd.PullRequestsFetched:
		fetched := cmd.PullRequestsFetched(msg)
		m.PullRequests = fetched.PullRequests
		return m, nil
	}
	return m, nil
}

func (m BootstrapModel) View() string {
	if m.IsComplete() {
		return ""
	}

	if m.Repository == nil {
		return "Getting repository information..."
	}
	nwo := cmd.ToNwo(m.Repository)

	s := strings.Builder{}

	if m.PullRequests == nil {
		return fmt.Sprintf("Looking for pull requests in %s...", nwo)
	} else {
		s.WriteString("Select a pull request to split:")
		for i := 0; i < len(m.PullRequests); i++ {
			pr := m.PullRequests[i]
			s.WriteString(fmt.Sprintf("\n%d\t%s\t%s", pr.Number, pr.Author, pr.Title))
		}
	}

	return s.String()
}
