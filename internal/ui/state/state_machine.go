// internal/ui/state/state_machine.go
package state

// StateType represents the current state of the application flow
type StateType int

const (
	StateIdle StateType = iota
	StateSpeciesSelection
	StateSubtypeSelection
	StateClassSelection
	StateSubclassSelection
	StateLevelUp
	StateFeatSelection
	StateAbilityChoice
	StateDivineOrderSelection
	StateSpellSelection
	StateEldritchKnightSpellSelection
	StateCharacterCreationWizard
)

// StateMachine manages complex selection flows and temporary state
type StateMachine struct {
	currentState StateType
	context      map[string]interface{}
}

// NewStateMachine creates a new state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		currentState: StateIdle,
		context:      make(map[string]interface{}),
	}
}

// Transition transitions to a new state with optional context
func (sm *StateMachine) Transition(newState StateType, context map[string]interface{}) {
	sm.currentState = newState
	for k, v := range context {
		sm.context[k] = v
	}
}

// GetState returns the current state
func (sm *StateMachine) GetState() StateType {
	return sm.currentState
}

// GetContext retrieves a value from the context
func (sm *StateMachine) GetContext(key string) interface{} {
	return sm.context[key]
}

// SetContext sets a value in the context
func (sm *StateMachine) SetContext(key string, value interface{}) {
	sm.context[key] = value
}

// ClearContext clears a specific context key
func (sm *StateMachine) ClearContext(key string) {
	delete(sm.context, key)
}

// Clear clears all state and context
func (sm *StateMachine) Clear() {
	sm.currentState = StateIdle
	sm.context = make(map[string]interface{})
}

// IsIdle returns true if the state machine is in idle state
func (sm *StateMachine) IsIdle() bool {
	return sm.currentState == StateIdle
}
