<div>

# Volt

**A blazingly fast, terminal-native HTTP client with Vim keybindings**

<br>

[![GitHub release](https://img.shields.io/github/release/owenHochwald/Volt.svg)](https://github.com/owenHochwald/Volt/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/owenHochwald/volt)](https://goreportcard.com/report/github.com/owenHochwald/volt)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

[Installation](#installation) • [Why Volt?](#why-volt)

![Demo](demo.gif)

</div>

---

## Overview

Volt is a **high-performance, keyboard-driven HTTP client** that lives in your terminal. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework, it delivers a responsive interface without sacrificing the speed and simplicity of the command line.

**Perfect for developers who:**
- Live in the terminal and hate context switching
- Want Postman's features without the Electron bloat
- Love Vim keybindings and keyboard-driven workflows
- Need a fast, scriptable HTTP client with a beautiful UI

## Why Volt?

|  | Postman | Insomnia | HTTPie | curl | **Volt** |
|---|:---:|:---:|:---:|:---:|:---:|
| **Terminal-native** | ❌ | ❌ | ✅ | ✅ | ✅ |
| **Interactive TUI** | ❌ | ❌ | ❌ | ❌ | ✅ |
| **Vim keybindings** | ❌ | ❌ | ❌ | ❌ | ✅ |
| **Syntax highlighting** | ✅ | ✅ | ✅ | ❌ | ✅ |
| **Save collections** | ✅ | ✅ | ❌ | ❌ | ✅ |
| **Zero install** | ❌ | ❌ | ❌ | ✅ | ✅ |
| **Memory footprint** | ~500MB | ~300MB | ~50MB | <5MB | **~15MB** |
| **Startup time** | ~3s | ~2s | <1s | instant | **instant** |




---

## Installation

**Prerequisites:** Go 1.25 or higher

### Option 1: Download and Build from Release

1. **Download the source code** from the [releases page](https://github.com/owenHochwald/Volt/releases/latest)
    - Download either the `.tar.gz` or `.zip` source code archive

2. **Extract the archive**
   ```bash
   # For .tar.gz
   tar -xzf volt-v0.1.0.tar.gz
   cd volt-0.1.0

   # For .zip
   unzip volt-v0.1.0.zip
   cd volt-0.1.0
   ```

3. **Build and install**
   ```bash
   # Build and install system-wide to /usr/local/bin
   make install

   # Or just build locally without installing
   make build
   ./volt
   ```

### Option 2: Clone and Build (Development)

If you want to contribute or build from the latest source:

```bash
# Clone the repository
git clone https://github.com/owenHochwald/volt.git
cd volt

# Build and install (makes 'volt' available system-wide)
make install

# Or just build locally
make build
./volt
```

## Quick Start

Once installed, simply run:

```bash
volt
```


## CLI Load Testing

Volt also includes a powerful little HTTP load testing tool for direct access, accessible via the `bench` subcommand.

```bash
volt bench [flags]
```

### Examples

```bash
# Basic throughput test
volt bench -url http://localhost:8080 -c 100 -d 30s

# POST request with custom headers
volt bench -url http://localhost:8080/api -m POST \
  -b '{"test":true}' -H "Content-Type: application/json"

# JSON output to file for CI/CD
volt bench -url http://localhost:8080 -c 50 -d 60s -json -o results.json

# Rate-limited testing
volt bench -url http://localhost:8080 -c 10 -d 30s -rate 1000

# For help!
volt bench -h
```

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](./LICENSE) file for details.

## Star History

If you find Volt useful, please consider giving it a star ⭐ on GitHub!

---

