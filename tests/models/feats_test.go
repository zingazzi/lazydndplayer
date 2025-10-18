// tests/models/feats_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// TestApplyFeatBenefits_SingleAbility tests applying a feat with a single ability increase
func TestApplyFeatBenefits_SingleAbility(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	initialCharisma := char.AbilityScores.Charisma

	// Actor feat: +1 Charisma
	feat := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}

	err := models.ApplyFeatBenefits(char, feat, "")
	if err != nil {
		t.Fatalf("ApplyFeatBenefits failed: %v", err)
	}

	// Verify ability increased
	if char.AbilityScores.Charisma != initialCharisma+1 {
		t.Errorf("Expected Charisma %d, got %d", initialCharisma+1, char.AbilityScores.Charisma)
	}

	// Verify benefit tracked
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Actor")
	if len(benefits) != 1 {
		t.Errorf("Expected 1 benefit, got %d", len(benefits))
	}

	if len(benefits) > 0 {
		if benefits[0].Type != models.BenefitAbilityScore {
			t.Errorf("Expected BenefitAbilityScore, got %s", benefits[0].Type)
		}
		if benefits[0].Target != "Charisma" {
			t.Errorf("Expected Charisma, got %s", benefits[0].Target)
		}
		if benefits[0].Value != 1 {
			t.Errorf("Expected value 1, got %d", benefits[0].Value)
		}
	}
}

// TestApplyFeatBenefits_MultipleChoice tests applying a feat with ability choices
func TestApplyFeatBenefits_MultipleChoice(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	initialStrength := char.AbilityScores.Strength

	// Athlete feat: +1 Strength or Dexterity
	feat := models.Feat{
		Name: "Athlete",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Choices: []string{"Strength", "Dexterity"},
			Amount:  1,
		},
	}

	// Choose Strength
	err := models.ApplyFeatBenefits(char, feat, "Strength")
	if err != nil {
		t.Fatalf("ApplyFeatBenefits failed: %v", err)
	}

	// Verify ability increased
	if char.AbilityScores.Strength != initialStrength+1 {
		t.Errorf("Expected Strength %d, got %d", initialStrength+1, char.AbilityScores.Strength)
	}

	// Verify benefit tracked with correct ability
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Athlete")
	if len(benefits) != 1 {
		t.Errorf("Expected 1 benefit, got %d", len(benefits))
	}

	if len(benefits) > 0 {
		if benefits[0].Target != "Strength" {
			t.Errorf("Expected Strength, got %s", benefits[0].Target)
		}
	}
}

// TestRemoveFeatBenefits_SingleAbility tests removing a feat with single ability
func TestRemoveFeatBenefits_SingleAbility(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	initialCharisma := char.AbilityScores.Charisma

	// Apply Actor feat
	feat := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}
	models.ApplyFeatBenefits(char, feat, "")

	// Verify it was applied
	if char.AbilityScores.Charisma != initialCharisma+1 {
		t.Fatalf("Feat was not applied correctly")
	}

	// Remove the feat
	err := models.RemoveFeatBenefits(char, feat)
	if err != nil {
		t.Fatalf("RemoveFeatBenefits failed: %v", err)
	}

	// Verify ability restored
	if char.AbilityScores.Charisma != initialCharisma {
		t.Errorf("Expected Charisma restored to %d, got %d", initialCharisma, char.AbilityScores.Charisma)
	}

	// Verify benefit removed from tracker
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Actor")
	if len(benefits) != 0 {
		t.Errorf("Expected 0 benefits after removal, got %d", len(benefits))
	}
}

// TestRemoveFeatBenefits_MultipleChoice tests removing a feat with ability choice
func TestRemoveFeatBenefits_MultipleChoice(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	initialStrength := char.AbilityScores.Strength
	initialDexterity := char.AbilityScores.Dexterity

	// Apply Athlete feat with Strength choice
	feat := models.Feat{
		Name: "Athlete",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Choices: []string{"Strength", "Dexterity"},
			Amount:  1,
		},
	}
	models.ApplyFeatBenefits(char, feat, "Strength")

	// Verify Strength increased, Dexterity unchanged
	if char.AbilityScores.Strength != initialStrength+1 {
		t.Fatalf("Strength not increased")
	}
	if char.AbilityScores.Dexterity != initialDexterity {
		t.Fatalf("Dexterity should not have changed")
	}

	// Remove the feat
	models.RemoveFeatBenefits(char, feat)

	// Verify Strength restored, Dexterity still unchanged
	if char.AbilityScores.Strength != initialStrength {
		t.Errorf("Expected Strength restored to %d, got %d", initialStrength, char.AbilityScores.Strength)
	}
	if char.AbilityScores.Dexterity != initialDexterity {
		t.Errorf("Expected Dexterity unchanged at %d, got %d", initialDexterity, char.AbilityScores.Dexterity)
	}

	// Verify no benefits remain
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Athlete")
	if len(benefits) != 0 {
		t.Errorf("Expected 0 benefits after removal, got %d", len(benefits))
	}
}

// TestApplyFeatBenefits_Tough tests the Tough feat (special HP bonus)
func TestApplyFeatBenefits_Tough(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	char.Level = 5
	char.MaxHP = 40
	char.CurrentHP = 40
	initialHP := char.MaxHP

	feat := models.Feat{
		Name: "Tough",
	}

	models.ApplyFeatBenefits(char, feat, "")

	// Tough gives +2 HP per level
	expectedHP := initialHP + (char.Level * 2)
	if char.MaxHP != expectedHP {
		t.Errorf("Expected MaxHP %d, got %d", expectedHP, char.MaxHP)
	}

	// Verify benefit tracked
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Tough")
	if len(benefits) != 1 {
		t.Errorf("Expected 1 benefit, got %d", len(benefits))
	}
}

// TestRemoveFeatBenefits_Tough tests removing Tough feat
func TestRemoveFeatBenefits_Tough(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	char.Level = 5
	char.MaxHP = 40
	char.CurrentHP = 40
	initialHP := char.MaxHP

	feat := models.Feat{
		Name: "Tough",
	}

	models.ApplyFeatBenefits(char, feat, "")
	models.RemoveFeatBenefits(char, feat)

	// HP should be restored
	if char.MaxHP != initialHP {
		t.Errorf("Expected MaxHP restored to %d, got %d", initialHP, char.MaxHP)
	}

	// Verify no benefits remain
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Tough")
	if len(benefits) != 0 {
		t.Errorf("Expected 0 benefits after removal, got %d", len(benefits))
	}
}

// TestApplyFeatBenefits_Mobile tests the Mobile feat (speed bonus)
func TestApplyFeatBenefits_Mobile(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	char.Speed = 30
	initialSpeed := char.Speed

	feat := models.Feat{
		Name: "Mobile",
	}

	models.ApplyFeatBenefits(char, feat, "")

	// Mobile gives +10 speed
	if char.Speed != initialSpeed+10 {
		t.Errorf("Expected Speed %d, got %d", initialSpeed+10, char.Speed)
	}

	// Verify benefit tracked
	benefits := char.BenefitTracker.GetBenefitsBySource("feat", "Mobile")
	if len(benefits) != 1 {
		t.Errorf("Expected 1 benefit, got %d", len(benefits))
	}
}

// TestMultipleFeats tests adding and removing multiple feats
func TestMultipleFeats(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	initialCharisma := char.AbilityScores.Charisma
	initialStrength := char.AbilityScores.Strength

	// Add Actor (+1 Charisma)
	actor := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}
	models.ApplyFeatBenefits(char, actor, "")

	// Add Athlete (+1 Strength)
	athlete := models.Feat{
		Name: "Athlete",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Choices: []string{"Strength", "Dexterity"},
			Amount:  1,
		},
	}
	models.ApplyFeatBenefits(char, athlete, "Strength")

	// Verify both applied
	if char.AbilityScores.Charisma != initialCharisma+1 {
		t.Errorf("Actor not applied correctly")
	}
	if char.AbilityScores.Strength != initialStrength+1 {
		t.Errorf("Athlete not applied correctly")
	}

	// Verify both tracked
	allBenefits := char.BenefitTracker.Benefits
	if len(allBenefits) != 2 {
		t.Errorf("Expected 2 benefits total, got %d", len(allBenefits))
	}

	// Remove Actor only
	models.RemoveFeatBenefits(char, actor)

	// Verify Actor removed, Athlete remains
	if char.AbilityScores.Charisma != initialCharisma {
		t.Errorf("Actor not removed correctly")
	}
	if char.AbilityScores.Strength != initialStrength+1 {
		t.Errorf("Athlete should still be active")
	}

	// Verify only Athlete benefits remain
	actorBenefits := char.BenefitTracker.GetBenefitsBySource("feat", "Actor")
	athleteBenefits := char.BenefitTracker.GetBenefitsBySource("feat", "Athlete")
	if len(actorBenefits) != 0 {
		t.Errorf("Expected 0 Actor benefits, got %d", len(actorBenefits))
	}
	if len(athleteBenefits) != 1 {
		t.Errorf("Expected 1 Athlete benefit, got %d", len(athleteBenefits))
	}
}

// TestHasAbilityChoice tests the helper function
func TestHasAbilityChoice(t *testing.T) {
	// Feat with choices
	athleteFeat := models.Feat{
		Name: "Athlete",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Choices: []string{"Strength", "Dexterity"},
			Amount:  1,
		},
	}

	if !models.HasAbilityChoice(athleteFeat) {
		t.Error("Athlete should have ability choice")
	}

	// Feat without choices
	actorFeat := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}

	if models.HasAbilityChoice(actorFeat) {
		t.Error("Actor should not have ability choice")
	}

	// Feat with no ability increases
	alertFeat := models.Feat{
		Name: "Alert",
	}

	if models.HasAbilityChoice(alertFeat) {
		t.Error("Alert should not have ability choice")
	}
}

// TestGetAbilityChoices tests getting available choices
func TestGetAbilityChoices(t *testing.T) {
	feat := models.Feat{
		Name: "Athlete",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Choices: []string{"Strength", "Dexterity"},
			Amount:  1,
		},
	}

	choices := models.GetAbilityChoices(feat)
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

// TestAbilityScoreMax tests that ability scores don't exceed 20
func TestAbilityScoreMax(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	char.AbilityScores.Charisma = 20

	feat := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}

	models.ApplyFeatBenefits(char, feat, "")

	// Should cap at 20
	if char.AbilityScores.Charisma != 20 {
		t.Errorf("Charisma should cap at 20, got %d", char.AbilityScores.Charisma)
	}
}

// TestAbilityScoreMin tests that ability scores don't go below 1 on removal
func TestAbilityScoreMin(t *testing.T) {
	char := models.NewCharacter()
	char.BenefitTracker = models.NewBenefitTracker()
	char.AbilityScores.Charisma = 1

	feat := models.Feat{
		Name: "Actor",
		AbilityIncreases: &models.FeatAbilityIncrease{
			Ability: "Charisma",
			Amount:  1,
		},
	}

	models.ApplyFeatBenefits(char, feat, "")
	char.AbilityScores.Charisma = 1 // Manually set to 1
	models.RemoveFeatBenefits(char, feat)

	// Should not go below 1
	if char.AbilityScores.Charisma < 1 {
		t.Errorf("Charisma should not go below 1, got %d", char.AbilityScores.Charisma)
	}
}
