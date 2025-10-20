// internal/ui/components/attackmenu.go
package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// AttackMenu represents a small contextual menu for attack options
type AttackMenu struct {
	visible       bool
	attack        *models.Attack
	selectedIndex int
	options       []string
}

// NewAttackMenu creates a new attack menu
func NewAttackMenu() *AttackMenu {
	return &AttackMenu{
		visible:       false,
		selectedIndex: 0,
	}
}

// Show displays the menu for a specific attack
func (am *AttackMenu) Show(attack *models.Attack) {
	am.visible = true
	am.attack = attack
	am.selectedIndex = 0

	// Build dynamic options based on attack type
	am.options = []string{
		"Attack with Advantage",
		"Attack with Disadvantage",
		"Attack (Normal)",
		"─────────────────────────",
	}

	// Check if weapon is versatile
	if attack.VersatileDamage != "" {
		// Versatile weapon - show 1-hand and 2-hand options
		am.options = append(am.options,
			fmt.Sprintf("1-Hand Damage (%s)", attack.DamageDice),
			fmt.Sprintf("1-Hand Critical (%s x2)", attack.DamageDice),
			fmt.Sprintf("2-Hands Damage (%s)", attack.VersatileDamage),
			fmt.Sprintf("2-Hands Critical (%s x2)", attack.VersatileDamage),
		)
	} else {
		// Regular weapon - show normal damage options
		am.options = append(am.options,
			fmt.Sprintf("Damage (%s)", attack.DamageDice),
			fmt.Sprintf("Critical Hit (%s x2)", attack.DamageDice),
		)
	}
}

// Hide hides the menu
func (am *AttackMenu) Hide() {
	am.visible = false
	am.attack = nil
}

// IsVisible returns true if the menu is visible
func (am *AttackMenu) IsVisible() bool {
	return am.visible
}

// Next moves selection down
func (am *AttackMenu) Next() {
	am.selectedIndex++
	if am.selectedIndex >= len(am.options) {
		am.selectedIndex = len(am.options) - 1
	}
	// Skip separator
	if am.options[am.selectedIndex] == "─────────────────────" {
		am.selectedIndex++
		if am.selectedIndex >= len(am.options) {
			am.selectedIndex = len(am.options) - 1
		}
	}
}

// Prev moves selection up
func (am *AttackMenu) Prev() {
	am.selectedIndex--
	if am.selectedIndex < 0 {
		am.selectedIndex = 0
	}
	// Skip separator
	if am.options[am.selectedIndex] == "─────────────────────" {
		am.selectedIndex--
		if am.selectedIndex < 0 {
			am.selectedIndex = 0
		}
	}
}

// GetSelectedOption returns the currently selected option
func (am *AttackMenu) GetSelectedOption() string {
	if am.selectedIndex >= 0 && am.selectedIndex < len(am.options) {
		return am.options[am.selectedIndex]
	}
	return ""
}

// GetAttack returns the current attack
func (am *AttackMenu) GetAttack() *models.Attack {
	return am.attack
}

// View renders the attack menu
func (am *AttackMenu) View(screenWidth, screenHeight int) string {
	if !am.visible || am.attack == nil {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237")).
		Padding(0, 1)

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	var lines []string

	// Title with attack name
	lines = append(lines, titleStyle.Render(fmt.Sprintf("⚔️  %s", am.attack.Name)))
	lines = append(lines, normalStyle.Render(am.attack.GetAttackSummary()))
	lines = append(lines, "")

	// Options
	for i, option := range am.options {
		if option == "─────────────────────" {
			lines = append(lines, separatorStyle.Render(option))
		} else if i == am.selectedIndex {
			lines = append(lines, selectedStyle.Render("▶ "+option))
		} else {
			lines = append(lines, normalStyle.Render("  "+option))
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Render("↑/↓: Navigate • Enter: Select • ESC: Cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Small box with border
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2)

	box := boxStyle.Render(content)

	// Center the box
	return lipgloss.Place(
		screenWidth,
		screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
