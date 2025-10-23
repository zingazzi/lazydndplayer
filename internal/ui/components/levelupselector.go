// internal/ui/components/levelupselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// LevelUpState represents the current state of the level-up process
type LevelUpState int

const (
	LevelUpSelectClass LevelUpState = iota
	LevelUpConfirm
	LevelUpSelectSkills
	LevelUpSelectSubclass
	LevelUpSelectFightingStyle
	LevelUpComplete
)

// LevelUpSelector handles the level-up process
type LevelUpSelector struct {
	visible           bool
	character         *models.Character
	state             LevelUpState
	selectedClass     string
	availableClasses  []models.Class
	cursor            int
	preview           *models.LevelUpResult
	takeAverage       bool // For HP rolling
	selectedSkills    []models.SkillType
	selectedSubclass  string
	message           string
	backup            *models.Character // For rollback
	subclassSelector  *SubclassSelector
}

// NewLevelUpSelector creates a new level-up selector
func NewLevelUpSelector(char *models.Character) *LevelUpSelector {
	return &LevelUpSelector{
		character:        char,
		takeAverage:      true, // Default to taking average HP
		subclassSelector: NewSubclassSelector(char),
	}
}

// Show displays the level-up selector
func (ls *LevelUpSelector) Show() {
	ls.visible = true
	ls.state = LevelUpSelectClass
	ls.cursor = 0
	ls.message = ""

	// Create backup for rollback
	ls.backup = &models.Character{}
	*ls.backup = *ls.character

	// Load available classes (filtered by multiclass prerequisites)
	ls.loadAvailableClasses()
}

// Hide hides the level-up selector
func (ls *LevelUpSelector) Hide() {
	ls.visible = false
	ls.state = LevelUpSelectClass
	ls.selectedClass = ""
	ls.preview = nil
	ls.selectedSkills = nil
}

// IsVisible returns whether the selector is visible
func (ls *LevelUpSelector) IsVisible() bool {
	return ls.visible
}

// loadAvailableClasses loads classes that can be leveled up
func (ls *LevelUpSelector) loadAvailableClasses() {
	ls.availableClasses = []models.Class{}

	// Add existing classes first (can always level up in current classes)
	for _, cl := range ls.character.Classes {
		classData := models.GetClassByName(cl.ClassName)
		if classData != nil {
			ls.availableClasses = append(ls.availableClasses, *classData)
		}
	}

	// Add new classes that meet prerequisites
	availableNew := models.GetAvailableClasses(ls.character)
	for _, newClass := range availableNew {
		// Check if not already in character's classes
		hasClass := false
		for _, cl := range ls.character.Classes {
			if cl.ClassName == newClass.Name {
				hasClass = true
				break
			}
		}
		if !hasClass {
			ls.availableClasses = append(ls.availableClasses, newClass)
		}
	}
}

// Update handles key presses
func (ls *LevelUpSelector) Update(msg tea.Msg) (LevelUpSelector, tea.Cmd) {
	if !ls.visible {
		return *ls, nil
	}

	// If subclass selector is visible, delegate to it
	if ls.subclassSelector.IsVisible() {
		updated, cmd := ls.subclassSelector.Update(msg)
		ls.subclassSelector = &updated

		// Check if a subclass was selected
		if !ls.subclassSelector.IsVisible() {
			selectedSubclass := ls.subclassSelector.GetSelectedSubclass()
			if selectedSubclass != nil {
				ls.selectedSubclass = selectedSubclass.Name

				// Apply the subclass to the character immediately
				classLevelData := ls.character.GetClassLevelStruct(ls.selectedClass)
				if classLevelData != nil {
					classLevelData.Subclass = ls.selectedSubclass
				}

				// Continue with level-up process
				if ls.preview != nil && ls.preview.RequiresSkills {
					ls.state = LevelUpSelectSkills
				} else {
					ls.state = LevelUpComplete
				}
			} else {
				// User pressed ESC, cancel the level-up
				*ls.character = *ls.backup
				ls.Hide()
			}
		}
		return *ls, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch ls.state {
		case LevelUpSelectClass:
			return ls.handleClassSelection(msg)
		case LevelUpConfirm:
			return ls.handleConfirmation(msg)
		case LevelUpSelectSubclass:
			return ls.handleSubclassSelection(msg)
		case LevelUpSelectSkills:
			return ls.handleSkillSelection(msg)
		case LevelUpComplete:
			if msg.String() == "enter" || msg.String() == "esc" {
				ls.Hide()
			}
		}
	}

	return *ls, nil
}

// handleClassSelection handles class selection keys
func (ls *LevelUpSelector) handleClassSelection(msg tea.KeyMsg) (LevelUpSelector, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if ls.cursor > 0 {
			ls.cursor--
		}
	case "down", "j":
		if ls.cursor < len(ls.availableClasses)-1 {
			ls.cursor++
		}
	case "enter":
		if ls.cursor < len(ls.availableClasses) {
			ls.selectedClass = ls.availableClasses[ls.cursor].Name

			// Generate preview
			preview, err := models.GetLevelUpPreview(ls.character, ls.selectedClass)
			if err != nil {
				ls.message = fmt.Sprintf("Error: %s", err.Error())
				return *ls, nil
			}
			ls.preview = preview
			ls.state = LevelUpConfirm
		}
	case "esc":
		ls.Hide()
	}
	return *ls, nil
}

// handleConfirmation handles confirmation screen keys
func (ls *LevelUpSelector) handleConfirmation(msg tea.KeyMsg) (LevelUpSelector, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		// Proceed with level up
		options := models.LevelUpOptions{
			ClassName:   ls.selectedClass,
			TakeAverage: ls.takeAverage,
		}

		result, err := models.LevelUp(ls.character, options)
		if err != nil {
			ls.message = fmt.Sprintf("Level up failed: %s", err.Error())
			ls.state = LevelUpSelectClass
			return *ls, nil
		}

		ls.preview = result

		// Check if subclass selection is needed
		if result.RequiresSubclass {
			// Determine which class level we're at
			classLevelData := ls.character.GetClassLevelStruct(ls.selectedClass)
			if classLevelData != nil {
				ls.subclassSelector.Show(ls.selectedClass, classLevelData.Level)
				ls.state = LevelUpSelectSubclass
				return *ls, nil
			}
		}

		// Check if skill selections are needed
		if result.RequiresSkills {
			ls.state = LevelUpSelectSkills
			return *ls, nil
		}

		// If no additional selections needed, we're done
		ls.state = LevelUpComplete

	case "r":
		// Toggle between roll and average
		ls.takeAverage = !ls.takeAverage
		// Regenerate preview with new HP option
		preview, _ := models.GetLevelUpPreview(ls.character, ls.selectedClass)
		if preview != nil {
			ls.preview = preview
		}

	case "n", "esc":
		// Cancel and go back
		ls.state = LevelUpSelectClass
		ls.selectedClass = ""
		ls.preview = nil
	}
	return *ls, nil
}

// handleSubclassSelection handles subclass selection
func (ls *LevelUpSelector) handleSubclassSelection(msg tea.KeyMsg) (LevelUpSelector, tea.Cmd) {
	// This state shouldn't be reached directly; the subclass selector is handled above
	// But if we're here, proceed to skills or complete
	switch msg.String() {
	case "enter":
		if ls.preview.RequiresSkills {
			ls.state = LevelUpSelectSkills
		} else {
			ls.state = LevelUpComplete
		}
	case "esc":
		// Rollback
		*ls.character = *ls.backup
		ls.Hide()
	}
	return *ls, nil
}

// handleSkillSelection handles skill selection (placeholder)
func (ls *LevelUpSelector) handleSkillSelection(msg tea.KeyMsg) (LevelUpSelector, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// For now, skip to complete
		ls.state = LevelUpComplete
	case "esc":
		// Rollback
		*ls.character = *ls.backup
		ls.Hide()
	}
	return *ls, nil
}

// Rollback restores the character to the backup state
func (ls *LevelUpSelector) Rollback() {
	if ls.backup != nil {
		*ls.character = *ls.backup
	}
}

// View renders the level-up selector
func (ls *LevelUpSelector) View() string {
	// If subclass selector is visible, show it
	if ls.subclassSelector.IsVisible() {
		return ls.subclassSelector.View()
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var content string

	switch ls.state {
	case LevelUpSelectClass:
		content += titleStyle.Render("LEVEL UP - SELECT CLASS") + "\n\n"
		content += normalStyle.Render(fmt.Sprintf("Current: %s (Level %d)",
			ls.character.GetClassDisplayString(), ls.character.TotalLevel)) + "\n\n"

		if len(ls.availableClasses) == 0 {
			content += dimStyle.Render("No classes available") + "\n"
		} else {
			for i, class := range ls.availableClasses {
				cursor := "  "
				style := normalStyle
				if i == ls.cursor {
					cursor = "❯ "
					style = selectedStyle
				}

				// Check if it's an existing class or new
				isExisting := ls.character.HasClass(class.Name)
				suffix := ""
				if isExisting {
					currentLevel := ls.character.GetClassLevel(class.Name)
					suffix = fmt.Sprintf(" (Level %d → %d)", currentLevel, currentLevel+1)
				} else {
					suffix = " [NEW CLASS]"
				}

				content += style.Render(fmt.Sprintf("%s%s%s", cursor, class.Name, suffix)) + "\n"
			}
		}

		content += "\n" + dimStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Cancel")

	case LevelUpConfirm:
		content += titleStyle.Render("LEVEL UP - CONFIRM") + "\n\n"

		if ls.preview != nil {
			content += selectedStyle.Render(fmt.Sprintf("Class: %s", ls.selectedClass)) + "\n"
			content += normalStyle.Render(fmt.Sprintf("Level: %d → %d",
				ls.preview.NewClassLevel-1, ls.preview.NewClassLevel)) + "\n"
			content += normalStyle.Render(fmt.Sprintf("Total Level: %d → %d",
				ls.character.TotalLevel, ls.preview.NewTotalLevel)) + "\n\n"

			content += titleStyle.Render("GAINS:") + "\n"

			// HP
			hpMethod := "average"
			if !ls.takeAverage {
				hpMethod = "rolled"
			}
			content += normalStyle.Render(fmt.Sprintf("  HP: +%d (%s)", ls.preview.HPGained, hpMethod)) + "\n"

			// Features
			if len(ls.preview.FeaturesGained) > 0 {
				content += normalStyle.Render("  Features:") + "\n"
				for _, feature := range ls.preview.FeaturesGained {
					content += dimStyle.Render(fmt.Sprintf("    • %s", feature)) + "\n"
				}
			}

			// Proficiencies
			if len(ls.preview.ProficienciesGained) > 0 {
				content += normalStyle.Render("  Proficiencies:") + "\n"
				for _, prof := range ls.preview.ProficienciesGained {
					content += dimStyle.Render(fmt.Sprintf("    • %s", prof)) + "\n"
				}
			}

			// Additional selections needed
			if ls.preview.RequiresSkills {
				content += "\n" + dimStyle.Render("  ⚠ Skill selection required after confirmation") + "\n"
			}
			if ls.preview.RequiresSubclass {
				content += dimStyle.Render("  ⚠ Subclass selection required after confirmation") + "\n"
			}
			if ls.preview.RequiresSpells {
				content += dimStyle.Render("  ⚠ Spell selection required after confirmation") + "\n"
			}
		}

		content += "\n" + dimStyle.Render("Y/Enter: Confirm • R: Toggle Roll/Average • N/Esc: Cancel")

	case LevelUpSelectSubclass:
		content += titleStyle.Render("LEVEL UP - SELECT SUBCLASS") + "\n\n"
		content += normalStyle.Render("Subclass selection in progress...") + "\n\n"
		content += dimStyle.Render("(Use subclass selector)")

	case LevelUpSelectSkills:
		content += titleStyle.Render("LEVEL UP - SELECT SKILLS") + "\n\n"
		content += normalStyle.Render("Skill selection UI (to be implemented)") + "\n\n"
		content += dimStyle.Render("Enter: Continue • Esc: Cancel")

	case LevelUpComplete:
		content += titleStyle.Render("LEVEL UP COMPLETE!") + "\n\n"

		if ls.preview != nil {
			content += selectedStyle.Render(fmt.Sprintf("You are now a Level %d %s!",
				ls.preview.NewTotalLevel, ls.character.GetClassDisplayString())) + "\n\n"

			// Show subclass if one was selected
			if ls.selectedSubclass != "" {
				content += normalStyle.Render(fmt.Sprintf("Subclass: %s", ls.selectedSubclass)) + "\n\n"
			}

			content += normalStyle.Render("Summary:") + "\n"
			content += dimStyle.Render(fmt.Sprintf("  HP: %d → %d (+%d)",
				ls.character.MaxHP-ls.preview.HPGained, ls.character.MaxHP, ls.preview.HPGained)) + "\n"

			if len(ls.preview.FeaturesGained) > 0 {
				content += dimStyle.Render(fmt.Sprintf("  New Features: %s",
					strings.Join(ls.preview.FeaturesGained, ", "))) + "\n"
			}
		}

		content += "\n" + dimStyle.Render("Enter: Close")
	}

	if ls.message != "" {
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(ls.message)
	}

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(70)

	return lipgloss.Place(
		80,
		30,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content),
	)
}
