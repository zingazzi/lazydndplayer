// internal/ui/components/wildshapeformlevelupselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// WildShapeFormLevelUpSelector handles wild shape form selection during level up (similar to cantrip selector)
type WildShapeFormLevelUpSelector struct {
	visible          bool
	availableForms   []models.BeastDefinition
	selectedForms    []string // Names of selected forms
	cursor           int
	maxForms         int
	druidLevel       int
	character        *models.Character
}

// NewWildShapeFormLevelUpSelector creates a new wild shape form level-up selector
func NewWildShapeFormLevelUpSelector(char *models.Character) *WildShapeFormLevelUpSelector {
	return &WildShapeFormLevelUpSelector{
		visible:         false,
		availableForms:  []models.BeastDefinition{},
		selectedForms:   []string{},
		cursor:          0,
		character:       char,
	}
}

// Show displays the wild shape form selector for level up
func (wsfs *WildShapeFormLevelUpSelector) Show(druidLevel int, maxForms int) {
	wsfs.druidLevel = druidLevel
	wsfs.maxForms = maxForms
	wsfs.visible = true
	wsfs.cursor = 0

	// Pre-select existing forms from character
	wsfs.selectedForms = make([]string, len(wsfs.character.WildShape.KnownForms))
	copy(wsfs.selectedForms, wsfs.character.WildShape.KnownForms)

	// Load available forms based on druid level
	wsfs.loadForms()
}

// Hide hides the selector
func (wsfs *WildShapeFormLevelUpSelector) Hide() {
	wsfs.visible = false
}

// IsVisible returns whether the selector is visible
func (wsfs *WildShapeFormLevelUpSelector) IsVisible() bool {
	return wsfs.visible
}

// loadForms loads available wild shape forms based on druid level
func (wsfs *WildShapeFormLevelUpSelector) loadForms() {
	// Get available wild shape forms based on druid level
	wsfs.availableForms = models.GetWildShapeForms(wsfs.druidLevel)
}

// Update handles key presses
func (wsfs *WildShapeFormLevelUpSelector) Update(msg tea.Msg) (WildShapeFormLevelUpSelector, tea.Cmd) {
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
			if len(wsfs.selectedForms) == wsfs.maxForms {
				// Confirmed
				return *wsfs, nil
			}
		case "esc":
			wsfs.Hide()
		}
	}

	return *wsfs, nil
}

// Next moves cursor down
func (wsfs *WildShapeFormLevelUpSelector) Next() {
	if wsfs.cursor < len(wsfs.availableForms)-1 {
		wsfs.cursor++
	}
}

// Prev moves cursor up
func (wsfs *WildShapeFormLevelUpSelector) Prev() {
	if wsfs.cursor > 0 {
		wsfs.cursor--
	}
}

// ToggleSelection toggles selection of the current form
func (wsfs *WildShapeFormLevelUpSelector) ToggleSelection() {
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
func (wsfs *WildShapeFormLevelUpSelector) IsSelected(formName string) bool {
	for _, name := range wsfs.selectedForms {
		if name == formName {
			return true
		}
	}
	return false
}

// CanConfirm returns true if the correct number of forms is selected
func (wsfs *WildShapeFormLevelUpSelector) CanConfirm() bool {
	return len(wsfs.selectedForms) == wsfs.maxForms
}

// GetSelectedForms returns the selected form names
func (wsfs *WildShapeFormLevelUpSelector) GetSelectedForms() []string {
	return wsfs.selectedForms
}

// GetMaxForms returns the maximum number of forms
func (wsfs *WildShapeFormLevelUpSelector) GetMaxForms() int {
	return wsfs.maxForms
}

// GetSelectedCount returns the number of selected forms
func (wsfs *WildShapeFormLevelUpSelector) GetSelectedCount() int {
	return len(wsfs.selectedForms)
}

// View renders the selector
func (wsfs *WildShapeFormLevelUpSelector) View() string {
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

	// Determine CR limit and flying restriction
	var maxCR float64
	canFly := false
	if wsfs.druidLevel >= 8 {
		maxCR = 1.0
		canFly = true
	} else if wsfs.druidLevel >= 4 {
		maxCR = 0.5
	} else {
		maxCR = 0.25
	}

	// Left column - form list
	var leftLines []string
	titleText := fmt.Sprintf("Select %d wild shape forms  (%d/%d)", wsfs.maxForms, len(wsfs.selectedForms), wsfs.maxForms)
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

		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("AC: %d", form.AC)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("HP: %d", form.HP)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Speed: %s", form.Speed)))
		rightLines = append(rightLines, "")

		// Ability Scores
		rightLines = append(rightLines, headerStyle.Render("Ability Scores:"))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("STR: %d (%+d)", form.AbilityScores.Strength, (&form.AbilityScores).GetModifier(models.Strength))))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("DEX: %d (%+d)", form.AbilityScores.Dexterity, (&form.AbilityScores).GetModifier(models.Dexterity))))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("CON: %d (%+d)", form.AbilityScores.Constitution, (&form.AbilityScores).GetModifier(models.Constitution))))
		rightLines = append(rightLines, "")

		// Actions
		if len(form.Actions) > 0 {
			rightLines = append(rightLines, headerStyle.Render("Actions:"))
			for _, action := range form.Actions {
				attackBonusStr := fmt.Sprintf("%+d", action.AttackBonus)
				damageBonusStr := fmt.Sprintf("%+d", action.DamageBonus)
				rightLines = append(rightLines, normalStyle.Render(action.Name))
				rightLines = append(rightLines, dimStyle.Render(fmt.Sprintf("  Attack: %s to hit", attackBonusStr)))
				rightLines = append(rightLines, dimStyle.Render(fmt.Sprintf("  Damage: %s%s %s", action.DamageDice, damageBonusStr, action.DamageType)))
			}
		}
	}

	rightColumn := strings.Join(rightLines, "\n")

	// Combine columns
	leftBox := lipgloss.NewStyle().Width(40).Render(leftColumn)
	rightBox := lipgloss.NewStyle().Width(42).Padding(0, 1).Render(rightColumn)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// Build final content
	var content strings.Builder
	content.WriteString(titleStyle.Render("🐺 SELECT WILD SHAPE FORMS"))
	content.WriteString("\n\n")
	content.WriteString(columns)
	content.WriteString("\n\n")

	// Help text
	if wsfs.CanConfirm() {
		content.WriteString(helpStyle.Render("↑/↓: Navigate • Space: Toggle • Enter: Confirm • Esc: Cancel"))
	} else {
		content.WriteString(helpStyle.Render(fmt.Sprintf("↑/↓: Navigate • Space: Toggle • Esc: Cancel (need %d more)", wsfs.maxForms-len(wsfs.selectedForms))))
	}

	// Create popup with larger size
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(90)

	return lipgloss.Place(
		100,
		30,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
