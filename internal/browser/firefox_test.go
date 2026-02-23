package browser

import (
	"context"
	"regexp"
	"testing"
)

func newTestFirefox(t *testing.T, rows []firefoxRow) *Firefox {
	t.Helper()
	path := createTestDB(t, firefoxSchema)
	seedFirefox(t, path, rows)
	return NewFirefox(path)
}

// Firefox timestamps: microseconds since Unix epoch
var firefoxTestRows = []firefoxRow{
	{id: 1, url: strPtr("https://example.com"), title: "Example", visitDate: 1717200000 * 1_000_000},
	{id: 2, url: strPtr("https://golang.org"), title: "Go", visitDate: 1717286400 * 1_000_000},
	{id: 3, url: strPtr("https://github.com"), title: "GitHub", visitDate: 1717203600 * 1_000_000},
}

func TestFirefoxList(t *testing.T) {
	f := newTestFirefox(t, firefoxTestRows)
	ctx := context.Background()

	entries, err := f.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("List() returned %d entries, want 3", len(entries))
	}
	for _, e := range entries {
		if e.Browser != "firefox" {
			t.Errorf("entry.Browser = %q, want %q", e.Browser, "firefox")
		}
	}
}

func TestFirefoxListPattern(t *testing.T) {
	f := newTestFirefox(t, firefoxTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`golang`)
	entries, err := f.List(ctx, ListOptions{Pattern: re})
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

func TestFirefoxDelete(t *testing.T) {
	f := newTestFirefox(t, firefoxTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`example\.com`)
	result, err := f.Delete(ctx, re, false)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	if result.Matched != 1 || result.Deleted != 1 {
		t.Errorf("Delete() = {Matched: %d, Deleted: %d}, want {1, 1}", result.Matched, result.Deleted)
	}

	entries, _ := f.List(ctx, ListOptions{})
	if len(entries) != 2 {
		t.Errorf("after delete, List() returned %d entries, want 2", len(entries))
	}
}
