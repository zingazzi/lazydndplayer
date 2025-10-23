// internal/models/dice_roller.go
package models

import (
	"math/rand"
	"time"
)

// DiceRoller is an interface for rolling dice
type DiceRoller interface {
	Roll(sides int) int
	RollMultiple(count, sides int) []int
}

// StandardDiceRoller is the default implementation using math/rand
type StandardDiceRoller struct {
	rng *rand.Rand
}

// NewStandardDiceRoller creates a new dice roller with a time-based seed
func NewStandardDiceRoller() *StandardDiceRoller {
	return &StandardDiceRoller{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewSeededDiceRoller creates a new dice roller with a specific seed (for testing)
func NewSeededDiceRoller(seed int64) *StandardDiceRoller {
	return &StandardDiceRoller{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Roll rolls a single die with the specified number of sides
func (dr *StandardDiceRoller) Roll(sides int) int {
	if sides <= 0 {
		return 0
	}
	return dr.rng.Intn(sides) + 1
}

// RollMultiple rolls multiple dice and returns all results
func (dr *StandardDiceRoller) RollMultiple(count, sides int) []int {
	results := make([]int, count)
	for i := 0; i < count; i++ {
		results[i] = dr.Roll(sides)
	}
	return results
}

// defaultDiceRoller is a package-level instance for backward compatibility
var defaultDiceRoller DiceRoller = NewStandardDiceRoller()

// SetDefaultDiceRoller allows tests to inject a mock dice roller
func SetDefaultDiceRoller(roller DiceRoller) {
	defaultDiceRoller = roller
}

// GetDefaultDiceRoller returns the current default dice roller
func GetDefaultDiceRoller() DiceRoller {
	return defaultDiceRoller
}

