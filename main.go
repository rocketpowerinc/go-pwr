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

// --- Optimized Syntax Highlighting ---
var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))   // Blue
	commentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))  // Grey
	stringStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))  // Orange
)

// Pre-compiled keyword maps for better performance
var (
	bashKeywords = []string{"if", "then", "else", "fi", "for", "in", "do", "done", "echo", "exit", "while", "case", "function", "export", "source"}
	psKeywords   = []string{"Write-Host", "if", "else", "foreach", "function", "return", "break", "param", "process", "begin", "end", "Get-", "Set-", "New-"}
	batchKeywords = []string{"echo", "set", "if", "else", "goto", "call", "exit", "for", "in", "do"}
)

func highlightScript(content, ext string) string {
	if content == "" {
		return content
	}
	
	lines := strings.Split(content, "\n")
	var keywords []string
	
	switch strings.ToLower(ext) {
	case ".sh":
		keywords = bashKeywords
	case ".ps1":
		keywords = psKeywords
	case ".bat", ".cmd":
		keywords = batchKeywords
	default:
		return content // No highlighting for unknown extensions
	}
	
	for i, line := range lines {
		// Handle comments first (different comment chars for different languages)
		commentChar := "#"
		if ext == ".bat" || ext == ".cmd" {
			commentChar = "REM"
		}
		
		var commentIdx int
		if commentChar == "#" {
			commentIdx = strings.Index(line, "#")
		} else {
			commentIdx = strings.Index(strings.ToUpper(line), "REM")
		}
		
		isComment := commentIdx == 0 || (commentIdx > 0 && strings.TrimSpace(line[:commentIdx]) == "")
		
		if commentIdx != -1 {
			comment := line[commentIdx:]
			line = line[:commentIdx] + commentStyle.Render(comment)
		}
		
		if !isComment {
			// Apply keyword highlighting
			for _, kw := range keywords {
				// Use word boundaries to avoid partial matches
				if strings.Contains(line, kw) {
					line = strings.ReplaceAll(line, kw, keywordStyle.Render(kw))
				}
			}
			
			// Handle string highlighting
			line = highlightStrings(line)
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func highlightStrings(line string) string {
	var out strings.Builder
	inString := false
	inSingleQuote := false
	
	for _, r := range line {
		switch r {
		case '"':
			if !inSingleQuote {
				inString = !inString
			}
			out.WriteString(stringStyle.Render(string(r)))
		case '\'':
			if !inString {
				inSingleQuote = !inSingleQuote
			}
			out.WriteString(stringStyle.Render(string(r)))
		default:
			if inString || inSingleQuote {
				out.WriteString(stringStyle.Render(string(r)))
			} else {
				out.WriteRune(r)
			}
		}
	}
	return out.String()
}

// getScriptbinPath returns the path where scriptbin should be stored
// It uses a consistent location in the user's home directory
func getScriptbinPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	// Create a .go-pwr directory in the user's home for storing scriptbin
	appDir := filepath.Join(homeDir, ".go-pwr")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory: %v", err)
	}
	
	return filepath.Join(appDir, "scriptbin"), nil
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
	listW := m.width / 3
	if listW < 20 {
		listW = 20
	}
	vpW := m.width - listW
	if vpW < 20 {
		vpW = 20
	}
	m.list.SetSize(listW, m.height-10)
	m.vp.Width = vpW
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
			x := 1
			for i, tab := range m.tabs {
				end := x + len(tab) + 2
				if msg.X >= x && msg.X < end {
					m.activeTab = i
					if i == 0 {
						m.list.SetItems(m.scriptItems)
					}
					break
				}
				x = end + 2
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
		case "right":
			// Traverse into selected directory if possible
			if m.activeTab == 0 && m.focus == focusList {
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
			if m.focus == focusList {
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
			} else if m.focus == focusPreview {
				m.vp.LineUp(1)
			}
			return m, cmd
		case "down":
			if m.focus == focusList {
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
			} else if m.focus == focusPreview {
				m.vp.LineDown(1)
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

	centerStyle := lipgloss.NewStyle().Align(lipgloss.Center).Height(m.height-10)
	breadcrumb := lipgloss.NewStyle().Faint(true).Render(m.currentPath)

	var body string
	if m.activeTab == 0 {
		leftW := m.width / 3
		rightW := m.width - leftW - 4
		panelHeight := m.height - 10

		left := borderStyle.Width(leftW).Height(panelHeight).Render(
			breadcrumb + m.list.View(),
		)

		rightBorderStyle := borderStyle.Copy().BorderRight(true)
		right := rightBorderStyle.Width(rightW).Height(panelHeight).Render(centerStyle.Render(m.vp.View()))

		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
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

	footer := lipgloss.NewStyle().Foreground(pink).MarginTop(1).Align(lipgloss.Center).Render("← → or 🖱️ Click Tabs • ↑↓ Select • Enter Preview • r Run Script • q Quit")

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

