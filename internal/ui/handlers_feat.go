// internal/ui/handlers_feat.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// handleFeatSelectorKeys handles feat selector specific keys
func (m *Model) handleFeatSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.featSelector.Prev()
	case "down", "j":
		m.featSelector.Next()
	case "pgup", "ctrl+u":
		m.featSelector.PageUp()
	case "pgdown", "ctrl+d":
		m.featSelector.PageDown()
	case "left", "h":
		m.featSelector.PrevCategory()
	case "right", "l":
		m.featSelector.NextCategory()
	case "enter":
		selectedFeat := m.featSelector.GetSelectedFeat()
		if selectedFeat != nil {
			// Check if we're in delete mode
			if m.featSelector.IsDeleteMode() {
				// Use service to remove feat
				m.message = m.featService.RemoveFeat(m.character, selectedFeat)
				m.storage.Save(m.character)
				m.featSelector.Hide()
			} else {
				// Check if the feat can be selected (prerequisites met)
				if !m.featSelector.CanSelectCurrentFeat() {
					m.message = fmt.Sprintf("Cannot select %s: Prerequisites not met!", selectedFeat.Name)
					return m, nil
				}

				// Use service to check if feat can be applied
				msg, err := m.featService.ApplyFeat(m.character, selectedFeat, "")
				if err != nil {
					m.message = err.Error()
					m.featSelector.Hide()
				} else {
					// Check if this feat has ability choices
					if m.featService.RequiresAbilityChoice(selectedFeat) {
						// Store the feat and show ability choice selector
						m.SetPendingFeat(selectedFeat)
						m.featSelector.Hide()
						choices := m.featService.GetAbilityChoices(selectedFeat)
						m.abilityChoiceSelector.Show(selectedFeat.Name, choices, m.character)
						m.message = "Choose which ability to increase"
					} else {
						// Feat applied automatically (no ability choice)
						m.message = msg
						m.storage.Save(m.character)
						m.featSelector.Hide()
					}
				}
			}
		}
	case "esc":
		m.featSelector.Hide()
		m.message = "Feat selection cancelled"
	}
	return m, nil
}

// handleAbilityChoiceSelectorKeys handles keyboard input for the ability choice selector
// (used for both feat and origin ability choices)
func (m *Model) handleAbilityChoiceSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.abilityChoiceSelector.Prev()
	case "down", "j":
		m.abilityChoiceSelector.Next()
	case "enter":
		chosenAbility := m.abilityChoiceSelector.GetSelectedAbility()
		if chosenAbility == "" {
			return m, nil
		}

		// Handle feat ability choice
		if pendingFeat := m.GetPendingFeat(); pendingFeat != nil {
			// Use service to apply feat with chosen ability
			msg, err := m.featService.ApplyFeat(m.character, pendingFeat, chosenAbility)
			if err != nil {
				m.message = err.Error()
			} else {
				m.message = msg
				m.storage.Save(m.character)
			}
			m.ClearPendingFeat()
			m.abilityChoiceSelector.Hide()
		}

		// Handle origin ability choice
		if pendingOrigin := m.GetPendingOrigin(); pendingOrigin != nil {
			// Use service to apply origin with chosen ability
			m.message = m.originService.ApplyOrigin(m.character, pendingOrigin, chosenAbility)
			m.storage.Save(m.character)
			m.ClearPendingOrigin()
			m.abilityChoiceSelector.Hide()
		}
	case "esc":
		// Cancel ability choice
		if pendingFeat := m.GetPendingFeat(); pendingFeat != nil {
			// Use service to cancel feat selection
			m.featService.CancelFeatSelection(m.character, pendingFeat)
			m.storage.Save(m.character)
			m.message = "Feat selection cancelled"
			m.ClearPendingFeat()
		}

		if m.GetPendingOrigin() != nil {
			m.message = "Origin selection cancelled"
			m.ClearPendingOrigin()
		}

		m.abilityChoiceSelector.Hide()
	}
	return m, nil
}
