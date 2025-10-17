// internal/ui/components/skillselector.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// D&D 5e skills
var dndSkills = []string{
	"Acrobatics",
	"Animal Handling",
	"Arcana",
	"Athletics",
	"Deception",
	"History",
	"Insight",
	"Intimidation",
	"Investigation",
	"Medicine",
	"Nature",
	"Perception",
	"Performance",
	"Persuasion",
	"Religion",
	"Sleight of Hand",
	"Stealth",
	"Survival",
}

// SkillSelector handles skill selection UI
type SkillSelector struct {
	skills        []string
	selectedIndex int
	viewport      viewport.Model
	visible       bool
}

// NewSkillSelector creates a new skill selector
func NewSkillSelector() *SkillSelector {
	return &SkillSelector{
		skills:        dndSkills,
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the skill selector
func (ss *SkillSelector) Show() {
	ss.visible = true
	ss.selectedIndex = 0
}

// Hide hides the skill selector
func (ss *SkillSelector) Hide() {
	ss.visible = false
}

// IsVisible returns whether the selector is visible
func (ss *SkillSelector) IsVisible() bool {
	return ss.visible
}

// Next moves to the next skill
func (ss *SkillSelector) Next() {
	if ss.selectedIndex < len(ss.skills)-1 {
		ss.selectedIndex++
		ss.viewport.LineDown(1)
	}
}

// Prev moves to the previous skill
func (ss *SkillSelector) Prev() {
	if ss.selectedIndex > 0 {
		ss.selectedIndex--
		ss.viewport.LineUp(1)
	}
}

// GetSelectedSkill returns the currently selected skill
func (ss *SkillSelector) GetSelectedSkill() string {
	if ss.selectedIndex >= 0 && ss.selectedIndex < len(ss.skills) {
		return ss.skills[ss.selectedIndex]
	}
	return ""
}

// View renders the skill selector
func (ss *SkillSelector) View(screenWidth, screenHeight int) string {
	if !ss.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render("SELECT SKILL PROFICIENCY"))
	content = append(content, "")

	// Skill list
	var skillList []string
	for i, skill := range ss.skills {
		skillLine := " " + skill
		if i == ss.selectedIndex {
			skillList = append(skillList, selectedStyle.Render(skillLine))
		} else {
			skillList = append(skillList, normalStyle.Render(skillLine))
		}
	}

	// Create viewport if needed
	listHeight := 20
	if ss.viewport.Width == 0 {
		ss.viewport = viewport.New(40, listHeight)
		ss.viewport.Style = lipgloss.NewStyle()
	}

	ss.viewport.SetContent(strings.Join(skillList, "\n"))
	content = append(content, ss.viewport.View())
	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// Update handles viewport updates
func (ss *SkillSelector) Update(msg tea.Msg) {
	var cmd tea.Cmd
	ss.viewport, cmd = ss.viewport.Update(msg)
	_ = cmd
}

