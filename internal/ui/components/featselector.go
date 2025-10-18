// internal/ui/components/featselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// FeatSelector is a component for selecting feats
type FeatSelector struct {
	feats         []models.Feat
	selectedIndex int
	visible       bool
	title         string
	character     *models.Character
	filterOrigin  bool // Only show origin-appropriate feats
	deleteMode    bool // If true, shows known feats for deletion
}

// NewFeatSelector creates a new feat selector
func NewFeatSelector() *FeatSelector {
	return &FeatSelector{
		feats:         []models.Feat{},
		selectedIndex: 0,
		visible:       false,
		title:         "SELECT FEAT",
		filterOrigin:  false,
		deleteMode:    false,
	}
}

// Show displays the feat selector with available feats for the character
func (f *FeatSelector) Show(char *models.Character, originFeat bool) {
	f.character = char
	f.filterOrigin = originFeat
	f.deleteMode = false

	if originFeat {
		f.title = "SELECT ORIGIN FEAT (HUMAN)"
		// For origin feats, show all feats the character can take
		f.feats = models.GetFeatsForCharacter(char)
	} else {
		f.title = "SELECT FEAT"
		f.feats = models.GetFeatsForCharacter(char)
	}

	f.selectedIndex = 0
	f.visible = true
}

// ShowForDeletion displays the feat selector with known feats (for deleting)
func (f *FeatSelector) ShowForDeletion(char *models.Character) {
	f.character = char
	f.filterOrigin = false
	f.deleteMode = true
	f.title = "SELECT FEAT TO REMOVE"

	// Get all feats and filter to only show ones the character has
	allFeats := models.GetAllFeats()
	f.feats = []models.Feat{}

	for _, featName := range char.Feats {
		// Find the feat details
		for _, feat := range allFeats {
			if feat.Name == featName {
				f.feats = append(f.feats, feat)
				break
			}
		}
	}

	f.selectedIndex = 0

	// If no feats, close selector
	if len(f.feats) == 0 {
		f.visible = false
	} else {
		f.visible = true
	}
}

// Hide hides the feat selector
func (f *FeatSelector) Hide() {
	f.visible = false
	f.deleteMode = false
}

// IsDeleteMode returns whether the selector is in delete mode
func (f *FeatSelector) IsDeleteMode() bool {
	return f.deleteMode
}

// IsVisible returns whether the selector is visible
func (f *FeatSelector) IsVisible() bool {
	return f.visible
}

// Next moves to the next feat
func (f *FeatSelector) Next() {
	if f.selectedIndex < len(f.feats)-1 {
		f.selectedIndex++
	}
}

// Prev moves to the previous feat
func (f *FeatSelector) Prev() {
	if f.selectedIndex > 0 {
		f.selectedIndex--
	}
}

// PageDown moves down by 5 feats
func (f *FeatSelector) PageDown() {
	f.selectedIndex += 5
	if f.selectedIndex >= len(f.feats) {
		f.selectedIndex = len(f.feats) - 1
	}
}

// PageUp moves up by 5 feats
func (f *FeatSelector) PageUp() {
	f.selectedIndex -= 5
	if f.selectedIndex < 0 {
		f.selectedIndex = 0
	}
}

// GetSelectedFeat returns the currently selected feat
func (f *FeatSelector) GetSelectedFeat() *models.Feat {
	if f.selectedIndex >= 0 && f.selectedIndex < len(f.feats) {
		return &f.feats[f.selectedIndex]
	}
	return nil
}

// Update handles input for the feat selector
func (f *FeatSelector) Update(msg tea.KeyMsg) tea.Cmd {
	if !f.visible {
		return nil
	}

	switch msg.String() {
	case "up", "k":
		f.Prev()
	case "down", "j":
		f.Next()
	case "pgup", "ctrl+u":
		f.PageUp()
	case "pgdown", "ctrl+d":
		f.PageDown()
	case "esc":
		f.Hide()
	}

	return nil
}

// View renders the feat selector
func (f *FeatSelector) View(width, height int) string {
	if !f.visible {
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

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 2)

	prerequisiteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Italic(true)

	// Border style
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2)

	// Build content
	var content strings.Builder

	content.WriteString(titleStyle.Render(f.title) + "\n\n")

	// Calculate visible range (show 10 feats at a time)
	visibleStart := f.selectedIndex - 5
	if visibleStart < 0 {
		visibleStart = 0
	}
	visibleEnd := visibleStart + 10
	if visibleEnd > len(f.feats) {
		visibleEnd = len(f.feats)
	}

	// Show feats
	for i := visibleStart; i < visibleEnd; i++ {
		feat := f.feats[i]

		if i == f.selectedIndex {
			// Selected feat - show full details
			content.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s", feat.Name)) + "\n")

			// Show prerequisite
			if feat.Prerequisite != "None" && feat.Prerequisite != "" {
				content.WriteString(detailStyle.Render(
					prerequisiteStyle.Render(fmt.Sprintf("Requires: %s", feat.Prerequisite)),
				) + "\n")
			}

			// Show ability increases
			if feat.AbilityIncreases != nil {
				if len(feat.AbilityIncreases.Choices) > 0 {
					content.WriteString(detailStyle.Render(
						fmt.Sprintf("  +%d to one: %s",
							feat.AbilityIncreases.Amount,
							strings.Join(feat.AbilityIncreases.Choices, " or ")),
					) + "\n")
				} else if feat.AbilityIncreases.Ability != "" {
					content.WriteString(detailStyle.Render(
						fmt.Sprintf("  +%d %s",
							feat.AbilityIncreases.Amount,
							feat.AbilityIncreases.Ability),
					) + "\n")
				}
			}

			// Show description
			desc := feat.Description
			if len(desc) > 80 {
				desc = desc[:77] + "..."
			}
			content.WriteString(detailStyle.Render(fmt.Sprintf("  %s", desc)) + "\n")

			// Show first benefit
			if len(feat.Benefits) > 0 {
				benefit := feat.Benefits[0]
				if len(benefit) > 70 {
					benefit = benefit[:67] + "..."
				}
				content.WriteString(detailStyle.Render(fmt.Sprintf("  • %s", benefit)) + "\n")
			}

			content.WriteString("\n")
		} else {
			// Other feats - show name only
			featLine := fmt.Sprintf("  %s", feat.Name)
			if feat.Prerequisite != "None" && feat.Prerequisite != "" {
				featLine += " *"
			}
			content.WriteString(normalStyle.Render(featLine) + "\n")
		}
	}

	// Show scroll indicator
	if len(f.feats) > 10 {
		scrollPercent := int(float64(f.selectedIndex) / float64(len(f.feats)-1) * 100)
		content.WriteString(fmt.Sprintf("\n%d/%d feats (%d%%)", f.selectedIndex+1, len(f.feats), scrollPercent))
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	content.WriteString("\n\n")

	// Show different help text based on mode
	if f.deleteMode {
		content.WriteString(helpStyle.Render("[↑/↓] Navigate • [Enter] Remove • [Esc] Cancel"))
	} else {
		content.WriteString(helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))
		content.WriteString("\n")
		content.WriteString(helpStyle.Render("* = Has prerequisites"))
	}

	// Wrap in border
	boxWidth := width - 20
	boxHeight := height - 10
	if boxWidth < 60 {
		boxWidth = 60
	}
	if boxHeight < 20 {
		boxHeight = 20
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
