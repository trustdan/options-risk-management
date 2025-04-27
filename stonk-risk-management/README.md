# Options Trading Dashboard

A comprehensive tool for options traders to track trades, manage risk, and analyze performance.

## Build Instructions

### Prerequisites

You'll need the following tools installed to build the application:

- Go (v1.21+)
- Node.js (v20.10+)
- Wails CLI (v2+)
- Inno Setup (v6+, for creating Windows installers)

You can install all these prerequisites automatically by running:

```
.\install-build-tools.bat
```

This will install all required tools for you. It's recommended to run this script with administrator privileges.

### Building the Application

To build the application, simply run:

```
.\build.bat
```

This will:
1. Compile the Go backend
2. Build the frontend
3. Package the application
4. Create an installer (if Inno Setup is installed)

The built executable will be available in the `build/bin` directory, and the installer (if created) will be in the `installer` directory.

### Troubleshooting

If the build fails, try these steps:

1. Ensure you've installed all prerequisites (run `install-build-tools.bat`)
2. Make sure you've refreshed your terminal after installing tools
3. Try running `wails build -skipbindings` directly to bypass frontend binding generation
4. Check that your `%GOPATH%\bin` directory is in your PATH
5. If Node.js issues occur, try running `npm install` in the `frontend` directory

## Development

For development, you can use:

```
wails dev
```

This will start the application in development mode with hot reloading. 