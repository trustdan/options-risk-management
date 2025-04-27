@echo off
echo Options Risk Management Build Script
echo ================================
echo.

powershell -ExecutionPolicy Bypass -File "%~dp0build.ps1"

echo.
echo Press any key to exit...
pause > nul 