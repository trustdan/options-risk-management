@echo off
echo Options Trading Dashboard Build Script
echo ================================
echo.

powershell -ExecutionPolicy Bypass -File "%~dp0build.ps1"

echo.
echo Press any key to exit...
pause > nul 