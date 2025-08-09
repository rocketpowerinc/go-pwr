# Project Structure

This document outlines the structure and organization of the go-pwr project.

## Directory Structure

```
go-pwr/
├── cmd/                    # Entry points for applications
│   └── go-pwr/            # Main application entry point
│       └── main.go        # Application bootstrap
├── internal/              # Private application code
│   ├── app/               # Application core logic
│   │   └── app.go         # Main application runner
│   ├── config/            # Configuration management
│   │   └── config.go      # Config loading and path resolution
│   ├── git/               # Git repository operations
│   │   └── repository.go  # Repository cloning and management
│   ├── scripts/           # Script discovery and management
│   │   └── scripts.go     # Script items, caching, and file operations
│   └── ui/                # User interface components
│       ├── components/    # Reusable UI components
│       │   └── components.go  # List delegates and component creation
│       ├── styles/        # UI styling and themes
│       │   └── themes.go  # Color schemes and style definitions
│       ├── model.go       # Main UI model and state
│       └── update.go      # UI update logic and event handling
├── pkg/                   # Public library code
│   └── platform/          # Platform-specific utilities
│       └── platform.go    # OS detection and script execution
├── build/                 # Build artifacts (created by make)
├── Makefile              # Build automation
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # Project documentation
```

## Package Responsibilities

### `cmd/go-pwr`
- **Purpose**: Application entry point
- **Responsibilities**: Bootstrap the application and handle top-level errors

### `internal/app`
- **Purpose**: Application orchestration
- **Responsibilities**: Coordinate configuration, repository setup, and UI startup

### `internal/config`
- **Purpose**: Configuration management
- **Responsibilities**: Load configuration, resolve paths, handle environment-specific settings

### `internal/git`
- **Purpose**: Git operations
- **Responsibilities**: Clone and update the scriptbin repository

### `internal/scripts`
- **Purpose**: Script management
- **Responsibilities**: Discover scripts, cache content, sanitize file content

### `internal/ui`
- **Purpose**: User interface
- **Responsibilities**: Handle all UI rendering, event processing, and state management

### `internal/ui/components`
- **Purpose**: Reusable UI components
- **Responsibilities**: List delegates, component creation helpers

### `internal/ui/styles`
- **Purpose**: UI styling
- **Responsibilities**: Color schemes, themes, and style definitions

### `pkg/platform`
- **Purpose**: Platform utilities
- **Responsibilities**: OS detection, platform-specific script execution

## Design Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Dependency Direction**: Dependencies flow inward (cmd → internal → pkg)
3. **Encapsulation**: Internal packages are not exposed outside the module
4. **Testability**: Clear interfaces make testing easier
5. **Maintainability**: Logical organization makes code easier to understand and modify

## Build and Development

Use the provided Makefile for common development tasks:

- `make build` - Build the application
- `make install` - Install to GOPATH/bin
- `make dev` - Build and run in development mode
- `make test` - Run tests
- `make fmt` - Format code
- `make clean` - Clean build artifacts

## Adding New Features

When adding new features:

1. **UI Features**: Add to `internal/ui/` packages
2. **Core Logic**: Add to appropriate `internal/` packages
3. **Platform Features**: Add to `pkg/platform/`
4. **External Utilities**: Consider if they belong in `pkg/` or `internal/`

Follow the existing patterns and maintain the separation of concerns.
