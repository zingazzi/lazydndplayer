// internal/ui/components/classselector.go
package components

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ClassSelector is a component for selecting a class
type ClassSelector struct {
	character     *models.Character
	classes       []models.Class
	selectedIndex int
	visible       bool
	isMulticlass  bool // Whether this is for multiclassing or first class
}

// NewClassSelector creates a new class selector
func NewClassSelector(char *models.Character) *ClassSelector {
	return &ClassSelector{
		character:     char,
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the class selector
func (c *ClassSelector) Show() {
	c.visible = true
	c.selectedIndex = 0

	// Determine if this is for multiclassing
	// Character is multiclassing if they have a level AND already have a class
	c.isMulticlass = c.character.TotalLevel > 0 && c.character.Class != ""
	debug.Log("ClassSelector.Show() - TotalLevel=%d, Class='%s', isMulticlass=%v", c.character.TotalLevel, c.character.Class, c.isMulticlass)

	// Load appropriate classes
	if c.isMulticlass {
		// Show only classes that meet prerequisites
		c.classes = models.GetAvailableClasses(c.character)
		debug.Log("ClassSelector.Show() - Multiclass mode, loaded %d classes", len(c.classes))
	} else {
		// Show all classes for first class selection
		c.classes = models.GetAllClasses()
		debug.Log("ClassSelector.Show() - First class mode, loaded %d classes", len(c.classes))
	}

	// Debug: Check if classes loaded successfully
	if len(c.classes) == 0 {
		debug.Log("WARNING: ClassSelector.Show() - No classes loaded! isMulticlass=%v", c.isMulticlass)
	} else {
		debug.Log("ClassSelector - First class: %s, Last class: %s", c.classes[0].Name, c.classes[len(c.classes)-1].Name)
	}
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
		debug.Log("ClassSelector.Next() - selectedIndex now: %d (total: %d)", c.selectedIndex, len(c.classes))
	} else {
		debug.Log("ClassSelector.Next() - Already at end (index: %d, total: %d)", c.selectedIndex, len(c.classes))
	}
}

// Prev moves to the previous class
func (c *ClassSelector) Prev() {
	if c.selectedIndex > 0 {
		c.selectedIndex--
		debug.Log("ClassSelector.Prev() - selectedIndex now: %d (total: %d)", c.selectedIndex, len(c.classes))
	} else {
		debug.Log("ClassSelector.Prev() - Already at start")
	}
}

// GetSelectedClass returns the currently selected class name
func (c *ClassSelector) GetSelectedClass() string {
	debug.Log("GetSelectedClass() - selectedIndex=%d, classes.len=%d", c.selectedIndex, len(c.classes))
	if c.selectedIndex >= 0 && c.selectedIndex < len(c.classes) {
		className := c.classes[c.selectedIndex].Name
		debug.Log("GetSelectedClass() - Returning: %s", className)
		return className
	}
	debug.Log("GetSelectedClass() - Returning empty (index out of bounds)")
	return ""
}

// Update handles keyboard input for the class selector
func (c *ClassSelector) Update(msg tea.Msg) (ClassSelector, tea.Cmd) {
	if !c.visible {
		return *c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			c.Prev()
		case "down", "j":
			c.Next()
		case "enter", "esc":
			return *c, nil
		}
	}

	return *c, nil
}

// getFeaturesToGain returns the names of features that will be gained at the specified level
func (c *ClassSelector) getFeaturesToGain(className string, level int) []string {
	featureNames := []string{}
	
	// Read class JSON file
	classFilePath := fmt.Sprintf("data/classes/%s.json", strings.ToLower(className))
	data, err := os.ReadFile(classFilePath)
	if err != nil {
		debug.Log("getFeaturesToGain: Error reading class file: %v", err)
		return featureNames
	}
	
	// Parse level_progression
	var classData struct {
		LevelProgression []struct {
			Level    int `json:"level"`
			Features []struct {
				Name string `json:"name"`
			} `json:"features"`
		} `json:"level_progression"`
	}
	
	if err := json.Unmarshal(data, &classData); err != nil {
		debug.Log("getFeaturesToGain: Error unmarshaling class data: %v", err)
		return featureNames
	}
	
	// Find features for the specified level
	for _, levelData := range classData.LevelProgression {
		if levelData.Level == level {
			for _, feature := range levelData.Features {
				// Skip subclass choice features
				if feature.Name == "Divine Domain" || feature.Name == "Sorcerous Origin" ||
				   feature.Name == "Otherworldly Patron" || feature.Name == "Monastic Tradition" ||
				   feature.Name == "Primal Path" || strings.Contains(strings.ToLower(feature.Name), "subclass") {
					continue
				}
				featureNames = append(featureNames, feature.Name)
			}
			break
		}
	}
	
	return featureNames
}

// View renders the class selector
func (c *ClassSelector) View(width, height int) string {
	if !c.visible || len(c.classes) == 0 {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	// Build content
	var content strings.Builder

	if c.isMulticlass {
		content.WriteString(titleStyle.Render("SELECT CLASS TO MULTICLASS INTO") + "\n")
		content.WriteString(dimStyle.Render(fmt.Sprintf("Current: %s (Level %d)",
			c.character.GetClassDisplayString(), c.character.TotalLevel)) + "\n\n")
	} else {
		content.WriteString(titleStyle.Render("SELECT YOUR CLASS") + "\n\n")
	}

	// Two-column layout: class list on left, details on right
	var leftContent strings.Builder
	var rightContent strings.Builder

	// Build class list (left side)
	for i, class := range c.classes {
		cursor := "  "
		style := normalStyle
		if i == c.selectedIndex {
			cursor = "❯ "
			style = selectedStyle
		}

		// Check if character already has this class
		hasClass := c.character.HasClass(class.Name)
		suffix := ""
		if hasClass {
			currentLevel := c.character.GetClassLevel(class.Name)
			suffix = fmt.Sprintf(" (Level %d)", currentLevel)
		}

		leftContent.WriteString(style.Render(fmt.Sprintf("%s%s%s", cursor, class.Name, suffix)) + "\n")
	}

	// Build details (right side)
	if c.selectedIndex >= 0 && c.selectedIndex < len(c.classes) {
		currentClass := c.classes[c.selectedIndex]

		// Class name and primary ability
		rightContent.WriteString(selectedStyle.Render(currentClass.Name) + "\n")
		rightContent.WriteString(dimStyle.Render(fmt.Sprintf("Hit Die: d%d • Primary: %s",
			currentClass.HitDie, currentClass.PrimaryAbility)) + "\n\n")

		// Description
		rightContent.WriteString(normalStyle.Render(currentClass.Description) + "\n\n")

		// Show prerequisites if multiclassing
		if c.isMulticlass {
			canMulticlass, reason := models.CanMulticlassInto(c.character, currentClass.Name)
			if canMulticlass {
				rightContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓ Prerequisites met") + "\n")
			} else {
				rightContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ " + reason) + "\n")
			}
			rightContent.WriteString("\n")
		}

		// Show proficiencies that will be granted
		rightContent.WriteString(titleStyle.Render("PROFICIENCIES GRANTED:") + "\n")

		if c.isMulticlass && !c.character.HasClass(currentClass.Name) {
			// Show limited multiclass proficiencies
			multiclassProf := models.GetMulticlassProficiencies(currentClass.Name)
			if len(multiclassProf) > 0 {
				rightContent.WriteString(warningStyle.Render("Multiclass (Limited):") + "\n")
				for _, prof := range multiclassProf {
					rightContent.WriteString(dimStyle.Render("  • " + prof) + "\n")
				}
			} else {
				rightContent.WriteString(dimStyle.Render("  • No additional proficiencies") + "\n")
			}
		} else if c.character.HasClass(currentClass.Name) {
			// Already has this class
			rightContent.WriteString(dimStyle.Render("  • Continuing existing class") + "\n")
		} else {
			// First class - full proficiencies
			rightContent.WriteString(normalStyle.Render("Full Proficiencies:") + "\n")

			// Armor
			if len(currentClass.ArmorProficiencies) > 0 {
				rightContent.WriteString(dimStyle.Render("  Armor: " + strings.Join(currentClass.ArmorProficiencies, ", ")) + "\n")
			}

			// Weapons
			if len(currentClass.WeaponProficiencies) > 0 {
				rightContent.WriteString(dimStyle.Render("  Weapons: " + strings.Join(currentClass.WeaponProficiencies, ", ")) + "\n")
			}

			// Saving throws
			if len(currentClass.SavingThrows) > 0 {
				rightContent.WriteString(dimStyle.Render("  Saves: " + strings.Join(currentClass.SavingThrows, ", ")) + "\n")
			}

			// Skills
			if currentClass.SkillChoices != nil {
				rightContent.WriteString(dimStyle.Render(fmt.Sprintf("  Skills: Choose %d", currentClass.SkillChoices.Choose)) + "\n")
			}
		}

		// Show features that will be gained
		nextLevel := 1
		if c.character.HasClass(currentClass.Name) {
			nextLevel = c.character.GetClassLevel(currentClass.Name) + 1
		}
		
		features := c.getFeaturesToGain(currentClass.Name, nextLevel)
		if len(features) > 0 {
			rightContent.WriteString("\n" + titleStyle.Render(fmt.Sprintf("FEATURES AT LEVEL %d:", nextLevel)) + "\n")
			for _, featureName := range features {
				rightContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("  ✦ " + featureName) + "\n")
			}
		}

		// Show spellcasting info
		if currentClass.Spellcasting != nil {
			rightContent.WriteString("\n" + titleStyle.Render("SPELLCASTING:") + "\n")
			rightContent.WriteString(dimStyle.Render(fmt.Sprintf("  Ability: %s", currentClass.Spellcasting.Ability)) + "\n")
			if currentClass.Spellcasting.CantripsKnown > 0 {
				rightContent.WriteString(dimStyle.Render(fmt.Sprintf("  Cantrips: %d", currentClass.Spellcasting.CantripsKnown)) + "\n")
			}
		}
	}

	// Join left and right in two columns
	leftBox := lipgloss.NewStyle().
		Width(25).
		Height(25).
		Padding(1).
		Render(leftContent.String())

	rightBox := lipgloss.NewStyle().
		Width(65).
		Height(25).
		Padding(1).
		Render(rightContent.String())

	twoColumns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	content.WriteString(twoColumns)

	content.WriteString("\n" + dimStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Cancel"))

	// Popup style
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(95)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
