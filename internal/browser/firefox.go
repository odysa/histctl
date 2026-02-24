package browser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Firefox struct {
	dbOverride string
}

func NewFirefox(dbOverride string) *Firefox {
	return &Firefox{dbOverride: dbOverride}
}

func (f *Firefox) Name() string       { return "firefox" }
func (f *Firefox) ProcessName() string { return firefoxProcessName }

func (f *Firefox) DBPath() (string, error) {
	if f.dbOverride != "" {
		return f.dbOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	profilesDir := filepath.Join(home, firefoxProfileBase)
	matches, err := filepath.Glob(filepath.Join(profilesDir, "*.default-release"))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		matches, _ = filepath.Glob(filepath.Join(profilesDir, "*.default"))
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no firefox profile found in %s", profilesDir)
	}
	p := filepath.Join(matches[0], "places.sqlite")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("firefox places.sqlite not found: %w", err)
	}
	return p, nil
}

const firefoxListQuery = `
	SELECT p.id, p.url, p.title, v.visit_date
	FROM moz_places p
	JOIN moz_historyvisits v ON v.place_id = p.id
	ORDER BY v.visit_date DESC`

func firefoxScanRow(rows *sql.Rows) (HistoryEntry, bool, error) {
	var id int64
	var url sql.NullString
	var title sql.NullString
	var visitDate sql.NullInt64

	if err := rows.Scan(&id, &url, &title, &visitDate); err != nil {
		return HistoryEntry{}, false, err
	}
	if !url.Valid {
		return HistoryEntry{}, false, nil
	}
	return HistoryEntry{
		URL:       url.String,
		Title:     title.String,
		VisitTime: FirefoxToTime(visitDate.Int64),
		Browser:   "firefox",
		ItemID:    id,
	}, true, nil
}

func (f *Firefox) List(ctx context.Context, opts ListOptions) ([]HistoryEntry, error) {
	dbPath, err := f.DBPath()
	if err != nil {
		return nil, err
	}
	return listEntries(ctx, dbPath, "firefox", firefoxListQuery, firefoxScanRow, opts)
}

func (f *Firefox) Delete(ctx context.Context, pattern *regexp.Regexp, dryRun bool) (DeleteResult, error) {
	dbPath, err := f.DBPath()
	if err != nil {
		return DeleteResult{}, err
	}
	entries, err := f.List(ctx, ListOptions{Pattern: pattern})
	if err != nil {
		return DeleteResult{}, err
	}
	return deleteEntries(ctx, dbPath, "firefox", entries, dryRun, func(ctx context.Context, tx *sql.Tx, e HistoryEntry) error {
		if _, err := tx.ExecContext(ctx, "DELETE FROM moz_historyvisits WHERE place_id = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete visits: %w", err)
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM moz_places WHERE id = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete places: %w", err)
		}
		return nil
	})
}
