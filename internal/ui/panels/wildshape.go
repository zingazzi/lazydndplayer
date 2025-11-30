// internal/ui/panels/wildshape.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// WildShapeViewMode represents the current view mode
type WildShapeViewMode int

const (
	WildShapeViewList WildShapeViewMode = iota
	WildShapeViewDetail
)

// WildShapePanel displays wild shape forms and information
type WildShapePanel struct {
	character    *models.Character
	viewMode     WildShapeViewMode
	selectedIndex int // Selected form index in list view
	viewport     viewport.Model
	ready        bool
}

// NewWildShapePanel creates a new wild shape panel
func NewWildShapePanel(char *models.Character) *WildShapePanel {
	return &WildShapePanel{
		character:    char,
		viewMode:     WildShapeViewList,
		selectedIndex: 0,
	}
}

// GetSelectedForm returns the currently selected wild shape form
func (p *WildShapePanel) GetSelectedForm() string {
	forms := p.getKnownForms()
	if p.selectedIndex < 0 || p.selectedIndex >= len(forms) {
		return ""
	}
	return forms[p.selectedIndex]
}

// GetViewMode returns the current view mode
func (p *WildShapePanel) GetViewMode() WildShapeViewMode {
	return p.viewMode
}

// GetSelectedIndex returns the selected form index
func (p *WildShapePanel) GetSelectedIndex() int {
	return p.selectedIndex
}

// SetSelectedIndex sets the selected form index
func (p *WildShapePanel) SetSelectedIndex(index int) {
	forms := p.getKnownForms()
	if index >= 0 && index < len(forms) {
		p.selectedIndex = index
	}
}

// Next moves to next form in list view
func (p *WildShapePanel) Next() {
	if p.viewMode == WildShapeViewList {
		forms := p.getKnownForms()
		if len(forms) > 0 {
			p.selectedIndex = (p.selectedIndex + 1) % len(forms)
		}
	}
}

// Prev moves to previous form in list view
func (p *WildShapePanel) Prev() {
	if p.viewMode == WildShapeViewList {
		forms := p.getKnownForms()
		if len(forms) > 0 {
			p.selectedIndex--
			if p.selectedIndex < 0 {
				p.selectedIndex = len(forms) - 1
			}
		}
	}
}

// EnterDetailView switches to detail view for selected form
func (p *WildShapePanel) EnterDetailView() {
	forms := p.getKnownForms()
	if len(forms) > 0 && p.selectedIndex >= 0 && p.selectedIndex < len(forms) {
		p.viewMode = WildShapeViewDetail
	}
}

// ExitDetailView returns to list view
func (p *WildShapePanel) ExitDetailView() {
	p.viewMode = WildShapeViewList
}

// getKnownForms returns the list of known wild shape forms
func (p *WildShapePanel) getKnownForms() []string {
	// Get forms from character's wild shape data
	if len(p.character.WildShape.KnownForms) > 0 {
		return p.character.WildShape.KnownForms
	}

	// Fallback: return empty list if no forms stored
	return []string{}
}

// getWildShapeUses returns current and max wild shape uses
func (p *WildShapePanel) getWildShapeUses() (current int, max int) {
	// Find Wild Shape feature
	for _, feature := range p.character.Features.Features {
		if feature.Name == "Wild Shape" {
			return feature.CurrentUses, feature.MaxUses
		}
	}
	return 0, 0
}

// getMaxCR returns the maximum challenge rating for wild shape
func (p *WildShapePanel) getMaxCR() float64 {
	druidLevel := p.character.GetClassLevel("Druid")
	if druidLevel >= 8 {
		return 1.0
	} else if druidLevel >= 4 {
		return 0.5
	}
	return 0.25
}

// canFly returns whether the druid can transform into flying beasts
func (p *WildShapePanel) canFly() bool {
	druidLevel := p.character.GetClassLevel("Druid")
	return druidLevel >= 8
}

// View renders the wild shape panel
func (p *WildShapePanel) View(width, height int) string {
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

	druidLevel := p.character.GetClassLevel("Druid")
	if druidLevel < 2 {
		lines = append(lines, titleStyle.Render("WILD SHAPE"))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("Wild Shape is available at Druid level 2."))
	} else {
		if p.viewMode == WildShapeViewList {
			// LIST VIEW
			lines = append(lines, titleStyle.Render("WILD SHAPE"))
			lines = append(lines, "")

			// Wild Shape uses
			currentUses, maxUses := p.getWildShapeUses()
			usesStr := fmt.Sprintf("Uses: %d/%d", currentUses, maxUses)
			lines = append(lines, labelStyle.Render(usesStr))
			lines = append(lines, "")

			// CR limit and flying info
			maxCR := p.getMaxCR()
			crStr := fmt.Sprintf("Max CR: %.2f", maxCR)
			lines = append(lines, labelStyle.Render(crStr))
			if p.canFly() {
				lines = append(lines, valueStyle.Render("Can transform into flying beasts"))
			} else {
				lines = append(lines, dimStyle.Render("Cannot transform into flying beasts (available at level 8)"))
			}
			lines = append(lines, "")

			// Known forms
			forms := p.getKnownForms()
			knownFormsCount := len(forms)
			lines = append(lines, labelStyle.Render(fmt.Sprintf("KNOWN FORMS (%d):", knownFormsCount)))
			lines = append(lines, "")

			if len(forms) == 0 {
				lines = append(lines, dimStyle.Render("No forms known."))
			} else {
				// Ensure selectedIndex is valid
				if p.selectedIndex >= len(forms) {
					p.selectedIndex = len(forms) - 1
				}
				if p.selectedIndex < 0 {
					p.selectedIndex = 0
				}

				// Display form list
				for i, form := range forms {
					cursor := "  "
					style := valueStyle
					if i == p.selectedIndex {
						cursor = "❯ "
						style = selectedStyle
					}
					line := fmt.Sprintf("%s%s", cursor, style.Render(form))
					lines = append(lines, line)
				}
			}

			lines = append(lines, "")
			lines = append(lines, dimStyle.Render("Enter: View details • a: Manage forms • r: Long rest (regain uses)"))
		} else {
			// DETAIL VIEW
			formName := p.GetSelectedForm()
			if formName == "" {
				p.viewMode = WildShapeViewList
				return p.View(width, height)
			}

			lines = append(lines, titleStyle.Render(fmt.Sprintf("WILD SHAPE: %s", formName)))
			lines = append(lines, "")

			// Try to load beast data
			beasts, err := models.LoadBeastsFromJSON("data/beasts.json")
			if err == nil {
				// Find the beast
				var beast *models.BeastDefinition
				for i := range beasts {
					if beasts[i].Name == formName {
						beast = &beasts[i]
						break
					}
				}

				if beast != nil {
					// Initialize HP if not set (use beast's base HP)
					if p.character.WildShape.MaxHP == 0 {
						p.character.WildShape.MaxHP = beast.HP
						if p.character.WildShape.CurrentHP == 0 {
							p.character.WildShape.CurrentHP = beast.HP
						}
					}

					// Stats Section
					lines = append(lines, labelStyle.Render("STATS:"))
					lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("AC:"), valueStyle.Render(fmt.Sprintf("%d", beast.AC))))

					// Display HP (current/max) and temp HP if any
					hpStr := fmt.Sprintf("%d/%d", p.character.WildShape.CurrentHP, p.character.WildShape.MaxHP)
					if p.character.WildShape.TempHP > 0 {
						hpStr += fmt.Sprintf(" (+%d temp)", p.character.WildShape.TempHP)
					}
					lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("HP:"), valueStyle.Render(hpStr)))

					lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Speed:"), valueStyle.Render(beast.Speed)))
					lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render("Senses:"), valueStyle.Render(beast.Senses)))
					lines = append(lines, "")

					// Ability Scores
					lines = append(lines, labelStyle.Render("ABILITY SCORES:"))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("STR:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Strength)), beast.AbilityScores.GetModifier(models.Strength)))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("DEX:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Dexterity)), beast.AbilityScores.GetModifier(models.Dexterity)))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CON:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Constitution)), beast.AbilityScores.GetModifier(models.Constitution)))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("INT:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Intelligence)), beast.AbilityScores.GetModifier(models.Intelligence)))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("WIS:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Wisdom)), beast.AbilityScores.GetModifier(models.Wisdom)))
					lines = append(lines, fmt.Sprintf("%s %s (%+d)", labelStyle.Render("CHA:"), valueStyle.Render(fmt.Sprintf("%d", beast.AbilityScores.Charisma)), beast.AbilityScores.GetModifier(models.Charisma)))
					lines = append(lines, "")

					// Traits
					if len(beast.Traits) > 0 {
						lines = append(lines, labelStyle.Render("TRAITS:"))
						for _, trait := range beast.Traits {
							lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render("• "+trait)))
						}
						lines = append(lines, "")
					}

					// Actions Section
					if len(beast.Actions) > 0 {
						lines = append(lines, labelStyle.Render("ACTIONS:"))
						for _, action := range beast.Actions {
							attackBonusStr := fmt.Sprintf("%+d", action.AttackBonus)
							damageBonusStr := fmt.Sprintf("%+d", action.DamageBonus)
							lines = append(lines, fmt.Sprintf("  %s", valueStyle.Render(action.Name)))
							lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Attack:"), valueStyle.Render(attackBonusStr+" to hit")))
							lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Damage:"), valueStyle.Render(fmt.Sprintf("%s%s %s", action.DamageDice, damageBonusStr, action.DamageType))))
							lines = append(lines, fmt.Sprintf("    %s %s", dimStyle.Render("Range:"), valueStyle.Render(action.Range)))
							if action.Description != "" {
								lines = append(lines, fmt.Sprintf("    %s", dimStyle.Render(action.Description)))
							}
						}
						lines = append(lines, "")
					}
				} else {
					lines = append(lines, dimStyle.Render(fmt.Sprintf("Beast data not found for: %s", formName)))
				}
			} else {
				lines = append(lines, dimStyle.Render("Could not load beast data."))
			}

			// Help text - contextualized with form name
			lines = append(lines, dimStyle.Render(fmt.Sprintf("Esc: Back to list")))
			lines = append(lines, dimStyle.Render(fmt.Sprintf("HP: h (edit +/-) • + (add 1) • - (remove 1) for %s", formName)))
			lines = append(lines, dimStyle.Render(fmt.Sprintf("Temp HP: t (edit +/-) for %s", formName)))
		}
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
func (p *WildShapePanel) Update(char *models.Character) {
	p.character = char
	// Ensure selectedIndex is valid
	forms := p.getKnownForms()
	if p.selectedIndex >= len(forms) {
		p.selectedIndex = len(forms) - 1
		if p.selectedIndex < 0 {
			p.selectedIndex = 0
		}
	}
}

// UpdateViewport handles viewport updates
func (p *WildShapePanel) UpdateViewport(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

// ScrollDown scrolls the viewport down
func (p *WildShapePanel) ScrollDown() {
	p.viewport.LineDown(1)
}

// ScrollUp scrolls the viewport up
func (p *WildShapePanel) ScrollUp() {
	p.viewport.LineUp(1)
}

// PageDown scrolls down by half a page
func (p *WildShapePanel) PageDown() {
	p.viewport.HalfViewDown()
}

// PageUp scrolls up by half a page
func (p *WildShapePanel) PageUp() {
	p.viewport.HalfViewUp()
}
