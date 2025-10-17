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

// HelpSection represents a section of help bindings
type HelpSection struct {
	Title    string
	Bindings []HelpBinding
}

// Help represents a help overlay
type Help struct {
	Visible bool
}

// NewHelp creates a new help overlay
func NewHelp() *Help {
	return &Help{
		Visible: false,
	}
}

// Toggle toggles the help overlay
func (h *Help) Toggle() {
	h.Visible = !h.Visible
}

// GetGeneralBindings returns common keyboard shortcuts
func GetGeneralBindings() []HelpBinding {
	return []HelpBinding{
		{"f", "Cycle focus (Main → CharStats → Actions → Dice)"},
		{"Tab / Shift+Tab", "Switch tabs (when in Main panel)"},
		{"1-4", "Quick tab select (when in Main panel)"},
		{"↑/↓ or j/k", "Navigate lists"},
		{"s", "Save character"},
		{"?", "Toggle help"},
		{"q", "Quit application"},
	}
}

// GetOverviewBindings returns overview panel bindings
func GetOverviewBindings() []HelpBinding {
	return []HelpBinding{
		{"Shift+R", "Long rest (restore HP, spells, abilities)"},
	}
}

// GetStatsBindings returns stats panel bindings
func GetStatsBindings() []HelpBinding {
	return []HelpBinding{
		{"e", "Edit ability scores (not yet implemented)"},
		{"Shift+R", "Long rest"},
	}
}

// GetSkillsBindings returns skills panel bindings
func GetSkillsBindings() []HelpBinding {
	return []HelpBinding{
		{"↑/↓ or j/k", "Navigate skills"},
		{"r", "Roll selected skill check"},
		{"e", "Toggle proficiency (None → Proficient → Expertise)"},
		{"Shift+R", "Long rest"},
	}
}

// GetInventoryBindings returns inventory panel bindings
func GetInventoryBindings() []HelpBinding {
	return []HelpBinding{
		{"↑/↓ or j/k", "Navigate items"},
		{"a", "Add new item"},
		{"e", "Toggle equipped status"},
		{"d", "Delete selected item"},
		{"Shift+R", "Long rest"},
	}
}

// GetSpellsBindings returns spells panel bindings
func GetSpellsBindings() []HelpBinding {
	return []HelpBinding{
		{"a", "Add new spell"},
		{"r", "Restore spell slots (short rest)"},
		{"Shift+R", "Long rest (restore all slots)"},
	}
}

// GetCharacterStatsBindings returns character stats panel bindings
func GetCharacterStatsBindings() []HelpBinding {
	return []HelpBinding{
		{"n", "Edit character name"},
		{"r", "Select species (from D&D 5e 2024 species)"},
		{"h", "Adjust HP (popup)"},
		{"+/-", "Quick HP adjust (±1)"},
		{"i", "Roll initiative (1d20 + DEX)"},
		{"Shift+I", "Toggle Inspiration"},
	}
}

// GetTraitsBindings returns traits panel bindings
func GetTraitsBindings() []HelpBinding {
	return []HelpBinding{
		{"↑/↓ or j/k", "Navigate items"},
		{"Ctrl+D/U", "Page down/up"},
		{"Ctrl+E/Y", "Scroll down/up"},
		{"a", "Add language/feat"},
		{"d", "Delete selected"},
	}
}

// GetActionsBindings returns actions panel bindings (bottom panel)
func GetActionsBindings() []HelpBinding {
	return []HelpBinding{
		{"↑/↓ or j/k", "Navigate actions"},
		{"Enter", "Activate selected action"},
	}
}

// GetFeaturesBindings returns features panel bindings (main panel)
func GetFeaturesBindings() []HelpBinding {
	return []HelpBinding{
		{"↑/↓ or j/k", "Navigate features"},
		{"Ctrl+D/U", "Page down/up"},
		{"Ctrl+E/Y", "Scroll down/up"},
		{"u", "Use feature (consume charge)"},
		{"+/=", "Restore one use"},
		{"d", "Delete feature"},
		{"a", "Add feature"},
		{"r", "Short rest (recover short rest features)"},
		{"Shift+R", "Long rest (recover all features)"},
	}
}

// GetDiceBindings returns dice roller panel bindings
func GetDiceBindings(mode string) []HelpBinding {
	bindings := []HelpBinding{
		{"Enter", "Start typing dice notation"},
		{"h", "Browse roll history"},
		{"r", "Reroll last dice"},
	}

	switch mode {
	case "input":
		return []HelpBinding{
			{"Enter", "Roll the dice"},
			{"Esc", "Cancel and return to idle"},
			{"", ""},
			{"Examples:", ""},
			{"1d20", "Roll a d20"},
			{"2d6+3", "Roll 2d6 and add 3"},
			{"1d20 adv", "Roll with advantage"},
			{"1d20 dis", "Roll with disadvantage"},
			{"2d8+3d4+2", "Roll multiple dice types"},
			{"1d20+3, 2d6", "Multiple separate rolls"},
		}
	case "history":
		return []HelpBinding{
			{"↑/↓ or j/k", "Navigate roll history"},
			{"Enter", "Reroll selected dice"},
			{"Esc", "Return to idle"},
		}
	}

	return bindings
}

// ViewWithContext renders the help overlay with contextual bindings
func (h *Help) ViewWithContext(width, height int, panelName string, contextBindings []HelpBinding) string {
	if !h.Visible {
		return ""
	}

	// Create sections
	sections := []HelpSection{
		{
			Title:    "GENERAL",
			Bindings: GetGeneralBindings(),
		},
	}

	// Add context-specific section if provided
	if len(contextBindings) > 0 {
		sections = append(sections, HelpSection{
			Title:    strings.ToUpper(panelName),
			Bindings: contextBindings,
		})
	}

	// Render sections
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Underline(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Width(22)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	var lines []string
	lines = append(lines, titleStyle.Render("⌨  KEYBOARD SHORTCUTS"))
	lines = append(lines, "")

	for _, section := range sections {
		lines = append(lines, sectionTitleStyle.Render(section.Title))
		lines = append(lines, "")

		for _, binding := range section.Bindings {
			if binding.Key == "" && binding.Desc == "" {
				lines = append(lines, "")
				continue
			}
			if binding.Key == "" {
				// Section header within a section
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Italic(true).
					Render("  "+binding.Desc))
				continue
			}
			key := keyStyle.Render("  " + binding.Key)
			desc := descStyle.Render(binding.Desc)
			lines = append(lines, key+" "+desc)
		}
		lines = append(lines, "")
	}

	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Render("Press ? or Esc to close"))

	content := strings.Join(lines, "\n")

	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		MaxWidth(width - 8).
		MaxHeight(height - 4)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, helpStyle.Render(content))
}
