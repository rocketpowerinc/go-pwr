package ui

import (
	"strings"

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
		switch msg.String() {
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
			if m.activeTab == 0 && m.focus == FocusList && len(m.parentPaths) > 0 {
				m.navigateToParent()
				return m, nil
			}
		case "right":
			if m.activeTab == 0 && m.focus == FocusList {
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

	m.list, _ = m.list.Update(msg)
	m.optionsRightList, _ = m.optionsRightList.Update(msg)
	m.vp, _ = m.vp.Update(msg)
	return m, cmd
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
			if sel.IsDirectory() {
				m.navigateIntoDirectory(sel)
			} else {
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

	// Render footer
	footer := lipgloss.NewStyle().
		Foreground(m.theme.Current.Primary).
		MarginTop(1).
		Align(lipgloss.Center).
		Render("'Tab' Switch Tabs • '←↑↓→' Navigate • 'Ctrl + ←→' Switch Panes • 'Enter' Run/Select • 'q' Quit")

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, body, footer)
}

// renderScriptsTab renders the scripts tab.
func (m Model) renderScriptsTab() string {
	leftPanelWidth := (m.width / 3) - 2
	rightPanelWidth := ((m.width * 2) / 3) - 2
	panelHeight := m.height - 10

	// Create breadcrumb
	maxBreadcrumbWidth := leftPanelWidth - 8
	if maxBreadcrumbWidth < 10 {
		maxBreadcrumbWidth = 10
	}

	truncatedPath := m.currentPath
	if len(truncatedPath) > maxBreadcrumbWidth {
		truncatedPath = "..." + truncatedPath[len(truncatedPath)-maxBreadcrumbWidth+3:]
	}

	breadcrumb := lipgloss.NewStyle().
		Faint(true).
		Width(maxBreadcrumbWidth).
		Render(truncatedPath)

	leftContent := breadcrumb + "\n" + m.list.View()
	rightContent := m.vp.View()

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

	aboutContent := "A cross-platform script browser powered by RocketPowerInc.\n\nMade with Bubble Tea, Lipgloss, and Go.\n\nVisit us at https://github.com/rocketpowerinc"

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
