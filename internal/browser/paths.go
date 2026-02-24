package browser

import (
	"path/filepath"
	"runtime"
)

var knownBrowsers = []string{"chrome", "edge", "firefox"}

var (
	chromeDBSubPath    string
	chromeProcessName  string
	edgeDBSubPath      string
	edgeProcessName    string
	firefoxProfileBase string
	firefoxProcessName string
)

func init() {
	switch runtime.GOOS {
	case "darwin":
		chromeDBSubPath = filepath.Join("Library", "Application Support", "Google", "Chrome", "Default", "History")
		edgeDBSubPath = filepath.Join("Library", "Application Support", "Microsoft Edge", "Default", "History")
		firefoxProfileBase = filepath.Join("Library", "Application Support", "Firefox", "Profiles")
		chromeProcessName = "Google Chrome"
		edgeProcessName = "Microsoft Edge"
		firefoxProcessName = "firefox"
	case "windows":
		chromeDBSubPath = filepath.Join("AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")
		edgeDBSubPath = filepath.Join("AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "History")
		firefoxProfileBase = filepath.Join("AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
		chromeProcessName = "chrome.exe"
		edgeProcessName = "msedge.exe"
		firefoxProcessName = "firefox.exe"
	default: // linux
		chromeDBSubPath = filepath.Join(".config", "google-chrome", "Default", "History")
		edgeDBSubPath = filepath.Join(".config", "microsoft-edge", "Default", "History")
		firefoxProfileBase = ".mozilla/firefox"
		chromeProcessName = "chrome"
		edgeProcessName = "microsoft-edge"
		firefoxProcessName = "firefox"
	}
}
