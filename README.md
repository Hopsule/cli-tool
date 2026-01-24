# Hopsule CLI

> Decision & Memory Layer for AI teams & coding tools

[![Release](https://img.shields.io/github/v/release/Hopsule/cli-tool)](https://github.com/Hopsule/cli-tool/releases)
[![License](https://img.shields.io/github/license/Hopsule/cli-tool)](LICENSE)

## ✨ Features

- 🎨 **Minimal Dashboard** - Clean, focused interface
- ⌨️ **Keyboard Navigation** - Arrow keys for command selection
- 🖤 **Monochrome Theme** - Works in dark and light terminals
- 🚀 **Essential Commands** - init, login, connect, dev
- 📦 **Easy Install** - Homebrew or direct download

## 🚀 Quick Start

### Installation

#### Homebrew (macOS/Linux)
```bash
brew install hopsule/tap/hopsule
```

#### Manual Installation
```bash
# macOS ARM64 (M1/M2/M3)
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-arm64.tar.gz | tar xz
mv decision-darwin-arm64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

### Usage

```bash
# Launch dashboard
hopsule
```

**Output:**
```
        ████      ████                       Hopsule
        ████      ████                       Decision & Memory Layer
            ████████                         for AI teams & coding tools
            ████████
        ████          ████                   v0.4.2
        ████          ████                   ─────────────────────────────
                                             Get started
                                             > hopsule init     (create config)
                                               hopsule login    (authenticate)
                                               hopsule connect  (link repo)
                                               hopsule dev      (interactive TUI)

                                             view all commands
```

**Keyboard Shortcuts:**
- `↑/↓` - Navigate commands
- `Enter` - Execute selected command
- `q` - Quit
- `?` - Help

## 🎯 Commands

| Command | Description |
|---------|-------------|
| `hopsule` | Launch interactive dashboard |
| `hopsule init` | Create configuration |
| `hopsule login` | Authenticate with decision-api |
| `hopsule connect` | Link repository |
| `hopsule dev` | Start interactive development mode |
| `hopsule --help` | Show help |
| `hopsule --version` | Show version |

## 📋 Requirements

- **decision-api** running and accessible
- **JWT Token** for authentication
- **Project ID** for your project

## ⚙️ Configuration

### Interactive Setup
```bash
hopsule init
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

## 🎨 Design Philosophy

**Minimal & Focused**
- Show only what matters
- No visual clutter
- Guide users to key actions
- Professional and elegant

**Universal Compatibility**
- Monochrome theme (black/white/gray)
- Works in dark and light terminals
- Adaptive colors via lipgloss
- Clean ASCII logo

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

## 📦 Release History

- **v0.4.0** - Minimal dashboard design
- **v0.3.0** - Monochrome theme + ASCII logo
- **v0.2.1** - Panic fix (lipgloss.Width)
- **v0.2.0** - Interactive TUI with bubbletea
- **v0.1.1** - Dashboard UI
- **v0.1.0** - Initial release

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🔗 Links

- [Decision API](https://github.com/Hopsule/api)
- [Web App](https://github.com/Hopsule/web-app)
- [Releases](https://github.com/Hopsule/cli-tool/releases)
- [Organization](https://github.com/Hopsule)

## 📞 Support

- GitHub Issues: [Report a bug](https://github.com/Hopsule/cli-tool/issues)
- Organization: [Hopsule](https://github.com/Hopsule)

---

Made with ❤️ by the Hopsule team
