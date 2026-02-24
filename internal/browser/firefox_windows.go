//go:build windows

package browser

import "path/filepath"

var firefoxProfileBase = filepath.Join("AppData", "Roaming", "Mozilla", "Firefox", "Profiles")

const firefoxProcessName = "firefox.exe"
