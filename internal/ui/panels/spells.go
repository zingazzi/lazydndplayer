// internal/ui/panels/spells.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellsPanel displays character spells
type SpellsPanel struct {
	character     *models.Character
	selectedIndex int
}

// NewSpellsPanel creates a new spells panel
func NewSpellsPanel(char *models.Character) *SpellsPanel {
	return &SpellsPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// View renders the spells panel
func (p *SpellsPanel) View(width, height int) string {
	char := p.character

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	var lines []string
	lines = append(lines, titleStyle.Render("SPELLS"))
	lines = append(lines, "")

	// Spellcasting info
	if char.SpellBook.SpellcastingMod != "" {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Render(fmt.Sprintf("Spellcasting Ability: %s", char.SpellBook.SpellcastingMod)))
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Render(fmt.Sprintf("Spell Save DC: %d  Spell Attack: +%d",
				char.SpellBook.SpellSaveDC, char.SpellBook.SpellAttackBonus)))
		lines = append(lines, "")
	}

	// Spell slots
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Render("SPELL SLOTS"))

	slotLevels := []struct {
		level int
		slot  models.SpellSlot
	}{
		{1, char.SpellBook.Slots.Level1},
		{2, char.SpellBook.Slots.Level2},
		{3, char.SpellBook.Slots.Level3},
		{4, char.SpellBook.Slots.Level4},
		{5, char.SpellBook.Slots.Level5},
		{6, char.SpellBook.Slots.Level6},
		{7, char.SpellBook.Slots.Level7},
		{8, char.SpellBook.Slots.Level8},
		{9, char.SpellBook.Slots.Level9},
	}

	for _, sl := range slotLevels {
		if sl.slot.Maximum > 0 {
			slots := strings.Repeat("●", sl.slot.Current) + strings.Repeat("○", sl.slot.Maximum-sl.slot.Current)
			lines = append(lines, fmt.Sprintf("Level %d: %s (%d/%d)",
				sl.level, slots, sl.slot.Current, sl.slot.Maximum))
		}
	}
	lines = append(lines, "")

	// Spells list
	if len(char.SpellBook.Spells) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No spells learned"))
	} else {
		// Group by level
		spellsByLevel := make(map[int][]models.Spell)
		for _, spell := range char.SpellBook.Spells {
			spellsByLevel[spell.Level] = append(spellsByLevel[spell.Level], spell)
		}

		// Display cantrips first, then by level
		for level := 0; level <= 9; level++ {
			spells := spellsByLevel[level]
			if len(spells) == 0 {
				continue
			}

			levelTitle := "Cantrips"
			if level > 0 {
				levelTitle = fmt.Sprintf("Level %d", level)
			}
			lines = append(lines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true).
				Render(levelTitle))

			for _, spell := range spells {
				prepMarker := " "
				if spell.Prepared || level == 0 {
					prepMarker = "●"
				}
				ritualMarker := ""
				if spell.Ritual {
					ritualMarker = " (R)"
				}

				line := fmt.Sprintf("%s %s%s", prepMarker, spell.Name, ritualMarker)
				lines = append(lines, normalStyle.Render(line))
			}
			lines = append(lines, "")
		}
	}

	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("● = Prepared  (R) = Ritual"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'a' to add spell"))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'r' to rest (restore slots)"))

	content := strings.Join(lines, "\n")

	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// Update handles updates for the spells panel
func (p *SpellsPanel) Update(char *models.Character) {
	p.character = char
}
