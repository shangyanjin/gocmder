package ui

import (
	"github.com/gdamore/tcell/v2"
)

// KeyMap defines all keyboard shortcuts and their handlers
type KeyMap struct {
	handlers map[tcell.Key]func()
}

// NewKeyMap creates a new keymap instance
func NewKeyMap() *KeyMap {
	km := &KeyMap{
		handlers: make(map[tcell.Key]func()),
	}
	km.setupDefaultKeys()
	return km
}

// setupDefaultKeys sets up default keyboard shortcuts
func (km *KeyMap) setupDefaultKeys() {
	// Navigation keys
	km.handlers[tcell.KeyUp] = func() {
		// Handle up arrow
	}

	km.handlers[tcell.KeyDown] = func() {
		// Handle down arrow
	}

	km.handlers[tcell.KeyLeft] = func() {
		// Handle left arrow
	}

	km.handlers[tcell.KeyRight] = func() {
		// Handle right arrow
	}

	// Function keys
	km.handlers[tcell.KeyF1] = func() {
		// Switch panels
	}

	km.handlers[tcell.KeyF2] = func() {
		// Focus search
	}

	km.handlers[tcell.KeyF3] = func() {
		// Focus keys list
	}

	km.handlers[tcell.KeyF4] = func() {
		// Focus commands
	}

	km.handlers[tcell.KeyF5] = func() {
		// Focus results
	}

	km.handlers[tcell.KeyF9] = func() {
		// Focus output
	}

	// Control keys
	km.handlers[tcell.KeyTab] = func() {
		// Switch focus
	}

	km.handlers[tcell.KeyEnter] = func() {
		// Execute action
	}

	km.handlers[tcell.KeyEsc] = func() {
		// Quit application
	}

	km.handlers[tcell.KeyCtrlQ] = func() {
		// Quit application
	}

	km.handlers[tcell.KeyCtrlC] = func() {
		// Quit application
	}

	// Number keys for quick actions
	km.handlers[tcell.KeyRune] = func() {
		// Handle rune keys (1, 2, 3, etc.)
	}
}

// HandleKey processes a key press
func (km *KeyMap) HandleKey(key tcell.Key, runeKey rune) bool {
	// Handle special rune keys
	if key == tcell.KeyRune {
		switch runeKey {
		case '1':
			// Select all tools
			return true
		case '2':
			// Clear all tools
			return true
		case '3':
			// Enable all settings
			return true
		case 'q':
			// Quit - this should be handled by the application
			return false // Let the app handle quit
		case 'j':
			// Down (vim-style) - delegate to actual handler
			if handler, exists := km.handlers[tcell.KeyDown]; exists {
				handler()
			}
			return true
		case 'k':
			// Up (vim-style) - delegate to actual handler
			if handler, exists := km.handlers[tcell.KeyUp]; exists {
				handler()
			}
			return true
		case 'h':
			// Left (vim-style) - delegate to actual handler
			if handler, exists := km.handlers[tcell.KeyLeft]; exists {
				handler()
			}
			return true
		case 'l':
			// Right (vim-style) - delegate to actual handler
			if handler, exists := km.handlers[tcell.KeyRight]; exists {
				handler()
			}
			return true
		default:
			// Unhandled rune keys
			return false
		}
	}

	// Handle regular keys - only if handler exists and is not empty
	if handler, exists := km.handlers[key]; exists {
		// Check if this is a critical key that should not be consumed
		if key == tcell.KeyEsc || key == tcell.KeyCtrlQ || key == tcell.KeyCtrlC {
			return false // Let the app handle quit keys
		}
		handler()
		return true
	}

	return false
}

// SetHandler sets a custom handler for a key
func (km *KeyMap) SetHandler(key tcell.Key, handler func()) {
	km.handlers[key] = handler
}

// GetHelpText returns help text for all available keys
func (km *KeyMap) GetHelpText() string {
	return `Keyboard Shortcuts:
Navigation:
  ↑/↓/←/→  - Navigate
  j/k/h/l   - Vim-style navigation
  Tab       - Switch panels
  
Function Keys:
  F1        - Switch panels
  F2        - Focus search
  F3        - Focus keys list
  F4        - Focus commands
  F5        - Focus results
  F9        - Focus output
  
Actions:
  Enter     - Execute action
  Space     - Toggle selection
  1         - Select all tools
  2         - Clear all tools
  3         - Enable all settings
  
Quit:
  Esc       - Quit
  Ctrl+Q    - Quit
  Ctrl+C    - Quit
  q         - Quit`
}
