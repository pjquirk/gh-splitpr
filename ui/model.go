package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type Model struct {
	verbose      bool
	repository   gh.Repository
	pullRequests []pullRequest
}

func NewModel(repo gh.Repository) Model {
	return Model{
		verbose:      false,
		repository:   repo,
		pullRequests: nil,
	}
}

func (m Model) Init() tea.Cmd {
	// Find the available pull requests
	return listPullRequests
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		//m.err = msg
		return m, tea.Quit

	case pullRequests:
		prs := pullRequests(msg)
		m.pullRequests = prs.pullRequests
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
	if m.pullRequests == nil {
		return fmt.Sprintf("Looking for pull requests in %s/%s...", m.repository.Owner(), m.repository.Name())
	}

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Found %d pull requests:", len(m.pullRequests)))
	for i := 0; i < len(m.pullRequests); i++ {
		pr := m.pullRequests[i]
		s.WriteString(fmt.Sprintf("\n%d\t%s\t%s", pr.number, pr.author, pr.title))
	}

	s.WriteString("\nPress q to quit.")

	return s.String()
}

// -----------

type errMsg struct {
	err error
}

type pullRequest struct {
	number int
	title  string
	author string
}

type pullRequests struct {
	pullRequests []pullRequest
}
type pullRequestRaw struct {
	Number int       `json:"number"`
	Title  string    `json:"title"`
	Author authorRaw `json:"author"`
}

type authorRaw struct {
	Login string `json:"login"`
}

func ConvertPullRequests(src []pullRequestRaw) []pullRequest {
	dst := make([]pullRequest, len(src))
	for i, t := range src {
		dst[i] = pullRequest{
			number: t.Number,
			title:  t.Title,
			author: t.Author.Login,
		}
	}
	return dst
}

func listPullRequests() tea.Msg {
	stdOut, _, err := gh.Exec("pr", "list", "--json", "number,title,author")
	if err != nil {
		return errMsg{err}
	}

	var rawPullRequests []pullRequestRaw
	err = json.Unmarshal(stdOut.Bytes(), &rawPullRequests)
	if err != nil {
		return errMsg{err}
	}

	return pullRequests{ConvertPullRequests(rawPullRequests)}
}
