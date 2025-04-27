#!/bin/bash

# Options Trading Dashboard - macOS Build Script
# This script builds both the application and creates a DMG

# Stop on any error
set -e

# Configuration - using the name from wails.json
APP_NAME="Options Trading Dashboard"

# Read the version from the version.js file
VERSION_JS_PATH="frontend/src/version.js"
if [ -f "$VERSION_JS_PATH" ]; then
    APP_VERSION=$(grep -o 'VERSION = "[^"]*"' "$VERSION_JS_PATH" | cut -d'"' -f2)
    echo "Using version $APP_VERSION from version.js"
else
    APP_VERSION="1.0.0"
    echo "version.js not found, using default version: $APP_VERSION"
fi

OUTPUT_DIR="build"
EXECUTABLE_NAME="options-trading-dashboard"
DMG_NAME="options-trading-dashboard-$APP_VERSION"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"
echo "Output directory: $OUTPUT_DIR"

# Step 1: Build the Wails application for macOS
echo "Building Options Trading Dashboard application for macOS..."

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "Wails CLI not found. Please install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Run wails build with macOS configuration
wails build -platform darwin/universal -skipbindings
echo "Application built successfully!"

# Step 2: Check if create-dmg is installed
if ! command -v create-dmg &> /dev/null; then
    echo "create-dmg not found. Installing with Homebrew..."
    if ! command -v brew &> /dev/null; then
        echo "Homebrew not found. Please install Homebrew first: https://brew.sh/"
        exit 1
    fi
    brew install create-dmg
fi

# Step 3: Create the DMG
echo "Creating DMG package..."

# The app is located in build/bin, but in a macOS app package
APP_FILE="$OUTPUT_DIR/bin/$EXECUTABLE_NAME.app"
DMG_FILE="$OUTPUT_DIR/$DMG_NAME.dmg"

# Remove any existing DMG
if [ -f "$DMG_FILE" ]; then
    echo "Removing existing DMG file..."
    rm "$DMG_FILE"
fi

# Create the DMG
create-dmg \
    --volname "$APP_NAME" \
    --volicon "$APP_FILE/Contents/Resources/iconfile.icns" \
    --background "build/resources/dmg-background.png" \
    --window-pos 200 120 \
    --window-size 800 400 \
    --icon-size 100 \
    --icon "$EXECUTABLE_NAME.app" 200 190 \
    --hide-extension "$EXECUTABLE_NAME.app" \
    --app-drop-link 600 185 \
    "$DMG_FILE" \
    "$APP_FILE" || {
        echo "Error: DMG creation failed. Attempting simplified DMG creation..."
        create-dmg \
            --volname "$APP_NAME" \
            --volicon "$APP_FILE/Contents/Resources/iconfile.icns" \
            --window-pos 200 120 \
            --window-size 800 400 \
            --icon-size 100 \
            --icon "$EXECUTABLE_NAME.app" 200 190 \
            --hide-extension "$EXECUTABLE_NAME.app" \
            "$DMG_FILE" \
            "$APP_FILE" || {
                echo "Simplified DMG creation also failed. Creating minimal DMG..."
                create-dmg \
                    --volname "$APP_NAME" \
                    "$DMG_FILE" \
                    "$APP_FILE"
            }
    }

# Check if DMG was created successfully
if [ -f "$DMG_FILE" ]; then
    echo "DMG package created successfully at: $DMG_FILE"
else
    echo "Failed to create DMG package."
    exit 1
fi

echo "Build process completed!" 