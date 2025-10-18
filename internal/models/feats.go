// internal/models/feats.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// FeatAbilityIncrease represents ability score increases from a feat
type FeatAbilityIncrease struct {
	Ability string   `json:"ability,omitempty"` // Single ability (e.g., "Charisma")
	Choices []string `json:"choices,omitempty"` // Multiple choice (e.g., ["Strength", "Dexterity"])
	Amount  int      `json:"amount"`            // Amount to increase
}

// Feat represents a character feat from D&D 5e 2024
type Feat struct {
	Name               string                `json:"name"`
	Category           string                `json:"category"`
	Prerequisite       string                `json:"prerequisite"`
	Repeatable         bool                  `json:"repeatable"`
	Benefits           []string              `json:"benefits"`
	Description        string                `json:"description"`
	AbilityIncreases   *FeatAbilityIncrease  `json:"ability_increases,omitempty"`
	SkillProficiencies []string              `json:"skill_proficiencies,omitempty"`
	Languages          []string              `json:"languages,omitempty"`
	GrantsSpells       []string              `json:"grants_spells,omitempty"`
	Note               string                `json:"note,omitempty"`
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

// HasAbilityChoice returns true if the feat requires choosing an ability
func HasAbilityChoice(feat Feat) bool {
	return feat.AbilityIncreases != nil && len(feat.AbilityIncreases.Choices) > 0
}

// GetAbilityChoices returns the list of abilities that can be chosen for a feat
func GetAbilityChoices(feat Feat) []string {
	if feat.AbilityIncreases != nil && len(feat.AbilityIncreases.Choices) > 0 {
		return feat.AbilityIncreases.Choices
	}
	return []string{}
}

// ApplyFeatBenefits applies the mechanical benefits of a feat to a character
// chosenAbility is used when the feat has multiple ability choices
func ApplyFeatBenefits(char *Character, feat Feat, chosenAbility string) error {
	source := BenefitSource{
		Type: "feat",
		Name: feat.Name,
	}

	applier := NewBenefitApplier(char)

	// Apply ability increases
	if feat.AbilityIncreases != nil {
		if len(feat.AbilityIncreases.Choices) > 0 {
			// Multiple choice - use chosenAbility
			if chosenAbility != "" {
				applier.AddAbilityScore(source, chosenAbility, feat.AbilityIncreases.Amount)
			}
		} else if feat.AbilityIncreases.Ability != "" {
			// Single ability
			applier.AddAbilityScore(source, feat.AbilityIncreases.Ability, feat.AbilityIncreases.Amount)
		}
	}

	// Apply skill proficiencies
	for _, skill := range feat.SkillProficiencies {
		applier.AddSkillProficiency(source, skill)
	}

	// Apply languages
	for _, lang := range feat.Languages {
		applier.AddLanguage(source, lang)
	}

	// Apply special benefits
	featNameLower := strings.ToLower(feat.Name)

	// Tough feat: +2 HP per level
	if featNameLower == "tough" {
		applier.AddHP(source, char.Level*2)
	}

	// Mobile feat: +10 speed
	if featNameLower == "mobile" {
		applier.AddSpeed(source, 10)
	}

	// Update derived stats after applying benefits
	char.UpdateDerivedStats()
	return nil
}

// RemoveFeatBenefits removes the mechanical benefits of a feat from a character
func RemoveFeatBenefits(char *Character, feat Feat) error {
	remover := NewBenefitRemover(char)
	return remover.RemoveAllBenefits("feat", feat.Name)
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

	if feat.AbilityIncreases != nil {
		sb.WriteString("Ability Increases:\n")
		if len(feat.AbilityIncreases.Choices) > 0 {
			sb.WriteString(fmt.Sprintf("  • Choose one: %s (+%d)\n",
				strings.Join(feat.AbilityIncreases.Choices, " or "),
				feat.AbilityIncreases.Amount))
		} else if feat.AbilityIncreases.Ability != "" {
			sb.WriteString(fmt.Sprintf("  • %s: +%d\n",
				feat.AbilityIncreases.Ability,
				feat.AbilityIncreases.Amount))
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
