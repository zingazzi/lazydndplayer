// internal/models/companion.go
package models

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/debug"
)

// CompanionType represents the type of beast companion
type CompanionType string

const (
	CompanionTypeLand CompanionType = "Beast of Land"
	CompanionTypeSky  CompanionType = "Beast of Sky"
	CompanionTypeSea  CompanionType = "Beast of Sea"
)

// CompanionAbilityScores represents the ability scores of a companion
type CompanionAbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// GetModifier returns the ability modifier for a given ability
func (cas *CompanionAbilityScores) GetModifier(ability AbilityType) int {
	var score int
	switch ability {
	case Strength:
		score = cas.Strength
	case Dexterity:
		score = cas.Dexterity
	case Constitution:
		score = cas.Constitution
	case Intelligence:
		score = cas.Intelligence
	case Wisdom:
		score = cas.Wisdom
	case Charisma:
		score = cas.Charisma
	default:
		return 0
	}
	return CalculateModifier(score)
}

// Companion represents a ranger's beast companion
type Companion struct {
	Type           CompanionType          `json:"type"`
	CurrentHP      int                    `json:"current_hp"`
	MaxHP          int                    `json:"max_hp"`
	AC             int                    `json:"ac"`
	Speed          string                 `json:"speed"` // e.g., "40ft, climb 40ft"
	Senses         string                 `json:"senses"` // e.g., "Darkvision 60ft, Passive Perception 12"
	AbilityScores  CompanionAbilityScores `json:"ability_scores"`
	Traits         []string               `json:"traits"`
	SpecialNotes   string                 `json:"special_notes,omitempty"` // For Beast of Land charge note
}

// CompanionMechanics handles companion-related calculations and state
type CompanionMechanics struct {
	char *Character
}

// NewCompanionMechanics creates a new companion mechanics handler
func NewCompanionMechanics(char *Character) *CompanionMechanics {
	return &CompanionMechanics{char: char}
}

// GetCompanion returns the character's companion, or nil if none
func (cm *CompanionMechanics) GetCompanion() *Companion {
	if cm.char.Companion == nil {
		return nil
	}
	return cm.char.Companion
}

// InitializeCompanion initializes a companion based on type and ranger level
func InitializeCompanion(companionType CompanionType, rangerLevel int) *Companion {
	companion := &Companion{
		Type:  companionType,
		AC:    13, // Fixed AC for all types
		Traits: []string{"Primal Bond"}, // All companions have Primal Bond
	}

	switch companionType {
	case CompanionTypeLand:
		companion.AbilityScores = CompanionAbilityScores{
			Strength:     12,
			Dexterity:    14,
			Constitution: 15,
			Intelligence: 8,
			Wisdom:       14,
			Charisma:     11,
		}
		companion.MaxHP = 5 + (5 * rangerLevel)
		companion.CurrentHP = companion.MaxHP
		companion.Speed = "40ft, climb 40ft"
		companion.Senses = "Darkvision 60ft, Passive Perception 12"
		companion.SpecialNotes = "Beast Strike: If the companion moves at least 20ft before attacking, add 1d6 damage and the target is knocked prone."

	case CompanionTypeSky:
		companion.AbilityScores = CompanionAbilityScores{
			Strength:     6,
			Dexterity:    16,
			Constitution: 13,
			Intelligence: 8,
			Wisdom:       14,
			Charisma:     11,
		}
		companion.MaxHP = 4 + (5 * rangerLevel)
		companion.CurrentHP = companion.MaxHP
		companion.Speed = "10ft, fly 60ft"
		companion.Senses = "Darkvision 60ft, Passive Perception 12"
		companion.Traits = append(companion.Traits, "Fly By: Movement doesn't provoke opportunity attacks")

	case CompanionTypeSea:
		companion.AbilityScores = CompanionAbilityScores{
			Strength:     14,
			Dexterity:    14,
			Constitution: 15,
			Intelligence: 8,
			Wisdom:       14,
			Charisma:     11,
		}
		companion.MaxHP = 5 + (5 * rangerLevel)
		companion.CurrentHP = companion.MaxHP
		companion.Speed = "5ft, swim 60ft"
		companion.Senses = "Darkvision 90ft, Passive Perception 12"
		companion.Traits = append(companion.Traits, "Amphibious: Can breathe both in air and water")
		companion.SpecialNotes = "Beast Strike: Target is grappled on hit."
	}

	return companion
}

// CalculateCompanionAttackBonus calculates the attack bonus for companion's Beast Strike
func (cm *CompanionMechanics) CalculateCompanionAttackBonus() int {
	if cm.char.Companion == nil {
		return 0
	}

	// Get ranger level
	rangerLevel := cm.char.GetClassLevel("Ranger")
	if rangerLevel == 0 {
		return 0
	}

	// Attack bonus = Ranger spell attack modifier (proficiency + WIS mod)
	profBonus := CalculateProficiencyBonus(rangerLevel)
	wisMod := cm.char.AbilityScores.GetModifier(Wisdom)

	return profBonus + wisMod
}

// CalculateCompanionDamageBonus calculates the damage bonus for companion's Beast Strike
func (cm *CompanionMechanics) CalculateCompanionDamageBonus() int {
	if cm.char.Companion == nil {
		return 0
	}

	// Damage bonus = base bonus (varies by type) + ranger WIS modifier
	var baseBonus int
	wisMod := cm.char.AbilityScores.GetModifier(Wisdom)

	switch cm.char.Companion.Type {
	case CompanionTypeLand:
		// 1d8+2 + WIS mod
		baseBonus = 2
	case CompanionTypeSky:
		// 1d4+3 + WIS mod
		baseBonus = 3
	case CompanionTypeSea:
		// 1d6+2 + WIS mod
		baseBonus = 2
	default:
		baseBonus = 0
	}

	return baseBonus + wisMod
}

// GetCompanionBeastStrikeDamageDice returns the damage dice for Beast Strike
func (cm *CompanionMechanics) GetCompanionBeastStrikeDamageDice() string {
	if cm.char.Companion == nil {
		return "1d4"
	}

	switch cm.char.Companion.Type {
	case CompanionTypeLand:
		return "1d8"
	case CompanionTypeSky:
		return "1d4"
	case CompanionTypeSea:
		return "1d6"
	default:
		return "1d4"
	}
}

// GetCompanionAttack creates an Attack struct for the companion's Beast Strike
func (cm *CompanionMechanics) GetCompanionAttack() *Attack {
	if cm.char.Companion == nil {
		return nil
	}

	attackBonus := cm.CalculateCompanionAttackBonus()
	damageBonus := cm.CalculateCompanionDamageBonus()
	damageDice := cm.GetCompanionBeastStrikeDamageDice()

	attackName := fmt.Sprintf("%s - Beast Strike", cm.char.Companion.Type)

	attack := &Attack{
		Name:        attackName,
		AttackBonus: attackBonus,
		DamageDice:  damageDice,
		DamageBonus: damageBonus,
		DamageType:  "bludgeoning", // Default, could vary by type
		IsWeapon:    false,
		Range:       "5 ft.",
		Properties:  []string{},
	}

	// Add special notes as properties for display
	if cm.char.Companion.SpecialNotes != "" {
		attack.Properties = append(attack.Properties, cm.char.Companion.SpecialNotes)
	}

	return attack
}

// UpdateCompanionHP updates companion HP when ranger levels up
func (cm *CompanionMechanics) UpdateCompanionHP() {
	if cm.char.Companion == nil {
		return
	}

	rangerLevel := cm.char.GetClassLevel("Ranger")
	if rangerLevel == 0 {
		return
	}

	// Calculate new max HP
	var baseHP int
	switch cm.char.Companion.Type {
	case CompanionTypeLand, CompanionTypeSea:
		baseHP = 5
	case CompanionTypeSky:
		baseHP = 4
	default:
		baseHP = 5
	}

	newMaxHP := baseHP + (5 * rangerLevel)
	oldMaxHP := cm.char.Companion.MaxHP

	// Update max HP
	cm.char.Companion.MaxHP = newMaxHP

	// Adjust current HP proportionally if it was at max
	if cm.char.Companion.CurrentHP == oldMaxHP {
		cm.char.Companion.CurrentHP = newMaxHP
	} else if cm.char.Companion.CurrentHP > newMaxHP {
		// If current HP exceeds new max, cap it
		cm.char.Companion.CurrentHP = newMaxHP
	}

	debug.Log("Updated companion HP: %d/%d (was %d/%d)", cm.char.Companion.CurrentHP, cm.char.Companion.MaxHP, cm.char.Companion.CurrentHP, oldMaxHP)
}

// GetCompanionAbilityCheckBonus calculates ability check bonus with Primal Bond
func (cm *CompanionMechanics) GetCompanionAbilityCheckBonus(ability AbilityType) int {
	if cm.char.Companion == nil {
		return 0
	}

	// Base ability modifier
	abilityMod := cm.char.Companion.AbilityScores.GetModifier(ability)

	// Add ranger proficiency bonus (Primal Bond)
	rangerLevel := cm.char.GetClassLevel("Ranger")
	if rangerLevel > 0 {
		profBonus := CalculateProficiencyBonus(rangerLevel)
		return abilityMod + profBonus
	}

	return abilityMod
}

// GetCompanionSavingThrowBonus calculates saving throw bonus with Primal Bond
func (cm *CompanionMechanics) GetCompanionSavingThrowBonus(ability AbilityType) int {
	if cm.char.Companion == nil {
		return 0
	}

	// Base ability modifier
	abilityMod := cm.char.Companion.AbilityScores.GetModifier(ability)

	// Add ranger proficiency bonus (Primal Bond)
	rangerLevel := cm.char.GetClassLevel("Ranger")
	if rangerLevel > 0 {
		profBonus := CalculateProficiencyBonus(rangerLevel)
		return abilityMod + profBonus
	}

	return abilityMod
}

// HasCompanion returns true if the character has a companion
func (char *Character) HasCompanion() bool {
	return char.Companion != nil
}

// GetCompanionMechanics returns the companion mechanics handler
func (char *Character) GetCompanionMechanics() *CompanionMechanics {
	return NewCompanionMechanics(char)
}

// IsBeastMaster returns true if the character is a Beast Master ranger
func (char *Character) IsBeastMaster() bool {
	for _, classLevel := range char.Classes {
		if classLevel.ClassName == "Ranger" && classLevel.Subclass == "Beast Master" {
			return true
		}
	}
	return false
}
