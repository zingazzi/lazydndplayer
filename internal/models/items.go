// internal/models/items.go
package models

import (
	"encoding/json"
	"os"
	"strings"
)

// ItemDefinition represents an item template from the items database
type ItemDefinition struct {
	Name                 string   `json:"name"`
	Category             string   `json:"category"`
	Subcategory          string   `json:"subcategory"`
	Weight               float64  `json:"weight"`
	PriceGP              float64  `json:"price_gp"`
	Description          string   `json:"description,omitempty"`
	Damage               string   `json:"damage,omitempty"`
	DamageType           string   `json:"damage_type,omitempty"`
	AC                   string   `json:"ac,omitempty"`
	Range                string   `json:"range,omitempty"`
	Properties           []string `json:"properties,omitempty"`
	StealthDisadvantage  bool     `json:"stealth_disadvantage,omitempty"`
	StrengthReq          int      `json:"strength_req,omitempty"`
	Equippable           bool     `json:"equippable,omitempty"`
	Mastery              string   `json:"mastery,omitempty"` // Weapon mastery property
}

// ItemsDatabase holds all item definitions
type ItemsDatabase struct {
	Weapons          []ItemDefinition `json:"weapons"`
	Armor            []ItemDefinition `json:"armor"`
	AdventuringGear  []ItemDefinition `json:"adventuring_gear"`
	Potions          []ItemDefinition `json:"potions"`
	MagicItems       []ItemDefinition `json:"magic_items"`
	Ammunition       []ItemDefinition `json:"ammunition"`
}

var itemsDB *ItemsDatabase

// LoadItemsFromJSON loads item definitions from JSON file
func LoadItemsFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	itemsDB = &ItemsDatabase{}
	return json.Unmarshal(data, itemsDB)
}

// GetAllItemDefinitions returns all item definitions as a flat list
func GetAllItemDefinitions() []ItemDefinition {
	if itemsDB == nil {
		// Try to load if not loaded yet
		LoadItemsFromJSON("data/items.json")
	}

	var allItems []ItemDefinition
	if itemsDB != nil {
		allItems = append(allItems, itemsDB.Weapons...)
		allItems = append(allItems, itemsDB.Armor...)
		allItems = append(allItems, itemsDB.AdventuringGear...)
		allItems = append(allItems, itemsDB.Potions...)
		allItems = append(allItems, itemsDB.MagicItems...)
		allItems = append(allItems, itemsDB.Ammunition...)
	}
	return allItems
}

// GetItemsByCategory returns items filtered by category
func GetItemsByCategory(category string) []ItemDefinition {
	if itemsDB == nil {
		LoadItemsFromJSON("data/items.json")
	}

	category = strings.ToLower(category)
	switch category {
	case "weapon", "weapons":
		return itemsDB.Weapons
	case "armor":
		return itemsDB.Armor
	case "gear", "adventuring gear":
		return itemsDB.AdventuringGear
	case "potion", "potions":
		return itemsDB.Potions
	case "magic", "magic items":
		return itemsDB.MagicItems
	case "ammunition":
		return itemsDB.Ammunition
	case "all", "":
		return GetAllItemDefinitions()
	default:
		return []ItemDefinition{}
	}
}

// SearchItems searches for items by name (fuzzy search)
func SearchItems(query string, category string) []ItemDefinition {
	items := GetItemsByCategory(category)
	if query == "" {
		return items
	}

	query = strings.ToLower(query)
	var results []ItemDefinition

	// Exact matches first
	for _, item := range items {
		if strings.ToLower(item.Name) == query {
			results = append(results, item)
		}
	}

	// Then prefix matches
	for _, item := range items {
		itemName := strings.ToLower(item.Name)
		if itemName != query && strings.HasPrefix(itemName, query) {
			results = append(results, item)
		}
	}

	// Then contains matches
	for _, item := range items {
		itemName := strings.ToLower(item.Name)
		if !strings.HasPrefix(itemName, query) && strings.Contains(itemName, query) {
			results = append(results, item)
		}
	}

	// Fuzzy matches (subsequence matching)
	for _, item := range items {
		itemName := strings.ToLower(item.Name)
		if !strings.Contains(itemName, query) && fuzzyMatch(query, itemName) {
			results = append(results, item)
		}
	}

	return results
}

// fuzzyMatch checks if all characters in query appear in order in target
func fuzzyMatch(query, target string) bool {
	queryIdx := 0
	for i := 0; i < len(target) && queryIdx < len(query); i++ {
		if target[i] == query[queryIdx] {
			queryIdx++
		}
	}
	return queryIdx == len(query)
}

// GetItemDefinitionByName returns an item definition by exact name
func GetItemDefinitionByName(name string) *ItemDefinition {
	items := GetAllItemDefinitions()
	for _, item := range items {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// IsEquippable returns whether an item can be equipped
func IsEquippable(item ItemDefinition) bool {
	category := strings.ToLower(item.Category)
	return category == "weapon" || category == "armor" || item.Equippable
}

// ConvertToInventoryItem converts an ItemDefinition to an Item for inventory
func ConvertToInventoryItem(def ItemDefinition, quantity int) Item {
	// Determine item type based on category
	var itemType ItemType
	switch strings.ToLower(def.Category) {
	case "weapon":
		itemType = Weapon
	case "armor":
		itemType = Armor
	case "potion":
		itemType = Potion
	case "magic":
		itemType = Magic
	case "gear":
		itemType = Gear
	case "ammunition":
		itemType = "Ammunition"
	default:
		itemType = Other
	}

	// Build description from item properties
	description := def.Description
	if def.Damage != "" {
		if description != "" {
			description += " | "
		}
		description += "Damage: " + def.Damage + " " + def.DamageType
	}
	if def.AC != "" {
		if description != "" {
			description += " | "
		}
		description += "AC: " + def.AC
	}
	if len(def.Properties) > 0 {
		if description != "" {
			description += " | "
		}
		description += strings.Join(def.Properties, ", ")
	}

	return Item{
		Name:        def.Name,
		Type:        itemType,
		Quantity:    quantity,
		Weight:      def.Weight,
		Description: description,
		Equipped:    false,
		Value:       int(def.PriceGP),
	}
}
