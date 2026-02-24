//go:build windows

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join("AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")

const chromeProcessName = "chrome.exe"
