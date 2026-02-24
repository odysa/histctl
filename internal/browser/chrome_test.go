package browser

import (
	"context"
	"regexp"
	"testing"
)

func newTestChrome(t *testing.T, rows []chromeRow) *Chrome {
	t.Helper()
	path := createTestDB(t, chromeSchema)
	seedChrome(t, path, rows)
	return NewChrome(path)
}

// Chrome timestamps: (unixSeconds + 11644473600) * 1_000_000
var chromeTestRows = []chromeRow{
	{id: 1, url: strPtr("https://example.com"), title: "Example", visitTime: (1717200000 + 11644473600) * 1_000_000},
	{id: 2, url: strPtr("https://golang.org"), title: "Go", visitTime: (1717286400 + 11644473600) * 1_000_000},
	{id: 3, url: strPtr("https://github.com"), title: "GitHub", visitTime: (1717203600 + 11644473600) * 1_000_000},
}

func TestChromeList(t *testing.T) {
	c := newTestChrome(t, chromeTestRows)
	ctx := context.Background()

	entries, err := c.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("List() returned %d entries, want 3", len(entries))
	}
	for _, e := range entries {
		if e.Browser != "chrome" {
			t.Errorf("entry.Browser = %q, want %q", e.Browser, "chrome")
		}
	}
}

func TestChromeListPattern(t *testing.T) {
	c := newTestChrome(t, chromeTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`golang`)
	entries, err := c.List(ctx, ListOptions{Pattern: re})
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

func TestChromeDelete(t *testing.T) {
	c := newTestChrome(t, chromeTestRows)
	ctx := context.Background()

	re := regexp.MustCompile(`example\.com`)
	result, err := c.Delete(ctx, re, false)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	if result.Matched != 1 || result.Deleted != 1 {
		t.Errorf("Delete() = {Matched: %d, Deleted: %d}, want {1, 1}", result.Matched, result.Deleted)
	}

	entries, _ := c.List(ctx, ListOptions{})
	if len(entries) != 2 {
		t.Errorf("after delete, List() returned %d entries, want 2", len(entries))
	}
}

func TestEdgeName(t *testing.T) {
	e := NewEdge("")
	if e.Name() != "edge" {
		t.Errorf("Edge.Name() = %q, want %q", e.Name(), "edge")
	}
	if e.ProcessName() != edgeProcessName {
		t.Errorf("Edge.ProcessName() = %q, want %q", e.ProcessName(), edgeProcessName)
	}
}

func TestEdgeList(t *testing.T) {
	// Edge reuses Chrome's schema and logic — verify it works with the correct browser name
	path := createTestDB(t, chromeSchema)
	seedChrome(t, path, chromeTestRows)
	e := &Edge{Chrome: &Chrome{
		name:       "edge",
		dbOverride: path,
	}}
	ctx := context.Background()

	entries, err := e.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("Edge.List() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("Edge.List() returned %d entries, want 3", len(entries))
	}
	for _, entry := range entries {
		if entry.Browser != "edge" {
			t.Errorf("entry.Browser = %q, want %q", entry.Browser, "edge")
		}
	}
}
