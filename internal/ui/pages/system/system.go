package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// System implements the system info page primitive
type System struct {
	*tview.Box

	title     string
	textView  *tview.TextView
	osName    string
	osVersion string
	arch      string
	cpuCount  int
	goVersion string
}

// NewSystem returns system info page view
func NewSystem() *System {
	system := &System{
		Box:       tview.NewBox(),
		title:     "system information",
		osName:    runtime.GOOS,
		osVersion: "Unknown",
		arch:      runtime.GOARCH,
		cpuCount:  runtime.NumCPU(),
		goVersion: runtime.Version(),
	}

	system.textView = tview.NewTextView()
	system.textView.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(system.title)))
	system.textView.SetBorderColor(style.BorderColor)
	system.textView.SetTitleColor(style.FgColor)
	system.textView.SetBackgroundColor(style.BgColor)
	system.textView.SetTextColor(style.FgColor)
	system.textView.SetBorder(true)
	system.textView.SetDynamicColors(true)
	system.textView.SetScrollable(true)

	system.updateInfo()

	return system
}

// GetTitle returns primitive title
func (s *System) GetTitle() string {
	return s.title
}

// HasFocus returns whether or not this primitive has focus
func (s *System) HasFocus() bool {
	return s.textView.HasFocus() || s.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (s *System) Focus(delegate func(p tview.Primitive)) {
	delegate(s.textView)
}

// SetAppFocusHandler sets application focus handler
func (s *System) SetAppFocusHandler(handler func()) {
	// System page doesn't need app focus handler
}

// HideAllDialogs hides all sub dialogs (none for system page)
func (s *System) HideAllDialogs() {
	// No dialogs in system page
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (s *System) SubDialogHasFocus() bool {
	return false
}

// updateInfo updates the system information display
func (s *System) updateInfo() {
	s.textView.Clear()

	headerColor := style.GetColorHex(style.PageHeaderFgColor)
	valueColor := style.GetColorHex(style.FgColor)
	highlightColor := style.GetColorHex(style.StatusInstalledColor)

	info := fmt.Sprintf(`
[%s::b]Operating System:[-::-]
  OS:           [%s]%s[-]
  Architecture: [%s]%s[-]

[%s::b]Hardware:[-::-]
  CPU Cores:    [%s]%d[-]

[%s::b]Runtime:[-::-]
  Go Version:   [%s]%s[-]

[%s::b]Application:[-::-]
  Name:         [%s]GoCmder[-]
  Description:  Developer environment setup tool

[%s::b]Capabilities:[-::-]
  • Install development tools (Git, VSCode, Go, Node.js, etc.)
  • Configure system settings
  • Manage PATH environment
  • Setup user directories
  • Power configuration

[%s::b]Keyboard Shortcuts:[-::-]
  [%s]Space[-]     Toggle selection
  [%s]a[-]         Select all
  [%s]i[-]         Install selected (Tools page)
  [%s]Enter[-]     Apply settings (Settings page)
  [%s]Tab[-]       Switch between pages
  [%s]q[-]         Quit application
  [%s]h[-]         Show help

`,
		headerColor, valueColor, s.osName, valueColor, s.arch,
		headerColor, valueColor, s.cpuCount,
		headerColor, valueColor, s.goVersion,
		headerColor, highlightColor,
		headerColor,
		headerColor, highlightColor, highlightColor, highlightColor,
		highlightColor, highlightColor, highlightColor, highlightColor,
	)

	fmt.Fprint(s.textView, info)
}

// Refresh refreshes the system information
func (s *System) Refresh() {
	s.cpuCount = runtime.NumCPU()
	s.goVersion = runtime.Version()
	s.updateInfo()
}

// InputHandler returns the input handler for this primitive
func (s *System) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if s.textView.HasFocus() {
			if textViewHandler := s.textView.InputHandler(); textViewHandler != nil {
				textViewHandler(event, setFocus)
				return
			}
		}
	})
}

// Draw draws this primitive onto the screen
func (s *System) Draw(screen tcell.Screen) {
	s.Box.DrawForSubclass(screen, s)
	x, y, width, height := s.GetInnerRect()

	s.textView.SetRect(x, y, width, height)
	s.textView.Draw(screen)
}
