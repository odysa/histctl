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
Sleep 800ms

# ── Navigate down ────────────────────────────────────────────────────
Down
Sleep 150ms
Down
Sleep 150ms
Down
Sleep 150ms
Down
Sleep 150ms
Down
Sleep 150ms
Down
Sleep 500ms

# ── Switch browser tabs ─────────────────────────────────────────────
Tab
Sleep 600ms
Tab
Sleep 600ms
Tab
Sleep 600ms
Tab
Sleep 800ms

# ── Search ───────────────────────────────────────────────────────────
Type "/"
Sleep 300ms
Type "github"
Sleep 500ms
Enter
Sleep 1s

# ── Select entries ───────────────────────────────────────────────────
Space
Sleep 300ms
Space
Sleep 300ms
Space
Sleep 500ms

# ── Select all ───────────────────────────────────────────────────────
Type "a"
Sleep 800ms

# ── Delete confirmation popup ────────────────────────────────────────
Type "d"
Sleep 1.2s

# ── Cancel delete ────────────────────────────────────────────────────
Escape
Sleep 800ms

# ── Deselect all ─────────────────────────────────────────────────────
Type "a"
Sleep 500ms

# ── Clear search ─────────────────────────────────────────────────────
Type "/"
Sleep 200ms
Ctrl+U
Enter
Sleep 800ms

# ── Show full help ───────────────────────────────────────────────────
Type "?"
Sleep 1s
Type "?"
Sleep 500ms

# ── Quit ─────────────────────────────────────────────────────────────
Type "q"
Sleep 300ms
