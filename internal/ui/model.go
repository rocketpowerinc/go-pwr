// Package ui provides the main user interface for go-pwr.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/rocketpowerinc/go-pwr/internal/config"
	"github.com/rocketpowerinc/go-pwr/internal/scripts"
	"github.com/rocketpowerinc/go-pwr/internal/ui/components"
	"github.com/rocketpowerinc/go-pwr/internal/ui/styles"
	"github.com/rocketpowerinc/go-pwr/pkg/platform"
)

// FocusArea represents which UI area has focus.
type FocusArea int

const (
	FocusList FocusArea = iota
	FocusPreview
)

// ParentNav tracks navigation state for going back to parent directories.
type ParentNav struct {
	Path  string
	Index int
}

// Model represents the main UI model.
type Model struct {
	config           *config.Config
	theme            *styles.Theme
	width            int
	height           int
	activeTab        int
	tabs             []string
	focus            FocusArea
	currentPath      string
	parentPaths      []ParentNav
	cache            *scripts.Cache
	selectedCategory string

	// Lists and viewport
	list             list.Model  // Main left panel list (scripts or categories)
	optionsRightList list.Model  // Right panel list for options
	categoryList     list.Model  // Category list for options tab
	vp               viewport.Model

	// Items
	scriptItems       []list.Item
	optionCategories  []list.Item
	colorSchemeItems  []list.Item

	// Delegates
	scriptDelegate   *components.ScriptDelegate
	optionDelegate   *components.OptionDelegate
	categoryDelegate *components.CategoryDelegate
}

// NewModel creates a new UI model.
func NewModel(cfg *config.Config) *Model {
	theme := styles.NewTheme(styles.OceanBreeze)
	cache := scripts.NewCache()

	// Create delegates
	scriptDelegate := components.NewScriptDelegate(theme)
	optionDelegate := components.NewOptionDelegate(theme)
	categoryDelegate := components.NewCategoryDelegate(theme)

	// Get initial script items
	scriptItems := scripts.GetItems(cfg.ScriptbinPath)

	// Create option categories
	optionCategories := []list.Item{
		components.CategoryItem{
			Name:     "Color Schemes",
			Desc:     "Change the application's color theme",
			Category: "color_schemes",
		},
	}

	// Create color scheme items
	schemes := styles.AllSchemes()
	colorSchemeItems := make([]list.Item, len(schemes))
	for i, scheme := range schemes {
		colorSchemeItems[i] = components.OptionItem{
			Name:   scheme.Name,
			Desc:   "Apply this color scheme",
			Action: "color_scheme",
		}
	}

	// Create lists
	scriptList := components.CreateList(scriptItems, scriptDelegate)
	categoryList := components.CreateList(optionCategories, categoryDelegate)
	optionsRightList := components.CreateList(colorSchemeItems, optionDelegate)

	// Create viewport
	vp := viewport.New(0, 0)
	vp.SetContent("Select a script to preview...")

	// Set initial preview if there are scripts
	if len(scriptItems) > 0 {
		if s, ok := scriptItems[0].(scripts.Item); ok && s.IsScript() {
			content := scripts.ReadContent(s.Description(), cache)
			vp.SetContent(scripts.SanitizeContent(content))
		}
	}

	return &Model{
		config:            cfg,
		theme:             theme,
		tabs:              []string{"Scripts", "Options", "About"},
		activeTab:         0,
		focus:             FocusList,
		currentPath:       cfg.ScriptbinPath,
		parentPaths:       []ParentNav{},
		cache:             cache,
		selectedCategory:  "",
		list:              scriptList,
		optionsRightList:  optionsRightList,
		categoryList:      categoryList,
		vp:                vp,
		scriptItems:       scriptItems,
		optionCategories:  optionCategories,
		colorSchemeItems:  colorSchemeItems,
		scriptDelegate:    scriptDelegate,
		optionDelegate:    optionDelegate,
		categoryDelegate:  categoryDelegate,
	}
}

// Start starts the UI.
func Start(cfg *config.Config) error {
	model := NewModel(cfg)
	
	program := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)
	
	return program.Start()
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// setSizes sets the sizes of UI components based on window dimensions.
func (m *Model) setSizes() {
	// Account for borders in width calculations to prevent overflow
	leftPanelWidth := (m.width / 3) - 2   // -2 for left panel border
	rightPanelWidth := ((m.width * 2) / 3) - 2 // -2 for right panel border

	// Content area accounting for borders and padding
	leftContentWidth := leftPanelWidth - 8
	rightContentWidth := rightPanelWidth - 8

	// Ensure positive values
	if leftContentWidth < 5 {
		leftContentWidth = 5
	}
	if rightContentWidth < 5 {
		rightContentWidth = 5
	}

	m.list.SetSize(leftContentWidth, m.height-10)
	m.optionsRightList.SetSize(rightContentWidth, m.height-10)
	m.vp.Width = rightContentWidth
	m.vp.Height = m.height - 10
}

// switchTab switches to the specified tab.
func (m *Model) switchTab(tabIndex int) {
	m.activeTab = tabIndex
	m.focus = FocusList // Reset focus when switching tabs

	switch tabIndex {
	case 0: // Scripts tab
		m.list.SetItems(m.scriptItems)
		if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsScript() {
			content := scripts.ReadContent(sel.Description(), m.cache)
			m.vp.SetContent(scripts.SanitizeContent(content))
		} else {
			m.vp.SetContent("Select a script to preview...")
		}
	case 1: // Options tab
		// Copy items from categoryList to main list for display
		items := make([]list.Item, len(m.optionCategories))
		copy(items, m.optionCategories)
		m.list.SetItems(items)
		m.selectedCategory = ""
		m.vp.SetContent("Select an option category from the left to see available settings.")
	case 2: // About tab
		m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
	}
}

// navigateIntoDirectory navigates into a selected directory.
func (m *Model) navigateIntoDirectory(item scripts.Item) {
	if !item.IsDirectory() {
		return
	}

	m.parentPaths = append(m.parentPaths, ParentNav{Path: m.currentPath, Index: m.list.Index()})
	m.currentPath = item.Description()
	m.scriptItems = scripts.GetItems(item.Description())
	m.list.SetItems(m.scriptItems)
	m.list.ResetSelected()

	// Show preview for first item if it's a script
	if len(m.scriptItems) > 0 {
		if first, ok := m.scriptItems[0].(scripts.Item); ok && first.IsScript() {
			content := scripts.ReadContent(first.Description(), m.cache)
			m.vp.SetContent(scripts.SanitizeContent(content))
		} else {
			m.vp.SetContent("Select a script to preview...")
		}
	} else {
		m.vp.SetContent("Select a script to preview...")
	}
}

// navigateToParent navigates back to the parent directory.
func (m *Model) navigateToParent() {
	if len(m.parentPaths) == 0 {
		return
	}

	parent := m.parentPaths[len(m.parentPaths)-1]
	m.parentPaths = m.parentPaths[:len(m.parentPaths)-1]
	m.scriptItems = scripts.GetItems(parent.Path)
	m.list.SetItems(m.scriptItems)
	m.list.Select(parent.Index) // Restore previous selection
	m.currentPath = parent.Path

	// Show preview for selected item
	if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsScript() {
		content := scripts.ReadContent(sel.Description(), m.cache)
		m.vp.SetContent(scripts.SanitizeContent(content))
	} else {
		m.vp.SetContent("Select a script to preview...")
	}
}

// executeScript runs the selected script.
func (m *Model) executeScript(item scripts.Item) {
	if !item.IsScript() {
		return
	}

	m.vp.SetContent("Running script in a new terminal window...")
	go func() {
		if err := platform.ExecuteScript(item.Description(), item.Title()); err != nil {
			// Could add error handling here, maybe show error in UI
		}
	}()
}

// updatePreview updates the preview pane with the content of the selected item.
func (m *Model) updatePreview() {
	if m.activeTab != 0 {
		return
	}

	if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsScript() {
		content := scripts.ReadContent(sel.Description(), m.cache)
		m.vp.SetContent(scripts.SanitizeContent(content))
	} else {
		m.vp.SetContent("Select a script to preview...")
	}
}

// applyColorScheme applies a color scheme.
func (m *Model) applyColorScheme(schemeName string) {
	schemes := styles.AllSchemes()
	for _, scheme := range schemes {
		if strings.Contains(schemeName, scheme.Name) {
			m.theme.UpdateScheme(scheme)
			m.vp.SetContent("Color scheme changed to " + scheme.Name + "! Press Tab to switch tabs and see the changes.")
			break
		}
	}
}
