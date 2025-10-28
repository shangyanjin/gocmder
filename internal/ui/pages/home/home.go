package home

import (
	"fmt"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// Home implements the home page primitive
type Home struct {
	*tview.Box

	title    string
	textView *tview.TextView
	flex     *tview.Flex
}

// NewHome returns home page view
func NewHome() *Home {
	home := &Home{
		Box:   tview.NewBox(),
		title: "home",
	}

	// Create text view for welcome content
	home.textView = tview.NewTextView()
	home.textView.SetDynamicColors(true)
	home.textView.SetBackgroundColor(style.BgColor)
	home.textView.SetTextColor(style.FgColor)
	home.textView.SetBorder(false)

	// Create welcome panel
	welcomePanel := tview.NewTextView()
	welcomePanel.SetDynamicColors(true)
	welcomePanel.SetBackgroundColor(style.BgColor)
	welcomePanel.SetTextColor(style.FgColor)
	welcomePanel.SetBorder(true)
	welcomePanel.SetBorderColor(style.BorderColor)
	welcomePanel.SetTitle(" Welcome to GoCmder ")

	// Create quick actions panel
	actionsPanel := tview.NewTextView()
	actionsPanel.SetDynamicColors(true)
	actionsPanel.SetBackgroundColor(style.BgColor)
	actionsPanel.SetTextColor(style.FgColor)
	actionsPanel.SetBorder(true)
	actionsPanel.SetBorderColor(style.BorderColor)
	actionsPanel.SetTitle(" Quick Actions ")

	// Create keyboard shortcuts panel
	shortcutsPanel := tview.NewTextView()
	shortcutsPanel.SetDynamicColors(true)
	shortcutsPanel.SetBackgroundColor(style.BgColor)
	shortcutsPanel.SetTextColor(style.FgColor)
	shortcutsPanel.SetBorder(true)
	shortcutsPanel.SetBorderColor(style.BorderColor)
	shortcutsPanel.SetTitle(" Function Keys ")

	// Fill content
	highlightColor := style.GetColorHex(style.StatusInstalledColor)
	titleColor := style.GetColorHex(style.PageHeaderFgColor)

	welcomeText := fmt.Sprintf(`
[%s::b]GoCmder - Developer Environment Setup Tool[-::-]

Platform: [%s]%s/%s[-]
Go Version: [%s]%s[-]

This tool helps you quickly rebuild your development 
environment after system reinstall or setup a new machine.

Navigate through different pages using function keys (F1-F9)
or use the Tab key to cycle through pages.
`,
		titleColor,
		highlightColor, runtime.GOOS, runtime.GOARCH,
		highlightColor, runtime.Version(),
	)

	actionsText := fmt.Sprintf(`
[%s]F2[-] Terminal
    Embedded terminal for command execution

[%s]F4[-] Database Manager
    Connect to MySQL/PostgreSQL, execute SQL,
    browse tables and view results

[%s]F6[-] Development Tools
    Install Git, VSCode, Go, Node.js, PostgreSQL,
    MySQL, Redis and other development tools

[%s]F7[-] System Settings
    Configure PATH, power settings, user directories

[%s]F8[-] System Information
    View system information and capabilities
`,
		highlightColor, highlightColor, highlightColor, highlightColor, highlightColor,
	)

	shortcutsText := fmt.Sprintf(`
[%s]F1[-]     Home (this page)
[%s]F2[-]     Terminal
[%s]F3[-]     Reserved (File Manager)
[%s]F4[-]     Database Manager
[%s]F5[-]     Reserved (Editor)
[%s]F6[-]     Development Tools
[%s]F7[-]     System Settings
[%s]F8[-]     System Information
[%s]F9[-]     Reserved

[%s]Tab[-]    Cycle through pages
[%s]q[-]      Quit application
`,
		highlightColor, highlightColor, highlightColor, highlightColor,
		highlightColor, highlightColor, highlightColor, highlightColor,
		highlightColor, highlightColor, highlightColor,
	)

	fmt.Fprint(welcomePanel, welcomeText)
	fmt.Fprint(actionsPanel, actionsText)
	fmt.Fprint(shortcutsPanel, shortcutsText)

	// Create layout
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPanel.AddItem(welcomePanel, 0, 1, false)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPanel.AddItem(actionsPanel, 0, 1, false)
	rightPanel.AddItem(shortcutsPanel, 13, 0, false)

	home.flex = tview.NewFlex().SetDirection(tview.FlexColumn)
	home.flex.AddItem(leftPanel, 0, 1, false)
	home.flex.AddItem(rightPanel, 0, 1, false)

	return home
}

// GetTitle returns primitive title
func (h *Home) GetTitle() string {
	return h.title
}

// HasFocus returns whether or not this primitive has focus
func (h *Home) HasFocus() bool {
	return h.flex.HasFocus() || h.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (h *Home) Focus(delegate func(p tview.Primitive)) {
	delegate(h.flex)
}

// SetAppFocusHandler sets application focus handler
func (h *Home) SetAppFocusHandler(handler func()) {
	// Home page doesn't need app focus handler
}

// HideAllDialogs hides all sub dialogs (none for home page)
func (h *Home) HideAllDialogs() {
	// No dialogs in home page
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (h *Home) SubDialogHasFocus() bool {
	return false
}

// InputHandler returns the input handler for this primitive
func (h *Home) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return h.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if h.flex.HasFocus() {
			if flexHandler := h.flex.InputHandler(); flexHandler != nil {
				flexHandler(event, setFocus)
				return
			}
		}
	})
}

// Draw draws this primitive onto the screen
func (h *Home) Draw(screen tcell.Screen) {
	h.Box.DrawForSubclass(screen, h)
	x, y, width, height := h.GetInnerRect()

	h.flex.SetRect(x, y, width, height)
	h.flex.Draw(screen)
}
