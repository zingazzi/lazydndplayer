// internal/ui/components/expertiseselector.go
package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ExpertiseSelector allows selecting skills for expertise
type ExpertiseSelector struct {
	character      *models.Character
	visible        bool
	selectedIndex  int
	maxExpertise   int // Maximum number of skills that can have expertise
	availableSkills []expertiseOption
	selectedSkills map[string]bool // Track which skills are selected for expertise
}

type expertiseOption struct {
	name        string
	ability     string
	proficiency models.ProficiencyLevel
}

// NewExpertiseSelector creates a new expertise selector
func NewExpertiseSelector(char *models.Character) *ExpertiseSelector {
	return &ExpertiseSelector{
		character:       char,
		visible:         false,
		selectedIndex:   0,
		selectedSkills:  make(map[string]bool),
	}
}

// Show displays the selector with the specified maximum number of expertise skills
func (s *ExpertiseSelector) Show(maxExpertise int) {
	s.visible = true
	s.selectedIndex = 0
	s.maxExpertise = maxExpertise
	s.selectedSkills = make(map[string]bool)

	// Initialize with currently expertised skills
	for _, skill := range s.character.Skills.List {
		if skill.Proficiency == models.Expertise {
			s.selectedSkills[string(skill.Name)] = true
		}
	}

	s.buildAvailableSkills()
}

// Hide closes the selector
func (s *ExpertiseSelector) Hide() {
	s.visible = false
}

// IsVisible returns whether the selector is visible
func (s *ExpertiseSelector) IsVisible() bool {
	return s.visible
}

// buildAvailableSkills builds the list of skills that can have expertise
func (s *ExpertiseSelector) buildAvailableSkills() {
	s.availableSkills = []expertiseOption{}

	// Get all skills that are proficient (not NotProficient or already Expertise)
	for _, skill := range s.character.Skills.List {
		if skill.Proficiency == models.Proficient {
			s.availableSkills = append(s.availableSkills, expertiseOption{
				name:        string(skill.Name),
				ability:     string(skill.Ability),
				proficiency: skill.Proficiency,
			})
		}
	}
}

// ToggleSelection toggles the selection of the current skill
func (s *ExpertiseSelector) ToggleSelection() bool {
	if s.selectedIndex < 0 || s.selectedIndex >= len(s.availableSkills) {
		return false
	}

	skillName := s.availableSkills[s.selectedIndex].name

	if s.selectedSkills[skillName] {
		// Deselect
		delete(s.selectedSkills, skillName)
		return true
	} else {
		// Check if we can select more
		if len(s.selectedSkills) < s.maxExpertise {
			s.selectedSkills[skillName] = true
			return true
		}
		return false
	}
}

// Next moves to the next skill
func (s *ExpertiseSelector) Next() {
	if s.selectedIndex < len(s.availableSkills)-1 {
		s.selectedIndex++
	}
}

// Prev moves to the previous skill
func (s *ExpertiseSelector) Prev() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// CanConfirm returns true if the selection is valid
func (s *ExpertiseSelector) CanConfirm() bool {
	return len(s.selectedSkills) == s.maxExpertise
}

// GetSelectedSkills returns the currently selected skills
func (s *ExpertiseSelector) GetSelectedSkills() map[string]bool {
	return s.selectedSkills
}

// GetMaxExpertise returns the maximum number of expertise skills
func (s *ExpertiseSelector) GetMaxExpertise() int {
	return s.maxExpertise
}

// ApplySelections applies the expertise selections to the character
func (s *ExpertiseSelector) ApplySelections() {
	// Set selected skills to Expertise level
	for skillName := range s.selectedSkills {
		s.character.Skills.SetProficiency(models.SkillType(skillName), models.Expertise)
	}
}

// View renders the expertise selector
func (s *ExpertiseSelector) View(width, height int) string {
	if !s.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Italic(true)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedSkillStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")) // Green for selected skills

	var lines []string
	lines = append(lines, titleStyle.Render("EXPERTISE SELECTION"))
	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render(fmt.Sprintf("Choose %d skill(s) for expertise (double proficiency bonus):", s.maxExpertise)))
	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render(fmt.Sprintf("Selected: %d/%d", len(s.selectedSkills), s.maxExpertise)))
	lines = append(lines, "")

	// Show available skills
	for i, skill := range s.availableSkills {
		marker := "  "
		style := normalStyle

		if s.selectedSkills[skill.name] {
			marker = "✓ "
			style = selectedSkillStyle
		}

		line := fmt.Sprintf("%s%-20s (%s)", marker, skill.name, skill.ability)

		if i == s.selectedIndex {
			lines = append(lines, selectedStyle.Render("▶ "+line))
		} else {
			lines = append(lines, style.Render("  "+line))
		}
	}

	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render("Controls: ↑/↓ Navigate, Space Toggle, Enter Confirm, Esc Cancel"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// Update handles key input for the expertise selector
func (s *ExpertiseSelector) Update(msg tea.Msg) tea.Cmd {
	if !s.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.Prev()
		case "down", "j":
			s.Next()
		case " ":
			s.ToggleSelection()
		case "enter":
			if s.CanConfirm() {
				s.ApplySelections()
				s.Hide()
				return tea.Quit
			}
		case "esc":
			s.Hide()
			return tea.Quit
		}
	}

	return nil
}
