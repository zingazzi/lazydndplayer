// internal/ui/components/spellselector.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellSelector handles spell selection UI for species traits
type SpellSelector struct {
	spells        []models.Spell
	selectedIndex int
	viewport      viewport.Model
	visible       bool
	title         string
}

// NewSpellSelector creates a new spell selector
func NewSpellSelector() *SpellSelector {
	return &SpellSelector{
		spells:        []models.Spell{},
		selectedIndex: 0,
		visible:       false,
		title:         "SELECT SPELL",
	}
}

// SetSpells sets the available spells to choose from
func (ss *SpellSelector) SetSpells(spells []models.Spell, title string) {
	ss.spells = spells
	ss.title = title
	ss.selectedIndex = 0
}

// Show displays the spell selector
func (ss *SpellSelector) Show() {
	ss.visible = true
	ss.selectedIndex = 0
}

// Hide hides the spell selector
func (ss *SpellSelector) Hide() {
	ss.visible = false
}

// IsVisible returns whether the selector is visible
func (ss *SpellSelector) IsVisible() bool {
	return ss.visible
}

// Next moves to the next spell
func (ss *SpellSelector) Next() {
	if ss.selectedIndex < len(ss.spells)-1 {
		ss.selectedIndex++
		ss.viewport.LineDown(1)
	}
}

// Prev moves to the previous spell
func (ss *SpellSelector) Prev() {
	if ss.selectedIndex > 0 {
		ss.selectedIndex--
		ss.viewport.LineUp(1)
	}
}

// GetSelectedSpell returns the currently selected spell
func (ss *SpellSelector) GetSelectedSpell() models.Spell {
	if ss.selectedIndex >= 0 && ss.selectedIndex < len(ss.spells) {
		return ss.spells[ss.selectedIndex]
	}
	return models.Spell{}
}

// View renders the spell selector
func (ss *SpellSelector) View(screenWidth, screenHeight int) string {
	if !ss.visible || len(ss.spells) == 0 {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	schoolStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render(ss.title))
	content = append(content, "")

	// Spell list with details
	var spellList []string
	for i, spell := range ss.spells {
		spellName := spell.Name
		school := string(spell.School)
		components := spell.GetComponentsString()

		var spellLine string
		if i == ss.selectedIndex {
			spellLine = selectedStyle.Render("→ " + spellName)
			spellList = append(spellList, spellLine)
			// Add details for selected spell
			spellList = append(spellList, schoolStyle.Render("  "+school+" • "+components))
			// Wrap description
			wrapped := wrapTextForSpellSelector(spell.Description, 60)
			for _, line := range wrapped {
				spellList = append(spellList, descStyle.Render("  "+line))
			}
			spellList = append(spellList, "")
		} else {
			spellLine = normalStyle.Render("  " + spellName)
			spellList = append(spellList, spellLine)
		}
	}

	// Create viewport if needed
	listHeight := 22
	if ss.viewport.Width == 0 {
		ss.viewport = viewport.New(70, listHeight)
		ss.viewport.Style = lipgloss.NewStyle()
	}

	ss.viewport.SetContent(strings.Join(spellList, "\n"))
	content = append(content, ss.viewport.View())
	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// Update handles viewport updates
func (ss *SpellSelector) Update(msg tea.Msg) {
	var cmd tea.Cmd
	ss.viewport, cmd = ss.viewport.Update(msg)
	_ = cmd
}

// wrapTextForSpellSelector wraps text to a specified width
func wrapTextForSpellSelector(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}
