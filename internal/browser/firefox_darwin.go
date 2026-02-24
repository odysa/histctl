//go:build darwin

package browser

import "path/filepath"

var firefoxProfileBase = filepath.Join("Library", "Application Support", "Firefox", "Profiles")

const firefoxProcessName = "firefox"
