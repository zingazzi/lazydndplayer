// internal/ui/handlers_wildshape_levelup.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// handleWildShapeFormLevelUpSelectorKeys handles wild shape form level-up selector specific keys
func (m *Model) handleWildShapeFormLevelUpSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleWildShapeFormLevelUpSelectorKeys: key=%s", msg.String())

	var cmd tea.Cmd

	switch msg.String() {
	case "up", "k":
		m.wildShapeFormLevelUpSelector.Prev()
		return m, nil
	case "down", "j":
		m.wildShapeFormLevelUpSelector.Next()
		return m, nil
	case " ":
		m.wildShapeFormLevelUpSelector.ToggleSelection()
		return m, nil
	case "enter":
		if m.wildShapeFormLevelUpSelector.CanConfirm() {
			selectedForms := m.wildShapeFormLevelUpSelector.GetSelectedForms()
			debug.Log("Wild shape form selector confirmed: %d forms selected", len(selectedForms))

			// Update character's wild shape forms
			m.character.WildShape.KnownForms = make([]string, len(selectedForms))
			copy(m.character.WildShape.KnownForms, selectedForms)

			// Update max forms
			druidLevel := m.character.GetClassLevel("Druid")
			var maxForms int
			if druidLevel >= 8 {
				maxForms = 8
			} else if druidLevel >= 4 {
				maxForms = 6
			} else if druidLevel >= 2 {
				maxForms = 4
			} else {
				maxForms = 4
			}
			m.character.WildShape.MaxForms = maxForms

			m.wildShapeFormLevelUpSelector.Hide()

			// Save character
			m.storage.Save(m.character)
			m.message = fmt.Sprintf("Wild shape forms selected! (%d/%d forms)", len(selectedForms), maxForms)
		} else {
			selectedCount := m.wildShapeFormLevelUpSelector.GetSelectedCount()
			maxCount := m.wildShapeFormLevelUpSelector.GetMaxForms()
			m.message = fmt.Sprintf("Please select %d form(s) (%d/%d selected)", maxCount, selectedCount, maxCount)
		}
		return m, cmd
	case "esc":
		m.wildShapeFormLevelUpSelector.Hide()
		m.message = "Wild shape form selection cancelled"
		return m, nil
	}

	return m, nil
}
