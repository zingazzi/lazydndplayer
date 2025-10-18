// internal/models/benefit_remover.go
package models

import (
	"fmt"
	"strings"
)

// BenefitRemover provides modular functions for removing benefits
type BenefitRemover struct {
	char *Character
}

// NewBenefitRemover creates a new benefit remover for a character
func NewBenefitRemover(char *Character) *BenefitRemover {
	// Ensure BenefitTracker is initialized
	if char.BenefitTracker == nil {
		char.BenefitTracker = NewBenefitTracker()
	}
	return &BenefitRemover{char: char}
}

// RemoveAllBenefits removes all benefits from a specific source
func (br *BenefitRemover) RemoveAllBenefits(sourceType, sourceName string) error {
	benefits := br.char.BenefitTracker.RemoveBenefitsBySource(sourceType, sourceName)

	// Remove each benefit
	for _, benefit := range benefits {
		switch benefit.Type {
		case BenefitAbilityScore:
			br.removeAbilityScore(benefit)
		case BenefitSkill:
			br.removeSkillProficiency(benefit)
		case BenefitLanguage:
			br.removeLanguage(benefit)
		case BenefitResistance:
			br.removeResistance(benefit)
		case BenefitSpeed:
			br.removeSpeed(benefit)
		case BenefitHP:
			br.removeHP(benefit)
		case BenefitSpell:
			br.removeSpell(benefit)
		case BenefitFeature:
			br.removeFeature(benefit)
		case BenefitInitiative:
			br.removeInitiative(benefit)
		case BenefitAC:
			br.removeACBonus(benefit)
		case BenefitPassive:
			br.removePassiveBonus(benefit)
		case BenefitTool:
			br.removeToolProficiency(benefit)
		}
	}

	br.char.UpdateDerivedStats()
	return nil
}

func (br *BenefitRemover) removeAbilityScore(benefit GrantedBenefit) {
	abilityLower := strings.ToLower(benefit.Target)

	switch {
	case strings.Contains(abilityLower, "strength"):
		br.char.AbilityScores.Strength = max(br.char.AbilityScores.Strength-benefit.Value, 1)
	case strings.Contains(abilityLower, "dexterity"):
		br.char.AbilityScores.Dexterity = max(br.char.AbilityScores.Dexterity-benefit.Value, 1)
	case strings.Contains(abilityLower, "constitution"):
		br.char.AbilityScores.Constitution = max(br.char.AbilityScores.Constitution-benefit.Value, 1)
	case strings.Contains(abilityLower, "intelligence"):
		br.char.AbilityScores.Intelligence = max(br.char.AbilityScores.Intelligence-benefit.Value, 1)
	case strings.Contains(abilityLower, "wisdom"):
		br.char.AbilityScores.Wisdom = max(br.char.AbilityScores.Wisdom-benefit.Value, 1)
	case strings.Contains(abilityLower, "charisma"):
		br.char.AbilityScores.Charisma = max(br.char.AbilityScores.Charisma-benefit.Value, 1)
	}
}

func (br *BenefitRemover) removeSkillProficiency(benefit GrantedBenefit) {
	skill := br.char.Skills.GetSkill(SkillType(benefit.Target))
	if skill != nil {
		// Restore to original level (stored in Value)
		skill.Proficiency = ProficiencyLevel(benefit.Value)
	}
}

func (br *BenefitRemover) removeLanguage(benefit GrantedBenefit) {
	// Check if any other source also grants this language
	remainingSources := 0
	for _, b := range br.char.BenefitTracker.Benefits {
		if b.Type == BenefitLanguage && strings.EqualFold(b.Target, benefit.Target) {
			remainingSources++
		}
	}

	// Only remove if no other source grants it
	if remainingSources == 0 {
		for i, lang := range br.char.Languages {
			if strings.EqualFold(lang, benefit.Target) {
				br.char.Languages = append(br.char.Languages[:i], br.char.Languages[i+1:]...)
				break
			}
		}
	}
}

func (br *BenefitRemover) removeResistance(benefit GrantedBenefit) {
	// Check if any other source also grants this resistance
	remainingSources := 0
	for _, b := range br.char.BenefitTracker.Benefits {
		if b.Type == BenefitResistance && strings.EqualFold(b.Target, benefit.Target) {
			remainingSources++
		}
	}

	// Only remove if no other source grants it
	if remainingSources == 0 {
		for i, res := range br.char.Resistances {
			if strings.EqualFold(res, benefit.Target) {
				br.char.Resistances = append(br.char.Resistances[:i], br.char.Resistances[i+1:]...)
				break
			}
		}
	}
}

func (br *BenefitRemover) removeSpeed(benefit GrantedBenefit) {
	br.char.Speed -= benefit.Value
	if br.char.Speed < 0 {
		br.char.Speed = 0
	}
}

func (br *BenefitRemover) removeHP(benefit GrantedBenefit) {
	br.char.MaxHP -= benefit.Value
	if br.char.MaxHP < 1 {
		br.char.MaxHP = 1
	}
	if br.char.CurrentHP > br.char.MaxHP {
		br.char.CurrentHP = br.char.MaxHP
	}
}

func (br *BenefitRemover) removeSpell(benefit GrantedBenefit) {
	// Remove from species spells tracking
	for i, spell := range br.char.SpeciesSpells {
		if strings.EqualFold(spell, benefit.Target) {
			br.char.SpeciesSpells = append(br.char.SpeciesSpells[:i], br.char.SpeciesSpells[i+1:]...)
			break
		}
	}

	// Remove from spellbook
	for i, spell := range br.char.SpellBook.Spells {
		if strings.EqualFold(spell.Name, benefit.Target) {
			br.char.SpellBook.Spells = append(br.char.SpellBook.Spells[:i], br.char.SpellBook.Spells[i+1:]...)
			break
		}
	}
}

func (br *BenefitRemover) removeFeature(benefit GrantedBenefit) {
	// Find and remove the feature by name and source
	sourceName := fmt.Sprintf("%s: %s", benefit.Source.Type, benefit.Source.Name)

	for i := len(br.char.Features.Features) - 1; i >= 0; i-- {
		feature := br.char.Features.Features[i]
		if feature.Name == benefit.Target && feature.Source == sourceName {
			br.char.Features.RemoveFeature(i)
			break
		}
	}
}

func (br *BenefitRemover) removeInitiative(benefit GrantedBenefit) {
	br.char.InitiativeBonus -= benefit.Value
	if br.char.InitiativeBonus < 0 {
		br.char.InitiativeBonus = 0
	}
}

func (br *BenefitRemover) removeACBonus(benefit GrantedBenefit) {
	br.char.ACBonus -= benefit.Value
	if br.char.ACBonus < 0 {
		br.char.ACBonus = 0
	}
}

func (br *BenefitRemover) removePassiveBonus(benefit GrantedBenefit) {
	switch benefit.Target {
	case "Perception":
		br.char.PassivePerceptionBonus -= benefit.Value
		if br.char.PassivePerceptionBonus < 0 {
			br.char.PassivePerceptionBonus = 0
		}
	case "Investigation":
		br.char.PassiveInvestigationBonus -= benefit.Value
		if br.char.PassiveInvestigationBonus < 0 {
			br.char.PassiveInvestigationBonus = 0
		}
	case "Insight":
		br.char.PassiveInsightBonus -= benefit.Value
		if br.char.PassiveInsightBonus < 0 {
			br.char.PassiveInsightBonus = 0
		}
	}
}

func (br *BenefitRemover) removeToolProficiency(benefit GrantedBenefit) {
	// Check if any other source also grants this tool proficiency
	remainingSources := 0
	for _, b := range br.char.BenefitTracker.Benefits {
		if b.Type == BenefitTool && strings.EqualFold(b.Target, benefit.Target) {
			remainingSources++
		}
	}

	// Only remove if no other source grants it
	if remainingSources == 0 {
		for i, tool := range br.char.ToolProficiencies {
			if strings.EqualFold(tool, benefit.Target) {
				br.char.ToolProficiencies = append(br.char.ToolProficiencies[:i], br.char.ToolProficiencies[i+1:]...)
				break
			}
		}
	}
}
