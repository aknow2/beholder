# Beholder Installation Script for Windows
# Usage: Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression

$ErrorActionPreference = "Stop"

# Configuration
$RepoOwner = "aknow2"
$RepoName = "beholder"
$InstallDir = "$env:USERPROFILE\.beholder\bin"
$DataDir = "$env:USERPROFILE\.beholder"
$ManifestFile = "$DataDir\install-manifest.txt"

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

# Detect platform (Windows only)
function Test-Platform {
    Write-Info "Detected platform: Windows (amd64)"
    $script:Platform = "windows"
    $script:Arch = "amd64"
}

# Check for existing installation
function Test-ExistingInstallation {
    $script:ExistingVersion = ""
    $script:UpgradeMode = $false
    
    if (Test-Path $ManifestFile) {
        # Read version from manifest
        $manifest = Get-Content $ManifestFile -Raw
        $script:ExistingVersion = ($manifest | Select-String "INSTALLED_VERSION=(.+)").Matches.Groups[1].Value
    }
    elseif (Get-Command beholder -ErrorAction SilentlyContinue) {
        # Try to get version from binary
        try {
            $versionOutput = & beholder --version 2>&1
            $script:ExistingVersion = ($versionOutput | Select-String "v[0-9.]+").Matches.Value
            if (-not $ExistingVersion) { $script:ExistingVersion = "unknown" }
        }
        catch {
            $script:ExistingVersion = "unknown"
        }
    }
    
    if ($ExistingVersion) {
        Write-Info "Existing installation found: $ExistingVersion"
        $script:UpgradeMode = $true
    }
}

# Get latest version from GitHub API
function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $apiUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
        $response = Invoke-RestMethod -Uri $apiUrl -UseBasicParsing
        $script:Version = $response.tag_name
        
        if (-not $Version) {
            Write-Error "Failed to fetch latest version"
        }
        
        Write-Info "Latest version: $Version"
    }
    catch {
        Write-Error "Failed to fetch version: $_"
    }
}

# Download binary with retry logic
function Get-Binary {
    $binaryName = "beholder-$Version-$Platform-$Arch.exe"
    $downloadUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/$Version/$binaryName"
    $script:TempFile = "$env:TEMP\beholder-$PID.exe"
    
    Write-Info "Downloading from: $downloadUrl"
    
    $maxRetries = 3
    for ($i = 1; $i -le $maxRetries; $i++) {
        try {
            Invoke-WebRequest -Uri $downloadUrl -OutFile $TempFile -UseBasicParsing
            Write-Success "Download complete"
            return
        }
        catch {
            if ($i -lt $maxRetries) {
                Write-Warn "Download failed, retrying ($i/$maxRetries)..."
                Start-Sleep -Seconds 2
            }
            else {
                Write-Error "Failed to download binary after $maxRetries attempts: $_"
            }
        }
    }
}

# Validate downloaded binary
function Test-Binary {
    Write-Info "Validating binary..."
    
    if (-not (Test-Path $TempFile)) {
        Write-Error "Binary file not found: $TempFile"
    }
    
    $fileSize = (Get-Item $TempFile).Length
    if ($fileSize -lt 1000000) {
        Write-Error "Binary file is too small ($fileSize bytes), download may be corrupted"
    }
    
    try {
        $versionOutput = & $TempFile --version 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Binary validation failed: cannot execute --version"
        }
        Write-Success "Binary validated"
    }
    catch {
        Write-Error "Binary validation failed: $_"
    }
}

# Install binary to target directory
function Install-Binary {
    Write-Info "Installing to $InstallDir..."
    
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    $targetPath = "$InstallDir\beholder.exe"
    
    if (Test-Path $targetPath) {
        if ($UpgradeMode) {
            Write-Info "Upgrading from $ExistingVersion to $Version..."
        }
        else {
            Write-Warn "Existing installation found, replacing..."
        }
        Remove-Item $targetPath -Force
    }
    
    Move-Item $TempFile $targetPath -Force
    
    if ($UpgradeMode) {
        Write-Success "Upgrade complete"
    }
    else {
        Write-Success "Binary installed"
    }
}

# Update user PATH environment variable
function Update-Path {
    Write-Info "Updating PATH..."
    
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($currentPath -notlike "*$InstallDir*") {
        $newPath = "$currentPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        
        # Update current session PATH
        $env:Path = "$env:Path;$InstallDir"
        
        Write-Success "Added $InstallDir to user PATH"
        Write-Warn "Restart your terminal for PATH changes to take full effect"
    }
    else {
        Write-Info "PATH already includes $InstallDir"
    }
}

# Create installation manifest
function New-Manifest {
    Write-Info "Creating installation manifest..."
    
    if (-not (Test-Path $DataDir)) {
        New-Item -ItemType Directory -Path $DataDir -Force | Out-Null
    }
    
    $timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
    
    $manifestContent = @"
INSTALLED_VERSION=$Version
INSTALLATION_PATH=$InstallDir\beholder.exe
INSTALLATION_DATE=$timestamp
USER_DATA_PATH=$DataDir
"@
    
    Set-Content -Path $ManifestFile -Value $manifestContent
    
    Write-Success "Manifest created at $ManifestFile"
}

# Main installation flow
function Main {
    Write-Host ""
    Write-Host "Beholder Installer" -ForegroundColor Cyan
    Write-Host "==================" -ForegroundColor Cyan
    Write-Host ""
    
    Test-Platform
    Test-ExistingInstallation
    Get-LatestVersion
    
    # Version comparison feedback
    if ($UpgradeMode) {
        if ($ExistingVersion -eq $Version) {
            Write-Info "Already at latest version ($Version), reinstalling..."
        }
        elseif ($ExistingVersion -eq "unknown") {
            Write-Info "Upgrading to $Version..."
        }
        else {
            Write-Info "Upgrade available: $ExistingVersion → $Version"
        }
    }
    
    Get-Binary
    Test-Binary
    Install-Binary
    Update-Path
    New-Manifest
    
    Write-Host ""
    Write-Success "Installation complete!"
    Write-Host ""
    Write-Info "Verify installation: beholder --version"
    Write-Info "Get started: beholder help"
    Write-Host ""
}

Main
