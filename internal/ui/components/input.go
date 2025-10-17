// internal/ui/components/input.go
package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputDialog represents an input dialog
type InputDialog struct {
	Title       string
	Placeholder string
	Input       textinput.Model
	Active      bool
	Value       string
}

// NewInputDialog creates a new input dialog
func NewInputDialog(title, placeholder string) InputDialog {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 100
	ti.Width = 40

	return InputDialog{
		Title:       title,
		Placeholder: placeholder,
		Input:       ti,
		Active:      false,
	}
}

// Activate activates the input dialog
func (d *InputDialog) Activate() {
	d.Active = true
	d.Input.Focus()
	d.Input.SetValue("")
}

// Deactivate deactivates the input dialog
func (d *InputDialog) Deactivate() {
	d.Active = false
	d.Input.Blur()
	d.Value = d.Input.Value()
}

// Update updates the input dialog
func (d *InputDialog) Update(msg tea.Msg) tea.Cmd {
	if !d.Active {
		return nil
	}

	var cmd tea.Cmd
	d.Input, cmd = d.Input.Update(msg)
	return cmd
}

// View renders the input dialog
func (d *InputDialog) View() string {
	if !d.Active {
		return ""
	}

	return d.Title + "\n" + d.Input.View()
}
