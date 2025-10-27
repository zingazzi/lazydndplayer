package components

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type SchoolSpellSelector struct {
	visible         bool
	character       *models.Character
	availableSpells []models.Spell
	selectedSpells  []models.Spell
	school          string
	maxLevel        int
	maxCount        int
	cursorIndex     int
	viewport        viewport.Model
}

func NewSchoolSpellSelector(char *models.Character) *SchoolSpellSelector {
	vp := viewport.New(80, 20)
	return &SchoolSpellSelector{
		character: char,
		viewport:  vp,
	}
}

func (s *SchoolSpellSelector) Show(school string, maxLevel int, maxCount int) {
	debug.Log("=== SCHOOL SPELL SELECTOR SHOW ===")
	debug.Log("  School: %s, Max Level: %d, Max Count: %d", school, maxLevel, maxCount)

	s.school = school
	s.maxLevel = maxLevel
	s.maxCount = maxCount
	s.cursorIndex = 0
	s.selectedSpells = []models.Spell{}
	s.visible = true

	// Load all spells from JSON
	allSpells, err := s.loadSpells()
	if err != nil {
		debug.Log("  ERROR loading spells: %v", err)
		return
	}

	// Filter spells
	s.availableSpells = []models.Spell{}
	for _, spell := range allSpells {
		// Check if spell matches criteria
		if !s.matchesFilter(spell) {
			continue
		}

		// Check if already in spellbook
		if s.isAlreadyKnown(spell) {
			debug.Log("  Skipping already known spell: %s", spell.Name)
			continue
		}

		s.availableSpells = append(s.availableSpells, spell)
	}

	debug.Log("  Found %d available spells", len(s.availableSpells))
}

func (s *SchoolSpellSelector) matchesFilter(spell models.Spell) bool {
	// Check school
	if !strings.EqualFold(string(spell.School), s.school) {
		return false
	}

	// Check level (1 to maxLevel)
	if spell.Level < 1 || spell.Level > s.maxLevel {
		return false
	}

	// Check if it's a Wizard spell
	hasWizard := false
	for _, class := range spell.Classes {
		if strings.EqualFold(class, "wizard") {
			hasWizard = true
			break
		}
	}

	return hasWizard
}

func (s *SchoolSpellSelector) isAlreadyKnown(spell models.Spell) bool {
	for _, known := range s.character.SpellBook.Spells {
		if strings.EqualFold(known.Name, spell.Name) {
			return true
		}
	}
	return false
}

func (s *SchoolSpellSelector) loadSpells() ([]models.Spell, error) {
	file, err := os.ReadFile("data/spells.json")
	if err != nil {
		return nil, err
	}

	var spells []models.Spell
	if err := json.Unmarshal(file, &spells); err != nil {
		return nil, err
	}

	return spells, nil
}

func (s *SchoolSpellSelector) IsVisible() bool {
	return s.visible
}

func (s *SchoolSpellSelector) Hide() {
	s.visible = false
}

func (s *SchoolSpellSelector) GetSelectedSpells() []models.Spell {
	return s.selectedSpells
}

func (s *SchoolSpellSelector) GetRemainingCount() int {
	return s.maxCount - len(s.selectedSpells)
}

func (s *SchoolSpellSelector) GetSchool() string {
	return s.school
}

func (s *SchoolSpellSelector) ToggleSelection() {
	if s.cursorIndex < 0 || s.cursorIndex >= len(s.availableSpells) {
		return
	}

	selectedSpell := s.availableSpells[s.cursorIndex]

	// Check if already selected
	alreadySelected := false
	selectedIndex := -1
	for i, spell := range s.selectedSpells {
		if spell.Name == selectedSpell.Name {
			alreadySelected = true
			selectedIndex = i
			break
		}
	}

	if alreadySelected {
		// Remove from selection
		s.selectedSpells = append(s.selectedSpells[:selectedIndex], s.selectedSpells[selectedIndex+1:]...)
		debug.Log("  Deselected spell: %s", selectedSpell.Name)
	} else {
		// Add to selection if not at max
		if len(s.selectedSpells) < s.maxCount {
			s.selectedSpells = append(s.selectedSpells, selectedSpell)
			debug.Log("  Selected spell: %s", selectedSpell.Name)
		}
	}
}

func (s *SchoolSpellSelector) isSelected(spell models.Spell) bool {
	for _, selected := range s.selectedSpells {
		if selected.Name == spell.Name {
			return true
		}
	}
	return false
}

func (s *SchoolSpellSelector) Update(msg tea.Msg) tea.Cmd {
	if !s.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursorIndex > 0 {
				s.cursorIndex--
			}
		case "down", "j":
			if s.cursorIndex < len(s.availableSpells)-1 {
				s.cursorIndex++
			}
		case " ":
			s.ToggleSelection()
		case "enter":
			if len(s.selectedSpells) > 0 {
				debug.Log("  Confirming selection: %d spells", len(s.selectedSpells))
				s.Hide()
			}
		case "esc":
			s.selectedSpells = []models.Spell{}
			s.Hide()
		}
	}

	return nil
}

func (s *SchoolSpellSelector) View(width, height int) string {
	if !s.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("15")).
		Bold(true)

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	schoolTitle := strings.Title(s.school)
	title := fmt.Sprintf("%s Spells (Level 1-%d)", schoolTitle, s.maxLevel)

	var leftLines []string
	leftLines = append(leftLines, titleStyle.Render(title))
	leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Select %d spells • Space: Select • Enter: Confirm • Esc: Cancel", s.maxCount)))
	leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Selected: %d/%d", len(s.selectedSpells), s.maxCount)))
	leftLines = append(leftLines, "")

	// Left panel: spell list
	for i, spell := range s.availableSpells {
		checkbox := "[ ]"
		if s.isSelected(spell) {
			checkbox = "[●]"
		}

		line := fmt.Sprintf("%s %s (Lv%d)", checkbox, spell.Name, spell.Level)

		if i == s.cursorIndex {
			line = selectedStyle.Render(line)
		}

		leftLines = append(leftLines, line)
	}

	leftPanel := strings.Join(leftLines, "\n")

	// Right panel: spell description
	var rightLines []string
	if s.cursorIndex >= 0 && s.cursorIndex < len(s.availableSpells) {
		spell := s.availableSpells[s.cursorIndex]

		rightLines = append(rightLines, lipgloss.NewStyle().Bold(true).Render(spell.Name))
		rightLines = append(rightLines, dimStyle.Render(fmt.Sprintf("Level %d %s", spell.Level, strings.Title(string(spell.School)))))
		rightLines = append(rightLines, "")
		rightLines = append(rightLines, fmt.Sprintf("Range: %s", spell.Range))
		rightLines = append(rightLines, fmt.Sprintf("Duration: %s", spell.Duration))

		// Components
		componentStr := "Components: " + spell.GetComponentsString()
		rightLines = append(rightLines, componentStr)

		rightLines = append(rightLines, "")

		// Description (wrap text)
		desc := spell.Description
		if len(desc) > 300 {
			desc = desc[:300] + "..."
		}
		rightLines = append(rightLines, desc)
	}

	rightPanel := strings.Join(rightLines, "\n")

	// Two-column layout
	leftWidth := width / 2
	rightWidth := width - leftWidth

	leftBox := lipgloss.NewStyle().
		Width(leftWidth).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		Render(leftPanel)

	rightBox := lipgloss.NewStyle().
		Width(rightWidth).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		Render(rightPanel)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}
