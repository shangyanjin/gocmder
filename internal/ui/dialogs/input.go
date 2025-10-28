package dialogs

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	inputDialogWidth  = 60
	inputDialogHeight = 10
)

const (
	inputFieldFocus = 0 + iota
	inputFormFocus
)

// SimpleInputDialog is a simple input dialog primitive
type SimpleInputDialog struct {
	*tview.Box

	layout        *tview.Flex
	form          *tview.Form
	inputField    *tview.InputField
	title         string
	display       bool
	focusElement  int
	selectHandler func()
	cancelHandler func()
}

// NewSimpleInputDialog returns a new simple input dialog primitive
func NewSimpleInputDialog(title string) *SimpleInputDialog {
	bgColor := style.DialogBgColor

	inputField := tview.NewInputField()
	inputField.SetBackgroundColor(bgColor)
	inputField.SetFieldBackgroundColor(style.BgColor)
	inputField.SetLabel("Input: ")
	inputField.SetLabelColor(style.DialogFgColor)
	inputField.SetFieldTextColor(style.FgColor)

	form := tview.NewForm().
		AddButton("OK", nil).
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)

	form.SetBackgroundColor(bgColor)
	form.SetButtonBackgroundColor(style.ButtonBgColor)

	inputLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	inputLayout.AddItem(inputField, 0, 1, true)
	inputLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	inputLayout.SetBackgroundColor(bgColor)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(inputLayout, 1, 0, true)
	layout.AddItem(form, DialogFormHeight, 0, true)
	layout.SetBorder(true)
	layout.SetTitle(" " + title + " ")
	layout.SetBorderColor(style.DialogBorderColor)
	layout.SetBackgroundColor(bgColor)

	dialog := &SimpleInputDialog{
		Box:          tview.NewBox(),
		layout:       layout,
		form:         form,
		inputField:   inputField,
		title:        title,
		display:      false,
		focusElement: inputFieldFocus,
	}

	return dialog
}

// Display displays this primitive
func (d *SimpleInputDialog) Display() {
	d.display = true
	d.focusElement = inputFieldFocus
}

// IsDisplay returns true if primitive is shown
func (d *SimpleInputDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *SimpleInputDialog) Hide() {
	d.display = false
	d.inputField.SetText("")
	d.focusElement = inputFieldFocus
}

// SetTitle sets input dialog title
func (d *SimpleInputDialog) SetTitle(title string) {
	d.title = title
	d.layout.SetTitle(" " + title + " ")
}

// SetLabel sets input field label
func (d *SimpleInputDialog) SetLabel(label string) {
	d.inputField.SetLabel(label)
}

// GetText returns the input field text
func (d *SimpleInputDialog) GetText() string {
	return d.inputField.GetText()
}

// SetText sets the input field text
func (d *SimpleInputDialog) SetText(text string) {
	d.inputField.SetText(text)
}

// HasFocus returns whether or not this primitive has focus
func (d *SimpleInputDialog) HasFocus() bool {
	return d.display && (d.inputField.HasFocus() || d.form.HasFocus())
}

// Focus is called when this primitive receives focus
func (d *SimpleInputDialog) Focus(delegate func(p tview.Primitive)) {
	if d.focusElement == inputFieldFocus {
		delegate(d.inputField)
		return
	}

	delegate(d.form)
}

// InputHandler returns input handler function for this primitive
func (d *SimpleInputDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == utils.CloseDialogKey.Key {
			if d.cancelHandler != nil {
				d.cancelHandler()
			}
			return
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
			if d.focusElement == inputFormFocus {
				setFocus(d.form)
			} else {
				setFocus(d.inputField)
			}
			return
		}

		if d.inputField.HasFocus() {
			if event.Key() == tcell.KeyEnter {
				if d.selectHandler != nil {
					d.selectHandler()
				}
				return
			}

			if inputHandler := d.inputField.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)
				return
			}
		}

		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
	})
}

// SetSelectedFunc sets form OK button selected function
func (d *SimpleInputDialog) SetSelectedFunc(handler func()) *SimpleInputDialog {
	d.selectHandler = handler
	okButton := d.form.GetButton(0)
	okButton.SetSelectedFunc(handler)
	return d
}

// SetCancelFunc sets form Cancel button selected function
func (d *SimpleInputDialog) SetCancelFunc(handler func()) *SimpleInputDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(1)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetRect sets rects for this primitive
func (d *SimpleInputDialog) SetRect(x, y, width, height int) {
	ws := (width - inputDialogWidth) / 2
	hs := (height - inputDialogHeight) / 2
	dy := y + hs
	bWidth := inputDialogWidth
	bHeight := inputDialogHeight

	if inputDialogWidth > width {
		ws = 0
		bWidth = width - 1
	}

	if inputDialogHeight >= height {
		dy = y + 1
		bHeight = height - 1
	}

	d.Box.SetRect(x+ws, dy, bWidth, bHeight)

	x, y, width, height = d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen
func (d *SimpleInputDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	d.layout.Draw(screen)
}

func (d *SimpleInputDialog) setFocusElement() {
	if d.focusElement == inputFieldFocus {
		d.focusElement = inputFormFocus
	} else {
		d.focusElement = inputFieldFocus
	}
}
