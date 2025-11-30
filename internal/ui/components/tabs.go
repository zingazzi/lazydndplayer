// internal/ui/components/tabs.go
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
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
			{Label: "Stats", Key: "1"},
			{Label: "Skills", Key: "2"},
			{Label: "Inventory", Key: "3"},
			{Label: "Spells", Key: "4"},
			{Label: "Features", Key: "5"},
			{Label: "Traits", Key: "6"},
			{Label: "Origin", Key: "7"},
			{Label: "Companion", Key: "8"},
		},
		SelectedIndex: 0,
	}
}

// UpdateTabsForCharacter updates tabs based on character class
// This allows conditional tabs like Wild Shape for Druids and Companion for Rangers/companions
func (t *Tabs) UpdateTabsForCharacter(char *models.Character) {
	baseTabs := []Tab{
		{Label: "Stats", Key: "1"},
		{Label: "Skills", Key: "2"},
		{Label: "Inventory", Key: "3"},
		{Label: "Spells", Key: "4"},
		{Label: "Features", Key: "5"},
		{Label: "Traits", Key: "6"},
		{Label: "Origin", Key: "7"},
	}

	// Add Companion tab if character has companions or is a Ranger
	hasCompanions := len(char.Companions) > 0 || char.HasClass("Ranger")
	if hasCompanions {
		baseTabs = append(baseTabs, Tab{Label: "Companion", Key: "8"})
	}

	// Add Wild Shape tab if character is a Druid (level 2+)
	hasDruid := char.HasClass("Druid") && char.GetClassLevel("Druid") >= 2
	if hasDruid {
		key := "8"
		if hasCompanions {
			key = "9"
		}
		baseTabs = append(baseTabs, Tab{Label: "Wild Shape", Key: key})
	}

	// Preserve selected index if possible
	oldSelectedIndex := t.SelectedIndex
	t.Items = baseTabs

	// Adjust selected index if it's out of bounds
	if oldSelectedIndex >= len(t.Items) {
		t.SelectedIndex = len(t.Items) - 1
		if t.SelectedIndex < 0 {
			t.SelectedIndex = 0
		}
	} else {
		t.SelectedIndex = oldSelectedIndex
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

	return lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("240")).
		Render(tabBar)
}
