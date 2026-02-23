package browser

import (
	"context"
	"regexp"
	"testing"
	"time"
)

func newTestSafari(t *testing.T, rows []safariRow) *Safari {
	t.Helper()
	path := createTestDB(t, safariSchema)
	seedSafari(t, path, rows)
	return NewSafari(path)
}

var safariTestRows = []safariRow{
	{id: 1, url: strPtr("https://example.com"), title: "Example", visitTime: 738892800},     // 2024-06-01
	{id: 2, url: strPtr("https://golang.org"), title: "Go", visitTime: 738892800 + 86400},   // 2024-06-02
	{id: 3, url: strPtr("https://github.com"), title: "GitHub", visitTime: 738892800 + 3600}, // 2024-06-01 01:00
}

func TestSafariList(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	entries, err := s.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("List() returned %d entries, want 3", len(entries))
	}
	for _, e := range entries {
		if e.Browser != "safari" {
			t.Errorf("entry.Browser = %q, want %q", e.Browser, "safari")
		}
	}
}

func TestSafariListPattern(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`golang`)
	entries, err := s.List(ctx, ListOptions{Pattern: re})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("List() returned %d entries, want 1", len(entries))
	}
	if entries[0].URL != "https://golang.org" {
		t.Errorf("URL = %q, want %q", entries[0].URL, "https://golang.org")
	}
}

func TestSafariListLimit(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	entries, err := s.List(ctx, ListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("List() returned %d entries, want 2", len(entries))
	}
}

func TestSafariListSinceUntil(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	since := time.Date(2024, 6, 1, 0, 30, 0, 0, time.UTC)
	until := time.Date(2024, 6, 1, 23, 59, 59, 0, time.UTC)
	entries, err := s.List(ctx, ListOptions{Since: since, Until: until})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	// Only the 2024-06-01 01:00 entry (github.com) should match
	if len(entries) != 1 {
		t.Fatalf("List() returned %d entries, want 1", len(entries))
	}
	if entries[0].URL != "https://github.com" {
		t.Errorf("URL = %q, want %q", entries[0].URL, "https://github.com")
	}
}

func TestSafariListNullURL(t *testing.T) {
	rows := []safariRow{
		{id: 1, url: nil, title: "Null URL", visitTime: 738892800},
		{id: 2, url: strPtr("https://example.com"), title: "Valid", visitTime: 738892800},
	}
	s := newTestSafari(t, rows)
	ctx := context.Background()

	entries, err := s.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("List() returned %d entries, want 1 (null URL should be skipped)", len(entries))
	}
}

func TestSafariDelete(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`example\.com`)
	result, err := s.Delete(ctx, re, false)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	if result.Matched != 1 || result.Deleted != 1 {
		t.Errorf("Delete() = {Matched: %d, Deleted: %d}, want {1, 1}", result.Matched, result.Deleted)
	}

	// Verify it was actually deleted
	entries, _ := s.List(ctx, ListOptions{})
	if len(entries) != 2 {
		t.Errorf("after delete, List() returned %d entries, want 2", len(entries))
	}
}

func TestSafariDeleteDryRun(t *testing.T) {
	s := newTestSafari(t, safariTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`example\.com`)
	result, err := s.Delete(ctx, re, true)
	if err != nil {
		t.Fatalf("Delete(dryRun) error: %v", err)
	}
	if result.Matched != 1 {
		t.Errorf("Delete(dryRun) Matched = %d, want 1", result.Matched)
	}
	if result.Deleted != 0 {
		t.Errorf("Delete(dryRun) Deleted = %d, want 0", result.Deleted)
	}

	// Verify nothing was deleted
	entries, _ := s.List(ctx, ListOptions{})
	if len(entries) != 3 {
		t.Errorf("after dry run, List() returned %d entries, want 3", len(entries))
	}
}
