// tests/models/feats_load_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// TestLoadAthleteFeat tests that Athlete feat loads correctly with ability choices
func TestLoadAthleteFeat(t *testing.T) {
	feats := models.GetAllFeats()

	var athleteFeat *models.Feat
	for i := range feats {
		if feats[i].Name == "Athlete" {
			athleteFeat = &feats[i]
			break
		}
	}

	if athleteFeat == nil {
		t.Fatal("Athlete feat not found")
	}

	// Verify ability increases structure
	if athleteFeat.AbilityIncreases == nil {
		t.Fatal("AbilityIncreases is nil")
	}

	if len(athleteFeat.AbilityIncreases.Choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(athleteFeat.AbilityIncreases.Choices))
	}

	if athleteFeat.AbilityIncreases.Amount != 1 {
		t.Errorf("Expected amount 1, got %d", athleteFeat.AbilityIncreases.Amount)
	}

	// Verify HasAbilityChoice works
	if !models.HasAbilityChoice(*athleteFeat) {
		t.Error("HasAbilityChoice should return true for Athlete")
	}

	// Verify GetAbilityChoices works
	choices := models.GetAbilityChoices(*athleteFeat)
	if len(choices) != 2 {
		t.Errorf("Expected 2 choices, got %d", len(choices))
	}

	expectedChoices := map[string]bool{
		"Strength":  false,
		"Dexterity": false,
	}

	for _, choice := range choices {
		if _, exists := expectedChoices[choice]; !exists {
			t.Errorf("Unexpected choice: %s", choice)
		}
		expectedChoices[choice] = true
	}

	for choice, found := range expectedChoices {
		if !found {
			t.Errorf("Missing expected choice: %s", choice)
		}
	}
}

// TestLoadActorFeat tests that Actor feat loads correctly without ability choices
func TestLoadActorFeat(t *testing.T) {
	feats := models.GetAllFeats()

	var actorFeat *models.Feat
	for i := range feats {
		if feats[i].Name == "Actor" {
			actorFeat = &feats[i]
			break
		}
	}

	if actorFeat == nil {
		t.Fatal("Actor feat not found")
	}

	// Verify ability increases structure
	if actorFeat.AbilityIncreases == nil {
		t.Fatal("AbilityIncreases is nil")
	}

	if actorFeat.AbilityIncreases.Ability != "Charisma" {
		t.Errorf("Expected Charisma, got %s", actorFeat.AbilityIncreases.Ability)
	}

	if actorFeat.AbilityIncreases.Amount != 1 {
		t.Errorf("Expected amount 1, got %d", actorFeat.AbilityIncreases.Amount)
	}

	// Verify HasAbilityChoice returns false
	if models.HasAbilityChoice(*actorFeat) {
		t.Error("HasAbilityChoice should return false for Actor")
	}
}
