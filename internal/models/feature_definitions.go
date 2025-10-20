// internal/models/feature_definitions.go
package models

import "fmt"

// FeatureDefinition defines a limited-use feature that can be granted by feats or species
type FeatureDefinition struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	MaxUses       string   `json:"max_uses"`       // Can be number or formula like "proficiency", "level"
	RestType      RestType `json:"rest_type"`
	UsesFormula   string   `json:"uses_formula"`   // Formula for max uses (e.g., "proficiency", "level", "1")
	EffectFormula string   `json:"effect_formula"` // Formula for effect (e.g., "2d6", "level")
}

// CalculateMaxUses calculates the maximum uses based on character level and proficiency
func (fd *FeatureDefinition) CalculateMaxUses(char *Character) int {
	switch fd.UsesFormula {
	case "proficiency":
		return char.ProficiencyBonus
	case "level":
		return char.Level
	case "":
		// If no formula, try parsing max_uses as a number
		if fd.MaxUses != "" {
			var uses int
			if _, err := fmt.Sscanf(fd.MaxUses, "%d", &uses); err == nil {
				return uses
			}
		}
		return 1
	default:
		// Try to parse as a number
		var uses int
		if _, err := fmt.Sscanf(fd.UsesFormula, "%d", &uses); err == nil {
			return uses
		}
		return 1
	}
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
	}
}
