package ui

import (
	"os"
	"runtime"

	"github.com/mattn/go-runewidth"
)

// Checkbox symbols with safe fallbacks.
var (
	CheckboxOnASCII  = "(x)" // ASCII compatible
	CheckboxOffASCII = "(o)" // ASCII compatible
)

// widthOf returns visual width in terminal cells.
func widthOf(s string) int {
	// Make ambiguous width treated as 1 cell to avoid double-width in CJK terminals.
	runewidth.DefaultCondition.EastAsianWidth = false
	return runewidth.StringWidth(s)
}

func preferASCII() bool {
	// Allow manual override: export GCMDER_ASCII=1
	if os.Getenv("GCMDER_ASCII") == "1" {
		return true
	}
	// On legacy Windows conhost (non Windows Terminal), fallback to ASCII.
	if runtime.GOOS == "windows" && os.Getenv("WT_SESSION") == "" {
		return true
	}
	return false
}

func init() {
	// Use ASCII compatible symbols for maximum terminal compatibility
	// ASCII symbols (escaped for tview)
	CheckboxOnASCII = "(*)"
	CheckboxOffASCII = "( )"
}
