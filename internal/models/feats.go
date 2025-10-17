// internal/models/feats.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Feat represents a character feat from D&D 5e 2024
type Feat struct {
	Name             string                 `json:"name"`
	Category         string                 `json:"category"`
	Prerequisite     string                 `json:"prerequisite"`
	Repeatable       bool                   `json:"repeatable"`
	Benefits         []string               `json:"benefits"`
	Description      string                 `json:"description"`
	AbilityIncreases map[string]int         `json:"ability_increases"`
	GrantsSpells     []string               `json:"grants_spells,omitempty"`
	Note             string                 `json:"note,omitempty"`
}

// FeatsData represents the structure of feats.json
type FeatsData struct {
	Feats []Feat `json:"feats"`
}

var cachedFeats *FeatsData

// LoadFeatsFromJSON loads all feats from the JSON file
func LoadFeatsFromJSON() (*FeatsData, error) {
	if cachedFeats != nil {
		return cachedFeats, nil
	}

	file, err := os.ReadFile("data/feats.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read feats.json: %w", err)
	}

	var data FeatsData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to parse feats.json: %w", err)
	}

	cachedFeats = &data
	return cachedFeats, nil
}

// GetAllFeats returns all available feats
func GetAllFeats() []Feat {
	data, err := LoadFeatsFromJSON()
	if err != nil {
		fmt.Printf("Error loading feats: %v\n", err)
		return []Feat{}
	}
	return data.Feats
}

// GetFeatByName returns a specific feat by name
func GetFeatByName(name string) *Feat {
	feats := GetAllFeats()
	nameLower := strings.ToLower(name)
	
	for _, feat := range feats {
		if strings.ToLower(feat.Name) == nameLower {
			return &feat
		}
	}
	return nil
}

// GetFeatsForCharacter returns feats that the character can take
// based on prerequisites
func GetFeatsForCharacter(char *Character) []Feat {
	allFeats := GetAllFeats()
	availableFeats := []Feat{}
	
	for _, feat := range allFeats {
		if CanTakeFeat(char, feat) {
			availableFeats = append(availableFeats, feat)
		}
	}
	
	return availableFeats
}

// CanTakeFeat checks if a character meets the prerequisites for a feat
func CanTakeFeat(char *Character, feat Feat) bool {
	prereq := strings.ToLower(feat.Prerequisite)
	
	// No prerequisite
	if prereq == "none" || prereq == "" {
		return true
	}
	
	// Check for ability score requirements
	if strings.Contains(prereq, "strength") {
		if strings.Contains(prereq, "13") && char.AbilityScores.Strength < 13 {
			return false
		}
	}
	if strings.Contains(prereq, "dexterity") {
		if strings.Contains(prereq, "13") && char.AbilityScores.Dexterity < 13 {
			return false
		}
	}
	if strings.Contains(prereq, "intelligence") {
		if strings.Contains(prereq, "13") && char.AbilityScores.Intelligence < 13 {
			return false
		}
	}
	if strings.Contains(prereq, "wisdom") {
		if strings.Contains(prereq, "13") && char.AbilityScores.Wisdom < 13 {
			return false
		}
	}
	if strings.Contains(prereq, "charisma") {
		if strings.Contains(prereq, "13") && char.AbilityScores.Charisma < 13 {
			return false
		}
	}
	
	// Check for spellcasting requirement
	if strings.Contains(prereq, "cast at least one spell") {
		// Check if character has any spells or is a spellcasting class
		hasSpells := len(char.SpellBook.Spells) > 0
		isSpellcaster := char.Class == "Wizard" || char.Class == "Cleric" || 
		                char.Class == "Druid" || char.Class == "Bard" || 
		                char.Class == "Sorcerer" || char.Class == "Warlock"
		if !hasSpells && !isSpellcaster {
			return false
		}
	}
	
	// Check for armor proficiency requirements
	// TODO: Implement armor proficiency tracking in character model
	
	return true
}

// HasFeat checks if a character already has a specific feat
func HasFeat(char *Character, featName string) bool {
	featNameLower := strings.ToLower(featName)
	for _, feat := range char.Feats {
		if strings.ToLower(feat) == featNameLower {
			return true
		}
	}
	return false
}

// AddFeatToCharacter adds a feat to the character's feat list
func AddFeatToCharacter(char *Character, featName string) error {
	feat := GetFeatByName(featName)
	if feat == nil {
		return fmt.Errorf("feat %s not found", featName)
	}
	
	// Check if character can take this feat
	if !CanTakeFeat(char, *feat) {
		return fmt.Errorf("character does not meet prerequisites for %s", featName)
	}
	
	// Check if feat is already taken and not repeatable
	if HasFeat(char, featName) && !feat.Repeatable {
		return fmt.Errorf("feat %s is already taken and is not repeatable", featName)
	}
	
	// Add feat to character's feat list
	char.Feats = append(char.Feats, featName)
	
	return nil
}

// ApplyFeatBenefits applies the mechanical benefits of a feat to a character
// This includes ability score increases, HP increases, etc.
func ApplyFeatBenefits(char *Character, feat Feat) {
	// Apply ability score increases
	for ability, increase := range feat.AbilityIncreases {
		abilityLower := strings.ToLower(ability)
		
		// Handle multiple choice abilities (e.g., "Strength or Dexterity")
		if strings.Contains(abilityLower, "or") {
			// This will need to be chosen by the player in the UI
			// For now, we'll skip automatic application
			continue
		}
		
		// Apply specific ability increases
		if strings.Contains(abilityLower, "strength") {
			char.AbilityScores.Strength = min(char.AbilityScores.Strength+increase, 20)
		} else if strings.Contains(abilityLower, "dexterity") {
			char.AbilityScores.Dexterity = min(char.AbilityScores.Dexterity+increase, 20)
		} else if strings.Contains(abilityLower, "constitution") {
			char.AbilityScores.Constitution = min(char.AbilityScores.Constitution+increase, 20)
		} else if strings.Contains(abilityLower, "intelligence") {
			char.AbilityScores.Intelligence = min(char.AbilityScores.Intelligence+increase, 20)
		} else if strings.Contains(abilityLower, "wisdom") {
			char.AbilityScores.Wisdom = min(char.AbilityScores.Wisdom+increase, 20)
		} else if strings.Contains(abilityLower, "charisma") {
			char.AbilityScores.Charisma = min(char.AbilityScores.Charisma+increase, 20)
		}
	}
	
	// Apply special benefits
	featNameLower := strings.ToLower(feat.Name)
	
	// Tough feat: +2 HP per level
	if featNameLower == "tough" {
		char.MaxHP += char.Level * 2
		char.CurrentHP = char.MaxHP
	}
	
	// Mobile feat: +10 speed
	if featNameLower == "mobile" {
		char.Speed += 10
	}
	
	// Durable feat: Already has Constitution increase handled above
	
	// Alert feat: +5 initiative (we'll add this when implementing initiative modifiers)
	
	// Update derived stats after applying benefits
	char.UpdateDerivedStats()
}

// GetFeatCategories returns all unique feat categories
func GetFeatCategories() []string {
	feats := GetAllFeats()
	categoryMap := make(map[string]bool)
	
	for _, feat := range feats {
		categoryMap[feat.Category] = true
	}
	
	categories := []string{}
	for category := range categoryMap {
		categories = append(categories, category)
	}
	
	return categories
}

// FormatFeatForDisplay returns a formatted string of feat information
func FormatFeatForDisplay(feat Feat) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("=== %s ===\n\n", feat.Name))
	sb.WriteString(fmt.Sprintf("Category: %s\n", feat.Category))
	sb.WriteString(fmt.Sprintf("Prerequisite: %s\n", feat.Prerequisite))
	
	if feat.Repeatable {
		sb.WriteString("Repeatable: Yes\n")
	}
	
	sb.WriteString(fmt.Sprintf("\nDescription:\n%s\n\n", feat.Description))
	
	if len(feat.Benefits) > 0 {
		sb.WriteString("Benefits:\n")
		for _, benefit := range feat.Benefits {
			sb.WriteString(fmt.Sprintf("  • %s\n", benefit))
		}
		sb.WriteString("\n")
	}
	
	if len(feat.AbilityIncreases) > 0 {
		sb.WriteString("Ability Increases:\n")
		for ability, increase := range feat.AbilityIncreases {
			sb.WriteString(fmt.Sprintf("  • %s: +%d\n", ability, increase))
		}
		sb.WriteString("\n")
	}
	
	if len(feat.GrantsSpells) > 0 {
		sb.WriteString("Grants Spells:\n")
		for _, spell := range feat.GrantsSpells {
			sb.WriteString(fmt.Sprintf("  • %s\n", spell))
		}
		sb.WriteString("\n")
	}
	
	if feat.Note != "" {
		sb.WriteString(fmt.Sprintf("Note: %s\n", feat.Note))
	}
	
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

