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
func (a *AbilityScores) GetModifier(ability AbilityType) int {
	return CalculateModifier(a.GetScore(ability))
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
