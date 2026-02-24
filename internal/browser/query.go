package browser

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// rowScanner converts a database row into a HistoryEntry.
// Returns false if the row should be skipped.
type rowScanner func(*sql.Rows) (HistoryEntry, bool, error)

// rowDeleter deletes a single entry within a transaction.
type rowDeleter func(ctx context.Context, tx *sql.Tx, e HistoryEntry) error

func listEntries(ctx context.Context, dbPath, name, query string, scan rowScanner, opts ListOptions) ([]HistoryEntry, error) {
	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("open %s db: %w", name, err)
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query %s history: %w", name, err)
	}
	defer rows.Close()

	var entries []HistoryEntry
	for rows.Next() {
		entry, ok, err := scan(rows)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if opts.Pattern != nil && !opts.Pattern.MatchString(entry.URL) {
			continue
		}
		if !opts.Since.IsZero() && entry.VisitTime.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && entry.VisitTime.After(opts.Until) {
			continue
		}
		entries = append(entries, entry)
		if opts.Limit > 0 && len(entries) >= opts.Limit {
			break
		}
	}
	return entries, rows.Err()
}

func deleteEntries(ctx context.Context, dbPath, name string, entries []HistoryEntry, dryRun bool, del rowDeleter) (DeleteResult, error) {
	result := DeleteResult{Matched: len(entries)}
	if dryRun || len(entries) == 0 {
		return result, nil
	}

	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=rw")
	if err != nil {
		return result, fmt.Errorf("open %s db for writing: %w", name, err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	for _, e := range entries {
		if err := del(ctx, tx, e); err != nil {
			return result, err
		}
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}
	result.Deleted = len(entries)
	return result, nil
}
