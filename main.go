package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type scriptItem struct {
	name string
	path string
}

func (s scriptItem) Title() string       { return s.name }
func (s scriptItem) Description() string { return s.path }
func (s scriptItem) FilterValue() string { return s.name }

type focusArea int

const (
	focusList focusArea = iota
	focusPreview
)

type parentNav struct {
	path  string
	index int
}

// Cache for script contents to avoid repeated file reads
type scriptCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

func newScriptCache() *scriptCache {
	return &scriptCache{
		cache: make(map[string]string),
	}
}

func (sc *scriptCache) get(path string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	content, exists := sc.cache[path]
	return content, exists
}

func (sc *scriptCache) set(path, content string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache[path] = content
}

func (sc *scriptCache) clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache = make(map[string]string)
}

type model struct {
	list        list.Model
	vp          viewport.Model
	width       int
	height      int
	activeTab   int
	tabs        []string
	scriptItems []list.Item
	focus       focusArea
	currentPath string
	parentPaths []parentNav
	cache       *scriptCache
}

type outputMsg string

var pink = lipgloss.Color("205")
var purple = lipgloss.Color("93")

var (
	borderStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).BorderForeground(pink)
	tabActiveStyle   = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1).Foreground(pink)
	tabInactiveStyle = lipgloss.NewStyle().Faint(true).Padding(0, 1).Foreground(pink)
	tabBarStyle      = lipgloss.NewStyle().MarginBottom(1).Foreground(pink)
)

func highlightScript(content, ext string) string {
	if content == "" {
		return content
	}

	// No syntax highlighting - just return clean content
	return sanitizeContentForBorders(content)
}

// sanitizeContentForBorders removes or escapes characters that can break border rendering
// This is especially important for shell scripts which may contain special characters
func sanitizeContentForBorders(content string) string {
	// Simple, reliable content cleaning
	content = strings.ReplaceAll(content, "\r\n", "\n") // Normalize line endings
	content = strings.ReplaceAll(content, "\r", "\n")   // Handle old Mac line endings
	content = strings.ReplaceAll(content, "\x00", "")   // Remove null bytes

	// Remove any stray ANSI sequences that might interfere
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Simple line length limit to prevent layout issues
		if len(line) > 200 { // Conservative limit
			lines[i] = line[:200] + "..."
		}
	}

	return strings.Join(lines, "\n")
}

// getScriptbinPath returns the path where scriptbin should be stored
// It tries to find the go-pwr installation directory and places scriptbin there
func getScriptbinPath() (string, error) {
	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)

	// Try to find the go-pwr source directory
	// Check if we're running from the development directory (has main.go)
	if _, err := os.Stat(filepath.Join(exeDir, "main.go")); err == nil {
		// We're in the development directory
		return filepath.Join(exeDir, "scriptbin"), nil
	}

	// Check if there's a go-pwr directory in common locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Common locations where go-pwr might be cloned
	possiblePaths := []string{
		filepath.Join(homeDir, "go-pwr"),
		filepath.Join(homeDir, "Github-pwr", "go-pwr"),
		filepath.Join(homeDir, "projects", "go-pwr"),
		filepath.Join(homeDir, "src", "go-pwr"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(filepath.Join(path, "main.go")); err == nil {
			return filepath.Join(path, "scriptbin"), nil
		}
	}

	// If we can't find the source directory, create scriptbin next to the executable
	// This handles the case where go-pwr is installed via go install
	return filepath.Join(exeDir, "scriptbin"), nil
}

func ensureRepo() error {
	root, err := getScriptbinPath()
	if err != nil {
		return err
	}

	// Always remove and re-clone for fresh content
	if _, err := os.Stat(root); err == nil {
		if err := os.RemoveAll(root); err != nil {
			return fmt.Errorf("failed to remove old scriptbin: %v", err)
		}
	}

	cmd := exec.Command("git", "clone", "https://github.com/rocketpowerinc/scriptbin.git", root)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v\n%s", err, string(out))
	}

	return nil
}

func getScriptItems(root string) []list.Item {
	var items []list.Item
	entries, err := os.ReadDir(root)
	if err != nil {
		return items
	}

	// Pre-allocate slice with estimated capacity
	items = make([]list.Item, 0, len(entries))

	// Sort entries: directories first, then files, both alphabetically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}
		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files and directories
		if strings.HasPrefix(name, ".") {
			continue
		}

		path := filepath.Join(root, name)
		if entry.IsDir() {
			items = append(items, scriptItem{name: name + "/", path: path})
		} else {
			// Only include supported script files
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".sh" || ext == ".ps1" || ext == ".bat" || ext == ".cmd" {
				items = append(items, scriptItem{name: name, path: path})
			}
		}
	}
	return items
}

func readScript(path string, cache *scriptCache) string {
	// Check cache first
	if content, exists := cache.get(path); exists {
		return content
	}

	data, err := os.ReadFile(path)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading file: %v", err)
		cache.set(path, errMsg) // Cache errors too to avoid repeated attempts
		return errMsg
	}

	content := string(data)
	cache.set(path, content)
	return content
}

func (m *model) setSizes() {
	// SIMPLE STATIC LAYOUT - same calculations as View()
	// Account for borders in width calculations to prevent overflow
	leftPanelWidth := (m.width / 3) - 2  // -2 for left panel border
	rightPanelWidth := ((m.width * 2) / 3) - 2  // -2 for right panel border

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
	m.vp.Width = rightContentWidth
	m.vp.Height = m.height - 10
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.setSizes()
		return m, nil

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Only handle clicks on the first line (tab bar area)
			if msg.Y == 0 {
				x := 0
				for i, tab := range m.tabs {
					// Calculate the exact rendered width of this tab
					style := tabInactiveStyle.Copy()
					if i == m.activeTab {
						style = style.Inherit(tabActiveStyle)
					}
					renderedTab := style.Render(tab)
					tabWidth := lipgloss.Width(renderedTab)

					// Check if click is within this specific tab's bounds
					if msg.X >= x && msg.X < x+tabWidth {
						m.activeTab = i
						if i == 0 {
							m.list.SetItems(m.scriptItems)
							if sel, ok := m.list.SelectedItem().(scriptItem); ok {
								if isScriptFile(sel.name) {
									ext := filepath.Ext(sel.path)
									m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
								} else {
									m.vp.SetContent("Select a script to preview...")
								}
							}
						} else {
							m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
						}
						break
					}

					// Move to next tab position (tab width + 2 spaces separator)
					x += tabWidth + 2
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			if m.activeTab == 0 {
				m.list.SetItems(m.scriptItems)
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
						ext := filepath.Ext(sel.path)
						m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
					} else {
						m.vp.SetContent("Select a script to preview...")
					}
				}
			} else {
				m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
			}
		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			if m.activeTab == 0 {
				m.list.SetItems(m.scriptItems)
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
						ext := filepath.Ext(sel.path)
						m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
					} else {
						m.vp.SetContent("Select a script to preview...")
					}
				}
			} else {
				m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
			}
		case "ctrl+tab":
			// Switch focus between list and preview panes (only in Scripts tab)
			if m.activeTab == 0 {
				if m.focus == focusList {
					m.focus = focusPreview
				} else {
					m.focus = focusList
				}
			}
		case "r":
			if m.activeTab == 0 && m.focus == focusList {
				m.list, _ = m.list.Update(msg) // Ensure list state is updated
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					fmt.Println("Selected script:", sel.name, sel.path) // Debug print
					if !strings.HasSuffix(sel.name, "/") {
						m.vp.SetContent("Running script in a new terminal window...")
						go func() {
							var cmd *exec.Cmd
							if isWindows() {
								if strings.HasSuffix(sel.name, ".ps1") {
									cmd = exec.Command("cmd", "/C", "start", "powershell", "-NoExit", "-Command", "Clear-Host; "+sel.path+"; Write-Host ''; Read-Host 'Press Enter to exit'")
								} else {
									cmd = exec.Command("cmd", "/C", "start", "cmd", "/K", "cls && bash -l "+sel.path+" & pause")
								}
							} else if isMac() {
								// Improved macOS terminal handling
								scriptCmd := sel.path
								if strings.HasSuffix(sel.name, ".ps1") {
									scriptCmd = "pwsh " + sel.path
								} else {
									scriptCmd = "bash " + sel.path
								}
								osaCmd := fmt.Sprintf(`tell application "Terminal"
    do script "clear; %s; echo; read -n 1 -s -r -p 'Press any key to exit...'"
    activate
end tell`, scriptCmd)
								cmd = exec.Command("osascript", "-e", osaCmd)
							} else if isLinux() {
								// Linux: try common terminals
								term := ""
								for _, candidate := range []string{"gnome-terminal", "konsole", "x-terminal-emulator"} {
									if _, err := exec.LookPath(candidate); err == nil {
										term = candidate
										break
									}
								}
								if term == "" {
									fmt.Println("No supported terminal emulator found.")
									return
								}
								if strings.HasSuffix(sel.name, ".ps1") {
									cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; pwsh "+sel.path+"; echo; read -p 'Press Enter to exit'")
								} else {
									cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; bash "+sel.path+"; echo; read -p 'Press Enter to exit'")
								}
							}
							err := cmd.Start()
							if err != nil {
								fmt.Printf("\nError opening terminal: %v\n", err)
							}
						}()
					}
				}
			}
		case "left":
			// Check if we're in Scripts tab and can navigate tabs
			if m.focus == focusPreview || m.currentPath == "" || len(m.parentPaths) == 0 {
				// Navigate to previous tab
				m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
				if m.activeTab == 0 {
					m.list.SetItems(m.scriptItems)
					if sel, ok := m.list.SelectedItem().(scriptItem); ok {
						if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
							ext := filepath.Ext(sel.path)
							m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
					}
				} else {
					m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
				}
			} else {
				// Go back to parent directory if possible
				if m.activeTab == 0 && m.focus == focusList && m.currentPath != "" && len(m.parentPaths) > 0 {
					parent := m.parentPaths[len(m.parentPaths)-1]
					m.parentPaths = m.parentPaths[:len(m.parentPaths)-1]
					m.scriptItems = getScriptItems(parent.path)
					m.list.SetItems(m.scriptItems)
					m.list.Select(parent.index) // Restore previous selection
					m.currentPath = parent.path
					// Show preview for selected item
					if sel, ok := m.list.SelectedItem().(scriptItem); ok {
						if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
							ext := filepath.Ext(sel.path)
							m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
					} else {
						m.vp.SetContent("Select a script to preview...")
					}
					return m, nil
				}
			}
		case "right":
			// Check if we can navigate into directory or should switch tabs
			canNavigateDir := false
			if m.activeTab == 0 && m.focus == focusList {
				if sel, ok := m.list.SelectedItem().(scriptItem); ok && strings.HasSuffix(sel.name, "/") {
					canNavigateDir = true
				}
			}
			
			if canNavigateDir {
				// Traverse into selected directory
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					if strings.HasSuffix(sel.name, "/") {
						m.parentPaths = append(m.parentPaths, parentNav{path: m.currentPath, index: m.list.Index()})
						m.currentPath = sel.path
						m.scriptItems = getScriptItems(sel.path)
						m.list.SetItems(m.scriptItems)
						m.list.ResetSelected()
						// Show preview for first item if it's a script
						if len(m.scriptItems) > 0 {
							if first, ok := m.scriptItems[0].(scriptItem); ok {
								if strings.HasSuffix(first.name, ".sh") || strings.HasSuffix(first.name, ".ps1") {
									ext := filepath.Ext(first.path)
									m.vp.SetContent(highlightScript(readScript(first.path, m.cache), ext))
								} else {
									m.vp.SetContent("Select a script to preview...")
								}
							}
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
						return m, nil
					}
				}
			} else {
				// Navigate to next tab
				m.activeTab = (m.activeTab + 1) % len(m.tabs)
				if m.activeTab == 0 {
					m.list.SetItems(m.scriptItems)
					if sel, ok := m.list.SelectedItem().(scriptItem); ok {
						if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
							ext := filepath.Ext(sel.path)
							m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
					}
				} else {
					m.vp.SetContent("A cross-platform script browser powered by RocketPowerInc.")
				}
			}
		case "enter":
			if m.activeTab == 0 && m.focus == focusList {
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					fmt.Println("Selected script:", sel.name, sel.path) // Debug print
					if !strings.HasSuffix(sel.name, "/") {
						// Preview the script file
						return m, func() tea.Msg {
							ext := filepath.Ext(sel.path)
							return outputMsg(highlightScript(readScript(sel.path, m.cache), ext))
						}
					}
				}
			}
		case "up":
			if m.focus == focusList && m.activeTab == 0 {
				prevIndex := m.list.Index()
				m.list, cmd = m.list.Update(msg)
				newIndex := m.list.Index()
				if prevIndex != newIndex {
					if sel, ok := m.list.SelectedItem().(scriptItem); ok {
						if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
							ext := filepath.Ext(sel.path)
							m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
					}
				}
			} else if m.focus == focusPreview && m.activeTab == 0 {
				m.vp.LineUp(1)
			}
			return m, cmd
		case "down":
			if m.focus == focusList && m.activeTab == 0 {
				prevIndex := m.list.Index()
				m.list, cmd = m.list.Update(msg)
				newIndex := m.list.Index()
				if prevIndex != newIndex {
					if sel, ok := m.list.SelectedItem().(scriptItem); ok {
						if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
							ext := filepath.Ext(sel.path)
							m.vp.SetContent(highlightScript(readScript(sel.path, m.cache), ext))
						} else {
							m.vp.SetContent("Select a script to preview...")
						}
					}
				}
			} else if m.focus == focusPreview && m.activeTab == 0 {
				m.vp.LineDown(1)
			}
			return m, cmd
		case "page_up":
			if m.focus == focusPreview && m.activeTab == 0 {
				for i := 0; i < 10; i++ {
					m.vp.LineUp(1)
				}
			}
			return m, cmd
		case "page_down":
			if m.focus == focusPreview && m.activeTab == 0 {
				for i := 0; i < 10; i++ {
					m.vp.LineDown(1)
				}
			}
			return m, cmd
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case outputMsg:
		m.vp.SetContent(strings.TrimSpace(string(msg)))
	}

	m.list, _ = m.list.Update(msg)
	m.vp, _ = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var tabLabels []string
	for i, name := range m.tabs {
		label := name
		style := tabInactiveStyle.Copy()
		if i == m.activeTab {
			style = style.Inherit(tabActiveStyle)
		}
		tabLabels = append(tabLabels, style.Render(label))
	}
	tabBar := tabBarStyle.Render(strings.Join(tabLabels, "  "))

	var body string
	if m.activeTab == 0 {
		// SIMPLE STATIC LAYOUT - Back to basics that worked
		// Account for borders in width calculations to prevent overflow
		leftPanelWidth := (m.width / 3) - 2  // -2 for left panel border
		rightPanelWidth := ((m.width * 2) / 3) - 2  // -2 for right panel border
		panelHeight := m.height - 10

		// Truncate breadcrumb to prevent it from affecting panel size
		maxBreadcrumbWidth := leftPanelWidth - 8 // Account for border + padding
		if maxBreadcrumbWidth < 10 {
			maxBreadcrumbWidth = 10
		}

		truncatedPath := m.currentPath
		if len(truncatedPath) > maxBreadcrumbWidth {
			// Truncate from the beginning, keeping the end
			truncatedPath = "..." + truncatedPath[len(truncatedPath)-maxBreadcrumbWidth+3:]
		}

		// Enforce maximum width to guarantee no overflow
		breadcrumb := lipgloss.NewStyle().
			Faint(true).
			Width(maxBreadcrumbWidth).
			Render(truncatedPath)

		// Content with controlled breadcrumb that cannot affect panel sizing
		leftContent := breadcrumb + "\n" + m.list.View()
		rightContent := m.vp.View()

		// Create panels with explicit dimensions - simple and stable
		// Add visual feedback for focused pane
		leftBorderColor := pink
		rightBorderColor := pink
		if m.activeTab == 0 {
			if m.focus == focusList {
				leftBorderColor = purple // Highlight focused pane
				rightBorderColor = lipgloss.Color("244") // Dim unfocused pane
			} else {
				leftBorderColor = lipgloss.Color("244") // Dim unfocused pane
				rightBorderColor = purple // Highlight focused pane
			}
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

		// Simple join - no Place() complications
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	} else {
		grey := lipgloss.Color("244")
		aboutStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(grey).
			Foreground(grey).
			Align(lipgloss.Center, lipgloss.Center).
			Width(m.width / 2).
			Height(m.height / 2)
		body = lipgloss.Place(m.width, m.height-10, lipgloss.Center, lipgloss.Center,
			aboutStyle.Render("A cross-platform script browser powered by RocketPowerInc.\n\nMade with Bubble Tea, Lipgloss, and Go. \n\nVisit us at https://github.com/rocketpowerinc"),
		)
	}

	footer := lipgloss.NewStyle().Foreground(pink).MarginTop(1).Align(lipgloss.Center).Render("Tab/Shift+Tab Switch Tabs • ← → Navigate/Switch Tabs • ↑↓ Select/Scroll • Ctrl+Tab Switch Panes • Enter Preview • r Run Script • q Quit")

	if m.activeTab == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			tabBar,
			body,
			footer,
		)
	} else {
		return lipgloss.JoinVertical(lipgloss.Left,
			tabBar,
			body,
		)
	}
}

type scriptDelegate struct{}

func (d scriptDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s, ok := item.(scriptItem)
	if !ok {
		return
	}
	var style lipgloss.Style
	if strings.HasSuffix(s.name, "/") {
		style = lipgloss.NewStyle().Foreground(purple)
	} else {
		style = lipgloss.NewStyle().Foreground(pink)
	}
	if index == m.Index() {
		style = style.Bold(true).Underline(true)
	}
	fmt.Fprint(w, style.Render(s.name))
}

func (d scriptDelegate) Height() int               { return 1 }
func (d scriptDelegate) Spacing() int              { return 0 }
func (d scriptDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

// Helper function to check if a file is a supported script type
func isScriptFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".sh" || ext == ".ps1" || ext == ".bat" || ext == ".cmd"
}

func main() {
	if err := ensureRepo(); err != nil {
		fmt.Println("Error cloning repo:", err)
		os.Exit(1)
	}

	// Get the proper scriptbin path
	scriptbinPath, err := getScriptbinPath()
	if err != nil {
		fmt.Println("Error getting scriptbin path:", err)
		os.Exit(1)
	}

	tabs := []string{"Scripts", "About"}
	scriptItems := getScriptItems(scriptbinPath)
	cache := newScriptCache()

	listModel := list.New(scriptItems, scriptDelegate{}, 0, 0)
	listModel.Title = "" // Remove the "Scripts" title
	listModel.SetShowHelp(false)
	listModel.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Select a script to preview...")

	// Initial preview for first script
	if len(scriptItems) > 0 {
		if s, ok := scriptItems[0].(scriptItem); ok {
			if isScriptFile(s.name) {
				ext := filepath.Ext(s.path)
				vp.SetContent(highlightScript(readScript(s.path, cache), ext))
			}
		}
	}

	m := model{
		list:        listModel,
		vp:          vp,
		tabs:        tabs,
		scriptItems: scriptItems,
		activeTab:   0,
		focus:       focusList, // Initialize focus to the list pane
		currentPath: scriptbinPath,
		parentPaths: []parentNav{},
		cache:       cache,
	}

	if err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

