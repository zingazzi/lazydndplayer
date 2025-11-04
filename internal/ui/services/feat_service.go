// internal/ui/services/feat_service.go
package services

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// FeatService handles feat-related business logic
type FeatService struct{}

// NewFeatService creates a new feat service
func NewFeatService() *FeatService {
	return &FeatService{}
}

// ApplyFeat applies a feat to a character
// Returns an error message if the feat cannot be applied, or success message
func (s *FeatService) ApplyFeat(character *models.Character, feat *models.Feat, chosenAbility string) (string, error) {
	// Check if character already has this feat (and it's not repeatable)
	if models.HasFeat(character, feat.Name) && !feat.Repeatable {
		return "", fmt.Errorf("you already have %s and it's not repeatable", feat.Name)
	}

	// Check prerequisites
	if !s.CheckPrerequisites(character, feat) {
		return "", fmt.Errorf("cannot select %s: prerequisites not met", feat.Name)
	}

	// Add feat to character
	err := models.AddFeatToCharacter(character, feat.Name)
	if err != nil {
		return "", fmt.Errorf("error adding feat: %w", err)
	}

	// Apply feat benefits
	if chosenAbility != "" {
		models.ApplyFeatBenefits(character, *feat, chosenAbility)
		return fmt.Sprintf("Feat gained: %s (+1 %s)!", feat.Name, chosenAbility), nil
	} else {
		models.ApplyFeatBenefits(character, *feat, "")
		return fmt.Sprintf("Feat gained: %s!", feat.Name), nil
	}
}

// RemoveFeat removes a feat from a character
// Returns success message
func (s *FeatService) RemoveFeat(character *models.Character, feat *models.Feat) string {
	// Remove from character's feat list
	for i, featName := range character.Feats {
		if featName == feat.Name {
			character.Feats = append(character.Feats[:i], character.Feats[i+1:]...)
			break
		}
	}

	// Remove feat benefits
	models.RemoveFeatBenefits(character, *feat)

	return fmt.Sprintf("Feat removed: %s (benefits reversed)", feat.Name)
}

// CheckPrerequisites checks if a character meets the prerequisites for a feat
func (s *FeatService) CheckPrerequisites(character *models.Character, feat *models.Feat) bool {
	// This is a placeholder - actual prerequisite checking would be implemented here
	// For now, return true as prerequisite checking is handled by the selector component
	// The selector component's CanSelectCurrentFeat() method handles this
	return true
}

// RequiresAbilityChoice checks if a feat requires an ability score choice
func (s *FeatService) RequiresAbilityChoice(feat *models.Feat) bool {
	return models.HasAbilityChoice(*feat)
}

// GetAbilityChoices returns the ability choices available for a feat
func (s *FeatService) GetAbilityChoices(feat *models.Feat) []string {
	return models.GetAbilityChoices(*feat)
}

// CancelFeatSelection cancels a pending feat selection and removes it from the character
func (s *FeatService) CancelFeatSelection(character *models.Character, feat *models.Feat) {
	// Remove the feat from character since we're cancelling
	for i, featName := range character.Feats {
		if featName == feat.Name {
			character.Feats = append(character.Feats[:i], character.Feats[i+1:]...)
			break
		}
	}
}
