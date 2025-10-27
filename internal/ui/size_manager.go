package ui

import (
	"fmt"
	"time"
)

// SizeManager manages dynamic sizing and layout calculations
type SizeManager struct {
	width       int
	height      int
	minWidth    int
	minHeight   int
	lastUpdate  time.Time
	resizeCount int
}

// NewSizeManager creates a new size manager
func NewSizeManager() *SizeManager {
	return &SizeManager{
		width:       80,
		height:      24,
		minWidth:    60,
		minHeight:   20,
		lastUpdate:  time.Now(),
		resizeCount: 0,
	}
}

// UpdateSize updates the current size
func (sm *SizeManager) UpdateSize(width, height int) bool {
	if sm.width == width && sm.height == height {
		return false // No change
	}

	sm.width = width
	sm.height = height
	sm.lastUpdate = time.Now()
	sm.resizeCount++

	return true // Changed
}

// GetSize returns current width and height
func (sm *SizeManager) GetSize() (int, int) {
	return sm.width, sm.height
}

// GetWidth returns current width
func (sm *SizeManager) GetWidth() int {
	return sm.width
}

// GetHeight returns current height
func (sm *SizeManager) GetHeight() int {
	return sm.height
}

// IsValidSize checks if current size meets minimum requirements
func (sm *SizeManager) IsValidSize() bool {
	return sm.width >= sm.minWidth && sm.height >= sm.minHeight
}

// GetErrorMessage returns error message if size is invalid
func (sm *SizeManager) GetErrorMessage() string {
	if !sm.IsValidSize() {
		return "Terminal too small. Please resize to at least 60x20"
	}
	return ""
}

// CalculatePanelWidths calculates left and right panel widths
func (sm *SizeManager) CalculatePanelWidths() (leftWidth, rightWidth int) {
	// Reserve space for borders and padding
	totalWidth := sm.width - 4 // Account for borders

	// Left panel: 1/3 of available width
	leftWidth = totalWidth / 3
	if leftWidth < 20 {
		leftWidth = 20
	}

	// Right panel: remaining width
	rightWidth = totalWidth - leftWidth - 1 // -1 for spacing
	if rightWidth < 30 {
		rightWidth = 30
	}

	return leftWidth, rightWidth
}

// CalculatePanelHeights calculates panel heights
func (sm *SizeManager) CalculatePanelHeights() (contentHeight int) {
	// Reserve space for header and footer
	headerHeight := 3
	footerHeight := 1

	contentHeight = sm.height - headerHeight - footerHeight
	if contentHeight < 10 {
		contentHeight = 10
	}

	return contentHeight
}

// CalculateSectionHeights calculates heights for different sections
func (sm *SizeManager) CalculateSectionHeights() (searchH, keysH, infoH, metaH, valueH, outputH int) {
	contentHeight := sm.CalculatePanelHeights()

	// Left panel sections
	searchH = 3
	infoH = 3
	keysH = contentHeight - searchH - infoH

	if keysH < 5 {
		keysH = 5
	}

	// Right panel sections
	metaH = 4
	remainingH := contentHeight - metaH
	valueH = remainingH / 2
	outputH = remainingH - valueH

	if valueH < 3 {
		valueH = 3
	}
	if outputH < 3 {
		outputH = 3
	}

	return searchH, keysH, infoH, metaH, valueH, outputH
}

// GetResizeInfo returns resize information
func (sm *SizeManager) GetResizeInfo() string {
	return fmt.Sprintf("Size: %dx%d | Resizes: %d | Last: %s",
		sm.width, sm.height,
		sm.resizeCount,
		sm.lastUpdate.Format("15:04:05"))
}

// SetMinimumSize sets minimum size requirements
func (sm *SizeManager) SetMinimumSize(width, height int) {
	sm.minWidth = width
	sm.minHeight = height
}
