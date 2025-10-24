// internal/models/feature_definitions.go
package models

import (
	"fmt"
	"strings"
)

// FeatureDefinition defines a limited-use feature that can be granted by feats or species
type FeatureDefinition struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	MaxUses       string                 `json:"max_uses"`       // Can be number or formula like "proficiency", "level"
	RestType      RestType               `json:"rest_type"`
	UsesFormula   string                 `json:"uses_formula"`   // Formula for max uses (e.g., "proficiency", "level", "1")
	EffectFormula string                 `json:"effect_formula"` // Formula for effect (e.g., "2d6", "level")
	Mechanics     map[string]interface{} `json:"mechanics"`      // Additional mechanics data (e.g., weapons_mastered, ki_points)
}

// CalculateMaxUses calculates the maximum uses based on character level and proficiency
func (fd *FeatureDefinition) CalculateMaxUses(char *Character) int {
	formula := fd.UsesFormula
	if formula == "" {
		formula = fd.MaxUses
	}

	switch formula {
	case "proficiency":
		return char.ProficiencyBonus
	case "level":
		return char.Level
	case "":
		return 1
	default:
		// Check if formula references a scaling table (ends with "_scaling")
		featureName := fd.Name
		if strings.HasSuffix(formula, "_scaling") {
			// Use the feature name for lookup
			if uses := GetFeatureScaling(char.Class, featureName, char.Level); uses > 0 {
				return uses
			}
		}

		// Try to parse as a number
		var uses int
		if _, err := fmt.Sscanf(formula, "%d", &uses); err == nil {
			return uses
		}

		// Try simple arithmetic formulas
		if parsed := parseFormulaExpression(formula, char); parsed > 0 {
			return parsed
		}

		return 1
	}
}

// parseFormulaExpression handles simple arithmetic formulas
func parseFormulaExpression(formula string, char *Character) int {
	// Handle ability modifier formulas
	switch formula {
	case "wisdom_mod":
		return char.AbilityScores.GetModifier("Wisdom")
	case "charisma_mod":
		return char.AbilityScores.GetModifier("Charisma")
	case "intelligence_mod":
		return char.AbilityScores.GetModifier("Intelligence")
	case "constitution_mod":
		return char.AbilityScores.GetModifier("Constitution")
	case "dexterity_mod":
		return char.AbilityScores.GetModifier("Dexterity")
	case "strength_mod":
		return char.AbilityScores.GetModifier("Strength")
	}

	// For now, return 0 for complex expressions we don't support yet
	// Future: Could implement a full expression parser
	return 0
}

// ToFeature converts a definition to an actual Feature instance
func (fd *FeatureDefinition) ToFeature(char *Character, source string) Feature {
	maxUses := fd.CalculateMaxUses(char)

	// Initialize CurrentUses to MaxUses (feature starts fully charged)
	currentUses := maxUses
	if maxUses == 0 {
		currentUses = 0 // Passive features stay at 0
	}

	return Feature{
		Name:        fd.Name,
		Description: fd.Description,
		MaxUses:     maxUses,
		CurrentUses: currentUses,
		RestType:    fd.RestType,
		Source:      source,
		Mechanics:   fd.Mechanics,
	}
}
