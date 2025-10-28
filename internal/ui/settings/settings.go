package settings

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/models"
	"github.com/shangyanjin/gocmder/internal/ui/dialogs"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	settingsNameColIndex = 0 + iota
	settingsDescriptionColIndex
	settingsSelectedColIndex
)

// Settings implements the settings page primitive
type Settings struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	confirmDialog   *dialogs.ConfirmDialog
	messageDialog   *dialogs.MessageDialog
	settingsList    settingsListReport
	confirmData     string
	applyHandler    func(settings []models.Setting)
	appFocusHandler func()
}

type settingsListReport struct {
	mu     sync.Mutex
	report []models.Setting
}

// NewSettings returns settings page view
func NewSettings() *Settings {
	settings := &Settings{
		Box:           tview.NewBox(),
		title:         "system settings",
		headers:       []string{"setting", "description", "selected"},
		errorDialog:   dialogs.NewErrorDialog(),
		confirmDialog: dialogs.NewConfirmDialog(),
		messageDialog: dialogs.NewMessageDialog(""),
		settingsList:  settingsListReport{},
	}

	settings.table = tview.NewTable()
	settings.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(settings.title)))
	settings.table.SetBorderColor(style.BorderColor)
	settings.table.SetTitleColor(style.FgColor)
	settings.table.SetBackgroundColor(style.BgColor)
	settings.table.SetBorder(true)

	// Set headers
	for i := range settings.headers {
		settings.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(settings.headers[i]))).
				SetExpansion(1).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	settings.table.SetFixed(1, 1)
	settings.table.SetSelectable(true, false)

	// Set error dialog functions with focus restoration
	settings.errorDialog.SetDoneFunc(func() {
		settings.errorDialog.Hide()
		if settings.appFocusHandler != nil {
			settings.appFocusHandler()
		}
	})

	// Set message dialog functions with focus restoration
	settings.messageDialog.SetCancelFunc(func() {
		settings.messageDialog.Hide()
		if settings.appFocusHandler != nil {
			settings.appFocusHandler()
		}
	})

	// Set confirm dialog functions with focus restoration
	settings.confirmDialog.SetSelectedFunc(func() {
		settings.confirmDialog.Hide()
		if settings.confirmData == "apply" {
			settings.applySettings()
		}
		if settings.appFocusHandler != nil {
			settings.appFocusHandler()
		}
	})
	settings.confirmDialog.SetCancelFunc(func() {
		settings.confirmDialog.Hide()
		if settings.appFocusHandler != nil {
			settings.appFocusHandler()
		}
	})

	return settings
}

// SetAppFocusHandler sets application focus handler
func (s *Settings) SetAppFocusHandler(handler func()) {
	s.appFocusHandler = handler
}

// GetTitle returns primitive title
func (s *Settings) GetTitle() string {
	return s.title
}

// HasFocus returns whether or not this primitive has focus
func (s *Settings) HasFocus() bool {
	if s.table.HasFocus() || s.errorDialog.HasFocus() {
		return true
	}

	if s.confirmDialog.HasFocus() || s.messageDialog.HasFocus() {
		return true
	}

	if s.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (s *Settings) SubDialogHasFocus() bool {
	if s.errorDialog.HasFocus() || s.confirmDialog.HasFocus() {
		return true
	}

	if s.messageDialog.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus
func (s *Settings) Focus(delegate func(p tview.Primitive)) {
	if s.errorDialog.IsDisplay() {
		delegate(s.errorDialog)
		return
	}

	if s.confirmDialog.IsDisplay() {
		delegate(s.confirmDialog)
		return
	}

	if s.messageDialog.IsDisplay() {
		delegate(s.messageDialog)
		return
	}

	delegate(s.table)
}

// HideAllDialogs hides all sub dialogs
func (s *Settings) HideAllDialogs() {
	if s.errorDialog.IsDisplay() {
		s.errorDialog.Hide()
	}

	if s.confirmDialog.IsDisplay() {
		s.confirmDialog.Hide()
	}

	if s.messageDialog.IsDisplay() {
		s.messageDialog.Hide()
	}
}

// SetApplyHandler sets the handler for applying settings
func (s *Settings) SetApplyHandler(handler func(settings []models.Setting)) {
	s.applyHandler = handler
}

// UpdateData updates the settings list data
func (s *Settings) UpdateData(settingsData []models.Setting) {
	s.settingsList.mu.Lock()
	defer s.settingsList.mu.Unlock()

	s.settingsList.report = settingsData
	s.updateTable()
}

// updateTable updates the table display
func (s *Settings) updateTable() {
	// Clear existing rows (except header)
	for row := s.table.GetRowCount() - 1; row > 0; row-- {
		s.table.RemoveRow(row)
	}

	// Define setting descriptions
	descriptions := map[string]string{
		"Add to PATH":              "Add development tools to system PATH",
		"Configure Power Settings": "Set power configuration for development",
		"Set User Directories":     "Configure user folders (Desktop, Documents, etc.)",
	}

	// Add setting rows
	for i, setting := range s.settingsList.report {
		row := i + 1

		// Setting name
		s.table.SetCell(row, settingsNameColIndex,
			tview.NewTableCell(setting.Name).
				SetTextColor(style.FgColor).
				SetAlign(tview.AlignLeft))

		// Description
		description := descriptions[setting.Name]
		if description == "" {
			description = "System configuration"
		}

		s.table.SetCell(row, settingsDescriptionColIndex,
			tview.NewTableCell(description).
				SetTextColor(style.FgColor).
				SetAlign(tview.AlignLeft))

		// Selected
		selectedText := "[ ]"
		if setting.Selected {
			selectedText = "[X]"
		}

		s.table.SetCell(row, settingsSelectedColIndex,
			tview.NewTableCell(selectedText).
				SetTextColor(style.StatusSelectedColor).
				SetAlign(tview.AlignCenter))
	}
}

// ToggleSelection toggles the selection of the current setting
func (s *Settings) ToggleSelection() {
	if s.table.GetRowCount() <= 1 {
		return
	}

	row, _ := s.table.GetSelection()
	if row < 1 {
		return
	}

	s.settingsList.mu.Lock()
	defer s.settingsList.mu.Unlock()

	settingIndex := row - 1
	if settingIndex < len(s.settingsList.report) {
		s.settingsList.report[settingIndex].Selected = !s.settingsList.report[settingIndex].Selected
		s.updateTable()
		s.table.Select(row, 0)
	}
}

// SelectAll selects all settings
func (s *Settings) SelectAll() {
	s.settingsList.mu.Lock()
	defer s.settingsList.mu.Unlock()

	for i := range s.settingsList.report {
		s.settingsList.report[i].Selected = true
	}
	s.updateTable()
}

// DeselectAll deselects all settings
func (s *Settings) DeselectAll() {
	s.settingsList.mu.Lock()
	defer s.settingsList.mu.Unlock()

	for i := range s.settingsList.report {
		s.settingsList.report[i].Selected = false
	}
	s.updateTable()
}

// GetSelectedSettings returns the list of selected settings
func (s *Settings) GetSelectedSettings() []models.Setting {
	s.settingsList.mu.Lock()
	defer s.settingsList.mu.Unlock()

	var selected []models.Setting
	for _, setting := range s.settingsList.report {
		if setting.Selected {
			selected = append(selected, setting)
		}
	}

	return selected
}

// ShowApplyConfirmation shows apply confirmation dialog
func (s *Settings) ShowApplyConfirmation() {
	selected := s.GetSelectedSettings()
	if len(selected) == 0 {
		s.errorDialog.SetTitle("Error")
		s.errorDialog.SetText("No settings selected")
		s.errorDialog.Display()
		return
	}

	s.confirmData = "apply"
	s.confirmDialog.SetTitle("Confirm Settings")
	s.confirmDialog.SetText(fmt.Sprintf("Apply %d selected setting(s)?", len(selected)))
	s.confirmDialog.Display()
}

// applySettings applies selected settings
func (s *Settings) applySettings() {
	selected := s.GetSelectedSettings()
	if len(selected) == 0 {
		return
	}

	if s.applyHandler != nil {
		s.applyHandler(selected)
	}
}

// InputHandler returns the input handler for this primitive
func (s *Settings) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Handle dialog input first if dialog has focus
		if s.SubDialogHasFocus() {
			// Pass event to the dialog that has focus
			if s.errorDialog.HasFocus() {
				if handler := s.errorDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if s.confirmDialog.HasFocus() {
				if handler := s.confirmDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			} else if s.messageDialog.HasFocus() {
				if handler := s.messageDialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
			return
		}

		// Handle table input
		if s.table.HasFocus() {
			// Space key toggles selection
			if event.Key() == utils.ToggleKey.Key && event.Rune() == utils.ToggleKey.Rune {
				s.ToggleSelection()
				return
			}

			// 'a' key selects all
			if event.Rune() == utils.SelectAllKey.Rune {
				s.SelectAll()
				return
			}

			// Enter key applies settings
			if event.Key() == utils.ConfirmKey.Key {
				s.ShowApplyConfirmation()
				return
			}

			// Default table handler
			if tableHandler := s.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
				return
			}
		}
	})
}

// Draw draws this primitive onto the screen
func (s *Settings) Draw(screen tcell.Screen) {
	s.Box.DrawForSubclass(screen, s)
	x, y, width, height := s.GetInnerRect()

	// Draw table
	s.table.SetRect(x, y, width, height)
	s.table.Draw(screen)

	// Draw dialogs
	if s.errorDialog.IsDisplay() {
		s.errorDialog.SetRect(x, y, width, height)
		s.errorDialog.Draw(screen)
	}

	if s.confirmDialog.IsDisplay() {
		s.confirmDialog.SetRect(x, y, width, height)
		s.confirmDialog.Draw(screen)
	}

	if s.messageDialog.IsDisplay() {
		s.messageDialog.SetRect(x, y, width, height)
		s.messageDialog.Draw(screen)
	}
}
