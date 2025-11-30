// internal/ui/components/wildshapeformselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// WildShapeFormSelector handles wild shape form selection for druids (multi-select, similar to cantrip selector)
type WildShapeFormSelector struct {
	visible         bool
	availableForms  []models.BeastDefinition
	selectedForms   []string // Names of selected forms
	cursor          int
	maxForms        int
	druidLevel      int
	character       *models.Character
}

// NewWildShapeFormSelector creates a new wild shape form selector
func NewWildShapeFormSelector(char *models.Character) *WildShapeFormSelector {
	return &WildShapeFormSelector{
		visible:        false,
		availableForms: []models.BeastDefinition{},
		selectedForms:  []string{},
		cursor:         0,
		character:      char,
	}
}

// Show displays the wild shape form selector with all available beasts
func (wsfs *WildShapeFormSelector) Show() {
	wsfs.visible = true
	wsfs.cursor = 0

	// Get druid level
	druidLevel := wsfs.character.GetClassLevel("Druid")
	wsfs.druidLevel = druidLevel
	if druidLevel < 2 {
		wsfs.availableForms = []models.BeastDefinition{}
		wsfs.maxForms = 0
		return
	}

	// Determine max forms based on druid level
	if druidLevel >= 8 {
		wsfs.maxForms = 8
	} else if druidLevel >= 4 {
		wsfs.maxForms = 6
	} else if druidLevel >= 2 {
		wsfs.maxForms = 4
	} else {
		wsfs.maxForms = 4
	}

	// Pre-select existing forms from character
	wsfs.selectedForms = make([]string, len(wsfs.character.WildShape.KnownForms))
	copy(wsfs.selectedForms, wsfs.character.WildShape.KnownForms)

	// Load available forms based on druid level
	wsfs.loadForms()
}

// Hide hides the wild shape form selector
func (wsfs *WildShapeFormSelector) Hide() {
	wsfs.visible = false
	wsfs.cursor = 0
	wsfs.selectedForms = []string{}
	wsfs.availableForms = []models.BeastDefinition{}
}

// IsVisible returns whether the selector is visible
func (wsfs *WildShapeFormSelector) IsVisible() bool {
	return wsfs.visible
}

// loadForms loads available wild shape forms based on druid level
func (wsfs *WildShapeFormSelector) loadForms() {
	// Get available wild shape forms based on druid level
	wsfs.availableForms = models.GetWildShapeForms(wsfs.druidLevel)
}

// Next moves cursor down
func (wsfs *WildShapeFormSelector) Next() {
	if wsfs.cursor < len(wsfs.availableForms)-1 {
		wsfs.cursor++
	}
}

// Prev moves cursor up
func (wsfs *WildShapeFormSelector) Prev() {
	if wsfs.cursor > 0 {
		wsfs.cursor--
	}
}

// ToggleSelection toggles the selection of the form at the cursor
func (wsfs *WildShapeFormSelector) ToggleSelection() {
	if len(wsfs.availableForms) == 0 {
		return
	}

	selected := wsfs.availableForms[wsfs.cursor]

	// Check if already selected
	for i, name := range wsfs.selectedForms {
		if name == selected.Name {
			// Deselect
			wsfs.selectedForms = append(wsfs.selectedForms[:i], wsfs.selectedForms[i+1:]...)
			return
		}
	}

	// Select if not at max
	if len(wsfs.selectedForms) < wsfs.maxForms {
		wsfs.selectedForms = append(wsfs.selectedForms, selected.Name)
	}
}

// IsSelected checks if a form is selected
func (wsfs *WildShapeFormSelector) IsSelected(formName string) bool {
	for _, name := range wsfs.selectedForms {
		if name == formName {
			return true
		}
	}
	return false
}

// CanConfirm returns true if at least 1 form is selected (up to maxForms)
func (wsfs *WildShapeFormSelector) CanConfirm() bool {
	return len(wsfs.selectedForms) > 0 && len(wsfs.selectedForms) <= wsfs.maxForms
}

// GetSelectedForms returns the list of selected form names
func (wsfs *WildShapeFormSelector) GetSelectedForms() []string {
	return wsfs.selectedForms
}

// GetMaxForms returns the maximum number of forms allowed
func (wsfs *WildShapeFormSelector) GetMaxForms() int {
	return wsfs.maxForms
}

// GetSelectedCount returns the number of selected forms
func (wsfs *WildShapeFormSelector) GetSelectedCount() int {
	return len(wsfs.selectedForms)
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
		case " ":
			wsfs.ToggleSelection()
		case "enter":
			if wsfs.CanConfirm() {
				// Confirmed - handler will process
				return *wsfs, nil
			}
		case "esc":
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

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	// Get restrictions for display
	maxCR := 0.25
	canFly := false
	if wsfs.druidLevel >= 8 {
		maxCR = 1.0
		canFly = true
	} else if wsfs.druidLevel >= 4 {
		maxCR = 0.5
	}

	// Left column - form list
	var leftLines []string
	titleText := fmt.Sprintf("Select 1-%d wild shape forms  (%d/%d)", wsfs.maxForms, len(wsfs.selectedForms), wsfs.maxForms)
	crText := fmt.Sprintf("Max CR: %.2f", maxCR)
	if canFly {
		crText += ", Flying allowed"
	}
	leftLines = append(leftLines, headerStyle.Render(titleText))
	leftLines = append(leftLines, dimStyle.Render(crText))
	leftLines = append(leftLines, "")

	// Show forms
	for i, form := range wsfs.availableForms {
		cursor := "  "
		if i == wsfs.cursor {
			cursor = "❯ "
		}

		checkbox := "[ ]"
		style := normalStyle
		if wsfs.IsSelected(form.Name) {
			checkbox = "[✓]"
			style = selectedStyle
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, form.Name)
		leftLines = append(leftLines, style.Render(line))
	}

	leftColumn := strings.Join(leftLines, "\n")

	// Right column - form details
	var rightLines []string
	if len(wsfs.availableForms) > 0 && wsfs.cursor < len(wsfs.availableForms) {
		form := wsfs.availableForms[wsfs.cursor]

		rightLines = append(rightLines, headerStyle.Render(form.Name))
		rightLines = append(rightLines, dimStyle.Render(fmt.Sprintf("CR: %.2f", form.CR)))
		rightLines = append(rightLines, "")

		rightLines = append(rightLines, statStyle.Render("STATS:"))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("AC: %d", form.AC)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("HP: %d", form.HP)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Speed: %s", form.Speed)))
		rightLines = append(rightLines, "")

		rightLines = append(rightLines, statStyle.Render("ABILITY SCORES:"))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("STR: %d (%+d)", form.AbilityScores.Strength, form.AbilityScores.GetModifier(models.Strength))))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("DEX: %d (%+d)", form.AbilityScores.Dexterity, form.AbilityScores.GetModifier(models.Dexterity))))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("CON: %d (%+d)", form.AbilityScores.Constitution, form.AbilityScores.GetModifier(models.Constitution))))
		rightLines = append(rightLines, "")

		// Actions
		if len(form.Actions) > 0 {
			rightLines = append(rightLines, statStyle.Render("ACTIONS:"))
			for _, action := range form.Actions {
				attackBonusStr := fmt.Sprintf("%+d", action.AttackBonus)
				damageBonusStr := fmt.Sprintf("%+d", action.DamageBonus)
				rightLines = append(rightLines, normalStyle.Render(action.Name))
				rightLines = append(rightLines, fmt.Sprintf("  %s %s to hit", dimStyle.Render("Attack:"), normalStyle.Render(attackBonusStr)))
				rightLines = append(rightLines, fmt.Sprintf("  %s %s%s %s", dimStyle.Render("Damage:"), normalStyle.Render(action.DamageDice), normalStyle.Render(damageBonusStr), normalStyle.Render(action.DamageType)))
				rightLines = append(rightLines, fmt.Sprintf("  %s %s", dimStyle.Render("Range:"), normalStyle.Render(action.Range)))
				if action.Description != "" {
					rightLines = append(rightLines, fmt.Sprintf("  %s", dimStyle.Render(action.Description)))
				}
			}
		}
	}

	rightColumn := strings.Join(rightLines, "\n")

	// Combine columns
	leftBox := lipgloss.NewStyle().Width(35).Render(leftColumn)
	rightBox := lipgloss.NewStyle().Width(40).Padding(0, 1).Render(rightColumn)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// Build final content
	var content strings.Builder
	content.WriteString(titleStyle.Render(fmt.Sprintf("🐺 SELECT WILD SHAPE FORMS")))
	content.WriteString("\n\n")
	content.WriteString(columns)
	content.WriteString("\n\n")

	// Help text
	if wsfs.CanConfirm() {
		content.WriteString(helpStyle.Render("↑/↓: Navigate • Space: Toggle • Enter: Confirm • Esc: Cancel"))
	} else {
		if len(wsfs.selectedForms) == 0 {
			content.WriteString(helpStyle.Render(fmt.Sprintf("↑/↓: Navigate • Space: Toggle • Esc: Cancel (select at least 1 form, up to %d)", wsfs.maxForms)))
		} else {
			content.WriteString(helpStyle.Render(fmt.Sprintf("↑/↓: Navigate • Space: Toggle • Esc: Cancel (can select up to %d more)", wsfs.maxForms-len(wsfs.selectedForms))))
		}
	}

	// Create popup with larger size
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(82)

	return lipgloss.Place(
		100,
		30,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
