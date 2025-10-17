// internal/dice/roller.go
package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RollType represents how to roll (normal, advantage, disadvantage)
type RollType string

const (
	Normal       RollType = "normal"
	Advantage    RollType = "advantage"
	Disadvantage RollType = "disadvantage"
)

// RollResult represents the result of a dice roll
type RollResult struct {
	Expression string   `json:"expression"`
	Rolls      []int    `json:"rolls"`
	Modifier   int      `json:"modifier"`
	Total      int      `json:"total"`
	RollType   RollType `json:"roll_type"`
	Timestamp  time.Time `json:"timestamp"`
}

// String returns a formatted string of the roll result
func (r *RollResult) String() string {
	rollsStr := make([]string, len(r.Rolls))
	for i, roll := range r.Rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	typeStr := ""
	if r.RollType == Advantage {
		typeStr = " (advantage)"
	} else if r.RollType == Disadvantage {
		typeStr = " (disadvantage)"
	}

	modStr := ""
	if r.Modifier > 0 {
		modStr = fmt.Sprintf(" + %d", r.Modifier)
	} else if r.Modifier < 0 {
		modStr = fmt.Sprintf(" - %d", -r.Modifier)
	}

	return fmt.Sprintf("%s%s: [%s]%s = %d", r.Expression, typeStr, strings.Join(rollsStr, ", "), modStr, r.Total)
}

// Roll rolls dice based on expression (e.g., "2d6+3", "1d20")
func Roll(expression string, rollType RollType) (*RollResult, error) {
	// Parse the expression
	dice, sides, modifier, err := parseExpression(expression)
	if err != nil {
		return nil, err
	}

	// Roll the dice
	rolls := rollDice(dice, sides, rollType)

	// Calculate total
	total := sum(rolls) + modifier

	return &RollResult{
		Expression: expression,
		Rolls:      rolls,
		Modifier:   modifier,
		Total:      total,
		RollType:   rollType,
		Timestamp:  time.Now(),
	}, nil
}

// parseExpression parses dice notation (e.g., "2d6+3")
func parseExpression(expr string) (dice, sides, modifier int, err error) {
	expr = strings.ToLower(strings.ReplaceAll(expr, " ", ""))

	// Match patterns like 2d6+3, 1d20, d20-2, etc.
	re := regexp.MustCompile(`^(\d*)d(\d+)([\+\-]\d+)?$`)
	matches := re.FindStringSubmatch(expr)

	if matches == nil {
		return 0, 0, 0, fmt.Errorf("invalid dice expression: %s", expr)
	}

	// Number of dice (default 1 if not specified)
	if matches[1] == "" {
		dice = 1
	} else {
		dice, err = strconv.Atoi(matches[1])
		if err != nil {
			return 0, 0, 0, err
		}
	}

	// Sides
	sides, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, 0, err
	}

	// Modifier (optional)
	if matches[3] != "" {
		modifier, err = strconv.Atoi(matches[3])
		if err != nil {
			return 0, 0, 0, err
		}
	}

	return dice, sides, modifier, nil
}

// rollDice performs the actual dice rolling
func rollDice(count, sides int, rollType RollType) []int {
	if rollType == Advantage || rollType == Disadvantage {
		// For advantage/disadvantage, roll twice and take best/worst
		roll1 := rollSingleDie(sides)
		roll2 := rollSingleDie(sides)

		if rollType == Advantage {
			if roll1 > roll2 {
				return []int{roll1, roll2} // Keep higher
			}
			return []int{roll2, roll1}
		} else {
			if roll1 < roll2 {
				return []int{roll1, roll2} // Keep lower
			}
			return []int{roll2, roll1}
		}
	}

	// Normal roll
	rolls := make([]int, count)
	for i := 0; i < count; i++ {
		rolls[i] = rollSingleDie(sides)
	}
	return rolls
}

// rollSingleDie rolls a single die
func rollSingleDie(sides int) int {
	return rand.Intn(sides) + 1
}

// sum calculates the sum of rolls (for advantage/disadvantage, only counts first roll)
func sum(rolls []int) int {
	if len(rolls) == 0 {
		return 0
	}
	// For advantage/disadvantage with 2 rolls, only count the first (best/worst)
	if len(rolls) == 2 {
		return rolls[0]
	}
	// For normal rolls, sum all
	total := 0
	for _, roll := range rolls {
		total += roll
	}
	return total
}

// RollHistory maintains a history of recent rolls
type RollHistory struct {
	Rolls    []RollResult
	MaxSize  int
}

// NewRollHistory creates a new roll history
func NewRollHistory(maxSize int) *RollHistory {
	return &RollHistory{
		Rolls:   []RollResult{},
		MaxSize: maxSize,
	}
}

// Add adds a roll to the history
func (h *RollHistory) Add(result RollResult) {
	h.Rolls = append([]RollResult{result}, h.Rolls...)
	if len(h.Rolls) > h.MaxSize {
		h.Rolls = h.Rolls[:h.MaxSize]
	}
}

// Clear clears the history
func (h *RollHistory) Clear() {
	h.Rolls = []RollResult{}
}

// GetRecent returns the most recent n rolls
func (h *RollHistory) GetRecent(n int) []RollResult {
	if n > len(h.Rolls) {
		n = len(h.Rolls)
	}
	return h.Rolls[:n]
}
