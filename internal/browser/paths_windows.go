//go:build windows

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join("AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")
var edgeDBSubPath = filepath.Join("AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "History")
var firefoxProfileBase = filepath.Join("AppData", "Roaming", "Mozilla", "Firefox", "Profiles")

const (
	chromeProcessName  = "chrome.exe"
	edgeProcessName    = "msedge.exe"
	firefoxProcessName = "firefox.exe"
)
