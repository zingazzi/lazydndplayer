// internal/ui/panels/skills.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SkillsPanel displays character skills
type SkillsPanel struct {
	character     *models.Character
	selectedIndex int
	viewport      viewport.Model
	ready         bool
}

// NewSkillsPanel creates a new skills panel
func NewSkillsPanel(char *models.Character) *SkillsPanel {
	return &SkillsPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// View renders the skills panel
func (p *SkillsPanel) View(width, height int) string {
	char := p.character

	// Initialize viewport if not ready
	if !p.ready {
		p.viewport = viewport.New(width, height)
		p.ready = true
	}
	p.viewport.Width = width
	p.viewport.Height = height

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	var lines []string
	lines = append(lines, titleStyle.Render("SKILLS"))
	lines = append(lines, "")

	for i, skill := range char.Skills.List {
		abilityMod := char.AbilityScores.GetModifier(skill.Ability)
		totalBonus := skill.CalculateBonus(abilityMod, char.ProficiencyBonus)

		profMarker := "  "
		if skill.Proficiency == models.Proficient {
			profMarker = "● "
		} else if skill.Proficiency == models.Expertise {
			profMarker = "◆ "
		}

		line := fmt.Sprintf("%s%-20s (%s) %+3d",
			profMarker,
			skill.Name,
			skill.Ability,
			totalBonus,
		)

		if i == p.selectedIndex {
			lines = append(lines, selectedStyle.Render(line))
		} else {
			lines = append(lines, normalStyle.Render(line))
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("  = Not proficient"))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("● = Proficient"))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("◆ = Expertise (double proficiency)"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'r' to roll selected skill"))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'e' to toggle proficiency"))

	content := strings.Join(lines, "\n")
	p.viewport.SetContent(content)

	// Add scroll indicators at bottom of viewport
	scrollInfo := ""
	if p.viewport.ScrollPercent() < 1.0 {
		scrollInfo = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("↓ Scroll: %d%%", int(p.viewport.ScrollPercent()*100)))
	}

	return p.viewport.View() + scrollInfo
}

// Update handles updates for the skills panel
func (p *SkillsPanel) Update(char *models.Character) {
	p.character = char
}

// Next moves to next skill
func (p *SkillsPanel) Next() {
	if p.selectedIndex < len(p.character.Skills.List)-1 {
		p.selectedIndex++
		// Auto-scroll viewport to keep selection visible
		p.viewport.LineDown(1)
	}
}

// Prev moves to previous skill
func (p *SkillsPanel) Prev() {
	if p.selectedIndex > 0 {
		p.selectedIndex--
		// Auto-scroll viewport to keep selection visible
		p.viewport.LineUp(1)
	}
}

// GetSelectedSkill returns the currently selected skill
func (p *SkillsPanel) GetSelectedSkill() *models.Skill {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Skills.List) {
		return &p.character.Skills.List[p.selectedIndex]
	}
	return nil
}

// ToggleProficiency toggles proficiency for selected skill
func (p *SkillsPanel) ToggleProficiency() {
	if skill := p.GetSelectedSkill(); skill != nil {
		switch skill.Proficiency {
		case models.NotProficient:
			skill.Proficiency = models.Proficient
		case models.Proficient:
			skill.Proficiency = models.Expertise
		case models.Expertise:
			skill.Proficiency = models.NotProficient
		}
	}
}
