// internal/ui/handlers_spell.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// handleCantripSelectorKeys handles cantrip selector specific keys
func (m *Model) handleCantripSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleCantripSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.cantripSelector, cmd = m.cantripSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		if m.cantripSelector.CanConfirm() {
			selectedCantrips := m.cantripSelector.GetSelectedCantrips()
			debug.Log("Cantrip selector confirmed: %d cantrips selected", len(selectedCantrips))

			// Update character's cantrips
			m.character.SpellBook.Cantrips = make([]string, len(selectedCantrips))
			copy(m.character.SpellBook.Cantrips, selectedCantrips)

			m.cantripSelector.Hide()

			// Check if we need weapon mastery selection next
			masteryCount := m.getWeaponMasteryCount()
			debug.Log("After cantrip selection, checking weapon mastery: count=%d", masteryCount)

			if masteryCount > 0 {
				// Show weapon mastery selector
				debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
				m.weaponMasterySelector.Show(masteryCount)
				m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
			} else {
				// Complete class selection
				debug.Log("Saving character and completing class selection")
				m.pendingChanges.Clear() // Clear backup on successful completion
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Cantrip selection complete! (Total: %d cantrips)", len(selectedCantrips))
			}
		} else {
			selectedCount := m.cantripSelector.GetSelectedCount()
			maxCount := m.cantripSelector.GetMaxCantrips()
			m.message = fmt.Sprintf("Please select %d cantrip(s) (%d/%d selected)", maxCount, selectedCount, maxCount)
		}
	case "esc":
		debug.Log("Cantrip selection cancelled, rolling back changes")
		m.cantripSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.storage.Save(m.character)
		m.message = "Cantrip selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSpellSelectorKeys handles spell selector specific keys
func (m *Model) handleSpellSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.spellSelector.Prev()
	case "down", "j":
		m.spellSelector.Next()
	case "enter":
		selectedSpell := m.spellSelector.GetSelectedSpell()
		if selectedSpell.Name != "" {
			// Check if we're in Eldritch Knight spell selection mode
			if m.character.IsEldritchKnight() && m.eldritchKnightSpellsSelected < 3 {
				// Check if spell is already selected
				alreadySelected := false
				for _, spell := range m.eldritchKnightSpells {
					if spell.Name == selectedSpell.Name {
						alreadySelected = true
						break
					}
				}

				if !alreadySelected {
					m.eldritchKnightSpells = append(m.eldritchKnightSpells, selectedSpell)
					m.eldritchKnightSpellsSelected++
					debug.Log("Eldritch Knight spell %d/3 selected: %s (School: %s)", m.eldritchKnightSpellsSelected, selectedSpell.Name, selectedSpell.School)

					if m.eldritchKnightSpellsSelected < 3 {
						// Continue selecting more spells
						m.spellSelector.Hide()

						// Load all Wizard level 1 spells again
						allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
						if err != nil {
							debug.Log("Error loading spells: %v", err)
							m.message = "Error loading spells"
							return m, nil
						}

						// Filter to Wizard, level 1 spells
						wizardSpells := []models.Spell{}
						for _, spell := range allSpells {
							if spell.Level == 1 {
								for _, class := range spell.Classes {
									if class == "Wizard" {
										wizardSpells = append(wizardSpells, spell)
										break
									}
								}
							}
						}

						m.spellSelector.SetSpells(wizardSpells, fmt.Sprintf("SELECT SPELL %d/3 (Wizard Level 1)", m.eldritchKnightSpellsSelected+1))
						m.spellSelector.Show()
						m.message = fmt.Sprintf("Spell %d/3 selected: %s", m.eldritchKnightSpellsSelected, selectedSpell.Name)
						return m, nil
					} else {
						// All 3 spells selected, validate and add
						m.spellSelector.Hide()

						// Count Abjuration/Evocation spells
						abjEvocCount := 0
						for _, spell := range m.eldritchKnightSpells {
							if spell.School == "Abjuration" || spell.School == "Evocation" {
								abjEvocCount++
							}
						}

						if abjEvocCount < 2 {
							m.message = fmt.Sprintf("ERROR: Need at least 2 Abjuration/Evocation spells (found %d). Please restart.", abjEvocCount)
							m.eldritchKnightSpells = []models.Spell{}
							m.eldritchKnightSpellsSelected = 0
							return m, nil
						}

						// Add all spells to spellbook
						for _, spell := range m.eldritchKnightSpells {
							m.character.SpellBook.AddSpell(spell)
							debug.Log("Added Eldritch Knight spell: %s", spell.Name)
						}

						// Clear temporary state
						m.eldritchKnightSpells = []models.Spell{}
						m.eldritchKnightSpellsSelected = 0

						// Check for weapon mastery
						masteryCount := m.getWeaponMasteryCount()
						debug.Log("After EK spells, checking weapon mastery: count=%d", masteryCount)

						if masteryCount > 0 {
							// Show weapon mastery selector
							debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
							m.weaponMasterySelector.Show(masteryCount)
							m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
						} else {
							// Complete class selection
							debug.Log("Saving character and completing class selection")
							m.pendingChanges.Clear() // Clear backup on successful completion
							m.storage.Save(m.character)
							m.message = "Eldritch Knight setup complete! 3 spells learned."
						}
						return m, nil
					}
				} else {
					m.message = fmt.Sprintf("You already selected %s", selectedSpell.Name)
					return m, nil
				}
			}

			// Normal spell selection (for species/origin)
			// Check if character already has this spell
			hasSpell := false
			for _, existing := range m.character.SpellBook.Spells {
				if existing.Name == selectedSpell.Name {
					hasSpell = true
					break
				}
			}

			if !hasSpell {
				// Add spell to spellbook and track it as a species spell
				m.character.SpellBook.AddSpell(selectedSpell)
				m.character.SpeciesSpells = append(m.character.SpeciesSpells, selectedSpell.Name)
				m.message = fmt.Sprintf("Spell learned: %s", selectedSpell.Name)
			} else {
				m.message = fmt.Sprintf("You already know %s", selectedSpell.Name)
			}

			// After spell selection, check if we need feat selection
			species := models.GetSpeciesByName(m.character.Race)
			if species != nil && models.HasFeatChoice(species) {
				m.spellSelector.Hide()
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
				return m, nil
			}

			// Regular spell selection (for High Elf wizard cantrip)
			if selectedSpell.Level == 0 {
				m.character.SpellBook.Cantrips = append(m.character.SpellBook.Cantrips, selectedSpell.Name)
				m.storage.Save(m.character)
				m.spellSelector.Hide()

				// Check if we need feat selection
				if species != nil && models.HasFeatChoice(species) {
					m.featSelector.Show(m.character, true)
					m.message = "Select your origin feat..."
				} else {
					m.message = fmt.Sprintf("Wizard cantrip learned: %s", selectedSpell.Name)
				}
			} else {
				// Save character after spell selection (final step)
				m.storage.Save(m.character)
			}
		}
		m.spellSelector.Hide()
	case "esc":
		// Check if we're cancelling Eldritch Knight spell selection
		if m.character.IsEldritchKnight() && m.eldritchKnightSpellsSelected > 0 {
			debug.Log("Eldritch Knight spell selection cancelled, rolling back changes")
			m.eldritchKnightSpells = []models.Spell{}
			m.eldritchKnightSpellsSelected = 0
			m.spellSelector.Hide()
			m.pendingChanges.RestoreClass(m.character)
			m.storage.Save(m.character)
			m.message = "Eldritch Knight selection cancelled - restored previous state"
		} else {
			m.spellSelector.Hide()
			m.message = "Spell selection cancelled"
		}
	}
	return m, nil
}

// handleLeveledSpellSelectorKeys handles keyboard input for the leveled spell selector
func (m *Model) handleLeveledSpellSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleLeveledSpellSelectorKeys: key=%s", msg.String())

	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.leveledSpellSelector, cmd = m.leveledSpellSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		if m.leveledSpellSelector.GetRemainingCount() == 0 {
			selectedSpells := m.leveledSpellSelector.GetSelectedSpells()
			debug.Log("=== SPELL SELECTOR CONFIRMED: Selected spells: %v", selectedSpells)

			// Add spells to character's spellbook
			addedCount := 0
			for _, spellName := range selectedSpells {
				// Load spell details and add to spellbook
				allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
				if err != nil {
					debug.Log("=== ERROR loading spells: %v", err)
					m.message = "Error loading spells"
					return m, cmd
				}

				for _, spell := range allSpells {
					if spell.Name == spellName {
						// Set spell as known but not prepared
						spell.Known = true
						spell.Prepared = false
						m.character.SpellBook.Spells = append(m.character.SpellBook.Spells, spell)
						addedCount++
						debug.Log("=== Added spell to spellbook: %s (Known=%v, Prepared=%v)", spell.Name, spell.Known, spell.Prepared)
						break
					}
				}
			}

			debug.Log("=== Total spells in spellbook after adding: %d", len(m.character.SpellBook.Spells))

			m.leveledSpellSelector.Hide()

			// Save character
			m.pendingChanges.Clear()
			m.storage.Save(m.character)
			m.message = fmt.Sprintf("Spell selection complete! Added %d spells to spellbook (Total: %d)", addedCount, len(m.character.SpellBook.Spells))
		} else {
			needed := m.leveledSpellSelector.GetRemainingCount()
			m.message = fmt.Sprintf("Please select %d more spell(s)", needed)
		}
	case "esc":
		// Cancel - rollback to previous state
		debug.Log("Spell selection cancelled, rolling back changes")
		m.leveledSpellSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.storage.Save(m.character)
		m.message = "Spell selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSchoolSpellSelectorKeys handles keyboard input for the school spell selector
func (m *Model) handleSchoolSpellSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleSchoolSpellSelectorKeys: key=%s", msg.String())

	// Delegate navigation and selection to the component's Update method
	cmd := m.schoolSpellSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		if m.schoolSpellSelector.GetRemainingCount() == 0 {
			selectedSpells := m.schoolSpellSelector.GetSelectedSpells()
			debug.Log("=== SCHOOL SPELL SELECTOR CONFIRMED: Selected spells: %v", len(selectedSpells))

			// Add spells to character's spellbook
			addedCount := 0
			for _, spell := range selectedSpells {
				// Set spell as known but not prepared
				spell.Known = true
				spell.Prepared = false
				m.character.SpellBook.Spells = append(m.character.SpellBook.Spells, spell)
				addedCount++
				debug.Log("=== Added spell to spellbook: %s (Level %d, School: %s)", spell.Name, spell.Level, spell.School)
			}

			m.schoolSpellSelector.Hide()

			// Save character
			m.storage.Save(m.character)

			// Show message popup about adding 2 more spells
			debug.Log("=== SHOWING MESSAGE POPUP after Savant selection")
			m.messagePopup.Show("Wizard Spellbook", "Add 2 spells to your spellbook.\n\nPress 'v' in the Spells panel to open your spellbook.")
			debug.Log("=== Message popup visible: %v", m.messagePopup.IsVisible())
			return m, cmd
		} else {
			needed := m.schoolSpellSelector.GetRemainingCount()
			m.message = fmt.Sprintf("Please select %d more spell(s)", needed)
		}
	case "esc":
		debug.Log("School spell selection cancelled")
		m.schoolSpellSelector.Hide()
		m.message = "Spell selection cancelled"
	}

	return m, cmd
}

// handleSpellbookEditorKeys handles keyboard input for the spellbook editor
func (m *Model) handleSpellbookEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleSpellbookEditorKeys: key=%s", msg.String())

	// Delegate to the component's Update method
	var cmd tea.Cmd
	*m.spellbookEditor, cmd = m.spellbookEditor.Update(tea.KeyMsg(msg))

	// Save character after any changes
	if msg.String() == " " || msg.String() == "a" {
		m.storage.Save(m.character)
		preparedCount := 0
		for _, spell := range m.character.SpellBook.Spells {
			if spell.Prepared && spell.Level > 0 {
				preparedCount++
			}
		}
		m.message = fmt.Sprintf("Spellbook updated (Prepared: %d/%d)", preparedCount, m.character.SpellBook.MaxPreparedSpells)
	}

	if msg.String() == "esc" {
		m.spellbookEditor.Hide()
		m.storage.Save(m.character)
		m.message = "Spellbook closed"
	}

	return m, cmd
}

// handleSpellPrepSelectorKeys handles keyboard input for the spell prep selector
func (m *Model) handleSpellPrepSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	*m.spellPrepSelector, cmd = m.spellPrepSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		m.spellPrepSelector.Hide()
		m.storage.Save(m.character)
		m.message = "Spell preparation updated!"
	case "esc":
		m.spellPrepSelector.Hide()
		m.message = "Spell preparation cancelled"
	}

	return m, cmd
}

// handleSlotRestorerKeys handles keyboard input for the slot restorer
func (m *Model) handleSlotRestorerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	*m.slotRestorer, cmd = m.slotRestorer.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter", " ":
		// Restore one slot
		restoredMsg, success := m.slotRestorer.RestoreSlot()
		if success {
			m.storage.Save(m.character)
			m.message = restoredMsg
		} else {
			m.message = restoredMsg
		}
	case "esc":
		m.slotRestorer.Hide()
		m.message = "Slot restoration cancelled"
	}

	return m, cmd
}
