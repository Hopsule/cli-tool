# Hopsule CLI — Installation Guide

## Quick Install

Pick your preferred method:

| Method | Platform | Command |
|--------|----------|---------|
| **npm** | All | `npm install -g hopsule` |
| **Homebrew** | macOS / Linux | `brew install hopsule/tap/hopsule` |
| **Scoop** | Windows | `scoop bucket add hopsule https://github.com/Hopsule/scoop-bucket && scoop install hopsule` |
| **Chocolatey** | Windows | `choco install hopsule` |
| **Shell script** | macOS / Linux | `curl -fsSL https://raw.githubusercontent.com/Hopsule/cli-tool/main/install.sh \| bash` |
| **Go install** | All (requires Go) | `go install github.com/Hopsule/cli-tool@latest` |

---

## Detailed Instructions

### npm (All Platforms)

Works on macOS, Linux, and Windows. Requires [Node.js](https://nodejs.org/) 16+.

```bash
npm install -g hopsule
```

The npm package automatically downloads the correct Go binary for your platform during installation.

**Update:**
```bash
npm update -g hopsule
```

**Uninstall:**
```bash
npm uninstall -g hopsule
```

---

### Homebrew (macOS / Linux)

```bash
brew install hopsule/tap/hopsule
```

**Update:**
```bash
brew update && brew upgrade hopsule
```

**Uninstall:**
```bash
brew uninstall hopsule
```

---

### Scoop (Windows)

```powershell
scoop bucket add hopsule https://github.com/Hopsule/scoop-bucket
scoop install hopsule
```

**Update:**
```powershell
scoop update hopsule
```

**Uninstall:**
```powershell
scoop uninstall hopsule
```

---

### Chocolatey (Windows)

```powershell
choco install hopsule
```

**Update:**
```powershell
choco upgrade hopsule
```

**Uninstall:**
```powershell
choco uninstall hopsule
```

---

### Shell Script (macOS / Linux)

One-line install that detects your OS and architecture, downloads the latest release, and installs to `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/Hopsule/cli-tool/main/install.sh | bash
```

This script:
- Detects OS (macOS / Linux) and architecture (amd64 / arm64)
- Downloads the latest release from GitHub
- Verifies SHA256 checksum
- Installs to `/usr/local/bin` (may request `sudo`)

---

### Manual Download

Download pre-built binaries from [GitHub Releases](https://github.com/Hopsule/cli-tool/releases/latest):

| Platform | Architecture | File |
|----------|-------------|------|
| macOS | Apple Silicon (M1/M2/M3/M4) | `hopsule-darwin-arm64.tar.gz` |
| macOS | Intel | `hopsule-darwin-amd64.tar.gz` |
| Linux | x86_64 | `hopsule-linux-amd64.tar.gz` |
| Linux | ARM64 | `hopsule-linux-arm64.tar.gz` |
| Windows | x86_64 | `hopsule-windows-amd64.zip` |
| Windows | ARM64 | `hopsule-windows-arm64.zip` |

**macOS / Linux:**
```bash
# Example for macOS Apple Silicon
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/hopsule-darwin-arm64.tar.gz | tar xz
sudo mv hopsule /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri https://github.com/Hopsule/cli-tool/releases/latest/download/hopsule-windows-amd64.zip -OutFile hopsule.zip
Expand-Archive hopsule.zip -DestinationPath .
Move-Item hopsule.exe C:\Windows\System32\  # or add to PATH
```

---

### Build from Source

Requires [Go 1.24+](https://go.dev/dl/).

```bash
git clone https://github.com/Hopsule/cli-tool.git
cd cli-tool
go build -o hopsule .
sudo mv hopsule /usr/local/bin/  # or add to PATH
```

---

## Verify Installation

```bash
hopsule --version
# hopsule version 0.9.0 (commit: ..., built: ...)
```

---

## Getting Started

### 1. Launch the Dashboard

```bash
hopsule
```

This opens the interactive TUI with organization and project navigation.

### 2. Authenticate

```bash
hopsule login
```

Opens your browser for device authentication. Once approved, your token is stored locally.

### 3. Explore

Use the interactive dashboard to navigate organizations, projects, decisions, memories, and capsules. Or use CLI commands directly:

```bash
hopsule list          # List decisions
hopsule create        # Create a decision
hopsule status        # Project status
hopsule whoami        # Current user
hopsule --help        # All commands
```

---

## Configuration

Config file: `~/.decision-cli/config.yaml`

```yaml
api_url: https://api.hopsule.com
project: your-project-id
organization: your-org-name
token: your-auth-token
```

Environment variables (override config file):

```bash
export DECISION_API_URL=https://api.hopsule.com
export DECISION_PROJECT=your-project-id
export DECISION_TOKEN=your-auth-token
```

---

## Troubleshooting

### "command not found" after install

**macOS/Linux:** Open a new terminal or run `exec $SHELL` to reload PATH.

**npm:** Ensure npm global bin is in PATH:
```bash
npm bin -g
# Add the output directory to your PATH if needed
```

**Homebrew on Apple Silicon:**
```bash
eval "$(/opt/homebrew/bin/brew shellenv)"
```

### Connection issues

```bash
# Check API health
curl https://api.hopsule.com/health

# Test with explicit flags
hopsule list --api-url https://api.hopsule.com --token YOUR_TOKEN
```

### Reset configuration

```bash
rm -rf ~/.decision-cli/config.yaml
hopsule config
```

---

## Supported Platforms

| OS | Architecture | Status |
|----|-------------|--------|
| macOS 11+ | Apple Silicon (arm64) | Supported |
| macOS 11+ | Intel (amd64) | Supported |
| Linux | x86_64 (amd64) | Supported |
| Linux | ARM64 | Supported |
| Windows 10+ | x86_64 (amd64) | Supported |
| Windows 10+ | ARM64 | Supported |

---

## Links

- [GitHub Repository](https://github.com/Hopsule/cli-tool)
- [Releases](https://github.com/Hopsule/cli-tool/releases)
- [Issues](https://github.com/Hopsule/cli-tool/issues)
- [Website](https://hopsule.com)
