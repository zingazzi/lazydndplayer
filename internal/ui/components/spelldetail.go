// internal/ui/components/spelldetail.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellDetailPopup displays detailed information about a spell
type SpellDetailPopup struct {
	visible bool
	spell   *models.Spell
}

// NewSpellDetailPopup creates a new spell detail popup
func NewSpellDetailPopup() *SpellDetailPopup {
	return &SpellDetailPopup{
		visible: false,
	}
}

// Show displays the popup with the given spell
func (sdp *SpellDetailPopup) Show(spell models.Spell) {
	sdp.spell = &spell
	sdp.visible = true
}

// Hide hides the popup
func (sdp *SpellDetailPopup) Hide() {
	sdp.visible = false
}

// IsVisible returns whether the popup is visible
func (sdp *SpellDetailPopup) IsVisible() bool {
	return sdp.visible
}

// View renders the spell detail popup
func (sdp *SpellDetailPopup) View() string {
	if !sdp.visible || sdp.spell == nil {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var content strings.Builder

	// Title
	levelStr := "Cantrip"
	if sdp.spell.Level > 0 {
		levelStr = fmt.Sprintf("Level %d", sdp.spell.Level)
	}
	content.WriteString(titleStyle.Render(sdp.spell.Name))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render(fmt.Sprintf("%s • %s", levelStr, sdp.spell.School)))
	content.WriteString("\n\n")

	// Casting details
	content.WriteString(labelStyle.Render("Casting Time: "))
	content.WriteString(valueStyle.Render(sdp.spell.CastingTime))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Range: "))
	content.WriteString(valueStyle.Render(sdp.spell.Range))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Components: "))
	content.WriteString(valueStyle.Render(sdp.spell.GetComponentsString()))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Duration: "))
	content.WriteString(valueStyle.Render(sdp.spell.Duration))
	content.WriteString("\n")

	if sdp.spell.Ritual {
		content.WriteString(labelStyle.Render("Ritual: "))
		content.WriteString(valueStyle.Render("Yes"))
		content.WriteString("\n")
	}

	if sdp.spell.Concentration {
		content.WriteString(labelStyle.Render("Concentration: "))
		content.WriteString(valueStyle.Render("Yes"))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Description
	content.WriteString(labelStyle.Render("Description:"))
	content.WriteString("\n")

	// Wrap description text at 60 characters
	descWords := strings.Fields(sdp.spell.Description)
	line := ""
	for _, word := range descWords {
		if len(line)+len(word)+1 > 60 {
			content.WriteString(valueStyle.Render(line))
			content.WriteString("\n")
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	if line != "" {
		content.WriteString(valueStyle.Render(line))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(dimStyle.Render("Press ESC to close"))

	// Create bordered popup
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(70)

	return lipgloss.Place(
		120, 40,
		lipgloss.Center, lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
