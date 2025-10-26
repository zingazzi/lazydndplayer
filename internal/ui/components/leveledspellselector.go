// internal/ui/components/leveledspellselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// LeveledSpellSelector allows selecting multiple leveled spells (for Wizard spellbook)
type LeveledSpellSelector struct {
	visible          bool
	availableSpells  []models.Spell
	selectedSpells   []string // Names of selected spells
	cursor           int
	maxSpells        int
	className        string
	spellLevel       int // Which level of spells to show
	character        *models.Character
	viewport         viewport.Model
}

func NewLeveledSpellSelector(char *models.Character) *LeveledSpellSelector {
	vp := viewport.New(70, 20)
	return &LeveledSpellSelector{
		visible:         false,
		availableSpells: []models.Spell{},
		selectedSpells:  []string{},
		cursor:          0,
		character:       char,
		viewport:        vp,
	}
}

func (ls *LeveledSpellSelector) Show(className string, spellLevel int, maxSpells int) {
	ls.className = className
	ls.spellLevel = spellLevel
	ls.maxSpells = maxSpells
	ls.visible = true
	ls.cursor = 0
	ls.selectedSpells = []string{}

	// Load spells for this class and level
	ls.loadSpells()
}

func (ls *LeveledSpellSelector) Hide() {
	ls.visible = false
}

func (ls *LeveledSpellSelector) IsVisible() bool {
	return ls.visible
}

func (ls *LeveledSpellSelector) loadSpells() {
	// Load all spells from data file
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		ls.availableSpells = []models.Spell{}
		return
	}

	ls.availableSpells = []models.Spell{}
	classLower := strings.ToLower(ls.className)

	for _, spell := range allSpells {
		// Only spells of the specified level
		if spell.Level != ls.spellLevel {
			continue
		}

		// Check if this class can learn this spell
		for _, spellClass := range spell.Classes {
			if strings.ToLower(spellClass) == classLower {
				ls.availableSpells = append(ls.availableSpells, spell)
				break
			}
		}
	}
}

func (ls *LeveledSpellSelector) Update(msg tea.Msg) (LeveledSpellSelector, tea.Cmd) {
	if !ls.visible {
		return *ls, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			ls.Prev()
		case "down", "j":
			ls.Next()
		case " ":
			ls.ToggleSelection()
		case "enter":
			if len(ls.selectedSpells) == ls.maxSpells {
				// Confirmed
				return *ls, nil
			}
		case "esc":
			ls.Hide()
		}
	}

	return *ls, nil
}

func (ls *LeveledSpellSelector) Next() {
	if ls.cursor < len(ls.availableSpells)-1 {
		ls.cursor++
		ls.viewport.LineDown(1)
	}
}

func (ls *LeveledSpellSelector) Prev() {
	if ls.cursor > 0 {
		ls.cursor--
		ls.viewport.LineUp(1)
	}
}

func (ls *LeveledSpellSelector) ToggleSelection() {
	if len(ls.availableSpells) == 0 {
		return
	}

	selected := ls.availableSpells[ls.cursor]

	// Check if already selected
	for i, name := range ls.selectedSpells {
		if name == selected.Name {
			// Deselect
			ls.selectedSpells = append(ls.selectedSpells[:i], ls.selectedSpells[i+1:]...)
			return
		}
	}

	// Select if not at max
	if len(ls.selectedSpells) < ls.maxSpells {
		ls.selectedSpells = append(ls.selectedSpells, selected.Name)
	}
}

func (ls *LeveledSpellSelector) IsSelected(spellName string) bool {
	for _, name := range ls.selectedSpells {
		if name == spellName {
			return true
		}
	}
	return false
}

func (ls *LeveledSpellSelector) GetSelectedSpells() []string {
	return ls.selectedSpells
}

func (ls *LeveledSpellSelector) GetRemainingCount() int {
	return ls.maxSpells - len(ls.selectedSpells)
}

func (ls *LeveledSpellSelector) View(width, height int) string {
	if !ls.visible {
		return ""
	}

	// Update viewport size
	ls.viewport.Width = width - 4
	ls.viewport.Height = height - 8

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content strings.Builder

	title := fmt.Sprintf("SELECT LEVEL %d SPELLS FOR %s", ls.spellLevel, strings.ToUpper(ls.className))
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	remaining := ls.GetRemainingCount()
	progressText := fmt.Sprintf("Selected: %d/%d", len(ls.selectedSpells), ls.maxSpells)
	if remaining > 0 {
		progressText += fmt.Sprintf(" (Select %d more)", remaining)
	} else {
		progressText += " (Press Enter to confirm)"
	}
	content.WriteString(dimStyle.Render(progressText))
	content.WriteString("\n\n")

	if len(ls.availableSpells) == 0 {
		content.WriteString(dimStyle.Render("No spells available"))
	} else {
		for i, spell := range ls.availableSpells {
			cursor := "  "
			checkbox := "☐ "
			style := normalStyle

			if i == ls.cursor {
				cursor = "→ "
				style = selectedStyle
			}

			if ls.IsSelected(spell.Name) {
				checkbox = "☑ "
			}

			// Add spell markers
			markers := ""
			if spell.Concentration {
				markers += " (C)"
			}
			if spell.Ritual {
				markers += " (R)"
			}

			spellLine := fmt.Sprintf("%s%s%s%s", cursor, checkbox, spell.Name, markers)
			content.WriteString(style.Render(spellLine))
			content.WriteString("\n")

			// Show description for selected spell
			if i == ls.cursor {
				wrappedDesc := wrapSpellSelectorText(spell.Description, ls.viewport.Width-4)
				for _, line := range wrappedDesc {
					content.WriteString(dimStyle.Render("    "+line))
					content.WriteString("\n")
				}
				content.WriteString("\n")
			}
		}
	}

	content.WriteString("\n")
	content.WriteString(dimStyle.Render("↑/↓: Navigate | Space: Toggle | Enter: Confirm | Esc: Cancel"))

	// Set viewport content
	ls.viewport.SetContent(content.String())

	// Render viewport with border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		borderStyle.Render(ls.viewport.View()),
	)
}

// wrapSpellSelectorText wraps text to specified width
func wrapSpellSelectorText(text string, width int) []string {
	if width <= 0 {
		width = 60
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
