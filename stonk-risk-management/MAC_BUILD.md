# Building for macOS

This document explains how to build the Options Trading Dashboard application for macOS.

## Prerequisites

1. **Wails CLI**: The application uses Wails to create a desktop application
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **Go**: Make sure you have Go 1.18 or higher installed
   ```bash
   brew install go
   ```

3. **Node.js**: Required for the frontend build (v16 or higher recommended)
   ```bash
   brew install node
   ```

## Build Options

There are three build scripts available for macOS:

### Option 1: Complete Build (Recommended)

This option uses built-in macOS tools to create a DMG with proper appearance and can use Python to generate a background image if available.

```bash
cd stonk-risk-management
./build-mac-complete.sh
```

For the best experience with a custom background image, install the Python Pillow library:
```bash
pip3 install Pillow
```

The script will:
1. Generate a background image using Python (if Pillow is installed)
2. Build the Wails application for macOS
3. Create a DMG file with proper layout and icon positioning
4. Apply the background image if available
5. Place the resulting DMG in the `build` directory

### Option 2: Simple Build

This option uses macOS built-in tools to create a basic DMG without requiring additional dependencies.

```bash
cd stonk-risk-management
./build-mac-simple.sh
```

The script will:
1. Build the Wails application for macOS
2. Create a DMG file using hdiutil with a basic layout
3. Place the resulting DMG in the `build` directory

### Option 3: Enhanced Build (Using create-dmg)

This option creates a more polished DMG with customized appearance but requires the `create-dmg` tool.

1. Install the required tool:
   ```bash
   brew install create-dmg
   ```

2. Run the build script:
   ```bash
   cd stonk-risk-management
   ./build-mac.sh
   ```

## Output

All scripts will produce a DMG file in the `build` directory with a name like:
```
options-trading-dashboard-0.1.7.dmg
```

The version number is automatically extracted from `frontend/src/version.js`.

## Troubleshooting

### Permission Issues

If you get a "permission denied" error, make sure the scripts are executable:
```bash
chmod +x build-mac.sh
chmod +x build-mac-simple.sh
chmod +x build-mac-complete.sh
```

### Universal Build Issues

If you encounter issues with universal builds, you can modify the scripts to build for your specific architecture by changing:
```bash
wails build -platform darwin/universal
```

To either:
```bash
wails build -platform darwin/amd64  # For Intel Macs
```

Or:
```bash
wails build -platform darwin/arm64  # For Apple Silicon (M1/M2) Macs
```

### Background Image Issues

If the background image doesn't appear in the DMG:
1. Make sure Python and Pillow are installed:
   ```bash
   pip3 install Pillow
   ```
2. The background image is created in `build/resources/dmg-background.png`
3. Check that the DMG window is large enough to display the background

### Code Signing

The current scripts don't handle code signing. For distribution outside of development:

1. Get an Apple Developer certificate
2. Use `codesign` to sign the application
3. Consider using `notarytool` for notarization 