package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type commitItem struct {
	sha     string
	comment string
	checked bool
}

func (i commitItem) FilterValue() string { return i.comment }

func newCommitItem(commit cmd.Commit) commitItem {
	return commitItem{
		sha:     commit.Sha,
		comment: commit.Comment,
		checked: false,
	}
}

func newCommitItems(commits []cmd.Commit) []list.Item {
	items := make([]list.Item, len(commits))
	for i, v := range commits {
		items[i] = newCommitItem(v)
	}
	return items
}

var (
	commitItemStyle         = defaultItemStyle
	commitCheckedItemStyle  = defaultItemStyle.Background(lipgloss.Color("240"))
	commitSelectedItemStyle = defaultSelectedItemStyle
)

type commitItemDelegate struct{}

func (d commitItemDelegate) Height() int                               { return 1 }
func (d commitItemDelegate) Spacing() int                              { return 0 }
func (d commitItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d commitItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(commitItem)
	if !ok {
		return
	}

	// Trim all but the first 8 chars of the SHA
	str := fmt.Sprintf("%.8s - %s", i.sha, i.comment)

	style := commitItemStyle
	if index == m.Index() {
		style = commitSelectedItemStyle
	}
	fn := func(s string) string {
		if i.checked {
			return style.Render("[X] " + s)
		} else {
			return style.Render("[ ] " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type SplitModel struct {
	Repository    gh.Repository
	PullRequestId int

	commits        []cmd.Commit
	commitSelector list.Model
}

func NewSplitModel() SplitModel {
	return SplitModel{
		Repository:     nil,
		PullRequestId:  -1,
		commits:        []cmd.Commit{},
		commitSelector: newListModel(5, commitItemDelegate{}, commitItemStyle),
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
		command  tea.Cmd
		commands []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.commitSelector.SetWidth(msg.Width)

	case cmd.PullRequestSelected:
		prSelected := cmd.PullRequestSelected(msg)
		m.Repository = prSelected.Repository
		m.PullRequestId = prSelected.PullRequestId
		nwo := cmd.ToNwo(m.Repository)
		m.commitSelector.Title = fmt.Sprintf("Getting commits for %s %d...", nwo, m.PullRequestId)

		// Start fetching commits
		fetchCommits := func() tea.Msg {
			return cmd.FetchCommits(m.Repository, m.PullRequestId)
		}
		commands = append(commands, fetchCommits, m.commitSelector.StartSpinner())

	case cmd.CommitsFetched:
		fetched := cmd.CommitsFetched(msg)
		m.commits = fetched.Commits
		m.commitSelector.Title = "Select commits to move to another branch:"
		m.commitSelector.StopSpinner()
		items := newCommitItems(m.commits)
		commands = append(commands, m.commitSelector.SetItems(items))

	case tea.KeyMsg:
		if m.commitSelector.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.commitSelector.SelectedItem().(commitItem)
			if ok {
				index := m.commitSelector.Index()
				i.checked = !i.checked
				commands = append(commands, m.commitSelector.SetItem(index, i))
			}
		}
	}

	m.commitSelector, command = m.commitSelector.Update(msg)
	commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m SplitModel) View() string {
	if m.IsComplete() {
		return ""
	}

	return m.commitSelector.View()
}
