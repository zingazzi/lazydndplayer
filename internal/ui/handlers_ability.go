// internal/ui/handlers_ability.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// handleAbilityRollerKeys handles ability roller specific keys
func (m *Model) handleAbilityRollerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.abilityRoller.Prev()
	case "down", "j":
		m.abilityRoller.Next()
	case "tab":
		m.abilityRoller.SwitchFocus()
		m.message = "Switched focus"
	case "space":
		m.abilityRoller.ToggleType()
	case "enter":
		// Roll the dice!
		expr := m.abilityRoller.GetRollExpression(m.character)
		description := m.abilityRoller.GetRollDescription(m.character)
		m.dicePanel.Roll(expr)
		m.message = fmt.Sprintf("%s: %s", description, m.dicePanel.LastMessage)
		m.abilityRoller.Hide()
	case "esc":
		m.abilityRoller.Hide()
		m.message = "Roll cancelled"
	}
	return m, nil
}
