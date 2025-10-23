// internal/models/calculations.go
package models

// This file contains pure calculation functions that don't mutate state.
// These functions can be easily tested and reused.

// CalculateAbilityIncrease calculates a new ability score after an increase.
// Ability scores are capped at maxScore (typically 20).
func CalculateAbilityIncrease(current, increase, maxScore int) int {
	newValue := current + increase
	if newValue > maxScore {
		return maxScore
	}
	return newValue
}

// CalculateProficiencyBonus calculates proficiency bonus based on total character level.
func CalculateProficiencyBonus(level int) int {
	if level <= 0 {
		return 2
	}
	return 2 + (level-1)/4
}

// CalculateAbilityModifier calculates the ability modifier for a given score.
func CalculateAbilityModifier(score int) int {
	return (score - 10) / 2
}

// CalculatePassiveScore calculates a passive skill score.
// passiveScore = 10 + abilityModifier + proficiencyBonus (if proficient) + otherBonuses
func CalculatePassiveScore(abilityModifier int, isProficient bool, proficiencyBonus int, otherBonuses int) int {
	base := 10 + abilityModifier
	if isProficient {
		base += proficiencyBonus
	}
	return base + otherBonuses
}

// CalculateACWithArmor calculates AC based on armor type, dexterity, and other bonuses.
// armorBase: Base AC from armor (e.g., 11 for Leather, 16 for Chainmail)
// dexMod: Dexterity modifier
// maxDexBonus: Maximum dex bonus allowed by armor (-1 for unlimited)
// otherBonuses: Additional AC bonuses (shield, fighting style, etc.)
func CalculateACWithArmor(armorBase, dexMod, maxDexBonus, otherBonuses int) int {
	ac := armorBase

	if maxDexBonus == -1 {
		// No limit on dex bonus (light armor or no armor)
		ac += dexMod
	} else if maxDexBonus > 0 {
		// Limited dex bonus (medium armor)
		if dexMod > maxDexBonus {
			ac += maxDexBonus
		} else {
			ac += dexMod
		}
	}
	// Heavy armor: no dex bonus (maxDexBonus == 0)

	return ac + otherBonuses
}

// CalculateUnarmoredAC calculates AC when not wearing armor.
// base: Usually 10
// dexMod: Dexterity modifier (always applies)
// otherMods: Other ability modifiers (e.g., Wisdom for Monk, Constitution for Barbarian)
// bonuses: Shields, fighting styles, etc.
func CalculateUnarmoredAC(base, dexMod int, otherMods []int, bonuses int) int {
	ac := base + dexMod
	for _, mod := range otherMods {
		ac += mod
	}
	return ac + bonuses
}

// CalculateCarryCapacity calculates how much a character can carry.
// Standard rule: Strength score × 15
func CalculateCarryCapacity(strengthScore int) float64 {
	return float64(strengthScore * 15)
}

// CalculateInitiativeModifier calculates the initiative bonus.
// Typically dexterity modifier + any bonuses (e.g., Alert feat adds +5)
func CalculateInitiativeModifier(dexMod, bonuses int) int {
	return dexMod + bonuses
}

// CalculateSpellSaveDC calculates spell save DC.
// DC = 8 + proficiencyBonus + spellcastingAbilityModifier
func CalculateSpellSaveDC(proficiencyBonus, spellcastingMod int) int {
	return 8 + proficiencyBonus + spellcastingMod
}

// CalculateSpellAttackBonus calculates spell attack bonus.
// Bonus = proficiencyBonus + spellcastingAbilityModifier
func CalculateSpellAttackBonus(proficiencyBonus, spellcastingMod int) int {
	return proficiencyBonus + spellcastingMod
}

// CalculateMaxPreparedSpellsFromFormula calculates max prepared spells from a formula.
// Common formulas: "level + wisdom", "level/2 + charisma", etc.
// This is a simplified version - the full implementation is in character.go
func CalculateMaxPreparedSpellsFromFormula(level, abilityMod int) int {
	result := level + abilityMod
	if result < 1 {
		return 1
	}
	return result
}

// CalculateSkillModifier calculates the total modifier for a skill check.
func CalculateSkillModifier(abilityMod int, proficiency ProficiencyLevel, proficiencyBonus int) int {
	modifier := abilityMod

	switch proficiency {
	case Proficient:
		modifier += proficiencyBonus
	case Expertise:
		modifier += proficiencyBonus * 2
	}

	return modifier
}

// CalculateHPRatio calculates the ratio of current to max HP (for proportional adjustments).
func CalculateHPRatio(current, max int) float64 {
	if max == 0 {
		return 1.0
	}
	return float64(current) / float64(max)
}

// ApplyHPRatio applies an HP ratio to a new max HP value.
func ApplyHPRatio(newMax int, ratio float64) int {
	newCurrent := int(float64(newMax) * ratio)
	if newCurrent < 1 {
		return 1
	}
	return newCurrent
}
