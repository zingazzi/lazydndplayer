// internal/models/character.go
package models

// Character represents a D&D 5e character
type Character struct {
	// Basic Info
	Name       string `json:"name"`
	Race       string `json:"race"`
	Class      string `json:"class"`
	Background string `json:"background"`
	Alignment  string `json:"alignment"`

	// Level & Experience
	Level      int `json:"level"`
	Experience int `json:"experience"`

	// Hit Points
	MaxHP     int `json:"max_hp"`
	CurrentHP int `json:"current_hp"`
	TempHP    int `json:"temp_hp"`

	// Armor Class & Speed
	ArmorClass int `json:"armor_class"`
	Speed      int `json:"speed"`

	// Core Stats
	AbilityScores AbilityScores `json:"ability_scores"`
	SavingThrows  SavingThrows  `json:"saving_throws"`

	// Skills & Proficiencies
	Skills            Skills `json:"skills"`
	ProficiencyBonus  int    `json:"proficiency_bonus"`

	// Combat & Features
	Initiative int        `json:"initiative"`
	Actions    ActionList `json:"actions"`

	// Equipment & Inventory
	Inventory Inventory `json:"inventory"`

	// Magic
	SpellBook SpellBook `json:"spellbook"`

	// Traits
	Languages           []string       `json:"languages"`
	Feats               []string       `json:"feats"`
	Resistances         []string       `json:"resistances"`
	Darkvision          int            `json:"darkvision"` // Range in feet, 0 if none
	SpeciesTraits       []SpeciesTrait `json:"species_traits"`
	SpeciesSkills       []SkillType    `json:"species_skills"` // Track which skills came from species
	SpeciesSpells       []string       `json:"species_spells"` // Track which spells came from species

	// Misc
	Notes string `json:"notes"`
}

// NewCharacter creates a new character with default values
func NewCharacter() *Character {
	char := &Character{
		Name:       "New Character",
		Race:       "Human",
		Class:      "Fighter",
		Background: "Folk Hero",
		Level:      1,
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
		Inventory: Inventory{
			Items:         []Item{},
			CarryCapacity: 150,
		},
		SpellBook: SpellBook{
			Spells: []Spell{},
			SpellcastingMod: Intelligence,
		},
		Languages:     []string{"Common"},
		Feats:         []string{},
		Resistances:   []string{},
		Darkvision:    0,
		SpeciesTraits: []SpeciesTrait{},
		SpeciesSkills: []SkillType{},
		SpeciesSpells: []string{},
	}
	char.UpdateDerivedStats()
	return char
}

// UpdateDerivedStats updates calculated values based on ability scores
func (c *Character) UpdateDerivedStats() {
	// Update initiative
	c.Initiative = c.AbilityScores.GetModifier(Dexterity)

	// Update carry capacity
	c.Inventory.CarryCapacity = CalculateCarryCapacity(c.AbilityScores.Strength)

	// Update proficiency bonus based on level
	c.ProficiencyBonus = CalculateProficiencyBonus(c.Level)

	// Update spell save DC and attack bonus if spellcaster
	if c.SpellBook.SpellcastingMod != "" {
		mod := c.AbilityScores.GetModifier(c.SpellBook.SpellcastingMod)
		c.SpellBook.SpellSaveDC = 8 + c.ProficiencyBonus + mod
		c.SpellBook.SpellAttackBonus = c.ProficiencyBonus + mod
	}
}

// CalculateProficiencyBonus returns proficiency bonus for a given level
func CalculateProficiencyBonus(level int) int {
	if level <= 4 {
		return 2
	} else if level <= 8 {
		return 3
	} else if level <= 12 {
		return 4
	} else if level <= 16 {
		return 5
	}
	return 6
}

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
}

// LongRest performs a long rest
func (c *Character) LongRest() {
	c.CurrentHP = c.MaxHP
	c.TempHP = 0
	c.Actions.LongRest()
	c.SpellBook.LongRest()
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
