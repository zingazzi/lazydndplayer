// test_features.go - Quick test of species features
package main

import (
	"fmt"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func main() {
	// Create a new character
	char := models.NewCharacter()
	fmt.Printf("✓ Created new character\n")
	fmt.Printf("  Features array: %v (length: %d)\n\n", char.Features.Features, len(char.Features.Features))

	// Apply Aasimar species
	fmt.Println("Applying Aasimar species...")
	models.ApplySpeciesToCharacter(char, "Aasimar")

	fmt.Printf("\n✓ Species applied: %s\n", char.Race)
	fmt.Printf("  Speed: %d ft\n", char.Speed)
	fmt.Printf("  Languages: %v\n", char.Languages)
	fmt.Printf("  Darkvision: %d ft\n", char.Darkvision)
	fmt.Printf("  Resistances: %v\n", char.Resistances)
	fmt.Printf("  Species Traits: %d\n", len(char.SpeciesTraits))

	// Check features
	fmt.Printf("\n✓ Features: %d\n", len(char.Features.Features))
	for i, feature := range char.Features.Features {
		fmt.Printf("  %d. %s (%d/%d uses) - %s\n",
			i+1, feature.Name, feature.CurrentUses, feature.MaxUses, feature.RestType)
		fmt.Printf("     Source: %s\n", feature.Source)
		fmt.Printf("     Description: %s\n", feature.Description)
	}

	if len(char.Features.Features) == 0 {
		fmt.Println("\n❌ ERROR: No features were added!")
		fmt.Println("\nDebugging info:")

		species := models.GetSpeciesByName("Aasimar")
		if species == nil {
			fmt.Println("  - Species 'Aasimar' not found!")
		} else {
			fmt.Printf("  - Species found: %s\n", species.Name)
			fmt.Printf("  - Traits count: %d\n", len(species.Traits))
			for i, trait := range species.Traits {
				fmt.Printf("    %d. %s (is_feature: %v)\n", i+1, trait.Name, trait.IsFeature)
			}
		}
	} else {
		fmt.Println("\n✅ SUCCESS: Features were added!")
	}
}
