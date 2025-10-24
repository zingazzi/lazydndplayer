package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// MasteryDetailPopup displays details about a weapon mastery
type MasteryDetailPopup struct {
	visible     bool
	weaponName  string
	masteryType string
	description string
}

// NewMasteryDetailPopup creates a new mastery detail popup
func NewMasteryDetailPopup() *MasteryDetailPopup {
	return &MasteryDetailPopup{
		visible: false,
	}
}

// Show displays the popup with mastery details
func (p *MasteryDetailPopup) Show(weaponName, masteryType string) {
	p.visible = true
	p.weaponName = weaponName
	p.masteryType = masteryType
	p.description = models.GetMasteryDescription(masteryType)
}

// Hide hides the popup
func (p *MasteryDetailPopup) Hide() {
	p.visible = false
}

// IsVisible returns whether the popup is visible
func (p *MasteryDetailPopup) IsVisible() bool {
	return p.visible
}

// View renders the popup
func (p *MasteryDetailPopup) View(width, height int) string {
	if !p.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		MarginBottom(1)

	weaponStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	// Build content
	var content []string
	content = append(content, titleStyle.Render(p.masteryType))
	content = append(content, weaponStyle.Render("Weapon: "+p.weaponName))
	content = append(content, "")

	// Wrap description
	wrapped := wrapMasteryDescText(p.description, 70)
	for _, line := range wrapped {
		content = append(content, descStyle.Render(line))
	}

	content = append(content, "")
	content = append(content, helpStyle.Render("Press Esc to close"))

	contentStr := strings.Join(content, "\n")

	// Create popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(2, 4).
		Width(80).
		Align(lipgloss.Center)

	popup := popupStyle.Render(contentStr)

	// Center the popup on screen
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}

// wrapMasteryDescText wraps text to fit within a given width
func wrapMasteryDescText(text string, width int) []string {
	if len(text) == 0 {
		return []string{}
	}

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
