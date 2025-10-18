// internal/ui/components/originselector.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// OriginSelector is a component for selecting an origin
type OriginSelector struct {
	visible       bool
	selectedIndex int
	origins       []models.Origin
	character     *models.Character
}

// NewOriginSelector creates a new origin selector
func NewOriginSelector() *OriginSelector {
	return &OriginSelector{
		visible: false,
	}
}

// Show displays the origin selector
func (s *OriginSelector) Show(char *models.Character) {
	s.character = char
	s.origins = models.GetAllOrigins()
	s.selectedIndex = 0
	s.visible = true
}

// Hide hides the origin selector
func (s *OriginSelector) Hide() {
	s.visible = false
}

// IsVisible returns whether the selector is visible
func (s *OriginSelector) IsVisible() bool {
	return s.visible
}

// Next selects the next origin
func (s *OriginSelector) Next() {
	if s.selectedIndex < len(s.origins)-1 {
		s.selectedIndex++
	}
}

// Prev selects the previous origin
func (s *OriginSelector) Prev() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// GetSelected returns the currently selected origin
func (s *OriginSelector) GetSelected() *models.Origin {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.origins) {
		return &s.origins[s.selectedIndex]
	}
	return nil
}

// View renders the origin selector
func (s *OriginSelector) View(width, height int) string {
	if !s.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	abilityStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	featStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Bold(true)

	skillStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	// Layout: two columns
	leftWidth := 30
	rightWidth := width - leftWidth - 6

	// Left column: origin list
	var leftContent []string
	leftContent = append(leftContent, titleStyle.Render("SELECT ORIGIN"))
	leftContent = append(leftContent, "")

	// Calculate scroll offset to keep selected item in view
	visibleItems := height - 8
	scrollOffset := 0
	if s.selectedIndex > visibleItems/2 {
		scrollOffset = s.selectedIndex - visibleItems/2
	}
	if scrollOffset+visibleItems > len(s.origins) {
		scrollOffset = len(s.origins) - visibleItems
	}
	if scrollOffset < 0 {
		scrollOffset = 0
	}

	// Show origins with scroll
	for i := scrollOffset; i < len(s.origins) && i < scrollOffset+visibleItems; i++ {
		origin := s.origins[i]
		if i == s.selectedIndex {
			leftContent = append(leftContent, selectedStyle.Render("→ "+origin.Name))
		} else {
			leftContent = append(leftContent, normalStyle.Render("  "+origin.Name))
		}
	}

	// Right column: selected origin details
	var rightContent []string
	selectedOrigin := s.GetSelected()
	if selectedOrigin != nil {
		rightContent = append(rightContent, sectionTitleStyle.Render(selectedOrigin.Name))
		rightContent = append(rightContent, "")

		// Description
		wrapped := wrapOriginText(selectedOrigin.Description, rightWidth-4)
		for _, line := range wrapped {
			rightContent = append(rightContent, descStyle.Render(line))
		}
		rightContent = append(rightContent, "")

		// Ability Increases
		if selectedOrigin.AbilityIncreases != nil {
			if len(selectedOrigin.AbilityIncreases.Choices) > 0 {
				rightContent = append(rightContent, abilityStyle.Render("ABILITY CHOICE:"))
				rightContent = append(rightContent, normalStyle.Render(
					"  +"+string(rune(selectedOrigin.AbilityIncreases.Amount+'0'))+
						" to "+strings.Join(selectedOrigin.AbilityIncreases.Choices, " or ")))
			} else if selectedOrigin.AbilityIncreases.Ability != "" {
				rightContent = append(rightContent, abilityStyle.Render("ABILITY:"))
				rightContent = append(rightContent, normalStyle.Render(
					"  +"+string(rune(selectedOrigin.AbilityIncreases.Amount+'0'))+
						" "+selectedOrigin.AbilityIncreases.Ability))
			}
			rightContent = append(rightContent, "")
		}

		// Granted Feat
		if selectedOrigin.Feat != "" {
			rightContent = append(rightContent, featStyle.Render("GRANTED FEAT:"))
			rightContent = append(rightContent, normalStyle.Render("  "+selectedOrigin.Feat))
			rightContent = append(rightContent, "")
		}

		// Skill Proficiencies
		if len(selectedOrigin.SkillProficiencies) > 0 {
			rightContent = append(rightContent, skillStyle.Render("SKILL PROFICIENCIES:"))
			for _, skill := range selectedOrigin.SkillProficiencies {
				rightContent = append(rightContent, normalStyle.Render("  • "+skill))
			}
			rightContent = append(rightContent, "")
		}

		// Tool Proficiencies
		if len(selectedOrigin.ToolProficiencies) > 0 {
			rightContent = append(rightContent, skillStyle.Render("TOOL PROFICIENCIES:"))
			for _, tool := range selectedOrigin.ToolProficiencies {
				rightContent = append(rightContent, normalStyle.Render("  • "+tool))
			}
			rightContent = append(rightContent, "")
		}

		// Equipment
		if len(selectedOrigin.Equipment) > 0 {
			rightContent = append(rightContent, sectionTitleStyle.Render("EQUIPMENT:"))
			for _, item := range selectedOrigin.Equipment {
				rightContent = append(rightContent, normalStyle.Render("  • "+item))
			}
		}
	}

	// Combine columns
	leftCol := lipgloss.NewStyle().
		Width(leftWidth).
		Render(strings.Join(leftContent, "\n"))

	rightCol := lipgloss.NewStyle().
		Width(rightWidth).
		Render(strings.Join(rightContent, "\n"))

	combined := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// Add help text
	helpText := helpStyle.Render("[↑/↓] Navigate  [Enter] Select  [Esc] Cancel")

	content := lipgloss.JoinVertical(lipgloss.Left, combined, "", helpText)

	// Wrap in styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(content)

	// Center on screen
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}

// wrapOriginText wraps text to a specified width
func wrapOriginText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
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

