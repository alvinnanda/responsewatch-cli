# Install script for ResponseWatch CLI (Windows)
# Usage: iwr -useb https://response-watch.web.app/install.ps1 | iex

$ErrorActionPreference = 'Stop'

$APP_NAME = "rwcli"
$INSTALL_DIR = "$env:LOCALAPPDATA\Programs"
$VERSION = "latest"

# Detect architecture
$ARCH = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $ARCH = "arm64"
}

$OS = "windows"
$BINARY_NAME = "${APP_NAME}_${OS}_${ARCH}.exe"

Write-Host "Installing $APP_NAME for Windows/$ARCH..." -ForegroundColor Cyan

# Create install directory
if (!(Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Download URL
$DOWNLOAD_URL = "https://response-watch.web.app/cli/$BINARY_NAME"

Write-Host "Downloading from $DOWNLOAD_URL..." -ForegroundColor Yellow

$TMP_FILE = "$env:TEMP\$BINARY_NAME"

try {
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TMP_FILE -UseBasicParsing
} catch {
    Write-Host "Error: Failed to download $DOWNLOAD_URL" -ForegroundColor Red
    Write-Host "Error details: $_" -ForegroundColor Red
    exit 1
}

# Move to install directory
$INSTALL_PATH = "$INSTALL_DIR\$APP_NAME.exe"
Move-Item -Path $TMP_FILE -Destination $INSTALL_PATH -Force

Write-Host "Installed to $INSTALL_PATH" -ForegroundColor Green

# Add to PATH if not already there
$USER_PATH = [Environment]::GetEnvironmentVariable("Path", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
    Write-Host "Adding to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$USER_PATH;$INSTALL_DIR", "User")
    Write-Host "Please restart your terminal or run 'refreshenv' to update PATH" -ForegroundColor Yellow
}

Write-Host "" -ForegroundColor Green
Write-Host "✓ Installation complete!" -ForegroundColor Green
Write-Host "Run '$APP_NAME --help' to get started" -ForegroundColor Green
