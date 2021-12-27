package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type Model struct {
	verbose bool
	repository gh.Repository
	pullRequestId int
}

func NewModel(repo gh.Repository) Model {
	return Model{
		verbose: false,
		repository: repo,
		pullRequestId: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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
	s := fmt.Sprintf("Looking for pull requests in %s/%s", m.repository.Owner(), m.repository.Name())
	return s + "\nPress q to quit."
}