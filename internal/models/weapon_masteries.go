package models

// WeaponMasteryDescription describes a weapon mastery property
type WeaponMasteryDescription struct {
	Name        string
	Description string
}

// GetWeaponMasteryDescriptions returns all weapon mastery descriptions
func GetWeaponMasteryDescriptions() map[string]WeaponMasteryDescription {
	return map[string]WeaponMasteryDescription{
		"Cleave": {
			Name:        "Cleave",
			Description: "If you hit a creature with a melee attack using this weapon, you can make an attack roll with the weapon against a second creature within 5 feet of the first that is also within your reach. On a hit, the second creature takes the weapon's damage, but don't add your ability modifier to that damage unless that modifier is negative. You can make this extra attack only once per turn.",
		},
		"Graze": {
			Name:        "Graze",
			Description: "If your attack roll with this weapon misses a creature, you can deal damage to that creature equal to the ability modifier you used to make the attack roll. This damage is the same type dealt by the weapon, and the damage can be increased only by increasing the ability modifier.",
		},
		"Nick": {
			Name:        "Nick",
			Description: "When you make the extra attack of the Light property, you can make it as part of the Attack action instead of as a Bonus Action. You can make this extra attack only once per turn.",
		},
		"Push": {
			Name:        "Push",
			Description: "If you hit a creature with this weapon, you can push the creature up to 10 feet straight away from yourself if it is no more than one size larger than you.",
		},
		"Sap": {
			Name:        "Sap",
			Description: "If you hit a creature with this weapon, that creature has Disadvantage on its next attack roll before the start of your next turn.",
		},
		"Slow": {
			Name:        "Slow",
			Description: "If you hit a creature with this weapon and deal damage to it, you can reduce its Speed by 10 feet until the start of your next turn. If the creature is hit more than once by weapons that have this property, the Speed reduction doesn't exceed 10 feet.",
		},
		"Topple": {
			Name:        "Topple",
			Description: "If you hit a creature with this weapon, you can force the creature to make a Constitution saving throw (DC 8 + your Proficiency Bonus + the ability modifier used to make the attack roll). On a failed save, the creature has the Prone condition.",
		},
		"Vex": {
			Name:        "Vex",
			Description: "If you hit a creature with this weapon and deal damage to the creature, you have Advantage on your next attack roll against that creature before the end of your next turn.",
		},
		"Ensnare": {
			Name:        "Ensnare",
			Description: "When you hit a Huge or smaller creature with this weapon, that creature is Restrained until it escapes. The creature can use an action to make a DC 10 Strength (Athletics) or Dexterity (Acrobatics) check, escaping on a success. Dealing 5 slashing damage to the net (AC 10) also frees the target without harming it and destroys the net.",
		},
	}
}

// GetMasteryDescription returns the description for a specific mastery
func GetMasteryDescription(masteryName string) string {
	masteries := GetWeaponMasteryDescriptions()
	if mastery, ok := masteries[masteryName]; ok {
		return mastery.Description
	}
	return ""
}
