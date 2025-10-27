<#
.SYNOPSIS
    Batch redirect Windows user folders (Desktop, Documents, Downloads, etc.) to specified locations.

.DESCRIPTION
    This script will:
    1. Auto-detect the current user's Downloads folder GUID
    2. Create target directories
    3. Update Windows Registry (User Shell Folders)
    4. Move existing files to new locations
    5. Refresh Windows Explorer

.PARAMETER TargetRoot
    Target root directory, default is D:\Personal

.PARAMETER NoMove
    Update registry only, do not move files

.EXAMPLE
    .\Set-PersonalFolders.ps1 -TargetRoot "D:\Personal"
    .\Set-PersonalFolders.ps1 -NoMove -WhatIf
#>

[CmdletBinding(SupportsShouldProcess=$true)]
param(
    [string]$TargetRoot = "D:\Personal",
    [switch]$NoMove
)

$ErrorActionPreference = 'Stop'

Write-Host "Personal Folders Redirection Tool" -ForegroundColor Green
Write-Host "Target Root: $TargetRoot" -ForegroundColor Gray

# Validate target drive
$targetDrive = ([IO.Path]::GetPathRoot($TargetRoot))
if (-not (Test-Path $targetDrive)) {
    Write-Error "Target drive '$targetDrive' does not exist. Please attach or create the drive before running."
    exit 1
}

$regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders"

$downloadsGuid = $null
try {
    $allKeys = Get-Item -Path $regPath | Select-Object -ExpandProperty Property
    foreach ($key in $allKeys) {
        $value = (Get-ItemProperty -Path $regPath -Name $key).$key
        $expandedValue = [Environment]::ExpandEnvironmentVariables($value)
        if ($expandedValue -match 'Downloads' -or $expandedValue -like '*下载*') {
            $downloadsGuid = $key
            break
        }
    }
} catch { }

# Folder mapping: Display name → Registry key name
$folders = @{
    "Desktop"      = "Desktop"
    "Documents"    = "Personal"
    "Music"        = "My Music"
    "Pictures"     = "My Pictures"
    "Videos"       = "My Video"
    "Favorites"    = "Favorites"
    "Contacts"     = "Contacts"
    "Links"        = "Links"
    "Searches"     = "Searches"
    "Saved Games"  = "SavedGames"
}

if ($downloadsGuid) {
    $folders["Downloads"] = $downloadsGuid
}

$customPaths = @{
    "Downloads" = "d:/Downloads"
}


# Ensure base directory exists
if (-not (Test-Path $TargetRoot)) {
    Write-Host "Creating base target directory: $TargetRoot" -ForegroundColor Yellow
    New-Item -Path $TargetRoot -ItemType Directory -Force | Out-Null
}

foreach ($name in $folders.Keys) {
    $regKey = $folders[$name]

    $regItem = Get-ItemProperty -Path $regPath -Name $regKey -ErrorAction SilentlyContinue
    if (-not $regItem) { continue }

    $rawValue = $regItem.$regKey
    if (-not $rawValue) { continue }

    $oldPath = [Environment]::ExpandEnvironmentVariables($rawValue)
    if ($customPaths.ContainsKey($name)) {
        $newPath = $customPaths[$name]
    } else {
        $newPath = Join-Path -Path $TargetRoot -ChildPath $name
    }

    Write-Host "`nProcessing [$name]" -ForegroundColor Cyan
    Write-Host "  Old Path: $oldPath"
    Write-Host "  New Path: $newPath"

    $oldNorm = ($oldPath.TrimEnd('\/'))
    $newNorm = ($newPath.TrimEnd('\/'))

    if (-not (Test-Path $newPath)) {
        New-Item -ItemType Directory -Force -Path $newPath | Out-Null
    }

    $doMove = $true
    if ($NoMove) { $doMove = $false }
    if ($oldNorm -ieq $newNorm) { $doMove = $false }

    if ($doMove -and (Test-Path $oldPath)) {
        try {
            $hasItems = (Get-ChildItem -LiteralPath $oldPath -Force -ErrorAction SilentlyContinue | Select-Object -First 1)
            if ($hasItems) {
                if ($PSCmdlet.ShouldProcess("$oldPath -> $newPath", "Move contents", "Move")) {
                    Write-Host "  Moving files..."
                    Move-Item -LiteralPath (Join-Path -Path $oldPath -ChildPath '*') -Destination $newPath -Force -ErrorAction Stop
                }
            } else {
                Write-Host "  Source is empty, skipping move."
            }
        } catch {
            Write-Warning "  Move failed: $_"
        }
    } else {
        Write-Host "  Skipping move."
    }

    try {
        if ($PSCmdlet.ShouldProcess("$regKey", "Set registry to $newPath", "Registry")) {
            Set-ItemProperty -Path $regPath -Name $regKey -Value $newPath -Type ExpandString -ErrorAction Stop
            Write-Host "  Registry updated: $regKey -> $newPath"
        }
    } catch {
        Write-Warning "  Failed to update registry: $_"
    }
}

Write-Host "`nRefreshing Windows Explorer settings..."
try {
    Start-Process -FilePath "RUNDLL32.EXE" -ArgumentList "USER32.DLL,UpdatePerUserSystemParameters" -WindowStyle Hidden -ErrorAction SilentlyContinue | Out-Null
} catch { }

Write-Host "`nCompleted. You may need to sign out or restart to fully apply changes." -ForegroundColor Green
