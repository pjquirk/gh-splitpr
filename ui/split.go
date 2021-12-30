package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type SplitModel struct {
	Repository    gh.Repository
	PullRequestId int
	Commits       []cmd.Commit
}

func NewSplitModel() SplitModel {
	return SplitModel{
		Repository:    nil,
		PullRequestId: -1,
		Commits:       nil,
	}
}

func (m SplitModel) IsComplete() bool {
	return false
}

func (m SplitModel) Init() tea.Cmd {
	return nil
}

func (m SplitModel) Update(msg tea.Msg) (SplitModel, tea.Cmd) {
	var (
		//command  tea.Cmd
		commands []tea.Cmd
	)

	// switch msg := msg.(type) {
	// }

	//m.prSelector, command = m.prSelector.Update(msg)
	//commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m SplitModel) View() string {
	if m.IsComplete() {
		return ""
	}

	//return m.prSelector.View()
	return "Here wo go!"
}
