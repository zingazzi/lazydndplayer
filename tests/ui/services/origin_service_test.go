// tests/ui/services/origin_service_test.go
package services_test

import (
	"strings"
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/services"
)

func TestNewOriginService(t *testing.T) {
	service := services.NewOriginService()
	if service == nil {
		t.Fatal("NewOriginService() returned nil")
	}
}

func TestOriginServiceRequiresAbilityChoice(t *testing.T) {
	service := services.NewOriginService()

	// Create origin with ability choice
	originWithChoice := &models.Origin{
		Name: "Test Origin",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount:  1,
			Choices: []string{"Strength", "Dexterity"},
		},
	}

	// Create origin without ability choice (fixed ability)
	originWithoutChoice := &models.Origin{
		Name: "Simple Origin",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount: 1,
			Ability: "Strength",
		},
	}

	if !service.RequiresAbilityChoice(originWithChoice) {
		t.Error("Expected origin with choice to require ability choice")
	}

	if service.RequiresAbilityChoice(originWithoutChoice) {
		t.Error("Expected origin without choice to not require ability choice")
	}
}

func TestOriginServiceGetAbilityChoices(t *testing.T) {
	service := services.NewOriginService()

	origin := &models.Origin{
		Name: "Test Origin",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount:  1,
			Choices: []string{"Intelligence", "Wisdom", "Charisma"},
		},
	}

	choices := service.GetAbilityChoices(origin)
	if len(choices) != 3 {
		t.Errorf("Expected 3 ability choices, got %d", len(choices))
	}

	expectedChoices := []string{"Intelligence", "Wisdom", "Charisma"}
	for _, expected := range expectedChoices {
		found := false
		for _, choice := range choices {
			if choice == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected choice %s not found in %v", expected, choices)
		}
	}
}

func TestOriginServiceApplyOrigin(t *testing.T) {
	service := services.NewOriginService()

	char := models.NewCharacter()
	char.Origin = "Old Origin"

	origin := &models.Origin{
		Name: "New Origin",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount: 1,
			Ability: "Strength",
		},
	}

	message := service.ApplyOrigin(char, origin, "")
	if !strings.Contains(message, "New Origin") {
		t.Errorf("Expected message to contain origin name, got %s", message)
	}

	if char.Origin != "New Origin" {
		t.Errorf("Expected character origin to be 'New Origin', got %s", char.Origin)
	}
}

func TestOriginServiceApplyOriginWithChoice(t *testing.T) {
	service := services.NewOriginService()

	char := models.NewCharacter()

	origin := &models.Origin{
		Name: "Chosen Origin",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount:  1,
			Choices: []string{"Strength", "Dexterity"},
		},
	}

	message := service.ApplyOrigin(char, origin, "Strength")
	if !strings.Contains(message, "Chosen Origin") {
		t.Errorf("Expected message to contain origin name, got %s", message)
	}
	if !strings.Contains(message, "+1 Strength") {
		t.Errorf("Expected message to contain ability increase, got %s", message)
	}

	if char.Origin != "Chosen Origin" {
		t.Errorf("Expected character origin to be 'Chosen Origin', got %s", char.Origin)
	}
}
