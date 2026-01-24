# Decision CLI

CLI tool for decision-first workflow management.

## Architecture

**IMPORTANT**: This CLI tool is a **CLIENT** that communicates with `decision-api`, the single authoritative source for all decisions.

### Authority Model

- **decision-api** is the ONLY authoritative component
- This CLI MUST NOT bypass decision-api
- This CLI MUST NOT create authority independently
- All state changes flow through decision-api

The CLI:
- ✅ Makes HTTP requests to decision-api
- ✅ Displays decisions, memories, and context packs
- ✅ Requests state changes from decision-api
- ❌ Does NOT own data
- ❌ Does NOT decide anything
- ❌ Does NOT create authority

## Installation

### From Source

```bash
git clone https://github.com/Cagangedik/cli-tool.git
cd cli-tool
go build -o decision ./cmd/decision
sudo mv decision /usr/local/bin/
```

### Using go install

```bash
go install github.com/Cagangedik/cli-tool/cmd/decision@latest
```

### Binary Download

Download pre-built binaries from [Releases](https://github.com/Cagangedik/cli-tool/releases).

## Configuration

### Initial Setup

Run the interactive configuration:

```bash
decision config
```

This will prompt for:
- API URL (default: `http://localhost:8080`)
- Authentication token
- Default project ID

Configuration is saved to `~/.decision-cli/config.yaml`.

### Environment Variables

You can also set configuration via environment variables:

```bash
export DECISION_API_URL="https://api.example.com"
export DECISION_TOKEN="your-token-here"
export DECISION_PROJECT="project-id"
```

### Command-Line Flags

All commands support global flags that override config:

```bash
decision list --api-url https://api.example.com --project <id> --token <token>
```

## Commands

### `decision list`

List all decisions for the current project.

```bash
decision list
decision list --project <project-id>
```

### `decision get <decision-id>`

Get detailed information about a specific decision.

```bash
decision get <decision-id>
decision get <decision-id> --output json
```

### `decision create`

Interactively create a new decision (will be in DRAFT status).

```bash
decision create
```

Prompts for:
- Statement (title)
- Rationale (description)

### `decision accept <decision-id>`

Accept a decision, moving it from DRAFT/PENDING to ACCEPTED status.

```bash
decision accept <decision-id>
```

**Note**: Only decision-api can accept decisions. This command requests acceptance from the API.

### `decision deprecate <decision-id>`

Deprecate a decision, moving it to DEPRECATED status.

```bash
decision deprecate <decision-id>
```

### `decision sync`

Sync local state with remote decision-api.

```bash
decision sync
```

### `decision status`

Show current project status and statistics.

```bash
decision status
decision status --output json
```

### `decision config`

Interactively configure CLI settings.

```bash
decision config
```

### `decision --version`

Show version information.

```bash
decision --version
```

## Usage Examples

### Basic Workflow

```bash
# Configure CLI
decision config

# List decisions
decision list

# Create a new decision
decision create

# Accept a decision
decision accept <decision-id>

# Check project status
decision status
```

### With Flags

```bash
# List decisions for a specific project
decision list --project abc123 --api-url https://api.example.com

# Get decision details in JSON format
decision get <decision-id> --output json --project abc123
```

## Development

### Prerequisites

- Go 1.23 or later
- Access to decision-api instance

### Building

```bash
go build -o decision ./cmd/decision
```

### Running Tests

```bash
go test ./...
```

### Project Structure

```
cli-tool/
├── cmd/
│   └── decision/
│       └── main.go          # CLI entry point
├── internal/
│   ├── api/
│   │   └── client.go        # decision-api HTTP client
│   ├── commands/
│   │   ├── list.go          # List command
│   │   ├── get.go           # Get command
│   │   ├── create.go        # Create command
│   │   ├── accept.go        # Accept command
│   │   ├── deprecate.go     # Deprecate command
│   │   ├── sync.go          # Sync command
│   │   ├── status.go        # Status command
│   │   └── config.go        # Config command
│   └── config/
│       └── config.go        # Configuration management
├── go.mod
├── go.sum
└── README.md
```

## Architecture Principles

This CLI tool adheres to the decision-first architecture:

1. **Single Authority**: decision-api is the only authoritative source
2. **Client Role**: CLI is a client, not a decision-maker
3. **No Bypass**: All operations go through decision-api
4. **Reflective**: CLI reflects authority, does not create it

## License

[Add license information]

## Contributing

[Add contributing guidelines]
