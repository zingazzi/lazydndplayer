// internal/models/beasts.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// BeastDefinition represents a regular beast companion definition from the Monster Manual
type BeastDefinition struct {
	Name          string                 `json:"name"`           // e.g., "Mastiff", "Mule", "Black Bear"
	CR            float64                `json:"cr,omitempty"`  // Challenge Rating (for wild shape filtering)
	AC            int                    `json:"ac"`            // Armor Class
	HP            int                    `json:"hp"`            // Hit Points (fixed, doesn't scale)
	Speed         string                 `json:"speed"`         // e.g., "40 ft."
	Senses        string                 `json:"senses"`        // e.g., "Passive Perception 12"
	AbilityScores CompanionAbilityScores `json:"ability_scores"` // Ability scores
	Traits        []string               `json:"traits"`        // Special traits
	Actions       []BeastAction          `json:"actions"`       // Available actions (simple attacks)
	CanFly        bool                   `json:"can_fly,omitempty"` // Whether the beast has a flying speed
	CanSwim       bool                   `json:"can_swim,omitempty"` // Whether the beast has a swimming speed
}

// BeastAction represents a simple attack action for a regular beast
type BeastAction struct {
	Name        string `json:"name"`         // e.g., "Bite"
	AttackBonus int    `json:"attack_bonus"` // Fixed attack bonus
	DamageDice  string `json:"damage_dice"`  // e.g., "1d6"
	DamageBonus int    `json:"damage_bonus"` // Fixed damage bonus
	DamageType  string `json:"damage_type"`  // e.g., "piercing"
	Range       string `json:"range"`        // e.g., "5 ft."
	Description string `json:"description"`  // Optional description
}

// BeastsData represents the root structure of the beasts JSON file
type BeastsData struct {
	Beasts []BeastDefinition `json:"beasts"`
}

var cachedBeasts []BeastDefinition

// LoadBeastsFromJSON loads beast data from the JSON file
func LoadBeastsFromJSON(filepath string) ([]BeastDefinition, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read beasts file: %w", err)
	}

	var data BeastsData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to parse beasts JSON: %w", err)
	}

	return data.Beasts, nil
}

// GetBeastDefinition returns a beast definition by name
func GetBeastDefinition(name string) *BeastDefinition {
	beasts := GetAllBeastDefinitions()
	for _, beast := range beasts {
		if beast.Name == name {
			return &beast
		}
	}
	return nil
}

// GetBeastsByCR returns all beasts with a CR less than or equal to maxCR
func GetBeastsByCR(maxCR float64, allowFlying bool, allowSwimming bool) []BeastDefinition {
	beasts := GetAllBeastDefinitions()
	var filtered []BeastDefinition

	for _, beast := range beasts {
		// Check CR limit
		if beast.CR > maxCR {
			continue
		}

		// Check flying restriction
		if !allowFlying && beast.CanFly {
			continue
		}

		// Check swimming restriction
		if !allowSwimming && beast.CanSwim {
			continue
		}

		filtered = append(filtered, beast)
	}

	return filtered
}

// GetWildShapeForms returns beasts suitable for wild shape based on druid level
func GetWildShapeForms(druidLevel int) []BeastDefinition {
	var maxCR float64
	allowFlying := false
	allowSwimming := false

	if druidLevel >= 8 {
		maxCR = 1.0
		allowFlying = true
		allowSwimming = true
	} else if druidLevel >= 4 {
		maxCR = 0.5
		allowSwimming = true
	} else if druidLevel >= 2 {
		maxCR = 0.25
	} else {
		return []BeastDefinition{} // No wild shape at level 1
	}

	return GetBeastsByCR(maxCR, allowFlying, allowSwimming)
}

// GetAllBeastDefinitions returns all available regular beast definitions
// Loads from JSON file with fallback to hardcoded data
func GetAllBeastDefinitions() []BeastDefinition {
	// Return cached beasts if already loaded
	if len(cachedBeasts) > 0 {
		return cachedBeasts
	}

	// Try to load from JSON file
	beasts, err := LoadBeastsFromJSON("data/beasts.json")
	if err == nil {
		cachedBeasts = beasts
		return cachedBeasts
	}

	// Fallback to hardcoded beasts if JSON file is not available
	cachedBeasts = []BeastDefinition{
		{
			Name:   "Mastiff",
			CR:     0.125,
			AC:     12,
			HP:     5,
			Speed:  "40 ft.",
			Senses: "Passive Perception 12",
			CanFly: false,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     13,
				Dexterity:    14,
				Constitution: 12,
				Intelligence: 3,
				Wisdom:       12,
				Charisma:     7,
			},
			Traits: []string{
				"Keen Hearing and Smell: The mastiff has advantage on Wisdom (Perception) checks that rely on hearing or smell.",
			},
			Actions: []BeastAction{
				{
					Name:        "Bite",
					AttackBonus: 3,
					DamageDice:  "1d6",
					DamageBonus: 1,
					DamageType:  "piercing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
			},
		},
		{
			Name:   "Mule",
			CR:     0.125,
			AC:     10,
			HP:     11,
			Speed:  "40 ft.",
			Senses: "Passive Perception 10",
			CanFly: false,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     14,
				Dexterity:    10,
				Constitution: 13,
				Intelligence: 2,
				Wisdom:       10,
				Charisma:     5,
			},
			Traits: []string{
				"Beast of Burden: The mule is considered to be a Large animal for the purpose of determining its carrying capacity.",
				"Sure-Footed: The mule has advantage on Strength and Dexterity saving throws made against effects that would knock it prone.",
			},
			Actions: []BeastAction{
				{
					Name:        "Hooves",
					AttackBonus: 2,
					DamageDice:  "1d4",
					DamageBonus: 2,
					DamageType:  "bludgeoning",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
			},
		},
		{
			Name:   "Black Bear",
			CR:     0.5,
			AC:     11,
			HP:     19,
			Speed:  "40 ft., climb 30 ft.",
			Senses: "Passive Perception 13",
			CanFly: false,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     15,
				Dexterity:    10,
				Constitution: 14,
				Intelligence: 2,
				Wisdom:       12,
				Charisma:     7,
			},
			Traits: []string{
				"Keen Smell: The bear has advantage on Wisdom (Perception) checks that rely on smell.",
			},
			Actions: []BeastAction{
				{
					Name:        "Multiattack",
					AttackBonus: 0,
					DamageDice:  "",
					DamageBonus: 0,
					DamageType:  "",
					Range:       "",
					Description: "The bear makes two attacks: one with its bite and one with its claws.",
				},
				{
					Name:        "Bite",
					AttackBonus: 3,
					DamageDice:  "1d6",
					DamageBonus: 2,
					DamageType:  "piercing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
				{
					Name:        "Claws",
					AttackBonus: 3,
					DamageDice:  "2d4",
					DamageBonus: 2,
					DamageType:  "slashing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
			},
		},
		{
			Name:   "Wolf",
			CR:     0.25,
			AC:     13,
			HP:     11,
			Speed:  "40 ft.",
			Senses: "Passive Perception 13",
			CanFly: false,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     12,
				Dexterity:    15,
				Constitution: 12,
				Intelligence: 3,
				Wisdom:       12,
				Charisma:     6,
			},
			Traits: []string{
				"Keen Hearing and Smell: The wolf has advantage on Wisdom (Perception) checks that rely on hearing or smell.",
				"Pack Tactics: The wolf has advantage on an attack roll against a creature if at least one of the wolf's allies is within 5 feet of the creature and the ally isn't incapacitated.",
			},
			Actions: []BeastAction{
				{
					Name:        "Bite",
					AttackBonus: 4,
					DamageDice:  "2d4",
					DamageBonus: 2,
					DamageType:  "piercing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack. If the target is a creature, it must succeed on a DC 11 Strength saving throw or be knocked prone.",
				},
			},
		},
		{
			Name:   "Hawk",
			CR:     0,
			AC:     13,
			HP:     1,
			Speed:  "10 ft., fly 60 ft.",
			Senses: "Passive Perception 14",
			CanFly: true,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     5,
				Dexterity:    16,
				Constitution: 8,
				Intelligence: 2,
				Wisdom:       14,
				Charisma:     6,
			},
			Traits: []string{
				"Keen Sight: The hawk has advantage on Wisdom (Perception) checks that rely on sight.",
			},
			Actions: []BeastAction{
				{
					Name:        "Talons",
					AttackBonus: 5,
					DamageDice:  "1d1",
					DamageBonus: 0,
					DamageType:  "slashing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
			},
		},
		{
			Name:   "Cat",
			CR:     0,
			AC:     12,
			HP:     2,
			Speed:  "40 ft., climb 30 ft.",
			Senses: "Passive Perception 13",
			CanFly: false,
			CanSwim: false,
			AbilityScores: CompanionAbilityScores{
				Strength:     3,
				Dexterity:    15,
				Constitution: 10,
				Intelligence: 3,
				Wisdom:       12,
				Charisma:     7,
			},
			Traits: []string{
				"Keen Smell: The cat has advantage on Wisdom (Perception) checks that rely on smell.",
			},
			Actions: []BeastAction{
				{
					Name:        "Claws",
					AttackBonus: 0,
					DamageDice:  "1d1",
					DamageBonus: 0,
					DamageType:  "slashing",
					Range:       "5 ft.",
					Description: "Melee Weapon Attack",
				},
			},
		},
	}
	return cachedBeasts
}
