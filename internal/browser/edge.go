package browser

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
