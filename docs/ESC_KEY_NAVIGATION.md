# ESC Key Navigation Implementation

## Overview
Implemented ESC key navigation to return to Home page from all other pages without exiting the application.

## Changes Made

### 1. Modified `internal/ui/app.go`

#### globalInputHandler Function
Added ESC key handling logic:

```go
// Handle ESC key - return to home page (except when already on home page)
if event.Key() == utils.CloseDialogKey.Key {
    if a.currentPageIdx != homePageIndex {
        a.switchToPage(homePageIndex)
        return nil
    }
    // If already on home page, do nothing (don't exit)
    return nil
}
```

**Key Features:**
- Checks if any dialog has focus first (dialogs handle ESC to close themselves)
- If on any page except Home, ESC returns to Home page
- If already on Home page, ESC does nothing (no application exit)
- Dialog ESC handling takes priority over page navigation

#### updateHelpBar Function
Updated to show "ESC Home" on all pages except the Home page:

```go
case terminalPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Enter[-] Execute"
case databasePageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Ctrl+N[-] Connect | ..."
case toolsPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | ..."
case settingsPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]Space[-] Toggle | ..."
case systemPageIndex:
    pageHelp = " | [" + highlightColor + "]ESC[-] Home | [" + highlightColor + "]↑/↓[-] Scroll"
```

## Navigation Flow

### ESC Key Priority Order:
1. **Dialog Open**: ESC closes the dialog (stays on current page)
2. **Non-Home Page**: ESC returns to Home page
3. **Home Page**: ESC does nothing (application stays open)

### Example Flow:
```
Home (F1) → Database (F4) → [ESC] → Home (F1)
Home (F1) → Tools (F6) → [ESC] → Home (F1)
Database → Ctrl+N (open dialog) → [ESC] (close dialog) → [ESC] → Home
```

## User Experience

### Before Changes:
- ESC key had no consistent behavior across pages
- Users might accidentally exit when pressing ESC
- No clear way to return to Home page quickly

### After Changes:
- ✅ Consistent ESC behavior across all pages
- ✅ Safe navigation - ESC never exits the application
- ✅ Intuitive - ESC returns to the main/home page
- ✅ Dialog-aware - ESC closes dialogs first
- ✅ Visual feedback - Help bar shows "ESC Home" on relevant pages

## Keyboard Shortcuts Summary

| Key | Action | Notes |
|-----|--------|-------|
| **ESC** | Return to Home page | Does nothing on Home page; closes dialogs first |
| **F1** | Go to Home page | Direct navigation |
| **F2** | Go to Terminal page | Direct navigation |
| **F4** | Go to Database page | Direct navigation |
| **F6** | Go to Tools page | Direct navigation |
| **F7** | Go to Settings page | Direct navigation |
| **F8** | Go to System page | Direct navigation |
| **Tab** | Cycle through pages | Sequential navigation |
| **q** | Quit application | Exits completely |

## Technical Details

### Implementation Pattern:
The ESC key handler is placed early in the global input handler chain:
1. Check for dialog focus → if dialog is open, pass event to dialog
2. Check for ESC key → handle navigation
3. Handle other global keys (q, Tab, Function keys)

This ensures proper event bubbling and prevents conflicts.

### Code Quality:
- ✅ No breaking changes to existing functionality
- ✅ Maintains existing dialog behavior
- ✅ Follows project's code style and patterns
- ✅ Uses existing constants (`utils.CloseDialogKey`)
- ✅ Builds successfully with no errors
- ✅ All pages tested and functional

## Testing

See [TEST_ESC_NAVIGATION.md](../TEST_ESC_NAVIGATION.md) for detailed test cases and verification steps.

## Compatibility

- **Windows**: Tested and working
- **Linux**: Should work (uses tcell library)
- **macOS**: Should work (uses tcell library)

The implementation uses the cross-platform `tcell` library which handles ESC key (KeyEscape) consistently across all platforms.

