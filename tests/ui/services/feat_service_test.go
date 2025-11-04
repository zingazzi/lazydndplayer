// tests/ui/services/feat_service_test.go
package services_test

import (
	"strings"
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/services"
)

func TestNewFeatService(t *testing.T) {
	service := services.NewFeatService()
	if service == nil {
		t.Fatal("NewFeatService() returned nil")
	}
}

func TestFeatServiceRequiresAbilityChoice(t *testing.T) {
	service := services.NewFeatService()

	// Create a feat with ability increases
	featWithChoice := &models.Feat{
		Name: "Test Feat",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount:  1,
			Choices: []string{"Strength", "Dexterity"},
		},
	}

	// Create a feat without ability increases
	featWithoutChoice := &models.Feat{
		Name: "Simple Feat",
	}

	if !service.RequiresAbilityChoice(featWithChoice) {
		t.Error("Expected feat with ability increases to require ability choice")
	}

	if service.RequiresAbilityChoice(featWithoutChoice) {
		t.Error("Expected feat without ability increases to not require ability choice")
	}
}

func TestFeatServiceGetAbilityChoices(t *testing.T) {
	service := services.NewFeatService()

	feat := &models.Feat{
		Name: "Test Feat",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Amount:  1,
			Choices: []string{"Strength", "Dexterity", "Constitution"},
		},
	}

	choices := service.GetAbilityChoices(feat)
	if len(choices) != 3 {
		t.Errorf("Expected 3 ability choices, got %d", len(choices))
	}

	expectedChoices := []string{"Strength", "Dexterity", "Constitution"}
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

func TestFeatServiceRemoveFeat(t *testing.T) {
	service := services.NewFeatService()

	char := models.NewCharacter()
	char.Feats = []string{"Test Feat", "Another Feat"}

	feat := &models.Feat{
		Name: "Test Feat",
		Benefits: []string{"+1 HP"},
	}

	message := service.RemoveFeat(char, feat)

	if !strings.Contains(message, "Test Feat") {
		t.Errorf("Expected message to contain feat name, got %s", message)
	}

	// Check that feat was removed from character
	featFound := false
	for _, featName := range char.Feats {
		if featName == "Test Feat" {
			featFound = true
			break
		}
	}
	if featFound {
		t.Error("Expected feat to be removed from character's feat list")
	}

	// Other feat should still be present
	otherFeatFound := false
	for _, featName := range char.Feats {
		if featName == "Another Feat" {
			otherFeatFound = true
			break
		}
	}
	if !otherFeatFound {
		t.Error("Expected other feat to still be present")
	}
}

func TestFeatServiceCancelFeatSelection(t *testing.T) {
	service := services.NewFeatService()

	char := models.NewCharacter()
	char.Feats = []string{"Pending Feat", "Other Feat"}

	feat := &models.Feat{
		Name: "Pending Feat",
	}

	service.CancelFeatSelection(char, feat)

	// Check that feat was removed
	featFound := false
	for _, featName := range char.Feats {
		if featName == "Pending Feat" {
			featFound = true
			break
		}
	}
	if featFound {
		t.Error("Expected feat to be removed when cancelling selection")
	}
}
