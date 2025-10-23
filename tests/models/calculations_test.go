// tests/models/calculations_test.go
package models_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func TestCalculateAbilityIncrease(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		increase int
		max      int
		want     int
	}{
		{"Normal increase", 10, 2, 20, 12},
		{"At max", 20, 2, 20, 20},
		{"Over max", 18, 5, 20, 20},
		{"Negative increase", 15, -2, 20, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.CalculateAbilityIncrease(tt.current, tt.increase, tt.max)
			if got != tt.want {
				t.Errorf("CalculateAbilityIncrease(%d, %d, %d) = %d, want %d",
					tt.current, tt.increase, tt.max, got, tt.want)
			}
		})
	}
}

func TestCalculateProficiencyBonus(t *testing.T) {
	tests := []struct {
		level int
		want  int
	}{
		{1, 2}, {4, 2},
		{5, 3}, {8, 3},
		{9, 4}, {12, 4},
		{13, 5}, {16, 5},
		{17, 6}, {20, 6},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := models.CalculateProficiencyBonus(tt.level)
			if got != tt.want {
				t.Errorf("CalculateProficiencyBonus(%d) = %d, want %d",
					tt.level, got, tt.want)
			}
		})
	}
}

func TestCalculateAbilityModifier(t *testing.T) {
	tests := []struct {
		score int
		want  int
	}{
		{1, -5}, {8, -1}, {9, -1},
		{10, 0}, {11, 0},
		{12, 1}, {13, 1},
		{20, 5}, {30, 10},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := models.CalculateAbilityModifier(tt.score)
			if got != tt.want {
				t.Errorf("CalculateAbilityModifier(%d) = %d, want %d",
					tt.score, got, tt.want)
			}
		})
	}
}

func TestCalculatePassiveScore(t *testing.T) {
	tests := []struct {
		name         string
		abilityMod   int
		isProficient bool
		profBonus    int
		otherBonuses int
		want         int
	}{
		{"Not proficient", 3, false, 2, 0, 13},
		{"Proficient", 3, true, 2, 0, 15},
		{"With bonuses", 3, true, 2, 5, 20},
		{"Negative modifier", -1, false, 2, 0, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.CalculatePassiveScore(tt.abilityMod, tt.isProficient, tt.profBonus, tt.otherBonuses)
			if got != tt.want {
				t.Errorf("CalculatePassiveScore(%d, %v, %d, %d) = %d, want %d",
					tt.abilityMod, tt.isProficient, tt.profBonus, tt.otherBonuses, got, tt.want)
			}
		})
	}
}

func TestCalculateACWithArmor(t *testing.T) {
	tests := []struct {
		name         string
		armorBase    int
		dexMod       int
		maxDexBonus  int
		otherBonuses int
		want         int
	}{
		{"Leather armor (light)", 11, 3, -1, 0, 14},
		{"Chainmail (heavy)", 16, 3, 0, 0, 16},
		{"Scale mail (medium)", 14, 3, 2, 0, 16},
		{"Scale mail low dex", 14, 1, 2, 0, 15},
		{"With shield", 11, 2, -1, 2, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.CalculateACWithArmor(tt.armorBase, tt.dexMod, tt.maxDexBonus, tt.otherBonuses)
			if got != tt.want {
				t.Errorf("CalculateACWithArmor(%d, %d, %d, %d) = %d, want %d",
					tt.armorBase, tt.dexMod, tt.maxDexBonus, tt.otherBonuses, got, tt.want)
			}
		})
	}
}

func TestCalculateUnarmoredAC(t *testing.T) {
	tests := []struct {
		name      string
		base      int
		dexMod    int
		otherMods []int
		bonuses   int
		want      int
	}{
		{"No armor", 10, 2, nil, 0, 12},
		{"Monk (Wis)", 10, 3, []int{2}, 0, 15},
		{"Barbarian (Con)", 10, 2, []int{3}, 0, 15},
		{"With shield", 10, 2, nil, 2, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.CalculateUnarmoredAC(tt.base, tt.dexMod, tt.otherMods, tt.bonuses)
			if got != tt.want {
				t.Errorf("CalculateUnarmoredAC(%d, %d, %v, %d) = %d, want %d",
					tt.base, tt.dexMod, tt.otherMods, tt.bonuses, got, tt.want)
			}
		})
	}
}

func TestCalculateCarryCapacity(t *testing.T) {
	tests := []struct {
		strength int
		want     float64
	}{
		{10, 150.0},
		{15, 225.0},
		{20, 300.0},
		{8, 120.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := models.CalculateCarryCapacity(tt.strength)
			if got != tt.want {
				t.Errorf("CalculateCarryCapacity(%d) = %.1f, want %.1f",
					tt.strength, got, tt.want)
			}
		})
	}
}

func TestCalculateSpellSaveDC(t *testing.T) {
	tests := []struct {
		profBonus      int
		spellcastingMod int
		want           int
	}{
		{2, 3, 13}, // 8 + 2 + 3
		{4, 5, 17}, // 8 + 4 + 5
		{3, 1, 12}, // 8 + 3 + 1
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := models.CalculateSpellSaveDC(tt.profBonus, tt.spellcastingMod)
			if got != tt.want {
				t.Errorf("CalculateSpellSaveDC(%d, %d) = %d, want %d",
					tt.profBonus, tt.spellcastingMod, got, tt.want)
			}
		})
	}
}

func TestCalculateHPRatio(t *testing.T) {
	tests := []struct {
		name    string
		current int
		max     int
		want    float64
	}{
		{"Full health", 20, 20, 1.0},
		{"Half health", 10, 20, 0.5},
		{"Quarter health", 5, 20, 0.25},
		{"Zero max", 10, 0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.CalculateHPRatio(tt.current, tt.max)
			if got != tt.want {
				t.Errorf("CalculateHPRatio(%d, %d) = %.2f, want %.2f",
					tt.current, tt.max, got, tt.want)
			}
		})
	}
}

func TestApplyHPRatio(t *testing.T) {
	tests := []struct {
		name   string
		newMax int
		ratio  float64
		want   int
	}{
		{"Full ratio", 30, 1.0, 30},
		{"Half ratio", 30, 0.5, 15},
		{"Very low", 30, 0.01, 1}, // Minimum 1
		{"Zero ratio", 30, 0.0, 1}, // Minimum 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.ApplyHPRatio(tt.newMax, tt.ratio)
			if got != tt.want {
				t.Errorf("ApplyHPRatio(%d, %.2f) = %d, want %d",
					tt.newMax, tt.ratio, got, tt.want)
			}
		})
	}
}

