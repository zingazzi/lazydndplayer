// internal/ui/components/companionselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CompanionSelector handles companion selection (regular beasts and Beast Master special beasts)
type CompanionSelector struct {
	visible      bool
	character    *models.Character
	cursor       int
	selectedBeast string
	beastNames   []string // Regular beast names
	showBeastMaster bool  // Whether to show Beast Master special beasts
}

// NewCompanionSelector creates a new companion selector
func NewCompanionSelector(char *models.Character) *CompanionSelector {
	allBeasts := models.GetAllBeastDefinitions()
	beastNames := make([]string, 0, len(allBeasts))
	for _, beast := range allBeasts {
		beastNames = append(beastNames, beast.Name)
	}

	return &CompanionSelector{
		character:  char,
		beastNames: beastNames,
		cursor:     0,
	}
}

// Show displays the companion selector
func (cs *CompanionSelector) Show() {
	cs.visible = true
	cs.cursor = 0
	cs.selectedBeast = ""
	cs.showBeastMaster = cs.character.IsBeastMaster()
}

// Hide hides the companion selector
func (cs *CompanionSelector) Hide() {
	cs.visible = false
	cs.cursor = 0
	cs.selectedBeast = ""
}

// IsVisible returns whether the selector is visible
func (cs *CompanionSelector) IsVisible() bool {
	return cs.visible
}

// GetSelectedBeast returns the currently selected beast name
func (cs *CompanionSelector) GetSelectedBeast() string {
	allOptions := cs.getAvailableOptions()
	if cs.cursor >= 0 && cs.cursor < len(allOptions) {
		return allOptions[cs.cursor]
	}
	return ""
}

// getAvailableOptions returns all available companion options
func (cs *CompanionSelector) getAvailableOptions() []string {
	options := make([]string, 0)

	// Add regular beasts
	options = append(options, cs.beastNames...)

	// Add Beast Master special beasts if applicable
	if cs.showBeastMaster {
		options = append(options, string(models.CompanionTypeLand))
		options = append(options, string(models.CompanionTypeSky))
		options = append(options, string(models.CompanionTypeSea))
	}

	return options
}

// Next moves cursor down
func (cs *CompanionSelector) Next() {
	allOptions := cs.getAvailableOptions()
	if cs.cursor < len(allOptions)-1 {
		cs.cursor++
	}
}

// Prev moves cursor up
func (cs *CompanionSelector) Prev() {
	if cs.cursor > 0 {
		cs.cursor--
	}
}

// Select confirms the current selection
func (cs *CompanionSelector) Select() bool {
	allOptions := cs.getAvailableOptions()
	if cs.cursor >= 0 && cs.cursor < len(allOptions) {
		cs.selectedBeast = allOptions[cs.cursor]
		return true
	}
	return false
}

// Update handles key presses
func (cs *CompanionSelector) Update(msg tea.Msg) (CompanionSelector, tea.Cmd) {
	if !cs.visible {
		return *cs, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			cs.Prev()
		case "down", "j":
			cs.Next()
		case "enter":
			if cs.Select() {
				// Don't clear selectedBeast here - let the handler process it
				// The handler will call Hide() after adding the companion
				cs.visible = false
			}
		case "esc":
			cs.selectedBeast = ""
			cs.Hide()
		}
	}

	return *cs, nil
}

// View renders the companion selector
func (cs *CompanionSelector) View() string {
	if !cs.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	var content strings.Builder

	content.WriteString(titleStyle.Render("SELECT COMPANION") + "\n\n")

	// Two-column layout: list on left, details on right
	var leftContent strings.Builder
	var rightContent strings.Builder

	allOptions := cs.getAvailableOptions()

	// Build companion list (left side)
	for i, option := range allOptions {
		cursor := "  "
		style := normalStyle
		if i == cs.cursor {
			cursor = "❯ "
			style = selectedStyle
		}
		leftContent.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, option)) + "\n")
	}

	// Build details (right side)
	if cs.cursor >= 0 && cs.cursor < len(allOptions) {
		selectedOption := allOptions[cs.cursor]

		// Check if it's a Beast Master special beast
		if selectedOption == string(models.CompanionTypeLand) ||
			selectedOption == string(models.CompanionTypeSky) ||
			selectedOption == string(models.CompanionTypeSea) {
			// Show Beast Master special beast details
			companionType := models.CompanionType(selectedOption)
			rangerLevel := cs.character.GetClassLevel("Ranger")
			if rangerLevel == 0 {
				rangerLevel = 3 // Default for selection
			}
			companion := models.InitializeCompanion(companionType, rangerLevel)

			if companion != nil {
				rightContent.WriteString(statStyle.Render("BEAST MASTER SPECIAL BEAST") + "\n\n")
				rightContent.WriteString(fmt.Sprintf("AC: %d\n", companion.AC))
				rightContent.WriteString(fmt.Sprintf("HP: %d\n", companion.MaxHP))
				rightContent.WriteString(fmt.Sprintf("Speed: %s\n", companion.Speed))
				rightContent.WriteString(fmt.Sprintf("Senses: %s\n", companion.Senses))
				rightContent.WriteString("\n")

				// Beast Strike details
				rightContent.WriteString(statStyle.Render("BEAST STRIKE:") + "\n")
				mechanics := models.NewCompanionMechanics(cs.character)
				attack := mechanics.GetCompanionAttack(companion)
				if attack != nil {
					attackBonusStr := fmt.Sprintf("%+d", attack.AttackBonus)
					damageBonusStr := fmt.Sprintf("%+d", attack.DamageBonus)
					rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Attack: %s to hit", attackBonusStr)) + "\n")
					rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Damage: %s%s %s", attack.DamageDice, damageBonusStr, attack.DamageType)) + "\n")
					rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Range: %s", attack.Range)) + "\n")
					if companion.SpecialNotes != "" {
						rightContent.WriteString("\n")
						rightContent.WriteString(dimStyle.Render(companion.SpecialNotes) + "\n")
					}
				}
			}
		} else {
			// Show regular beast details
			beastDef := models.GetBeastDefinition(selectedOption)
			if beastDef != nil {
				rightContent.WriteString(statStyle.Render("REGULAR BEAST") + "\n\n")
				rightContent.WriteString(fmt.Sprintf("AC: %d\n", beastDef.AC))
				rightContent.WriteString(fmt.Sprintf("HP: %d\n", beastDef.HP))
				rightContent.WriteString(fmt.Sprintf("Speed: %s\n", beastDef.Speed))
				rightContent.WriteString(fmt.Sprintf("Senses: %s\n", beastDef.Senses))
				rightContent.WriteString("\n")

				// Ability Scores
				rightContent.WriteString(statStyle.Render("ABILITY SCORES:") + "\n")
				rightContent.WriteString(fmt.Sprintf("STR: %d (%+d)\n", beastDef.AbilityScores.Strength, beastDef.AbilityScores.GetModifier(models.Strength)))
				rightContent.WriteString(fmt.Sprintf("DEX: %d (%+d)\n", beastDef.AbilityScores.Dexterity, beastDef.AbilityScores.GetModifier(models.Dexterity)))
				rightContent.WriteString(fmt.Sprintf("CON: %d (%+d)\n", beastDef.AbilityScores.Constitution, beastDef.AbilityScores.GetModifier(models.Constitution)))
				rightContent.WriteString("\n")

				// Actions
				if len(beastDef.Actions) > 0 {
					rightContent.WriteString(statStyle.Render("ACTIONS:") + "\n")
					for _, action := range beastDef.Actions {
						attackBonusStr := fmt.Sprintf("%+d", action.AttackBonus)
						damageBonusStr := fmt.Sprintf("%+d", action.DamageBonus)
						rightContent.WriteString(normalStyle.Render(action.Name) + "\n")
						rightContent.WriteString(fmt.Sprintf("  Attack: %s to hit\n", attackBonusStr))
						rightContent.WriteString(fmt.Sprintf("  Damage: %s%s %s\n", action.DamageDice, damageBonusStr, action.DamageType))
						rightContent.WriteString(fmt.Sprintf("  Range: %s\n", action.Range))
						if action.Description != "" {
							rightContent.WriteString(fmt.Sprintf("  %s\n", dimStyle.Render(action.Description)))
						}
					}
				}
			}
		}
	}

	// Combine left and right content
	leftStr := leftContent.String()
	rightStr := rightContent.String()

	if rightStr != "" {
		content.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftStr, "  ", rightStr))
	} else {
		content.WriteString(leftStr)
	}

	content.WriteString("\n\n")
	content.WriteString(dimStyle.Render("Enter: Select • Esc: Cancel"))

	return content.String()
}
