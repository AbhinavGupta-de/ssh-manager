# SSH Manager Installation Script for Windows
# Run with: PowerShell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

Write-Host "🔑 SSH Manager Installation" -ForegroundColor Blue
Write-Host "================================" -ForegroundColor Blue

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✓ Go is installed ($goVersion)" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed!" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Build the binary
Write-Host "`nBuilding SSH Manager..." -ForegroundColor Yellow
try {
    go build -o sshm.exe main.go
    Write-Host "✓ Build successful" -ForegroundColor Green
} catch {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}

# Determine installation directory
$installDir = "$env:USERPROFILE\bin"

# Create bin directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "✓ Created directory: $installDir" -ForegroundColor Green
}

# Copy the binary
Write-Host "`nInstalling to $installDir..." -ForegroundColor Yellow
try {
    Copy-Item sshm.exe "$installDir\sshm.exe" -Force
    Write-Host "✓ Installed to $installDir\sshm.exe" -ForegroundColor Green
} catch {
    Write-Host "❌ Installation failed!" -ForegroundColor Red
    exit 1
}

# Check if directory is in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    Write-Host "`n⚠️  $installDir is not in your PATH" -ForegroundColor Yellow
    Write-Host "Adding to PATH..." -ForegroundColor Blue
    
    try {
        $newPath = $currentPath + ";$installDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "✓ Updated PATH" -ForegroundColor Green
        Write-Host "⚠️  Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
    } catch {
        Write-Host "❌ Failed to update PATH" -ForegroundColor Red
        Write-Host "Please manually add $installDir to your PATH" -ForegroundColor Yellow
    }
} else {
    Write-Host "✓ $installDir is already in PATH" -ForegroundColor Green
}

# Clean up build artifact
Remove-Item sshm.exe -ErrorAction SilentlyContinue

Write-Host "`n✅ Installation complete!" -ForegroundColor Green
Write-Host "`nUsage:" -ForegroundColor Blue
Write-Host "  sshm new              # Create a new profile"
Write-Host "  sshm list             # List all profiles"
Write-Host "  sshm switch <n>    # Switch profiles"
Write-Host "  sshm current          # Show current profile"
Write-Host ""
Write-Host "💡 Try: sshm new" -ForegroundColor Yellow
Write-Host "`n⚠️  Remember to restart your terminal!" -ForegroundColor Yellow