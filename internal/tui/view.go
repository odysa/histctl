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
		loadingText := lipgloss.NewStyle().Foreground(Accent).Italic(true).Render("Loading history...")
		sections = append(sections, fmt.Sprintf("\n  %s %s\n", m.spinner.View(), loadingText))
	} else if len(m.filteredEntries) == 0 {
		emptyMsg := "No history entries"
		if m.searchText != "" {
			emptyMsg = fmt.Sprintf("No results for \"%s\"", m.searchText)
		}
		sections = append(sections, lipgloss.NewStyle().
			Foreground(Subtle).
			Italic(true).
			Width(m.width).
			Align(lipgloss.Center).
			Padding(3, 0).
			Render(emptyMsg))
	} else {
		sections = append(sections, m.table.View())
	}

	sections = append(sections, m.renderStatusBar())

	if m.showHelp {
		sections = append(sections, m.help.View(m.keys))
	} else {
		sections = append(sections, m.renderShortHelp())
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	if m.state == stateConfirmDelete {
		content = m.overlayDialog(content)
	}

	return content
}

func (m Model) renderTitle() string {
	return TitleStyle.Render(" histctl ")
}

func (m Model) renderHeader() string {
	title := m.renderTitle()

	pills := make([]string, len(m.browserNames))
	for i := range m.browserNames {
		pills[i] = m.renderPill(i)
	}

	right := lipgloss.JoinHorizontal(lipgloss.Top, pills...)
	gap := m.width - lipgloss.Width(title) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		strings.Repeat(" ", gap),
		right,
	)

	divider := lipgloss.NewStyle().Foreground(Muted).Render(strings.Repeat("─", m.width))
	return header + "\n" + divider
}

func (m Model) renderPill(index int) string {
	name := m.browserNames[index]
	color := BrowserColor(name)
	if name == "all" {
		color = Accent
	}
	count := m.browserCount(name)
	label := fmt.Sprintf("%s %d", name, count)
	if index == m.activeBrowser {
		return BrowserPillActive(color).Render(label)
	}
	return BrowserPillInactive(color).Render(label)
}

func (m Model) browserPillAt(x, y int) int {
	if y != 0 {
		return -1
	}

	titleW := lipgloss.Width(m.renderTitle())
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
	dot := lipgloss.NewStyle().Foreground(Muted).Render(" · ")
	var parts []string

	if len(m.filteredEntries) > 0 {
		cursor := m.table.Cursor() + 1
		parts = append(parts, lipgloss.NewStyle().Foreground(Accent).Bold(true).Render(
			fmt.Sprintf("%d/%d", cursor, len(m.filteredEntries))))
	}

	entryText := fmt.Sprintf("%d entries", len(m.filteredEntries))
	if len(m.allEntries) != len(m.filteredEntries) {
		entryText += fmt.Sprintf(" of %d", len(m.allEntries))
	}
	parts = append(parts, lipgloss.NewStyle().Foreground(Subtle).Render(entryText))

	if len(m.selected) > 0 {
		parts = append(parts, SelectedCountStyle.Render(
			fmt.Sprintf("%d selected", len(m.selected))))
	}

	if m.statusMsg != "" {
		parts = append(parts, m.statusMsg)
	}

	if m.err != nil {
		parts = append(parts, ErrorStyle.Render(m.err.Error()))
	}

	divider := lipgloss.NewStyle().Foreground(Muted).Render(strings.Repeat("─", m.width))
	return divider + "\n" + StatusBarStyle.Render("  "+strings.Join(parts, dot))
}

func (m Model) overlayDialog(bg string) string {
	title := lipgloss.NewStyle().Bold(true).Foreground(Danger).
		Render(fmt.Sprintf("Delete %d entries?", len(m.selected)))

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3D3D5C")).
		Padding(0, 1)
	descStyle := lipgloss.NewStyle().Foreground(Subtle)
	sep := lipgloss.NewStyle().Foreground(Muted).Render("    ")
	hint := keyStyle.Render("y") + " " + descStyle.Render("confirm") +
		sep +
		keyStyle.Render("esc") + " " + descStyle.Render("cancel")

	dialog := DialogStyle.Render(lipgloss.JoinVertical(lipgloss.Center, title, "", hint))

	dimStyle := lipgloss.NewStyle().Foreground(Muted)
	bgLines := strings.Split(bg, "\n")
	for i, line := range bgLines {
		bgLines[i] = dimStyle.Render(line)
	}

	for len(bgLines) < m.height {
		bgLines = append(bgLines, "")
	}

	dialogLines := strings.Split(dialog, "\n")
	dH := len(dialogLines)
	dW := lipgloss.Width(dialog)
	startRow := (m.height - dH) / 2
	startCol := (m.width - dW) / 2
	if startCol < 0 {
		startCol = 0
	}

	for i, dLine := range dialogLines {
		row := startRow + i
		if row < 0 || row >= len(bgLines) {
			continue
		}
		bgLine := bgLines[row]
		bgW := lipgloss.Width(bgLine)
		if bgW < startCol {
			bgLine += strings.Repeat(" ", startCol-bgW)
		}
		left := ansiTruncate(bgLine, startCol)
		bgLines[row] = left + dLine
	}

	return strings.Join(bgLines[:m.height], "\n")
}

func ansiTruncate(s string, n int) string {
	var result strings.Builder
	visible := 0
	inEsc := false
	for _, r := range s {
		if inEsc {
			result.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			result.WriteRune(r)
			continue
		}
		if visible >= n {
			break
		}
		result.WriteRune(r)
		visible++
	}
	for visible < n {
		result.WriteByte(' ')
		visible++
	}
	return result.String()
}

func (m Model) renderShortHelp() string {
	bindings := m.keys.ShortHelp()
	bracket := lipgloss.NewStyle().Foreground(Muted)
	var parts []string
	for _, b := range bindings {
		k := bracket.Render("[") + HelpKeyStyle.Render(b.Help().Key) + bracket.Render("]")
		d := HelpDescStyle.Render(b.Help().Desc)
		parts = append(parts, k+" "+d)
	}
	sep := HelpSepStyle.Render("  ")
	return lipgloss.NewStyle().MarginLeft(2).Render(
		strings.Join(parts, sep))
}

