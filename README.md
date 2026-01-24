# Hopsule CLI

> Decision-first workflow management CLI with interactive terminal UI

[![Release](https://img.shields.io/github/v/release/Hopsule/cli-tool)](https://github.com/Hopsule/cli-tool/releases)
[![License](https://img.shields.io/github/license/Hopsule/cli-tool)](LICENSE)

## ✨ Features

- 🎨 **Interactive TUI** - Daytona-style terminal interface
- ⚡ **Arrow Key Navigation** - Navigate commands with ↑/↓
- ⌨️ **Keyboard Shortcuts** - Full keyboard control (Enter, q, ?)
- 📊 **Real-time Status** - Live connection and sync status
- 🎯 **Command Execution** - Execute commands directly from TUI
- 🔧 **Configuration** - Easy setup with `hopsule config`
- 🌈 **Colored Output** - Beautiful terminal styling

## 🚀 Quick Start

### Installation

#### Homebrew (macOS/Linux)
```bash
brew tap hopsule/tap
brew install hopsule
```

#### Manual Installation
```bash
# macOS ARM64 (M1/M2/M3)
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-arm64.tar.gz | tar xz
mv decision-darwin-arm64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

### Usage

#### Interactive Mode
```bash
# Launch interactive dashboard
hopsule
```

**Keyboard Shortcuts:**
- `↑/↓` or `k/j` - Navigate commands
- `Enter` - Execute selected command
- `q` - Quit
- `?` - Toggle help

#### Direct Commands
```bash
# Configure CLI
hopsule config

# List all decisions
hopsule list

# Create new decision
hopsule create

# Get decision details
hopsule get <decision-id>

# Accept decision
hopsule accept <decision-id>

# Deprecate decision
hopsule deprecate <decision-id>

# Show project status
hopsule status

# Sync with decision-api
hopsule sync
```

## 📋 Requirements

- **decision-api** running and accessible
- **JWT Token** for authentication
- **Project ID** for your project

## ⚙️ Configuration

### Interactive Setup
```bash
hopsule config
```

### Manual Configuration

Config file: `~/.decision-cli/config.yaml`

```yaml
api_url: http://localhost:8080
project: your-project-id
organization: your-org-name
token: your-jwt-token
```

### Environment Variables
```bash
export DECISION_API_URL=http://localhost:8080
export DECISION_PROJECT=your-project-id
export DECISION_TOKEN=your-jwt-token
```

## 🎯 Commands

| Command | Description |
|---------|-------------|
| `hopsule` | Launch interactive dashboard |
| `hopsule config` | Configure CLI settings |
| `hopsule list` | List all decisions |
| `hopsule get <id>` | Get decision details |
| `hopsule create` | Create new decision |
| `hopsule accept <id>` | Accept a decision |
| `hopsule deprecate <id>` | Deprecate a decision |
| `hopsule status` | Show project status |
| `hopsule sync` | Sync with decision-api |
| `hopsule --help` | Show help |
| `hopsule --version` | Show version |

## 🏗️ Architecture

Hopsule CLI is a **client-only tool** that communicates with `decision-api`:

- ✅ **Strictly Advisory** - Cannot create authority independently
- ✅ **API-First** - All operations go through decision-api
- ✅ **No Direct Database Access** - Only communicates via API
- ✅ **Stateless** - Configuration stored locally, state in API

### Authority Model

```
┌─────────────────┐
│   Hopsule CLI   │  ← Client (No Authority)
└────────┬────────┘
         │
         │ API Calls
         │
         ▼
┌─────────────────┐
│  decision-api   │  ← Single Authority
└─────────────────┘
```

## 🛠️ Development

### Prerequisites
- Go 1.24+
- decision-api running locally

### Build from Source
```bash
git clone https://github.com/Hopsule/cli-tool.git
cd cli-tool
go build -o decision ./cmd/decision
./decision
```

### Run Tests
```bash
go test ./...
```

### Install Locally
```bash
go install ./cmd/decision
```

## 📦 Release Process

1. Update version
2. Create git tag
3. Build binaries
4. Create GitHub release
5. Update Homebrew formula

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🔗 Links

- [Decision API](https://github.com/Hopsule/api)
- [Web App](https://github.com/Hopsule/web-app)
- [Documentation](https://github.com/Hopsule/cli-tool/wiki)
- [Releases](https://github.com/Hopsule/cli-tool/releases)

## 📞 Support

- GitHub Issues: [Report a bug](https://github.com/Hopsule/cli-tool/issues)
- Organization: [Hopsule](https://github.com/Hopsule)

---

Made with ❤️ by the Hopsule team
