// internal/ui/panels/dice.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/dice"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// DicePanelMode represents the current mode of the dice panel
type DicePanelMode int

const (
	DiceModeIdle DicePanelMode = iota
	DiceModeInput
	DiceModeHistory
)

// DicePanel displays dice roller
type DicePanel struct {
	character            *models.Character
	input                textinput.Model
	history              *dice.RollHistory
	LastMessage          string
	mode                 DicePanelMode
	historySelectedIndex int
	viewport             viewport.Model
	ready                bool
}

// NewDicePanel creates a new dice panel
func NewDicePanel(char *models.Character) *DicePanel {
	ti := textinput.New()
	ti.Placeholder = "Type dice and press Enter..."
	ti.CharLimit = 50
	ti.Width = 30
	// Don't auto-focus - user will press Enter to activate

	return &DicePanel{
		character:            char,
		input:                ti,
		history:              dice.NewRollHistory(20),
		LastMessage:          "",
		mode:                 DiceModeIdle,
		historySelectedIndex: 0,
		viewport:             viewport.New(0, 0),
		ready:                false,
	}
}

// View renders the dice panel
func (p *DicePanel) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	// Title with mode indicator
	modeIndicator := ""
	switch p.mode {
	case DiceModeInput:
		modeIndicator = " [INPUT]"
	case DiceModeHistory:
		modeIndicator = " [HISTORY]"
	}

	var headerLines []string
	headerLines = append(headerLines, titleStyle.Render("DICE ROLLER"+modeIndicator))
	headerLines = append(headerLines, "")

	// Input field (highlighted when in input mode)
	inputLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	if p.mode == DiceModeInput {
		inputLabelStyle = inputLabelStyle.Bold(true).Foreground(lipgloss.Color("42"))
	}
	headerLines = append(headerLines, inputLabelStyle.Render("Enter dice notation:"))
	headerLines = append(headerLines, p.input.View())
	headerLines = append(headerLines, "")

	// Last message
	if p.LastMessage != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
		headerLines = append(headerLines, messageStyle.Render(p.LastMessage))
		headerLines = append(headerLines, "")
	}

	// Roll history label (highlighted when in history mode)
	historyLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	if p.mode == DiceModeHistory {
		historyLabelStyle = historyLabelStyle.Foreground(lipgloss.Color("42"))
	}
	headerLines = append(headerLines, historyLabelStyle.Render("ROLL HISTORY"))

	// Mode hints
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	var hint string
	switch p.mode {
	case DiceModeIdle:
		hint = hintStyle.Render("[Enter] Input • [h] History • [r] Reroll last")
	case DiceModeInput:
		hint = hintStyle.Render("[Enter] Roll • [Esc] Back")
	case DiceModeHistory:
		hint = hintStyle.Render("[↑/↓] Navigate • [Enter] Reroll • [Esc] Back")
	}

	// Calculate header and footer heights
	headerContent := strings.Join(headerLines, "\n")
	headerHeight := strings.Count(headerContent, "\n") + 1
	hintHeight := 2 // blank line + hint line

	// Available height for history viewport
	viewportHeight := height - headerHeight - hintHeight - 2 // -2 for padding
	if viewportHeight < 3 {
		viewportHeight = 3
	}

	// Initialize or update viewport size
	if !p.ready || p.viewport.Width != width-4 || p.viewport.Height != viewportHeight {
		p.viewport = viewport.New(width-4, viewportHeight)
		p.ready = true
	}

	// Build history content
	var historyLines []string
	recentRolls := p.history.GetRecent(20)
	if len(recentRolls) == 0 {
		historyLines = append(historyLines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No rolls yet"))
	} else {
		for i, roll := range recentRolls {
			rollStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

			// Highlight selected in history mode
			if p.mode == DiceModeHistory && i == p.historySelectedIndex {
				rollStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("205")).
					Bold(true).
					Reverse(true)
			} else {
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
			}

			historyLines = append(historyLines, rollStyle.Render(roll.String()))
		}
	}

	// Set viewport content
	p.viewport.SetContent(strings.Join(historyLines, "\n"))

	// Add scroll indicators
	historyContent := p.viewport.View()
	if p.viewport.TotalLineCount() > p.viewport.Height {
		scrollIndicator := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
			fmt.Sprintf("↑↓ %d/%d", p.viewport.YOffset+1, p.viewport.TotalLineCount()-p.viewport.Height+1),
		)
		historyContent += "\n" + scrollIndicator
	}

	// Combine all sections
	content := headerContent + "\n" + historyContent + "\n\n" + hint

	return content
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

// Roll performs a dice roll (supports comma-separated multiple rolls)
func (p *DicePanel) Roll(expression string) {
	// Check if it's a comma-separated list
	if strings.Contains(expression, ",") {
		results, err := dice.RollMultiple(expression)
		if err != nil {
			p.LastMessage = fmt.Sprintf("Error: %s", err.Error())
			return
		}

		// Add all results to history
		var messages []string
		for _, result := range results {
			p.history.Add(*result)
			messages = append(messages, result.String())
		}
		p.LastMessage = strings.Join(messages, "\n")
		p.input.SetValue("")
		return
	}

	// Single roll (always use Normal, adv/dis is in the notation)
	result, err := dice.Roll(expression, dice.Normal)
	if err != nil {
		p.LastMessage = fmt.Sprintf("Error: %s", err.Error())
		return
	}

	p.history.Add(*result)
	p.LastMessage = result.String()
	p.input.SetValue("")
}

// GetInput returns the current input value
func (p *DicePanel) GetInput() string {
	return p.input.Value()
}

// SetMode changes the panel mode
func (p *DicePanel) SetMode(mode DicePanelMode) {
	p.mode = mode
	if mode == DiceModeInput {
		p.input.Focus()
	} else {
		p.input.Blur()
	}

	// When entering history mode, reset to top and first item
	if mode == DiceModeHistory {
		p.historySelectedIndex = 0
		p.viewport.GotoTop()
	}
}

// GetMode returns the current mode
func (p *DicePanel) GetMode() DicePanelMode {
	return p.mode
}

// HistoryNext moves selection down in history
func (p *DicePanel) HistoryNext() {
	recentRolls := p.history.GetRecent(20)
	if len(recentRolls) > 0 && p.historySelectedIndex < len(recentRolls)-1 {
		p.historySelectedIndex++
		p.viewport.LineDown(1)
	}
}

// HistoryPrev moves selection up in history
func (p *DicePanel) HistoryPrev() {
	if p.historySelectedIndex > 0 {
		p.historySelectedIndex--
		p.viewport.LineUp(1)
	}
}

// RerollLast rerolls the last roll
func (p *DicePanel) RerollLast() {
	recentRolls := p.history.GetRecent(1)
	if len(recentRolls) > 0 {
		p.Roll(recentRolls[0].Expression)
	}
}

// RerollSelected rerolls the selected history item
func (p *DicePanel) RerollSelected() {
	recentRolls := p.history.GetRecent(20)
	if len(recentRolls) > 0 && p.historySelectedIndex < len(recentRolls) {
		p.Roll(recentRolls[p.historySelectedIndex].Expression)
	}
}
