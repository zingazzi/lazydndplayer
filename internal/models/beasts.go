// internal/models/beasts.go
package models

// BeastDefinition represents a regular beast companion definition from the Monster Manual
type BeastDefinition struct {
	Name          string                 // e.g., "Mastiff", "Mule", "Black Bear"
	AC            int                    // Armor Class
	HP            int                    // Hit Points (fixed, doesn't scale)
	Speed         string                 // e.g., "40 ft."
	Senses        string                 // e.g., "Passive Perception 12"
	AbilityScores CompanionAbilityScores // Ability scores
	Traits        []string               // Special traits
	Actions       []BeastAction          // Available actions (simple attacks)
}

// BeastAction represents a simple attack action for a regular beast
type BeastAction struct {
	Name        string // e.g., "Bite"
	AttackBonus int    // Fixed attack bonus
	DamageDice  string // e.g., "1d6"
	DamageBonus int    // Fixed damage bonus
	DamageType  string // e.g., "piercing"
	Range       string // e.g., "5 ft."
	Description string // Optional description
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

// GetAllBeastDefinitions returns all available regular beast definitions
func GetAllBeastDefinitions() []BeastDefinition {
	return []BeastDefinition{
		{
			Name: "Mastiff",
			AC:   12,
			HP:   5,
			Speed: "40 ft.",
			Senses: "Passive Perception 12",
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
			Name: "Mule",
			AC:   10,
			HP:   11,
			Speed: "40 ft.",
			Senses: "Passive Perception 10",
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
			Name: "Black Bear",
			AC:   11,
			HP:   19,
			Speed: "40 ft., climb 30 ft.",
			Senses: "Passive Perception 13",
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
			Name: "Wolf",
			AC:   13,
			HP:   11,
			Speed: "40 ft.",
			Senses: "Passive Perception 13",
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
			Name: "Hawk",
			AC:   13,
			HP:   1,
			Speed: "10 ft., fly 60 ft.",
			Senses: "Passive Perception 14",
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
			Name: "Cat",
			AC:   12,
			HP:   2,
			Speed: "40 ft., climb 30 ft.",
			Senses: "Passive Perception 13",
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
}
