// internal/ui/components/spellselector.go
package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellSelector allows selecting a single spell from a list (used for species/origin bonuses)
type SpellSelector struct {
	visible      bool
	spells       []models.Spell
	selectedIdx  int
	viewport     viewport.Model
	title        string
}

// NewSpellSelector creates a new spell selector
func NewSpellSelector() *SpellSelector {
	vp := viewport.New(50, 15)
	return &SpellSelector{
		visible:     false,
		spells:      []models.Spell{},
		selectedIdx: 0,
		viewport:    vp,
		title:       "SELECT SPELL",
	}
}

// SetSpells sets the available spells and title
func (ss *SpellSelector) SetSpells(spells []models.Spell, title string) {
	ss.spells = spells
	ss.title = title
	ss.selectedIdx = 0
}

// Show displays the selector
func (ss *SpellSelector) Show() {
	ss.visible = true
}

// Hide hides the selector
func (ss *SpellSelector) Hide() {
	ss.visible = false
}

// IsVisible returns visibility state
func (ss *SpellSelector) IsVisible() bool {
	return ss.visible
}

// Next moves to next spell
func (ss *SpellSelector) Next() {
	if ss.selectedIdx < len(ss.spells)-1 {
		ss.selectedIdx++
		ss.viewport.ScrollDown(1)
	}
}

// Prev moves to previous spell
func (ss *SpellSelector) Prev() {
	if ss.selectedIdx > 0 {
		ss.selectedIdx--
		ss.viewport.ScrollUp(1)
	}
}

// GetSelectedSpell returns the currently selected spell
func (ss *SpellSelector) GetSelectedSpell() models.Spell {
	if ss.selectedIdx >= 0 && ss.selectedIdx < len(ss.spells) {
		return ss.spells[ss.selectedIdx]
	}
	return models.Spell{}
}

// View renders the spell selector
func (ss *SpellSelector) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var content string
	content += titleStyle.Render(ss.title) + "\n\n"

	if len(ss.spells) == 0 {
		content += dimStyle.Render("No spells available")
	} else {
		for i, spell := range ss.spells {
			cursor := "  "
			style := normalStyle
			if i == ss.selectedIdx {
				cursor = "❯ "
				style = selectedStyle
			}

			schoolInfo := fmt.Sprintf(" (%s)", spell.School)
			line := fmt.Sprintf("%s%s%s", cursor, spell.Name, dimStyle.Render(schoolInfo))
			content += style.Render(line) + "\n"
		}
	}

	content += "\n" + dimStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Cancel")

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(width).
		Height(height)

	return lipgloss.Place(
		width+20,
		height+10,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content),
	)
}
