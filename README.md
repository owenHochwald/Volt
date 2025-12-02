<div>

# Volt

**A blazingly fast, terminal-native HTTP client with Vim keybindings**

<br>

[![GitHub release](https://img.shields.io/github/release/owenHochwald/Volt.svg)](https://github.com/owenHochwald/Volt/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/owenHochwald/volt)](https://goreportcard.com/report/github.com/owenHochwald/volt)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

[Installation](#installation) • [Why Volt?](#why-volt)

<!-- Add your demo.gif here once created -->
<!-- ![Demo](demo.gif) -->

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

### What makes Volt different?

**Speed First** - Built in Go with a concurrent architecture that handles HTTP requests asynchronously. Peak performance: **213K req/sec** with 4µs latency.

**Keyboard-Driven** - Navigate, edit, and send requests entirely with your keyboard. Vim-inspired motions (h/j/k/l) eliminate the need for a mouse.

**Terminal Native** - No Electron, no browser overhead. A single 15MB binary that integrates seamlessly into your terminal workflow.

**Beautiful Responses** - Syntax-highlighted JSON/XML/HTML with color-coded status codes, automatic formatting, and scrollable views.

## Performance

Volt is built for raw speed. The core engine is designed to minimize memory allocations and maximize throughput.
### Benchmarks (Apple M4)

| Concurrency | Req/Sec | Allocations |
|-------------|--------|-------------|
| 10          | 141,533 | 0 B/op      |
| 50          | 208,035   | 0 B/op      |
| **100**     | **213,885**   | **0 B/op**  |
| 500         | 92,891   | 4 B/op      |

> **Note:** At peak performance (100 workers), Volt can generate **1 Million requests in < 5 seconds**.

*Benchmarks run against a zero-latency local endpoint to measure engine overhead.*

### Real-World Benchmarks vs. `hey` ⚡


| Connections (C) | Volt Req/Sec | `hey` Req/Sec | **Volt Advantage** |
|:---------------:|:------------:|:-------------:|:------------------:|
| 10              | 137.9        | 145.3         | 0.95x              |
| **50** | **561.8** | 144.0         | **3.9x faster** |
| 100             | 877.5        | 294.9         | 3.0x faster        |

* You can test this with the included `benchmark-remote.sh` script.*


---

## Installation

### (WIP) Homebrew (macOS & Linux) 

```bash
brew install owenHochwald/volt/volt
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/owenHochwald/Volt/releases/latest).

#### macOS

```bash
# Intel Macs
curl -LO https://github.com/owenHochwald/Volt/releases/latest/download/volt_0.1.0_darwin_amd64.tar.gz
tar -xzf volt_0.1.0_darwin_amd64.tar.gz
chmod +x volt
sudo mv volt /usr/local/bin/

# Apple Silicon Macs (M1/M2/M3/M4)
curl -LO https://github.com/owenHochwald/Volt/releases/latest/download/volt_0.1.0_darwin_arm64.tar.gz
tar -xzf volt_0.1.0_darwin_arm64.tar.gz
chmod +x volt
sudo mv volt /usr/local/bin/
```

#### Linux

```bash
# x86_64
curl -LO https://github.com/owenHochwald/Volt/releases/latest/download/volt_0.1.0_linux_amd64.tar.gz
tar -xzf volt_0.1.0_linux_amd64.tar.gz
chmod +x volt
sudo mv volt /usr/local/bin/

# ARM64
curl -LO https://github.com/owenHochwald/Volt/releases/latest/download/volt_0.1.0_linux_arm64.tar.gz
tar -xzf volt_0.1.0_linux_arm64.tar.gz
chmod +x volt
sudo mv volt /usr/local/bin/
```

#### Windows

Download `volt_0.1.0_windows_amd64.zip` from the [releases page](https://github.com/owenHochwald/Volt/releases/latest), extract it, and add `volt.exe` to your PATH.

### Build from Source

**Prerequisites:** Go 1.25 or higher

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

This will launch the TUI interface where you can start making HTTP requests.

### Vim Philosophy

Volt embraces Vim's modal efficiency:
- Navigation uses h/j/k/l where applicable
- Common actions are single-key commands
- Esc always returns you to a safe state
- Focus is indicated visually, eliminating mode confusion


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

