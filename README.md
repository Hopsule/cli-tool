# Hopsule CLI

> The command-line interface for developers working with decision-first, context-aware, portable memory systems.

[![Release](https://img.shields.io/github/v/release/Hopsule/cli-tool)](https://github.com/Hopsule/cli-tool/releases)
[![License](https://img.shields.io/github/license/Hopsule/cli-tool)](LICENSE)

## Overview

**Hopsule CLI** is a powerful command-line tool designed for developers who work with decision-first workflow management. It provides an intuitive interface to interact with the `decision-api`, enabling you to manage decisions, track project status, and maintain organizational memory directly from your terminal.

The CLI features an interactive Terminal User Interface (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) (v0.7.5), making it easy to navigate and execute commands without memorizing complex syntax.

## Features

- 🎨 **Interactive TUI (v0.7.5)** - Beautiful, keyboard-navigable interface with ASCII art logo
- 📦 **Homebrew Installation** - One-command installation on macOS
- 🎯 **Decision Management** - Create, list, accept, and deprecate decisions
- 🔐 **Authentication** - Secure JWT token-based authentication
- 📊 **Project Status** - View comprehensive project statistics
- 🔄 **Sync Capabilities** - Synchronize with remote decision-api
- ⚙️ **Flexible Configuration** - Config file and environment variable support
- 🎨 **Monochrome Theme** - Works beautifully in both dark and light terminals
- ⌨️ **Keyboard Navigation** - Arrow keys and vim-style navigation (j/k)

## Installation

### Option 1: Homebrew (Recommended - macOS)

```bash
brew install hopsule/tap/hopsule
```

This installs the latest stable release from the Hopsule Homebrew tap.

### Option 2: Manual Installation

#### macOS

**Apple Silicon (M1/M2/M3)**
```bash
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-arm64.tar.gz | tar xz
mv decision-darwin-arm64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

**Intel (x86_64)**
```bash
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-amd64.tar.gz | tar xz
mv decision-darwin-amd64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

#### Windows

**PowerShell**
```powershell
# Download and extract
Invoke-WebRequest -Uri https://github.com/Hopsule/cli-tool/releases/latest/download/decision-windows-amd64.zip -OutFile hopsule.zip
Expand-Archive hopsule.zip -DestinationPath .
Rename-Item decision-windows-amd64.exe hopsule.exe

# Run
.\hopsule.exe
```

**Manual Download**
1. Visit [Releases](https://github.com/Hopsule/cli-tool/releases/latest)
2. Download `decision-windows-amd64.zip`
3. Extract and rename `decision-windows-amd64.exe` to `hopsule.exe`
4. Add to PATH (optional, for global access)

### Verification

After installation, verify it works:

```bash
hopsule --version
```

## Quick Start

### 1. Launch the Interactive Dashboard

Simply run:

```bash
hopsule
```

This launches the interactive TUI dashboard with a beautiful ASCII logo and command menu:

```
  ──────────────────────────────────────────────────────────────────

  ⢠⣶⣶⣶⣶⣶⣶⣶⣆     	 ⣴⣶⣶⣶⣶⣶⣶⣶⡄
  ⢸⣿⣯⣷⣿⢿⣾⣟⣿    	 	 ⣿⣿⣽⣾⡿⣷⡿⣯⡇  Hopsule
  ⢸⣿⣾⣟⣿⣟⣯⣿⣿      	 ⣿⣷⣿⣻⡿⣟⣿⣿⡇  Decision & Memory Layer
  ⢸⣿⡾⣿⣯⣿⣟⣿⣾     	 ⣿⣷⡿⣟⣿⣿⣟⣷⡇  for AI teams & coding tools
  ⠘⢿⣻⣿⣽⣾⣿⣻⣽    		 ⣿⣷⣿⣿⢿⣷⣿⠿⠃
           ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀      v0.7.5
           ⣿⡿⣾⣿⡾⣿⣾⢿⣾⡿⣾⣿⣿⣄  ─────────────────────────────
           ⣿⣿⣻⣷⣿⣿⣻⣿⣯⣿⡿⣯⣿⣻⣿⣄  Get started
           ⣿⣽⣿⣽⣾⡿⣯⣷⣿⣯⣿⣿⣻⣿⣽⡿⣿⣄  > hopsule init     (create config)
  ⢠⣾⣿⣿⢿⣿⡿⣿⣿     	 ⣿⣽⣿⣽⣷⣿⡿⣟⡇    hopsule login    (authenticate)
  ⢸⣿⢷⣿⣿⣻⣿⣻⣽     	 ⣿⣿⣽⣯⣿⣾⣿⣿⡇    hopsule connect  (link repo)
  ⢸⣿⢿⣻⣾⣿⣻⣿⣻     	 ⣿⣷⣿⣯⣿⣷⣿⣾⡇    hopsule dev      (interactive TUI)
  ⢸⣿⢿⣿⣻⣽⣿⣽⣿      	 ⣿⣷⡿⣷⣿⣾⣷⡿⡇
  ⠘⠿⠿⠻⠿⠻⠽⠿⠾      	 ⠻⠷⠿⠿⠻⠾⠟⠿⠁  view all commands

  ──────────────────────────────────────────────────────────────────
```

**Keyboard Shortcuts:**
- `↑/↓` or `k/j` - Navigate commands
- `Enter` - Execute selected command
- `q` - Quit
- `?` - Show help

### 2. Configure the CLI

Run the interactive configuration:

```bash
hopsule config
```

You'll be prompted for:
- **API URL**: The decision-api endpoint (default: `http://localhost:8080`)
- **Token**: Your JWT authentication token
- **Default Project ID**: Your project identifier

Alternatively, see the [Configuration](#configuration) section for manual setup.

### 3. Start Using Commands

Once configured, you can use any of the available commands. See the [Command Reference](#command-reference) section below.

## Command Reference

### Core Commands

#### `hopsule`
Launch the interactive dashboard (default command when no subcommand is provided).

```bash
hopsule
```

#### `hopsule config`
Interactively configure CLI settings (API URL, token, default project).

```bash
hopsule config
```

**What it does:**
- Prompts for API URL (with default: `http://localhost:8080`)
- Prompts for authentication token (masks existing token)
- Prompts for default project ID
- Saves configuration to `~/.decision-cli/config.yaml`

#### `hopsule init`
Alias for `hopsule config` - creates initial configuration file.

```bash
hopsule init
```

#### `hopsule login`
Authenticate with decision-api (configure your JWT token).

**Note:** Currently maps to `config` command. Future versions will include dedicated authentication flow.

```bash
hopsule login
```

#### `hopsule connect`
Link repository to a project.

**Note:** Planned feature for repository linking.

```bash
hopsule connect
```

#### `hopsule dev`
Start interactive development mode (TUI).

**Note:** Currently launches the same interactive dashboard as `hopsule`.

```bash
hopsule dev
```

### Decision Management Commands

#### `hopsule list`
List all decisions for the current project.

```bash
hopsule list
hopsule list --project my-project-id
hopsule list --api-url http://localhost:8080 --token your-token
```

**Output:**
```
ID          TITLE                                    STATUS    CREATED
---         -----                                    ------    -------
abc123def... Use TypeScript for...                   ACCEPTED  2024-01-15T10:30:00Z
xyz789ghi... Enforce code review...                  PENDING   2024-01-16T14:20:00Z
```

**Flags:**
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

#### `hopsule create`
Interactively create a new decision (will be in DRAFT status).

```bash
hopsule create
hopsule create --project my-project-id
```

**Interactive prompts:**
1. **Statement**: The decision statement (required)
2. **Rationale**: Multi-line rationale (end with empty line)

**Example:**
```bash
$ hopsule create
Statement: Use TypeScript for all new frontend code
Rationale (multi-line, end with empty line):
TypeScript provides type safety and better IDE support.
It reduces runtime errors and improves code maintainability.
End with empty line:

Decision created successfully!
ID: abc123def456
Status: DRAFT
```

**Flags:**
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

#### `hopsule get <decision-id>`
Get detailed information about a specific decision.

```bash
hopsule get abc123def456
hopsule get abc123def456 --output json
```

**Output (text format):**
```
ID: abc123def456
Statement: Use TypeScript for all new frontend code
Status: ACCEPTED
Created: 2024-01-15T10:30:00Z
Updated: 2024-01-15T11:00:00Z
Accepted: 2024-01-15T11:00:00Z by user@example.com
Tags: frontend, typescript

Rationale:
TypeScript provides type safety and better IDE support.
It reduces runtime errors and improves code maintainability.
```

**Flags:**
- `--output, -o` - Output format: `text` (default) or `json`
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

#### `hopsule accept <decision-id>`
Accept a decision, moving it from DRAFT/PENDING to ACCEPTED status.

```bash
hopsule accept abc123def456
```

**Output:**
```
Decision accepted successfully!
ID: abc123def456
Status: ACCEPTED
```

**Flags:**
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

#### `hopsule deprecate <decision-id>`
Deprecate a decision, moving it to DEPRECATED status.

```bash
hopsule deprecate abc123def456
```

**Output:**
```
Decision deprecated successfully!
ID: abc123def456
Status: DEPRECATED
```

**Flags:**
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

### Project Management Commands

#### `hopsule status`
Show current project status and statistics.

```bash
hopsule status
hopsule status --output json
```

**Output (text format):**
```
Project: my-project-id

Total Decisions: 42
  Accepted:   30
  Pending:    5
  Draft:      5
  Deprecated: 2
```

**Flags:**
- `--output, -o` - Output format: `text` (default) or `json`
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

#### `hopsule sync`
Sync local state with the remote decision-api.

```bash
hopsule sync
```

**What it does:**
- Tests connection to decision-api
- Verifies authentication
- Currently a no-op (future: cache management)

**Flags:**
- `--project` - Override default project ID
- `--api-url` - Override default API URL
- `--token` - Override default token

### Global Flags

All commands support these global flags:

- `--api-url` - Override the API URL from config
- `--project` - Override the default project ID from config
- `--token` - Override the authentication token from config

### Help and Version

```bash
hopsule --help          # Show help
hopsule --version       # Show version information
hopsule <command> --help # Show help for specific command
```

## Configuration

### Config File Location

The CLI stores configuration in:

```
~/.decision-cli/config.yaml
```

### Config File Structure

```yaml
api_url: http://localhost:8080
project: your-project-id
organization: your-org-name
token: your-jwt-token
```

### Environment Variables

You can also configure via environment variables (takes precedence over config file):

```bash
export DECISION_API_URL=http://localhost:8080
export DECISION_PROJECT=your-project-id
export DECISION_TOKEN=your-jwt-token
```

**Environment Variable Mapping:**
- `DECISION_API_URL` → `api_url`
- `DECISION_PROJECT` → `project`
- `DECISION_TOKEN` → `token`

### Configuration Precedence

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Config file** (`~/.decision-cli/config.yaml`)
4. **Defaults** (API URL defaults to `http://localhost:8080`)

### Manual Configuration

You can manually create/edit the config file:

```bash
mkdir -p ~/.decision-cli
cat > ~/.decision-cli/config.yaml << EOF
api_url: http://localhost:8080
project: my-project-id
token: your-jwt-token-here
EOF
```

## Requirements

- **decision-api** - The authoritative API server must be running and accessible
- **JWT Token** - Valid authentication token for decision-api
- **Project ID** - A valid project identifier in your organization

### Supported Platforms

- ✅ **macOS 11+** (Apple Silicon M1/M2/M3, Intel x86_64)
- ✅ **Windows 10/11** (AMD64, ARM64)
- ⏳ **Linux** (coming soon)

## Architecture

### Authority Model

Hopsule CLI is a **client-only tool** that communicates with `decision-api`. It follows the decision-first architecture principles:

```
┌─────────────────┐
│   Hopsule CLI   │  ← Client (No Authority)
└────────┬────────┘
         │
         │ API Calls (HTTP/REST)
         │
         ▼
┌─────────────────┐
│  decision-api   │  ← Single Authority
└─────────────────┘
```

**Key Principles:**
- ✅ **Strictly Advisory** - CLI cannot create authority independently
- ✅ **API-First** - All operations go through decision-api
- ✅ **No Direct Database Access** - Only communicates via API
- ✅ **Stateless** - Configuration stored locally, state in API

The CLI respects the authoritative nature of decision-api and never bypasses it or creates decisions independently.

## Development

### Prerequisites

- **Go 1.24+** - Required for building from source
- **decision-api** - Should be running locally for testing

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Hopsule/cli-tool.git
cd cli-tool

# Build the binary
go build -o hopsule ./cmd/decision

# Run it
./hopsule
```

### Development Workflow

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build for your platform
go build -o hopsule ./cmd/decision

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o hopsule-linux-amd64 ./cmd/decision
GOOS=darwin GOARCH=arm64 go build -o hopsule-darwin-arm64 ./cmd/decision
GOOS=windows GOARCH=amd64 go build -o hopsule-windows-amd64.exe ./cmd/decision
```

### Project Structure

```
cli-tool/
├── cmd/
│   └── decision/
│       └── main.go          # Entry point, Cobra root command
├── internal/
│   ├── api/
│   │   └── client.go        # HTTP client for decision-api
│   ├── commands/
│   │   ├── accept.go        # Accept decision command
│   │   ├── config.go        # Configuration command
│   │   ├── create.go        # Create decision command
│   │   ├── deprecate.go     # Deprecate decision command
│   │   ├── get.go           # Get decision command
│   │   ├── list.go          # List decisions command
│   │   ├── status.go        # Status command
│   │   └── sync.go          # Sync command
│   ├── config/
│   │   └── config.go        # Configuration management
│   └── ui/
│       ├── dashboard.go     # Dashboard UI (future)
│       └── interactive.go   # Interactive TUI (Bubble Tea)
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── .goreleaser.yml          # Release configuration
└── README.md                 # This file
```

### Dependencies

Key dependencies:
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework (v0.7.5)
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling

## Troubleshooting

### Command Not Found After Installation

**macOS (Homebrew):**
```bash
# Open a new terminal window
# Or reload shell
exec zsh

# Or use full path
/opt/homebrew/bin/hopsule  # Apple Silicon
/usr/local/bin/hopsule     # Intel
```

**Windows:**
- Ensure the binary is in your PATH
- Or use the full path: `C:\path\to\hopsule.exe`

### Connection Issues

**Check if decision-api is running:**
```bash
curl http://localhost:8080/health
```

**Test with explicit flags:**
```bash
hopsule list --api-url http://localhost:8080 --project <id> --token <token>
```

**Verify configuration:**
```bash
cat ~/.decision-cli/config.yaml
```

### Authentication Errors

**Invalid token:**
- Ensure your JWT token is valid and not expired
- Get a new token from your decision-api administrator
- Update config: `hopsule config`

**Token not found:**
- Run `hopsule config` to set your token
- Or set `DECISION_TOKEN` environment variable

### Reset Configuration

```bash
# Remove config file
rm -rf ~/.decision-cli/config.yaml

# Reconfigure
hopsule config
```

## Release History

- **v0.7.5** - Current stable release with interactive TUI
- **v0.4.0** - Minimal dashboard design
- **v0.3.0** - Monochrome theme + ASCII logo
- **v0.2.1** - Panic fix (lipgloss.Width)
- **v0.2.0** - Interactive TUI with bubbletea
- **v0.1.1** - Dashboard UI
- **v0.1.0** - Initial release

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Contribution Guidelines

1. Follow Go best practices and conventions
2. Maintain compatibility with decision-api
3. Respect the authority model (CLI is advisory only)
4. Add tests for new features
5. Update documentation as needed

## License

MIT License - see [LICENSE](LICENSE) file for details

## Links

- **Repository**: [github.com/Hopsule/cli-tool](https://github.com/Hopsule/cli-tool)
- **Decision API**: [github.com/Hopsule/api](https://github.com/Hopsule/api)
- **Web App**: [github.com/Hopsule/web-app](https://github.com/Hopsule/web-app)
- **Releases**: [github.com/Hopsule/cli-tool/releases](https://github.com/Hopsule/cli-tool/releases)
- **Organization**: [github.com/Hopsule](https://github.com/Hopsule)

## Support

- **GitHub Issues**: [Report a bug](https://github.com/Hopsule/cli-tool/issues)
- **Documentation**: See [INSTALL.md](INSTALL.md) for detailed installation guide

---

Made with ❤️ by the Hopsule team
