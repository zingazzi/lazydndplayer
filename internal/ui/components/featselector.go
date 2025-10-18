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
	feats          []models.Feat
	allFeats       []models.Feat // All feats before filtering
	selectedIndex  int
	visible        bool
	title          string
	character      *models.Character
	filterOrigin   bool   // Only show origin-appropriate feats
	deleteMode     bool   // If true, shows known feats for deletion
	categoryFilter string // Current category filter ("All", "General", etc.)
	categories     []string
}

// NewFeatSelector creates a new feat selector
func NewFeatSelector() *FeatSelector {
	return &FeatSelector{
		feats:          []models.Feat{},
		allFeats:       []models.Feat{},
		selectedIndex:  0,
		visible:        false,
		title:          "SELECT FEAT",
		filterOrigin:   false,
		deleteMode:     false,
		categoryFilter: "All",
		categories:     []string{"All", "General", "Combat", "Magic", "Skill"},
	}
}

// Show displays the feat selector with available feats for the character
func (f *FeatSelector) Show(char *models.Character, originFeat bool) {
	f.character = char
	f.filterOrigin = originFeat
	f.deleteMode = false
	f.categoryFilter = "All"

	if originFeat {
		f.title = "SELECT ORIGIN FEAT (HUMAN)"
	} else {
		f.title = "SELECT FEAT"
	}

	// Load ALL feats (not just available ones)
	f.allFeats = models.GetAllFeats()

	// Remove already-taken non-repeatable feats
	availableFeats := []models.Feat{}
	for _, feat := range f.allFeats {
		// Skip if character already has this feat and it's not repeatable
		if models.HasFeat(char, feat.Name) && !feat.Repeatable {
			continue
		}
		availableFeats = append(availableFeats, feat)
	}
	f.allFeats = availableFeats

	// Apply category filter
	f.applyFilter()

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

// applyFilter filters feats based on the current category filter
func (f *FeatSelector) applyFilter() {
	if f.categoryFilter == "All" {
		f.feats = f.allFeats
	} else {
		filtered := []models.Feat{}
		for _, feat := range f.allFeats {
			if feat.Category == f.categoryFilter {
				filtered = append(filtered, feat)
			}
		}
		f.feats = filtered
	}

	// Reset selection if out of bounds
	if f.selectedIndex >= len(f.feats) {
		f.selectedIndex = 0
	}
}

// NextCategory cycles to the next category filter
func (f *FeatSelector) NextCategory() {
	for i, cat := range f.categories {
		if cat == f.categoryFilter {
			f.categoryFilter = f.categories[(i+1)%len(f.categories)]
			break
		}
	}
	f.applyFilter()
}

// PrevCategory cycles to the previous category filter
func (f *FeatSelector) PrevCategory() {
	for i, cat := range f.categories {
		if cat == f.categoryFilter {
			newIndex := (i - 1 + len(f.categories)) % len(f.categories)
			f.categoryFilter = f.categories[newIndex]
			break
		}
	}
	f.applyFilter()
}

// CanSelectCurrentFeat returns true if the currently selected feat can be taken
func (f *FeatSelector) CanSelectCurrentFeat() bool {
	if f.selectedIndex < 0 || f.selectedIndex >= len(f.feats) {
		return false
	}

	if f.deleteMode {
		return true // In delete mode, all feats shown can be selected
	}

	feat := f.feats[f.selectedIndex]
	return models.CanTakeFeat(f.character, feat)
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
	case "left", "h":
		f.PrevCategory()
	case "right", "l":
		f.NextCategory()
	case "esc":
		f.Hide()
	}

	return nil
}

// View renders the feat selector with two-column layout
func (f *FeatSelector) View(width, height int) string {
	if !f.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	featNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	selectedFeatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalFeatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	disabledFeatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Strikethrough(true)

	selectedDisabledFeatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Strikethrough(true).
		Background(lipgloss.Color("237"))

	categoryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Italic(true)

	prerequisiteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	prerequisiteUnmetStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	abilityStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	benefitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content []string
	content = append(content, titleStyle.Render(f.title))

	// Category filter display
	if !f.deleteMode {
		categoryFilter := fmt.Sprintf("Category: ← %s → (%d feats)", f.categoryFilter, len(f.feats))
		content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Render(categoryFilter))
	}
	content = append(content, "")

	// Left side: Feat list (scrollable)
	var featList []string

	// Calculate visible range for scrolling
	visibleHeight := 25 // Number of feats visible at once
	visibleStart := f.selectedIndex - visibleHeight/2
	if visibleStart < 0 {
		visibleStart = 0
	}
	visibleEnd := visibleStart + visibleHeight
	if visibleEnd > len(f.feats) {
		visibleEnd = len(f.feats)
		visibleStart = visibleEnd - visibleHeight
		if visibleStart < 0 {
			visibleStart = 0
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		feat := f.feats[i]

		// Check if character can take this feat
		canTake := f.deleteMode || models.CanTakeFeat(f.character, feat)

		featLine := fmt.Sprintf(" %s", feat.Name)

		// Add indicator for prerequisites
		if feat.Prerequisite != "None" && feat.Prerequisite != "" {
			if canTake {
				featLine += " ✓"
			} else {
				featLine += " 🔒"
			}
		}

		// Apply appropriate style
		if i == f.selectedIndex {
			if canTake {
				featList = append(featList, selectedFeatStyle.Render(featLine))
			} else {
				featList = append(featList, selectedDisabledFeatStyle.Render(featLine))
			}
		} else {
			if canTake {
				featList = append(featList, normalFeatStyle.Render(featLine))
			} else {
				featList = append(featList, disabledFeatStyle.Render(featLine))
			}
		}
	}

	// Add scroll indicator
	if len(f.feats) > visibleHeight {
		scrollPercent := int(float64(f.selectedIndex) / float64(len(f.feats)-1) * 100)
		scrollInfo := fmt.Sprintf(" %d/%d (%d%%)", f.selectedIndex+1, len(f.feats), scrollPercent)
		featList = append(featList, "")
		featList = append(featList, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(scrollInfo))
	}

	// Right side: Selected feat details
	selectedFeat := f.GetSelectedFeat()
	var featDetails []string
	if selectedFeat != nil {
		// Feat name
		featDetails = append(featDetails, featNameStyle.Render(selectedFeat.Name))
		featDetails = append(featDetails, "")

		// Category and repeatability
		categoryInfo := selectedFeat.Category
		if selectedFeat.Repeatable {
			categoryInfo += " (Repeatable)"
		}
		featDetails = append(featDetails, categoryStyle.Render(categoryInfo))
		featDetails = append(featDetails, "")

		// Prerequisite
		if selectedFeat.Prerequisite != "None" && selectedFeat.Prerequisite != "" {
			canTake := f.deleteMode || models.CanTakeFeat(f.character, *selectedFeat)
			prereqText := "Prerequisite: " + selectedFeat.Prerequisite
			if canTake {
				featDetails = append(featDetails, prerequisiteStyle.Render(prereqText+" ✓"))
			} else {
				featDetails = append(featDetails, prerequisiteUnmetStyle.Render(prereqText+" ✗ NOT MET"))
			}
			featDetails = append(featDetails, "")
		}

		// Ability increases
		if selectedFeat.AbilityIncreases != nil {
			var abilityText string
			if len(selectedFeat.AbilityIncreases.Choices) > 0 {
				abilityText = fmt.Sprintf("+%d to one: %s",
					selectedFeat.AbilityIncreases.Amount,
					strings.Join(selectedFeat.AbilityIncreases.Choices, " or "))
			} else if selectedFeat.AbilityIncreases.Ability != "" {
				abilityText = fmt.Sprintf("+%d %s",
					selectedFeat.AbilityIncreases.Amount,
					selectedFeat.AbilityIncreases.Ability)
			}
			if abilityText != "" {
				featDetails = append(featDetails, abilityStyle.Render(abilityText))
				featDetails = append(featDetails, "")
			}
		}

		// Skill proficiencies
		if len(selectedFeat.SkillProficiencies) > 0 {
			featDetails = append(featDetails, abilityStyle.Render("Skill Proficiency: "+strings.Join(selectedFeat.SkillProficiencies, ", ")))
			featDetails = append(featDetails, "")
		}

		// Languages
		if len(selectedFeat.Languages) > 0 {
			featDetails = append(featDetails, abilityStyle.Render("Languages: "+strings.Join(selectedFeat.Languages, ", ")))
			featDetails = append(featDetails, "")
		}

		// Description
		featDetails = append(featDetails, descStyle.Render(selectedFeat.Description))
		featDetails = append(featDetails, "")

		// Benefits
		if len(selectedFeat.Benefits) > 0 {
			featDetails = append(featDetails, sectionTitleStyle.Render("BENEFITS:"))
			for _, benefit := range selectedFeat.Benefits {
				featDetails = append(featDetails, "")
				// Wrap benefit text
				wrapped := wrapFeatText(benefit, 50)
				for _, line := range wrapped {
					featDetails = append(featDetails, benefitStyle.Render("• "+line))
				}
			}
		}
	}

	// Combine list and details side by side
	listBox := lipgloss.NewStyle().
		Width(30).
		Height(28).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Render(strings.Join(featList, "\n"))

	detailsBox := lipgloss.NewStyle().
		Width(55).
		Height(28).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Render(strings.Join(featDetails, "\n"))

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listBox,
		" ",
		detailsBox,
	)

	content = append(content, mainContent)
	content = append(content, "")

	// Help text
	if f.deleteMode {
		content = append(content, helpStyle.Render("[↑/↓] Navigate • [PgUp/PgDn] Scroll • [Enter] Remove • [Esc] Cancel"))
	} else {
		content = append(content, helpStyle.Render("[↑/↓] Navigate • [←/→] Category • [PgUp/PgDn] Scroll • [Enter] Select • [Esc] Cancel"))
		content = append(content, helpStyle.Render("✓ = Prerequisites met  •  🔒 = Prerequisites NOT met (unselectable)"))
	}

	// Wrap in a styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}

// wrapFeatText wraps text to a specified width for feat selector
func wrapFeatText(text string, width int) []string {
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
