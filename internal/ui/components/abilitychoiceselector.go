// internal/ui/components/abilitychoiceselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// AbilityChoiceSelector handles selecting an ability for feats with choices
type AbilityChoiceSelector struct {
	visible       bool
	choices       []string
	selectedIndex int
	featName      string
	character     *models.Character
}

// NewAbilityChoiceSelector creates a new ability choice selector
func NewAbilityChoiceSelector() *AbilityChoiceSelector {
	return &AbilityChoiceSelector{
		visible:       false,
		choices:       []string{},
		selectedIndex: 0,
	}
}

// Show displays the ability choice selector
func (a *AbilityChoiceSelector) Show(featName string, choices []string, char *models.Character) {
	a.visible = true
	a.featName = featName
	a.choices = choices
	a.selectedIndex = 0
	a.character = char
}

// Hide hides the ability choice selector
func (a *AbilityChoiceSelector) Hide() {
	a.visible = false
}

// IsVisible returns whether the selector is visible
func (a *AbilityChoiceSelector) IsVisible() bool {
	return a.visible
}

// Next moves to next choice
func (a *AbilityChoiceSelector) Next() {
	if a.selectedIndex < len(a.choices)-1 {
		a.selectedIndex++
	}
}

// Prev moves to previous choice
func (a *AbilityChoiceSelector) Prev() {
	if a.selectedIndex > 0 {
		a.selectedIndex--
	}
}

// GetSelectedAbility returns the currently selected ability
func (a *AbilityChoiceSelector) GetSelectedAbility() string {
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.choices) {
		return a.choices[a.selectedIndex]
	}
	return ""
}

// Update handles input for the ability choice selector
func (a *AbilityChoiceSelector) Update(msg tea.KeyMsg) tea.Cmd {
	if !a.visible {
		return nil
	}

	switch msg.String() {
	case "up", "k":
		a.Prev()
	case "down", "j":
		a.Next()
	case "esc":
		a.Hide()
	}

	return nil
}

// View renders the ability choice selector
func (a *AbilityChoiceSelector) View(screenWidth, screenHeight int) string {
	if !a.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render(fmt.Sprintf("SELECT ABILITY FOR %s", strings.ToUpper(a.featName))))
	content = append(content, "")
	content = append(content, normalStyle.Render("Choose which ability score to increase:"))
	content = append(content, "")

	// Show current ability scores to help with decision
	if a.character != nil {
		for i, choice := range a.choices {
			// Get current score for this ability
			var currentScore int
			choiceLower := strings.ToLower(choice)

			switch {
			case strings.Contains(choiceLower, "strength"):
				currentScore = a.character.AbilityScores.Strength
			case strings.Contains(choiceLower, "dexterity"):
				currentScore = a.character.AbilityScores.Dexterity
			case strings.Contains(choiceLower, "constitution"):
				currentScore = a.character.AbilityScores.Constitution
			case strings.Contains(choiceLower, "intelligence"):
				currentScore = a.character.AbilityScores.Intelligence
			case strings.Contains(choiceLower, "wisdom"):
				currentScore = a.character.AbilityScores.Wisdom
			case strings.Contains(choiceLower, "charisma"):
				currentScore = a.character.AbilityScores.Charisma
			}

			line := fmt.Sprintf("  %s (currently %d)", choice, currentScore)
			if i == a.selectedIndex {
				content = append(content, selectedStyle.Render(line))
			} else {
				content = append(content, normalStyle.Render(line))
			}
		}
	}

	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}
