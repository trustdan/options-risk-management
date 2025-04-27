# Options Trading Dashboard - Build Tools Installer
# This script helps install necessary build tools like Inno Setup

Write-Host "Options Trading Dashboard - Build Tools Installer" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Check if script is running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "This script should be run as Administrator for best results." -ForegroundColor Yellow
    Write-Host "Some installations may fail or require manual intervention." -ForegroundColor Yellow
    Write-Host ""
    
    $continue = Read-Host "Continue anyway? (Y/N)"
    if ($continue -ne "Y" -and $continue -ne "y") {
        exit
    }
}

# Create a temp directory
$tempDir = Join-Path $env:TEMP "options-trading-dashboard-setup"
if (-not (Test-Path $tempDir)) {
    New-Item -Path $tempDir -ItemType Directory | Out-Null
}

# Function to check if a command exists
function Check-Command($cmdName) {
    return [bool](Get-Command -Name $cmdName -ErrorAction SilentlyContinue)
}

# 1. Check and install Go if needed
Write-Host "Checking for Go installation..." -ForegroundColor Cyan
if (Check-Command "go") {
    $goVersion = (go version) -replace "go version go([0-9]+\.[0-9]+).*", '$1'
    Write-Host "Go is already installed (version $goVersion)" -ForegroundColor Green
} else {
    Write-Host "Go is not installed. Installing..." -ForegroundColor Yellow
    
    # Download Go installer - updated to newer version
    $goInstallerUrl = "https://go.dev/dl/go1.21.6.windows-amd64.msi"
    $goInstallerPath = Join-Path $tempDir "go_installer.msi"
    
    try {
        Invoke-WebRequest -Uri $goInstallerUrl -OutFile $goInstallerPath
        Write-Host "Installing Go..."
        Start-Process -FilePath "msiexec.exe" -ArgumentList "/i", $goInstallerPath, "/qb" -Wait
        Write-Host "Go installed successfully" -ForegroundColor Green
        
        # Refresh PATH environment variable to make Go available immediately
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    }
    catch {
        Write-Host "Failed to download or install Go: $_" -ForegroundColor Red
        Write-Host "Please install Go manually from https://golang.org/dl/" -ForegroundColor Yellow
    }
}

# 2. Check and install Node.js if needed
Write-Host "Checking for Node.js installation..." -ForegroundColor Cyan
if (Check-Command "node") {
    $nodeVersion = (node -v)
    Write-Host "Node.js is already installed (version $nodeVersion)" -ForegroundColor Green
} else {
    Write-Host "Node.js is not installed. Installing..." -ForegroundColor Yellow
    
    # Download Node.js installer - updated to LTS version
    $nodeInstallerUrl = "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi"
    $nodeInstallerPath = Join-Path $tempDir "node_installer.msi"
    
    try {
        Invoke-WebRequest -Uri $nodeInstallerUrl -OutFile $nodeInstallerPath
        Write-Host "Installing Node.js..."
        Start-Process -FilePath "msiexec.exe" -ArgumentList "/i", $nodeInstallerPath, "/qb" -Wait
        Write-Host "Node.js installed successfully" -ForegroundColor Green
        
        # Refresh PATH environment variable to make Node.js available immediately
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    }
    catch {
        Write-Host "Failed to download or install Node.js: $_" -ForegroundColor Red
        Write-Host "Please install Node.js manually from https://nodejs.org/" -ForegroundColor Yellow
    }
}

# 3. Check and install Wails if needed
Write-Host "Checking for Wails installation..." -ForegroundColor Cyan
if (Check-Command "wails") {
    $wailsVersion = (wails version)
    Write-Host "Wails is already installed (version $wailsVersion)" -ForegroundColor Green
    
    # Update Wails to latest version
    Write-Host "Updating Wails to latest version..." -ForegroundColor Cyan
    try {
        Start-Process -FilePath "go" -ArgumentList "install", "github.com/wailsapp/wails/v2/cmd/wails@latest" -Wait -NoNewWindow
        $newWailsVersion = (wails version)
        Write-Host "Wails updated to version $newWailsVersion" -ForegroundColor Green
    }
    catch {
        Write-Host "Failed to update Wails: $_" -ForegroundColor Red
    }
} else {
    Write-Host "Wails is not installed. Installing..." -ForegroundColor Yellow
    
    try {
        Write-Host "Installing Wails via Go..."
        # Make sure GOPATH is in PATH
        $goPath = Join-Path $env:USERPROFILE "go\bin"
        if (-not ($env:Path -like "*$goPath*")) {
            $env:Path += ";$goPath"
        }
        
        Start-Process -FilePath "go" -ArgumentList "install", "github.com/wailsapp/wails/v2/cmd/wails@latest" -Wait -NoNewWindow
        
        # Verify installation
        if (Check-Command "wails") {
            $wailsVersion = (wails version)
            Write-Host "Wails installed successfully (version $wailsVersion)" -ForegroundColor Green
        } else {
            Write-Host "Wails installation completed but wails command not found in PATH" -ForegroundColor Yellow
            Write-Host "You may need to restart your terminal or add Go bin directory to your PATH" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "Failed to install Wails: $_" -ForegroundColor Red
        Write-Host "Please install Wails manually with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Yellow
    }
}

# 4. Check and install Inno Setup if needed
Write-Host "Checking for Inno Setup installation..." -ForegroundColor Cyan
$innoSetupCompilerPaths = @(
    "C:\Program Files (x86)\Inno Setup 6\ISCC.exe",
    "C:\Program Files\Inno Setup 6\ISCC.exe",
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
)

$innoSetupFound = $false
foreach ($path in $innoSetupCompilerPaths) {
    if (Test-Path $path) {
        Write-Host "Inno Setup is already installed at: $path" -ForegroundColor Green
        $innoSetupFound = $true
        break
    }
}

if (-not $innoSetupFound) {
    Write-Host "Inno Setup is not installed. Installing..." -ForegroundColor Yellow
    
    # Download Inno Setup installer - using direct link
    $innoSetupUrl = "https://files.jrsoftware.org/is/6/innosetup-6.2.2.exe"
    $innoSetupPath = Join-Path $tempDir "inno_setup_installer.exe"
    
    try {
        Invoke-WebRequest -Uri $innoSetupUrl -OutFile $innoSetupPath
        Write-Host "Installing Inno Setup..."
        Start-Process -FilePath $innoSetupPath -ArgumentList "/VERYSILENT", "/SUPPRESSMSGBOXES", "/NORESTART", "/NOCANCEL", "/SP-", "/CLOSEAPPLICATIONS", "/RESTARTAPPLICATIONS", "/NOICONS", "/TYPE=full" -Wait
        Write-Host "Inno Setup installed successfully" -ForegroundColor Green
    }
    catch {
        Write-Host "Failed to download or install Inno Setup: $_" -ForegroundColor Red
        Write-Host "Please install Inno Setup manually from https://jrsoftware.org/isinfo.php" -ForegroundColor Yellow
    }
}

# Clean up
Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Setup completed!" -ForegroundColor Green
Write-Host "You may need to restart your terminal or computer for some changes to take effect." -ForegroundColor Cyan
Write-Host "After restarting, you can build the application using build.bat" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") 