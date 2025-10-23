// internal/services/character_service.go
package services

import (
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CharacterService provides high-level character management operations.
// This service layer orchestrates model operations and provides a clean API for the UI.
type CharacterService struct {
	classRepo models.ClassRepository
	validator models.CharacterValidator
	hpCalc    models.HPCalculator
}

// NewCharacterService creates a new character service with default implementations.
func NewCharacterService() *CharacterService {
	return &CharacterService{
		classRepo: &models.DefaultClassRepository{},
		validator: &models.DefaultCharacterValidator{},
		hpCalc:    models.NewDefaultHPCalculator(models.GetDefaultDiceRoller()),
	}
}

// NewCharacterServiceWithDeps creates a service with injected dependencies (for testing).
func NewCharacterServiceWithDeps(
	classRepo models.ClassRepository,
	validator models.CharacterValidator,
	hpCalc models.HPCalculator,
) *CharacterService {
	return &CharacterService{
		classRepo: classRepo,
		validator: validator,
		hpCalc:    hpCalc,
	}
}

// LevelUpCharacter handles the complete level-up process.
func (cs *CharacterService) LevelUpCharacter(char *models.Character, options models.LevelUpOptions) (*models.LevelUpResult, error) {
	// Delegate to the model's LevelUp function
	// In the future, this could add additional business logic, logging, events, etc.
	return models.LevelUp(char, options)
}

// CanLevelUp checks if a character meets the requirements to level up in a class.
func (cs *CharacterService) CanLevelUp(char *models.Character, className string) (bool, string) {
	return cs.validator.CanLevelUp(char, className)
}

// GetAvailableClasses returns classes the character can multiclass into.
func (cs *CharacterService) GetAvailableClasses(char *models.Character) []models.Class {
	return models.GetAvailableClasses(char)
}

// ApplyClass changes the character's primary class.
func (cs *CharacterService) ApplyClass(char *models.Character, className string) error {
	return models.ApplyClassToCharacter(char, className)
}

// RollHP rolls HP for a level-up.
func (cs *CharacterService) RollHP(hitDie, conMod int, takeAverage bool) int {
	return cs.hpCalc.RollHP(hitDie, conMod, takeAverage)
}

// CalculateMaxHP calculates the maximum HP for a character.
func (cs *CharacterService) CalculateMaxHP(char *models.Character, class *models.Class) int {
	return cs.hpCalc.CalculateMaxHP(char, class)
}

// ApplyFeat applies a feat to a character.
func (cs *CharacterService) ApplyFeat(char *models.Character, feat models.Feat, choices map[string]interface{}) error {
	return models.ApplyFeatBenefits(char, feat, choices)
}

// RemoveFeat removes a feat from a character.
func (cs *CharacterService) RemoveFeat(char *models.Character, featName string) error {
	return models.ReverseFeatBenefits(char, featName)
}

// ApplySpecies applies a species to a character.
func (cs *CharacterService) ApplySpecies(char *models.Character, species models.SpeciesInfo, choices map[string]interface{}) error {
	return models.ApplySpeciesBenefits(char, species, choices)
}

// ApplyOrigin applies an origin to a character.
func (cs *CharacterService) ApplyOrigin(char *models.Character, origin models.Origin, choices map[string]interface{}) error {
	return models.ApplyOriginBenefits(char, origin, choices)
}

// RemoveOrigin removes an origin from a character.
func (cs *CharacterService) RemoveOrigin(char *models.Character, originName string) error {
	return models.ReverseOriginBenefits(char, originName)
}

// UpdateDerivedStats recalculates all character derived statistics.
func (cs *CharacterService) UpdateDerivedStats(char *models.Character) {
	char.UpdateDerivedStats()
}

// GetClass retrieves class data by name.
func (cs *CharacterService) GetClass(className string) *models.Class {
	return cs.classRepo.GetClassByName(className)
}

// GetAllClasses retrieves all available classes.
func (cs *CharacterService) GetAllClasses() []models.Class {
	return cs.classRepo.GetAllClasses()
}

