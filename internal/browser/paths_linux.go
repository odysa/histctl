//go:build linux

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join(".config", "google-chrome", "Default", "History")
var edgeDBSubPath = filepath.Join(".config", "microsoft-edge", "Default", "History")
var firefoxProfileBase = ".mozilla/firefox"

const (
	chromeProcessName  = "chrome"
	edgeProcessName    = "microsoft-edge"
	firefoxProcessName = "firefox"
)
