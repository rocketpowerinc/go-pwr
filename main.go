package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

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

type model struct {
	list        list.Model
	vp          viewport.Model
	width       int
	height      int
	activeTab   int
	tabs        []string
	scriptItems []list.Item
	focus       focusArea
	currentPath string // Track current directory
	parentPaths []parentNav // Track parent directories and selected index
}

type outputMsg string

var pink = lipgloss.Color("205") // ANSI pink
var purple = lipgloss.Color("93") // ANSI purple

var tabColors = []lipgloss.TerminalColor{
	pink, // Scripts - pink
	pink, // About - pink
}

var (
	borderStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).BorderForeground(pink)
	tabActiveStyle   = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1).Foreground(pink)
	tabInactiveStyle = lipgloss.NewStyle().Faint(true).Padding(0, 1).Foreground(pink)
	tabBarStyle      = lipgloss.NewStyle().MarginBottom(1).Foreground(pink)
	headerStyle      = lipgloss.NewStyle().Bold(true).MarginBottom(1).Foreground(pink)
	tabLabelStyle    = lipgloss.NewStyle().Bold(true).MarginBottom(1).Foreground(pink)
)

// --- Syntax Highlighting Helper ---
func highlightScript(content, ext string) string {
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33")) // Blue
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244")) // Grey
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("202")) // Orange

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Highlight comments first
		commentIdx := strings.Index(line, "#")
		isComment := commentIdx == 0 || (commentIdx > 0 && strings.TrimSpace(line[:commentIdx]) == "")
		if commentIdx != -1 {
			comment := line[commentIdx:]
			line = line[:commentIdx] + commentStyle.Render(comment)
		}
		// Highlight keywords (skip comments)
		if !isComment {
			if ext == ".sh" {
				for _, kw := range []string{"if", "then", "else", "fi", "for", "in", "do", "done", "echo", "exit"} {
					line = strings.ReplaceAll(line, kw, keywordStyle.Render(kw))
				}
			} else if ext == ".ps1" {
				for _, kw := range []string{"Write-Host", "if", "else", "foreach", "function", "return", "break"} {
					line = strings.ReplaceAll(line, kw, keywordStyle.Render(kw))
				}
			}
			// Highlight strings in double quotes (skip comments)
			var out strings.Builder
			inString := false
			for _, r := range line {
				if r == '"' {
					inString = !inString
					out.WriteString(stringStyle.Render(string(r)))
				} else if inString {
					out.WriteString(stringStyle.Render(string(r)))
				} else {
					out.WriteRune(r)
				}
			}
			line = out.String()
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func ensureRepo() error {
	root := filepath.Clean("scriptbin")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "https://github.com/rocketpowerinc/scriptbin.git")
		return cmd.Run()
	}
	return nil
}

func getScriptItems(root string) []list.Item {
	var items []list.Item
	entries, err := ioutil.ReadDir(root)
	if err != nil {
		return items
	}
	// Sort: folders first, then files, both alphabetically
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
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			items = append(items, scriptItem{name: entry.Name() + "/", path: path})
		} else if strings.HasSuffix(entry.Name(), ".sh") || strings.HasSuffix(entry.Name(), ".ps1") {
			items = append(items, scriptItem{name: entry.Name(), path: path})
		}
	}
	return items
}

func readScript(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	return string(data)
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
				// Show preview for selected item when switching back to Scripts tab
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					if strings.HasSuffix(sel.name, ".sh") || strings.HasSuffix(sel.name, ".ps1") {
						ext := filepath.Ext(sel.path)
						m.vp.SetContent(highlightScript(readScript(sel.path), ext))
					} else {
						m.vp.SetContent("Select a script to preview...")
					}
				}
			} else {
				m.vp.SetContent("A cross-platform script browser powered by Bubble Tea.")
			}
		case "r":
			if m.activeTab == 0 && m.focus == focusList {
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					if !strings.HasSuffix(sel.name, "/") {
						m.vp.SetContent("Running script in a new terminal window...")
						go func() {
							var cmd *exec.Cmd
							if isWindows() {
								if strings.HasSuffix(sel.name, ".ps1") {
									// PowerShell: add Read-Host to pause
									cmd = exec.Command("cmd", "/C", "start", "powershell", "-NoExit", "-Command", sel.path+"; Write-Host ''; Read-Host 'Press Enter to exit'")
								} else {
									// Batch: add pause
									cmd = exec.Command("cmd", "/C", "start", "cmd", "/K", "sh "+sel.path+" & pause")
								}
							} else {
								// Try x-terminal-emulator, gnome-terminal, xterm
								term := "x-terminal-emulator"
								if _, err := exec.LookPath(term); err != nil {
									term = "gnome-terminal"
								}
								if _, err := exec.LookPath(term); err != nil {
									term = "xterm"
								}
								if strings.HasSuffix(sel.name, ".ps1") {
									cmd = exec.Command(term, "-e", "bash", "-c", "pwsh "+sel.path+"; echo; read -p 'Press Enter to exit'")
								} else {
									cmd = exec.Command(term, "-e", "bash", "-c", "sh "+sel.path+"; echo; read -p 'Press Enter to exit'")
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
						m.vp.SetContent(highlightScript(readScript(sel.path), ext))
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
									m.vp.SetContent(highlightScript(readScript(first.path), ext))
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
					if !strings.HasSuffix(sel.name, "/") {
						// Preview the script file
						return m, func() tea.Msg {
							ext := filepath.Ext(sel.path)
							return outputMsg(highlightScript(readScript(sel.path), ext))
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
							m.vp.SetContent(highlightScript(readScript(sel.path), ext))
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
							m.vp.SetContent(highlightScript(readScript(sel.path), ext))
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
        style := tabInactiveStyle.Copy()
        if i == m.activeTab {
            style = style.Inherit(tabActiveStyle)
        }
        tabLabels = append(tabLabels, style.Render(name))
    }
    tabBar := tabBarStyle.Render(strings.Join(tabLabels, "  "))

    centerStyle := lipgloss.NewStyle().Align(lipgloss.Center).Height(m.height-10)

    // Breadcrumb path above list
    breadcrumb := lipgloss.NewStyle().Faint(true).Render(m.currentPath)

    var body string
    if m.activeTab == 0 {
        // Divide width equally for 3 panels
        panelW := (m.width - 4) / 3 // -4 for border padding
        if panelW < 20 {
            panelW = 20
        }
        panelHeight := m.height - 10

        left := borderStyle.Width(panelW).Height(panelHeight).Render(
            breadcrumb + "\n" + centerStyle.Render(m.list.View()),
        )

        middle := borderStyle.Width(panelW).Height(panelHeight).Align(lipgloss.Center, lipgloss.Center).
            Render(
                lipgloss.NewStyle().
                    Foreground(lipgloss.Color("226")).
                    Bold(true).
                    Align(lipgloss.Center, lipgloss.Center).
                    Height(panelHeight).
                    Render("G\nO\n-\nP\nW\nR"),
            )

        right := borderStyle.Width(panelW).Height(panelHeight).Render(centerStyle.Render(m.vp.View()))

        body = lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)
    } else {
        body = borderStyle.Render(centerStyle.Render("A cross-platform script browser powered by Bubble Tea."))
    }

    footer := lipgloss.NewStyle().Foreground(pink).MarginTop(1).Align(lipgloss.Center).Render("← → or 🖱️ Click Tabs • ↑↓ Select • Enter Preview • r Run Script • q Quit")

    return lipgloss.JoinVertical(lipgloss.Left,
        tabBar,
        body,
        footer,
    )
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
	fmt.Fprint(w, style.Render(s.name)) // No .Width(w)
}

func (d scriptDelegate) Height() int               { return 1 }
func (d scriptDelegate) Spacing() int              { return 0 }
func (d scriptDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func isWindows() bool {
    return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows")
}

func main() {
	if err := ensureRepo(); err != nil {
		fmt.Println("Error cloning repo:", err)
		os.Exit(1)
	}

	tabs := []string{"Scripts", "About"}
	scriptItems := getScriptItems(filepath.Clean("scriptbin"))

	listModel := list.New(scriptItems, scriptDelegate{}, 0, 0)
	listModel.Title = "Scripts"
	listModel.SetShowHelp(false)
	listModel.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Select a script to preview...")

	// Initial preview for first script
	if len(scriptItems) > 0 {
		if s, ok := scriptItems[0].(scriptItem); ok {
			if strings.HasSuffix(s.name, ".sh") || strings.HasSuffix(s.name, ".ps1") {
				ext := filepath.Ext(s.path)
				vp.SetContent(highlightScript(readScript(s.path), ext))
			}
		}
	}

	m := model{
		list:        listModel,
		vp:          vp,
		tabs:        tabs,
		scriptItems: scriptItems,
		activeTab:   0,
		currentPath: filepath.Clean("scriptbin"),
		parentPaths: []parentNav{},
	}

	if err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

