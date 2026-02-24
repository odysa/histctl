package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Select    key.Binding
	All       key.Binding
	Delete    key.Binding
	DeleteAll key.Binding
	Search    key.Binding
	Tab       key.Binding
	Apply     key.Binding
	Cancel    key.Binding
	Help      key.Binding
	Quit      key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Select:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "select")),
		All:       key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all")),
		Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		DeleteAll: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete all")),
		Search:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Tab:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "browser")),
		Apply:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Select, k.Delete, k.DeleteAll, k.Tab, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.All},
		{k.Search, k.Delete, k.DeleteAll, k.Tab},
		{k.Help, k.Quit},
	}
}
