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
	character      *models.Character
	selectedAbility int // 0-5 for the selected ability
}

// NewStatsPanel creates a new stats panel
func NewStatsPanel(char *models.Character) *StatsPanel {
	return &StatsPanel{
		character:      char,
		selectedAbility: 0,
	}
}

// Next moves to next ability
func (p *StatsPanel) Next() {
	p.selectedAbility = (p.selectedAbility + 1) % 6
}

// Prev moves to previous ability
func (p *StatsPanel) Prev() {
	p.selectedAbility--
	if p.selectedAbility < 0 {
		p.selectedAbility = 5
	}
}

// GetSelectedAbility returns the currently selected ability
func (p *StatsPanel) GetSelectedAbility() models.AbilityType {
	abilities := []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}
	return abilities[p.selectedAbility]
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

	modPositiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	modNegativeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	var lines []string
	lines = append(lines, titleStyle.Render("ABILITY SCORES"))
	lines = append(lines, "")

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)
	lines = append(lines, headerStyle.Render(fmt.Sprintf("%-14s %5s  %5s  %6s", "Ability", "Score", "Mod", "Save")))
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

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237")).
		Bold(true)

	for i, ability := range abilities {
		score := char.AbilityScores.GetScore(ability)
		modifier := char.AbilityScores.GetModifier(ability)

		// Check if proficient in this saving throw (from class)
		isProficient := false
		abilityName := abilityNames[ability] // Use full name (e.g., "Dexterity")
		for _, prof := range char.SavingThrowProficiencies {
			if strings.EqualFold(prof, abilityName) {
				isProficient = true
				break
			}
		}

		saveBonus := modifier
		if isProficient {
			saveBonus += char.ProficiencyBonus
		}

		modStr := fmt.Sprintf("%+d", modifier)
		modStyle := modPositiveStyle
		if modifier < 0 {
			modStyle = modNegativeStyle
		}

		saveStr := fmt.Sprintf("%+d", saveBonus)
		saveStyle := modPositiveStyle
		if saveBonus < 0 {
			saveStyle = modNegativeStyle
		}

		// Add proficiency indicator
		if isProficient {
			saveStr = "⦿" + saveStr
		} else {
			saveStr = " " + saveStr
		}

		line := fmt.Sprintf("%-14s %5d  %s  %s",
			labelStyle.Render(abilityNames[ability]),
			score,
			modStyle.Render(fmt.Sprintf("%5s", modStr)),
			saveStyle.Render(fmt.Sprintf("%6s", saveStr)),
		)

		// Highlight selected ability
		if i == p.selectedAbility {
			lines = append(lines, selectedStyle.Render("▶ "+line))
		} else {
			lines = append(lines, "  "+line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("⦿ = Proficient in saving throw"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/↓: Select  •  'r' Roll stats  •  'e' Edit modifiers  •  't' Saving throw  •  'a' Ability check"))

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
