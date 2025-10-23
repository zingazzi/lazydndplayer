// internal/models/levelup.go
package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LevelUpResult contains information about what was gained on level up
type LevelUpResult struct {
	ClassName        string
	NewClassLevel    int
	NewTotalLevel    int
	HPGained         int
	FeaturesGained   []string
	ProficienciesGained []string
	SpellSlotsGained string
	IsNewClass       bool
	RequiresSubclass bool
	RequiresSkills   bool
	RequiresSpells   bool
}

// LevelUpOptions contains options for leveling up
type LevelUpOptions struct {
	ClassName      string
	TakeAverage    bool // If true, take average HP instead of rolling
	SelectedSkills []SkillType
	Subclass       string
	FightingStyle  string
}

// CanLevelUp checks if a character can level up in a specific class
func CanLevelUp(char *Character, className string) (bool, string) {
	// Check if it's a new class (multiclassing)
	isNewClass := !char.HasClass(className)

	if isNewClass {
		// Check multiclass prerequisites
		canMulticlass, reason := CanMulticlassInto(char, className)
		if !canMulticlass {
			return false, reason
		}
	}

	// Character can always level up in existing class
	return true, ""
}

// LevelUp levels up a character in the specified class
func LevelUp(char *Character, options LevelUpOptions) (*LevelUpResult, error) {
	debug.Log("LevelUp: class=%s, character=%s", options.ClassName, char.Name)

	// Validate
	canLevel, reason := CanLevelUp(char, options.ClassName)
	if !canLevel {
		return nil, fmt.Errorf("cannot level up: %s", reason)
	}

	// Load class data
	classData := GetClassByName(options.ClassName)
	if classData == nil {
		return nil, fmt.Errorf("class not found: %s", options.ClassName)
	}

	isNewClass := !char.HasClass(options.ClassName)
	var classLevel *ClassLevel
	newClassLevel := 1

	if isNewClass {
		// Add new class
		classLevel = &ClassLevel{
			ClassName:     options.ClassName,
			Level:         1,
			Subclass:      options.Subclass,
			FightingStyle: options.FightingStyle,
		}
		char.Classes = append(char.Classes, *classLevel)
		debug.Log("Added new class: %s level 1", options.ClassName)
	} else {
		// Level up existing class
		for i := range char.Classes {
			if char.Classes[i].ClassName == options.ClassName {
				char.Classes[i].Level++
				classLevel = &char.Classes[i]
				newClassLevel = classLevel.Level

				// Update subclass if provided (for classes that get it at higher levels)
				if options.Subclass != "" {
					char.Classes[i].Subclass = options.Subclass
				}
				if options.FightingStyle != "" {
					char.Classes[i].FightingStyle = options.FightingStyle
				}
				break
			}
		}
		debug.Log("Leveled up %s to level %d", options.ClassName, newClassLevel)
	}

	// Calculate new total level
	newTotalLevel := char.CalculateTotalLevel()

	// Roll/calculate HP
	hpGained := RollHP(classData.HitDie, char.AbilityScores.GetModifier(Constitution), options.TakeAverage)
	char.MaxHP += hpGained
	char.CurrentHP += hpGained
	debug.Log("HP gained: %d (new max: %d)", hpGained, char.MaxHP)

	// Apply level benefits
	result := &LevelUpResult{
		ClassName:     options.ClassName,
		NewClassLevel: newClassLevel,
		NewTotalLevel: newTotalLevel,
		HPGained:      hpGained,
		IsNewClass:    isNewClass,
	}

	// Apply proficiencies (only for new class, and only limited ones for multiclass)
	if isNewClass {
		if len(char.Classes) == 1 {
			// First class - grant all starting proficiencies
			result.ProficienciesGained = applyFullProficiencies(char, classData)
		} else {
			// Multiclassing - grant limited proficiencies
			result.ProficienciesGained = applyMulticlassProficiencies(char, options.ClassName)
		}
	}

	// Grant features for this level
	features := GrantLevelFeatures(char, classData, newClassLevel)
	result.FeaturesGained = features

	// Handle skill choices
	if isNewClass && classData.SkillChoices != nil {
		result.RequiresSkills = true
		if len(options.SelectedSkills) > 0 {
			ApplySkillChoices(char, options.SelectedSkills, options.ClassName)
		}
	}

	// Update spellcasting if applicable
	if classData.Spellcasting != nil {
		UpdateSpellcasting(char, classData, newClassLevel)
		result.SpellSlotsGained = "Spell slots updated"

		// Check if spell selection is needed
		casterInfo := GetClassCasterInfo(options.ClassName)
		if casterInfo != nil {
			if casterInfo.Method == KnownCaster || casterInfo.Method == SpellbookCaster {
				result.RequiresSpells = true
			}
		}
	}

	// Check if subclass selection is required
	if options.Subclass == "" {
		result.RequiresSubclass = RequiresSubclassAtLevel(classData, newClassLevel)
	}

	// Update derived stats
	char.Level = newTotalLevel
	char.TotalLevel = newTotalLevel
	char.UpdateDerivedStats()

	debug.Log("LevelUp complete: total level %d", newTotalLevel)
	return result, nil
}

// RollHP rolls hit points for a level up
func RollHP(hitDie int, constitutionMod int, takeAverage bool) int {
	var roll int

	if takeAverage {
		// Take average: (die / 2) + 1
		roll = (hitDie / 2) + 1
	} else {
		// Roll the die
		roll = rand.Intn(hitDie) + 1
	}

	// Add constitution modifier (minimum 1 HP per level)
	hpGained := roll + constitutionMod
	if hpGained < 1 {
		hpGained = 1
	}

	return hpGained
}

// GrantLevelFeatures grants features for a specific class level
func GrantLevelFeatures(char *Character, class *Class, level int) []string {
	grantedFeatures := []string{}

	// Add level 1 features if this is level 1
	if level == 1 && class.Level1Features != nil {
		for _, featureDef := range class.Level1Features {
			// Convert FeatureDefinition to Feature
			feature := Feature{
				Name:        featureDef.Name,
				Description: featureDef.Description,
				MaxUses:     ParseUsesFormula(featureDef.UsesFormula, char),
				CurrentUses: ParseUsesFormula(featureDef.UsesFormula, char),
				RestType:    featureDef.RestType,
				Source:      "Class: " + class.Name,
			}
			char.Features.AddFeature(feature)
			grantedFeatures = append(grantedFeatures, feature.Name)
			debug.Log("Granted feature: %s", feature.Name)
		}
	}

	return grantedFeatures
}

// ApplySkillChoices applies selected skill proficiencies
func ApplySkillChoices(char *Character, skills []SkillType, source string) {
	for _, skillType := range skills {
		// Find the skill in the character's skill list and set it to proficient
		for i := range char.Skills.List {
			if char.Skills.List[i].Name == skillType {
				char.Skills.List[i].Proficiency = Proficient
				break
			}
		}
		char.ClassSkills = append(char.ClassSkills, skillType)
		debug.Log("Granted skill proficiency: %s from %s", skillType, source)
	}
}

// applyFullProficiencies grants all starting proficiencies for a first class
func applyFullProficiencies(char *Character, class *Class) []string {
	granted := []string{}

	// Armor proficiencies
	for _, armor := range class.ArmorProficiencies {
		if !contains(char.ArmorProficiencies, armor) {
			char.ArmorProficiencies = append(char.ArmorProficiencies, armor)
			granted = append(granted, armor)
		}
	}

	// Weapon proficiencies
	for _, weapon := range class.WeaponProficiencies {
		if !contains(char.WeaponProficiencies, weapon) {
			char.WeaponProficiencies = append(char.WeaponProficiencies, weapon)
			granted = append(granted, weapon)
		}
	}

	// Tool proficiencies
	for _, tool := range class.ToolProficiencies {
		if !contains(char.ToolProficiencies, tool) {
			char.ToolProficiencies = append(char.ToolProficiencies, tool)
			granted = append(granted, tool)
		}
	}

	// Saving throw proficiencies
	for _, save := range class.SavingThrows {
		if !contains(char.SavingThrowProficiencies, save) {
			char.SavingThrowProficiencies = append(char.SavingThrowProficiencies, save)
			granted = append(granted, save+" Save")
		}
	}

	debug.Log("Applied full proficiencies: %v", granted)
	return granted
}

// applyMulticlassProficiencies grants limited proficiencies when multiclassing
func applyMulticlassProficiencies(char *Character, className string) []string {
	granted := []string{}
	multiclassProfs := GetMulticlassProficiencies(className)

	for _, prof := range multiclassProfs {
		// Check armor proficiencies
		if !contains(char.ArmorProficiencies, prof) &&
		   (prof == "Light" || prof == "Medium" || prof == "Heavy" || prof == "Shields") {
			char.ArmorProficiencies = append(char.ArmorProficiencies, prof)
			char.MulticlassProficiencies = append(char.MulticlassProficiencies, prof)
			granted = append(granted, prof)
			continue
		}

		// Check weapon proficiencies
		if !contains(char.WeaponProficiencies, prof) {
			char.WeaponProficiencies = append(char.WeaponProficiencies, prof)
			char.MulticlassProficiencies = append(char.MulticlassProficiencies, prof)
			granted = append(granted, prof)
		}
	}

	debug.Log("Applied multiclass proficiencies: %v", granted)
	return granted
}

// UpdateSpellcasting updates spell slots and cantrips known for the new level
func UpdateSpellcasting(char *Character, class *Class, classLevel int) {
	casterInfo := GetClassCasterInfo(class.Name)
	if casterInfo == nil {
		return
	}

	// Update cantrips known
	if cantrips, exists := casterInfo.CantripsKnownByLevel[classLevel]; exists {
		char.SpellBook.CantripsKnown = cantrips
		debug.Log("Cantrips known updated to: %d", cantrips)
	}

	// Update spell slots based on multiclass calculation
	char.SpellBook.Slots = CalculateMulticlassSpellSlots(char.Classes)

	// For Warlock, handle Pact Magic separately
	if class.Name == "Warlock" {
		warlockLevel := char.GetClassLevel("Warlock")
		slots, slotLevel := GetWarlockPactSlots(warlockLevel)
		debug.Log("Warlock Pact Magic: %d slots of level %d", slots, slotLevel)
		// Store this information (would need to extend SpellBook structure)
	}

	// Update spellcasting ability if not set
	if char.SpellBook.SpellcastingMod == "" {
		char.SpellBook.SpellcastingMod = casterInfo.SpellcastingAbility
	}

	// Update prepared spells limit for prepared casters
	if casterInfo.Method == PreparedCaster {
		char.SpellBook.IsPreparedCaster = true
		char.SpellBook.PreparationFormula = casterInfo.PreparationFormula
		// Max prepared will be calculated in UpdateDerivedStats
	}
}

// RequiresSubclassAtLevel checks if a class requires subclass selection at a given level
func RequiresSubclassAtLevel(class *Class, level int) bool {
	// Check if class has subclasses defined
	if len(class.Subclasses) == 0 {
		return false
	}

	// Classes that get subclass at level 1
	level1Subclasses := []string{"Cleric", "Sorcerer", "Warlock"}
	for _, name := range level1Subclasses {
		if class.Name == name && level == 1 {
			return true
		}
	}

	// Classes that get subclass at level 2
	level2Subclasses := []string{"Druid", "Wizard"}
	for _, name := range level2Subclasses {
		if class.Name == name && level == 2 {
			return true
		}
	}

	// Classes that get subclass at level 3
	level3Subclasses := []string{"Bard", "Fighter", "Monk", "Paladin", "Ranger", "Rogue", "Barbarian"}
	for _, name := range level3Subclasses {
		if class.Name == name && level == 3 {
			return true
		}
	}

	return false
}

// ParseUsesFormula parses a uses formula string and returns the value
func ParseUsesFormula(formula string, char *Character) int {
	if formula == "" || formula == "0" {
		return 0
	}
	if formula == "1" {
		return 1
	}
	// Handle formulas like "proficiency_bonus" or "wisdom_mod"
	// For now, simple handling
	return 1
}

// Helper functions removed since RestType is already the correct type

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// GetLevelUpPreview generates a preview of what will be gained on level up
func GetLevelUpPreview(char *Character, className string) (*LevelUpResult, error) {
	classData := GetClassByName(className)
	if classData == nil {
		return nil, fmt.Errorf("class not found: %s", className)
	}

	isNewClass := !char.HasClass(className)
	newClassLevel := 1
	if !isNewClass {
		newClassLevel = char.GetClassLevel(className) + 1
	}

	newTotalLevel := char.CalculateTotalLevel()
	if isNewClass {
		newTotalLevel++
	} else {
		newTotalLevel++
	}

	// Calculate average HP gain
	avgHP := RollHP(classData.HitDie, char.AbilityScores.GetModifier(Constitution), true)

	preview := &LevelUpResult{
		ClassName:        className,
		NewClassLevel:    newClassLevel,
		NewTotalLevel:    newTotalLevel,
		HPGained:         avgHP,
		IsNewClass:       isNewClass,
		RequiresSubclass: RequiresSubclassAtLevel(classData, newClassLevel),
		RequiresSkills:   isNewClass && classData.SkillChoices != nil,
	}

	// Check if spell selection is needed
	if classData.Spellcasting != nil {
		casterInfo := GetClassCasterInfo(className)
		if casterInfo != nil && (casterInfo.Method == KnownCaster || casterInfo.Method == SpellbookCaster) {
			preview.RequiresSpells = true
		}
	}

	return preview, nil
}
