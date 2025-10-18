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
		// Unarmored: 10 + Dex modifier
		baseAC = 10 + dexMod
	}

	// Add shield bonus (+2)
	if equippedShield != nil {
		baseAC += 2
	}

	// Add any AC bonuses from feats/magic items
	baseAC += char.ACBonus

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
