package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// EmptyBoxSpace returns an empty box with the specified background color
func EmptyBoxSpace(bgColor tcell.Color) *tview.Box {
	box := tview.NewBox()
	box.SetBackgroundColor(bgColor)
	return box
}

// AlignStringRight aligns a string to the right within the specified width
func AlignStringRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	spaces := width - len(s)
	result := ""
	for i := 0; i < spaces; i++ {
		result += " "
	}
	return result + s
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// CreateStyledTable creates a new table with standard styling
func CreateStyledTable(title string, headers []string) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle(title)
	table.SetTitleColor(style.FgColor)
	table.SetBorderColor(style.BorderColor)
	table.SetBackgroundColor(style.BgColor)

	// Set headers
	for i, header := range headers {
		cell := tview.NewTableCell(header)
		cell.SetExpansion(1)
		cell.SetBackgroundColor(style.PageHeaderBgColor)
		cell.SetTextColor(style.PageHeaderFgColor)
		cell.SetAlign(tview.AlignLeft)
		cell.SetSelectable(false)
		table.SetCell(0, i, cell)
	}

	table.SetFixed(1, 1)
	table.SetSelectable(true, false)

	return table
}
