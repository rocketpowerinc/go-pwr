package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rocketpowerinc/go-pwr/internal/config"
	"github.com/rocketpowerinc/go-pwr/internal/ui/styles"
)

// RepositoryInput represents a repository URL input component
type RepositoryInput struct {
	textInput textinput.Model
	theme     *styles.Theme
	active    bool
	errorMsg  string
}

// NewRepositoryInput creates a new repository input component
func NewRepositoryInput(theme *styles.Theme) *RepositoryInput {
	ti := textinput.New()
	ti.Placeholder = "https://github.com/username/repo.git"
	ti.CharLimit = 200
	ti.Width = 50 // Default width, will be adjusted dynamically

	return &RepositoryInput{
		textInput: ti,
		theme:     theme,
		active:    false,
		errorMsg:  "",
	}
}

// SetWidth adjusts the width of the repository input
func (ri *RepositoryInput) SetWidth(width int) {
	// Ensure minimum width for usability
	if width < 20 {
		width = 20
	}
	// Maximum width to prevent excessive length
	if width > 100 {
		width = 100
	}
	ri.textInput.Width = width
}

// SetActive sets whether the repository input is active
func (ri *RepositoryInput) SetActive(active bool) {
	ri.active = active
	if active {
		ri.textInput.Focus()
		ri.errorMsg = "" // Clear error when activating
	} else {
		ri.textInput.Blur()
		ri.textInput.SetCursor(0)
	}
}

// Update updates the repository input
func (ri *RepositoryInput) Update(msg tea.Msg) tea.Cmd {
	// Don't process escape key - let parent handle it
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "escape" {
		return nil
	}

	var cmd tea.Cmd
	ri.textInput, cmd = ri.textInput.Update(msg)

	// Clear error message when user types
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyRunes {
		ri.errorMsg = ""
	}

	return cmd
}

// View renders the repository input
func (ri *RepositoryInput) View() string {
	inputView := ri.textInput.View()

	// Add error message if present
	if ri.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		return inputView + "\n" + errorStyle.Render("❌ " + ri.errorMsg)
	}

	return inputView
}

// Value returns the current input value
func (ri *RepositoryInput) Value() string {
	return ri.textInput.Value()
}

// SetValue sets the input value
func (ri *RepositoryInput) SetValue(value string) {
	ri.textInput.SetValue(value)
}

// Reset clears the input
func (ri *RepositoryInput) Reset() {
	ri.textInput.SetValue("")
	ri.errorMsg = ""
}

// Validate validates the current repository URL and sets error message if invalid
func (ri *RepositoryInput) Validate() bool {
	url := ri.textInput.Value()
	if url == "" {
		ri.errorMsg = "Repository URL cannot be empty"
		return false
	}

	if err := config.ValidateRepoURL(url); err != nil {
		ri.errorMsg = err.Error()
		return false
	}

	ri.errorMsg = ""
	return true
}

// SetError sets a custom error message
func (ri *RepositoryInput) SetError(errorMsg string) {
	ri.errorMsg = errorMsg
}

// ClearError clears the error message
func (ri *RepositoryInput) ClearError() {
	ri.errorMsg = ""
}
