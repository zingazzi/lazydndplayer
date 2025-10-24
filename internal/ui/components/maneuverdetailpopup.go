// internal/ui/components/maneuverdetailpopup.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ManeuverDetailPopup displays detailed information about a Battle Master maneuver
type ManeuverDetailPopup struct {
	visible      bool
	maneuverName string
}

// NewManeuverDetailPopup creates a new maneuver detail popup
func NewManeuverDetailPopup() *ManeuverDetailPopup {
	return &ManeuverDetailPopup{
		visible: false,
	}
}

// Show displays the popup with the given maneuver
func (p *ManeuverDetailPopup) Show(maneuverName string) {
	p.maneuverName = maneuverName
	p.visible = true
}

// Hide hides the popup
func (p *ManeuverDetailPopup) Hide() {
	p.visible = false
}

// IsVisible returns whether the popup is visible
func (p *ManeuverDetailPopup) IsVisible() bool {
	return p.visible
}

// View renders the popup
func (p *ManeuverDetailPopup) View(width, height int) string {
	if !p.visible {
		return ""
	}

	// Get maneuver definition
	maneuver := models.GetManeuverByName(p.maneuverName)
	if maneuver == nil {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208")).
		Padding(0, 1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("248")).
		Padding(1, 2)

	// Build content
	var content strings.Builder
	content.WriteString(titleStyle.Render(maneuver.Name))
	content.WriteString("\n\n")

	// Type
	content.WriteString(labelStyle.Render("Type: "))
	content.WriteString(valueStyle.Render(maneuver.Type))
	content.WriteString("\n\n")

	// Description (wrapped)
	wrappedDesc := wrapManeuverDetailText(maneuver.Description, width-8)
	content.WriteString(descStyle.Render(wrappedDesc))

	content.WriteString("\n\n")
	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("  Press ESC to close"))

	// Create popup box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(1, 2).
		Width(width - 4).
		MaxWidth(100)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		boxStyle.Render(content.String()),
	)
}

// wrapManeuverDetailText wraps text to fit within a given width
func wrapManeuverDetailText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
