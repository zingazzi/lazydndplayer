// internal/ui/components/featuredetailpopup.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FeatureDetailPopup displays detailed information about a feature
type FeatureDetailPopup struct {
	visible     bool
	featureName string
	description string
	uses        string
	restType    string
}

// NewFeatureDetailPopup creates a new feature detail popup
func NewFeatureDetailPopup() *FeatureDetailPopup {
	return &FeatureDetailPopup{
		visible: false,
	}
}

// Show displays the popup with feature details
func (f *FeatureDetailPopup) Show(name, description, uses, restType string) {
	f.visible = true
	f.featureName = name
	f.description = description
	f.uses = uses
	f.restType = restType
}

// Hide closes the popup
func (f *FeatureDetailPopup) Hide() {
	f.visible = false
}

// IsVisible returns whether the popup is currently visible
func (f *FeatureDetailPopup) IsVisible() bool {
	return f.visible
}

// View renders the popup
func (f *FeatureDetailPopup) View(width, height int) string {
	if !f.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	var content strings.Builder

	content.WriteString(titleStyle.Render(f.featureName))
	content.WriteString("\n\n")

	// Show uses if applicable
	if f.uses != "" {
		content.WriteString(labelStyle.Render("Uses: "))
		content.WriteString(valueStyle.Render(f.uses))
		content.WriteString("\n")
	}

	// Show rest type if applicable
	if f.restType != "" && f.restType != "None" {
		content.WriteString(labelStyle.Render("Recovery: "))
		restTypeText := string(f.restType)
		content.WriteString(valueStyle.Render(restTypeText))
		content.WriteString("\n")
	}

	if f.uses != "" || f.restType != "" {
		content.WriteString("\n")
	}

	// Description
	content.WriteString(labelStyle.Render("Description:"))
	content.WriteString("\n")

	// Wrap description text
	wrappedDesc := wrapFeatureText(f.description, width-10)
	for _, line := range wrappedDesc {
		content.WriteString(descStyle.Render(line))
		content.WriteString("\n")
	}

	content.WriteString("\n")
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

// wrapFeatureText wraps text to fit within the given width
func wrapFeatureText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		if len(testLine) > width {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
