package models

// Maneuver represents a Battle Master maneuver
type Maneuver struct {
	Name        string
	Description string
	Type        string // "Attack", "Reaction", "Bonus Action", "Special"
}

// AllManeuvers contains all available Battle Master maneuvers
var AllManeuvers = []Maneuver{
	{
		Name:        "Disarming Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to attempt to disarm the target, forcing it to drop one item of your choice that it's holding. You add the superiority die to the attack's damage roll, and the target must make a Strength saving throw. On a failed save, it drops the object you choose. The object lands at its feet.",
		Type:        "Attack",
	},
	{
		Name:        "Distracting Strike",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to distract the creature, giving your allies an opening. You add the superiority die to the attack's damage roll. The next attack roll against the target by an attacker other than you has advantage if the attack is made before the start of your next turn.",
		Type:        "Attack",
	},
	{
		Name:        "Feinting Attack",
		Description: "You can expend one superiority die and use a bonus action on your turn to feint, choosing one creature within 5 feet of you as your target. You have advantage on your next attack roll this turn against that creature. If that attack hits, add the superiority die to the attack's damage roll.",
		Type:        "Bonus Action",
	},
	{
		Name:        "Goading Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to attempt to goad the target into attacking you. You add the superiority die to the attack's damage roll, and the target must make a Wisdom saving throw. On a failed save, the target has disadvantage on all attack rolls against targets other than you until the end of your next turn.",
		Type:        "Attack",
	},
	{
		Name:        "Lunging Attack",
		Description: "When you make a melee weapon attack on your turn, you can expend one superiority die to increase your reach for that attack by 5 feet. If you hit, you add the superiority die to the attack's damage roll.",
		Type:        "Attack",
	},
	{
		Name:        "Maneuvering Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to maneuver one of your comrades into a more advantageous position. You add the superiority die to the attack's damage roll, and you choose a friendly creature who can see or hear you. That creature can use its reaction to move up to half its speed without provoking opportunity attacks from the target of your attack.",
		Type:        "Attack",
	},
	{
		Name:        "Menacing Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to attempt to frighten the target. You add the superiority die to the attack's damage roll, and the target must make a Wisdom saving throw. On a failed save, it is frightened of you until the end of your next turn.",
		Type:        "Attack",
	},
	{
		Name:        "Parry",
		Description: "When another creature damages you with a melee attack, you can use your reaction and expend one superiority die to reduce the damage by the number you roll on your superiority die + your Dexterity modifier.",
		Type:        "Reaction",
	},
	{
		Name:        "Precision Attack",
		Description: "When you make a weapon attack roll against a creature, you can expend one superiority die to add it to the roll. You can use this maneuver before or after making the attack roll, but before any effects of the attack are applied.",
		Type:        "Attack",
	},
	{
		Name:        "Pushing Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to attempt to drive the target back. You add the superiority die to the attack's damage roll, and if the target is Large or smaller, it must make a Strength saving throw. On a failed save, you push the target up to 15 feet away from you.",
		Type:        "Attack",
	},
	{
		Name:        "Rally",
		Description: "On your turn, you can use a bonus action and expend one superiority die to bolster the resolve of one of your companions. When you do so, choose a friendly creature who can see or hear you. That creature gains temporary hit points equal to the superiority die roll + your Charisma modifier.",
		Type:        "Bonus Action",
	},
	{
		Name:        "Riposte",
		Description: "When a creature misses you with a melee attack, you can use your reaction and expend one superiority die to make a melee weapon attack against the creature. If you hit, you add the superiority die to the attack's damage roll.",
		Type:        "Reaction",
	},
	{
		Name:        "Sweeping Attack",
		Description: "When you hit a creature with a melee weapon attack, you can expend one superiority die to attempt to damage another creature with the same attack. Choose another creature within 5 feet of the original target and within your reach. If the original attack roll would hit the second creature, it takes damage equal to the number you roll on your superiority die. The damage is of the same type dealt by the original attack.",
		Type:        "Attack",
	},
	{
		Name:        "Trip Attack",
		Description: "When you hit a creature with a weapon attack, you can expend one superiority die to attempt to knock the target down. You add the superiority die to the attack's damage roll, and if the target is Large or smaller, it must make a Strength saving throw. On a failed save, you knock the target prone.",
		Type:        "Attack",
	},
	{
		Name:        "Commander's Strike",
		Description: "When you take the Attack action on your turn, you can forgo one of your attacks and use a bonus action to direct one of your companions to strike. When you do so, choose a friendly creature who can see or hear you and expend one superiority die. That creature can immediately use its reaction to make one weapon attack, adding the superiority die to the attack's damage roll.",
		Type:        "Bonus Action",
	},
	{
		Name:        "Evasive Footwork",
		Description: "When you move, you can expend one superiority die, rolling the die and adding the number rolled to your AC until you stop moving.",
		Type:        "Special",
	},
}

// GetManeuverByName returns a maneuver by name
func GetManeuverByName(name string) *Maneuver {
	for i := range AllManeuvers {
		if AllManeuvers[i].Name == name {
			return &AllManeuvers[i]
		}
	}
	return nil
}

// GetAllManeuverNames returns a list of all maneuver names
func GetAllManeuverNames() []string {
	names := make([]string, len(AllManeuvers))
	for i, m := range AllManeuvers {
		names[i] = m.Name
	}
	return names
}
