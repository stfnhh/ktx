package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
██╗  ██╗████████╗██╗  ██╗
██║ ██╔╝╚══██╔══╝╚██╗██╔╝
█████╔╝    ██║    ╚███╔╝ 
██╔═██╗    ██║    ██╔██╗ 
██║  ██╗   ██║   ██╔╝ ██╗
╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
`

type model struct {
	contexts       []string
	filtered       []string
	cursor         int
	selected       string
	filter         string
	currentContext string
	quitting       bool
	width          int
	height         int
}

func loadContexts() ([]string, error) {
	out, err := exec.Command("kubectl", "config", "get-contexts", "-o", "name").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

func getCurrentContext() string {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func filterContexts(contexts []string, filter string) []string {
	if filter == "" {
		return contexts
	}

	var result []string
	filter = strings.ToLower(filter)

	for _, c := range contexts {
		if strings.Contains(strings.ToLower(c), filter) {
			result = append(result, c)
		}
	}
	return result
}

func initialModel() model {
	contexts, err := loadContexts()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	current := getCurrentContext()

	return model{
		contexts:       contexts,
		filtered:       contexts,
		currentContext: current,
		cursor:         0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.filtered = filterContexts(m.contexts, m.filter)
				m.cursor = 0
			}

		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor]
			}
			m.quitting = true
			return m, tea.Quit

		default:
			// typing filter
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.filtered = filterContexts(m.contexts, m.filter)
				m.cursor = 0
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var list strings.Builder

	for i, c := range m.filtered {
		prefix := "  "

		if c == m.currentContext {
			prefix = "● "
		}

		if i == m.cursor {
			prefix = "➜ "
		}

		list.WriteString(prefix + c + "\n")
	}

	header := lipgloss.NewStyle().
		Width(60).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#326CE5")).
		Bold(true).
		Render(logo)

	current := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("current: %s", m.currentContext))

	filterLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Render(fmt.Sprintf("Filter: %s", m.filter))

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", 60))

	content := lipgloss.NewStyle().
		Width(60).
		Render(
			header +
				divider + "\n" +
				current +
				"\n\n" +
				filterLine +
				"\n\n" +
				list.String() +
				"\n↑↓ navigate • type to filter • Enter select • q quit",
		)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	final := m.(model)

	if final.selected != "" {
		cmd := exec.Command("kubectl", "config", "use-context", final.selected)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}