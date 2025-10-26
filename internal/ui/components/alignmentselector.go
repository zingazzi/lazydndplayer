// internal/ui/components/alignmentselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Alignment represents a D&D alignment with description
type Alignment struct {
	Name        string
	Short       string
	Description string
}

// AlignmentSelector allows selection of character alignment
type AlignmentSelector struct {
	visible       bool
	selectedIndex int
	alignments    []Alignment
}

// NewAlignmentSelector creates a new alignment selector
func NewAlignmentSelector() *AlignmentSelector {
	alignments := []Alignment{
		{Name: "Lawful Good", Short: "LG", Description: "Acts with compassion and honor, respecting law and tradition while helping others"},
		{Name: "Neutral Good", Short: "NG", Description: "Does their best to help others, regardless of law or chaos"},
		{Name: "Chaotic Good", Short: "CG", Description: "Acts as their conscience directs with little regard for what others expect"},
		{Name: "Lawful Neutral", Short: "LN", Description: "Acts in accordance with law, tradition, or personal codes"},
		{Name: "True Neutral", Short: "N", Description: "Prefers to stay out of moral questions and doesn't take sides"},
		{Name: "Chaotic Neutral", Short: "CN", Description: "Follows their whims, holding personal freedom above all else"},
		{Name: "Lawful Evil", Short: "LE", Description: "Methodically takes what they want within the limits of a code"},
		{Name: "Neutral Evil", Short: "NE", Description: "Does whatever they can get away with, without compassion or qualms"},
		{Name: "Chaotic Evil", Short: "CE", Description: "Acts with arbitrary violence, spurred by greed, hatred, or bloodlust"},
	}

	return &AlignmentSelector{
		visible:       false,
		selectedIndex: 4, // Default to True Neutral
		alignments:    alignments,
	}
}

// Show displays the selector
func (a *AlignmentSelector) Show() {
	a.visible = true
}

// Hide closes the selector
func (a *AlignmentSelector) Hide() {
	a.visible = false
}

// IsVisible returns whether the selector is currently visible
func (a *AlignmentSelector) IsVisible() bool {
	return a.visible
}

// Next moves selection down
func (a *AlignmentSelector) Next() {
	if a.selectedIndex < len(a.alignments)-1 {
		a.selectedIndex++
	}
}

// Prev moves selection up
func (a *AlignmentSelector) Prev() {
	if a.selectedIndex > 0 {
		a.selectedIndex--
	}
}

// GetSelectedAlignment returns the currently selected alignment
func (a *AlignmentSelector) GetSelectedAlignment() string {
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.alignments) {
		return a.alignments[a.selectedIndex].Name
	}
	return "True Neutral"
}

// View renders the selector
func (a *AlignmentSelector) View(width, height int) string {
	if !a.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	var content strings.Builder

	content.WriteString(titleStyle.Render("SELECT ALIGNMENT"))
	content.WriteString("\n\n")

	// Render as vertical list
	for i, alignment := range a.alignments {
		cursor := "  "
		style := normalStyle
		if i == a.selectedIndex {
			cursor = "→ "
			style = selectedStyle
		}

		line := fmt.Sprintf("%-15s (%s)", alignment.Name, alignment.Short)
		content.WriteString(style.Render(cursor + line))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Show description of selected alignment
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.alignments) {
		selectedAlignment := a.alignments[a.selectedIndex]
		// Wrap description to fit width
		wrappedDesc := wrapAlignmentText(selectedAlignment.Description, 70)
		for _, line := range wrappedDesc {
			content.WriteString(descStyle.Render(line))
			content.WriteString("\n")
		}
	}
	content.WriteString("\n")

	content.WriteString(hintStyle.Render("↑/↓: Navigate | Enter: Select | Esc: Cancel"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(80)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

// wrapAlignmentText wraps text to fit within the given width
func wrapAlignmentText(text string, width int) []string {
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
