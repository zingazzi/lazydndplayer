// test_dwarven_toughness.go - Test Dwarven Toughness feature
package main

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/leveling"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func main() {
	fmt.Println("=== Testing Dwarven Toughness ===\n")

	// Create a level 1 Dwarf Fighter
	char := models.NewCharacter()
	char.Name = "Thorin"
	char.Class = "Fighter"
	char.Level = 1
	char.MaxHP = 10 // Base HP for level 1
	char.CurrentHP = 10

	fmt.Printf("Character: %s the %s\n", char.Name, char.Class)
	fmt.Printf("Level: %d\n", char.Level)
	fmt.Printf("HP before species: %d/%d\n", char.CurrentHP, char.MaxHP)
	fmt.Printf("Species HP Bonus: %d\n\n", char.SpeciesHPBonus)

	// Apply Dwarf species
	fmt.Println("Applying Dwarf species...")
	models.ApplySpeciesToCharacter(char, "Dwarf")

	fmt.Printf("Species: %s\n", char.Race)
	fmt.Printf("HP after species: %d/%d\n", char.CurrentHP, char.MaxHP)
	fmt.Printf("Species HP Bonus: %d\n", char.SpeciesHPBonus)
	fmt.Printf("Expected: MaxHP should be 11 (10 base + 1 from Dwarven Toughness)\n\n")

	// Level up to 2
	fmt.Println("Leveling up to 2...")
	options := leveling.LevelUpOptions{
		HPIncrease: 6, // d10 roll + CON
	}
	leveling.PerformLevelUp(char, options)

	fmt.Printf("Level: %d\n", char.Level)
	fmt.Printf("HP after level up: %d/%d\n", char.CurrentHP, char.MaxHP)
	fmt.Printf("Species HP Bonus: %d\n", char.SpeciesHPBonus)
	fmt.Printf("Expected: MaxHP should be 18 (10 base + 6 level up + 2 from Dwarven Toughness)\n\n")

	// Level up to 3
	fmt.Println("Leveling up to 3...")
	options = leveling.LevelUpOptions{
		HPIncrease: 7,
	}
	leveling.PerformLevelUp(char, options)

	fmt.Printf("Level: %d\n", char.Level)
	fmt.Printf("HP after level up: %d/%d\n", char.CurrentHP, char.MaxHP)
	fmt.Printf("Species HP Bonus: %d\n", char.SpeciesHPBonus)
	fmt.Printf("Expected: MaxHP should be 26 (10 + 6 + 7 + 3 from Dwarven Toughness)\n\n")

	// Test changing species (should remove bonus)
	fmt.Println("Changing species to Human...")
	models.ApplySpeciesToCharacter(char, "Human")

	fmt.Printf("Species: %s\n", char.Race)
	fmt.Printf("HP after species change: %d/%d\n", char.CurrentHP, char.MaxHP)
	fmt.Printf("Species HP Bonus: %d\n", char.SpeciesHPBonus)
	fmt.Printf("Expected: MaxHP should be 23 (26 - 3 from removing Dwarven Toughness)\n\n")

	// Summary
	if char.MaxHP == 23 && char.SpeciesHPBonus == 0 {
		fmt.Println("✅ SUCCESS: Dwarven Toughness working correctly!")
	} else {
		fmt.Printf("❌ ERROR: Expected MaxHP=23, got %d\n", char.MaxHP)
		fmt.Printf("❌ ERROR: Expected SpeciesHPBonus=0, got %d\n", char.SpeciesHPBonus)
	}
}
