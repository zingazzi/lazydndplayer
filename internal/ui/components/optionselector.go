// internal/ui/components/optionselector.go
package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OptionSelector handles simple option selection from a list
type OptionSelector struct {
	visible       bool
	title         string
	options       []string
	selectedIndex int
	selectedValue string
}

// NewOptionSelector creates a new option selector
func NewOptionSelector() *OptionSelector {
	return &OptionSelector{
		visible:       false,
		selectedIndex: 0,
	}
}

// Show displays the option selector with a title and list of options
func (os *OptionSelector) Show(title string, options []string) {
	os.visible = true
	os.title = title
	os.options = options
	os.selectedIndex = 0
	os.selectedValue = ""
}

// Hide hides the option selector
func (os *OptionSelector) Hide() {
	os.visible = false
	os.selectedIndex = 0
	os.selectedValue = ""
}

// IsVisible returns whether the selector is visible
func (os *OptionSelector) IsVisible() bool {
	return os.visible
}

// GetSelectedValue returns the currently selected option
func (os *OptionSelector) GetSelectedValue() string {
	return os.selectedValue
}

// Next moves to the next option
func (os *OptionSelector) Next() {
	if os.selectedIndex < len(os.options)-1 {
		os.selectedIndex++
	}
}

// Prev moves to the previous option
func (os *OptionSelector) Prev() {
	if os.selectedIndex > 0 {
		os.selectedIndex--
	}
}

// Select confirms the current selection
func (os *OptionSelector) Select() bool {
	if os.selectedIndex >= 0 && os.selectedIndex < len(os.options) {
		os.selectedValue = os.options[os.selectedIndex]
		return true
	}
	return false
}

// Update handles key presses
func (os *OptionSelector) Update(msg tea.Msg) (OptionSelector, tea.Cmd) {
	if !os.visible {
		return *os, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			os.Prev()
		case "down", "j":
			os.Next()
		case "enter":
			if os.Select() {
				os.visible = false
			}
		case "esc":
			os.selectedValue = ""
			os.Hide()
		}
	}

	return *os, nil
}

// View renders the option selector
func (os *OptionSelector) View(width, height int) string {
	if !os.visible || len(os.options) == 0 {
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

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	var content strings.Builder

	content.WriteString(titleStyle.Render(os.title))
	content.WriteString("\n\n")

	for i, option := range os.options {
		cursor := "  "
		style := normalStyle

		if i == os.selectedIndex {
			cursor = "→ "
			style = selectedStyle
		}

		// Wrap long options
		displayText := option
		if len(displayText) > 70 {
			displayText = displayText[:67] + "..."
		}

		content.WriteString(style.Render(cursor + displayText))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(hintStyle.Render("↑/↓: Navigate | Enter: Select | Esc: Cancel"))

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
