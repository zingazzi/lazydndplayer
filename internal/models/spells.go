// internal/models/spells.go
package models

import (
	"fmt"
	"strings"
)

// SpellSchool represents the school of magic
type SpellSchool string

const (
	Abjuration    SpellSchool = "Abjuration"
	Conjuration   SpellSchool = "Conjuration"
	Divination    SpellSchool = "Divination"
	Enchantment   SpellSchool = "Enchantment"
	Evocation     SpellSchool = "Evocation"
	Illusion      SpellSchool = "Illusion"
	Necromancy    SpellSchool = "Necromancy"
	Transmutation SpellSchool = "Transmutation"
)

// Spell represents a single spell
type Spell struct {
	Name           string      `json:"name"`
	Level          int         `json:"level"` // 0 for cantrips
	School         SpellSchool `json:"school"`
	CastingTime    string      `json:"casting_time"`
	ActionType     string      `json:"actionType"`     // action, bonus action, reaction
	Range          string      `json:"range"`
	Components     interface{} `json:"components"`     // Can be string or []string
	Material       string      `json:"material"`       // Material component description
	Duration       string      `json:"duration"`
	Concentration  bool        `json:"concentration"`
	Description    string      `json:"description"`
	CantripUpgrade string      `json:"cantripUpgrade"` // Cantrip scaling info
	Prepared       bool        `json:"prepared"`       // For prepared casters
	Known          bool        `json:"known"`          // For known casters
	Ritual         bool        `json:"ritual"`
	Classes        []string    `json:"classes"` // Classes that can learn this spell
}

// GetComponentsString returns components as a formatted string
func (s *Spell) GetComponentsString() string {
	switch v := s.Components.(type) {
	case string:
		return v
	case []interface{}:
		parts := make([]string, len(v))
		for i, comp := range v {
			parts[i] = strings.ToUpper(fmt.Sprint(comp))
		}
		return strings.Join(parts, ", ")
	case []string:
		parts := make([]string, len(v))
		for i, comp := range v {
			parts[i] = strings.ToUpper(comp)
		}
		return strings.Join(parts, ", ")
	default:
		return "V, S"
	}
}

// SpellSlots represents available spell slots per level
type SpellSlots struct {
	Level1 SpellSlot `json:"level_1"`
	Level2 SpellSlot `json:"level_2"`
	Level3 SpellSlot `json:"level_3"`
	Level4 SpellSlot `json:"level_4"`
	Level5 SpellSlot `json:"level_5"`
	Level6 SpellSlot `json:"level_6"`
	Level7 SpellSlot `json:"level_7"`
	Level8 SpellSlot `json:"level_8"`
	Level9 SpellSlot `json:"level_9"`
}

// SpellSlot represents slots for a specific spell level
type SpellSlot struct {
	Maximum int `json:"maximum"`
	Current int `json:"current"` // Remaining slots
}

// UseSlot uses one spell slot if available
func (s *SpellSlot) UseSlot() bool {
	if s.Current > 0 {
		s.Current--
		return true
	}
	return false
}

// RestoreSlot restores one spell slot
func (s *SpellSlot) RestoreSlot() {
	if s.Current < s.Maximum {
		s.Current++
	}
}

// RestoreAll restores all spell slots (long rest)
func (s *SpellSlot) RestoreAll() {
	s.Current = s.Maximum
}

// SpellBook holds all character spells and slots
type SpellBook struct {
	Spells            []Spell     `json:"spells"`
	Slots             SpellSlots  `json:"slots"`
	SpellcastingMod   AbilityType `json:"spellcasting_mod"`    // INT, WIS, or CHA
	SpellSaveDC       int         `json:"spell_save_dc"`
	SpellAttackBonus  int         `json:"spell_attack_bonus"`
	IsPreparedCaster  bool        `json:"is_prepared_caster"`  // true for Druid, Cleric, Paladin, Wizard
	MaxPreparedSpells int         `json:"max_prepared_spells"` // Number of spells that can be prepared
	PreparationFormula string     `json:"preparation_formula"` // Formula for calculating max prepared spells (e.g., "wisdom+level")
	CantripsKnown     int         `json:"cantrips_known"`      // Number of cantrips known (for level)
	Cantrips          []string    `json:"cantrips"`            // List of selected cantrip names
}

// GetSlotByLevel returns a pointer to the spell slot for a given level
func (sb *SpellBook) GetSlotByLevel(level int) *SpellSlot {
	switch level {
	case 1:
		return &sb.Slots.Level1
	case 2:
		return &sb.Slots.Level2
	case 3:
		return &sb.Slots.Level3
	case 4:
		return &sb.Slots.Level4
	case 5:
		return &sb.Slots.Level5
	case 6:
		return &sb.Slots.Level6
	case 7:
		return &sb.Slots.Level7
	case 8:
		return &sb.Slots.Level8
	case 9:
		return &sb.Slots.Level9
	default:
		return nil
	}
}

// LongRest restores all spell slots
func (sb *SpellBook) LongRest() {
	sb.Slots.Level1.RestoreAll()
	sb.Slots.Level2.RestoreAll()
	sb.Slots.Level3.RestoreAll()
	sb.Slots.Level4.RestoreAll()
	sb.Slots.Level5.RestoreAll()
	sb.Slots.Level6.RestoreAll()
	sb.Slots.Level7.RestoreAll()
	sb.Slots.Level8.RestoreAll()
	sb.Slots.Level9.RestoreAll()
}

// AddSpell adds a spell to the spellbook
func (sb *SpellBook) AddSpell(spell Spell) {
	sb.Spells = append(sb.Spells, spell)
}

// RemoveSpell removes a spell from the spellbook
func (sb *SpellBook) RemoveSpell(name string) bool {
	for i, spell := range sb.Spells {
		if spell.Name == name {
			sb.Spells = append(sb.Spells[:i], sb.Spells[i+1:]...)
			return true
		}
	}
	return false
}

// GetPreparedSpells returns all prepared spells
func (sb *SpellBook) GetPreparedSpells() []Spell {
	prepared := []Spell{}
	for _, spell := range sb.Spells {
		if spell.Prepared || spell.Level == 0 { // Cantrips are always prepared
			prepared = append(prepared, spell)
		}
	}
	return prepared
}
