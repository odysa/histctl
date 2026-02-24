#!/bin/bash
# Creates a temporary HOME with fake browser history databases for the demo.
set -euo pipefail

DEMO_HOME=$(mktemp -d)
export HOME="$DEMO_HOME"

# ── Chrome ───────────────────────────────────────────────────────────
CHROME_DIR="$DEMO_HOME/Library/Application Support/Google/Chrome/Default"
mkdir -p "$CHROME_DIR"
CHROME_DB="$CHROME_DIR/History"

sqlite3 "$CHROME_DB" <<'SQL'
CREATE TABLE urls (id INTEGER PRIMARY KEY, url TEXT, title TEXT, visit_count INTEGER DEFAULT 0);
CREATE TABLE visits (id INTEGER PRIMARY KEY, url INTEGER REFERENCES urls(id), visit_time INTEGER);
INSERT INTO urls VALUES(1, 'https://github.com/charmbracelet/bubbletea', 'Bubble Tea - TUI Framework', 0);
INSERT INTO urls VALUES(2, 'https://go.dev/doc/install', 'Download and install - The Go Programming Language', 0);
INSERT INTO urls VALUES(3, 'https://news.ycombinator.com/', 'Hacker News', 0);
INSERT INTO urls VALUES(4, 'https://stackoverflow.com/questions/golang-channels', 'How do Go channels work? - Stack Overflow', 0);
INSERT INTO urls VALUES(5, 'https://pkg.go.dev/github.com/spf13/cobra', 'cobra package - github.com/spf13/cobra', 0);
INSERT INTO urls VALUES(6, 'https://github.com/odysa/histctl', 'histctl - Browser history manager', 0);
INSERT INTO urls VALUES(7, 'https://example.com/login', 'Login Page', 0);
INSERT INTO urls VALUES(8, 'https://example.com/dashboard', 'Dashboard - Example App', 0);
-- Chrome timestamps: (unixSec + 11644473600) * 1000000
INSERT INTO visits VALUES(1,  1, 13416303600000000);
INSERT INTO visits VALUES(2,  2, 13416233400000000);
INSERT INTO visits VALUES(3,  3, 13416128100000000);
INSERT INTO visits VALUES(4,  4, 13416068700000000);
INSERT INTO visits VALUES(5,  5, 13415962800000000);
INSERT INTO visits VALUES(6,  6, 13415876400000000);
INSERT INTO visits VALUES(7,  7, 13415790000000000);
INSERT INTO visits VALUES(8,  8, 13415703600000000);
SQL

# ── Safari ───────────────────────────────────────────────────────────
SAFARI_DIR="$DEMO_HOME/Library/Safari"
mkdir -p "$SAFARI_DIR"
SAFARI_DB="$SAFARI_DIR/History.db"

sqlite3 "$SAFARI_DB" <<'SQL'
CREATE TABLE history_items (id INTEGER PRIMARY KEY, url TEXT);
CREATE TABLE history_visits (id INTEGER PRIMARY KEY, history_item INTEGER REFERENCES history_items(id), title TEXT, visit_time REAL);
INSERT INTO history_items VALUES(1, 'https://developer.apple.com/swift/');
INSERT INTO history_items VALUES(2, 'https://www.rust-lang.org/learn');
INSERT INTO history_items VALUES(3, 'https://github.com/odysa/histctl/issues');
INSERT INTO history_items VALUES(4, 'https://docs.github.com/en/actions');
INSERT INTO history_items VALUES(5, 'https://claude.ai/');
INSERT INTO history_items VALUES(6, 'https://example.com/settings');
-- Safari timestamps: unixSec - 978307200
INSERT INTO history_visits VALUES(1, 1, 'Swift - Apple Developer', 793522800);
INSERT INTO history_visits VALUES(2, 2, 'Learn Rust', 793452600);
INSERT INTO history_visits VALUES(3, 3, 'Issues - odysa/histctl', 793347300);
INSERT INTO history_visits VALUES(4, 4, 'GitHub Actions Documentation', 793287900);
INSERT INTO history_visits VALUES(5, 5, 'Claude', 793182000);
INSERT INTO history_visits VALUES(6, 6, 'Settings - Example App', 793095600);
SQL

# ── Firefox ──────────────────────────────────────────────────────────
FF_PROFILE="$DEMO_HOME/Library/Application Support/Firefox/Profiles/demo.default-release"
mkdir -p "$FF_PROFILE"
FF_DB="$FF_PROFILE/places.sqlite"

sqlite3 "$FF_DB" <<'SQL'
CREATE TABLE moz_places (id INTEGER PRIMARY KEY, url TEXT, title TEXT);
CREATE TABLE moz_historyvisits (id INTEGER PRIMARY KEY, place_id INTEGER REFERENCES moz_places(id), visit_date INTEGER);
INSERT INTO moz_places VALUES(1, 'https://www.mozilla.org/en-US/firefox/', 'Firefox Browser');
INSERT INTO moz_places VALUES(2, 'https://crates.io/', 'crates.io: Rust Package Registry');
INSERT INTO moz_places VALUES(3, 'https://github.com/trending', 'Trending repositories on GitHub');
INSERT INTO moz_places VALUES(4, 'https://en.wikipedia.org/wiki/Unix', 'Unix - Wikipedia');
INSERT INTO moz_places VALUES(5, 'https://example.com/api/docs', 'API Documentation - Example');
-- Firefox timestamps: unixSec * 1000000
INSERT INTO moz_historyvisits VALUES(1, 1, 1771830000000000);
INSERT INTO moz_historyvisits VALUES(2, 2, 1771759800000000);
INSERT INTO moz_historyvisits VALUES(3, 3, 1771654500000000);
INSERT INTO moz_historyvisits VALUES(4, 4, 1771595100000000);
INSERT INTO moz_historyvisits VALUES(5, 5, 1771489200000000);
SQL

# No output — sourced silently by demo tape
