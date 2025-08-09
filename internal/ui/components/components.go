// Package components provides reusable UI components for go-pwr.
package components

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/rocketpowerinc/go-pwr/internal/scripts"
	"github.com/rocketpowerinc/go-pwr/internal/ui/styles"
)

// ScriptDelegate handles rendering of script items in lists.
type ScriptDelegate struct {
	theme *styles.Theme
}

// NewScriptDelegate creates a new script delegate.
func NewScriptDelegate(theme *styles.Theme) *ScriptDelegate {
	return &ScriptDelegate{theme: theme}
}

// Render renders a script item.
func (d *ScriptDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s, ok := item.(scripts.Item)
	if !ok {
		return
	}

	var style lipgloss.Style
	if s.IsDirectory() {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Secondary)
	} else {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Primary)
	}

	if index == m.Index() {
		style = style.Bold(true).Underline(true)
	}

	fmt.Fprint(w, style.Render(s.Title()))
}

func (d *ScriptDelegate) Height() int                             { return 1 }
func (d *ScriptDelegate) Spacing() int                            { return 0 }
func (d *ScriptDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// OptionItem represents an option in the options menu.
type OptionItem struct {
	Name        string
	Desc        string
	Action      string
}

func (o OptionItem) Title() string       { return o.Name }
func (o OptionItem) Description() string { return o.Desc }
func (o OptionItem) FilterValue() string { return o.Name }

// OptionDelegate handles rendering of option items.
type OptionDelegate struct {
	theme *styles.Theme
}

// NewOptionDelegate creates a new option delegate.
func NewOptionDelegate(theme *styles.Theme) *OptionDelegate {
	return &OptionDelegate{theme: theme}
}

// Render renders an option item.
func (d *OptionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	o, ok := item.(OptionItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	if index == m.Index() {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Accent).Bold(true).Underline(true)
	} else {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Primary)
	}

	fmt.Fprint(w, style.Render(o.Title()))
}

func (d *OptionDelegate) Height() int                             { return 1 }
func (d *OptionDelegate) Spacing() int                            { return 0 }
func (d *OptionDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// CategoryItem represents a category in the options menu.
type CategoryItem struct {
	Name        string
	Desc        string
	Category    string
}

func (oc CategoryItem) Title() string       { return oc.Name }
func (oc CategoryItem) Description() string { return oc.Desc }
func (oc CategoryItem) FilterValue() string { return oc.Name }

// CategoryDelegate handles rendering of category items.
type CategoryDelegate struct {
	theme *styles.Theme
}

// NewCategoryDelegate creates a new category delegate.
func NewCategoryDelegate(theme *styles.Theme) *CategoryDelegate {
	return &CategoryDelegate{theme: theme}
}

// Render renders a category item.
func (d *CategoryDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	oc, ok := item.(CategoryItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	if index == m.Index() {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Accent).Bold(true).Underline(true)
	} else {
		style = lipgloss.NewStyle().Foreground(d.theme.Current.Primary)
	}

	fmt.Fprint(w, style.Render(oc.Title()))
}

func (d *CategoryDelegate) Height() int                             { return 1 }
func (d *CategoryDelegate) Spacing() int                            { return 0 }
func (d *CategoryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// CreateList creates a configured list widget.
func CreateList(items []list.Item, delegate list.ItemDelegate) list.Model {
	listModel := list.New(items, delegate, 0, 0)
	listModel.Title = ""
	listModel.SetShowHelp(false)
	listModel.SetFilteringEnabled(false)
	listModel.SetShowStatusBar(false)
	listModel.SetShowPagination(false)
	listModel.DisableQuitKeybindings()
	listModel.Styles.Title = lipgloss.NewStyle()
	listModel.Styles.PaginationStyle = lipgloss.NewStyle()
	listModel.Styles.HelpStyle = lipgloss.NewStyle()
	return listModel
}
