// internal/ui/components/backstoryeditor.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// BackstoryEditor provides a large text editor for character backstory
type BackstoryEditor struct {
	visible  bool
	textarea textarea.Model
}

// NewBackstoryEditor creates a new backstory editor
func NewBackstoryEditor() *BackstoryEditor {
	ta := textarea.New()
	ta.Placeholder = "Enter your character's backstory..."
	ta.Focus()
	ta.CharLimit = 2000
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.ShowLineNumbers = false

	return &BackstoryEditor{
		visible:  false,
		textarea: ta,
	}
}

// Show displays the editor with existing text
func (b *BackstoryEditor) Show(currentText string) {
	b.visible = true
	b.textarea.SetValue(currentText)
	b.textarea.Focus()
}

// Hide closes the editor
func (b *BackstoryEditor) Hide() {
	b.visible = false
}

// IsVisible returns whether the editor is currently visible
func (b *BackstoryEditor) IsVisible() bool {
	return b.visible
}

// SetValue sets the textarea value
func (b *BackstoryEditor) SetValue(value string) {
	b.textarea.SetValue(value)
}

// GetValue returns the current textarea value
func (b *BackstoryEditor) GetValue() string {
	return b.textarea.Value()
}

// Update updates the textarea
func (b *BackstoryEditor) Update(msg interface{}) {
	var cmd interface{}
	b.textarea, cmd = b.textarea.Update(msg)
	_ = cmd
}

// View renders the editor
func (b *BackstoryEditor) View(width, height int) string {
	if !b.visible {
		return ""
	}

	// Calculate dimensions (80% of screen)
	editorWidth := int(float64(width) * 0.8)
	editorHeight := int(float64(height) * 0.8)

	if editorWidth < 60 {
		editorWidth = 60
	}
	if editorHeight < 15 {
		editorHeight = 15
	}

	// Update textarea dimensions
	b.textarea.SetWidth(editorWidth - 6) // Account for border and padding
	b.textarea.SetHeight(editorHeight - 8) // Account for title, hints, border

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	charCount := len(b.textarea.Value())
	charLimit := b.textarea.CharLimit

	var content strings.Builder

	content.WriteString(titleStyle.Render("EDIT BACKSTORY"))
	content.WriteString("\n\n")
	content.WriteString(b.textarea.View())
	content.WriteString("\n\n")
	content.WriteString(hintStyle.Render(fmt.Sprintf("Characters: %d/%d", charCount, charLimit)))
	content.WriteString("\n")
	content.WriteString(hintStyle.Render("Ctrl+Enter: Save | Esc: Cancel"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(editorWidth)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(content.String()),
	)
}

