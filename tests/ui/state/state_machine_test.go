// tests/ui/state/state_machine_test.go
package state_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/ui/state"
)

func TestNewStateMachine(t *testing.T) {
	sm := state.NewStateMachine()

	if sm == nil {
		t.Fatal("NewStateMachine() returned nil")
	}

	if sm.GetState() != state.StateIdle {
		t.Errorf("Expected initial state to be StateIdle, got %v", sm.GetState())
	}

	if !sm.IsIdle() {
		t.Error("Expected state machine to be idle initially")
	}
}

func TestTransition(t *testing.T) {
	sm := state.NewStateMachine()

	// Test transition to a new state
	sm.Transition(state.StateFeatSelection, nil)

	if sm.GetState() != state.StateFeatSelection {
		t.Errorf("Expected state to be StateFeatSelection, got %v", sm.GetState())
	}

	if sm.IsIdle() {
		t.Error("State machine should not be idle after transition")
	}
}

func TestTransitionWithContext(t *testing.T) {
	sm := state.NewStateMachine()

	context := map[string]interface{}{
		"pendingFeat": "Test Feat",
		"count":       42,
	}

	sm.Transition(state.StateFeatSelection, context)

	if sm.GetState() != state.StateFeatSelection {
		t.Errorf("Expected state to be StateFeatSelection, got %v", sm.GetState())
	}

	// Check context values
	feat := sm.GetContext("pendingFeat")
	if feat != "Test Feat" {
		t.Errorf("Expected context 'pendingFeat' to be 'Test Feat', got %v", feat)
	}

	count := sm.GetContext("count")
	if count != 42 {
		t.Errorf("Expected context 'count' to be 42, got %v", count)
	}
}

func TestSetContext(t *testing.T) {
	sm := state.NewStateMachine()

	sm.SetContext("key1", "value1")
	sm.SetContext("key2", 123)

	if sm.GetContext("key1") != "value1" {
		t.Errorf("Expected context 'key1' to be 'value1', got %v", sm.GetContext("key1"))
	}

	if sm.GetContext("key2") != 123 {
		t.Errorf("Expected context 'key2' to be 123, got %v", sm.GetContext("key2"))
	}
}

func TestClearContext(t *testing.T) {
	sm := state.NewStateMachine()

	sm.SetContext("key1", "value1")
	sm.SetContext("key2", "value2")

	sm.ClearContext("key1")

	if sm.GetContext("key1") != nil {
		t.Error("Expected context 'key1' to be cleared")
	}

	if sm.GetContext("key2") != "value2" {
		t.Error("Expected context 'key2' to still exist")
	}
}

func TestClear(t *testing.T) {
	sm := state.NewStateMachine()

	sm.Transition(state.StateFeatSelection, map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	sm.Clear()

	if sm.GetState() != state.StateIdle {
		t.Errorf("Expected state to be StateIdle after Clear(), got %v", sm.GetState())
	}

	if !sm.IsIdle() {
		t.Error("Expected state machine to be idle after Clear()")
	}

	if sm.GetContext("key1") != nil {
		t.Error("Expected all context to be cleared")
	}
}

func TestMultipleTransitions(t *testing.T) {
	sm := state.NewStateMachine()

	// Test state transitions
	states := []state.StateType{
		state.StateSpeciesSelection,
		state.StateClassSelection,
		state.StateSubclassSelection,
		state.StateLevelUp,
		state.StateFeatSelection,
	}

	for _, expectedState := range states {
		sm.Transition(expectedState, nil)
		if sm.GetState() != expectedState {
			t.Errorf("Expected state %v, got %v", expectedState, sm.GetState())
		}
	}
}

func TestContextOverwrite(t *testing.T) {
	sm := state.NewStateMachine()

	sm.SetContext("key", "value1")
	sm.SetContext("key", "value2")

	if sm.GetContext("key") != "value2" {
		t.Errorf("Expected context 'key' to be overwritten to 'value2', got %v", sm.GetContext("key"))
	}
}

func TestTransitionOverwritesContext(t *testing.T) {
	sm := state.NewStateMachine()

	sm.SetContext("key1", "old1")
	sm.SetContext("key2", "old2")

	sm.Transition(state.StateFeatSelection, map[string]interface{}{
		"key1": "new1",
		"key3": "new3",
	})

	if sm.GetContext("key1") != "new1" {
		t.Error("Expected key1 to be overwritten")
	}

	if sm.GetContext("key2") != "old2" {
		t.Error("Expected key2 to remain unchanged")
	}

	if sm.GetContext("key3") != "new3" {
		t.Error("Expected key3 to be added")
	}
}
