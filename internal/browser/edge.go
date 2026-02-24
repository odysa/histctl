package browser

import "path/filepath"

type Edge struct {
	*Chrome
}

func NewEdge(dbOverride string) *Edge {
	return &Edge{
		Chrome: &Chrome{
			name:        "edge",
			processName: "Microsoft Edge",
			dbSubPath:   filepath.Join("Library", "Application Support", "Microsoft Edge", "Default", "History"),
			dbOverride:  dbOverride,
		},
	}
}
