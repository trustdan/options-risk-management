# Options Trading Dashboard Build Script
# This script builds both the application and installer

# Stop on any error
$ErrorActionPreference = "Stop"

# Configuration - using the name from wails.json
$appName = "Options Trading Dashboard"
$appVersion = "1.0.0"
$outputDir = "build"
$installerDir = "installer"
$executableName = "options-trading-dashboard" # Changed to match wails.json outputfilename
$installerName = "options-trading-dashboard" # Changed installer name as requested

# Create output directories if they don't exist
if (-not (Test-Path $outputDir)) {
    New-Item -Path $outputDir -ItemType Directory | Out-Null
    Write-Host "Created build directory at $outputDir"
}

if (-not (Test-Path $installerDir)) {
    New-Item -Path $installerDir -ItemType Directory | Out-Null
    Write-Host "Created installer directory at $installerDir"
}

# Step 1: Build the Wails application
Write-Host "Building Options Trading Dashboard application..." -ForegroundColor Cyan

# Check if wails is installed
if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
    Write-Host "Wails CLI not found. Please install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Red
    Write-Host "Alternatively, run the install-build-tools.bat script to install all dependencies." -ForegroundColor Yellow
    exit 1
}

# Run wails build
try {
    # Add -skipbindings flag to fix common build issues
    wails build -skipbindings
    Write-Host "Application built successfully!" -ForegroundColor Green
}
catch {
    Write-Host "Failed to build application: $_" -ForegroundColor Red
    exit 1
}

# Step 2: Create Inno Setup script
$innoSetupScript = "$installerDir\installer.iss"

Write-Host "Creating Inno Setup script..." -ForegroundColor Cyan

$innoContent = @"
#define MyAppName "$appName"
#define MyAppVersion "$appVersion"
#define MyAppPublisher "Options Trading Dashboard"
#define MyAppURL "https://github.com/yourusername/options-trading-dashboard"
#define MyAppExeName "$executableName.exe"

[Setup]
AppId={{7D49E80F-96A9-4AC7-B7C9-B95CD2219C3A}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
LicenseFile=..\LICENSE
OutputDir=..\$installerDir
OutputBaseFilename=$installerName-setup-{#MyAppVersion}
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "..\build\bin\$executableName.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion
; Add any additional files or directories here

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent
"@

# Write the Inno Setup script file
$innoContent | Out-File -FilePath $innoSetupScript -Encoding utf8

# Step 3: Compile the installer with Inno Setup
Write-Host "Creating installer..." -ForegroundColor Cyan

# Check multiple potential Inno Setup installation paths
$innoSetupCompilerPaths = @(
    "C:\Program Files (x86)\Inno Setup 6\ISCC.exe",
    "C:\Program Files\Inno Setup 6\ISCC.exe",
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
)

$innoSetupCompiler = $null
foreach ($path in $innoSetupCompilerPaths) {
    if (Test-Path $path) {
        $innoSetupCompiler = $path
        break
    }
}

if (-not $innoSetupCompiler) {
    Write-Host "Inno Setup compiler not found in standard locations." -ForegroundColor Yellow
    Write-Host "Please install Inno Setup from: https://jrsoftware.org/isinfo.php" -ForegroundColor Yellow
    Write-Host "Or run the install-build-tools.bat script included with this project." -ForegroundColor Yellow
    
    $userChoice = Read-Host "Do you want to continue without building the installer? (Y/N)"
    if ($userChoice -ne "Y" -and $userChoice -ne "y") {
        exit 1
    }
} else {
    # Ensure the installer directory exists and is empty for a clean build
    if (Test-Path "$installerDir\$installerName-setup-$appVersion.exe") {
        Write-Host "Removing existing installer..." -ForegroundColor Yellow
        Remove-Item "$installerDir\$installerName-setup-$appVersion.exe" -Force
    }
    
    # Compile the installer
    try {
        Write-Host "Found Inno Setup compiler at: $innoSetupCompiler" -ForegroundColor Green
        & $innoSetupCompiler $innoSetupScript
        
        if ($LASTEXITCODE -eq 0) {
            # Verify the installer was created
            if (Test-Path "$installerDir\$installerName-setup-$appVersion.exe") {
                Write-Host "Installer created successfully at: $installerDir\$installerName-setup-$appVersion.exe" -ForegroundColor Green
            } else {
                Write-Host "Warning: Installer compilation reported success but the file was not found at the expected location." -ForegroundColor Yellow
            }
        } else {
            Write-Host "Failed to create installer." -ForegroundColor Red
            exit 1
        }
    }
    catch {
        Write-Host "Error running Inno Setup compiler: $_" -ForegroundColor Red
        exit 1
    }
}

Write-Host "Build process completed!" -ForegroundColor Green 