package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rocketpowerinc/go-pwr/internal/config"
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
			if m.repositoryInputActive && m.activeTab == 1 {
				m.repositoryInputActive = false
				m.repositoryInput.SetActive(false)
				m.focus = FocusPreview
				return m, nil
			} else if m.repositoryViewActive && m.activeTab == 1 {
				m.repositoryViewActive = false
				m.focus = FocusPreview
				return m, nil
			} else if m.repositoryResetActive && m.activeTab == 1 {
				m.repositoryResetActive = false
				m.focus = FocusPreview
				return m, nil
			} else if m.searchActive && m.activeTab == 0 {
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

		// Handle repository input if repository input is active
		if m.repositoryInputActive && m.activeTab == 1 {
			switch msg.String() {
			case "enter":
				// Validate and save the repository URL
				if m.repositoryInput.Validate() {
					url := m.repositoryInput.Value()
					
					// Show loading message
					m.vp.SetContent("Saving and loading new repository...\n\nPlease wait while we switch to: " + url)
					
					if err := config.SaveRepoURL(url); err != nil {
						m.repositoryInput.SetError("Failed to save repository: " + err.Error())
						return m, nil
					} else {
						// Update the config immediately
						m.config.RepoURL = url
						m.repositoryInputActive = false
						m.repositoryInput.SetActive(false)
						m.repositoryResetActive = true // Show result in dedicated screen
						m.repositoryViewActive = false
						m.focus = FocusPreview
						
						// Try to refresh repository immediately
						if err := m.refreshRepository(); err != nil {
							m.vp.SetContent("✅ Repository saved, but failed to load scripts: " + err.Error() + "\n\nNew repository: " + url + "\n\nPlease restart the application to see the changes.")
						} else {
							m.vp.SetContent(fmt.Sprintf("✅ Custom Repository Successfully Set!\n\n🔄 New Repository URL:\n%s\n\n📁 Scripts Location:\n%s\n\n✨ Scripts have been refreshed and are ready to use.\nSwitch to the Scripts tab to see your custom content.", 
								url, m.config.ScriptbinPath))
						}
						
						// Update the repository items to show the new current repo
						m.updateRepositoryItems()
						return m, nil
					}
				}
				return m, nil
			default:
				// Pass all other keys (except escape, handled above) to repository input
				cmd = m.repositoryInput.Update(msg)
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
		case "ctrl+tab":
			// Tab switching - same as regular tab but with ctrl modifier
			m.switchTab((m.activeTab + 1) % len(m.tabs))
		case "ctrl+left", "cmd+left", "alt+left", "shift+left", "ctrl+h":
			if m.activeTab == 0 || m.activeTab == 1 {
				m.focus = FocusList
			}
		case "ctrl+right", "cmd+right", "alt+right", "shift+right", "ctrl+l":
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

	// Update repository input if active
	if m.repositoryInputActive {
		m.repositoryInput.Update(msg)
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
	} else if m.focus == FocusPreview && m.activeTab == 1 && !m.repositoryInputActive {
		m.optionsRightList, cmd = m.optionsRightList.Update(msg)
	} else if m.focus == FocusRepositoryInput && m.activeTab == 1 {
		// Repository input handles its own navigation
		cmd = m.repositoryInput.Update(msg)
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
				m.vp.SetContent("Select a color scheme from the list on the right to apply it instantly!\n\nUse Ctrl+Right or Ctrl+L to switch to the right pane.")
				// Automatically switch focus to right pane for easier navigation
				m.focus = FocusPreview
			} else if sel.Category == "repository" {
				m.optionsRightList.SetItems(m.repositoryItems)
				m.vp.SetContent("Configure your script repository below:\n\nUse the options on the right to manage your repository.\n\nUse Ctrl+Right or Ctrl+L to switch to the right pane.")
				// Automatically switch focus to right pane for easier navigation
				m.focus = FocusPreview
			}
		}
	} else if m.activeTab == 1 && m.focus == FocusPreview {
		// Options tab - apply selected option
		if m.selectedCategory == "color_schemes" {
			if sel, ok := m.optionsRightList.SelectedItem().(components.OptionItem); ok {
				m.applyColorScheme(sel.Name)
			}
		} else if m.selectedCategory == "repository" {
			if sel, ok := m.optionsRightList.SelectedItem().(components.OptionItem); ok {
				m.handleRepositoryAction(sel.Action)
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
	} else if m.activeTab == 1 && m.repositoryInputActive {
		footerText = "'Enter' Save Repository • 'Esc' Cancel • Type repository URL"
	} else if m.activeTab == 1 && (m.repositoryViewActive || m.repositoryResetActive) {
		footerText = "'Esc' Back to Repository Options • 'Tab' Switch Tabs • 'q' Quit"
	} else if m.activeTab == 0 {
		if m.width < 80 {
			// Short footer for small terminals
			footerText = "'Tab' Tabs • '↑↓' Navigate • 'Enter' Run • 'Ctrl+F' Search • 'q' Quit"
		} else if m.width < 120 {
			// Medium footer for medium terminals
			footerText = "'Tab' Switch • '←↑↓→' Navigate • 'Enter' Run/Select • 'Ctrl+F' Search • 'Ctrl+R' Recursive • 'Ctrl+H/L' Switch Panes • 'q' Quit"
		} else {
			// Full footer for large terminals
			footerText = "'Tab' Switch Tabs • '←↑↓→' Navigate • 'Ctrl+H/L' Switch Panes • 'Enter' Run/Select • 'Ctrl+F' Search • 'Ctrl+R' Toggle Recursive • 'q' Quit"
		}
	} else {
		if m.width < 80 {
			footerText = "'Tab' Tabs • '↑↓' Navigate • 'Enter' Select • 'q' Quit"
		} else {
			footerText = "'Tab' Switch Tabs • '←↑↓→' Navigate • 'Ctrl+H/L' Switch Panes • 'Enter' Run/Select • 'q' Quit"
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

	leftContent := m.list.View()

	var rightContent string
	if m.selectedCategory == "color_schemes" {
		rightContent = m.optionsRightList.View()
	} else if m.selectedCategory == "repository" {
		if m.repositoryInputActive {
			// Show repository input interface
			m.repositoryInput.SetWidth(rightPanelWidth - 8)
			inputView := m.repositoryInput.View()
			rightContent = "Enter Repository URL:\n\n" + inputView + "\n\n" + m.vp.View()
		} else if m.repositoryViewActive || m.repositoryResetActive {
			// Show dedicated screen for repository info or reset result
			rightContent = m.vp.View() + "\n\n" + "Press Esc to go back to repository options"
		} else {
			rightContent = m.optionsRightList.View()
		}
	} else {
		rightContent = m.vp.View()
	}

	// Focus highlighting
	leftBorderColor := m.theme.Current.Primary
	rightBorderColor := m.theme.Current.Primary
	if m.focus == FocusList {
		leftBorderColor = m.theme.Current.Accent
		rightBorderColor = lipgloss.Color("244")
	} else if m.focus == FocusRepositoryInput {
		leftBorderColor = lipgloss.Color("244")
		rightBorderColor = m.theme.Current.Accent
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
	panelHeight := m.height - 6

	// ASCII art for the logo - raw string without pre-styling
	asciiArt := "      ██████╗  ██████╗  ██████╗██╗  ██╗███████╗████████╗\n" +
		"      ██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝\n" +
		"      ██████╔╝██║   ██║██║     █████╔╝ █████╗     ██║   \n" +
		"      ██╔══██╗██║   ██║██║     ██╔═██╗ ██╔══╝     ██║   \n" +
		"      ██║  ██║╚██████╔╝╚██████╗██║  ██╗███████╗   ██║   \n" +
		"      ╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   \n" +
		"                                                           \n" +
		" ██████╗  ██████╗ ██╗    ██╗███████╗██████╗      ██╗███╗   ██╗ ██████╗\n" +
		" ██╔══██╗██╔═══██╗██║    ██║██╔════╝██╔══██╗     ██║████╗  ██║██╔════╝\n" +
		" ██████╔╝██║   ██║██║ █╗ ██║█████╗  ██████╔╝     ██║██╔██╗ ██║██║     \n" +
		" ██╔═══╝ ██║   ██║██║███╗██║██╔══╝  ██╔══██╗     ██║██║╚██╗██║██║     \n" +
		" ██║     ╚██████╔╝╚███╔███╔╝███████╗██║  ██║     ██║██║ ╚████║╚██████╗\n" +
		" ╚═╝      ╚═════╝  ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝"

	// Description text
	description := "\n\nA cross-platform script browser powered by RocketPowerInc.\n\nBuilt with Go and powered by Charm_ Bubble Tea framework.\n\nVisit us at https://github.com/rocketpowerinc"

	aboutContent := asciiArt + description

	aboutPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Current.Primary).
		Width(m.width - 4).
		Height(panelHeight).
		Foreground(m.theme.Current.Accent).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center).
		Render(aboutContent)

	return aboutPanel
}

// handleRepositoryAction handles repository-related actions
func (m *Model) handleRepositoryAction(action string) {
	switch action {
	case "set_repo":
		// Activate repository input
		m.repositoryInputActive = true
		m.repositoryViewActive = false
		m.repositoryResetActive = false
		m.repositoryInput.SetActive(true)
		m.repositoryInput.Reset()
		// Pre-fill with current URL if it's not the default
		if m.config.RepoURL != config.GetDefaultRepoURL() {
			m.repositoryInput.SetValue(m.config.RepoURL)
		}
		m.focus = FocusRepositoryInput
		m.vp.SetContent("Enter a Git repository URL ending with .git\n\nSupported formats:\n- https://github.com/username/repo.git\n- https://gitlab.com/username/repo.git\n- git@github.com:username/repo.git\n\nPress Enter to save, Esc to cancel")
	case "reset_repo":
		// Activate repository reset view
		m.repositoryResetActive = true
		m.repositoryInputActive = false
		m.repositoryViewActive = false
		m.focus = FocusPreview
		
		// Show loading message
		m.vp.SetContent("Resetting repository to default...\n\nPlease wait while we switch back to RocketPowerInc scriptbin.")
		
		if err := config.ResetToDefaultRepo(); err != nil {
			m.vp.SetContent("❌ Failed to reset repository: " + err.Error())
		} else {
			// Update the config immediately
			defaultRepo := config.GetDefaultRepoURL()
			m.config.RepoURL = defaultRepo
			
			// Try to refresh repository immediately
			if err := m.refreshRepository(); err != nil {
				m.vp.SetContent("✅ Repository reset to default, but failed to refresh scripts: " + err.Error() + "\n\nPlease restart the application to see the changes.")
			} else {
				m.vp.SetContent(fmt.Sprintf("✅ Repository Successfully Reset!\n\n🔄 Reset to Default Repository:\n%s\n\n📁 Scripts Location:\n%s\n\n✨ Scripts have been refreshed and are ready to use.\nSwitch to the Scripts tab to see the default content.", 
					defaultRepo, m.config.ScriptbinPath))
			}
			
			// Update the repository items to show the new current repo
			m.updateRepositoryItems()
		}
	case "view_repo":
		// Activate repository view
		m.repositoryViewActive = true
		m.repositoryInputActive = false
		m.repositoryResetActive = false
		m.focus = FocusPreview
		
		// Provide detailed current repository information with visual formatting
		defaultRepo := config.GetDefaultRepoURL()
		isDefault := m.config.RepoURL == defaultRepo
		
		var headerSection, statusSection, detailsSection, pathSection string
		
		if isDefault {
			headerSection = "🏠 DEFAULT REPOSITORY"
			statusSection = "✅ Status: Using RocketPowerInc's Official Scriptbin"
		} else {
			headerSection = "🔧 CUSTOM REPOSITORY"
			statusSection = "⚙️  Status: Using Custom Repository"
		}
		
		detailsSection = fmt.Sprintf("🌐 Current Repository URL:\n%s\n\n🏠 Default Repository URL:\n%s", 
			m.config.RepoURL, defaultRepo)
		
		pathSection = fmt.Sprintf("📁 Local Scripts Path:\n%s\n\n💡 This is where go-pwr loads scripts from.", 
			m.config.ScriptbinPath)
		
		// Repository type information
		var repoTypeInfo string
		if isDefault {
			repoTypeInfo = "ℹ️  Repository Type: Official RocketPowerInc scriptbin\n   Contains curated, tested scripts for various platforms."
		} else {
			repoTypeInfo = "ℹ️  Repository Type: Custom\n   You can switch back to default using 'Reset to Default'."
		}
		
		m.vp.SetContent(fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s", 
			headerSection, statusSection, detailsSection, pathSection, repoTypeInfo))
	}
}

// updateRepositoryItems updates the repository items list with current config
func (m *Model) updateRepositoryItems() {
	m.repositoryItems = []list.Item{
		components.OptionItem{
			Name:   "Set Custom Repository",
			Desc:   "Use your own script repository",
			Action: "set_repo",
		},
		components.OptionItem{
			Name:   "Reset to Default",
			Desc:   "Use RocketPowerInc's scriptbin",
			Action: "reset_repo",
		},
		components.OptionItem{
			Name:   "Current Repository",
			Desc:   m.config.RepoURL,
			Action: "view_repo",
		},
	}
	// Update the right list if repository category is selected
	if m.selectedCategory == "repository" {
		m.optionsRightList.SetItems(m.repositoryItems)
	}
}
