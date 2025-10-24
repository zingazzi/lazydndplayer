// internal/models/feature_scaling.go
package models

// FeatureScaling defines level-based scaling for class features that don't follow simple formulas
// Map structure: [ClassName][FeatureName][Level] = Uses
var FeatureScaling = map[string]map[string]map[int]int{
	"Barbarian": {
		"Rage": {
			1: 2, 2: 2,
			3: 3, 4: 3, 5: 3,
			6: 4, 7: 4, 8: 4, 9: 4, 10: 4, 11: 4,
			12: 5, 13: 5, 14: 5, 15: 5, 16: 5,
			17: 6, 18: 6, 19: 6,
			20: 999, // Unlimited at level 20
		},
	},
	"Fighter": {
		"Action Surge": {
			2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1, 10: 1,
			11: 1, 12: 1, 13: 1, 14: 1, 15: 1, 16: 1,
			17: 2, 18: 2, 19: 2, 20: 2,
		},
		"Second Wind": {
			1: 2, 2: 2, // 2 uses at levels 1-2
			3: 3, 4: 3, 5: 3, 6: 3, 7: 3, 8: 3, 9: 3, // 3 uses at levels 3-9
			10: 4, 11: 4, 12: 4, 13: 4, 14: 4, 15: 4, 16: 4, 17: 4, 18: 4, 19: 4, 20: 4, // 4 uses at levels 10-20
		},
		"Indomitable": {
			9: 1, 10: 1, 11: 1, 12: 1,
			13: 2, 14: 2, 15: 2, 16: 2,
			17: 3, 18: 3, 19: 3, 20: 3,
		},
	},
	"Psi Warrior": {
		"Psionic Power": {
			3: 4, 4: 4, // 4 dice (d6)
			5: 6, 6: 6, 7: 6, 8: 6, // 6 dice (d8)
			9: 8, 10: 8, // 8 dice (d8)
			11: 8, 12: 8, // 8 dice (d10)
			13: 10, 14: 10, 15: 10, 16: 10, // 10 dice (d10)
			17: 12, 18: 12, 19: 12, 20: 12, // 12 dice (d12)
		},
	},
	"Battle Master": {
		"Combat Superiority": {
			3: 4, 4: 4, 5: 4, 6: 4, // 4 dice (d8)
			7: 5, 8: 5, 9: 5, // 5 dice (d8)
			10: 5, 11: 5, 12: 5, 13: 5, 14: 5, // 5 dice (d10)
			15: 6, 16: 6, 17: 6, // 6 dice (d10)
			18: 6, 19: 6, 20: 6, // 6 dice (d12)
		},
	},
	"Monk": {
		"Ki": {
			2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10,
			11: 11, 12: 12, 13: 13, 14: 14, 15: 15, 16: 16, 17: 17, 18: 18, 19: 19, 20: 20,
		},
	},
	"Bard": {
		"Bardic Inspiration": {
			1: 0, // Not available yet
			2: 3, 3: 3, 4: 3, // Proficiency bonus (assume +2)
			5: 3, 6: 3, 7: 3, 8: 3, // Proficiency bonus +3
			9: 4, 10: 4, 11: 4, 12: 4, // Proficiency bonus +4
			13: 5, 14: 5, 15: 5, 16: 5, // Proficiency bonus +5
			17: 6, 18: 6, 19: 6, 20: 6, // Proficiency bonus +6
		},
	},
	"Druid": {
		"Wild Shape": {
			2: 2, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2, 8: 2, 9: 2, 10: 2,
			11: 2, 12: 2, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2, 18: 2, 19: 2, 20: 2,
		},
	},
	"Cleric": {
		"Channel Divinity": {
			2: 1, 3: 1, 4: 1, 5: 1,
			6: 2, 7: 2, 8: 2, 9: 2, 10: 2, 11: 2, 12: 2, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2,
			18: 3, 19: 3, 20: 3,
		},
	},
	"Paladin": {
		"Channel Divinity": {
			3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1, 16: 1, 17: 1, 18: 1, 19: 1, 20: 1,
		},
		"Lay on Hands": {
			// This is a pool, not uses. Value = 5 × level
			1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30, 7: 35, 8: 40, 9: 45, 10: 50,
			11: 55, 12: 60, 13: 65, 14: 70, 15: 75, 16: 80, 17: 85, 18: 90, 19: 95, 20: 100,
		},
	},
	"Sorcerer": {
		"Sorcery Points": {
			2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10,
			11: 11, 12: 12, 13: 13, 14: 14, 15: 15, 16: 16, 17: 17, 18: 18, 19: 19, 20: 20,
		},
	},
	"Wizard": {
		"Arcane Recovery": {
			// Once per day, recovers spell slots up to half wizard level (rounded up)
			1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1, 10: 1,
			11: 1, 12: 1, 13: 1, 14: 1, 15: 1, 16: 1, 17: 1, 18: 1, 19: 1, 20: 1,
		},
	},
	"Warlock": {
		// Warlock spell slots recovered on short rest - handled separately in spellcasting
	},
	"Ranger": {
		// Most ranger features are passive or use spell slots
	},
	"Rogue": {
		// Sneak Attack is per-turn, not limited use
	},
}

// GetFeatureScaling looks up the scaled value for a feature at a specific level
func GetFeatureScaling(className string, featureName string, level int) int {
	if classMap, ok := FeatureScaling[className]; ok {
		if featureMap, ok := classMap[featureName]; ok {
			if uses, ok := featureMap[level]; ok {
				return uses
			}
		}
	}
	return 0 // Not found in scaling table
}

// GetRestTypeForFeature returns the rest type for a feature (some features change rest type at higher levels)
func GetRestTypeForFeature(className string, featureName string, level int) RestType {
	// Special cases where rest type changes with level
	if className == "Bard" && featureName == "Bardic Inspiration" {
		if level >= 5 {
			return ShortRest
		}
		return LongRest
	}

	// Default - will be overridden by JSON definition
	return None
}

// GetPsiDiceSize returns the die size for Psi Warrior at a given Fighter level
func GetPsiDiceSize(fighterLevel int) string {
	switch {
	case fighterLevel >= 17:
		return "d12"
	case fighterLevel >= 11:
		return "d10"
	case fighterLevel >= 5:
		return "d8"
	default:
		return "d6"
	}
}

// GetSuperiorityDiceSize returns the die size for Battle Master at a given Fighter level
func GetSuperiorityDiceSize(fighterLevel int) string {
	switch {
	case fighterLevel >= 18:
		return "d12"
	case fighterLevel >= 10:
		return "d10"
	default:
		return "d8"
	}
}
