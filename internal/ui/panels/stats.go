// internal/ui/panels/stats.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// StatsPanel displays ability scores and saving throws
type StatsPanel struct {
	character     *models.Character
	selectedIndex int
}

// NewStatsPanel creates a new stats panel
func NewStatsPanel(char *models.Character) *StatsPanel {
	return &StatsPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// View renders the stats panel
func (p *StatsPanel) View(width, height int) string {
	char := p.character

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Width(15)

	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(5)

	modPositiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	modNegativeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	var lines []string
	lines = append(lines, titleStyle.Render("ABILITY SCORES"))
	lines = append(lines, "")

	abilities := []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}

	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength",
		models.Dexterity:    "Dexterity",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom",
		models.Charisma:     "Charisma",
	}

	for _, ability := range abilities {
		score := char.AbilityScores.GetScore(ability)
		modifier := char.AbilityScores.GetModifier(ability)
		saveProficient := char.SavingThrows.IsProficient(ability)

		modStr := fmt.Sprintf("%+d", modifier)
		modStyle := modPositiveStyle
		if modifier < 0 {
			modStyle = modNegativeStyle
		}

		saveBonus := modifier
		if saveProficient {
			saveBonus += char.ProficiencyBonus
		}

		profMarker := " "
		if saveProficient {
			profMarker = "●"
		}

		line := fmt.Sprintf("%s %s %s (%s)  Save: %+d %s",
			labelStyle.Render(abilityNames[ability]),
			scoreStyle.Render(fmt.Sprintf("%d", score)),
			modStyle.Render(modStr),
			string(ability),
			saveBonus,
			profMarker,
		)

		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("● = Proficient in saving throw"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'e' to edit ability scores"))

	content := strings.Join(lines, "\n")

	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// Update handles updates for the stats panel
func (p *StatsPanel) Update(char *models.Character) {
	p.character = char
}
