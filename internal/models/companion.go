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

// CompanionCategory distinguishes between regular companions and Beast Master special beasts
type CompanionCategory string

const (
	RegularCompanion   CompanionCategory = "Regular"        // Regular beast from Monster Manual
	BeastMasterSpecial CompanionCategory = "BeastMaster"    // Beast Master special beast (Beast of Land/Sky/Sea)
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

// Companion represents a character's companion (regular beast or Beast Master special beast)
type Companion struct {
	Category       CompanionCategory      `json:"category"`        // RegularCompanion or BeastMasterSpecial
	Type           CompanionType          `json:"type,omitempty"`  // For Beast Master special beasts (Beast of Land/Sky/Sea)
	BeastName      string                 `json:"beast_name,omitempty"` // For regular companions (e.g., "Mastiff", "Mule", "Black Bear")
	Name           string                 `json:"name,omitempty"`  // Custom name given by player
	CurrentHP      int                    `json:"current_hp"`
	MaxHP          int                    `json:"max_hp"`
	AC             int                    `json:"ac"`
	Speed          string                 `json:"speed"` // e.g., "40ft, climb 40ft"
	Senses         string                 `json:"senses"` // e.g., "Darkvision 60ft, Passive Perception 12"
	AbilityScores  CompanionAbilityScores `json:"ability_scores"`
	Traits         []string               `json:"traits"`
	SpecialNotes   string                 `json:"special_notes,omitempty"` // For Beast of Land charge note
	Actions        []BeastAction          `json:"actions,omitempty"`       // For regular companions
}

// CompanionMechanics handles companion-related calculations and state
type CompanionMechanics struct {
	char *Character
}

// NewCompanionMechanics creates a new companion mechanics handler
func NewCompanionMechanics(char *Character) *CompanionMechanics {
	return &CompanionMechanics{char: char}
}

// GetCompanion returns the character's first companion, or nil if none (for backward compatibility)
func (cm *CompanionMechanics) GetCompanion() *Companion {
	if len(cm.char.Companions) > 0 {
		return &cm.char.Companions[0]
	}
	// Fallback to old field for backward compatibility
	if cm.char.Companion != nil {
		return cm.char.Companion
	}
	return nil
}

// GetCompanionByIndex returns the companion at the given index
func (cm *CompanionMechanics) GetCompanionByIndex(index int) *Companion {
	return cm.char.GetCompanion(index)
}

// InitializeCompanion initializes a Beast Master special companion based on type and ranger level
func InitializeCompanion(companionType CompanionType, rangerLevel int) *Companion {
	companion := &Companion{
		Category: BeastMasterSpecial,
		Type:     companionType,
		AC:       13, // Fixed AC for all types
		Traits:   []string{"Primal Bond"}, // All companions have Primal Bond
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

// InitializeRegularCompanion initializes a regular companion from a beast definition
func InitializeRegularCompanion(beastName string) *Companion {
	beastDef := GetBeastDefinition(beastName)
	if beastDef == nil {
		return nil
	}

	companion := &Companion{
		Category:      RegularCompanion,
		BeastName:     beastName,
		AC:            beastDef.AC,
		MaxHP:         beastDef.HP,
		CurrentHP:     beastDef.HP,
		Speed:         beastDef.Speed,
		Senses:        beastDef.Senses,
		AbilityScores: beastDef.AbilityScores,
		Traits:        beastDef.Traits,
		Actions:       beastDef.Actions,
	}

	return companion
}

// GetDisplayName returns the display name for the companion (custom name, or beast name, or type)
func (c *Companion) GetDisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.BeastName != "" {
		return c.BeastName
	}
	return string(c.Type)
}

// IsBeastMasterSpecial returns true if this is a Beast Master special beast
func (c *Companion) IsBeastMasterSpecial() bool {
	return c.Category == BeastMasterSpecial
}

// IsRegularCompanion returns true if this is a regular companion
func (c *Companion) IsRegularCompanion() bool {
	return c.Category == RegularCompanion
}

// CalculateCompanionAttackBonus calculates the attack bonus for a companion's attack
// For Beast Master special beasts: scales with ranger level
// For regular companions: uses fixed attack bonus from beast definition
func (cm *CompanionMechanics) CalculateCompanionAttackBonus(companion *Companion) int {
	if companion == nil {
		return 0
	}

	// Regular companions use fixed attack bonus from their actions
	if companion.IsRegularCompanion() {
		if len(companion.Actions) > 0 {
			return companion.Actions[0].AttackBonus
		}
		return 0
	}

	// Beast Master special beasts scale with ranger level
	if companion.IsBeastMasterSpecial() {
		rangerLevel := cm.char.GetClassLevel("Ranger")
		if rangerLevel == 0 {
			return 0
		}

		// Attack bonus = Ranger spell attack modifier (proficiency + WIS mod)
		profBonus := CalculateProficiencyBonus(rangerLevel)
		wisMod := cm.char.AbilityScores.GetModifier(Wisdom)

		return profBonus + wisMod
	}

	return 0
}

// CalculateCompanionDamageBonus calculates the damage bonus for a companion's attack
// For Beast Master special beasts: scales with ranger level
// For regular companions: uses fixed damage bonus from beast definition
func (cm *CompanionMechanics) CalculateCompanionDamageBonus(companion *Companion) int {
	if companion == nil {
		return 0
	}

	// Regular companions use fixed damage bonus from their actions
	if companion.IsRegularCompanion() {
		if len(companion.Actions) > 0 {
			return companion.Actions[0].DamageBonus
		}
		return 0
	}

	// Beast Master special beasts scale with ranger level
	if companion.IsBeastMasterSpecial() {
		// Damage bonus = base bonus (varies by type) + ranger WIS modifier
		var baseBonus int
		wisMod := cm.char.AbilityScores.GetModifier(Wisdom)

		switch companion.Type {
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

	return 0
}

// GetCompanionDamageDice returns the damage dice for a companion's attack
// For Beast Master special beasts: returns Beast Strike dice
// For regular companions: returns damage dice from beast definition
func (cm *CompanionMechanics) GetCompanionDamageDice(companion *Companion) string {
	if companion == nil {
		return "1d4"
	}

	// Regular companions use fixed damage dice from their actions
	if companion.IsRegularCompanion() {
		if len(companion.Actions) > 0 {
			return companion.Actions[0].DamageDice
		}
		return "1d4"
	}

	// Beast Master special beasts use Beast Strike dice
	if companion.IsBeastMasterSpecial() {
		switch companion.Type {
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

	return "1d4"
}

// GetCompanionAttack creates an Attack struct for a companion's attack
// For Beast Master special beasts: creates Beast Strike attack
// For regular companions: creates attack from beast definition
func (cm *CompanionMechanics) GetCompanionAttack(companion *Companion) *Attack {
	if companion == nil {
		return nil
	}

	// Regular companions use their defined actions
	if companion.IsRegularCompanion() {
		if len(companion.Actions) == 0 {
			return nil
		}
		beastAction := companion.Actions[0] // Use first action

		attack := &Attack{
			Name:        beastAction.Name,
			AttackBonus: beastAction.AttackBonus,
			DamageDice:  beastAction.DamageDice,
			DamageBonus: beastAction.DamageBonus,
			DamageType:  beastAction.DamageType,
			IsWeapon:    false,
			Range:       beastAction.Range,
			Properties:  []string{},
		}

		if beastAction.Description != "" {
			attack.Properties = append(attack.Properties, beastAction.Description)
		}

		return attack
	}

	// Beast Master special beasts use Beast Strike
	if companion.IsBeastMasterSpecial() {
		attackBonus := cm.CalculateCompanionAttackBonus(companion)
		damageBonus := cm.CalculateCompanionDamageBonus(companion)
		damageDice := cm.GetCompanionDamageDice(companion)

		attackName := fmt.Sprintf("%s - Beast Strike", companion.Type)

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
		if companion.SpecialNotes != "" {
			attack.Properties = append(attack.Properties, companion.SpecialNotes)
		}

		return attack
	}

	return nil
}

// UpdateCompanionHP updates Beast Master special beast HP when ranger levels up
// Regular companions don't scale, so this only affects Beast Master special beasts
func (cm *CompanionMechanics) UpdateCompanionHP() {
	rangerLevel := cm.char.GetClassLevel("Ranger")
	if rangerLevel == 0 {
		return
	}

	// Update all Beast Master special beasts
	for i := range cm.char.Companions {
		companion := &cm.char.Companions[i]
		if !companion.IsBeastMasterSpecial() {
			continue // Regular companions don't scale
		}

		// Calculate new max HP
		var baseHP int
		switch companion.Type {
		case CompanionTypeLand, CompanionTypeSea:
			baseHP = 5
		case CompanionTypeSky:
			baseHP = 4
		default:
			baseHP = 5
		}

		newMaxHP := baseHP + (5 * rangerLevel)
		oldMaxHP := companion.MaxHP

		// Update max HP
		companion.MaxHP = newMaxHP

		// Adjust current HP proportionally if it was at max
		if companion.CurrentHP == oldMaxHP {
			companion.CurrentHP = newMaxHP
		} else if companion.CurrentHP > newMaxHP {
			// If current HP exceeds new max, cap it
			companion.CurrentHP = newMaxHP
		}

		debug.Log("Updated companion HP: %d/%d (was %d/%d)", companion.CurrentHP, companion.MaxHP, companion.CurrentHP, oldMaxHP)
	}
}

// GetCompanionAbilityCheckBonus calculates ability check bonus for a companion
// Beast Master special beasts get Primal Bond bonus, regular companions don't
func (cm *CompanionMechanics) GetCompanionAbilityCheckBonus(companion *Companion, ability AbilityType) int {
	if companion == nil {
		return 0
	}

	// Base ability modifier
	abilityMod := companion.AbilityScores.GetModifier(ability)

	// Beast Master special beasts get Primal Bond (ranger proficiency bonus)
	if companion.IsBeastMasterSpecial() {
		rangerLevel := cm.char.GetClassLevel("Ranger")
		if rangerLevel > 0 {
			profBonus := CalculateProficiencyBonus(rangerLevel)
			return abilityMod + profBonus
		}
	}

	return abilityMod
}

// GetCompanionSavingThrowBonus calculates saving throw bonus for a companion
// Beast Master special beasts get Primal Bond bonus, regular companions don't
func (cm *CompanionMechanics) GetCompanionSavingThrowBonus(companion *Companion, ability AbilityType) int {
	if companion == nil {
		return 0
	}

	// Base ability modifier
	abilityMod := companion.AbilityScores.GetModifier(ability)

	// Beast Master special beasts get Primal Bond (ranger proficiency bonus)
	if companion.IsBeastMasterSpecial() {
		rangerLevel := cm.char.GetClassLevel("Ranger")
		if rangerLevel > 0 {
			profBonus := CalculateProficiencyBonus(rangerLevel)
			return abilityMod + profBonus
		}
	}

	return abilityMod
}

// HasCompanion returns true if the character has at least one companion
func (char *Character) HasCompanion() bool {
	return len(char.Companions) > 0 || char.Companion != nil
}

// HasCompanions returns true if the character has at least one companion
func (char *Character) HasCompanions() bool {
	return len(char.Companions) > 0
}

// GetCompanion returns the companion at the given index, or nil if invalid
func (char *Character) GetCompanion(index int) *Companion {
	if index < 0 || index >= len(char.Companions) {
		return nil
	}
	return &char.Companions[index]
}

// AddCompanion adds a companion to the character's companions list
func (char *Character) AddCompanion(companion *Companion) {
	if companion == nil {
		return
	}
	if char.Companions == nil {
		char.Companions = []Companion{}
	}
	char.Companions = append(char.Companions, *companion)
}

// RemoveCompanion removes a companion at the given index
func (char *Character) RemoveCompanion(index int) bool {
	if index < 0 || index >= len(char.Companions) {
		return false
	}
	char.Companions = append(char.Companions[:index], char.Companions[index+1:]...)
	return true
}

// GetBeastMasterCompanion returns the Beast Master special beast companion, if any
func (char *Character) GetBeastMasterCompanion() *Companion {
	for i := range char.Companions {
		if char.Companions[i].IsBeastMasterSpecial() {
			return &char.Companions[i]
		}
	}
	return nil
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
