//go:build linux

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join(".config", "google-chrome", "Default", "History")

const chromeProcessName = "chrome"
