// internal/ui/components/tabs.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab
type Tab struct {
	Label string
	Key   string
}

// Tabs represents a tab navigation component
type Tabs struct {
	Items         []Tab
	SelectedIndex int
}

// NewTabs creates a new tabs component
func NewTabs() *Tabs {
	return &Tabs{
		Items: []Tab{
			{Label: "Overview", Key: "1"},
			{Label: "Stats", Key: "2"},
			{Label: "Skills", Key: "3"},
			{Label: "Inventory", Key: "4"},
			{Label: "Spells", Key: "5"},
		},
		SelectedIndex: 0,
	}
}

// Next moves to the next tab
func (t *Tabs) Next() {
	t.SelectedIndex = (t.SelectedIndex + 1) % len(t.Items)
}

// Prev moves to the previous tab
func (t *Tabs) Prev() {
	t.SelectedIndex--
	if t.SelectedIndex < 0 {
		t.SelectedIndex = len(t.Items) - 1
	}
}

// SetIndex sets the selected tab index
func (t *Tabs) SetIndex(index int) {
	if index >= 0 && index < len(t.Items) {
		t.SelectedIndex = index
	}
}

// GetSelected returns the currently selected tab
func (t *Tabs) GetSelected() Tab {
	if t.SelectedIndex >= 0 && t.SelectedIndex < len(t.Items) {
		return t.Items[t.SelectedIndex]
	}
	return t.Items[0]
}

// View renders the tabs
func (t *Tabs) View(width int) string {
	activeTabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("205")).
		Padding(0, 2)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("237")).
		Padding(0, 2)

	tabSeparator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(" ")

	var tabs []string
	for i, tab := range t.Items {
		if i == t.SelectedIndex {
			tabs = append(tabs, activeTabStyle.Render(tab.Label))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(tab.Label))
		}
	}

	tabBar := strings.Join(tabs, tabSeparator)

	// Add keyboard hints
	hints := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 2).
		Render("[f] Focus • [Tab] Switch • [s] Save • [?] Help • [q] Quit")

	// Combine tab bar and hints
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabBar,
		lipgloss.NewStyle().Width(width-lipgloss.Width(tabBar)-lipgloss.Width(hints)).Render(""),
		hints,
	)

	return lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		Render(combined)
}
