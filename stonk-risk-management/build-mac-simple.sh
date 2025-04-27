#!/bin/bash

# Options Trading Dashboard - Simple macOS Build Script
# This script builds the application and creates a basic DMG using built-in macOS tools

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

# Step 2: Create the DMG using hdiutil (built-in macOS tool)
echo "Creating DMG package using hdiutil..."

# The app is located in build/bin, but in a macOS app package
APP_FILE="$OUTPUT_DIR/bin/$EXECUTABLE_NAME.app"
DMG_FILE="$OUTPUT_DIR/$DMG_NAME.dmg"
TMP_DMG="$OUTPUT_DIR/tmp-$DMG_NAME.dmg"

# Remove any existing DMG
if [ -f "$DMG_FILE" ]; then
    echo "Removing existing DMG file..."
    rm "$DMG_FILE"
fi

if [ -f "$TMP_DMG" ]; then
    rm "$TMP_DMG"
fi

# Create a temporary directory for DMG contents
TMP_DIR=$(mktemp -d)
MOUNT_POINT="$TMP_DIR/volume"
mkdir -p "$MOUNT_POINT"

# Create a README file that instructs users how to install
cat > "$TMP_DIR/README.txt" << EOL
# Options Trading Dashboard

To install:
1. Drag the app icon to the Applications folder
2. Eject the disk image

That's it!
EOL

# Copy the app and README to the temporary location
cp -R "$APP_FILE" "$MOUNT_POINT/"
cp "$TMP_DIR/README.txt" "$MOUNT_POINT/"

# Create a symbolic link to Applications folder
echo "Creating Applications folder shortcut..."
ln -s /Applications "$MOUNT_POINT/Applications"

# Create the temporary DMG
echo "Creating DMG file..."
hdiutil create -volname "$APP_NAME" -srcfolder "$MOUNT_POINT" -ov -format UDRW "$TMP_DMG"

# Mount the temporary DMG
MOUNT_DEV=$(hdiutil attach -readwrite -noverify -noautoopen "$TMP_DMG" | head -n 1 | awk '{print $1}')
sleep 2

# Set the DMG appearance (window location, size, background, etc.)
echo "Configuring DMG appearance..."
MOUNT_VOL="/Volumes/$APP_NAME"

# Create .DS_Store file with appearance settings
echo "Setting Finder window position and size..."
osascript << EOT
tell application "Finder"
    tell disk "$APP_NAME"
        open
        set current view of container window to icon view
        set toolbar visible of container window to false
        set statusbar visible of container window to false
        set the bounds of container window to {400, 100, 920, 440}
        set theViewOptions to the icon view options of container window
        set arrangement of theViewOptions to not arranged
        set icon size of theViewOptions to 128
        set position of item "$EXECUTABLE_NAME.app" of container window to {160, 180}
        set position of item "Applications" of container window to {360, 180}
        set position of item "README.txt" of container window to {260, 300}
        close
        open
        update without registering applications
        delay 5
        close
    end tell
end tell
EOT

# Finalize and convert to compressed readonly DMG
echo "Finalizing DMG..."
sync
sleep 2
hdiutil detach "$MOUNT_DEV"
hdiutil convert "$TMP_DMG" -format UDZO -o "$DMG_FILE"
rm -f "$TMP_DMG"

# Clean up
rm -rf "$TMP_DIR"

# Check if DMG was created successfully
if [ -f "$DMG_FILE" ]; then
    echo "DMG package created successfully at: $DMG_FILE"
else
    echo "Failed to create DMG package."
    exit 1
fi

echo "Build process completed!" 