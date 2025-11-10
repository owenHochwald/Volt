# Volt
[![Go Report Card](https://goreportcard.com/badge/github.com/owenHochwald/volt)](https://goreportcard.com/report/github.com/owenHochwald/volt)

A high-performance, terminal-based HTTP client built for developers who live in the command line.

## Overview

Volt is a fast, keyboard-driven API testing tool that brings modern HTTP clients directly to your terminal. Built with Go and the Bubble Tea TUI framework, Volt delivers a responsive, concurrent interface without sacrificing the speed and simplicity of the command line.

## Why Volt?

Modern API clients are powerful but heavy, slow for developer speed, and unintuitive. Volt takes a different approach:

**Speed First**: Volt is built in Go with a concurrent architecture that handles HTTP requests asynchronously. The message-passing concurrency model of Bubble Tea ensures the UI remains responsive even during long-running requests.

**Keyboard-Driven**: Navigate, edit, and send requests entirely with your keyboard. Vim-inspired motions (h/j/k/l) and intuitive shortcuts eliminate the need for a mouse and keep you in flow.

**Terminal Native**: No Electron, no browser overhead. Volt is a single statically-compiled binary that integrates naturally into your terminal workflow. 

**Developer-Focused**: Clean, syntax-highlighted responses. Request validation with immediate feedback. Persistent storage for your API collections. Built by developers, for developers.

## Features

**Current**
- HTTP methods: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, etc.
- Request builder with URL, headers, and JSON body support
- Three-panel layout: saved requests, request editor, and response viewer
- Key-value parser for headers and body with error reporting
- Request validation and immediate feedback
- Concurrent request handling with performance metrics
- Vim-style keyboard navigation
- **Beautiful response viewer:**
  - Syntax-highlighted responses with Chroma (JSON, XML, HTML)
  - Automatic JSON formatting with proper indentation
  - Color-coded HTTP status codes (green=2xx, orange=3xx, red=4xx/5xx)
  - Response timing and size metrics
  - Scrollable viewport for large responses (j/k navigation)
  - Multiple content-type support (JSON, XML, HTML, plain text)

**Planned**
- SQLite-based persistent storage for request collections
- Response history and diff viewer
- Environment variables and templating
- Custom themes and configuration files
- GraphQL support
- Enhanced response viewer (headers tab, copy/download actions)

## Quick Start

### Prerequisites

- Go 1.25 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/owenHochwald/volt.git
cd volt

# Build the binary
make build

# Run Volt
./volt
```

### Build Commands

```bash
make build       # Build for current platform
make build-mac   # Build for macOS (amd64)
make run         # Build and run
make clean       # Clean build artifacts
```

## Keyboard Shortcuts

Volt is designed to be used entirely with your keyboard.

### Global Navigation
- `q` or `Ctrl+C` - Quit application
- `Shift+Tab` - Cycle between panels (sidebar → request → response)
- `Esc` - Return to sidebar from request panel

### Sidebar Panel
- `j` / `↓` - Move down in list
- `k` / `↑` - Move up in list
- `Enter` or `Space` - Select request and open in editor
- `n` - Create new request
- `d` - Delete selected request
- `/` - Filter/search requests

### Request Editor
- `Tab` / `↑` / `↓` - Navigate between form fields
- `h` / `←` - Cycle to previous HTTP method
- `l` / `→` - Cycle to next HTTP method
- `Enter` - Send request (when on submit button)
- Standard editing keys work in text inputs and text areas

### Response Viewer
- `j` / `↓` - Scroll down through response
- `k` / `↑` - Scroll up through response
- `d` - Scroll down half page
- `u` - Scroll up half page
- `g` - Jump to top of response
- `G` - Jump to bottom of response

### Vim Philosophy

Volt embraces Vim's modal efficiency:
- Navigation uses h/j/k/l where applicable
- Common actions are single-key commands
- Esc always returns you to a safe state
- Focus is indicated visually, eliminating mode confusion

## Architecture

Volt is structured around clean separation of concerns with a concurrent message-passing architecture.


### Design Decisions

**Zero External Runtime Dependencies**: Volt compiles to a single binary with no external dependencies. The UI is rendered using ANSI escape sequences, ensuring compatibility across terminals.

**Stateless HTTP Client**: The HTTP client is designed to be stateless and concurrent. Multiple requests can be in flight simultaneously without blocking the UI, and request metrics (timing, response size) are collected automatically.

**Structured Parsing**: Headers and request bodies use a custom key-value parser that provides immediate validation feedback. Parse errors are surfaced inline, helping developers fix issues before sending requests.

**Validation First**: Requests are validated before sending (URL format, method validity, header count limits, body size limits). This catches common errors early and provides clear error messages.

### Concurrency Model

Volt leverages Go's concurrency primitives and Bubble Tea's message-passing architecture:

- HTTP requests execute in goroutines and send results back as messages
- The UI remains responsive during long-running operations
- Request timing is measured with microsecond precision
- Multiple panels can process input independently based on focus state

This architecture ensures Volt stays fast even when working with slow APIs or large responses.

## Customization

Volt is built to be customizable:

**Color Schemes**: HTTP methods are color-coded (GET=green, POST=orange, PUT=blue, PATCH=purple, DELETE=red). Styles are defined in `internal/ui/styles.go` and can be easily modified.

**Keybindings**: Key mappings are defined in `internal/app/commands.go`. Change them to match your preferences.

**Request Defaults**: Default request templates and validation rules are in `internal/http/request.go`.

Future versions will support configuration files for themes, keybindings, and default settings.

## Roadmap

### Near Term
- Complete persistent storage implementation (SQLite)
- Request history and quick recall
- Environment variable support
- Collection management and organization

### Medium Term
- Collection import/export (Postman/Insomnia formats)
- Enhanced response viewer (headers, cookies, timeline tabs)
- Custom themes and configuration files

### Long Term
- GraphQL support
- WebSocket testing
- Plugin system for extensibility

## Contributing

Volt is in active development. Contributions, bug reports, and feature requests are welcome.

## License

MIT License - see LICENSE file for details.