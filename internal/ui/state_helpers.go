// internal/ui/state_helpers.go
package ui

import (
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/state"
)

// State helper methods to transition from direct field access to StateMachine
// These methods provide backward compatibility while using StateMachine internally

// PendingFeat methods
func (m *Model) GetPendingFeat() *models.Feat {
	if feat, ok := m.stateMachine.GetContext("pendingFeat").(*models.Feat); ok {
		return feat
	}
	return m.pendingFeat // Fallback to legacy field
}

func (m *Model) SetPendingFeat(feat *models.Feat) {
	m.stateMachine.SetContext("pendingFeat", feat)
	m.stateMachine.Transition(state.StateFeatSelection, map[string]interface{}{"pendingFeat": feat})
	m.pendingFeat = feat // Keep legacy field in sync for now
}

func (m *Model) ClearPendingFeat() {
	m.stateMachine.ClearContext("pendingFeat")
	m.stateMachine.Transition(state.StateIdle, nil)
	m.pendingFeat = nil // Keep legacy field in sync for now
}

// PendingOrigin methods
func (m *Model) GetPendingOrigin() *models.Origin {
	if origin, ok := m.stateMachine.GetContext("pendingOrigin").(*models.Origin); ok {
		return origin
	}
	return m.pendingOrigin // Fallback to legacy field
}

func (m *Model) SetPendingOrigin(origin *models.Origin) {
	m.stateMachine.SetContext("pendingOrigin", origin)
	m.stateMachine.Transition(state.StateAbilityChoice, map[string]interface{}{"pendingOrigin": origin})
	m.pendingOrigin = origin // Keep legacy field in sync for now
}

func (m *Model) ClearPendingOrigin() {
	m.stateMachine.ClearContext("pendingOrigin")
	m.stateMachine.Transition(state.StateIdle, nil)
	m.pendingOrigin = nil // Keep legacy field in sync for now
}

// DivineOrder methods
func (m *Model) GetPendingDivineOrder() string {
	if order, ok := m.stateMachine.GetContext("pendingDivineOrder").(string); ok {
		return order
	}
	return m.pendingDivineOrder // Fallback to legacy field
}

func (m *Model) SetPendingDivineOrder(order string) {
	m.stateMachine.SetContext("pendingDivineOrder", order)
	m.stateMachine.Transition(state.StateDivineOrderSelection, map[string]interface{}{"pendingDivineOrder": order})
	m.pendingDivineOrder = order // Keep legacy field in sync for now
}

func (m *Model) GetPendingDivineOrderSkill() string {
	if skill, ok := m.stateMachine.GetContext("pendingDivineOrderSkill").(string); ok {
		return skill
	}
	return m.pendingDivineOrderSkill // Fallback to legacy field
}

func (m *Model) SetPendingDivineOrderSkill(skill string) {
	m.stateMachine.SetContext("pendingDivineOrderSkill", skill)
	m.pendingDivineOrderSkill = skill // Keep legacy field in sync for now
}

func (m *Model) IsDivineOrderSelectorVisible() bool {
	if visible, ok := m.stateMachine.GetContext("divineOrderSelectorVisible").(bool); ok {
		return visible
	}
	return m.divineOrderSelectorVisible // Fallback to legacy field
}

func (m *Model) SetDivineOrderSelectorVisible(visible bool) {
	m.stateMachine.SetContext("divineOrderSelectorVisible", visible)
	if visible {
		m.stateMachine.Transition(state.StateDivineOrderSelection, map[string]interface{}{"divineOrderSelectorVisible": true})
	} else {
		m.stateMachine.ClearContext("divineOrderSelectorVisible")
		if m.stateMachine.GetState() == state.StateDivineOrderSelection {
			m.stateMachine.Transition(state.StateIdle, nil)
		}
	}
	m.divineOrderSelectorVisible = visible // Keep legacy field in sync for now
}

// EldritchKnightSpells methods
func (m *Model) GetEldritchKnightSpellsSelected() int {
	if count, ok := m.stateMachine.GetContext("eldritchKnightSpellsSelected").(int); ok {
		return count
	}
	return m.eldritchKnightSpellsSelected // Fallback to legacy field
}

func (m *Model) SetEldritchKnightSpellsSelected(count int) {
	m.stateMachine.SetContext("eldritchKnightSpellsSelected", count)
	m.stateMachine.Transition(state.StateEldritchKnightSpellSelection, map[string]interface{}{"eldritchKnightSpellsSelected": count})
	m.eldritchKnightSpellsSelected = count // Keep legacy field in sync for now
}

func (m *Model) GetEldritchKnightSpells() []models.Spell {
	if spells, ok := m.stateMachine.GetContext("eldritchKnightSpells").([]models.Spell); ok {
		return spells
	}
	return m.eldritchKnightSpells // Fallback to legacy field
}

func (m *Model) SetEldritchKnightSpells(spells []models.Spell) {
	m.stateMachine.SetContext("eldritchKnightSpells", spells)
	m.eldritchKnightSpells = spells // Keep legacy field in sync for now
}

func (m *Model) ClearEldritchKnightSpells() {
	m.stateMachine.ClearContext("eldritchKnightSpells")
	m.stateMachine.ClearContext("eldritchKnightSpellsSelected")
	m.stateMachine.Transition(state.StateIdle, nil)
	m.eldritchKnightSpells = nil // Keep legacy field in sync for now
	m.eldritchKnightSpellsSelected = 0
}

// StudentOfWarToolSelected methods
func (m *Model) IsStudentOfWarToolSelected() bool {
	if selected, ok := m.stateMachine.GetContext("studentOfWarToolSelected").(bool); ok {
		return selected
	}
	return m.studentOfWarToolSelected // Fallback to legacy field
}

func (m *Model) SetStudentOfWarToolSelected(selected bool) {
	m.stateMachine.SetContext("studentOfWarToolSelected", selected)
	m.studentOfWarToolSelected = selected // Keep legacy field in sync for now
}

// InputPopupContext methods
func (m *Model) GetInputPopupContext() string {
	if context, ok := m.stateMachine.GetContext("inputPopupContext").(string); ok {
		return context
	}
	return m.inputPopupContext // Fallback to legacy field
}

func (m *Model) SetInputPopupContext(context string) {
	m.stateMachine.SetContext("inputPopupContext", context)
	m.inputPopupContext = context // Keep legacy field in sync for now
}

func (m *Model) ClearInputPopupContext() {
	m.stateMachine.ClearContext("inputPopupContext")
	m.inputPopupContext = "" // Keep legacy field in sync for now
}
