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

	title          string
	systemInfoView *tview.TextView
	helpView       *tview.TextView
	flex           *tview.Flex
}

// NewHome returns home page view
func NewHome() *Home {
	home := &Home{
		Box:   tview.NewBox(),
		title: "home",
	}

	// Create system info panel
	home.systemInfoView = tview.NewTextView()
	home.systemInfoView.SetDynamicColors(true)
	home.systemInfoView.SetBackgroundColor(style.BgColor)
	home.systemInfoView.SetTextColor(style.FgColor)
	home.systemInfoView.SetBorder(true)
	home.systemInfoView.SetBorderColor(style.BorderColor)
	home.systemInfoView.SetTitle(" System Dashboard ")
	home.systemInfoView.SetScrollable(true)

	// Create help panel
	home.helpView = tview.NewTextView()
	home.helpView.SetDynamicColors(true)
	home.helpView.SetBackgroundColor(style.BgColor)
	home.helpView.SetTextColor(style.FgColor)
	home.helpView.SetBorder(true)
	home.helpView.SetBorderColor(style.BorderColor)
	home.helpView.SetTitle(" Help & Quick Reference ")
	home.helpView.SetScrollable(true)

	// Update initial content
	home.updateSystemInfo()
	home.updateHelp()

	// Create layout - left: system info, right: help
	home.flex = tview.NewFlex().SetDirection(tview.FlexColumn)
	home.flex.AddItem(home.systemInfoView, 0, 1, false)
	home.flex.AddItem(home.helpView, 0, 1, false)

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

// updateSystemInfo updates the system information display
func (h *Home) updateSystemInfo() {
	h.systemInfoView.Clear()

	headerColor := style.GetColorHex(style.PageHeaderFgColor)
	valueColor := style.GetColorHex(style.FgColor)
	highlightColor := style.GetColorHex(style.StatusInstalledColor)

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info := fmt.Sprintf(`
[%s::b]GoCmder - Developer Environment Setup Tool v0.2.0[-::-]
[%s]GitHub: https://github.com/shangyanjin/gocmder[-]

[%s::b]System Information:[-::-]
  OS:           [%s]%s[-]
  Architecture: [%s]%s[-]
  CPU Cores:    [%s]%d[-]
  Go Version:   [%s]%s[-]

[%s::b]Runtime Statistics:[-::-]
  Goroutines:   [%s]%d[-]
  Memory Alloc: [%s]%.2f MB[-]
  Total Alloc:  [%s]%.2f MB[-]
  Sys Memory:   [%s]%.2f MB[-]
  GC Runs:      [%s]%d[-]

[%s::b]Application Status:[-::-]
  Status:       [%s]Running[-]

[%s::b]Quick Actions:[-::-]
  • Press [%s]F2[-] to open Terminal
  • Press [%s]F4[-] to manage Databases
  • Press [%s]F6[-] to install Development Tools
  • Press [%s]F7[-] to configure System Settings
  • Press [%s]F8[-] to view detailed System Information
`,
		headerColor,
		valueColor,
		headerColor,
		valueColor, runtime.GOOS,
		valueColor, runtime.GOARCH,
		valueColor, runtime.NumCPU(),
		valueColor, runtime.Version(),
		headerColor,
		valueColor, runtime.NumGoroutine(),
		valueColor, float64(m.Alloc)/1024/1024,
		valueColor, float64(m.TotalAlloc)/1024/1024,
		valueColor, float64(m.Sys)/1024/1024,
		valueColor, m.NumGC,
		headerColor,
		highlightColor,
		headerColor,
		highlightColor, highlightColor, highlightColor, highlightColor, highlightColor,
	)

	fmt.Fprint(h.systemInfoView, info)
}

// updateHelp updates the help information display
func (h *Home) updateHelp() {
	h.helpView.Clear()

	headerColor := style.GetColorHex(style.PageHeaderFgColor)
	highlightColor := style.GetColorHex(style.StatusInstalledColor)

	help := fmt.Sprintf(`
[%s::b]Function Keys:[-::-]
  [%s]F1[-]     Home (this page)
  [%s]F2[-]     Terminal - Embedded terminal
  [%s]F3[-]     Reserved (File Manager)
  [%s]F4[-]     Database Manager
  [%s]F5[-]     Reserved (Editor)
  [%s]F6[-]     Development Tools
  [%s]F7[-]     System Settings
  [%s]F8[-]     System Information
  [%s]F9[-]     Reserved

[%s::b]Global Shortcuts:[-::-]
  [%s]Tab[-]       Cycle through pages
  [%s]ESC[-]       Return to Home page
  [%s]q[-]         Quit application

[%s::b]Database Manager (F4):[-::-]
  [%s]Ctrl+N[-]    New connection
  [%s]ALT+M[-]     MySQL preset
  [%s]ALT+P[-]     PostgreSQL preset
  [%s]ALT+L[-]     SQLite preset
  [%s]ALT+S[-]     Save connection
  [%s]ALT+C[-]     Connect & close

[%s::b]Development Tools (F6):[-::-]
  [%s]Space[-]     Toggle selection
  [%s]a[-]         Select all
  [%s]i[-]         Install selected
  [%s]r[-]         Refresh list

[%s::b]Settings (F7):[-::-]
  [%s]Space[-]     Toggle selection
  [%s]Enter[-]     Apply settings

[%s::b]About:[-::-]
  This tool helps you quickly rebuild
  your development environment after
  system reinstall or setup a new machine.
`,
		headerColor,
		highlightColor, highlightColor, highlightColor, highlightColor, highlightColor,
		highlightColor, highlightColor, highlightColor, highlightColor,
		headerColor,
		highlightColor, highlightColor, highlightColor,
		headerColor,
		highlightColor, highlightColor, highlightColor, highlightColor, highlightColor, highlightColor,
		headerColor,
		highlightColor, highlightColor, highlightColor, highlightColor,
		headerColor,
		highlightColor, highlightColor,
		headerColor,
	)

	fmt.Fprint(h.helpView, help)
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
