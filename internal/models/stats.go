// internal/models/stats.go
package models

import "math"

// AbilityType represents one of the six core abilities
type AbilityType string

const (
	Strength     AbilityType = "STR"
	Dexterity    AbilityType = "DEX"
	Constitution AbilityType = "CON"
	Intelligence AbilityType = "INT"
	Wisdom       AbilityType = "WIS"
	Charisma     AbilityType = "CHA"
)

// AbilityScores holds all six ability scores
type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// Base values (before racial/feat bonuses)
	StrengthBase     int `json:"strength_base,omitempty"`
	DexterityBase    int `json:"dexterity_base,omitempty"`
	ConstitutionBase int `json:"constitution_base,omitempty"`
	IntelligenceBase int `json:"intelligence_base,omitempty"`
	WisdomBase       int `json:"wisdom_base,omitempty"`
	CharismaBase     int `json:"charisma_base,omitempty"`

	// Extra bonuses (from species, feats, etc.)
	StrengthExtra     int `json:"strength_extra,omitempty"`
	DexterityExtra    int `json:"dexterity_extra,omitempty"`
	ConstitutionExtra int `json:"constitution_extra,omitempty"`
	IntelligenceExtra int `json:"intelligence_extra,omitempty"`
	WisdomExtra       int `json:"wisdom_extra,omitempty"`
	CharismaExtra     int `json:"charisma_extra,omitempty"`
}

// GetScore returns the score for a given ability
func (a *AbilityScores) GetScore(ability AbilityType) int {
	switch ability {
	case Strength:
		return a.Strength
	case Dexterity:
		return a.Dexterity
	case Constitution:
		return a.Constitution
	case Intelligence:
		return a.Intelligence
	case Wisdom:
		return a.Wisdom
	case Charisma:
		return a.Charisma
	default:
		return 10
	}
}

// SetScore sets the score for a given ability
func (a *AbilityScores) SetScore(ability AbilityType, score int) {
	switch ability {
	case Strength:
		a.Strength = score
	case Dexterity:
		a.Dexterity = score
	case Constitution:
		a.Constitution = score
	case Intelligence:
		a.Intelligence = score
	case Wisdom:
		a.Wisdom = score
	case Charisma:
		a.Charisma = score
	}
}

// CalculateModifier calculates the ability modifier for a given score
func CalculateModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2.0))
}

// GetModifier returns the modifier for a given ability
func (a *AbilityScores) GetModifier(ability interface{}) int {
	var abilityType AbilityType

	// Handle both string and AbilityType
	switch v := ability.(type) {
	case string:
		// Convert string to AbilityType
		switch v {
		case "Strength", "STR":
			abilityType = Strength
		case "Dexterity", "DEX":
			abilityType = Dexterity
		case "Constitution", "CON":
			abilityType = Constitution
		case "Intelligence", "INT":
			abilityType = Intelligence
		case "Wisdom", "WIS":
			abilityType = Wisdom
		case "Charisma", "CHA":
			abilityType = Charisma
		default:
			return 0
		}
	case AbilityType:
		abilityType = v
	default:
		return 0
	}

	return CalculateModifier(a.GetScore(abilityType))
}

// SavingThrows tracks proficiency in saving throws
type SavingThrows struct {
	Strength     bool `json:"strength"`
	Dexterity    bool `json:"dexterity"`
	Constitution bool `json:"constitution"`
	Intelligence bool `json:"intelligence"`
	Wisdom       bool `json:"wisdom"`
	Charisma     bool `json:"charisma"`
}

// IsProficient checks if character is proficient in a saving throw
func (s *SavingThrows) IsProficient(ability AbilityType) bool {
	switch ability {
	case Strength:
		return s.Strength
	case Dexterity:
		return s.Dexterity
	case Constitution:
		return s.Constitution
	case Intelligence:
		return s.Intelligence
	case Wisdom:
		return s.Wisdom
	case Charisma:
		return s.Charisma
	default:
		return false
	}
}

// GetBaseScore returns the base score for a given ability
func (a *AbilityScores) GetBaseScore(ability AbilityType) int {
	switch ability {
	case Strength:
		return a.StrengthBase
	case Dexterity:
		return a.DexterityBase
	case Constitution:
		return a.ConstitutionBase
	case Intelligence:
		return a.IntelligenceBase
	case Wisdom:
		return a.WisdomBase
	case Charisma:
		return a.CharismaBase
	default:
		return 10
	}
}

// SetBaseScore sets the base score for a given ability and updates total
func (a *AbilityScores) SetBaseScore(ability AbilityType, score int) {
	switch ability {
	case Strength:
		a.StrengthBase = score
		a.Strength = score + a.StrengthExtra
	case Dexterity:
		a.DexterityBase = score
		a.Dexterity = score + a.DexterityExtra
	case Constitution:
		a.ConstitutionBase = score
		a.Constitution = score + a.ConstitutionExtra
	case Intelligence:
		a.IntelligenceBase = score
		a.Intelligence = score + a.IntelligenceExtra
	case Wisdom:
		a.WisdomBase = score
		a.Wisdom = score + a.WisdomExtra
	case Charisma:
		a.CharismaBase = score
		a.Charisma = score + a.CharismaExtra
	}
}

// GetExtraScore returns the extra bonus for a given ability
func (a *AbilityScores) GetExtraScore(ability AbilityType) int {
	switch ability {
	case Strength:
		return a.StrengthExtra
	case Dexterity:
		return a.DexterityExtra
	case Constitution:
		return a.ConstitutionExtra
	case Intelligence:
		return a.IntelligenceExtra
	case Wisdom:
		return a.WisdomExtra
	case Charisma:
		return a.CharismaExtra
	default:
		return 0
	}
}

// SetExtraScore sets the extra bonus for a given ability and updates total
func (a *AbilityScores) SetExtraScore(ability AbilityType, extra int) {
	switch ability {
	case Strength:
		a.StrengthExtra = extra
		a.Strength = a.StrengthBase + extra
	case Dexterity:
		a.DexterityExtra = extra
		a.Dexterity = a.DexterityBase + extra
	case Constitution:
		a.ConstitutionExtra = extra
		a.Constitution = a.ConstitutionBase + extra
	case Intelligence:
		a.IntelligenceExtra = extra
		a.Intelligence = a.IntelligenceBase + extra
	case Wisdom:
		a.WisdomExtra = extra
		a.Wisdom = a.WisdomBase + extra
	case Charisma:
		a.CharismaExtra = extra
		a.Charisma = a.CharismaBase + extra
	}
}

// RecalculateTotals recalculates all total scores from base and extra
func (a *AbilityScores) RecalculateTotals() {
	a.Strength = a.StrengthBase + a.StrengthExtra
	a.Dexterity = a.DexterityBase + a.DexterityExtra
	a.Constitution = a.ConstitutionBase + a.ConstitutionExtra
	a.Intelligence = a.IntelligenceBase + a.IntelligenceExtra
	a.Wisdom = a.WisdomBase + a.WisdomExtra
	a.Charisma = a.CharismaBase + a.CharismaExtra
}
