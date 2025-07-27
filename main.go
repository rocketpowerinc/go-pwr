package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

var pink = lipgloss.Color("205") // ANSI pink, or use "#ff69b4" for hex

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

func ensureRepo() error {
	if _, err := os.Stat("scriptbin"); os.IsNotExist(err) {
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
	for _, entry := range entries {
		// Skip hidden files/folders
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
	// Always split the terminal in half, minimum width 20 per pane
	halfW := m.width / 2
	if halfW < 20 {
		halfW = 20
	}
	m.list.SetSize(halfW-4, m.height-10)
	m.vp.Width = halfW - 4
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
			// Only allow tab key to switch tabs
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			if m.activeTab == 0 {
				m.list.SetItems(m.scriptItems)
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
				m.vp.SetContent("Select a script to preview...")
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
									m.vp.SetContent(readScript(first.path))
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
							return outputMsg(readScript(sel.path))
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
							m.vp.SetContent(readScript(sel.path))
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
							m.vp.SetContent(readScript(sel.path))
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

	var body string
	if m.activeTab == 0 {
		left := borderStyle.Render(m.list.View())
		right := borderStyle.Render(m.vp.View())
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	} else {
		body = borderStyle.Render("A cross-platform script browser powered by Bubble Tea.")
	}

	footer := lipgloss.NewStyle().Foreground(pink).MarginTop(1).Render("← → or 🖱️ Click Tabs • ↑↓ Select • Enter Preview • q Quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("🧬 ScriptBin Browser"),
		tabBar,
		body,
		footer,
	)
}

func main() {
	if err := ensureRepo(); err != nil {
		fmt.Println("Error cloning repo:", err)
		os.Exit(1)
	}

	tabs := []string{"Scripts", "About"}
	scriptItems := getScriptItems("scriptbin")

	listModel := list.New(scriptItems, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Scripts"
	listModel.SetShowHelp(false)
	listModel.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Select a script to preview...")

	m := model{
		list:        listModel,
		vp:          vp,
		tabs:        tabs,
		scriptItems: scriptItems,
		activeTab:   0,
		currentPath: "scriptbin",
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

