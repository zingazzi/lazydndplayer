// internal/ui/components/spellprepselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type SpellPrepSelector struct {
	visible         bool
	allClericSpells []models.Spell // All cleric spells from data file
	filteredSpells  []models.Spell // Filtered based on current tab
	selectedIndex   int
	filterLevel     int  // -1 = cantrips, 1-9 = specific level
	currentTab      int  // 0 = cantrip, 1-9 = spell level (for tab navigation)
	character       *models.Character
}

func NewSpellPrepSelector(char *models.Character) *SpellPrepSelector {
	return &SpellPrepSelector{
		visible:         false,
		allClericSpells: []models.Spell{},
		filteredSpells:  []models.Spell{},
		selectedIndex:   0,
		filterLevel:     -1, // Start with cantrips
		currentTab:      0,  // Start on cantrip tab
		character:       char,
	}
}

func (sps *SpellPrepSelector) Show() {
	sps.visible = true
	sps.selectedIndex = 0
	sps.currentTab = 0    // Start on cantrip tab
	sps.filterLevel = -1  // Filter to cantrips
	sps.loadClericSpells()
	sps.applyFilter()
}

func (sps *SpellPrepSelector) Hide() {
	sps.visible = false
}

func (sps *SpellPrepSelector) IsVisible() bool {
	return sps.visible
}

func (sps *SpellPrepSelector) loadClericSpells() {
	if sps.character.Class == "" {
		sps.allClericSpells = []models.Spell{}
		return
	}

	// Load all spells from data file
	allSpells, err := models.LoadSpellsFromJSON("data/spells.json")
	if err != nil {
		sps.allClericSpells = []models.Spell{}
		return
	}

	sps.allClericSpells = []models.Spell{}

	// Filter by Cleric class
	for _, spell := range allSpells {
		// Check if this spell is for Cleric
		isForCleric := false
		for _, spellClass := range spell.Classes {
			if strings.ToLower(spellClass) == "cleric" {
				isForCleric = true
				break
			}
		}

		if !isForCleric {
			continue
		}

		// Check if character can cast this spell level
		maxSpellLevel := sps.getMaxSpellLevel()
		if spell.Level == 0 || spell.Level <= maxSpellLevel {
			sps.allClericSpells = append(sps.allClericSpells, spell)
		}
	}
}

func (sps *SpellPrepSelector) applyFilter() {
	sps.filteredSpells = []models.Spell{}

	for _, spell := range sps.allClericSpells {
		// Apply level filter
		if sps.filterLevel == -1 && spell.Level != 0 {
			continue // Show only cantrips
		}
		if sps.filterLevel > 0 && spell.Level != sps.filterLevel {
			continue // Show only specific level
		}

		sps.filteredSpells = append(sps.filteredSpells, spell)
	}

	// Reset selection if out of bounds
	if sps.selectedIndex >= len(sps.filteredSpells) {
		sps.selectedIndex = 0
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
		case "tab", "ctrl+i":
			sps.NextTab()
		case "shift+tab", "backtab":
			sps.PrevTab()
		case "right", "l":
			sps.NextTab()
		case "left", "h":
			sps.PrevTab()
		case " ":
			if sps.currentTab == 0 {
				// Cantrip tab - toggle add/remove
				sps.ToggleCantrip()
			} else {
				// Spell level tab - toggle prepare/unprepare
				sps.TogglePrepared()
			}
		case "c":
			// Jump to cantrip tab
			sps.SetTab(0)
		case "0":
			// Jump to cantrips tab
			sps.SetTab(0)
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Filter by spell level
			level := int(msg.String()[0] - '0')
			if level <= sps.getMaxSpellLevel() {
				sps.SetTab(level)
			}
		case "a":
			// Add cantrip (only works on cantrip tab)
			if sps.currentTab == 0 {
				sps.AddCantrip()
			}
		case "d", "x", "delete":
			// Remove cantrip (only works on cantrip tab)
			if sps.currentTab == 0 {
				sps.RemoveCantrip()
			}
		case "enter", "esc":
			return *sps, nil
		}
	}

	return *sps, nil
}

func (sps *SpellPrepSelector) Next() {
	if sps.selectedIndex < len(sps.filteredSpells)-1 {
		sps.selectedIndex++
	}
}

func (sps *SpellPrepSelector) Prev() {
	if sps.selectedIndex > 0 {
		sps.selectedIndex--
	}
}

func (sps *SpellPrepSelector) SetTab(tab int) {
	if tab == 0 {
		sps.selectedIndex = 0
		sps.currentTab = 0
		sps.filterLevel = -1 // Cantrips
	} else if tab <= 9 {
		sps.selectedIndex = 0
		sps.currentTab = tab
		sps.filterLevel = tab
	}
	sps.applyFilter()
}

func (sps *SpellPrepSelector) NextTab() {
	sps.currentTab++
	maxLevel := sps.getMaxSpellLevel()
	if sps.currentTab > maxLevel {
		sps.currentTab = 0 // Wrap to cantrips
	}

	// Update filterLevel based on tab
	if sps.currentTab == 0 {
		sps.filterLevel = -1 // Cantrips
	} else {
		sps.filterLevel = sps.currentTab
	}

	sps.selectedIndex = 0
	sps.applyFilter()
}

func (sps *SpellPrepSelector) PrevTab() {
	sps.currentTab--
	if sps.currentTab < 0 {
		maxLevel := sps.getMaxSpellLevel()
		sps.currentTab = maxLevel // Wrap to highest level
	}

	// Update filterLevel based on tab
	if sps.currentTab == 0 {
		sps.filterLevel = -1 // Cantrips
	} else {
		sps.filterLevel = sps.currentTab
	}

	sps.selectedIndex = 0
	sps.applyFilter()
}

func (sps *SpellPrepSelector) IsCantripKnown(cantripName string) bool {
	for _, name := range sps.character.SpellBook.Cantrips {
		if name == cantripName {
			return true
		}
	}
	return false
}

func (sps *SpellPrepSelector) ToggleCantrip() {
	if len(sps.filteredSpells) == 0 || sps.selectedIndex >= len(sps.filteredSpells) {
		return
	}

	selectedCantrip := sps.filteredSpells[sps.selectedIndex]
	if selectedCantrip.Level != 0 {
		return // Not a cantrip
	}

	if sps.IsCantripKnown(selectedCantrip.Name) {
		sps.RemoveCantrip()
	} else {
		sps.AddCantrip()
	}
}

func (sps *SpellPrepSelector) AddCantrip() {
	if len(sps.filteredSpells) == 0 || sps.selectedIndex >= len(sps.filteredSpells) {
		return
	}

	selectedCantrip := sps.filteredSpells[sps.selectedIndex]
	if selectedCantrip.Level != 0 {
		return // Not a cantrip
	}

	// Check if already known
	if sps.IsCantripKnown(selectedCantrip.Name) {
		return
	}

	// Add cantrip
	sps.character.SpellBook.Cantrips = append(sps.character.SpellBook.Cantrips, selectedCantrip.Name)
}

func (sps *SpellPrepSelector) RemoveCantrip() {
	if len(sps.filteredSpells) == 0 || sps.selectedIndex >= len(sps.filteredSpells) {
		return
	}

	selectedCantrip := sps.filteredSpells[sps.selectedIndex]
	if selectedCantrip.Level != 0 {
		return // Not a cantrip
	}

	// Remove cantrip
	for i, name := range sps.character.SpellBook.Cantrips {
		if name == selectedCantrip.Name {
			sps.character.SpellBook.Cantrips = append(sps.character.SpellBook.Cantrips[:i], sps.character.SpellBook.Cantrips[i+1:]...)
			return
		}
	}
}

func (sps *SpellPrepSelector) TogglePrepared() {
	if len(sps.filteredSpells) == 0 || sps.selectedIndex >= len(sps.filteredSpells) {
		return
	}

	selectedSpell := sps.filteredSpells[sps.selectedIndex]

	// Cantrips can't be prepared
	if selectedSpell.Level == 0 {
		return
	}

	// All cleric spells are considered "known" - find or create in spellbook
	spellIndex := -1
	for i, spell := range sps.character.SpellBook.Spells {
		if spell.Name == selectedSpell.Name {
			spellIndex = i
			break
		}
	}

	if spellIndex >= 0 {
		// Spell exists in spellbook - toggle prepared
		// Don't allow toggling if spell is always prepared (domain spell)
		if sps.character.SpellBook.Spells[spellIndex].AlwaysPrepared {
			return // Domain spells are always prepared, can't be unprepared
		}
		sps.character.SpellBook.Spells[spellIndex].Prepared = !sps.character.SpellBook.Spells[spellIndex].Prepared
	} else {
		// Spell not in spellbook yet - add it as known and prepared
		// But first check if we can prepare more
		preparedCount := sps.getPreparedCount()
		if preparedCount >= sps.character.SpellBook.MaxPreparedSpells {
			return // Can't prepare more
		}

		selectedSpell.Prepared = true
		selectedSpell.Known = true
		sps.character.SpellBook.Spells = append(sps.character.SpellBook.Spells, selectedSpell)
	}
}

func (sps *SpellPrepSelector) IsSpellKnown(spellName string) bool {
	// All cleric spells are automatically "known"
	// Check if it's in the loaded spells
	for _, spell := range sps.allClericSpells {
		if spell.Name == spellName {
			return true
		}
	}
	return false
}

func (sps *SpellPrepSelector) IsSpellPrepared(spellName string) bool {
	for _, spell := range sps.character.SpellBook.Spells {
		if spell.Name == spellName && spell.Prepared {
			return true
		}
	}
	return false
}

func (sps *SpellPrepSelector) getPreparedCount() int {
	count := 0
	for _, spell := range sps.character.SpellBook.Spells {
		if spell.Prepared && !spell.AlwaysPrepared && spell.Level > 0 {
			count++
		}
	}
	return count
}

func (sps *SpellPrepSelector) GetKnownCountByLevel(level int) int {
	count := 0
	if level == 0 {
		// Cantrips
		return len(sps.character.SpellBook.Cantrips)
	}
	// Leveled spells - count all cleric spells of this level
	for _, spell := range sps.allClericSpells {
		if spell.Level == level {
			count++
		}
	}
	return count
}

func (sps *SpellPrepSelector) GetPreparedCountByLevel(level int) int {
	count := 0
	for _, spell := range sps.character.SpellBook.Spells {
		if spell.Level == level && spell.Prepared && !spell.AlwaysPrepared {
			count++
		}
	}
	return count
}


func (sps *SpellPrepSelector) View() string {
	if !sps.visible {
		return ""
	}

	// Use same structure as spellbook editor
	width := 100
	height := 30

	// Two-panel layout like spellbook editor
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

	// Tab styles - same as spellbook editor
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

	// Tab bar - same as spellbook editor
	tabs := []string{"C", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var tabBar strings.Builder
	maxLevel := sps.getMaxSpellLevel()
	for i, tab := range tabs {
		// Only show tabs up to max level
		tabLevel := i
		if i > 0 && tabLevel > maxLevel {
			break
		}

		if i == sps.currentTab {
			tabBar.WriteString(activeTabStyle.Render(tab))
		} else {
			tabBar.WriteString(inactiveTabStyle.Render(tab))
		}
		if i < len(tabs)-1 && (i == 0 || i <= maxLevel) {
			tabBar.WriteString(" ")
		}
	}
	leftLines = append(leftLines, tabBar.String())
	leftLines = append(leftLines, "")

	// Show level-specific information
	if sps.currentTab == 0 {
		// Cantrip tab - show cantrips known
		knownCantrips := sps.GetKnownCountByLevel(0)
		maxCantrips := sps.character.SpellBook.CantripsKnown
		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Known: %d/%d", knownCantrips, maxCantrips)))
	} else {
		// Spell level tab - show known (all cleric spells) and prepared for THIS level only
		level := sps.currentTab
		preparedThisLevel := sps.GetPreparedCountByLevel(level)
		knownThisLevel := sps.GetKnownCountByLevel(level)
		maxPrepared := sps.character.SpellBook.MaxPreparedSpells

		// If maxPrepared is 0, recalculate it
		if maxPrepared == 0 && sps.character.SpellBook.PreparationFormula != "" {
			debug.Log("MaxPrepared is 0, recalculating...")
			maxPrepared = sps.character.CalculateMaxPreparedSpells(sps.character.SpellBook.PreparationFormula)
			sps.character.SpellBook.MaxPreparedSpells = maxPrepared
		}

		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Known: %d", knownThisLevel)))
		leftLines = append(leftLines, dimStyle.Render(fmt.Sprintf("Prepared: %d/%d", preparedThisLevel, maxPrepared)))
	}
	leftLines = append(leftLines, "")

	if len(sps.filteredSpells) == 0 {
		leftLines = append(leftLines, dimStyle.Render("No spells available"))
	} else {
		// Show spells with simple status markers - same as spellbook editor
		for i, spell := range sps.filteredSpells {
			style := normalStyle
			cursor := "  "

			if i == sps.selectedIndex {
				cursor = "→ "
				style = selectedStyle
			}

			// Simple status: ● for known, ✓ for prepared
			prefix := "  "
			if sps.currentTab == 0 {
				// Cantrip tab
				if sps.IsCantripKnown(spell.Name) {
					prefix = "✓ "
					if i != sps.selectedIndex {
						style = knownStyle
					}
				}
			} else {
				// Spell level tab - all spells are "known" for clerics
				if sps.IsSpellKnown(spell.Name) {
					// Check if spell is always prepared (domain spell)
					isAlwaysPrepared := false
					for _, sbSpell := range sps.character.SpellBook.Spells {
						if sbSpell.Name == spell.Name && sbSpell.AlwaysPrepared {
							isAlwaysPrepared = true
							break
						}
					}

					if isAlwaysPrepared {
						prefix = "★ " // Star for always-prepared (domain spells)
						if i != sps.selectedIndex {
							style = preparedStyle
						}
					} else if sps.IsSpellPrepared(spell.Name) {
						prefix = "✓ "
						if i != sps.selectedIndex {
							style = preparedStyle
						}
					} else {
						prefix = "● "
						if i != sps.selectedIndex {
							style = knownStyle
						}
					}
				}
			}

			line := cursor + prefix + spell.Name
			leftLines = append(leftLines, style.Render(line))
		}
	}

	leftLines = append(leftLines, "")
	leftLines = append(leftLines, dimStyle.Render("──────────────────────────"))
	if sps.currentTab == 0 {
		// Cantrip controls
		leftLines = append(leftLines, dimStyle.Render("[←/→] Switch level"))
		leftLines = append(leftLines, dimStyle.Render("[↑/↓] Navigate"))
		leftLines = append(leftLines, dimStyle.Render("[Space] Add/Remove"))
		leftLines = append(leftLines, dimStyle.Render("[a] Add  [d/x] Remove"))
		leftLines = append(leftLines, dimStyle.Render("[Esc] Close"))
	} else {
		// Spell level controls
		leftLines = append(leftLines, dimStyle.Render("[←/→] Switch level"))
		leftLines = append(leftLines, dimStyle.Render("[↑/↓] Navigate"))
		leftLines = append(leftLines, dimStyle.Render("[Space] Prepare/Unprepare"))
		leftLines = append(leftLines, dimStyle.Render("[Esc] Close"))
	}

	// Build right panel (spell details) - same as spellbook editor
	var rightLines []string
	if sps.selectedIndex >= 0 && sps.selectedIndex < len(sps.filteredSpells) {
		spell := sps.filteredSpells[sps.selectedIndex]

		// Title
		rightLines = append(rightLines, titleStyle.Render(spell.Name))
		rightLines = append(rightLines, "")

		// Level and school
		levelText := "Cantrip"
		if spell.Level > 0 {
			levelText = fmt.Sprintf("Level %d %s", spell.Level, string(spell.School))
		} else {
			levelText = fmt.Sprintf("%s Cantrip", string(spell.School))
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

		// Components
		rightLines = append(rightLines, "")
		rightLines = append(rightLines, sectionTitleStyle.Render("COMPONENTS:"))
		rightLines = append(rightLines, "")
		rightLines = append(rightLines, normalStyle.Render(spell.GetComponentsString()))

		// Status - same as spellbook editor but adapted for cleric
		rightLines = append(rightLines, "")
		rightLines = append(rightLines, sectionTitleStyle.Render("STATUS:"))
		if sps.currentTab == 0 {
			// Cantrip
			if sps.IsCantripKnown(spell.Name) {
				rightLines = append(rightLines, knownStyle.Render("✓ Known"))
			} else {
				rightLines = append(rightLines, dimStyle.Render("Not known"))
				rightLines = append(rightLines, dimStyle.Render("  Press [Space] or [a] to add"))
			}
		} else {
			// Leveled spell - all cleric spells are automatically known
			if sps.IsSpellPrepared(spell.Name) {
				rightLines = append(rightLines, preparedStyle.Render("✓ Known (Prepared)"))
			} else {
				rightLines = append(rightLines, knownStyle.Render("● Known (Not Prepared)"))
				rightLines = append(rightLines, dimStyle.Render("  Press [Space] to prepare"))
			}
		}
	}

	// Combine panels with simple box borders - same as spellbook editor
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
