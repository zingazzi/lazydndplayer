// internal/models/armor_calculation.go
package models

import (
	"strconv"
	"strings"
)

// ArmorType represents the type of armor
type ArmorType int

const (
	ArmorTypeNone ArmorType = iota
	ArmorTypeLight
	ArmorTypeMedium
	ArmorTypeHeavy
	ArmorTypeShield
)

// CalculateAC calculates the character's Armor Class based on equipped armor and shield
func CalculateAC(char *Character) int {
	dexMod := char.AbilityScores.GetModifier("Dexterity")

	// Find equipped armor and shield
	var equippedArmor *Item
	var equippedShield *Item

	for i := range char.Inventory.Items {
		item := &char.Inventory.Items[i]
		if !item.Equipped {
			continue
		}

		// Check if it's armor or shield
		if item.Type == Armor {
			itemDef := GetItemDefinitionByName(item.Name)
			if itemDef != nil {
				if strings.ToLower(itemDef.Subcategory) == "shield" {
					equippedShield = item
				} else {
					equippedArmor = item
				}
			}
		}
	}

	// Calculate base AC
	baseAC := 10

	if equippedArmor != nil {
		// Get armor definition to determine AC and type
		armorDef := GetItemDefinitionByName(equippedArmor.Name)
		if armorDef != nil {
			armorType := getArmorType(armorDef.Subcategory)
			armorBaseAC := parseArmorAC(armorDef.AC)

			switch armorType {
			case ArmorTypeLight:
				// Light armor: Base AC + full Dex modifier
				baseAC = armorBaseAC + dexMod

			case ArmorTypeMedium:
				// Medium armor: Base AC + Dex modifier (max +2)
				dexBonus := dexMod
				if dexBonus > 2 {
					dexBonus = 2
				}
				baseAC = armorBaseAC + dexBonus

			case ArmorTypeHeavy:
				// Heavy armor: Base AC only (no Dex modifier)
				baseAC = armorBaseAC
			}
		}
	} else {
		// Unarmored: Check for special unarmored defense
		if char.IsMonk() && char.HasFeature("Unarmored Defense") {
			// Monk Unarmored Defense: 10 + Dex modifier + Wis modifier
			monk := char.GetMonkMechanics()
			baseAC = monk.CalculateUnarmoredAC()
		} else if char.Class == "Barbarian" {
			// Barbarian Unarmored Defense: 10 + Dex modifier + Con modifier
			conMod := char.AbilityScores.GetModifier("Constitution")
			baseAC = 10 + dexMod + conMod
		} else {
			// Standard unarmored: 10 + Dex modifier
			baseAC = 10 + dexMod
		}
	}

	// Add shield bonus (+2)
	if equippedShield != nil {
		baseAC += 2
	}

	// Add any AC bonuses from feats/magic items
	baseAC += char.ACBonus

	// Apply conditional fighting style bonuses
	baseAC += GetFightingStyleACBonus(char, equippedArmor != nil)

	// Apply conditional feat bonuses (e.g., Dual Wielder)
	baseAC += GetFeatACBonus(char)

	return baseAC
}

// getArmorType determines the armor type from subcategory
func getArmorType(subcategory string) ArmorType {
	subcategory = strings.ToLower(subcategory)
	switch subcategory {
	case "light":
		return ArmorTypeLight
	case "medium":
		return ArmorTypeMedium
	case "heavy":
		return ArmorTypeHeavy
	case "shield":
		return ArmorTypeShield
	default:
		return ArmorTypeNone
	}
}

// parseArmorAC extracts the base AC value from the AC string
// Examples: "11 + Dex" -> 11, "16" -> 16, "14 + Dex (max 2)" -> 14
func parseArmorAC(acString string) int {
	// Remove everything after +, spaces, and parse the number
	acString = strings.TrimSpace(acString)

	// Split by + or ( to get just the number
	parts := strings.FieldsFunc(acString, func(r rune) bool {
		return r == '+' || r == '('
	})

	if len(parts) == 0 {
		return 10
	}

	// Parse the first part as the base AC
	baseAC, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 10
	}

	return baseAC
}

// UnequipOtherArmor unequips any other armor when equipping new armor
// (only one armor piece can be equipped at a time, shields are separate)
func UnequipOtherArmor(char *Character, itemToEquip *Item) {
	// Get the item definition to check if it's armor or shield
	itemDef := GetItemDefinitionByName(itemToEquip.Name)
	if itemDef == nil || itemDef.Category != "armor" {
		return
	}

	isShield := strings.ToLower(itemDef.Subcategory) == "shield"

	// Unequip other armor (but not shields if we're equipping armor, and vice versa)
	for i := range char.Inventory.Items {
		item := &char.Inventory.Items[i]

		// Skip if not equipped or if it's the same item
		if !item.Equipped || item == itemToEquip {
			continue
		}

		// Check if this is armor
		if item.Type == Armor {
			otherItemDef := GetItemDefinitionByName(item.Name)
			if otherItemDef != nil {
				otherIsShield := strings.ToLower(otherItemDef.Subcategory) == "shield"

				// Unequip if:
				// - We're equipping armor and this is also armor (not shield)
				// - We're equipping shield and this is also shield
				if isShield == otherIsShield {
					item.Equipped = false
				}
			}
		}
	}
}

// GetEquippedArmorInfo returns information about equipped armor and shield
func GetEquippedArmorInfo(char *Character) (armor string, shield string) {
	armor = "None"
	shield = "None"

	for i := range char.Inventory.Items {
		item := &char.Inventory.Items[i]
		if !item.Equipped || item.Type != Armor {
			continue
		}

		itemDef := GetItemDefinitionByName(item.Name)
		if itemDef != nil {
			if strings.ToLower(itemDef.Subcategory) == "shield" {
				shield = item.Name
			} else {
				armor = item.Name
			}
		}
	}

	return armor, shield
}

// GetFightingStyleACBonus returns conditional AC bonus from fighting style
func GetFightingStyleACBonus(char *Character, isWearingArmor bool) int {
	if char.FightingStyle == "" {
		return 0
	}

	style := GetFightingStyleByName(char.FightingStyle)
	if style == nil {
		return 0
	}

	// Check for conditional AC bonus
	if bonus, ok := style.Benefits["ac_bonus_conditional"].(float64); ok {
		if condition, hasCondition := style.Benefits["condition"].(string); hasCondition {
			switch condition {
			case "wearing_armor":
				// Defense fighting style: +1 AC only when wearing armor
				if isWearingArmor {
					return int(bonus)
				}
			case "dual_wielding":
				// Two-Weapon Fighting style (if implemented)
				// Check if wielding two weapons
				// TODO: implement when weapon tracking is added
			}
		}
	}

	return 0
}

// GetFeatACBonus returns conditional AC bonus from feats
func GetFeatACBonus(char *Character) int {
	totalBonus := 0

	// Check for Dual Wielder feat
	for _, featName := range char.Feats {
		if featName == "Dual Wielder" {
			// Check if wielding two melee weapons
			if IsDualWieldingMelee(char) {
				totalBonus += 1
			}
		}
	}

	return totalBonus
}

// IsDualWieldingMelee checks if the character is wielding two separate melee weapons
func IsDualWieldingMelee(char *Character) bool {
	equippedMeleeWeapons := 0

	for i := range char.Inventory.Items {
		item := &char.Inventory.Items[i]
		if !item.Equipped || item.Type != Weapon {
			continue
		}

		// Check if it's a melee weapon
		weaponDef := GetItemDefinitionByName(item.Name)
		if weaponDef != nil {
			// Check if it's not ranged
			if !strings.Contains(strings.ToLower(weaponDef.Subcategory), "ranged") {
				equippedMeleeWeapons++
			}
		}
	}

	// Must have exactly 2 melee weapons equipped
	return equippedMeleeWeapons == 2
}
