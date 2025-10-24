// internal/models/levelup.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

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

// DeLevelResult contains information about what was removed when de-leveling
type DeLevelResult struct {
	ClassName       string
	OldClassLevel   int
	NewClassLevel   int
	NewTotalLevel   int
	HPLost          int
	FeaturesRemoved []string
	ClassRemoved    bool // True if the class was completely removed
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

// LevelUp levels up a character in the specified class.
// This is the main entry point that orchestrates the level-up process.
func LevelUp(char *Character, options LevelUpOptions) (*LevelUpResult, error) {
	debug.Log("LevelUp: class=%s, character=%s", options.ClassName, char.Name)

	// Step 1: Validate the level-up request
	classData, err := validateLevelUp(char, options.ClassName)
	if err != nil {
		return nil, err
	}

	// Step 2: Update class level (add new class or increment existing)
	classLevel, isNewClass := updateClassLevel(char, classData, options)
	newTotalLevel := char.CalculateTotalLevel()

	// Step 2.5: Update character level BEFORE granting features (needed for level-based formulas)
	char.Level = newTotalLevel
	char.TotalLevel = newTotalLevel
	debug.Log("Updated character level to %d before granting features", newTotalLevel)

	// Step 3: Calculate and apply HP gain
	hpGained := applyHPGain(char, classData, options.TakeAverage)

	// Step 4: Initialize result structure
	result := &LevelUpResult{
		ClassName:     options.ClassName,
		NewClassLevel: classLevel.Level,
		NewTotalLevel: newTotalLevel,
		HPGained:      hpGained,
		IsNewClass:    isNewClass,
	}

	// Step 5: Apply proficiencies for new classes
	if isNewClass {
		result.ProficienciesGained = applyClassProficiencies(char, classData, options.ClassName)
	}

	// Step 6: Grant class features for this level
	result.FeaturesGained = GrantLevelFeatures(char, classData, classLevel.Level)

	// Step 7: Handle skill selection requirements
	result.RequiresSkills = handleSkillChoices(char, classData, options, isNewClass)

	// Step 8: Update spellcasting capabilities
	result.SpellSlotsGained, result.RequiresSpells = updateClassSpellcasting(char, classData, classLevel.Level, options.ClassName)

	// Step 9: Check if subclass selection is required or grant subclass features
	result.RequiresSubclass = checkSubclassRequirement(classData, classLevel.Level, options.Subclass)

	// Step 9.5: If subclass was provided, grant subclass features
	if options.Subclass != "" {
		debug.Log("LevelUp: Granting subclass features for %s at level %d", options.Subclass, classLevel.Level)
		subclassFeatures := GrantSubclassFeatures(char, options.ClassName, options.Subclass, classLevel.Level)
		result.FeaturesGained = append(result.FeaturesGained, subclassFeatures...)

		// Update character's subclass
		for i := range char.Classes {
			if char.Classes[i].ClassName == options.ClassName {
				char.Classes[i].Subclass = options.Subclass
				break
			}
		}
	}

	// Step 10: Finalize by updating all derived stats
	finalizeLevelUp(char, newTotalLevel)

	debug.Log("LevelUp complete: total level %d", newTotalLevel)
	return result, nil
}

// validateLevelUp checks if the character can level up and loads class data.
func validateLevelUp(char *Character, className string) (*Class, error) {
	canLevel, reason := CanLevelUp(char, className)
	if !canLevel {
		return nil, fmt.Errorf("cannot level up: %s", reason)
	}

	classData := GetClassByName(className)
	if classData == nil {
		return nil, fmt.Errorf("class not found: %s", className)
	}

	return classData, nil
}

// updateClassLevel adds a new class or increments an existing class level.
// Returns the ClassLevel pointer and whether this is a new class.
func updateClassLevel(char *Character, classData *Class, options LevelUpOptions) (*ClassLevel, bool) {
	isNewClass := !char.HasClass(options.ClassName)

	if isNewClass {
		// Add new class at level 1
		newClass := ClassLevel{
			ClassName:     options.ClassName,
			Level:         1,
			Subclass:      options.Subclass,
			FightingStyle: options.FightingStyle,
		}
		char.Classes = append(char.Classes, newClass)
		debug.Log("Added new class: %s level 1", options.ClassName)
		return &char.Classes[len(char.Classes)-1], true
	}

	// Level up existing class
	for i := range char.Classes {
		if char.Classes[i].ClassName == options.ClassName {
			char.Classes[i].Level++

			// Update subclass/fighting style if provided
			if options.Subclass != "" {
				char.Classes[i].Subclass = options.Subclass
			}
			if options.FightingStyle != "" {
				char.Classes[i].FightingStyle = options.FightingStyle
			}

			debug.Log("Leveled up %s to level %d", options.ClassName, char.Classes[i].Level)
			return &char.Classes[i], false
		}
	}

	// Should never reach here due to validation
	return nil, false
}

// applyHPGain rolls HP and adds it to the character.
// Returns the amount of HP gained.
func applyHPGain(char *Character, classData *Class, takeAverage bool) int {
	conMod := char.AbilityScores.GetModifier(Constitution)
	hpGained := RollHP(classData.HitDie, conMod, takeAverage)

	char.MaxHP += hpGained
	char.CurrentHP += hpGained

	debug.Log("HP gained: %d (new max: %d)", hpGained, char.MaxHP)
	return hpGained
}

// applyClassProficiencies applies either full or multiclass proficiencies.
// Returns a list of proficiencies granted.
func applyClassProficiencies(char *Character, classData *Class, className string) []string {
	if len(char.Classes) == 1 {
		// First class - grant all starting proficiencies
		return applyFullProficiencies(char, classData)
	}
	// Multiclassing - grant limited proficiencies
	return applyMulticlassProficiencies(char, className)
}

// handleSkillChoices applies selected skills or marks that skill selection is required.
// Returns true if skill selection is still needed.
func handleSkillChoices(char *Character, classData *Class, options LevelUpOptions, isNewClass bool) bool {
	if !isNewClass || classData.SkillChoices == nil {
		return false
	}

	if len(options.SelectedSkills) > 0 {
		ApplySkillChoices(char, options.SelectedSkills, options.ClassName)
		return false
	}

	return true // Skills selection required but not provided
}

// updateClassSpellcasting updates spellcasting capabilities for spellcasting classes.
// Returns a message about spell slots and whether spell selection is required.
func updateClassSpellcasting(char *Character, classData *Class, classLevel int, className string) (string, bool) {
	if classData.Spellcasting == nil {
		return "", false
	}

	UpdateSpellcasting(char, classData, classLevel)

	// Check if spell selection is needed
	casterInfo := GetClassCasterInfo(className)
	requiresSpells := false
	if casterInfo != nil {
		if casterInfo.Method == KnownCaster || casterInfo.Method == SpellbookCaster {
			requiresSpells = true
		}
	}

	return "Spell slots updated", requiresSpells
}

// checkSubclassRequirement determines if subclass selection is needed at this level.
func checkSubclassRequirement(classData *Class, classLevel int, selectedSubclass string) bool {
	if selectedSubclass != "" {
		return false // Already selected
	}
	return RequiresSubclassAtLevel(classData, classLevel)
}

// finalizeLevelUp updates all derived character statistics.
func finalizeLevelUp(char *Character, newTotalLevel int) {
	char.Level = newTotalLevel
	char.TotalLevel = newTotalLevel

	// Ensure character has minimum XP for this level
	requiredXP := GetLevelXP(newTotalLevel)
	if char.Experience < requiredXP {
		debug.Log("finalizeLevelUp: Adjusting XP from %d to %d (required for level %d)",
			char.Experience, requiredXP, newTotalLevel)
		char.Experience = requiredXP
	}

	char.UpdateDerivedStats()
}

// RollHP rolls hit points for a level up using the default dice roller
func RollHP(hitDie int, constitutionMod int, takeAverage bool) int {
	return RollHPWithDiceRoller(hitDie, constitutionMod, takeAverage, GetDefaultDiceRoller())
}

// RollHPWithDiceRoller rolls hit points with a specified dice roller (for testing)
func RollHPWithDiceRoller(hitDie int, constitutionMod int, takeAverage bool, roller DiceRoller) int {
	var roll int

	if takeAverage {
		// Take average: (die / 2) + 1
		roll = (hitDie / 2) + 1
	} else {
		// Roll the die using injected roller
		roll = roller.Roll(hitDie)
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

	debug.Log("GrantLevelFeatures: Granting features for %s at level %d", class.Name, level)

	// Read class JSON to get level_progression
	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(class.Name))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("GrantLevelFeatures: Error reading class file: %v", err)
		// Fallback to Level1Features if level 1 and it exists
		if level == 1 && class.Level1Features != nil {
			for _, featureDef := range class.Level1Features {
				source := fmt.Sprintf("Class: %s", class.Name)
				feature := featureDef.ToFeature(char, source)
				char.Features.AddFeature(feature)
				grantedFeatures = append(grantedFeatures, feature.Name)
			}
		}
		return grantedFeatures
	}

	// Parse level_progression
	var classWithProgression struct {
		LevelProgression []struct {
			Level    int                   `json:"level"`
			Features []FeatureDefinition `json:"features"`
		} `json:"level_progression"`
	}

	if err := json.Unmarshal(data, &classWithProgression); err != nil {
		debug.Log("GrantLevelFeatures: Error unmarshaling class data: %v", err)
		return grantedFeatures
	}

	// Find features for the specified level
	for _, levelData := range classWithProgression.LevelProgression {
		if levelData.Level == level {
			debug.Log("GrantLevelFeatures: Found %d features for level %d", len(levelData.Features), level)
			for _, featureDef := range levelData.Features {
				// Skip subclass choice features
				if featureDef.Name == "Divine Domain" || featureDef.Name == "Sorcerous Origin" ||
				   featureDef.Name == "Otherworldly Patron" || featureDef.Name == "Monastic Tradition" ||
				   featureDef.Name == "Primal Path" || strings.Contains(strings.ToLower(featureDef.Description), "subclass") {
					debug.Log("  Skipping subclass choice feature: %s", featureDef.Name)
					continue
				}

				// Handle special Ki improvements
				if featureDef.Name == "Ki Improvement" {
					// Update existing Ki feature instead of adding a new one
					for i := range char.Features.Features {
						if char.Features.Features[i].Name == "Ki" {
							oldUses := char.Features.Features[i].MaxUses
							newUses := featureDef.CalculateMaxUses(char)
							char.Features.Features[i].MaxUses = newUses
							char.Features.Features[i].CurrentUses = newUses
							debug.Log("  Updated Ki: %d -> %d", oldUses, newUses)
							grantedFeatures = append(grantedFeatures, "Ki (increased to "+fmt.Sprint(newUses)+")")
							break
						}
					}
					continue
				}

				// Handle special Focus Point improvements (Monk)
				if featureDef.Name == "Focus Point Improvement" {
					// Update existing Focus Points feature instead of adding a new one
					for i := range char.Features.Features {
						if char.Features.Features[i].Name == "Focus Points" {
							oldUses := char.Features.Features[i].MaxUses
							newUses := featureDef.CalculateMaxUses(char)
							char.Features.Features[i].MaxUses = newUses
							char.Features.Features[i].CurrentUses = newUses
							debug.Log("  Updated Focus Points: %d -> %d", oldUses, newUses)
							grantedFeatures = append(grantedFeatures, "Focus Points (increased to "+fmt.Sprint(newUses)+")")
							break
						}
					}
					continue
				}

				source := fmt.Sprintf("Class: %s", class.Name)
				feature := featureDef.ToFeature(char, source)
				debug.Log("  Adding feature: %s (uses: %d/%d, rest: %s)",
					feature.Name, feature.CurrentUses, feature.MaxUses, feature.RestType)
				char.Features.AddFeature(feature)
				grantedFeatures = append(grantedFeatures, feature.Name)
			}
			break
		}
	}

	debug.Log("GrantLevelFeatures: Granted %d features", len(grantedFeatures))
	return grantedFeatures
}

// GrantSubclassFeatures grants features from a selected subclass at a specific level
func GrantSubclassFeatures(char *Character, className string, subclassName string, level int) []string {
	grantedFeatures := []string{}

	debug.Log("GrantSubclassFeatures: className=%s, subclassName=%s, level=%d", className, subclassName, level)

	// Read class JSON to get subclasses
	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(className))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("GrantSubclassFeatures: Error reading class file: %v", err)
		return grantedFeatures
	}

	// Parse subclasses
	var classWithSubclasses struct {
		Subclasses []struct {
			Name             string `json:"name"`
			SubclassLevel    int    `json:"subclass_level"`
			FeaturesByLevel  map[string][]FeatureDefinition `json:"features_by_level"`
		} `json:"subclasses"`
	}

	if err := json.Unmarshal(data, &classWithSubclasses); err != nil {
		debug.Log("GrantSubclassFeatures: Error unmarshaling: %v", err)
		return grantedFeatures
	}

	// Find the selected subclass
	for _, subclass := range classWithSubclasses.Subclasses {
		if subclass.Name == subclassName {
			debug.Log("GrantSubclassFeatures: Found subclass %s", subclassName)

			// Grant features for the specified level
			levelKey := fmt.Sprintf("%d", level)
			if features, ok := subclass.FeaturesByLevel[levelKey]; ok {
				debug.Log("GrantSubclassFeatures: Found %d features for level %d", len(features), level)
				for _, featureDef := range features {
					source := fmt.Sprintf("Subclass: %s", subclassName)
					feature := featureDef.ToFeature(char, source)
					debug.Log("  Adding subclass feature: %s (uses: %d/%d, rest: %s)",
						feature.Name, feature.CurrentUses, feature.MaxUses, feature.RestType)
					char.Features.AddFeature(feature)
					grantedFeatures = append(grantedFeatures, feature.Name)

					// Apply special benefits for certain subclass features
					applySubclassFeatureBenefits(char, feature.Name, subclassName)
				}
			} else {
				debug.Log("GrantSubclassFeatures: No features found for level %d", level)
			}
			break
		}
	}

	debug.Log("GrantSubclassFeatures: Granted %d subclass features", len(grantedFeatures))
	return grantedFeatures
}

// applySubclassFeatureBenefits applies special benefits for specific subclass features
func applySubclassFeatureBenefits(char *Character, featureName string, subclassName string) {
	source := BenefitSource{
		Type: "Subclass",
		Name: subclassName,
	}
	applier := NewBenefitApplier(char)

	debug.Log("applySubclassFeatureBenefits: Checking feature '%s' for subclass '%s'", featureName, subclassName)

	switch featureName {
	case "Implements of Mercy":
		// Warrior of Mercy: Grant Medicine, Insight proficiency, and Herbalism Kit
		debug.Log("  Applying Implements of Mercy benefits")

		// Add Medicine skill proficiency
		if err := applier.AddSkillProficiency(source, "Medicine"); err != nil {
			debug.Log("  Error adding Medicine proficiency: %v", err)
		} else {
			debug.Log("  Granted Medicine proficiency")
		}

		// Add Insight skill proficiency
		if err := applier.AddSkillProficiency(source, "Insight"); err != nil {
			debug.Log("  Error adding Insight proficiency: %v", err)
		} else {
			debug.Log("  Granted Insight proficiency")
		}

		// Add Herbalism Kit tool proficiency
		if err := applier.AddToolProficiency(source, "Herbalism Kit"); err != nil {
			debug.Log("  Error adding Herbalism Kit proficiency: %v", err)
		} else {
			debug.Log("  Granted Herbalism Kit proficiency")
		}

	case "Shadow Arts":
		// Warrior of Shadow: Grant Darkvision (60ft or +60ft)
		debug.Log("  Applying Shadow Arts benefits")

		if char.Darkvision == 0 {
			// No darkvision, grant 60ft
			char.Darkvision = 60
			debug.Log("  Granted Darkvision 60ft")
		} else {
			// Already has darkvision, increase by 60ft
			oldDarkvision := char.Darkvision
			char.Darkvision += 60
			debug.Log("  Increased Darkvision from %dft to %dft", oldDarkvision, char.Darkvision)
		}

		// Track darkvision benefit
		char.BenefitTracker.AddBenefit(GrantedBenefit{
			Source:      source,
			Type:        "Darkvision",
			Target:      "Darkvision",
			Value:       60,
			Description: "Darkvision 60ft (or +60ft)",
		})

		// Add "Darkness" spell
		darknessSpell := GetSpellByName("Darkness")
		if darknessSpell != nil {
			darknessSpell.Prepared = true
			char.SpellBook.AddSpell(*darknessSpell)
			debug.Log("  Granted Darkness spell")
		} else {
			debug.Log("  Warning: Darkness spell not found")
		}

		// Add "Minor Illusion" cantrip
		minorIllusion := GetSpellByName("Minor Illusion")
		if minorIllusion != nil && minorIllusion.Level == 0 {
			char.SpellBook.Cantrips = append(char.SpellBook.Cantrips, minorIllusion.Name)
			char.SpellBook.CantripsKnown++
			debug.Log("  Granted Minor Illusion cantrip")
		} else {
			debug.Log("  Warning: Minor Illusion cantrip not found")
		}

	case "Elemental Attunement":
		// Warrior of Elements: Grant Elementalism spell
		debug.Log("  Applying Elemental Attunement benefits")

		elementalismSpell := GetSpellByName("Elementalism")
		if elementalismSpell != nil && elementalismSpell.Level == 0 {
			char.SpellBook.Cantrips = append(char.SpellBook.Cantrips, elementalismSpell.Name)
			char.SpellBook.CantripsKnown++
			debug.Log("  Granted Elementalism cantrip")
		} else {
			debug.Log("  Warning: Elementalism cantrip not found")
		}

	case "Psionic Power":
		// Psi Warrior: Initialize Psi Dice
		debug.Log("  Applying Psionic Power benefits - Initializing Psi Dice")
		fighterLevel := char.GetFighterLevel()

		// Get psi dice count and size from scaling tables
		psiDiceCount := GetFeatureScaling("Psi Warrior", "Psionic Power", fighterLevel)
		psiDiceSize := GetPsiDiceSize(fighterLevel)

		char.PsiDice.Max = psiDiceCount
		char.PsiDice.Current = psiDiceCount // Start fully charged
		char.PsiDice.Size = psiDiceSize

		debug.Log("  Initialized Psi Dice: %d%s (%d/%d)", psiDiceCount, psiDiceSize, char.PsiDice.Current, char.PsiDice.Max)

	case "Combat Superiority":
		// Battle Master: Initialize Superiority Dice
		debug.Log("  Applying Combat Superiority benefits - Initializing Superiority Dice")
		fighterLevel := char.GetFighterLevel()

		// Get superiority dice count and size from scaling tables
		supDiceCount := GetFeatureScaling("Battle Master", "Combat Superiority", fighterLevel)
		supDiceSize := GetSuperiorityDiceSize(fighterLevel)

		char.SuperiorityDice.Max = supDiceCount
		char.SuperiorityDice.Current = supDiceCount // Start fully charged
		char.SuperiorityDice.Size = supDiceSize

		debug.Log("  Initialized Superiority Dice: %d%s (%d/%d)", supDiceCount, supDiceSize, char.SuperiorityDice.Current, char.SuperiorityDice.Max)

	case "Spellcasting":
		// Eldritch Knight: Set up spellcasting
		if subclassName == "Eldritch Knight" {
			debug.Log("  Applying Eldritch Knight Spellcasting benefits")
			fighterLevel := char.GetFighterLevel()

			// Initialize spell slots for level 3
			if fighterLevel == 3 {
				// Set spellcasting ability
				char.SpellBook.SpellcastingMod = Intelligence

				// Grant 2 level 1 spell slots
				char.SpellBook.Slots.Level1.Maximum = 2
				char.SpellBook.Slots.Level1.Current = 2

				debug.Log("  Set up Eldritch Knight spellcasting: 2x 1st level slots, Intelligence ability")

				// Note: Cantrip and spell selection will be prompted in the UI
			}
		}

	case "Student of War":
		// Battle Master: This will be handled in the UI to prompt for tool/skill selection
		debug.Log("  Student of War benefits will be handled in UI")
	}
}

// removeSubclassFeatureBenefits removes special benefits for specific subclass features
func removeSubclassFeatureBenefits(char *Character, featureName string, subclassName string) {
	source := BenefitSource{
		Type: "Subclass",
		Name: subclassName,
	}
	remover := NewBenefitRemover(char)

	debug.Log("removeSubclassFeatureBenefits: Removing benefits for '%s' from subclass '%s'", featureName, subclassName)

	switch featureName {
	case "Implements of Mercy":
		// Remove Medicine, Insight proficiency, and Herbalism Kit
		debug.Log("  Removing Implements of Mercy benefits")
		remover.RemoveAllBenefits(source.Type, source.Name)

	case "Shadow Arts":
		// Remove Darkvision increase
		debug.Log("  Removing Shadow Arts benefits")

		// Reduce darkvision by 60ft
		if char.Darkvision >= 60 {
			oldDarkvision := char.Darkvision
			char.Darkvision -= 60
			debug.Log("  Reduced Darkvision from %dft to %dft", oldDarkvision, char.Darkvision)
		}

		// Remove Darkness spell
		for i, spell := range char.SpellBook.Spells {
			if spell.Name == "Darkness" {
				char.SpellBook.Spells = append(char.SpellBook.Spells[:i], char.SpellBook.Spells[i+1:]...)
				debug.Log("  Removed Darkness spell")
				break
			}
		}

		// Remove Minor Illusion cantrip
		for i, cantripName := range char.SpellBook.Cantrips {
			if cantripName == "Minor Illusion" {
				char.SpellBook.Cantrips = append(char.SpellBook.Cantrips[:i], char.SpellBook.Cantrips[i+1:]...)
				if char.SpellBook.CantripsKnown > 0 {
					char.SpellBook.CantripsKnown--
				}
				debug.Log("  Removed Minor Illusion cantrip")
				break
			}
		}

		// Remove benefits from tracker
		remover.RemoveAllBenefits(source.Type, source.Name)

	case "Elemental Attunement":
		// Remove Elementalism cantrip
		debug.Log("  Removing Elemental Attunement benefits")

		for i, cantripName := range char.SpellBook.Cantrips {
			if cantripName == "Elementalism" {
				char.SpellBook.Cantrips = append(char.SpellBook.Cantrips[:i], char.SpellBook.Cantrips[i+1:]...)
				if char.SpellBook.CantripsKnown > 0 {
					char.SpellBook.CantripsKnown--
				}
				debug.Log("  Removed Elementalism cantrip")
				break
			}
		}

	case "Psionic Power":
		// Psi Warrior: Clear Psi Dice
		debug.Log("  Removing Psionic Power benefits - Clearing Psi Dice")
		char.PsiDice.Max = 0
		char.PsiDice.Current = 0
		char.PsiDice.Size = ""
		debug.Log("  Cleared Psi Dice")

	case "Combat Superiority":
		// Battle Master: Clear Superiority Dice and Maneuvers
		debug.Log("  Removing Combat Superiority benefits - Clearing Superiority Dice and Maneuvers")
		char.SuperiorityDice.Max = 0
		char.SuperiorityDice.Current = 0
		char.SuperiorityDice.Size = ""
		char.Maneuvers = []string{}
		debug.Log("  Cleared Superiority Dice and Maneuvers")

	case "Spellcasting":
		// Eldritch Knight: Remove spellcasting (if this is level 3 Fighter removal)
		if subclassName == "Eldritch Knight" {
			debug.Log("  Removing Eldritch Knight Spellcasting benefits")

			// Clear spell slots
			char.SpellBook.Slots.Level1.Maximum = 0
			char.SpellBook.Slots.Level1.Current = 0

			// Remove Eldritch Knight cantrips and spells
			// Note: This is a simplified approach. In a real implementation,
			// you might want to track which spells came from Eldritch Knight specifically
			char.SpellBook.Cantrips = []string{}
			char.SpellBook.CantripsKnown = 0
			char.SpellBook.Spells = []Spell{}

			debug.Log("  Cleared Eldritch Knight spellcasting")
		}

	case "Student of War":
		// Battle Master: Remove tool/skill proficiency
		// Note: This is handled by the benefit tracker in RemoveAllBenefits
		debug.Log("  Removing Student of War benefits")
		remover.RemoveAllBenefits(source.Type, source.Name)
	}
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

	// Update cantrips known from ALL spellcasting classes (multiclass totals)
	totalCantrips := GetMaxCantripsKnown(char.Classes)
	if totalCantrips > 0 {
		char.SpellBook.CantripsKnown = totalCantrips
		debug.Log("Total cantrips known (multiclass): %d", totalCantrips)
	}

	// Update spell slots based on MULTICLASS calculation
	// This automatically handles:
	// - Full casters: level 1:1 (Bard, Cleric, Druid, Sorcerer, Wizard)
	// - Half casters: level/2 (Paladin, Ranger)
	// - Third casters: level/3 (Eldritch Knight, Arcane Trickster)
	// - Excludes Warlock Pact Magic
	char.SpellBook.Slots = CalculateMulticlassSpellSlots(char.Classes)
	debug.Log("Multiclass spell slots calculated for total level %d", char.TotalLevel)

	// Warlock Pact Magic handled separately - doesn't combine!
	warlockLevel := char.GetClassLevel("Warlock")
	if warlockLevel > 0 {
		warlockSlots, warlockSlotLevel := GetWarlockPactSlots(warlockLevel)
		debug.Log("Warlock Pact Magic: %d slots of level %d", warlockSlots, warlockSlotLevel)

		// If character is ONLY a Warlock, use Pact Magic slots
		if len(char.Classes) == 1 && char.HasClass("Warlock") {
			// Pure Warlock: overwrite with pact slots
			char.SpellBook.Slots = SpellSlots{}
			switch warlockSlotLevel {
			case 1:
				char.SpellBook.Slots.Level1 = SpellSlot{Maximum: warlockSlots, Current: warlockSlots}
			case 2:
				char.SpellBook.Slots.Level2 = SpellSlot{Maximum: warlockSlots, Current: warlockSlots}
			case 3:
				char.SpellBook.Slots.Level3 = SpellSlot{Maximum: warlockSlots, Current: warlockSlots}
			case 4:
				char.SpellBook.Slots.Level4 = SpellSlot{Maximum: warlockSlots, Current: warlockSlots}
			case 5:
				char.SpellBook.Slots.Level5 = SpellSlot{Maximum: warlockSlots, Current: warlockSlots}
			}
			debug.Log("Pure Warlock: using Pact Magic slots only")
		} else if len(char.Classes) > 1 {
			// Multiclass with Warlock: Pact slots are SEPARATE from standard slots
			// TODO: Add PactMagicSlots field to SpellBook to track both
			debug.Log("Warlock multiclass: Pact Magic tracked separately (not yet fully implemented)")
		}
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

// DeLevel removes one level from the specified class
func DeLevel(char *Character, className string) (*DeLevelResult, error) {
	debug.Log("DeLevel: Removing level from %s", className)

	// Find the class
	classIndex := -1
	var classLevel *ClassLevel
	for i := range char.Classes {
		if char.Classes[i].ClassName == className {
			classIndex = i
			classLevel = &char.Classes[i]
			break
		}
	}

	if classIndex == -1 {
		return nil, fmt.Errorf("character does not have %s class", className)
	}

	if classLevel.Level < 1 {
		return nil, fmt.Errorf("class %s is already at level 0", className)
	}

	oldClassLevel := classLevel.Level
	oldTotalLevel := char.TotalLevel

	// Load class data
	classData := GetClassByName(className)
	if classData == nil {
		return nil, fmt.Errorf("failed to load class data for %s", className)
	}

	result := &DeLevelResult{
		ClassName:       className,
		OldClassLevel:   oldClassLevel,
		FeaturesRemoved: []string{},
	}

	// Remove features from this level
	featuresRemoved := removeFeaturesForLevel(char, classData, oldClassLevel)
	result.FeaturesRemoved = featuresRemoved

	// Remove subclass features if this was the subclass level
	if classLevel.Subclass != "" {
		subclassLevel := getSubclassLevelForClass(classData)
		if oldClassLevel == subclassLevel {
			subclassFeatures := removeSubclassFeaturesForLevel(char, className, classLevel.Subclass, oldClassLevel)
			result.FeaturesRemoved = append(result.FeaturesRemoved, subclassFeatures...)
			classLevel.Subclass = ""
		}
	}

	// Remove ASI choice if this level had one
	if CheckASIAvailable(className, oldClassLevel) {
		debug.Log("DeLevel: Removing ASI from level %d", oldClassLevel)
		if err := RemoveASIChoice(char, className, oldClassLevel); err != nil {
			debug.Log("DeLevel: Error removing ASI: %v", err)
		}
	}

	// Reduce HP (average for the class)
	conMod := char.AbilityScores.GetModifier("Constitution")
	avgHP := (classData.HitDie / 2) + 1 + conMod
	if avgHP < 1 {
		avgHP = 1
	}
	char.MaxHP -= avgHP
	if char.CurrentHP > char.MaxHP {
		char.CurrentHP = char.MaxHP
	}
	result.HPLost = avgHP

	// Decrease class level
	classLevel.Level--
	result.NewClassLevel = classLevel.Level

	// If class level is now 0, remove the class entirely
	if classLevel.Level == 0 {
		char.Classes = append(char.Classes[:classIndex], char.Classes[classIndex+1:]...)
		result.ClassRemoved = true
		debug.Log("DeLevel: Removed %s class entirely", className)
	}

	// Update total level
	char.TotalLevel = char.CalculateTotalLevel()
	char.Level = char.TotalLevel
	result.NewTotalLevel = char.TotalLevel

	// Update proficiency bonus
	char.ProficiencyBonus = CalculateProficiencyBonus(char.TotalLevel)

	// Ensure XP matches new level
	requiredXP := GetLevelXP(char.TotalLevel)
	if char.Experience > requiredXP {
		// Keep current XP if higher (they might be ready to level up again)
		debug.Log("DeLevel: Keeping current XP %d (required: %d)", char.Experience, requiredXP)
	}

	// Update derived stats
	char.UpdateDerivedStats()

	debug.Log("DeLevel complete: %s %d→%d, total level %d→%d",
		className, oldClassLevel, result.NewClassLevel, oldTotalLevel, result.NewTotalLevel)

	return result, nil
}

// removeFeaturesForLevel removes all features granted at a specific class level
func removeFeaturesForLevel(char *Character, class *Class, level int) []string {
	removed := []string{}

	// Read class JSON to get features for this level
	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(class.Name))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("removeFeaturesForLevel: Error reading class file: %v", err)
		return removed
	}

	var classWithProgression struct {
		LevelProgression []struct {
			Level    int                   `json:"level"`
			Features []FeatureDefinition `json:"features"`
		} `json:"level_progression"`
	}

	if err := json.Unmarshal(data, &classWithProgression); err != nil {
		debug.Log("removeFeaturesForLevel: Error unmarshaling: %v", err)
		return removed
	}

	// Find features for this level
	var featuresToRemove []string
	for _, levelData := range classWithProgression.LevelProgression {
		if levelData.Level == level {
			for _, featureDef := range levelData.Features {
				// Skip improvement features
				if featureDef.Name == "Ki Improvement" || featureDef.Name == "Focus Point Improvement" {
					continue
				}
				featuresToRemove = append(featuresToRemove, featureDef.Name)
			}
			break
		}
	}

	// Remove features from character
	for _, featureName := range featuresToRemove {
		for i := len(char.Features.Features) - 1; i >= 0; i-- {
			if char.Features.Features[i].Name == featureName {
				debug.Log("removeFeaturesForLevel: Removing feature %s", featureName)
				char.Features.Features = append(char.Features.Features[:i], char.Features.Features[i+1:]...)
				removed = append(removed, featureName)
				break // Remove only one instance
			}
		}
	}

	return removed
}

// removeSubclassFeaturesForLevel removes subclass features for a specific level
func removeSubclassFeaturesForLevel(char *Character, className, subclassName string, level int) []string {
	removed := []string{}

	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(className))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("removeSubclassFeaturesForLevel: Error reading class file: %v", err)
		return removed
	}

	var classWithSubclasses struct {
		Subclasses []struct {
			Name             string                              `json:"name"`
			FeaturesByLevel  map[string][]FeatureDefinition `json:"features_by_level"`
		} `json:"subclasses"`
	}

	if err := json.Unmarshal(data, &classWithSubclasses); err != nil {
		debug.Log("removeSubclassFeaturesForLevel: Error unmarshaling: %v", err)
		return removed
	}

	// Find the subclass
	for _, subclass := range classWithSubclasses.Subclasses {
		if subclass.Name == subclassName {
			levelKey := fmt.Sprintf("%d", level)
			if features, ok := subclass.FeaturesByLevel[levelKey]; ok {
				for _, featureDef := range features {
					// Remove benefits associated with this feature BEFORE removing the feature
					removeSubclassFeatureBenefits(char, featureDef.Name, subclassName)

					// Remove this feature
					for i := len(char.Features.Features) - 1; i >= 0; i-- {
						if char.Features.Features[i].Name == featureDef.Name {
							debug.Log("removeSubclassFeaturesForLevel: Removing subclass feature %s", featureDef.Name)
							char.Features.Features = append(char.Features.Features[:i], char.Features.Features[i+1:]...)
							removed = append(removed, featureDef.Name)
							break
						}
					}
				}
			}
			break
		}
	}

	return removed
}

// getSubclassLevelForClass returns the level at which a class gets its subclass
func getSubclassLevelForClass(class *Class) int {
	if len(class.Subclasses) > 0 {
		return class.Subclasses[0].SubclassLevel
	}
	return 3 // Default to 3 for most classes
}
