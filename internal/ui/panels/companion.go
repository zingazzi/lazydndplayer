// internal/ui/panels/companion.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CompanionPanel displays companion stats, traits, and actions
type CompanionPanel struct {
	character    *models.Character
	selectedItem int // 0 = stats, 1 = actions
}

// NewCompanionPanel creates a new companion panel
func NewCompanionPanel(char *models.Character) *CompanionPanel {
	return &CompanionPanel{
		character:    char,
		selectedItem: 0,
	}
}

// Next moves to next item
func (p *CompanionPanel) Next() {
	p.selectedItem = (p.selectedItem + 1) % 2
}

// Prev moves to previous item
func (p *CompanionPanel) Prev() {
	p.selectedItem--
	if p.selectedItem < 0 {
		p.selectedItem = 1
	}
}

// GetSelectedAttack returns the selected attack (if on actions section)
func (p *CompanionPanel) GetSelectedAttack() *models.Attack {
	if p.character.Companion == nil {
		return nil
	}
	mechanics := models.NewCompanionMechanics(p.character)
	return mechanics.GetCompanionAttack()
}

// View renders the companion panel
func (p *CompanionPanel) View(width, height int) string {
	if p.character.Companion == nil {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 0, 1, 0)

		dimStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		var lines []string
		lines = append(lines, titleStyle.Render("COMPANION"))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("No companion selected."))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("Beast Master rangers can select a companion at level 3."))

		return strings.Join(lines, "\n")
	}

	companion := p.character.Companion
	mechanics := models.NewCompanionMechanics(p.character)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	var lines []string

	// Title
	lines = append(lines, titleStyle.Render(fmt.Sprintf("COMPANION: %s", string(companion.Type))))
	lines = append(lines, "")

	// Stats Section
	sectionStyle := labelStyle
	if p.selectedItem == 0 {
		sectionStyle = selectedStyle
	}
	lines = append(lines, sectionStyle.Render("STATS:"))
	lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("AC:"), valueStyle.Render(fmt.Sprintf("%d", companion.AC))))
	lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("HP:"), valueStyle.Render(fmt.Sprintf("%d/%d", companion.CurrentHP, companion.MaxHP))))
	lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Speed:"), valueStyle.Render(companion.Speed)))
	lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Senses:"), valueStyle.Render(companion.Senses)))
	lines = append(lines, "")

	// Ability Scores
	lines = append(lines, sectionStyle.Render("ABILITY SCORES:"))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("STR:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Strength)), companion.AbilityScores.GetModifier(models.Strength)))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("DEX:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Dexterity)), companion.AbilityScores.GetModifier(models.Dexterity)))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CON:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Constitution)), companion.AbilityScores.GetModifier(models.Constitution)))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("INT:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Intelligence)), companion.AbilityScores.GetModifier(models.Intelligence)))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("WIS:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Wisdom)), companion.AbilityScores.GetModifier(models.Wisdom)))
	lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CHA:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Charisma)), companion.AbilityScores.GetModifier(models.Charisma)))
	lines = append(lines, "")

	// Traits
	lines = append(lines, sectionStyle.Render("TRAITS:"))
	for _, trait := range companion.Traits {
		lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render("• "+trait)))
	}
	if companion.SpecialNotes != "" {
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("Note: "+companion.SpecialNotes))
	}
	lines = append(lines, "")

	// Actions Section
	sectionStyle = labelStyle
	if p.selectedItem == 1 {
		sectionStyle = selectedStyle
	}
	lines = append(lines, sectionStyle.Render("ACTIONS:"))
	attack := mechanics.GetCompanionAttack()
	if attack != nil {
		attackBonusStr := fmt.Sprintf("%+d", attack.AttackBonus)
		damageBonusStr := fmt.Sprintf("%+d", attack.DamageBonus)
		lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render(attack.Name)))
		lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Attack:"), valueStyle.Render(attackBonusStr+" to hit")))
		lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Damage:"), valueStyle.Render(fmt.Sprintf("%s%s %s", attack.DamageDice, damageBonusStr, attack.DamageType))))
		lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Range:"), valueStyle.Render(attack.Range)))
	}

	return strings.Join(lines, "\n")
}

// Update updates the panel with character data
func (p *CompanionPanel) Update(char *models.Character) {
	p.character = char
}
