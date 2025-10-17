// internal/models/species.go
package models

import (
	"encoding/json"
	"os"
	"strings"
)

// SpeciesTrait represents a special trait or ability of a species
type SpeciesTrait struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SpeciesInfo contains all information about a D&D 5e 2024 species
type SpeciesInfo struct {
	Name        string         `json:"name"`
	Size        string         `json:"size"`
	Speed       int            `json:"speed"`
	Traits      []SpeciesTrait `json:"traits"`
	Languages   []string       `json:"languages"`
	Resistances []string       `json:"resistances"`
	Darkvision  int            `json:"darkvision"` // Range in feet, 0 if none
	Description string         `json:"description"`
}

// SpeciesData represents the root structure of the species JSON file
type SpeciesData struct {
	Species []SpeciesInfo `json:"species"`
}

var cachedSpecies []SpeciesInfo

// LoadSpeciesFromJSON loads species data from the JSON file
func LoadSpeciesFromJSON(filepath string) ([]SpeciesInfo, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data SpeciesData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return data.Species, nil
}

// GetAllSpecies returns all available species in D&D 5e 2024
func GetAllSpecies() []SpeciesInfo {
	// Return cached species if already loaded
	if len(cachedSpecies) > 0 {
		return cachedSpecies
	}

	// Try to load from JSON file
	species, err := LoadSpeciesFromJSON("data/species.json")
	if err == nil {
		cachedSpecies = species
		return cachedSpecies
	}

	// Fallback to hardcoded species if JSON file is not available
	cachedSpecies = []SpeciesInfo{
		{
			Name:        "Aasimar",
			Size:        "Medium",
			Speed:       30,
			Description: "Aasimar bear a divine light within them, descended from celestial beings.",
			Traits: []SpeciesTrait{
				{Name: "Celestial Resistance", Description: "You have resistance to necrotic damage and radiant damage."},
				{Name: "Healing Hands", Description: "You can touch a creature and restore hit points equal to your level. Once per long rest."},
				{Name: "Light Bearer", Description: "You know the Light cantrip."},
			},
			Languages:   []string{"Common", "Celestial"},
			Resistances: []string{"Necrotic", "Radiant"},
			Darkvision:  60,
		},
		{
			Name:        "Dragonborn",
			Size:        "Medium",
			Speed:       30,
			Description: "Dragonborn are draconic humanoids with breath weapons and draconic ancestry.",
			Traits: []SpeciesTrait{
				{Name: "Draconic Ancestry", Description: "Choose a dragon type, gaining resistance to its damage type."},
				{Name: "Breath Weapon", Description: "Exhale destructive energy based on your draconic ancestry."},
				{Name: "Damage Resistance", Description: "Resistance to the damage type of your draconic ancestry."},
			},
			Languages:   []string{"Common", "Draconic"},
			Resistances: []string{"Varies (by dragon type)"},
			Darkvision:  60,
		},
		{
			Name:        "Dwarf",
			Size:        "Medium",
			Speed:       25,
			Description: "Dwarves are solid and enduring, known for their craftsmanship and valor.",
			Traits: []SpeciesTrait{
				{Name: "Dwarven Resilience", Description: "Advantage on saves against poison, resistance to poison damage."},
				{Name: "Dwarven Toughness", Description: "Your hit point maximum increases by 1, and it increases by 1 again whenever you gain a level."},
				{Name: "Stonecunning", Description: "Proficiency in History checks related to the origin of stonework."},
			},
			Languages:   []string{"Common", "Dwarvish"},
			Resistances: []string{"Poison"},
			Darkvision:  60,
		},
		{
			Name:        "Elf (Drow)",
			Size:        "Medium",
			Speed:       30,
			Description: "Drow are elves adapted to the Underdark, with innate magical abilities and superior darkvision.",
			Traits: []SpeciesTrait{
				{Name: "Fey Ancestry", Description: "Advantage on saves against being charmed, and magic can't put you to sleep."},
				{Name: "Keen Senses", Description: "Proficiency in Perception."},
				{Name: "Trance", Description: "You don't sleep but meditate for 4 hours instead of sleeping for 8."},
				{Name: "Drow Magic", Description: "You know the Dancing Lights cantrip. At 3rd level, you can cast Faerie Fire. At 5th level, you can cast Darkness."},
				{Name: "Sunlight Sensitivity", Description: "Disadvantage on attack rolls and Perception checks in direct sunlight."},
			},
			Languages:   []string{"Common", "Elven"},
			Resistances: []string{},
			Darkvision:  120,
		},
		{
			Name:        "Elf (High Elf)",
			Size:        "Medium",
			Speed:       30,
			Description: "High elves are scholarly and magical, with a keen mind and mastery of wizardry.",
			Traits: []SpeciesTrait{
				{Name: "Fey Ancestry", Description: "Advantage on saves against being charmed, and magic can't put you to sleep."},
				{Name: "Keen Senses", Description: "Proficiency in Perception."},
				{Name: "Trance", Description: "You don't sleep but meditate for 4 hours instead of sleeping for 8."},
				{Name: "Cantrip", Description: "You know one cantrip from the Wizard spell list. Intelligence is your spellcasting ability for it."},
				{Name: "Extra Language", Description: "You can speak, read, and write one extra language of your choice."},
			},
			Languages:   []string{"Common", "Elven", "One additional language"},
			Resistances: []string{},
			Darkvision:  60,
		},
		{
			Name:        "Elf (Wood Elf)",
			Size:        "Medium",
			Speed:       35,
			Description: "Wood elves are swift and stealthy, at home in the wilderness and attuned to nature.",
			Traits: []SpeciesTrait{
				{Name: "Fey Ancestry", Description: "Advantage on saves against being charmed, and magic can't put you to sleep."},
				{Name: "Keen Senses", Description: "Proficiency in Perception."},
				{Name: "Trance", Description: "You don't sleep but meditate for 4 hours instead of sleeping for 8."},
				{Name: "Fleet of Foot", Description: "Your walking speed increases to 35 feet."},
				{Name: "Mask of the Wild", Description: "You can attempt to hide even when only lightly obscured by foliage, rain, snow, mist, or other natural phenomena."},
			},
			Languages:   []string{"Common", "Elven"},
			Resistances: []string{},
			Darkvision:  60,
		},
		{
			Name:        "Dwarf",
			Size:        "Medium",
			Speed:       25,
			Description: "Dwarves are solid and enduring, known for their craftsmanship and valor.",
			Traits: []SpeciesTrait{
				{Name: "Dwarven Resilience", Description: "Advantage on saves against poison, resistance to poison damage."},
				{Name: "Dwarven Toughness", Description: "Your hit point maximum increases by 1, and it increases by 1 again whenever you gain a level."},
				{Name: "Stonecunning", Description: "Proficiency in History checks related to the origin of stonework."},
			},
			Languages: []string{"Common", "Dwarvish"},
		},
		{
			Name:        "Gnome",
			Size:        "Small",
			Speed:       25,
			Description: "Gnomes are small, inventive, and curious, with a love of tinkering and discovery.",
			Traits: []SpeciesTrait{
				{Name: "Gnome Cunning", Description: "Advantage on Intelligence, Wisdom, and Charisma saves against magic."},
			},
			Languages:   []string{"Common", "Gnomish"},
			Resistances: []string{},
			Darkvision:  60,
		},
		{
			Name:        "Goliath",
			Size:        "Medium",
			Speed:       30,
			Description: "Goliaths are massive, powerful beings native to mountain peaks and high altitudes.",
			Traits: []SpeciesTrait{
				{Name: "Large Form", Description: "You have advantage on saving throws against being pushed or knocked prone."},
				{Name: "Mountain Born", Description: "You have resistance to cold damage and are acclimated to high altitude."},
				{Name: "Powerful Build", Description: "Count as one size larger for carrying capacity and push/drag/lift."},
				{Name: "Stone's Endurance", Description: "You can reduce damage taken by 1d12 + Constitution modifier. Once per long rest."},
			},
			Languages:   []string{"Common", "Giant"},
			Resistances: []string{"Cold"},
			Darkvision:  0,
		},
		{
			Name:        "Halfling",
			Size:        "Small",
			Speed:       25,
			Description: "Halflings are a cheerful and lucky people who prefer comfort and community.",
			Traits: []SpeciesTrait{
				{Name: "Brave", Description: "Advantage on saves against being frightened."},
				{Name: "Halfling Nimbleness", Description: "You can move through the space of any creature larger than you."},
				{Name: "Luck", Description: "When you roll a 1 on d20, you can reroll and must use the new roll."},
				{Name: "Naturally Stealthy", Description: "You can attempt to hide even when obscured only by a larger creature."},
			},
			Languages:   []string{"Common", "Halfling"},
			Resistances: []string{},
			Darkvision:  0,
		},
		{
			Name:        "Human",
			Size:        "Medium",
			Speed:       30,
			Description: "Adaptable and ambitious, humans are the most diverse of the common species.",
			Traits: []SpeciesTrait{
				{Name: "Resourceful", Description: "You gain Inspiration whenever you finish a Long Rest."},
				{Name: "Skillful", Description: "You gain proficiency in one skill of your choice."},
				{Name: "Versatile", Description: "You gain an Origin feat."},
			},
			Languages:   []string{"Common", "One additional language of your choice"},
			Resistances: []string{},
			Darkvision:  0,
		},
		{
			Name:        "Orc",
			Size:        "Medium",
			Speed:       30,
			Description: "Orcs are fierce warriors with a strong sense of honor and strength.",
			Traits: []SpeciesTrait{
				{Name: "Adrenaline Rush", Description: "You can take the Dash action as a bonus action."},
				{Name: "Powerful Build", Description: "Count as one size larger for carrying capacity and push/drag/lift."},
				{Name: "Relentless Endurance", Description: "When reduced to 0 HP, drop to 1 HP instead. Once per long rest."},
			},
			Languages:   []string{"Common", "Orc"},
			Resistances: []string{},
			Darkvision:  60,
		},
		{
			Name:        "Tiefling (Abyssal)",
			Size:        "Medium",
			Speed:       30,
			Description: "Abyssal tieflings are touched by the chaotic energies of the Abyss, embodying primal and destructive magic.",
			Traits: []SpeciesTrait{
				{Name: "Abyssal Legacy", Description: "You know the Dancing Lights cantrip. At 3rd level, you can cast Burning Hands. At 5th level, you can cast Alter Self."},
				{Name: "Abyssal Resistance", Description: "Resistance to poison damage."},
			},
			Languages:   []string{"Common", "Abyssal"},
			Resistances: []string{"Poison"},
			Darkvision:  60,
		},
		{
			Name:        "Tiefling (Chthonic)",
			Size:        "Medium",
			Speed:       30,
			Description: "Chthonic tieflings are connected to the Lower Planes, with powers over death and darkness.",
			Traits: []SpeciesTrait{
				{Name: "Chthonic Legacy", Description: "You know the Chill Touch cantrip. At 3rd level, you can cast False Life. At 5th level, you can cast Ray of Enfeeblement."},
				{Name: "Necrotic Resistance", Description: "Resistance to necrotic damage."},
			},
			Languages:   []string{"Common", "Infernal"},
			Resistances: []string{"Necrotic"},
			Darkvision:  60,
		},
		{
			Name:        "Tiefling (Infernal)",
			Size:        "Medium",
			Speed:       30,
			Description: "Infernal tieflings bear the classic marks of their heritage from the Nine Hells, with fire and command.",
			Traits: []SpeciesTrait{
				{Name: "Infernal Legacy", Description: "You know the Thaumaturgy cantrip. At 3rd level, you can cast Hellish Rebuke. At 5th level, you can cast Darkness."},
				{Name: "Hellish Resistance", Description: "Resistance to fire damage."},
			},
			Languages:   []string{"Common", "Infernal"},
			Resistances: []string{"Fire"},
			Darkvision:  60,
		},
	}

	return cachedSpecies
}

// GetSpeciesByName returns species info for a given species name
func GetSpeciesByName(name string) *SpeciesInfo {
	for _, species := range GetAllSpecies() {
		if species.Name == name {
			return &species
		}
	}
	return nil
}

// ApplySpeciesToCharacter applies species bonuses and traits to a character
func ApplySpeciesToCharacter(char *Character, speciesName string) {
	species := GetSpeciesByName(speciesName)
	if species == nil {
		return
	}

	// Remove old species skill proficiencies first
	RemoveSpeciesSkillProficiencies(char)

	// Update basic stats
	char.Race = species.Name
	char.Speed = species.Speed

	// Update languages from species
	char.Languages = make([]string, len(species.Languages))
	copy(char.Languages, species.Languages)

	// Update resistances from species
	char.Resistances = make([]string, len(species.Resistances))
	copy(char.Resistances, species.Resistances)

	// Update darkvision from species
	char.Darkvision = species.Darkvision

	// Update species traits
	char.SpeciesTraits = make([]SpeciesTrait, len(species.Traits))
	for i, trait := range species.Traits {
		char.SpeciesTraits[i] = SpeciesTrait{
			Name:        trait.Name,
			Description: trait.Description,
		}
	}

	// Clear old species skills list
	char.SpeciesSkills = []SkillType{}

	// Apply skill proficiencies from species traits
	ApplySpeciesSkillProficiencies(char, species)

	// Note: In D&D 5e 2024, ability score increases are more flexible
	// and often chosen by the player rather than fixed by species.
	// We'll handle this through the character creation/leveling system.
}

// RemoveSpeciesSkillProficiencies removes all skill proficiencies granted by the previous species
func RemoveSpeciesSkillProficiencies(char *Character) {
	for _, skillType := range char.SpeciesSkills {
		char.Skills.SetProficiency(skillType, NotProficient)
	}
	char.SpeciesSkills = []SkillType{}
}

// ApplySpeciesSkillProficiencies applies skill proficiencies from species traits
func ApplySpeciesSkillProficiencies(char *Character, species *SpeciesInfo) {
	for _, trait := range species.Traits {
		traitNameLower := strings.ToLower(trait.Name)
		traitDescLower := strings.ToLower(trait.Description)

		// Check for specific skill proficiencies mentioned in traits
		if strings.Contains(traitNameLower, "keen senses") || strings.Contains(traitDescLower, "proficiency in perception") {
			char.Skills.SetProficiency(Perception, Proficient)
			// Track that this skill came from species
			char.SpeciesSkills = append(char.SpeciesSkills, Perception)
		}
	}
}

// AddSpeciesSkillChoice adds a skill proficiency from a species choice
func AddSpeciesSkillChoice(char *Character, skillType SkillType) {
	char.Skills.SetProficiency(skillType, Proficient)
	// Track that this skill came from species
	char.SpeciesSkills = append(char.SpeciesSkills, skillType)
}

// HasSkillChoice checks if a species grants a skill choice
func HasSkillChoice(species *SpeciesInfo) bool {
	for _, trait := range species.Traits {
		traitNameLower := strings.ToLower(trait.Name)
		traitDescLower := strings.ToLower(trait.Description)
		if (strings.Contains(traitNameLower, "skillful") || strings.Contains(traitDescLower, "proficiency in one skill")) &&
		   strings.Contains(traitDescLower, "choice") {
			return true
		}
	}
	return false
}
