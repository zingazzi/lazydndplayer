// internal/ui/panels/companion.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CompanionPanel displays companion stats, traits, and actions
type CompanionPanel struct {
	character    *models.Character
	selectedItem int // 0 = stats, 1 = actions
	viewport     viewport.Model
	ready        bool
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
	// Use all available height for the viewport
	viewportHeight := height

	if !p.ready {
		p.viewport = viewport.New(width, viewportHeight)
		p.viewport.Style = lipgloss.NewStyle()
		p.ready = true
	}

	if p.viewport.Width != width || p.viewport.Height != viewportHeight {
		p.viewport.Width = width
		p.viewport.Height = viewportHeight
	}

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

		contentStr := strings.Join(lines, "\n")
		p.viewport.SetContent(contentStr)
		return p.viewport.View()
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
	lines = append(lines, "")

	// Help text - show different message based on whether companion exists
	helpText := "c: Select beast • r: Roll attack • h: Edit HP"
	if companion != nil {
		helpText = "c: Change beast • r: Roll attack • h: Edit HP"
	}
	lines = append(lines, dimStyle.Render(helpText))

	contentStr := strings.Join(lines, "\n")
	p.viewport.SetContent(contentStr)

	// Render viewport
	viewportContent := p.viewport.View()

	// Overlay scroll indicator if content is scrollable
	if p.viewport.TotalLineCount() > p.viewport.Height {
		scrollPercentage := int(p.viewport.ScrollPercent() * 100)
		scrollInfo := fmt.Sprintf("[%d%%]", scrollPercentage)

		scrollStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Right)

		lines := strings.Split(viewportContent, "\n")
		if len(lines) > 0 {
			paddedLine := lipgloss.NewStyle().Width(width).Render(lines[len(lines)-1])
			lines[len(lines)-1] = lipgloss.JoinHorizontal(lipgloss.Top, paddedLine)
			lines = append(lines[:len(lines)-1],
				lipgloss.PlaceHorizontal(width, lipgloss.Right, scrollStyle.Render(scrollInfo)))
			viewportContent = strings.Join(lines, "\n")
		}
	}

	return viewportContent
}

// Update updates the panel with character data
func (p *CompanionPanel) Update(char *models.Character) {
	p.character = char
}

// Update handles viewport updates
func (p *CompanionPanel) UpdateViewport(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

// ScrollDown scrolls the viewport down
func (p *CompanionPanel) ScrollDown() {
	p.viewport.LineDown(1)
}

// ScrollUp scrolls the viewport up
func (p *CompanionPanel) ScrollUp() {
	p.viewport.LineUp(1)
}

// PageDown scrolls down by half a page
func (p *CompanionPanel) PageDown() {
	p.viewport.HalfViewDown()
}

// PageUp scrolls up by half a page
func (p *CompanionPanel) PageUp() {
	p.viewport.HalfViewUp()
}
