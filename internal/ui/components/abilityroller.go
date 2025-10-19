// internal/ui/components/abilityroller.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// RollType represents the type of roll
type RollType int

const (
	RollAbilityCheck RollType = iota
	RollSavingThrow
)

// AbilityRoller is a component for rolling ability checks and saving throws
type AbilityRoller struct {
	visible         bool
	selectedAbility int  // 0-5 for STR, DEX, CON, INT, WIS, CHA
	selectedType    int  // 0 = Ability Check, 1 = Saving Throw
	focusOnAbility  bool // true = selecting ability, false = selecting type
}

// NewAbilityRoller creates a new ability roller
func NewAbilityRoller() *AbilityRoller {
	return &AbilityRoller{
		visible:         false,
		selectedAbility: 0,
		selectedType:    0,
		focusOnAbility:  true,
	}
}

// Show displays the ability roller
func (a *AbilityRoller) Show() {
	a.visible = true
	a.selectedAbility = 0
	a.selectedType = 0
	a.focusOnAbility = true
}

// Hide hides the ability roller
func (a *AbilityRoller) Hide() {
	a.visible = false
}

// IsVisible returns whether the ability roller is visible
func (a *AbilityRoller) IsVisible() bool {
	return a.visible
}

// Next moves to next item in current section
func (a *AbilityRoller) Next() {
	if a.focusOnAbility {
		a.selectedAbility = (a.selectedAbility + 1) % 6
	} else {
		a.selectedType = (a.selectedType + 1) % 2
	}
}

// Prev moves to previous item in current section
func (a *AbilityRoller) Prev() {
	if a.focusOnAbility {
		a.selectedAbility--
		if a.selectedAbility < 0 {
			a.selectedAbility = 5
		}
	} else {
		a.selectedType--
		if a.selectedType < 0 {
			a.selectedType = 1
		}
	}
}

// NextAbility moves to next ability (deprecated, use Next)
func (a *AbilityRoller) NextAbility() {
	if a.focusOnAbility {
		a.selectedAbility = (a.selectedAbility + 1) % 6
	}
}

// PrevAbility moves to previous ability (deprecated, use Prev)
func (a *AbilityRoller) PrevAbility() {
	if a.focusOnAbility {
		a.selectedAbility--
		if a.selectedAbility < 0 {
			a.selectedAbility = 5
		}
	}
}

// ToggleType switches between ability check and saving throw
func (a *AbilityRoller) ToggleType() {
	a.selectedType = (a.selectedType + 1) % 2
}

// SwitchFocus toggles between ability and type selection
func (a *AbilityRoller) SwitchFocus() {
	a.focusOnAbility = !a.focusOnAbility
}

// GetSelectedAbility returns the currently selected ability
func (a *AbilityRoller) GetSelectedAbility() models.AbilityType {
	abilities := []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}
	return abilities[a.selectedAbility]
}

// GetSelectedType returns the currently selected roll type
func (a *AbilityRoller) GetSelectedType() RollType {
	if a.selectedType == 0 {
		return RollAbilityCheck
	}
	return RollSavingThrow
}

// GetRollExpression returns the dice expression for the roll
func (a *AbilityRoller) GetRollExpression(char *models.Character) string {
	ability := a.GetSelectedAbility()
	rollType := a.GetSelectedType()

	modifier := char.AbilityScores.GetModifier(ability)

	if rollType == RollSavingThrow {
		// Add proficiency bonus if proficient (check class proficiencies)
		isProficient := false
		for _, prof := range char.SavingThrowProficiencies {
			if strings.EqualFold(prof, string(ability)) {
				isProficient = true
				break
			}
		}

		if isProficient {
			modifier += char.ProficiencyBonus
		}
	}

	return fmt.Sprintf("1d20%+d", modifier)
}

// GetRollDescription returns a description of what's being rolled
func (a *AbilityRoller) GetRollDescription(char *models.Character) string {
	ability := a.GetSelectedAbility()
	rollType := a.GetSelectedType()

	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength",
		models.Dexterity:    "Dexterity",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom",
		models.Charisma:     "Charisma",
	}

	if rollType == RollAbilityCheck {
		return fmt.Sprintf("%s Check", abilityNames[ability])
	}
	return fmt.Sprintf("%s Saving Throw", abilityNames[ability])
}

// View renders the ability roller
func (a *AbilityRoller) View(width, height int, char *models.Character) string {
	if !a.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	unselectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	unfocusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	var lines []string
	lines = append(lines, titleStyle.Render("ROLL ABILITY"))
	lines = append(lines, "")

	// Ability selection
	var abilityHeader string
	if a.focusOnAbility {
		abilityHeader = focusedStyle.Render("► SELECT ABILITY:")
	} else {
		abilityHeader = unfocusedStyle.Render("  SELECT ABILITY:")
	}
	lines = append(lines, abilityHeader)
	lines = append(lines, "")

	abilities := []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}

	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength    ",
		models.Dexterity:    "Dexterity   ",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom      ",
		models.Charisma:     "Charisma    ",
	}

	for i, ability := range abilities {
		score := char.AbilityScores.GetScore(ability)
		modifier := char.AbilityScores.GetModifier(ability)
		isProficient := char.SavingThrows.IsProficient(ability)

		profMarker := " "
		if isProficient {
			profMarker = "●"
		}

		line := fmt.Sprintf("%s %2d (%+d) %s", abilityNames[ability], score, modifier, profMarker)

		if i == a.selectedAbility && a.focusOnAbility {
			lines = append(lines, selectedStyle.Render(fmt.Sprintf("  ▶ %s", line)))
		} else {
			lines = append(lines, unselectedStyle.Render(fmt.Sprintf("    %s", line)))
		}
	}

	lines = append(lines, "")

	// Type selection
	var typeHeader string
	if !a.focusOnAbility {
		typeHeader = focusedStyle.Render("► SELECT TYPE:")
	} else {
		typeHeader = unfocusedStyle.Render("  SELECT TYPE:")
	}
	lines = append(lines, typeHeader)
	lines = append(lines, "")

	// Calculate what the roll will be
	ability := a.GetSelectedAbility()
	baseModifier := char.AbilityScores.GetModifier(ability)
	isProficient := char.SavingThrows.IsProficient(ability)

	// Ability Check
	checkLine := fmt.Sprintf("Ability Check     - 1d20%+d (modifier only)", baseModifier)
	if a.selectedType == 0 && !a.focusOnAbility {
		lines = append(lines, selectedStyle.Render(fmt.Sprintf("  ▶ %s", checkLine)))
	} else {
		lines = append(lines, unselectedStyle.Render(fmt.Sprintf("    %s", checkLine)))
	}

	// Saving Throw
	saveModifier := baseModifier
	profText := ""
	if isProficient {
		saveModifier += char.ProficiencyBonus
		profText = fmt.Sprintf(" + %d prof", char.ProficiencyBonus)
	}
	saveLine := fmt.Sprintf("Saving Throw      - 1d20%+d (modifier%s)", saveModifier, profText)
	if a.selectedType == 1 && !a.focusOnAbility {
		lines = append(lines, selectedStyle.Render(fmt.Sprintf("  ▶ %s", saveLine)))
	} else {
		lines = append(lines, unselectedStyle.Render(fmt.Sprintf("    %s", saveLine)))
	}

	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render("● = Proficient in saving throw"))
	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render("↑/↓: Navigate  Tab: Switch Section  Enter: Roll  Esc: Cancel"))

	content := strings.Join(lines, "\n")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(width - 20)

	box := boxStyle.Render(content)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}
