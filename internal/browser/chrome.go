package browser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Chrome struct {
	name        string
	processName string
	dbSubPath   string
	dbOverride  string
}

func NewChrome(dbOverride string) *Chrome {
	return &Chrome{
		name:        "chrome",
		processName: chromeProcessName,
		dbSubPath:   chromeDBSubPath,
		dbOverride:  dbOverride,
	}
}

func (c *Chrome) Name() string       { return c.name }
func (c *Chrome) ProcessName() string { return c.processName }

func (c *Chrome) DBPath() (string, error) {
	if c.dbOverride != "" {
		return c.dbOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(home, c.dbSubPath)
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("%s history not found: %w", c.name, err)
	}
	return p, nil
}

const chromeListQuery = `
	SELECT u.id, u.url, u.title, v.visit_time
	FROM urls u
	JOIN visits v ON v.url = u.id
	ORDER BY v.visit_time DESC`

func chromeScanRow(name string) rowScanner {
	return func(rows *sql.Rows) (HistoryEntry, bool, error) {
		var id int64
		var url sql.NullString
		var title sql.NullString
		var visitTime sql.NullInt64

		if err := rows.Scan(&id, &url, &title, &visitTime); err != nil {
			return HistoryEntry{}, false, err
		}
		if !url.Valid {
			return HistoryEntry{}, false, nil
		}
		return HistoryEntry{
			URL:       url.String,
			Title:     title.String,
			VisitTime: ChromeToTime(visitTime.Int64),
			Browser:   name,
			ItemID:    id,
		}, true, nil
	}
}

func (c *Chrome) List(ctx context.Context, opts ListOptions) ([]HistoryEntry, error) {
	dbPath, err := c.DBPath()
	if err != nil {
		return nil, err
	}
	return listEntries(ctx, dbPath, c.name, chromeListQuery, chromeScanRow(c.name), opts)
}

func (c *Chrome) Delete(ctx context.Context, pattern *regexp.Regexp, dryRun bool) (DeleteResult, error) {
	dbPath, err := c.DBPath()
	if err != nil {
		return DeleteResult{}, err
	}
	entries, err := c.List(ctx, ListOptions{Pattern: pattern})
	if err != nil {
		return DeleteResult{}, err
	}
	return deleteEntries(ctx, dbPath, c.name, entries, dryRun, func(ctx context.Context, tx *sql.Tx, e HistoryEntry) error {
		if _, err := tx.ExecContext(ctx, "DELETE FROM visits WHERE url = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete visits: %w", err)
		}
		// keyword_search_terms may not exist in all versions
		tx.ExecContext(ctx, "DELETE FROM keyword_search_terms WHERE url_id = ?", e.ItemID)
		if _, err := tx.ExecContext(ctx, "DELETE FROM urls WHERE id = ?", e.ItemID); err != nil {
			return fmt.Errorf("delete urls: %w", err)
		}
		return nil
	})
}
