Output demo/tui.gif

Set Shell "bash"
Set Theme "Catppuccin Mocha"
Set FontSize 14
Set Width 1100
Set Height 650
Set Padding 15
Set TypingSpeed 30ms
Set PlaybackSpeed 1.5

# Setup: add to PATH, create fake browser history
Hide
Type "export PATH=$PWD:$PATH && source demo/setup.sh && clear"
Enter
Sleep 2s
Show

# ── Launch TUI ───────────────────────────────────────────────────────
Type "histctl"
Sleep 300ms
Enter
Sleep 1s

# ── Navigate down (shows cursor position in status bar) ─────────────
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 600ms

# ── Switch browser tabs (shows entry counts) ────────────────────────
Tab
Sleep 700ms
Tab
Sleep 700ms
Tab
Sleep 700ms
Tab
Sleep 900ms

# ── Search with no results (shows empty state) ─────────────────────
Type "/"
Sleep 300ms
Type "nonexistent"
Sleep 500ms
Enter
Sleep 1.2s

# ── Clear search and search for github ──────────────────────────────
Type "/"
Sleep 200ms
Ctrl+U
Sleep 200ms
Type "github"
Sleep 500ms
Enter
Sleep 1s

# ── Select entries (shows > prefix on URL) ──────────────────────────
Space
Sleep 300ms
Space
Sleep 300ms
Space
Sleep 600ms

# ── Select all ──────────────────────────────────────────────────────
Type "a"
Sleep 800ms

# ── Delete confirmation popup ───────────────────────────────────────
Type "d"
Sleep 1.2s

# ── Cancel delete ───────────────────────────────────────────────────
Escape
Sleep 800ms

# ── Deselect all ────────────────────────────────────────────────────
Type "a"
Sleep 500ms

# ── Clear search ────────────────────────────────────────────────────
Type "/"
Sleep 200ms
Ctrl+U
Enter
Sleep 800ms

# ── Show full help ──────────────────────────────────────────────────
Type "?"
Sleep 1s
Type "?"
Sleep 500ms

# ── Quit ────────────────────────────────────────────────────────────
Type "q"
Sleep 300ms
