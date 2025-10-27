# Changelog

All notable changes to GoCmder project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Linux Kernel Configuration Style UI** (`internal/ui/dashboard.go`)
  - Two-panel layout (left: schemes, right: details) for better organization
  - Predefined installation schemes: Minimal, Go Developer, Node Developer, Backend, Full Stack, Custom, Personal Settings
  - Customizable schemes: users can modify any predefined scheme, automatically switching to Custom
  - Interactive detail panel showing all tools and settings with checkbox indicators
  - Color-coded status: green for selected items, yellow for unselected items
  - Tab key navigation between left (schemes) and right (details) panels
  - Space key to toggle tool/setting selection directly in right panel
  - Automatic state preservation when switching from predefined schemes to Custom
- **Bootstrap Architecture** (`internal/bootstrap/`)
  - `Application` struct for centralized application lifecycle management
  - Unified entry point for all application components
  - Chainable API for flexible configuration (e.g., `New().SetLogger(lg).Setup()`)
- **Enhanced Logger System** (`internal/logger/`)
  - Multi-level logging: Debug, Info, Warn, Error
  - Dual output: console and file (daily log files in `logs/` directory)
  - Thread-safe logging with mutex synchronization
  - Structured logging with timestamps
  - Global logger functions for easy access
  - Instance methods for fine-grained control
  - Log file management with automatic directory creation
- **Application Integration**
  - Logger initialization in main and bootstrap
  - Graceful shutdown with log cleanup (`defer logger.Close()`)
  - Automatic logger recovery in application setup

### Changed
- Refactored models package into multiple files: models.go (data structures), constants.go (indices), config.go (configuration logic)
- Refactored UI to use scheme-based navigation instead of direct tool selection
- `InstallConfig` now includes `Schemes` slice and `CurrentScheme` index for state management
- Dashboard now supports two navigation modes: "scheme" mode (selecting predefined schemes) and "tools" mode (direct tool customization)
- ESC key returns from Tools Mode to Scheme Mode (Ctrl+C now sole exit mechanism)
- Custom scheme displays actual user-modified state instead of predefined state
- Tool and setting selection now preserves other selected items when switching from predefined schemes
- Refactored `main.go` to use structured logging instead of standard `log` package
- Renamed field `tvApp` to `tviewApp` in Application struct for clarity
- Application now manages logger instance as a core component
- Improved error handling with descriptive log messages
- Simplified application startup with bootstrap pattern

### Improved
- Dashboard-style TUI interface with tview
- Tool selection panel with version and size information
- System settings configuration panel
- Quick action shortcuts (Select All, Clear All, Enable Settings)
- Real-time selection counters and status messages
- Support for 7 development tools:
  - Git 2.43.0
  - VSCode 1.84.2
  - Go 1.21.3
  - Node.js 20.10.0
  - PostgreSQL 16.0
  - MySQL 8.1.0
  - Redis 3.0.504
- System configuration features:
  - Add custom paths to system PATH
  - Configure power settings (monitor timeout, sleep mode)
- Local file priority with automatic download fallback
- Comprehensive project structure with organized directories
- Core business logic modules:
  - UI components (dashboard)
  - Installer implementations (Windows)
  - Setup and recovery logic
  - System and platform detection
  - Configuration management
  - Logging system (enhanced)
  - Bootstrap initialization
- Download management for platform-specific installers
- Backup management for user data and configurations
- Documentation in README.md

### Planned
- macOS installer support
- Linux installer support
- Real-time installation progress tracking
- Installation history tracking
- Uninstall functionality
- Configuration profiles (save/load presets)
- Update checker for tool versions
- Theme support (Dark/Light mode)
- Configuration file support (YAML)
- Advanced system detection
- Rollback functionality
- Async logger with buffering
- Log rotation strategy

## [0.1.0] - 2025-10-26

### Initial Release
- Project structure established following Go best practices
- Basic TUI framework setup with tview
- Dashboard UI prototype
- Data models for tools and settings
- Windows installer skeleton
- Documentation and README

### Project Structure
- `/internal/` - Core business logic
  - `ui/` - Dashboard components
  - `installer/` - Platform-specific installers
  - `setup/` - Setup and recovery logic
  - `detect/` - System detection
  - `config/` - Configuration management
  - `logger/` - Logging
  - `bootstrap/` - Application initialization
  - `models/` - Data models
- `/scripts/` - Shell/PowerShell scripts
- `/docs/` - Documentation
- `/downloads/` - Platform-specific downloads
- `/backups/` - User data backups

---

## Development Notes

### Current Status
- **Phase**: Early Development (v0.1.0)
- **Platform**: Windows (primary), macOS and Linux planned
- **UI Framework**: tview (tcell)
- **Language**: Go 1.24+
- **Architecture**: Bootstrap pattern with lifecycle management

### Next Steps
1. Implement core installer logic for Windows
2. Add system detection module
3. Implement configuration manager (YAML support)
4. Enhance logging with async operations
5. Create setup and recovery logic
6. Extend to macOS support
7. Extend to Linux support

### Known Limitations
- Currently Windows-only
- No actual installation execution (UI only)
- No configuration file support yet
- No progress tracking
- No installation history
- Logger is synchronous (no buffering)

### Architecture Decisions
- **Bootstrap Pattern**: Centralized application lifecycle management
- **Logger-First Approach**: Logging initialized early in main()
- **Modular Design**: Separate concerns into internal packages
- **Thread-Safe Logger**: Mutex-protected for concurrent access
- **Dual Output**: Both console and file logging for debugging
- **Daily Log Files**: Automatic log rotation by date
- **TUI Framework**: tview for cross-platform terminal UI
- **Directory Structure**: Organized by functionality and purpose

### Testing Strategy
- Unit tests for models and utilities
- Integration tests for installer logic
- Manual testing for UI components
- Platform-specific testing for each OS
- Logger functionality testing with file verification

### Performance Considerations
- Lazy loading of tools and settings
- Efficient file downloads with caching
- Minimal memory footprint for TUI
- Async operations for long-running tasks
- Logger mutex to prevent concurrent write conflicts

---

## Contributing

When adding new features or changes:
1. Update this CHANGELOG.md
2. Follow the format: Added, Changed, Deprecated, Removed, Fixed, Security
3. Include version number and date for releases
4. Link to relevant issues/PRs if applicable

## Versioning

This project uses Semantic Versioning:
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

---

## Resources

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [tview Documentation](https://github.com/rivo/tview)
- [tcell Documentation](https://github.com/gdamore/tcell)
