// internal/ui/components/weaponmasteryselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// WeaponMasterySelector allows selecting weapons to master
type WeaponMasterySelector struct {
	character      *models.Character
	visible        bool
	selectedIndex  int
	maxMasteries   int // Maximum number of weapons that can be mastered
	availableWeapons []weaponMasteryOption
	selectedWeapons map[string]bool // Track which weapons are selected
}

type weaponMasteryOption struct {
	name    string
	mastery string
}

// NewWeaponMasterySelector creates a new weapon mastery selector
func NewWeaponMasterySelector(char *models.Character) *WeaponMasterySelector {
	return &WeaponMasterySelector{
		character:       char,
		visible:         false,
		selectedIndex:   0,
		selectedWeapons: make(map[string]bool),
	}
}

// Show displays the selector with the specified maximum number of masteries
func (s *WeaponMasterySelector) Show(maxMasteries int) {
	s.visible = true
	s.selectedIndex = 0
	s.maxMasteries = maxMasteries
	s.selectedWeapons = make(map[string]bool)

	// Initialize with currently mastered weapons
	for _, weapon := range s.character.MasteredWeapons {
		s.selectedWeapons[weapon] = true
	}

	s.buildAvailableWeapons()
}

// Hide closes the selector
func (s *WeaponMasterySelector) Hide() {
	s.visible = false
}

// IsVisible returns whether the selector is visible
func (s *WeaponMasterySelector) IsVisible() bool {
	return s.visible
}

// buildAvailableWeapons builds the list of weapons the character can master
func (s *WeaponMasterySelector) buildAvailableWeapons() {
	s.availableWeapons = []weaponMasteryOption{}

	// Get all weapon items
	allItems := models.GetAllItemDefinitions()

	// Build a map of weapons with their mastery properties
	for _, itemDef := range allItems {
		// Only consider weapons
		if itemDef.Category != "weapon" {
			continue
		}

		// Check if character has proficiency with this weapon
		if !s.hasProficiency(itemDef.Subcategory) {
			continue
		}

		s.availableWeapons = append(s.availableWeapons, weaponMasteryOption{
			name:    itemDef.Name,
			mastery: itemDef.Mastery,
		})
	}
}

// hasProficiency checks if the character has proficiency with a weapon type
func (s *WeaponMasterySelector) hasProficiency(subcategory string) bool {
	subcategory = strings.ToLower(subcategory)

	for _, prof := range s.character.WeaponProficiencies {
		profLower := strings.ToLower(prof)

		// Check if subcategory contains the proficiency
		// e.g., "simple melee" contains "simple", "martial ranged" contains "martial"
		if strings.Contains(subcategory, profLower) {
			return true
		}
	}

	return false
}

// ToggleSelection toggles the selection of the current weapon
func (s *WeaponMasterySelector) ToggleSelection() bool {
	if s.selectedIndex < 0 || s.selectedIndex >= len(s.availableWeapons) {
		return false
	}

	weaponName := s.availableWeapons[s.selectedIndex].name

	if s.selectedWeapons[weaponName] {
		// Deselect
		delete(s.selectedWeapons, weaponName)
		return true
	} else {
		// Check if we can select more
		if len(s.selectedWeapons) < s.maxMasteries {
			s.selectedWeapons[weaponName] = true
			return true
		}
		return false
	}
}

// Next moves to the next weapon
func (s *WeaponMasterySelector) Next() {
	if s.selectedIndex < len(s.availableWeapons)-1 {
		s.selectedIndex++
	}
}

// Prev moves to the previous weapon
func (s *WeaponMasterySelector) Prev() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// CanConfirm returns true if the selection is valid
func (s *WeaponMasterySelector) CanConfirm() bool {
	return len(s.selectedWeapons) > 0 && len(s.selectedWeapons) <= s.maxMasteries
}

// GetSelectedWeapons returns the list of selected weapon names
func (s *WeaponMasterySelector) GetSelectedWeapons() []string {
	weapons := []string{}
	for weapon := range s.selectedWeapons {
		weapons = append(weapons, weapon)
	}
	return weapons
}

// Update handles keyboard input for the weapon mastery selector
func (s *WeaponMasterySelector) Update(msg tea.Msg) (WeaponMasterySelector, tea.Cmd) {
	if !s.visible {
		return *s, nil
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
		case "enter", "esc":
			return *s, nil
		}
	}

	return *s, nil
}

// View renders the selector with description on the right
func (s *WeaponMasterySelector) View() string {
	if !s.visible {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	checkboxStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	masteryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	// Left side - weapon list
	var leftLines []string
	title := fmt.Sprintf("Select Weapon Masteries (%d/%d selected)", len(s.selectedWeapons), s.maxMasteries)
	leftLines = append(leftLines, titleStyle.Render(title))
	leftLines = append(leftLines, "")

	if len(s.availableWeapons) == 0 {
		leftLines = append(leftLines, normalStyle.Render("  No weapons available (no weapon proficiencies)"))
	} else {
		for i, weapon := range s.availableWeapons {
			checkbox := "[ ]"
			if s.selectedWeapons[weapon.name] {
				checkbox = checkboxStyle.Render("[✓]")
			}

			masteryText := ""
			if weapon.mastery != "" {
				masteryText = " " + masteryStyle.Render(fmt.Sprintf("(%s)", weapon.mastery))
			}

			line := fmt.Sprintf("%s %s%s", checkbox, weapon.name, masteryText)

			if i == s.selectedIndex {
				leftLines = append(leftLines, selectedStyle.Render("▶ "+line))
			} else {
				leftLines = append(leftLines, normalStyle.Render("  "+line))
			}
		}
	}

	leftLines = append(leftLines, "")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if s.CanConfirm() {
		leftLines = append(leftLines, helpStyle.Render("↑/↓: Navigate • Space: Toggle"))
		leftLines = append(leftLines, helpStyle.Render("Enter: Confirm • Esc: Cancel"))
	} else {
		leftLines = append(leftLines, helpStyle.Render("↑/↓: Navigate • Space: Toggle"))
		leftLines = append(leftLines, helpStyle.Render("Esc: Cancel"))
	}

	leftContent := strings.Join(leftLines, "\n")

	// Right side - mastery description
	var rightLines []string
	rightLines = append(rightLines, titleStyle.Render("MASTERY DESCRIPTION"))
	rightLines = append(rightLines, "")

	if len(s.availableWeapons) > 0 && s.selectedIndex >= 0 && s.selectedIndex < len(s.availableWeapons) {
		currentWeapon := s.availableWeapons[s.selectedIndex]
		if currentWeapon.mastery != "" {
			rightLines = append(rightLines, masteryStyle.Render(currentWeapon.mastery))
			rightLines = append(rightLines, "")

			// Get mastery description
			description := models.GetMasteryDescription(currentWeapon.mastery)
			if description != "" {
				// Wrap description to fit in the box (45 chars width)
				wrapped := wrapMasteryText(description, 43)
				rightLines = append(rightLines, normalStyle.Render(wrapped))
			}
		} else {
			rightLines = append(rightLines, normalStyle.Render("No mastery property"))
		}
	} else {
		rightLines = append(rightLines, normalStyle.Render("Select a weapon to view its mastery"))
	}

	rightContent := strings.Join(rightLines, "\n")

	// Create two boxes side by side
	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(50).
		Height(25)

	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		Width(50).
		Height(25)

	// Join boxes horizontally
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftBox.Render(leftContent),
		rightBox.Render(rightContent),
	)

	return combined
}

// wrapMasteryText wraps text to the specified width
func wrapMasteryText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine)
			result.WriteString("\n")
			currentLine = word
		}
	}

	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}
