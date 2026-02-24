package browser

// edgeDBSubPath and edgeProcessName are defined in edge_{darwin,linux,windows}.go

type Edge struct {
	*Chrome
}

func NewEdge(dbOverride string) *Edge {
	return &Edge{
		Chrome: &Chrome{
			name:        "edge",
			processName: edgeProcessName,
			dbSubPath:   edgeDBSubPath,
			dbOverride:  dbOverride,
		},
	}
}
