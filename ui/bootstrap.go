package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type BootstrapModel struct {
	Repository    gh.Repository
	PullRequestId int

	pullRequests []cmd.PullRequest
	prSelector   *selector.Model
}

func (m BootstrapModel) IsComplete() bool {
	return m.Repository != nil && m.PullRequestId > 0
}

func (m BootstrapModel) Init() tea.Cmd {
	return nil
}

func (m BootstrapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		command  tea.Cmd
		commands []tea.Cmd
	)

	// Handle selected PR first since its different
	if msg == common.DONE {
		selected := m.prSelector.Selected()
		pr := selected.(cmd.PullRequest)
		m.PullRequestId = pr.Number
		return m, nil
	} else {
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
			m.pullRequests = fetched.PullRequests

			data := make([]interface{}, len(m.pullRequests))
			for i, v := range m.pullRequests {
				data[i] = v
			}
			m.prSelector = &selector.Model{
				Data:       data,
				PerPage:    5,
				HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select a pull request to split:"),
				// [1] The title of the pull request (@author)
				SelectedFunc: func(m selector.Model, obj interface{}, index int) string {
					pr := obj.(cmd.PullRequest)
					return common.FontColor(fmt.Sprintf("[%d] %s (@%s)", pr.Number, pr.Title, pr.Author), selector.ColorSelected)
				},
				//  1. The title of the pull request (@author)
				UnSelectedFunc: func(m selector.Model, obj interface{}, index int) string {
					pr := obj.(cmd.PullRequest)
					return common.FontColor(fmt.Sprintf(" %d. %s (%s)", pr.Number, pr.Title, pr.Author), selector.ColorUnSelected)
				},
			}
		}
	}

	if m.prSelector != nil {
		m.prSelector, command = m.prSelector.Update(msg)
		commands = append(commands, command)
	}
	return m, tea.Batch(commands...)
}

func (m BootstrapModel) View() string {
	if m.IsComplete() {
		return ""
	}

	if m.Repository == nil {
		return "Getting repository information..."
	}
	nwo := cmd.ToNwo(m.Repository)

	if m.pullRequests == nil {
		return fmt.Sprintf("Looking for pull requests in %s...", nwo)
	} else {
		return m.prSelector.View()
	}
}
