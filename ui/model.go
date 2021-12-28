package ui

import (
	"strings"

	"github.com/pjquirk/gh-splitpr/cmd"

	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	verbose   bool
	err       string
	bootstrap BootstrapModel
}

func NewModel() Model {
	return Model{
		verbose:   false,
		bootstrap: BootstrapModel{},
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
		m.err = fmt.Sprintf("Error: %s", errMsg.Error)
		return m, tea.Quit

	case cmd.UsageShown:
		usageShown := cmd.UsageShown(msg)
		m.err = usageShown.Usage
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			if !m.bootstrap.IsComplete() {
				return m.bootstrap.Update(msg)
			}
		}
	default:
		if !m.bootstrap.IsComplete() {
			return m.bootstrap.Update(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.err) > 0 {
		return m.err
	}

	s := strings.Builder{}

	if !m.bootstrap.IsComplete() {
		s.WriteString(m.bootstrap.View())
	} else {
		nwo := cmd.ToNwo(m.bootstrap.Repository)
		s.WriteString(fmt.Sprintf("Getting commit information for #%d in %s...", m.bootstrap.PullRequestId, nwo))
	}

	s.WriteString("\n\nPress q to quit")
	return s.String()
}
