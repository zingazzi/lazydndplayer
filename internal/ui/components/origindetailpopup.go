// internal/ui/components/origindetailpopup.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// OriginDetailPopup displays full origin information
type OriginDetailPopup struct {
	visible bool
	origin  *models.Origin
}

// NewOriginDetailPopup creates a new origin detail popup
func NewOriginDetailPopup() *OriginDetailPopup {
	return &OriginDetailPopup{
		visible: false,
	}
}

// Show displays the popup with origin details
func (o *OriginDetailPopup) Show(origin *models.Origin) {
	o.visible = true
	o.origin = origin
}

// Hide closes the popup
func (o *OriginDetailPopup) Hide() {
	o.visible = false
}

// IsVisible returns whether the popup is currently visible
func (o *OriginDetailPopup) IsVisible() bool {
	return o.visible
}

// View renders the popup
func (o *OriginDetailPopup) View(width, height int) string {
	if !o.visible || o.origin == nil {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	var content strings.Builder

	content.WriteString(titleStyle.Render(o.origin.Name))
	content.WriteString("\n\n")

	// Description
	wrappedDesc := wrapDetailText(o.origin.Description, width-10)
	for _, line := range wrappedDesc {
		content.WriteString(descStyle.Render(line))
		content.WriteString("\n")
	}
	content.WriteString("\n")

	// Ability Increases
	if o.origin.AbilityIncreases != nil {
		content.WriteString(sectionTitleStyle.Render("Ability Increases"))
		content.WriteString("\n")
		if len(o.origin.AbilityIncreases.Choices) > 0 {
			content.WriteString(normalStyle.Render(fmt.Sprintf("  +%d to one of: %s",
				o.origin.AbilityIncreases.Amount,
				strings.Join(o.origin.AbilityIncreases.Choices, ", "))))
		} else if o.origin.AbilityIncreases.Ability != "" {
			content.WriteString(normalStyle.Render(fmt.Sprintf("  +%d %s",
				o.origin.AbilityIncreases.Amount,
				o.origin.AbilityIncreases.Ability)))
		}
		content.WriteString("\n\n")
	}

	// Feat
	if o.origin.Feat != "" {
		content.WriteString(sectionTitleStyle.Render("Feat"))
		content.WriteString("\n")
		content.WriteString(normalStyle.Render("  " + o.origin.Feat))
		content.WriteString("\n\n")
	}

	// Skill Proficiencies
	if len(o.origin.SkillProficiencies) > 0 {
		content.WriteString(sectionTitleStyle.Render("Skill Proficiencies"))
		content.WriteString("\n")
		content.WriteString(normalStyle.Render("  " + strings.Join(o.origin.SkillProficiencies, ", ")))
		content.WriteString("\n\n")
	}

	// Tool Proficiencies
	if len(o.origin.ToolProficiencies) > 0 {
		content.WriteString(sectionTitleStyle.Render("Tool Proficiencies"))
		content.WriteString("\n")
		content.WriteString(normalStyle.Render("  " + strings.Join(o.origin.ToolProficiencies, ", ")))
		content.WriteString("\n\n")
	}

	// Equipment
	if len(o.origin.Equipment) > 0 {
		content.WriteString(sectionTitleStyle.Render("Equipment"))
		content.WriteString("\n")
		for _, item := range o.origin.Equipment {
			content.WriteString(normalStyle.Render("  • " + item))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	content.WriteString(hintStyle.Render("Press Esc to close"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(width - 4)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

