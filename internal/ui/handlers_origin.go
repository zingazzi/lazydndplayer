// internal/ui/handlers_origin.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// handleOriginSelectorKeys handles keyboard input for the origin selector
func (m *Model) handleOriginSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.originSelector.Prev()
	case "down", "j":
		m.originSelector.Next()
	case "enter":
		selectedOrigin := m.originSelector.GetSelected()
		if selectedOrigin != nil {
			debug.Log("handleOriginSelectorKeys: Origin selected=%s, IsInWizard=%v, wizardStep=%d", selectedOrigin.Name, m.IsInWizard(), m.GetWizardStep())
			// Check if origin has ability choice using service
			if m.originService.RequiresAbilityChoice(selectedOrigin) {
				debug.Log("handleOriginSelectorKeys: Origin requires ability choice, showing ability choice selector")
				// Store origin temporarily and show ability choice selector
				m.SetPendingOrigin(selectedOrigin)
				m.originSelector.Hide()
				choices := m.originService.GetAbilityChoices(selectedOrigin)
				m.abilityChoiceSelector.Show(selectedOrigin.Name, choices, m.character)
				m.message = "Choose an ability score to increase..."
			} else {
				// Apply origin directly (no choice needed) using service
				debug.Log("handleOriginSelectorKeys: Origin does not require ability choice, applying directly")
				originMsg := m.originService.ApplyOrigin(m.character, selectedOrigin, "")
				m.storage.Save(m.character)
				m.originSelector.Hide()

				// If in wizard mode, advance to next step (this will set the correct wizard message)
				if m.IsInWizard() {
					debug.Log("handleOriginSelectorKeys: In wizard mode, advancing from Origin step")
					m.advanceWizardStep()
				} else {
					m.message = originMsg
				}
			}
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		m.originSelector.Hide()
		m.message = "Origin selection cancelled"
	}
	return m, nil
}

// handleAlignmentSelectorKeys handles alignment selector input (Origin panel only)
func (m *Model) handleAlignmentSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		m.alignmentSelector.Hide()
		m.message = ""
	case "up", "k":
		m.alignmentSelector.Prev()
	case "down", "j":
		m.alignmentSelector.Next()
	case "enter":
		selected := m.alignmentSelector.GetSelectedAlignment()
		m.character.Alignment = selected
		m.alignmentSelector.Hide()
		m.storage.Save(m.character)

		// If in wizard mode, complete wizard
		if m.IsInWizard() {
			m.completeWizard()
		} else {
			m.message = fmt.Sprintf("Alignment set to %s", selected)
		}
	}
	return m, nil
}

// handleTraitSelectorKeys handles trait selector input (Origin panel only)
func (m *Model) handleTraitSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If in custom mode, handle it separately to allow all keys including 'p' and space
	if m.traitSelector.IsCustomMode() {
		switch msg.String() {
		case "esc":
			m.traitSelector.Hide()
			m.message = ""
		case "enter":
			// Save custom value and add to the list
			customValue := m.traitSelector.GetSelectedTrait()
			if customValue != "" {
				// Add to appropriate array
				switch m.traitSelector.TraitType {
				case models.TraitPersonality:
					m.character.Personality = append(m.character.Personality, customValue)
					m.message = "Custom personality trait added!"
				case models.TraitIdeal:
					m.character.Ideal = append(m.character.Ideal, customValue)
					m.message = "Custom ideal added!"
				case models.TraitBond:
					m.character.Bond = append(m.character.Bond, customValue)
					m.message = "Custom bond added!"
				case models.TraitFlaw:
					m.character.Flaw = append(m.character.Flaw, customValue)
					m.message = "Custom flaw added!"
				}
				m.traitSelector.Hide()
				m.storage.Save(m.character)
			}
		default:
			// Pass all other keys to text input (including 'p', space, etc.)
			m.traitSelector.Update(msg)
		}
		return m, nil
	}

	// Normal multi-select mode
	switch msg.String() {
	case "esc":
		m.traitSelector.Hide()
		m.message = ""
	case "up", "k":
		m.traitSelector.Prev()
	case "down", "j":
		m.traitSelector.Next()
	case " ":
		// Space key toggles the current item
		m.traitSelector.ToggleItem()
	case "enter":
		// Check if [Custom] is selected to enter custom mode
		m.traitSelector.ToggleCustomMode()

		// If not entering custom mode, save all selected items
		if !m.traitSelector.IsCustomMode() {
			selectedItems := m.traitSelector.GetSelectedItems()

			// Save the trait based on current type
			switch m.traitSelector.TraitType {
			case models.TraitPersonality:
				m.character.Personality = selectedItems
				m.message = "Personality traits saved!"
			case models.TraitIdeal:
				m.character.Ideal = selectedItems
				m.message = "Ideals saved!"
			case models.TraitBond:
				m.character.Bond = selectedItems
				m.message = "Bonds saved!"
			case models.TraitFlaw:
				m.character.Flaw = selectedItems
				m.message = "Flaws saved!"
			}
			m.traitSelector.Hide()
			m.storage.Save(m.character)
		}
	}
	return m, nil
}

// handleBackstoryEditorKeys handles backstory editor input (Origin panel only)
func (m *Model) handleBackstoryEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+enter":
		// Save and close
		m.character.Backstory = m.backstoryEditor.GetValue()
		m.backstoryEditor.Hide()
		m.message = "Backstory saved!"
		m.storage.Save(m.character)
	default:
		// Update textarea
		m.backstoryEditor.Update(msg)
	}
	return m, nil
}

// handleInputPopupKeys handles input popup input (used in Origin panel for height/weight)
func (m *Model) handleInputPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.inputPopup.Hide()
		m.ClearInputPopupContext()
		m.message = "Cancelled"
	case "enter":
		// Save based on context
		value := m.inputPopup.GetValue()
		switch m.GetInputPopupContext() {
		case "height":
			m.character.Height = value
			m.message = "Height updated!"
		case "weight":
			m.character.Weight = value
			m.message = "Weight updated!"
		}
		m.inputPopup.Hide()
		m.ClearInputPopupContext()
		m.storage.Save(m.character)
	default:
		// Update text input
		m.inputPopup.Update(msg)
	}
	return m, nil
}
