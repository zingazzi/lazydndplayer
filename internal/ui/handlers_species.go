// internal/ui/handlers_species.go
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// handleSpeciesSelectorKeys handles species selector specific keys
func (m *Model) handleSpeciesSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleSpeciesSelectorKeys: key=%s", msg.String())
	switch msg.String() {
	case "up", "k":
		m.speciesSelector.Prev()
	case "down", "j":
		m.speciesSelector.Next()
	case "enter":
		selectedSpecies := m.speciesSelector.GetSelectedSpecies()
		debug.Log("handleSpeciesSelectorKeys: Enter pressed, selectedSpecies=%v", selectedSpecies)
		if selectedSpecies != nil {
			debug.Log("handleSpeciesSelectorKeys: Selected species name=%s, HasSubtypes=%v", selectedSpecies.Name, selectedSpecies.HasSubtypes)
			// Check if species has subtypes
			if selectedSpecies.HasSubtypes && len(selectedSpecies.Subtypes) > 0 {
				// Show subtype selector
				subtypes := make([]components.SpeciesSubtype, len(selectedSpecies.Subtypes))
				for i, st := range selectedSpecies.Subtypes {
					subtypes[i] = components.SpeciesSubtype{
						Name:        st.Name,
						Description: st.Description,
						Modifier:    st.Modifier,
					}
					// For Dragonborn, show damage type
					if st.DamageType != "" {
						subtypes[i].Modifier = fmt.Sprintf("%s damage", st.DamageType)
					}
				}
				m.subtypeSelector.Show(selectedSpecies.Name, subtypes)
				m.message = fmt.Sprintf("Select %s subtype...", selectedSpecies.Name)
				m.speciesSelector.Hide()
				return m, nil
			}

			// No subtypes, apply species directly
			debug.Log("handleSpeciesSelectorKeys: No subtypes, applying species directly: %s", selectedSpecies.Name)
			oldSpecies := m.character.Race
			models.ApplySpeciesToCharacter(m.character, selectedSpecies.Name)
			debug.Log("handleSpeciesSelectorKeys: Species applied, new Race=%s", m.character.Race)

			// Check if we need to select additional languages
			needsLanguageSelection := false
			for _, lang := range m.character.Languages {
				if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
					needsLanguageSelection = true
					break
				}
			}

			// Check for various selections needed
			needsSkillSelection := models.HasSkillChoice(selectedSpecies)
			needsSpellSelection := models.HasSpellChoice(selectedSpecies)
			needsFeatSelection := models.HasFeatChoice(selectedSpecies)

			debug.Log("handleSpeciesSelectorKeys: Selection needs - language=%v, skill=%v, spell=%v, feat=%v",
				needsLanguageSelection, needsSkillSelection, needsSpellSelection, needsFeatSelection)

			if needsLanguageSelection {
				// Filter out languages the character already knows
				debug.Log("handleSpeciesSelectorKeys: Showing language selector")
				m.languageSelector.SetExcludeLanguages(m.character.Languages)
				m.languageSelector.Show()
				m.speciesSelector.Hide()
				m.message = "Select your additional language..."
				return m, nil
			} else if needsSkillSelection {
				debug.Log("handleSpeciesSelectorKeys: Showing skill selector")
				m.skillSelector.Show()
				m.speciesSelector.Hide()
				m.message = "Select your skill proficiency..."
				return m, nil
			} else if needsSpellSelection {
				// Show wizard cantrip selector
				debug.Log("handleSpeciesSelectorKeys: Showing spell selector")
				cantrips := models.GetWizardCantrips()
				m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
				m.spellSelector.Show()
				m.speciesSelector.Hide()
				m.message = "Select your wizard cantrip..."
				return m, nil
			} else if needsFeatSelection {
				// Show feat selector for origin feat
				debug.Log("handleSpeciesSelectorKeys: Showing feat selector")
				m.featSelector.Show(m.character, true)
				m.speciesSelector.Hide()
				m.message = "Select your origin feat..."
				return m, nil
			} else {
				// Species selection complete (no additional selections needed)
				debug.Log("handleSpeciesSelectorKeys: Species selection complete, no additional selections needed")
				if m.IsInWizard() {
					// In wizard mode, check if species setup is complete and advance
					debug.Log("handleSpeciesSelectorKeys: In wizard mode, checking if species setup complete")
					m.speciesSelector.Hide()
					m.checkAndAdvanceWizardAfterSpeciesSetup()
					return m, nil
				} else {
					m.message = fmt.Sprintf("Species changed from %s to %s. Speed updated to %d ft.", oldSpecies, selectedSpecies.Name, m.character.Speed)
					// Save character when species change is complete (no additional selections)
					m.storage.Save(m.character)
					m.speciesSelector.Hide()
				}
			}
		} else {
			debug.Log("handleSpeciesSelectorKeys: ERROR - selectedSpecies is nil!")
			m.message = "Error: No species selected"
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			// In wizard, ESC cancels the entire wizard
			m.cancelWizard()
			return m, nil
		}
		m.speciesSelector.Hide()
		m.message = "Species selection cancelled"
	}
	return m, nil
}

// handleSubtypeSelectorKeys handles subtype selector specific keys
func (m *Model) handleSubtypeSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.subtypeSelector.Prev()
	case "down", "j":
		m.subtypeSelector.Next()
	case "enter":
		selectedSubtype := m.subtypeSelector.GetSelectedSubtype()
		if selectedSubtype != nil {
			// ALWAYS use the species name from the selector (just selected)
			// Don't use m.character.Race as it might still have the old species
			speciesName := m.subtypeSelector.SpeciesName

			// Apply species with the selected subtype
			models.ApplySpeciesWithSubtype(m.character, speciesName, selectedSubtype.Name)

			// Get species info for additional checks
			species := models.GetSpeciesByName(speciesName)

			// Check if we need to select additional languages
			needsLanguageSelection := false
			for _, lang := range m.character.Languages {
				if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
					needsLanguageSelection = true
					break
				}
			}

			// Check for various selections needed
			needsSkillSelection := species != nil && models.HasSkillChoice(species)
			needsSpellSelection := species != nil && models.HasSpellChoice(species)
			needsFeatSelection := species != nil && models.HasFeatChoice(species)

			if needsLanguageSelection {
				// Filter out languages the character already knows
				m.languageSelector.SetExcludeLanguages(m.character.Languages)
				m.languageSelector.Show()
				m.message = "Select your additional language..."
			} else if needsSkillSelection {
				m.skillSelector.Show()
				m.message = "Select your skill proficiency..."
			} else if needsSpellSelection {
				// Show wizard cantrip selector (High Elf)
				cantrips := models.GetWizardCantrips()
				m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
				m.spellSelector.Show()
				m.message = "Select your wizard cantrip..."
			} else if needsFeatSelection {
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
			} else {
				// Species selection complete (no additional selections needed)
				if m.IsInWizard() {
					// In wizard mode, check if species setup is complete and advance
					m.checkAndAdvanceWizardAfterSpeciesSetup()
					return m, nil
				} else {
					m.message = fmt.Sprintf("%s (%s) selected! Speed: %d ft, Darkvision: %d ft",
						speciesName, selectedSubtype.Name, m.character.Speed, m.character.Darkvision)
					// Save character when species change is complete (no additional selections)
					m.storage.Save(m.character)
				}
			}
		}
		m.subtypeSelector.Hide()
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			// In wizard, ESC cancels the entire wizard
			m.cancelWizard()
			return m, nil
		}
		m.subtypeSelector.Hide()
		// Show species selector again to let user pick different species
		m.speciesSelector.Show()
		m.message = "Subtype selection cancelled"
	}
	return m, nil
}
