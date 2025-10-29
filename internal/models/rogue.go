// internal/models/rogue.go
package models

import (
	"fmt"
)

// RogueMechanics handles Rogue-specific calculations and features
type RogueMechanics struct {
	character *Character
}

// NewRogueMechanics creates a new RogueMechanics instance
func NewRogueMechanics(char *Character) *RogueMechanics {
	return &RogueMechanics{character: char}
}

// GetSneakAttackDice returns the number of dice for Sneak Attack based on rogue level
// Scaling: 1d6 (1-2), 2d6 (3-4), 3d6 (5-6), 4d6 (7-8), 5d6 (9-10),
// 6d6 (11-12), 7d6 (13-14), 8d6 (15-16), 9d6 (17-18), 10d6 (19-20)
func (r *RogueMechanics) GetSneakAttackDice() int {
	rogueLevel := r.getRogueLevel()
	switch {
	case rogueLevel >= 19:
		return 10 // Levels 19-20
	case rogueLevel >= 17:
		return 9  // Levels 17-18
	case rogueLevel >= 15:
		return 8  // Levels 15-16
	case rogueLevel >= 13:
		return 7  // Levels 13-14
	case rogueLevel >= 11:
		return 6  // Levels 11-12
	case rogueLevel >= 9:
		return 5  // Levels 9-10
	case rogueLevel >= 7:
		return 4  // Levels 7-8
	case rogueLevel >= 5:
		return 3  // Levels 5-6
	case rogueLevel >= 3:
		return 2  // Levels 3-4
	case rogueLevel >= 1:
		return 1  // Levels 1-2
	default:
		return 0  // Not a rogue
	}
}

// GetSneakAttackDamage returns the damage dice string for Sneak Attack
func (r *RogueMechanics) GetSneakAttackDamage() string {
	dice := r.GetSneakAttackDice()
	if dice == 0 {
		return "0d6"
	}
	return fmt.Sprintf("%dd6", dice)
}

// getRogueLevel returns the character's Rogue class level
func (r *RogueMechanics) getRogueLevel() int {
	for _, classLevel := range r.character.Classes {
		if classLevel.ClassName == "Rogue" {
			return classLevel.Level
		}
	}
	return 0
}

// HasSneakAttackCondition checks if the character meets conditions for Sneak Attack
func (r *RogueMechanics) HasSneakAttackCondition() bool {
	// This would check for advantage, finesse/ranged weapon, etc.
	// For now, just return true if they have the feature
	return r.character.HasFeature("Sneak Attack")
}

// GetCunningActionOptions returns available Cunning Action options
func (r *RogueMechanics) GetCunningActionOptions() []string {
	if !r.character.HasFeature("Cunning Action") {
		return []string{}
	}

	options := []string{"Dash", "Disengage", "Hide"}

	// Add Steady Aim if available (level 3+)
	if r.character.HasFeature("Steady Aim") {
		options = append(options, "Steady Aim")
	}

	return options
}
