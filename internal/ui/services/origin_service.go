// internal/ui/services/origin_service.go
package services

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// OriginService handles origin-related business logic
type OriginService struct{}

// NewOriginService creates a new origin service
func NewOriginService() *OriginService {
	return &OriginService{}
}

// ApplyOrigin applies an origin to a character, removing the old one if present
// Returns success message
func (s *OriginService) ApplyOrigin(character *models.Character, origin *models.Origin, chosenAbility string) string {
	// Remove old origin first
	if character.Origin != "" {
		oldOrigin := models.GetOriginByName(character.Origin)
		if oldOrigin != nil {
			models.RemoveOriginBenefits(character, *oldOrigin)
		}
	}

	// Apply new origin
	character.Origin = origin.Name
	if chosenAbility != "" {
		models.ApplyOriginBenefits(character, *origin, chosenAbility)
		return fmt.Sprintf("Origin changed to: %s (+1 %s)!", origin.Name, chosenAbility)
	} else {
		models.ApplyOriginBenefits(character, *origin, "")
		return fmt.Sprintf("Origin changed to: %s!", origin.Name)
	}
}

// RequiresAbilityChoice checks if an origin requires an ability score choice
func (s *OriginService) RequiresAbilityChoice(origin *models.Origin) bool {
	return models.HasOriginAbilityChoice(*origin)
}

// GetAbilityChoices returns the ability choices available for an origin
func (s *OriginService) GetAbilityChoices(origin *models.Origin) []string {
	return models.GetOriginAbilityChoices(*origin)
}
