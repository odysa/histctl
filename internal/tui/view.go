package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sections []string

	sections = append(sections, m.renderHeader())
	sections = append(sections, m.renderSearchBar())

	if m.state == stateLoading {
		sections = append(sections, fmt.Sprintf("\n  %s Loading history...\n", m.spinner.View()))
	} else if m.state == stateConfirmDelete {
		sections = append(sections, m.table.View())
		confirm := lipgloss.NewStyle().
			Bold(true).
			Foreground(Danger).
			MarginTop(1).
			Render(fmt.Sprintf("  Delete %d entries? (y/N)", len(m.selected)))
		sections = append(sections, confirm)
	} else {
		sections = append(sections, m.table.View())
	}

	sections = append(sections, m.renderStatusBar())

	if m.showHelp {
		sections = append(sections, m.help.View(m.keys))
	} else {
		sections = append(sections, m.renderShortHelp())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	title := TitleStyle.Render(" histctl ")

	pills := make([]string, len(m.browserNames))
	for i := range m.browserNames {
		pills[i] = m.renderPill(i)
	}

	right := lipgloss.JoinHorizontal(lipgloss.Top, pills...)
	gap := m.width - lipgloss.Width(title) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		strings.Repeat(" ", gap),
		right,
	)
}

func (m Model) renderPill(index int) string {
	name := m.browserNames[index]
	color := BrowserColor(name)
	if name == "all" {
		color = Accent
	}
	if index == m.activeBrowser {
		return BrowserPillActive(color).Render(name)
	}
	return BrowserPillInactive(color).Render(name)
}

// browserPillAt returns the browser index at the given screen position, or -1.
func (m Model) browserPillAt(x, y int) int {
	if y != 0 {
		return -1
	}

	titleW := lipgloss.Width(TitleStyle.Render(" histctl "))
	var totalPillW int
	pillWidths := make([]int, len(m.browserNames))
	for i := range m.browserNames {
		w := lipgloss.Width(m.renderPill(i))
		pillWidths[i] = w
		totalPillW += w
	}

	gap := m.width - titleW - totalPillW - 2
	if gap < 1 {
		gap = 1
	}

	pos := titleW + gap
	for i, w := range pillWidths {
		if x >= pos && x < pos+w {
			return i
		}
		pos += w
	}
	return -1
}

func (m Model) isSearchBarAt(y int) bool {
	headerH := lipgloss.Height(m.renderHeader())
	searchH := lipgloss.Height(m.renderSearchBar())
	return y >= headerH && y < headerH+searchH
}

func (m Model) renderSearchBar() string {
	style := SearchBarStyle
	if m.state == stateSearching {
		style = SearchBarActiveStyle
	}

	var content string
	if m.state == stateSearching {
		content = m.searchInput.View()
	} else if m.searchText != "" {
		content = SearchLabelStyle.Render("/ ") +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Render(m.searchText)
	} else {
		content = lipgloss.NewStyle().Foreground(Subtle).Render("/ search with regex...")
	}

	w := m.width - 4
	if w < 10 {
		w = 10
	}
	return style.Width(w).Render(content)
}

func (m Model) renderStatusBar() string {
	var parts []string

	parts = append(parts, lipgloss.NewStyle().Foreground(Subtle).Render(
		fmt.Sprintf("%d entries", len(m.filteredEntries))))

	if len(m.allEntries) != len(m.filteredEntries) {
		parts = append(parts, lipgloss.NewStyle().Foreground(Subtle).Render(
			fmt.Sprintf("of %d", len(m.allEntries))))
	}

	if len(m.selected) > 0 {
		parts = append(parts, SelectedCountStyle.Render(
			fmt.Sprintf("│ %d selected", len(m.selected))))
	}

	if m.statusMsg != "" {
		parts = append(parts, "│", m.statusMsg)
	}

	if m.err != nil {
		parts = append(parts, "│", ErrorStyle.Render(m.err.Error()))
	}

	return StatusBarStyle.Render("  " + strings.Join(parts, " "))
}

func (m Model) renderShortHelp() string {
	bindings := m.keys.ShortHelp()
	var parts []string
	for _, b := range bindings {
		k := HelpKeyStyle.Render(b.Help().Key)
		d := HelpDescStyle.Render(b.Help().Desc)
		parts = append(parts, k+" "+d)
	}
	sep := HelpSepStyle.Render(" · ")
	return lipgloss.NewStyle().Foreground(Subtle).MarginLeft(2).Render(
		strings.Join(parts, sep))
}
