package ui

import (
	"github.com/rivo/tview"
)

// LayoutManager manages the overall layout structure
type LayoutManager struct {
	root        *tview.Flex
	leftPanel   *tview.Flex
	rightPanel  *tview.Flex
	components  *PanelComponents
	sizeManager *SizeManager
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(components *PanelComponents, sizeManager *SizeManager) *LayoutManager {
	lm := &LayoutManager{
		components:  components,
		sizeManager: sizeManager,
	}
	lm.setupLayout()
	return lm
}

// setupLayout sets up the main layout structure
func (lm *LayoutManager) setupLayout() {
	// Create main layout (vertical split)
	lm.root = tview.NewFlex().SetDirection(tview.FlexRow)

	// Create main content area (horizontal split)
	mainContent := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create left panel with dynamic sizing
	lm.leftPanel = lm.createLeftPanel()

	// Create right panel with dynamic sizing
	lm.rightPanel = lm.createRightPanel()

	// Add panels to main content with dynamic proportions
	mainContent.AddItem(lm.leftPanel, 0, 1, true)
	mainContent.AddItem(lm.rightPanel, 0, 2, false)

	// Add all sections to root with dynamic heights
	lm.root.AddItem(lm.components.CreateHeader(), 3, 0, false)
	lm.root.AddItem(mainContent, 0, 1, true)
	lm.root.AddItem(lm.components.CreateFooter(), 1, 0, false)
}

// createLeftPanel creates the left panel with Search, Keys, and Info
func (lm *LayoutManager) createLeftPanel() *tview.Flex {
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Calculate dynamic heights with safety checks
	searchH, keysH, infoH, _, _, _ := lm.sizeManager.CalculateSectionHeights()

	// Ensure minimum heights
	if searchH < 1 {
		searchH = 1
	}
	if keysH < 3 {
		keysH = 3
	}
	if infoH < 1 {
		infoH = 1
	}

	// Add sections to left panel with safe heights
	leftPanel.AddItem(lm.components.CreateSearchPanel(), searchH, 0, false)
	leftPanel.AddItem(lm.components.CreateKeysPanel(), 0, 1, true) // Use flexible sizing
	leftPanel.AddItem(lm.components.CreateInfoPanel(), infoH, 0, false)

	return leftPanel
}

// createRightPanel creates the right panel with Meta, Value, and Output
func (lm *LayoutManager) createRightPanel() *tview.Flex {
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Calculate dynamic heights with safety checks
	_, _, _, metaH, valueH, outputH := lm.sizeManager.CalculateSectionHeights()

	// Ensure minimum heights
	if metaH < 1 {
		metaH = 1
	}
	if valueH < 3 {
		valueH = 3
	}
	if outputH < 3 {
		outputH = 3
	}

	// Add sections to right panel with safe heights
	rightPanel.AddItem(lm.components.CreateMetaPanel(), metaH, 0, false)
	rightPanel.AddItem(lm.components.CreateValuePanel(), 0, 1, false)  // Use flexible sizing
	rightPanel.AddItem(lm.components.CreateOutputPanel(), 0, 1, false) // Use flexible sizing

	return rightPanel
}

// GetRoot returns the root primitive
func (lm *LayoutManager) GetRoot() tview.Primitive {
	return lm.root
}

// GetComponents returns the panel components
func (lm *LayoutManager) GetComponents() *PanelComponents {
	return lm.components
}
