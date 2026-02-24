package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/odysa/histctl/internal/browser"
)

type state int

const (
	stateViewing state = iota
	stateSearching
	stateLoading
	stateConfirmDelete
)

type historyLoadedMsg struct {
	entries []browser.HistoryEntry
	err     error
}

type deleteResultMsg struct {
	result browser.DeleteResult
	err    error
}

type searchResultsMsg struct {
	entries []browser.HistoryEntry
	err     error
}

type Model struct {
	browsers      []browser.Browser
	activeBrowser int // index into browserNames; 0 = all
	browserNames  []string

	allEntries      []browser.HistoryEntry
	searchEntries   []browser.HistoryEntry // non-nil when DB search is active
	filteredEntries []browser.HistoryEntry
	selected        map[int]bool // index in filteredEntries

	table       table.Model
	searchInput textinput.Model
	spinner     spinner.Model
	help        help.Model
	keys        KeyMap

	state      state
	searchText string
	width      int
	height     int
	err        error
	statusMsg  string
	showHelp   bool
}

func NewModel(browsers []browser.Browser) Model {
	si := textinput.New()
	si.Placeholder = "regex pattern..."
	si.PromptStyle = SearchPromptStyle
	si.Prompt = "/ "
	si.CharLimit = 200

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(Accent)

	names := []string{"all"}
	for _, b := range browsers {
		names = append(names, b.Name())
	}

	t := table.New(
		table.WithFocused(true),
		table.WithHeight(20),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Muted).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(DimBg)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(Accent).
		Bold(false)
	s.Cell = s.Cell.
		Foreground(lipgloss.Color("#C0CAF5"))
	t.SetStyles(s)

	return Model{
		browsers:      browsers,
		activeBrowser: 0,
		browserNames:  names,
		selected:      make(map[int]bool),
		table:         t,
		searchInput:   si,
		spinner:       sp,
		help:          help.New(),
		keys:          DefaultKeyMap(),
		state:         stateLoading,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadHistory(),
	)
}

func (m *Model) applyFilters() {
	m.filteredEntries = m.filteredEntries[:0]

	source := m.allEntries
	if m.searchEntries != nil {
		source = m.searchEntries
	}

	activeName := m.browserNames[m.activeBrowser]

	for _, e := range source {
		if activeName != "all" && e.Browser != activeName {
			continue
		}
		m.filteredEntries = append(m.filteredEntries, e)
	}

	m.updateTableRows()
}

func (m *Model) updateTableRows() {
	rows := make([]table.Row, len(m.filteredEntries))
	for i, e := range m.filteredEntries {
		urlW := m.urlWidth()
		url := truncate(e.URL, urlW)
		if m.selected[i] {
			url = "> " + truncate(e.URL, urlW-2)
		}
		title := truncate(e.Title, 30)
		rows[i] = table.Row{
			url,
			title,
			relativeTime(e.VisitTime),
			e.Browser,
		}
	}
	m.table.SetRows(rows)
}

func (m *Model) resizeTable() {
	h := m.height - 12
	if h < 5 {
		h = 5
	}
	m.table.SetHeight(h)
	m.table.SetColumns(m.columns())
	m.updateTableRows()
}

func (m *Model) urlWidth() int {
	w := m.width - 30 - 12 - 8 - 10
	if w < 20 {
		w = 20
	}
	return w
}

func (m *Model) columns() []table.Column {
	return []table.Column{
		{Title: "URL", Width: m.urlWidth()},
		{Title: "Title", Width: 30},
		{Title: "Time", Width: 12},
		{Title: "Browser", Width: 8},
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}

func relativeTime(t time.Time) string {
	now := time.Now()
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 48*time.Hour:
		return "yesterday"
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case t.Year() == now.Year():
		return t.Local().Format("Jan 02")
	default:
		return t.Local().Format("Jan 2006")
	}
}

func (m Model) browserCount(name string) int {
	source := m.allEntries
	if m.searchEntries != nil {
		source = m.searchEntries
	}
	if name == "all" {
		return len(source)
	}
	count := 0
	for _, e := range source {
		if e.Browser == name {
			count++
		}
	}
	return count
}

func Run(browsers []browser.Browser) error {
	m := NewModel(browsers)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
