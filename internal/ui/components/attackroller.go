// internal/ui/components/attackroller.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/dice"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// AttackRoller represents the attack rolling interface
type AttackRoller struct {
	visible        bool
	selectedAttack *models.Attack
	selectedIndex  int
	availableAttacks []models.Attack
	state          string // "select_attack", "select_roll_type", "show_result"
	rollType       string // "normal", "advantage", "disadvantage"
	lastResult     string
	character      *models.Character
}

// NewAttackRoller creates a new attack roller
func NewAttackRoller() *AttackRoller {
	return &AttackRoller{
		visible:       false,
		selectedIndex: 0,
		state:         "select_attack",
		rollType:      "normal",
	}
}

// Show displays the attack roller with available attacks
func (ar *AttackRoller) Show(char *models.Character) {
	ar.character = char
	ar.visible = true
	ar.state = "select_attack"
	ar.selectedIndex = 0
	ar.lastResult = ""

	// Generate attacks based on equipped weapons
	attackList := models.GenerateAttacks(char)
	ar.availableAttacks = attackList.Attacks
}

// Hide hides the attack roller
func (ar *AttackRoller) Hide() {
	ar.visible = false
	ar.selectedAttack = nil
	ar.state = "select_attack"
	ar.lastResult = ""
}

// IsVisible returns true if the roller is visible
func (ar *AttackRoller) IsVisible() bool {
	return ar.visible
}

// Next moves selection down
func (ar *AttackRoller) Next() {
	if ar.state == "select_attack" {
		if ar.selectedIndex < len(ar.availableAttacks)-1 {
			ar.selectedIndex++
		}
	}
}

// Prev moves selection up
func (ar *AttackRoller) Prev() {
	if ar.state == "select_attack" {
		if ar.selectedIndex > 0 {
			ar.selectedIndex--
		}
	}
}

// SelectAttack confirms attack selection
func (ar *AttackRoller) SelectAttack() {
	if ar.selectedIndex >= 0 && ar.selectedIndex < len(ar.availableAttacks) {
		ar.selectedAttack = &ar.availableAttacks[ar.selectedIndex]
		ar.state = "select_roll_type"
	}
}

// RollAttack performs an attack roll
func (ar *AttackRoller) RollAttack() string {
	if ar.selectedAttack == nil {
		return ""
	}

	var rollType dice.RollType
	var advantageStr string

	switch ar.rollType {
	case "advantage":
		rollType = dice.Advantage
		advantageStr = "Advantage"
	case "disadvantage":
		rollType = dice.Disadvantage
		advantageStr = "Disadvantage"
	default:
		rollType = dice.Normal
		advantageStr = ""
	}

	result, err := dice.Roll("1d20", rollType)
	if err != nil {
		return fmt.Sprintf("Error rolling: %v", err)
	}

	// For advantage/disadvantage, only the first roll in the array is used
	roll := 0
	if len(result.Rolls) > 0 {
		roll = result.Rolls[0]
	}

	total := roll + ar.selectedAttack.AttackBonus

	return ar.selectedAttack.FormatAttackRoll(roll, total, advantageStr)
}

// RollDamage performs a damage roll
func (ar *AttackRoller) RollDamage() string {
	if ar.selectedAttack == nil {
		return ""
	}

	result, err := dice.Roll(ar.selectedAttack.DamageDice, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling damage: %v", err)
	}

	total := result.Total + ar.selectedAttack.DamageBonus

	return ar.selectedAttack.FormatDamageRoll(result.Rolls, total)
}

// GetSelectedAttack returns the currently selected attack
func (ar *AttackRoller) GetSelectedAttack() *models.Attack {
	return ar.selectedAttack
}

// GetState returns the current state
func (ar *AttackRoller) GetState() string {
	return ar.state
}

// SetRollType sets the roll type (normal, advantage, disadvantage)
func (ar *AttackRoller) SetRollType(rollType string) {
	ar.rollType = rollType
}

// SetSelectedIndex sets the selected attack index
func (ar *AttackRoller) SetSelectedIndex(index int) {
	if index >= 0 && index < len(ar.availableAttacks) {
		ar.selectedIndex = index
	}
}

// GetAvailableAttacks returns the list of available attacks
func (ar *AttackRoller) GetAvailableAttacks() []models.Attack {
	return ar.availableAttacks
}

// View renders the attack roller
func (ar *AttackRoller) View(width, height int) string {
	if !ar.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237")).
		Padding(0, 2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(1, 2)

	var content string

	switch ar.state {
	case "select_attack":
		lines := []string{
			titleStyle.Render("⚔️  SELECT ATTACK"),
			"",
		}

		for i, attack := range ar.availableAttacks {
			line := fmt.Sprintf("%-20s  %s", attack.Name, attack.GetAttackSummary())
			if i == ar.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, normalStyle.Render("  "+line))
			}
		}

		lines = append(lines, "")
		lines = append(lines, helpStyle.Render("↑/↓: Select • Enter: Choose • ESC: Cancel"))

		content = strings.Join(lines, "\n")

	case "select_roll_type":
		lines := []string{
			titleStyle.Render(fmt.Sprintf("⚔️  %s", ar.selectedAttack.Name)),
			"",
			normalStyle.Render(fmt.Sprintf("Attack Bonus: +%d", ar.selectedAttack.AttackBonus)),
			normalStyle.Render(fmt.Sprintf("Damage: %s+%d %s", ar.selectedAttack.DamageDice, ar.selectedAttack.DamageBonus, ar.selectedAttack.DamageType)),
			"",
			titleStyle.Render("Choose Action:"),
			"",
			normalStyle.Render("  'a' - Roll Attack (1d20 + bonus)"),
			normalStyle.Render("  'd' - Roll Damage"),
			normalStyle.Render("  'v' - Roll Attack with Advantage"),
			normalStyle.Render("  'x' - Roll Attack with Disadvantage"),
			"",
			helpStyle.Render("ESC: Back"),
		}

		content = strings.Join(lines, "\n")
	}

	// Create a box around the content
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Width(width - 4).
		Padding(1)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content),
	)
}
