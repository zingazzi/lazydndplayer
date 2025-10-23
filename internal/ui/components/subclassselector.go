// internal/ui/components/subclassselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SubclassSelector handles subclass selection for classes that require it at specific levels
type SubclassSelector struct {
	visible     bool
	character   *models.Character
	className   string
	classLevel  int
	subclasses  []models.Subclass
	cursor      int
	selectedIdx int
}

// NewSubclassSelector creates a new subclass selector
func NewSubclassSelector(char *models.Character) *SubclassSelector {
	return &SubclassSelector{
		character:   char,
		selectedIdx: -1,
	}
}

// Show displays the subclass selector for a specific class and level
func (ss *SubclassSelector) Show(className string, classLevel int) {
	ss.visible = true
	ss.className = className
	ss.classLevel = classLevel
	ss.cursor = 0
	ss.selectedIdx = -1

	// Load subclasses for this class
	class := models.GetClassByName(className)
	if class != nil {
		ss.subclasses = class.Subclasses
	} else {
		ss.subclasses = []models.Subclass{}
	}
}

// Hide hides the subclass selector
func (ss *SubclassSelector) Hide() {
	ss.visible = false
	ss.className = ""
	ss.subclasses = nil
	ss.cursor = 0
	ss.selectedIdx = -1
}

// IsVisible returns whether the selector is visible
func (ss *SubclassSelector) IsVisible() bool {
	return ss.visible
}

// GetSelectedSubclass returns the currently selected subclass
func (ss *SubclassSelector) GetSelectedSubclass() *models.Subclass {
	if ss.selectedIdx >= 0 && ss.selectedIdx < len(ss.subclasses) {
		return &ss.subclasses[ss.selectedIdx]
	}
	return nil
}

// Next moves cursor down
func (ss *SubclassSelector) Next() {
	if ss.cursor < len(ss.subclasses)-1 {
		ss.cursor++
	}
}

// Prev moves cursor up
func (ss *SubclassSelector) Prev() {
	if ss.cursor > 0 {
		ss.cursor--
	}
}

// Select confirms the current selection
func (ss *SubclassSelector) Select() bool {
	if ss.cursor >= 0 && ss.cursor < len(ss.subclasses) {
		ss.selectedIdx = ss.cursor
		return true
	}
	return false
}

// Update handles key presses
func (ss *SubclassSelector) Update(msg tea.Msg) (SubclassSelector, tea.Cmd) {
	if !ss.visible {
		return *ss, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			ss.Prev()
		case "down", "j":
			ss.Next()
		case "enter":
			if ss.Select() {
				// Hide the selector after successful selection
				ss.visible = false
			}
		case "esc":
			// Cancel without selecting
			ss.selectedIdx = -1
			ss.Hide()
		}
	}

	return *ss, nil
}

// View renders the subclass selector
func (ss *SubclassSelector) View() string {
	if !ss.visible || len(ss.subclasses) == 0 {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(60).
		Padding(0, 1)

	// Get label for subclass type
	subclassLabel := "Subclass"
	switch ss.className {
	case "Cleric":
		subclassLabel = "Domain"
	case "Warlock":
		subclassLabel = "Patron"
	case "Sorcerer":
		subclassLabel = "Sorcerous Origin"
	case "Druid":
		subclassLabel = "Circle"
	case "Wizard":
		subclassLabel = "Arcane Tradition"
	case "Fighter":
		subclassLabel = "Martial Archetype"
	case "Rogue":
		subclassLabel = "Roguish Archetype"
	case "Bard":
		subclassLabel = "College"
	case "Paladin":
		subclassLabel = "Sacred Oath"
	case "Ranger":
		subclassLabel = "Ranger Archetype"
	case "Barbarian":
		subclassLabel = "Primal Path"
	case "Monk":
		subclassLabel = "Monastic Tradition"
	}

	var content strings.Builder

	content.WriteString(titleStyle.Render(fmt.Sprintf("SELECT %s - %s", strings.ToUpper(ss.className), strings.ToUpper(subclassLabel))) + "\n\n")

	// Two-column layout: list on left, description on right
	var leftContent strings.Builder
	var rightContent strings.Builder

	// Build subclass list (left side)
	for i, subclass := range ss.subclasses {
		cursor := "  "
		style := normalStyle
		if i == ss.cursor {
			cursor = "❯ "
			style = selectedStyle
		}
		leftContent.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, subclass.Name)) + "\n")
	}

	// Build description (right side)
	if ss.cursor >= 0 && ss.cursor < len(ss.subclasses) {
		currentSubclass := ss.subclasses[ss.cursor]

		// Description
		rightContent.WriteString(selectedStyle.Render(currentSubclass.Name) + "\n\n")
		rightContent.WriteString(descStyle.Render(currentSubclass.Description) + "\n\n")

		// Features
		if len(currentSubclass.Features) > 0 {
			rightContent.WriteString(titleStyle.Render("FEATURES:") + "\n")
			for _, feature := range currentSubclass.Features {
				rightContent.WriteString(normalStyle.Render("• "+feature.Name) + "\n")
				if feature.Description != "" {
					rightContent.WriteString(dimStyle.Render("  "+feature.Description) + "\n")
				}
			}
		}

		// Expanded spell list (for classes like Cleric, Warlock)
		if len(currentSubclass.ExpandedSpells) > 0 {
			rightContent.WriteString("\n" + titleStyle.Render("EXPANDED SPELLS:") + "\n")
			rightContent.WriteString(dimStyle.Render(strings.Join(currentSubclass.ExpandedSpells, ", ")) + "\n")
		}
	}

	// Join left and right in two columns
	leftBox := lipgloss.NewStyle().
		Width(30).
		Height(20).
		Padding(1).
		Render(leftContent.String())

	rightBox := lipgloss.NewStyle().
		Width(65).
		Height(20).
		Padding(1).
		Render(rightContent.String())

	twoColumns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	content.WriteString(twoColumns)

	content.WriteString("\n" + dimStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Cancel"))

	// Popup style
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(100)

	return lipgloss.Place(
		120,
		35,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
