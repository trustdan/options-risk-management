#!/bin/bash

# Options Trading Dashboard - Complete macOS Build Script
# This script builds the application and creates a DMG with proper appearance settings
# No external dependencies required beyond standard macOS tools

# Stop on any error
set -e

# Configuration - using the name from wails.json
APP_NAME="Options Trading Dashboard"
VOLUME_NAME="${APP_NAME}_$$"  # Add process ID to make unique

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

# Step 1: Create a background image if Pillow is available
BACKGROUND_IMAGE="$OUTPUT_DIR/resources/dmg-background.png"
mkdir -p "$OUTPUT_DIR/resources"

if command -v python3 &> /dev/null; then
    echo "Checking for Python Pillow library..."
    
    # Check if Pillow is installed
    if python3 -c "import PIL" &> /dev/null; then
        echo "Creating background image using Python..."
        python3 create_background.py
    else
        echo "Python Pillow library not found. DMG will be created without a custom background."
    fi
else
    echo "Python not found. DMG will be created without a custom background."
fi

# Step 2: Build the Wails application for macOS
echo "Building Options Trading Dashboard application for macOS..."

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "Wails CLI not found. Please install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Run wails build with macOS configuration
wails build -platform darwin/universal -skipbindings
echo "Application built successfully!"

# Step 3: Create the DMG using hdiutil (built-in macOS tool)
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

# Unmount any existing volumes with similar names
echo "Cleaning up any existing mounts..."
df | grep "$APP_NAME" | awk '{print $1}' | while read -r device; do
    hdiutil detach "$device" -force 2>/dev/null || true
done

# Create the temporary DMG
echo "Creating temporary DMG file..."
# Using VOLUME_NAME defined at the beginning of the script
SIZE=$(du -sm "$MOUNT_POINT" | cut -f1)
SIZE=$((SIZE + 50))  # Add 50MB for safety

hdiutil create -verbose -size ${SIZE}m -volname "$VOLUME_NAME" -fs HFS+ \
    -srcfolder "$MOUNT_POINT" -ov -format UDRW "$TMP_DMG"

# Mount the temporary DMG
echo "Mounting DMG for customization..."
MOUNT_OUTPUT=$(hdiutil attach -verbose -readwrite -noverify -mountpoint "/Volumes/$VOLUME_NAME" "$TMP_DMG")
MOUNT_DEV=$(echo "$MOUNT_OUTPUT" | grep "Apple_HFS" | awk '{print $1}')
MOUNT_VOL="/Volumes/$VOLUME_NAME"

# Wait for mount and verify
sleep 2
echo "Mount device: $MOUNT_DEV"
echo "Mount point: $MOUNT_VOL"
ls -la "$MOUNT_VOL"

if [ ! -d "$MOUNT_VOL" ]; then
    echo "Error: Failed to mount DMG at $MOUNT_VOL"
    exit 1
fi

# Test write permissions with sudo
if ! sudo touch "$MOUNT_VOL/.test_write" 2>/dev/null; then
    echo "Error: Unable to write to mounted DMG even with elevated permissions"
    sudo hdiutil detach "$MOUNT_DEV" -force 2>/dev/null || true
    exit 1
fi
sudo rm -f "$MOUNT_VOL/.test_write"

# Set the DMG appearance (window location, size, background, etc.)
echo "Configuring DMG appearance..."

# Apply background if it exists
if [ -f "$BACKGROUND_IMAGE" ]; then
    echo "Adding background image to DMG..."
    
    # Copy background image to volume (with sudo for permissions)
    sudo mkdir -p "$MOUNT_VOL/.background"
    sudo cp "$BACKGROUND_IMAGE" "$MOUNT_VOL/.background/background.png"
    sudo chmod 755 "$MOUNT_VOL/.background"
    
    # Set background image in Finder preferences
    BACKGROUND_SCRIPT="tell application \"Finder\"
        tell disk \"$VOLUME_NAME\"
            open
            set current view of container window to icon view
            set toolbar visible of container window to false
            set statusbar visible of container window to false
            set the bounds of container window to {400, 100, 1200, 500}
            set theViewOptions to the icon view options of container window
            set background picture of theViewOptions to file \".background:background.png\"
            set arrangement of theViewOptions to not arranged
            set icon size of theViewOptions to 128
            set position of item \"$EXECUTABLE_NAME.app\" of container window to {180, 200}
            set position of item \"Applications\" of container window to {480, 200}
            set position of item \"README.txt\" of container window to {330, 320}
            close
            open
            update without registering applications
            delay 5
            close
        end tell
    end tell"
    
    echo "$BACKGROUND_SCRIPT" | osascript
else
    # Create .DS_Store file with appearance settings (no background)
    echo "Setting Finder window position and size (no background)..."
    osascript << EOT
    tell application "Finder"
        tell disk "$VOLUME_NAME"
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
fi

# Finalize and convert to compressed readonly DMG
echo "Finalizing DMG..."
sync
sleep 2

# Ensure proper unmounting
echo "Unmounting DMG..."
sudo hdiutil detach "$MOUNT_DEV" -force || {
    echo "Failed to unmount $MOUNT_DEV, trying alternative detachment..."
    diskutil unmount force "$MOUNT_DEV" || true
    sleep 2
}

# Convert to final DMG
hdiutil convert "$TMP_DMG" -format UDZO -o "$DMG_FILE"
rm -f "$TMP_DMG"

# Clean up
rm -rf "$TMP_DIR"

# Check if DMG was created successfully
if [ -f "$DMG_FILE" ]; then
    echo "DMG package created successfully at: $DMG_FILE"
    echo "Open the DMG file to verify its appearance and functionality."
else
    echo "Failed to create DMG package."
    exit 1
fi

echo "Build process completed!" 