# histctl

**Search, browse, and bulk-delete browser history from the terminal.**

Works with Chrome, Edge, Firefox, and Safari across macOS, Linux, and Windows.

![demo](demo/tui.gif)

## Features

- Interactive TUI with regex search, multi-select, and tab switching
- Scriptable CLI with JSON output for automation
- Automatic backups before every delete
- Cross-browser — one tool for all your browsers

## Install

### Homebrew (macOS / Linux)

```sh
brew install odysa/tap/histctl
```

### curl (macOS / Linux)

```sh
curl -fsSL https://raw.githubusercontent.com/odysa/histctl/main/install.sh | sh
```

### PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/odysa/histctl/main/install.ps1 | iex
```

### Go

```
go install github.com/odysa/histctl@latest
```

<details>
<summary>Build from source</summary>

```
git clone https://github.com/odysa/histctl.git
cd histctl
go build -o histctl .
```

</details>

## Usage

### TUI

Run `histctl` with no arguments to launch the interactive interface.

```
histctl
```

| Key | Action |
|-----|--------|
| `/` | Search with regex |
| `space` | Toggle selection |
| `a` | Select / deselect all |
| `d` | Delete selected |
| `tab` | Switch browser |
| `↑/k` `↓/j` | Navigate |
| `?` | Help |
| `q` | Quit |

Mouse scrolling and clicking browser tabs are also supported.

### CLI

```sh
# Search history
histctl list [pattern]         # regex search (case-insensitive)
histctl list -n 100            # limit results
histctl list --json            # JSON output

# Delete history
histctl delete <pattern>              # interactive confirmation
histctl delete <pattern> -d           # dry run — preview matches
histctl delete <pattern> -y           # skip confirmation
histctl delete <pattern> --no-backup  # skip backup
```

| Flag | Description |
|------|-------------|
| `-b, --browser` | Target browser: `safari\|chrome\|edge\|firefox\|all` (default: `all`) |
| `-n, --limit` | Max entries (default: `50`) |
| `--json` | JSON output |
| `-d, --dry-run` | Preview without deleting |
| `-y, --yes` | Skip confirmation |
| `--no-backup` | Skip backup |

## Notes

- Safari requires **Full Disk Access** for your terminal (System Settings > Privacy & Security > Full Disk Access)
- Backups are saved as `<db-path>.<timestamp>.bak` before each delete
- Browsers are auto-detected based on installed database files
- Close the target browser before deleting history
- Go 1.25+ required only for `go install` or building from source
