// internal/models/speciesspells.go
package models

import (
	"regexp"
	"strings"
)

// GetSpeciesSpells returns a map of spell definitions for species-granted spells
func GetSpeciesSpells() map[string]Spell {
	return map[string]Spell{
		// Cantrips
		"Light": {
			Name:        "Light",
			Level:       0,
			School:      Evocation,
			CastingTime: "1 action",
			Range:       "Touch",
			Components:  "V, M",
			Duration:    "1 hour",
			Description: "You touch one object that is no larger than 10 feet in any dimension. Until the spell ends, the object sheds bright light in a 20-foot radius and dim light for an additional 20 feet.",
			Known:       true,
			Prepared:    true,
		},
		"Dancing Lights": {
			Name:        "Dancing Lights",
			Level:       0,
			School:      Evocation,
			CastingTime: "1 action",
			Range:       "120 feet",
			Components:  "V, S, M",
			Duration:    "Concentration, up to 1 minute",
			Description: "You create up to four torch-sized lights within range, making them appear as torches, lanterns, or glowing orbs that hover in the air for the duration.",
			Known:       true,
			Prepared:    true,
		},
		"Chill Touch": {
			Name:        "Chill Touch",
			Level:       0,
			School:      Necromancy,
			CastingTime: "1 action",
			Range:       "120 feet",
			Components:  "V, S",
			Duration:    "1 round",
			Description: "You create a ghostly, skeletal hand in the space of a creature within range. Make a ranged spell attack. On a hit, the target takes 1d8 necrotic damage.",
			Known:       true,
			Prepared:    true,
		},
		"Thaumaturgy": {
			Name:        "Thaumaturgy",
			Level:       0,
			School:      Transmutation,
			CastingTime: "1 action",
			Range:       "30 feet",
			Components:  "V",
			Duration:    "Up to 1 minute",
			Description: "You manifest a minor wonder, a sign of supernatural power, within range. You create one of the following effects: your voice booms up to three times as loud, you cause flames to flicker, you cause harmless tremors, or you create an instantaneous sound.",
			Known:       true,
			Prepared:    true,
		},
		// 1st Level Spells
		"Faerie Fire": {
			Name:        "Faerie Fire",
			Level:       1,
			School:      Evocation,
			CastingTime: "1 action",
			Range:       "60 feet",
			Components:  "V",
			Duration:    "Concentration, up to 1 minute",
			Description: "Each object in a 20-foot cube within range is outlined in blue, green, or violet light. Any creature in the area when the spell is cast is also outlined if it fails a Dexterity saving throw. For the duration, objects and affected creatures shed dim light in a 10-foot radius. Attack rolls against affected creatures have advantage.",
			Known:       true,
			Prepared:    false,
		},
		"Burning Hands": {
			Name:        "Burning Hands",
			Level:       1,
			School:      Evocation,
			CastingTime: "1 action",
			Range:       "Self (15-foot cone)",
			Components:  "V, S",
			Duration:    "Instantaneous",
			Description: "As you hold your hands with thumbs touching and fingers spread, a thin sheet of flames shoots forth from your outstretched fingertips. Each creature in a 15-foot cone must make a Dexterity saving throw. A creature takes 3d6 fire damage on a failed save, or half as much damage on a successful one.",
			Known:       true,
			Prepared:    false,
		},
		"False Life": {
			Name:        "False Life",
			Level:       1,
			School:      Necromancy,
			CastingTime: "1 action",
			Range:       "Self",
			Components:  "V, S, M",
			Duration:    "1 hour",
			Description: "Bolstering yourself with a necromantic facsimile of life, you gain 1d4 + 4 temporary hit points for the duration.",
			Known:       true,
			Prepared:    false,
		},
		"Hellish Rebuke": {
			Name:        "Hellish Rebuke",
			Level:       1,
			School:      Evocation,
			CastingTime: "1 reaction",
			Range:       "60 feet",
			Components:  "V, S",
			Duration:    "Instantaneous",
			Description: "You point your finger, and the creature that damaged you is momentarily surrounded by hellish flames. The creature must make a Dexterity saving throw. It takes 2d10 fire damage on a failed save, or half as much damage on a successful one.",
			Known:       true,
			Prepared:    false,
		},
		// 2nd Level Spells
		"Darkness": {
			Name:        "Darkness",
			Level:       2,
			School:      Evocation,
			CastingTime: "1 action",
			Range:       "60 feet",
			Components:  "V, M",
			Duration:    "Concentration, up to 10 minutes",
			Description: "Magical darkness spreads from a point you choose within range to fill a 15-foot radius sphere for the duration. The darkness spreads around corners. Darkvision can't penetrate this darkness.",
			Known:       true,
			Prepared:    false,
		},
		"Alter Self": {
			Name:        "Alter Self",
			Level:       2,
			School:      Transmutation,
			CastingTime: "1 action",
			Range:       "Self",
			Components:  "V, S",
			Duration:    "Concentration, up to 1 hour",
			Description: "You assume a different form. When you cast the spell, choose one of the following options: Aquatic Adaptation, Change Appearance, or Natural Weapons.",
			Known:       true,
			Prepared:    false,
		},
		"Ray of Enfeeblement": {
			Name:        "Ray of Enfeeblement",
			Level:       2,
			School:      Necromancy,
			CastingTime: "1 action",
			Range:       "60 feet",
			Components:  "V, S",
			Duration:    "Concentration, up to 1 minute",
			Description: "A black beam of enervating energy springs from your finger toward a creature within range. Make a ranged spell attack. On a hit, the target deals only half damage with weapon attacks that use Strength until the spell ends.",
			Known:       true,
			Prepared:    false,
		},
	}
}

// ParseSpeciesSpells extracts spell information from species traits
func ParseSpeciesSpells(trait SpeciesTrait) []SpellGrant {
	grants := []SpellGrant{}
	desc := strings.ToLower(trait.Description)
	name := strings.ToLower(trait.Name)

	// Check for spell grants in the description
	if strings.Contains(desc, "cantrip") || strings.Contains(name, "legacy") || strings.Contains(name, "bearer") || strings.Contains(name, "magic") {
		// Extract spell names and their level requirements
		spellNames := []string{
			"Light", "Dancing Lights", "Chill Touch", "Thaumaturgy",
			"Faerie Fire", "Burning Hands", "False Life", "Hellish Rebuke",
			"Darkness", "Alter Self", "Ray of Enfeeblement",
		}

		for _, spellName := range spellNames {
			spellLower := strings.ToLower(spellName)
			if strings.Contains(desc, spellLower) || strings.Contains(name, spellLower) {
				// Determine level requirement
				levelReq := 1 // Default to level 1 for cantrips

				// Check for level requirements
				if strings.Contains(desc, "3rd level") && strings.Index(desc, "3rd level") < strings.Index(desc, spellLower) {
					levelReq = 3
				} else if strings.Contains(desc, "5th level") && strings.Index(desc, "5th level") < strings.Index(desc, spellLower) {
					levelReq = 5
				}

				grants = append(grants, SpellGrant{
					SpellName:    spellName,
					LevelRequired: levelReq,
				})
			}
		}
	}

	return grants
}

// SpellGrant represents a spell granted by a species trait
type SpellGrant struct {
	SpellName     string
	LevelRequired int
}

// GetSpellsForLevel returns all species spells available at a given level
func GetSpellsForLevel(trait SpeciesTrait, characterLevel int) []string {
	grants := ParseSpeciesSpells(trait)
	available := []string{}

	for _, grant := range grants {
		if characterLevel >= grant.LevelRequired {
			available = append(available, grant.SpellName)
		}
	}

	return available
}

// ExtractSpellsFromDescription extracts spells using regex pattern matching
func ExtractSpellsFromDescription(description string) []SpellGrant {
	grants := []SpellGrant{}

	// Pattern: "spell name" followed by optional "at Xth level"
	patterns := map[string]int{
		`(?i)you know the (\w+(?:\s+\w+)*) cantrip`:                         1,
		`(?i)at 3rd level.*?you can cast (\w+(?:\s+\w+)*)`:                  3,
		`(?i)at 5th level.*?you can cast (\w+(?:\s+\w+)*)`:                  5,
		`(?i)3rd level.*?cast (\w+(?:\s+\w+)*)`:                             3,
		`(?i)5th level.*?cast (\w+(?:\s+\w+)*)`:                             5,
	}

	for pattern, level := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(description, -1)
		for _, match := range matches {
			if len(match) > 1 {
				spellName := capitalizeSpellName(match[1])
				grants = append(grants, SpellGrant{
					SpellName:     spellName,
					LevelRequired: level,
				})
			}
		}
	}

	return grants
}

// capitalizeSpellName properly capitalizes spell names
func capitalizeSpellName(name string) string {
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
