// internal/ui/components/featdetailpopup.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// FeatDetailPopup displays detailed information about a feat the character has
type FeatDetailPopup struct {
	visible   bool
	featName  string
	feat      *models.Feat
	character *models.Character
}

// NewFeatDetailPopup creates a new feat detail popup
func NewFeatDetailPopup() *FeatDetailPopup {
	return &FeatDetailPopup{
		visible: false,
	}
}

// Show displays the feat detail popup
func (f *FeatDetailPopup) Show(featName string, char *models.Character) {
	f.featName = featName
	f.character = char
	f.feat = models.GetFeatByName(featName)
	f.visible = true
}

// Hide hides the feat detail popup
func (f *FeatDetailPopup) Hide() {
	f.visible = false
}

// IsVisible returns whether the popup is visible
func (f *FeatDetailPopup) IsVisible() bool {
	return f.visible
}

// View renders the feat detail popup
func (f *FeatDetailPopup) View(width, height int) string {
	if !f.visible || f.feat == nil {
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

	categoryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Italic(true)

	prerequisiteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	choiceStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	benefitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content []string

	// Title
	content = append(content, titleStyle.Render(f.feat.Name))
	content = append(content, "")

	// Category and repeatability
	categoryInfo := f.feat.Category
	if f.feat.Repeatable {
		categoryInfo += " (Repeatable)"
	}
	content = append(content, categoryStyle.Render(categoryInfo))
	content = append(content, "")

	// Prerequisite
	if f.feat.Prerequisite != "None" && f.feat.Prerequisite != "" {
		content = append(content, prerequisiteStyle.Render("Prerequisite: "+f.feat.Prerequisite))
		content = append(content, "")
	}

	// Get benefits from BenefitTracker to show what choices were made
	benefits := f.character.BenefitTracker.GetBenefitsBySource("feat", f.featName)

	// Show granted benefits (including choices made)
	if len(benefits) > 0 {
		content = append(content, sectionTitleStyle.Render("GRANTED BENEFITS:"))

		for _, benefit := range benefits {
			switch benefit.Type {
			case models.BenefitAbilityScore:
				content = append(content, choiceStyle.Render(fmt.Sprintf("  ✓ +%d %s", benefit.Value, benefit.Target)))
			case models.BenefitSkill:
				content = append(content, choiceStyle.Render(fmt.Sprintf("  ✓ Proficiency: %s", benefit.Target)))
			case models.BenefitLanguage:
				content = append(content, choiceStyle.Render(fmt.Sprintf("  ✓ Language: %s", benefit.Target)))
			case models.BenefitHP:
				content = append(content, choiceStyle.Render(fmt.Sprintf("  ✓ +%d HP", benefit.Value)))
			case models.BenefitSpeed:
				content = append(content, choiceStyle.Render(fmt.Sprintf("  ✓ +%d ft Speed", benefit.Value)))
			}
		}
		content = append(content, "")
	}

	// Description
	content = append(content, sectionTitleStyle.Render("DESCRIPTION:"))
	content = append(content, "")
	wrapped := wrapDetailText(f.feat.Description, 70)
	for _, line := range wrapped {
		content = append(content, descStyle.Render(line))
	}
	content = append(content, "")

	// Benefits
	if len(f.feat.Benefits) > 0 {
		content = append(content, sectionTitleStyle.Render("FEAT BENEFITS:"))
		for _, benef := range f.feat.Benefits {
			content = append(content, "")
			wrapped := wrapDetailText(benef, 70)
			for i, line := range wrapped {
				if i == 0 {
					content = append(content, benefitStyle.Render("• "+line))
				} else {
					content = append(content, benefitStyle.Render("  "+line))
				}
			}
		}
		content = append(content, "")
	}

	// Help text
	content = append(content, "")
	content = append(content, helpStyle.Render("[Esc] Close"))

	// Wrap in styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Background(lipgloss.Color("235")).
		Width(78)

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}

// wrapDetailText wraps text to a specified width for feat detail popup
func wrapDetailText(text string, width int) []string {
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
