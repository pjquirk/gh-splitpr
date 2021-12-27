package ui

import (
	"github.com/pjquirk/gh-splitpr/cmd"

	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type Model struct {
	verbose      bool
	err          error
	repository   gh.Repository
	pullRequests []cmd.PullRequest
}

func NewModel() Model {
	return Model{
		verbose:      false,
		repository:   nil,
		pullRequests: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return cmd.ParseOptions
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case cmd.ErrMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		errMsg := cmd.ErrMsg(msg)
		m.err = errMsg.Error
		return m, tea.Quit

	case cmd.OptionsParsed:
		options := cmd.OptionsParsed(msg)
		m.repository = options.Repository
		//m.pullRequestId = options.PullRequest
		fetchPullRequests := func() tea.Msg {
			return cmd.FetchPullRequests(m.repository)
		}
		return m, fetchPullRequests

	case cmd.PullRequestsFetched:
		fetched := cmd.PullRequestsFetched(msg)
		m.pullRequests = fetched.PullRequests
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "?":
			return m, nil
		}

	}
	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err)
	}

	if m.repository == nil {
		return "Getting repository information..."
	}
	if m.pullRequests == nil {
		return fmt.Sprintf("Looking for pull requests in %s/%s...", m.repository.Owner(), m.repository.Name())
	}

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Found %d pull requests:", len(m.pullRequests)))
	for i := 0; i < len(m.pullRequests); i++ {
		pr := m.pullRequests[i]
		s.WriteString(fmt.Sprintf("\n%d\t%s\t%s", pr.Number, pr.Author, pr.Title))
	}

	s.WriteString("\nPress q to quit.")

	return s.String()
}

// -----------
