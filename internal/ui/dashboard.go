package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/models"
)

// Dashboard represents the main dashboard
type Dashboard struct {
	app           *tview.Application
	config        *models.InstallConfig
	currentPanel  string
	selectedIndex int
	message       string
	messageType   string

	// Navigation state for two-level menu
	navigationMode string // "scheme" or "tools"
	focusedPanel   string // "left" or "right" - which panel has focus

	// UI Components
	infoPanel      *tview.TextView
	titlePanel     *tview.TextView
	optionsTable   *tview.Table // Left panel
	detailTable    *tview.Table // Right panel - now interactive
	setupContainer *tview.Flex
	outputPanel    *tview.TextView
	inputContainer *tview.Flex
	inputField     *tview.InputField
	inputInfoText  *tview.TextView
	rightPanels    *tview.Flex
	middleRow      *tview.Flex
	mainLayout     *tview.Flex
}

// NewDashboard creates a new dashboard
func NewDashboard(app *tview.Application) *Dashboard {
	dashboard := &Dashboard{
		app:            app,
		config:         models.NewInstallConfig(),
		currentPanel:   "tools",
		selectedIndex:  0,
		message:        "",
		messageType:    "",
		navigationMode: "scheme", // Start with scheme selection
		focusedPanel:   "left",   // Start with left panel focused
	}

	dashboard.setupUI()
	return dashboard
}

// setupUI sets up the user interface
func (d *Dashboard) setupUI() {
	// Initialize system info cache
	InitializeSystemInfoCache()

	// Create all UI components
	d.createComponents()
	d.setupLayout()
	d.setupKeyHandlers()
	d.updateContent()

	// Set initial focus to options table and select first item
	d.app.SetFocus(d.optionsTable)
	if d.optionsTable.GetRowCount() > 0 {
		d.optionsTable.Select(0, 0)
	}
}

// createComponents creates all UI components
func (d *Dashboard) createComponents() {
	// Create title panel (top)
	d.titlePanel = tview.NewTextView()
	d.titlePanel.SetDynamicColors(true).
		SetText("[yellow]GoCmder Setup[white] - [cyan]Scheme Mode[white]")

	// Create left panel - list of schemes or tools
	d.optionsTable = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(0, 0)
	d.optionsTable.SetSelectedStyle(
		tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorDarkBlue),
	)

	leftPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.optionsTable, 0, 1, true)
	leftPanel.SetBorder(true).SetTitle("Options")

	// Create right panel - details table (now interactive)
	d.detailTable = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(0, 0)
	d.detailTable.SetSelectedStyle(
		tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorDarkBlue),
	)

	// Wrap right panel in Flex and add border
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.detailTable, 0, 1, false)
	rightPanel.SetBorder(true).SetTitle("Details")

	// Create middle layout - left and right panels
	middleLayout := tview.NewFlex().
		AddItem(leftPanel, 25, 0, true).
		AddItem(rightPanel, 0, 1, false)

	// Create footer/status bar
	d.outputPanel = tview.NewTextView()
	d.outputPanel.SetDynamicColors(true).
		SetText("[cyan][Y/n/?][white] Space:Select  Tab:Switch  F5:Run  F9:Output  Esc:Return  Ctrl+C:Exit")

	// Create main layout
	d.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.titlePanel, 1, 0, false).
		AddItem(middleLayout, 0, 1, true).
		AddItem(d.outputPanel, 1, 0, false)

	// Keep input field for command execution (hidden until F5)
	d.inputField = tview.NewInputField()
	d.inputField.SetLabel("> ").SetFieldWidth(0)

	// Keep other components for potential use
	d.setupContainer = leftPanel
	d.rightPanels = rightPanel
	d.inputContainer = tview.NewFlex().AddItem(d.inputField, 0, 1, true)
	d.middleRow = middleLayout
}

// setupLayout sets up the layout structure
func (d *Dashboard) setupLayout() {
	// Layout is already set up in createComponents
}

// setupKeyHandlers sets up keyboard event handlers
func (d *Dashboard) setupKeyHandlers() {
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEsc:
			// ESC: Return to Scheme Mode if in Tools Mode, otherwise no action
			if d.navigationMode == "tools" {
				d.navigationMode = "scheme"
				d.updateMessage("Returned to Scheme Mode (use Ctrl+C to exit)", "info")
				d.updateContent()
				return nil
			}
			// In Scheme Mode, ESC does nothing - use Ctrl+C to exit
			return nil

		case event.Modifiers()&tcell.ModCtrl != 0 && event.Rune() == 'c':
			// Ctrl+C: Direct exit
			d.app.Stop()
			return nil

		case event.Key() == tcell.KeyTab:
			// Tab to switch focus between left and right panels
			if d.focusedPanel == "left" {
				d.focusedPanel = "right"
				d.app.SetFocus(d.detailTable)
				d.updateMessage("Switched to Details panel", "info")
			} else {
				d.focusedPanel = "left"
				d.app.SetFocus(d.optionsTable)
				d.updateMessage("Switched to Options panel", "info")
			}
			return nil

		case event.Key() == tcell.KeyRune && event.Rune() == ' ':
			if d.app.GetFocus() == d.optionsTable {
				d.toggleCurrentOption()
				return nil
			} else if d.app.GetFocus() == d.detailTable {
				// Right panel - toggle detail item (works in both Scheme and Tools mode)
				d.toggleDetailItem()
				return nil
			}

		case event.Key() == tcell.KeyF9:
			d.app.SetFocus(d.outputPanel)
			return nil

		case event.Key() == tcell.KeyF2:
			d.app.SetFocus(d.optionsTable)
			return nil

		case event.Key() == tcell.KeyF5:
			d.app.SetFocus(d.inputField)
			return nil

		case event.Key() == tcell.KeyF3:
			// Manual refresh of system information
			UpdateSystemInfoAsync()
			d.updateMessage("System info refresh triggered", "info")
			return nil

		case event.Key() == tcell.KeyUp:
			if d.app.GetFocus() == d.optionsTable {
				d.cycleUp()
				return nil
			}

		case event.Key() == tcell.KeyDown:
			if d.app.GetFocus() == d.optionsTable {
				d.cycleDown()
				return nil
			}

		case event.Modifiers()&tcell.ModCtrl != 0 && event.Rune() >= '1' && event.Rune() <= '6':
			optionIndex := int(event.Rune() - '1')
			if optionIndex < d.optionsTable.GetRowCount() {
				d.optionsTable.Select(optionIndex, 0)
				d.app.SetFocus(d.optionsTable)
				d.updateMessage(fmt.Sprintf("Selected option %d", optionIndex+1), "info")
			}
			return nil
		}
		return event
	})
}

// cycleUp cycles selection up with wrapping from first to last
func (d *Dashboard) cycleUp() {
	currentRow, _ := d.optionsTable.GetSelection()
	totalRows := d.optionsTable.GetRowCount()

	if totalRows == 0 {
		return
	}

	nextRow := currentRow - 1
	if nextRow < 0 {
		nextRow = totalRows - 1
	}

	d.optionsTable.Select(nextRow, 0)
	// Dynamically update right panel when selection changes
	d.updateDetailTable()
}

// cycleDown cycles selection down with wrapping from last to first
func (d *Dashboard) cycleDown() {
	currentRow, _ := d.optionsTable.GetSelection()
	totalRows := d.optionsTable.GetRowCount()

	if totalRows == 0 {
		return
	}

	nextRow := currentRow + 1
	if nextRow >= totalRows {
		nextRow = 0
	}

	d.optionsTable.Select(nextRow, 0)
	// Dynamically update right panel when selection changes
	d.updateDetailTable()
}

// toggleCurrentOption toggles the selection state of the currently focused item or applies a scheme
func (d *Dashboard) toggleCurrentOption() {
	currentRow, _ := d.optionsTable.GetSelection()

	// Handle scheme mode - select a predefined scheme
	if d.navigationMode == "scheme" {
		if currentRow >= 0 && currentRow < len(d.config.Schemes) {
			scheme := d.config.Schemes[currentRow]

			// Check if Exit option is selected
			if scheme.Name == "Exit" {
				d.updateMessage("Exiting application...", "info")
				d.app.Stop()
				return
			}

			// Check if Personal Settings is selected - apply personal settings configuration
			if scheme.Name == "Personal Settings" {
				d.config.ApplyScheme(currentRow)
				d.updateMessage("Personal Settings configured: Add PATH, PowerConfig, SetUserDirs enabled", "success")
				d.updateContent()
				return
			}

			d.config.ApplyScheme(currentRow)
			d.updateMessage(fmt.Sprintf("Applied scheme: %s", scheme.Name), "success")
			d.updateContent()
		}
		return
	}

	// Handle tools mode - toggle tool or setting selection
	total := len(d.config.Tools)

	if currentRow < total {
		// Toggle tool
		d.config.Tools[currentRow].Selected = !d.config.Tools[currentRow].Selected
		tool := d.config.Tools[currentRow]
		status := "selected"
		if !tool.Selected {
			status = "deselected"
		}
		d.updateMessage(fmt.Sprintf("Tool '%s' %s", tool.Name, status), "info")
		d.refreshOptionsTable()
	} else {
		// Toggle setting
		settingIdx := currentRow - total
		if settingIdx < len(d.config.Settings) {
			d.config.Settings[settingIdx].Selected = !d.config.Settings[settingIdx].Selected
			setting := d.config.Settings[settingIdx]
			status := "enabled"
			if !setting.Selected {
				status = "disabled"
			}
			d.updateMessage(fmt.Sprintf("Setting '%s' %s", setting.Name, status), "info")
			d.refreshOptionsTable()
		}
	}
}

// refreshOptionsTable updates the options table display
func (d *Dashboard) refreshOptionsTable() {
	currentRow, _ := d.optionsTable.GetSelection()
	d.updateOptionsTable()
	if currentRow >= 0 {
		d.optionsTable.Select(currentRow, 0)
	}
}

// updateContent updates all panel content
func (d *Dashboard) updateContent() {
	d.updateTitlePanel()
	d.updateOptionsTable()
	d.updateDetailTable() // Update right panel details
	d.updateOutputPanel()
}

// updateTitlePanel updates the title panel with current mode
func (d *Dashboard) updateTitlePanel() {
	modeText := "Scheme Mode"
	if d.navigationMode == "tools" {
		modeText = "Tools Mode"
	}
	d.titlePanel.SetText(fmt.Sprintf("[yellow]GoCmder Setup[white] - [cyan]%s[white]", modeText))
}

// updateOptionsTable populates the left panel list based on navigation mode
func (d *Dashboard) updateOptionsTable() {
	currentRow, _ := d.optionsTable.GetSelection()

	// Clear existing content
	d.optionsTable.Clear()

	if d.navigationMode == "scheme" {
		// Display schemes
		for i, scheme := range d.config.Schemes {
			schemeIndicator := " "
			if i == d.config.CurrentScheme {
				schemeIndicator = "✓"
			}

			// Scheme name with indicator
			nameCell := tview.NewTableCell(fmt.Sprintf("%s %s", schemeIndicator, scheme.Name)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignLeft)
			d.optionsTable.SetCell(i, 0, nameCell)
		}
	} else {
		// Display tools (Tools mode)
		for i, tool := range d.config.Tools {
			checkbox := CheckboxOffASCII
			if tool.Selected {
				checkbox = CheckboxOnASCII
			}

			var checkboxColor tcell.Color
			var nameColor tcell.Color
			if tool.Selected {
				checkboxColor = tcell.ColorGreen
				nameColor = tcell.ColorGreen
			} else {
				checkboxColor = tcell.ColorYellow
				nameColor = tcell.ColorWhite
			}

			checkboxCell := tview.NewTableCell(checkbox).
				SetTextColor(checkboxColor).
				SetAlign(tview.AlignLeft)
			d.optionsTable.SetCell(i, 0, checkboxCell)

			nameCell := tview.NewTableCell(tool.Name).
				SetTextColor(nameColor).
				SetAlign(tview.AlignLeft)
			d.optionsTable.SetCell(i, 1, nameCell)
		}

		// Add settings
		for i, setting := range d.config.Settings {
			row := len(d.config.Tools) + i
			checkbox := CheckboxOffASCII
			if setting.Selected {
				checkbox = CheckboxOnASCII
			}

			var checkboxColor tcell.Color
			var nameColor tcell.Color
			if setting.Selected {
				checkboxColor = tcell.ColorGreen
				nameColor = tcell.ColorGreen
			} else {
				checkboxColor = tcell.ColorYellow
				nameColor = tcell.ColorWhite
			}

			checkboxCell := tview.NewTableCell(checkbox).
				SetTextColor(checkboxColor).
				SetAlign(tview.AlignLeft)
			d.optionsTable.SetCell(row, 0, checkboxCell)

			nameCell := tview.NewTableCell(setting.Name).
				SetTextColor(nameColor).
				SetAlign(tview.AlignLeft)
			d.optionsTable.SetCell(row, 1, nameCell)
		}
	}

	// Restore selection if valid, otherwise select first item
	if currentRow >= 0 && currentRow < d.optionsTable.GetRowCount() {
		d.optionsTable.Select(currentRow, 0)
	} else if d.optionsTable.GetRowCount() > 0 {
		d.optionsTable.Select(0, 0)
	}
}

// updateOutputPanel updates the output panel (status bar)
func (d *Dashboard) updateOutputPanel() {
	toolsCount := d.config.GetSelectedToolsCount()
	settingsCount := d.config.GetSelectedSettingsCount()

	statusText := fmt.Sprintf("[cyan][Y/n/?][white] Space:Select  Tab:Switch  F5:Run  F9:Output  Esc:Return  Ctrl+C:Exit  |  Selected: %d tools, %d settings",
		toolsCount, settingsCount)

	d.outputPanel.SetText(statusText)
}

// updateMessage updates the status message
func (d *Dashboard) updateMessage(msg string, msgType string) {
	d.message = msg
	d.messageType = msgType

	color := "gray"
	prefix := "[INFO]"
	switch msgType {
	case "success":
		color = "green"
		prefix = "[OK]"
	case "error":
		color = "red"
		prefix = "[ERROR]"
	case "info":
		color = "blue"
		prefix = "[INFO]"
	}

	if msg != "" {
		// Update output panel with formatted message
		currentText := d.outputPanel.GetText(false)
		d.outputPanel.SetText(fmt.Sprintf("[%s]%s %s[white]\n%s",
			color, prefix, msg, currentText))
	}
	d.updateOutputPanel()
}

// toggleDetailItem handles toggling items in the detail panel
func (d *Dashboard) toggleDetailItem() {
	// Find Custom scheme index
	customSchemeIdx := -1
	for i, scheme := range d.config.Schemes {
		if scheme.Name == "Custom" {
			customSchemeIdx = i
			break
		}
	}

	currentRow, _ := d.detailTable.GetSelection()
	total := len(d.config.Tools)

	// In scheme mode, account for headers (scheme name + "Tools:" header)
	// Row 0: scheme name, Row 1: "Tools:" header, Rows 2+: tools
	var toolRowOffset int
	var settingRowStart int

	if d.navigationMode == "scheme" {
		toolRowOffset = 2                           // Skip scheme name and "Tools:" header
		settingRowStart = toolRowOffset + total + 1 // +1 for "Settings:" header
	} else {
		toolRowOffset = 0
		settingRowStart = total
	}

	// If in scheme mode and not already in Custom, apply current scheme before making changes
	if d.navigationMode == "scheme" {
		currentSchemeRow, _ := d.optionsTable.GetSelection()
		if currentSchemeRow >= 0 && currentSchemeRow < len(d.config.Schemes) {
			currentScheme := d.config.Schemes[currentSchemeRow]
			if currentScheme.Name != "Custom" {
				// Apply current scheme to Selected state first
				d.config.ApplyScheme(currentSchemeRow)
			}
		}
	}

	// Check if selection is in tools section
	if d.navigationMode == "tools" || (currentRow >= toolRowOffset && currentRow < (toolRowOffset+total)) {
		// Toggle tool
		var toolIdx int
		if d.navigationMode == "scheme" {
			toolIdx = currentRow - toolRowOffset
		} else {
			toolIdx = currentRow
		}

		if toolIdx >= 0 && toolIdx < len(d.config.Tools) {
			d.config.Tools[toolIdx].Selected = !d.config.Tools[toolIdx].Selected
			tool := d.config.Tools[toolIdx]
			status := "selected"
			if !tool.Selected {
				status = "deselected"
			}
			d.updateMessage(fmt.Sprintf("Tool '%s' %s - Switched to Custom", tool.Name, status), "info")

			// Auto-switch to Custom scheme
			if customSchemeIdx >= 0 {
				d.config.CurrentScheme = customSchemeIdx
				d.optionsTable.Select(customSchemeIdx, 0)
			}
			d.updateDetailTable()
		}
	} else if d.navigationMode == "scheme" || (currentRow >= settingRowStart) {
		// Toggle setting
		var settingIdx int
		if d.navigationMode == "scheme" {
			settingIdx = currentRow - settingRowStart
		} else {
			settingIdx = currentRow - total
		}

		if settingIdx >= 0 && settingIdx < len(d.config.Settings) {
			d.config.Settings[settingIdx].Selected = !d.config.Settings[settingIdx].Selected
			setting := d.config.Settings[settingIdx]
			status := "enabled"
			if !setting.Selected {
				status = "disabled"
			}
			d.updateMessage(fmt.Sprintf("Setting '%s' %s - Switched to Custom", setting.Name, status), "info")

			// Auto-switch to Custom scheme
			if customSchemeIdx >= 0 {
				d.config.CurrentScheme = customSchemeIdx
				d.optionsTable.Select(customSchemeIdx, 0)
			}
			d.updateDetailTable()
		}
	}
}

// updateDetailTable updates the detail table with current scheme/tool details
func (d *Dashboard) updateDetailTable() {
	d.detailTable.Clear()

	if d.navigationMode == "scheme" {
		// In scheme mode, show all tools and settings for customization
		// Pre-check items included in the selected scheme
		currentRow, _ := d.optionsTable.GetSelection()
		if currentRow >= 0 && currentRow < len(d.config.Schemes) {
			scheme := d.config.Schemes[currentRow]

			// Add scheme name as header
			headerCell := tview.NewTableCell(fmt.Sprintf("[bold]%s[white]", scheme.Name)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignLeft)
			d.detailTable.SetCell(0, 0, headerCell)

			row := 1

			// Add tools header
			toolHeaderCell := tview.NewTableCell("[yellow]Tools:[white]").
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft)
			d.detailTable.SetCell(row, 0, toolHeaderCell)
			row++

			// Create a map for quick lookup of included tool indices
			includedTools := make(map[int]bool)
			for _, idx := range scheme.ToolIndices {
				includedTools[idx] = true
			}

			// Show all tools with checkboxes
			for i, tool := range d.config.Tools {
				checkbox := CheckboxOffASCII
				textColor := tcell.ColorYellow

				// For Custom scheme: show actual selected state
				// For other schemes: show items included in the scheme
				if scheme.Name == "Custom" {
					if tool.Selected {
						checkbox = CheckboxOnASCII
						textColor = tcell.ColorGreen
					}
				} else {
					if includedTools[i] {
						checkbox = CheckboxOnASCII
						textColor = tcell.ColorGreen
					}
				}

				toolCell := tview.NewTableCell(fmt.Sprintf("  %s %s", checkbox, tool.Name)).
					SetTextColor(textColor).
					SetAlign(tview.AlignLeft)
				d.detailTable.SetCell(row, 0, toolCell)
				row++
			}

			// Add settings header
			settingHeaderCell := tview.NewTableCell("[yellow]Settings:[white]").
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft)
			d.detailTable.SetCell(row, 0, settingHeaderCell)
			row++

			// Create a map for quick lookup of included setting indices
			includedSettings := make(map[int]bool)
			for _, idx := range scheme.SettingIndices {
				includedSettings[idx] = true
			}

			// Show all settings with checkboxes
			for i, setting := range d.config.Settings {
				checkbox := CheckboxOffASCII
				textColor := tcell.ColorYellow

				// For Custom scheme: show actual selected state
				// For other schemes: show items included in the scheme
				if scheme.Name == "Custom" {
					if setting.Selected {
						checkbox = CheckboxOnASCII
						textColor = tcell.ColorGreen
					}
				} else {
					if includedSettings[i] {
						checkbox = CheckboxOnASCII
						textColor = tcell.ColorGreen
					}
				}

				settingCell := tview.NewTableCell(fmt.Sprintf("  %s %s", checkbox, setting.Name)).
					SetTextColor(textColor).
					SetAlign(tview.AlignLeft)
				d.detailTable.SetCell(row, 0, settingCell)
				row++
			}
		}
		return
	}

	// In tools mode, show all tools and settings as interactive list
	// Add tools
	for i, tool := range d.config.Tools {
		checkbox := CheckboxOffASCII
		if tool.Selected {
			checkbox = CheckboxOnASCII
		}

		var checkboxColor tcell.Color
		var nameColor tcell.Color
		if tool.Selected {
			checkboxColor = tcell.ColorGreen
			nameColor = tcell.ColorGreen
		} else {
			checkboxColor = tcell.ColorYellow
			nameColor = tcell.ColorWhite
		}

		checkboxCell := tview.NewTableCell(checkbox).
			SetTextColor(checkboxColor).
			SetAlign(tview.AlignLeft)
		d.detailTable.SetCell(i, 0, checkboxCell)

		nameCell := tview.NewTableCell(tool.Name).
			SetTextColor(nameColor).
			SetAlign(tview.AlignLeft)
		d.detailTable.SetCell(i, 1, nameCell)
	}

	// Add settings
	for i, setting := range d.config.Settings {
		row := len(d.config.Tools) + i
		checkbox := CheckboxOffASCII
		if setting.Selected {
			checkbox = CheckboxOnASCII
		}

		var checkboxColor tcell.Color
		var nameColor tcell.Color
		if setting.Selected {
			checkboxColor = tcell.ColorGreen
			nameColor = tcell.ColorGreen
		} else {
			checkboxColor = tcell.ColorYellow
			nameColor = tcell.ColorWhite
		}

		checkboxCell := tview.NewTableCell(checkbox).
			SetTextColor(checkboxColor).
			SetAlign(tview.AlignLeft)
		d.detailTable.SetCell(row, 0, checkboxCell)

		nameCell := tview.NewTableCell(setting.Name).
			SetTextColor(nameColor).
			SetAlign(tview.AlignLeft)
		d.detailTable.SetCell(row, 1, nameCell)
	}
}

// GetRoot returns the root primitive
func (d *Dashboard) GetRoot() tview.Primitive {
	return d.mainLayout
}
