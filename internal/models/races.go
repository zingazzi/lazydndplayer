// internal/models/races.go
package models

// RaceTrait represents a special trait or ability of a race
type RaceTrait struct {
	Name        string
	Description string
}

// RaceInfo contains all information about a D&D 5e 2024 race
type RaceInfo struct {
	Name        string
	Size        string
	Speed       int
	Traits      []RaceTrait
	Languages   []string
	Description string
}

// GetAllRaces returns all available races in D&D 5e 2024
func GetAllRaces() []RaceInfo {
	return []RaceInfo{
		{
			Name:        "Human",
			Size:        "Medium",
			Speed:       30,
			Description: "Adaptable and ambitious, humans are the most diverse of the common races.",
			Traits: []RaceTrait{
				{Name: "Resourceful", Description: "You gain Inspiration whenever you finish a Long Rest."},
				{Name: "Skillful", Description: "You gain proficiency in one skill of your choice."},
				{Name: "Versatile", Description: "You gain an Origin feat."},
			},
			Languages: []string{"Common", "One additional language of your choice"},
		},
		{
			Name:        "Elf",
			Size:        "Medium",
			Speed:       30,
			Description: "Elves are magical people of otherworldly grace, living in harmony with nature.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Elven Lineage", Description: "Choose High Elf, Wood Elf, or Drow subrace."},
				{Name: "Fey Ancestry", Description: "Advantage on saves against being charmed, and magic can't put you to sleep."},
				{Name: "Keen Senses", Description: "Proficiency in Perception."},
				{Name: "Trance", Description: "You don't sleep but meditate for 4 hours instead of sleeping for 8."},
			},
			Languages: []string{"Common", "Elven"},
		},
		{
			Name:        "Dwarf",
			Size:        "Medium",
			Speed:       25,
			Description: "Dwarves are solid and enduring, known for their craftsmanship and valor.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Dwarven Resilience", Description: "Advantage on saves against poison, resistance to poison damage."},
				{Name: "Dwarven Toughness", Description: "Your hit point maximum increases by 1, and it increases by 1 again whenever you gain a level."},
				{Name: "Stonecunning", Description: "Proficiency in History checks related to the origin of stonework."},
			},
			Languages: []string{"Common", "Dwarvish"},
		},
		{
			Name:        "Halfling",
			Size:        "Small",
			Speed:       25,
			Description: "Halflings are a cheerful and lucky people who prefer comfort and community.",
			Traits: []RaceTrait{
				{Name: "Brave", Description: "Advantage on saves against being frightened."},
				{Name: "Halfling Nimbleness", Description: "You can move through the space of any creature larger than you."},
				{Name: "Luck", Description: "When you roll a 1 on d20, you can reroll and must use the new roll."},
				{Name: "Naturally Stealthy", Description: "You can attempt to hide even when obscured only by a larger creature."},
			},
			Languages: []string{"Common", "Halfling"},
		},
		{
			Name:        "Dragonborn",
			Size:        "Medium",
			Speed:       30,
			Description: "Dragonborn are draconic humanoids with breath weapons and draconic ancestry.",
			Traits: []RaceTrait{
				{Name: "Draconic Ancestry", Description: "Choose a dragon type, gaining resistance to its damage type."},
				{Name: "Breath Weapon", Description: "Exhale destructive energy based on your draconic ancestry."},
				{Name: "Damage Resistance", Description: "Resistance to the damage type of your draconic ancestry."},
			},
			Languages: []string{"Common", "Draconic"},
		},
		{
			Name:        "Gnome",
			Size:        "Small",
			Speed:       25,
			Description: "Gnomes are small, inventive, and curious, with a love of tinkering and discovery.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Gnome Cunning", Description: "Advantage on Intelligence, Wisdom, and Charisma saves against magic."},
			},
			Languages: []string{"Common", "Gnomish"},
		},
		{
			Name:        "Half-Elf",
			Size:        "Medium",
			Speed:       30,
			Description: "Half-elves combine human and elven traits, belonging fully to neither race.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Fey Ancestry", Description: "Advantage on saves against being charmed, and magic can't put you to sleep."},
				{Name: "Skill Versatility", Description: "Proficiency in two skills of your choice."},
			},
			Languages: []string{"Common", "Elven", "One additional language"},
		},
		{
			Name:        "Half-Orc",
			Size:        "Medium",
			Speed:       30,
			Description: "Half-orcs combine orcish strength and human versatility.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Adrenaline Rush", Description: "You can take the Dash action as a bonus action. You can use this trait a number of times equal to your proficiency bonus."},
				{Name: "Relentless Endurance", Description: "When reduced to 0 HP but not killed, drop to 1 HP instead. Once per long rest."},
			},
			Languages: []string{"Common", "Orc"},
		},
		{
			Name:        "Tiefling",
			Size:        "Medium",
			Speed:       30,
			Description: "Tieflings bear the marks of their infernal heritage, with innate magical abilities.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Fiendish Legacy", Description: "You know the Thaumaturgy cantrip and gain additional spells as you level."},
				{Name: "Hellish Resistance", Description: "Resistance to fire damage."},
			},
			Languages: []string{"Common", "Infernal"},
		},
		{
			Name:        "Orc",
			Size:        "Medium",
			Speed:       30,
			Description: "Orcs are fierce warriors with a strong sense of honor and strength.",
			Traits: []RaceTrait{
				{Name: "Darkvision", Description: "You can see in dim light within 60 feet as if it were bright light."},
				{Name: "Adrenaline Rush", Description: "You can take the Dash action as a bonus action."},
				{Name: "Powerful Build", Description: "Count as one size larger for carrying capacity and push/drag/lift."},
				{Name: "Relentless Endurance", Description: "When reduced to 0 HP, drop to 1 HP instead. Once per long rest."},
			},
			Languages: []string{"Common", "Orc"},
		},
	}
}

// GetRaceByName returns race info for a given race name
func GetRaceByName(name string) *RaceInfo {
	for _, race := range GetAllRaces() {
		if race.Name == name {
			return &race
		}
	}
	return nil
}

// ApplyRaceToCharacter applies race bonuses and traits to a character
func ApplyRaceToCharacter(char *Character, raceName string) {
	race := GetRaceByName(raceName)
	if race == nil {
		return
	}

	// Update basic stats
	char.Race = race.Name
	char.Speed = race.Speed

	// Note: In D&D 5e 2024, ability score increases are more flexible
	// and often chosen by the player rather than fixed by race.
	// We'll handle this through the character creation/leveling system.
}
