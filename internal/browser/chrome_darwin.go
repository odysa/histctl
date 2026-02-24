//go:build darwin

package browser

import "path/filepath"

var chromeDBSubPath = filepath.Join("Library", "Application Support", "Google", "Chrome", "Default", "History")

const chromeProcessName = "Google Chrome"
