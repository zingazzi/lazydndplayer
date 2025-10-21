// internal/ui/components/cantripselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type CantripSelector struct {
	visible            bool
	availableCantrips  []models.Spell
	selectedCantrips   []string // Names of selected cantrips
	cursor             int
	maxCantrips        int
	className          string
	character          *models.Character
}

func NewCantripSelector(char *models.Character) *CantripSelector {
	return &CantripSelector{
		visible:           false,
		availableCantrips: []models.Spell{},
		selectedCantrips:  []string{},
		cursor:            0,
		character:         char,
	}
}

func (cs *CantripSelector) Show(className string, maxCantrips int) {
	cs.className = className
	cs.maxCantrips = maxCantrips
	cs.visible = true
	cs.cursor = 0

	// Pre-select existing cantrips from character
	cs.selectedCantrips = make([]string, len(cs.character.SpellBook.Cantrips))
	copy(cs.selectedCantrips, cs.character.SpellBook.Cantrips)

	// Load cantrips for this class
	cs.loadCantrips()
}

func (cs *CantripSelector) Hide() {
	cs.visible = false
}

func (cs *CantripSelector) IsVisible() bool {
	return cs.visible
}

func (cs *CantripSelector) loadCantrips() {
	// Load all spells from data file
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		// Fallback to empty list if loading fails
		cs.availableCantrips = []models.Spell{}
		return
	}

	cs.availableCantrips = []models.Spell{}
	classLower := strings.ToLower(cs.className)

	for _, spell := range allSpells {
		// Only cantrips (level 0)
		if spell.Level != 0 {
			continue
		}

		// Check if this class can learn this spell
		for _, spellClass := range spell.Classes {
			if strings.ToLower(spellClass) == classLower {
				cs.availableCantrips = append(cs.availableCantrips, spell)
				break
			}
		}
	}
}

func (cs *CantripSelector) Update(msg tea.Msg) (CantripSelector, tea.Cmd) {
	if !cs.visible {
		return *cs, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			cs.Prev()
		case "down", "j":
			cs.Next()
		case " ":
			cs.ToggleSelection()
		case "enter":
			if len(cs.selectedCantrips) == cs.maxCantrips {
				// Confirmed
				return *cs, nil
			}
		case "esc":
			cs.Hide()
		}
	}

	return *cs, nil
}

func (cs *CantripSelector) Next() {
	if cs.cursor < len(cs.availableCantrips)-1 {
		cs.cursor++
	}
}

func (cs *CantripSelector) Prev() {
	if cs.cursor > 0 {
		cs.cursor--
	}
}

func (cs *CantripSelector) ToggleSelection() {
	if len(cs.availableCantrips) == 0 {
		return
	}

	selected := cs.availableCantrips[cs.cursor]

	// Check if already selected
	for i, name := range cs.selectedCantrips {
		if name == selected.Name {
			// Deselect
			cs.selectedCantrips = append(cs.selectedCantrips[:i], cs.selectedCantrips[i+1:]...)
			return
		}
	}

	// Select if not at max
	if len(cs.selectedCantrips) < cs.maxCantrips {
		cs.selectedCantrips = append(cs.selectedCantrips, selected.Name)
	}
}

func (cs *CantripSelector) IsSelected(spellName string) bool {
	for _, name := range cs.selectedCantrips {
		if name == spellName {
			return true
		}
	}
	return false
}

func (cs *CantripSelector) CanConfirm() bool {
	return len(cs.selectedCantrips) == cs.maxCantrips
}

func (cs *CantripSelector) GetSelectedCantrips() []string {
	return cs.selectedCantrips
}

func (cs *CantripSelector) GetMaxCantrips() int {
	return cs.maxCantrips
}

func (cs *CantripSelector) GetSelectedCount() int {
	return len(cs.selectedCantrips)
}

func (cs *CantripSelector) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Left column - cantrip list
	var leftLines []string
	leftLines = append(leftLines, headerStyle.Render(fmt.Sprintf("Select %d cantrips  (%d/%d)", cs.maxCantrips, len(cs.selectedCantrips), cs.maxCantrips)))
	leftLines = append(leftLines, "")

	// Show cantrips
	for i, cantrip := range cs.availableCantrips {
		cursor := "  "
		if i == cs.cursor {
			cursor = "❯ "
		}

		checkbox := "[ ]"
		style := normalStyle
		if cs.IsSelected(cantrip.Name) {
			checkbox = "[✓]"
			style = selectedStyle
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, cantrip.Name)
		leftLines = append(leftLines, style.Render(line))
	}

	leftColumn := strings.Join(leftLines, "\n")

	// Right column - cantrip description
	var rightLines []string
	if len(cs.availableCantrips) > 0 && cs.cursor < len(cs.availableCantrips) {
		cantrip := cs.availableCantrips[cs.cursor]

		rightLines = append(rightLines, headerStyle.Render(cantrip.Name))
		rightLines = append(rightLines, dimStyle.Render(string(cantrip.School)+" Cantrip"))
		rightLines = append(rightLines, "")

		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Casting Time: %s", cantrip.CastingTime)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Range: %s", cantrip.Range)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Components: %s", cantrip.Components)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Duration: %s", cantrip.Duration)))
		rightLines = append(rightLines, "")

		// Wrap description
		descWords := strings.Fields(cantrip.Description)
		line := ""
		for _, word := range descWords {
			if len(line)+len(word)+1 > 35 {
				rightLines = append(rightLines, normalStyle.Render(line))
				line = word
			} else {
				if line != "" {
					line += " "
				}
				line += word
			}
		}
		if line != "" {
			rightLines = append(rightLines, normalStyle.Render(line))
		}
	}

	rightColumn := strings.Join(rightLines, "\n")

	// Combine columns
	leftBox := lipgloss.NewStyle().Width(35).Render(leftColumn)
	rightBox := lipgloss.NewStyle().Width(40).Padding(0, 1).Render(rightColumn)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// Build final content
	var content strings.Builder
	content.WriteString(titleStyle.Render(fmt.Sprintf("⚡ SELECT CANTRIPS FOR %s", strings.ToUpper(cs.className))))
	content.WriteString("\n\n")
	content.WriteString(columns)
	content.WriteString("\n\n")

	// Help text
	if cs.CanConfirm() {
		content.WriteString(helpStyle.Render("↑/↓: Navigate • Space: Toggle • Enter: Confirm • Esc: Cancel"))
	} else {
		content.WriteString(helpStyle.Render(fmt.Sprintf("↑/↓: Navigate • Space: Toggle • Esc: Cancel (need %d more)", cs.maxCantrips-len(cs.selectedCantrips))))
	}

	// Create popup with larger size
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(82)

	return lipgloss.Place(
		100,
		30,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
