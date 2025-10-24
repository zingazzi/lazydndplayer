// internal/models/character.go
package models

import "strings"

// Character represents a D&D 5e character
type Character struct {
	// Basic Info
	Name       string `json:"name"`
	Race          string `json:"race"`
	Subtype       string `json:"subtype,omitempty"`        // For species with subtypes (Elf, Tiefling, Dragonborn)
	Class         string `json:"class"`                    // Backward compatibility - formatted as "Fighter 3 / Druid 2"
	Classes       []ClassLevel `json:"classes"`            // Multiclass support
	FightingStyle string `json:"fighting_style,omitempty"` // Fighting style for Fighter, Paladin, Ranger (deprecated - use ClassLevel.FightingStyle)
	Background    string `json:"background"`
	Origin        string `json:"origin"`    // Character origin (2024 rules)
	Alignment     string `json:"alignment"`

	// Level & Experience
	Level      int `json:"level"`      // Total character level (sum of all class levels)
	TotalLevel int `json:"total_level"` // Explicit total level field
	Experience int `json:"experience"`

	// Hit Points
	MaxHP           int `json:"max_hp"`
	CurrentHP       int `json:"current_hp"`
	TempHP          int `json:"temp_hp"`
	SpeciesHPBonus  int `json:"species_hp_bonus"` // HP bonus from species (e.g., Dwarven Toughness)

	// Armor Class & Speed
	ArmorClass int `json:"armor_class"`
	AC         int `json:"ac"` // Calculated AC (kept for compatibility)
	Speed      int `json:"speed"`

	// Core Stats
	AbilityScores AbilityScores `json:"ability_scores"`
	SavingThrows  SavingThrows  `json:"saving_throws"`

	// Skills & Proficiencies
	Skills                    Skills   `json:"skills"`
	ProficiencyBonus          int      `json:"proficiency_bonus"`
	ClassSkills               []SkillType `json:"class_skills"`             // Track which skills came from class
	ToolProficiencies         []string `json:"tool_proficiencies"`          // Tool proficiencies from origins/classes
	ArmorProficiencies        []string `json:"armor_proficiencies"`         // Armor proficiencies from class
	WeaponProficiencies       []string `json:"weapon_proficiencies"`        // Weapon proficiencies from class
	SavingThrowProficiencies  []string `json:"saving_throw_proficiencies"`  // Saving throw proficiencies from class
	MasteredWeapons           []string `json:"mastered_weapons"`            // Weapons with mastery properties
	MulticlassProficiencies   []string `json:"multiclass_proficiencies"`    // Track proficiencies from multiclassing

	// Combat & Features
	Initiative       int         `json:"initiative"`
	InitiativeBonus  int         `json:"initiative_bonus"`  // Bonus to initiative from feats
	ACBonus          int         `json:"ac_bonus"`          // Bonus to AC from feats
	PassivePerceptionBonus  int  `json:"passive_perception_bonus"`  // Bonus to passive Perception
	PassiveInvestigationBonus int `json:"passive_investigation_bonus"` // Bonus to passive Investigation
	PassiveInsightBonus      int  `json:"passive_insight_bonus"`       // Bonus to passive Insight
	Actions          ActionList  `json:"actions"`
	Features         FeatureList `json:"features"`

	// Equipment & Inventory
	Inventory Inventory `json:"inventory"`

	// Magic
	SpellBook SpellBook `json:"spellbook"`

	// Traits
	Languages           []string        `json:"languages"`
	Feats               []string        `json:"feats"`
	Resistances         []string        `json:"resistances"`
	BenefitTracker      *BenefitTracker `json:"benefit_tracker"` // Unified benefit tracking
	Darkvision          int            `json:"darkvision"` // Range in feet, 0 if none
	SpeciesTraits       []SpeciesTrait `json:"species_traits"`
	SpeciesSkills       []SkillType    `json:"species_skills"` // Track which skills came from species
	SpeciesSpells       []string       `json:"species_spells"` // Track which spells came from species

	// Choice Tracking (for easy rollback)
	Choices *CharacterChoices `json:"choices"` // Tracks all user choices for rollback

	// Inspiration
	Inspiration bool `json:"inspiration"` // Can be used to gain advantage on rolls

	// Misc
	Notes string `json:"notes"`
}

// SyncClassData ensures Classes array and Class string are in sync
func (c *Character) SyncClassData() {
	// If we have Classes array, sync Class string
	if len(c.Classes) > 0 {
		c.Class = c.GetClassDisplayString()
		c.TotalLevel = c.CalculateTotalLevel()
		c.Level = c.TotalLevel // Keep Level in sync for compatibility
	} else if c.Class != "" {
		// Backward compatibility: if only Class string exists, create Classes array
		// This handles old save files
		// For now, assume single class at current level
		c.Classes = []ClassLevel{
			{
				ClassName:     c.Class,
				Level:         c.Level,
				FightingStyle: c.FightingStyle,
			},
		}
		c.TotalLevel = c.Level
	}
}

// NewCharacter creates a new character with default values
func NewCharacter() *Character {
	char := &Character{
		Name:       "New Character",
		Race:       "Human",
		Class:      "",           // Empty - user must select a class
		Classes:    []ClassLevel{}, // Empty - populated when class is selected
		Background: "Folk Hero",
		Level:      1,            // Start at level 1
		TotalLevel: 1,            // Total level starts at 1
		MaxHP:      10,
		CurrentHP:  10,
		ArmorClass: 10,
		Speed:      30,
		AbilityScores: AbilityScores{
			Strength:     10,
			Dexterity:    10,
			Constitution: 10,
			Intelligence: 10,
			Wisdom:       10,
			Charisma:     10,
		},
		Skills:           NewDefaultSkills(),
		ProficiencyBonus: 2,
		Actions:          NewDefaultActions(),
		Features:         *NewFeatureList(),
		Inventory: Inventory{
			Items:         []Item{},
			CarryCapacity: 150,
		},
		SpellBook: SpellBook{
			Spells: []Spell{},
			SpellcastingMod: Intelligence,
		},
		Languages:         []string{"Common"},
		Feats:             []string{},
		Resistances:       []string{},
		ToolProficiencies:        []string{},
		ArmorProficiencies:       []string{},
		WeaponProficiencies:      []string{},
		SavingThrowProficiencies: []string{},
		MulticlassProficiencies:  []string{},
		BenefitTracker:    NewBenefitTracker(),
		Darkvision:        0,
		SpeciesTraits: []SpeciesTrait{},
		SpeciesSkills: []SkillType{},
		SpeciesSpells: []string{},
		ClassSkills:   []SkillType{},
		Choices:       NewCharacterChoices(),
	}
	char.UpdateDerivedStats()
	return char
}

// UpdateDerivedStats updates calculated values based on ability scores
func (c *Character) UpdateDerivedStats() {
	// Sync multiclass data first
	c.SyncClassData()

	// Update initiative (DEX modifier + bonus from feats)
	c.Initiative = c.AbilityScores.GetModifier(Dexterity) + c.InitiativeBonus

	// Update carry capacity
	c.Inventory.CarryCapacity = CalculateCarryCapacity(c.AbilityScores.Strength)

	// Update proficiency bonus based on level
	c.ProficiencyBonus = CalculateProficiencyBonus(c.Level)

	// Update spell save DC and attack bonus if spellcaster
	if c.SpellBook.SpellcastingMod != "" {
		mod := c.AbilityScores.GetModifier(c.SpellBook.SpellcastingMod)
		c.SpellBook.SpellSaveDC = 8 + c.ProficiencyBonus + mod
		c.SpellBook.SpellAttackBonus = c.ProficiencyBonus + mod

		// Update max prepared spells if prepared caster
		if c.SpellBook.IsPreparedCaster && c.SpellBook.PreparationFormula != "" {
			c.SpellBook.MaxPreparedSpells = c.CalculateMaxPreparedSpells(c.SpellBook.PreparationFormula)
		}
	}

	// Update AC based on equipped armor and shield
	c.AC = CalculateAC(c)
	c.ArmorClass = c.AC // Keep both for compatibility

	// Apply Monk speed bonus if applicable
	if c.IsMonk() && c.HasFeature("Unarmored Movement") {
		// Check if not wearing armor or shield
		hasArmorOrShield := false
		for _, item := range c.Inventory.Items {
			if item.Equipped && item.Type == Armor {
				hasArmorOrShield = true
				break
			}
		}

		if !hasArmorOrShield {
			monk := c.GetMonkMechanics()
			c.Speed += monk.GetUnarmoredMovementBonus()
		}
	}
}

// CalculateMaxPreparedSpells calculates the maximum number of spells that can be prepared
func (c *Character) CalculateMaxPreparedSpells(formula string) int {
	parts := strings.Split(strings.ToLower(formula), "+")
	total := 0

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "level":
			total += c.Level
		case "wisdom", "wis":
			total += c.AbilityScores.GetModifier("Wisdom")
		case "intelligence", "int":
			total += c.AbilityScores.GetModifier("Intelligence")
		case "charisma", "cha":
			total += c.AbilityScores.GetModifier("Charisma")
		}
	}

	// Minimum of 1
	if total < 1 {
		total = 1
	}

	return total
}

// CalculateProficiencyBonus returns proficiency bonus for a given level
// TakeDamage applies damage to the character
func (c *Character) TakeDamage(damage int) {
	// Apply to temp HP first
	if c.TempHP > 0 {
		if damage <= c.TempHP {
			c.TempHP -= damage
			return
		}
		damage -= c.TempHP
		c.TempHP = 0
	}

	// Apply remaining damage to current HP
	c.CurrentHP -= damage
	if c.CurrentHP < 0 {
		c.CurrentHP = 0
	}
}

// Heal restores HP
func (c *Character) Heal(amount int) {
	c.CurrentHP += amount
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
}

// ShortRest performs a short rest
func (c *Character) ShortRest() {
	c.Actions.ShortRest()
	c.Features.ShortRestRecover()
}

// LongRest performs a long rest
func (c *Character) LongRest() {
	c.CurrentHP = c.MaxHP
	c.TempHP = 0
	c.Actions.LongRest()
	c.SpellBook.LongRest()
	c.Features.LongRestRecover()

	// Humans regain Inspiration on long rest (Resourceful trait)
	if c.Race == "Human" {
		c.Inspiration = true
	}
}

// ExperienceThresholds returns XP needed for each level (1-20)
var ExperienceThresholds = []int{
	0,      // Level 1
	300,    // Level 2
	900,    // Level 3
	2700,   // Level 4
	6500,   // Level 5
	14000,  // Level 6
	23000,  // Level 7
	34000,  // Level 8
	48000,  // Level 9
	64000,  // Level 10
	85000,  // Level 11
	100000, // Level 12
	120000, // Level 13
	140000, // Level 14
	165000, // Level 15
	195000, // Level 16
	225000, // Level 17
	265000, // Level 18
	305000, // Level 19
	355000, // Level 20
}

// CanLevelUp checks if character has enough XP to level up
func (c *Character) CanLevelUp() bool {
	if c.Level >= 20 {
		return false
	}
	return c.Experience >= ExperienceThresholds[c.Level]
}

// GetNextLevelXP returns XP needed for next level
func (c *Character) GetNextLevelXP() int {
	if c.Level >= 20 {
		return ExperienceThresholds[19]
	}
	return ExperienceThresholds[c.Level]
}

// IsMonk checks if character has Monk levels
func (c *Character) IsMonk() bool {
	for _, classLevel := range c.Classes {
		if classLevel.ClassName == "Monk" {
			return true
		}
	}
	return false
}

// GetMonkMechanics returns Monk-specific mechanics handler
func (c *Character) GetMonkMechanics() *MonkMechanics {
	return NewMonkMechanics(c)
}

// HasFeature checks if character has a feature with the given name
func (c *Character) HasFeature(featureName string) bool {
	for _, feature := range c.Features.Features {
		if feature.Name == featureName {
			return true
		}
	}
	return false
}

// GetLevelXP returns the XP required to reach a given level
func GetLevelXP(level int) int {
	xpTable := map[int]int{
		1:  0,
		2:  300,
		3:  900,
		4:  2700,
		5:  6500,
		6:  14000,
		7:  23000,
		8:  34000,
		9:  48000,
		10: 64000,
		11: 85000,
		12: 100000,
		13: 120000,
		14: 140000,
		15: 165000,
		16: 195000,
		17: 225000,
		18: 265000,
		19: 305000,
		20: 355000,
	}
	if xp, ok := xpTable[level]; ok {
		return xp
	}
	return 355000 // Max level XP
}
