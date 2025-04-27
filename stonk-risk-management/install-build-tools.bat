@echo off
echo Options Trading Dashboard - Build Tools Installer
echo ================================================
echo.

:: Check for administrative privileges
net session >nul 2>&1
if %errorLevel% == 0 (
    echo Running with Administrator privileges
) else (
    echo This script needs to be run as Administrator to install software properly.
    echo Attempting to elevate privileges...
    echo.
    
    :: Attempt to run the script with elevated permissions and bypass execution policy
    powershell -Command "Start-Process cmd -ArgumentList '/c cd /d \"%~dp0\" && powershell -ExecutionPolicy Bypass -File \"%~dp0install-build-tools.ps1\"' -Verb RunAs" 
    exit /b
)

:: Run the PowerShell script with bypass execution policy
powershell -ExecutionPolicy Bypass -File "%~dp0install-build-tools.ps1" 