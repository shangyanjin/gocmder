# GoCmder UI Layout - Linux Kernel Configuration Style

## Overview

The GoCmder application now features a Linux kernel configuration style interface with a left-right layout for intuitive navigation and selection.

## Layout Structure

```
┌─────────────────────────────────────────────────────────────────┐
│ GoCmder Setup - Scheme Mode                                     │
├─────────────────────────┬──────────────────────────────────────┤
│ Options                 │ Details                              │
├─────────────────────────┼──────────────────────────────────────┤
│                         │                                      │
│ ✓ Minimal               │ Minimal                              │
│ [ ] Go Developer        │ Basic tools: Git, VSCode             │
│ [ ] Node Developer      │                                      │
│ [ ] Backend             │ Tools (2):                           │
│ [ ] Full Stack          │   • Git                              │
│ [ ] Custom              │   • VSCode                           │
│                         │                                      │
│                         │ Settings (1):                        │
│                         │   • Add PATH                         │
│                         │                                      │
│                         │ [Press Space to apply this scheme]   │
│                         │                                      │
├─────────────────────────┴──────────────────────────────────────┤
│ [Y/n/?] Space:Select Tab:Mode F5:Run F9:Output Esc:Exit        │
└─────────────────────────────────────────────────────────────────┘
```

Features:
- Top: Title bar (1 line) with current mode
- Left: Options panel with border (fixed 25 columns)
- Right: Details panel with border (flexible width)
- Bottom: Status bar (1 line) with shortcuts

## Components

### Top: Title Panel
- Shows current navigation mode (Scheme Mode / Tools Mode)
- Always visible for context

### Left: Options Panel (25 columns)
- Lists available schemes or tools depending on mode
- Shows checkbox status (✓ selected, [ ] unselected)
- Color coding:
  - Green: Selected items
  - Yellow: Unselected items
  - Blue: Scheme names (Scheme Mode)
- Current selection: Highlighted with dark blue background

### Right: Details Panel (remaining width)
- **Scheme Mode**: Shows scheme description and included tools/settings
- **Tools Mode**: Shows tool version, size, or setting details
- Provides context for selected item

### Bottom: Status Bar (1 line)
- Quick keyboard shortcut reference
- Shows count of selected tools and settings
- Displays [Y/n/?] style prompt

## Navigation Modes

### Scheme Mode (Default)
Select from predefined schemes:
1. **Minimal** - Basic tools (Git, VSCode)
2. **Go Developer** - Go stack (Git, VSCode, Go)
3. **Node Developer** - Node stack (Git, VSCode, Node.js)
4. **Backend** - Backend tools (Go, PostgreSQL, MySQL, Redis)
5. **Full Stack** - All tools (Git, VSCode, Go, Node.js, PostgreSQL, MySQL, Redis)
6. **Custom** - Manual selection (empty by default)
7. **Personal Settings** - Configure personal environment (Add PATH, PowerConfig, SetUserDirs)
8. **Exit** - Exit the application

### Tools Mode
View and customize individual tool selection:
- Lists all tools and settings
- Toggle each item independently
- Fine-tune the configuration from any scheme

## Keyboard Shortcuts

| Key | Function |
|-----|----------|
| ↑ / ↓ | Navigate options (with circular wrapping) |
| Tab | Switch between Scheme Mode and Tools Mode |
| Space | Apply scheme / Toggle tool or setting |
| F2 | Focus options panel |
| F3 | Refresh system information |
| F5 | Run / Execute installation |
| F9 | Show output panel |
| **Esc** | **Return to Scheme Mode (only works in Tools Mode)** |
| **Ctrl+C** | **Exit application (works from any mode)** |

## Color Scheme

- **Cyan**: Headings, mode indicators, status symbols
- **Yellow**: Selection highlight, unselected items
- **Green**: Selected/enabled items
- **Blue**: Scheme names, category headers
- **White**: Default text
- **Dark Blue**: Selection background

## Typical Workflow

1. Start in **Scheme Mode** (default)
2. **Navigate** schemes with ↑ ↓ arrow keys
3. **Preview** tools and settings in right panel
4. **Apply** with Space to auto-select all tools for that scheme
5. **Switch** to Tools Mode with Tab for fine-tuning
6. **Adjust** individual selections if needed
7. **Return** to Scheme Mode with Esc (from Tools Mode only)
8. **Execute** with F5 to start installation
9. **Exit** with Ctrl+C (works from any mode)

## Design Benefits

Inspired by Linux kernel `make menuconfig`:
- **Clarity**: Clear visual separation of selection and details
- **Efficiency**: Quick access to predefined schemes
- **Flexibility**: Easy customization in Tools mode
- **Familiarity**: Similar interaction model for Linux developers
- **Simplicity**: Minimal navigation complexity
