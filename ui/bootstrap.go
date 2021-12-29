package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
	"github.com/pjquirk/gh-splitpr/cmd"
)

type prItem struct {
	number int
	title  string
	author string
	filter string
}

func (i prItem) FilterValue() string { return i.filter }

func newPrItem(pr cmd.PullRequest) prItem {
	return prItem{
		number: pr.Number,
		title:  pr.Title,
		author: pr.Author,
		filter: fmt.Sprintf("%s %s", pr.Author, pr.Title),
	}
}

func newPrItems(pullRequests []cmd.PullRequest) []list.Item {
	items := make([]list.Item, len(pullRequests))
	for i, v := range pullRequests {
		items[i] = newPrItem(v)
	}
	return items
}

type prItemDelegate struct{}

func (d prItemDelegate) Height() int                               { return 1 }
func (d prItemDelegate) Spacing() int                              { return 0 }
func (d prItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d prItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(prItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("#%d - @%s - %s", i.number, i.author, i.title)

	fn := defaultItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return defaultSelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type BootstrapModel struct {
	Repository    gh.Repository
	PullRequestId int

	pullRequests []cmd.PullRequest
	prSelector   list.Model
}

func NewBootstrapModel() BootstrapModel {
	return BootstrapModel{
		Repository:    nil,
		PullRequestId: -1,
		pullRequests:  []cmd.PullRequest{},
		prSelector:    newListModel(5, defaultItemStyle),
	}
}

func (m BootstrapModel) IsComplete() bool {
	return m.Repository != nil && m.PullRequestId > 0
}

func (m BootstrapModel) Init() tea.Cmd {
	return nil
}

func (m BootstrapModel) Update(msg tea.Msg) (BootstrapModel, tea.Cmd) {
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
		items := newPrItems(m.pullRequests)
		commands = append(commands, m.prSelector.SetItems(items))

	case tea.KeyMsg:
		if m.prSelector.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.prSelector.SelectedItem().(prItem)
			if ok {
				m.PullRequestId = i.number
			}
			return m, nil
		}
	}

	m.prSelector, command = m.prSelector.Update(msg)
	commands = append(commands, command)

	return m, tea.Batch(commands...)
}

func (m BootstrapModel) View() string {
	if m.IsComplete() {
		return ""
	}

	//nwo := cmd.ToNwo(m.Repository)
	return m.prSelector.View()
}
