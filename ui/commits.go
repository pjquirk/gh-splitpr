package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	commitSelectedItemStyle = defaultSelectedItemStyle.Copy().PaddingLeft(4)
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

type listKeyMap struct {
	toggleItem      key.Binding
	finishSelection key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleItem: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle item"),
		),
		finishSelection: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "finish"),
		),
	}
}

type CommitsModel struct {
	Repository    gh.Repository
	PullRequestId int

	commits         []cmd.Commit
	selectedCommits []cmd.Commit
	commitSelector  list.Model
	keys            *listKeyMap
}

func NewCommitsModel() CommitsModel {
	commitsModel := CommitsModel{
		Repository:      nil,
		PullRequestId:   -1,
		commits:         []cmd.Commit{},
		selectedCommits: nil,
		commitSelector:  newListModel(5, commitItemDelegate{}, commitItemStyle),
		keys:            newListKeyMap(),
	}
	commitsModel.commitSelector.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			commitsModel.keys.toggleItem,
			commitsModel.keys.finishSelection,
		}
	}
	commitsModel.commitSelector.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			commitsModel.keys.toggleItem,
			commitsModel.keys.finishSelection,
		}
	}
	return commitsModel
}

func (m CommitsModel) IsComplete() bool {
	return m.selectedCommits != nil
}

func (m CommitsModel) Init() tea.Cmd {
	return nil
}

func (m CommitsModel) Update(msg tea.Msg) (CommitsModel, tea.Cmd) {
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

		switch {
		case key.Matches(msg, m.keys.toggleItem):
			i, ok := m.commitSelector.SelectedItem().(commitItem)
			if ok {
				index := m.commitSelector.Index()
				i.checked = !i.checked
				commands = append(commands, m.commitSelector.SetItem(index, i))
			}

		case key.Matches(msg, m.keys.finishSelection):
			m.selectedCommits = selectedCommits(m)
			commands = append(commands, func() tea.Msg {
				return cmd.CommitsSelected{
					Repository:    m.Repository,
					PullRequestId: m.PullRequestId,
					Commits:       m.selectedCommits,
				}
			})
		}
	}

	m.commitSelector, command = m.commitSelector.Update(msg)
	commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m CommitsModel) View() string {
	if m.IsComplete() {
		return ""
	}

	return m.commitSelector.View()
}

func selectedCommits(m CommitsModel) (selected []cmd.Commit) {
	allItems := m.commitSelector.Items()
	for _, v := range allItems {
		commit := v.(commitItem)
		if commit.checked {
			selected = append(selected, cmd.Commit{
				Sha:     commit.sha,
				Comment: commit.comment,
			})
		}
	}
	return
}
