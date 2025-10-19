// internal/models/fighting_styles.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// FightingStyle represents a D&D 5e fighting style
type FightingStyle struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Requirements []string          `json:"requirements"`
	Benefits     map[string]any    `json:"benefits"`
}

// FightingStylesData represents the structure of fighting_styles.json
type FightingStylesData struct {
	FightingStyles []FightingStyle `json:"fighting_styles"`
}

var cachedFightingStyles *FightingStylesData

// LoadFightingStylesFromJSON loads all fighting styles from the JSON file
func LoadFightingStylesFromJSON(filepath string) (*FightingStylesData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fighting styles file: %w", err)
	}

	var stylesData FightingStylesData
	if err := json.Unmarshal(data, &stylesData); err != nil {
		return nil, fmt.Errorf("failed to parse fighting styles JSON: %w", err)
	}

	cachedFightingStyles = &stylesData
	return &stylesData, nil
}

// GetAllFightingStyles returns all available fighting styles
func GetAllFightingStyles() []FightingStyle {
	if cachedFightingStyles == nil {
		_, err := LoadFightingStylesFromJSON("data/fighting_styles.json")
		if err != nil {
			fmt.Printf("Error loading fighting styles: %v\n", err)
			return []FightingStyle{}
		}
	}
	return cachedFightingStyles.FightingStyles
}

// GetFightingStyleByName returns a specific fighting style by name
func GetFightingStyleByName(name string) *FightingStyle {
	styles := GetAllFightingStyles()
	for i := range styles {
		if styles[i].Name == name {
			return &styles[i]
		}
	}
	return nil
}

// ApplyFightingStyle applies a fighting style to a character
func ApplyFightingStyle(char *Character, styleName string) error {
	style := GetFightingStyleByName(styleName)
	if style == nil {
		return fmt.Errorf("fighting style %s not found", styleName)
	}

	// Set the fighting style name
	char.FightingStyle = styleName

	// Apply benefits through benefit tracker
	applier := NewBenefitApplier(char)
	source := BenefitSource{
		Type: "fighting_style",
		Name: styleName,
	}

	// Parse and apply benefits based on type
	if acBonus, ok := style.Benefits["ac_bonus"].(float64); ok {
		applier.AddACBonus(source, int(acBonus))
	}

	// Note: Other benefits like attack bonuses, special abilities, etc.
	// are tracked in the fighting style name and applied during combat
	// For now, we just store the name and AC bonus if applicable

	return nil
}

// RemoveFightingStyle removes a fighting style from a character
func RemoveFightingStyle(char *Character) error {
	if char.FightingStyle == "" {
		return nil // Nothing to remove
	}

	// Remove benefits
	remover := NewBenefitRemover(char)
	remover.RemoveAllBenefits("fighting_style", char.FightingStyle)

	// Clear the fighting style
	char.FightingStyle = ""

	return nil
}
