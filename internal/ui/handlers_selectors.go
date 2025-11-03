// internal/ui/handlers_selectors.go
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// handleLanguageSelectorKeys handles language selector specific keys
func (m *Model) handleLanguageSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.languageSelector.Prev()
	case "down", "j":
		m.languageSelector.Next()
	case "enter":
		selectedLanguage := m.languageSelector.GetSelectedLanguage()
		if selectedLanguage != "" {
			// Check if we're in delete mode
			if m.languageSelector.IsDeleteMode() {
				// Remove the language
				for i, lang := range m.character.Languages {
					if lang == selectedLanguage {
						m.character.Languages = append(m.character.Languages[:i], m.character.Languages[i+1:]...)
						break
					}
				}
				m.message = fmt.Sprintf("Language removed: %s", selectedLanguage)
				m.storage.Save(m.character)
				m.languageSelector.Hide()
			} else {
				// Add mode: Check if adding a new language (from Traits panel) or replacing placeholder (from species selection)
				foundPlaceholder := false
				for i, lang := range m.character.Languages {
					if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
						m.character.Languages[i] = selectedLanguage
						foundPlaceholder = true
						break
					}
				}

				// If no placeholder found, just append the new language
				if !foundPlaceholder {
					m.character.Languages = append(m.character.Languages, selectedLanguage)
				}

				// After language selection, check for skill, spell, or feat selection (only during species selection)
				species := models.GetSpeciesByName(m.character.Race)
				if species != nil && foundPlaceholder {
					if models.HasSkillChoice(species) {
						m.skillSelector.Show()
						m.message = "Select your skill proficiency..."
					} else if models.HasSpellChoice(species) {
						// Show wizard cantrip selector for High Elf
						cantrips := models.GetWizardCantrips()
						m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
						m.spellSelector.Show()
						m.message = "Select your wizard cantrip..."
					} else if models.HasFeatChoice(species) {
						// Show feat selector for origin feat
						m.featSelector.Show(m.character, true)
						m.message = "Select your origin feat..."
					} else {
						m.message = fmt.Sprintf("Language selected: %s (Total languages: %d)", selectedLanguage, len(m.character.Languages))
						// Save when selection is complete (no more selections needed)
						m.storage.Save(m.character)
					}
				} else {
					m.message = fmt.Sprintf("Language learned: %s!", selectedLanguage)
					// Save when adding a new language (not replacing placeholder)
					m.storage.Save(m.character)
				}
				m.languageSelector.Hide()
			}
		}
	case "esc":
		m.languageSelector.Hide()
		m.message = "Language selection cancelled"
	}
	return m, nil
}

// handleToolSelectorKeys handles tool selector specific keys
func (m *Model) handleToolSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.toolSelector.Prev()
	case "down", "j":
		m.toolSelector.Next()
	case "enter":
		selectedTool := m.toolSelector.GetSelected()
		if selectedTool != "" {
			// Check if we're in delete mode
			if m.toolSelector.IsDeleteMode() {
				// Remove the tool proficiency using benefit remover
				// Find and remove from ToolProficiencies
				for i, tool := range m.character.ToolProficiencies {
					if tool == selectedTool {
						m.character.ToolProficiencies = append(m.character.ToolProficiencies[:i], m.character.ToolProficiencies[i+1:]...)
						break
					}
				}

				// Also remove from BenefitTracker (all sources that granted this tool)
				// This is more complex - we need to iterate through all benefits
				allBenefits := m.character.BenefitTracker.Benefits
				for _, benefit := range allBenefits {
					if benefit.Type == models.BenefitTool && benefit.Target == selectedTool {
						m.character.BenefitTracker.RemoveBenefitsBySource(benefit.Source.Type, benefit.Source.Name)
						break
					}
				}

				m.message = fmt.Sprintf("Tool proficiency removed: %s", selectedTool)
				m.storage.Save(m.character)
				m.toolSelector.Hide()
			} else {
				// Add mode: Add tool proficiency directly (not from origin)
				// We'll add it as a "manual" benefit (or Student of War)
				sourceType := "manual"
				sourceName := "Tool Proficiency"
				if !m.studentOfWarToolSelected && m.character.HasFeature("Student of War") {
					sourceType = "subclass_feature"
					sourceName = "Student of War"
				}

				source := models.BenefitSource{Type: sourceType, Name: sourceName}
				applier := models.NewBenefitApplier(m.character)
				applier.AddToolProficiency(source, selectedTool)

				m.toolSelector.Hide()

				// Check if we need to prompt for Fighter skill (Student of War)
				if !m.studentOfWarToolSelected && m.character.HasFeature("Student of War") {
					m.studentOfWarToolSelected = true

					// Prompt for Fighter skill selection (use standard skill selector)
					m.skillSelector.Show()
					m.message = "Select a skill from the Fighter skill list for Student of War..."
					return m, nil
				}

				m.message = fmt.Sprintf("Tool proficiency learned: %s!", selectedTool)
				m.storage.Save(m.character)
			}
		}
	case "esc":
		m.toolSelector.Hide()
		m.message = "Tool selection cancelled"
	}
	return m, nil
}

// handleWeaponMasterySelectorKeys handles weapon mastery selector specific keys
func (m *Model) handleWeaponMasterySelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.weaponMasterySelector, cmd = m.weaponMasterySelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case " ":
		// Check if toggle was successful and provide feedback
		if len(m.weaponMasterySelector.GetSelectedWeapons()) >= m.getWeaponMasteryCount() {
			selectedCount := len(m.weaponMasterySelector.GetSelectedWeapons())
			maxCount := m.getWeaponMasteryCount()
			if selectedCount > maxCount {
				m.message = "Maximum weapons already selected"
			}
		}
	case "enter":
		// Confirm selection
		if m.weaponMasterySelector.CanConfirm() {
			m.character.MasteredWeapons = m.weaponMasterySelector.GetSelectedWeapons()
			m.weaponMasterySelector.Hide()

			// Clear pending changes and complete class setup
			m.pendingChanges.Clear()
			m.storage.Save(m.character)

			debug.Log("Weapon mastery selection complete: %v", m.character.MasteredWeapons)
			m.message = fmt.Sprintf("Weapon mastery complete! Mastered %d weapons. Class setup complete. (HP: %d/%d)",
				len(m.character.MasteredWeapons), m.character.CurrentHP, m.character.MaxHP)
		} else {
			m.message = "Please select at least one weapon"
		}
	case "esc":
		// Cancel
		m.weaponMasterySelector.Hide()
		m.message = "Cancelled"
	}

	return m, cmd
}

// handleExpertiseSelectorKeys handles expertise selector specific keys
func (m *Model) handleExpertiseSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate navigation and selection to the component's Update method
	cmd := m.expertiseSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case " ":
		// Check if toggle was successful and provide feedback
		selectedCount := len(m.expertiseSelector.GetSelectedSkills())
		maxCount := m.expertiseSelector.GetMaxExpertise()
		if selectedCount > maxCount {
			m.message = "Maximum skills already selected"
		}
	case "enter":
		// Confirm selection
		if m.expertiseSelector.CanConfirm() {
			m.expertiseSelector.ApplySelections()
			m.expertiseSelector.Hide()

			// Clear pending changes and complete class setup
			m.pendingChanges.Clear()
			m.storage.Save(m.character)

			debug.Log("Expertise selection complete: %d skills selected", len(m.expertiseSelector.GetSelectedSkills()))
			m.message = fmt.Sprintf("Expertise complete! Selected %d skills. Class setup complete. (HP: %d/%d)",
				len(m.expertiseSelector.GetSelectedSkills()), m.character.CurrentHP, m.character.MaxHP)
		} else {
			selectedCount := len(m.expertiseSelector.GetSelectedSkills())
			maxCount := m.expertiseSelector.GetMaxExpertise()
			m.message = fmt.Sprintf("Please select %d skill(s) for expertise (%d/%d selected)", maxCount, selectedCount, maxCount)
		}
	case "esc":
		// Cancel selection
		m.expertiseSelector.Hide()
		m.message = "Expertise selection cancelled"
	}
	return m, cmd
}

// handleManeuverSelectorKeys handles maneuver selector specific keys
func (m *Model) handleManeuverSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.maneuverSelector, cmd = m.maneuverSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		// Confirm selection
		if m.maneuverSelector.CanConfirm() {
			// Check if this is initial selection (maneuvers array was empty)
			isInitialSelection := len(m.character.Maneuvers) == 0

			m.character.Maneuvers = m.maneuverSelector.GetSelectedManeuvers()
			m.maneuverSelector.Hide()

			debug.Log("Maneuver selection complete: %v", m.character.Maneuvers)

			// Only trigger Student of War flow on initial selection, not when editing
			if isInitialSelection && m.character.HasFeature("Student of War") && !m.studentOfWarToolSelected {
				// Set flag for Student of War flow
				m.studentOfWarToolSelected = false

				// Prompt for artisan tool selection (use standard tool selector)
				m.toolSelector.Show()
				m.message = "Select an artisan tool for Student of War..."
				return m, cmd
			}

			m.message = fmt.Sprintf("Battle Master maneuvers learned: %v", m.character.Maneuvers)
			m.storage.Save(m.character)
		} else {
			maneuversNeeded := 3 // Default for level 3
			if feature := m.character.GetFeature("Combat Superiority"); feature != nil && feature.Mechanics != nil {
				if count, ok := feature.Mechanics["maneuvers_known"].(float64); ok {
					maneuversNeeded = int(count)
				}
			}
			m.message = fmt.Sprintf("Please select %d maneuvers", maneuversNeeded)
		}
	case "esc":
		// Cancel
		m.maneuverSelector.Hide()
		m.message = "Maneuver selection cancelled"
	}

	return m, cmd
}

// handleItemSelectorKeys handles item selector specific keys
func (m *Model) handleItemSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle the key once and store the result
	cmd := m.itemSelector.HandleKey(msg)

	// Check if we're confirming quantity entry (only in quantity mode)
	if msg.String() == "enter" && m.itemSelector.IsInQuantityMode() {
		selectedDef, quantity := m.itemSelector.GetSelectedItem()
		if selectedDef != nil {
			// Convert to inventory item and add
			item := models.ConvertToInventoryItem(*selectedDef, quantity)
			m.character.Inventory.AddItem(item)
			m.message = fmt.Sprintf("Added %dx %s to inventory", quantity, item.Name)
			m.storage.Save(m.character)
			m.itemSelector.Hide()
		}
	}

	return m, cmd
}
