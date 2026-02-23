# histctl

A CLI tool to search, visualize, and delete browser history across Safari, Chrome, Edge, and Firefox on macOS.

## Install

```
go install github.com/odysa/histctl@latest
```

Or build from source:

```
git clone https://github.com/odysa/histctl.git
cd histctl
go build -o histctl .
```

## Usage

Launch the interactive TUI:

```
histctl
```

### List history

```
histctl list [pattern]         # regex search (case-insensitive)
histctl list -n 100            # limit results
histctl list --json            # JSON output
```

### Delete history

```
histctl delete <pattern>       # interactive confirmation
histctl delete <pattern> -d    # dry run — preview matches
histctl delete <pattern> -y    # skip confirmation
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

- macOS
- Go 1.25+
- Safari history requires **Full Disk Access** for your terminal (System Settings > Privacy & Security > Full Disk Access)
- Close the target browser before deleting history

## Notes

- Backups are saved as `<db-path>.<timestamp>.bak` before each delete (unless `--no-backup`)
- All regex patterns are case-insensitive
- The TUI mode shows history from all installed browsers
