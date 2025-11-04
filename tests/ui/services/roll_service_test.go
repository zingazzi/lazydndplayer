// tests/ui/services/roll_service_test.go
package services_test

import (
	"strings"
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/services"
)

func TestNewRollService(t *testing.T) {
	service := services.NewRollService()
	if service == nil {
		t.Fatal("NewRollService() returned nil")
	}
}

func TestCalculateSavingThrowRoll(t *testing.T) {
	service := services.NewRollService()

	char := models.NewCharacter()
	char.AbilityScores.SetBaseScore(models.Strength, 16) // +3 modifier
	char.ProficiencyBonus = 2

	// Test saving throw without proficiency
	expression, message := service.CalculateSavingThrowRoll(char, models.Strength)
	if !strings.Contains(expression, "1d20") {
		t.Errorf("Expected expression to contain '1d20', got %s", expression)
	}
	if !strings.Contains(expression, "+3") {
		t.Errorf("Expected expression to contain '+3' (base modifier), got %s", expression)
	}
	if !strings.Contains(message, "Strength") {
		t.Errorf("Expected message to contain 'Strength', got %s", message)
	}

	// Test saving throw with proficiency
	char.SavingThrowProficiencies = []string{"Strength"}
	expression, message = service.CalculateSavingThrowRoll(char, models.Strength)
	if !strings.Contains(expression, "+5") {
		t.Errorf("Expected expression to contain '+5' (3 base + 2 proficiency), got %s", expression)
	}
	if !strings.Contains(message, "proficient") {
		t.Errorf("Expected message to contain 'proficient', got %s", message)
	}
}

func TestCalculateSavingThrowRollAllAbilities(t *testing.T) {
	service := services.NewRollService()

	char := models.NewCharacter()
	char.AbilityScores.SetBaseScore(models.Strength, 16)
	char.AbilityScores.SetBaseScore(models.Dexterity, 14)
	char.AbilityScores.SetBaseScore(models.Constitution, 12)
	char.AbilityScores.SetBaseScore(models.Intelligence, 10)
	char.AbilityScores.SetBaseScore(models.Wisdom, 8)
	char.AbilityScores.SetBaseScore(models.Charisma, 6)
	char.ProficiencyBonus = 2

	abilities := []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}

	for _, ability := range abilities {
		expression, message := service.CalculateSavingThrowRoll(char, ability)
		if expression == "" {
			t.Errorf("Expected non-empty expression for %v", ability)
		}
		if message == "" {
			t.Errorf("Expected non-empty message for %v", ability)
		}
		if !strings.Contains(expression, "1d20") {
			t.Errorf("Expected expression to contain '1d20' for %v, got %s", ability, expression)
		}
	}
}

func TestCalculateAbilityCheckRoll(t *testing.T) {
	service := services.NewRollService()

	char := models.NewCharacter()
	char.AbilityScores.SetBaseScore(models.Dexterity, 18) // +4 modifier
	char.ProficiencyBonus = 3

	// Ability checks don't include proficiency
	expression, message := service.CalculateAbilityCheckRoll(char, models.Dexterity)
	if !strings.Contains(expression, "1d20") {
		t.Errorf("Expected expression to contain '1d20', got %s", expression)
	}
	if !strings.Contains(expression, "+4") {
		t.Errorf("Expected expression to contain '+4' (base modifier only), got %s", expression)
	}
	if strings.Contains(expression, "+7") {
		t.Errorf("Expected expression NOT to contain '+7' (no proficiency), got %s", expression)
	}
	if !strings.Contains(message, "Dexterity") {
		t.Errorf("Expected message to contain 'Dexterity', got %s", message)
	}
	if !strings.Contains(message, "ability check") {
		t.Errorf("Expected message to contain 'ability check', got %s", message)
	}
}

func TestCalculateAbilityCheckRollNegativeModifier(t *testing.T) {
	service := services.NewRollService()

	char := models.NewCharacter()
	char.AbilityScores.SetBaseScore(models.Strength, 6) // -2 modifier

	expression, message := service.CalculateAbilityCheckRoll(char, models.Strength)
	if !strings.Contains(expression, "-2") {
		t.Errorf("Expected expression to contain '-2' for negative modifier, got %s", expression)
	}
	if message == "" {
		t.Error("Expected non-empty message")
	}
}

func TestCalculateSavingThrowRollNegativeModifier(t *testing.T) {
	service := services.NewRollService()

	char := models.NewCharacter()
	char.AbilityScores.SetBaseScore(models.Charisma, 5) // -3 modifier
	char.ProficiencyBonus = 2

	// Without proficiency
	expression, _ := service.CalculateSavingThrowRoll(char, models.Charisma)
	if !strings.Contains(expression, "-3") {
		t.Errorf("Expected expression to contain '-3' for negative modifier, got %s", expression)
	}

	// With proficiency
	char.SavingThrowProficiencies = []string{"Charisma"}
	expression, _ = service.CalculateSavingThrowRoll(char, models.Charisma)
	if !strings.Contains(expression, "-1") {
		t.Errorf("Expected expression to contain '-1' (-3 + 2 proficiency), got %s", expression)
	}
}
