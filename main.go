package main

import (
	"github.com/pjquirk/gh-splitpr/ui"

	"fmt"
	"os"

	"github.com/cli/go-gh"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	repo, err := gh.CurrentRepository()
	if err != nil {
		fmt.Printf("Must be run from within a clone of a GitHub repository: %v", err)
        os.Exit(1)
	}

    p := tea.NewProgram(ui.NewModel(repo))
    if err := p.Start(); err != nil {
        fmt.Printf("Fatal error encountered: %v", err)
        os.Exit(1)
    }
}
