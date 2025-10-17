// internal/leveling/guide.go
package leveling

import (
	"encoding/json"
	"fmt"
	"os"
)

// LevelInfo contains all information about what's gained at a specific level
type LevelInfo struct {
	ProficiencyBonus int               `json:"proficiency_bonus"`
	Features         []string          `json:"features"`
	SpellSlots       map[string]int    `json:"spell_slots"`
	SpellsKnown      int               `json:"spells_known"`
	SpellsToAdd      int               `json:"spells_to_add"`
	Notes            string            `json:"notes"`
}

// ClassInfo contains leveling information for a class
type ClassInfo struct {
	HitDie int                    `json:"hit_die"`
	Levels map[string]LevelInfo   `json:"levels"`
}

// SpeciesScaling contains species-specific scaling information
type SpeciesScaling struct {
	HPPerLevel int               `json:"hp_per_level"`
	Features   map[string]string `json:"features"`
	Note       string            `json:"note"`
}

// LevelingGuide contains the complete leveling guide
type LevelingGuide struct {
	Classes         map[string]ClassInfo         `json:"classes"`
	SpeciesScaling  map[string]SpeciesScaling    `json:"species_scaling"`
	GeneralRules    map[string]interface{}       `json:"general_rules"`
}

var cachedGuide *LevelingGuide

// LoadLevelingGuide loads the leveling guide from JSON
func LoadLevelingGuide() (*LevelingGuide, error) {
	if cachedGuide != nil {
		return cachedGuide, nil
	}

	file, err := os.ReadFile("data/leveling_guide.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read leveling guide: %w", err)
	}

	var guide LevelingGuide
	if err := json.Unmarshal(file, &guide); err != nil {
		return nil, fmt.Errorf("failed to parse leveling guide: %w", err)
	}

	cachedGuide = &guide
	return cachedGuide, nil
}

// GetLevelInfo returns information about a specific class and level
func GetLevelInfo(class string, level int) (*LevelInfo, error) {
	guide, err := LoadLevelingGuide()
	if err != nil {
		return nil, err
	}

	classInfo, ok := guide.Classes[class]
	if !ok {
		return nil, fmt.Errorf("class %s not found in leveling guide", class)
	}

	levelKey := fmt.Sprintf("%d", level)
	levelInfo, ok := classInfo.Levels[levelKey]
	if !ok {
		// Return empty info for levels not defined
		return &LevelInfo{}, nil
	}

	return &levelInfo, nil
}

// GetSpeciesScaling returns scaling information for a species
func GetSpeciesScaling(species string) *SpeciesScaling {
	guide, err := LoadLevelingGuide()
	if err != nil {
		return nil
	}

	if scaling, ok := guide.SpeciesScaling[species]; ok {
		return &scaling
	}

	return nil
}

// GetHitDieForClass returns the hit die for a class from the guide
func GetHitDieForClass(class string) int {
	guide, err := LoadLevelingGuide()
	if err != nil {
		return GetHitDie(class) // Fallback to hardcoded
	}

	if classInfo, ok := guide.Classes[class]; ok {
		return classInfo.HitDie
	}

	return GetHitDie(class) // Fallback
}

// FormatLevelUpSummary creates a formatted summary of what's gained at a level
func FormatLevelUpSummary(class string, level int) string {
	info, err := GetLevelInfo(class, level)
	if err != nil {
		return fmt.Sprintf("Level %d: No information available", level)
	}

	summary := fmt.Sprintf("=== LEVEL %d ===\n\n", level)

	if info.ProficiencyBonus > 0 {
		summary += fmt.Sprintf("Proficiency Bonus: +%d\n\n", info.ProficiencyBonus)
	}

	if len(info.Features) > 0 {
		summary += "Features Gained:\n"
		for _, feature := range info.Features {
			summary += fmt.Sprintf("  • %s\n", feature)
		}
		summary += "\n"
	}

	if len(info.SpellSlots) > 0 {
		summary += "Spell Slots:\n"
		for level, slots := range info.SpellSlots {
			summary += fmt.Sprintf("  Level %s: %d slots\n", level, slots)
		}
		summary += "\n"
	}

	if info.SpellsKnown > 0 {
		summary += fmt.Sprintf("Spells Known: %d\n\n", info.SpellsKnown)
	}

	if info.SpellsToAdd > 0 {
		summary += fmt.Sprintf("Add %d spells to your spellbook\n\n", info.SpellsToAdd)
	}

	if info.Notes != "" {
		summary += fmt.Sprintf("Notes:\n%s\n\n", info.Notes)
	}

	return summary
}
