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

// CompanionViewMode represents the current view mode
type CompanionViewMode int

const (
	CompanionViewList CompanionViewMode = iota
	CompanionViewDetail
)

// CompanionPanel displays companion list and details
type CompanionPanel struct {
	character    *models.Character
	viewMode     CompanionViewMode
	selectedIndex int // Selected companion index in list view
	viewport     viewport.Model
	ready        bool
}

// NewCompanionPanel creates a new companion panel
func NewCompanionPanel(char *models.Character) *CompanionPanel {
	return &CompanionPanel{
		character:    char,
		viewMode:     CompanionViewList,
		selectedIndex: 0,
	}
}

// GetSelectedCompanion returns the currently selected companion
func (p *CompanionPanel) GetSelectedCompanion() *models.Companion {
	if p.selectedIndex < 0 || p.selectedIndex >= len(p.character.Companions) {
		return nil
	}
	return &p.character.Companions[p.selectedIndex]
}

// GetViewMode returns the current view mode
func (p *CompanionPanel) GetViewMode() CompanionViewMode {
	return p.viewMode
}

// GetSelectedIndex returns the selected companion index
func (p *CompanionPanel) GetSelectedIndex() int {
	return p.selectedIndex
}

// SetSelectedIndex sets the selected companion index
func (p *CompanionPanel) SetSelectedIndex(index int) {
	if index >= 0 && index < len(p.character.Companions) {
		p.selectedIndex = index
	}
}

// GetSelectedAttack returns the selected companion's attack
func (p *CompanionPanel) GetSelectedAttack() *models.Attack {
	companion := p.GetSelectedCompanion()
	if companion == nil {
		return nil
	}
	mechanics := models.NewCompanionMechanics(p.character)
	return mechanics.GetCompanionAttack(companion)
}

// Next moves to next companion in list view
func (p *CompanionPanel) Next() {
	if p.viewMode == CompanionViewList {
		if len(p.character.Companions) > 0 {
			p.selectedIndex = (p.selectedIndex + 1) % len(p.character.Companions)
		}
	}
}

// Prev moves to previous companion in list view
func (p *CompanionPanel) Prev() {
	if p.viewMode == CompanionViewList {
		if len(p.character.Companions) > 0 {
			p.selectedIndex--
			if p.selectedIndex < 0 {
				p.selectedIndex = len(p.character.Companions) - 1
			}
		}
	}
}

// EnterDetailView switches to detail view for selected companion
func (p *CompanionPanel) EnterDetailView() {
	if len(p.character.Companions) > 0 && p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Companions) {
		p.viewMode = CompanionViewDetail
	}
}

// ExitDetailView returns to list view
func (p *CompanionPanel) ExitDetailView() {
	p.viewMode = CompanionViewList
}

// View renders the companion panel
func (p *CompanionPanel) View(width, height int) string {
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

	if p.viewMode == CompanionViewList {
		// LIST VIEW
		lines = append(lines, titleStyle.Render("COMPANIONS"))
		lines = append(lines, "")

		if len(p.character.Companions) == 0 {
			lines = append(lines, dimStyle.Render("No companions."))
			lines = append(lines, "")
			lines = append(lines, dimStyle.Render("Press 'a' to add a companion."))
		} else {
			// Ensure selectedIndex is valid
			if p.selectedIndex >= len(p.character.Companions) {
				p.selectedIndex = len(p.character.Companions) - 1
			}
			if p.selectedIndex < 0 {
				p.selectedIndex = 0
			}

			// Display companion list
			for i, companion := range p.character.Companions {
				cursor := "  "
				style := valueStyle
				if i == p.selectedIndex {
					cursor = "❯ "
					style = selectedStyle
				}

				displayName := companion.GetDisplayName()
				// Add indicator for Beast Master special beasts
				if companion.IsBeastMasterSpecial() {
					displayName += " (Beast Master)"
				}

				hpStr := fmt.Sprintf("HP: %d/%d", companion.CurrentHP, companion.MaxHP)
				line := fmt.Sprintf("%s%s - %s", cursor, style.Render(displayName), valueStyle.Render(hpStr))
				lines = append(lines, line)
			}

			lines = append(lines, "")
			lines = append(lines, dimStyle.Render("Enter: View details • a: Add • d: Remove • +/-: HP • r: Roll attack"))
		}
	} else {
		// DETAIL VIEW
		companion := p.GetSelectedCompanion()
		if companion == nil {
			p.viewMode = CompanionViewList
			return p.View(width, height)
		}

		mechanics := models.NewCompanionMechanics(p.character)

		displayName := companion.GetDisplayName()
		if companion.IsBeastMasterSpecial() {
			displayName += " (Beast Master)"
		}

		lines = append(lines, titleStyle.Render(fmt.Sprintf("COMPANION: %s", displayName)))
		lines = append(lines, "")

		// Stats Section
		lines = append(lines, labelStyle.Render("STATS:"))
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("AC:"), valueStyle.Render(fmt.Sprintf("%d", companion.AC))))
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("HP:"), valueStyle.Render(fmt.Sprintf("%d/%d", companion.CurrentHP, companion.MaxHP))))
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Speed:"), valueStyle.Render(companion.Speed)))
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Senses:"), valueStyle.Render(companion.Senses)))
		lines = append(lines, "")

		// Ability Scores
		lines = append(lines, labelStyle.Render("ABILITY SCORES:"))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("STR:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Strength)), companion.AbilityScores.GetModifier(models.Strength)))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("DEX:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Dexterity)), companion.AbilityScores.GetModifier(models.Dexterity)))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CON:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Constitution)), companion.AbilityScores.GetModifier(models.Constitution)))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("INT:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Intelligence)), companion.AbilityScores.GetModifier(models.Intelligence)))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("WIS:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Wisdom)), companion.AbilityScores.GetModifier(models.Wisdom)))
		lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CHA:"), valueStyle.Render(fmt.Sprintf("%d", companion.AbilityScores.Charisma)), companion.AbilityScores.GetModifier(models.Charisma)))
		lines = append(lines, "")

		// Traits
		if len(companion.Traits) > 0 {
			lines = append(lines, labelStyle.Render("TRAITS:"))
			for _, trait := range companion.Traits {
				lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render("• "+trait)))
			}
			lines = append(lines, "")
		}

		if companion.SpecialNotes != "" {
			lines = append(lines, dimStyle.Render("Note: "+companion.SpecialNotes))
			lines = append(lines, "")
		}

		// Actions Section
		lines = append(lines, labelStyle.Render("ACTIONS:"))
		attack := mechanics.GetCompanionAttack(companion)
		if attack != nil {
			attackBonusStr := fmt.Sprintf("%+d", attack.AttackBonus)
			damageBonusStr := fmt.Sprintf("%+d", attack.DamageBonus)
			lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render(attack.Name)))
			lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Attack:"), valueStyle.Render(attackBonusStr+" to hit")))
			lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Damage:"), valueStyle.Render(fmt.Sprintf("%s%s %s", attack.DamageDice, damageBonusStr, attack.DamageType))))
			lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Range:"), valueStyle.Render(attack.Range)))
			if len(attack.Properties) > 0 {
				for _, prop := range attack.Properties {
					lines = append(lines, fmt.Sprintf("    %s", dimStyle.Render(prop)))
				}
			}
		}
		lines = append(lines, "")

		// Help text
		lines = append(lines, dimStyle.Render("Esc: Back to list • n: Rename • h: Edit HP • r: Roll attack"))
	}

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
	// Ensure selectedIndex is valid
	if p.selectedIndex >= len(p.character.Companions) {
		p.selectedIndex = len(p.character.Companions) - 1
		if p.selectedIndex < 0 {
			p.selectedIndex = 0
		}
	}
}

// UpdateViewport handles viewport updates
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
