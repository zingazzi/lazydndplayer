// internal/models/spellcasting_types.go
package models

import (
	"math"
)

// CasterType represents different spellcasting progression types
type CasterType string

const (
	FullCaster  CasterType = "full"    // Bard, Cleric, Druid, Sorcerer, Wizard
	HalfCaster  CasterType = "half"    // Paladin, Ranger
	ThirdCaster CasterType = "third"   // Eldritch Knight, Arcane Trickster
	PactCaster  CasterType = "pact"    // Warlock (Pact Magic - separate system)
	NonCaster   CasterType = "none"    // Non-spellcasting classes
)

// SpellcastingMethod represents how a caster learns and prepares spells
type SpellcastingMethod string

const (
	PreparedCaster SpellcastingMethod = "prepared" // Cleric, Druid, Wizard - know all, prepare subset
	KnownCaster    SpellcastingMethod = "known"    // Bard, Sorcerer, Ranger, Warlock - know limited spells
	SpellbookCaster SpellcastingMethod = "spellbook" // Wizard - prepare from spellbook
)

// ClassSpellcastingInfo contains spellcasting progression for a class
type ClassSpellcastingInfo struct {
	ClassName           string
	CasterType          CasterType
	Method              SpellcastingMethod
	SpellcastingAbility AbilityType
	RitualCasting       bool
	// For known casters
	CantripsKnownByLevel map[int]int // level -> cantrips known
	SpellsKnownByLevel   map[int]int // level -> spells known
	// For prepared casters
	PreparationFormula string // e.g., "wisdom+level"
}

// GetClassCasterInfo returns spellcasting info for a class
func GetClassCasterInfo(className string) *ClassSpellcastingInfo {
	casterInfo := map[string]*ClassSpellcastingInfo{
		"Bard": {
			ClassName:           "Bard",
			CasterType:          FullCaster,
			Method:              KnownCaster,
			SpellcastingAbility: Charisma,
			RitualCasting:       true,
			CantripsKnownByLevel: map[int]int{
				1: 2, 2: 2, 3: 2, 4: 3, 5: 3, 6: 3, 7: 3, 8: 3, 9: 3, 10: 4,
				11: 4, 12: 4, 13: 4, 14: 4, 15: 4, 16: 4, 17: 4, 18: 4, 19: 4, 20: 4,
			},
			SpellsKnownByLevel: map[int]int{
				1: 4, 2: 5, 3: 6, 4: 7, 5: 8, 6: 9, 7: 10, 8: 11, 9: 12, 10: 14,
				11: 15, 12: 15, 13: 16, 14: 18, 15: 19, 16: 19, 17: 20, 18: 22, 19: 22, 20: 22,
			},
		},
		"Cleric": {
			ClassName:           "Cleric",
			CasterType:          FullCaster,
			Method:              PreparedCaster,
			SpellcastingAbility: Wisdom,
			RitualCasting:       true,
			PreparationFormula:  "wisdom+level",
			CantripsKnownByLevel: map[int]int{
				1: 3, 2: 3, 3: 3, 4: 4, 5: 4, 6: 4, 7: 4, 8: 4, 9: 4, 10: 5,
				11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
			},
		},
		"Druid": {
			ClassName:           "Druid",
			CasterType:          FullCaster,
			Method:              PreparedCaster,
			SpellcastingAbility: Wisdom,
			RitualCasting:       true,
			PreparationFormula:  "wisdom+level",
			CantripsKnownByLevel: map[int]int{
				1: 2, 2: 2, 3: 2, 4: 3, 5: 3, 6: 3, 7: 3, 8: 3, 9: 3, 10: 4,
				11: 4, 12: 4, 13: 4, 14: 4, 15: 4, 16: 4, 17: 4, 18: 4, 19: 4, 20: 4,
			},
		},
		"Sorcerer": {
			ClassName:           "Sorcerer",
			CasterType:          FullCaster,
			Method:              KnownCaster,
			SpellcastingAbility: Charisma,
			RitualCasting:       false,
			CantripsKnownByLevel: map[int]int{
				1: 4, 2: 4, 3: 4, 4: 5, 5: 5, 6: 5, 7: 5, 8: 5, 9: 5, 10: 6,
				11: 6, 12: 6, 13: 6, 14: 6, 15: 6, 16: 6, 17: 6, 18: 6, 19: 6, 20: 6,
			},
			SpellsKnownByLevel: map[int]int{
				1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10, 10: 11,
				11: 12, 12: 12, 13: 13, 14: 13, 15: 14, 16: 14, 17: 15, 18: 15, 19: 15, 20: 15,
			},
		},
		"Wizard": {
			ClassName:           "Wizard",
			CasterType:          FullCaster,
			Method:              SpellbookCaster,
			SpellcastingAbility: Intelligence,
			RitualCasting:       true,
			PreparationFormula:  "intelligence+level",
			CantripsKnownByLevel: map[int]int{
				1: 3, 2: 3, 3: 3, 4: 4, 5: 4, 6: 4, 7: 4, 8: 4, 9: 4, 10: 5,
				11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
			},
		},
		"Warlock": {
			ClassName:           "Warlock",
			CasterType:          PactCaster,
			Method:              KnownCaster,
			SpellcastingAbility: Charisma,
			RitualCasting:       false,
			CantripsKnownByLevel: map[int]int{
				1: 2, 2: 2, 3: 2, 4: 3, 5: 3, 6: 3, 7: 3, 8: 3, 9: 3, 10: 4,
				11: 4, 12: 4, 13: 4, 14: 4, 15: 4, 16: 4, 17: 4, 18: 4, 19: 4, 20: 4,
			},
			SpellsKnownByLevel: map[int]int{
				1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10, 10: 10,
				11: 11, 12: 11, 13: 12, 14: 12, 15: 13, 16: 13, 17: 14, 18: 14, 19: 15, 20: 15,
			},
		},
		"Paladin": {
			ClassName:           "Paladin",
			CasterType:          HalfCaster,
			Method:              PreparedCaster,
			SpellcastingAbility: Charisma,
			RitualCasting:       false,
			PreparationFormula:  "charisma+(level/2)",
		},
		"Ranger": {
			ClassName:           "Ranger",
			CasterType:          HalfCaster,
			Method:              KnownCaster,
			SpellcastingAbility: Wisdom,
			RitualCasting:       false,
			SpellsKnownByLevel: map[int]int{
				2: 2, 3: 3, 4: 3, 5: 4, 6: 4, 7: 5, 8: 5, 9: 6, 10: 6,
				11: 7, 12: 7, 13: 8, 14: 8, 15: 9, 16: 9, 17: 10, 18: 10, 19: 11, 20: 11,
			},
		},
	}

	if info, exists := casterInfo[className]; exists {
		return info
	}
	return nil
}

// CalculateMulticlassSpellSlots calculates spell slots for a multiclass character
// following D&D 5e multiclass spellcasting rules
func CalculateMulticlassSpellSlots(classes []ClassLevel) SpellSlots {
	var fullCasterLevels int
	var halfCasterLevels int
	var thirdCasterLevels int
	var warlockLevels int

	// Sum up caster levels by type
	for _, cl := range classes {
		info := GetClassCasterInfo(cl.ClassName)
		if info == nil {
			continue
		}

		switch info.CasterType {
		case FullCaster:
			fullCasterLevels += cl.Level
		case HalfCaster:
			halfCasterLevels += cl.Level
		case ThirdCaster:
			thirdCasterLevels += cl.Level
		case PactCaster:
			warlockLevels += cl.Level
		}
	}

	// Calculate effective caster level
	// Full casters count as full level, half as level/2, third as level/3
	effectiveCasterLevel := fullCasterLevels
	effectiveCasterLevel += int(math.Floor(float64(halfCasterLevels) / 2.0))
	effectiveCasterLevel += int(math.Floor(float64(thirdCasterLevels) / 3.0))

	// Warlock (Pact Magic) doesn't combine with other casters
	// It's tracked separately
	slots := GetSpellSlotsByLevel(effectiveCasterLevel)

	return slots
}

// GetSpellSlotsByLevel returns the standard spell slots for a given caster level
func GetSpellSlotsByLevel(level int) SpellSlots {
	if level <= 0 {
		return SpellSlots{}
	}

	// Standard D&D 5e spell slot progression for full casters
	slotTable := []SpellSlots{
		{}, // Level 0 (invalid)
		{Level1: SpellSlot{Maximum: 2, Current: 2}},   // Level 1
		{Level1: SpellSlot{Maximum: 3, Current: 3}},   // Level 2
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 2, Current: 2}}, // Level 3
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}}, // Level 4
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 2, Current: 2}}, // Level 5
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}}, // Level 6
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 1, Current: 1}}, // Level 7
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 2, Current: 2}}, // Level 8
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 1, Current: 1}}, // Level 9
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}}, // Level 10
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}}, // Level 11
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}}, // Level 12
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}}, // Level 13
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}}, // Level 14
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}, Level8: SpellSlot{Maximum: 1, Current: 1}}, // Level 15
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}, Level8: SpellSlot{Maximum: 1, Current: 1}}, // Level 16
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 2, Current: 2}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}, Level8: SpellSlot{Maximum: 1, Current: 1}, Level9: SpellSlot{Maximum: 1, Current: 1}}, // Level 17
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 3, Current: 3}, Level6: SpellSlot{Maximum: 1, Current: 1}, Level7: SpellSlot{Maximum: 1, Current: 1}, Level8: SpellSlot{Maximum: 1, Current: 1}, Level9: SpellSlot{Maximum: 1, Current: 1}}, // Level 18
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 3, Current: 3}, Level6: SpellSlot{Maximum: 2, Current: 2}, Level7: SpellSlot{Maximum: 1, Current: 1}, Level8: SpellSlot{Maximum: 1, Current: 1}, Level9: SpellSlot{Maximum: 1, Current: 1}}, // Level 19
		{Level1: SpellSlot{Maximum: 4, Current: 4}, Level2: SpellSlot{Maximum: 3, Current: 3}, Level3: SpellSlot{Maximum: 3, Current: 3}, Level4: SpellSlot{Maximum: 3, Current: 3}, Level5: SpellSlot{Maximum: 3, Current: 3}, Level6: SpellSlot{Maximum: 2, Current: 2}, Level7: SpellSlot{Maximum: 2, Current: 2}, Level8: SpellSlot{Maximum: 1, Current: 1}, Level9: SpellSlot{Maximum: 1, Current: 1}}, // Level 20
	}

	if level > 20 {
		level = 20
	}

	return slotTable[level]
}

// GetWarlockPactSlots returns Warlock-specific Pact Magic slots
func GetWarlockPactSlots(level int) (slots int, slotLevel int) {
	if level <= 0 {
		return 0, 0
	}

	// Warlock Pact Magic progression
	warlockTable := []struct {
		slots int
		level int
	}{
		{0, 0},  // Level 0
		{1, 1},  // Level 1
		{2, 1},  // Level 2
		{2, 2},  // Level 3
		{2, 2},  // Level 4
		{2, 3},  // Level 5
		{2, 3},  // Level 6
		{2, 4},  // Level 7
		{2, 4},  // Level 8
		{2, 5},  // Level 9
		{2, 5},  // Level 10
		{3, 5},  // Level 11
		{3, 5},  // Level 12
		{3, 5},  // Level 13
		{3, 5},  // Level 14
		{3, 5},  // Level 15
		{3, 5},  // Level 16
		{4, 5},  // Level 17
		{4, 5},  // Level 18
		{4, 5},  // Level 19
		{4, 5},  // Level 20
	}

	if level > 20 {
		level = 20
	}

	return warlockTable[level].slots, warlockTable[level].level
}

// IsSpellcaster checks if a character has any spellcasting classes
func IsSpellcaster(classes []ClassLevel) bool {
	for _, cl := range classes {
		info := GetClassCasterInfo(cl.ClassName)
		if info != nil && info.CasterType != NonCaster {
			return true
		}
	}
	return false
}

// GetSpellcastingAbility returns the spellcasting ability for a character's primary spellcasting class
func GetSpellcastingAbility(classes []ClassLevel) AbilityType {
	// Return the spellcasting ability of the first spellcasting class
	for _, cl := range classes {
		info := GetClassCasterInfo(cl.ClassName)
		if info != nil && info.CasterType != NonCaster {
			return info.SpellcastingAbility
		}
	}
	return ""
}

// CanCastRituals checks if a character can cast rituals
func CanCastRituals(classes []ClassLevel) bool {
	for _, cl := range classes {
		info := GetClassCasterInfo(cl.ClassName)
		if info != nil && info.RitualCasting {
			return true
		}
	}
	return false
}

// GetMaxCantripsKnown returns the total number of cantrips a multiclass character knows
func GetMaxCantripsKnown(classes []ClassLevel) int {
	maxCantrips := 0
	for _, cl := range classes {
		info := GetClassCasterInfo(cl.ClassName)
		if info != nil && info.CantripsKnownByLevel != nil {
			if cantrips, exists := info.CantripsKnownByLevel[cl.Level]; exists {
				maxCantrips += cantrips
			}
		}
	}
	return maxCantrips
}

// SpellcastingSummary provides a summary of character's spellcasting capabilities
type SpellcastingSummary struct {
	IsSpellcaster       bool
	SpellcastingAbility AbilityType
	CanCastRituals      bool
	SpellSlots          SpellSlots
	WarlockSlots        int // Pact Magic slots
	WarlockSlotLevel    int // Level of Pact Magic slots
	MaxCantrips         int
	Method              SpellcastingMethod
}

// GetSpellcastingSummary returns a summary of the character's spellcasting capabilities
func GetSpellcastingSummary(char *Character) SpellcastingSummary {
	summary := SpellcastingSummary{
		IsSpellcaster:       IsSpellcaster(char.Classes),
		SpellcastingAbility: GetSpellcastingAbility(char.Classes),
		CanCastRituals:      CanCastRituals(char.Classes),
		MaxCantrips:         GetMaxCantripsKnown(char.Classes),
	}

	// Calculate spell slots for non-Warlock spellcasters
	summary.SpellSlots = CalculateMulticlassSpellSlots(char.Classes)

	// Check for Warlock levels (Pact Magic)
	warlockLevel := char.GetClassLevel("Warlock")
	if warlockLevel > 0 {
		summary.WarlockSlots, summary.WarlockSlotLevel = GetWarlockPactSlots(warlockLevel)
	}

	// Determine primary spellcasting method
	if char.HasClass("Wizard") {
		summary.Method = SpellbookCaster
	} else if char.HasClass("Cleric") || char.HasClass("Druid") || char.HasClass("Paladin") {
		summary.Method = PreparedCaster
	} else {
		summary.Method = KnownCaster
	}

	return summary
}
