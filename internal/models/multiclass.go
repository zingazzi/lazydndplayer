// internal/models/multiclass.go
package models

import "fmt"

// ClassLevel represents a single class and its level for multiclassing
type ClassLevel struct {
	ClassName     string `json:"class_name"`
	Level         int    `json:"level"`
	Subclass      string `json:"subclass,omitempty"`       // Domain/Patron/Origin/etc
	FightingStyle string `json:"fighting_style,omitempty"` // For Fighter, Paladin, Ranger
}

// MulticlassPrerequisites defines ability score requirements for multiclassing
var MulticlassPrerequisites = map[string]map[string]int{
	"Barbarian": {"Strength": 13},
	"Bard":      {"Charisma": 13},
	"Cleric":    {"Wisdom": 13},
	"Druid":     {"Wisdom": 13},
	"Fighter":   {"Strength": 13, "Dexterity": 13}, // Either Str OR Dex
	"Monk":      {"Dexterity": 13, "Wisdom": 13},   // Both required
	"Paladin":   {"Strength": 13, "Charisma": 13},  // Both required
	"Ranger":    {"Dexterity": 13, "Wisdom": 13},   // Both required
	"Rogue":     {"Dexterity": 13},
	"Sorcerer":  {"Charisma": 13},
	"Warlock":   {"Charisma": 13},
	"Wizard":    {"Intelligence": 13},
}

// MulticlassProficiencies defines limited proficiencies granted when multiclassing
// (not the full starting proficiencies)
var MulticlassProficiencies = map[string][]string{
	"Barbarian": {"Light Armor", "Medium Armor", "Shields", "Simple", "Martial"},
	"Bard":      {"Light Armor", "Simple", "Hand Crossbows", "Longswords", "Rapiers", "Shortswords"},
	"Cleric":    {"Light Armor", "Medium Armor", "Shields"},
	"Druid":     {"Light Armor", "Medium Armor", "Shields"},
	"Fighter":   {"Light Armor", "Medium Armor", "Heavy Armor", "Shields", "Simple", "Martial"},
	"Monk":      {"Simple", "Shortswords"},
	"Paladin":   {"Light Armor", "Medium Armor", "Heavy Armor", "Shields", "Simple", "Martial"},
	"Ranger":    {"Light Armor", "Medium Armor", "Shields", "Simple", "Martial"},
	"Rogue":     {"Light Armor", "Simple", "Hand Crossbows", "Longswords", "Rapiers", "Shortswords"},
	"Sorcerer":  {},
	"Warlock":   {"Light Armor", "Simple"},
	"Wizard":    {},
}

// CanMulticlassInto checks if a character meets prerequisites for a class
func CanMulticlassInto(char *Character, className string) (bool, string) {
	prereqs, exists := MulticlassPrerequisites[className]
	if !exists {
		return false, fmt.Sprintf("Unknown class: %s", className)
	}

	// Check if prerequisites are met
	// For classes with multiple prerequisites, check the logic
	switch className {
	case "Fighter":
		// Fighter requires Str 13 OR Dex 13 (either one)
		strScore := char.AbilityScores.Strength
		dexScore := char.AbilityScores.Dexterity
		if strScore >= 13 || dexScore >= 13 {
			return true, ""
		}
		return false, "Requires Strength 13 or Dexterity 13"

	case "Monk", "Paladin", "Ranger":
		// These classes require BOTH abilities at 13
		for ability, required := range prereqs {
			score := char.GetAbilityScore(ability)
			if score < required {
				return false, fmt.Sprintf("Requires %s %d", ability, required)
			}
		}
		return true, ""

	default:
		// Single ability requirement
		for ability, required := range prereqs {
			score := char.GetAbilityScore(ability)
			if score < required {
				return false, fmt.Sprintf("Requires %s %d", ability, required)
			}
		}
		return true, ""
	}
}

// GetAbilityScore returns the total score for a given ability name
func (c *Character) GetAbilityScore(abilityName string) int {
	switch abilityName {
	case "Strength":
		return c.AbilityScores.Strength
	case "Dexterity":
		return c.AbilityScores.Dexterity
	case "Constitution":
		return c.AbilityScores.Constitution
	case "Intelligence":
		return c.AbilityScores.Intelligence
	case "Wisdom":
		return c.AbilityScores.Wisdom
	case "Charisma":
		return c.AbilityScores.Charisma
	default:
		return 0
	}
}

// GetAvailableClasses returns classes the character can multiclass into
func GetAvailableClasses(char *Character) []Class {
	allClasses := GetAllClasses()
	available := make([]Class, 0)

	for _, class := range allClasses {
		canMulticlass, _ := CanMulticlassInto(char, class.Name)
		if canMulticlass {
			available = append(available, class)
		}
	}

	return available
}

// GetMulticlassProficiencies returns the limited proficiencies granted when multiclassing
func GetMulticlassProficiencies(className string) []string {
	if profs, exists := MulticlassProficiencies[className]; exists {
		return profs
	}
	return []string{}
}

// GetClassLevel returns the level in a specific class, or 0 if not present
func (c *Character) GetClassLevel(className string) int {
	for _, cl := range c.Classes {
		if cl.ClassName == className {
			return cl.Level
		}
	}
	return 0
}

// HasClass checks if character has any levels in a class
func (c *Character) HasClass(className string) bool {
	return c.GetClassLevel(className) > 0
}

// GetPrimaryClass returns the class with the highest level (first if tied)
func (c *Character) GetPrimaryClass() *ClassLevel {
	if len(c.Classes) == 0 {
		return nil
	}

	primary := &c.Classes[0]
	for i := range c.Classes {
		if c.Classes[i].Level > primary.Level {
			primary = &c.Classes[i]
		}
	}
	return primary
}

// GetClassDisplayString returns formatted class string (e.g., "Fighter 3 / Druid 2")
func (c *Character) GetClassDisplayString() string {
	if len(c.Classes) == 0 {
		return "No Class"
	}

	if len(c.Classes) == 1 {
		cl := c.Classes[0]
		if cl.Subclass != "" {
			return fmt.Sprintf("%s (%s) %d", cl.ClassName, cl.Subclass, cl.Level)
		}
		return fmt.Sprintf("%s %d", cl.ClassName, cl.Level)
	}

	// Multiple classes
	result := ""
	for i, cl := range c.Classes {
		if i > 0 {
			result += " / "
		}
		result += fmt.Sprintf("%s %d", cl.ClassName, cl.Level)
	}
	return result
}

// CalculateTotalLevel returns the sum of all class levels
func (c *Character) CalculateTotalLevel() int {
	total := 0
	for _, cl := range c.Classes {
		total += cl.Level
	}
	return total
}
