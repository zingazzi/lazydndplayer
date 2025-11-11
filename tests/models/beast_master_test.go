package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// TestBeastMasterLevelUp_BeastSelection tests that Beast Master rangers get a beast companion at level 3
func TestBeastMasterLevelUp_BeastSelection(t *testing.T) {
	char := models.NewCharacter()
	char.Name = "Test Ranger"
	char.Classes = []models.ClassLevel{
		{ClassName: "Ranger", Level: 2, Subclass: ""},
	}
	char.AbilityScores.Dexterity = 13
	char.AbilityScores.Wisdom = 13
	char.Level = 2
	char.TotalLevel = 2

	// Level up to Ranger 3
	options := models.LevelUpOptions{
		ClassName:   "Ranger",
		TakeAverage: true,
	}

	result, err := models.LevelUp(char, options)
	if err != nil {
		t.Fatalf("Level up failed: %v", err)
	}

	// Verify level increased
	if result.NewClassLevel != 3 {
		t.Errorf("Expected Ranger level 3, got %d", result.NewClassLevel)
	}

	// Verify subclass selection is required
	if !result.RequiresSubclass {
		t.Error("Expected RequiresSubclass to be true for Ranger level 3")
	}
}

// TestBeastMaster_InitializeCompanion tests companion initialization
func TestBeastMaster_InitializeCompanion(t *testing.T) {
	tests := []struct {
		beastType    models.CompanionType
		rangerLevel  int
		expectedHP   int
		expectedAC   int
		description  string
	}{
		{
			beastType:   models.CompanionTypeLand,
			rangerLevel: 3,
			expectedHP:  20, // 5 + (5 * 3)
			expectedAC:  13,
			description: "Beast of Land at level 3",
		},
		{
			beastType:   models.CompanionTypeSky,
			rangerLevel: 3,
			expectedHP:  19, // 4 + (5 * 3)
			expectedAC:  13,
			description: "Beast of Sky at level 3",
		},
		{
			beastType:   models.CompanionTypeSea,
			rangerLevel: 3,
			expectedHP:  20, // 5 + (5 * 3)
			expectedAC:  13,
			description: "Beast of Sea at level 3",
		},
		{
			beastType:   models.CompanionTypeLand,
			rangerLevel: 5,
			expectedHP:  30, // 5 + (5 * 5)
			expectedAC:  13,
			description: "Beast of Land at level 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			companion := models.InitializeCompanion(tt.beastType, tt.rangerLevel)

			if companion == nil {
				t.Fatal("Expected companion to be initialized, got nil")
			}

			if companion.Type != tt.beastType {
				t.Errorf("Expected type %s, got %s", tt.beastType, companion.Type)
			}

			if companion.MaxHP != tt.expectedHP {
				t.Errorf("Expected MaxHP %d, got %d", tt.expectedHP, companion.MaxHP)
			}

			if companion.CurrentHP != tt.expectedHP {
				t.Errorf("Expected CurrentHP %d, got %d", tt.expectedHP, companion.CurrentHP)
			}

			if companion.AC != tt.expectedAC {
				t.Errorf("Expected AC %d, got %d", tt.expectedAC, companion.AC)
			}

			// Verify Primal Bond trait exists
			hasPrimalBond := false
			for _, trait := range companion.Traits {
				if trait == "Primal Bond" {
					hasPrimalBond = true
					break
				}
			}
			if !hasPrimalBond {
				t.Error("Expected companion to have Primal Bond trait")
			}
		})
	}
}

// TestBeastMaster_IsBeastMaster tests the IsBeastMaster method
func TestBeastMaster_IsBeastMaster(t *testing.T) {
	tests := []struct {
		classes      []models.ClassLevel
		expected     bool
		description  string
	}{
		{
			classes: []models.ClassLevel{
				{ClassName: "Ranger", Level: 3, Subclass: "Beast Master"},
			},
			expected:    true,
			description: "Ranger with Beast Master subclass",
		},
		{
			classes: []models.ClassLevel{
				{ClassName: "Ranger", Level: 3, Subclass: "Hunter"},
			},
			expected:    false,
			description: "Ranger with Hunter subclass",
		},
		{
			classes: []models.ClassLevel{
				{ClassName: "Fighter", Level: 1, Subclass: ""},
			},
			expected:    false,
			description: "Non-Ranger class",
		},
		{
			classes: []models.ClassLevel{
				{ClassName: "Ranger", Level: 3, Subclass: "Beast Master"},
				{ClassName: "Fighter", Level: 1, Subclass: ""},
			},
			expected:    true,
			description: "Multiclass with Beast Master Ranger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			char := models.NewCharacter()
			char.Classes = tt.classes

			result := char.IsBeastMaster()
			if result != tt.expected {
				t.Errorf("Expected IsBeastMaster() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestBeastMaster_CompanionHPUpdate tests that companion HP updates when ranger levels up
func TestBeastMaster_CompanionHPUpdate(t *testing.T) {
	char := models.NewCharacter()
	char.Name = "Test Ranger"
	char.Classes = []models.ClassLevel{
		{ClassName: "Ranger", Level: 3, Subclass: "Beast Master"},
	}
	char.AbilityScores.Dexterity = 13
	char.AbilityScores.Wisdom = 13
	char.Level = 3
	char.TotalLevel = 3

	// Initialize companion at level 3
	companion := models.InitializeCompanion(models.CompanionTypeLand, 3)
	char.Companion = companion

	initialHP := companion.MaxHP
	if initialHP != 20 {
		t.Fatalf("Expected initial HP 20, got %d", initialHP)
	}

	// Level up to Ranger 4
	options := models.LevelUpOptions{
		ClassName:   "Ranger",
		TakeAverage: true,
	}

	_, err := models.LevelUp(char, options)
	if err != nil {
		t.Fatalf("Level up failed: %v", err)
	}

	// Verify companion HP was updated
	if char.Companion == nil {
		t.Fatal("Expected companion to still exist after level up")
	}

	expectedHP := 25 // 5 + (5 * 4)
	if char.Companion.MaxHP != expectedHP {
		t.Errorf("Expected companion MaxHP %d after level up, got %d", expectedHP, char.Companion.MaxHP)
	}

	// Verify current HP was also updated (should be at max if it was at max before)
	if char.Companion.CurrentHP != expectedHP {
		t.Errorf("Expected companion CurrentHP %d after level up, got %d", expectedHP, char.Companion.CurrentHP)
	}
}

// TestBeastMaster_CompanionMechanics tests companion mechanics calculations
func TestBeastMaster_CompanionMechanics(t *testing.T) {
	char := models.NewCharacter()
	char.Classes = []models.ClassLevel{
		{ClassName: "Ranger", Level: 3, Subclass: "Beast Master"},
	}
	char.AbilityScores.Wisdom = 16 // +3 modifier
	char.Level = 3
	char.TotalLevel = 3

	companion := models.InitializeCompanion(models.CompanionTypeLand, 3)
	char.Companion = companion

	mechanics := models.NewCompanionMechanics(char)

	// Test attack bonus calculation (proficiency + WIS mod)
	// At level 3, proficiency is +2, WIS mod is +3, so attack bonus should be +5
	attackBonus := mechanics.CalculateCompanionAttackBonus()
	expectedAttackBonus := 5 // +2 proficiency + +3 WIS
	if attackBonus != expectedAttackBonus {
		t.Errorf("Expected attack bonus %d, got %d", expectedAttackBonus, attackBonus)
	}

	// Test damage bonus calculation (base + WIS mod)
	// Beast of Land: 1d8+2 + WIS mod, so damage bonus should be 2 + 3 = 5
	damageBonus := mechanics.CalculateCompanionDamageBonus()
	expectedDamageBonus := 5 // 2 base + 3 WIS
	if damageBonus != expectedDamageBonus {
		t.Errorf("Expected damage bonus %d, got %d", expectedDamageBonus, damageBonus)
	}

	// Test ability check bonus (ability mod + proficiency from Primal Bond)
	// STR 12 = +1 mod, proficiency +2, so should be +3
	abilityCheckBonus := mechanics.GetCompanionAbilityCheckBonus(models.Strength)
	expectedAbilityCheckBonus := 3 // +1 STR mod + +2 proficiency
	if abilityCheckBonus != expectedAbilityCheckBonus {
		t.Errorf("Expected ability check bonus %d, got %d", expectedAbilityCheckBonus, abilityCheckBonus)
	}
}
