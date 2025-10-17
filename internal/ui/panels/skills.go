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

	// Build title
	var contentLines []string
	contentLines = append(contentLines, titleStyle.Render("SKILLS"))
	contentLines = append(contentLines, "")

	// Split skills into two columns
	numSkills := len(char.Skills.List)
	midpoint := (numSkills + 1) / 2

	columnWidth := width / 2
	if columnWidth > 35 {
		columnWidth = 35
	}

	// Build each row with two columns
	for i := 0; i < midpoint; i++ {
		leftSkill := char.Skills.List[i]
		leftAbilityMod := char.AbilityScores.GetModifier(leftSkill.Ability)
		leftBonus := leftSkill.CalculateBonus(leftAbilityMod, char.ProficiencyBonus)

		leftProfMarker := "  "
		if leftSkill.Proficiency == models.Proficient {
			leftProfMarker = "● "
		} else if leftSkill.Proficiency == models.Expertise {
			leftProfMarker = "◆ "
		}

		leftLine := fmt.Sprintf("%s%-15s (%s) %+3d",
			leftProfMarker,
			leftSkill.Name,
			leftSkill.Ability,
			leftBonus,
		)

		// Style left column
		var leftStyled string
		if i == p.selectedIndex {
			leftStyled = selectedStyle.Width(columnWidth).Render(leftLine)
		} else {
			leftStyled = normalStyle.Width(columnWidth).Render(leftLine)
		}

		// Right column (if exists)
		rightStyled := ""
		rightIdx := i + midpoint
		if rightIdx < numSkills {
			rightSkill := char.Skills.List[rightIdx]
			rightAbilityMod := char.AbilityScores.GetModifier(rightSkill.Ability)
			rightBonus := rightSkill.CalculateBonus(rightAbilityMod, char.ProficiencyBonus)

			rightProfMarker := "  "
			if rightSkill.Proficiency == models.Proficient {
				rightProfMarker = "● "
			} else if rightSkill.Proficiency == models.Expertise {
				rightProfMarker = "◆ "
			}

			rightLine := fmt.Sprintf("%s%-15s (%s) %+3d",
				rightProfMarker,
				rightSkill.Name,
				rightSkill.Ability,
				rightBonus,
			)

			if rightIdx == p.selectedIndex {
				rightStyled = selectedStyle.Width(columnWidth).Render(rightLine)
			} else {
				rightStyled = normalStyle.Width(columnWidth).Render(rightLine)
			}
		}

		// Join columns
		row := lipgloss.JoinHorizontal(lipgloss.Left, leftStyled, rightStyled)
		contentLines = append(contentLines, row)
	}

	contentLines = append(contentLines, "")
	contentLines = append(contentLines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("  = Not proficient  ● = Proficient  ◆ = Expertise"))
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'r' to roll selected skill | Press 'e' to toggle proficiency"))

	content := strings.Join(contentLines, "\n")
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
