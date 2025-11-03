// internal/ui/handlers_message.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// handleMessagePopupKeys handles message popup keyboard input
func (m *Model) handleMessagePopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cmd := m.messagePopup.Update(msg)
	return m, cmd
}

// handleRestPopupKeys handles rest popup keyboard input
func (m *Model) handleRestPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cmd := m.restPopup.Update(msg)

	// Check if rest was confirmed
	if m.restPopup.IsConfirmed() {
		healing := m.restPopup.GetHealing()
		if healing > 0 {
			m.message = fmt.Sprintf("Rest complete! Restored %d HP. HP: %d/%d",
				healing, m.character.CurrentHP, m.character.MaxHP)
		} else {
			m.message = fmt.Sprintf("Rest complete! HP: %d/%d", m.character.CurrentHP, m.character.MaxHP)
		}
		m.storage.Save(m.character)
		m.restPopup.Hide()
	} else if m.restPopup.IsCancelled() {
		m.message = "Rest cancelled"
		m.restPopup.Hide()
	}

	return m, cmd
}
