package dialogs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/ui/style"
	"github.com/shangyanjin/gocmder/internal/ui/utils"
)

const (
	progressDialogWidth  = 70
	progressDialogHeight = 10
)

// ProgressDialog is a progress dialog primitive
type ProgressDialog struct {
	*tview.Box

	layout      *tview.Flex
	textView    *tview.TextView
	progressBar *tview.TextView
	title       string
	display     bool
	progress    int
	maxProgress int
}

// NewProgressDialog returns a new progress dialog primitive
func NewProgressDialog() *ProgressDialog {
	bgColor := style.DialogBgColor

	textView := tview.NewTextView()
	textView.SetBackgroundColor(bgColor)
	textView.SetTextColor(style.DialogFgColor)
	textView.SetDynamicColors(true)

	progressBar := tview.NewTextView()
	progressBar.SetBackgroundColor(bgColor)
	progressBar.SetTextColor(style.ProgressBarFgColor)
	progressBar.SetDynamicColors(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(textView, 3, 0, false)
	layout.AddItem(progressBar, 1, 0, false)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.SetBorder(true)
	layout.SetBorderColor(style.DialogBorderColor)
	layout.SetBackgroundColor(bgColor)

	return &ProgressDialog{
		Box:         tview.NewBox(),
		layout:      layout,
		textView:    textView,
		progressBar: progressBar,
		title:       "Progress",
		display:     false,
		maxProgress: 100,
	}
}

// Display displays this primitive
func (d *ProgressDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ProgressDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ProgressDialog) Hide() {
	d.display = false
	d.progress = 0
	d.textView.Clear()
	d.progressBar.Clear()
}

// SetTitle sets progress dialog title
func (d *ProgressDialog) SetTitle(title string) {
	d.title = title
	d.layout.SetTitle(fmt.Sprintf(" %s ", title))
}

// SetText sets progress dialog text
func (d *ProgressDialog) SetText(text string) {
	d.textView.Clear()
	fmt.Fprintf(d.textView, "%s", text)
}

// SetProgress sets the current progress value
func (d *ProgressDialog) SetProgress(current, max int) {
	d.progress = current
	d.maxProgress = max
	d.updateProgressBar()
}

// updateProgressBar updates the progress bar display
func (d *ProgressDialog) updateProgressBar() {
	d.progressBar.Clear()

	if d.maxProgress <= 0 {
		return
	}

	percent := (d.progress * 100) / d.maxProgress
	barWidth := 50
	filled := (percent * barWidth) / 100

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += fmt.Sprintf("] %d%%", percent)

	fmt.Fprintf(d.progressBar, "[%s]%s", style.GetColorHex(style.ProgressBarFgColor), bar)
}

// HasFocus returns whether or not this primitive has focus
func (d *ProgressDialog) HasFocus() bool {
	return d.display && d.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ProgressDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.Box)
}

// SetRect sets rects for this primitive
func (d *ProgressDialog) SetRect(x, y, width, height int) {
	ws := (width - progressDialogWidth) / 2
	hs := (height - progressDialogHeight) / 2
	dy := y + hs
	bWidth := progressDialogWidth
	bHeight := progressDialogHeight

	if progressDialogWidth > width {
		ws = 0
		bWidth = width - 1
	}

	if progressDialogHeight >= height {
		dy = y + 1
		bHeight = height - 1
	}

	d.Box.SetRect(x+ws, dy, bWidth, bHeight)

	x, y, width, height = d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen
func (d *ProgressDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	d.layout.Draw(screen)
}
