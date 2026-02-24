package browser

import (
	"context"
	"regexp"
	"time"
)

// HistoryEntry is a single history visit across all browsers.
type HistoryEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitTime  time.Time `json:"visit_time"`
	VisitCount int       `json:"visit_count,omitempty"`
	Browser    string    `json:"browser"`
	ItemID     int64     `json:"-"` // internal: used for deletion
}

// ListOptions controls filtering when listing history.
type ListOptions struct {
	Pattern *regexp.Regexp
	Limit   int
	Since   time.Time
	Until   time.Time
}

// DeleteResult summarizes a delete operation.
type DeleteResult struct {
	Matched int
	Deleted int
}

// Browser is the interface every browser backend must implement.
type Browser interface {
	Name() string
	DBPath() (string, error)
	ProcessName() string
	List(ctx context.Context, opts ListOptions) ([]HistoryEntry, error)
	Delete(ctx context.Context, pattern *regexp.Regexp, dryRun bool) (DeleteResult, error)
}
