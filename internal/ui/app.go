package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/models"
	"github.com/shangyanjin/gocmder/internal/ui/database"
	"github.com/shangyanjin/gocmder/internal/ui/home"
	"github.com/shangyanjin/gocmder/internal/ui/settings"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/system"
	"github.com/shangyanjin/gocmder/internal/ui/terminal"
	"github.com/shangyanjin/gocmder/internal/ui/tools"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	homePageIndex = 0 + iota
	terminalPageIndex
	databasePageIndex
	toolsPageIndex
	settingsPageIndex
	systemPageIndex
	// F3 reserved for File Manager
	// F4 = Database (index 2)
	// F5 reserved for Editor
	// F6 = Tools (index 3)
	// F7 = Settings (index 4)
	// F8 = System (index 5)
	// F9 = Reserved
)

// UIPage represents a page in the UI
type UIPage interface {
	tview.Primitive
	GetTitle() string
	HasFocus() bool
	SubDialogHasFocus() bool
	HideAllDialogs()
	SetAppFocusHandler(handler func())
}

// App represents the main UI application
type App struct {
	app            *tview.Application
	pages          *tview.Pages
	layout         *tview.Flex
	mainLayout     *tview.Flex
	infoBar        *tview.TextView
	helpBar        *tview.TextView
	homePage       *home.Home
	terminalPage   *terminal.Terminal
	databasePage   *database.Database
	toolsPage      *tools.Tools
	settingsPage   *settings.Settings
	systemPage     *system.System
	currentPageIdx int
	pageList       []UIPage
	installHandler func(toolName string)
	applyHandler   func(settings []models.Setting)
}

// NewApp creates a new UI application
func NewApp(tviewApp *tview.Application) *App {
	uiApp := &App{
		app:            tviewApp,
		pages:          tview.NewPages(),
		currentPageIdx: homePageIndex,
	}

	// Create pages
	uiApp.homePage = home.NewHome()
	uiApp.terminalPage = terminal.NewTerminal()
	uiApp.databasePage = database.NewDatabase()
	uiApp.toolsPage = tools.NewTools()
	uiApp.settingsPage = settings.NewSettings()
	uiApp.systemPage = system.NewSystem()

	// Set page list
	uiApp.pageList = []UIPage{
		uiApp.homePage,
		uiApp.terminalPage,
		uiApp.databasePage,
		uiApp.toolsPage,
		uiApp.settingsPage,
		uiApp.systemPage,
	}

	// Set app focus handlers
	for _, page := range uiApp.pageList {
		page.SetAppFocusHandler(func() {
			uiApp.app.SetFocus(page)
		})
	}

	// Create info bar
	uiApp.infoBar = tview.NewTextView()
	uiApp.infoBar.SetBackgroundColor(style.InfoBarBgColor)
	uiApp.infoBar.SetTextColor(style.InfoBarFgColor)
	uiApp.infoBar.SetDynamicColors(true)
	uiApp.infoBar.SetText(" GoCmder - Developer Environment Setup Tool")

	// Create help bar
	uiApp.helpBar = tview.NewTextView()
	uiApp.helpBar.SetBackgroundColor(style.InfoBarBgColor)
	uiApp.helpBar.SetTextColor(style.InfoBarFgColor)
	uiApp.helpBar.SetDynamicColors(true)
	uiApp.updateHelpBar()

	// Add pages
	uiApp.pages.AddPage("home", uiApp.homePage, true, true)
	uiApp.pages.AddPage("terminal", uiApp.terminalPage, true, false)
	uiApp.pages.AddPage("database", uiApp.databasePage, true, false)
	uiApp.pages.AddPage("tools", uiApp.toolsPage, true, false)
	uiApp.pages.AddPage("settings", uiApp.settingsPage, true, false)
	uiApp.pages.AddPage("system", uiApp.systemPage, true, false)

	// Create main layout
	uiApp.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	uiApp.mainLayout.AddItem(uiApp.pages, 0, 1, true)

	// Create overall layout with info bars
	uiApp.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	uiApp.layout.AddItem(uiApp.infoBar, 1, 0, false)
	uiApp.layout.AddItem(uiApp.mainLayout, 0, 1, true)
	uiApp.layout.AddItem(uiApp.helpBar, 1, 0, false)

	// Set input capture for global key bindings
	uiApp.layout.SetInputCapture(uiApp.globalInputHandler)

	return uiApp
}

// GetRoot returns the root primitive
func (a *App) GetRoot() tview.Primitive {
	return a.layout
}

// SetInstallHandler sets the handler for tool installation
func (a *App) SetInstallHandler(handler func(toolName string)) {
	a.installHandler = handler
	a.toolsPage.SetInstallHandler(handler)
}

// SetApplySettingsHandler sets the handler for applying settings
func (a *App) SetApplySettingsHandler(handler func(settings []models.Setting)) {
	a.applyHandler = handler
	a.settingsPage.SetApplyHandler(handler)
}

// SetRefreshHandler sets the handler for refreshing tool list
func (a *App) SetRefreshHandler(handler func()) {
	a.toolsPage.SetRefreshHandler(handler)
}

// UpdateToolsData updates the tools page data
func (a *App) UpdateToolsData(toolsData []models.Tool) {
	a.toolsPage.UpdateData(toolsData)
}

// UpdateSettingsData updates the settings page data
func (a *App) UpdateSettingsData(settingsData []models.Setting) {
	a.settingsPage.UpdateData(settingsData)
}

// RefreshSystemInfo refreshes the system info page
func (a *App) RefreshSystemInfo() {
	a.systemPage.Refresh()
}

// globalInputHandler handles global key bindings
func (a *App) globalInputHandler(event *tcell.EventKey) *tcell.EventKey {
	// Check if any dialog is displayed
	if a.getCurrentPage().SubDialogHasFocus() {
		return event
	}

	// Handle ESC key - return to home page (except when already on home page)
	if event.Key() == utils.CloseDialogKey.Key {
		if a.currentPageIdx != homePageIndex {
			a.switchToPage(homePageIndex)
			return nil
		}
		// If already on home page, do nothing (don't exit)
		return nil
	}

	// Handle quit key
	if event.Rune() == utils.QuitKey.Rune {
		a.app.Stop()
		return nil
	}

	// Handle Tab key for page switching
	if event.Key() == utils.SwitchFocusKey.Key {
		a.switchToNextPage()
		return nil
	}

	// Handle Function key navigation (F1-F9)
	switch event.Key() {
	case tcell.KeyF1:
		a.switchToPage(homePageIndex)
		return nil
	case tcell.KeyF2:
		a.switchToPage(terminalPageIndex)
		return nil
	case tcell.KeyF3:
		// Reserved: File Manager
		return nil
	case tcell.KeyF4:
		a.switchToPage(databasePageIndex)
		return nil
	case tcell.KeyF5:
		// Reserved: Editor
		return nil
	case tcell.KeyF6:
		a.switchToPage(toolsPageIndex)
		return nil
	case tcell.KeyF7:
		a.switchToPage(settingsPageIndex)
		return nil
	case tcell.KeyF8:
		a.switchToPage(systemPageIndex)
		return nil
	case tcell.KeyF9:
		// Reserved for future
		return nil
	}

	return event
}

// switchToNextPage switches to the next page
func (a *App) switchToNextPage() {
	a.currentPageIdx = (a.currentPageIdx + 1) % len(a.pageList)
	a.updateCurrentPage()
}

// switchToPage switches to a specific page
func (a *App) switchToPage(pageIdx int) {
	if pageIdx >= 0 && pageIdx < len(a.pageList) {
		a.currentPageIdx = pageIdx
		a.updateCurrentPage()
	}
}

// updateCurrentPage updates the displayed page
func (a *App) updateCurrentPage() {
	// Hide all dialogs on current page
	for _, page := range a.pageList {
		page.HideAllDialogs()
	}

	// Switch page
	switch a.currentPageIdx {
	case homePageIndex:
		a.pages.SwitchToPage("home")
	case terminalPageIndex:
		a.pages.SwitchToPage("terminal")
	case databasePageIndex:
		a.pages.SwitchToPage("database")
	case toolsPageIndex:
		a.pages.SwitchToPage("tools")
	case settingsPageIndex:
		a.pages.SwitchToPage("settings")
	case systemPageIndex:
		a.pages.SwitchToPage("system")
	}

	// Update help bar
	a.updateHelpBar()

	// Set focus to current page
	a.app.SetFocus(a.getCurrentPage())
}

// getCurrentPage returns the current page
func (a *App) getCurrentPage() UIPage {
	if a.currentPageIdx >= 0 && a.currentPageIdx < len(a.pageList) {
		return a.pageList[a.currentPageIdx]
	}
	return a.pageList[0]
}

// updateHelpBar updates the help bar text
func (a *App) updateHelpBar() {
	highlightColor := style.GetColorHex(style.StatusInstalledColor)

	baseHelp := " [" + highlightColor + "]F1[-] Home | [" + highlightColor + "]F2[-] Term | [" + highlightColor + "]F4[-] DB | [" + highlightColor + "]F6[-] Tools | [" + highlightColor + "]F7[-] Settings | [" + highlightColor + "]Tab[-] Next | [" + highlightColor + "]q[-] Quit"

	var pageHelp string
	switch a.currentPageIdx {
	case homePageIndex:
		pageHelp = ""
	case terminalPageIndex:
		pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Enter[-] Execute"
	case databasePageIndex:
		pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Ctrl+N[-] Connect | [" + highlightColor + "]Ctrl+R[-] Execute | [" + highlightColor + "]Ctrl+←/→[-] Switch Panel"
	case toolsPageIndex:
		pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | [" + highlightColor + "]a[-] All | [" + highlightColor + "]i[-] Install"
	case settingsPageIndex:
		pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | [" + highlightColor + "]a[-] All | [" + highlightColor + "]Enter[-] Apply"
	case systemPageIndex:
		pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]↑/↓[-] Scroll"
	}

	a.helpBar.SetText(baseHelp + pageHelp)
}

// ShowError displays an error message
func (a *App) ShowError(title, message string) {
	// Show error on current page if it has error dialog capability
	// For now, we can add this to specific pages as needed
}

// ShowMessage displays a message
func (a *App) ShowMessage(title, message string) {
	// Show message on current page if it has message dialog capability
	// For now, we can add this to specific pages as needed
}
