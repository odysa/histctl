#!/usr/bin/env pwsh
# Install script for histctl on Windows

$ErrorActionPreference = "Stop"

$repo = "odysa/histctl"

# Detect architecture
$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "ARM64" { "arm64" }
    "AMD64" { "amd64" }
    default {
        Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
        exit 1
    }
}

$binary = "histctl-windows-${arch}.exe"
$url = "https://github.com/${repo}/releases/latest/download/${binary}"

# Install to LocalAppData bin directory
$installDir = Join-Path $env:LOCALAPPDATA "histctl"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

$dest = Join-Path $installDir "histctl.exe"

Write-Host "Downloading histctl for windows/${arch}..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing

# Add to PATH if not already present
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "Added $installDir to user PATH." -ForegroundColor Yellow
    Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}

Write-Host "Installed histctl to $dest" -ForegroundColor Green
