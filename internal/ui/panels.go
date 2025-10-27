package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shangyanjin/gocmder/internal/models"
)

// PanelComponents holds all UI components
type PanelComponents struct {
	// Left panel components
	SearchBox *tview.InputField
	KeysList  *tview.List
	InfoText  *tview.TextView

	// Right panel components
	MetaText   *tview.TextView
	ValueText  *tview.TextView
	OutputText *tview.TextView

	// Layout components
	Header *tview.TextView
	Footer *tview.TextView
}

// NewPanelComponents creates new panel components
func NewPanelComponents() *PanelComponents {
	return &PanelComponents{}
}

// CreateSearchPanel creates the search panel
func (pc *PanelComponents) CreateSearchPanel() *tview.Flex {
	pc.SearchBox = tview.NewInputField().
		SetLabel("Key: ").
		SetFieldWidth(20).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				// Handle search - this will be set by the dashboard
			}
		})

	searchTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Search (F2)[white]")

	searchSection := tview.NewFlex().SetDirection(tview.FlexRow)
	searchSection.AddItem(searchTitle, 1, 0, false)
	searchSection.AddItem(pc.SearchBox, 1, 0, true)

	return searchSection
}

// CreateKeysPanel creates the keys list panel
func (pc *PanelComponents) CreateKeysPanel() *tview.Flex {
	pc.KeysList = tview.NewList().
		SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			// Handle selection - this will be set by the dashboard
		})

	keysTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Keys (F3)[white]")

	keysSection := tview.NewFlex().SetDirection(tview.FlexRow)
	keysSection.AddItem(keysTitle, 1, 0, false)
	keysSection.AddItem(pc.KeysList, 0, 1, true)

	return keysSection
}

// CreateInfoPanel creates the info panel
func (pc *PanelComponents) CreateInfoPanel() *tview.Flex {
	pc.InfoText = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false)

	infoTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Info[white]")

	infoSection := tview.NewFlex().SetDirection(tview.FlexRow)
	infoSection.AddItem(infoTitle, 1, 0, false)
	infoSection.AddItem(pc.InfoText, 2, 0, false)

	return infoSection
}

// CreateMetaPanel creates the meta information panel
func (pc *PanelComponents) CreateMetaPanel() *tview.Flex {
	pc.MetaText = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false)

	metaTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Meta[white]")

	metaSection := tview.NewFlex().SetDirection(tview.FlexRow)
	metaSection.AddItem(metaTitle, 1, 0, false)
	metaSection.AddItem(pc.MetaText, 4, 0, false)

	return metaSection
}

// CreateValuePanel creates the value display panel
func (pc *PanelComponents) CreateValuePanel() *tview.Flex {
	pc.ValueText = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false)

	valueTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Value[white]")

	valueSection := tview.NewFlex().SetDirection(tview.FlexRow)
	valueSection.AddItem(valueTitle, 1, 0, false)
	valueSection.AddItem(pc.ValueText, 0, 1, false)

	return valueSection
}

// CreateOutputPanel creates the output/log panel
func (pc *PanelComponents) CreateOutputPanel() *tview.Flex {
	pc.OutputText = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(false).
		SetScrollable(true)

	outputTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Output (F9)[white]")

	outputSection := tview.NewFlex().SetDirection(tview.FlexRow)
	outputSection.AddItem(outputTitle, 1, 0, false)
	outputSection.AddItem(pc.OutputText, 0, 1, false)

	return outputSection
}

// CreateHeader creates the header section
func (pc *PanelComponents) CreateHeader() *tview.TextView {
	pc.Header = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(false)

	return pc.Header
}

// CreateFooter creates the footer section
func (pc *PanelComponents) CreateFooter() *tview.TextView {
	pc.Footer = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[white]Total matched: 9 | Focus: TOOLS | [Tab] Switch | [Space] Toggle | [Enter] Install | [q] Quit[white]")

	return pc.Footer
}

// UpdateHeader updates the header with current information
func (pc *PanelComponents) UpdateHeader(config *models.InstallConfig, sizeManager *SizeManager) {
	selectedTools := config.GetSelectedToolsCount()
	selectedSettings := config.GetSelectedSettingsCount()

	// Get system info for header
	systemSummary := GetSystemInfoCompact()

	// Get resize info
	resizeInfo := sizeManager.GetResizeInfo()

	headerText := fmt.Sprintf(
		"[white]GoCmderVersion: 1.0.0 | %s[yellow]\n",
		systemSummary,
	)

	keyspace := fmt.Sprintf(
		"[white]KeySpace: tools=%d, settings=%d, selected=%d | %s[yellow]\n",
		len(config.Tools), len(config.Settings), selectedTools+selectedSettings, resizeInfo,
	)

	hint := "[white]* Press F1 to switch between panels, Esc to quit[yellow]"

	pc.Header.SetText(headerText + keyspace + hint)
}

// UpdateKeysList populates the keys list with tools and settings
func (pc *PanelComponents) UpdateKeysList(config *models.InstallConfig) {
	pc.KeysList.Clear()

	// Add tools
	for i, tool := range config.Tools {
		checkbox := "☐"
		if tool.Selected {
			checkbox = "☑"
		}

		itemText := fmt.Sprintf("%d | %s %s", i+1, checkbox, tool.Name)
		pc.KeysList.AddItem(itemText, "", 0, nil)
	}

	// Add settings
	for i, setting := range config.Settings {
		checkbox := "☐"
		if setting.Selected {
			checkbox = "☑"
		}

		itemText := fmt.Sprintf("%d | %s %s", len(config.Tools)+i+1, checkbox, setting.Name)
		pc.KeysList.AddItem(itemText, "", 0, nil)
	}
}

// UpdateInfo updates the info section
func (pc *PanelComponents) UpdateInfo(config *models.InstallConfig) {
	total := len(config.Tools) + len(config.Settings)
	infoText := fmt.Sprintf("Total matched: %d", total)
	pc.InfoText.SetText(infoText)
}

// UpdateMeta updates the meta information
func (pc *PanelComponents) UpdateMeta(config *models.InstallConfig, currentPanel string, selectedIndex int) {
	if currentPanel == "tools" && selectedIndex < len(config.Tools) {
		tool := config.Tools[selectedIndex]
		metaText := fmt.Sprintf("Key: %s\nType: tool\nTTL: -1s", tool.Name)
		pc.MetaText.SetText(metaText)
	} else if currentPanel == "settings" && selectedIndex < len(config.Settings) {
		setting := config.Settings[selectedIndex]
		metaText := fmt.Sprintf("Key: %s\nType: setting\nTTL: -1s", setting.Name)
		pc.MetaText.SetText(metaText)
	} else {
		// Show system information when nothing is selected
		systemSummary := GetSystemInfoCompact()
		pc.MetaText.SetText(systemSummary)
	}
}

// UpdateValue updates the value display
func (pc *PanelComponents) UpdateValue() {
	// Show detailed system information with safe formatting
	systemInfo := GetSystemInfoText()

	// Limit the length to prevent overflow
	maxLength := 2000 // Reasonable limit
	if len(systemInfo) > maxLength {
		systemInfo = systemInfo[:maxLength] + "\n... (truncated)"
	}

	pc.ValueText.SetText(systemInfo)
}

// UpdateOutput updates the output section
func (pc *PanelComponents) UpdateOutput(message, messageType, currentPanel string) {
	if message != "" {
		timestamp := time.Now().Format("2006-01-02T15:04:05+08:00")
		var prefix string

		switch messageType {
		case "success":
			prefix = "✓"
		case "error":
			prefix = "✗"
		default:
			prefix = "ℹ"
		}

		logLine := fmt.Sprintf("[%s] %s %s, type=%s, ttl=-1s",
			timestamp, prefix, message, currentPanel)
		pc.OutputText.SetText(logLine)
	} else {
		// Show default log entries
		logs := []string{
			"[2025-01-26T20:30:00+08:00] query tools OK, type=tool, ttl=-1s",
			"[2025-01-26T20:30:01+08:00] query settings OK, type=setting, ttl=-1s",
			"[2025-01-26T20:30:02+08:00] query config OK, type=config, ttl=-1s",
		}

		pc.OutputText.SetText(strings.Join(logs, "\n"))
	}
}
