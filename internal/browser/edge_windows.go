//go:build windows

package browser

import "path/filepath"

var edgeDBSubPath = filepath.Join("AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "History")

const edgeProcessName = "msedge.exe"
