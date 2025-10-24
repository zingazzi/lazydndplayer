// internal/models/asi.go
package models

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// CheckASIAvailable returns true if this level grants an Ability Score Improvement
func CheckASIAvailable(className string, level int) bool {
	// Standard ASI levels for all classes (4, 8, 12, 16)
	// Fighters get additional ASI at levels 6 and 14, Rogues at 10
	standardASILevels := level == 4 || level == 8 || level == 12 || level == 16

	// Class-specific additional ASI levels
	switch className {
	case "Fighter":
		return standardASILevels || level == 6 || level == 14
	case "Rogue":
		return standardASILevels || level == 10
	default:
		return standardASILevels
	}
}

// GetAvailableFeatsForASI returns feats the character can take during ASI
// Checks: ability prerequisites, repeatable limits, conflicting feats
func GetAvailableFeatsForASI(char *Character) []Feat {
	allFeats := GetAllFeats()
	availableFeats := []Feat{}

	for _, feat := range allFeats {
		// Check if character can take this feat (prerequisites)
		if !CanTakeFeat(char, feat) {
			continue
		}

		// Skip if character already has this feat and it's not repeatable
		if HasFeat(char, feat.Name) && !feat.Repeatable {
			continue
		}

		// Check for conflicting feats
		if HasConflictingFeat(char, feat.Name) {
			continue
		}

		availableFeats = append(availableFeats, feat)
	}

	debug.Log("GetAvailableFeatsForASI: Found %d available feats", len(availableFeats))
	return availableFeats
}

// ApplyASIChoice applies the selected ASI choice to the character
func ApplyASIChoice(char *Character, className string, level int, choice ASIChoice) error {
	debug.Log("ApplyASIChoice: className=%s, level=%d, type=%s", className, level, choice.Type)

	// Get the class level struct
	classLevel := char.GetClassLevelStruct(className)
	if classLevel == nil {
		return fmt.Errorf("character does not have %s class", className)
	}

	// Initialize ASIChoices map if needed
	if classLevel.ASIChoices == nil {
		classLevel.ASIChoices = make(map[int]ASIChoice)
	}

	// Apply the choice
	if choice.Type == "ability" {
		// Apply ability score boosts
		source := BenefitSource{
			Type: "ASI",
			Name: fmt.Sprintf("%s Level %d", className, level),
		}
		applier := NewBenefitApplier(char)

		for _, boost := range choice.AbilityBoosts {
			debug.Log("  Applying +%d to %s", boost.Amount, boost.Ability)
			if err := applier.AddAbilityScore(source, boost.Ability, boost.Amount); err != nil {
				return fmt.Errorf("failed to apply ability boost: %w", err)
			}
		}
	} else if choice.Type == "feat" {
		// Apply feat
		debug.Log("  Applying feat: %s", choice.FeatName)
		feat := GetFeatByName(choice.FeatName)
		if feat == nil {
			return fmt.Errorf("feat not found: %s", choice.FeatName)
		}

		// Apply feat benefits
		source := BenefitSource{
			Type: "ASI Feat",
			Name: fmt.Sprintf("%s Level %d", className, level),
		}
		if err := ApplyFeat(char, *feat, source); err != nil {
			return fmt.Errorf("failed to apply feat: %w", err)
		}

		// Add feat to character's feat list
		char.Feats = append(char.Feats, choice.FeatName)
	} else {
		return fmt.Errorf("invalid ASI choice type: %s", choice.Type)
	}

	// Store the choice
	classLevel.ASIChoices[level] = choice
	debug.Log("ApplyASIChoice: Successfully applied and stored choice")

	// Update derived stats
	char.UpdateDerivedStats()

	return nil
}

// RemoveASIChoice removes ASI effects when de-leveling
func RemoveASIChoice(char *Character, className string, level int) error {
	debug.Log("RemoveASIChoice: className=%s, level=%d", className, level)

	// Get the class level struct
	classLevel := char.GetClassLevelStruct(className)
	if classLevel == nil {
		return fmt.Errorf("character does not have %s class", className)
	}

	// Get the choice
	choice, exists := classLevel.ASIChoices[level]
	if !exists {
		debug.Log("RemoveASIChoice: No ASI choice found for level %d", level)
		return nil // No choice to remove
	}

	// Remove the effects
	if choice.Type == "ability" {
		// Remove ability score boosts
		source := BenefitSource{
			Type: "ASI",
			Name: fmt.Sprintf("%s Level %d", className, level),
		}
		remover := NewBenefitRemover(char)

		debug.Log("  Removing ability boosts from ASI")
		if err := remover.RemoveAllBenefits(source.Type, source.Name); err != nil {
			return fmt.Errorf("failed to remove ability boosts: %w", err)
		}
	} else if choice.Type == "feat" {
		// Remove feat
		debug.Log("  Removing feat: %s", choice.FeatName)

		// Remove feat benefits
		source := BenefitSource{
			Type: "ASI Feat",
			Name: fmt.Sprintf("%s Level %d", className, level),
		}
		remover := NewBenefitRemover(char)
		if err := remover.RemoveAllBenefits(source.Type, source.Name); err != nil {
			return fmt.Errorf("failed to remove feat benefits: %w", err)
		}

		// Remove feat from character's feat list
		for i, featName := range char.Feats {
			if featName == choice.FeatName {
				char.Feats = append(char.Feats[:i], char.Feats[i+1:]...)
				debug.Log("  Removed feat from Feats array")
				break
			}
		}
	}

	// Delete the choice from the map
	delete(classLevel.ASIChoices, level)
	debug.Log("RemoveASIChoice: Successfully removed ASI choice")

	// Update derived stats
	char.UpdateDerivedStats()

	return nil
}

// ValidateAbilityBoosts checks if ability boosts are valid (total +2, no more than +2 per ability, respects 20 cap)
func ValidateAbilityBoosts(char *Character, boosts []AbilityBoost) error {
	if len(boosts) == 0 {
		return fmt.Errorf("no ability boosts specified")
	}

	// Calculate total boost
	totalBoost := 0
	abilityTotals := make(map[string]int)

	for _, boost := range boosts {
		if boost.Amount < 1 {
			return fmt.Errorf("ability boost amount must be positive")
		}
		totalBoost += boost.Amount
		abilityTotals[boost.Ability] += boost.Amount
	}

	// Check total is +2
	if totalBoost != 2 {
		return fmt.Errorf("total ability boost must be +2, got %d", totalBoost)
	}

	// Check no ability gets more than +2
	for ability, total := range abilityTotals {
		if total > 2 {
			return fmt.Errorf("%s cannot receive more than +2", ability)
		}

		// Check ability score won't exceed 20 (standard cap)
		currentScore := char.GetAbilityScore(ability)
		if currentScore+total > 20 {
			return fmt.Errorf("%s would exceed maximum of 20 (current: %d)", ability, currentScore)
		}
	}

	return nil
}
