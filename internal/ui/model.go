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
	FocusSearch
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
	allScriptItems    []list.Item // Unfiltered backup for search
	optionCategories  []list.Item
	colorSchemeItems  []list.Item

	// Search
	searchInput  *components.SearchInput
	searchActive bool
	recursiveMode bool // Toggle for recursive vs directory view

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

	// Create search input
	searchInput := components.NewSearchInput(theme)

	// Set initial preview if there are scripts
	if len(scriptItems) > 0 {
		if s, ok := scriptItems[0].(scripts.Item); ok && s.IsScript() {
			content := scripts.ReadContentWithHighlighting(s.Description(), cache)
			vp.SetContent(content)
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
		allScriptItems:    scriptItems, // Keep backup for search
		optionCategories:  optionCategories,
		colorSchemeItems:  colorSchemeItems,
		searchInput:       searchInput,
		searchActive:      false,
		recursiveMode:     false, // Start in directory mode
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
	// Handle different terminal sizes more conservatively
	var leftPanelWidth, rightPanelWidth int
	
	if m.width < 50 {
		// Very small terminals - use almost full width
		leftPanelWidth = m.width - 2
		rightPanelWidth = 0
	} else if m.width < 80 {
		// Small terminals - single panel
		leftPanelWidth = m.width - 4
		rightPanelWidth = 0
	} else if m.width < 120 {
		// Medium terminals - adjust ratios
		leftPanelWidth = (m.width * 2) / 5 - 1
		rightPanelWidth = (m.width * 3) / 5 - 1
	} else {
		// Large terminals - standard layout
		leftPanelWidth = (m.width / 3) - 2
		rightPanelWidth = ((m.width * 2) / 3) - 2
	}

	// Content area accounting for borders and padding
	var leftContentWidth, rightContentWidth int
	
	if m.width < 40 {
		// Very small - no padding/borders
		leftContentWidth = leftPanelWidth
		rightContentWidth = rightPanelWidth
	} else if m.width < 60 {
		// Small - minimal padding
		leftContentWidth = leftPanelWidth - 4
		rightContentWidth = rightPanelWidth - 4
	} else {
		// Standard - full padding
		leftContentWidth = leftPanelWidth - 8
		rightContentWidth = rightPanelWidth - 8
	}

	// Ensure positive values
	if leftContentWidth < 5 {
		leftContentWidth = 5
	}
	if rightContentWidth < 5 {
		rightContentWidth = 5
	}

	// Adjust list heights based on available space
	listHeight := m.height - 10
	if m.height < 20 {
		listHeight = m.height - 6 // Less space for footer/header in small terminals
	}

	m.list.SetSize(leftContentWidth, listHeight)
	
	// Only set right panel sizes if we have space for it
	if rightPanelWidth > 0 {
		m.optionsRightList.SetSize(rightContentWidth, listHeight)
		m.vp.Width = rightContentWidth
		m.vp.Height = listHeight
	} else {
		// In single-panel mode, make viewport smaller for overlay display
		m.vp.Width = leftContentWidth
		m.vp.Height = listHeight / 2
	}
}

// switchTab switches to the specified tab.
func (m *Model) switchTab(tabIndex int) {
	m.activeTab = tabIndex
	m.focus = FocusList // Reset focus when switching tabs

	switch tabIndex {
	case 0: // Scripts tab
		m.list.SetItems(m.scriptItems)
		if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsScript() {
			content := scripts.ReadContentWithHighlighting(sel.Description(), m.cache)
			m.vp.SetContent(content)
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
	if !item.IsDirectory() || m.recursiveMode {
		return // Don't navigate directories in recursive mode
	}

	m.parentPaths = append(m.parentPaths, ParentNav{Path: m.currentPath, Index: m.list.Index()})
	m.currentPath = item.Description()
	newItems := scripts.GetItems(item.Description())
	m.scriptItems = newItems
	m.allScriptItems = newItems // Update backup as well
	m.list.SetItems(m.scriptItems)
	m.list.ResetSelected()

	// Clear search when navigating
	m.searchInput.Reset()
	m.searchActive = false

	// Show preview for first item if it's a script
	if len(m.scriptItems) > 0 {
		if first, ok := m.scriptItems[0].(scripts.Item); ok && first.IsScript() {
			content := scripts.ReadContentWithHighlighting(first.Description(), m.cache)
			m.vp.SetContent(content)
		} else {
			m.vp.SetContent("Select a script to preview...")
		}
	} else {
		m.vp.SetContent("Select a script to preview...")
	}
}

// navigateToParent navigates back to the parent directory.
func (m *Model) navigateToParent() {
	if len(m.parentPaths) == 0 || m.recursiveMode {
		return // Don't navigate in recursive mode
	}

	parent := m.parentPaths[len(m.parentPaths)-1]
	m.parentPaths = m.parentPaths[:len(m.parentPaths)-1]
	newItems := scripts.GetItems(parent.Path)
	m.scriptItems = newItems
	m.allScriptItems = newItems // Update backup as well
	m.list.SetItems(m.scriptItems)
	m.list.Select(parent.Index) // Restore previous selection
	m.currentPath = parent.Path

	// Clear search when navigating
	m.searchInput.Reset()
	m.searchActive = false

	// Show preview for selected item
	if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsScript() {
		content := scripts.ReadContentWithHighlighting(sel.Description(), m.cache)
		m.vp.SetContent(content)
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
		content := scripts.ReadContentWithHighlighting(sel.Description(), m.cache)
		m.vp.SetContent(content)
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
