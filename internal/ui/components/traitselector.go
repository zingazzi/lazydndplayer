// internal/ui/components/traitselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// TraitSelector allows selection of character traits (personality, ideal, bond, flaw)
type TraitSelector struct {
	visible        bool
	TraitType      models.TraitType
	options        []string
	selectedIndex  int
	customMode     bool
	customInput    textinput.Model
	selectedItems  []string // Items that are currently selected (multi-select mode)
}

// NewTraitSelector creates a new trait selector
func NewTraitSelector() *TraitSelector {
	ti := textinput.New()
	ti.Placeholder = "Enter custom trait..."
	ti.CharLimit = 200
	ti.Width = 60

	return &TraitSelector{
		visible:       false,
		selectedIndex: 0,
		customMode:    false,
		customInput:   ti,
	}
}

// Show displays the selector with the specified trait type and current selections
func (t *TraitSelector) Show(traitType models.TraitType, currentSelections []string) {
	t.visible = true
	t.TraitType = traitType

	// Get predefined options
	predefinedOptions := models.GetTraitOptions(traitType)

	// Build options list: predefined options + custom items (excluding "[Custom]" entry)
	t.options = []string{}

	// Add predefined options (except "[Custom]")
	for _, opt := range predefinedOptions {
		if opt != "[Custom]" {
			t.options = append(t.options, opt)
		}
	}

	// Add custom items that are in currentSelections but not in predefined options
	for _, selection := range currentSelections {
		isCustom := true
		for _, predefined := range predefinedOptions {
			if selection == predefined {
				isCustom = false
				break
			}
		}
		if isCustom && selection != "" {
			t.options = append(t.options, selection)
		}
	}

	// Add "[Custom]" option at the end
	t.options = append(t.options, "[Custom]")

	t.selectedIndex = 0
	t.customMode = false
	t.customInput.SetValue("")
	t.selectedItems = make([]string, len(currentSelections))
	copy(t.selectedItems, currentSelections)
}

// Hide closes the selector
func (t *TraitSelector) Hide() {
	t.visible = false
	t.customMode = false
}

// IsVisible returns whether the selector is currently visible
func (t *TraitSelector) IsVisible() bool {
	return t.visible
}

// Next moves selection down
func (t *TraitSelector) Next() {
	if !t.customMode && t.selectedIndex < len(t.options)-1 {
		t.selectedIndex++
	}
}

// Prev moves selection up
func (t *TraitSelector) Prev() {
	if !t.customMode && t.selectedIndex > 0 {
		t.selectedIndex--
	}
}

// ToggleCustomMode switches to custom input mode if [Custom] is selected
func (t *TraitSelector) ToggleCustomMode() {
	if t.options[t.selectedIndex] == "[Custom]" {
		t.customMode = true
		t.customInput.Focus()
	}
}

// GetSelectedTrait returns the selected trait text (for custom mode)
func (t *TraitSelector) GetSelectedTrait() string {
	if t.customMode {
		return t.customInput.Value()
	}
	if t.selectedIndex >= 0 && t.selectedIndex < len(t.options) {
		selected := t.options[t.selectedIndex]
		if selected == "[Custom]" {
			return ""
		}
		return selected
	}
	return ""
}

// GetSelectedItems returns all selected items (for multi-select mode)
func (t *TraitSelector) GetSelectedItems() []string {
	result := make([]string, len(t.selectedItems))
	copy(result, t.selectedItems)
	return result
}

// IsCustomMode returns whether the selector is in custom input mode
func (t *TraitSelector) IsCustomMode() bool {
	return t.customMode
}

// ToggleItem toggles the selection of the current item
func (t *TraitSelector) ToggleItem() {
	if t.selectedIndex < 0 || t.selectedIndex >= len(t.options) {
		return
	}

	currentOption := t.options[t.selectedIndex]
	if currentOption == "[Custom]" {
		return // Don't toggle [Custom] option
	}

	// Check if already selected
	found := false
	foundIndex := -1
	for i, item := range t.selectedItems {
		if item == currentOption {
			found = true
			foundIndex = i
			break
		}
	}

	if found {
		// Remove from selection
		t.selectedItems = append(t.selectedItems[:foundIndex], t.selectedItems[foundIndex+1:]...)
	} else {
		// Add to selection
		t.selectedItems = append(t.selectedItems, currentOption)
	}
}

// IsItemSelected checks if the current item is selected
func (t *TraitSelector) IsItemSelected(item string) bool {
	for _, selected := range t.selectedItems {
		if selected == item {
			return true
		}
	}
	return false
}

// Update handles text input updates
func (t *TraitSelector) Update(msg tea.Msg) {
	if t.customMode {
		var cmd tea.Cmd
		t.customInput, cmd = t.customInput.Update(msg)
		_ = cmd
	}
}

// View renders the selector
func (t *TraitSelector) View(width, height int) string {
	if !t.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	var content strings.Builder

	traitName := models.GetTraitTypeName(t.TraitType)
	content.WriteString(titleStyle.Render(fmt.Sprintf("SELECT %s", strings.ToUpper(traitName))))
	content.WriteString("\n\n")

	if t.customMode {
		// Custom input mode
		content.WriteString(normalStyle.Render("Enter your custom trait:"))
		content.WriteString("\n\n")
		content.WriteString(t.customInput.View())
		content.WriteString("\n\n")
		content.WriteString(hintStyle.Render("Enter: Confirm | Esc: Cancel"))
	} else {
		// Selection mode with multi-select
		for i, option := range t.options {
			cursor := "  "
			checkbox := "☐ "
			style := normalStyle

			if i == t.selectedIndex {
				cursor = "→ "
				style = selectedStyle
			}

			// Check if this item is selected
			if t.IsItemSelected(option) && option != "[Custom]" {
				checkbox = "☑ "
			}

			// Don't show checkbox for [Custom] option
			if option == "[Custom]" {
				checkbox = ""
			}

			// Show full text without truncation
			displayText := option

			content.WriteString(style.Render(cursor + checkbox + displayText))
			content.WriteString("\n")
		}

		content.WriteString("\n")
		content.WriteString(hintStyle.Render("↑/↓: Navigate | Space: Toggle | Enter: Confirm | Esc: Cancel"))
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(80)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}
