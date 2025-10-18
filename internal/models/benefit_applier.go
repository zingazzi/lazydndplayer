// internal/models/benefit_applier.go
package models

import (
	"fmt"
	"strings"
)

// BenefitApplier provides modular functions for applying benefits
type BenefitApplier struct {
	char *Character
}

// NewBenefitApplier creates a new benefit applier for a character
func NewBenefitApplier(char *Character) *BenefitApplier {
	// Ensure BenefitTracker is initialized
	if char.BenefitTracker == nil {
		char.BenefitTracker = NewBenefitTracker()
	}
	return &BenefitApplier{char: char}
}

// AddAbilityScore increases an ability score and tracks it
func (ba *BenefitApplier) AddAbilityScore(source BenefitSource, ability string, increase int) error {
	abilityLower := strings.ToLower(ability)

	// Apply the increase
	switch {
	case strings.Contains(abilityLower, "strength"):
		ba.char.AbilityScores.Strength = min(ba.char.AbilityScores.Strength+increase, 20)
		ability = "Strength"
	case strings.Contains(abilityLower, "dexterity"):
		ba.char.AbilityScores.Dexterity = min(ba.char.AbilityScores.Dexterity+increase, 20)
		ability = "Dexterity"
	case strings.Contains(abilityLower, "constitution"):
		ba.char.AbilityScores.Constitution = min(ba.char.AbilityScores.Constitution+increase, 20)
		ability = "Constitution"
	case strings.Contains(abilityLower, "intelligence"):
		ba.char.AbilityScores.Intelligence = min(ba.char.AbilityScores.Intelligence+increase, 20)
		ability = "Intelligence"
	case strings.Contains(abilityLower, "wisdom"):
		ba.char.AbilityScores.Wisdom = min(ba.char.AbilityScores.Wisdom+increase, 20)
		ability = "Wisdom"
	case strings.Contains(abilityLower, "charisma"):
		ba.char.AbilityScores.Charisma = min(ba.char.AbilityScores.Charisma+increase, 20)
		ability = "Charisma"
	default:
		return fmt.Errorf("unknown ability: %s", ability)
	}

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitAbilityScore,
		Target:      ability,
		Value:       increase,
		Description: fmt.Sprintf("+%d %s", increase, ability),
	})

	return nil
}

// AddSkillProficiency adds a skill proficiency and tracks it
func (ba *BenefitApplier) AddSkillProficiency(source BenefitSource, skillName string) error {
	skill := ba.char.Skills.GetSkill(SkillType(skillName))
	if skill == nil {
		return fmt.Errorf("unknown skill: %s", skillName)
	}

	// Store original level
	originalLevel := skill.Proficiency

	// Increase proficiency level
	if skill.Proficiency < Proficient {
		skill.Proficiency = Proficient
	} else if skill.Proficiency == Proficient {
		skill.Proficiency = Expertise
	}

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitSkill,
		Target:      skillName,
		Value:       int(originalLevel), // Store original level to revert correctly
		Description: fmt.Sprintf("%s proficiency", skillName),
	})

	return nil
}

// AddLanguage adds a language and tracks it
func (ba *BenefitApplier) AddLanguage(source BenefitSource, language string) error {
	// Check if already known
	for _, lang := range ba.char.Languages {
		if strings.EqualFold(lang, language) {
			// Already known, but still track it for this source
			ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
				Source:      source,
				Type:        BenefitLanguage,
				Target:      language,
				Value:       1,
				Description: fmt.Sprintf("%s language", language),
			})
			return nil
		}
	}

	// Add language
	ba.char.Languages = append(ba.char.Languages, language)

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitLanguage,
		Target:      language,
		Value:       1,
		Description: fmt.Sprintf("%s language", language),
	})

	return nil
}

// AddResistance adds a damage resistance and tracks it
func (ba *BenefitApplier) AddResistance(source BenefitSource, damageType string) error {
	// Check if already has resistance
	hasResistance := false
	for _, res := range ba.char.Resistances {
		if strings.EqualFold(res, damageType) {
			hasResistance = true
			break
		}
	}

	// Add resistance if not already present
	if !hasResistance {
		ba.char.Resistances = append(ba.char.Resistances, damageType)
	}

	// Track the benefit (even if already had it)
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitResistance,
		Target:      damageType,
		Value:       1,
		Description: fmt.Sprintf("%s resistance", damageType),
	})

	return nil
}

// AddSpeed increases speed and tracks it
func (ba *BenefitApplier) AddSpeed(source BenefitSource, increase int) error {
	ba.char.Speed += increase

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitSpeed,
		Target:      "Speed",
		Value:       increase,
		Description: fmt.Sprintf("+%d speed", increase),
	})

	return nil
}

// AddHP increases max HP and tracks it
func (ba *BenefitApplier) AddHP(source BenefitSource, increase int) error {
	ba.char.MaxHP += increase
	ba.char.CurrentHP = ba.char.MaxHP

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitHP,
		Target:      "Max HP",
		Value:       increase,
		Description: fmt.Sprintf("+%d HP", increase),
	})

	return nil
}

// AddSpell adds a spell and tracks it
func (ba *BenefitApplier) AddSpell(source BenefitSource, spellName string) error {
	// Add to species spells tracking
	ba.char.SpeciesSpells = append(ba.char.SpeciesSpells, spellName)

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitSpell,
		Target:      spellName,
		Value:       1,
		Description: fmt.Sprintf("%s spell", spellName),
	})

	return nil
}

// AddFeature adds a limited-use feature and tracks it
func (ba *BenefitApplier) AddFeature(source BenefitSource, featureDef FeatureDefinition) error {
	// Convert definition to actual feature
	feature := featureDef.ToFeature(ba.char, fmt.Sprintf("%s: %s", source.Type, source.Name))

	// Add to character's feature list
	ba.char.Features.AddFeature(feature)

	// Track in benefit system (use feature name as target)
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        "feature", // New benefit type for features
		Target:      feature.Name,
		Value:       feature.MaxUses,
		Description: fmt.Sprintf("Feature: %s (%d uses, %s)", feature.Name, feature.MaxUses, feature.RestType),
	})

	return nil
}
