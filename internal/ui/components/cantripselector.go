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
	cs.selectedCantrips = []string{}

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

	var lines []string
	lines = append(lines, titleStyle.Render(fmt.Sprintf("⚡ SELECT CANTRIPS FOR %s", strings.ToUpper(cs.className))))
	lines = append(lines, "")
	lines = append(lines, headerStyle.Render(fmt.Sprintf("Select %d cantrips  (%d/%d selected)", cs.maxCantrips, len(cs.selectedCantrips), cs.maxCantrips)))
	lines = append(lines, "")

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

		// Show school
		school := string(cantrip.School)
		line := fmt.Sprintf("%s%s %s  %s", cursor, checkbox, cantrip.Name, dimStyle.Render(fmt.Sprintf("(%s)", school)))
		lines = append(lines, style.Render(line))
	}

	lines = append(lines, "")

	// Help text
	if cs.CanConfirm() {
		lines = append(lines, helpStyle.Render("↑/↓: Navigate • Space: Toggle • Enter: Confirm • Esc: Cancel"))
	} else {
		lines = append(lines, helpStyle.Render(fmt.Sprintf("↑/↓: Navigate • Space: Toggle • Esc: Cancel (need %d more)", cs.maxCantrips-len(cs.selectedCantrips))))
	}

	content := strings.Join(lines, "\n")

	// Create popup with medium size
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(60)

	return lipgloss.Place(
		80,
		24,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content),
	)
}
