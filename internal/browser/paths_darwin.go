//go:build darwin

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join("Library", "Application Support", "Google", "Chrome", "Default", "History")
var edgeDBSubPath = filepath.Join("Library", "Application Support", "Microsoft Edge", "Default", "History")
var firefoxProfileBase = filepath.Join("Library", "Application Support", "Firefox", "Profiles")

const (
	chromeProcessName  = "Google Chrome"
	edgeProcessName    = "Microsoft Edge"
	firefoxProcessName = "firefox"
)
