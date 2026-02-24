package browser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Safari struct {
	dbOverride string
}

func NewSafari(dbOverride string) *Safari {
	return &Safari{dbOverride: dbOverride}
}

func (s *Safari) Name() string       { return "safari" }
func (s *Safari) ProcessName() string { return "Safari" }

func (s *Safari) DBPath() (string, error) {
	if s.dbOverride != "" {
		return s.dbOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(home, "Library", "Safari", "History.db")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("safari history not found: %w", err)
	}
	return p, nil
}

const safariListQuery = `
	SELECT hi.id, hi.url, hv.title, hv.visit_time
	FROM history_items hi
	JOIN history_visits hv ON hv.history_item = hi.id
	ORDER BY hv.visit_time DESC`

func safariScanRow(rows *sql.Rows) (HistoryEntry, bool, error) {
	var id int64
	var url sql.NullString
	var title sql.NullString
	var visitTime sql.NullFloat64

	if err := rows.Scan(&id, &url, &title, &visitTime); err != nil {
		return HistoryEntry{}, false, err
	}
	if !url.Valid {
		return HistoryEntry{}, false, nil
	}
	return HistoryEntry{
		URL:       url.String,
		Title:     title.String,
		VisitTime: WebKitToTime(visitTime.Float64),
		Browser:   "safari",
		ItemID:    id,
	}, true, nil
}

func (s *Safari) List(ctx context.Context, opts ListOptions) ([]HistoryEntry, error) {
	dbPath, err := s.DBPath()
	if err != nil {
		return nil, err
	}
	entries, err := listEntries(ctx, dbPath, "safari", safariListQuery, safariScanRow, opts)
	if err != nil && strings.Contains(err.Error(), "unable to open database") {
		return nil, fmt.Errorf("cannot read Safari history — grant Full Disk Access to your terminal in System Settings > Privacy & Security > Full Disk Access")
	}
	return entries, err
}

func (s *Safari) Delete(ctx context.Context, pattern *regexp.Regexp, dryRun bool) (DeleteResult, error) {
	dbPath, err := s.DBPath()
	if err != nil {
		return DeleteResult{}, err
	}
	entries, err := s.List(ctx, ListOptions{Pattern: pattern})
	if err != nil {
		return DeleteResult{}, err
	}
	return deleteEntries(ctx, dbPath, "safari", entries, dryRun, func(ctx context.Context, tx *sql.Tx, e HistoryEntry) error {
		if _, err := tx.ExecContext(ctx, "DELETE FROM history_visits WHERE history_item = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete visits: %w", err)
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM history_items WHERE id = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete items: %w", err)
		}
		return nil
	})
}
