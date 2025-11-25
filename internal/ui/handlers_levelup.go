// internal/ui/handlers_levelup.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// handleLevelUpSelectorKeys handles level-up selector keys
func (m *Model) handleLevelUpSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for Deft Explorer/Scholar BEFORE processing keys - these need to intercept
	// Check if Deft Explorer skill/language selection is needed (Ranger level 2)
	if m.character.HasClass("Ranger") {
		rangerLevel := m.character.GetClassLevel("Ranger")
		if rangerLevel == 2 {
			// Check if Deft Explorer feature exists and needs skill/language selection
			for _, feature := range m.character.Features.Features {
				if feature.Name == "Deft Explorer" && feature.Mechanics != nil {
					if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "skill_choice" {
						// Check if we've already selected the skill
						skillCount := 1 // default
						if count, ok := feature.Mechanics["skill_count"].(float64); ok {
							skillCount = int(count)
						}
						langCount := 2 // default
						if count, ok := feature.Mechanics["languages"].(float64); ok {
							langCount = int(count)
						}

						// Check if skill was already selected (stored in feature mechanics)
						skillSelected := false
						if selected, ok := feature.Mechanics["skill_selected"].(bool); ok {
							skillSelected = selected
						}

						// Check how many languages were selected
						languagesSelected := 0
						if count, ok := feature.Mechanics["languages_selected"].(float64); ok {
							languagesSelected = int(count)
						}

						if !skillSelected {
							// Prompt for skill selection
							rangerSkills := []string{"Animal Handling", "Athletics", "Insight", "Investigation", "Nature", "Perception", "Stealth", "Survival"}
							debug.Log("Deft Explorer feature needs skill selection: %v", rangerSkills)
							m.levelUpSelector.Hide() // Hide level up selector so skill selector can receive keys
							m.classSkillSelector.Show("Ranger", rangerSkills, skillCount, m.character)
							m.message = "Select 1 skill from the Ranger skill list for Deft Explorer..."
							return m, nil
						} else if languagesSelected < langCount {
							// Prompt for language selection (multi-select mode)
							remainingCount := langCount - languagesSelected
							debug.Log("Deft Explorer feature needs language selection: %d remaining of %d total", remainingCount, langCount)
							m.levelUpSelector.Hide() // Hide level up selector so language selector can receive keys
							m.languageSelector.ShowForMultiSelect(m.character.Languages, remainingCount, true) // excludeCommon=true
							m.message = fmt.Sprintf("Select %d language(s) for Deft Explorer (Common excluded)...", remainingCount)
							return m, nil
						}
					}
				}
			}
		}
	}

	// Check if Scholar skill selection is needed (Wizard level 2)
	if m.character.HasClass("Wizard") {
		wizardLevel := m.character.GetClassLevel("Wizard")
		if wizardLevel == 2 {
			// Check if Scholar feature exists and needs skill selection
			for _, feature := range m.character.Features.Features {
				if feature.Name == "Scholar" && feature.Mechanics != nil {
					if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "skill_choice" {
						if skillList, ok := feature.Mechanics["skill_list"].([]interface{}); ok {
							skills := []string{}
							for _, s := range skillList {
								if skillStr, ok := s.(string); ok {
									skills = append(skills, skillStr)
								}
							}
							debug.Log("Scholar feature needs skill selection: %v", skills)
							m.levelUpSelector.Hide() // Hide level up selector so skill selector can receive keys
							m.classSkillSelector.Show("Wizard", skills, 1, m.character)
							m.message = "Select 1 skill for Scholar feature..."
							return m, nil
						}
					}
				}
			}
		}
	}

	updated, cmd := m.levelUpSelector.Update(msg)
	m.levelUpSelector = &updated

	// Save character if level-up is complete
	if !m.levelUpSelector.IsVisible() {
		m.storage.Save(m.character)
		m.character.UpdateDerivedStats()

		// Check if Battle Master needs maneuver selection
		if m.levelUpSelector.NeedsManeuverSelection {
			debug.Log("Battle Master detected - prompting for maneuver selection")
			m.levelUpSelector.NeedsManeuverSelection = false // Clear flag
			m.maneuverSelector.Show(3) // Battle Master starts with 3 maneuvers
			m.message = "Select 3 maneuvers for Battle Master..."
			return m, cmd
		}

		// Check if Eldritch Knight needs cantrip selection
		if m.levelUpSelector.NeedsCantripSelection {
			debug.Log("Eldritch Knight detected - prompting for cantrip selection")
			m.levelUpSelector.NeedsCantripSelection = false // Clear flag
			m.cantripSelector.Show("Wizard", 2) // Eldritch Knight starts with 2 cantrips
			m.message = "Select 2 cantrips from the Wizard spell list..."
			return m, cmd
		}

		// Check if Beast Master needs beast selection
		if m.levelUpSelector.NeedsBeastSelection {
			debug.Log("Beast Master detected - prompting for beast selection")
			m.levelUpSelector.NeedsBeastSelection = false // Clear flag
			m.levelUpSelector.Hide() // Hide level up selector so beast selector can receive keys
			m.beastSelector.Show()
			m.message = "Select your beast companion..."
			return m, cmd
		}

		// Check if Fey Wanderer needs Fey Gift selection
		if m.levelUpSelector.NeedsFeyGiftSelection {
			debug.Log("Fey Wanderer detected - prompting for Fey Gift selection")
			m.levelUpSelector.NeedsFeyGiftSelection = false // Clear flag
			m.levelUpSelector.Hide() // Hide level up selector so option selector can receive keys
			feyGiftOptions := []string{
				"Illusory butterflies flutter around you while you take a short or long rest",
				"Fresh, seasonal flowers sprout from your hair each dawn",
				"You faintly smell of cinnamon, lavender, nutmeg, or another comforting herb or spice",
				"Your shadow dances while no one is looking directly at it",
				"Horns or antlers sprout from your head",
				"Your skin and hair change color to match the season at each dawn",
			}
			m.optionSelector.Show("Select Fey Gift", feyGiftOptions)
			m.message = "Select your Fey Gift..."
			return m, cmd
		}

		// Check if Fey Wanderer needs Otherworldly Glamour skill selection
		if m.levelUpSelector.NeedsGlamourSkillSelection {
			debug.Log("Fey Wanderer detected - prompting for Otherworldly Glamour skill selection")
			m.levelUpSelector.NeedsGlamourSkillSelection = false // Clear flag
			m.levelUpSelector.Hide() // Hide level up selector so skill selector can receive keys
			glamourSkills := []string{"Deception", "Performance", "Persuasion"}
			m.classSkillSelector.Show("Ranger", glamourSkills, 1, m.character)
			m.message = "Select 1 skill for Otherworldly Glamour (Deception, Performance, or Persuasion)..."
			return m, cmd
		}

		// Check if Hunter needs Hunter's Prey selection
		if m.levelUpSelector.NeedsHunterPreySelection {
			debug.Log("Hunter detected - prompting for Hunter's Prey selection")
			m.levelUpSelector.NeedsHunterPreySelection = false // Clear flag
			m.levelUpSelector.Hide() // Hide level up selector so option selector can receive keys
			hunterPreyOptions := []string{
				"Colossus Slayer",
				"Giant Killer",
				"Horde Breaker",
			}
			m.optionSelector.Show("Select Hunter's Prey", hunterPreyOptions)
			m.message = "Select your Hunter's Prey option..."
			return m, cmd
		}

		// Check if Savant spell selection is needed (Wizard level 3 with subclass)
		if m.character.HasClass("Wizard") {
			wizardLevel := m.character.GetClassLevel("Wizard")
			if wizardLevel == 3 {
				debug.Log("Wizard level 3 detected - checking for Savant spell selection")
				// Check for Savant features that require spell selection
				for _, feature := range m.character.Features.Features {
					if feature.Mechanics != nil {
						if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "spell_selection" {
							// Get school, count, and maxLevel from mechanics
							school, _ := feature.Mechanics["spell_school"].(string)
							spellCount := 2 // default
							if count, ok := feature.Mechanics["spell_count"].(float64); ok {
								spellCount = int(count)
							}
							maxLevel := 2 // default
							if level, ok := feature.Mechanics["max_spell_level"].(float64); ok {
								maxLevel = int(level)
							}

							debug.Log("Triggering school spell selector: school=%s, count=%d, maxLevel=%d", school, spellCount, maxLevel)
							m.schoolSpellSelector.Show(school, maxLevel, spellCount)
							m.message = fmt.Sprintf("Select %d %s spells for your spellbook...", spellCount, school)
							return m, cmd
						}
					}
				}
			}
		}

		// Wizard spell learning - show message popup for levels 4+
		if m.character.HasClass("Wizard") {
			wizardLevel := m.character.GetClassLevel("Wizard")
			if wizardLevel > 3 {
				// Level 4+: Show message popup
				maxSpellLevel := (wizardLevel + 1) / 2
				if maxSpellLevel > 9 {
					maxSpellLevel = 9
				}
				m.messagePopup.Show("Wizard Spellbook", fmt.Sprintf("Add 2 spells (up to level %d) to your spellbook.\n\nPress 'v' in the Spells panel to open your spellbook.", maxSpellLevel))
				return m, cmd
			} else {
				m.message = "Character updated!"
			}
		} else {
			m.message = "Character updated!"
		}

		// Ensure focusArea is set to FocusMain after level-up completes
		m.focusArea = FocusMain
	}

	return m, cmd
}

// handleDeLevelSelectorKeys handles de-level selector keys
func (m *Model) handleDeLevelSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleDeLevelSelectorKeys: key=%s, showConfirm=%v", msg.String(), m.deLevelSelector.IsShowingConfirmation())

	switch msg.String() {
	case "up", "k":
		if !m.deLevelSelector.IsShowingConfirmation() {
			m.deLevelSelector.Prev()
		}
	case "down", "j":
		if !m.deLevelSelector.IsShowingConfirmation() {
			m.deLevelSelector.Next()
		}
	case "enter":
		if m.deLevelSelector.IsShowingConfirmation() {
			// Second enter - execute de-level (already done in ShowConfirmation)
			m.deLevelSelector.ConfirmDeLevel()
			m.storage.Save(m.character)
			m.character.UpdateDerivedStats()

			// Show result message
			previewResult := m.deLevelSelector.GetPreviewResult()
			if previewResult != nil {
				if previewResult.ClassRemoved {
					m.message = fmt.Sprintf("%s class removed entirely (was level %d). Total level: %d",
						previewResult.ClassName, previewResult.OldClassLevel, previewResult.NewTotalLevel)
				} else {
					m.message = fmt.Sprintf("%s level reduced from %d to %d. Total level: %d",
						previewResult.ClassName, previewResult.OldClassLevel, previewResult.NewClassLevel, previewResult.NewTotalLevel)
				}
			} else {
				m.message = "De-level complete!"
			}
		} else {
			// First enter - show confirmation
			err := m.deLevelSelector.ShowConfirmation()
			if err != nil {
				m.message = fmt.Sprintf("Error: %v", err)
				m.deLevelSelector.Hide()
			} else {
				m.message = "Confirm de-level (Enter) or cancel (Esc)"
			}
		}
	case "esc":
		if m.deLevelSelector.IsShowingConfirmation() {
			// Cancel confirmation, go back to class list
			m.deLevelSelector.CancelConfirmation()
			m.message = "De-level cancelled"
		} else {
			// Cancel entire de-level process
			m.deLevelSelector.Hide()
			m.message = "De-level cancelled"
		}
	}

	return m, nil
}
