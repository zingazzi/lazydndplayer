// internal/models/attacks.go
package models

import (
	"fmt"
	"strings"
)

// Attack represents an attack option for a character
type Attack struct {
	Name                 string `json:"name"`
	AttackBonus          int    `json:"attack_bonus"`           // Modifier to attack roll
	DamageDice           string `json:"damage_dice"`            // e.g., "1d8", "2d6"
	VersatileDamage      string `json:"versatile_damage"`       // e.g., "1d10" for two-handed
	DamageBonus          int    `json:"damage_bonus"`           // Modifier to damage (one-handed)
	TwoHandDamageBonus   int    `json:"two_hand_damage_bonus"`  // Modifier when using two hands (no Dueling bonus)
	DamageType           string `json:"damage_type"`            // e.g., "slashing", "bludgeoning"
	IsWeapon             bool   `json:"is_weapon"`              // True if from equipped weapon
	WeaponName           string `json:"weapon_name"`            // Name of weapon if applicable
	Range                string `json:"range,omitempty"`        // e.g., "5 ft.", "range 80/320 ft."
	Properties           []string `json:"properties,omitempty"`  // Weapon properties (finesse, versatile, etc.)
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
	dexMod := char.AbilityScores.GetModifier("Dexterity")
	profBonus := char.ProficiencyBonus

	// Calculate unarmed strike damage and attack bonus
	unarmedDamage := "1"      // Default: 1 + STR mod
	attackMod := strMod       // Default: use STR
	damageBonus := strMod     // Default: use STR

	// Check for Monk Martial Arts
	if char.IsMonk() && char.HasFeature("Martial Arts") {
		monk := char.GetMonkMechanics()
		unarmedDamage = monk.GetMartialArtsDie() // 1d6->1d8->1d10->1d12

		// Monks can use Dex or Str (use whichever is higher)
		if dexMod > strMod {
			attackMod = dexMod
			damageBonus = dexMod
		}
	} else if char.FightingStyle == "Unarmed Fighting" {
		// Check if character has any weapons or shield equipped
		hasWeaponOrShield := false
		for i := range char.Inventory.Items {
			item := &char.Inventory.Items[i]
			if item.Equipped && (item.Type == Weapon || (item.Type == Armor && strings.Contains(strings.ToLower(item.Name), "shield"))) {
				hasWeaponOrShield = true
				break
			}
		}

		// Unarmed Fighting: 1d6 normally, 1d8 if no weapons or shield
		if hasWeaponOrShield {
			unarmedDamage = "1d6"
		} else {
			unarmedDamage = "1d8"
		}
	}

	attacks.Attacks = append(attacks.Attacks, Attack{
		Name:         "Unarmed Strike",
		AttackBonus:  attackMod + profBonus,
		DamageDice:   unarmedDamage,
		DamageBonus:  damageBonus,
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

		// Apply Archery fighting style bonus for ranged weapons
		isRanged := strings.Contains(strings.ToLower(weaponDef.Subcategory), "ranged")
		if char.FightingStyle == "Archery" && isRanged {
			attackBonus += 2
		}

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

		// Calculate damage bonus (ability modifier + fighting style bonuses)
		oneHandDamageBonus := abilityMod
		twoHandDamageBonus := abilityMod // For versatile weapons used two-handed

		// Check for Dueling fighting style
		// Dueling: +2 damage when wielding a melee weapon in one hand with no weapon in the other hand
		// IMPORTANT: Does NOT apply when using a versatile weapon with two hands!
		if char.FightingStyle == "Dueling" {
			// Check if this is a melee weapon (not ranged)
			isMelee := !strings.Contains(strings.ToLower(weaponDef.Subcategory), "ranged")

			// Check if it's a two-handed weapon (can't use Dueling at all)
			isTwoHanded := false
			for _, prop := range weaponDef.Properties {
				if strings.EqualFold(strings.TrimSpace(prop), "two-handed") {
					isTwoHanded = true
					break
				}
			}

			// Only apply Dueling bonus if it's a melee weapon that can be used one-handed
			if isMelee && !isTwoHanded {
				// Check if no other weapon is equipped (shield is okay)
				otherWeaponCount := 0
				for j := range char.Inventory.Items {
					otherItem := &char.Inventory.Items[j]
					if otherItem.Equipped && otherItem.Type == Weapon && otherItem.Name != item.Name {
						otherWeaponCount++
					}
				}

				// Apply Dueling +2 bonus ONLY to one-handed attacks
				if otherWeaponCount == 0 {
					oneHandDamageBonus += 2
					// twoHandDamageBonus stays as abilityMod (no Dueling bonus when using two hands)
				}
			}
		}

		attacks.Attacks = append(attacks.Attacks, Attack{
			Name:               item.Name,
			AttackBonus:        attackBonus,
			DamageDice:         damageDice,
			VersatileDamage:    versatileDamage,
			DamageBonus:        oneHandDamageBonus,
			TwoHandDamageBonus: twoHandDamageBonus,
			DamageType:         damageType,
			IsWeapon:           true,
			WeaponName:         item.Name,
			Range:              weaponDef.Range,
			Properties:         weaponDef.Properties,
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
