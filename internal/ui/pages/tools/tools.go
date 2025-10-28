package tools

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/models"
	"github.com/shangyanjin/gocmder/internal/ui/components/dialogs"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	toolsNameColIndex = 0 + iota
	toolsVersionColIndex
	toolsSizeColIndex
	toolsStatusColIndex
	toolsSelectedColIndex
)

// Tools implements the tools page primitive
type Tools struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	confirmDialog   *dialogs.ConfirmDialog
	messageDialog   *dialogs.MessageDialog
	progressDialog  *dialogs.ProgressDialog
	inputDialog     *dialogs.SimpleInputDialog
	toolsList       toolsListReport
	selectedID      int
	confirmData     string
	installHandler  func(toolName string)
	refreshHandler  func()
	appFocusHandler func()
}

type toolsListReport struct {
	mu     sync.Mutex
	report []models.Tool
}

// NewTools returns tools page view
func NewTools() *Tools {
	tools := &Tools{
		Box:            tview.NewBox(),
		title:          "development tools",
		headers:        []string{"tool", "version", "size", "status", "selected"},
		errorDialog:    dialogs.NewErrorDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		progressDialog: dialogs.NewProgressDialog(),
		inputDialog:    dialogs.NewSimpleInputDialog(""),
		toolsList:      toolsListReport{},
	}

	tools.table = tview.NewTable()
	tools.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(tools.title)))
	tools.table.SetBorderColor(style.BorderColor)
	tools.table.SetTitleColor(style.FgColor)
	tools.table.SetBackgroundColor(style.BgColor)
	tools.table.SetBorder(true)

	// Set headers
	for i := range tools.headers {
		tools.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(tools.headers[i]))).
				SetExpansion(1).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	tools.table.SetFixed(1, 1)
	tools.table.SetSelectable(true, false)

	// Set error dialog functions with focus restoration
	tools.errorDialog.SetDoneFunc(func() {
		tools.errorDialog.Hide()
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})

	// Set message dialog functions with focus restoration
	tools.messageDialog.SetCancelFunc(func() {
		tools.messageDialog.Hide()
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})

	// Set confirm dialog functions with focus restoration
	tools.confirmDialog.SetSelectedFunc(func() {
		tools.confirmDialog.Hide()
		switch tools.confirmData {
		case "install":
			tools.installSelected()
		case "install_all":
			tools.installAll()
		}
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})
	tools.confirmDialog.SetCancelFunc(func() {
		tools.confirmDialog.Hide()
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})

	// Set input dialog functions with focus restoration
	tools.inputDialog.SetSelectedFunc(func() {
		tools.inputDialog.Hide()
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})
	tools.inputDialog.SetCancelFunc(func() {
		tools.inputDialog.Hide()
		if tools.appFocusHandler != nil {
			tools.appFocusHandler()
		}
	})

	return tools
}

// SetAppFocusHandler sets application focus handler
func (t *Tools) SetAppFocusHandler(handler func()) {
	t.appFocusHandler = handler
}

// GetTitle returns primitive title
func (t *Tools) GetTitle() string {
	return t.title
}

// HasFocus returns whether or not this primitive has focus
func (t *Tools) HasFocus() bool {
	if t.table.HasFocus() || t.errorDialog.HasFocus() {
		return true
	}

	if t.confirmDialog.HasFocus() || t.messageDialog.HasFocus() {
		return true
	}

	if t.progressDialog.HasFocus() || t.inputDialog.HasFocus() {
		return true
	}

	if t.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (t *Tools) SubDialogHasFocus() bool {
	if t.errorDialog.HasFocus() || t.confirmDialog.HasFocus() {
		return true
	}

	if t.messageDialog.HasFocus() || t.progressDialog.HasFocus() {
		return true
	}

	if t.inputDialog.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus
func (t *Tools) Focus(delegate func(p tview.Primitive)) {
	if t.errorDialog.IsDisplay() {
		delegate(t.errorDialog)
		return
	}

	if t.confirmDialog.IsDisplay() {
		delegate(t.confirmDialog)
		return
	}

	if t.messageDialog.IsDisplay() {
		delegate(t.messageDialog)
		return
	}

	if t.progressDialog.IsDisplay() {
		delegate(t.progressDialog)
		return
	}

	if t.inputDialog.IsDisplay() {
		delegate(t.inputDialog)
		return
	}

	delegate(t.table)
}

// HideAllDialogs hides all sub dialogs
func (t *Tools) HideAllDialogs() {
	if t.errorDialog.IsDisplay() {
		t.errorDialog.Hide()
	}

	if t.confirmDialog.IsDisplay() {
		t.confirmDialog.Hide()
	}

	if t.messageDialog.IsDisplay() {
		t.messageDialog.Hide()
	}

	if t.progressDialog.IsDisplay() {
		t.progressDialog.Hide()
	}

	if t.inputDialog.IsDisplay() {
		t.inputDialog.Hide()
	}
}

// SetInstallHandler sets the handler for tool installation
func (t *Tools) SetInstallHandler(handler func(toolName string)) {
	t.installHandler = handler
}

// SetRefreshHandler sets the handler for refreshing tool list
func (t *Tools) SetRefreshHandler(handler func()) {
	t.refreshHandler = handler
}

// UpdateData updates the tools list data
func (t *Tools) UpdateData(toolsData []models.Tool) {
	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	t.toolsList.report = toolsData
	t.updateTable()
}

// updateTable updates the table display
func (t *Tools) updateTable() {
	// Clear existing rows (except header)
	for row := t.table.GetRowCount() - 1; row > 0; row-- {
		t.table.RemoveRow(row)
	}

	// Add tool rows
	for i, tool := range t.toolsList.report {
		row := i + 1

		// Tool name
		t.table.SetCell(row, toolsNameColIndex,
			tview.NewTableCell(tool.Name).
				SetTextColor(style.FgColor).
				SetAlign(tview.AlignLeft))

		// Version
		t.table.SetCell(row, toolsVersionColIndex,
			tview.NewTableCell(tool.Version).
				SetTextColor(style.FgColor).
				SetAlign(tview.AlignLeft))

		// Size
		t.table.SetCell(row, toolsSizeColIndex,
			tview.NewTableCell(tool.Size).
				SetTextColor(style.FgColor).
				SetAlign(tview.AlignLeft))

		// Status
		statusText := "Not Installed"
		statusColor := style.StatusNotInstalledColor
		if tool.Installed {
			statusText = "Installed"
			statusColor = style.StatusInstalledColor
		}

		t.table.SetCell(row, toolsStatusColIndex,
			tview.NewTableCell(statusText).
				SetTextColor(statusColor).
				SetAlign(tview.AlignLeft))

		// Selected
		selectedText := "[ ]"
		if tool.Selected {
			selectedText = "[X]"
		}

		t.table.SetCell(row, toolsSelectedColIndex,
			tview.NewTableCell(selectedText).
				SetTextColor(style.StatusSelectedColor).
				SetAlign(tview.AlignCenter))
	}
}

// ToggleSelection toggles the selection of the current tool
func (t *Tools) ToggleSelection() {
	if t.table.GetRowCount() <= 1 {
		return
	}

	row, _ := t.table.GetSelection()
	if row < 1 {
		return
	}

	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	toolIndex := row - 1
	if toolIndex < len(t.toolsList.report) {
		t.toolsList.report[toolIndex].Selected = !t.toolsList.report[toolIndex].Selected
		t.updateTable()
		t.table.Select(row, 0)
	}
}

// SelectAll selects all tools
func (t *Tools) SelectAll() {
	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	for i := range t.toolsList.report {
		t.toolsList.report[i].Selected = true
	}
	t.updateTable()
}

// DeselectAll deselects all tools
func (t *Tools) DeselectAll() {
	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	for i := range t.toolsList.report {
		t.toolsList.report[i].Selected = false
	}
	t.updateTable()
}

// GetSelectedTools returns the list of selected tools
func (t *Tools) GetSelectedTools() []models.Tool {
	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	var selected []models.Tool
	for _, tool := range t.toolsList.report {
		if tool.Selected {
			selected = append(selected, tool)
		}
	}

	return selected
}

// ShowInstallConfirmation shows installation confirmation dialog
func (t *Tools) ShowInstallConfirmation() {
	selected := t.GetSelectedTools()
	if len(selected) == 0 {
		t.errorDialog.SetTitle("Error")
		t.errorDialog.SetText("No tools selected for installation")
		t.errorDialog.Display()
		return
	}

	t.confirmData = "install"
	t.confirmDialog.SetTitle("Confirm Installation")
	t.confirmDialog.SetText(fmt.Sprintf("Install %d selected tool(s)?", len(selected)))
	t.confirmDialog.Display()
}

// installSelected installs selected tools
func (t *Tools) installSelected() {
	selected := t.GetSelectedTools()
	if len(selected) == 0 {
		return
	}

	if t.installHandler != nil {
		for _, tool := range selected {
			t.installHandler(tool.Name)
		}
	}
}

// installAll installs all tools
func (t *Tools) installAll() {
	t.toolsList.mu.Lock()
	defer t.toolsList.mu.Unlock()

	if t.installHandler != nil {
		for _, tool := range t.toolsList.report {
			t.installHandler(tool.Name)
		}
	}
}

// InputHandler returns the input handler for this primitive
func (t *Tools) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Handle dialog input first if dialog has focus
		if t.SubDialogHasFocus() {
			// Pass event to the dialog that has focus
			if t.errorDialog.HasFocus() {
				if handler := t.errorDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if t.confirmDialog.HasFocus() {
				if handler := t.confirmDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if t.messageDialog.HasFocus() {
				if handler := t.messageDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if t.progressDialog.HasFocus() {
				if handler := t.progressDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if t.inputDialog.HasFocus() {
				if handler := t.inputDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
			return
		}

		// Handle table input
		if t.table.HasFocus() {
			// Space key toggles selection
			if event.Key() == utils.ToggleKey.Key && event.Rune() == utils.ToggleKey.Rune {
				t.ToggleSelection()
				return
			}

			// 'a' key selects all
			if event.Rune() == utils.SelectAllKey.Rune {
				t.SelectAll()
				return
			}

			// 'i' key shows install confirmation
			if event.Rune() == utils.InstallKey.Rune {
				t.ShowInstallConfirmation()
				return
			}

			// 'r' key refreshes
			if event.Rune() == utils.RefreshKey.Rune {
				if t.refreshHandler != nil {
					t.refreshHandler()
				}
				return
			}

			// Default table handler
			if tableHandler := t.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
				return
			}
		}
	})
}

// Draw draws this primitive onto the screen
func (t *Tools) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)
	x, y, width, height := t.GetInnerRect()

	// Draw table
	t.table.SetRect(x, y, width, height)
	t.table.Draw(screen)

	// Draw dialogs
	if t.errorDialog.IsDisplay() {
		t.errorDialog.SetRect(x, y, width, height)
		t.errorDialog.Draw(screen)
	}

	if t.confirmDialog.IsDisplay() {
		t.confirmDialog.SetRect(x, y, width, height)
		t.confirmDialog.Draw(screen)
	}

	if t.messageDialog.IsDisplay() {
		t.messageDialog.SetRect(x, y, width, height)
		t.messageDialog.Draw(screen)
	}

	if t.progressDialog.IsDisplay() {
		t.progressDialog.SetRect(x, y, width, height)
		t.progressDialog.Draw(screen)
	}

	if t.inputDialog.IsDisplay() {
		t.inputDialog.SetRect(x, y, width, height)
		t.inputDialog.Draw(screen)
	}
}
