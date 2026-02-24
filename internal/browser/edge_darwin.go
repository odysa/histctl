//go:build darwin

package browser

import "path/filepath"

var edgeDBSubPath = filepath.Join("Library", "Application Support", "Microsoft Edge", "Default", "History")

const edgeProcessName = "Microsoft Edge"
