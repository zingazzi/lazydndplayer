// internal/ui/handlers_panel.go
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/dice"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
)

// handleMainPanelKeys handles keys when main panel has focus
func (m *Model) handleMainPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleMainPanelKeys: key=%s, currentPanel=%d (0=Stats,1=Skills,2=Inv,3=Spells,4=Features,5=Traits,6=Origin)", msg.String(), m.currentPanel)
	switch m.currentPanel {
	case StatsPanel:
		return m.handleStatsPanel(msg)
	case SkillsPanel:
		return m.handleSkillsPanel(msg)
	case InventoryPanel:
		return m.handleInventoryPanel(msg)
	case SpellsPanel:
		return m.handleSpellsPanel(msg)
	case FeaturesPanel:
		return m.handleFeaturesPanel(msg)
	case TraitsPanel:
		return m.handleTraitsPanel(msg)
	case OriginPanel:
		return m.handleOriginPanel(msg)
	}
	return m, nil
}

// handleStatsPanel handles stats panel specific keys
func (m *Model) handleStatsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.statsPanel.Prev()
	case "down", "j":
		m.statsPanel.Next()
	case "e":
		// Go directly to extras/modifier editing
		m.statGenerator.ShowExtrasOnly(&m.character.AbilityScores)
		m.message = "Edit ability modifiers..."
	case "r":
		// Open full stat generator for rolling/assigning stats
		m.statGenerator.Show(&m.character.AbilityScores)
		m.message = "Generate ability scores..."
	case "t":
		// Roll saving throw for selected ability
		selectedAbility := m.statsPanel.GetSelectedAbility()
		m.rollSavingThrow(selectedAbility)
	case "a":
		// Roll ability check for selected ability
		selectedAbility := m.statsPanel.GetSelectedAbility()
		m.rollAbilityCheck(selectedAbility)
	}
	return m, nil
}

// rollSavingThrow rolls a saving throw for the given ability
func (m *Model) rollSavingThrow(ability models.AbilityType) {
	expression, message := m.rollService.CalculateSavingThrowRoll(m.character, ability)
	m.dicePanel.Roll(expression)
	m.message = message
}

// rollAbilityCheck rolls an ability check for the given ability
func (m *Model) rollAbilityCheck(ability models.AbilityType) {
	expression, message := m.rollService.CalculateAbilityCheckRoll(m.character, ability)
	m.dicePanel.Roll(expression)
	m.message = message
}

// handleSkillsPanel handles skills panel specific keys
func (m *Model) handleSkillsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.skillsPanel.Prev()
	case "down", "j":
		m.skillsPanel.Next()
	case "e":
		m.skillsPanel.ToggleProficiency()
		m.message = "Proficiency toggled"
	case "r":
		if skill := m.skillsPanel.GetSelectedSkill(); skill != nil {
			abilityMod := m.character.AbilityScores.GetModifier(skill.Ability)
			bonus := skill.CalculateBonus(abilityMod, m.character.ProficiencyBonus)
			expr := fmt.Sprintf("1d20%+d", bonus)
			m.dicePanel.Roll(expr)
			m.message = fmt.Sprintf("Rolling %s: %s", skill.Name, m.dicePanel.LastMessage)
		}
	}
	return m, nil
}

// handleInventoryPanel handles inventory panel specific keys
func (m *Model) handleInventoryPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.inventoryPanel.Prev()
	case "down", "j":
		m.inventoryPanel.Next()
	case "enter":
		// Show item details
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			m.itemDetailPopup.Show(item)
			m.message = "Viewing item details..."
		}
	case "e":
		// Toggle equipped status for selected item
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			// Check if item is equippable
			def := models.GetItemDefinitionByName(item.Name)
			if def != nil && models.IsEquippable(*def) {
				// Only check proficiency when EQUIPPING (not unequipping)
				if !item.Equipped {
					// Check armor proficiency
					if item.Type == models.Armor {
						// Get armor subcategory (Light, Medium, Heavy, Shield)
						armorType := def.Subcategory

						// Check proficiency
						if !models.HasArmorProficiency(m.character, armorType) {
							m.message = fmt.Sprintf("Cannot equip %s: Not proficient with %s armor!", item.Name, armorType)
							return m, nil
						}

						// Unequip other armor pieces first
						models.UnequipOtherArmor(m.character, item)
					}

					// Check weapon proficiency
					if item.Type == models.Weapon {
						// Get weapon subcategory (simple melee, martial melee, etc.)
						weaponType := def.Subcategory

						// Check proficiency
						if !models.HasWeaponProficiency(m.character, weaponType) {
							m.message = fmt.Sprintf("Cannot equip %s: Not proficient with %s weapons!", item.Name, weaponType)
							return m, nil
						}
					}
				}

				m.inventoryPanel.ToggleEquipped()

				// Recalculate AC after equipping/unequipping
				m.character.UpdateDerivedStats()

				if item.Equipped {
					m.message = fmt.Sprintf("%s equipped (AC: %d)", item.Name, m.character.AC)
				} else {
					m.message = fmt.Sprintf("%s unequipped (AC: %d)", item.Name, m.character.AC)
				}
				m.storage.Save(m.character)
			} else {
				m.message = "This item cannot be equipped"
			}
		}
	case "d":
		// Delete selected item (decrease quantity by 1 or remove if quantity is 1)
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			wasEquipped := item.Equipped
			itemType := item.Type
			if item.Quantity > 1 {
				item.Quantity--
				m.message = fmt.Sprintf("%s quantity decreased to %d", item.Name, item.Quantity)
			} else {
				itemName := item.Name
				m.inventoryPanel.DeleteSelected()
				m.message = fmt.Sprintf("%s removed from inventory", itemName)
			}

			// Recalculate AC if armor was equipped
			if wasEquipped && itemType == models.Armor {
				m.character.UpdateDerivedStats()
				m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
			}
			m.storage.Save(m.character)
		}
	case "D":
		// Delete all of selected item
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			itemName := item.Name
			wasEquipped := item.Equipped
			itemType := item.Type
			m.inventoryPanel.DeleteSelected()
			m.message = fmt.Sprintf("All %s removed from inventory", itemName)

			// Recalculate AC if armor was equipped
			if wasEquipped && itemType == models.Armor {
				m.character.UpdateDerivedStats()
				m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
			}
			m.storage.Save(m.character)
		}
	case "a":
		// Open item selector to add items
		m.itemSelector.Show(m.character)
		m.message = "Select item category..."
	}
	return m, nil
}

// handleSpellsPanel handles spells panel specific keys
func (m *Model) handleSpellsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.spellsPanel.HandleKey(msg)
	case "down", "j":
		m.spellsPanel.HandleKey(msg)
	case "pgup":
		m.spellsPanel.HandleKey(msg)
	case "pgdown":
		m.spellsPanel.HandleKey(msg)
	case "r":
		// Rest - restore all spell slots
		m.spellsPanel.Rest()
		m.storage.Save(m.character)
		m.message = "Spell slots restored!"
	case "c":
		// Change cantrips
		if m.character.SpellBook.IsPreparedCaster {
			m.cantripSelector.Show(m.character.Class, m.character.SpellBook.CantripsKnown)
			m.message = fmt.Sprintf("Select %d cantrips...", m.character.SpellBook.CantripsKnown)
		} else {
			m.message = "Only prepared casters can change cantrips this way"
		}
	case "v", "b":
		// Open spellbook editor (for Wizards) or spell prep selector (for other prepared casters)
		debug.Log("=== KEY 'v' or 'b' PRESSED IN SPELLS PANEL")
		debug.Log("=== IsSpellbookCaster: %v, HasWizard: %v, IsPreparedCaster: %v",
			m.character.SpellBook.IsSpellbookCaster,
			m.character.HasClass("Wizard"),
			m.character.SpellBook.IsPreparedCaster)
		debug.Log("=== Current spells in spellbook: %d", len(m.character.SpellBook.Spells))

		if m.character.SpellBook.IsSpellbookCaster && (m.character.HasClass("Wizard") || m.character.IsArcaneTrickster()) {
			debug.Log("=== OPENING SPELLBOOK EDITOR")
			m.spellbookEditor.Show()
			m.message = "Managing spellbook... (Space: Prepare | a: Add | d/x: Remove | c: Cantrips | 0-9: Filter)"
			debug.Log("=== Spellbook editor visible: %v", m.spellbookEditor.IsVisible())
		} else if m.character.SpellBook.IsPreparedCaster || m.character.HasClass("Paladin") {
			// Paladins are prepared casters (like Clerics)
			// Check if Paladin has spellcasting (level 1+)
			if m.character.HasClass("Paladin") && m.character.GetClassLevel("Paladin") >= 1 {
				// Ensure spellcasting is initialized
				if !m.character.SpellBook.IsPreparedCaster {
					// Re-initialize spellcasting if needed
					class := models.GetClassByName("Paladin")
					if class != nil {
						models.InitializeSpellcasting(m.character, class)
						m.character.UpdateDerivedStats() // Recalculate prepared spell limit
					}
				}
			}
			debug.Log("=== OPENING SPELL PREP SELECTOR")
			m.spellPrepSelector.Show()
			m.message = "Managing spellbook... (Tab: Switch tabs • Space: Prepare/Add • a/d: Add/Remove cantrips)"
		} else {
			debug.Log("=== Not a prepared caster - showing error message")
			m.message = "Only prepared casters can prepare spells"
		}
	case "s":
		// Open slot restorer
		if m.character.SpellBook.SpellcastingMod != "" {
			m.slotRestorer.Show()
			m.message = "Select spell slot to restore..."
		} else {
			m.message = "Not a spellcaster"
		}
	case "a":
		// Add new spell to spellbook (Wizard only)
		debug.Log("=== KEY 'a' PRESSED IN SPELLS PANEL")
		debug.Log("=== IsSpellbookCaster: %v, HasWizard: %v",
			m.character.SpellBook.IsSpellbookCaster,
			m.character.HasClass("Wizard"))

		if m.character.SpellBook.IsSpellbookCaster && (m.character.HasClass("Wizard") || m.character.IsArcaneTrickster()) {
			var maxSpellLevel int
			if m.character.HasClass("Wizard") {
				wizardLevel := m.character.GetClassLevel("Wizard")
				maxSpellLevel = (wizardLevel + 1) / 2
			} else if m.character.IsArcaneTrickster() {
				rogueLevel := m.character.GetRogueLevel()
				maxSpellLevel = (rogueLevel + 2) / 3 // Third caster progression
			}
			if maxSpellLevel > 9 {
				maxSpellLevel = 9
			}
			debug.Log("=== SHOWING SPELLBOOK EDITOR: MaxSpellLevel=%d", maxSpellLevel)
			// Open spellbook editor to add spells
			m.spellbookEditor.Show()
			m.message = fmt.Sprintf("Add new spells to your spellbook (up to level %d)...", maxSpellLevel)
			debug.Log("=== Spellbook editor visible: %v", m.spellbookEditor.IsVisible())
		} else {
			debug.Log("=== Not a spellbook caster")
			m.message = "Only spellbook casters can add spells to their spellbook"
		}
	case "enter":
		// View spell details
		spell := m.spellsPanel.GetSelectedSpell()
		if spell != nil {
			m.spellDetailPopup.Show(*spell)
			m.message = "Viewing spell details..."
		}
	case "u":
		// Consume spell slot
		if m.spellsPanel.ConsumeSpellSlot() {
			m.storage.Save(m.character)
			m.message = "Spell slot consumed"
		} else {
			m.message = "Cannot consume spell slot (already at 0 or no slot selected)"
		}
	case "U":
		// Restore spell slot
		if m.spellsPanel.RestoreSpellSlot() {
			m.storage.Save(m.character)
			m.message = "Spell slot restored"
		} else {
			m.message = "Cannot restore spell slot (already at max or no slot selected)"
		}
	}
	return m, nil
}

// handleFeaturesPanel handles features panel specific keys
func (m *Model) handleFeaturesPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.featuresPanel.Prev()
	case "down", "j":
		m.featuresPanel.Next()
	case "ctrl+u", "pgup":
		m.featuresPanel.PageUp()
	case "ctrl+d", "pgdown":
		m.featuresPanel.PageDown()
	case "ctrl+y":
		m.featuresPanel.ScrollUp()
	case "ctrl+e":
		m.featuresPanel.ScrollDown()
	case "enter":
		// Show popup for selected item (any type)
		selectedItem := m.featuresPanel.GetSelectedItem()
		if selectedItem != nil {
			switch selectedItem.ItemType {
			case "consumable":
				// Show consumable detail popup
				item := selectedItem.Consumable
				if item != nil {
					restTypeStr := "Unknown"
					switch item.RestType {
					case models.ShortRest:
						restTypeStr = "Short Rest"
					case models.LongRest:
						restTypeStr = "Long Rest"
					case models.Daily:
						restTypeStr = "Daily"
					}
					m.consumableDetailPopup.Show(item.Name, item.Current, item.Max, restTypeStr, item.Description)
					m.message = "Viewing details..."
				}
			case "passive", "rage_effect":
				// Show feature detail popup
				usesStr := ""
				restTypeStr := ""

				// For passive features, show uses info if applicable
				if selectedItem.Feature != nil {
					if selectedItem.Feature.MaxUses > 0 {
						usesStr = fmt.Sprintf("%d/%d", selectedItem.Feature.CurrentUses, selectedItem.Feature.MaxUses)
					}
					restTypeStr = string(selectedItem.Feature.RestType)
				}

				m.featureDetailPopup.Show(selectedItem.Name, selectedItem.Description, usesStr, restTypeStr)
				m.message = "Viewing feature details..."
			}
		} else {
			m.message = "No item selected"
		}
		return m, nil
	case "u":
		// Consume selected item
		item := m.featuresPanel.GetSelectedConsumable()
		if item == nil {
			m.message = "No item selected"
			return m, nil
		}

		if item.Current <= 0 {
			m.message = fmt.Sprintf("%s has no uses remaining", item.Name)
			return m, nil
		}

		// Handle by type
		switch item.ItemType {
		case "resource":
			switch item.ResourceType {
			case "focus_points":
				monk := m.character.GetMonkMechanics()
				if monk.SpendFocusPoint(1) {
					current, max := monk.GetFocusPoints()
					m.message = fmt.Sprintf("Focus Point spent. Current: %d/%d", current, max)
				} else {
					m.message = "Not enough Focus Points"
				}
			case "psi_dice":
				m.character.PsiDice.Current--
				m.message = fmt.Sprintf("Psi Die spent. Current: %d/%d", m.character.PsiDice.Current, m.character.PsiDice.Max)
			case "superiority_dice":
				m.character.SuperiorityDice.Current--
				m.message = fmt.Sprintf("Superiority Die spent. Current: %d/%d", m.character.SuperiorityDice.Current, m.character.SuperiorityDice.Max)
			case "warrior_dice":
				m.character.WarriorDice.Current--
				m.message = fmt.Sprintf("Warrior Die spent. Current: %d/%d", m.character.WarriorDice.Current, m.character.WarriorDice.Max)
			case "soulknife_psi_dice":
				m.character.SoulknifePsiDice.Current--
				m.message = fmt.Sprintf("Psionic Energy die spent. Current: %d/%d", m.character.SoulknifePsiDice.Current, m.character.SoulknifePsiDice.Max)
			}
			m.storage.Save(m.character)
		case "feature":
			// Handle feature consumption with special effects
			if item.Feature != nil {
				switch item.Feature.Name {
				case "Second Wind":
					// Roll 1d10 + fighter level
					result, err := dice.Roll("1d10", dice.Normal)
					if err == nil {
						fighterLevel := m.character.GetFighterLevel()
						healing := result.Total + fighterLevel
						m.character.CurrentHP += healing
						if m.character.CurrentHP > m.character.MaxHP {
							m.character.CurrentHP = m.character.MaxHP
						}
						m.message = fmt.Sprintf("Second Wind: Healed %d HP (1d10[%d] + %d level)", healing, result.Total, fighterLevel)
					} else {
						fighterLevel := m.character.GetFighterLevel()
						m.message = fmt.Sprintf("Second Wind: Healed %d HP", fighterLevel)
						m.character.CurrentHP += fighterLevel
					}
				default:
					m.message = fmt.Sprintf("%s used", item.Feature.Name)
				}
				// Decrement uses
				m.featuresPanel.UseFeature()
				m.storage.Save(m.character)
			}
		}
		return m, nil
	case "U":
		// Restore selected item
		item := m.featuresPanel.GetSelectedConsumable()
		if item == nil {
			m.message = "No item selected"
			return m, nil
		}

		if item.Current >= item.Max {
			m.message = fmt.Sprintf("%s already at maximum", item.Name)
			return m, nil
		}

		// Handle by type
		switch item.ItemType {
		case "resource":
			switch item.ResourceType {
			case "focus_points":
				monk := m.character.GetMonkMechanics()
				monk.RestoreFocusPoints(1)
				current, max := monk.GetFocusPoints()
				m.message = fmt.Sprintf("Focus Point restored. Current: %d/%d", current, max)
			case "psi_dice":
				m.character.PsiDice.Current++
				m.message = fmt.Sprintf("Psi Die restored. Current: %d/%d", m.character.PsiDice.Current, m.character.PsiDice.Max)
			case "superiority_dice":
				m.character.SuperiorityDice.Current++
				m.message = fmt.Sprintf("Superiority Die restored. Current: %d/%d", m.character.SuperiorityDice.Current, m.character.SuperiorityDice.Max)
			case "warrior_dice":
				m.character.WarriorDice.Current++
				m.message = fmt.Sprintf("Warrior Die restored. Current: %d/%d", m.character.WarriorDice.Current, m.character.WarriorDice.Max)
			case "soulknife_psi_dice":
				m.character.SoulknifePsiDice.Current++
				m.message = fmt.Sprintf("Psionic Energy die restored. Current: %d/%d", m.character.SoulknifePsiDice.Current, m.character.SoulknifePsiDice.Max)
			}
			m.storage.Save(m.character)
		case "feature":
			m.featuresPanel.RestoreFeature()
			m.message = fmt.Sprintf("%s restored", item.Name)
			m.storage.Save(m.character)
		}
		return m, nil
	case "r":
		// Short rest
		m.character.ShortRest()
		m.message = "Short rest completed - features recovered"
		m.storage.Save(m.character)
	case "R":
		// Long rest
		m.character.LongRest()
		m.message = "Long rest completed - all features recovered"
		m.storage.Save(m.character)
	}
	return m, nil
}

// handleTraitsPanel handles traits panel specific keys
func (m *Model) handleTraitsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleTraitsPanel: key=%s", msg.String())
	switch msg.String() {
	case "up", "k":
		// Navigate to previous selectable item
		debug.Log("handleTraitsPanel: calling Prev()")
		m.traitsPanel.Prev()
	case "down", "j":
		// Navigate to next selectable item
		debug.Log("handleTraitsPanel: calling Next()")
		m.traitsPanel.Next()
	case "shift+up":
		// Scroll viewport up
		debug.Log("handleTraitsPanel: scrolling viewport up")
		m.traitsPanel.ScrollUp()
	case "shift+down":
		// Scroll viewport down
		debug.Log("handleTraitsPanel: scrolling viewport down")
		m.traitsPanel.ScrollDown()
	case "ctrl+u", "pgup":
		// Page up
		m.traitsPanel.PageUp()
	case "ctrl+d", "pgdown":
		// Page down
		m.traitsPanel.PageDown()
	case "g":
		// Go to top
		m.traitsPanel.GotoTop()
	case "G":
		// Go to bottom (Shift+G)
		m.traitsPanel.GotoBottom()
	case "enter":
		// Show feat detail popup if on a feat
		if m.traitsPanel.IsOnFeat() {
			featName := m.traitsPanel.GetSelectedFeat()
			if featName != "" {
				m.featDetailPopup.Show(featName, m.character)
				m.message = "Viewing feat details..."
			}
		}
		// Show mastery detail popup if on a weapon mastery
		if m.traitsPanel.IsOnMastery() {
			weaponName, masteryType := m.traitsPanel.GetSelectedMastery()
			if weaponName != "" && masteryType != "" {
				m.masteryDetailPopup.Show(weaponName, masteryType)
				m.message = "Viewing weapon mastery details..."
			}
		}
		// Show maneuver detail popup if on a maneuver
		if m.traitsPanel.IsOnManeuver() {
			maneuverName := m.traitsPanel.GetSelectedManeuver()
			if maneuverName != "" {
				m.maneuverDetailPopup.Show(maneuverName)
				m.message = "Viewing maneuver details..."
			}
		}
	case "l":
		// Add language
		m.languageSelector.SetExcludeLanguages(m.character.Languages)
		m.languageSelector.Show()
		m.message = "Select a language to learn..."
	case "L": // Shift+L
		// Remove language
		if len(m.character.Languages) == 0 {
			m.message = "No languages to remove"
		} else {
			m.languageSelector.ShowForDeletion(m.character.Languages)
			m.message = "Select a language to remove..."
		}
	case "f":
		// Add feat
		m.featSelector.Show(m.character, false) // false = not an origin feat
		m.message = "Select a feat to acquire..."
	case "F": // Shift+F
		// Remove feat
		if len(m.character.Feats) == 0 {
			m.message = "No feats to remove"
		} else {
			m.featSelector.ShowForDeletion(m.character)
			m.message = "Select a feat to remove..."
		}
	case "s":
		// Change fighting style
		if m.character.FightingStyle != "" {
			// Character already has a fighting style, allow changing it
			m.fightingStyleSelector.Show(m.character.Class)
			m.message = "Select a new fighting style..."
		} else {
			m.message = "You don't have a Fighting Style to change"
		}
	case "m":
		// Manage weapon mastery
		debug.Log("handleTraitsPanel: 'm' key pressed - checking weapon mastery")
		// Check if character has weapon mastery feature
		masteryCount := getWeaponMasteryCount(m.character)
		debug.Log("handleTraitsPanel: masteryCount=%d", masteryCount)
		if masteryCount > 0 {
			debug.Log("handleTraitsPanel: Showing weapon mastery selector")
			m.weaponMasterySelector.Show(masteryCount)
			m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
		} else {
			debug.Log("handleTraitsPanel: No weapon mastery feature found")
			m.message = "You don't have the Weapon Mastery feature"
		}
	case "e":
		// Manage expertise
		debug.Log("handleTraitsPanel: 'e' key pressed - checking expertise")
		expertiseCount := getExpertiseCount(m.character)
		debug.Log("handleTraitsPanel: expertiseCount=%d", expertiseCount)
		if expertiseCount > 0 {
			debug.Log("handleTraitsPanel: Showing expertise selector")
			m.expertiseSelector.Show(expertiseCount)
			m.message = fmt.Sprintf("Select %d skill(s) for expertise...", expertiseCount)
		} else {
			debug.Log("handleTraitsPanel: No expertise feature found")
			m.message = "You don't have the Expertise feature"
		}
	case "n":
		// Manage Battle Master maneuvers (edit only, no Student of War benefits)
		debug.Log("handleTraitsPanel: 'n' key pressed - checking Battle Master maneuvers")
		// Check if character is a Battle Master
		if m.character.IsBattleMaster() {
			// Get maneuver count from Combat Superiority feature
			maneuverCount := 3 // Default for level 3
			if feature := m.character.GetFeature("Combat Superiority"); feature != nil && feature.Mechanics != nil {
				if count, ok := feature.Mechanics["maneuvers_known"].(float64); ok {
					maneuverCount = int(count)
				}
			}
			debug.Log("handleTraitsPanel: Showing maneuver selector for %d maneuvers", maneuverCount)
			// Clear Student of War flag to prevent prompting for tool/skill
			m.studentOfWarToolSelected = false
			// Show selector first (which clears selections), THEN load current maneuvers
			m.maneuverSelector.Show(maneuverCount)
			m.maneuverSelector.SetSelectedManeuvers(m.character.Maneuvers)
			debug.Log("handleTraitsPanel: Loaded %d existing maneuvers into selector", len(m.character.Maneuvers))
			m.message = fmt.Sprintf("Select up to %d maneuvers...", maneuverCount)
		} else {
			debug.Log("handleTraitsPanel: Not a Battle Master")
			m.message = "Only Battle Masters can learn maneuvers"
		}
	case "d", "x":
		m.traitsPanel.RemoveSelected()
		m.message = "Item removed"
		m.storage.Save(m.character)
	}
	return m, nil
}

// handleOriginPanel handles origin panel specific keys
func (m *Model) handleOriginPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.originPanel.ScrollUp()
	case "down", "j":
		m.originPanel.ScrollDown()
	case "ctrl+u", "pgup":
		m.originPanel.PageUp()
	case "ctrl+d", "pgdown":
		m.originPanel.PageDown()
	case "o":
		// Open origin selector
		m.originSelector.Show(m.character)
		m.message = "Select an origin..."
	case "enter":
		// Show origin details
		if m.character.Origin != "" {
			origin := models.GetOriginByName(m.character.Origin)
			if origin != nil {
				m.originDetailPopup.Show(origin)
			}
		}
	case "a":
		// Open alignment selector
		m.alignmentSelector.Show()
		m.message = "Select an alignment..."
	case "h":
		// Edit height
		m.inputPopup.Show("Edit Height", m.character.Height, "Enter height (e.g., 6'2\")...")
		m.SetInputPopupContext("height")
		m.message = "Editing height..."
	case "w":
		// Edit weight
		m.inputPopup.Show("Edit Weight", m.character.Weight, "Enter weight (e.g., 180 lbs)...")
		m.SetInputPopupContext("weight")
		m.message = "Editing weight..."
	case "t":
		// Edit personality
		m.traitSelector.Show(models.TraitPersonality, m.character.Personality)
		m.message = "Select or create personality traits (Space to toggle, Enter to confirm)..."
	case "i":
		// Edit ideal
		m.traitSelector.Show(models.TraitIdeal, m.character.Ideal)
		m.message = "Select or create ideals (Space to toggle, Enter to confirm)..."
	case "b":
		// Edit bond
		m.traitSelector.Show(models.TraitBond, m.character.Bond)
		m.message = "Select or create bonds (Space to toggle, Enter to confirm)..."
	case "f":
		// Edit flaw
		m.traitSelector.Show(models.TraitFlaw, m.character.Flaw)
		m.message = "Select or create flaws (Space to toggle, Enter to confirm)..."
	case "s":
		// Change species
		m.speciesSelector.Show()
		m.message = "Select a species..."
	case "S":
		// Save character (Shift+S)
		m.storage.Save(m.character)
		m.message = "Character saved!"
	}
	return m, nil
}

// handleActionsPanelKeys handles keys when actions panel has focus
func (m *Model) handleActionsPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if attack menu is active
	if m.attackMenu.IsVisible() {
		return m.handleAttackMenuKeys(msg)
	}

	switch msg.String() {
	case "up", "k":
		m.actionsPanel.Prev()
	case "down", "j":
		m.actionsPanel.Next()
	case "r":
		// Roll normal attack
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "normal")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "a":
		// Roll attack with advantage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "advantage")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "x":
		// Roll attack with disadvantage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "disadvantage")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "d":
		// Roll damage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollDamageDirect(attack)
				m.dicePanel.LastMessage = result
				m.message = result
			}
		} else if m.actionsPanel.IsSneakAttackSelected() {
			// Roll Sneak Attack damage
			sneakAttackDice := m.actionsPanel.GetSneakAttackDice()
			m.dicePanel.Roll(sneakAttackDice)
			m.message = fmt.Sprintf("Sneak Attack damage: %s", m.dicePanel.LastMessage)
		}
	case "enter":
		// Show attack menu if attack is selected
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				debug.Log("Opening attack menu for: %s", attack.Name)
				m.attackMenu.Show(attack)
				m.message = "Select attack option..."
			} else {
				debug.Log("No attack selected (attack is nil)")
			}
		} else if m.actionsPanel.IsSneakAttackSelected() {
			// Roll Sneak Attack damage directly
			sneakAttackDice := m.actionsPanel.GetSneakAttackDice()
			m.dicePanel.Roll(sneakAttackDice)
			m.message = fmt.Sprintf("Sneak Attack damage: %s", m.dicePanel.LastMessage)
		} else if m.actionsPanel.IsSpellSelected() {
			// Cast spell if spell is selected
			spell := m.actionsPanel.GetSelectedSpell()
			if spell != nil {
				debug.Log("Attempting to cast spell: %s", spell.Name)
				msg, success := m.actionsPanel.CastSelectedSpell()
				if success {
					m.storage.Save(m.character)
					m.dicePanel.LastMessage = msg
					m.message = msg
					debug.Log("Spell cast successfully: %s", spell.Name)
				} else {
					m.message = msg
					debug.Log("Failed to cast spell: %s", msg)
				}
			}
		} else {
			// For non-attack/spell actions
			m.message = "Other actions not fully implemented"
		}
	}
	return m, nil
}

// handleAttackMenuKeys handles keys when attack menu is visible
func (m *Model) handleAttackMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleAttackMenuKeys: key=%s", msg.String())

	switch msg.String() {
	case "up", "k":
		m.attackMenu.Prev()
		debug.Log("Attack menu: moved up")
	case "down", "j":
		m.attackMenu.Next()
		debug.Log("Attack menu: moved down")
	case "enter":
		option := m.attackMenu.GetSelectedOption()
		attack := m.attackMenu.GetAttack()
		debug.Log("Attack menu: enter pressed, option=%s, attack=%v", option, attack != nil)

		if attack != nil {
			var result string
			switch {
			case option == "Attack with Advantage":
				result = m.rollAttackDirect(attack, "advantage")
			case option == "Attack with Disadvantage":
				result = m.rollAttackDirect(attack, "disadvantage")
			case option == "Attack (Normal)":
				result = m.rollAttackDirect(attack, "normal")
			case strings.HasPrefix(option, "1-Hand Damage"):
				result = m.rollDamageDirect(attack)
			case strings.HasPrefix(option, "1-Hand Critical"):
				result = m.rollCriticalDamage(attack, attack.DamageDice)
			case strings.HasPrefix(option, "2-Hands Damage"):
				result = m.rollVersatileDamage(attack)
			case strings.HasPrefix(option, "2-Hands Critical"):
				result = m.rollCriticalDamage(attack, attack.VersatileDamage)
			case strings.HasPrefix(option, "Damage"):
				result = m.rollDamageDirect(attack)
			case strings.HasPrefix(option, "Critical Hit"):
				result = m.rollCriticalDamage(attack, attack.DamageDice)
			default:
				result = fmt.Sprintf("Unknown option: %s", option)
			}
			debug.Log("Attack result: %s", result)
			m.dicePanel.LastMessage = result
			m.message = result
		}
		m.attackMenu.Hide()
		debug.Log("Attack menu hidden")
	case "esc":
		debug.Log("Attack menu: cancelled")
		m.attackMenu.Hide()
		m.message = ""
	}
	return m, nil
}

// rollAttackDirect performs an attack roll directly without popup
func (m *Model) rollAttackDirect(attack *models.Attack, rollType string) string {
	var diceRollType dice.RollType
	var advantageStr string

	switch rollType {
	case "advantage":
		diceRollType = dice.Advantage
		advantageStr = "Advantage"
	case "disadvantage":
		diceRollType = dice.Disadvantage
		advantageStr = "Disadvantage"
	default:
		diceRollType = dice.Normal
		advantageStr = ""
	}

	result, err := dice.Roll("1d20", diceRollType)
	if err != nil {
		return fmt.Sprintf("Error rolling: %v", err)
	}

	roll := 0
	if len(result.Rolls) > 0 {
		roll = result.Rolls[0]
	}

	total := roll + attack.AttackBonus

	// Check for critical hit based on character's critical range
	critRange := m.character.GetCriticalRange()
	if roll >= critRange {
		critText := "CRITICAL HIT!"
		if critRange < 20 {
			critText = "CRITICAL HIT! (19-20 range)"
		}
		return fmt.Sprintf("%s: %s [%d] + %d = %d %s",
			attack.Name, critText, roll, attack.AttackBonus, total, advantageStr)
	}

	return attack.FormatAttackRoll(roll, total, advantageStr)
}

// rollDamageDirect performs a damage roll directly without popup
func (m *Model) rollDamageDirect(attack *models.Attack) string {
	result, err := dice.Roll(attack.DamageDice, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling damage: %v", err)
	}

	total := result.Total + attack.DamageBonus
	return attack.FormatDamageRoll(result.Rolls, total)
}

// rollVersatileDamage performs a damage roll with two-handed versatile damage
func (m *Model) rollVersatileDamage(attack *models.Attack) string {
	result, err := dice.Roll(attack.VersatileDamage, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling damage: %v", err)
	}

	// Use TwoHandDamageBonus for two-handed attacks (no Dueling bonus)
	total := result.Total + attack.TwoHandDamageBonus
	return fmt.Sprintf("%s (2-Hands): Damage = %v +%d = %d %s",
		attack.Name, result.Rolls, attack.TwoHandDamageBonus, total, attack.DamageType)
}

// rollCriticalDamage performs a critical hit damage roll (double dice)
func (m *Model) rollCriticalDamage(attack *models.Attack, damageDice string) string {
	// Double the damage dice for critical hits
	// Parse dice notation (e.g., "1d8" -> "2d8", "2d6" -> "4d6")
	critDice := damageDice

	// Simple parsing: if it starts with a number, double it
	parts := strings.Split(damageDice, "d")
	if len(parts) == 2 {
		numDice := 1
		if parts[0] != "" {
			if n, err := fmt.Sscanf(parts[0], "%d", &numDice); err == nil && n == 1 {
				critDice = fmt.Sprintf("%dd%s", numDice*2, parts[1])
			}
		} else {
			// "d8" format, assume 1d8
			critDice = fmt.Sprintf("2d%s", parts[1])
		}
	}

	result, err := dice.Roll(critDice, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling critical damage: %v", err)
	}

	// Determine which damage bonus to use
	damageBonus := attack.DamageBonus
	label := "Critical Hit"

	// If this is a two-handed critical (versatile weapon), use TwoHandDamageBonus
	if attack.VersatileDamage != "" && damageDice == attack.VersatileDamage {
		label = "Critical Hit (2-Hands)"
		damageBonus = attack.TwoHandDamageBonus
	}

	total := result.Total + damageBonus

	return fmt.Sprintf("%s %s: Damage = %v +%d = %d %s",
		attack.Name, label, result.Rolls, damageBonus, total, attack.DamageType)
}

// handleAttackRollerKeys handles keys when attack roller is visible
func (m *Model) handleAttackRollerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	state := m.attackRoller.GetState()

	switch state {
	case "select_attack":
		switch msg.String() {
		case "up", "k":
			m.attackRoller.Prev()
		case "down", "j":
			m.attackRoller.Next()
		case "enter":
			m.attackRoller.SelectAttack()
			m.message = "Choose action: 'a'ttack, 'd'amage, ad'v'antage, disadvantage (x)"
		case "esc":
			m.attackRoller.Hide()
			m.message = ""
		}
	case "select_roll_type":
		switch msg.String() {
		case "a":
			// Roll normal attack
			m.attackRoller.SetRollType("normal")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "d":
			// Roll damage
			result := m.attackRoller.RollDamage()
			m.dicePanel.LastMessage = result
			m.message = result
		case "v":
			// Roll attack with advantage
			m.attackRoller.SetRollType("advantage")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "x":
			// Roll attack with disadvantage
			m.attackRoller.SetRollType("disadvantage")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "esc":
			m.attackRoller.Hide()
			m.message = ""
		}
	}

	return m, nil
}

// handleDicePanelKeys handles keys when dice panel has focus
func (m *Model) handleDicePanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	mode := m.dicePanel.GetMode()

	switch mode {
	case panels.DiceModeIdle:
		// Idle mode - waiting for user to choose action
		switch msg.String() {
		case "enter":
			m.dicePanel.SetMode(panels.DiceModeInput)
			m.message = "Enter dice notation and press Enter to roll"
		case "h":
			m.dicePanel.SetMode(panels.DiceModeHistory)
			m.message = "Navigate history with ↑/↓, press Enter to reroll"
		case "r":
			m.dicePanel.RerollLast()
			m.message = "Rerolled last dice"
		}
		return m, nil

	case panels.DiceModeInput:
		// Input mode - typing dice notation
		switch msg.String() {
		case "esc":
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = ""
			return m, nil
		case "enter":
			if m.dicePanel.GetInput() != "" {
				m.dicePanel.Roll(m.dicePanel.GetInput())
				m.dicePanel.SetMode(panels.DiceModeIdle)
				m.message = ""
			}
			return m, nil
		}
		// Pass all other keys to input
		return m, m.dicePanel.Update(msg)

	case panels.DiceModeHistory:
		// History mode - browsing previous rolls
		switch msg.String() {
		case "esc":
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = ""
		case "up", "k":
			m.dicePanel.HistoryPrev()
		case "down", "j":
			m.dicePanel.HistoryNext()
		case "enter":
			m.dicePanel.RerollSelected()
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = "Rerolled selected dice"
		}
		return m, nil
	}

	return m, nil
}

// handleCharStatsPanelKeys handles character stats panel specific keys
func (m *Model) handleCharStatsPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	editMode := m.characterStatsPanel.GetEditMode()

	// If in edit mode, handle save/cancel
	if editMode != panels.CharStatsNormal {
		switch msg.String() {
		case "enter":
			if editMode == panels.CharStatsEditName {
				m.characterStatsPanel.SaveName()
				m.message = "Name updated"
			} else if editMode == panels.CharStatsEditRace {
				m.characterStatsPanel.SaveRace()
				m.message = "Race updated"
			} else if editMode == panels.CharStatsEditHP {
				amount, err := m.characterStatsPanel.SaveHP()
				if err != nil {
					m.message = fmt.Sprintf("Invalid HP value: %v", err)
				} else {
					m.message = fmt.Sprintf("HP adjusted by %+d. Current: %d/%d", amount, m.character.CurrentHP, m.character.MaxHP)
				}
			} else if editMode == panels.CharStatsEditXP {
				amount, err := m.characterStatsPanel.SaveXP()
				if err != nil {
					m.message = fmt.Sprintf("Invalid XP value: %v", err)
				} else {
					nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
					if m.character.Experience >= nextLevelXP {
						m.message = fmt.Sprintf("XP adjusted by %+d. Current: %d (LEVEL UP AVAILABLE!)", amount, m.character.Experience)
					} else {
						m.message = fmt.Sprintf("XP adjusted by %+d. Current: %d (next: %d)", amount, m.character.Experience, nextLevelXP)
					}
				}
			}
			return m, nil
		case "esc":
			m.characterStatsPanel.CancelEdit()
			m.message = "Edit cancelled"
			return m, nil
		default:
			// Pass key to input field
			return m, m.characterStatsPanel.HandleInput(msg)
		}
	}

	// In normal mode, allow class change
	switch msg.String() {
	case "c":
		// Backup current class state before opening selector
		m.pendingChanges.BackupClass(m.character)
		debug.Log("Backed up class state: %s", m.character.Class)
		m.classSelector.Show()
		m.message = "Select a class..."
		return m, nil
	case "l":
		// Open level-up selector (lowercase l)
		m.levelUpSelector.Show()
		m.message = "Level up your character..."
		return m, nil
	case "L":
		// Open de-level selector (Shift+L / uppercase L)
		if m.character.TotalLevel > 1 && len(m.character.Classes) > 0 {
			m.deLevelSelector.Show()
			m.message = "Select class to remove a level from..."
		} else {
			m.message = "Cannot de-level: Character is already at minimum level"
		}
		return m, nil
	case "]", "}":
		// Add 100 XP
		m.character.Experience += 100
		nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
		if m.character.Experience >= nextLevelXP {
			m.message = fmt.Sprintf("XP: %d (LEVEL UP AVAILABLE!)", m.character.Experience)
		} else {
			m.message = fmt.Sprintf("XP: %d (next level: %d)", m.character.Experience, nextLevelXP)
		}
		return m, nil
	case "[", "{":
		// Remove 100 XP (minimum 0)
		if m.character.Experience >= 100 {
			m.character.Experience -= 100
		} else {
			m.character.Experience = 0
		}
		nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
		m.message = fmt.Sprintf("XP: %d (next level: %d)", m.character.Experience, nextLevelXP)
		return m, nil
	}

	// Normal mode - handle actions
	switch msg.String() {
	case "n":
		m.characterStatsPanel.EditName()
		m.message = "Editing name..."
	case "r":
		// Short rest
		m.restPopup.Show(components.ShortRestType)
		m.message = "Taking a short rest..."
	case "R":
		// Long rest
		m.restPopup.Show(components.LongRestType)
		m.message = "Taking a long rest..."
	case "h":
		m.characterStatsPanel.EditHP()
		m.message = "Enter HP change (+/- amount)..."
	case "x":
		m.characterStatsPanel.EditXP()
		m.message = "Enter XP change (+/- amount)..."
	case "+", "=":
		m.characterStatsPanel.AddHP(1)
		m.message = fmt.Sprintf("HP: %d/%d", m.character.CurrentHP, m.character.MaxHP)
	case "-", "_":
		m.characterStatsPanel.RemoveHP(1)
		m.message = fmt.Sprintf("HP: %d/%d", m.character.CurrentHP, m.character.MaxHP)
	case "i":
		// Roll initiative
		initMod := m.characterStatsPanel.GetInitiativeModifier()
		expr := fmt.Sprintf("1d20%+d", initMod)
		m.dicePanel.Roll(expr)
		m.message = fmt.Sprintf("Initiative rolled: %s", m.dicePanel.LastMessage)
	case "I":
		// Toggle inspiration
		m.characterStatsPanel.ToggleInspiration()
		if m.character.Inspiration {
			m.message = "✨ Inspiration gained!"
		} else {
			m.message = "Inspiration used"
		}
		m.storage.Save(m.character)
	case "s":
		// Change species
		m.speciesSelector.Show()
		m.message = "Select a species..."
	case "S":
		// Save character (Shift+S)
		m.storage.Save(m.character)
		m.message = "Character saved!"
	}
	return m, nil
}

// handleStatGeneratorKeys handles stat generator specific keys (Stats panel related)
func (m *Model) handleStatGeneratorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if we're in editing mode for extras
	if m.statGenerator.IsVisible() {
		// Special handling for extras editing mode
		editingExtra := m.statGenerator.IsEditingExtra()

		switch msg.String() {
		case "up", "k":
			if !editingExtra {
				m.statGenerator.Prev()
			}
		case "down", "j":
			if !editingExtra {
				m.statGenerator.Next()
			}
		case "esc":
			// Check if in wizard mode
			if m.IsInWizard() {
				// In wizard, ESC cancels the entire wizard
				m.cancelWizard()
				return m, nil
			}

			// Cancel extra input or go back
			if editingExtra {
				m.statGenerator.CancelExtra()
			} else {
				m.statGenerator.GoBack()
				if !m.statGenerator.IsVisible() {
					m.message = "Stat generation cancelled"
				}
			}
		case "enter":
			// Save extra or continue to next step
			if editingExtra {
				m.statGenerator.SaveExtra()
			} else if m.statGenerator.CanContinue() {
				// Check if we're at the final step and on confirm button
				m.statGenerator.Continue()
				if !m.statGenerator.IsVisible() {
					// Apply stats and close
					m.statGenerator.ApplyToCharacter(m.character)
					m.storage.Save(m.character)

					// If in wizard mode, advance to next step
					if m.IsInWizard() {
						m.advanceWizardStep()
					} else {
						m.message = "Ability scores updated!"
					}
				}
			} else {
				m.message = "Please assign all stats before continuing"
			}
		case "e":
			// Edit extra in extras state
			if !editingExtra {
				m.statGenerator.StartEditingExtra()
			}
		case "1", "2", "3", "4", "5", "6":
			if !editingExtra {
				// Only assign stats for 4d6 and Standard Array methods
				method := m.statGenerator.GetMethod()
				state := m.statGenerator.GetState()
				if state == components.StateAssignStats &&
					(method == components.Method4d6DropLowest || method == components.MethodStandardArray) {
					idx := int(msg.String()[0] - '1')
					m.statGenerator.ToggleAssignment(idx)
				}
			}
		case "+", "=":
			if !editingExtra {
				// Increase in point buy state or extras
				state := m.statGenerator.GetState()
				if state == components.StateSetExtras {
					m.statGenerator.IncreaseExtra()
				} else {
					m.statGenerator.IncreasePointBuy()
				}
			}
		case "-", "_":
			if !editingExtra {
				// Decrease in point buy state or extras
				state := m.statGenerator.GetState()
				if state == components.StateSetExtras {
					m.statGenerator.DecreaseExtra()
				} else {
					m.statGenerator.DecreasePointBuy()
				}
			}
		case "backspace", "delete":
			// Delete character in extra input
			if editingExtra {
				m.statGenerator.DeleteExtraInput()
			}
		default:
			// Handle typing for extra input
			if editingExtra && len(msg.String()) == 1 {
				char := []rune(msg.String())[0]
				if (char >= '0' && char <= '9') || char == '+' || char == '-' {
					m.statGenerator.HandleExtraInput(char)
				}
			}
		}
	}
	return m, nil
}

// getLevelXP returns the XP required to reach a given level
func getLevelXP(level int) int {
	xpTable := map[int]int{
		1:  0,
		2:  300,
		3:  900,
		4:  2700,
		5:  6500,
		6:  14000,
		7:  23000,
		8:  34000,
		9:  48000,
		10: 64000,
		11: 85000,
		12: 100000,
		13: 120000,
		14: 140000,
		15: 165000,
		16: 195000,
		17: 225000,
		18: 265000,
		19: 305000,
		20: 355000,
	}
	if xp, ok := xpTable[level]; ok {
		return xp
	}
	return 355000 // Max level XP
}
