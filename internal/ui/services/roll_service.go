// internal/ui/services/roll_service.go
package services

import (
	"fmt"
	"strings"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// RollService handles dice rolling and roll calculations
type RollService struct{}

// NewRollService creates a new roll service
func NewRollService() *RollService {
	return &RollService{}
}

// CalculateSavingThrowRoll calculates the dice expression for a saving throw
// Returns the expression (e.g., "1d20+5") and a message
func (s *RollService) CalculateSavingThrowRoll(character *models.Character, ability models.AbilityType) (expression string, message string) {
	modifier := character.AbilityScores.GetModifier(ability)
	abilityFullName := s.getAbilityFullName(ability)

	// Check if proficient in this saving throw
	isProficient := false
	for _, prof := range character.SavingThrowProficiencies {
		if strings.EqualFold(prof, abilityFullName) {
			isProficient = true
			break
		}
	}

	// Add proficiency bonus if proficient
	if isProficient {
		modifier += character.ProficiencyBonus
	}

	// Roll 1d20 + modifier
	expression = fmt.Sprintf("1d20%+d", modifier)

	profStr := ""
	if isProficient {
		profStr = " (proficient)"
	}
	message = fmt.Sprintf("Rolled %s saving throw%s: %s", abilityFullName, profStr, expression)

	return expression, message
}

// CalculateAbilityCheckRoll calculates the dice expression for an ability check
// Returns the expression (e.g., "1d20+3") and a message
func (s *RollService) CalculateAbilityCheckRoll(character *models.Character, ability models.AbilityType) (expression string, message string) {
	modifier := character.AbilityScores.GetModifier(ability)
	abilityFullName := s.getAbilityFullName(ability)

	// Apply Jack of All Trades bonus (half proficiency, rounded down) for ability checks without proficiency
	jackOfAllTradesBonus := models.GetJackOfAllTradesBonus(character, models.NotProficient)
	modifier += jackOfAllTradesBonus

	// Roll 1d20 + modifier
	expression = fmt.Sprintf("1d20%+d", modifier)
	message = fmt.Sprintf("Rolled %s ability check: %s", abilityFullName, expression)
	if jackOfAllTradesBonus > 0 {
		message += fmt.Sprintf(" (Jack of All Trades: +%d)", jackOfAllTradesBonus)
	}

	return expression, message
}

// getAbilityFullName converts AbilityType to full name string
func (s *RollService) getAbilityFullName(ability models.AbilityType) string {
	switch ability {
	case models.Strength:
		return "Strength"
	case models.Dexterity:
		return "Dexterity"
	case models.Constitution:
		return "Constitution"
	case models.Intelligence:
		return "Intelligence"
	case models.Wisdom:
		return "Wisdom"
	case models.Charisma:
		return "Charisma"
	default:
		return "Unknown"
	}
}
