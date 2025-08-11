package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rocketpowerinc/go-pwr/internal/ui/styles"
)

// SearchInput represents a search input component
type SearchInput struct {
	textInput textinput.Model
	theme     *styles.Theme
	active    bool
}

// NewSearchInput creates a new search input component
func NewSearchInput(theme *styles.Theme) *SearchInput {
	ti := textinput.New()
	ti.Placeholder = "Search by tags (e.g., bash linux ubuntu)..."
	ti.CharLimit = 100
	ti.Width = 50 // Default width, will be adjusted dynamically
	
	return &SearchInput{
		textInput: ti,
		theme:     theme,
		active:    false,
	}
}

// SetWidth adjusts the width of the search input
func (si *SearchInput) SetWidth(width int) {
	// Ensure minimum width for usability
	if width < 8 {
		width = 8
	}
	// Maximum width to prevent excessive length
	if width > 80 {
		width = 80
	}
	si.textInput.Width = width
}

// ViewMinimal renders a minimal search input for very small spaces
func (si *SearchInput) ViewMinimal() string {
	value := si.textInput.Value()
	if value == "" {
		if si.active {
			return "🔍 Search..."
		}
		return "🔍 Search"
	}
	
	// Truncate if too long
	if len(value) > 15 {
		value = value[:12] + "..."
	}
	
	if si.active {
		return "🔍 " + value + "_"
	}
	return "🔍 " + value
}

// SetActive sets whether the search input is active
func (si *SearchInput) SetActive(active bool) {
	si.active = active
	if active {
		si.textInput.Focus()
	} else {
		si.textInput.Blur()
		// Also reset cursor to beginning when deactivating
		si.textInput.SetCursor(0)
	}
}

// Update updates the search input
func (si *SearchInput) Update(msg tea.Msg) tea.Cmd {
	// Don't process escape key - let parent handle it
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "escape" {
		return nil
	}
	
	var cmd tea.Cmd
	si.textInput, cmd = si.textInput.Update(msg)
	return cmd
}

// View renders the search input
func (si *SearchInput) View() string {
	// Simple approach: just show the prefix and the actual text value, no fancy styling
	prefix := "🔍 "
	
	if si.active {
		// When active, show the value and cursor
		value := si.textInput.Value()
		if value == "" {
			return prefix + "Search by tags..."
		}
		return prefix + value + "_"
	} else {
		// When inactive, just show the value or placeholder
		value := si.textInput.Value()
		if value == "" {
			return prefix + "Search by tags"
		}
		return prefix + value
	}
}

// Value returns the current input value
func (si *SearchInput) Value() string {
	return si.textInput.Value()
}

// SetValue sets the input value
func (si *SearchInput) SetValue(value string) {
	si.textInput.SetValue(value)
}

// Reset clears the input
func (si *SearchInput) Reset() {
	si.textInput.SetValue("")
}
