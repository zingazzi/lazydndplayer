// internal/models/feature_validation.go
package models

import (
	"fmt"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// DuplicateFeatureRule defines how to handle duplicate features
type DuplicateFeatureRule string

const (
	// Stack: Multiple instances stack (e.g., HP bonuses)
	FeatureStack DuplicateFeatureRule = "stack"
	// Replace: New instance replaces old (e.g., base AC calculations)
	FeatureReplace DuplicateFeatureRule = "replace"
	// Keep: Keep existing, ignore new (e.g., proficiencies)
	FeatureKeep DuplicateFeatureRule = "keep"
	// UseHighest: Keep the highest value (e.g., Extra Attack)
	FeatureUseHighest DuplicateFeatureRule = "highest"
	// AllowMultiple: Allow multiple distinct instances (e.g., Fighting Styles)
	FeatureAllowMultiple DuplicateFeatureRule = "multiple"
)

// FeatureRules defines rules for handling duplicate features
var FeatureRules = map[string]DuplicateFeatureRule{
	// Combat features
	"Extra Attack":        FeatureUseHighest, // Don't stack, use highest
	"Action Surge":        FeatureUseHighest, // Don't stack
	"Second Wind":         FeatureKeep,       // Don't gain twice
	"Rage":                FeatureKeep,       // Don't gain twice
	"Unarmored Defense":   FeatureKeep,       // Only one AC calculation method
	"Fighting Style":      FeatureAllowMultiple, // Can have different styles from different classes

	// Spellcasting features
	"Spellcasting":        FeatureAllowMultiple, // Combine spell lists
	"Pact Magic":          FeatureAllowMultiple, // Warlock multiclass

	// HP bonuses
	"Durable":             FeatureStack, // HP bonuses stack
	"Tough":               FeatureStack, // HP bonuses stack

	// Proficiencies
	"Armor Proficiency":   FeatureKeep, // Already have
	"Weapon Proficiency":  FeatureKeep, // Already have
	"Tool Proficiency":    FeatureKeep, // Already have
	"Skill Proficiency":   FeatureKeep, // Already have

	// Special features
	"Ability Score Improvement": FeatureAllowMultiple, // Each class grants separately
	"Feat":                FeatureAllowMultiple, // Can have multiple feats
}

// GetFeatureRule returns the rule for handling a duplicate feature
func GetFeatureRule(featureName string) DuplicateFeatureRule {
	// Normalize feature name (case-insensitive, trim spaces)
	normalized := strings.TrimSpace(featureName)

	// Check exact match first
	if rule, exists := FeatureRules[normalized]; exists {
		return rule
	}

	// Check partial matches for common patterns
	lowerName := strings.ToLower(normalized)
	if strings.Contains(lowerName, "extra attack") {
		return FeatureUseHighest
	}
	if strings.Contains(lowerName, "proficiency") {
		return FeatureKeep
	}
	if strings.Contains(lowerName, "fighting style") {
		return FeatureAllowMultiple
	}

	// Default: allow multiple (most features are class-specific)
	return FeatureAllowMultiple
}

// CanAddFeature checks if a feature can be added based on duplicate rules
func CanAddFeature(char *Character, featureName string, sourceClass string) (bool, string) {
	rule := GetFeatureRule(featureName)

	// Check if character already has this feature
	hasFeature := false
	existingSource := ""

	for _, feature := range char.Features.Features {
		if strings.EqualFold(feature.Name, featureName) {
			hasFeature = true
			existingSource = feature.Source
			break
		}
	}

	if !hasFeature {
		// Doesn't have it yet, always OK to add
		return true, ""
	}

	// Has the feature - check the rule
	switch rule {
	case FeatureStack:
		// Always OK to add (e.g., HP bonuses)
		return true, ""

	case FeatureReplace:
		// Replace the old one (rare, mostly for conflicting AC calculations)
		return true, fmt.Sprintf("Replacing %s from %s", featureName, existingSource)

	case FeatureKeep:
		// Already have it, don't add again
		return false, fmt.Sprintf("You already have %s from %s", featureName, existingSource)

	case FeatureUseHighest:
		// Don't add again, but note that we're keeping the better version
		return false, fmt.Sprintf("You already have %s (keeping the better version)", featureName)

	case FeatureAllowMultiple:
		// Check if it's from a different source
		if existingSource != sourceClass {
			// Different source, OK to add
			return true, fmt.Sprintf("Adding %s from %s (you also have it from %s)", featureName, sourceClass, existingSource)
		} else {
			// Same source, don't duplicate
			return false, fmt.Sprintf("You already have %s from %s", featureName, existingSource)
		}

	default:
		// Unknown rule, default to allow
		return true, ""
	}
}

// ValidateFightingStyle checks if a character can select a fighting style
func ValidateFightingStyle(char *Character, style string, sourceClass string) (bool, string) {
	// Check if character already has this fighting style from ANY class
	for _, cl := range char.Classes {
		if cl.FightingStyle == style {
			if cl.ClassName == sourceClass {
				return false, fmt.Sprintf("You already have %s fighting style from %s", style, sourceClass)
			} else {
				return false, fmt.Sprintf("You already have %s fighting style from %s (can't select same style from multiple classes)", style, cl.ClassName)
			}
		}
	}

	// Check if this class already has a fighting style
	for _, cl := range char.Classes {
		if cl.ClassName == sourceClass && cl.FightingStyle != "" {
			return false, fmt.Sprintf("%s already has %s fighting style", sourceClass, cl.FightingStyle)
		}
	}

	return true, ""
}

// RemoveConflictingFeature removes a feature that should be replaced
func RemoveConflictingFeature(char *Character, featureName string) {
	newFeatures := []Feature{}
	for _, feature := range char.Features.Features {
		if !strings.EqualFold(feature.Name, featureName) {
			newFeatures = append(newFeatures, feature)
		} else {
			debug.Log("Removed conflicting feature: %s", featureName)
		}
	}
	char.Features.Features = newFeatures
}

// ValidateExtraAttack handles the Extra Attack feature specially
// Multiple classes grant it, but it doesn't stack (except Fighter 11th level)
func ValidateExtraAttack(char *Character, attackCount int) int {
	// Check existing Extra Attack benefits
	maxAttacks := attackCount

	for _, benefit := range char.BenefitTracker.Benefits {
		if benefit.Type == BenefitFeature && strings.Contains(benefit.Target, "Extra Attack") {
			if benefit.Value > maxAttacks {
				maxAttacks = benefit.Value
			}
		}
	}

	debug.Log("Extra Attack validation: requested %d, max %d", attackCount, maxAttacks)
	return maxAttacks
}

// CheckProficiencyDuplicate checks if a proficiency is already granted
func CheckProficiencyDuplicate(char *Character, profType string, profName string) bool {
	switch profType {
	case "armor":
		for _, prof := range char.ArmorProficiencies {
			if strings.EqualFold(prof, profName) {
				return true
			}
		}
	case "weapon":
		for _, prof := range char.WeaponProficiencies {
			if strings.EqualFold(prof, profName) {
				return true
			}
		}
	case "tool":
		for _, prof := range char.ToolProficiencies {
			if strings.EqualFold(prof, profName) {
				return true
			}
		}
	case "skill":
		for _, skill := range char.Skills.List {
			if strings.EqualFold(string(skill.Name), profName) && skill.Proficiency >= Proficient {
				return true
			}
		}
	}
	return false
}
