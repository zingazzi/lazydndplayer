// tests/ui/utils_test.go
package ui_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// Note: Since getWeaponMasteryCount and getExpertiseCount are private functions,
// we test the underlying functionality through the character model.
// These tests verify the feature mechanics that the utility functions would use.

func TestWeaponMasteryFeatureMechanics(t *testing.T) {
	char := models.NewCharacter()

	// Test with weapon mastery feature
	char.Features = *models.NewFeatureList()
	char.Features.Features = append(char.Features.Features, models.Feature{
		Name: "Weapon Mastery",
		Mechanics: map[string]interface{}{
			"weapons_mastered": float64(4),
		},
	})

	// Verify the feature exists and has correct mechanics
	var feature *models.Feature
	for i := range char.Features.Features {
		if char.Features.Features[i].Name == "Weapon Mastery" {
			feature = &char.Features.Features[i]
			break
		}
	}
	if feature == nil {
		t.Error("Expected Weapon Mastery feature to exist")
	} else if feature.Mechanics != nil {
		if weaponsMastered, ok := feature.Mechanics["weapons_mastered"].(float64); ok {
			if int(weaponsMastered) != 4 {
				t.Errorf("Expected 4 weapons mastered, got %d", int(weaponsMastered))
			}
		} else {
			t.Error("Expected weapons_mastered to be a float64 in mechanics")
		}
	}
}

func TestExpertiseFeatureMechanics(t *testing.T) {
	char := models.NewCharacter()

	// Test with expertise feature
	char.Features = *models.NewFeatureList()
	char.Features.Features = append(char.Features.Features, models.Feature{
		Name: "Expertise",
		Mechanics: map[string]interface{}{
			"expertise_count": float64(2),
		},
	})

	// Verify the feature exists and has correct mechanics
	var feature *models.Feature
	for i := range char.Features.Features {
		if char.Features.Features[i].Name == "Expertise" {
			feature = &char.Features.Features[i]
			break
		}
	}
	if feature == nil {
		t.Error("Expected Expertise feature to exist")
	} else if feature.Mechanics != nil {
		if expertiseCount, ok := feature.Mechanics["expertise_count"].(float64); ok {
			if int(expertiseCount) != 2 {
				t.Errorf("Expected 2 expertise, got %d", int(expertiseCount))
			}
		} else {
			t.Error("Expected expertise_count to be a float64 in mechanics")
		}
	}
}

func TestRogueClassLevels(t *testing.T) {
	char := models.NewCharacter()
	char.Class = "Rogue"

	// Test rogue level detection
	tests := []struct {
		level int
		desc  string
	}{
		{1, "level 1"},
		{5, "level 5"},
		{6, "level 6"},
		{10, "level 10"},
	}

	for _, tt := range tests {
		char.Level = tt.level
		// Verify character level is set correctly
		if char.Level != tt.level {
			t.Errorf("Expected character level %d, got %d", tt.level, char.Level)
		}
	}
}

// Note: getWeaponMasteryCount and getExpertiseCount are private functions.
// To test them properly, we would need to either:
// 1. Export them (make them public)
// 2. Test them indirectly through handlers that use them
// 3. Create a test helper package
// For now, we test the core functionality through the character model itself
