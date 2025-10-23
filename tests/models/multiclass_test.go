// tests/models/multiclass_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func TestCanMulticlassInto_MeetsPrerequisites(t *testing.T) {
	char := models.NewCharacter()
	char.AbilityScores.Strength = 13
	char.AbilityScores.Dexterity = 13

	tests := []struct {
		className string
		canDo     bool
	}{
		{"Fighter", true},     // Needs 13 Str or Dex
		{"Barbarian", true},   // Needs 13 Str
		{"Rogue", true},       // Needs 13 Dex
		{"Wizard", false},     // Needs 13 Int (char has 10)
		{"Cleric", false},     // Needs 13 Wis (char has 10)
	}

	for _, tt := range tests {
		t.Run(tt.className, func(t *testing.T) {
			can, _ := models.CanMulticlassInto(char, tt.className)
			if can != tt.canDo {
				t.Errorf("CanMulticlassInto(%s) = %v, want %v",
					tt.className, can, tt.canDo)
			}
		})
	}
}

func TestGetAvailableClasses_FiltersPrerequisites(t *testing.T) {
	char := models.NewCharacter()
	char.AbilityScores.Strength = 13
	char.AbilityScores.Wisdom = 13
	// All other abilities at 10 (default)

	available := models.GetAvailableClasses(char)

	// Should only include classes with met prerequisites
	hasClasses := make(map[string]bool)
	for _, class := range available {
		hasClasses[class.Name] = true
	}

	// Should have classes that need Str or Wis
	if !hasClasses["Fighter"] {
		t.Error("Should include Fighter (has Str 13)")
	}
	if !hasClasses["Druid"] {
		t.Error("Should include Druid (has Wis 13)")
	}

	// Should NOT have classes that need other abilities
	if hasClasses["Wizard"] {
		t.Error("Should not include Wizard (needs Int 13)")
	}
}

func TestHasClass(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 3},
		{ClassName: "Wizard", Level: 2},
	}

	if !char.HasClass("Fighter") {
		t.Error("HasClass should return true for Fighter")
	}
	if !char.HasClass("Wizard") {
		t.Error("HasClass should return true for Wizard")
	}
	if char.HasClass("Rogue") {
		t.Error("HasClass should return false for Rogue")
	}
}

func TestGetClassLevel(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 3},
		{ClassName: "Wizard", Level: 2},
	}

	if level := char.GetClassLevel("Fighter"); level != 3 {
		t.Errorf("GetClassLevel(Fighter) = %d, want 3", level)
	}
	if level := char.GetClassLevel("Wizard"); level != 2 {
		t.Errorf("GetClassLevel(Wizard) = %d, want 2", level)
	}
	if level := char.GetClassLevel("Rogue"); level != 0 {
		t.Errorf("GetClassLevel(Rogue) = %d, want 0", level)
	}
}

func TestGetClassDisplayString(t *testing.T) {
	tests := []struct {
		name    string
		classes []models.ClassLevel
		want    string
	}{
		{
			name:    "Single class",
			classes: []models.ClassLevel{{ClassName: "Fighter", Level: 3}},
			want:    "Fighter 3",
		},
		{
			name: "Multiclass",
			classes: []models.ClassLevel{
				{ClassName: "Fighter", Level: 3},
				{ClassName: "Wizard", Level: 2},
			},
			want: "Fighter 3 / Wizard 2",
		},
		{
			name: "Triple multiclass",
			classes: []models.ClassLevel{
				{ClassName: "Fighter", Level: 5},
				{ClassName: "Rogue", Level: 3},
				{ClassName: "Wizard", Level: 2},
			},
			want: "Fighter 5 / Rogue 3 / Wizard 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := models.NewCharacter()
			char.Classes = tt.classes

			got := char.GetClassDisplayString()
			if got != tt.want {
				t.Errorf("GetClassDisplayString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalculateTotalLevel(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Fighter", Level: 3},
		{ClassName: "Wizard", Level: 2},
		{ClassName: "Rogue", Level: 1},
	}

	total := char.CalculateTotalLevel()
	want := 6

	if total != want {
		t.Errorf("CalculateTotalLevel() = %d, want %d", total, want)
	}
}

func TestGetPrimaryClass(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Wizard", Level: 2},
		{ClassName: "Fighter", Level: 5},
		{ClassName: "Rogue", Level: 1},
	}

	primary := char.GetPrimaryClass()

	if primary.ClassName != "Fighter" {
		t.Errorf("GetPrimaryClass() = %s, want Fighter (highest level)", primary.ClassName)
	}
	if primary.Level != 5 {
		t.Errorf("GetPrimaryClass().Level = %d, want 5", primary.Level)
	}
}

