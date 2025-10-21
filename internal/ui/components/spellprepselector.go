// internal/ui/components/spellprepselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type SpellPrepSelector struct {
	visible       bool
	allSpells     []models.Spell // All available spells for the class
	selectedLevel int            // 0 = all, 1-9 = filter by level
	cursor        int
	character     *models.Character
}

func NewSpellPrepSelector(char *models.Character) *SpellPrepSelector {
	return &SpellPrepSelector{
		visible:       false,
		allSpells:     []models.Spell{},
		selectedLevel: 0, // Show all levels by default
		cursor:        0,
		character:     char,
	}
}

func (sps *SpellPrepSelector) Show() {
	sps.visible = true
	sps.cursor = 0
	sps.selectedLevel = 0
	sps.loadSpells()
}

func (sps *SpellPrepSelector) Hide() {
	sps.visible = false
}

func (sps *SpellPrepSelector) IsVisible() bool {
	return sps.visible
}

func (sps *SpellPrepSelector) loadSpells() {
	if sps.character.Class == "" {
		sps.allSpells = []models.Spell{}
		return
	}

	// Load all spells from data file
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		sps.allSpells = []models.Spell{}
		return
	}

	// Filter by class and level
	classLower := strings.ToLower(sps.character.Class)
	maxSpellLevel := sps.getMaxSpellLevel()

	sps.allSpells = []models.Spell{}
	for _, spell := range allSpells {
		// Skip cantrips
		if spell.Level == 0 {
			continue
		}

		// Check if this spell is for the character's class
		isForClass := false
		for _, spellClass := range spell.Classes {
			if strings.ToLower(spellClass) == classLower {
				isForClass = true
				break
			}
		}

		if !isForClass {
			continue
		}

		// Check if character can cast this spell level
		if spell.Level > maxSpellLevel {
			continue
		}

		// Apply level filter if set
		if sps.selectedLevel > 0 && spell.Level != sps.selectedLevel {
			continue
		}

		sps.allSpells = append(sps.allSpells, spell)
	}
}

func (sps *SpellPrepSelector) getMaxSpellLevel() int {
	level := sps.character.Level
	if level >= 17 {
		return 9
	} else if level >= 15 {
		return 8
	} else if level >= 13 {
		return 7
	} else if level >= 11 {
		return 6
	} else if level >= 9 {
		return 5
	} else if level >= 7 {
		return 4
	} else if level >= 5 {
		return 3
	} else if level >= 3 {
		return 2
	} else if level >= 1 {
		return 1
	}
	return 0
}

func (sps *SpellPrepSelector) Update(msg tea.Msg) (SpellPrepSelector, tea.Cmd) {
	if !sps.visible {
		return *sps, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			sps.Prev()
		case "down", "j":
			sps.Next()
		case " ":
			sps.TogglePrepare()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Filter by spell level
			level := int(msg.String()[0] - '0')
			if level <= sps.getMaxSpellLevel() {
				if sps.selectedLevel == level {
					sps.selectedLevel = 0 // Toggle off filter
				} else {
					sps.selectedLevel = level
				}
				sps.loadSpells()
				sps.cursor = 0
			}
		case "0":
			// Show all levels
			sps.selectedLevel = 0
			sps.loadSpells()
			sps.cursor = 0
		case "enter", "esc":
			return *sps, nil
		}
	}

	return *sps, nil
}

func (sps *SpellPrepSelector) Next() {
	if sps.cursor < len(sps.allSpells)-1 {
		sps.cursor++
	}
}

func (sps *SpellPrepSelector) Prev() {
	if sps.cursor > 0 {
		sps.cursor--
	}
}

func (sps *SpellPrepSelector) TogglePrepare() {
	if len(sps.allSpells) == 0 || sps.cursor >= len(sps.allSpells) {
		return
	}

	selectedSpell := sps.allSpells[sps.cursor]

	// Check if spell is already prepared
	isPrepared := false
	spellIndex := -1
	for i, spell := range sps.character.SpellBook.Spells {
		if spell.Name == selectedSpell.Name {
			isPrepared = spell.Prepared
			spellIndex = i
			break
		}
	}

	if isPrepared {
		// Unprepare
		if spellIndex >= 0 {
			sps.character.SpellBook.Spells[spellIndex].Prepared = false
		}
	} else {
		// Check if can prepare more
		preparedCount := sps.getPreparedCount()
		if preparedCount >= sps.character.SpellBook.MaxPreparedSpells {
			return // Can't prepare more
		}

		// Prepare
		if spellIndex >= 0 {
			sps.character.SpellBook.Spells[spellIndex].Prepared = true
		} else {
			// Add to spellbook and prepare
			selectedSpell.Prepared = true
			sps.character.SpellBook.Spells = append(sps.character.SpellBook.Spells, selectedSpell)
		}
	}
}

func (sps *SpellPrepSelector) getPreparedCount() int {
	count := 0
	for _, spell := range sps.character.SpellBook.Spells {
		if spell.Prepared && spell.Level > 0 {
			count++
		}
	}
	return count
}

func (sps *SpellPrepSelector) isSpellPrepared(spellName string) bool {
	for _, spell := range sps.character.SpellBook.Spells {
		if spell.Name == spellName && spell.Prepared {
			return true
		}
	}
	return false
}

func (sps *SpellPrepSelector) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	preparedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Left column - spell list
	var leftLines []string
	preparedCount := sps.getPreparedCount()
	filterInfo := "All Levels"
	if sps.selectedLevel > 0 {
		filterInfo = fmt.Sprintf("Level %d", sps.selectedLevel)
	}
	leftLines = append(leftLines, headerStyle.Render(fmt.Sprintf("Prepared: %d/%d | %s", preparedCount, sps.character.SpellBook.MaxPreparedSpells, filterInfo)))
	leftLines = append(leftLines, "")

	// Show spells
	if len(sps.allSpells) == 0 {
		leftLines = append(leftLines, dimStyle.Render("  No spells available"))
	} else {
		for i, spell := range sps.allSpells {
			cursor := "  "
			style := normalStyle
			if i == sps.cursor {
				cursor = "❯ "
				style = selectedStyle
			}

			prepMarker := "[ ]"
			if sps.isSpellPrepared(spell.Name) {
				prepMarker = "[✓]"
				if i != sps.cursor {
					style = preparedStyle
				}
			}

			ritualMarker := ""
			if spell.Ritual {
				ritualMarker = " (R)"
			}

			line := fmt.Sprintf("%s%s %s%s", cursor, prepMarker, spell.Name, ritualMarker)
			leftLines = append(leftLines, style.Render(line))
		}
	}

	leftColumn := strings.Join(leftLines, "\n")

	// Right column - spell description
	var rightLines []string
	if len(sps.allSpells) > 0 && sps.cursor < len(sps.allSpells) {
		spell := sps.allSpells[sps.cursor]

		rightLines = append(rightLines, headerStyle.Render(spell.Name))
		rightLines = append(rightLines, dimStyle.Render(fmt.Sprintf("Level %d %s", spell.Level, spell.School)))
		rightLines = append(rightLines, "")

		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Casting Time: %s", spell.CastingTime)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Range: %s", spell.Range)))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Components: %s", spell.GetComponentsString())))
		rightLines = append(rightLines, normalStyle.Render(fmt.Sprintf("Duration: %s", spell.Duration)))
		if spell.Concentration {
			rightLines = append(rightLines, normalStyle.Render("Concentration: Yes"))
		}
		rightLines = append(rightLines, "")

		// Wrap description
		descWords := strings.Fields(spell.Description)
		line := ""
		for _, word := range descWords {
			if len(line)+len(word)+1 > 35 {
				rightLines = append(rightLines, normalStyle.Render(line))
				line = word
			} else {
				if line != "" {
					line += " "
				}
				line += word
			}
		}
		if line != "" {
			rightLines = append(rightLines, normalStyle.Render(line))
		}
	}

	rightColumn := strings.Join(rightLines, "\n")

	// Combine columns
	leftBox := lipgloss.NewStyle().Width(35).Render(leftColumn)
	rightBox := lipgloss.NewStyle().Width(40).Padding(0, 1).Render(rightColumn)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// Build final content
	var content strings.Builder
	content.WriteString(titleStyle.Render(fmt.Sprintf("📖 PREPARE SPELLS - %s", strings.ToUpper(sps.character.Class))))
	content.WriteString("\n\n")
	content.WriteString(columns)
	content.WriteString("\n\n")

	// Help text
	content.WriteString(helpStyle.Render("↑/↓: Navigate • Space: Prepare/Unprepare • 0-9: Filter Level • Enter: Done • Esc: Cancel"))

	// Create popup with larger size
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(82)

	return lipgloss.Place(
		100,
		30,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
