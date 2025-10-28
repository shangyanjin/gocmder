package terminal

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/pages/terminal/vterm"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// Terminal implements the terminal page primitive
type Terminal struct {
	*tview.Box

	title       string
	vtermDialog *vterm.VtermDialog
}

// NewTerminal returns terminal page view
func NewTerminal() *Terminal {
	terminal := &Terminal{
		Box:   tview.NewBox(),
		title: "terminal",
	}

	terminal.vtermDialog = vterm.NewVtermDialog()
	terminal.vtermDialog.SetCancelFunc(func() {
		// Terminal page doesn't close dialog, it's always visible
	})

	return terminal
}

// GetTitle returns primitive title
func (t *Terminal) GetTitle() string {
	return t.title
}

// HasFocus returns whether or not this primitive has focus
func (t *Terminal) HasFocus() bool {
	return t.vtermDialog.HasFocus() || t.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (t *Terminal) Focus(delegate func(p tview.Primitive)) {
	delegate(t.vtermDialog)
}

// SetAppFocusHandler sets application focus handler
func (t *Terminal) SetAppFocusHandler(handler func()) {
	// Terminal page doesn't need app focus handler
}

// HideAllDialogs hides all sub dialogs (none for terminal page)
func (t *Terminal) HideAllDialogs() {
	// Terminal is always visible when on this page
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (t *Terminal) SubDialogHasFocus() bool {
	return false
}

// InputHandler returns the input handler for this primitive
func (t *Terminal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if t.vtermDialog.HasFocus() {
			if vtermHandler := t.vtermDialog.InputHandler(); vtermHandler != nil {
				vtermHandler(event, setFocus)
				return
			}
		}
	})
}

// Draw draws this primitive onto the screen
func (t *Terminal) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)
	x, y, width, height := t.GetInnerRect()

	// Add title bar
	titleBar := tview.NewTextView()
	titleBar.SetBackgroundColor(style.PageHeaderBgColor)
	titleBar.SetTextColor(style.PageHeaderFgColor)
	titleBar.SetText(fmt.Sprintf(" [::b]%s[0] ", strings.ToUpper(t.title)))
	titleBar.SetRect(x, y, width, 1)
	titleBar.Draw(screen)

	// Draw vterm dialog (takes remaining space)
	t.vtermDialog.Display()
	t.vtermDialog.SetRect(x, y+1, width, height-1)
	t.vtermDialog.Draw(screen)
}
