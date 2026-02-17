# Hopsule CLI

> The command-line interface for developers working with decision-first, context-aware, portable memory systems.

[![Release](https://img.shields.io/github/v/release/Hopsule/cli-tool)](https://github.com/Hopsule/cli-tool/releases)
[![License](https://img.shields.io/github/license/Hopsule/cli-tool)](LICENSE)

## Overview

**Hopsule CLI** is a powerful command-line tool designed for developers who work with decision-first workflow management. It provides an intuitive interface to interact with the `decision-api`, enabling you to manage decisions, track project status, and maintain organizational memory directly from your terminal.

The CLI features an interactive Terminal User Interface (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), making it easy to navigate and execute commands without memorizing complex syntax.

## Features

- **Interactive TUI** - Full-screen, keyboard-navigable interface with organization/project selection
- **Card & List Views** - Toggle between card grid and paginated datatable with `v`
- **Hopper AI Chat** - Built-in AI assistant with project context awareness
- **Decision Management** - Create, list, accept, deprecate decisions with status filters
- **Memory Management** - Create, edit, delete project memories
- **Task Management** - List and Kanban views with status toggles
- **Capsule Browsing** - View and search context packs
- **Device Auth Login** - Browser-based secure authentication flow
- **Organization & Project Navigation** - Multi-org, multi-project support
- **Dynamic Search** - Real-time `/` search across all entity views
- **Homebrew Installation** - One-command installation on macOS
- **Flexible Configuration** - Config file, environment variable, and flag support
- **Adaptive Theme** - Works in both dark and light terminals

## Installation

Install Hopsule CLI using any of the methods below. See [INSTALL.md](INSTALL.md) for the full guide.

| Method | Command |
|--------|---------|
| **npm** | `npm install -g hopsule` |
| **Homebrew** | `brew install hopsule/tap/hopsule` |
| **Scoop** (Windows) | `scoop bucket add hopsule https://github.com/Hopsule/scoop-bucket && scoop install hopsule` |
| **Chocolatey** (Windows) | `choco install hopsule` |
| **Shell script** | `curl -fsSL https://raw.githubusercontent.com/Hopsule/cli-tool/main/install.sh \| bash` |

### Verify

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
- `↑/↓` or `k/j` - Navigate items
- `←/→` or `h/l` - Navigate grid / change page (list view)
- `Enter` - Select / open detail
- `/` - Search within decisions, memories, tasks, capsules
- `v` - Toggle Card / List view
- `tab` - Toggle status filter (decisions) / view mode (tasks)
- `n` - Create new item
- `e` - Edit selected item
- `d` - Delete selected item
- `a` - Accept decision
- `x` - Deprecate decision
- `q` / `Esc` - Go back / quit

### 2. Configure the CLI

Run the interactive configuration:

```bash
hopsule config
```

You'll be prompted for:
- **API URL**: The decision-api endpoint (default: `https://api.hopsule.com`)
- **Token**: Your JWT authentication token
- **Default Project ID**: Your project identifier

Alternatively, see the [Configuration](#configuration) section for manual setup.

### 3. Start Using Commands

Once configured, you can use any of the available commands. See the [Command Reference](#command-reference) section below.

## Command Reference

### Core Commands

#### `hopsule`
Launch the interactive TUI (default command when no subcommand is provided). Provides full organization/project navigation, decision/memory/task/capsule management, and Hopper AI chat.

```bash
hopsule
```

#### `hopsule login`
Authenticate via browser-based device auth flow. Opens your browser, waits for approval, and stores the token locally.

```bash
hopsule login
```

#### `hopsule logout`
Clear stored authentication credentials.

```bash
hopsule logout
```

#### `hopsule whoami`
Display current authenticated user info.

```bash
hopsule whoami
```

#### `hopsule orgs`
List organizations you belong to.

```bash
hopsule orgs
```

#### `hopsule projects`
List projects in an organization.

```bash
hopsule projects
```

#### `hopsule init`
Link the current directory to a Hopsule project. Creates a local `.hopsule.yaml` project config.

```bash
hopsule init
```

#### `hopsule config`
Interactively configure CLI settings (API URL, token, default project).

```bash
hopsule config
```

### Decision Management Commands

#### `hopsule list`
List all decisions for the current project.

```bash
hopsule list
hopsule list --project my-project-id
hopsule list --api-url https://api.hopsule.com --token your-token
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
api_url: https://api.hopsule.com
project: your-project-id
organization: your-org-name
token: your-jwt-token
```

### Environment Variables

You can also configure via environment variables (takes precedence over config file):

```bash
export DECISION_API_URL=https://api.hopsule.com
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
4. **Defaults** (API URL defaults to `https://api.hopsule.com`)

### Manual Configuration

You can manually create/edit the config file:

```bash
mkdir -p ~/.decision-cli
cat > ~/.decision-cli/config.yaml << EOF
api_url: https://api.hopsule.com
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
- ✅ **Linux** (AMD64, ARM64)
- ✅ **Windows 10/11** (AMD64, ARM64)

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
go build -o hopsule .

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
go build -o hopsule .

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o hopsule-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o hopsule-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o hopsule-windows-amd64.exe .
```

### Project Structure

```
cli-tool/
├── main.go                        # Entry point, Cobra root command + interactive TUI launcher
├── internal/
│   ├── api/
│   │   └── client.go              # HTTP client for decision-api
│   ├── commands/                   # All CLI subcommands (login, list, create, etc.)
│   ├── config/                     # Configuration management
│   └── ui/                         # Bubble Tea TUI (styles, model, views, cards, list, etc.)
├── npm/                            # npm package (postinstall binary downloader)
│   ├── package.json
│   ├── install.js
│   └── bin/                        # Shell/CMD wrappers
├── choco/                          # Chocolatey package
│   ├── hopsule.nuspec
│   └── tools/                      # Install/uninstall scripts
├── .github/workflows/              # CI/CD (ci.yml, release.yml)
├── .goreleaser.yml                 # GoReleaser config (Homebrew + Scoop + binaries)
├── install.sh                      # Universal curl installer
├── LICENSE                         # MIT License
├── INSTALL.md                      # Comprehensive installation guide
└── README.md                       # This file
```

### Dependencies

Key dependencies:
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[color](https://github.com/fatih/color)** - ANSI color output

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
curl https://api.hopsule.com/health
```

**Test with explicit flags:**
```bash
hopsule list --api-url https://api.hopsule.com --project <id> --token <token>
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

- **v0.9.0** - Production distribution: npm, Homebrew, Scoop, Chocolatey, curl installer, CI/CD, Windows/Linux support
- **v0.8.0** - Card/List toggle, paginated datatable, Hopper AI chat, modular UI refactor
- **v0.7.5** - Interactive TUI with org/project navigation, device auth login
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
