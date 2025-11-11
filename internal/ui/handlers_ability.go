// internal/ui/handlers_ability.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// handleAbilityRollerKeys handles ability roller specific keys
// Returns (model, cmd, handled) where handled indicates if the key was processed
func (m *Model) handleAbilityRollerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch msg.String() {
	case "up", "k":
		m.abilityRoller.Prev()
		return m, nil, true
	case "down", "j":
		m.abilityRoller.Next()
		return m, nil, true
	case "tab":
		m.abilityRoller.SwitchFocus()
		m.message = "Switched focus"
		return m, nil, true
	case "space":
		m.abilityRoller.ToggleType()
		return m, nil, true
	case "enter":
		// Roll the dice!
		expr := m.abilityRoller.GetRollExpression(m.character)
		description := m.abilityRoller.GetRollDescription(m.character)
		m.dicePanel.Roll(expr)
		m.message = fmt.Sprintf("%s: %s", description, m.dicePanel.LastMessage)
		m.abilityRoller.Hide()
		return m, nil, true
	case "esc":
		m.abilityRoller.Hide()
		m.message = "Roll cancelled"
		return m, nil, true
	}
	return m, nil, false
}
