// internal/models/choices.go
package models

// CharacterChoices tracks all choices made during character creation/modification
// This allows for easy rollback when changing species, class, origin, or feats
type CharacterChoices struct {
	Species SpeciesChoices  `json:"species"`
	Origin  OriginChoices   `json:"origin"`
	Class   ClassChoices    `json:"class"`
	Feats   []FeatChoice    `json:"feats"`
}

// SpeciesChoices tracks choices made when selecting a species
type SpeciesChoices struct {
	Name              string   `json:"name"`                // e.g., "Elf"
	Subtype           string   `json:"subtype,omitempty"`   // e.g., "High Elf"
	SelectedLanguages []string `json:"selected_languages"`  // Languages chosen from choices
	SelectedSkills    []string `json:"selected_skills"`     // Skills chosen from choices
	AbilityIncreases  []string `json:"ability_increases"`   // Which abilities got increases (e.g., ["Dexterity", "Intelligence"])
}

// OriginChoices tracks choices made when selecting an origin
type OriginChoices struct {
	Name              string   `json:"name"`                // e.g., "Acolyte"
	SelectedFeat      string   `json:"selected_feat"`       // Feat granted by origin (if any)
	SelectedLanguages []string `json:"selected_languages"`  // Languages chosen
	SelectedSkills    []string `json:"selected_skills"`     // Skills chosen
	SelectedTools     []string `json:"selected_tools"`      // Tools chosen
	AbilityIncreases  []string `json:"ability_increases"`   // Which abilities got increases
	FeatChoices       FeatChoiceDetail `json:"feat_choices,omitempty"` // Choices made within the feat
}

// ClassChoices tracks choices made when selecting a class
type ClassChoices struct {
	Name              string              `json:"name"`                // e.g., "Fighter"
	FightingStyle     string              `json:"fighting_style"`      // Fighting style chosen (Fighter, Paladin, Ranger)
	LevelChoices      []LevelChoice       `json:"level_choices"`       // Choices made at each level
}

// LevelChoice tracks choices made at a specific class level
type LevelChoice struct {
	Level             int      `json:"level"`               // Character level
	SelectedSkills    []string `json:"selected_skills"`     // Skills chosen at this level
	SelectedSubclass  string   `json:"selected_subclass"`   // Subclass chosen (if any)
	AbilityIncreases  []string `json:"ability_increases"`   // ASI or which abilities increased
	SelectedFeat      string   `json:"selected_feat"`       // Feat chosen instead of ASI (if any)
	OtherChoices      map[string]string `json:"other_choices,omitempty"` // Other class-specific choices
}

// FeatChoice tracks choices made when selecting a feat
type FeatChoice struct {
	FeatName      string           `json:"feat_name"`       // Name of the feat
	LevelGained   int              `json:"level_gained"`    // Level when feat was gained
	Source        string           `json:"source"`          // "origin", "class_asi", "species", etc.
	Choices       FeatChoiceDetail `json:"choices"`         // Specific choices made within the feat
}

// FeatChoiceDetail tracks the specific choices made within a feat
type FeatChoiceDetail struct {
	SelectedAbility    string   `json:"selected_ability,omitempty"`    // Which ability was increased (if choice)
	SelectedAbilities  []string `json:"selected_abilities,omitempty"`  // Multiple abilities (if multiple increases)
	SelectedSkills     []string `json:"selected_skills,omitempty"`     // Skills chosen
	SelectedLanguages  []string `json:"selected_languages,omitempty"`  // Languages chosen
	SelectedTools      []string `json:"selected_tools,omitempty"`      // Tools chosen
	OtherChoices       map[string]string `json:"other_choices,omitempty"` // Other feat-specific choices
}

// NewCharacterChoices creates a new empty CharacterChoices
func NewCharacterChoices() *CharacterChoices {
	return &CharacterChoices{
		Feats: []FeatChoice{},
		Class: ClassChoices{
			LevelChoices: []LevelChoice{},
		},
	}
}

// RecordSpeciesChoice records a species selection and its choices
func (cc *CharacterChoices) RecordSpeciesChoice(name, subtype string, languages, skills, abilities []string) {
	cc.Species = SpeciesChoices{
		Name:              name,
		Subtype:           subtype,
		SelectedLanguages: languages,
		SelectedSkills:    skills,
		AbilityIncreases:  abilities,
	}
}

// RecordOriginChoice records an origin selection and its choices
func (cc *CharacterChoices) RecordOriginChoice(name string, feat string, languages, skills, tools, abilities []string, featChoices FeatChoiceDetail) {
	cc.Origin = OriginChoices{
		Name:              name,
		SelectedFeat:      feat,
		SelectedLanguages: languages,
		SelectedSkills:    skills,
		SelectedTools:     tools,
		AbilityIncreases:  abilities,
		FeatChoices:       featChoices,
	}
}

// RecordClassChoice records a class selection
func (cc *CharacterChoices) RecordClassChoice(name, fightingStyle string) {
	cc.Class = ClassChoices{
		Name:          name,
		FightingStyle: fightingStyle,
		LevelChoices:  []LevelChoice{},
	}
}

// RecordLevelChoice records choices made at a specific level
func (cc *CharacterChoices) RecordLevelChoice(level int, skills []string, subclass string, abilities []string, feat string, other map[string]string) {
	levelChoice := LevelChoice{
		Level:            level,
		SelectedSkills:   skills,
		SelectedSubclass: subclass,
		AbilityIncreases: abilities,
		SelectedFeat:     feat,
		OtherChoices:     other,
	}

	// Update or append level choice
	found := false
	for i := range cc.Class.LevelChoices {
		if cc.Class.LevelChoices[i].Level == level {
			cc.Class.LevelChoices[i] = levelChoice
			found = true
			break
		}
	}
	if !found {
		cc.Class.LevelChoices = append(cc.Class.LevelChoices, levelChoice)
	}
}

// RecordFeatChoice records a feat selection and its choices
func (cc *CharacterChoices) RecordFeatChoice(featName string, level int, source string, choices FeatChoiceDetail) {
	feat := FeatChoice{
		FeatName:    featName,
		LevelGained: level,
		Source:      source,
		Choices:     choices,
	}
	cc.Feats = append(cc.Feats, feat)
}

// RemoveFeatChoice removes a feat choice by name
func (cc *CharacterChoices) RemoveFeatChoice(featName string) bool {
	for i, feat := range cc.Feats {
		if feat.FeatName == featName {
			cc.Feats = append(cc.Feats[:i], cc.Feats[i+1:]...)
			return true
		}
	}
	return false
}

// GetFeatChoice retrieves a feat choice by name
func (cc *CharacterChoices) GetFeatChoice(featName string) *FeatChoice {
	for i := range cc.Feats {
		if cc.Feats[i].FeatName == featName {
			return &cc.Feats[i]
		}
	}
	return nil
}

// ClearSpeciesChoices clears all species-related choices
func (cc *CharacterChoices) ClearSpeciesChoices() {
	cc.Species = SpeciesChoices{}
}

// ClearOriginChoices clears all origin-related choices
func (cc *CharacterChoices) ClearOriginChoices() {
	cc.Origin = OriginChoices{}
}

// ClearClassChoices clears all class-related choices
func (cc *CharacterChoices) ClearClassChoices() {
	cc.Class = ClassChoices{
		LevelChoices: []LevelChoice{},
	}
}
