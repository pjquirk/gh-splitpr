package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	titleBarStyle            = list.DefaultStyles().TitleBar
	titleStyle               = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle          = list.DefaultStyles().PaginationStyle.PaddingLeft(4).PaddingBottom(0)
	helpStyle                = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	defaultItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	defaultSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func newListModel(pageSize int, itemStyle lipgloss.Style) list.Model {
	linesPerItem := (itemStyle.GetVerticalFrameSize() + 1)
	titleBarHeight := titleBarStyle.GetVerticalFrameSize()
	// Add one for each line of text
	titleHeight := titleStyle.GetVerticalFrameSize() + 1
	pageHeight := paginationStyle.GetVerticalFrameSize() + 1
	helpHeight := helpStyle.GetVerticalFrameSize() + 1
	height := titleBarHeight + titleHeight + pageHeight + helpHeight + (pageSize * linesPerItem)

	newModel := list.NewModel([]list.Item{}, prItemDelegate{}, 0, height)
	newModel.SetShowStatusBar(false)
	newModel.Styles.Title = titleStyle
	newModel.Styles.TitleBar = titleBarStyle
	newModel.Styles.PaginationStyle = paginationStyle
	newModel.Styles.HelpStyle = helpStyle
	return newModel
}
