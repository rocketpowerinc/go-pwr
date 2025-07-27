package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type action struct {
	title, desc, script string
}

func (a action) Title() string       { return a.title }
func (a action) Description() string { return a.desc }
func (a action) FilterValue() string { return a.title }

type model struct {
	list        list.Model
	vp          viewport.Model
	width       int
	height      int
	message     string
	activeTab   int
	tabs        []string
	tabContents [][]list.Item
}

type outputMsg string
type tabClickMsg int

// Colors per tab
var tabColors = []lipgloss.TerminalColor{
	lipgloss.Color("27"), // Windows - blue
	lipgloss.Color("8"),  // Mac - gray
	lipgloss.Color("40"), // Linux - green
}

// Tab label emojis
var tabIcons = []string{"🪟", "🍎", "🍎"}
#🐧

var (
	borderStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	tabActiveStyle  = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1)
	tabInactiveStyle = lipgloss.NewStyle().Faint(true).Padding(0, 1)
	tabBarStyle     = lipgloss.NewStyle().MarginBottom(1)
	headerStyle     = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	tabLabelStyle   = lipgloss.NewStyle().Bold(true).MarginBottom(1)
)

func run(script string) string {
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("[error] %v\n%s", err, out)
	}
	return string(out)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		halfW := m.width / 2
		m.list.SetSize(halfW-4, m.height-10)
		m.vp.Width = halfW - 4
		m.vp.Height = m.height - 10
		return m, nil

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Click tab headers
			x := 1
			for i, tab := range m.tabs {
				end := x + len(tab) + 2
				if msg.X >= x && msg.X < end {
					m.activeTab = i
					m.list.SetItems(m.tabContents[i])
					break
				}
				x = end + 2
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.activeTab > 0 {
				m.activeTab--
				m.list.SetItems(m.tabContents[m.activeTab])
			}
		case "right":
			if m.activeTab < len(m.tabs)-1 {
				m.activeTab++
				m.list.SetItems(m.tabContents[m.activeTab])
			}
		case "enter":
			sel := m.list.SelectedItem().(action)
			return m, func() tea.Msg {
				return outputMsg(run(sel.script))
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case outputMsg:
		m.message = string(msg)
		m.vp.SetContent(strings.TrimSpace(m.message))
	}

	m.list, cmd = m.list.Update(msg)
	m.vp, _ = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() string {
	tabColor := tabColors[m.activeTab]
	colorStyle := lipgloss.NewStyle().Foreground(tabColor)

	// Active tab label
	tabLabel := colorStyle.Copy().Bold(true).Render(fmt.Sprintf("%s %s", tabIcons[m.activeTab], m.tabs[m.activeTab]))

	// Tab bar
	var tabLabels []string
	for i, name := range m.tabs {
		style := tabInactiveStyle.Copy()
		if i == m.activeTab {
			style = style.Inherit(tabActiveStyle).Foreground(tabColors[i])
		}
		tabLabels = append(tabLabels, style.Render(name))
	}
	tabBar := tabBarStyle.Render(strings.Join(tabLabels, "  "))

	left := borderStyle.Copy().BorderForeground(tabColor).Render(m.list.View())
	right := borderStyle.Copy().BorderForeground(tabColor).Render(m.vp.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := lipgloss.NewStyle().Foreground(tabColor).MarginTop(1).Render("← → or 🖱️ Click Tabs • ↑↓ Select • Enter Run • q Quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Foreground(tabColor).Render("🧬 Cross-Platform Greeter"),
		tabBar,
		tabLabelStyle.Render(tabLabel),
		body,
		footer,
	)
}

func main() {
	tabs := []string{"Windows", "Mac", "Linux"}

	windowsActions := []list.Item{
		action{"🔒 Lock", "Simulated Windows lock", `echo "Windows Lock"`},
		action{"🚪 Logout", "Simulated logout", `echo "Windows Logout"`},
	}

	macActions := []list.Item{
		action{"🔒 Lock", "Simulated Mac lock", `echo "Mac Lock"`},
		action{"🔁 Restart", "Simulated Mac restart", `echo "Restarting Mac"`},
	}

	linuxActions := []list.Item{
		action{"🔒 Lock", "Lock session", "loginctl lock-session"},
		action{"🔁 Logout", "Logout user", "loginctl kill-user $USER"},
		action{"⟳ Reboot", "Reboot system", "sudo reboot"},
		action{"⏻ Poweroff", "Shutdown", "sudo poweroff"},
		action{"👋 Greet", "Welcome message", `echo "Welcome, $(whoami)! 🐧"`},
	}

	tabContents := [][]list.Item{windowsActions, macActions, linuxActions}

	listModel := list.New(tabContents[0], list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Actions"
	listModel.SetShowHelp(false)
	listModel.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Choose an action...")

	m := model{
		list:        listModel,
		vp:          vp,
		tabs:        tabs,
		tabContents: tabContents,
		activeTab:   0,
	}

	if err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(), // Enables mouse support!
	).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

