# gocmder

GoCmder is a cross-platform terminal UI application for developers to rebuild their environment after system reinstall. It provides an interactive dashboard for installing Git, VSCode, Go, Node, PgSQL, MySQL, Redis, and configuring system settings.

## Features

- **Interactive TUI Dashboard**: Terminal-based user interface with real-time system information
- **Developer Tools Installation**: One-click installation of development tools
- **System Configuration**: Automated PATH setup, power settings, and personal folders
- **Multi-platform Support**: Windows (primary), macOS and Linux (planned)
- **Scheme-based Setup**: Predefined installation schemes for different development scenarios
- **Custom Configuration**: Fine-grained control over tools and settings
- **Real-time System Info**: CPU, memory, and system information display

## Components

### 1. Go Application (`main.go`)

The main interactive application with TUI interface:

- Dashboard-style terminal UI using tview
- Real-time system information
- Interactive tool selection with checkboxes
- Multi-panel layout with navigation support
- Automatic installer downloads
- Logging and error handling

### 2. Shell Scripts (`/scripts`)

Traditional shell and PowerShell scripts for quick setup:

- **Windows**: Run `scripts/win/Set-PersonalFolders.ps1` or `scripts/win/Install-DevTools.ps1`
- **macOS**: Run `scripts/macos/set-personal-folders.sh`
- **Linux**: Run `scripts/linux/set-personal-folders.sh`

## Installation

### Build from Source (Recommended)

**Prerequisites:**
- Go 1.21 or higher
- Git

**All Platforms (Windows / Linux / macOS):**

```bash
# Clone the repository
git clone https://github.com/shangyanjin/gocmder.git
cd gocmder

# Download dependencies
go mod download

# Build the application
go build -o gocmder .
```

**Windows:**
```powershell
# Build for Windows
go build -o gocmder.exe .

# Run
.\gocmder.exe
```

**macOS / Linux:**
```bash
# Build the application
go build -o gocmder .

# Run
./gocmder
```

## Usage

Choose your preferred method:

**Option 1: Shell Scripts (Quick)**
- Windows: `scripts/win/Set-PersonalFolders.ps1`
- macOS: `scripts/macos/set-personal-folders.sh`
- Linux: `scripts/linux/set-personal-folders.sh`

**Option 2: GoCmder (Interactive)**
- Build and run the Go application for an interactive dashboard experience

## Project Structure

```
gocmder/
├── main.go                    # Application entry point
├── internal/                  # Core business logic
│   ├── bootstrap/            # Application initialization
│   ├── config/               # Configuration management
│   ├── detect/               # System detection
│   ├── installer/            # Tool installation logic
│   ├── logger/               # Logging system
│   ├── models/               # Data models
│   ├── setup/                # Setup utilities
│   └── ui/                   # User interface components
├── scripts/                   # Shell/PowerShell scripts
│   ├── win/                  # Windows scripts
│   ├── macos/                # macOS scripts
│   └── linux/                # Linux scripts
├── docs/                      # Documentation
├── downloads/                 # Installer downloads
├── backups/                  # User backups
├── logs/                      # Application logs
├── README.md                  # This file
├── LICENSE                    # License
└── go.mod                     # Go module definition
```

## ⚠️ Disclaimer

**THIS IS A BETA/TEST VERSION - NOT FOR PRODUCTION USE**

- This software is provided "AS IS" without warranty of any kind
- No guarantees of reliability, security, or data integrity
- Use at your own risk - the author assumes no liability for any damages
- Recommended for testing and development environments only
- Always backup your data before use

## Repository

- **GitHub**: [https://github.com/shangyanjin/gocmder](https://github.com/shangyanjin/gocmder)

## Documentation

- [CHANGELOG](docs/CHANGELOG.md) - Version history and development notes
- [UI Layout](docs/UI_LAYOUT.md) - User interface documentation
- [Shell Scripts](scripts/) - Traditional automation scripts

## Author

Maintained by [shangyanjin](https://github.com/shangyanjin)
