// internal/ui/components/alignmentselector.go
package components

import (
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
	visible      bool
	selectedRow  int // 0-2 (Good/Neutral/Evil)
	selectedCol  int // 0-2 (Lawful/Neutral/Chaotic)
	alignments   [][]Alignment
}

// NewAlignmentSelector creates a new alignment selector
func NewAlignmentSelector() *AlignmentSelector {
	alignments := [][]Alignment{
		// Good row
		{
			{Name: "Lawful Good", Short: "LG", Description: "Acts with compassion and honor, respecting law and tradition while helping others"},
			{Name: "Neutral Good", Short: "NG", Description: "Does their best to help others, regardless of law or chaos"},
			{Name: "Chaotic Good", Short: "CG", Description: "Acts as their conscience directs with little regard for what others expect"},
		},
		// Neutral row
		{
			{Name: "Lawful Neutral", Short: "LN", Description: "Acts in accordance with law, tradition, or personal codes"},
			{Name: "True Neutral", Short: "N", Description: "Prefers to stay out of moral questions and doesn't take sides"},
			{Name: "Chaotic Neutral", Short: "CN", Description: "Follows their whims, holding personal freedom above all else"},
		},
		// Evil row
		{
			{Name: "Lawful Evil", Short: "LE", Description: "Methodically takes what they want within the limits of a code"},
			{Name: "Neutral Evil", Short: "NE", Description: "Does whatever they can get away with, without compassion or qualms"},
			{Name: "Chaotic Evil", Short: "CE", Description: "Acts with arbitrary violence, spurred by greed, hatred, or bloodlust"},
		},
	}

	return &AlignmentSelector{
		visible:      false,
		selectedRow:  1, // Start at True Neutral
		selectedCol:  1,
		alignments:   alignments,
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

// MoveUp moves selection up
func (a *AlignmentSelector) MoveUp() {
	if a.selectedRow > 0 {
		a.selectedRow--
	}
}

// MoveDown moves selection down
func (a *AlignmentSelector) MoveDown() {
	if a.selectedRow < 2 {
		a.selectedRow++
	}
}

// MoveLeft moves selection left
func (a *AlignmentSelector) MoveLeft() {
	if a.selectedCol > 0 {
		a.selectedCol--
	}
}

// MoveRight moves selection right
func (a *AlignmentSelector) MoveRight() {
	if a.selectedCol < 2 {
		a.selectedCol++
	}
}

// GetSelectedAlignment returns the currently selected alignment
func (a *AlignmentSelector) GetSelectedAlignment() string {
	return a.alignments[a.selectedRow][a.selectedCol].Name
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
		Background(lipgloss.Color("236")).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Width(60).
		Align(lipgloss.Center)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	var content strings.Builder

	content.WriteString(titleStyle.Render("SELECT ALIGNMENT"))
	content.WriteString("\n\n")

	// Render 3x3 grid
	columnHeaders := []string{"Lawful", "Neutral", "Chaotic"}
	rowHeaders := []string{"Good", "Neutral", "Evil"}

	// Column headers
	content.WriteString("           ")
	for _, header := range columnHeaders {
		headerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Width(20).
			Align(lipgloss.Center)
		content.WriteString(headerStyle.Render(header))
	}
	content.WriteString("\n")

	// Grid rows
	for row := 0; row < 3; row++ {
		// Row header
		rowHeaderStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Width(10).
			Align(lipgloss.Right)
		content.WriteString(rowHeaderStyle.Render(rowHeaders[row]))
		content.WriteString(" ")

		// Cells
		for col := 0; col < 3; col++ {
			alignment := a.alignments[row][col]
			cellText := alignment.Short

			var cellStyle lipgloss.Style
			if row == a.selectedRow && col == a.selectedCol {
				cellStyle = selectedStyle
			} else {
				cellStyle = normalStyle
			}

			content.WriteString(cellStyle.Render(cellText))
		}
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Show description of selected alignment
	selectedAlignment := a.alignments[a.selectedRow][a.selectedCol]
	content.WriteString(descStyle.Render(selectedAlignment.Description))
	content.WriteString("\n\n")

	content.WriteString(hintStyle.Render("Arrow keys: Navigate | Enter: Select | Esc: Cancel"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(70)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

