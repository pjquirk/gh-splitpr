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

type item struct {
	number int
	title  string
	author string
	filter string
}

func newItem(pr cmd.PullRequest) item {
	return item{
		number: pr.Number,
		title:  pr.Title,
		author: pr.Author,
		filter: fmt.Sprintf("%s %s", pr.Author, pr.Title),
	}
}

func newItems(pullRequests []cmd.PullRequest) []list.Item {
	items := make([]list.Item, len(pullRequests))
	for i, v := range pullRequests {
		items[i] = newItem(v)
	}
	return items
}

func (i item) FilterValue() string { return i.filter }

var (
	titleBarStyle     = list.DefaultStyles().TitleBar
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4).PaddingBottom(0)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("#%d - @%s - %s", i.number, i.author, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

func newListModel() list.Model {
	pageSize := 5 * (itemStyle.GetVerticalFrameSize() + 1)
	titleBarHeight := titleBarStyle.GetVerticalFrameSize()
	// Add one for each line of text
	titleHeight := titleStyle.GetVerticalFrameSize() + 1
	pageHeight := paginationStyle.GetVerticalFrameSize() + 1
	helpHeight := helpStyle.GetVerticalFrameSize() + 1
	height := titleBarHeight + titleHeight + pageHeight + helpHeight + pageSize

	newModel := list.NewModel([]list.Item{}, itemDelegate{}, 0, height)
	newModel.SetShowStatusBar(false)
	newModel.Styles.Title = titleStyle
	newModel.Styles.TitleBar = titleBarStyle
	newModel.Styles.PaginationStyle = paginationStyle
	newModel.Styles.HelpStyle = helpStyle
	return newModel
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
		prSelector:    newListModel(),
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
		commands = append(commands, m.prSelector.SetItems(newItems(m.pullRequests)))

	case tea.KeyMsg:
		if m.prSelector.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.prSelector.SelectedItem().(item)
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
