package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rocketpowerinc/go-pwr/internal/scripts"
	"github.com/rocketpowerinc/go-pwr/internal/ui/components"
)

// Update handles UI state updates.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.setSizes()
		return m, nil

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft && msg.Y == 0 {
			// Handle tab clicks
			x := 0
			for i, tab := range m.tabs {
				style := m.theme.TabInactive.Copy()
				if i == m.activeTab {
					style = style.Inherit(m.theme.TabActive)
				}
				renderedTab := style.Render(tab)
				tabWidth := lipgloss.Width(renderedTab)

				if msg.X >= x && msg.X < x+tabWidth {
					m.switchTab(i)
					break
				}
				x += tabWidth + 2
			}
		}

	case tea.KeyMsg:
		// Handle escape key with multiple detection methods
		if msg.String() == "escape" || msg.Type == tea.KeyEscape {
			if m.searchActive && m.activeTab == 0 {
				m.searchActive = false
				m.searchInput.SetActive(false)
				m.focus = FocusList
				// Force a view update to ensure UI changes
				return m, nil
			} else if m.activeTab == 0 {
				// Clear search and show all items
				m.searchInput.Reset()
				m.applySearch()
				return m, nil
			}
		}
		
		// Handle search input if search is active
		if m.searchActive && m.activeTab == 0 {
			switch msg.String() {
			case "enter":
				m.searchActive = false
				m.searchInput.SetActive(false)
				m.focus = FocusList
				return m, nil
			default:
				// Pass all other keys (except escape, handled above) to search input
				cmd = m.searchInput.Update(msg)
				// Live search as user types
				m.applySearch()
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+f", "/":
			if m.activeTab == 0 {
				m.searchActive = true
				m.searchInput.SetActive(true)
				m.focus = FocusSearch
				return m, nil
			}
		case "ctrl+r":
			if m.activeTab == 0 {
				// Toggle recursive mode
				m.recursiveMode = !m.recursiveMode
				m.refreshView()
				return m, nil
			}
		case "tab":
			m.switchTab((m.activeTab + 1) % len(m.tabs))
		case "shift+tab":
			m.switchTab((m.activeTab - 1 + len(m.tabs)) % len(m.tabs))
		case "ctrl+left", "cmd+left", "alt+left":
			if m.activeTab == 0 || m.activeTab == 1 {
				m.focus = FocusList
			}
		case "ctrl+right", "cmd+right", "alt+right":
			if m.activeTab == 0 || m.activeTab == 1 {
				m.focus = FocusPreview
			}
		case "left":
			if m.activeTab == 0 && m.focus == FocusList && len(m.parentPaths) > 0 && !m.recursiveMode {
				m.navigateToParent()
				return m, nil
			}
		case "right":
			if m.activeTab == 0 && m.focus == FocusList && !m.recursiveMode {
				if sel, ok := m.list.SelectedItem().(scripts.Item); ok && sel.IsDirectory() {
					m.navigateIntoDirectory(sel)
					return m, nil
				}
			}
		case "up":
			return m.handleUpDown(msg, true)
		case "down":
			return m.handleUpDown(msg, false)
		case "enter":
			return m.handleEnter()
		case "page_up":
			if m.focus == FocusPreview && m.activeTab == 0 {
				for i := 0; i < 10; i++ {
					m.vp.LineUp(1)
				}
			}
		case "page_down":
			if m.focus == FocusPreview && m.activeTab == 0 {
				for i := 0; i < 10; i++ {
					m.vp.LineDown(1)
				}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Update search input if active
	if m.searchActive {
		m.searchInput.Update(msg)
	}

	m.list, _ = m.list.Update(msg)
	m.optionsRightList, _ = m.optionsRightList.Update(msg)
	m.vp, _ = m.vp.Update(msg)
	return m, cmd
}

// refreshView refreshes the current view based on recursive mode
func (m *Model) refreshView() {
	if m.recursiveMode {
		// Get all scripts recursively
		allItems := scripts.GetAllScriptsRecursively(m.config.ScriptbinPath)
		m.allScriptItems = allItems
		m.scriptItems = allItems
	} else {
		// Get items from current directory only
		items := scripts.GetItems(m.currentPath)
		m.allScriptItems = items
		m.scriptItems = items
	}
	
	// Apply current search if any
	m.applySearch()
}

// applySearch filters the script items based on the search input
func (m *Model) applySearch() {
	searchTerm := strings.TrimSpace(m.searchInput.Value())
	
	if searchTerm == "" {
		// Show all items when search is empty
		m.scriptItems = m.allScriptItems
	} else {
		// Split search term into individual tags
		searchTags := strings.Fields(strings.ToLower(searchTerm))
		if m.recursiveMode {
			// In recursive mode, only show scripts (no directories)
			var scriptOnlyItems []list.Item
			for _, item := range m.allScriptItems {
				if scriptItem, ok := item.(scripts.Item); ok && scriptItem.IsScript() {
					scriptOnlyItems = append(scriptOnlyItems, item)
				}
			}
			m.scriptItems = scripts.FilterItemsByTags(scriptOnlyItems, searchTags)
		} else {
			m.scriptItems = scripts.FilterItemsByTags(m.allScriptItems, searchTags)
		}
	}
	
	m.list.SetItems(m.scriptItems)
	if len(m.scriptItems) > 0 {
		m.list.Select(0)
		m.updatePreview()
	} else {
		if m.recursiveMode {
			m.vp.SetContent("No scripts found matching your search criteria in any directory.")
		} else {
			m.vp.SetContent("No scripts found matching your search criteria.")
		}
	}
}

// handleUpDown handles up/down key navigation.
func (m Model) handleUpDown(msg tea.KeyMsg, isUp bool) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.focus == FocusList && m.activeTab == 0 {
		prevIndex := m.list.Index()
		m.list, cmd = m.list.Update(msg)
		newIndex := m.list.Index()
		if prevIndex != newIndex {
			m.updatePreview()
		}
	} else if m.focus == FocusList && m.activeTab == 1 {
		m.list, cmd = m.list.Update(msg)
	} else if m.focus == FocusPreview && m.activeTab == 0 {
		if isUp {
			m.vp.LineUp(1)
		} else {
			m.vp.LineDown(1)
		}
	} else if m.focus == FocusPreview && m.activeTab == 1 {
		m.optionsRightList, cmd = m.optionsRightList.Update(msg)
	}

	return m, cmd
}

// handleEnter handles enter key presses.
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.activeTab == 0 && m.focus == FocusList {
		// Scripts tab - run script or navigate directory
		if sel, ok := m.list.SelectedItem().(scripts.Item); ok {
			if sel.IsDirectory() && !m.recursiveMode {
				m.navigateIntoDirectory(sel)
			} else if sel.IsScript() {
				m.executeScript(sel)
			}
		}
	} else if m.activeTab == 1 && m.focus == FocusList {
		// Options tab - select category
		if sel, ok := m.list.SelectedItem().(components.CategoryItem); ok {
			m.selectedCategory = sel.Category
			if sel.Category == "color_schemes" {
				m.optionsRightList.SetItems(m.colorSchemeItems)
				m.vp.SetContent("Select a color scheme from the list on the left to apply it instantly!")
			}
		}
	} else if m.activeTab == 1 && m.focus == FocusPreview {
		// Options tab - apply selected option
		if m.selectedCategory == "color_schemes" {
			if sel, ok := m.optionsRightList.SelectedItem().(components.OptionItem); ok {
				m.applyColorScheme(sel.Name)
			}
		}
	}
	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	// Handle extremely small terminals
	if m.width < 30 || m.height < 10 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			Width(m.width).
			Height(m.height).
			Render("Terminal too small!\nMinimum: 30x10\nCurrent: " + 
				fmt.Sprintf("%dx%d", m.width, m.height))
	}

	// Render tab bar
	var tabLabels []string
	for i, name := range m.tabs {
		style := m.theme.TabInactive.Copy()
		if i == m.activeTab {
			style = style.Inherit(m.theme.TabActive)
		}
		tabLabels = append(tabLabels, style.Render(name))
	}
	tabBar := m.theme.TabBar.Render(strings.Join(tabLabels, "  "))

	// Render body based on active tab
	var body string
	switch m.activeTab {
	case 0:
		body = m.renderScriptsTab()
	case 1:
		body = m.renderOptionsTab()
	case 2:
		body = m.renderAboutTab()
	}

	// Render footer - responsive based on terminal width
	var footerText string
	if m.activeTab == 0 && m.searchActive {
		footerText = "'Enter' Apply • 'Esc' Cancel • Type to search..."
	} else if m.activeTab == 0 {
		if m.width < 80 {
			// Short footer for small terminals
			footerText = "'Tab' Tabs • '↑↓' Navigate • 'Enter' Run • 'Ctrl+F' Search • 'q' Quit"
		} else if m.width < 120 {
			// Medium footer for medium terminals
			footerText = "'Tab' Switch • '←↑↓→' Navigate • 'Enter' Run/Select • 'Ctrl+F' Search • 'Ctrl+R' Recursive • 'q' Quit"
		} else {
			// Full footer for large terminals
			footerText = "'Tab' Switch Tabs • '←↑↓→' Navigate • 'Ctrl + ←→' Switch Panes • 'Enter' Run/Select • 'Ctrl+F' Search • 'Ctrl+R' Toggle Recursive • 'q' Quit"
		}
	} else {
		if m.width < 80 {
			footerText = "'Tab' Tabs • '↑↓' Navigate • 'Enter' Select • 'q' Quit"
		} else {
			footerText = "'Tab' Switch Tabs • '←↑↓→' Navigate • 'Ctrl + ←→' Switch Panes • 'Enter' Run/Select • 'q' Quit"
		}
	}
	
	footer := lipgloss.NewStyle().
		Foreground(m.theme.Current.Primary).
		MarginTop(1).
		Align(lipgloss.Center).
		Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, body, footer)
}

// renderScriptsTab renders the scripts tab.
func (m Model) renderScriptsTab() string {
	// Simplified responsive layout
	var leftPanelWidth, rightPanelWidth int
	
	if m.width < 80 {
		// Small terminals - single panel
		leftPanelWidth = m.width - 4
		rightPanelWidth = 0
	} else {
		// Large terminals - dual panel
		leftPanelWidth = (m.width / 3) - 2
		rightPanelWidth = ((m.width * 2) / 3) - 2
	}

	panelHeight := m.height - 10
	if m.height < 15 {
		panelHeight = m.height - 6
	}

	// Create breadcrumb
	maxBreadcrumbWidth := leftPanelWidth - 4
	if maxBreadcrumbWidth < 10 {
		maxBreadcrumbWidth = 10
	}

	var breadcrumbText string
	if m.recursiveMode {
		if m.width < 50 {
			breadcrumbText = "📁 All"
		} else {
			breadcrumbText = "📁 All Scripts (Recursive)"
		}
	} else {
		truncatedPath := m.currentPath
		if len(truncatedPath) > maxBreadcrumbWidth {
			truncatedPath = "..." + truncatedPath[len(truncatedPath)-maxBreadcrumbWidth+3:]
		}
		breadcrumbText = truncatedPath
	}

	breadcrumb := lipgloss.NewStyle().
		Faint(true).
		Width(maxBreadcrumbWidth).
		Render(breadcrumbText)

	// Search input - simplified approach
	var searchSection string
	availableSearchWidth := leftPanelWidth - 8 // Account for padding and prefix
	
	if (m.searchActive || m.searchInput.Value() != "") {
		m.searchInput.SetWidth(availableSearchWidth)
		if m.width < 50 {
			// Very small terminals use minimal view
			searchSection = "\n" + m.searchInput.ViewMinimal()
		} else {
			// All other terminals use standard view (no borders)
			searchSection = "\n" + m.searchInput.View()
		}
	}

	leftContent := breadcrumb + searchSection + "\n" + m.list.View()

	// Handle right panel for small terminals
	var rightContent string
	if rightPanelWidth > 0 {
		rightContent = m.vp.View()
	}

	// Focus highlighting
	leftBorderColor := m.theme.Current.Primary
	rightBorderColor := m.theme.Current.Primary
	if m.focus == FocusList || m.focus == FocusSearch {
		leftBorderColor = m.theme.Current.Accent
		rightBorderColor = lipgloss.Color("244")
	} else {
		leftBorderColor = lipgloss.Color("244")
		rightBorderColor = m.theme.Current.Accent
	}

	// Standard panel styling - no complex breakpoints
	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(leftBorderColor).
		Width(leftPanelWidth).
		Height(panelHeight).
		Padding(1, 2).
		Render(leftContent)

	// Return only left panel for small terminals
	if rightPanelWidth <= 0 {
		return leftPanel
	}

	// Right panel
	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(rightBorderColor).
		Width(rightPanelWidth).
		Height(panelHeight).
		Padding(1, 2).
		Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

// renderOptionsTab renders the options tab.
func (m Model) renderOptionsTab() string {
	leftPanelWidth := (m.width / 3) - 2
	rightPanelWidth := ((m.width * 2) / 3) - 2
	panelHeight := m.height - 10

	leftContent := "Themes\n" + m.list.View()

	var rightContent string
	if m.selectedCategory == "color_schemes" {
		rightContent = m.optionsRightList.View()
	} else {
		rightContent = m.vp.View()
	}

	// Focus highlighting
	leftBorderColor := m.theme.Current.Primary
	rightBorderColor := m.theme.Current.Primary
	if m.focus == FocusList {
		leftBorderColor = m.theme.Current.Accent
		rightBorderColor = lipgloss.Color("244")
	} else {
		leftBorderColor = lipgloss.Color("244")
		rightBorderColor = m.theme.Current.Accent
	}

	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(leftBorderColor).
		Width(leftPanelWidth).
		Height(panelHeight).
		Padding(1, 2).
		Render(leftContent)

	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(rightBorderColor).
		Width(rightPanelWidth).
		Height(panelHeight).
		Padding(1, 2).
		Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

// renderAboutTab renders the about tab.
func (m Model) renderAboutTab() string {
	grey := lipgloss.Color("244")
	panelHeight := m.height - 6

	// ASCII art for the logo
	asciiArt := "         ██████╗  ██████╗  ██████╗██╗  ██╗███████╗████████╗\n" +
		"         ██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝\n" +
		"         ██████╔╝██║   ██║██║     █████╔╝ █████╗     ██║   \n" +
		"         ██╔══██╗██║   ██║██║     ██╔═██╗ ██╔══╝     ██║   \n" +
		"         ██║  ██║╚██████╔╝╚██████╗██║  ██╗███████╗   ██║   \n" +
		"         ╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   \n" +
		"                                                           \n" +
		" ██████╗  ██████╗ ██╗    ██╗███████╗██████╗      ██╗███╗   ██╗ ██████╗\n" +
		" ██╔══██╗██╔═══██╗██║    ██║██╔════╝██╔══██╗     ██║████╗  ██║██╔════╝\n" +
		" ██████╔╝██║   ██║██║ █╗ ██║█████╗  ██████╔╝     ██║██╔██╗ ██║██║     \n" +
		" ██╔═══╝ ██║   ██║██║███╗██║██╔══╝  ██╔══██╗     ██║██║╚██╗██║██║     \n" +
		" ██║     ╚██████╔╝╚███╔███╔╝███████╗██║  ██║     ██║██║ ╚████║╚██████╗\n" +
		" ╚═╝      ╚═════╝  ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝"

	aboutContent := asciiArt + "\n\nA cross-platform script browser powered by RocketPowerInc.\n\nBuilt with Go and powered by Charm_ Bubble Tea framework.\n\nVisit us at https://github.com/rocketpowerinc"

	aboutPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(grey).
		Foreground(grey).
		Width(m.width - 4).
		Height(panelHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(aboutContent)

	return aboutPanel
}
