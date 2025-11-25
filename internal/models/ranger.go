// internal/models/ranger.go
package models

// RangerMechanics handles Ranger-specific calculations and features
type RangerMechanics struct {
	character *Character
}

// NewRangerMechanics creates a new RangerMechanics instance
func NewRangerMechanics(char *Character) *RangerMechanics {
	return &RangerMechanics{character: char}
}

// GetDreadfulStrikesDamage returns the damage dice string for Fey Wanderer Dreadful Strikes
// Scales from 1d4 at level 3 to 1d6 at level 11
func (r *RangerMechanics) GetDreadfulStrikesDamage() string {
	rangerLevel := r.getRangerLevel()
	if rangerLevel >= 11 {
		return "1d6"
	} else if rangerLevel >= 3 {
		return "1d4"
	}
	return "0d4"
}

// GetDreadfulStrikesDamageType returns the damage type for Dreadful Strikes
func (r *RangerMechanics) GetDreadfulStrikesDamageType() string {
	// Check if character has Fey Wanderer Dreadful Strikes
	if r.character.HasFeature("Dreadful Strikes") {
		// Check subclass
		for _, classLevel := range r.character.Classes {
			if classLevel.ClassName == "Ranger" && classLevel.Subclass == "Fey Wanderer" {
				return "psychic"
			}
		}
	}
	return ""
}

// HasDreadfulStrikes checks if the character has Dreadful Strikes feature
func (r *RangerMechanics) HasDreadfulStrikes() bool {
	return r.character.HasFeature("Dreadful Strikes")
}

// getRangerLevel returns the character's Ranger class level
func (r *RangerMechanics) getRangerLevel() int {
	for _, classLevel := range r.character.Classes {
		if classLevel.ClassName == "Ranger" {
			return classLevel.Level
		}
	}
	return 0
}

// GetRangerSubclass returns the Ranger subclass name, or empty string if none
func (r *RangerMechanics) GetRangerSubclass() string {
	for _, classLevel := range r.character.Classes {
		if classLevel.ClassName == "Ranger" {
			return classLevel.Subclass
		}
	}
	return ""
}

// IsFeyWanderer checks if character is a Fey Wanderer Ranger
func (r *RangerMechanics) IsFeyWanderer() bool {
	return r.GetRangerSubclass() == "Fey Wanderer"
}

// IsGloomStalker checks if character is a Gloom Stalker Ranger
func (r *RangerMechanics) IsGloomStalker() bool {
	return r.GetRangerSubclass() == "Gloom Stalker"
}

// IsHunter checks if character is a Hunter Ranger
func (r *RangerMechanics) IsHunter() bool {
	return r.GetRangerSubclass() == "Hunter"
}
