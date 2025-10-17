// internal/models/actions.go
package models

// ActionType represents the type of action
type ActionType string

const (
	StandardAction ActionType = "Action"
	BonusAction    ActionType = "Bonus Action"
	Reaction       ActionType = "Reaction"
	FreeAction     ActionType = "Free Action"
)

// Action represents a character action
type Action struct {
	Name         string     `json:"name"`
	Type         ActionType `json:"type"`
	Description  string     `json:"description"`
	UsesPerRest  int        `json:"uses_per_rest"`  // -1 for unlimited
	UsesRemaining int       `json:"uses_remaining"`
	RestType     string     `json:"rest_type"` // "short" or "long"
}

// CanUse checks if the action can be used
func (a *Action) CanUse() bool {
	return a.UsesPerRest == -1 || a.UsesRemaining > 0
}

// Use uses the action if available
func (a *Action) Use() bool {
	if !a.CanUse() {
		return false
	}
	if a.UsesPerRest != -1 {
		a.UsesRemaining--
	}
	return true
}

// RestoreUses restores uses after a rest
func (a *Action) RestoreUses() {
	a.UsesRemaining = a.UsesPerRest
}

// ActionList holds all character actions
type ActionList struct {
	Actions []Action `json:"actions"`
}

// NewDefaultActions creates default D&D 5e actions
func NewDefaultActions() ActionList {
	return ActionList{
		Actions: []Action{
			{Name: "Attack", Type: StandardAction, Description: "Make a weapon attack", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Cast a Spell", Type: StandardAction, Description: "Cast a spell with casting time of 1 action", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Dash", Type: StandardAction, Description: "Double movement speed", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Disengage", Type: StandardAction, Description: "Movement doesn't provoke opportunity attacks", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Dodge", Type: StandardAction, Description: "Attacks against you have disadvantage", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Help", Type: StandardAction, Description: "Help an ally with a task", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Hide", Type: StandardAction, Description: "Make a Stealth check", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Ready", Type: StandardAction, Description: "Prepare an action for later", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Search", Type: StandardAction, Description: "Make a Perception or Investigation check", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Use Object", Type: StandardAction, Description: "Interact with an object", UsesPerRest: -1, UsesRemaining: -1},
			{Name: "Opportunity Attack", Type: Reaction, Description: "Attack when enemy leaves reach", UsesPerRest: -1, UsesRemaining: -1},
		},
	}
}

// AddAction adds a custom action
func (al *ActionList) AddAction(action Action) {
	al.Actions = append(al.Actions, action)
}

// RemoveAction removes an action by name
func (al *ActionList) RemoveAction(name string) bool {
	for i, action := range al.Actions {
		if action.Name == name {
			al.Actions = append(al.Actions[:i], al.Actions[i+1:]...)
			return true
		}
	}
	return false
}

// ShortRest restores actions that recharge on short rest
func (al *ActionList) ShortRest() {
	for i := range al.Actions {
		if al.Actions[i].RestType == "short" || al.Actions[i].RestType == "long" {
			al.Actions[i].RestoreUses()
		}
	}
}

// LongRest restores all actions
func (al *ActionList) LongRest() {
	for i := range al.Actions {
		al.Actions[i].RestoreUses()
	}
}
