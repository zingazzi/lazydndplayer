// internal/ui/panels/characterstats.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CharacterStatsPanel displays key character statistics
type CharacterStatsPanel struct {
	character *models.Character
}

// NewCharacterStatsPanel creates a new character stats panel
func NewCharacterStatsPanel(char *models.Character) *CharacterStatsPanel {
	return &CharacterStatsPanel{
		character: char,
	}
}

// View renders the character stats panel
func (p *CharacterStatsPanel) View(width, height int) string {
	char := p.character

	// Styles
	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	raceStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Italic(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(12)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true)

	statBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Align(lipgloss.Center)

	criticalStatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	// Calculate initiative modifier (DEX modifier)
	initiativeMod := char.AbilityScores.GetModifier(models.Dexterity)

	// Build stat boxes for important stats (smaller for 2-row layout)
	boxWidth := 10

	hpBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("❤ HP") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%d/%d", char.CurrentHP, char.MaxHP)),
	)

	acBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("🛡 AC") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%d", char.ArmorClass)),
	)

	initBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render("⚡ INIT") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%+d", initiativeMod)),
	)

	speedBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Bold(true).Render("👣 SPD") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%dft", char.Speed)),
	)

	profBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("⭐ PRF") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("+%d", char.ProficiencyBonus)),
	)

	// Character name and race
	var lines []string
	lines = append(lines, nameStyle.Render("⚔ "+char.Name)+" "+raceStyle.Render(char.Race))
	lines = append(lines, "")

	// Class and level
	classInfo := fmt.Sprintf("%s, Level %d", char.Class, char.Level)
	lines = append(lines, labelStyle.Render("Class:")+" "+valueStyle.Render(classInfo))

	// XP information
	xpToNext := getLevelXP(char.Level+1) - char.Experience
	xpInfo := fmt.Sprintf("%d XP (next: %d)", char.Experience, xpToNext)
	lines = append(lines, labelStyle.Render("Experience:")+" "+valueStyle.Render(xpInfo))
	lines = append(lines, "")

	// Stat boxes in 2 rows
	// Row 1: HP, AC, INIT
	statBoxesRow1 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		hpBox,
		" ",
		acBox,
		" ",
		initBox,
	)
	lines = append(lines, statBoxesRow1)
	lines = append(lines, "")

	// Row 2: SPD, PROF
	statBoxesRow2 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		speedBox,
		" ",
		profBox,
	)
	lines = append(lines, statBoxesRow2)

	content := strings.Join(lines, "\n")

	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// Update handles updates for the character stats panel
func (p *CharacterStatsPanel) Update(char *models.Character) {
	p.character = char
}

// getLevelXP returns the XP required to reach a given level (simplified)
func getLevelXP(level int) int {
	xpTable := map[int]int{
		1:  0,
		2:  300,
		3:  900,
		4:  2700,
		5:  6500,
		6:  14000,
		7:  23000,
		8:  34000,
		9:  48000,
		10: 64000,
		11: 85000,
		12: 100000,
		13: 120000,
		14: 140000,
		15: 165000,
		16: 195000,
		17: 225000,
		18: 265000,
		19: 305000,
		20: 355000,
	}
	if xp, exists := xpTable[level]; exists {
		return xp
	}
	return 355000 // Max level
}
