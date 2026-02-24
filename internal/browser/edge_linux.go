//go:build linux

package browser

import "path/filepath"

var edgeDBSubPath = filepath.Join(".config", "microsoft-edge", "Default", "History")

const edgeProcessName = "microsoft-edge"
