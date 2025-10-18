// internal/models/origins.go
package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// Origin represents a character origin from D&D 5e 2024
type Origin struct {
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	AbilityIncreases   *FeatAbilityIncrease  `json:"ability_increases,omitempty"`
	Feat               string                `json:"feat"`                         // Granted feat
	SkillProficiencies []string              `json:"skill_proficiencies"`
	ToolProficiencies  []string              `json:"tool_proficiencies"`
	Equipment          []string              `json:"equipment"`
}

// OriginsData represents the structure of origins.json
type OriginsData struct {
	Origins []Origin `json:"origins"`
}

var cachedOrigins *OriginsData

// LoadOriginsFromJSON loads all origins from the JSON file
func LoadOriginsFromJSON(filepath string) (*OriginsData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read origins file: %w", err)
	}

	var originsData OriginsData
	if err := json.Unmarshal(data, &originsData); err != nil {
		return nil, fmt.Errorf("failed to parse origins JSON: %w", err)
	}

	cachedOrigins = &originsData
	return &originsData, nil
}

// GetAllOrigins returns all origins, loading from file if necessary
func GetAllOrigins() []Origin {
	if cachedOrigins == nil {
		_, err := LoadOriginsFromJSON("data/origins.json")
		if err != nil {
			return []Origin{}
		}
	}
	return cachedOrigins.Origins
}

// GetOriginByName returns an origin by name
func GetOriginByName(name string) *Origin {
	origins := GetAllOrigins()
	for i := range origins {
		if origins[i].Name == name {
			return &origins[i]
		}
	}
	return nil
}

// HasAbilityChoice returns true if the origin offers ability score choices
func HasOriginAbilityChoice(origin Origin) bool {
	return origin.AbilityIncreases != nil && len(origin.AbilityIncreases.Choices) > 0
}

// GetOriginAbilityChoices returns the ability choices for an origin
func GetOriginAbilityChoices(origin Origin) []string {
	if origin.AbilityIncreases != nil {
		return origin.AbilityIncreases.Choices
	}
	return []string{}
}

// ApplyOriginBenefits applies all benefits from an origin to a character
func ApplyOriginBenefits(char *Character, origin Origin, chosenAbility string) error {
	source := BenefitSource{
		Type: "origin",
		Name: origin.Name,
	}

	applier := NewBenefitApplier(char)

	// Apply ability increases
	if origin.AbilityIncreases != nil {
		if len(origin.AbilityIncreases.Choices) > 0 {
			// Multiple choice - use chosenAbility
			if chosenAbility != "" {
				applier.AddAbilityScore(source, chosenAbility, origin.AbilityIncreases.Amount)
			}
		} else if origin.AbilityIncreases.Ability != "" {
			// Single ability
			applier.AddAbilityScore(source, origin.AbilityIncreases.Ability, origin.AbilityIncreases.Amount)
		}
	}

	// Apply skill proficiencies
	for _, skill := range origin.SkillProficiencies {
		applier.AddSkillProficiency(source, skill)
	}

	// Apply tool proficiencies
	for _, tool := range origin.ToolProficiencies {
		applier.AddToolProficiency(source, tool)
	}

	// Apply the granted feat
	if origin.Feat != "" {
		// Check if character already has this feat
		alreadyHasFeat := false
		for _, existingFeat := range char.Feats {
			if existingFeat == origin.Feat {
				alreadyHasFeat = true
				break
			}
		}

		// Only add if not already present
		if !alreadyHasFeat {
			grantedFeat := GetFeatByName(origin.Feat)
			if grantedFeat != nil {
				// Add feat to character's feat list
				char.Feats = append(char.Feats, origin.Feat)

				// Apply feat benefits (this will track them separately)
				ApplyFeatBenefits(char, *grantedFeat, "")
			}
		}
	}

	// Update derived stats after applying benefits
	char.UpdateDerivedStats()
	return nil
}

// RemoveOriginBenefits removes the mechanical benefits of an origin from a character
func RemoveOriginBenefits(char *Character, origin Origin) error {
	// Check if this origin granted a feat
	if origin.Feat != "" {
		// Check if the feat was granted by this origin (tracked in BenefitTracker)
		// We track origin benefits separately, but feats are in the main feat list
		// Only remove the feat if it was added by the origin (not manually added before)

		// Count how many times this feat appears (should only be once, but let's be safe)
		featCount := 0
		for _, featName := range char.Feats {
			if featName == origin.Feat {
				featCount++
			}
		}

		// Only remove if we have the feat (but don't remove if user added it manually before selecting origin)
		// We'll check if the feat benefits are tracked from "feat" source
		grantedFeat := GetFeatByName(origin.Feat)
		if grantedFeat != nil && featCount > 0 {
			// Check if this feat has benefits tracked from "feat" source (meaning origin added it)
			benefitsFromFeat := char.BenefitTracker.GetBenefitsBySource("feat", origin.Feat)

			// If origin granted this feat, it would have created feat benefits
			// Only remove the feat from the list if we have those tracked benefits
			if len(benefitsFromFeat) > 0 {
				// Find and remove feat from list (only once)
				for i, featName := range char.Feats {
					if featName == origin.Feat {
						char.Feats = append(char.Feats[:i], char.Feats[i+1:]...)
						break
					}
				}

				// Remove feat benefits
				RemoveFeatBenefits(char, *grantedFeat)
			}
		}
	}

	// Remove origin-specific benefits
	remover := NewBenefitRemover(char)
	return remover.RemoveAllBenefits("origin", origin.Name)
}

// String returns a formatted string representation of an origin
func (o *Origin) String() string {
	result := fmt.Sprintf("%s\n%s\n", o.Name, o.Description)

	if o.AbilityIncreases != nil {
		if len(o.AbilityIncreases.Choices) > 0 {
			result += fmt.Sprintf("Ability Choice: +%d to %v\n",
				o.AbilityIncreases.Amount,
				o.AbilityIncreases.Choices)
		} else if o.AbilityIncreases.Ability != "" {
			result += fmt.Sprintf("Ability: +%d %s\n",
				o.AbilityIncreases.Amount,
				o.AbilityIncreases.Ability)
		}
	}

	if o.Feat != "" {
		result += fmt.Sprintf("Feat: %s\n", o.Feat)
	}

	if len(o.SkillProficiencies) > 0 {
		result += fmt.Sprintf("Skills: %v\n", o.SkillProficiencies)
	}

	if len(o.ToolProficiencies) > 0 {
		result += fmt.Sprintf("Tools: %v\n", o.ToolProficiencies)
	}

	return result
}
