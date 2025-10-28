package style

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Color definitions for the UI - GNOME style
var (
	// Background colors - GNOME dark theme inspired
	BgColor                  = tcell.NewRGBColor(46, 52, 54)   // #2e3436 GNOME dark base
	DialogBgColor            = tcell.NewRGBColor(32, 36, 38)   // #202426 darker dialog
	ErrorDialogBgColor       = tcell.NewRGBColor(164, 0, 0)    // #a40000 GNOME red
	PageHeaderBgColor        = tcell.NewRGBColor(52, 101, 164) // #3465a4 GNOME blue
	TableHeaderBgColor       = tcell.NewRGBColor(52, 101, 164) // #3465a4 GNOME blue
	ButtonBgColor            = tcell.NewRGBColor(52, 101, 164) // #3465a4 GNOME blue
	InfoBarBgColor           = tcell.NewRGBColor(36, 41, 43)   // #24292b dark bar
	ProgressBarBgColor       = tcell.NewRGBColor(85, 87, 83)   // #555753 dark gray
	ErrorDialogButtonBgColor = tcell.NewRGBColor(204, 0, 0)    // #cc0000 bright red

	// Foreground colors
	FgColor            = tcell.NewRGBColor(211, 215, 207) // #d3d7cf GNOME light text
	DialogFgColor      = tcell.NewRGBColor(238, 238, 236) // #eeeeec bright text
	PageHeaderFgColor  = tcell.NewRGBColor(238, 238, 236) // #eeeeec header text
	TableHeaderFgColor = tcell.NewRGBColor(238, 238, 236) // #eeeeec header text
	InfoBarFgColor     = tcell.NewRGBColor(211, 215, 207) // #d3d7cf info text
	ProgressBarFgColor = tcell.NewRGBColor(115, 210, 22)  // #73d216 GNOME green

	// Border colors
	BorderColor       = tcell.NewRGBColor(85, 87, 83)    // #555753 GNOME border
	DialogBorderColor = tcell.NewRGBColor(136, 138, 133) // #888a85 lighter border

	// Status colors - GNOME palette
	StatusInstalledColor    = tcell.NewRGBColor(115, 210, 22)  // #73d216 green
	StatusNotInstalledColor = tcell.NewRGBColor(252, 175, 62)  // #fcaf3e orange
	StatusErrorColor        = tcell.NewRGBColor(239, 41, 41)   // #ef2929 red
	StatusSelectedColor     = tcell.NewRGBColor(114, 159, 207) // #729fcf blue
)

// GetColorHex returns the hex string for a color
func GetColorHex(color tcell.Color) string {
	// For RGB colors, extract the hex value
	r, g, b := color.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
