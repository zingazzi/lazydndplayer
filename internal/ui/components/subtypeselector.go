// internal/ui/components/subtypeselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpeciesSubtype represents a subtype/variant of a species
type SpeciesSubtype struct {
	Name        string
	Description string
	Modifier    string // e.g., "+5 ft speed" for Wood Elf, damage type for Dragonborn
}

// SubtypeSelector is a component for selecting species subtypes
type SubtypeSelector struct {
	subtypes      []SpeciesSubtype
	selectedIndex int
	visible       bool
	title         string
	SpeciesName   string // Exported so app.go can access it
}

// NewSubtypeSelector creates a new subtype selector
func NewSubtypeSelector() *SubtypeSelector {
	return &SubtypeSelector{
		subtypes:      []SpeciesSubtype{},
		selectedIndex: 0,
		visible:       false,
		title:         "SELECT SUBTYPE",
	}
}

// Show displays the subtype selector with available subtypes
func (s *SubtypeSelector) Show(speciesName string, subtypes []SpeciesSubtype) {
	s.SpeciesName = speciesName
	s.subtypes = subtypes
	s.selectedIndex = 0
	s.visible = true

	switch speciesName {
	case "Elf":
		s.title = "SELECT ELF SUBTYPE"
	case "Tiefling":
		s.title = "SELECT TIEFLING LINEAGE"
	case "Dragonborn":
		s.title = "SELECT DRACONIC ANCESTRY"
	default:
		s.title = fmt.Sprintf("SELECT %s SUBTYPE", strings.ToUpper(speciesName))
	}
}

// Hide hides the subtype selector
func (s *SubtypeSelector) Hide() {
	s.visible = false
}

// IsVisible returns whether the selector is visible
func (s *SubtypeSelector) IsVisible() bool {
	return s.visible
}

// Next moves to the next subtype
func (s *SubtypeSelector) Next() {
	if s.selectedIndex < len(s.subtypes)-1 {
		s.selectedIndex++
	}
}

// Prev moves to the previous subtype
func (s *SubtypeSelector) Prev() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// GetSelectedSubtype returns the currently selected subtype
func (s *SubtypeSelector) GetSelectedSubtype() *SpeciesSubtype {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.subtypes) {
		return &s.subtypes[s.selectedIndex]
	}
	return nil
}

// Update handles input for the subtype selector
func (s *SubtypeSelector) Update(msg tea.KeyMsg) tea.Cmd {
	if !s.visible {
		return nil
	}

	switch msg.String() {
	case "up", "k":
		s.Prev()
	case "down", "j":
		s.Next()
	case "esc":
		s.Hide()
	}

	return nil
}

// View renders the subtype selector
func (s *SubtypeSelector) View(width, height int) string {
	if !s.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Padding(0, 2)

	modifierStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Padding(0, 2)

	// Border style
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2)

	// Build content
	var content strings.Builder

	content.WriteString(titleStyle.Render(s.title) + "\n\n")

	// Show all subtypes
	for i, subtype := range s.subtypes {
		if i == s.selectedIndex {
			// Selected subtype - show with details
			content.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s", subtype.Name)) + "\n")

			if subtype.Description != "" {
				desc := subtype.Description
				if len(desc) > 80 {
					desc = desc[:77] + "..."
				}
				content.WriteString(descStyle.Render(desc) + "\n")
			}

			if subtype.Modifier != "" {
				content.WriteString(modifierStyle.Render(fmt.Sprintf("→ %s", subtype.Modifier)) + "\n")
			}

			content.WriteString("\n")
		} else {
			// Other subtypes - show name only
			subtypeLine := fmt.Sprintf("  %s", subtype.Name)
			if subtype.Modifier != "" {
				subtypeLine += fmt.Sprintf(" (%s)", subtype.Modifier)
			}
			content.WriteString(normalStyle.Render(subtypeLine) + "\n")
		}
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	// Wrap in border
	boxWidth := width - 20
	boxHeight := height - 10
	if boxWidth < 60 {
		boxWidth = 60
	}
	if boxHeight < 15 {
		boxHeight = 15
	}

	contentStr := content.String()

	// Center the box
	box := borderStyle.Width(boxWidth).Height(boxHeight).Render(contentStr)

	// Center horizontally and vertically
	paddingTop := (height - boxHeight) / 2
	paddingLeft := (width - boxWidth) / 2

	if paddingTop < 0 {
		paddingTop = 0
	}
	if paddingLeft < 0 {
		paddingLeft = 0
	}

	centered := lipgloss.NewStyle().
		Padding(paddingTop, 0, 0, paddingLeft).
		Render(box)

	return centered
}
