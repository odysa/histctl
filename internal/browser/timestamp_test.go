package browser

import (
	"testing"
	"time"
)

func TestWebKitToTime(t *testing.T) {
	tests := []struct {
		name   string
		webkit float64
		want   time.Time
	}{
		{
			name:   "epoch",
			webkit: 0,
			want:   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "known date",
			webkit: 738892800, // 2024-06-01 00:00:00 UTC
			want:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "fractional seconds",
			webkit: 738892800.5,
			want:   time.Date(2024, 6, 1, 0, 0, 0, 500000000, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebKitToTime(tt.webkit).UTC()
			if !got.Equal(tt.want) {
				t.Errorf("WebKitToTime(%v) = %v, want %v", tt.webkit, got, tt.want)
			}
		})
	}
}

func TestChromeToTime(t *testing.T) {
	// Chrome epoch: 1601-01-01 00:00:00 UTC, stored as microseconds
	// offset to Unix epoch: 11644473600 seconds
	tests := []struct {
		name string
		ts   int64
		want time.Time
	}{
		{
			name: "unix epoch",
			ts:   11644473600 * 1_000_000,
			want: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "known date",
			ts:   (1717200000 + 11644473600) * 1_000_000, // 2024-06-01 00:00:00 UTC
			want: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChromeToTime(tt.ts).UTC()
			if !got.Equal(tt.want) {
				t.Errorf("ChromeToTime(%v) = %v, want %v", tt.ts, got, tt.want)
			}
		})
	}
}

func TestFirefoxToTime(t *testing.T) {
	tests := []struct {
		name string
		ts   int64
		want time.Time
	}{
		{
			name: "unix epoch",
			ts:   0,
			want: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "known date",
			ts:   1717200000 * 1_000_000, // 2024-06-01 00:00:00 UTC
			want: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirefoxToTime(tt.ts).UTC()
			if !got.Equal(tt.want) {
				t.Errorf("FirefoxToTime(%v) = %v, want %v", tt.ts, got, tt.want)
			}
		})
	}
}
