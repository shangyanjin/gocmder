package utils

import "github.com/gdamore/tcell/v2"

// Key represents a keyboard key combination
type Key struct {
	Key  tcell.Key
	Rune rune
}

// Common key bindings used throughout the application
var (
	CloseDialogKey = Key{Key: tcell.KeyEscape}
	SwitchFocusKey = Key{Key: tcell.KeyTab}
	ConfirmKey     = Key{Key: tcell.KeyEnter}
	HelpKey        = Key{Rune: 'h'}
	QuitKey        = Key{Rune: 'q'}
	CommandKey     = Key{Rune: 'c'}
	RefreshKey     = Key{Rune: 'r'}
	InstallKey     = Key{Rune: 'i'}
	SelectAllKey   = Key{Rune: 'a'}
	ToggleKey      = Key{Key: tcell.KeyRune, Rune: ' '}
)
