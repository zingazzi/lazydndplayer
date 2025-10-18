// internal/ui/components/classselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ClassInfo represents a D&D 5e class
type ClassInfo struct {
	Name          string
	Description   string
	HitDie        int
	PrimaryAbility string
}

// ClassSelector is a component for selecting a class
type ClassSelector struct {
	classes       []ClassInfo
	selectedIndex int
	visible       bool
}

// NewClassSelector creates a new class selector
func NewClassSelector() *ClassSelector {
	// Hard-coded list of 2024 D&D classes
	classes := []ClassInfo{
		{Name: "Barbarian", Description: "A fierce warrior of primitive background who can enter a battle rage", HitDie: 12, PrimaryAbility: "Strength"},
		{Name: "Bard", Description: "An inspiring magician whose power echoes the music of creation", HitDie: 8, PrimaryAbility: "Charisma"},
		{Name: "Cleric", Description: "A priestly champion who wields divine magic in service of a higher power", HitDie: 8, PrimaryAbility: "Wisdom"},
		{Name: "Druid", Description: "A priest of the Old Faith, wielding the powers of nature and adopting animal forms", HitDie: 8, PrimaryAbility: "Wisdom"},
		{Name: "Fighter", Description: "A master of martial combat, skilled with a variety of weapons and armor", HitDie: 10, PrimaryAbility: "Str/Dex"},
		{Name: "Monk", Description: "A master of martial arts, harnessing the power of the body", HitDie: 8, PrimaryAbility: "Dex/Wis"},
		{Name: "Paladin", Description: "A holy warrior bound to a sacred oath", HitDie: 10, PrimaryAbility: "Str/Cha"},
		{Name: "Ranger", Description: "A warrior who uses martial prowess and nature magic", HitDie: 10, PrimaryAbility: "Dex/Wis"},
		{Name: "Rogue", Description: "A scoundrel who uses stealth and trickery to overcome obstacles", HitDie: 8, PrimaryAbility: "Dexterity"},
		{Name: "Sorcerer", Description: "A spellcaster who draws on inherent magic from a gift or bloodline", HitDie: 6, PrimaryAbility: "Charisma"},
		{Name: "Warlock", Description: "A wielder of magic derived from a bargain with an extraplanar entity", HitDie: 8, PrimaryAbility: "Charisma"},
		{Name: "Wizard", Description: "A scholarly magic-user capable of manipulating the structures of reality", HitDie: 6, PrimaryAbility: "Intelligence"},
	}

	return &ClassSelector{
		classes:       classes,
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the class selector
func (c *ClassSelector) Show() {
	c.visible = true
	c.selectedIndex = 0
}

// Hide closes the class selector
func (c *ClassSelector) Hide() {
	c.visible = false
}

// IsVisible returns whether the selector is visible
func (c *ClassSelector) IsVisible() bool {
	return c.visible
}

// Next moves to the next class
func (c *ClassSelector) Next() {
	if c.selectedIndex < len(c.classes)-1 {
		c.selectedIndex++
	}
}

// Prev moves to the previous class
func (c *ClassSelector) Prev() {
	if c.selectedIndex > 0 {
		c.selectedIndex--
	}
}

// GetSelectedClass returns the currently selected class name
func (c *ClassSelector) GetSelectedClass() string {
	if c.selectedIndex >= 0 && c.selectedIndex < len(c.classes) {
		return c.classes[c.selectedIndex].Name
	}
	return ""
}

// View renders the class selector
func (c *ClassSelector) View(width, height int) string {
	if !c.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Width(70)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	hitDieStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	abilityStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	// Build content
	var content []string
	content = append(content, titleStyle.Render("SELECT CLASS"))
	content = append(content, "")

	// List all classes
	for i, class := range c.classes {
		className := fmt.Sprintf(" %s", class.Name)
		if i == c.selectedIndex {
			content = append(content, selectedStyle.Render(className))
		} else {
			content = append(content, normalStyle.Render(className))
		}
	}

	content = append(content, "")
	content = append(content, strings.Repeat("─", 70))
	content = append(content, "")

	// Show details of selected class
	if c.selectedIndex >= 0 && c.selectedIndex < len(c.classes) {
		selectedClass := c.classes[c.selectedIndex]
		content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render(selectedClass.Name))
		content = append(content, "")
		content = append(content, descStyle.Render(selectedClass.Description))
		content = append(content, "")
		content = append(content, hitDieStyle.Render(fmt.Sprintf("Hit Die: d%d", selectedClass.HitDie)))
		content = append(content, abilityStyle.Render(fmt.Sprintf("Primary Ability: %s", selectedClass.PrimaryAbility)))
	}

	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	// Wrap in a box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(76)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(strings.Join(content, "\n")),
	)
}

