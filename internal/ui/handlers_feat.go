// internal/ui/handlers_feat.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/models"
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
				// Remove the feat
				for i, featName := range m.character.Feats {
					if featName == selectedFeat.Name {
						m.character.Feats = append(m.character.Feats[:i], m.character.Feats[i+1:]...)
						break
					}
				}
				// Remove feat benefits (ability increases, HP, speed, etc.)
				models.RemoveFeatBenefits(m.character, *selectedFeat)

				m.message = fmt.Sprintf("Feat removed: %s (benefits reversed)", selectedFeat.Name)
				m.storage.Save(m.character)
				m.featSelector.Hide()
			} else {
				// Check if the feat can be selected (prerequisites met)
				if !m.featSelector.CanSelectCurrentFeat() {
					m.message = fmt.Sprintf("Cannot select %s: Prerequisites not met!", selectedFeat.Name)
					return m, nil
				}

				// Add mode: Check if character already has this feat
				if models.HasFeat(m.character, selectedFeat.Name) && !selectedFeat.Repeatable {
					m.message = fmt.Sprintf("You already have %s and it's not repeatable", selectedFeat.Name)
					m.featSelector.Hide()
				} else {
					// Add feat to character
					err := models.AddFeatToCharacter(m.character, selectedFeat.Name)
					if err != nil {
						m.message = fmt.Sprintf("Error adding feat: %v", err)
						m.featSelector.Hide()
					} else {
						// Check if this feat has ability choices
						if models.HasAbilityChoice(*selectedFeat) {
							// Store the feat and show ability choice selector
							m.pendingFeat = selectedFeat
							m.featSelector.Hide()
							choices := models.GetAbilityChoices(*selectedFeat)
							m.abilityChoiceSelector.Show(selectedFeat.Name, choices, m.character)
							m.message = "Choose which ability to increase"
						} else {
							// Apply feat benefits automatically (no ability choice)
							models.ApplyFeatBenefits(m.character, *selectedFeat, "")
							m.message = fmt.Sprintf("Feat gained: %s!", selectedFeat.Name)
							// Save character after feat selection
							m.storage.Save(m.character)
							m.featSelector.Hide()
						}
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
		if m.pendingFeat != nil {
			// Apply feat benefits with the chosen ability
			models.ApplyFeatBenefits(m.character, *m.pendingFeat, chosenAbility)
			m.message = fmt.Sprintf("Feat gained: %s (+1 %s)!", m.pendingFeat.Name, chosenAbility)
			m.storage.Save(m.character)
			m.pendingFeat = nil
			m.abilityChoiceSelector.Hide()
		}

		// Handle origin ability choice
		if m.pendingOrigin != nil {
			// Remove old origin first
			if m.character.Origin != "" {
				oldOrigin := models.GetOriginByName(m.character.Origin)
				if oldOrigin != nil {
					models.RemoveOriginBenefits(m.character, *oldOrigin)
				}
			}

			// Apply new origin with chosen ability
			m.character.Origin = m.pendingOrigin.Name
			models.ApplyOriginBenefits(m.character, *m.pendingOrigin, chosenAbility)
			m.message = fmt.Sprintf("Origin changed to: %s (+1 %s)!", m.pendingOrigin.Name, chosenAbility)
			m.storage.Save(m.character)
			m.pendingOrigin = nil
			m.abilityChoiceSelector.Hide()
		}
	case "esc":
		// Cancel ability choice
		if m.pendingFeat != nil {
			// Remove the feat from character since we're cancelling
			for i, featName := range m.character.Feats {
				if featName == m.pendingFeat.Name {
					m.character.Feats = append(m.character.Feats[:i], m.character.Feats[i+1:]...)
					break
				}
			}
			m.storage.Save(m.character)
			m.message = "Feat selection cancelled"
			m.pendingFeat = nil
		}

		if m.pendingOrigin != nil {
			m.message = "Origin selection cancelled"
			m.pendingOrigin = nil
		}

		m.abilityChoiceSelector.Hide()
	}
	return m, nil
}
