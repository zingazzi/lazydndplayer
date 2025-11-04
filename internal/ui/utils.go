// internal/ui/utils.go
package ui

import (
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// getWeaponMasteryCount returns the number of weapons the character can master
func getWeaponMasteryCount(character *models.Character) int {
	debug.Log("getWeaponMasteryCount: Checking for Weapon Mastery feature")
	debug.Log("getWeaponMasteryCount: Character class=%s, total features=%d", character.Class, len(character.Features.Features))

	for i, feature := range character.Features.Features {
		debug.Log("getWeaponMasteryCount: Feature[%d]='%s'", i, feature.Name)
		if feature.Name == "Weapon Mastery" {
			debug.Log("getWeaponMasteryCount: Found Weapon Mastery feature!")

			// Read weapons_mastered from feature mechanics
			if feature.Mechanics != nil {
				if weaponsMastered, ok := feature.Mechanics["weapons_mastered"].(float64); ok {
					count := int(weaponsMastered)
					debug.Log("getWeaponMasteryCount: Returning %d from feature mechanics", count)
					return count
				}
			}

			// Fallback: if no mechanics data, return 0
			debug.Log("getWeaponMasteryCount: No mechanics data found, returning 0")
			return 0
		}
	}
	debug.Log("getWeaponMasteryCount: Weapon Mastery feature not found, returning 0")
	return 0
}

// getExpertiseCount returns the number of skills the character can have expertise in
func getExpertiseCount(character *models.Character) int {
	debug.Log("getExpertiseCount: Checking for Expertise feature")

	// Check for Expertise feature
	for _, feature := range character.Features.Features {
		if feature.Name == "Expertise" && feature.Mechanics != nil {
			if expertiseCount, ok := feature.Mechanics["expertise_count"].(float64); ok {
				debug.Log("getExpertiseCount: Returning %d from feature mechanics", int(expertiseCount))
				return int(expertiseCount)
			}
		}
	}

	// Default: Rogue gets 2 expertise at level 1, 4 at level 6
	if character.IsRogue() {
		rogueLevel := character.GetRogueLevel()
		if rogueLevel >= 6 {
			debug.Log("getExpertiseCount: Rogue level %d, returning 4", rogueLevel)
			return 4
		} else if rogueLevel >= 1 {
			debug.Log("getExpertiseCount: Rogue level %d, returning 2", rogueLevel)
			return 2
		}
	}

	debug.Log("getExpertiseCount: No expertise feature found, returning 0")
	return 0
}
