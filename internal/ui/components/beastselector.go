// internal/ui/components/beastselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// BeastSelector handles beast companion selection for Beast Master rangers
type BeastSelector struct {
	visible      bool
	character    *models.Character
	cursor       int
	selectedType models.CompanionType
	beastTypes   []models.CompanionType
}

// NewBeastSelector creates a new beast selector
func NewBeastSelector(char *models.Character) *BeastSelector {
	return &BeastSelector{
		character:  char,
		beastTypes: []models.CompanionType{models.CompanionTypeLand, models.CompanionTypeSky, models.CompanionTypeSea},
		cursor:     0,
	}
}

// Show displays the beast selector
func (bs *BeastSelector) Show() {
	bs.visible = true
	bs.cursor = 0
	bs.selectedType = ""
}

// Hide hides the beast selector
func (bs *BeastSelector) Hide() {
	bs.visible = false
	bs.cursor = 0
	bs.selectedType = ""
}

// IsVisible returns whether the selector is visible
func (bs *BeastSelector) IsVisible() bool {
	return bs.visible
}

// GetSelectedType returns the currently selected beast type
func (bs *BeastSelector) GetSelectedType() models.CompanionType {
	if bs.cursor >= 0 && bs.cursor < len(bs.beastTypes) {
		return bs.beastTypes[bs.cursor]
	}
	return ""
}

// Next moves cursor down
func (bs *BeastSelector) Next() {
	if bs.cursor < len(bs.beastTypes)-1 {
		bs.cursor++
	}
}

// Prev moves cursor up
func (bs *BeastSelector) Prev() {
	if bs.cursor > 0 {
		bs.cursor--
	}
}

// Select confirms the current selection
func (bs *BeastSelector) Select() bool {
	if bs.cursor >= 0 && bs.cursor < len(bs.beastTypes) {
		bs.selectedType = bs.beastTypes[bs.cursor]
		return true
	}
	return false
}

// Update handles key presses
func (bs *BeastSelector) Update(msg tea.Msg) (BeastSelector, tea.Cmd) {
	if !bs.visible {
		return *bs, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			bs.Prev()
		case "down", "j":
			bs.Next()
		case "enter":
			if bs.Select() {
				bs.visible = false
			}
		case "esc":
			bs.selectedType = ""
			bs.Hide()
		}
	}

	return *bs, nil
}

// View renders the beast selector
func (bs *BeastSelector) View() string {
	if !bs.visible {
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

	content.WriteString(titleStyle.Render("SELECT BEAST COMPANION") + "\n\n")

	// Two-column layout: list on left, details on right
	var leftContent strings.Builder
	var rightContent strings.Builder

	// Get ranger level for calculations
	rangerLevel := bs.character.GetClassLevel("Ranger")
	if rangerLevel == 0 {
		rangerLevel = 3 // Default for selection
	}

	// Build beast list (left side)
	for i, beastType := range bs.beastTypes {
		cursor := "  "
		style := normalStyle
		if i == bs.cursor {
			cursor = "❯ "
			style = selectedStyle
		}
		leftContent.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, string(beastType))) + "\n")
	}

	// Build details (right side)
	if bs.cursor >= 0 && bs.cursor < len(bs.beastTypes) {
		currentType := bs.beastTypes[bs.cursor]
		companion := models.InitializeCompanion(currentType, rangerLevel)

		// Name
		rightContent.WriteString(selectedStyle.Render(string(currentType)) + "\n\n")

		// Stats
		rightContent.WriteString(statStyle.Render("STATS:") + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("AC: %d", companion.AC)) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("HP: %d", companion.MaxHP)) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Speed: %s", companion.Speed)) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Senses: %s", companion.Senses)) + "\n\n")

		// Ability Scores
		rightContent.WriteString(statStyle.Render("ABILITY SCORES:") + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("STR: %d (%+d)", companion.AbilityScores.Strength, companion.AbilityScores.GetModifier(models.Strength))) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("DEX: %d (%+d)", companion.AbilityScores.Dexterity, companion.AbilityScores.GetModifier(models.Dexterity))) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("CON: %d (%+d)", companion.AbilityScores.Constitution, companion.AbilityScores.GetModifier(models.Constitution))) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("INT: %d (%+d)", companion.AbilityScores.Intelligence, companion.AbilityScores.GetModifier(models.Intelligence))) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("WIS: %d (%+d)", companion.AbilityScores.Wisdom, companion.AbilityScores.GetModifier(models.Wisdom))) + "\n")
		rightContent.WriteString(normalStyle.Render(fmt.Sprintf("CHA: %d (%+d)", companion.AbilityScores.Charisma, companion.AbilityScores.GetModifier(models.Charisma))) + "\n\n")

		// Traits
		rightContent.WriteString(statStyle.Render("TRAITS:") + "\n")
		for _, trait := range companion.Traits {
			rightContent.WriteString(normalStyle.Render("• "+trait) + "\n")
		}
		rightContent.WriteString("\n")

		// Beast Strike details
		rightContent.WriteString(statStyle.Render("BEAST STRIKE:") + "\n")
		mechanics := models.NewCompanionMechanics(bs.character)
		attack := mechanics.GetCompanionAttack(companion)
		if attack != nil {
			attackBonusStr := fmt.Sprintf("%+d", attack.AttackBonus)
			damageBonusStr := fmt.Sprintf("%+d", attack.DamageBonus)
			rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Attack: %s to hit", attackBonusStr)) + "\n")
			rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Damage: %s%s %s", attack.DamageDice, damageBonusStr, attack.DamageType)) + "\n")
			rightContent.WriteString(normalStyle.Render(fmt.Sprintf("Range: %s", attack.Range)) + "\n")
			if companion.SpecialNotes != "" {
				rightContent.WriteString(dimStyle.Render("Note: "+companion.SpecialNotes) + "\n")
			}
		}
	}

	// Join left and right in two columns
	leftBox := lipgloss.NewStyle().
		Width(30).
		Height(25).
		Padding(1).
		Render(leftContent.String())

	rightBox := lipgloss.NewStyle().
		Width(65).
		Height(25).
		Padding(1).
		Render(rightContent.String())

	twoColumns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	content.WriteString(twoColumns)

	content.WriteString("\n" + dimStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Cancel"))

	// Popup style
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(100)

	return lipgloss.Place(
		120,
		35,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
