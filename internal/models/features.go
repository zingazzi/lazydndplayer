// internal/models/features.go
package models

import (
	"strconv"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// RestType indicates when a feature recharges
type RestType string

const (
	ShortRest RestType = "Short Rest"
	LongRest  RestType = "Long Rest"
	Daily     RestType = "Daily"
	None      RestType = "None" // Passive features
)

// Feature represents a limited-use ability (class features, racial abilities, etc.)
type Feature struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	MaxUses     int                    `json:"max_uses"`      // Maximum uses per rest
	CurrentUses int                    `json:"current_uses"`  // Current available uses
	RestType    RestType               `json:"rest_type"`     // When it recharges
	Source      string                 `json:"source"`        // e.g., "Class: Barbarian", "Species: Dragonborn"
	Mechanics   map[string]interface{} `json:"mechanics"`     // Additional mechanics data (e.g., weapons_mastered, ki_points)
	UsesFormula string                 `json:"uses_formula,omitempty"` // Original formula for recalculating (e.g., "wisdom_mod", "proficiency")
}

// FeatureList manages character features
type FeatureList struct {
	Features []Feature `json:"features"`
}

// NewFeatureList creates an empty feature list
func NewFeatureList() *FeatureList {
	return &FeatureList{
		Features: []Feature{},
	}
}

// AddFeature adds a new feature
func (fl *FeatureList) AddFeature(feature Feature) {
	// Set current uses to max if not specified
	if feature.CurrentUses == 0 && feature.MaxUses > 0 {
		feature.CurrentUses = feature.MaxUses
	}
	fl.Features = append(fl.Features, feature)
}

// UseFeature decrements the usage count for a feature
func (fl *FeatureList) UseFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}

	feature := &fl.Features[index]
	if feature.CurrentUses > 0 {
		feature.CurrentUses--

		// Check if this is Bardic Inspiration and restore Beguiling Magic uses
		if feature.Name == "Bardic Inspiration" {
			fl.restoreBeguilingMagicOnBardicInspiration()
		}

		return true
	}
	return false
}

// restoreBeguilingMagicOnBardicInspiration restores Beguiling Magic uses when Bardic Inspiration is used
func (fl *FeatureList) restoreBeguilingMagicOnBardicInspiration() {
	for i := range fl.Features {
		if fl.Features[i].Name == "Beguiling Magic" {
			if fl.Features[i].Mechanics != nil {
				if restoreOnBardic, ok := fl.Features[i].Mechanics["restore_on_bardic_inspiration"].(bool); ok && restoreOnBardic {
					// Restore Beguiling Magic uses
					fl.Features[i].CurrentUses = fl.Features[i].MaxUses
					debug.Log("Beguiling Magic uses restored from Bardic Inspiration: %d/%d",
						fl.Features[i].CurrentUses, fl.Features[i].MaxUses)
					break
				}
			}
		}
	}
}

// RestoreFeature restores one use of a feature
func (fl *FeatureList) RestoreFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}

	feature := &fl.Features[index]
	if feature.CurrentUses < feature.MaxUses {
		feature.CurrentUses++
		return true
	}
	return false
}

// ShortRestRecover recovers features that recharge on short rest
func (fl *FeatureList) ShortRestRecover() {
	for i := range fl.Features {
		if fl.Features[i].RestType == ShortRest || fl.Features[i].RestType == Daily {
			fl.Features[i].CurrentUses = fl.Features[i].MaxUses
		}
	}
}

// LongRestRecover recovers all rechargeable features
func (fl *FeatureList) LongRestRecover() {
	for i := range fl.Features {
		if fl.Features[i].RestType != None {
			fl.Features[i].CurrentUses = fl.Features[i].MaxUses
		}
	}
}

// RemoveFeature removes a feature by index
func (fl *FeatureList) RemoveFeature(index int) bool {
	if index < 0 || index >= len(fl.Features) {
		return false
	}
	fl.Features = append(fl.Features[:index], fl.Features[index+1:]...)
	return true
}

// UpdateFormulaBasedFeatures recalculates max uses for features that have formulas
func (fl *FeatureList) UpdateFormulaBasedFeatures(char *Character) {
	for i := range fl.Features {
		feature := &fl.Features[i]
		if feature.UsesFormula == "" {
			continue // No formula to recalculate
		}

		// Calculate new max uses
		newMaxUses := calculateFeatureUsesFromFormula(feature.UsesFormula, char)

		// Only update if the max changed
		if newMaxUses != feature.MaxUses {
			oldMax := feature.MaxUses
			feature.MaxUses = newMaxUses

			// Adjust current uses proportionally if feature was partially used
			// If at max, keep at max. Otherwise, try to maintain ratio
			if oldMax > 0 {
				if feature.CurrentUses >= oldMax {
					// Was at max, keep at new max
					feature.CurrentUses = newMaxUses
				} else if feature.CurrentUses > 0 {
					// Was partially used, maintain ratio
					ratio := float64(feature.CurrentUses) / float64(oldMax)
					feature.CurrentUses = int(float64(newMaxUses) * ratio)
					if feature.CurrentUses < 1 && newMaxUses > 0 {
						feature.CurrentUses = 1 // At least 1 if there are any uses available
					}
				}
			} else {
				// Was at 0, initialize to max
				feature.CurrentUses = newMaxUses
			}

			// Don't let current exceed max
			if feature.CurrentUses > feature.MaxUses {
				feature.CurrentUses = feature.MaxUses
			}
		}
	}
}

// calculateFeatureUsesFromFormula calculates uses from a formula string
func calculateFeatureUsesFromFormula(formula string, char *Character) int {
	formula = strings.TrimSpace(strings.ToLower(formula))

	switch formula {
	case "proficiency":
		return char.ProficiencyBonus
	case "level":
		return char.Level
	case "wisdom_mod":
		mod := char.AbilityScores.GetModifier("Wisdom")
		if mod < 1 {
			return 1 // Minimum of 1
		}
		return mod
	case "charisma_mod":
		mod := char.AbilityScores.GetModifier("Charisma")
		if mod < 1 {
			return 1
		}
		return mod
	case "intelligence_mod":
		mod := char.AbilityScores.GetModifier("Intelligence")
		if mod < 1 {
			return 1
		}
		return mod
	case "constitution_mod":
		mod := char.AbilityScores.GetModifier("Constitution")
		if mod < 1 {
			return 1
		}
		return mod
	case "dexterity_mod":
		mod := char.AbilityScores.GetModifier("Dexterity")
		if mod < 1 {
			return 1
		}
		return mod
	case "strength_mod":
		mod := char.AbilityScores.GetModifier("Strength")
		if mod < 1 {
			return 1
		}
		return mod
	default:
		// Try to parse as number
		if uses, err := strconv.Atoi(formula); err == nil {
			return uses
		}
		return 1
	}
}
