// internal/ui/components/wildshapeformselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// WildShapeFormSelector handles wild shape form selection for druids
type WildShapeFormSelector struct {
	visible      bool
	character    *models.Character
	cursor       int
	selectedForm string
	beastNames   []string // Available beast names (filtered by CR/flying restrictions)
}

// NewWildShapeFormSelector creates a new wild shape form selector
func NewWildShapeFormSelector(char *models.Character) *WildShapeFormSelector {
	return &WildShapeFormSelector{
		character:  char,
		beastNames: []string{},
		cursor:     0,
	}
}

// Show displays the wild shape form selector with filtered beasts
func (wsfs *WildShapeFormSelector) Show() {
	wsfs.visible = true
	wsfs.cursor = 0
	wsfs.selectedForm = ""

	// Get druid level
	druidLevel := wsfs.character.GetClassLevel("Druid")
	if druidLevel < 2 {
		wsfs.beastNames = []string{}
		return
	}

	// Get available wild shape forms based on druid level
	availableBeasts := models.GetWildShapeForms(druidLevel)

	// Get currently known forms to exclude them from the list
	knownForms := make(map[string]bool)
	if len(wsfs.character.WildShape.KnownForms) > 0 {
		for _, form := range wsfs.character.WildShape.KnownForms {
			knownForms[form] = true
		}
	}

	// Build list of available beasts (excluding already known forms)
	wsfs.beastNames = make([]string, 0, len(availableBeasts))
	for _, beast := range availableBeasts {
		if !knownForms[beast.Name] {
			wsfs.beastNames = append(wsfs.beastNames, beast.Name)
		}
	}
}

// Hide hides the wild shape form selector
func (wsfs *WildShapeFormSelector) Hide() {
	wsfs.visible = false
	wsfs.cursor = 0
	wsfs.selectedForm = ""
	wsfs.beastNames = []string{}
}

// IsVisible returns whether the selector is visible
func (wsfs *WildShapeFormSelector) IsVisible() bool {
	return wsfs.visible
}

// GetSelectedForm returns the currently selected beast name
func (wsfs *WildShapeFormSelector) GetSelectedForm() string {
	if wsfs.cursor >= 0 && wsfs.cursor < len(wsfs.beastNames) {
		return wsfs.beastNames[wsfs.cursor]
	}
	return ""
}

// Next moves cursor down
func (wsfs *WildShapeFormSelector) Next() {
	if wsfs.cursor < len(wsfs.beastNames)-1 {
		wsfs.cursor++
	}
}

// Prev moves cursor up
func (wsfs *WildShapeFormSelector) Prev() {
	if wsfs.cursor > 0 {
		wsfs.cursor--
	}
}

// Select confirms the current selection
func (wsfs *WildShapeFormSelector) Select() bool {
	if wsfs.cursor >= 0 && wsfs.cursor < len(wsfs.beastNames) {
		wsfs.selectedForm = wsfs.beastNames[wsfs.cursor]
		return true
	}
	return false
}

// Update handles key presses
func (wsfs *WildShapeFormSelector) Update(msg tea.Msg) (WildShapeFormSelector, tea.Cmd) {
	if !wsfs.visible {
		return *wsfs, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			wsfs.Prev()
		case "down", "j":
			wsfs.Next()
		case "enter":
			if wsfs.Select() {
				// Don't clear selectedForm here - let the handler process it
				// The handler will call Hide() after adding the form
				wsfs.visible = false
			}
		case "esc":
			wsfs.selectedForm = ""
			wsfs.Hide()
		}
	}

	return *wsfs, nil
}

// View renders the wild shape form selector as a popup
func (wsfs *WildShapeFormSelector) View(width, height int) string {
	if !wsfs.visible {
		return ""
	}

	// Popup styling
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(width - 4).
		Background(lipgloss.Color("235"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center).
		Width(width - 8).
		Padding(0, 0, 1, 0)

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

	// Get druid level for restrictions display
	druidLevel := wsfs.character.GetClassLevel("Druid")
	maxCR := 0.25
	canFly := false
	if druidLevel >= 8 {
		maxCR = 1.0
		canFly = true
	} else if druidLevel >= 4 {
		maxCR = 0.5
	}

	title := fmt.Sprintf("SELECT WILD SHAPE FORM (Max CR: %.2f", maxCR)
	if canFly {
		title += ", Flying allowed"
	}
	title += ")"

	content.WriteString(titleStyle.Render(title) + "\n\n")

	// Two-column layout: list on left, details on right
	var leftContent strings.Builder
	var rightContent strings.Builder

	// Build form list (left side)
	if len(wsfs.beastNames) == 0 {
		leftContent.WriteString(dimStyle.Render("No available forms (all forms already known)") + "\n")
	} else {
		for i, formName := range wsfs.beastNames {
			cursor := "  "
			style := normalStyle
			if i == wsfs.cursor {
				cursor = "❯ "
				style = selectedStyle
			}
			leftContent.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, formName)) + "\n")
		}
	}

	// Build details (right side)
	if wsfs.cursor >= 0 && wsfs.cursor < len(wsfs.beastNames) {
		selectedForm := wsfs.beastNames[wsfs.cursor]
		beastDef := models.GetBeastDefinition(selectedForm)

		if beastDef != nil {
			rightContent.WriteString(statStyle.Render("BEAST STATS") + "\n\n")
			rightContent.WriteString(fmt.Sprintf("CR: %.2f\n", beastDef.CR))
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

	// Render as popup with border
	popupContent := popupStyle.Render(content.String())

	// Center the popup on screen
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popupContent)
}
