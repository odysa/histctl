package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.resizeTable()
		return m, nil

	case historyLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateViewing
			return m, nil
		}
		m.allEntries = msg.entries
		m.searchEntries = nil
		if m.searchText != "" {
			return m, m.searchHistory(m.searchText)
		}
		m.state = stateViewing
		m.applyFilters()
		return m, nil

	case searchResultsMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateViewing
			return m, nil
		}
		m.searchEntries = msg.entries
		m.state = stateViewing
		m.applyFilters()
		return m, nil

	case deleteResultMsg:
		if msg.err != nil {
			m.statusMsg = ErrorStyle.Render(fmt.Sprintf("Delete failed: %v", msg.err))
		} else {
			m.statusMsg = lipgloss.NewStyle().Foreground(Success).Render(
				fmt.Sprintf("Deleted %d entries", msg.result.Deleted))
		}
		m.state = stateLoading
		return m, m.loadHistory()

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			if idx := m.browserPillAt(msg.X, msg.Y); idx >= 0 && idx != m.activeBrowser {
				m.activeBrowser = idx
				m.selected = make(map[int]bool)
				m.applyFilters()
				return m, nil
			}
			if m.state == stateViewing && m.isSearchBarAt(msg.Y) {
				return m.enterSearchMode()
			}
		}
		if m.state == stateViewing {
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateSearching:
			return m.updateSearching(msg)
		case stateConfirmDelete:
			return m.updateConfirmDelete(msg)
		case stateViewing:
			return m.updateViewing(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateViewing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Search):
		return m.enterSearchMode()

	case key.Matches(msg, m.keys.Select):
		cursor := m.table.Cursor()
		if cursor < len(m.filteredEntries) {
			if m.selected[cursor] {
				delete(m.selected, cursor)
			} else {
				m.selected[cursor] = true
			}
			m.updateTableRows()
			if cursor < len(m.filteredEntries)-1 {
				m.table.SetCursor(cursor + 1)
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.All):
		m.toggleSelectAll()
		m.updateTableRows()
		return m, nil

	case key.Matches(msg, m.keys.Delete):
		if len(m.selected) > 0 {
			m.state = stateConfirmDelete
		}
		return m, nil

	case key.Matches(msg, m.keys.Tab):
		m.activeBrowser = (m.activeBrowser + 1) % len(m.browserNames)
		m.selected = make(map[int]bool)
		m.applyFilters()
		return m, nil

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		return m, nil

	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
}

func (m Model) updateSearching(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.state = stateViewing
		m.searchInput.Blur()
		return m, nil

	case key.Matches(msg, m.keys.Apply):
		m.searchText = m.searchInput.Value()
		m.searchInput.Blur()
		m.selected = make(map[int]bool)
		if m.searchText == "" {
			m.searchEntries = nil
			m.state = stateViewing
			m.applyFilters()
			return m, nil
		}
		m.state = stateLoading
		return m, m.searchHistory(m.searchText)

	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
}

func (m Model) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.state = stateLoading
		return m, m.performDelete()
	default:
		m.state = stateViewing
		m.statusMsg = lipgloss.NewStyle().Foreground(Subtle).Render("Delete cancelled")
		return m, nil
	}
}

// Helpers

func (m Model) enterSearchMode() (tea.Model, tea.Cmd) {
	m.state = stateSearching
	m.searchInput.SetValue(m.searchText)
	m.searchInput.Focus()
	return m, m.searchInput.Cursor.BlinkCmd()
}

func (m *Model) selectAll() {
	for i := range m.filteredEntries {
		m.selected[i] = true
	}
}

func (m *Model) toggleSelectAll() {
	if len(m.selected) == len(m.filteredEntries) {
		m.selected = make(map[int]bool)
	} else {
		m.selectAll()
	}
}
