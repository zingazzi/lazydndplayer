// internal/ui/components/sidebar.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a sidebar menu item
type MenuItem struct {
	Label string
	Key   string
}

// Sidebar represents the navigation sidebar
type Sidebar struct {
	Items         []MenuItem
	SelectedIndex int
	Width         int
}

// NewSidebar creates a new sidebar
func NewSidebar() *Sidebar {
	return &Sidebar{
		Items: []MenuItem{
			{Label: "Overview", Key: "1"},
			{Label: "Stats", Key: "2"},
			{Label: "Skills", Key: "3"},
			{Label: "Inventory", Key: "4"},
			{Label: "Spells", Key: "5"},
			{Label: "Actions", Key: "6"},
			{Label: "Dice", Key: "7"},
		},
		SelectedIndex: 0,
		Width:         15,
	}
}

// Next moves to the next item
func (s *Sidebar) Next() {
	s.SelectedIndex = (s.SelectedIndex + 1) % len(s.Items)
}

// Prev moves to the previous item
func (s *Sidebar) Prev() {
	s.SelectedIndex--
	if s.SelectedIndex < 0 {
		s.SelectedIndex = len(s.Items) - 1
	}
}

// GetSelected returns the currently selected item
func (s *Sidebar) GetSelected() MenuItem {
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Items) {
		return s.Items[s.SelectedIndex]
	}
	return s.Items[0]
}

// View renders the sidebar
func (s *Sidebar) View(height int) string {
	var items []string

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(s.Width).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("237")).
		Width(s.Width).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	for i, item := range s.Items {
		label := fmt.Sprintf("[%s] %s", item.Key, item.Label)
		if i == s.SelectedIndex {
			items = append(items, selectedStyle.Render(label))
		} else {
			items = append(items, normalStyle.Render(label))
		}
	}

	// Add spacing
	for len(items) < height-5 {
		items = append(items, normalStyle.Render(""))
	}

	// Add shortcuts at bottom
	items = append(items, "")
	items = append(items, normalStyle.Render(keyStyle.Render("[q]")+" Quit"))
	items = append(items, normalStyle.Render(keyStyle.Render("[?]")+" Help"))
	items = append(items, normalStyle.Render(keyStyle.Render("[s]")+" Save"))

	sidebarStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(lipgloss.Color("240")).
		Width(s.Width + 2).
		Height(height - 3)

	return sidebarStyle.Render(strings.Join(items, "\n"))
}
