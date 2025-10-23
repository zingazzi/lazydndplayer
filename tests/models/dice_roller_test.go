// tests/models/dice_roller_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func TestStandardDiceRoller_Roll(t *testing.T) {
	roller := models.NewStandardDiceRoller()

	tests := []struct {
		name  string
		sides int
		want  bool // Should return valid result
	}{
		{"Roll d6", 6, true},
		{"Roll d20", 20, true},
		{"Roll d12", 12, true},
		{"Roll d0", 0, false}, // Invalid
		{"Roll negative", -5, false}, // Invalid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roller.Roll(tt.sides)

			if tt.want {
				if result < 1 || result > tt.sides {
					t.Errorf("Roll(%d) = %d, want between 1 and %d", tt.sides, result, tt.sides)
				}
			} else {
				if result != 0 {
					t.Errorf("Roll(%d) = %d, want 0 for invalid input", tt.sides, result)
				}
			}
		})
	}
}

func TestSeededDiceRoller_Deterministic(t *testing.T) {
	// With same seed, rolls should be deterministic
	seed := int64(42)
	roller1 := models.NewSeededDiceRoller(seed)
	roller2 := models.NewSeededDiceRoller(seed)

	for i := 0; i < 10; i++ {
		roll1 := roller1.Roll(20)
		roll2 := roller2.Roll(20)

		if roll1 != roll2 {
			t.Errorf("Seeded rolls differ: roll1=%d, roll2=%d", roll1, roll2)
		}
	}
}

func TestDiceRoller_RollMultiple(t *testing.T) {
	roller := models.NewSeededDiceRoller(123)

	results := roller.RollMultiple(4, 6)

	if len(results) != 4 {
		t.Errorf("RollMultiple(4, 6) returned %d results, want 4", len(results))
	}

	for i, result := range results {
		if result < 1 || result > 6 {
			t.Errorf("Result[%d] = %d, want between 1 and 6", i, result)
		}
	}
}

func TestDefaultDiceRoller(t *testing.T) {
	// Test that default roller works
	defaultRoller := models.GetDefaultDiceRoller()

	result := defaultRoller.Roll(20)
	if result < 1 || result > 20 {
		t.Errorf("Default roller produced invalid result: %d", result)
	}
}

func TestSetDefaultDiceRoller(t *testing.T) {
	// Create a mock roller
	mockRoller := models.NewSeededDiceRoller(999)

	// Set it as default
	original := models.GetDefaultDiceRoller()
	models.SetDefaultDiceRoller(mockRoller)

	// Verify it was set
	current := models.GetDefaultDiceRoller()
	if current != mockRoller {
		t.Error("SetDefaultDiceRoller did not update the default roller")
	}

	// Restore original
	models.SetDefaultDiceRoller(original)
}

