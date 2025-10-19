// internal/models/attacks.go
package models

import (
	"fmt"
	"strings"
)

// Attack represents an attack option for a character
type Attack struct {
	Name              string `json:"name"`
	AttackBonus       int    `json:"attack_bonus"`        // Modifier to attack roll
	DamageDice        string `json:"damage_dice"`         // e.g., "1d8", "2d6"
	VersatileDamage   string `json:"versatile_damage"`    // e.g., "1d10" for two-handed
	DamageBonus       int    `json:"damage_bonus"`        // Modifier to damage
	DamageType        string `json:"damage_type"`         // e.g., "slashing", "bludgeoning"
	IsWeapon          bool   `json:"is_weapon"`           // True if from equipped weapon
	WeaponName        string `json:"weapon_name"`         // Name of weapon if applicable
	Range             string `json:"range,omitempty"`     // e.g., "5 ft.", "range 80/320 ft."
	Properties        []string `json:"properties,omitempty"` // Weapon properties (finesse, versatile, etc.)
}

// AttackList holds all available attacks for a character
type AttackList struct {
	Attacks []Attack `json:"attacks"`
}

// GenerateAttacks creates a list of all available attacks for a character
func GenerateAttacks(char *Character) AttackList {
	attacks := AttackList{
		Attacks: []Attack{},
	}

	// Always add Unarmed Strike
	strMod := char.AbilityScores.GetModifier("Strength")
	profBonus := char.ProficiencyBonus

	attacks.Attacks = append(attacks.Attacks, Attack{
		Name:         "Unarmed Strike",
		AttackBonus:  strMod + profBonus,
		DamageDice:   "1",
		DamageBonus:  strMod,
		DamageType:   "bludgeoning",
		IsWeapon:     false,
		Range:        "5 ft.",
		Properties:   []string{},
	})

	// Add attacks from equipped weapons
	for i := range char.Inventory.Items {
		item := &char.Inventory.Items[i]
		if !item.Equipped || item.Type != Weapon {
			continue
		}

		// Get weapon definition
		weaponDef := GetItemDefinitionByName(item.Name)
		if weaponDef == nil {
			continue
		}

		// Calculate ability modifiers
		strMod := char.AbilityScores.GetModifier("Strength")
		dexMod := char.AbilityScores.GetModifier("Dexterity")

		// Determine which ability modifier to use
		abilityMod := strMod // Default to Strength

		// Check for Finesse property - use higher of STR or DEX
		hasFinesse := false
		for _, prop := range weaponDef.Properties {
			if strings.EqualFold(strings.TrimSpace(prop), "finesse") {
				hasFinesse = true
				break
			}
		}

		if hasFinesse {
			// Finesse: use DEX if it's higher than STR
			if dexMod > strMod {
				abilityMod = dexMod
			}
		} else if strings.Contains(strings.ToLower(weaponDef.Subcategory), "ranged") {
			// Ranged weapons use DEX
			abilityMod = dexMod
		}
		// Otherwise, use STR (melee weapons)

		attackBonus := abilityMod + profBonus

		// Parse damage dice from weapon
		damageDice := "1d4"
		if weaponDef.Damage != "" {
			damageDice = parseDamageDice(weaponDef.Damage)
		}

		// Determine damage type
		damageType := "bludgeoning"
		if weaponDef.DamageType != "" {
			damageType = strings.ToLower(weaponDef.DamageType)
		}

		// Check for versatile property (e.g., "Versatile (1d10)")
		versatileDamage := ""
		for _, prop := range weaponDef.Properties {
			if strings.HasPrefix(strings.ToLower(prop), "versatile") {
				// Extract versatile damage from property like "Versatile (1d10)"
				if start := strings.Index(prop, "("); start != -1 {
					if end := strings.Index(prop, ")"); end != -1 {
						versatileDamage = strings.TrimSpace(prop[start+1 : end])
					}
				}
			}
		}

		attacks.Attacks = append(attacks.Attacks, Attack{
			Name:            item.Name,
			AttackBonus:     attackBonus,
			DamageDice:      damageDice,
			VersatileDamage: versatileDamage,
			DamageBonus:     abilityMod,
			DamageType:      damageType,
			IsWeapon:        true,
			WeaponName:      item.Name,
			Range:           weaponDef.Range,
			Properties:      weaponDef.Properties,
		})
	}

	return attacks
}

// parseDamageDice extracts the dice notation from weapon damage string
// Examples: "1d8", "1d6 slashing", "2d6" -> returns just the dice part
func parseDamageDice(damageStr string) string {
	parts := strings.Fields(damageStr)
	if len(parts) == 0 {
		return "1d4"
	}

	// First part should be the dice notation
	dice := parts[0]

	// Validate it's a proper dice notation (XdY format)
	if strings.Contains(dice, "d") {
		return dice
	}

	return "1d4"
}

// FormatAttackRoll formats an attack roll display
func (a *Attack) FormatAttackRoll(roll int, total int, advantage string) string {
	advStr := ""
	if advantage != "" {
		advStr = fmt.Sprintf(" [%s]", advantage)
	}

	bonus := ""
	if a.AttackBonus >= 0 {
		bonus = fmt.Sprintf("+%d", a.AttackBonus)
	} else {
		bonus = fmt.Sprintf("%d", a.AttackBonus)
	}

	return fmt.Sprintf("%s: Attack Roll = %d %s = %d%s", a.Name, roll, bonus, total, advStr)
}

// FormatDamageRoll formats a damage roll display
func (a *Attack) FormatDamageRoll(rolls []int, total int) string {
	bonus := ""
	if a.DamageBonus >= 0 {
		bonus = fmt.Sprintf("+%d", a.DamageBonus)
	} else {
		bonus = fmt.Sprintf("%d", a.DamageBonus)
	}

	rollsStr := ""
	if len(rolls) > 0 {
		rollsStr = fmt.Sprintf("%v", rolls)
	}

	return fmt.Sprintf("%s: Damage = %s %s = %d %s", a.Name, rollsStr, bonus, total, a.DamageType)
}

// GetAttackSummary returns a one-line summary of the attack
func (a *Attack) GetAttackSummary() string {
	bonus := ""
	if a.AttackBonus >= 0 {
		bonus = fmt.Sprintf("+%d", a.AttackBonus)
	} else {
		bonus = fmt.Sprintf("%d", a.AttackBonus)
	}

	dmgBonus := ""
	if a.DamageBonus >= 0 {
		dmgBonus = fmt.Sprintf("+%d", a.DamageBonus)
	} else {
		dmgBonus = fmt.Sprintf("%d", a.DamageBonus)
	}

	return fmt.Sprintf("%s to hit, %s%s %s", bonus, a.DamageDice, dmgBonus, a.DamageType)
}
