// internal/ui/components/inputpopup.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputPopup provides a simple popup for text input
type InputPopup struct {
	visible     bool
	title       string
	textinput   textinput.Model
	placeholder string
}

// NewInputPopup creates a new input popup
func NewInputPopup() *InputPopup {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Width = 60

	return &InputPopup{
		visible:   false,
		textinput: ti,
	}
}

// Show displays the popup with title and initial value
func (i *InputPopup) Show(title, initialValue, placeholder string) {
	i.visible = true
	i.title = title
	i.placeholder = placeholder
	i.textinput.Placeholder = placeholder
	i.textinput.SetValue(initialValue)
	i.textinput.Focus()
	i.textinput.CursorEnd()
}

// Hide closes the popup
func (i *InputPopup) Hide() {
	i.visible = false
	i.textinput.Blur()
}

// IsVisible returns whether the popup is currently visible
func (i *InputPopup) IsVisible() bool {
	return i.visible
}

// GetValue returns the current input value
func (i *InputPopup) GetValue() string {
	return i.textinput.Value()
}

// Update updates the textinput
func (i *InputPopup) Update(msg tea.Msg) {
	var cmd tea.Cmd
	i.textinput, cmd = i.textinput.Update(msg)
	_ = cmd
}

// View renders the popup
func (i *InputPopup) View(width, height int) string {
	if !i.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	var content strings.Builder

	content.WriteString(titleStyle.Render(i.title))
	content.WriteString("\n\n")
	content.WriteString(i.textinput.View())
	content.WriteString("\n\n")
	content.WriteString(hintStyle.Render("Enter: Confirm | Esc: Cancel"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(70)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

