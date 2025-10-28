# GoCmder UI Structure Overview

## Component Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                    main.go                                      │
│                       ↓                                         │
│              bootstrap.Application                             │
│                       ↓                                         │
│                   ui.App (Main UI Coordinator)                 │
└─────────────────────────────────────────────────────────────────┘
                            ↓
    ┌──────────────────────┼──────────────────────┐
    ↓                      ↓                      ↓
┌─────────┐          ┌──────────┐          ┌──────────┐
│ Tools   │          │ Settings │          │  System  │
│  Page   │          │   Page   │          │   Page   │
└─────────┘          └──────────┘          └──────────┘
    │                     │                      │
    └──────────────────┬──┴──────────────────────┘
                       ↓
        ┌──────────────────────────────┐
        │      Common Components       │
        ├──────────────────────────────┤
        │  • ErrorDialog               │
        │  • ConfirmDialog             │
        │  • MessageDialog             │
        │  • ProgressDialog            │
        │  • SimpleInputDialog         │
        │  • Style (colors/themes)     │
        │  • Utils (helpers)           │
        │  • Keys (bindings)           │
        └──────────────────────────────┘
```

## Page Components Detail

### Tools Page
```
┌────────────────────────────────────────────────────────┐
│ [::b]DEVELOPMENT TOOLS[0]                              │
├────────────────────────────────────────────────────────┤
│ TOOL      │ VERSION  │ SIZE    │ STATUS        │ SEL   │
├───────────┼──────────┼─────────┼───────────────┼───────┤
│ Git       │ 2.43.0   │ ~50 MB  │ Not Installed │ [ ]   │
│ VSCode    │ 1.84.2   │ ~100 MB │ Not Installed │ [ ]   │
│ Go        │ 1.21.3   │ ~130 MB │ Not Installed │ [X]   │
│ Node.js   │ 20.10.0  │ ~40 MB  │ Installed     │ [ ]   │
│ ...       │ ...      │ ...     │ ...           │ ...   │
└────────────────────────────────────────────────────────┘

Key Bindings:
  Space = Toggle Selection
  a     = Select All
  i     = Install Selected
  r     = Refresh
```

### Settings Page
```
┌────────────────────────────────────────────────────────────────┐
│ [::b]SYSTEM SETTINGS[0]                                        │
├────────────────────────────────────────────────────────────────┤
│ SETTING              │ DESCRIPTION                     │ SEL   │
├──────────────────────┼─────────────────────────────────┼───────┤
│ Add to PATH          │ Add dev tools to system PATH    │ [ ]   │
│ Configure Power      │ Set power config for dev        │ [X]   │
│ Set User Directories │ Configure user folders          │ [ ]   │
└────────────────────────────────────────────────────────────────┘

Key Bindings:
  Space = Toggle Selection
  a     = Select All
  Enter = Apply Settings
```

### System Page
```
┌────────────────────────────────────────────────────────┐
│ [::b]SYSTEM INFORMATION[0]                             │
├────────────────────────────────────────────────────────┤
│                                                        │
│ Operating System:                                      │
│   OS:           windows                                │
│   Architecture: amd64                                  │
│                                                        │
│ Hardware:                                              │
│   CPU Cores:    8                                      │
│                                                        │
│ Runtime:                                               │
│   Go Version:   go1.21.3                               │
│                                                        │
│ Application:                                           │
│   Name:         GoCmder                                │
│   Description:  Developer environment setup tool       │
│                                                        │
│ Capabilities:                                          │
│   • Install development tools                          │
│   • Configure system settings                          │
│   • Manage PATH environment                            │
│   ...                                                  │
│                                                        │
│ Keyboard Shortcuts:                                    │
│   Space  - Toggle selection                            │
│   a      - Select all                                  │
│   ...                                                  │
│                                                        │
└────────────────────────────────────────────────────────┘

Key Bindings:
  ↑/↓ = Scroll
```

## Dialog Components

### Error Dialog
```
┌─────────────────────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │                                 │ │
│ │    ERROR TITLE                  │ │
│ │                                 │ │
│ │    Error message text here...   │ │
│ │                                 │ │
│ │           [ OK ]                │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### Confirm Dialog
```
┌─────────────────────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │                                 │ │
│ │    CONFIRM TITLE                │ │
│ │                                 │ │
│ │    Confirmation message?        │ │
│ │                                 │ │
│ │      [ Yes ]    [ No ]          │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### Progress Dialog
```
┌─────────────────────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │                                 │ │
│ │    Installing Git...            │ │
│ │                                 │ │
│ │    [==========          ] 45%   │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### Input Dialog
```
┌─────────────────────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │                                 │ │
│ │    INPUT TITLE                  │ │
│ │                                 │ │
│ │    Input: [______________]      │ │
│ │                                 │ │
│ │      [ OK ]    [ Cancel ]       │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

## Screen Layout

```
┌──────────────────────────────────────────────────────────────────┐
│ GoCmder - Developer Environment Setup Tool                       │ ← Info Bar
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│                                                                  │
│                      Main Content Area                           │
│                   (Active Page Display)                          │
│                                                                  │
│                                                                  │
│                                                                  │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│ Tab: Switch | 1: Tools | 2: Settings | 3: System | q: Quit      │ ← Help Bar
└──────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Tool Installation Flow
```
User selects tools → Presses 'i' → Confirm Dialog
                                         ↓
                                      Confirmed?
                                         ↓ Yes
                     Progress Dialog ← Bootstrap.handleInstallTool
                                         ↓
                                    Installer.InstallTool
                                         ↓
                                    Update Status
                                         ↓
                                    Refresh Display
```

### Settings Application Flow
```
User selects settings → Presses Enter → Confirm Dialog
                                              ↓
                                          Confirmed?
                                              ↓ Yes
                         Bootstrap.handleApplySettings
                                              ↓
                                         Apply Settings
                                              ↓
                                        Success Message
```

## File Organization

```
gocmder/
├── main.go
├── internal/
│   ├── bootstrap/
│   │   └── app.go              # Application bootstrap
│   ├── ui/
│   │   ├── app.go              # Main UI coordinator
│   │   ├── style/
│   │   │   └── style.go        # Colors and styling
│   │   ├── utils/
│   │   │   ├── ui_dialog.go    # Dialog interface
│   │   │   ├── keys.go         # Key bindings
│   │   │   └── utils.go        # Helper functions
│   │   ├── dialogs/
│   │   │   ├── dialogs.go      # Constants
│   │   │   ├── error.go
│   │   │   ├── confirm.go
│   │   │   ├── message.go
│   │   │   ├── progress.go
│   │   │   └── input.go
│   │   ├── tools/
│   │   │   └── tools.go        # Tools page
│   │   ├── settings/
│   │   │   └── settings.go     # Settings page
│   │   └── system/
│   │       └── system.go       # System info page
│   ├── models/
│   │   ├── models.go           # Data models
│   │   └── constants.go        # Constants
│   ├── config/
│   ├── logger/
│   ├── installer/
│   └── detect/
└── docs/
    ├── UI_ARCHITECTURE.md      # Detailed architecture
    └── UI_STRUCTURE.md         # This file
```

## Key Design Principles

1. **Separation of Concerns**: Each component has a single responsibility
2. **Interface-Based Design**: Components implement common interfaces
3. **Reusability**: Dialogs and utilities are shared across pages
4. **Consistency**: All pages follow the same pattern
5. **Modularity**: Easy to add/remove/modify components
6. **Maintainability**: Clear structure, similar to podman-tui
7. **Extensibility**: New pages and dialogs can be added easily

## Comparison with podman-tui

### Similarities
- Component-based architecture
- Reusable dialog system
- Page-based navigation
- Focus management pattern
- Input handling hierarchy
- Style centralization

### Adaptations for gocmder
- Simplified page structure (3 pages vs many in podman-tui)
- Domain-specific components (tools, settings vs containers, images, etc.)
- Streamlined for developer environment setup
- Windows-first design considerations
- Integrated logging system

