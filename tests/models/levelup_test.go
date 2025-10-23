// tests/models/levelup_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func TestRollHP_Average(t *testing.T) {
	tests := []struct {
		name   string
		hitDie int
		conMod int
		want   int
	}{
		{"d6 with +2 Con", 6, 2, 5},   // (6/2)+1+2 = 5
		{"d10 with +3 Con", 10, 3, 8}, // (10/2)+1+3 = 8
		{"d12 with -1 Con", 12, -1, 5}, // (12/2)+1-1 = 5
		{"d8 with +0 Con", 8, 0, 5},   // (8/2)+1+0 = 5
		{"d6 with -5 Con", 6, -5, 1},  // Minimum 1 HP
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.RollHP(tt.hitDie, tt.conMod, true)
			if got != tt.want {
				t.Errorf("RollHP(%d, %d, true) = %d, want %d",
					tt.hitDie, tt.conMod, got, tt.want)
			}
		})
	}
}

func TestRollHP_WithDiceRoller(t *testing.T) {
	// Use seeded roller for deterministic results
	roller := models.NewSeededDiceRoller(42)

	hitDie := 10
	conMod := 2

	// Roll with our seeded roller
	result := models.RollHPWithDiceRoller(hitDie, conMod, false, roller)

	// Verify it's in valid range
	if result < 1+conMod || result > hitDie+conMod {
		t.Errorf("RollHPWithDiceRoller produced invalid result: %d", result)
	}

	// Roll again with same seed should produce same result
	roller2 := models.NewSeededDiceRoller(42)
	result2 := models.RollHPWithDiceRoller(hitDie, conMod, false, roller2)

	if result != result2 {
		t.Errorf("Deterministic rolls differ: %d vs %d", result, result2)
	}
}

func TestRollHP_MinimumOne(t *testing.T) {
	// Even with very negative Con modifier, should always get at least 1 HP
	result := models.RollHP(6, -20, true)

	if result < 1 {
		t.Errorf("RollHP should return minimum 1 HP, got %d", result)
	}
}

func TestCanLevelUp_FirstClass(t *testing.T) {
	char := models.NewCharacter()
	char.Name = "Test Character"

	// Should be able to select any class as first class
	canLevel, reason := models.CanLevelUp(char, "Fighter")

	if !canLevel {
		t.Errorf("Should be able to select first class, but got: %s", reason)
	}
}

func TestCanLevelUp_ExistingClass(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 1},
	}
	char.AbilityScores.Strength = 13

	// Should be able to level up in existing class
	canLevel, reason := models.CanLevelUp(char, "Fighter")

	if !canLevel {
		t.Errorf("Should be able to level up existing class, but got: %s", reason)
	}
}

func TestCanLevelUp_MulticlassPrerequisites(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 1},
	}
	char.AbilityScores.Strength = 10 // Below requirement

	// Fighter requires 13 Strength to multiclass into
	canLevel, _ := models.CanLevelUp(char, "Barbarian")

	// Should fail due to prerequisites
	if canLevel {
		t.Error("Should not be able to multiclass without meeting prerequisites")
	}
}

func TestCanLevelUp_MulticlassWithPrerequisites(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 1},
	}
	char.AbilityScores.Strength = 13
	char.AbilityScores.Wisdom = 13

	// Should be able to multiclass into Druid with proper stats
	canLevel, reason := models.CanLevelUp(char, "Druid")

	if !canLevel {
		t.Errorf("Should be able to multiclass with prerequisites met, but got: %s", reason)
	}
}

