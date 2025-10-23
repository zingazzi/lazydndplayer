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

	// Determine which ability and calculate new value
	const maxAbilityScore = 20
	var newValue int
	var normalizedAbility string

	switch {
	case strings.Contains(abilityLower, "strength"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Strength, increase, maxAbilityScore)
		ba.char.AbilityScores.Strength = newValue
		normalizedAbility = "Strength"
	case strings.Contains(abilityLower, "dexterity"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Dexterity, increase, maxAbilityScore)
		ba.char.AbilityScores.Dexterity = newValue
		normalizedAbility = "Dexterity"
	case strings.Contains(abilityLower, "constitution"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Constitution, increase, maxAbilityScore)
		ba.char.AbilityScores.Constitution = newValue
		normalizedAbility = "Constitution"
	case strings.Contains(abilityLower, "intelligence"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Intelligence, increase, maxAbilityScore)
		ba.char.AbilityScores.Intelligence = newValue
		normalizedAbility = "Intelligence"
	case strings.Contains(abilityLower, "wisdom"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Wisdom, increase, maxAbilityScore)
		ba.char.AbilityScores.Wisdom = newValue
		normalizedAbility = "Wisdom"
	case strings.Contains(abilityLower, "charisma"):
		newValue = CalculateAbilityIncrease(ba.char.AbilityScores.Charisma, increase, maxAbilityScore)
		ba.char.AbilityScores.Charisma = newValue
		normalizedAbility = "Charisma"
	default:
		return fmt.Errorf("unknown ability: %s", ability)
	}

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitAbilityScore,
		Target:      normalizedAbility,
		Value:       increase,
		Description: fmt.Sprintf("+%d %s", increase, normalizedAbility),
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
		Type:        BenefitFeature,
		Target:      feature.Name,
		Value:       feature.MaxUses,
		Description: fmt.Sprintf("Feature: %s (%d uses, %s)", feature.Name, feature.MaxUses, feature.RestType),
	})

	return nil
}

// AddInitiative adds an initiative bonus and tracks it
func (ba *BenefitApplier) AddInitiative(source BenefitSource, bonus int) error {
	ba.char.InitiativeBonus += bonus

	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitInitiative,
		Target:      "initiative",
		Value:       bonus,
		Description: fmt.Sprintf("+%d initiative", bonus),
	})

	return nil
}

// AddACBonus adds an AC bonus and tracks it
func (ba *BenefitApplier) AddACBonus(source BenefitSource, bonus int) error {
	ba.char.ACBonus += bonus

	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitAC,
		Target:      "ac",
		Value:       bonus,
		Description: fmt.Sprintf("+%d AC", bonus),
	})

	return nil
}

// AddPassiveBonus adds a bonus to passive skills and tracks it
func (ba *BenefitApplier) AddPassiveBonus(source BenefitSource, skillName string, bonus int) error {
	switch skillName {
	case "Perception":
		ba.char.PassivePerceptionBonus += bonus
	case "Investigation":
		ba.char.PassiveInvestigationBonus += bonus
	case "Insight":
		ba.char.PassiveInsightBonus += bonus
	}

	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitPassive,
		Target:      skillName,
		Value:       bonus,
		Description: fmt.Sprintf("+%d passive %s", bonus, skillName),
	})

	return nil
}

// AddToolProficiency adds a tool proficiency and tracks it
func (ba *BenefitApplier) AddToolProficiency(source BenefitSource, toolName string) error {
	// Check if already proficient
	for _, tool := range ba.char.ToolProficiencies {
		if tool == toolName {
			// Already have this proficiency, don't add duplicate
			return nil
		}
	}

	// Add tool proficiency
	ba.char.ToolProficiencies = append(ba.char.ToolProficiencies, toolName)

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitTool,
		Target:      toolName,
		Value:       1,
		Description: fmt.Sprintf("Tool proficiency: %s", toolName),
	})

	return nil
}

// AddItem adds an item to inventory and tracks it
func (ba *BenefitApplier) AddItem(source BenefitSource, itemName string, quantity int) error {
	// Check if item is gold
	if strings.Contains(strings.ToLower(itemName), " gp") || strings.Contains(strings.ToLower(itemName), "gold") {
		// Parse gold amount (e.g., "50 GP" -> 50)
		amount := 0
		fmt.Sscanf(itemName, "%d", &amount)
		if amount > 0 {
			ba.char.Inventory.Gold += amount

			// Track the benefit
			ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
				Source:      source,
				Type:        BenefitItem,
				Target:      itemName,
				Value:       amount,
				Description: fmt.Sprintf("%d gold", amount),
			})
			return nil
		}
	}

	// Try to find the item in the items database
	itemDef := GetItemDefinitionByName(itemName)
	if itemDef != nil {
		// Convert to inventory item
		item := ConvertToInventoryItem(*itemDef, quantity)
		ba.char.Inventory.AddItem(item)
	} else {
		// Item not found in database, add as generic item
		item := Item{
			Name:        itemName,
			Type:        Other,
			Quantity:    quantity,
			Weight:      0,
			Description: fmt.Sprintf("Item from %s", source.Name),
			Equipped:    false,
		}
		ba.char.Inventory.AddItem(item)
	}

	// Track the benefit
	ba.char.BenefitTracker.AddBenefit(GrantedBenefit{
		Source:      source,
		Type:        BenefitItem,
		Target:      itemName,
		Value:       quantity,
		Description: fmt.Sprintf("%dx %s", quantity, itemName),
	})

	return nil
}
