package main

import (
	"github.com/pjquirk/gh-splitpr/ui"

	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
