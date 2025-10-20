// internal/models/pending_changes.go
package models

// PendingChanges holds temporary state while user is making changes
// Allows for rollback if they cancel
type PendingChanges struct {
	// Class changes
	OldClass              string         `json:"old_class,omitempty"`
	OldClassSkills        []SkillType    `json:"old_class_skills,omitempty"`
	OldFightingStyle      string         `json:"old_fighting_style,omitempty"`
	OldClassChoices       ClassChoices   `json:"old_class_choices,omitempty"`
	OldArmorProfs         []string       `json:"old_armor_profs,omitempty"`
	OldWeaponProfs        []string       `json:"old_weapon_profs,omitempty"`
	OldSavingThrowProfs   []string       `json:"old_saving_throw_profs,omitempty"`
	OldMaxHP              int            `json:"old_max_hp,omitempty"`
	OldCurrentHP          int            `json:"old_current_hp,omitempty"`

	// Species changes
	OldSpecies            string         `json:"old_species,omitempty"`
	OldSubtype            string         `json:"old_subtype,omitempty"`
	OldSpeciesSkills      []SkillType    `json:"old_species_skills,omitempty"`
	OldSpeciesChoices     SpeciesChoices `json:"old_species_choices,omitempty"`
	OldSpeed              int            `json:"old_speed,omitempty"`
	OldDarkvision         int            `json:"old_darkvision,omitempty"`

	// Origin changes
	OldOrigin             string         `json:"old_origin,omitempty"`
	OldOriginChoices      OriginChoices  `json:"old_origin_choices,omitempty"`

	// Feat changes (for when removing a feat)
	OldFeat               string         `json:"old_feat,omitempty"`
	OldFeatChoice         *FeatChoice    `json:"old_feat_choice,omitempty"`

	// Type of pending change
	ChangeType            string         `json:"change_type"` // "class", "species", "origin", "feat"
}

// NewPendingChanges creates a new empty pending changes structure
func NewPendingChanges() *PendingChanges {
	return &PendingChanges{}
}

// BackupClass saves current class state before making changes
// Also clears the current class data so selectors appear clean
func (pc *PendingChanges) BackupClass(char *Character) {
	pc.ChangeType = "class"
	pc.OldClass = char.Class
	pc.OldClassSkills = make([]SkillType, len(char.ClassSkills))
	copy(pc.OldClassSkills, char.ClassSkills)
	pc.OldFightingStyle = char.FightingStyle
	pc.OldClassChoices = char.Choices.Class
	pc.OldArmorProfs = make([]string, len(char.ArmorProficiencies))
	copy(pc.OldArmorProfs, char.ArmorProficiencies)
	pc.OldWeaponProfs = make([]string, len(char.WeaponProficiencies))
	copy(pc.OldWeaponProfs, char.WeaponProficiencies)
	pc.OldSavingThrowProfs = make([]string, len(char.SavingThrowProficiencies))
	copy(pc.OldSavingThrowProfs, char.SavingThrowProficiencies)
	pc.OldMaxHP = char.MaxHP
	pc.OldCurrentHP = char.CurrentHP

	// IMMEDIATELY clear old class data from character
	// This makes selectors appear clean with no old selections
	for _, skillType := range char.ClassSkills {
		skill := char.Skills.GetSkill(skillType)
		if skill != nil && skill.Proficiency > 0 {
			skill.Proficiency = 0
		}
	}
	char.ClassSkills = []SkillType{}
	char.Class = ""
	char.FightingStyle = ""
	char.ArmorProficiencies = []string{}
	char.WeaponProficiencies = []string{}
	char.SavingThrowProficiencies = []string{}
}

// RestoreClass restores class state from backup (rollback on cancel)
func (pc *PendingChanges) RestoreClass(char *Character) {
	if pc.ChangeType != "class" {
		return
	}

	char.Class = pc.OldClass
	char.ClassSkills = make([]SkillType, len(pc.OldClassSkills))
	copy(char.ClassSkills, pc.OldClassSkills)
	char.FightingStyle = pc.OldFightingStyle
	char.Choices.Class = pc.OldClassChoices
	char.ArmorProficiencies = make([]string, len(pc.OldArmorProfs))
	copy(char.ArmorProficiencies, pc.OldArmorProfs)
	char.WeaponProficiencies = make([]string, len(pc.OldWeaponProfs))
	copy(char.WeaponProficiencies, pc.OldWeaponProfs)
	char.SavingThrowProficiencies = make([]string, len(pc.OldSavingThrowProfs))
	copy(char.SavingThrowProficiencies, pc.OldSavingThrowProfs)
	char.MaxHP = pc.OldMaxHP
	char.CurrentHP = pc.OldCurrentHP

	// Restore skill proficiencies
	for _, skillType := range pc.OldClassSkills {
		skill := char.Skills.GetSkill(skillType)
		if skill != nil {
			skill.Proficiency = 1
		}
	}
}

// Clear clears all pending changes after commit
func (pc *PendingChanges) Clear() {
	*pc = PendingChanges{}
}

// HasPendingChanges returns true if there are pending changes
func (pc *PendingChanges) HasPendingChanges() bool {
	return pc.ChangeType != ""
}
