// internal/ui/components/classskillselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ClassSkillSelector handles skill selection when choosing a class
type ClassSkillSelector struct {
	availableSkills   []string        // Skills the class can choose from
	SelectedSkills    map[string]bool // Skills selected in this session (exported for access)
	existingSkills    map[string]bool // Skills character already has
	MaxChoices        int             // How many skills can be chosen (exported for access)
	currentIndex      int             // Current cursor position
	visible           bool
	ClassName         string          // Class name (exported for access)
}

// NewClassSkillSelector creates a new class skill selector
func NewClassSkillSelector() *ClassSkillSelector {
	return &ClassSkillSelector{
		SelectedSkills: make(map[string]bool),
		existingSkills: make(map[string]bool),
		visible:        false,
	}
}

// Show displays the skill selector for a specific class
func (css *ClassSkillSelector) Show(className string, availableSkills []string, maxChoices int, char *models.Character) {
	css.ClassName = className
	css.availableSkills = availableSkills
	css.MaxChoices = maxChoices
	css.currentIndex = 0
	css.SelectedSkills = make(map[string]bool)
	css.existingSkills = make(map[string]bool)

	// Mark existing proficiencies
	allSkills := []models.SkillType{
		models.Acrobatics, models.AnimalHandling, models.Arcana, models.Athletics,
		models.Deception, models.History, models.Insight, models.Intimidation,
		models.Investigation, models.Medicine, models.Nature, models.Perception,
		models.Performance, models.Persuasion, models.Religion, models.SleightOfHand,
		models.Stealth, models.Survival,
	}

	for _, skillType := range allSkills {
		skill := char.Skills.GetSkill(skillType)
		if skill != nil && skill.Proficiency > 0 {
			css.existingSkills[string(skillType)] = true
		}
	}

	css.visible = true
}

// Hide closes the skill selector
func (css *ClassSkillSelector) Hide() {
	css.visible = false
}

// IsVisible returns whether the selector is visible
func (css *ClassSkillSelector) IsVisible() bool {
	return css.visible
}

// Next moves to the next skill
func (css *ClassSkillSelector) Next() {
	if css.currentIndex < len(css.availableSkills)-1 {
		css.currentIndex++
	}
}

// Prev moves to the previous skill
func (css *ClassSkillSelector) Prev() {
	if css.currentIndex > 0 {
		css.currentIndex--
	}
}

// ToggleSkill toggles selection of the current skill
func (css *ClassSkillSelector) ToggleSkill() bool {
	if css.currentIndex < 0 || css.currentIndex >= len(css.availableSkills) {
		return false
	}

	skillName := css.availableSkills[css.currentIndex]

	// Can't select if already proficient
	if css.existingSkills[skillName] {
		return false
	}

	// Toggle selection
	if css.SelectedSkills[skillName] {
		delete(css.SelectedSkills, skillName)
		return true
	}

	// Check if we can select more
	if len(css.SelectedSkills) < css.MaxChoices {
		css.SelectedSkills[skillName] = true
		return true
	}

	return false
}

// GetSelectedSkills returns the list of selected skills
func (css *ClassSkillSelector) GetSelectedSkills() []string {
	skills := []string{}
	for skill := range css.SelectedSkills {
		skills = append(skills, skill)
	}
	return skills
}

// CanConfirm returns true if the correct number of skills have been selected
func (css *ClassSkillSelector) CanConfirm() bool {
	return len(css.SelectedSkills) == css.MaxChoices
}

// Update handles keyboard input for the class skill selector
func (css *ClassSkillSelector) Update(msg tea.Msg) (ClassSkillSelector, tea.Cmd) {
	if !css.visible {
		return *css, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			css.Prev()
		case "down", "j":
			css.Next()
		case " ":
			css.ToggleSkill()
		case "enter", "esc":
			return *css, nil
		}
	}

	return *css, nil
}

// View renders the skill selector
func (css *ClassSkillSelector) View(width, height int) string {
	if !css.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	proficientStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Italic(true)

	chosenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	counterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render(fmt.Sprintf("SELECT SKILLS FOR %s", strings.ToUpper(css.ClassName))))
	content = append(content, "")
	content = append(content, counterStyle.Render(fmt.Sprintf("Choose %d skill(s) • Selected: %d/%d", css.MaxChoices, len(css.SelectedSkills), css.MaxChoices)))
	content = append(content, "")

	// List all available skills
	for i, skill := range css.availableSkills {
		var line string
		alreadyProficient := css.existingSkills[skill]
		alreadySelected := css.SelectedSkills[skill]

		// Build the line
		if alreadyProficient {
			line = fmt.Sprintf("  [✓] %s (Already Proficient)", skill)
			if i == css.currentIndex {
				content = append(content, selectedStyle.Render(line))
			} else {
				content = append(content, proficientStyle.Render(line))
			}
		} else if alreadySelected {
			line = fmt.Sprintf("  [✓] %s", skill)
			if i == css.currentIndex {
				content = append(content, selectedStyle.Render(line))
			} else {
				content = append(content, chosenStyle.Render(line))
			}
		} else {
			line = fmt.Sprintf("  [ ] %s", skill)
			if i == css.currentIndex {
				content = append(content, selectedStyle.Render(line))
			} else {
				content = append(content, normalStyle.Render(line))
			}
		}
	}

	content = append(content, "")
	content = append(content, helpStyle.Render("Legend:"))
	content = append(content, helpStyle.Render("  [✓] = Selected or Already Proficient"))
	content = append(content, helpStyle.Render("  [ ] = Available"))
	content = append(content, "")

	if css.CanConfirm() {
		content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("[↑/↓] Navigate • [Space] Toggle • [Enter] Confirm • [Esc] Cancel"))
	} else {
		content = append(content, helpStyle.Render(fmt.Sprintf("[↑/↓] Navigate • [Space] Toggle • Select %d more skill(s) • [Esc] Cancel", css.MaxChoices-len(css.SelectedSkills))))
	}

	// Wrap in a box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(80)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(strings.Join(content, "\n")),
	)
}
