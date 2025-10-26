// internal/models/character_traits.go
package models

// TraitType represents the type of character trait
type TraitType string

const (
	TraitPersonality TraitType = "personality"
	TraitIdeal       TraitType = "ideal"
	TraitBond        TraitType = "bond"
	TraitFlaw        TraitType = "flaw"
)

// PersonalityTraits contains predefined personality trait options
var PersonalityTraits = []string{
	"I idolize a particular hero and constantly refer to their deeds",
	"I can find common ground between the fiercest enemies",
	"I'm confident in my own abilities and do what I can to instill confidence in others",
	"Thinking is for other people. I prefer action",
	"I misuse long words in an attempt to sound smarter",
	"I get bored easily. When am I going to get on with my destiny?",
	"I'm always polite and respectful",
	"I'm haunted by memories of war. I can't get the images of violence out of my mind",
	"[Custom]",
}

// IdealTraits contains predefined ideal trait options
var IdealTraits = []string{
	"Greater Good: My gifts are meant to be shared with all",
	"Freedom: Tyrants must not be allowed to oppress the people",
	"Charity: I always try to help those in need",
	"Sincerity: There's no good pretending to be something I'm not",
	"Power: If I can attain more power, no one will tell me what to do",
	"Might: The strongest are meant to rule",
	"Independence: I must prove that I can handle myself without the coddling of my family",
	"Redemption: There's a spark of good in everyone",
	"[Custom]",
}

// BondTraits contains predefined bond trait options
var BondTraits = []string{
	"I would die to recover an ancient relic that was lost in my homeland",
	"I will someday get revenge on the corrupt temple hierarchy who branded me a heretic",
	"I owe my life to the priest who took me in when my parents died",
	"Everything I do is for the common people",
	"I will do anything to protect the temple where I served",
	"My family, clan, or tribe is the most important thing in my life",
	"Nothing is more important than the other members of my family",
	"Someone I loved died because of a mistake I made. That will never happen again",
	"[Custom]",
}

// FlawTraits contains predefined flaw trait options
var FlawTraits = []string{
	"I judge others harshly, and myself even more severely",
	"I put too much trust in those who wield power within my temple's hierarchy",
	"Once I pick a goal, I become obsessed with it to the detriment of everything else",
	"I can't resist a pretty face",
	"I'm always in debt. I spend my ill-gotten gains on decadent luxuries faster than I bring them in",
	"The tyrant who rules my land will stop at nothing to see me killed",
	"I'm convinced of the significance of my destiny, and blind to my shortcomings",
	"I have a weakness for the vices of the city, especially hard drink",
	"[Custom]",
}

// GetTraitOptions returns the predefined options for a given trait type
func GetTraitOptions(traitType TraitType) []string {
	switch traitType {
	case TraitPersonality:
		return PersonalityTraits
	case TraitIdeal:
		return IdealTraits
	case TraitBond:
		return BondTraits
	case TraitFlaw:
		return FlawTraits
	default:
		return []string{"[Custom]"}
	}
}

// GetTraitTypeName returns a display name for a trait type
func GetTraitTypeName(traitType TraitType) string {
	switch traitType {
	case TraitPersonality:
		return "Personality Trait"
	case TraitIdeal:
		return "Ideal"
	case TraitBond:
		return "Bond"
	case TraitFlaw:
		return "Flaw"
	default:
		return "Trait"
	}
}

