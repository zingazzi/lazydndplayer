// internal/ui/panels/spells.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellsPanel displays character spells
type SpellsPanel struct {
	character *models.Character
	viewport  viewport.Model
	ready     bool
	allSpells []models.Spell // All spells available to the class
}

// NewSpellsPanel creates a new spells panel
func NewSpellsPanel(char *models.Character) *SpellsPanel {
	return &SpellsPanel{
		character: char,
		allSpells: []models.Spell{},
	}
}

// Init initializes the spells panel
func (p *SpellsPanel) Init() {
	p.viewport = viewport.New(80, 20)
	p.viewport.YPosition = 0
	p.ready = true
	p.loadAvailableSpells()
}

// loadAvailableSpells loads all spells available to the character's class
func (p *SpellsPanel) loadAvailableSpells() {
	if p.character.Class == "" {
		p.allSpells = []models.Spell{}
		return
	}

	// Load all spells from data file
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		p.allSpells = []models.Spell{}
		return
	}

	// Filter by class and level
	classLower := strings.ToLower(p.character.Class)
	maxSpellLevel := p.getMaxSpellLevel()

	p.allSpells = []models.Spell{}
	for _, spell := range allSpells {
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

		p.allSpells = append(p.allSpells, spell)
	}
}

// getMaxSpellLevel returns the maximum spell level the character can cast
func (p *SpellsPanel) getMaxSpellLevel() int {
	// Simplified: based on character level
	// Full casters get 9th level spells at level 17+
	level := p.character.Level

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

// View renders the spells panel
func (p *SpellsPanel) View(width, height int) string {
	if !p.ready {
		p.Init()
	}

	// Update viewport size if needed (maximize available space)
	viewportWidth := width - 4
	viewportHeight := height - 2
	if p.viewport.Width != viewportWidth || p.viewport.Height != viewportHeight {
		p.viewport.Width = viewportWidth
		p.viewport.Height = viewportHeight
	}

	char := p.character

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var lines []string
	lines = append(lines, titleStyle.Render("SPELLS"))
	lines = append(lines, "")

	// Spellcasting info
	if char.SpellBook.SpellcastingMod != "" {
		lines = append(lines, headerStyle.Render(fmt.Sprintf("Spellcasting: %s", char.SpellBook.SpellcastingMod)))
		lines = append(lines, normalStyle.Render(fmt.Sprintf("Spell Save DC: %d  Attack: +%d",
			char.SpellBook.SpellSaveDC, char.SpellBook.SpellAttackBonus)))
		lines = append(lines, "")
	}

	// For prepared casters, show preparation info
	if char.SpellBook.IsPreparedCaster {
		preparedCount := p.getPreparedCount()
		prepInfo := fmt.Sprintf("Prepared: %d/%d", preparedCount, char.SpellBook.MaxPreparedSpells)
		lines = append(lines, headerStyle.Render(prepInfo))
		lines = append(lines, "")
	}

	// Spell slots
	lines = append(lines, headerStyle.Render("SPELL SLOTS"))
	slotLines := p.renderSpellSlots()
	lines = append(lines, slotLines...)
	lines = append(lines, "")

	// Cantrips section
	lines = append(lines, headerStyle.Render("CANTRIPS"))
	if len(char.SpellBook.Cantrips) == 0 {
		lines = append(lines, dimStyle.Render("  No cantrips known"))
	} else {
		for _, cantripName := range char.SpellBook.Cantrips {
			lines = append(lines, normalStyle.Render(fmt.Sprintf("  ● %s", cantripName)))
		}
	}
	lines = append(lines, dimStyle.Render(fmt.Sprintf("  Press 'c' to change cantrips (%d known)", char.SpellBook.CantripsKnown)))
	lines = append(lines, "")

	// Prepared spells section (show only prepared spells)
	lines = append(lines, headerStyle.Render("PREPARED SPELLS"))

	if char.SpellBook.IsPreparedCaster {
		preparedCount := p.getPreparedCount()
		lines = append(lines, dimStyle.Render(fmt.Sprintf("  (%d/%d prepared)", preparedCount, char.SpellBook.MaxPreparedSpells)))
		lines = append(lines, "")
	}

	spellLines, totalCount := p.renderPreparedSpellsByLevel()
	lines = append(lines, spellLines...)

	if totalCount == 0 {
		lines = append(lines, dimStyle.Render("  No spells prepared"))
		lines = append(lines, dimStyle.Render("  Press 'v' to prepare spells"))
	}

	lines = append(lines, "")
	lines = append(lines, dimStyle.Render("Keys: 'v': Prepare Spells • 'c': Change Cantrips • 'r': Rest"))

	content := strings.Join(lines, "\n")

	// Set viewport content
	p.viewport.SetContent(content)

	// Render viewport with border
	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(p.viewport.View())
}

// renderSpellSlots renders the spell slots display
func (p *SpellsPanel) renderSpellSlots() []string {
	char := p.character
	var lines []string

	slotLevels := []struct {
		level int
		slot  *models.SpellSlot
	}{
		{1, &char.SpellBook.Slots.Level1},
		{2, &char.SpellBook.Slots.Level2},
		{3, &char.SpellBook.Slots.Level3},
		{4, &char.SpellBook.Slots.Level4},
		{5, &char.SpellBook.Slots.Level5},
		{6, &char.SpellBook.Slots.Level6},
		{7, &char.SpellBook.Slots.Level7},
		{8, &char.SpellBook.Slots.Level8},
		{9, &char.SpellBook.Slots.Level9},
	}

	hasSlots := false
	for _, sl := range slotLevels {
		if sl.slot.Maximum > 0 {
			hasSlots = true
			slots := strings.Repeat("●", sl.slot.Current) + strings.Repeat("○", sl.slot.Maximum-sl.slot.Current)
			lines = append(lines, fmt.Sprintf("  Level %d: %s (%d/%d)",
				sl.level, slots, sl.slot.Current, sl.slot.Maximum))
		}
	}

	if !hasSlots {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("  No spell slots available"))
	}

	return lines
}

// renderSpellsByLevel renders spells grouped by level
func (p *SpellsPanel) renderSpellsByLevel() ([]string, int) {
	var lines []string
	totalCount := 0

	// Group spells by level
	spellsByLevel := make(map[int][]models.Spell)
	for _, spell := range p.allSpells {
		if spell.Level > 0 { // Skip cantrips
			spellsByLevel[spell.Level] = append(spellsByLevel[spell.Level], spell)
		}
	}

	// Display by level
	for level := 1; level <= 9; level++ {
		spells := spellsByLevel[level]
		if len(spells) == 0 {
			continue
		}

		levelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)
		lines = append(lines, levelStyle.Render(fmt.Sprintf("Level %d:", level)))

		for _, spell := range spells {
			isPrepared := p.isSpellPrepared(spell.Name)
			prepMarker := " "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

			if isPrepared {
				prepMarker = "●"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
			}

			ritualMarker := ""
			if spell.Ritual {
				ritualMarker = " (R)"
			}

			line := fmt.Sprintf("  %s %s%s", prepMarker, spell.Name, ritualMarker)
			lines = append(lines, style.Render(line))
			totalCount++
		}
		lines = append(lines, "")
	}

	return lines, totalCount
}

// renderPreparedSpellsByLevel renders only prepared spells grouped by level
func (p *SpellsPanel) renderPreparedSpellsByLevel() ([]string, int) {
	var lines []string
	totalCount := 0

	// Group prepared spells by level
	spellsByLevel := make(map[int][]models.Spell)
	for _, spell := range p.character.SpellBook.Spells {
		if spell.Level > 0 && spell.Prepared { // Skip cantrips and unprepared
			spellsByLevel[spell.Level] = append(spellsByLevel[spell.Level], spell)
		}
	}

	// Display by level
	for level := 1; level <= 9; level++ {
		spells := spellsByLevel[level]
		if len(spells) == 0 {
			continue
		}

		levelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)
		lines = append(lines, levelStyle.Render(fmt.Sprintf("Level %d:", level)))

		preparedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

		for _, spell := range spells {
			ritualMarker := ""
			if spell.Ritual {
				ritualMarker = " (R)"
			}

			line := fmt.Sprintf("  ● %s%s", spell.Name, ritualMarker)
			lines = append(lines, preparedStyle.Render(line))
			totalCount++
		}
		lines = append(lines, "")
	}

	return lines, totalCount
}

// isSpellPrepared checks if a spell is prepared
func (p *SpellsPanel) isSpellPrepared(spellName string) bool {
	for _, spell := range p.character.SpellBook.Spells {
		if spell.Name == spellName && spell.Prepared {
			return true
		}
	}
	return false
}

// getPreparedCount returns the number of currently prepared spells
func (p *SpellsPanel) getPreparedCount() int {
	count := 0
	for _, spell := range p.character.SpellBook.Spells {
		if spell.Prepared && spell.Level > 0 { // Don't count cantrips
			count++
		}
	}
	return count
}

// TogglePrepare toggles the prepared status of the currently selected spell
func (p *SpellsPanel) TogglePrepare(spellName string) bool {
	// Check if spell exists in spellbook
	spellIndex := -1
	for i, spell := range p.character.SpellBook.Spells {
		if spell.Name == spellName {
			spellIndex = i
			break
		}
	}

	if spellIndex >= 0 {
		// Spell exists, toggle it
		spell := &p.character.SpellBook.Spells[spellIndex]
		if spell.Prepared {
			// Unprepare
			spell.Prepared = false
			return true
		} else {
			// Try to prepare
			if p.canPrepareMore() {
				spell.Prepared = true
				return true
			}
			return false // Can't prepare more
		}
	} else {
		// Spell doesn't exist in spellbook, add and prepare it
		if !p.canPrepareMore() {
			return false
		}

		// Find the spell in available spells
		for _, availSpell := range p.allSpells {
			if availSpell.Name == spellName {
				availSpell.Prepared = true
				p.character.SpellBook.Spells = append(p.character.SpellBook.Spells, availSpell)
				return true
			}
		}
	}

	return false
}

// canPrepareMore checks if the character can prepare more spells
func (p *SpellsPanel) canPrepareMore() bool {
	if !p.character.SpellBook.IsPreparedCaster {
		return true // Known casters can always learn more (up to limit)
	}

	preparedCount := p.getPreparedCount()
	return preparedCount < p.character.SpellBook.MaxPreparedSpells
}

// Rest restores all spell slots
func (p *SpellsPanel) Rest() {
	p.character.SpellBook.LongRest()
}

// Update handles updates for the spells panel
func (p *SpellsPanel) Update(char *models.Character) {
	p.character = char
	p.loadAvailableSpells()
}

// SetSize sets the viewport size
func (p *SpellsPanel) SetSize(width, height int) {
	if p.ready {
		p.viewport.Width = width - 4
		p.viewport.Height = height - 4
	}
}

// Handle messages for the spells panel
func (p *SpellsPanel) HandleKey(msg tea.KeyMsg) {
	if !p.ready {
		return
	}

	switch msg.String() {
	case "up", "k":
		p.viewport.LineUp(1)
	case "down", "j":
		p.viewport.LineDown(1)
	case "pgup":
		p.viewport.HalfViewUp()
	case "pgdown":
		p.viewport.HalfViewDown()
	}
}
