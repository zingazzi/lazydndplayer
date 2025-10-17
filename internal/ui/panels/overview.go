// internal/ui/panels/overview.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// OverviewPanel displays character overview
type OverviewPanel struct {
	character *models.Character
}

// NewOverviewPanel creates a new overview panel
func NewOverviewPanel(char *models.Character) *OverviewPanel {
	return &OverviewPanel{character: char}
}

// View renders the overview panel
func (p *OverviewPanel) View(width, height int) string {
	char := p.character

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	var lines []string
	lines = append(lines, titleStyle.Render("CHARACTER OVERVIEW"))
	lines = append(lines, "")

	// Basic Info
	lines = append(lines, labelStyle.Render("Name:")+" "+valueStyle.Render(char.Name))
	lines = append(lines, labelStyle.Render("Race:")+" "+valueStyle.Render(char.Race))
	lines = append(lines, labelStyle.Render("Class:")+" "+valueStyle.Render(char.Class))
	lines = append(lines, labelStyle.Render("Background:")+" "+valueStyle.Render(char.Background))
	lines = append(lines, labelStyle.Render("Alignment:")+" "+valueStyle.Render(char.Alignment))
	lines = append(lines, "")

	// Level & XP
	xpCurrent := char.Experience
	xpNext := char.GetNextLevelXP()
	xpProgress := fmt.Sprintf("%d / %d", xpCurrent, xpNext)
	lines = append(lines, labelStyle.Render("Level:")+" "+valueStyle.Render(fmt.Sprintf("%d", char.Level)))
	lines = append(lines, labelStyle.Render("Experience:")+" "+valueStyle.Render(xpProgress))

	if char.CanLevelUp() {
		levelUpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
		lines = append(lines, levelUpStyle.Render(">>> Ready to Level Up! Press 'l' <<<"))
	}
	lines = append(lines, "")

	// Combat Stats
	hpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	if float64(char.CurrentHP)/float64(char.MaxHP) <= 0.5 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	}
	if float64(char.CurrentHP)/float64(char.MaxHP) <= 0.25 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	}

	lines = append(lines, labelStyle.Render("Hit Points:")+" "+hpStyle.Bold(true).Render(fmt.Sprintf("%d / %d", char.CurrentHP, char.MaxHP)))
	if char.TempHP > 0 {
		lines = append(lines, labelStyle.Render("Temp HP:")+" "+valueStyle.Render(fmt.Sprintf("%d", char.TempHP)))
	}
	lines = append(lines, labelStyle.Render("Armor Class:")+" "+valueStyle.Render(fmt.Sprintf("%d", char.ArmorClass)))
	lines = append(lines, labelStyle.Render("Speed:")+" "+valueStyle.Render(fmt.Sprintf("%d ft", char.Speed)))
	lines = append(lines, labelStyle.Render("Initiative:")+" "+valueStyle.Render(fmt.Sprintf("%+d", char.Initiative)))
	lines = append(lines, labelStyle.Render("Proficiency Bonus:")+" "+valueStyle.Render(fmt.Sprintf("+%d", char.ProficiencyBonus)))
	lines = append(lines, "")

	// Quick Actions
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'e' to edit character details"))

	content := strings.Join(lines, "\n")

	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// Update handles updates for the overview panel
func (p *OverviewPanel) Update(char *models.Character) {
	p.character = char
}
