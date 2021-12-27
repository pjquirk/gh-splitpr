package main

import (
	"github.com/pjquirk/gh-splitpr/ui"

	"log"

	"github.com/cli/go-gh"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	repo, err := gh.CurrentRepository()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: args for the PR

	p := tea.NewProgram(ui.NewModel(repo))
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
