package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ManeuverSelector handles selection of Battle Master maneuvers
type ManeuverSelector struct {
	visible            bool
	maneuvers          []models.Maneuver
	selectedManeuvers  map[string]bool
	selectedIndex      int
	maxManeuvers       int
}

// NewManeuverSelector creates a new maneuver selector
func NewManeuverSelector() *ManeuverSelector {
	return &ManeuverSelector{
		visible:           false,
		maneuvers:         models.AllManeuvers,
		selectedManeuvers: make(map[string]bool),
		selectedIndex:     0,
		maxManeuvers:      0,
	}
}

// Show displays the maneuver selector
func (s *ManeuverSelector) Show(maxManeuvers int) {
	s.visible = true
	s.selectedIndex = 0
	s.maxManeuvers = maxManeuvers
	s.selectedManeuvers = make(map[string]bool)
}

// Hide hides the maneuver selector
func (s *ManeuverSelector) Hide() {
	s.visible = false
}

// IsVisible returns whether the selector is visible
func (s *ManeuverSelector) IsVisible() bool {
	return s.visible
}

// Next moves to the next maneuver
func (s *ManeuverSelector) Next() {
	if s.selectedIndex < len(s.maneuvers)-1 {
		s.selectedIndex++
	}
}

// Prev moves to the previous maneuver
func (s *ManeuverSelector) Prev() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// ToggleSelection toggles the current maneuver
func (s *ManeuverSelector) ToggleSelection() {
	maneuver := s.maneuvers[s.selectedIndex]
	if s.selectedManeuvers[maneuver.Name] {
		delete(s.selectedManeuvers, maneuver.Name)
	} else {
		// Only allow selection up to max
		if len(s.selectedManeuvers) < s.maxManeuvers {
			s.selectedManeuvers[maneuver.Name] = true
		}
	}
}

// CanConfirm returns true if the correct number of maneuvers are selected
func (s *ManeuverSelector) CanConfirm() bool {
	return len(s.selectedManeuvers) == s.maxManeuvers
}

// GetSelectedManeuvers returns the list of selected maneuvers
func (s *ManeuverSelector) GetSelectedManeuvers() []string {
	selected := make([]string, 0, len(s.selectedManeuvers))
	for name := range s.selectedManeuvers {
		selected = append(selected, name)
	}
	return selected
}

// SetSelectedManeuvers sets the currently selected maneuvers (for editing)
func (s *ManeuverSelector) SetSelectedManeuvers(maneuvers []string) {
	s.selectedManeuvers = make(map[string]bool)
	for _, m := range maneuvers {
		s.selectedManeuvers[m] = true
	}
}

// Update handles keyboard input
func (s *ManeuverSelector) Update(msg tea.Msg) (ManeuverSelector, tea.Cmd) {
	if !s.visible {
		return *s, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.Prev()
		case "down", "j":
			s.Next()
		case " ":
			s.ToggleSelection()
		case "enter", "esc":
			return *s, nil
		}
	}
	return *s, nil
}

// View renders the maneuver selector
func (s *ManeuverSelector) View() string {
	if !s.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	checkboxStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	typeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Left side - maneuver list
	var leftLines []string
	title := fmt.Sprintf("Select Battle Master Maneuvers (%d/%d selected)", len(s.selectedManeuvers), s.maxManeuvers)
	leftLines = append(leftLines, titleStyle.Render(title))
	leftLines = append(leftLines, "")

	for i, maneuver := range s.maneuvers {
		checkbox := "[ ]"
		if s.selectedManeuvers[maneuver.Name] {
			checkbox = checkboxStyle.Render("[✓]")
		}

		typeText := ""
		if maneuver.Type != "" {
			typeText = " " + typeStyle.Render(fmt.Sprintf("(%s)", maneuver.Type))
		}

		line := fmt.Sprintf("%s %s%s", checkbox, maneuver.Name, typeText)

		if i == s.selectedIndex {
			leftLines = append(leftLines, selectedStyle.Render("▶ "+line))
		} else {
			leftLines = append(leftLines, normalStyle.Render("  "+line))
		}
	}

	leftLines = append(leftLines, "")
	if s.CanConfirm() {
		leftLines = append(leftLines, helpStyle.Render("↑/↓: Navigate • Space: Toggle"))
		leftLines = append(leftLines, helpStyle.Render("Enter: Confirm • Esc: Cancel"))
	} else {
		leftLines = append(leftLines, helpStyle.Render("↑/↓: Navigate • Space: Toggle"))
		leftLines = append(leftLines, helpStyle.Render("Esc: Cancel"))
	}

	leftContent := strings.Join(leftLines, "\n")

	// Right side - maneuver description
	var rightLines []string
	rightLines = append(rightLines, titleStyle.Render("MANEUVER DESCRIPTION"))
	rightLines = append(rightLines, "")

	if s.selectedIndex >= 0 && s.selectedIndex < len(s.maneuvers) {
		currentManeuver := s.maneuvers[s.selectedIndex]

		// Show name and type
		rightLines = append(rightLines, selectedStyle.Render(currentManeuver.Name))
		if currentManeuver.Type != "" {
			rightLines = append(rightLines, typeStyle.Render(fmt.Sprintf("Type: %s", currentManeuver.Type)))
		}
		rightLines = append(rightLines, "")

		// Show description (wrapped)
		if currentManeuver.Description != "" {
			wrapped := wrapManeuverText(currentManeuver.Description, 45)
			rightLines = append(rightLines, normalStyle.Render(wrapped))
		}
	} else {
		rightLines = append(rightLines, normalStyle.Render("Select a maneuver to view its description"))
	}

	rightContent := strings.Join(rightLines, "\n")

	// Create two boxes side by side
	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(50).
		Height(28)

	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(52).
		Height(28)

	// Join boxes horizontally
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftBox.Render(leftContent),
		rightBox.Render(rightContent),
	)

	return combined
}

// wrapManeuverText wraps text to fit within a given width
func wrapManeuverText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine)
			result.WriteString("\n")
			currentLine = word
		}
	}

	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}
