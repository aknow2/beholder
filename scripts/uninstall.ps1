# Beholder Uninstallation Script for Windows
# Usage: Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.ps1 -UseBasicParsing | Invoke-Expression

$ErrorActionPreference = "Stop"

# Configuration
$DataDir = "$env:USERPROFILE\.beholder"
$ManifestFile = "$DataDir\install-manifest.txt"
$FallbackInstallDir = "$env:USERPROFILE\.beholder\bin"

# Helper functions
function Write-Info {
    param([string]$Message)
    Write-Host "→ $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
    exit 1
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Ask {
    param([string]$Message)
    Write-Host "? $Message" -ForegroundColor Yellow -NoNewline
}

# Read installation manifest
function Read-Manifest {
    if (Test-Path $ManifestFile) {
        Write-Info "Reading installation manifest..."
        
        $manifest = Get-Content $ManifestFile -Raw
        $script:InstalledVersion = ($manifest | Select-String "INSTALLED_VERSION=(.+)").Matches.Groups[1].Value
        $script:InstallationPath = ($manifest | Select-String "INSTALLATION_PATH=(.+)").Matches.Groups[1].Value
        $script:UserDataPath = ($manifest | Select-String "USER_DATA_PATH=(.+)").Matches.Groups[1].Value
        
        if ($InstallationPath) {
            Write-Info "Found installation: $InstallationPath"
            return
        }
    }
    
    # Fallback: try to find beholder in PATH
    Write-Warn "Installation manifest not found, searching for beholder..."
    
    $beholderCmd = Get-Command beholder -ErrorAction SilentlyContinue
    if ($beholderCmd) {
        $script:InstallationPath = $beholderCmd.Source
        Write-Info "Found beholder at: $InstallationPath"
        $script:UserDataPath = $DataDir
        return
    }
    
    Write-Error "Could not locate beholder installation. Please specify the path manually."
}

# Remove binary
function Remove-Binary {
    if (-not $InstallationPath) {
        Write-Warn "Installation path not specified, skipping binary removal"
        return
    }
    
    if (-not (Test-Path $InstallationPath)) {
        Write-Warn "Binary not found at $InstallationPath, may already be removed"
        return
    }
    
    Write-Info "Removing binary from $InstallationPath..."
    
    try {
        Remove-Item $InstallationPath -Force
        Write-Success "Binary removed"
    }
    catch {
        Write-Error "Failed to remove binary: $_"
    }
}

# Clean PATH from user environment
function Clear-Path {
    Write-Info "Cleaning PATH from user environment..."
    
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    # Check for beholder installation directory
    $pathToRemove = Split-Path $InstallationPath -Parent
    
    if ($currentPath -like "*$pathToRemove*") {
        $pathArray = $currentPath -split ';'
        $newPath = ($pathArray | Where-Object { $_ -ne $pathToRemove -and $_ -ne "$pathToRemove\" }) -join ';'
        
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        
        # Update current session PATH
        $env:Path = $env:Path -replace [regex]::Escape("$pathToRemove;"), ""
        
        Write-Success "Cleaned PATH from user environment"
    }
    else {
        Write-Info "No beholder PATH entry found in user environment"
    }
}

# Handle user data
function Remove-UserData {
    if (-not (Test-Path $UserDataPath)) {
        Write-Info "No user data found at $UserDataPath"
        return
    }
    
    # Check if running interactively
    if ([Environment]::UserInteractive -and -not $NonInteractive) {
        Write-Host ""
        Write-Ask "Remove user data in $UserDataPath? (y/N): "
        $response = Read-Host
        Write-Host ""
    }
    else {
        # Non-interactive: preserve data by default
        $response = "n"
        Write-Info "Running non-interactively, preserving user data"
    }
    
    switch -Regex ($response) {
        '^[Yy]' {
            Write-Info "Removing user data..."
            Remove-Item $UserDataPath -Recurse -Force
            Write-Success "User data removed"
        }
        default {
            Write-Info "User data preserved at $UserDataPath"
            Write-Warn "To remove manually later: Remove-Item -Recurse -Force $UserDataPath"
        }
    }
}

# Main uninstallation flow
function Main {
    Write-Host ""
    Write-Host "Beholder Uninstaller" -ForegroundColor Cyan
    Write-Host "====================" -ForegroundColor Cyan
    Write-Host ""
    
    Read-Manifest
    Remove-Binary
    Clear-Path
    Remove-UserData
    
    Write-Host ""
    Write-Success "Uninstallation complete!"
    Write-Host ""
    Write-Info "To verify removal: Get-Command beholder"
    if (Test-Path $UserDataPath) {
        Write-Info "User data preserved at: $UserDataPath"
    }
    Write-Host ""
}

Main
