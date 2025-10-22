// internal/ui/components/slotrestorer.go
package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// SlotRestorer allows restoring individual spell slots
type SlotRestorer struct {
	visible   bool
	character *models.Character
	cursor    int
	slots     []int // Available slot levels (1-9)
}

// NewSlotRestorer creates a new slot restorer
func NewSlotRestorer(char *models.Character) *SlotRestorer {
	return &SlotRestorer{
		visible:   false,
		character: char,
		cursor:    0,
		slots:     []int{},
	}
}

// Show displays the slot restorer
func (sr *SlotRestorer) Show() {
	sr.visible = true
	sr.cursor = 0
	sr.loadAvailableSlots()
}

// Hide hides the slot restorer
func (sr *SlotRestorer) Hide() {
	sr.visible = false
}

// IsVisible returns whether the slot restorer is visible
func (sr *SlotRestorer) IsVisible() bool {
	return sr.visible
}

// loadAvailableSlots finds all spell levels that have slots
func (sr *SlotRestorer) loadAvailableSlots() {
	sr.slots = []int{}

	allSlots := []struct {
		level int
		slot  *models.SpellSlot
	}{
		{1, &sr.character.SpellBook.Slots.Level1},
		{2, &sr.character.SpellBook.Slots.Level2},
		{3, &sr.character.SpellBook.Slots.Level3},
		{4, &sr.character.SpellBook.Slots.Level4},
		{5, &sr.character.SpellBook.Slots.Level5},
		{6, &sr.character.SpellBook.Slots.Level6},
		{7, &sr.character.SpellBook.Slots.Level7},
		{8, &sr.character.SpellBook.Slots.Level8},
		{9, &sr.character.SpellBook.Slots.Level9},
	}

	for _, s := range allSlots {
		if s.slot.Maximum > 0 {
			sr.slots = append(sr.slots, s.level)
		}
	}
}

// Update handles input
func (sr *SlotRestorer) Update(msg tea.Msg) (SlotRestorer, tea.Cmd) {
	if !sr.visible {
		return *sr, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			sr.Prev()
		case "down", "j":
			sr.Next()
		case "enter", " ":
			// Restore one slot for selected level
			return *sr, nil
		case "esc":
			sr.Hide()
		}
	}

	return *sr, nil
}

// Next moves cursor down
func (sr *SlotRestorer) Next() {
	if sr.cursor < len(sr.slots)-1 {
		sr.cursor++
	}
}

// Prev moves cursor up
func (sr *SlotRestorer) Prev() {
	if sr.cursor > 0 {
		sr.cursor--
	}
}

// GetSelectedLevel returns the currently selected spell level
func (sr *SlotRestorer) GetSelectedLevel() int {
	if sr.cursor >= 0 && sr.cursor < len(sr.slots) {
		return sr.slots[sr.cursor]
	}
	return 0
}

// RestoreSlot restores one slot for the selected level
func (sr *SlotRestorer) RestoreSlot() (string, bool) {
	level := sr.GetSelectedLevel()
	if level == 0 {
		return "No spell level selected", false
	}

	slot := sr.character.SpellBook.GetSlotByLevel(level)
	if slot == nil {
		return fmt.Sprintf("No slots at level %d", level), false
	}

	if slot.Current >= slot.Maximum {
		return fmt.Sprintf("Level %d slots already full (%d/%d)", level, slot.Current, slot.Maximum), false
	}

	slot.Current++
	return fmt.Sprintf("Restored 1 level %d slot (%d/%d)", level, slot.Current, slot.Maximum), true
}

// View renders the slot restorer
func (sr *SlotRestorer) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	fullStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var content string
	content += titleStyle.Render("RESTORE SPELL SLOT") + "\n\n"

	if len(sr.slots) == 0 {
		content += dimStyle.Render("No spell slots available") + "\n"
	} else {
		for i, level := range sr.slots {
			slot := sr.character.SpellBook.GetSlotByLevel(level)
			if slot == nil {
				continue
			}

			cursor := "  "
			style := normalStyle
			if i == sr.cursor {
				cursor = "❯ "
				style = selectedStyle
			}

			status := fmt.Sprintf("(%d/%d)", slot.Current, slot.Maximum)
			if slot.Current >= slot.Maximum {
				style = fullStyle
				status += " FULL"
			} else {
				status += fmt.Sprintf(" [+%d]", slot.Maximum-slot.Current)
			}

			line := fmt.Sprintf("%sLevel %d  %s", cursor, level, status)
			content += style.Render(line) + "\n"
		}
	}

	content += "\n" + dimStyle.Render("↑/↓: Select • Enter/Space: Restore +1 • Esc: Cancel")

	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(50)

	return lipgloss.Place(
		60,
		20,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content),
	)
}
