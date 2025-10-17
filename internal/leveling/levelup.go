// internal/leveling/levelup.go
package leveling

import (
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// LevelUpOptions represents choices during level up
type LevelUpOptions struct {
	HPIncrease       int                   // HP gained this level
	AbilityIncrease  map[models.AbilityType]int // Ability score improvements
	NewSpellSlots    models.SpellSlots     // New spell slots gained
	ClassFeatures    []string              // New class features gained
}

// CalculateHPOptions returns options for HP increase
func CalculateHPOptions(class string, conMod int) (rolled int, average int) {
	hitDie := GetHitDie(class)
	average = (hitDie / 2) + 1 + conMod
	// For actual rolled value, this would be done by the UI with dice roller
	return hitDie + conMod, average
}

// GetHitDie returns the hit die for a class
func GetHitDie(class string) int {
	hitDice := map[string]int{
		"Barbarian": 12,
		"Fighter":   10,
		"Paladin":   10,
		"Ranger":    10,
		"Bard":      8,
		"Cleric":    8,
		"Druid":     8,
		"Monk":      8,
		"Rogue":     8,
		"Warlock":   8,
		"Sorcerer":  6,
		"Wizard":    6,
	}

	if die, ok := hitDice[class]; ok {
		return die
	}
	return 8 // Default
}

// CanIncreaseAbilityScores checks if character can increase ability scores
func CanIncreaseAbilityScores(newLevel int) bool {
	// Every 4 levels (4, 8, 12, 16, 19 for fighters, 20 for rogues)
	return newLevel%4 == 0
}

// GetSpellSlotsForLevel returns spell slots for a given class and level
func GetSpellSlotsForLevel(class string, level int) models.SpellSlots {
	// Simplified spell slot progression for full casters
	// This is for Wizard/Cleric/Druid/Bard/Sorcerer
	fullCasterSlots := map[int]models.SpellSlots{
		1: {Level1: models.SpellSlot{Maximum: 2, Current: 2}},
		2: {Level1: models.SpellSlot{Maximum: 3, Current: 3}},
		3: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 2, Current: 2},
		},
		4: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
		},
		5: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 2, Current: 2},
		},
		6: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
		},
		7: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
			Level4: models.SpellSlot{Maximum: 1, Current: 1},
		},
		8: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
			Level4: models.SpellSlot{Maximum: 2, Current: 2},
		},
		9: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
			Level4: models.SpellSlot{Maximum: 3, Current: 3},
			Level5: models.SpellSlot{Maximum: 1, Current: 1},
		},
		10: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
			Level4: models.SpellSlot{Maximum: 3, Current: 3},
			Level5: models.SpellSlot{Maximum: 2, Current: 2},
		},
		11: {
			Level1: models.SpellSlot{Maximum: 4, Current: 4},
			Level2: models.SpellSlot{Maximum: 3, Current: 3},
			Level3: models.SpellSlot{Maximum: 3, Current: 3},
			Level4: models.SpellSlot{Maximum: 3, Current: 3},
			Level5: models.SpellSlot{Maximum: 2, Current: 2},
			Level6: models.SpellSlot{Maximum: 1, Current: 1},
		},
		// Continue for levels 12-20...
	}

	// Check if class is a full caster
	fullCasters := map[string]bool{
		"Wizard":   true,
		"Cleric":   true,
		"Druid":    true,
		"Bard":     true,
		"Sorcerer": true,
	}

	if fullCasters[class] {
		if slots, ok := fullCasterSlots[level]; ok {
			return slots
		}
	}

	// Return empty slots for non-casters or levels not defined
	return models.SpellSlots{}
}

// GetClassFeatures returns class features gained at a specific level
func GetClassFeatures(class string, level int) []string {
	// Simplified class features
	features := map[string]map[int][]string{
		"Fighter": {
			1: {"Fighting Style", "Second Wind"},
			2: {"Action Surge"},
			3: {"Martial Archetype"},
			5: {"Extra Attack"},
			9: {"Indomitable"},
		},
		"Wizard": {
			1: {"Spellcasting", "Arcane Recovery"},
			2: {"Arcane Tradition"},
			6: {"Arcane Tradition Feature"},
			10: {"Arcane Tradition Feature"},
		},
		"Rogue": {
			1: {"Expertise", "Sneak Attack", "Thieves' Cant"},
			2: {"Cunning Action"},
			3: {"Roguish Archetype"},
			5: {"Uncanny Dodge"},
			7: {"Evasion"},
		},
	}

	if classFeatures, ok := features[class]; ok {
		if levelFeatures, ok := classFeatures[level]; ok {
			return levelFeatures
		}
	}

	return []string{}
}

// PerformLevelUp applies level up changes to character
func PerformLevelUp(char *models.Character, options LevelUpOptions) {
	// Increase level
	char.Level++

	// Increase HP
	char.MaxHP += options.HPIncrease
	char.CurrentHP = char.MaxHP // Fully heal on level up

	// Apply ability score increases
	for ability, increase := range options.AbilityIncrease {
		currentScore := char.AbilityScores.GetScore(ability)
		char.AbilityScores.SetScore(ability, currentScore+increase)
	}

	// Update spell slots
	char.SpellBook.Slots = options.NewSpellSlots

	// Update derived stats
	char.UpdateDerivedStats()
}

// ValidateAbilityIncreases ensures ability increases are valid
func ValidateAbilityIncreases(increases map[models.AbilityType]int) bool {
	total := 0
	for _, increase := range increases {
		total += increase
		if increase < 0 {
			return false
		}
	}
	// Standard ASI gives 2 points total
	return total == 2
}
