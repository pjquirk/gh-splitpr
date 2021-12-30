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
	commits   CommitsModel
}

func NewModel() Model {
	return Model{
		verbose:   false,
		bootstrap: NewBootstrapModel(),
		commits:   NewCommitsModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return cmd.ParseOptions
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		command  tea.Cmd
		commands []tea.Cmd
	)

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
		}
	}

	if !m.bootstrap.IsComplete() {
		m.bootstrap, command = m.bootstrap.Update(msg)
		commands = append(commands, command)
	} else {
		m.commits, command = m.commits.Update(msg)
		commands = append(commands, command)
	}
	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	if len(m.err) > 0 {
		return m.err
	}

	s := strings.Builder{}

	if !m.bootstrap.IsComplete() {
		s.WriteString(m.bootstrap.View())
	} else {
		s.WriteString(m.commits.View())
	}

	view := s.String()
	return view
}
