package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConsumableDetailPopup displays details about a consumable resource or feature
type ConsumableDetailPopup struct {
	visible     bool
	name        string
	current     int
	max         int
	restType    string
	description string
}

// NewConsumableDetailPopup creates a new consumable detail popup
func NewConsumableDetailPopup() *ConsumableDetailPopup {
	return &ConsumableDetailPopup{
		visible: false,
	}
}

// Show displays the popup with consumable details
func (p *ConsumableDetailPopup) Show(name string, current, max int, restType, description string) {
	p.visible = true
	p.name = name
	p.current = current
	p.max = max
	p.restType = restType
	p.description = description
}

// Hide closes the popup
func (p *ConsumableDetailPopup) Hide() {
	p.visible = false
}

// IsVisible returns whether the popup is currently visible
func (p *ConsumableDetailPopup) IsVisible() bool {
	return p.visible
}

// View renders the popup
func (p *ConsumableDetailPopup) View(width, height int) string {
	if !p.visible {
		return ""
	}

	// Styles
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(width - 4)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	usageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	restStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("135")).
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render(p.name))
	content.WriteString("\n\n")

	// Usage bar
	usageBar := fmt.Sprintf("Uses: %d/%d", p.current, p.max)
	if p.current == 0 {
		usageBar += " [DEPLETED]"
	}
	content.WriteString(usageStyle.Render(usageBar))
	content.WriteString("\n")

	// Rest type
	content.WriteString(restStyle.Render(fmt.Sprintf("Recharges on: %s", p.restType)))
	content.WriteString("\n\n")

	// Description (wrapped)
	wrappedDesc := wrapConsumableText(p.description, width-8)
	content.WriteString(descStyle.Render(wrappedDesc))
	content.WriteString("\n\n")

	// Hint
	content.WriteString(hintStyle.Render("Press Enter or Esc to close"))

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

// wrapConsumableText wraps text to fit within the given width
func wrapConsumableText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// Check if adding this word would exceed the width
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		if len(testLine) > width {
			// Save current line and start new one
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		} else {
			// Add word to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}
	}

	// Add last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}
