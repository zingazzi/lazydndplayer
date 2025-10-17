// internal/ui/components/help.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpBinding represents a key binding with description
type HelpBinding struct {
	Key  string
	Desc string
}

// Help represents a help overlay
type Help struct {
	Visible  bool
	Bindings []HelpBinding
}

// NewHelp creates a new help overlay
func NewHelp() *Help {
	return &Help{
		Visible: false,
		Bindings: []HelpBinding{
			{"f", "Cycle focus (Main/Actions/Dice)"},
			{"↑/↓ or j/k", "Navigate lists"},
			{"Tab/Shift+Tab", "Switch tabs (when in Main)"},
			{"1-5", "Quick tab select (when in Main)"},
			{"1-7", "Quick dice (when in Dice focus)"},
			{"d + 1-7", "Quick dice (global)"},
			{"Enter", "Activate / Roll"},
			{"e", "Edit / Toggle"},
			{"a", "Add item/spell"},
			{"r", "Roll skill check"},
			{"Shift+R", "Long rest"},
			{"s", "Save"},
			{"?", "Help"},
			{"q", "Quit"},
		},
	}
}

// Toggle toggles the help overlay
func (h *Help) Toggle() {
	h.Visible = !h.Visible
}

// View renders the help overlay
func (h *Help) View(width, height int) string {
	if !h.Visible {
		return ""
	}

	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		Width(width - 4).
		Height(height - 4)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render("Keyboard Shortcuts")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	for _, binding := range h.Bindings {
		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Width(20)

		key := keyStyle.Render(binding.Key)
		desc := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(binding.Desc)
		lines = append(lines, key+" "+desc)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press ? or Esc to close"))

	content := strings.Join(lines, "\n")
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, helpStyle.Render(content))
}
