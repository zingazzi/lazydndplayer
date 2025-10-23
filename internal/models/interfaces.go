// internal/models/interfaces.go
package models

import "fmt"

// ClassRepository defines the interface for loading and retrieving class data
type ClassRepository interface {
	GetClassByName(name string) *Class
	GetAllClasses() []Class
	LoadClasses(dirPath string) error
}

// CharacterValidator defines validation rules for character operations
type CharacterValidator interface {
	CanMulticlassInto(char *Character, className string) (bool, string)
	CanLevelUp(char *Character, className string) (bool, string)
	ValidateLevelUpOptions(char *Character, opts LevelUpOptions) error
}

// HPCalculator defines HP calculation operations
type HPCalculator interface {
	RollHP(hitDie, conMod int, takeAverage bool) int
	CalculateMaxHP(char *Character, class *Class) int
	CalculateHPForLevel(char *Character, level int) int
}

// BenefitManager defines benefit application and removal operations
type BenefitManager interface {
	ApplyBenefit(char *Character, source BenefitSource, benefitType BenefitType, target string, value int) error
	RemoveBenefit(char *Character, source BenefitSource) error
	GetBenefits(char *Character, benefitType BenefitType) []GrantedBenefit
}

// FeatureGranter defines feature granting operations
type FeatureGranter interface {
	GrantFeature(char *Character, feature FeatureDefinition, source string) error
	RemoveFeature(char *Character, featureName string, source string) error
	HasFeature(char *Character, featureName string) bool
}

// SpellcastingManager defines spellcasting operations
type SpellcastingManager interface {
	InitializeSpellcasting(char *Character, class *Class) error
	UpdateSpellSlots(char *Character) error
	CanCastSpell(char *Character, spell Spell) bool
}

// ProficiencyChecker defines proficiency checking operations
type ProficiencyChecker interface {
	HasArmorProficiency(char *Character, armorType string) bool
	HasWeaponProficiency(char *Character, weaponType string) bool
	HasSkillProficiency(char *Character, skillName string) bool
	HasToolProficiency(char *Character, toolName string) bool
}

// DefaultClassRepository is the current implementation
type DefaultClassRepository struct{}

// GetClassByName retrieves a class by name (delegates to existing function)
func (r *DefaultClassRepository) GetClassByName(name string) *Class {
	return GetClassByName(name)
}

// GetAllClasses retrieves all classes (delegates to existing function)
func (r *DefaultClassRepository) GetAllClasses() []Class {
	return GetAllClasses()
}

// LoadClasses loads classes from a directory (delegates to existing function)
func (r *DefaultClassRepository) LoadClasses(dirPath string) error {
	_, err := LoadClassesFromJSON(dirPath)
	return err
}

// DefaultCharacterValidator is the current implementation
type DefaultCharacterValidator struct{}

// CanMulticlassInto checks multiclass prerequisites (delegates to existing function)
func (v *DefaultCharacterValidator) CanMulticlassInto(char *Character, className string) (bool, string) {
	return CanMulticlassInto(char, className)
}

// CanLevelUp checks if character can level up (delegates to existing function)
func (v *DefaultCharacterValidator) CanLevelUp(char *Character, className string) (bool, string) {
	return CanLevelUp(char, className)
}

// ValidateLevelUpOptions validates level up options
func (v *DefaultCharacterValidator) ValidateLevelUpOptions(char *Character, opts LevelUpOptions) error {
	// Basic validation - extend as needed
	if opts.ClassName == "" {
		return ErrInvalidClassName
	}
	canLevel, reason := v.CanLevelUp(char, opts.ClassName)
	if !canLevel {
		return &ValidationError{Field: "className", Message: reason}
	}
	return nil
}

// DefaultHPCalculator is the current implementation
type DefaultHPCalculator struct {
	roller DiceRoller
}

// NewDefaultHPCalculator creates a new HP calculator with a dice roller
func NewDefaultHPCalculator(roller DiceRoller) *DefaultHPCalculator {
	return &DefaultHPCalculator{roller: roller}
}

// RollHP rolls HP for a level (delegates to existing function)
func (calc *DefaultHPCalculator) RollHP(hitDie, conMod int, takeAverage bool) int {
	return RollHPWithDiceRoller(hitDie, conMod, takeAverage, calc.roller)
}

// CalculateMaxHP calculates maximum HP (delegates to existing function)
func (calc *DefaultHPCalculator) CalculateMaxHP(char *Character, class *Class) int {
	return CalculateMaxHP(char, class)
}

// CalculateHPForLevel calculates HP for a specific level
func (calc *DefaultHPCalculator) CalculateHPForLevel(char *Character, level int) int {
	// Simplified implementation - extend as needed
	return char.MaxHP
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// Common errors
var (
	ErrInvalidClassName = &ValidationError{Field: "className", Message: "class name cannot be empty"}
	ErrInvalidLevel     = &ValidationError{Field: "level", Message: "level must be positive"}
)
