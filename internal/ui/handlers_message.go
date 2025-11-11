// internal/ui/handlers_message.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// handleMessagePopupKeys handles message popup keyboard input
// Returns (model, cmd, handled) where handled indicates if the key was processed
func (m *Model) handleMessagePopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	// Message popup handles all keys (esc, enter, space)
	cmd := m.messagePopup.Update(msg)
	return m, cmd, true
}

// handleRestPopupKeys handles rest popup keyboard input
// Returns (model, cmd, handled) where handled indicates if the key was processed
func (m *Model) handleRestPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	// Handle 'c' key for changing companion during long rest
	if msg.String() == "c" && m.restPopup.IsVisible() {
		// Check if it's a long rest and character is Beast Master
		if m.restPopup.GetRestType() == components.LongRestType && m.character.IsBeastMaster() {
			m.restPopup.Hide()
			m.beastSelector.Show()
			m.message = "Select a new beast companion..."
			return m, nil, true
		}
	}

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
		return m, cmd, true
	} else if m.restPopup.IsCancelled() {
		m.message = "Rest cancelled"
		m.restPopup.Hide()
		return m, cmd, true
	}

	// Rest popup handles all keys when visible
	return m, cmd, true
}
