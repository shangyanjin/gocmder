package dialogs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// ConfirmDialog is a confirmation dialog primitive
type ConfirmDialog struct {
	*tview.Box

	modal         *tview.Modal
	title         string
	message       string
	display       bool
	selectHandler func()
	cancelHandler func()
}

// NewConfirmDialog returns a new confirmation dialog primitive
func NewConfirmDialog() *ConfirmDialog {
	bgColor := style.DialogBgColor
	dialog := &ConfirmDialog{
		Box:     tview.NewBox(),
		modal:   tview.NewModal().SetBackgroundColor(bgColor).AddButtons([]string{"Yes", "No"}),
		display: false,
	}

	dialog.modal.SetButtonBackgroundColor(style.ButtonBgColor)
	dialog.modal.SetBorderStyle(tcell.StyleDefault.
		Background(bgColor).
		Foreground(style.DialogFgColor))

	dialog.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			if dialog.selectHandler != nil {
				dialog.selectHandler()
			}
		} else {
			if dialog.cancelHandler != nil {
				dialog.cancelHandler()
			}
		}
	})

	return dialog
}

// Display displays this primitive
func (d *ConfirmDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ConfirmDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ConfirmDialog) Hide() {
	d.title = ""
	d.message = ""
	d.display = false
}

// SetText sets confirmation dialog message
func (d *ConfirmDialog) SetText(message string) {
	d.message = message
}

// SetTitle sets confirmation dialog title
func (d *ConfirmDialog) SetTitle(title string) {
	d.title = title
}

// HasFocus returns whether or not this primitive has focus
func (d *ConfirmDialog) HasFocus() bool {
	return d.display && d.modal.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ConfirmDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.modal)
}

// InputHandler returns input handler function for this primitive
func (d *ConfirmDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if modalHandler := d.modal.InputHandler(); modalHandler != nil {
			modalHandler(event, setFocus)
			return
		}
	})
}

// SetRect sets rects for this primitive
func (d *ConfirmDialog) SetRect(x, y, width, height int) {
	d.Box.SetRect(x, y, width, height)
}

// GetRect returns the current position of the primitive
func (d *ConfirmDialog) GetRect() (int, int, int, int) {
	return d.Box.GetRect()
}

// Draw draws this primitive onto the screen
func (d *ConfirmDialog) Draw(screen tcell.Screen) {
	hFgColor := style.FgColor
	headerColor := style.GetColorHex(hFgColor)

	var confirmMessage string

	if d.title != "" {
		confirmMessage = fmt.Sprintf("[%s::b]%s[-::-]\n", headerColor, d.title)
	}

	confirmMessage += d.message
	d.modal.SetText(confirmMessage)
	d.modal.Draw(screen)
}

// SetSelectedFunc sets form enter button selected function
func (d *ConfirmDialog) SetSelectedFunc(handler func()) *ConfirmDialog {
	d.selectHandler = handler
	return d
}

// SetCancelFunc sets form cancel button selected function
func (d *ConfirmDialog) SetCancelFunc(handler func()) *ConfirmDialog {
	d.cancelHandler = handler
	return d
}
