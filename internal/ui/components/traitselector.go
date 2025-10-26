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
	visible       bool
	TraitType     models.TraitType
	options       []string
	selectedIndex int
	customMode    bool
	customInput   textinput.Model
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

// Show displays the selector with the specified trait type
func (t *TraitSelector) Show(traitType models.TraitType) {
	t.visible = true
	t.TraitType = traitType
	t.options = models.GetTraitOptions(traitType)
	t.selectedIndex = 0
	t.customMode = false
	t.customInput.SetValue("")
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

// GetSelectedTrait returns the selected trait text
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
		// Selection mode
		for i, option := range t.options {
			cursor := "  "
			style := normalStyle
			if i == t.selectedIndex {
				cursor = "→ "
				style = selectedStyle
			}

			// Truncate long options for display
			displayText := option
			if len(displayText) > 70 {
				displayText = displayText[:67] + "..."
			}

			content.WriteString(style.Render(cursor + displayText))
			content.WriteString("\n")
		}

		content.WriteString("\n")
		content.WriteString(hintStyle.Render("↑/↓: Navigate | Enter: Select | Esc: Cancel"))
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

