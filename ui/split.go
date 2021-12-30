package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type SplitModel struct {
	Repository    gh.Repository
	PullRequestId int
	Commits       []cmd.Commit

	splitSettings *cmd.SplitSettings
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

	switch msg := msg.(type) {
	case cmd.CommitsSelected:
		commitsSelected := cmd.CommitsSelected(msg)
		m.Repository = commitsSelected.Repository
		m.PullRequestId = commitsSelected.PullRequestId
		m.Commits = commitsSelected.Commits
		commands = append(commands, func() tea.Msg {
			return cmd.GetSplitSettings(m.Repository, m.PullRequestId, m.Commits)
		})

	case cmd.SplitSettings:
		splitSettings := cmd.SplitSettings(msg)
		m.splitSettings = &splitSettings
	}

	//m.prSelector, command = m.prSelector.Update(msg)
	//commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m SplitModel) View() string {
	if m.IsComplete() {
		return ""
	}

	/* To perform the split, we need a working copy.  We have a few options:
	 * 1. Require the working directory to match the selected repository.
	 *    - We could remove the repository CLI option.
	 *    - We need to fetch latest
	 *    - May hit issues if the working dir has changes in it
	 * 2. Clone the repository into a temp directory
	 *    - gh repo clone <nwo> <tmp dir>
	 *    - Delete the directory afterwards
	 *
	 * Once we have a working copy with any branch checked out:
	 * - git checkout -b <branch name> <same base branch as original branch>
	 * - git cherry-pick sha1 sha2 sha3 sha4
	 * - Write PR body to a temp file
	 * - gh pr create --base <original base> --head <branch name> --title <PR title> --body-file <temp file with body>
	 * - Delete body file
	 *
	 * Potential problems:
	 * - The chosen commits may not apply cleanly
	 * - Need to check the remote name
	 */

	// Necessary inputs:
	// - New branch name (default to current branch + "split-ABCDEFGH"
	// - Base branch (default to base of given PR)
	// - New PR title (default to "Split from " + PR's title)
	if m.splitSettings != nil {
		b := strings.Builder{}
		b.WriteString("Settings:\n")
		b.WriteString(fmt.Sprintf("  BodyFile:       %s\n", m.splitSettings.BodyFile))
		b.WriteString(fmt.Sprintf("  CloneDir:       %s\n", m.splitSettings.CloneDir))
		b.WriteString(fmt.Sprintf("  BaseBranchName: %s\n", m.splitSettings.BaseBranchName))
		b.WriteString(fmt.Sprintf("  NewBranchName:  %s\n", m.splitSettings.NewBranchName))
		b.WriteString(fmt.Sprintf("  NewPrTitle:     %s\n", m.splitSettings.NewPrTitle))
		return b.String()
	}

	//return m.prSelector.View()
	return "Here wo go!"
}
