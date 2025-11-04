// internal/ui/services/class_service.go
package services

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ClassService handles class-related business logic
type ClassService struct{}

// NewClassService creates a new class service
func NewClassService() *ClassService {
	return &ClassService{}
}

// ApplyClass applies a class to a character
// Returns success message and error
func (s *ClassService) ApplyClass(character *models.Character, className string) (string, error) {
	err := models.ApplyClassToCharacter(character, className)
	if err != nil {
		return "", fmt.Errorf("error applying class: %w", err)
	}

	return fmt.Sprintf("Class changed to: %s (HP: %d/%d)", className, character.CurrentHP, character.MaxHP), nil
}

// RequiresSkillChoice checks if a class requires skill choices
func (s *ClassService) RequiresSkillChoice(className string) bool {
	classData := models.GetClassByName(className)
	if classData == nil {
		return false
	}
	return classData.SkillChoices != nil && classData.SkillChoices.Choose > 0
}

// GetSkillChoices returns the skill choices available for a class
// Returns the skill list as []string (as expected by the component)
func (s *ClassService) GetSkillChoices(className string) (from []string, choose int) {
	classData := models.GetClassByName(className)
	if classData == nil || classData.SkillChoices == nil {
		return nil, 0
	}
	return classData.SkillChoices.From, classData.SkillChoices.Choose
}

// ApplySubclass applies a subclass to a character's current class
// Returns success message and error
func (s *ClassService) ApplySubclass(character *models.Character, subclassName string) (string, error) {
	if len(character.Classes) == 0 {
		return "", fmt.Errorf("character has no classes")
	}

	// Update the most recent class
	character.Classes[len(character.Classes)-1].Subclass = subclassName

	// Grant subclass features for the current level
	className := character.Classes[len(character.Classes)-1].ClassName
	classLevel := character.Classes[len(character.Classes)-1].Level

	_ = models.GrantSubclassFeatures(character, className, subclassName, classLevel)
	// GrantSubclassFeatures returns []string (features granted), not an error

	return fmt.Sprintf("Subclass changed to: %s", subclassName), nil
}

// GetSubclasses returns available subclasses for a class at a given level
func (s *ClassService) GetSubclasses(className string, classLevel int) ([]models.Subclass, error) {
	classData := models.GetClassByName(className)
	if classData == nil {
		return nil, fmt.Errorf("class %s not found", className)
	}

	// Filter subclasses available at this level
	available := []models.Subclass{}
	for _, subclass := range classData.Subclasses {
		if subclass.SubclassLevel <= classLevel {
			available = append(available, subclass)
		}
	}

	return available, nil
}
