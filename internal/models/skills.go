// internal/models/skills.go
package models

// SkillType represents a D&D 5e skill
type SkillType string

const (
	Acrobatics     SkillType = "Acrobatics"
	AnimalHandling SkillType = "Animal Handling"
	Arcana         SkillType = "Arcana"
	Athletics      SkillType = "Athletics"
	Deception      SkillType = "Deception"
	History        SkillType = "History"
	Insight        SkillType = "Insight"
	Intimidation   SkillType = "Intimidation"
	Investigation  SkillType = "Investigation"
	Medicine       SkillType = "Medicine"
	Nature         SkillType = "Nature"
	Perception     SkillType = "Perception"
	Performance    SkillType = "Performance"
	Persuasion     SkillType = "Persuasion"
	Religion       SkillType = "Religion"
	SleightOfHand  SkillType = "Sleight of Hand"
	Stealth        SkillType = "Stealth"
	Survival       SkillType = "Survival"
)

// ProficiencyLevel represents the level of proficiency in a skill
type ProficiencyLevel int

const (
	NotProficient ProficiencyLevel = 0
	Proficient    ProficiencyLevel = 1
	Expertise     ProficiencyLevel = 2 // Double proficiency
)

// Skill represents a single skill with its proficiency level
type Skill struct {
	Name        SkillType        `json:"name"`
	Ability     AbilityType      `json:"ability"`
	Proficiency ProficiencyLevel `json:"proficiency"`
}

// Skills holds all character skills
type Skills struct {
	List []Skill `json:"skills"`
}

// NewDefaultSkills creates a new Skills struct with all 18 D&D 5e skills
func NewDefaultSkills() Skills {
	return Skills{
		List: []Skill{
			{Name: Acrobatics, Ability: Dexterity, Proficiency: NotProficient},
			{Name: AnimalHandling, Ability: Wisdom, Proficiency: NotProficient},
			{Name: Arcana, Ability: Intelligence, Proficiency: NotProficient},
			{Name: Athletics, Ability: Strength, Proficiency: NotProficient},
			{Name: Deception, Ability: Charisma, Proficiency: NotProficient},
			{Name: History, Ability: Intelligence, Proficiency: NotProficient},
			{Name: Insight, Ability: Wisdom, Proficiency: NotProficient},
			{Name: Intimidation, Ability: Charisma, Proficiency: NotProficient},
			{Name: Investigation, Ability: Intelligence, Proficiency: NotProficient},
			{Name: Medicine, Ability: Wisdom, Proficiency: NotProficient},
			{Name: Nature, Ability: Intelligence, Proficiency: NotProficient},
			{Name: Perception, Ability: Wisdom, Proficiency: NotProficient},
			{Name: Performance, Ability: Charisma, Proficiency: NotProficient},
			{Name: Persuasion, Ability: Charisma, Proficiency: NotProficient},
			{Name: Religion, Ability: Intelligence, Proficiency: NotProficient},
			{Name: SleightOfHand, Ability: Dexterity, Proficiency: NotProficient},
			{Name: Stealth, Ability: Dexterity, Proficiency: NotProficient},
			{Name: Survival, Ability: Wisdom, Proficiency: NotProficient},
		},
	}
}

// GetSkill returns a skill by name
func (s *Skills) GetSkill(name SkillType) *Skill {
	for i := range s.List {
		if s.List[i].Name == name {
			return &s.List[i]
		}
	}
	return nil
}

// SetProficiency sets the proficiency level for a skill
func (s *Skills) SetProficiency(name SkillType, level ProficiencyLevel) {
	for i := range s.List {
		if s.List[i].Name == name {
			s.List[i].Proficiency = level
			break
		}
	}
}

// CalculateBonus calculates the total bonus for a skill
func (s *Skill) CalculateBonus(abilityMod, proficiencyBonus int) int {
	bonus := abilityMod
	switch s.Proficiency {
	case Proficient:
		bonus += proficiencyBonus
	case Expertise:
		bonus += proficiencyBonus * 2
	}
	return bonus
}
