# GoCmder UI Architecture

## Overview

The GoCmder UI has been refactored based on the podman-tui architecture pattern, providing a modular, maintainable, and extensible terminal user interface using the tview library.

## Architecture Pattern

The UI follows a hierarchical component-based architecture similar to podman-tui:

```
internal/ui/
├── app.go                    # Main UI coordinator
├── style/                    # Color and styling definitions
│   └── style.go
├── utils/                    # Common utilities
│   ├── ui_dialog.go         # Dialog interface
│   ├── keys.go              # Key bindings
│   └── utils.go             # Helper functions
├── dialogs/                  # Reusable dialog components
│   ├── dialogs.go           # Common constants
│   ├── error.go             # Error dialog
│   ├── confirm.go           # Confirmation dialog
│   ├── message.go           # Message dialog
│   ├── progress.go          # Progress dialog
│   └── input.go             # Input dialog
├── tools/                    # Tools page
│   └── tools.go
├── settings/                 # Settings page
│   └── settings.go
└── system/                   # System info page
    └── system.go
```

## Key Components

### 1. Style Package (`ui/style/`)

Centralized color and styling definitions for consistent UI appearance:

- Background colors (dialogs, pages, headers)
- Foreground colors (text, headers, info bars)
- Border colors
- Status colors (installed, not installed, selected, error)
- Color helper functions

### 2. Utils Package (`ui/utils/`)

Common utilities used throughout the UI:

**UIDialog Interface:**
```go
type UIDialog interface {
    tview.Primitive
    IsDisplay() bool
    Hide()
}
```

**Key Bindings:**
- `CloseDialogKey`: ESC - Close dialogs
- `SwitchFocusKey`: Tab - Switch focus/pages
- `ConfirmKey`: Enter - Confirm actions
- `ToggleKey`: Space - Toggle selection
- `SelectAllKey`: 'a' - Select all items
- `InstallKey`: 'i' - Install selected
- `RefreshKey`: 'r' - Refresh data
- `QuitKey`: 'q' - Quit application

**Helper Functions:**
- `EmptyBoxSpace()`: Create empty box with background color
- `CreateStyledTable()`: Create table with standard styling
- `TruncateString()`: Truncate strings to fit display
- `AlignStringRight()`: Right-align strings

### 3. Dialogs Package (`ui/dialogs/`)

Reusable dialog components that implement the `UIDialog` interface:

**ErrorDialog:**
- Display error messages with title
- Modal dialog with OK button
- Red background for visibility

**ConfirmDialog:**
- Confirmation prompts with Yes/No buttons
- Customizable title and message
- Callback handlers for both buttons

**MessageDialog:**
- Display informational messages
- Modal dialog with OK button
- Customizable title

**ProgressDialog:**
- Show progress bar for long operations
- Display progress percentage
- Customizable message

**SimpleInputDialog:**
- Get text input from user
- Input field with label
- OK/Cancel buttons
- Focus management between field and buttons

### 4. Page Components

Each page implements a common interface:

```go
type UIPage interface {
    tview.Primitive
    GetTitle() string
    HasFocus() bool
    SubDialogHasFocus() bool
    HideAllDialogs()
    SetAppFocusHandler(handler func())
}
```

#### Tools Page (`ui/tools/`)

Displays and manages development tools:

**Features:**
- Table view with columns: Name, Version, Size, Status, Selected
- Toggle selection with Space key
- Select all with 'a' key
- Install selected tools with 'i' key
- Color-coded status (Green=Installed, Yellow=Not Installed)
- Confirmation dialog before installation
- Progress feedback during operations

**Data Structure:**
```go
type Tool struct {
    Name      string
    Version   string
    Size      string
    Selected  bool
    Installed bool
}
```

#### Settings Page (`ui/settings/`)

Manages system configuration settings:

**Features:**
- Table view with columns: Setting, Description, Selected
- Toggle selection with Space key
- Apply settings with Enter key
- Confirmation dialog before applying
- Settings descriptions for clarity

**Data Structure:**
```go
type Setting struct {
    Name     string
    Selected bool
}
```

**Available Settings:**
- Add to PATH
- Configure Power Settings
- Set User Directories

#### System Page (`ui/system/`)

Displays system information:

**Features:**
- TextView showing system details
- OS information (name, architecture)
- Hardware info (CPU cores)
- Runtime info (Go version)
- Application capabilities
- Keyboard shortcuts help
- Scrollable content

### 5. Main UI Coordinator (`ui/app.go`)

Central component that manages all pages and navigation:

**Features:**
- Multi-page navigation with Tab key
- Direct page access (1=Tools, 2=Settings, 3=System)
- Info bar showing application title
- Help bar with context-sensitive shortcuts
- Global key bindings (quit, navigation)
- Focus management across pages
- Dialog coordination

**Structure:**
```go
type App struct {
    app            *tview.Application
    pages          *tview.Pages
    layout         *tview.Flex
    infoBar        *tview.TextView
    helpBar        *tview.TextView
    toolsPage      *tools.Tools
    settingsPage   *settings.Settings
    systemPage     *system.System
    currentPageIdx int
    // handlers...
}
```

## Integration with Bootstrap

The bootstrap layer (`internal/bootstrap/app.go`) integrates the UI:

```go
type Application struct {
    Tview        *tview.Application
    Config       *config.Config
    UI           *ui.App
    Logger       *logger.Logger
    toolsData    []models.Tool
    settingsData []models.Setting
}
```

**Initialization Flow:**
1. Create Application instance
2. Initialize logger and config
3. Initialize default data (tools and settings)
4. Create UI with NewApp()
5. Set handlers (install, apply, refresh)
6. Update initial data
7. Set root primitive and run

**Handler Pattern:**
```go
// Tool installation handler
func (a *Application) handleInstallTool(toolName string) {
    a.logInfo("Install requested for tool: %s", toolName)
    // Implementation here
}

// Settings application handler
func (a *Application) handleApplySettings(settings []models.Setting) {
    a.logInfo("Apply settings requested for %d settings", len(settings))
    // Implementation here
}

// Refresh handler
func (a *Application) handleRefresh() {
    a.logInfo("Refresh requested")
    // Implementation here
}
```

## Key Features

### Focus Management

The architecture implements proper focus management:
- Each page manages its own dialogs' focus
- App coordinator handles page-level focus
- Dialogs take precedence when displayed
- Tab key switches between pages/elements

### Input Handling

Input is handled hierarchically:
1. Global key bindings (quit, page navigation)
2. Dialog-level key bindings (if dialog is visible)
3. Page-level key bindings
4. Primitive-level default handlers

### Data Flow

```
User Action → Page Component → Handler Callback → Bootstrap Layer → Business Logic
                                                          ↓
User Display ← Page Component ← Data Update ← Bootstrap Layer
```

### Dialog Pattern

All dialogs follow a consistent pattern:
1. Implement `UIDialog` interface
2. Support `Display()`, `IsDisplay()`, `Hide()` methods
3. Provide callback setters (`SetSelectedFunc`, `SetCancelFunc`)
4. Handle input with `InputHandler()`
5. Draw with proper positioning in parent's `Draw()` method

## Color Scheme

The UI uses a consistent color scheme:

- **Background**: Default/Black
- **Foreground**: White
- **Headers**: Dark Cyan background, Black text
- **Borders**: Gray/White
- **Status Colors**:
  - Installed: Green
  - Not Installed: Yellow
  - Error: Red
  - Selected: Light Cyan

## Keyboard Shortcuts

### Global
- `Tab`: Switch to next page
- `1`, `2`, `3`: Jump to specific page
- `q`: Quit application
- `ESC`: Close dialog

### Tools Page
- `Space`: Toggle tool selection
- `a`: Select all tools
- `i`: Install selected tools
- `r`: Refresh tool list
- `↑`/`↓`: Navigate list

### Settings Page
- `Space`: Toggle setting selection
- `a`: Select all settings
- `Enter`: Apply selected settings
- `↑`/`↓`: Navigate list

### System Page
- `↑`/`↓`: Scroll content
- `PageUp`/`PageDown`: Scroll page

## Extension Guide

### Adding a New Page

1. Create new package under `internal/ui/`
2. Implement the `UIPage` interface
3. Add page to `App.pageList` in `ui/app.go`
4. Add navigation key in `globalInputHandler`
5. Update help bar in `updateHelpBar()`

### Adding a New Dialog

1. Create new dialog in `internal/ui/dialogs/`
2. Implement `UIDialog` interface
3. Add to page component that will use it
4. Include in page's `HasFocus()` and `HideAllDialogs()`
5. Add to page's `Focus()` method
6. Draw in page's `Draw()` method

### Customizing Styles

Edit `internal/ui/style/style.go`:
1. Define new color constants
2. Add to `GetColorHex()` map if needed
3. Use throughout UI components

## Benefits of This Architecture

1. **Modularity**: Each component has clear responsibilities
2. **Reusability**: Dialogs and utilities are shared
3. **Maintainability**: Similar to podman-tui, well-understood pattern
4. **Testability**: Components can be tested independently
5. **Extensibility**: Easy to add new pages and dialogs
6. **Consistency**: Centralized styling and behavior
7. **Focus Management**: Proper handling of complex UI flows

## Future Enhancements

Potential improvements:
- Add more dialog types (multi-select, tree view)
- Implement command palette
- Add help page with interactive documentation
- Support themes/skins
- Add logging viewer
- Implement real-time system monitoring
- Add plugin architecture for custom pages

