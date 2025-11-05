// internal/ui/components/spellbookeditor.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SpellbookEditor allows managing spellbook for spellbook casters (Wizards and Arcane Tricksters)
type SpellbookEditor struct {
	visible       bool
	character     *models.Character
	allWizardSpells []models.Spell // All wizard spells from data file (used by Wizards and Arcane Tricksters)
	filteredSpells  []models.Spell // Filtered based on current filter
	selectedIndex   int
	filterLevel     int  // 0 = all, 1-9 = specific level, -1 = cantrips
	currentTab      int  // 0 = cantrip, 1-9 = spell level (for tab navigation)
	mode            string // "view" or "add"
}

func NewSpellbookEditor(char *models.Character) *SpellbookEditor {
	return &SpellbookEditor{
		visible:       false,
		character:     char,
		selectedIndex: 0,
		filterLevel:   -1, // Start with cantrips
		currentTab:    0,  // Start on cantrip tab
		mode:          "view",
	}
}

func (sb *SpellbookEditor) Show() {
	sb.visible = true
	sb.selectedIndex = 0
	sb.currentTab = 0    // Start on cantrip tab
	sb.filterLevel = -1  // Filter to cantrips
	sb.mode = "view"
	sb.loadWizardSpells()
	sb.applyFilter()
}

func (sb *SpellbookEditor) Hide() {
	sb.visible = false
}

func (sb *SpellbookEditor) IsVisible() bool {
	return sb.visible
}

func (sb *SpellbookEditor) loadWizardSpells() {
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		sb.allWizardSpells = []models.Spell{}
		return
	}

	sb.allWizardSpells = []models.Spell{}

	// Determine which class spells to load
	var targetClass string
	if sb.character.HasClass("Wizard") {
		targetClass = "wizard"
	} else if sb.character.HasClass("Bard") {
		targetClass = "bard"
	} else if sb.character.IsArcaneTrickster() {
		targetClass = "wizard" // Arcane Trickster uses wizard spells
	} else {
		// Default to wizard if no specific class found
		targetClass = "wizard"
	}

	for _, spell := range allSpells {
		// Check if it's a spell for the target class
		for _, class := range spell.Classes {
			if strings.ToLower(class) == targetClass {
				sb.allWizardSpells = append(sb.allWizardSpells, spell)
				break
			}
		}
	}
}

func (sb *SpellbookEditor) applyFilter() {
	sb.filteredSpells = []models.Spell{}

	for _, spell := range sb.allWizardSpells {
		// Apply level filter
		if sb.filterLevel == -1 && spell.Level != 0 {
			continue // Show only cantrips
		}
		if sb.filterLevel > 0 && spell.Level != sb.filterLevel {
			continue // Show only specific level
		}
		if sb.filterLevel == 0 {
			// Show all spells including cantrips
		}

		sb.filteredSpells = append(sb.filteredSpells, spell)
	}

	// Reset selection if out of bounds
	if sb.selectedIndex >= len(sb.filteredSpells) {
		sb.selectedIndex = 0
	}
}

func (sb *SpellbookEditor) Next() {
	if sb.selectedIndex < len(sb.filteredSpells)-1 {
		sb.selectedIndex++
	}
}

func (sb *SpellbookEditor) Prev() {
	if sb.selectedIndex > 0 {
		sb.selectedIndex--
	}
}

func (sb *SpellbookEditor) SetFilter(level int) {
	sb.filterLevel = level
	sb.selectedIndex = 0

	// Update current tab to match filter
	if level == -1 {
		sb.currentTab = 0 // Cantrip tab
	} else {
		sb.currentTab = level // Level 1-9 tab
	}

	sb.applyFilter()
}

func (sb *SpellbookEditor) NextTab() {
	sb.currentTab++
	if sb.currentTab > 9 {
		sb.currentTab = 0
	}

	// Update filter based on tab
	if sb.currentTab == 0 {
		sb.filterLevel = -1 // Cantrips
	} else {
		sb.filterLevel = sb.currentTab // Level 1-9
	}

	sb.selectedIndex = 0
	sb.applyFilter()
}

func (sb *SpellbookEditor) PrevTab() {
	sb.currentTab--
	if sb.currentTab < 0 {
		sb.currentTab = 9
	}

	// Update filter based on tab
	if sb.currentTab == 0 {
		sb.filterLevel = -1 // Cantrips
	} else {
		sb.filterLevel = sb.currentTab // Level 1-9
	}

	sb.selectedIndex = 0
	sb.applyFilter()
}

func (sb *SpellbookEditor) TogglePrepared() {
	if sb.selectedIndex < 0 || sb.selectedIndex >= len(sb.filteredSpells) {
		return
	}

	selectedSpell := sb.filteredSpells[sb.selectedIndex]

	// Cantrips can't be prepared/unprepared
	if selectedSpell.Level == 0 {
		return
	}

	// Must be in spellbook to prepare
	if !sb.IsSpellKnown(selectedSpell.Name) {
		return
	}

	// Find this spell in character's spellbook and toggle
	for i := range sb.character.SpellBook.Spells {
		if sb.character.SpellBook.Spells[i].Name == selectedSpell.Name {
			// Toggle preparation
			sb.character.SpellBook.Spells[i].Prepared = !sb.character.SpellBook.Spells[i].Prepared
			return
		}
	}
}

func (sb *SpellbookEditor) IsSpellKnown(spellName string) bool {
	// Check cantrips
	for _, name := range sb.character.SpellBook.Cantrips {
		if name == spellName {
			return true
		}
	}
	// Check leveled spells
	for _, spell := range sb.character.SpellBook.Spells {
		if spell.Name == spellName {
			return true
		}
	}
	return false
}

func (sb *SpellbookEditor) IsSpellPrepared(spellName string) bool {
	for _, spell := range sb.character.SpellBook.Spells {
		if spell.Name == spellName && spell.Prepared {
			return true
		}
	}
	return false
}

func (sb *SpellbookEditor) GetPreparedCount() int {
	count := 0
	for _, spell := range sb.character.SpellBook.Spells {
		if spell.Prepared && spell.Level > 0 {
			count++
		}
	}
	return count
}

func (sb *SpellbookEditor) GetPreparedCountByLevel(level int) int {
	count := 0
	for _, spell := range sb.character.SpellBook.Spells {
		if spell.Prepared && spell.Level == level {
			count++
		}
	}
	return count
}

func (sb *SpellbookEditor) GetKnownCountByLevel(level int) int {
	if level == 0 {
		// Cantrips
		return len(sb.character.SpellBook.Cantrips)
	}
	count := 0
	for _, spell := range sb.character.SpellBook.Spells {
		if spell.Level == level {
			count++
		}
	}
	return count
}

func (sb *SpellbookEditor) GetSpellSlots(level int) (current, max int) {
	switch level {
	case 1:
		return sb.character.SpellBook.Slots.Level1.Current, sb.character.SpellBook.Slots.Level1.Maximum
	case 2:
		return sb.character.SpellBook.Slots.Level2.Current, sb.character.SpellBook.Slots.Level2.Maximum
	case 3:
		return sb.character.SpellBook.Slots.Level3.Current, sb.character.SpellBook.Slots.Level3.Maximum
	case 4:
		return sb.character.SpellBook.Slots.Level4.Current, sb.character.SpellBook.Slots.Level4.Maximum
	case 5:
		return sb.character.SpellBook.Slots.Level5.Current, sb.character.SpellBook.Slots.Level5.Maximum
	case 6:
		return sb.character.SpellBook.Slots.Level6.Current, sb.character.SpellBook.Slots.Level6.Maximum
	case 7:
		return sb.character.SpellBook.Slots.Level7.Current, sb.character.SpellBook.Slots.Level7.Maximum
	case 8:
		return sb.character.SpellBook.Slots.Level8.Current, sb.character.SpellBook.Slots.Level8.Maximum
	case 9:
		return sb.character.SpellBook.Slots.Level9.Current, sb.character.SpellBook.Slots.Level9.Maximum
	}
	return 0, 0
}

func (sb *SpellbookEditor) AddSpellToSpellbook() {
	if sb.selectedIndex < 0 || sb.selectedIndex >= len(sb.filteredSpells) {
		return
	}

	selectedSpell := sb.filteredSpells[sb.selectedIndex]

	// Check if already known
	if sb.IsSpellKnown(selectedSpell.Name) {
		return
	}

	// Add to spellbook (cantrips or leveled spells)
	selectedSpell.Known = true
	selectedSpell.Prepared = false

	if selectedSpell.Level == 0 {
		// It's a cantrip - add to cantrips list
		sb.character.SpellBook.Cantrips = append(sb.character.SpellBook.Cantrips, selectedSpell.Name)
		// NOTE: Don't update CantripsKnown - it's the MAX from level progression, not current count
	} else {
		// It's a leveled spell - add to spells list
		sb.character.SpellBook.Spells = append(sb.character.SpellBook.Spells, selectedSpell)
	}
}

func (sb *SpellbookEditor) RemoveSpellFromSpellbook() {
	if sb.selectedIndex < 0 || sb.selectedIndex >= len(sb.filteredSpells) {
		return
	}

	selectedSpell := sb.filteredSpells[sb.selectedIndex]

	// Check if it's known
	if !sb.IsSpellKnown(selectedSpell.Name) {
		return
	}

	if selectedSpell.Level == 0 {
		// It's a cantrip - remove from cantrips list
		for i, name := range sb.character.SpellBook.Cantrips {
			if name == selectedSpell.Name {
				sb.character.SpellBook.Cantrips = append(sb.character.SpellBook.Cantrips[:i], sb.character.SpellBook.Cantrips[i+1:]...)
				// NOTE: Don't update CantripsKnown - it's the MAX from level progression, not current count
				return
			}
		}
	} else {
		// It's a leveled spell - remove from spells list
		for i := range sb.character.SpellBook.Spells {
			if sb.character.SpellBook.Spells[i].Name == selectedSpell.Name {
				sb.character.SpellBook.Spells = append(sb.character.SpellBook.Spells[:i], sb.character.SpellBook.Spells[i+1:]...)
				return
			}
		}
	}
}

func (sb *SpellbookEditor) IsSpellCantrip(spellName string) bool {
	for _, name := range sb.character.SpellBook.Cantrips {
		if name == spellName {
			return true
		}
	}
	return false
}

func (sb *SpellbookEditor) Update(msg tea.Msg) (SpellbookEditor, tea.Cmd) {
	if !sb.visible {
		return *sb, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			sb.Prev()
		case "down", "j":
			sb.Next()
		case "tab", "ctrl+i":
			sb.NextTab()
		case "shift+tab", "backtab":
			sb.PrevTab()
		case "right", "l":
			sb.NextTab()
		case "left", "h":
			sb.PrevTab()
		case "0":
			sb.SetFilter(0) // Show all
		case "c":
			sb.SetFilter(-1) // Show cantrips
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			level := int(msg.String()[0] - '0')
			sb.SetFilter(level)
		case " ":
			// Toggle prepared state (only for leveled spells, not cantrips)
			if sb.selectedIndex >= 0 && sb.selectedIndex < len(sb.filteredSpells) {
				spell := sb.filteredSpells[sb.selectedIndex]
				if spell.Level > 0 && sb.IsSpellKnown(spell.Name) {
					sb.TogglePrepared()
				}
			}
		case "a":
			// Add spell/cantrip to spellbook
			sb.AddSpellToSpellbook()
		case "d", "x", "delete":
			// Remove spell/cantrip from spellbook
			sb.RemoveSpellFromSpellbook()
		case "esc":
			sb.Hide()
		}
	}

	return *sb, nil
}

func (sb *SpellbookEditor) View(width, height int) string {
	if !sb.visible {
		return ""
	}

	// Two-panel layout like origin selector
	leftWidth := 35
	rightWidth := width - leftWidth - 8

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	knownStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	preparedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	concentrationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	ritualStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141"))

	// Tab styles
	activeTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	// Build left panel (spell list - simple names only)
	var leftLines []string
	leftLines = append(leftLines, titleStyle.Render("SPELLBOOK"))
	leftLines = append(leftLines, "")

	// Tab bar
	tabs := []string{"C", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var tabBar strings.Builder
	for i, tab := range tabs {
		if i == sb.currentTab {
			tabBar.WriteString(activeTabStyle.Render(tab))
		} else {
			tabBar.WriteString(inactiveTabStyle.Render(tab))
		}
		if i < len(tabs)-1 {
			tabBar.WriteString(" ")
		}
	}
	leftLines = append(leftLines, tabBar.String())
	leftLines = append(leftLines, "")

	// Show level-specific information
	if sb.currentTab == 0 {
		// Cantrip tab - show cantrips known
		knownCantrips := sb.GetKnownCountByLevel(0)
		maxCantrips := sb.character.SpellBook.CantripsKnown
		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Known: %d/%d", knownCantrips, maxCantrips)))
	} else {
		// Spell level tab - show known and prepared for THIS level only
		level := sb.currentTab
		preparedThisLevel := sb.GetPreparedCountByLevel(level)
		knownThisLevel := sb.GetKnownCountByLevel(level)
		maxPrepared := sb.character.SpellBook.MaxPreparedSpells

		// Debug logging
		debug.Log("=== SPELLBOOK DISPLAY DEBUG ===")
		debug.Log("  Current Tab/Level: %d", level)
		debug.Log("  PreparationFormula: '%s'", sb.character.SpellBook.PreparationFormula)
		debug.Log("  MaxPreparedSpells (stored): %d", maxPrepared)
		debug.Log("  Character Level: %d", sb.character.Level)
		debug.Log("  Intelligence Score: %d", sb.character.AbilityScores.Intelligence)
		debug.Log("  Intelligence Modifier: %d", sb.character.AbilityScores.GetModifier("Intelligence"))

		// If maxPrepared is 0, recalculate it
		if maxPrepared == 0 && sb.character.SpellBook.PreparationFormula != "" {
			debug.Log("  MaxPrepared is 0, recalculating...")
			maxPrepared = sb.character.CalculateMaxPreparedSpells(sb.character.SpellBook.PreparationFormula)
			sb.character.SpellBook.MaxPreparedSpells = maxPrepared
			debug.Log("  Recalculated MaxPrepared: %d", maxPrepared)
		} else if maxPrepared == 0 {
			debug.Log("  WARNING: MaxPrepared is 0 and no formula available!")
		}

		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Known: %d", knownThisLevel)))
		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Prepared: %d/%d", preparedThisLevel, maxPrepared)))
	}
	leftLines = append(leftLines, "")

	if len(sb.filteredSpells) == 0 {
		leftLines = append(leftLines, dimStyle.Render("No spells available"))
	} else {
		// Show spells with simple status markers
		for i, spell := range sb.filteredSpells {
			style := normalStyle
			cursor := "  "

			if i == sb.selectedIndex {
				cursor = "→ "
				style = selectedStyle
			}

			// Simple status: ● for known, ✓ for prepared
			prefix := "  "
			if sb.IsSpellKnown(spell.Name) {
				if sb.IsSpellPrepared(spell.Name) {
					prefix = "✓ "
					if i != sb.selectedIndex {
						style = preparedStyle
					}
				} else {
					prefix = "● "
					if i != sb.selectedIndex {
						style = knownStyle
					}
				}
			}

			line := cursor + prefix + spell.Name
			leftLines = append(leftLines, style.Render(line))
		}
	}

	leftLines = append(leftLines, "")
	leftLines = append(leftLines, dimStyle.Render("──────────────────────────"))
	leftLines = append(leftLines, dimStyle.Render("[←/→] Switch level"))
	leftLines = append(leftLines, dimStyle.Render("[↑/↓] Navigate"))
	leftLines = append(leftLines, dimStyle.Render("[Space] Prepare/Unprepare"))
	leftLines = append(leftLines, dimStyle.Render("[a] Add  [d/x] Remove"))
	leftLines = append(leftLines, dimStyle.Render("[Esc] Close"))

	// Build right panel (spell details)
	var rightLines []string
	if sb.selectedIndex >= 0 && sb.selectedIndex < len(sb.filteredSpells) {
		spell := sb.filteredSpells[sb.selectedIndex]

		// Title
		rightLines = append(rightLines, titleStyle.Render(spell.Name))
		rightLines = append(rightLines, "")

		// Level and school
		levelText := "Cantrip"
		if spell.Level > 0 {
			levelText = fmt.Sprintf("Level %d %s", spell.Level, spell.School)
		} else {
			levelText = fmt.Sprintf("%s Cantrip", spell.School)
		}
		rightLines = append(rightLines, sectionTitleStyle.Render(levelText))
		rightLines = append(rightLines, "")

		// Casting details
		rightLines = append(rightLines, normalStyle.Render("Casting Time: "+spell.CastingTime))
		rightLines = append(rightLines, normalStyle.Render("Range: "+spell.Range))
		rightLines = append(rightLines, normalStyle.Render("Duration: "+spell.Duration))
		rightLines = append(rightLines, "")

		// Special properties
		if spell.Concentration {
			rightLines = append(rightLines, concentrationStyle.Render("⚠ Requires Concentration"))
		}
		if spell.Ritual {
			rightLines = append(rightLines, ritualStyle.Render("◆ Can be cast as Ritual"))
		}
		if spell.Concentration || spell.Ritual {
			rightLines = append(rightLines, "")
		}

		// Description
		rightLines = append(rightLines, sectionTitleStyle.Render("DESCRIPTION:"))
		rightLines = append(rightLines, "")
		wrappedDesc := wrapSpellbookText(spell.Description, rightWidth-4)
		for _, line := range wrappedDesc {
			rightLines = append(rightLines, normalStyle.Render(line))
		}
		rightLines = append(rightLines, "")

		// Status
		rightLines = append(rightLines, sectionTitleStyle.Render("STATUS:"))
		if sb.IsSpellKnown(spell.Name) {
			if sb.IsSpellPrepared(spell.Name) {
				rightLines = append(rightLines, preparedStyle.Render("✓ In Spellbook (Prepared)"))
			} else {
				rightLines = append(rightLines, knownStyle.Render("✓ In Spellbook (Not Prepared)"))
				if spell.Level > 0 {
					rightLines = append(rightLines, dimStyle.Render("  Press [Space] to prepare"))
				}
			}
		} else {
			rightLines = append(rightLines, dimStyle.Render("Not in spellbook"))
			rightLines = append(rightLines, dimStyle.Render("  Press [a] to add"))
		}
	}

	// Combine panels with simple box borders
	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1).
		Width(leftWidth).
		Height(height - 4).
		Render(strings.Join(leftLines, "\n"))

	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1).
		Width(rightWidth).
		Height(height - 4).
		Render(strings.Join(rightLines, "\n"))

	combined := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		combined,
	)
}

// wrapSpellbookText wraps text to specified width
func wrapSpellbookText(text string, width int) []string {
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
