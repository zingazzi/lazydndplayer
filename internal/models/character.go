// internal/models/character.go
package models

import (
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// HitDicePool tracks hit dice for a specific class
type HitDicePool struct {
	Count   int `json:"count"`    // Number of hit dice available for this class
	MaxDice int `json:"max_dice"` // Max hit dice for this class (equal to class level)
	DieSize int `json:"die_size"` // Size of the die (6, 8, 10, 12)
}

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
	Height        string `json:"height"`         // e.g., "6'2\""
	Weight        string `json:"weight"`         // e.g., "180 lbs"
	Personality   []string `json:"personality"`    // Personality traits (can have multiple)
	Ideal         []string `json:"ideal"`          // Character ideals (can have multiple)
	Bond          []string `json:"bond"`           // Character bonds (can have multiple)
	Flaw          []string `json:"flaw"`           // Character flaws (can have multiple)
	Backstory     string `json:"backstory"`      // Character backstory

	// Level & Experience
	Level      int `json:"level"`      // Total character level (sum of all class levels)
	TotalLevel int `json:"total_level"` // Explicit total level field
	Experience int `json:"experience"`

	// Hit Points
	MaxHP           int `json:"max_hp"`
	CurrentHP       int `json:"current_hp"`
	TempHP          int `json:"temp_hp"`
	SpeciesHPBonus  int `json:"species_hp_bonus"` // HP bonus from species (e.g., Dwarven Toughness)

	// Hit Dice (for resting)
	HitDice struct {
		Current int `json:"current"` // Current hit dice available
		Max     int `json:"max"`     // Max hit dice (equal to total level)
	} `json:"hit_dice"`
	HitDiceByClass map[string]HitDicePool `json:"hit_dice_by_class"` // Track hit dice per class

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

	// Fighter Subclass Resources
	PsiDice struct {
		Current int    `json:"current"`
		Max     int    `json:"max"`
		Size    string `json:"size"` // "d6", "d8", etc.
	} `json:"psi_dice,omitempty"`
	SuperiorityDice struct {
		Current int    `json:"current"`
		Max     int    `json:"max"`
		Size    string `json:"size"` // "d8", "d10", "d12"
	} `json:"superiority_dice,omitempty"`
	Maneuvers []string `json:"maneuvers,omitempty"` // List of known Battle Master maneuvers

	// Barbarian Subclass Resources
	WarriorDice struct {
		Current int    `json:"current"`
		Max     int    `json:"max"`
		Size    string `json:"size"` // Always "d12" for Zealot
	} `json:"warrior_dice,omitempty"` // Path of the Zealot

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
		Alignment:  "Neutral",
		Height:     "",
		Weight:     "",
		Personality: []string{},
		Ideal:      []string{},
		Bond:       []string{},
		Flaw:       []string{},
		Backstory:  "",
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

// GetMonkLevel returns the Monk level, or 0 if not a Monk
func (c *Character) GetMonkLevel() int {
	return c.GetClassLevel("Monk")
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

// GetFeature returns a feature by name, or nil if not found
func (c *Character) GetFeature(featureName string) *Feature {
	for i := range c.Features.Features {
		if c.Features.Features[i].Name == featureName {
			return &c.Features.Features[i]
		}
	}
	return nil
}

// IsFighter checks if character has Fighter levels
func (c *Character) IsFighter() bool {
	for _, classLevel := range c.Classes {
		if classLevel.ClassName == "Fighter" {
			return true
		}
	}
	return false
}

// GetFighterLevel returns the Fighter level, or 0 if not a Fighter
func (c *Character) GetFighterLevel() int {
	return c.GetClassLevel("Fighter")
}

// GetFighterSubclass returns the Fighter subclass name, or empty string if none
func (c *Character) GetFighterSubclass() string {
	for _, classLevel := range c.Classes {
		if classLevel.ClassName == "Fighter" {
			return classLevel.Subclass
		}
	}
	return ""
}

// IsChampion checks if character is a Champion Fighter
func (c *Character) IsChampion() bool {
	return c.GetFighterSubclass() == "Champion"
}

// IsEldritchKnight checks if character is an Eldritch Knight Fighter
func (c *Character) IsEldritchKnight() bool {
	return c.GetFighterSubclass() == "Eldritch Knight"
}

// IsPsiWarrior checks if character is a Psi Warrior Fighter
func (c *Character) IsPsiWarrior() bool {
	return c.GetFighterSubclass() == "Psi Warrior"
}

// IsBattleMaster checks if character is a Battle Master Fighter
func (c *Character) IsBattleMaster() bool {
	return c.GetFighterSubclass() == "Battle Master"
}

// GetBarbarianSubclass returns the character's Barbarian subclass
func (c *Character) GetBarbarianSubclass() string {
	for _, classLevel := range c.Classes {
		if classLevel.ClassName == "Barbarian" {
			return classLevel.Subclass
		}
	}
	return ""
}

// IsZealot checks if character is a Path of the Zealot Barbarian
func (c *Character) IsZealot() bool {
	return c.GetBarbarianSubclass() == "Path of the Zealot"
}

// IsWildHeart checks if character is a Path of the Wild Heart Barbarian
func (c *Character) IsWildHeart() bool {
	return c.GetBarbarianSubclass() == "Path of the Wild Heart"
}

// IsBerserker checks if character is a Path of the Berserker Barbarian
func (c *Character) IsBerserker() bool {
	return c.GetBarbarianSubclass() == "Path of the Berserker"
}

// IsWorldTree checks if character is a Path of the World Tree Barbarian
func (c *Character) IsWorldTree() bool {
	return c.GetBarbarianSubclass() == "Path of the World Tree"
}

// GetBarbarianLevel returns the character's Barbarian level
func (c *Character) GetBarbarianLevel() int {
	for _, classLevel := range c.Classes {
		if classLevel.ClassName == "Barbarian" {
			return classLevel.Level
		}
	}
	return 0
}

// HasImprovedCritical checks if character has Improved Critical feature (Champion)
func (c *Character) HasImprovedCritical() bool {
	if !c.IsChampion() {
		return false
	}
	return c.HasFeature("Improved Critical")
}

// GetCriticalRange returns the critical hit range (19 or 20, default 20)
func (c *Character) GetCriticalRange() int {
	if c.HasImprovedCritical() {
		feature := c.GetFeature("Improved Critical")
		if feature != nil && feature.Mechanics != nil {
			if critRange, ok := feature.Mechanics["critical_range"].(float64); ok {
				return int(critRange)
			}
		}
		return 19 // Default for Champion
	}
	return 20 // Normal critical range
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

// InitializeHitDice sets up hit dice tracking based on character's classes
// This should be called after loading class data
func (c *Character) InitializeHitDice(classes map[string]*Class) {
	if c.HitDiceByClass == nil {
		c.HitDiceByClass = make(map[string]HitDicePool)
	}

	totalDice := 0
	for _, classLevel := range c.Classes {
		pool := c.HitDiceByClass[classLevel.ClassName]
		pool.MaxDice = classLevel.Level

		// Get hit die size from class definition
		if class, ok := classes[classLevel.ClassName]; ok {
			pool.DieSize = class.HitDie
		}

		// If count is 0 or greater than max, reset to max (for new characters or after long rest)
		if pool.Count == 0 || pool.Count > pool.MaxDice {
			pool.Count = pool.MaxDice
		}

		c.HitDiceByClass[classLevel.ClassName] = pool
		totalDice += pool.Count
	}

	c.HitDice.Max = c.Level
	c.HitDice.Current = totalDice

	debug.Log("Initialized hit dice: Current=%d, Max=%d", c.HitDice.Current, c.HitDice.Max)
}

// GetTotalHitDice returns the total number of hit dice available
func (c *Character) GetTotalHitDice() int {
	total := 0
	for _, pool := range c.HitDiceByClass {
		total += pool.Count
	}
	return total
}

// SpendHitDice spends hit dice and restores HP (returns HP restored)
func (c *Character) SpendHitDice(count int, roller DiceRoller) int {
	if count <= 0 {
		return 0
	}

	totalHealing := 0
	diceSpent := 0
	conMod := c.AbilityScores.GetModifier("Constitution")

	debug.Log("=== SPENDING HIT DICE ===")
	debug.Log("  Requested: %d dice", count)
	debug.Log("  Available: %d dice", c.GetTotalHitDice())
	debug.Log("  CON Modifier: %d", conMod)

	// Spend hit dice from each class pool
	for className, pool := range c.HitDiceByClass {
		if diceSpent >= count {
			break
		}

		diceToSpend := count - diceSpent
		if diceToSpend > pool.Count {
			diceToSpend = pool.Count
		}

		for i := 0; i < diceToSpend; i++ {
			roll := roller.Roll(pool.DieSize)
			healing := roll + conMod
			if healing < 1 {
				healing = 1 // Minimum 1 HP per hit die
			}
			totalHealing += healing
			debug.Log("  Spent 1d%d from %s: rolled %d + %d CON = %d HP", pool.DieSize, className, roll, conMod, healing)
		}

		pool.Count -= diceToSpend
		c.HitDiceByClass[className] = pool
		diceSpent += diceToSpend
	}

	// Update total hit dice count
	c.HitDice.Current = c.GetTotalHitDice()

	// Restore HP
	c.CurrentHP += totalHealing
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}

	debug.Log("  Total Healing: %d HP", totalHealing)
	debug.Log("  HP: %d/%d", c.CurrentHP, c.MaxHP)
	debug.Log("  Remaining Hit Dice: %d/%d", c.HitDice.Current, c.HitDice.Max)

	return totalHealing
}

// PerformShortRest performs a short rest, restoring hit dice-based HP and short rest features
func (c *Character) PerformShortRest(diceSpent int, roller DiceRoller) int {
	debug.Log("=== PERFORMING SHORT REST ===")

	// Spend hit dice for healing
	healing := c.SpendHitDice(diceSpent, roller)

	// Restore all short rest features
	for i := range c.Features.Features {
		if c.Features.Features[i].RestType == ShortRest || c.Features.Features[i].RestType == Daily {
			c.Features.Features[i].CurrentUses = c.Features.Features[i].MaxUses
			debug.Log("  Restored feature: %s (%d uses)", c.Features.Features[i].Name, c.Features.Features[i].MaxUses)
		}
	}

	// Note: Actions don't track uses separately, they're tied to features

	// Restore Psi Dice for Psi Warrior (1 die on short rest)
	if c.PsiDice.Max > 0 && c.PsiDice.Current < c.PsiDice.Max {
		c.PsiDice.Current++
		debug.Log("  Restored 1 Psi Die: %d/%d", c.PsiDice.Current, c.PsiDice.Max)
	}

	// Restore Superiority Dice for Battle Master (all dice on short rest)
	if c.SuperiorityDice.Max > 0 {
		c.SuperiorityDice.Current = c.SuperiorityDice.Max
		debug.Log("  Restored Superiority Dice: %d/%d", c.SuperiorityDice.Current, c.SuperiorityDice.Max)
	}

	// Restore Warrior Dice for Zealot (1 die on short rest)
	if c.WarriorDice.Max > 0 && c.WarriorDice.Current < c.WarriorDice.Max {
		c.WarriorDice.Current++
		debug.Log("  Restored 1 Warrior Die: %d/%d", c.WarriorDice.Current, c.WarriorDice.Max)
	}

	debug.Log("=== SHORT REST COMPLETE ===")
	return healing
}

// PerformLongRest performs a long rest, restoring all HP, hit dice, and features
func (c *Character) PerformLongRest() {
	debug.Log("=== PERFORMING LONG REST ===")

	// Restore all HP
	c.CurrentHP = c.MaxHP
	debug.Log("  HP fully restored: %d/%d", c.CurrentHP, c.MaxHP)

	// Restore hit dice (regain half of max, minimum 1)
	hitDiceToRestore := c.HitDice.Max / 2
	if hitDiceToRestore < 1 {
		hitDiceToRestore = 1
	}

	debug.Log("  Restoring %d hit dice (half of %d)", hitDiceToRestore, c.HitDice.Max)

	// Distribute restored hit dice across classes
	for className, pool := range c.HitDiceByClass {
		if hitDiceToRestore <= 0 {
			break
		}

		missing := pool.MaxDice - pool.Count
		toRestore := missing
		if toRestore > hitDiceToRestore {
			toRestore = hitDiceToRestore
		}

		pool.Count += toRestore
		c.HitDiceByClass[className] = pool
		hitDiceToRestore -= toRestore
		debug.Log("  Restored %d d%d hit dice for %s", toRestore, pool.DieSize, className)
	}

	c.HitDice.Current = c.GetTotalHitDice()
	debug.Log("  Total hit dice: %d/%d", c.HitDice.Current, c.HitDice.Max)

	// Restore all features (both short and long rest)
	for i := range c.Features.Features {
		if c.Features.Features[i].RestType != None {
			c.Features.Features[i].CurrentUses = c.Features.Features[i].MaxUses
			debug.Log("  Restored feature: %s (%d uses)", c.Features.Features[i].Name, c.Features.Features[i].MaxUses)
		}
	}

	// Note: Actions don't track uses separately, they're tied to features

	// Restore all spell slots
	c.SpellBook.Slots.Level1.Current = c.SpellBook.Slots.Level1.Maximum
	c.SpellBook.Slots.Level2.Current = c.SpellBook.Slots.Level2.Maximum
	c.SpellBook.Slots.Level3.Current = c.SpellBook.Slots.Level3.Maximum
	c.SpellBook.Slots.Level4.Current = c.SpellBook.Slots.Level4.Maximum
	c.SpellBook.Slots.Level5.Current = c.SpellBook.Slots.Level5.Maximum
	c.SpellBook.Slots.Level6.Current = c.SpellBook.Slots.Level6.Maximum
	c.SpellBook.Slots.Level7.Current = c.SpellBook.Slots.Level7.Maximum
	c.SpellBook.Slots.Level8.Current = c.SpellBook.Slots.Level8.Maximum
	c.SpellBook.Slots.Level9.Current = c.SpellBook.Slots.Level9.Maximum
	debug.Log("  All spell slots restored")

	// Restore Psi Dice for Psi Warrior
	if c.PsiDice.Max > 0 {
		c.PsiDice.Current = c.PsiDice.Max
		debug.Log("  Restored Psi Dice: %d/%d", c.PsiDice.Current, c.PsiDice.Max)
	}

	// Restore Superiority Dice for Battle Master
	if c.SuperiorityDice.Max > 0 {
		c.SuperiorityDice.Current = c.SuperiorityDice.Max
		debug.Log("  Restored Superiority Dice: %d/%d", c.SuperiorityDice.Current, c.SuperiorityDice.Max)
	}

	// Restore Warrior Dice for Zealot
	if c.WarriorDice.Max > 0 {
		c.WarriorDice.Current = c.WarriorDice.Max
		debug.Log("  Restored Warrior Dice: %d/%d", c.WarriorDice.Current, c.WarriorDice.Max)
	}

	debug.Log("=== LONG REST COMPLETE ===")
}
