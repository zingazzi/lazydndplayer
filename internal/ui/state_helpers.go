// internal/ui/state_helpers.go
package ui

import (
	"github.com/marcozingoni/lazydndplayer/internal/debug"
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
	debug.Log("SetPendingFeat: Setting pending feat=%s, IsInWizard=%v", feat.Name, m.IsInWizard())
	// If in wizard mode, preserve wizard state; otherwise transition to feat selection
	if m.IsInWizard() {
		// Keep wizard state, just add pendingFeat to context
		m.stateMachine.SetContext("pendingFeat", feat)
	} else {
		m.stateMachine.Transition(state.StateFeatSelection, map[string]interface{}{"pendingFeat": feat})
	}
	m.pendingFeat = feat // Keep legacy field in sync for now
}

func (m *Model) ClearPendingFeat() {
	debug.Log("ClearPendingFeat: Clearing pending feat, IsInWizard=%v", m.IsInWizard())
	m.stateMachine.ClearContext("pendingFeat")
	// If in wizard mode, preserve wizard state; otherwise transition to idle
	if !m.IsInWizard() {
		m.stateMachine.Transition(state.StateIdle, nil)
	}
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
	debug.Log("SetPendingOrigin: Setting pending origin=%s, IsInWizard=%v", origin.Name, m.IsInWizard())
	// If in wizard mode, preserve wizard state; otherwise transition to ability choice
	if m.IsInWizard() {
		// Keep wizard state, just add pendingOrigin to context
		m.stateMachine.SetContext("pendingOrigin", origin)
	} else {
		m.stateMachine.Transition(state.StateAbilityChoice, map[string]interface{}{"pendingOrigin": origin})
	}
	m.pendingOrigin = origin // Keep legacy field in sync for now
}

func (m *Model) ClearPendingOrigin() {
	debug.Log("ClearPendingOrigin: Clearing pending origin, IsInWizard=%v", m.IsInWizard())
	m.stateMachine.ClearContext("pendingOrigin")
	// If in wizard mode, preserve wizard state; otherwise transition to idle
	if !m.IsInWizard() {
		m.stateMachine.Transition(state.StateIdle, nil)
	}
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

// Wizard state helper methods
func (m *Model) IsInWizard() bool {
	isWizard := m.stateMachine.GetState() == state.StateCharacterCreationWizard
	// Only log occasionally to avoid spam
	// debug.Log("IsInWizard: state=%v, result=%v", m.stateMachine.GetState(), isWizard)
	return isWizard
}

func (m *Model) GetWizardStep() int {
	if step, ok := m.stateMachine.GetContext("wizardStep").(int); ok {
		return step
	}
	return 0
}

func (m *Model) SetWizardStep(step int) {
	debug.Log("SetWizardStep: Setting wizard step to %d", step)
	m.stateMachine.SetContext("wizardStep", step)
	m.stateMachine.Transition(state.StateCharacterCreationWizard, map[string]interface{}{
		"wizardStep": step,
	})
	debug.Log("SetWizardStep: Wizard step set, current state=%v, step=%d", m.stateMachine.GetState(), m.GetWizardStep())
}

func (m *Model) ClearWizardState() {
	debug.Log("ClearWizardState: Clearing wizard state")
	m.stateMachine.ClearContext("wizardStep")
	m.stateMachine.Transition(state.StateIdle, nil)
}
