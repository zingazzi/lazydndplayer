// internal/ui/panels/dice.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/dice"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// DicePanel displays dice roller
type DicePanel struct {
	character   *models.Character
	input       textinput.Model
	history     *dice.RollHistory
	rollType    dice.RollType
	LastMessage string
}

// NewDicePanel creates a new dice panel
func NewDicePanel(char *models.Character) *DicePanel {
	ti := textinput.New()
	ti.Placeholder = "e.g., 2d6+3, 1d20, d20"
	ti.CharLimit = 20
	ti.Width = 30

	return &DicePanel{
		character:   char,
		input:       ti,
		history:     dice.NewRollHistory(10),
		rollType:    dice.Normal,
		LastMessage: "",
	}
}

// View renders the dice panel
func (p *DicePanel) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	var lines []string
	lines = append(lines, titleStyle.Render("DICE ROLLER"))
	lines = append(lines, "")

	// Roll type selector
	rollTypeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	lines = append(lines, rollTypeStyle.Render(fmt.Sprintf("Roll Type: %s", p.rollType)))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'n' for normal, 'a' for advantage, 'd' for disadvantage"))
	lines = append(lines, "")

	// Input field
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Render("Enter dice notation:"))
	lines = append(lines, p.input.View())
	lines = append(lines, "")

	// Quick roll buttons with simpler layout
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Render("QUICK ROLLS"))

	quickRolls := []string{
		"d+4  d+6  d+8  d+10  d+12  d+20  d+100",
	}
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'd' then number key to roll"))
	lines = append(lines, strings.Join(quickRolls, "  "))
	lines = append(lines, "")

	// Last message
	if p.LastMessage != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
		lines = append(lines, messageStyle.Render(p.LastMessage))
		lines = append(lines, "")
	}

	// Roll history
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Render("ROLL HISTORY"))

	recentRolls := p.history.GetRecent(5)
	if len(recentRolls) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No rolls yet"))
	} else {
		for _, roll := range recentRolls {
			rollStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

			// Highlight critical hits (nat 20) and fails (nat 1) for d20 rolls
			if len(roll.Rolls) > 0 {
				if roll.Expression == "1d20" || roll.Expression == "d20" {
					if roll.Rolls[0] == 20 {
						rollStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("42")).
							Bold(true)
					} else if roll.Rolls[0] == 1 {
						rollStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("196")).
							Bold(true)
					}
				}
			}

			lines = append(lines, rollStyle.Render(roll.String()))
		}
	}

	content := strings.Join(lines, "\n")

	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// Update handles updates for the dice panel
func (p *DicePanel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return cmd
}

// UpdateCharacter updates the character reference
func (p *DicePanel) UpdateCharacter(char *models.Character) {
	p.character = char
}

// Focus focuses the input
func (p *DicePanel) Focus() {
	p.input.Focus()
}

// Blur blurs the input
func (p *DicePanel) Blur() {
	p.input.Blur()
}

// Roll performs a dice roll
func (p *DicePanel) Roll(expression string) {
	result, err := dice.Roll(expression, p.rollType)
	if err != nil {
		p.LastMessage = fmt.Sprintf("Error: %s", err.Error())
		return
	}

	p.history.Add(*result)
	p.LastMessage = result.String()
	p.input.SetValue("")
}

// RollQuick performs a quick roll
func (p *DicePanel) RollQuick(diceType string) {
	p.Roll(diceType)
}

// SetRollType sets the roll type
func (p *DicePanel) SetRollType(rollType dice.RollType) {
	p.rollType = rollType
}

// GetInput returns the current input value
func (p *DicePanel) GetInput() string {
	return p.input.Value()
}
