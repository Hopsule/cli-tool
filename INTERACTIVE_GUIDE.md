# Hopsule Interactive TUI Guide

## 🎨 Interactive Dashboard

When you run `hopsule` without any arguments, an interactive terminal UI launches with full keyboard control.

```bash
hopsule
```

## ⌨️ Keyboard Controls

### Navigation
| Key | Action |
|-----|--------|
| `↑` or `k` | Move selection up |
| `↓` or `j` | Move selection down |

### Actions
| Key | Action |
|-----|--------|
| `Enter` | Execute selected command |
| `q` | Quit the dashboard |
| `?` | Toggle help screen |

## 🎯 Dashboard Layout

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

        ████      ████                       Hopsule
        ████      ████                       The future of dev governance
            ████████                         
            ████████                         org: hopsule-inc  •  project: app
        ████          ████                   capture: ON  •  sync: ON  •  privacy: redacted
        ████          ████                   last sync: 12s  •  latency: 84ms
                                             
                                             v0.2.0
        ─────────────────────────────────────────────────────────────────────────────
        ✓ Connected

        Get started

        ❯ hopsule config         (Configure CLI settings)
          hopsule list           (List all decisions)
          hopsule create         (Create a new decision)
          hopsule get <id>       (Get decision details)
          hopsule accept <id>    (Accept a decision)
          hopsule deprecate <id> (Deprecate a decision)
          hopsule status         (Show project status)
          hopsule sync           (Sync with decision-api)

        API: https://api.hopsule.com
        Token: configured ✓

        ↑/↓: navigate  •  Enter: execute  •  q: quit  •  ?: help
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Layout Components

1. **Top Section**
   - ASCII art logo (left)
   - Project info (right)
     - Organization name
     - Project name
     - Capture/sync status
     - Last sync time & latency
     - Version

2. **Middle Section**
   - Connection status (✓ Connected / ⚠ Not configured)
   - Command list with descriptions
   - Selected command highlighted with `❯`

3. **Bottom Section**
   - API endpoint
   - Token status
   - Keyboard shortcuts reminder

## 🎬 Workflow Examples

### First Time Setup
```bash
# 1. Launch interactive dashboard
hopsule

# 2. Use ↓ to navigate to "hopsule config"
# 3. Press Enter to execute

# You'll be prompted for:
# - API URL (e.g., https://api.hopsule.com)
# - Project ID
# - Organization name
# - Auth token
```

### Daily Usage
```bash
# Launch dashboard
hopsule

# Navigate to desired command with ↑/↓
# Press Enter to execute

# Or use direct commands
hopsule list              # List decisions
hopsule create            # Create new decision
hopsule status            # Check status
```

### Getting Help
```bash
# In interactive mode, press ?
# Help screen shows:
# - All keyboard shortcuts
# - Available commands
# - Configuration info

# Press ? again to close help
```

## 🎨 Color Scheme

The TUI uses a carefully designed color scheme:

- **Accent (Green)**: Primary actions, success states
- **Title (Cyan)**: Section headers, important info
- **Info (Gray)**: Secondary information, descriptions
- **Warning (Yellow)**: Not configured, warnings
- **Error (Red)**: Errors, failed states

## 🚀 Advanced Features

### Vim-style Navigation
```bash
k  # Move up (same as ↑)
j  # Move down (same as ↓)
```

### Quick Quit
```bash
q         # Quit from main dashboard
Ctrl+C    # Force quit (emergency exit)
```

### Help Toggle
```bash
?  # Toggle help screen on/off
```

## 🐛 Troubleshooting

### TUI Not Showing
**Problem**: Plain text output instead of interactive TUI

**Solution**: Make sure you're in a real terminal (not a script or pipe)
```bash
# This works:
hopsule

# This won't show TUI:
echo | hopsule
```

### Colors Not Showing
**Problem**: No colors or broken display

**Solution**: Check terminal support
```bash
# Test your terminal
echo $TERM

# Most modern terminals support colors
# If not, try:
export TERM=xterm-256color
```

### Keyboard Not Working
**Problem**: Arrow keys not responding

**Solution**: Ensure terminal is in raw mode (automatic in real terminals)

## 💡 Tips & Tricks

1. **Fast Navigation**: Use `k`/`j` for Vim-style navigation
2. **Quick Config**: First time? Arrow down once and hit Enter for config
3. **Help Anytime**: Press `?` to see all shortcuts
4. **Clean Exit**: Always use `q` for graceful shutdown
5. **Status Check**: The top shows real-time sync status

## 🎯 Command Execution

When you press Enter on a selected command:

1. **Interactive Commands** (config, create):
   - Opens prompts for input
   - Guides you through the process

2. **Direct Commands** (list, status):
   - Executes immediately
   - Shows results
   - Returns to dashboard

3. **ID-Required Commands** (get, accept, deprecate):
   - Prompts for decision ID
   - Validates input
   - Executes action

## 📊 Status Indicators

| Indicator | Meaning |
|-----------|---------|
| ✓ Connected | API reachable, token valid |
| ⚠ Not configured | Need to run `hopsule config` |
| ON (green) | Feature active |
| OFF (red) | Feature inactive |
| configured ✓ | Token set |
| not set (yellow) | Token missing |

## 🎉 Best Practices

1. **Configure First**: Run config before other commands
2. **Use Interactive Mode**: Faster than typing commands
3. **Check Status**: Verify connection before operations
4. **Read Help**: Press `?` to discover features
5. **Clean Workflow**: Use TUI for exploration, direct commands for scripting

---

For more information, run `hopsule --help` or visit the [documentation](https://github.com/Hopsule/cli-tool).
