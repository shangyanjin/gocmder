package dialogs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
)

// MessageDialog is a message dialog primitive
type MessageDialog struct {
	*tview.Box

	modal         *tview.Modal
	title         string
	message       string
	display       bool
	cancelHandler func()
}

// NewMessageDialog returns a new message dialog primitive
func NewMessageDialog(title string) *MessageDialog {
	bgColor := style.DialogBgColor
	dialog := &MessageDialog{
		Box:     tview.NewBox(),
		modal:   tview.NewModal().SetBackgroundColor(bgColor).AddButtons([]string{"OK"}),
		title:   title,
		display: false,
	}

	dialog.modal.SetButtonBackgroundColor(style.ButtonBgColor)
	dialog.modal.SetBorderStyle(tcell.StyleDefault.
		Background(bgColor).
		Foreground(style.DialogFgColor))

	dialog.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if dialog.cancelHandler != nil {
			dialog.cancelHandler()
		}
	})

	return dialog
}

// Display displays this primitive
func (d *MessageDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *MessageDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *MessageDialog) Hide() {
	d.message = ""
	d.display = false
}

// SetText sets message dialog text
func (d *MessageDialog) SetText(message string) {
	d.message = message
}

// SetTitle sets message dialog title
func (d *MessageDialog) SetTitle(title string) {
	d.title = title
}

// HasFocus returns whether or not this primitive has focus
func (d *MessageDialog) HasFocus() bool {
	return d.display && d.modal.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *MessageDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.modal)
}

// InputHandler returns input handler function for this primitive
func (d *MessageDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if modalHandler := d.modal.InputHandler(); modalHandler != nil {
			modalHandler(event, setFocus)
			return
		}
	})
}

// SetRect sets rects for this primitive
func (d *MessageDialog) SetRect(x, y, width, height int) {
	d.Box.SetRect(x, y, width, height)
}

// GetRect returns the current position of the primitive
func (d *MessageDialog) GetRect() (int, int, int, int) {
	return d.Box.GetRect()
}

// Draw draws this primitive onto the screen
func (d *MessageDialog) Draw(screen tcell.Screen) {
	hFgColor := style.FgColor
	headerColor := style.GetColorHex(hFgColor)

	var messageText string

	if d.title != "" {
		messageText = fmt.Sprintf("[%s::b]%s[-::-]\n", headerColor, d.title)
	}

	messageText += d.message
	d.modal.SetText(messageText)
	d.modal.Draw(screen)
}

// SetCancelFunc sets modal cancel function
func (d *MessageDialog) SetCancelFunc(handler func()) *MessageDialog {
	d.cancelHandler = handler
	return d
}
