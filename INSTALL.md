# Hopsule CLI Installation Guide

## 📦 Installation

### Option 1: Homebrew (Recommended - macOS/Linux)

```bash
brew install hopsule/tap/hopsule
```

### Option 2: Manual Installation

#### macOS ARM64 (M1/M2/M3)
```bash
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-arm64.tar.gz | tar xz
mv decision-darwin-arm64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

#### macOS Intel (x86_64)
```bash
curl -L https://github.com/Hopsule/cli-tool/releases/latest/download/decision-darwin-amd64.tar.gz | tar xz
mv decision-darwin-amd64 /usr/local/bin/hopsule
chmod +x /usr/local/bin/hopsule
```

---

## 🚀 Quick Start

### 1. Show Dashboard
```bash
hopsule
```

Output:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

        ████      ████                       Hopsule
        ████      ████                       
            ████████                         
            ████████                         
        ████          ████                   
        ████          ████                   

        org: hopsule-inc  •  project: app
        capture: ON  •  sync: ON  •  privacy: redacted
        last sync: 12s  •  latency: 84ms
        ─────────────────────────────────────────────────────────────────────────────
        ✓ Connected

        Get started
        ❯ hopsule config  (configure cli)
          hopsule list    (list decisions)
          hopsule create  (create decision)
          hopsule status  (health check)

        API: https://api.hopsule.com
        Token: not set

        Run 'hopsule --help' for more commands
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 2. Configure CLI
```bash
hopsule config
```

You'll be prompted for:
- **API URL**: https://api.hopsule.com (or your decision-api URL)
- **Project ID**: Your project identifier
- **Organization**: Your organization name
- **Auth Token**: Your JWT authentication token

### 3. Use Commands
```bash
# List decisions
hopsule list

# Create decision
hopsule create

# Get decision details
hopsule get <id>

# Accept decision
hopsule accept <id>

# Deprecate decision
hopsule deprecate <id>

# Show status
hopsule status

# Get help
hopsule --help
```

---

## 🔧 Configuration

### Config File Location
```
~/.decision-cli/config.yaml
```

### Manual Configuration
```yaml
api_url: https://api.hopsule.com
project: your-project-id
organization: your-org-name
token: your-jwt-token
```

### Environment Variables
```bash
export DECISION_API_URL=https://api.hopsule.com
export DECISION_PROJECT=your-project-id
export DECISION_TOKEN=your-jwt-token
```

---

## 🔄 Updating

### Homebrew
```bash
brew update
brew upgrade hopsule
```

### Manual
Download the latest release and replace the binary.

---

## ✅ Verification

```bash
# Check version
hopsule --version

# Show dashboard
hopsule

# Test connection
hopsule status
```

---

## 🆘 Troubleshooting

### Command not found after install
```bash
# Open a new terminal window
# Or reload shell
exec zsh

# Or use full path
/opt/homebrew/bin/hopsule
```

### Connection issues
```bash
# Check if decision-api is running
curl https://api.hopsule.com/health

# Test with explicit flags
hopsule list --api-url https://api.hopsule.com --project <id> --token <token>
```

### Reset configuration
```bash
rm -rf ~/.decision-cli/config.yaml
hopsule config
```

---

## 📚 Documentation

- [CLI Tool Repository](https://github.com/Hopsule/cli-tool)
- [Decision API](https://github.com/Hopsule/api)
- [Full Documentation](https://github.com/Hopsule/cli-tool#readme)

---

## 🎉 Success!

You're ready to use Hopsule CLI for decision-first workflow management!

Run `hopsule` to see your dashboard.
