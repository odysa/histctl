package tui

import (
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

// Messages
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
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(Accent)

	names := []string{"all"}
	for _, b := range browsers {
		names = append(names, b.Name())
	}

	columns := []table.Column{
		{Title: " ", Width: 3},
		{Title: "URL", Width: 50},
		{Title: "Title", Width: 30},
		{Title: "Time", Width: 19},
		{Title: "Browser", Width: 8},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(Accent)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("57")).
		Bold(false)
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
		check := "  "
		if m.selected[i] {
			check = "✓ "
		}
		url := truncate(e.URL, m.urlWidth())
		title := truncate(e.Title, 30)
		rows[i] = table.Row{
			check,
			url,
			title,
			e.VisitTime.Local().Format("2006-01-02 15:04"),
			e.Browser,
		}
	}
	m.table.SetRows(rows)
}

func (m *Model) resizeTable() {
	h := m.height - 10
	if h < 5 {
		h = 5
	}
	m.table.SetHeight(h)

	urlW := m.urlWidth()
	columns := []table.Column{
		{Title: " ", Width: 3},
		{Title: "URL", Width: urlW},
		{Title: "Title", Width: 30},
		{Title: "Time", Width: 16},
		{Title: "Browser", Width: 8},
	}
	m.table.SetColumns(columns)
	m.updateTableRows()
}

func (m *Model) urlWidth() int {
	w := m.width - 3 - 30 - 16 - 8 - 10
	if w < 20 {
		w = 20
	}
	return w
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

func Run(browsers []browser.Browser) error {
	m := NewModel(browsers)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
