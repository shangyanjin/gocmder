# Developer Tools Installation Script for Windows
# Installs: Git, VSCode, Go, Node, PostgreSQL, MySQL, Redis
# Features: Local file priority, automatic download fallback, selective installation

param(
    [string]$LocalSourcePath = ".\installers",
    [switch]$SkipGit,
    [switch]$SkipVSCode,
    [switch]$SkipGo,
    [switch]$SkipNode,
    [switch]$SkipPostgreSQL,
    [switch]$SkipMySQL,
    [switch]$SkipRedis,
    [switch]$SkipAll
)

# Color output functions
function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ $Message" -ForegroundColor Cyan
}

function Write-Warning-Custom {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

# Check if running as administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Download file if not exists locally
function Get-InstallFile {
    param(
        [string]$FileName,
        [string]$DownloadUrl,
        [string]$LocalPath = $LocalSourcePath
    )
    
    $localFile = Join-Path $LocalPath $FileName
    
    if (Test-Path $localFile) {
        Write-Success "Found local file: $FileName"
        return $localFile
    }
    
    Write-Info "Local file not found, downloading: $FileName"
    
    if (-not (Test-Path $LocalPath)) {
        New-Item -ItemType Directory -Path $LocalPath -Force | Out-Null
    }
    
    try {
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $localFile -UseBasicParsing
        Write-Success "Downloaded: $FileName"
        return $localFile
    }
    catch {
        Write-Error-Custom "Failed to download $FileName : $_"
        return $null
    }
}

# Install Git
function Install-Git {
    if ($SkipGit) {
        Write-Info "Skipping Git installation"
        return
    }
    
    Write-Info "Installing Git..."
    $gitFile = Get-InstallFile "Git-2.43.0-64-bit.exe" "https://github.com/git-for-windows/git/releases/download/v2.43.0.windows.1/Git-2.43.0-64-bit.exe"
    
    if ($gitFile) {
        & $gitFile /VERYSILENT /NORESTART
        Write-Success "Git installed"
    }
}

# Install VSCode
function Install-VSCode {
    if ($SkipVSCode) {
        Write-Info "Skipping VSCode installation"
        return
    }
    
    Write-Info "Installing VSCode..."
    $vscodeFile = Get-InstallFile "VSCodeSetup-x64-1.84.2.exe" "https://aka.ms/win32-x64-user-stable"
    
    if ($vscodeFile) {
        & $vscodeFile /VERYSILENT /NORESTART
        Write-Success "VSCode installed"
    }
}

# Install Go
function Install-Go {
    if ($SkipGo) {
        Write-Info "Skipping Go installation"
        return
    }
    
    Write-Info "Installing Go..."
    $goFile = Get-InstallFile "go1.21.3.windows-amd64.msi" "https://go.dev/dl/go1.21.3.windows-amd64.msi"
    
    if ($goFile) {
        msiexec.exe /i $goFile /quiet /norestart
        Write-Success "Go installed"
    }
}

# Install Node.js
function Install-Node {
    if ($SkipNode) {
        Write-Info "Skipping Node.js installation"
        return
    }
    
    Write-Info "Installing Node.js..."
    $nodeFile = Get-InstallFile "node-v20.10.0-x64.msi" "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi"
    
    if ($nodeFile) {
        msiexec.exe /i $nodeFile /quiet /norestart
        Write-Success "Node.js installed"
    }
}

# Install PostgreSQL
function Install-PostgreSQL {
    if ($SkipPostgreSQL) {
        Write-Info "Skipping PostgreSQL installation"
        return
    }
    
    Write-Info "Installing PostgreSQL..."
    $pgFile = Get-InstallFile "postgresql-16.0-1-windows-x64.exe" "https://get.enterprisedb.com/postgresql/postgresql-16.0-1-windows-x64.exe"
    
    if ($pgFile) {
        & $pgFile --unattendedmodeui minimal --mode unattended --superpassword postgres
        Write-Success "PostgreSQL installed"
    }
}

# Install MySQL
function Install-MySQL {
    if ($SkipMySQL) {
        Write-Info "Skipping MySQL installation"
        return
    }
    
    Write-Info "Installing MySQL..."
    $mysqlFile = Get-InstallFile "mysql-8.1.0-winx64.msi" "https://dev.mysql.com/get/Downloads/MySQLInstaller/mysql-installer-community-8.1.0.0.msi"
    
    if ($mysqlFile) {
        msiexec.exe /i $mysqlFile /quiet /norestart
        Write-Success "MySQL installed"
    }
}

# Install Redis
function Install-Redis {
    if ($SkipRedis) {
        Write-Info "Skipping Redis installation"
        return
    }
    
    Write-Info "Installing Redis..."
    Write-Warning-Custom "Redis requires Windows Subsystem for Linux (WSL2) or third-party builds"
    Write-Info "Downloading Redis from memurai (Windows native build)..."
    
    $redisFile = Get-InstallFile "memurai-setup.exe" "https://github.com/microsoftarchive/redis/releases/download/win-3.0.504/Redis-x64-3.0.504.msi"
    
    if ($redisFile) {
        & $redisFile /VERYSILENT /NORESTART
        Write-Success "Redis installed"
    }
}

# Add paths to system PATH
function Add-SystemPaths {
    Write-Info "Adding custom paths to system PATH..."
    
    $pathsToAdd = @(
        "D:\Program Files\cmder",
        "D:\Program Files\sqlite",
        "D:\Program Files\ffmpeg\bin",
        "D:\Wnmp\mariadb-bins\default\bin",
        "D:\pgsql\bin",
        "D:\ollama"
    )
    
    $envPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    $pathModified = $false
    
    foreach ($path in $pathsToAdd) {
        if ($envPath -notlike "*$path*") {
            $envPath += ";$path"
            Write-Success "Added $path to system Path"
            $pathModified = $true
        } else {
            Write-Info "$path already exists in system Path"
        }
    }
    
    if ($pathModified) {
        [System.Environment]::SetEnvironmentVariable("Path", $envPath, "Machine")
        Write-Success "System PATH updated"
    }
}

# Configure power settings
function Configure-PowerSettings {
    Write-Info "Configuring power settings..."
    
    try {
        powercfg /change monitor-timeout-dc 180
        Write-Success "Set monitor timeout (battery) to 180 minutes"
        
        powercfg /change monitor-timeout-ac 180
        Write-Success "Set monitor timeout (AC) to 180 minutes"
        
        powercfg /change standby-timeout-dc 0
        Write-Success "Disabled sleep mode (battery)"
        
        powercfg /change standby-timeout-ac 0
        Write-Success "Disabled sleep mode (AC)"
    }
    catch {
        Write-Error-Custom "Failed to configure power settings: $_"
    }
}

# Main execution
function Main {
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Developer Tools Installation Script" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    
    if (-not (Test-Administrator)) {
        Write-Error-Custom "This script must be run as Administrator"
        exit 1
    }
    
    Write-Info "Local source path: $LocalSourcePath"
    Write-Info "Checking for local installers..."
    Write-Host ""
    
    if ($SkipAll) {
        Write-Warning-Custom "All installations skipped (SkipAll flag set)"
        return
    }
    
    Install-Git
    Install-VSCode
    Install-Go
    Install-Node
    Install-PostgreSQL
    Install-MySQL
    Install-Redis
    
    Write-Host ""
    Write-Info "Configuring system settings..."
    Write-Host ""
    
    Add-SystemPaths
    Configure-PowerSettings
    
    Write-Host ""
    Write-Success "Installation and configuration complete!"
    Write-Info "Some tools may require system restart to complete installation"
    Write-Info "Please restart your computer when ready"
}

Main
