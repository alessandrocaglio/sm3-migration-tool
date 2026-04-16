package tui

import (
	"fmt"
	"strings"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	rhRed    = lipgloss.Color("#EE0000")
	darkGray = lipgloss.Color("#151515")
	white    = lipgloss.Color("#FFFFFF")

	titleStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(rhRed).
			Padding(0, 1).
			Bold(true)

	findingStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(rhRed).
			Padding(0, 1).
			MarginBottom(1)
)

type model struct {
	findings []checkers.Finding
	cursor   int
}

func NewModel(findings []checkers.Finding) model {
	return model{
		findings: findings,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.findings)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(" RED HAT SERVICE MESH 3.0 MIGRATION ASSISTANT ") + "\n\n")

	if len(m.findings) == 0 {
		s.WriteString("No findings yet. Run a scan first.\n")
	} else {
		for i, f := range m.findings {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			severityColor := lipgloss.Color("#AAAAAA")
			switch f.Severity {
			case checkers.SeverityHigh:
				severityColor = lipgloss.Color("#FF0000")
			case checkers.SeverityMedium:
				severityColor = lipgloss.Color("#FFA500")
			}

			sevStyle := lipgloss.NewStyle().Foreground(severityColor).Bold(true)

			content := fmt.Sprintf("%s [%s] %s\n%s\nRemediation: %s",
				sevStyle.Render(string(f.Severity)),
				f.Kind,
				f.ResourceName,
				f.Message,
				f.Remediation,
			)

			if m.cursor == i {
				s.WriteString(findingStyle.BorderForeground(rhRed).Render(cursor+" "+content) + "\n")
			} else {
				s.WriteString("  " + content + "\n\n")
			}
		}
	}

	s.WriteString("\nPress q to quit | ↑/↓ to navigate\n")
	return s.String()
}

func Run(findings []checkers.Finding) error {
	p := tea.NewProgram(NewModel(findings), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
