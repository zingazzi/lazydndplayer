// internal/ui/panels/characterstats.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// CharStatsEditMode represents what is being edited
type CharStatsEditMode int

const (
	CharStatsNormal CharStatsEditMode = iota
	CharStatsEditName
	CharStatsEditRace
	CharStatsEditHP
	CharStatsEditXP
)

// CharacterStatsPanel displays key character statistics
type CharacterStatsPanel struct {
	character *models.Character
	editMode  CharStatsEditMode
	nameInput textinput.Model
	raceInput textinput.Model
	hpInput   textinput.Model
	xpInput   textinput.Model
}

// NewCharacterStatsPanel creates a new character stats panel
func NewCharacterStatsPanel(char *models.Character) *CharacterStatsPanel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Character name"
	nameInput.CharLimit = 30
	nameInput.Width = 30

	raceInput := textinput.New()
	raceInput.Placeholder = "Race"
	raceInput.CharLimit = 20
	raceInput.Width = 20

	hpInput := textinput.New()
	hpInput.Placeholder = "+5 or -3"
	hpInput.CharLimit = 5
	hpInput.Width = 15

	xpInput := textinput.New()
	xpInput.Placeholder = "+100 or -50"
	xpInput.CharLimit = 8
	xpInput.Width = 20

	return &CharacterStatsPanel{
		character: char,
		editMode:  CharStatsNormal,
		nameInput: nameInput,
		raceInput: raceInput,
		hpInput:   hpInput,
		xpInput:   xpInput,
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

	// Build stat boxes for important stats (smaller for 2-row layout)
	boxWidth := 10

	hpBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("❤ HP") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%d/%d", char.CurrentHP, char.MaxHP)),
	)

	// Focus Points box for Monk
	var fpBox string
	if char.IsMonk() {
		monk := char.GetMonkMechanics()
		currentFP, maxFP := monk.GetFocusPoints()
		if maxFP > 0 {
			fpBox = statBoxStyle.Copy().Width(boxWidth).Render(
				lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Render("✧ FP") + "\n" +
					criticalStatStyle.Render(fmt.Sprintf("%d/%d", currentFP, maxFP)),
			)
		}
	}

	acBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("🛡 AC") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%d", char.ArmorClass)),
	)

	initBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render("⚡ INIT") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%+d", char.Initiative)),
	)

	speedBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Bold(true).Render("👣 SPD") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("%dft", char.Speed)),
	)

	profBox := statBoxStyle.Copy().Width(boxWidth).Render(
		lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("⭐ PRF") + "\n" +
			criticalStatStyle.Render(fmt.Sprintf("+%d", char.ProficiencyBonus)),
	)

	// Calculate passive scores
	passivePerception := p.calculatePassiveScore(char, models.Perception)
	passiveInvestigation := p.calculatePassiveScore(char, models.Investigation)
	passiveInsight := p.calculatePassiveScore(char, models.Insight)

	// Passive stat styles
	passiveTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	passiveNumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Bold(true)

	// Character name and race (editable)
	var lines []string

	// Show input fields when editing, otherwise show static text
	if p.editMode == CharStatsEditName {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Name: ")+p.nameInput.View())
	} else if p.editMode == CharStatsEditRace {
		lines = append(lines, nameStyle.Render("⚔ "+char.Name)+" "+lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render("Race: ")+p.raceInput.View())
	} else {
		lines = append(lines, nameStyle.Render("⚔ "+char.Name)+" "+raceStyle.Render(char.Race))
	}
	lines = append(lines, "")

	// Class and level (class can be changed with 'c', level up with 'L')
	// Class info (fighting style is shown in Traits panel)
	classDisplay := char.GetClassDisplayString()
	if classDisplay == "" {
		classDisplay = "None"
	}
	classInfoFull := fmt.Sprintf("%s (Level %d)", classDisplay, char.TotalLevel)
	if p.editMode == CharStatsNormal {
		lines = append(lines, labelStyle.Render("Class:")+" "+valueStyle.Render(classInfoFull)+" "+lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(c:change, L:level up)"))
	} else {
		lines = append(lines, labelStyle.Render("Class:")+" "+valueStyle.Render(classInfoFull))
	}

	// XP information
	xpToNext := getLevelXP(char.TotalLevel+1) - char.Experience
	xpInfo := fmt.Sprintf("%d XP (next: %d)", char.Experience, xpToNext)

	// Add level-up indicator if enough XP, or warning if too low
	levelUpIndicator := ""
	requiredXP := getLevelXP(char.TotalLevel)
	if char.Experience >= getLevelXP(char.TotalLevel+1) {
		// Ready to level up
		levelUpIndicator = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("↑ LEVEL UP!")
	} else if char.Experience < requiredXP && char.TotalLevel > 1 {
		// XP too low for current level
		deficit := requiredXP - char.Experience
		levelUpIndicator = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render(fmt.Sprintf("⚠ LOW XP (need %d)", deficit))
	}

	if p.editMode == CharStatsNormal {
		lines = append(lines, labelStyle.Render("Experience:")+" "+valueStyle.Render(xpInfo)+levelUpIndicator+" "+lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("([/]: adjust)"))
	} else {
		lines = append(lines, labelStyle.Render("Experience:")+" "+valueStyle.Render(xpInfo)+levelUpIndicator)
	}
	lines = append(lines, "")

	// Stat boxes in 2 rows
	// Row 1: HP, [FP for Monk], AC, INIT
	statBoxesRow1Parts := []string{hpBox}
	if fpBox != "" {
		statBoxesRow1Parts = append(statBoxesRow1Parts, " ", fpBox)
	}
	statBoxesRow1Parts = append(statBoxesRow1Parts, " ", acBox, " ", initBox)
	statBoxesRow1 := lipgloss.JoinHorizontal(lipgloss.Top, statBoxesRow1Parts...)
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
	lines = append(lines, "")

	// Passive stats (one per line, text in gray, numbers in purple)
	lines = append(lines, passiveTextStyle.Render("Passive Perception ")+" "+passiveNumberStyle.Render(fmt.Sprintf("%d", passivePerception)))
	lines = append(lines, passiveTextStyle.Render("Passive Investigation ")+" "+passiveNumberStyle.Render(fmt.Sprintf("%d", passiveInvestigation)))
	lines = append(lines, passiveTextStyle.Render("Passive Insight ")+" "+passiveNumberStyle.Render(fmt.Sprintf("%d", passiveInsight)))
	lines = append(lines, "")

	// Inspiration
	inspirationIcon := "☐"
	inspirationColor := lipgloss.Color("240")
	if char.Inspiration {
		inspirationIcon = "☑"
		inspirationColor = lipgloss.Color("42") // Green when active
	}
	inspirationStyle := lipgloss.NewStyle().Foreground(inspirationColor).Bold(true)
	inspirationLabel := inspirationStyle.Render(fmt.Sprintf("%s Inspiration", inspirationIcon))

	// Add note for Humans
	if char.Race == "Human" {
		inspirationLabel += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render(" (auto-restored on rest)")
	}

	lines = append(lines, inspirationLabel)

	content := strings.Join(lines, "\n")

	return content
}

// Update handles updates for the character stats panel
func (p *CharacterStatsPanel) Update(char *models.Character) {
	p.character = char
}

// EditName starts editing the character name
func (p *CharacterStatsPanel) EditName() {
	p.editMode = CharStatsEditName
	p.nameInput.SetValue(p.character.Name)
	p.nameInput.Focus()
}

// EditRace starts editing the character race
func (p *CharacterStatsPanel) EditRace() {
	p.editMode = CharStatsEditRace
	p.raceInput.SetValue(p.character.Race)
	p.raceInput.Focus()
}

// EditHP starts editing HP with a popup
func (p *CharacterStatsPanel) EditHP() {
	p.editMode = CharStatsEditHP
	p.hpInput.SetValue("")
	p.hpInput.Focus()
}

// EditXP starts editing XP with a popup
func (p *CharacterStatsPanel) EditXP() {
	p.editMode = CharStatsEditXP
	p.xpInput.SetValue("")
	p.xpInput.Focus()
}

// SaveName saves the edited name
func (p *CharacterStatsPanel) SaveName() {
	p.character.Name = p.nameInput.Value()
	p.editMode = CharStatsNormal
	p.nameInput.Blur()
}

// SaveRace saves the edited race
func (p *CharacterStatsPanel) SaveRace() {
	p.character.Race = p.raceInput.Value()
	p.editMode = CharStatsNormal
	p.raceInput.Blur()
}

// SaveHP applies HP change from input
func (p *CharacterStatsPanel) SaveHP() (int, error) {
	value := p.hpInput.Value()
	if value == "" {
		return 0, fmt.Errorf("no value entered")
	}

	// Parse the value (supports +5, -3, or just 5)
	var amount int
	_, err := fmt.Sscanf(value, "%d", &amount)
	if err != nil {
		return 0, err
	}

	// Apply HP change
	p.character.CurrentHP += amount
	if p.character.CurrentHP > p.character.MaxHP {
		p.character.CurrentHP = p.character.MaxHP
	}
	if p.character.CurrentHP < 0 {
		p.character.CurrentHP = 0
	}

	p.editMode = CharStatsNormal
	p.hpInput.Blur()
	return amount, nil
}

// SaveXP applies XP change from input
func (p *CharacterStatsPanel) SaveXP() (int, error) {
	value := p.xpInput.Value()
	if value == "" {
		return 0, fmt.Errorf("no value entered")
	}

	// Parse the value (supports +100, -50, or just 500)
	var amount int
	_, err := fmt.Sscanf(value, "%d", &amount)
	if err != nil {
		return 0, err
	}

	// Apply XP change
	p.character.Experience += amount
	if p.character.Experience < 0 {
		p.character.Experience = 0
	}

	p.editMode = CharStatsNormal
	p.xpInput.Blur()
	return amount, nil
}

// CancelEdit cancels editing
func (p *CharacterStatsPanel) CancelEdit() {
	p.editMode = CharStatsNormal
	p.nameInput.Blur()
	p.raceInput.Blur()
	p.hpInput.Blur()
	p.xpInput.Blur()
}

// AddHP adds HP to the character
func (p *CharacterStatsPanel) AddHP(amount int) {
	p.character.CurrentHP += amount
	if p.character.CurrentHP > p.character.MaxHP {
		p.character.CurrentHP = p.character.MaxHP
	}
}

// RemoveHP removes HP from the character
func (p *CharacterStatsPanel) RemoveHP(amount int) {
	p.character.CurrentHP -= amount
	if p.character.CurrentHP < 0 {
		p.character.CurrentHP = 0
	}
}

// GetInitiativeModifier returns the initiative modifier (single source of truth from character model)
func (p *CharacterStatsPanel) GetInitiativeModifier() int {
	return p.character.Initiative
}

// GetEditMode returns the current edit mode
func (p *CharacterStatsPanel) GetEditMode() CharStatsEditMode {
	return p.editMode
}

// HandleInput handles text input updates
func (p *CharacterStatsPanel) HandleInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if p.editMode == CharStatsEditName {
		p.nameInput, cmd = p.nameInput.Update(msg)
	} else if p.editMode == CharStatsEditRace {
		p.raceInput, cmd = p.raceInput.Update(msg)
	} else if p.editMode == CharStatsEditHP {
		p.hpInput, cmd = p.hpInput.Update(msg)
	} else if p.editMode == CharStatsEditXP {
		p.xpInput, cmd = p.xpInput.Update(msg)
	}
	return cmd
}

// RenderHPPopup renders the HP input popup overlay
func (p *CharacterStatsPanel) RenderHPPopup(screenWidth, screenHeight int) string {
	if p.editMode != CharStatsEditHP {
		return ""
	}

	// Create popup content
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Adjust HP"),
		"",
		"Enter amount (e.g., +5 or -3):",
		p.hpInput.View(),
		"",
		helpStyle.Render("[Enter] Apply • [Esc] Cancel"),
	)

	popup := popupStyle.Render(content)

	// Center the popup on screen using Place
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// RenderXPPopup renders the XP input popup overlay
func (p *CharacterStatsPanel) RenderXPPopup(screenWidth, screenHeight int) string {
	if p.editMode != CharStatsEditXP {
		return ""
	}

	// Create popup content
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("42")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Calculate current and next level XP
	currentXP := p.character.Experience
	nextLevelXP := getLevelXP(p.character.TotalLevel + 1)
	xpToNext := nextLevelXP - currentXP

	// Build status line
	xpStatus := labelStyle.Render("Current: ") + valueStyle.Render(fmt.Sprintf("%d XP", currentXP))
	xpNext := labelStyle.Render("Next Level: ") + infoStyle.Render(fmt.Sprintf("%d XP", nextLevelXP)) +
		labelStyle.Render(fmt.Sprintf(" (need %d more)", xpToNext))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Adjust Experience Points"),
		"",
		xpStatus,
		xpNext,
		"",
		infoStyle.Render("Enter amount (e.g., +100 or -50):"),
		p.xpInput.View(),
		"",
		helpStyle.Render("[Enter] Apply • [Esc] Cancel"),
	)

	popup := popupStyle.Render(content)

	// Center the popup on screen using Place
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// ToggleInspiration toggles the inspiration state
func (p *CharacterStatsPanel) ToggleInspiration() {
	p.character.Inspiration = !p.character.Inspiration
}

// calculatePassiveScore calculates a passive score for a given skill
// Passive score = 10 + ability modifier + proficiency bonus (if proficient)
func (p *CharacterStatsPanel) calculatePassiveScore(char *models.Character, skillName models.SkillType) int {
	// Get the skill
	skill := char.Skills.GetSkill(skillName)
	if skill == nil {
		return 10 // Default if skill not found
	}

	// Get ability modifier for the skill
	abilityMod := char.AbilityScores.GetModifier(skill.Ability)

	// Calculate skill bonus (includes proficiency if applicable)
	skillBonus := skill.CalculateBonus(abilityMod, char.ProficiencyBonus)

	// Add feat bonuses
	featBonus := 0
	switch skillName {
	case models.Perception:
		featBonus = char.PassivePerceptionBonus
	case models.Investigation:
		featBonus = char.PassiveInvestigationBonus
	case models.Insight:
		featBonus = char.PassiveInsightBonus
	}

	// Passive score = 10 + skill bonus + feat bonus
	return 10 + skillBonus + featBonus
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
