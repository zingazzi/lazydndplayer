// internal/models/classes.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Class represents a D&D 5e class
type Class struct {
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	HitDie               int                 `json:"hit_die"`
	PrimaryAbility       string              `json:"primary_ability"`
	SavingThrows         []string            `json:"saving_throws"`
	ArmorProficiencies   []string            `json:"armor_proficiencies"`
	WeaponProficiencies  []string            `json:"weapon_proficiencies"`
	ToolProficiencies    []string            `json:"tool_proficiencies"`
	SkillChoices         *SkillChoiceInfo    `json:"skill_choices"`
	StartingEquipment    []string            `json:"starting_equipment"`
	Spellcasting         *SpellcastingInfo   `json:"spellcasting"`
	Level1Features       []FeatureDefinition `json:"level_1_features"`
}

// SkillChoiceInfo defines how many skills to choose and from which list
type SkillChoiceInfo struct {
	Choose int      `json:"choose"`
	From   []string `json:"from"`
}

// SpellcastingInfo contains spellcasting information for a class
type SpellcastingInfo struct {
	Ability        string         `json:"ability"`
	CantripsKnown  int            `json:"cantrips_known"`
	SpellsKnown    int            `json:"spells_known,omitempty"`
	SpellsPrepared string         `json:"spells_prepared,omitempty"`
	SpellSlots     map[string]int `json:"spell_slots"`
}

// ClassesData represents the structure of classes.json
type ClassesData struct {
	Classes []Class `json:"classes"`
}

var cachedClasses *ClassesData

// LoadClassesFromJSON loads all classes from the JSON file
func LoadClassesFromJSON(filepath string) (*ClassesData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read classes file: %w", err)
	}

	var classesData ClassesData
	if err := json.Unmarshal(data, &classesData); err != nil {
		return nil, fmt.Errorf("failed to parse classes JSON: %w", err)
	}

	cachedClasses = &classesData
	return &classesData, nil
}

// GetAllClasses returns all available classes
func GetAllClasses() []Class {
	if cachedClasses == nil {
		_, err := LoadClassesFromJSON("data/classes.json")
		if err != nil {
			fmt.Printf("Error loading classes: %v\n", err)
			return []Class{}
		}
	}
	return cachedClasses.Classes
}

// GetClassByName returns a specific class by name
func GetClassByName(name string) *Class {
	classes := GetAllClasses()
	for i := range classes {
		if classes[i].Name == name {
			return &classes[i]
		}
	}
	return nil
}

// CalculateMaxHP calculates the maximum HP for a character
// Formula: (Hit Die * Level) + (CON modifier * Level) + bonuses
// At level 1, use maximum hit die value instead of rolling
func CalculateMaxHP(char *Character, class *Class) int {
	if class == nil {
		return char.MaxHP // Return current if class not found
	}

	level := char.Level
	if level < 1 {
		level = 1
	}

	// Base HP calculation
	var baseHP int
	if level == 1 {
		// Level 1: Always use maximum hit die value
		baseHP = class.HitDie
	} else {
		// Higher levels: Hit die per level (simplified - should track actual rolls)
		// For now, use average: (HitDie/2 + 1) per level after 1st
		baseHP = class.HitDie // First level (max)
		averagePerLevel := (class.HitDie / 2) + 1
		baseHP += averagePerLevel * (level - 1)
	}

	// Add Constitution modifier per level
	conModifier := char.AbilityScores.GetModifier("Constitution")
	totalConBonus := conModifier * level

	// Add any HP bonuses from feats, species, etc.
	// These are tracked separately and added to base HP
	bonusHP := char.SpeciesHPBonus

	// Check benefit tracker for HP bonuses
	if char.BenefitTracker != nil {
		hpBenefits := char.BenefitTracker.GetBenefitsByType(BenefitHP)
		for _, benefit := range hpBenefits {
			bonusHP += benefit.Value
		}
	}

	totalHP := baseHP + totalConBonus + bonusHP

	// Minimum 1 HP
	if totalHP < 1 {
		totalHP = 1
	}

	return totalHP
}

// ApplyClassToCharacter applies a class to a character and updates HP
func ApplyClassToCharacter(char *Character, className string) error {
	class := GetClassByName(className)
	if class == nil {
		return fmt.Errorf("class %s not found", className)
	}

	// Update class name
	char.Class = className

	// Apply armor proficiencies
	char.ArmorProficiencies = make([]string, len(class.ArmorProficiencies))
	copy(char.ArmorProficiencies, class.ArmorProficiencies)

	// Apply weapon proficiencies
	char.WeaponProficiencies = make([]string, len(class.WeaponProficiencies))
	copy(char.WeaponProficiencies, class.WeaponProficiencies)

	// Apply saving throw proficiencies
	char.SavingThrowProficiencies = make([]string, len(class.SavingThrows))
	copy(char.SavingThrowProficiencies, class.SavingThrows)

	// Calculate and set HP
	newMaxHP := CalculateMaxHP(char, class)

	// If this is a new character or HP is at max, heal to full
	if char.CurrentHP == char.MaxHP || char.MaxHP == 0 {
		char.CurrentHP = newMaxHP
	} else {
		// Adjust current HP proportionally to avoid full heal exploit
		hpRatio := float64(char.CurrentHP) / float64(char.MaxHP)
		char.CurrentHP = int(float64(newMaxHP) * hpRatio)
		if char.CurrentHP < 1 {
			char.CurrentHP = 1
		}
	}

	char.MaxHP = newMaxHP

	// Update derived stats
	char.UpdateDerivedStats()

	return nil
}

// HasArmorProficiency checks if character is proficient with an armor type
func HasArmorProficiency(char *Character, armorType string) bool {
	// Normalize armor type for comparison (case-insensitive)
	normalizedArmorType := strings.ToLower(armorType)

	for _, prof := range char.ArmorProficiencies {
		normalizedProf := strings.ToLower(prof)

		// Handle "shield" vs "shields" plural
		if normalizedProf == "shields" {
			normalizedProf = "shield"
		}

		if normalizedProf == normalizedArmorType {
			return true
		}
	}
	return false
}

// HasWeaponProficiency checks if character is proficient with a weapon type
func HasWeaponProficiency(char *Character, weaponType string) bool {
	// Normalize weapon type for comparison (case-insensitive)
	normalizedWeaponType := strings.ToLower(weaponType)

	for _, prof := range char.WeaponProficiencies {
		normalizedProf := strings.ToLower(prof)

		// Check if the weapon type matches (e.g., "simple melee" contains "simple")
		if strings.Contains(normalizedWeaponType, normalizedProf) {
			return true
		}
	}
	return false
}
