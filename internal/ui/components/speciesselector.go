// internal/ui/components/speciesselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpeciesSelector handles species selection UI
type SpeciesSelector struct {
	species       []models.SpeciesInfo
	selectedIndex int
	visible       bool
}

// NewSpeciesSelector creates a new species selector
func NewSpeciesSelector() *SpeciesSelector {
	return &SpeciesSelector{
		species:       models.GetAllSpecies(),
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the species selector
func (ss *SpeciesSelector) Show() {
	ss.visible = true
}

// Hide hides the species selector
func (ss *SpeciesSelector) Hide() {
	ss.visible = false
}

// IsVisible returns whether the selector is visible
func (ss *SpeciesSelector) IsVisible() bool {
	return ss.visible
}

// Next moves to the next species
func (ss *SpeciesSelector) Next() {
	if ss.selectedIndex < len(ss.species)-1 {
		ss.selectedIndex++
	}
}

// Prev moves to the previous species
func (ss *SpeciesSelector) Prev() {
	if ss.selectedIndex > 0 {
		ss.selectedIndex--
	}
}

// GetSelectedSpecies returns the currently selected species
func (ss *SpeciesSelector) GetSelectedSpecies() *models.SpeciesInfo {
	if ss.selectedIndex >= 0 && ss.selectedIndex < len(ss.species) {
		return &ss.species[ss.selectedIndex]
	}
	return nil
}

// View renders the species selector
func (ss *SpeciesSelector) View(screenWidth, screenHeight int) string {
	if !ss.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	speciesNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	selectedSpeciesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalSpeciesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	traitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content []string
	content = append(content, titleStyle.Render("SELECT SPECIES"))
	content = append(content, "")

	// Left side: Species list
	var speciesList []string
	for i, species := range ss.species {
		speciesLine := fmt.Sprintf(" %s", species.Name)
		if i == ss.selectedIndex {
			speciesList = append(speciesList, selectedSpeciesStyle.Render(speciesLine))
		} else {
			speciesList = append(speciesList, normalSpeciesStyle.Render(speciesLine))
		}
	}

	// Right side: Selected species details
	selectedSpecies := ss.GetSelectedSpecies()
	var speciesDetails []string
	if selectedSpecies != nil {
		speciesDetails = append(speciesDetails, speciesNameStyle.Render(selectedSpecies.Name))
		speciesDetails = append(speciesDetails, "")
		speciesDetails = append(speciesDetails, descStyle.Render(selectedSpecies.Description))
		speciesDetails = append(speciesDetails, "")
		speciesDetails = append(speciesDetails, fmt.Sprintf("Size: %s  Speed: %d ft", selectedSpecies.Size, selectedSpecies.Speed))
		speciesDetails = append(speciesDetails, "")
		speciesDetails = append(speciesDetails, lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("TRAITS:"))
		for _, trait := range selectedSpecies.Traits {
			speciesDetails = append(speciesDetails, "")
			speciesDetails = append(speciesDetails, traitStyle.Render("• "+trait.Name))
			// Wrap description
			wrapped := wrapText(trait.Description, 50)
			for _, line := range wrapped {
				speciesDetails = append(speciesDetails, "  "+lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(line))
			}
		}
		speciesDetails = append(speciesDetails, "")
		speciesDetails = append(speciesDetails, lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("LANGUAGES:"))
		for _, lang := range selectedSpecies.Languages {
			speciesDetails = append(speciesDetails, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("• "+lang))
		}
	}

	// Combine list and details side by side
	listBox := lipgloss.NewStyle().
		Width(20).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Render(strings.Join(speciesList, "\n"))

	detailsBox := lipgloss.NewStyle().
		Width(55).
		Height(30).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Render(strings.Join(speciesDetails, "\n"))

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listBox,
		" ",
		detailsBox,
	)

	content = append(content, mainContent)
	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	// Wrap in a styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// wrapText wraps text to a specified width for species selector
func wrapText(text string, width int) []string {
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
