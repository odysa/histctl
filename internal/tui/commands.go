package tui

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/odysa/histctl/internal/backup"
	"github.com/odysa/histctl/internal/browser"
	"github.com/odysa/histctl/internal/process"
)

func (m Model) loadHistory() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		var all []browser.HistoryEntry
		for _, b := range m.browsers {
			entries, err := b.List(ctx, browser.ListOptions{Limit: 5000})
			if err != nil {
				continue
			}
			all = append(all, entries...)
		}
		sortEntries(all)
		return historyLoadedMsg{entries: all}
	}
}

func (m Model) searchHistory(pattern string) tea.Cmd {
	return func() tea.Msg {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return searchResultsMsg{err: fmt.Errorf("invalid regex: %s", pattern)}
		}
		ctx := context.Background()
		var all []browser.HistoryEntry
		for _, b := range m.browsers {
			entries, err := b.List(ctx, browser.ListOptions{Pattern: re})
			if err != nil {
				continue
			}
			all = append(all, entries...)
		}
		sortEntries(all)
		return searchResultsMsg{entries: all}
	}
}

func (m Model) performDelete() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		var totalResult browser.DeleteResult

		byBrowser := make(map[string][]browser.HistoryEntry)
		for idx := range m.selected {
			if idx < len(m.filteredEntries) {
				e := m.filteredEntries[idx]
				byBrowser[e.Browser] = append(byBrowser[e.Browser], e)
			}
		}

		for _, b := range m.browsers {
			entries, ok := byBrowser[b.Name()]
			if !ok {
				continue
			}

			running, err := process.IsRunning(b.ProcessName())
			if err != nil || running {
				return deleteResultMsg{err: fmt.Errorf("%s is running — close it first", b.Name())}
			}

			dbPath, err := b.DBPath()
			if err != nil {
				return deleteResultMsg{err: err}
			}
			if _, err := backup.Create(dbPath); err != nil {
				return deleteResultMsg{err: fmt.Errorf("backup failed: %w", err)}
			}

			urlPatterns := make([]string, len(entries))
			for i, e := range entries {
				urlPatterns[i] = regexp.QuoteMeta(e.URL)
			}
			pattern := regexp.MustCompile("^(" + strings.Join(urlPatterns, "|") + ")$")

			result, err := b.Delete(ctx, pattern, false)
			if err != nil {
				return deleteResultMsg{err: err}
			}
			totalResult.Matched += result.Matched
			totalResult.Deleted += result.Deleted
		}

		return deleteResultMsg{result: totalResult}
	}
}

func sortEntries(entries []browser.HistoryEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].VisitTime.After(entries[j].VisitTime)
	})
}
