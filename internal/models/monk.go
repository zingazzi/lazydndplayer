// internal/models/monk.go
package models

import (
	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// MonkMechanics handles Monk-specific calculations and features
type MonkMechanics struct {
	character *Character
}

// NewMonkMechanics creates a new MonkMechanics instance
func NewMonkMechanics(char *Character) *MonkMechanics {
	return &MonkMechanics{character: char}
}

// GetMartialArtsDie returns the appropriate die for the character's Monk level
// Scaling: 1d6 (1-3), 1d8 (4-9), 1d10 (10-16), 1d12 (17-20)
func (m *MonkMechanics) GetMartialArtsDie() string {
	monkLevel := m.getMonkLevel()
	switch {
	case monkLevel >= 17:
		return "1d12" // Levels 17-20
	case monkLevel >= 10:
		return "1d10" // Levels 10-16
	case monkLevel >= 4:
		return "1d8" // Levels 4-9
	default:
		return "1d6" // Levels 1-3
	}
}

// CalculateUnarmoredAC calculates AC when not wearing armor
// Formula: 10 + Dexterity modifier + Wisdom modifier
func (m *MonkMechanics) CalculateUnarmoredAC() int {
	dexMod := m.character.AbilityScores.GetModifier("Dexterity")
	wisMod := m.character.AbilityScores.GetModifier("Wisdom")
	ac := 10 + dexMod + wisMod
	debug.Log("Monk Unarmored AC: 10 + %d (Dex) + %d (Wis) = %d", dexMod, wisMod, ac)
	return ac
}

// GetUnarmoredMovementBonus returns the speed bonus for the Monk's level
func (m *MonkMechanics) GetUnarmoredMovementBonus() int {
	monkLevel := m.getMonkLevel()
	switch {
	case monkLevel >= 18:
		return 30
	case monkLevel >= 14:
		return 25
	case monkLevel >= 10:
		return 20
	case monkLevel >= 6:
		return 15
	case monkLevel >= 2:
		return 10
	default:
		return 0
	}
}

// GetFocusPoints returns the current and max Focus Points
func (m *MonkMechanics) GetFocusPoints() (current int, max int) {
	for _, feature := range m.character.Features.Features {
		if feature.Name == "Focus Points" {
			return feature.CurrentUses, feature.MaxUses
		}
	}
	return 0, 0
}

// SpendFocusPoint spends Focus Points if available
func (m *MonkMechanics) SpendFocusPoint(count int) bool {
	for i := range m.character.Features.Features {
		if m.character.Features.Features[i].Name == "Focus Points" {
			if m.character.Features.Features[i].CurrentUses >= count {
				m.character.Features.Features[i].CurrentUses -= count
				debug.Log("Spent %d Focus Point(s). Remaining: %d/%d",
					count,
					m.character.Features.Features[i].CurrentUses,
					m.character.Features.Features[i].MaxUses)
				return true
			}
			debug.Log("Not enough Focus Points. Have %d, need %d",
				m.character.Features.Features[i].CurrentUses, count)
			return false
		}
	}
	debug.Log("Focus Points feature not found")
	return false
}

// RestoreFocusPoints restores Focus Points (capped at max)
func (m *MonkMechanics) RestoreFocusPoints(count int) {
	for i := range m.character.Features.Features {
		if m.character.Features.Features[i].Name == "Focus Points" {
			oldValue := m.character.Features.Features[i].CurrentUses
			m.character.Features.Features[i].CurrentUses += count
			if m.character.Features.Features[i].CurrentUses > m.character.Features.Features[i].MaxUses {
				m.character.Features.Features[i].CurrentUses = m.character.Features.Features[i].MaxUses
			}
			debug.Log("Restored %d Focus Point(s). %d -> %d (max: %d)",
				count, oldValue,
				m.character.Features.Features[i].CurrentUses,
				m.character.Features.Features[i].MaxUses)
			return
		}
	}
	debug.Log("Focus Points feature not found")
}

// RestoreAllFocusPoints restores all Focus Points to maximum
func (m *MonkMechanics) RestoreAllFocusPoints() {
	for i := range m.character.Features.Features {
		if m.character.Features.Features[i].Name == "Focus Points" {
			m.character.Features.Features[i].CurrentUses = m.character.Features.Features[i].MaxUses
			debug.Log("Restored all Focus Points to %d", m.character.Features.Features[i].MaxUses)
			return
		}
	}
}

// UseUncannyMetabolism uses the Uncanny Metabolism feature
// Returns: HP restored, FP restored, success
func (m *MonkMechanics) UseUncannyMetabolism() (hpRestored int, fpRestored int, success bool) {
	// Find Uncanny Metabolism feature
	for i := range m.character.Features.Features {
		if m.character.Features.Features[i].Name == "Uncanny Metabolism" {
			if m.character.Features.Features[i].CurrentUses > 0 {
				// Use the feature
				m.character.Features.Features[i].CurrentUses--
				debug.Log("Using Uncanny Metabolism")

				// Restore all Focus Points
				for j := range m.character.Features.Features {
					if m.character.Features.Features[j].Name == "Focus Points" {
						fpRestored = m.character.Features.Features[j].MaxUses - m.character.Features.Features[j].CurrentUses
						m.character.Features.Features[j].CurrentUses = m.character.Features.Features[j].MaxUses
						debug.Log("Restored %d Focus Points", fpRestored)
						break
					}
				}

				// Restore HP = Monk level + Martial Arts die (roll it)
				monkLevel := m.getMonkLevel()
				martialArtsDieRoll := m.rollMartialArtsDie()
				hpRestored = monkLevel + martialArtsDieRoll
				debug.Log("Healing: %d (monk level) + %d (martial arts die) = %d HP", monkLevel, martialArtsDieRoll, hpRestored)

				m.character.CurrentHP += hpRestored
				if m.character.CurrentHP > m.character.MaxHP {
					actualHealed := hpRestored - (m.character.CurrentHP - m.character.MaxHP)
					m.character.CurrentHP = m.character.MaxHP
					debug.Log("HP restored: %d (capped at max HP)", actualHealed)
				}

				return hpRestored, fpRestored, true
			}
			debug.Log("Uncanny Metabolism already used")
			return 0, 0, false
		}
	}
	debug.Log("Uncanny Metabolism feature not found")
	return 0, 0, false
}

// getMonkLevel returns the character's Monk class level
func (m *MonkMechanics) getMonkLevel() int {
	for _, classLevel := range m.character.Classes {
		if classLevel.ClassName == "Monk" {
			return classLevel.Level
		}
	}
	return 0
}

// rollMartialArtsDie rolls the Martial Arts die and returns the result
func (m *MonkMechanics) rollMartialArtsDie() int {
	die := m.GetMartialArtsDie()
	roller := GetDefaultDiceRoller()

	switch die {
	case "1d6":
		return roller.Roll(6)
	case "1d8":
		return roller.Roll(8)
	case "1d10":
		return roller.Roll(10)
	case "1d12":
		return roller.Roll(12)
	default:
		return roller.Roll(6)
	}
}

// UpdateFocusPointsOnLevelUp updates Focus Points max when leveling up
func (m *MonkMechanics) UpdateFocusPointsOnLevelUp() {
	monkLevel := m.getMonkLevel()
	for i := range m.character.Features.Features {
		if m.character.Features.Features[i].Name == "Focus Points" {
			oldMax := m.character.Features.Features[i].MaxUses
			m.character.Features.Features[i].MaxUses = monkLevel
			// Also increase current by the difference so they don't lose points
			diff := monkLevel - oldMax
			if diff > 0 {
				m.character.Features.Features[i].CurrentUses += diff
			}
			debug.Log("Updated Focus Points: %d -> %d (current: %d)",
				oldMax, monkLevel, m.character.Features.Features[i].CurrentUses)
			return
		}
	}
}
