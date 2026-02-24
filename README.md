# histctl

A terminal tool to search, browse, and delete browser history across Chrome, Edge, Firefox, and Safari. Comes with an interactive TUI and a scriptable CLI.

Supports macOS, Linux, and Windows. Safari is macOS-only.

## Install

### curl (macOS / Linux)

```
curl -fsSL "https://github.com/odysa/histctl/releases/latest/download/histctl-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')" -o histctl
chmod +x histctl
sudo mv histctl /usr/local/bin/
```

### PowerShell (Windows)

```powershell
Invoke-WebRequest -Uri "https://github.com/odysa/histctl/releases/latest/download/histctl-windows-amd64.exe" -OutFile histctl.exe
Move-Item histctl.exe "$env:LOCALAPPDATA\Microsoft\WindowsApps\histctl.exe"
```

### Go

```
go install github.com/odysa/histctl@latest
```

### Build from source

```
git clone https://github.com/odysa/histctl.git
cd histctl
go build -o histctl .
```

## TUI

Run `histctl` with no arguments to launch the interactive interface.

```
histctl
```

### Keybindings

| Key | Action |
|-----|--------|
| `/` | Search with regex |
| `space` | Toggle selection on current row |
| `a` | Select / deselect all |
| `d` | Delete selected entries |
| `tab` | Switch browser |
| `↑/k` `↓/j` | Navigate rows |
| `?` | Toggle full help |
| `q` / `Ctrl+C` | Quit |

Mouse scrolling and clicking browser tabs are also supported.

Deleting entries shows a confirmation popup overlaid on the table. Press `y` to confirm or any other key to cancel.

## CLI

### List history

```
histctl list [pattern]         # regex search (case-insensitive)
histctl list -n 100            # limit results
histctl list --json            # JSON output
```

### Delete history

```
histctl delete <pattern>              # interactive confirmation
histctl delete <pattern> -d           # dry run — preview matches
histctl delete <pattern> -y           # skip confirmation
histctl delete <pattern> --no-backup  # skip backup
```

### Flags

| Flag | Description |
|------|-------------|
| `-b, --browser` | Target browser: `safari\|chrome\|edge\|firefox\|all` (default: `all`) |
| `-n, --limit` | Max entries to display (default: `50`) |
| `--json` | Output as JSON |
| `-d, --dry-run` | Preview matches without deleting |
| `-y, --yes` | Skip confirmation prompt |
| `--no-backup` | Skip creating a backup before delete |

## Requirements

- macOS, Linux, or Windows
- Go 1.25+ (only for `go install` or building from source)
- Safari history (macOS only) requires **Full Disk Access** for your terminal (System Settings > Privacy & Security > Full Disk Access)
- Close the target browser before deleting history

## Notes

- Backups are saved as `<db-path>.<timestamp>.bak` before each delete (unless `--no-backup`)
- All regex patterns are case-insensitive
- Supported browsers are auto-detected based on installed database files
- Safari is only available on macOS
