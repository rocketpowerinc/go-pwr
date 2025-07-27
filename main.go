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

type model struct {
	list        list.Model
	vp          viewport.Model
	width       int
	height      int
	activeTab   int
	tabs        []string
	scriptItems []list.Item
	focus       focusArea
}

type outputMsg string

var tabColors = []lipgloss.TerminalColor{
	lipgloss.Color("27"), // Scripts - blue
	lipgloss.Color("40"), // About - green
}
var tabIcons = []string{"📜", "ℹ️"}

var (
	borderStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	tabActiveStyle   = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1)
	tabInactiveStyle = lipgloss.NewStyle().Faint(true).Padding(0, 1)
	tabBarStyle      = lipgloss.NewStyle().MarginBottom(1)
	headerStyle      = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	tabLabelStyle    = lipgloss.NewStyle().Bold(true).MarginBottom(1)
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
		fullPath := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			// Add folder as an item (user can traverse into it)
			items = append(items, scriptItem{name: entry.Name() + "/", path: fullPath})
		} else if strings.HasSuffix(entry.Name(), ".sh") || strings.HasSuffix(entry.Name(), ".ps1") {
			items = append(items, scriptItem{name: entry.Name(), path: fullPath})
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
			if m.focus == focusList {
				m.focus = focusPreview
			} else {
				m.focus = focusList
			}
		case "left":
			if m.activeTab > 0 {
				m.activeTab--
				if m.activeTab == 0 {
					m.list.SetItems(m.scriptItems)
				}
			}
		case "right":
			if m.activeTab < len(m.tabs)-1 {
				m.activeTab++
			}
		case "enter":
			if m.activeTab == 0 && m.focus == focusList {
				if sel, ok := m.list.SelectedItem().(scriptItem); ok {
					return m, func() tea.Msg {
						return outputMsg(readScript(sel.path))
					}
				}
			}
		case "up":
			if m.focus == focusList {
				m.list, cmd = m.list.Update(msg)
			} else if m.focus == focusPreview {
				m.vp.LineUp(1)
			}
		case "down":
			if m.focus == focusList {
				m.list, cmd = m.list.Update(msg)
			} else if m.focus == focusPreview {
				m.vp.LineDown(1)
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case outputMsg:
		m.vp.SetContent(strings.TrimSpace(string(msg)))
	}

	// Always update both panes so they stay in sync
	m.list, _ = m.list.Update(msg)
	m.vp, _ = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() string {
	tabColor := tabColors[m.activeTab]
	colorStyle := lipgloss.NewStyle().Foreground(tabColor)
	tabLabel := colorStyle.Copy().Bold(true).Render(fmt.Sprintf("%s %s", tabIcons[m.activeTab], m.tabs[m.activeTab]))

	var tabLabels []string
	for i, name := range m.tabs {
		style := tabInactiveStyle.Copy()
		if i == m.activeTab {
			style = style.Inherit(tabActiveStyle).Foreground(tabColors[i])
		}
		tabLabels = append(tabLabels, style.Render(name))
	}
	tabBar := tabBarStyle.Render(strings.Join(tabLabels, "  "))

	var body string
	if m.activeTab == 0 {
		left := borderStyle.Copy().BorderForeground(tabColor).Render(m.list.View())
		right := borderStyle.Copy().BorderForeground(tabColor).Render(m.vp.View())
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	} else {
		body = borderStyle.Copy().BorderForeground(tabColor).Render("A cross-platform script browser powered by Bubble Tea.")
	}

	footer := lipgloss.NewStyle().Foreground(tabColor).MarginTop(1).Render("← → or 🖱️ Click Tabs • ↑↓ Select • Enter Preview • q Quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Foreground(tabColor).Render("🧬 ScriptBin Browser"),
		tabBar,
		tabLabelStyle.Render(tabLabel),
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
	}

	if err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

