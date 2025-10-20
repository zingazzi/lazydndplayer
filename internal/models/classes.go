// internal/models/classes.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
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
	Ability              string         `json:"ability"`
	RitualCasting        bool           `json:"ritual_casting,omitempty"`
	SpellsKnownFormula   string         `json:"spells_known_formula,omitempty"` // "all" for prepared casters, formula for known casters
	CantripsKnown        int            `json:"cantrips_known"`
	PreparationFormula   string         `json:"preparation_formula,omitempty"` // e.g., "wisdom+level"
	SpellSlots           map[string]int `json:"spell_slots,omitempty"`
}

// ClassesData represents the structure of classes.json
type ClassesData struct {
	Classes []Class `json:"classes"`
}

var cachedClasses *ClassesData

// LoadClassesFromJSON loads all classes from individual JSON files in the directory
func LoadClassesFromJSON(dirpath string) (*ClassesData, error) {
	// Read all files in the classes directory
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read classes directory: %w", err)
	}

	var classes []Class

	// Load each class file
	for _, file := range files {
		// Skip non-JSON files and the README
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dirpath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to read class file %s: %v\n", file.Name(), err)
			continue
		}

		var class Class
		if err := json.Unmarshal(data, &class); err != nil {
			fmt.Printf("Warning: failed to parse class file %s: %v\n", file.Name(), err)
			continue
		}

		classes = append(classes, class)
	}

	classesData := &ClassesData{Classes: classes}
	cachedClasses = classesData
	return classesData, nil
}

// GetAllClasses returns all available classes
func GetAllClasses() []Class {
	if cachedClasses == nil {
		_, err := LoadClassesFromJSON("data/classes")
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

	// Note: char.Class might be empty due to BackupClass clearing it
	// So we can't use char.Class to identify old features
	// Features must be removed by checking their Source field
	debug.Log("ApplyClassToCharacter: Applying %s (current char.Class='%s')", className, char.Class)

	// Remove previous class features by checking Source field
	RemoveAllClassFeatures(char)

	// Remove previous fighting style if changing class
	if char.FightingStyle != "" {
		RemoveFightingStyle(char)
	}

	// Update class name
	char.Class = className

	// Apply armor proficiencies (replaces old ones)
	char.ArmorProficiencies = make([]string, len(class.ArmorProficiencies))
	copy(char.ArmorProficiencies, class.ArmorProficiencies)

	// Apply weapon proficiencies (replaces old ones)
	char.WeaponProficiencies = make([]string, len(class.WeaponProficiencies))
	copy(char.WeaponProficiencies, class.WeaponProficiencies)

	// Apply saving throw proficiencies (replaces old ones)
	char.SavingThrowProficiencies = make([]string, len(class.SavingThrows))
	copy(char.SavingThrowProficiencies, class.SavingThrows)

	// Grant level 1 class features (this will add new class features)
	GrantClassFeatures(char, class)

	// Initialize spellcasting for spellcasting classes
	if class.Spellcasting != nil {
		InitializeSpellcasting(char, class)
	}

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

// GrantClassFeatures grants level 1 features from a class
func GrantClassFeatures(char *Character, class *Class) {
	// Add level 1 features (removal is now done in ApplyClassToCharacter)
	for _, featureDef := range class.Level1Features {
		// Convert definition to actual feature
		source := fmt.Sprintf("Class: %s", class.Name)
		feature := featureDef.ToFeature(char, source)
		char.Features.AddFeature(feature)
	}
}

// RemoveClassFeatures removes all features from a specific class
func RemoveClassFeatures(char *Character, className string) {
	debug.Log("RemoveClassFeatures: Called for class='%s'", className)

	if className == "" {
		debug.Log("RemoveClassFeatures: No class name provided, returning")
		return
	}

	// Remove features that came from the class
	sourcePrefix := fmt.Sprintf("Class: %s", className)
	newFeatures := []Feature{}

	debug.Log("RemoveClassFeatures: Looking for source='%s'", sourcePrefix)
	debug.Log("RemoveClassFeatures: Total features before: %d", len(char.Features.Features))

	for i, feature := range char.Features.Features {
		debug.Log("  Feature[%d]: Name='%s', Source='%s', Match=%v",
			i, feature.Name, feature.Source, feature.Source == sourcePrefix)
		if feature.Source != sourcePrefix {
			newFeatures = append(newFeatures, feature)
		} else {
			debug.Log("    -> REMOVING this feature")
		}
	}

	debug.Log("RemoveClassFeatures: Total features after: %d", len(newFeatures))
	char.Features.Features = newFeatures
}

// RemoveAllClassFeatures removes features from ANY class (by checking Source prefix)
func RemoveAllClassFeatures(char *Character) {
	debug.Log("RemoveAllClassFeatures: Called")
	newFeatures := []Feature{}

	debug.Log("RemoveAllClassFeatures: Total features before: %d", len(char.Features.Features))

	for i, feature := range char.Features.Features {
		isClassFeature := strings.HasPrefix(feature.Source, "Class: ")
		debug.Log("  Feature[%d]: Name='%s', Source='%s', IsClassFeature=%v",
			i, feature.Name, feature.Source, isClassFeature)
		if !isClassFeature {
			newFeatures = append(newFeatures, feature)
		} else {
			debug.Log("    -> REMOVING this class feature")
		}
	}

	debug.Log("RemoveAllClassFeatures: Total features after: %d", len(newFeatures))
	char.Features.Features = newFeatures
}

// InitializeSpellcasting sets up spellcasting for a class
func InitializeSpellcasting(char *Character, class *Class) {
	if class.Spellcasting == nil {
		return
	}

	debug.Log("InitializeSpellcasting: Setting up spellcasting for %s", class.Name)

	// Set spellcasting ability
	char.SpellBook.SpellcastingMod = AbilityType(class.Spellcasting.Ability)

	// Determine if this is a prepared caster (knows all spells but must prepare them)
	isPreparedCaster := class.Spellcasting.SpellsKnownFormula == "all"
	char.SpellBook.IsPreparedCaster = isPreparedCaster
	debug.Log("  IsPreparedCaster: %v", isPreparedCaster)

	// Set cantrips known
	char.SpellBook.CantripsKnown = class.Spellcasting.CantripsKnown
	debug.Log("  CantripsKnown: %d", char.SpellBook.CantripsKnown)

	// Calculate max prepared spells if this is a prepared caster
	if isPreparedCaster && class.Spellcasting.PreparationFormula != "" {
		maxPrepared := CalculateMaxPreparedSpells(char, class.Spellcasting.PreparationFormula)
		char.SpellBook.MaxPreparedSpells = maxPrepared
		debug.Log("  MaxPreparedSpells: %d (formula: %s)", maxPrepared, class.Spellcasting.PreparationFormula)
	}

	// Load spell slots from class progression
	LoadSpellSlotsForLevel(char, class, char.Level)

	// Update spell save DC and attack bonus
	char.UpdateDerivedStats()
}

// LoadSpellSlotsForLevel loads spell slots from class progression data
func LoadSpellSlotsForLevel(char *Character, class *Class, level int) {
	// Load class data from JSON to get spell slots
	classData := GetClassByName(class.Name)
	if classData == nil {
		debug.Log("LoadSpellSlotsForLevel: Could not load class data for %s", class.Name)
		return
	}

	// Find the level progression data
	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(class.Name))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("LoadSpellSlotsForLevel: Error reading class file: %v", err)
		return
	}

	var classWithProgression struct {
		LevelProgression []struct {
			Level            int                   `json:"level"`
			SpellcastingInfo *struct {
				SpellSlots map[string]int `json:"spell_slots"`
			} `json:"spellcasting_info"`
		} `json:"level_progression"`
	}

	if err := json.Unmarshal(data, &classWithProgression); err != nil {
		debug.Log("LoadSpellSlotsForLevel: Error unmarshaling class data: %v", err)
		return
	}

	// Find spell slots for current level
	for _, levelData := range classWithProgression.LevelProgression {
		if levelData.Level == level && levelData.SpellcastingInfo != nil {
			debug.Log("LoadSpellSlotsForLevel: Found spell slots for level %d", level)

			// Set spell slots
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["1"]; ok {
				char.SpellBook.Slots.Level1.Maximum = slots
				char.SpellBook.Slots.Level1.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["2"]; ok {
				char.SpellBook.Slots.Level2.Maximum = slots
				char.SpellBook.Slots.Level2.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["3"]; ok {
				char.SpellBook.Slots.Level3.Maximum = slots
				char.SpellBook.Slots.Level3.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["4"]; ok {
				char.SpellBook.Slots.Level4.Maximum = slots
				char.SpellBook.Slots.Level4.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["5"]; ok {
				char.SpellBook.Slots.Level5.Maximum = slots
				char.SpellBook.Slots.Level5.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["6"]; ok {
				char.SpellBook.Slots.Level6.Maximum = slots
				char.SpellBook.Slots.Level6.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["7"]; ok {
				char.SpellBook.Slots.Level7.Maximum = slots
				char.SpellBook.Slots.Level7.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["8"]; ok {
				char.SpellBook.Slots.Level8.Maximum = slots
				char.SpellBook.Slots.Level8.Current = slots
			}
			if slots, ok := levelData.SpellcastingInfo.SpellSlots["9"]; ok {
				char.SpellBook.Slots.Level9.Maximum = slots
				char.SpellBook.Slots.Level9.Current = slots
			}

			debug.Log("  Level 1 slots: %d", char.SpellBook.Slots.Level1.Maximum)
			break
		}
	}
}

// CalculateMaxPreparedSpells calculates how many spells can be prepared
// Formula examples: "wisdom+level", "intelligence+level"
func CalculateMaxPreparedSpells(char *Character, formula string) int {
	parts := strings.Split(strings.ToLower(formula), "+")
	total := 0

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "level":
			total += char.Level
		case "wisdom", "wis":
			total += char.AbilityScores.GetModifier("Wisdom")
		case "intelligence", "int":
			total += char.AbilityScores.GetModifier("Intelligence")
		case "charisma", "cha":
			total += char.AbilityScores.GetModifier("Charisma")
		}
	}

	// Minimum of 1
	if total < 1 {
		total = 1
	}

	return total
}
