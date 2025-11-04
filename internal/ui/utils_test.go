// internal/ui/utils_test.go
package ui

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func TestGetWeaponMasteryCount(t *testing.T) {
	char := models.NewCharacter()

	// Test with no weapon mastery feature
	count := getWeaponMasteryCount(char)
	if count != 0 {
		t.Errorf("Expected 0 for character with no weapon mastery, got %d", count)
	}

	// Test with weapon mastery feature
	char.Features = *models.NewFeatureList()
	char.Features.Features = append(char.Features.Features, models.Feature{
		Name: "Weapon Mastery",
		Mechanics: map[string]interface{}{
			"weapons_mastered": float64(4),
		},
	})

	count = getWeaponMasteryCount(char)
	if count != 4 {
		t.Errorf("Expected 4 weapons mastered, got %d", count)
	}
}

func TestGetExpertiseCount(t *testing.T) {
	char := models.NewCharacter()

	// Test with no expertise feature
	count := getExpertiseCount(char)
	if count != 0 {
		t.Errorf("Expected 0 for character with no expertise, got %d", count)
	}

	// Test with expertise feature
	char.Features = *models.NewFeatureList()
	char.Features.Features = append(char.Features.Features, models.Feature{
		Name: "Expertise",
		Mechanics: map[string]interface{}{
			"expertise_count": float64(2),
		},
	})

	count = getExpertiseCount(char)
	if count != 2 {
		t.Errorf("Expected 2 expertise, got %d", count)
	}
}

func TestGetExpertiseCountRogue(t *testing.T) {
	char := models.NewCharacter()
	char.Class = "Rogue"
	char.Level = 1

	// Test rogue at level 1 (should get 2 expertise)
	count := getExpertiseCount(char)
	if count != 2 {
		t.Errorf("Expected 2 expertise for level 1 rogue, got %d", count)
	}
}
