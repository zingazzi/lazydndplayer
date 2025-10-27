package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MessagePopup struct {
	visible bool
	message string
	title   string
}

func NewMessagePopup() *MessagePopup {
	return &MessagePopup{
		visible: false,
	}
}

func (m *MessagePopup) Show(title, message string) {
	m.title = title
	m.message = message
	m.visible = true
}

func (m *MessagePopup) Hide() {
	m.visible = false
}

func (m *MessagePopup) IsVisible() bool {
	return m.visible
}

func (m *MessagePopup) Update(msg tea.KeyMsg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg.String() {
	case "enter", "esc":
		m.Hide()
	}

	return nil
}

func (m *MessagePopup) View(width, height int) string {
	if !m.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		MarginBottom(1)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(60).
		Align(lipgloss.Center)

	content := titleStyle.Render(m.title) + "\n\n" +
		messageStyle.Render(m.message) + "\n\n" +
		hintStyle.Render("Press Enter to continue")

	box := boxStyle.Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
