// internal/ui/handlers_class.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// handleClassSelectorKeys handles class selector specific keys
func (m *Model) handleClassSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleClassSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.classSelector, cmd = m.classSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		selectedClassName := m.classSelector.GetSelectedClass()
		debug.Log("Class selector: enter pressed, selected=%s", selectedClassName)

		if selectedClassName != "" {
			// Get the full class data to check skill choices
			classData := models.GetClassByName(selectedClassName)
			debug.Log("Class data loaded: %v (nil=%v)", selectedClassName, classData == nil)

			if classData == nil {
				debug.Log("ERROR: Class %s not found!", selectedClassName)
				m.message = fmt.Sprintf("Error: Class %s not found", selectedClassName)
				m.classSelector.Hide()
				return m, nil
			}

			debug.Log("Class %s: SkillChoices=%v", selectedClassName, classData.SkillChoices)
			if classData.SkillChoices != nil {
				debug.Log("  Choose=%d, From=%v", classData.SkillChoices.Choose, classData.SkillChoices.From)
			}

			// Check if class has skill choices using service
			if m.classService.RequiresSkillChoice(selectedClassName) {
				// Show skill selector
				debug.Log("Showing skill selector for %s", selectedClassName)
				m.classSelector.Hide()
				from, choose := m.classService.GetSkillChoices(selectedClassName)
				m.classSkillSelector.Show(selectedClassName, from, choose, m.character)
				if m.IsInWizard() {
					m.message = fmt.Sprintf("Step 1/4: Class Selection - Select skills for %s class...", selectedClassName)
				} else {
					m.message = fmt.Sprintf("Select skills for %s class...", selectedClassName)
				}
			} else {
				// No skill choices, apply class directly using service
				debug.Log("No skill choices, applying class directly")
				msg, err := m.classService.ApplyClass(m.character, selectedClassName)
				if err != nil {
					debug.Log("ERROR applying class: %v", err)
					m.message = err.Error()
				} else {
					debug.Log("Class applied successfully: %s (HP: %d/%d)", selectedClassName, m.character.CurrentHP, m.character.MaxHP)
					m.message = msg
				}
				m.storage.Save(m.character)
				m.classSelector.Hide()

				// If in wizard mode, advance to next step
				if m.IsInWizard() {
					m.advanceWizardStep()
				}
			}
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		debug.Log("Class selector: cancelled - restoring previous state")
		// Restore previous class state
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.classSelector.Hide()
		m.storage.Save(m.character) // Save restored state
		m.message = "Class selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSubclassSelectorKeys handles subclass selector specific keys
func (m *Model) handleSubclassSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleSubclassSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.subclassSelector, cmd = m.subclassSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		selectedSubclass := m.subclassSelector.GetSelectedSubclass()
		if selectedSubclass != nil {
			debug.Log("Subclass selected: %s", selectedSubclass.Name)

			// Apply the subclass to the character's current class
			if len(m.character.Classes) > 0 {
				// Update the most recent class (should be the only one at level 1)
				m.character.Classes[len(m.character.Classes)-1].Subclass = selectedSubclass.Name
				debug.Log("Set character subclass to: %s", selectedSubclass.Name)

				// Grant subclass features for the current level
				className := m.character.Classes[len(m.character.Classes)-1].ClassName
				classLevel := m.character.Classes[len(m.character.Classes)-1].Level
				subclassFeatures := models.GrantSubclassFeatures(m.character, className, selectedSubclass.Name, classLevel)
				debug.Log("Granted %d subclass features: %v", len(subclassFeatures), subclassFeatures)

				m.subclassSelector.Hide()

				// Handle Fighter subclass-specific prompts
				debug.Log("Checking for Fighter subclass prompts: className='%s', subclass='%s'", className, selectedSubclass.Name)

				if className == "Fighter" {
					debug.Log("Fighter detected, checking subclass type")
					switch selectedSubclass.Name {
					case "Eldritch Knight":
						// Eldritch Knight: Prompt for 2 cantrips from Wizard list
						debug.Log("Eldritch Knight selected - prompting for cantrips")
						m.cantripSelector.Show("Wizard", 2)
						m.message = "Select 2 cantrips from the Wizard spell list..."
						return m, cmd
					case "Battle Master":
						// Battle Master: Prompt for 3 maneuvers
						debug.Log("Battle Master selected - prompting for maneuvers")
						m.maneuverSelector.Show(3)
						m.message = "Select 3 maneuvers for Battle Master..."
						return m, cmd
					default:
						debug.Log("Fighter subclass '%s' doesn't require special prompts", selectedSubclass.Name)
					}
				}

				// Handle Ranger Beast Master subclass-specific prompts
				if className == "Ranger" && selectedSubclass.Name == "Beast Master" {
					debug.Log("Beast Master selected - prompting for beast selection")
					m.beastSelector.Show()
					m.message = "Select your beast companion..."
					return m, cmd
				}

				// Handle Wizard subclass-specific prompts
				if className == "Wizard" {
					debug.Log("Wizard detected, checking for Savant spell selection")

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
				} else {
					debug.Log("Not a Fighter or Wizard, className='%s'", className)
				}

				// Check if we need cantrip selection next (for other spellcasters)
				classData := models.GetClassByName(m.character.Class)
			needsCantrips := classData != nil && classData.Spellcasting != nil && classData.Spellcasting.CantripsKnown > 0
			debug.Log("Class %s needs cantrips: %v", m.character.Class, needsCantrips)

			if needsCantrips {
				// Show cantrip selector for the class
				cantripCount := classData.Spellcasting.CantripsKnown
				debug.Log("Showing cantrip selector for %s: %d cantrips", m.character.Class, cantripCount)
				m.cantripSelector.Show(m.character.Class, cantripCount)
				m.message = fmt.Sprintf("Select %d cantrip(s) from the %s spell list...", cantripCount, m.character.Class)
				return m, cmd
			}

			// Check if we need weapon mastery selection
			masteryCount := getWeaponMasteryCount(m.character)
			debug.Log("After subclass, checking weapon mastery: count=%d", masteryCount)

				if masteryCount > 0 {
					// Show weapon mastery selector
					debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
					m.weaponMasterySelector.Show(masteryCount)
					m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
					return m, cmd
				}

				// Check if we need fighting style selection
				if m.character.HasClass("Fighter") || m.character.HasClass("Paladin") || m.character.HasClass("Ranger") {
					fighterLevel := m.character.GetClassLevel("Fighter")
					paladinLevel := m.character.GetClassLevel("Paladin")
					rangerLevel := m.character.GetClassLevel("Ranger")
					if fighterLevel == 1 || paladinLevel == 2 || rangerLevel == 2 {
						debug.Log("Showing fighting style selector")
						m.fightingStyleSelector.Show(m.character.Class)
						m.message = "Select your fighting style..."
						return m, cmd
					}
				}

				// Check if we need expertise selection (Rogue level 1)
				if m.character.HasClass("Rogue") {
					rogueLevel := m.character.GetClassLevel("Rogue")
					if rogueLevel == 1 {
						debug.Log("Rogue level 1 - showing expertise selector")
						m.expertiseSelector.Show(2) // 2 skills at level 1
						m.message = "Select 2 skills for expertise..."
						return m, cmd
					}
				}

				// Complete class selection
				debug.Log("Saving character and completing class selection")
				m.pendingChanges.Clear() // Clear backup on successful completion
				m.storage.Save(m.character)

				// If in wizard mode, check if we can advance
				if m.IsInWizard() {
					m.checkAndAdvanceWizardAfterClassSetup()
					// Return immediately to ensure wizard state is preserved
					return m, cmd
				} else {
					m.message = fmt.Sprintf("Subclass '%s' selected! Class setup complete. (HP: %d/%d)", selectedSubclass.Name, m.character.CurrentHP, m.character.MaxHP)
				}
			}
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		debug.Log("Subclass selector: cancelled - restoring previous state")
		m.subclassSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.storage.Save(m.character)
		m.message = "Subclass selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleBeastSelectorKeys handles beast selector specific keys
// Returns (model, cmd, handled) where handled indicates if the key was processed
func (m *Model) handleBeastSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	debug.Log("handleBeastSelectorKeys: key=%s", msg.String())

	var cmd tea.Cmd

	switch msg.String() {
	case "up", "k":
		m.beastSelector.Prev()
		return m, nil, true
	case "down", "j":
		m.beastSelector.Next()
		return m, nil, true
	case "enter":
		selectedType := m.beastSelector.GetSelectedType()
		debug.Log("Beast selector: enter pressed, selected=%s", selectedType)

		if selectedType != "" {
			// Get ranger level
			rangerLevel := m.character.GetClassLevel("Ranger")
			if rangerLevel == 0 {
				rangerLevel = 3 // Default for level 3 selection
			}

			// Check if this is from long rest (Beast Master companion already exists)
			// Store old companion before replacing
			oldBeastMasterCompanion := m.character.GetBeastMasterCompanion()
			wasChangingCompanion := oldBeastMasterCompanion != nil

			// Initialize companion
			companion := models.InitializeCompanion(selectedType, rangerLevel)

			// Remove old Beast Master companion if exists
			if oldBeastMasterCompanion != nil {
				for i := range m.character.Companions {
					if m.character.Companions[i].IsBeastMasterSpecial() {
						m.character.RemoveCompanion(i)
						break
					}
				}
			}

			// Add new companion
			m.character.AddCompanion(companion)

			debug.Log("Companion '%s' initialized for Ranger level %d", selectedType, rangerLevel)
			m.beastSelector.Hide()

			// If companion was already set, this is a change (from long rest or companion panel)
			if wasChangingCompanion {
				// Check if we're in a long rest
				if m.restPopup.IsVisible() && m.restPopup.GetRestType() == components.LongRestType {
					// This is a change during long rest - just update and return to rest popup
					m.storage.Save(m.character)
					m.message = fmt.Sprintf("Companion changed to '%s'. Press Enter to complete long rest.", selectedType)
				// Re-show rest popup
				m.restPopup.Show(components.LongRestType)
				return m, cmd, true
			} else {
				// This is a change from companion panel - just update and save
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Companion changed to '%s'!", selectedType)
				return m, cmd, true
			}
		}

		// Check if we're in a level-up flow - if so, complete level-up instead of class setup
		// During level-up, levelUpSelector is hidden when beast selector is shown, so we check:
		// 1. If levelUpSelector was visible (might still be)
		// 2. If character is Ranger level 3 with Beast Master and NOT in wizard mode
		//    (indicates level-up, not initial class selection)
		rangerClass := m.character.GetClassLevelStruct("Ranger")
		isLevelUp := m.levelUpSelector.IsVisible() ||
			(rangerClass != nil && rangerClass.Level == 3 &&
			 rangerClass.Subclass == "Beast Master" && !m.IsInWizard())

		if isLevelUp {
			debug.Log("Beast selection during level-up - completing level-up process")
			m.storage.Save(m.character)
			m.message = fmt.Sprintf("Beast companion '%s' selected! Level-up complete. (HP: %d/%d)", selectedType, m.character.CurrentHP, m.character.MaxHP)
			// Level-up completion will be handled by the level-up handler
			return m, cmd, true
		}

		// Continue with class setup flow (check for weapon mastery, fighting style, etc.)
		// This should only happen during initial class selection or wizard mode
		// Check if we need weapon mastery selection
		masteryCount := getWeaponMasteryCount(m.character)
		debug.Log("After beast selection, checking weapon mastery: count=%d", masteryCount)

		if masteryCount > 0 {
			// Show weapon mastery selector
			debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
			m.weaponMasterySelector.Show(masteryCount)
			m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
			return m, cmd, true
		}

		// Check if we need fighting style selection
		if m.character.HasClass("Fighter") || m.character.HasClass("Paladin") || m.character.HasClass("Ranger") {
			fighterLevel := m.character.GetClassLevel("Fighter")
			paladinLevel := m.character.GetClassLevel("Paladin")
			rangerLevel := m.character.GetClassLevel("Ranger")
			if fighterLevel == 1 || paladinLevel == 2 || rangerLevel == 2 {
				debug.Log("Showing fighting style selector")
				m.fightingStyleSelector.Show(m.character.Class)
				m.message = "Select your fighting style..."
				return m, cmd, true
			}
		}

		// Complete class selection
		debug.Log("Saving character and completing class selection")
		m.pendingChanges.Clear() // Clear backup on successful completion
		m.storage.Save(m.character)

		// If in wizard mode, check if we can advance
		if m.IsInWizard() {
			m.checkAndAdvanceWizardAfterClassSetup()
			// Return immediately to ensure wizard state is preserved
			return m, cmd, true
		} else {
			m.message = fmt.Sprintf("Beast companion '%s' selected! Class setup complete. (HP: %d/%d)", selectedType, m.character.CurrentHP, m.character.MaxHP)
		}
		return m, cmd, true
	}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil, true
		}

		debug.Log("Beast selector: cancelled - restoring previous state")
		m.beastSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.message = "Beast selection cancelled - restored previous state"
		return m, nil, true
	}

	return m, nil, false
}

// handleClassSkillSelectorKeys handles class skill selector specific keys
func (m *Model) handleClassSkillSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleClassSkillSelectorKeys: key=%s", msg.String())

	var cmd tea.Cmd

	switch msg.String() {
	case "up", "k":
		debug.Log("ClassSkillSelector: up/k pressed, calling Prev()")
		m.classSkillSelector.Prev()
		return m, nil
	case "down", "j":
		debug.Log("ClassSkillSelector: down/j pressed, calling Next()")
		m.classSkillSelector.Next()
		return m, nil
	case " ": // Space to toggle
		m.classSkillSelector.ToggleSkill()
		debug.Log("Skill selector: selected=%d/%d", len(m.classSkillSelector.SelectedSkills), m.classSkillSelector.MaxChoices)
		return m, nil
	case "enter":
		canConfirm := m.classSkillSelector.CanConfirm()
		debug.Log("Skill selector: enter pressed, canConfirm=%v", canConfirm)

		if canConfirm {
			selectedSkills := m.classSkillSelector.GetSelectedSkills()
			selectedClassName := m.classSkillSelector.ClassName
			debug.Log("Skill selector confirmed: class=%s, skills=%v", selectedClassName, selectedSkills)

			// Check if this is for a feature (Deft Explorer, Scholar, or Otherworldly Glamour) rather than class selection
			// If character already has this class, it's likely a feature selection
			isFeatureSelection := false
			var featureName string
			isLevelUpFlow := m.levelUpSelector.IsVisible() || m.levelUpSelector.NeedsGlamourSkillSelection

			if m.character.HasClass(selectedClassName) || isLevelUpFlow {
				// Check if we're selecting skills for Deft Explorer
				for _, feature := range m.character.Features.Features {
					if feature.Name == "Deft Explorer" && feature.Mechanics != nil {
						if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "skill_choice" {
							skillSelected, _ := feature.Mechanics["skill_selected"].(bool)
							if !skillSelected {
								isFeatureSelection = true
								featureName = "Deft Explorer"
								break
							}
						}
					}
					if feature.Name == "Scholar" && feature.Mechanics != nil {
						if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "skill_choice" {
							isFeatureSelection = true
							featureName = "Scholar"
							break
						}
					}
					// Check if we're selecting skills for Otherworldly Glamour (Fey Wanderer)
					if feature.Name == "Otherworldly Glamour" && feature.Mechanics != nil {
						if mechType, ok := feature.Mechanics["type"].(string); ok && mechType == "skill_choice" {
							// Check if skill_options exists (confirms this is Otherworldly Glamour)
							if _, ok := feature.Mechanics["skill_options"].([]interface{}); ok {
								// This is Otherworldly Glamour skill selection
								isFeatureSelection = true
								featureName = "Otherworldly Glamour"
								break
							}
						}
					}
				}

				// Also check if we're in level-up flow and selecting for Otherworldly Glamour
				if !isFeatureSelection && isLevelUpFlow && m.message == "Select 1 skill for Otherworldly Glamour (Deception, Performance, or Persuasion)..." {
					isFeatureSelection = true
					featureName = "Otherworldly Glamour"
				}
			}

			if isFeatureSelection {
				// Handle feature-based skill selection
				debug.Log("Feature-based skill selection for %s", featureName)
				for _, skillName := range selectedSkills {
					skillType := models.SkillType(skillName)
					skill := m.character.Skills.GetSkill(skillType)
					if skill != nil && skill.Proficiency == 0 {
						skill.Proficiency = 1 // Grant proficiency
						debug.Log("Granted proficiency in %s from %s", skillName, featureName)
					}
					// Track this skill as coming from class
					m.character.ClassSkills = append(m.character.ClassSkills, skillType)
				}

				// Mark skill as selected in feature mechanics
				for i := range m.character.Features.Features {
					if m.character.Features.Features[i].Name == featureName {
						if m.character.Features.Features[i].Mechanics == nil {
							m.character.Features.Features[i].Mechanics = make(map[string]interface{})
						}
						m.character.Features.Features[i].Mechanics["skill_selected"] = true
						break
					}
				}

				m.classSkillSelector.Hide()
				m.storage.Save(m.character)

				// For Otherworldly Glamour, complete level-up if in level-up flow
				if featureName == "Otherworldly Glamour" && isLevelUpFlow {
					debug.Log("Otherworldly Glamour skill selection complete - completing level-up")
					m.character.UpdateDerivedStats()
					m.focusArea = FocusMain
					m.message = fmt.Sprintf("Otherworldly Glamour skill selected! Level-up complete.")
					return m, cmd
				}

				// For Deft Explorer, check if we need language selection
				if featureName == "Deft Explorer" {
					// Find the feature to check language count
					for _, feature := range m.character.Features.Features {
						if feature.Name == "Deft Explorer" && feature.Mechanics != nil {
							langCount := 2 // default
							if count, ok := feature.Mechanics["languages"].(float64); ok {
								langCount = int(count)
							}
							languagesSelected := 0
							if count, ok := feature.Mechanics["languages_selected"].(float64); ok {
								languagesSelected = int(count)
							}

							if languagesSelected < langCount {
								// Prompt for language selection (multi-select mode)
								remainingCount := langCount - languagesSelected
								m.languageSelector.ShowForMultiSelect(m.character.Languages, remainingCount, true) // excludeCommon=true
								m.message = fmt.Sprintf("Select %d language(s) for Deft Explorer (Common excluded)...", remainingCount)
								return m, cmd
							}
						}
					}
				}

				// Check if we need fighting style selection (for Ranger level 2)
				if selectedClassName == "Ranger" {
					rangerLevel := m.character.GetClassLevel("Ranger")
					if rangerLevel == 2 {
						debug.Log("Ranger level 2 - checking for fighting style selection")
						m.levelUpSelector.Hide() // Hide level up selector so fighting style selector can receive keys
						m.fightingStyleSelector.Show("Ranger")
						m.message = "Select your fighting style..."
						return m, cmd
					}
				}

				m.message = fmt.Sprintf("%s skill selection complete!", featureName)
				return m, cmd
			}

			// Original class selection logic
			debug.Log("Applying class %s with skills: %v", selectedClassName, selectedSkills)

			// Apply the class first
			err := models.ApplyClassToCharacter(m.character, selectedClassName)
			if err != nil {
				debug.Log("ERROR applying class: %v", err)
				m.message = fmt.Sprintf("Error applying class: %v", err)
				m.classSkillSelector.Hide()
				return m, nil
			}
			debug.Log("Class %s applied successfully", selectedClassName)

			// Apply selected skills
			for _, skillName := range selectedSkills {
				skillType := models.SkillType(skillName)
				skill := m.character.Skills.GetSkill(skillType)
				if skill != nil && skill.Proficiency == 0 {
					skill.Proficiency = 1 // Grant proficiency
					debug.Log("Granted proficiency in %s", skillName)
				}
				// Track this skill as coming from class
				m.character.ClassSkills = append(m.character.ClassSkills, skillType)
				debug.Log("Tracked %s as class skill", skillName)
			}

			// Record choices for rollback
			debug.Log("Recording class choice: %s with skills %v", selectedClassName, selectedSkills)
			m.character.Choices.RecordClassChoice(selectedClassName, "")
			m.character.Choices.RecordLevelChoice(1, selectedSkills, "", []string{}, "", nil)

			m.classSkillSelector.Hide()

			// Check if we need subclass selection at level 1
			classData := models.GetClassByName(selectedClassName)
			if classData != nil && models.RequiresSubclassAtLevel(classData, 1) {
				debug.Log("Class requires subclass at level 1, showing subclass selector")
				m.subclassSelector.Show(selectedClassName, 1) // Level 1 subclass selection
				if m.IsInWizard() {
					m.message = fmt.Sprintf("Step 1/4: Class Selection - Select %s subclass...", selectedClassName)
				} else {
					m.message = fmt.Sprintf("Select %s subclass...", selectedClassName)
				}
				return m, cmd
			}

			// Check if we need Primal Order selection (Druid level 1)
			if selectedClassName == "Druid" {
				druidLevel := m.character.GetClassLevel("Druid")
				if druidLevel == 1 {
					// Check if Primal Order feature exists and hasn't been selected
					primalOrderApplied := false
					if m.character.BenefitTracker != nil {
						for _, benefit := range m.character.BenefitTracker.Benefits {
							if benefit.Source.Type == "class_feature" && benefit.Source.Name == "Primal Order" {
								primalOrderApplied = true
								break
							}
						}
					}
					if !primalOrderApplied {
						debug.Log("Druid level 1 - showing Primal Order selector")
						m.SetPrimalOrderSelectorVisible(true)
						m.message = "Select your Primal Order..."
						return m, cmd
					}
				}
			}

			// Check if we need cantrip selection
			if classData != nil && classData.Spellcasting != nil && classData.Spellcasting.CantripsKnown > 0 {
				cantripCount := classData.Spellcasting.CantripsKnown
				debug.Log("Showing cantrip selector for %s: %d cantrips", selectedClassName, cantripCount)
				m.cantripSelector.Show(selectedClassName, cantripCount) // Use class name as spell list
				m.message = fmt.Sprintf("Select %d cantrip(s) from the %s spell list...", cantripCount, selectedClassName)
				return m, cmd
			}

			// Check if we need weapon mastery selection
			masteryCount := getWeaponMasteryCount(m.character)
			debug.Log("After class skills, checking weapon mastery: count=%d", masteryCount)

			if masteryCount > 0 {
				// Show weapon mastery selector
				debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
				m.weaponMasterySelector.Show(masteryCount)
				m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
				return m, cmd
			}

			// Check if we need fighting style selection
			if selectedClassName == "Fighter" || selectedClassName == "Paladin" || selectedClassName == "Ranger" {
				fighterLevel := m.character.GetClassLevel("Fighter")
				paladinLevel := m.character.GetClassLevel("Paladin")
				rangerLevel := m.character.GetClassLevel("Ranger")
				if fighterLevel == 1 || paladinLevel == 2 || rangerLevel == 2 {
					debug.Log("Showing fighting style selector")
					m.fightingStyleSelector.Show(m.character.Class)
					m.message = "Select your fighting style..."
					return m, cmd
				}
			}

			// Check if we need expertise selection (Rogue level 1)
			if selectedClassName == "Rogue" {
				rogueLevel := m.character.GetClassLevel("Rogue")
				if rogueLevel == 1 {
					debug.Log("Rogue level 1 - showing expertise selector")
					m.expertiseSelector.Show(2) // 2 skills at level 1
					m.message = "Select 2 skills for expertise..."
					return m, cmd
				}
			}

			// Complete class selection
			debug.Log("Saving character and completing class selection")
			m.pendingChanges.Clear() // Clear backup on successful completion
			m.storage.Save(m.character)

			// If in wizard mode, advance to next step
			if m.IsInWizard() {
				m.advanceWizardStep()
			} else {
				m.message = fmt.Sprintf("Class %s selected! Class setup complete. (HP: %d/%d)", selectedClassName, m.character.CurrentHP, m.character.MaxHP)
			}
		} else {
			selectedCount := len(m.classSkillSelector.SelectedSkills)
			maxCount := m.classSkillSelector.MaxChoices
			m.message = fmt.Sprintf("Please select %d skill(s) (%d/%d selected)", maxCount, selectedCount, maxCount)
		}
	case "esc":
		debug.Log("Class skill selector: cancelled - restoring previous state")
		m.classSkillSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.storage.Save(m.character)
		m.message = "Class selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleFightingStyleSelectorKeys handles fighting style selector specific keys
func (m *Model) handleFightingStyleSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleFightingStyleSelectorKeys: key=%s", msg.String())

	var cmd tea.Cmd

	switch msg.String() {
	case "up", "k":
		m.fightingStyleSelector.Prev()
		return m, nil
	case "down", "j":
		m.fightingStyleSelector.Next()
		return m, nil
	case "enter":
		selectedStyle := m.fightingStyleSelector.GetSelectedStyle()
		debug.Log("Fighting style selector: enter pressed, selected=%s", selectedStyle)

		if selectedStyle != "" {
			// Apply fighting style
			err := models.ApplyFightingStyle(m.character, selectedStyle)
			if err != nil {
				debug.Log("ERROR applying fighting style: %v", err)
				m.message = fmt.Sprintf("Error applying fighting style: %v", err)
			} else {
				debug.Log("Fighting style '%s' applied successfully", selectedStyle)
				// Update the choice record with fighting style
				m.character.Choices.Class.FightingStyle = selectedStyle

				// Check if Blessed Warrior was selected (needs 2 cantrips)
				if selectedStyle == "Blessed Warrior" {
					debug.Log("Blessed Warrior selected - showing cantrip selector for 2 cantrips")
					m.cantripSelector.Show("Cleric", 2) // Blessed Warrior learns from cleric spell list
					m.message = "Select 2 cantrips from the cleric spell list for Blessed Warrior..."
					return m, cmd
				}

				// Check if Druidic Warrior was selected (needs 2 cantrips)
				if selectedStyle == "Druidic Warrior" {
					debug.Log("Druidic Warrior selected - showing cantrip selector for 2 cantrips")
					m.cantripSelector.Show("Druid", 2) // Druidic Warrior learns from druid spell list
					m.message = "Select 2 cantrips from the druid spell list for Druidic Warrior..."
					m.fightingStyleSelector.Hide()
					return m, cmd
				}

				// Check if this is Ranger level 2 and we just completed Deft Explorer
				rangerLevel := m.character.GetClassLevel("Ranger")
				if rangerLevel == 2 {
					// Check if Deft Explorer is complete
					deftExplorerComplete := false
					for _, feature := range m.character.Features.Features {
						if feature.Name == "Deft Explorer" && feature.Mechanics != nil {
							skillSelected, _ := feature.Mechanics["skill_selected"].(bool)
							if skillSelected {
								langCount := 2 // default
								if count, ok := feature.Mechanics["languages"].(float64); ok {
									langCount = int(count)
								}
								languagesSelected := 0
								if count, ok := feature.Mechanics["languages_selected"].(float64); ok {
									languagesSelected = int(count)
								}
								if languagesSelected >= langCount {
									deftExplorerComplete = true
								}
							}
						}
					}

					if deftExplorerComplete {
						// Level 2 setup complete (only if not Blessed/Druidic Warrior which need cantrips)
						m.fightingStyleSelector.Hide()
						m.storage.Save(m.character)
						m.message = fmt.Sprintf("Fighting style '%s' selected! Ranger level 2 setup complete.", selectedStyle)
						return m, cmd
					}
				}

				// Check if character also needs weapon mastery selection
				{
					masteryCount := getWeaponMasteryCount(m.character)
					debug.Log("After fighting style, checking weapon mastery: count=%d", masteryCount)

					if masteryCount > 0 {
						// Show weapon mastery selector
						debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
						m.weaponMasterySelector.Show(masteryCount)
						m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
						m.fightingStyleSelector.Hide()
					} else {
						// No weapon mastery needed, class setup complete
						m.pendingChanges.Clear() // Clear backup on successful completion
						m.fightingStyleSelector.Hide()

						// If in wizard mode, check if we can advance
						if m.IsInWizard() {
							m.checkAndAdvanceWizardAfterClassSetup()
							// Return immediately to ensure wizard state is preserved
							return m, cmd
						} else {
							m.message = fmt.Sprintf("Fighting style '%s' selected! Class setup complete. (HP: %d/%d)", selectedStyle, m.character.CurrentHP, m.character.MaxHP)
						}
					}
				}
			}
			m.storage.Save(m.character)
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		debug.Log("Fighting style selector: cancelled")
		m.fightingStyleSelector.Hide()
		m.message = "Fighting style selection cancelled"
	}

	return m, cmd
}

// handleDivineOrderSelectorKeys handles keyboard input for Divine Order selection
func (m *Model) handleDivineOrderSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Allow number keys 3-7 for tab navigation to pass through
	// Only handle keys specifically for Divine Order selection
	switch msg.String() {
	case "3", "4", "5", "6", "7":
		// Tab navigation keys - always allow these to pass through
		// (handled in main Update function's tab navigation check)
		// Return empty to allow fall-through to tab navigation
		return m, cmd
	case "1":
		// Protector selected
		m.SetPendingDivineOrder("Protector")
		debug.Log("Divine Order selected: Protector")
		err := models.ApplyDivineOrderBenefits(m.character, "Protector", "")
		if err != nil {
			m.message = fmt.Sprintf("Error applying Divine Order: %v", err)
			return m, cmd
		}
		m.SetDivineOrderSelectorVisible(false)
		m.message = "Divine Order: Protector selected (Martial weapons + Heavy armor, 3 cantrips)"

		// Check if cantrip selection is needed next
		needsCantrips := m.character.SpellBook.CantripsKnown > 0
		if needsCantrips {
			debug.Log("After Divine Order, showing cantrip selector for %d cantrips", m.character.SpellBook.CantripsKnown)
			m.cantripSelector.Show("Cleric", m.character.SpellBook.CantripsKnown)
			m.message = fmt.Sprintf("Select %d cantrips for Cleric...", m.character.SpellBook.CantripsKnown)
		} else {
			m.pendingChanges.Clear()
			m.storage.Save(m.character)

			// If in wizard mode, check if we can advance (after Divine Order Protector selection)
			if m.IsInWizard() {
				m.checkAndAdvanceWizardAfterClassSetup()
				// Return immediately to ensure wizard state is preserved
				return m, cmd
			}
		}
	case "2":
		// Thaumaturgic selected - need to choose skill first
		m.SetPendingDivineOrder("Thaumaturgic")
		debug.Log("Divine Order selected: Thaumaturgic - prompting for skill choice")
		m.message = "Choose skill for Thaumaturgic: [a] Arcana or [r] Religion"
	case "a":
		// Arcana chosen for Thaumaturgic
		if m.GetPendingDivineOrder() == "Thaumaturgic" {
			m.SetPendingDivineOrderSkill("Arcana")
			debug.Log("Thaumaturgic skill selected: Arcana")
			err := models.ApplyDivineOrderBenefits(m.character, "Thaumaturgic", "Arcana")
			if err != nil {
				m.message = fmt.Sprintf("Error applying Divine Order: %v", err)
				return m, cmd
			}
			m.SetDivineOrderSelectorVisible(false)
			m.message = "Divine Order: Thaumaturgic selected (Extra cantrip + Arcana expertise)"

			// Check if cantrip selection is needed next (Thaumaturgic gets +1 cantrip)
			cantripsKnown := m.character.SpellBook.CantripsKnown
			if cantripsKnown > 0 {
				// Thaumaturgic gets one extra cantrip
				cantripsKnown++
				debug.Log("Thaumaturgic: increasing cantrips from %d to %d", m.character.SpellBook.CantripsKnown, cantripsKnown)
				m.character.SpellBook.CantripsKnown = cantripsKnown
				m.cantripSelector.Show("Cleric", cantripsKnown)
				m.message = fmt.Sprintf("Select %d cantrips for Cleric (Thaumaturgic: +1 extra)...", cantripsKnown)
			} else {
				m.pendingChanges.Clear()
				m.storage.Save(m.character)
			}
		}
	case "r":
		// Religion chosen for Thaumaturgic
		if m.GetPendingDivineOrder() == "Thaumaturgic" {
			m.SetPendingDivineOrderSkill("Religion")
			debug.Log("Thaumaturgic skill selected: Religion")
			err := models.ApplyDivineOrderBenefits(m.character, "Thaumaturgic", "Religion")
			if err != nil {
				m.message = fmt.Sprintf("Error applying Divine Order: %v", err)
				return m, cmd
			}
			m.SetDivineOrderSelectorVisible(false)
			m.message = "Divine Order: Thaumaturgic selected (Extra cantrip + Religion expertise)"

			// Check if cantrip selection is needed next (Thaumaturgic gets +1 cantrip)
			cantripsKnown := m.character.SpellBook.CantripsKnown
			if cantripsKnown > 0 {
				// Thaumaturgic gets one extra cantrip
				cantripsKnown++
				debug.Log("Thaumaturgic: increasing cantrips from %d to %d", m.character.SpellBook.CantripsKnown, cantripsKnown)
				m.character.SpellBook.CantripsKnown = cantripsKnown
				m.cantripSelector.Show("Cleric", cantripsKnown)
				m.message = fmt.Sprintf("Select %d cantrips for Cleric (Thaumaturgic: +1 extra)...", cantripsKnown)
			} else {
				m.pendingChanges.Clear()
				m.storage.Save(m.character)

				// If in wizard mode, check if we can advance (after Divine Order selection)
				if m.IsInWizard() {
					m.checkAndAdvanceWizardAfterClassSetup()
					// Return immediately to ensure wizard state is preserved
					return m, cmd
				}
			}
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		debug.Log("Divine Order selection cancelled")
		m.SetDivineOrderSelectorVisible(false)
		m.stateMachine.ClearContext("pendingDivineOrder")
		m.stateMachine.ClearContext("pendingDivineOrderSkill")
		m.pendingDivineOrder = ""
		m.pendingDivineOrderSkill = ""
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.storage.Save(m.character)
		m.message = "Divine Order selection cancelled - restored previous state"
	}
	return m, cmd
}

// handlePrimalOrderSelectorKeys handles keyboard input for Primal Order selection
func (m *Model) handlePrimalOrderSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Allow number keys 3-7 for tab navigation to pass through
	// Only handle keys specifically for Primal Order selection
	switch msg.String() {
	case "3", "4", "5", "6", "7":
		// Tab navigation keys - always allow these to pass through
		// (handled in main Update function's tab navigation check)
		// Return empty to allow fall-through to tab navigation
		return m, cmd
	case "1":
		// Magician selected - need to choose skill first
		m.SetPendingPrimalOrder("Magician")
		debug.Log("Primal Order selected: Magician - prompting for skill choice")
		m.message = "Choose skill for Magician: [a] Arcana or [n] Nature"
	case "2":
		// Warden selected
		m.SetPendingPrimalOrder("Warden")
		debug.Log("Primal Order selected: Warden")
		err := models.ApplyPrimalOrderBenefits(m.character, "Warden", "")
		if err != nil {
			m.message = fmt.Sprintf("Error applying Primal Order: %v", err)
			return m, cmd
		}
		m.SetPrimalOrderSelectorVisible(false)
		m.message = "Primal Order: Warden selected (Martial weapons + Medium armor)"

		// Check if cantrip selection is needed next
		needsCantrips := m.character.SpellBook.CantripsKnown > 0
		if needsCantrips {
			debug.Log("After Primal Order, showing cantrip selector for %d cantrips", m.character.SpellBook.CantripsKnown)
			m.cantripSelector.Show("Druid", m.character.SpellBook.CantripsKnown)
			m.message = fmt.Sprintf("Select %d cantrips for Druid...", m.character.SpellBook.CantripsKnown)
		} else {
			m.pendingChanges.Clear()
			m.storage.Save(m.character)

			// If in wizard mode, check if we can advance (after Primal Order Warden selection)
			if m.IsInWizard() {
				m.checkAndAdvanceWizardAfterClassSetup()
				// Return immediately to ensure wizard state is preserved
				return m, cmd
			}
		}
	case "a":
		// Arcana chosen for Magician
		if m.GetPendingPrimalOrder() == "Magician" {
			m.SetPendingPrimalOrderSkill("Arcana")
			debug.Log("Magician skill selected: Arcana")
			err := models.ApplyPrimalOrderBenefits(m.character, "Magician", "Arcana")
			if err != nil {
				m.message = fmt.Sprintf("Error applying Primal Order: %v", err)
				return m, cmd
			}
			m.SetPrimalOrderSelectorVisible(false)
			m.message = "Primal Order: Magician selected (Extra cantrip + Wisdom modifier to Arcana checks, min +1)"

			// Check if cantrip selection is needed next (Magician gets +1 cantrip)
			cantripsKnown := m.character.SpellBook.CantripsKnown
			if cantripsKnown > 0 {
				// Magician gets one extra cantrip
				cantripsKnown++
				debug.Log("Magician: increasing cantrips from %d to %d", m.character.SpellBook.CantripsKnown, cantripsKnown)
				m.character.SpellBook.CantripsKnown = cantripsKnown
				m.cantripSelector.Show("Druid", cantripsKnown)
				m.message = fmt.Sprintf("Select %d cantrips for Druid (Magician: +1 extra)...", cantripsKnown)
			} else {
				m.pendingChanges.Clear()
				m.storage.Save(m.character)
			}
		}
	case "n":
		// Nature chosen for Magician
		if m.GetPendingPrimalOrder() == "Magician" {
			m.SetPendingPrimalOrderSkill("Nature")
			debug.Log("Magician skill selected: Nature")
			err := models.ApplyPrimalOrderBenefits(m.character, "Magician", "Nature")
			if err != nil {
				m.message = fmt.Sprintf("Error applying Primal Order: %v", err)
				return m, cmd
			}
			m.SetPrimalOrderSelectorVisible(false)
			m.message = "Primal Order: Magician selected (Extra cantrip + Wisdom modifier to Nature checks, min +1)"

			// Check if cantrip selection is needed next (Magician gets +1 cantrip)
			cantripsKnown := m.character.SpellBook.CantripsKnown
			if cantripsKnown > 0 {
				// Magician gets one extra cantrip
				cantripsKnown++
				debug.Log("Magician: increasing cantrips from %d to %d", m.character.SpellBook.CantripsKnown, cantripsKnown)
				m.character.SpellBook.CantripsKnown = cantripsKnown
				m.cantripSelector.Show("Druid", cantripsKnown)
				m.message = fmt.Sprintf("Select %d cantrips for Druid (Magician: +1 extra)...", cantripsKnown)
			} else {
				m.pendingChanges.Clear()
				m.storage.Save(m.character)

				// If in wizard mode, check if we can advance (after Primal Order selection)
				if m.IsInWizard() {
					m.checkAndAdvanceWizardAfterClassSetup()
					// Return immediately to ensure wizard state is preserved
					return m, cmd
				}
			}
		}
	case "esc":
		// Check if in wizard mode
		if m.IsInWizard() {
			m.cancelWizard()
			return m, nil
		}

		debug.Log("Primal Order selection cancelled")
		m.SetPrimalOrderSelectorVisible(false)
		m.stateMachine.ClearContext("pendingPrimalOrder")
		m.stateMachine.ClearContext("pendingPrimalOrderSkill")
		m.pendingPrimalOrder = ""
		m.pendingPrimalOrderSkill = ""
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.storage.Save(m.character)
		m.message = "Primal Order selection cancelled - restored previous state"
	}
	return m, cmd
}

// renderPrimalOrderSelector renders the Primal Order selection popup (Druid-specific)
func (m *Model) renderPrimalOrderSelector() string {
	popupMediumWidth := int(float64(m.width) * 0.60)
	popupMediumHeight := int(float64(m.height) * 0.50)
	if popupMediumWidth < 70 {
		popupMediumWidth = 70
	}
	if popupMediumHeight < 20 {
		popupMediumHeight = 20
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		MarginBottom(1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(popupMediumWidth - 8).
		Align(lipgloss.Center)

	var optionsText string
	if m.GetPendingPrimalOrder() == "Magician" {
		// Show skill selection for Magician
		optionsText = optionStyle.Render("Choose skill for Magician:") + "\n\n" +
			highlightStyle.Render("[a]") + " Arcana\n" +
			highlightStyle.Render("[n]") + " Nature\n\n" +
			hintStyle.Render("Press 'a' for Arcana or 'n' for Nature")
	} else {
		// Show Primal Order selection
		optionsText = titleStyle.Render("Select Primal Order") + "\n\n" +
			optionStyle.Render("Choose how you channel your connection to nature:") + "\n\n" +
			highlightStyle.Render("[1]") + " Magician\n" +
			"   Extra cantrip (+1)\n" +
			"   Add Wisdom modifier to Intelligence checks (Arcana or Nature, min +1)\n\n" +
			highlightStyle.Render("[2]") + " Warden\n" +
			"   Proficiency with martial weapons and medium armor\n\n" +
			hintStyle.Render("Press 1 for Magician or 2 for Warden • ESC to cancel")
	}

	box := boxStyle.Render(optionsText)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

// renderDivineOrderSelector renders the Divine Order selection popup (Cleric-specific)
func (m *Model) renderDivineOrderSelector() string {
	popupMediumWidth := int(float64(m.width) * 0.60)
	popupMediumHeight := int(float64(m.height) * 0.50)
	if popupMediumWidth < 70 {
		popupMediumWidth = 70
	}
	if popupMediumHeight < 20 {
		popupMediumHeight = 20
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		MarginBottom(1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(popupMediumWidth - 8).
		Align(lipgloss.Center)

	var optionsText string
	if m.GetPendingDivineOrder() == "Thaumaturgic" {
		// Show skill selection for Thaumaturgic
		optionsText = optionStyle.Render("Choose skill for Thaumaturgic:") + "\n\n" +
			highlightStyle.Render("[a]") + " Arcana\n" +
			highlightStyle.Render("[r]") + " Religion\n\n" +
			hintStyle.Render("Press 'a' for Arcana or 'r' for Religion")
	} else {
		// Show Divine Order selection
		optionsText = titleStyle.Render("Select Divine Order") + "\n\n" +
			optionStyle.Render("Choose how you channel your divine faith:") + "\n\n" +
			highlightStyle.Render("[1]") + " Protector\n" +
			"   Proficiency with martial weapons and heavy armor\n" +
			"   3 cantrips\n\n" +
			highlightStyle.Render("[2]") + " Thaumaturgic\n" +
			"   Extra cantrip (+1, for 4 total)\n" +
			"   Expertise in Arcana or Religion (your choice)\n\n" +
			hintStyle.Render("Press 1 for Protector or 2 for Thaumaturgic • ESC to cancel")
	}

	box := boxStyle.Render(optionsText)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

// handleOptionSelectorKeys handles option selector specific keys
func (m *Model) handleOptionSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleOptionSelectorKeys: key=%s", msg.String())

	var cmd tea.Cmd
	updated, cmd := m.optionSelector.Update(msg)
	m.optionSelector = &updated

	switch msg.String() {
	case "enter":
		selectedValue := m.optionSelector.GetSelectedValue()
		debug.Log("Option selector: enter pressed, selected=%s", selectedValue)

		if selectedValue != "" {
			m.optionSelector.Hide()

			// Determine selection type from message
			// Check if we're selecting Fey Gift
			if m.message == "Select your Fey Gift..." {
				// Add Fey Gift to character's Personality traits
				m.character.Personality = append(m.character.Personality, selectedValue)
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Fey Gift selected: %s", selectedValue)

				// Check if we need Otherworldly Glamour skill selection next
				if m.levelUpSelector.NeedsGlamourSkillSelection {
					m.levelUpSelector.NeedsGlamourSkillSelection = false
					glamourSkills := []string{"Deception", "Performance", "Persuasion"}
					m.classSkillSelector.Show("Ranger", glamourSkills, 1, m.character)
					m.message = "Select 1 skill for Otherworldly Glamour (Deception, Performance, or Persuasion)..."
					return m, cmd
				}
			} else if m.message == "Select your Hunter's Prey option..." {
				// Store Hunter's Prey selection in feature mechanics
				rangerClass := m.character.GetClassLevelStruct("Ranger")
				if rangerClass != nil {
					// Find the Hunter's Prey feature and update its mechanics
					for i := range m.character.Features.Features {
						feature := &m.character.Features.Features[i]
						if feature.Name == "Hunter's Prey" {
							if feature.Mechanics == nil {
								feature.Mechanics = make(map[string]interface{})
							}
							feature.Mechanics["selected_option"] = selectedValue
							debug.Log("Stored Hunter's Prey selection: %s", selectedValue)
							break
						}
					}
				}
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Hunter's Prey selected: %s", selectedValue)
			}

			// Complete level-up if needed
			if !m.levelUpSelector.IsVisible() {
				m.character.UpdateDerivedStats()
				m.message += " Level-up complete!"
				// Ensure focusArea is set to FocusMain after selection
				m.focusArea = FocusMain
			}
		}
	case "esc":
		m.optionSelector.Hide()
		m.message = "Selection cancelled"
		// Ensure focusArea is set to FocusMain after cancellation
		m.focusArea = FocusMain
	}

	return m, cmd
}
