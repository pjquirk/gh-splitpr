package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type commitItem struct {
	sha     string
	comment string
}

func (i commitItem) FilterValue() string { return i.comment }

func newCommitItem(pr cmd.PullRequest) commitItem {
	return commitItem{
		sha:     pr.Author,
		comment: pr.Title,
	}
}

func newCommitItems(pullRequests []cmd.PullRequest) []list.Item {
	items := make([]list.Item, len(pullRequests))
	for i, v := range pullRequests {
		items[i] = newCommitItem(v)
	}
	return items
}

var (
	commitItemStyle         = defaultItemStyle
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

	str := fmt.Sprintf("%s - %s", i.sha, i.comment)

	fn := commitItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return commitSelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type SplitModel struct {
	Repository    gh.Repository
	PullRequestId int

	pullRequests []cmd.PullRequest
	prSelector   list.Model
}

func NewSplitModel() SplitModel {
	return SplitModel{
		Repository:    nil,
		PullRequestId: -1,
		pullRequests:  []cmd.PullRequest{},
		prSelector:    newListModel(5, commitItemStyle),
	}
}

func (m SplitModel) IsComplete() bool {
	return m.Repository != nil && m.PullRequestId > 0
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
		m.prSelector.SetWidth(msg.Width)

	case cmd.OptionsParsed:
		options := cmd.OptionsParsed(msg)
		m.Repository = options.Repository
		m.PullRequestId = options.PullRequest
		if m.PullRequestId <= 0 {
			m.prSelector.Title = "Looking for pull requests..."
			// Start fetching pull requests
			fetchPullRequests := func() tea.Msg {
				return cmd.FetchPullRequests(m.Repository)
			}
			commands = append(commands, fetchPullRequests, m.prSelector.StartSpinner())
		}

	case cmd.PullRequestsFetched:
		fetched := cmd.PullRequestsFetched(msg)
		m.pullRequests = fetched.PullRequests
		m.prSelector.Title = "Select a pull request to split:"
		m.prSelector.StopSpinner()
		commands = append(commands, m.prSelector.SetItems(newCommitItems(m.pullRequests)))

	case tea.KeyMsg:
		if m.prSelector.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "enter":
			_, ok := m.prSelector.SelectedItem().(commitItem)
			if ok {
				//m.PullRequestId = i.n
			}
			return m, nil
		}
	}

	m.prSelector, command = m.prSelector.Update(msg)
	commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m SplitModel) View() string {
	if m.IsComplete() {
		return ""
	}

	//nwo := cmd.ToNwo(m.Repository)
	return m.prSelector.View()
}
