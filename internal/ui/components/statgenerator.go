// internal/ui/components/statgenerator.go
package components

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// StatGenMethod represents the method used to generate stats
type StatGenMethod int

const (
	Method4d6DropLowest StatGenMethod = iota
	MethodStandardArray
	MethodPointBuy
	MethodCustomValues
)

// StatGenState represents the current state of the stat generator
type StatGenState int

const (
	StateSelectMethod StatGenState = iota
	StateAssignStats
	StateSetExtras
)

// StatGenerator is a component for generating ability scores
type StatGenerator struct {
	visible       bool
	method        StatGenMethod
	state         StatGenState
	selectedIndex int

	// Generated or selected stats (not yet assigned)
	availableStats []int
	rollDetails    []string // Details of rolls for 4d6 method

	// Assignment tracking
	assignments map[models.AbilityType]int // ability -> index in availableStats (-1 if not assigned)

	// Point buy specific
	pointBuyScores map[models.AbilityType]int
	pointsSpent    int
	maxPoints      int

	// Extras editing
	extraValues map[models.AbilityType]int
	editingExtra bool
	extraInput   string

	// Track if user went directly to extras (e) or through full flow (r)
	directToExtras bool

	// Original character data (to restore on cancel)
	originalScores *models.AbilityScores
}

// Standard array values
var standardArray = []int{15, 14, 13, 12, 10, 8}

// Point buy costs (score -> cost)
var pointBuyCosts = map[int]int{
	8:  0,
	9:  1,
	10: 2,
	11: 3,
	12: 4,
	13: 5,
	14: 7,
	15: 9,
}

// NewStatGenerator creates a new stat generator
func NewStatGenerator() *StatGenerator {
	return &StatGenerator{
		visible:        false,
		method:         Method4d6DropLowest,
		state:          StateSelectMethod,
		selectedIndex:  0,
		availableStats: []int{},
		assignments:    make(map[models.AbilityType]int),
		pointBuyScores: make(map[models.AbilityType]int),
		extraValues:    make(map[models.AbilityType]int),
		maxPoints:      27,
	}
}

// Show displays the stat generator
func (s *StatGenerator) Show(currentScores *models.AbilityScores) {
	s.visible = true
	s.state = StateSelectMethod
	s.selectedIndex = 0
	s.availableStats = []int{}
	s.rollDetails = []string{}
	s.assignments = make(map[models.AbilityType]int)
	s.pointBuyScores = make(map[models.AbilityType]int)
	s.extraValues = make(map[models.AbilityType]int)
	s.pointsSpent = 0
	s.editingExtra = false
	s.extraInput = ""
	s.directToExtras = false

	// Save original scores for cancel
	s.originalScores = &models.AbilityScores{
		Strength:          currentScores.Strength,
		Dexterity:         currentScores.Dexterity,
		Constitution:      currentScores.Constitution,
		Intelligence:      currentScores.Intelligence,
		Wisdom:            currentScores.Wisdom,
		Charisma:          currentScores.Charisma,
		StrengthBase:      currentScores.StrengthBase,
		DexterityBase:     currentScores.DexterityBase,
		ConstitutionBase:  currentScores.ConstitutionBase,
		IntelligenceBase:  currentScores.IntelligenceBase,
		WisdomBase:        currentScores.WisdomBase,
		CharismaBase:      currentScores.CharismaBase,
		StrengthExtra:     currentScores.StrengthExtra,
		DexterityExtra:    currentScores.DexterityExtra,
		ConstitutionExtra: currentScores.ConstitutionExtra,
		IntelligenceExtra: currentScores.IntelligenceExtra,
		WisdomExtra:       currentScores.WisdomExtra,
		CharismaExtra:     currentScores.CharismaExtra,
	}

	// Load current extras
	abilities := []models.AbilityType{
		models.Strength, models.Dexterity, models.Constitution,
		models.Intelligence, models.Wisdom, models.Charisma,
	}
	for _, ability := range abilities {
		s.extraValues[ability] = currentScores.GetExtraScore(ability)
	}
}

// ShowExtrasOnly displays only the extras editing screen
func (s *StatGenerator) ShowExtrasOnly(currentScores *models.AbilityScores) {
	s.Show(currentScores)
	s.state = StateSetExtras
	s.selectedIndex = 0
	s.directToExtras = true // Mark that we came directly to extras
}

// Hide hides the stat generator
func (s *StatGenerator) Hide() {
	s.visible = false
}

// IsVisible returns whether the stat generator is visible
func (s *StatGenerator) IsVisible() bool {
	return s.visible
}

// IsEditingExtra returns whether we're currently editing an extra value
func (s *StatGenerator) IsEditingExtra() bool {
	return s.editingExtra
}

// GetMethod returns the current generation method
func (s *StatGenerator) GetMethod() StatGenMethod {
	return s.method
}

// GetState returns the current state
func (s *StatGenerator) GetState() StatGenState {
	return s.state
}

// Next moves selection down
func (s *StatGenerator) Next() {
	switch s.state {
	case StateSelectMethod:
		s.selectedIndex = (s.selectedIndex + 1) % 4 // 4 methods now
	case StateAssignStats:
		s.selectedIndex = (s.selectedIndex + 1) % 6
	case StateSetExtras:
		if !s.editingExtra {
			s.selectedIndex = (s.selectedIndex + 1) % 7 // 6 abilities + confirm button
		}
	}
}

// Prev moves selection up
func (s *StatGenerator) Prev() {
	switch s.state {
	case StateSelectMethod:
		s.selectedIndex--
		if s.selectedIndex < 0 {
			s.selectedIndex = 3 // 4 methods now
		}
	case StateAssignStats:
		s.selectedIndex--
		if s.selectedIndex < 0 {
			s.selectedIndex = 5
		}
	case StateSetExtras:
		if !s.editingExtra {
			s.selectedIndex--
			if s.selectedIndex < 0 {
				s.selectedIndex = 6
			}
		}
	}
}

// SelectMethod selects the stat generation method
func (s *StatGenerator) SelectMethod() {
	s.method = StatGenMethod(s.selectedIndex)

	switch s.method {
	case Method4d6DropLowest:
		// Roll stats
		s.availableStats, s.rollDetails = s.roll4d6DropLowest()
		s.state = StateAssignStats
		s.selectedIndex = 0
		// Initialize assignments
		abilities := []models.AbilityType{
			models.Strength, models.Dexterity, models.Constitution,
			models.Intelligence, models.Wisdom, models.Charisma,
		}
		for _, ability := range abilities {
			s.assignments[ability] = -1
		}

	case MethodStandardArray:
		// Use standard array
		s.availableStats = make([]int, len(standardArray))
		copy(s.availableStats, standardArray)
		s.state = StateAssignStats
		s.selectedIndex = 0
		// Initialize assignments
		abilities := []models.AbilityType{
			models.Strength, models.Dexterity, models.Constitution,
			models.Intelligence, models.Wisdom, models.Charisma,
		}
		for _, ability := range abilities {
			s.assignments[ability] = -1
		}

	case MethodPointBuy:
		// Initialize all stats to 8
		abilities := []models.AbilityType{
			models.Strength, models.Dexterity, models.Constitution,
			models.Intelligence, models.Wisdom, models.Charisma,
		}
		for _, ability := range abilities {
			s.pointBuyScores[ability] = 8
		}
		s.pointsSpent = 0
		s.state = StateAssignStats
		s.selectedIndex = 0

	case MethodCustomValues:
		// Initialize all stats to 10 (default)
		abilities := []models.AbilityType{
			models.Strength, models.Dexterity, models.Constitution,
			models.Intelligence, models.Wisdom, models.Charisma,
		}
		for _, ability := range abilities {
			s.pointBuyScores[ability] = 10 // Reuse pointBuyScores for custom
		}
		s.state = StateAssignStats
		s.selectedIndex = 0
	}
}

// roll4d6DropLowest rolls 4d6 and drops the lowest, repeated 6 times
func (s *StatGenerator) roll4d6DropLowest() ([]int, []string) {
	stats := make([]int, 6)
	details := make([]string, 6)

	for i := 0; i < 6; i++ {
		rolls := make([]int, 4)
		for j := 0; j < 4; j++ {
			rolls[j] = rand.Intn(6) + 1
		}
		sort.Ints(rolls)

		// Drop the lowest (first in sorted array)
		sum := rolls[1] + rolls[2] + rolls[3]
		stats[i] = sum

		details[i] = fmt.Sprintf("[%d, %d, %d, %d] drop %d = %d",
			rolls[0], rolls[1], rolls[2], rolls[3], rolls[0], sum)
	}

	return stats, details
}

// GetAbilityOrder returns abilities in display order
func (s *StatGenerator) GetAbilityOrder() []models.AbilityType {
	return []models.AbilityType{
		models.Strength,
		models.Dexterity,
		models.Constitution,
		models.Intelligence,
		models.Wisdom,
		models.Charisma,
	}
}

// ToggleAssignment assigns/unassigns a stat to the selected ability
func (s *StatGenerator) ToggleAssignment(statIndex int) {
	abilities := s.GetAbilityOrder()
	selectedAbility := abilities[s.selectedIndex]

	// Check if this stat is already assigned to another ability
	for ability, idx := range s.assignments {
		if idx == statIndex && ability != selectedAbility {
			// Unassign it first
			s.assignments[ability] = -1
		}
	}

	// If current ability has this stat, unassign it
	if s.assignments[selectedAbility] == statIndex {
		s.assignments[selectedAbility] = -1
	} else {
		// Assign the stat
		s.assignments[selectedAbility] = statIndex
	}
}

// IncreasePointBuy increases the selected ability score (point buy)
func (s *StatGenerator) IncreasePointBuy() {
	abilities := s.GetAbilityOrder()
	ability := abilities[s.selectedIndex]
	current := s.pointBuyScores[ability]

	if s.method == MethodCustomValues {
		// For custom values, allow any value up to 20
		if current >= 20 {
			return
		}
		s.pointBuyScores[ability] = current + 1
		return
	}

	// Point buy mode
	if current >= 15 {
		return // Max is 15
	}

	cost := pointBuyCosts[current+1] - pointBuyCosts[current]
	if s.pointsSpent+cost > s.maxPoints {
		return // Not enough points
	}

	s.pointBuyScores[ability] = current + 1
	s.pointsSpent += cost
}

// DecreasePointBuy decreases the selected ability score (point buy)
func (s *StatGenerator) DecreasePointBuy() {
	abilities := s.GetAbilityOrder()
	ability := abilities[s.selectedIndex]
	current := s.pointBuyScores[ability]

	if s.method == MethodCustomValues {
		// For custom values, allow any value down to 1
		if current <= 1 {
			return
		}
		s.pointBuyScores[ability] = current - 1
		return
	}

	// Point buy mode
	if current <= 8 {
		return // Min is 8
	}

	cost := pointBuyCosts[current] - pointBuyCosts[current-1]
	s.pointBuyScores[ability] = current - 1
	s.pointsSpent -= cost
}

// AllStatsAssigned checks if all stats have been assigned
func (s *StatGenerator) AllStatsAssigned() bool {
	for _, idx := range s.assignments {
		if idx == -1 {
			return false
		}
	}
	return true
}

// ApplyToCharacter applies the generated stats to the character
func (s *StatGenerator) ApplyToCharacter(char *models.Character) {
	abilities := s.GetAbilityOrder()

	switch s.method {
	case Method4d6DropLowest, MethodStandardArray:
		// Apply assignments
		for _, ability := range abilities {
			idx := s.assignments[ability]
			if idx >= 0 && idx < len(s.availableStats) {
				char.AbilityScores.SetBaseScore(ability, s.availableStats[idx])
			}
		}

	case MethodPointBuy, MethodCustomValues:
		// Apply point buy or custom scores
		for _, ability := range abilities {
			char.AbilityScores.SetBaseScore(ability, s.pointBuyScores[ability])
		}
	}

	// Apply extras
	for _, ability := range abilities {
		char.AbilityScores.SetExtraScore(ability, s.extraValues[ability])
	}

	char.AbilityScores.RecalculateTotals()
	char.UpdateDerivedStats()
}

// StartEditingExtra starts editing the extra value for selected ability
func (s *StatGenerator) StartEditingExtra() {
	if s.selectedIndex < 6 {
		abilities := s.GetAbilityOrder()
		ability := abilities[s.selectedIndex]
		s.editingExtra = true
		s.extraInput = fmt.Sprintf("%d", s.extraValues[ability])
	}
}

// IncreaseExtra increases the extra value for selected ability by 1
func (s *StatGenerator) IncreaseExtra() {
	if s.selectedIndex < 6 {
		abilities := s.GetAbilityOrder()
		ability := abilities[s.selectedIndex]
		if s.extraValues[ability] < 10 { // Max +10
			s.extraValues[ability]++
		}
	}
}

// DecreaseExtra decreases the extra value for selected ability by 1
func (s *StatGenerator) DecreaseExtra() {
	if s.selectedIndex < 6 {
		abilities := s.GetAbilityOrder()
		ability := abilities[s.selectedIndex]
		if s.extraValues[ability] > -5 { // Min -5
			s.extraValues[ability]--
		}
	}
}

// HandleExtraInput handles character input for extra editing
func (s *StatGenerator) HandleExtraInput(char rune) {
	s.extraInput += string(char)
}

// DeleteExtraInput deletes last character from extra input
func (s *StatGenerator) DeleteExtraInput() {
	if len(s.extraInput) > 0 {
		s.extraInput = s.extraInput[:len(s.extraInput)-1]
	}
}

// SaveExtra saves the edited extra value
func (s *StatGenerator) SaveExtra() {
	if s.selectedIndex < 6 {
		abilities := s.GetAbilityOrder()
		ability := abilities[s.selectedIndex]

		var value int
		fmt.Sscanf(s.extraInput, "%d", &value)

		// Clamp to reasonable range
		if value < -5 {
			value = -5
		}
		if value > 10 {
			value = 10
		}

		s.extraValues[ability] = value
	}
	s.editingExtra = false
	s.extraInput = ""
}

// CancelExtra cancels extra editing
func (s *StatGenerator) CancelExtra() {
	s.editingExtra = false
	s.extraInput = ""
}

// View renders the stat generator
func (s *StatGenerator) View(width, height int) string {
	if !s.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	unselectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(width - 20).
		Height(height - 6)

	var content string

	switch s.state {
	case StateSelectMethod:
		content = s.renderMethodSelection(titleStyle, selectedStyle, unselectedStyle, instructionStyle)
	case StateAssignStats:
		content = s.renderStatsAssignment(titleStyle, selectedStyle, unselectedStyle, instructionStyle)
	case StateSetExtras:
		content = s.renderExtrasEditing(titleStyle, selectedStyle, unselectedStyle, instructionStyle)
	}

	box := boxStyle.Render(content)

	// Center the box
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

func (s *StatGenerator) renderMethodSelection(titleStyle, selectedStyle, unselectedStyle, instructionStyle lipgloss.Style) string {
	var lines []string
	lines = append(lines, titleStyle.Render("SELECT STAT GENERATION METHOD"))
	lines = append(lines, "")

	methods := []string{
		"4d6 Drop Lowest - Roll four six-sided dice, drop the lowest",
		"Standard Array - Use pre-determined set [15, 14, 13, 12, 10, 8]",
		"Point Buy - Start with 8 in each stat, spend 27 points to increase",
		"Custom Values - Manually enter any values you want",
	}

	for i, method := range methods {
		if i == s.selectedIndex {
			lines = append(lines, selectedStyle.Render(fmt.Sprintf("▶ %s", method)))
		} else {
			lines = append(lines, unselectedStyle.Render(fmt.Sprintf("  %s", method)))
		}
	}

	lines = append(lines, "")
	lines = append(lines, instructionStyle.Render("↑/↓: Navigate  Enter: Select  Esc: Cancel"))

	return strings.Join(lines, "\n")
}

func (s *StatGenerator) renderStatsAssignment(titleStyle, selectedStyle, unselectedStyle, instructionStyle lipgloss.Style) string {
	var lines []string

	// Style for used stats
	usedStatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).  // Red
		Strikethrough(true)

	switch s.method {
	case Method4d6DropLowest:
		lines = append(lines, titleStyle.Render("4D6 DROP LOWEST - ASSIGN STATS"))
		lines = append(lines, "")
		lines = append(lines, "Roll Results:")

		// Create a map to track which stats are used
		usedStats := make(map[int]bool)
		for _, idx := range s.assignments {
			if idx >= 0 {
				usedStats[idx] = true
			}
		}

		for i, detail := range s.rollDetails {
			score := s.availableStats[i]
			modifier := models.CalculateModifier(score)
			line := fmt.Sprintf("  %d: %s (mod: %+d)", i+1, detail, modifier)

			// Mark used stats in red
			if usedStats[i] {
				lines = append(lines, usedStatStyle.Render(line)+" "+lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("(used)"))
			} else {
				lines = append(lines, line)
			}
		}
		lines = append(lines, "")
		lines = append(lines, s.renderAssignments(selectedStyle, unselectedStyle))
		lines = append(lines, "")
		lines = append(lines, instructionStyle.Render("↑/↓: Select Ability  1-6: Assign Stat  Enter: Continue  Esc: Back"))

	case MethodStandardArray:
		lines = append(lines, titleStyle.Render("STANDARD ARRAY - ASSIGN STATS"))
		lines = append(lines, "")
		lines = append(lines, "Available values (with modifiers):")

		// Create a map to track which stats are used
		usedStats := make(map[int]bool)
		for _, idx := range s.assignments {
			if idx >= 0 {
				usedStats[idx] = true
			}
		}

		for i, score := range s.availableStats {
			modifier := models.CalculateModifier(score)
			line := fmt.Sprintf("  %d: %d (%+d)", i+1, score, modifier)

			// Mark used stats in red
			if usedStats[i] {
				lines = append(lines, usedStatStyle.Render(line)+" "+lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("(used)"))
			} else {
				lines = append(lines, line)
			}
		}
		lines = append(lines, "")
		lines = append(lines, s.renderAssignments(selectedStyle, unselectedStyle))
		lines = append(lines, "")
		lines = append(lines, instructionStyle.Render("↑/↓: Select Ability  1-6: Assign Stat  Enter: Continue  Esc: Back"))

	case MethodPointBuy:
		lines = append(lines, titleStyle.Render("POINT BUY - ADJUST STATS"))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Points Remaining: %d / %d", s.maxPoints-s.pointsSpent, s.maxPoints))
		lines = append(lines, "Range: 8 (min) to 15 (max)")
		lines = append(lines, "")
		lines = append(lines, s.renderPointBuy(selectedStyle, unselectedStyle))
		lines = append(lines, "")
		lines = append(lines, instructionStyle.Render("↑/↓: Select Ability  +/-: Increase/Decrease  Enter: Continue  Esc: Back"))

	case MethodCustomValues:
		lines = append(lines, titleStyle.Render("CUSTOM VALUES - SET YOUR STATS"))
		lines = append(lines, "")
		lines = append(lines, "Set any values you want (min: 1, max: 20)")
		lines = append(lines, "")
		lines = append(lines, s.renderPointBuy(selectedStyle, unselectedStyle)) // Reuse same rendering
		lines = append(lines, "")
		lines = append(lines, instructionStyle.Render("↑/↓: Select Ability  +/-: Increase/Decrease  Enter: Continue  Esc: Back"))
	}

	return strings.Join(lines, "\n")
}

func (s *StatGenerator) renderAssignments(selectedStyle, unselectedStyle lipgloss.Style) string {
	var lines []string
	abilities := s.GetAbilityOrder()
	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength    ",
		models.Dexterity:    "Dexterity   ",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom      ",
		models.Charisma:     "Charisma    ",
	}

	for i, ability := range abilities {
		idx := s.assignments[ability]
		value := "---    "
		if idx >= 0 && idx < len(s.availableStats) {
			score := s.availableStats[idx]
			modifier := models.CalculateModifier(score)
			value = fmt.Sprintf("%2d (%+d)", score, modifier)
		}

		line := fmt.Sprintf("%s: %s", abilityNames[ability], value)
		if i == s.selectedIndex {
			lines = append(lines, selectedStyle.Render(fmt.Sprintf("▶ %s", line)))
		} else {
			lines = append(lines, unselectedStyle.Render(fmt.Sprintf("  %s", line)))
		}
	}

	return strings.Join(lines, "\n")
}

func (s *StatGenerator) renderPointBuy(selectedStyle, unselectedStyle lipgloss.Style) string {
	var lines []string
	abilities := s.GetAbilityOrder()
	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength    ",
		models.Dexterity:    "Dexterity   ",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom      ",
		models.Charisma:     "Charisma    ",
	}

	for i, ability := range abilities {
		score := s.pointBuyScores[ability]
		modifier := models.CalculateModifier(score)

		var line string
		if s.method == MethodCustomValues {
			line = fmt.Sprintf("%s: %2d (mod: %+d)", abilityNames[ability], score, modifier)
		} else {
			// Point buy - show cost
			cost := pointBuyCosts[score]
			line = fmt.Sprintf("%s: %2d (mod: %+d, cost: %2d)", abilityNames[ability], score, modifier, cost)
		}

		if i == s.selectedIndex {
			lines = append(lines, selectedStyle.Render(fmt.Sprintf("▶ %s", line)))
		} else {
			lines = append(lines, unselectedStyle.Render(fmt.Sprintf("  %s", line)))
		}
	}

	return strings.Join(lines, "\n")
}

func (s *StatGenerator) renderExtrasEditing(titleStyle, selectedStyle, unselectedStyle, instructionStyle lipgloss.Style) string {
	var lines []string
	lines = append(lines, titleStyle.Render("SET EXTRA BONUSES"))
	lines = append(lines, "")
	lines = append(lines, "Format: Base + Extra = Total (Modifier)")
	lines = append(lines, "")

	abilities := s.GetAbilityOrder()
	abilityNames := map[models.AbilityType]string{
		models.Strength:     "Strength    ",
		models.Dexterity:    "Dexterity   ",
		models.Constitution: "Constitution",
		models.Intelligence: "Intelligence",
		models.Wisdom:       "Wisdom      ",
		models.Charisma:     "Charisma    ",
	}

	for i, ability := range abilities {
		extra := s.extraValues[ability]

		// Get base score from current assignment/generation
		var baseScore int

		// If we came directly to extras, use original base scores
		if s.directToExtras && s.originalScores != nil {
			baseScore = s.originalScores.GetBaseScore(ability)
			if baseScore == 0 {
				// Fallback to total score if base not set
				baseScore = s.originalScores.GetScore(ability)
			}
		} else {
			// Otherwise use scores from current generation method
			switch s.method {
			case Method4d6DropLowest, MethodStandardArray:
				// Use assigned value
				idx := s.assignments[ability]
				if idx >= 0 && idx < len(s.availableStats) {
					baseScore = s.availableStats[idx]
				} else {
					baseScore = 10 // Default if not assigned
				}
			case MethodPointBuy, MethodCustomValues:
				// Use point buy value
				baseScore = s.pointBuyScores[ability]
			default:
				// Use original score if no method selected yet
				if s.originalScores != nil {
					baseScore = s.originalScores.GetBaseScore(ability)
					if baseScore == 0 {
						baseScore = s.originalScores.GetScore(ability)
					}
				} else {
					baseScore = 10
				}
			}
		}

		total := baseScore + extra
		modifier := models.CalculateModifier(total)

		var line string
		if s.editingExtra && i == s.selectedIndex {
			// Show input cursor
			line = fmt.Sprintf("%s: %2d + %s█ = ? (mod: ?)", abilityNames[ability], baseScore, s.extraInput)
		} else {
			// Show full calculation
			line = fmt.Sprintf("%s: %2d %+2d = %2d (mod: %+d)", abilityNames[ability], baseScore, extra, total, modifier)
		}

		if i == s.selectedIndex {
			lines = append(lines, selectedStyle.Render(fmt.Sprintf("▶ %s", line)))
		} else {
			lines = append(lines, unselectedStyle.Render(fmt.Sprintf("  %s", line)))
		}
	}

	lines = append(lines, "")

	// Confirm button
	if s.selectedIndex == 6 {
		lines = append(lines, selectedStyle.Render("▶ [CONFIRM]"))
	} else {
		lines = append(lines, unselectedStyle.Render("  [CONFIRM]"))
	}

	lines = append(lines, "")
	if s.editingExtra {
		lines = append(lines, instructionStyle.Render("Type number  Enter: Save  Esc: Cancel"))
	} else {
		var escAction string
		if s.directToExtras {
			escAction = "Back to Stats Panel"
		} else {
			escAction = "Back to Assignment"
		}
		lines = append(lines, instructionStyle.Render(fmt.Sprintf("↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: %s", escAction)))
	}

	return strings.Join(lines, "\n")
}

// CanContinue checks if user can proceed to next step
func (s *StatGenerator) CanContinue() bool {
	switch s.state {
	case StateSelectMethod:
		return true
	case StateAssignStats:
		if s.method == MethodPointBuy || s.method == MethodCustomValues {
			return true // Can always continue with point buy or custom
		}
		return s.AllStatsAssigned()
	case StateSetExtras:
		return true
	default:
		return false
	}
}

// Continue moves to the next state
func (s *StatGenerator) Continue() {
	switch s.state {
	case StateSelectMethod:
		s.SelectMethod()
	case StateAssignStats:
		s.state = StateSetExtras
		s.selectedIndex = 0
	case StateSetExtras:
		// Check if we're on the confirm button (index 6)
		if s.selectedIndex == 6 {
			// Close the generator (caller will apply stats)
			s.Hide()
		}
	}
}

// GoBack goes to the previous state
func (s *StatGenerator) GoBack() {
	switch s.state {
	case StateSelectMethod:
		s.Hide()
	case StateAssignStats:
		s.state = StateSelectMethod
		s.selectedIndex = int(s.method)
	case StateSetExtras:
		// If we came directly to extras (pressed 'e'), close the generator
		if s.directToExtras {
			s.Hide()
		} else {
			// Otherwise, go back to assignment screen (pressed 'r')
			s.state = StateAssignStats
			s.selectedIndex = 0
		}
	}
}

// Cancel cancels and restores original scores
func (s *StatGenerator) Cancel(char *models.Character) {
	if s.originalScores != nil {
		char.AbilityScores = *s.originalScores
		char.UpdateDerivedStats()
	}
	s.Hide()
}
