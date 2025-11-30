// internal/ui/component_router.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/handlers"
)

// routeComponentToHandler routes a component to its appropriate handler based on component type
// Returns true if the component was handled, false otherwise
func (m *Model) routeComponentToHandler(component ComponentHandler, msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	// Use type assertions to route to the appropriate handler
	// Components are checked in priority order, so we handle the first matching type

	switch comp := component.(type) {
	case *components.StatGenerator:
		// TODO: Update handler to return bool
		model, cmd := m.handleStatGeneratorKeys(msg)
		return model, cmd, true

	case *components.AbilityRoller:
		model, cmd, handled := m.handleAbilityRollerKeys(msg)
		return model, cmd, handled

	case *components.AttackRoller:
		// TODO: Update handler to return bool
		model, cmd := m.handleAttackRollerKeys(msg)
		return model, cmd, true

	case *components.RestPopup:
		model, cmd, handled := m.handleRestPopupKeys(msg)
		return model, cmd, handled

	case *components.MessagePopup:
		model, cmd, handled := m.handleMessagePopupKeys(msg)
		return model, cmd, handled

	case *components.SpellSelector:
		// TODO: Update handler to return bool
		model, cmd := m.handleSpellSelectorKeys(msg)
		return model, cmd, true

	case *components.FeatSelector:
		// TODO: Update handler to return bool
		model, cmd := m.handleFeatSelectorKeys(msg)
		return model, cmd, true

	case *components.FeatDetailPopup:
		if handlers.HandleFeatDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.MasteryDetailPopup:
		if handlers.HandleMasteryDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.ManeuverDetailPopup:
		if handlers.HandleManeuverDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.ConsumableDetailPopup:
		if handlers.HandleConsumableDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.FeatureDetailPopup:
		if handlers.HandleFeatureDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.ItemDetailPopup:
		if handlers.HandleItemDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.SpellDetailPopup:
		if handlers.HandleSpellDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.OriginSelector:
		model, cmd := m.handleOriginSelectorKeys(msg)
		return model, cmd, true

	case *components.AlignmentSelector:
		model, cmd := m.handleAlignmentSelectorKeys(msg)
		return model, cmd, true

	case *components.OriginDetailPopup:
		if handlers.HandleOriginDetailPopupKeys(msg, comp, &m.message) {
			return m, nil, true
		}

	case *components.InputPopup:
		model, cmd := m.handleInputPopupKeys(msg)
		return model, cmd, true

	case *components.AbilityChoiceSelector:
		model, cmd := m.handleAbilityChoiceSelectorKeys(msg)
		return model, cmd, true

	case *components.SubtypeSelector:
		model, cmd := m.handleSubtypeSelectorKeys(msg)
		return model, cmd, true

	case *components.SkillSelector:
		model, cmd := m.handleSkillSelectorKeys(msg)
		return model, cmd, true

	case *components.LanguageSelector:
		model, cmd := m.handleLanguageSelectorKeys(msg)
		return model, cmd, true

	case *components.ToolSelector:
		model, cmd := m.handleToolSelectorKeys(msg)
		return model, cmd, true

	case *components.WeaponMasterySelector:
		model, cmd := m.handleWeaponMasterySelectorKeys(msg)
		return model, cmd, true

	case *components.ExpertiseSelector:
		model, cmd := m.handleExpertiseSelectorKeys(msg)
		return model, cmd, true

	case *components.ManeuverSelector:
		model, cmd := m.handleManeuverSelectorKeys(msg)
		return model, cmd, true

	case *components.ItemSelector:
		model, cmd := m.handleItemSelectorKeys(msg)
		return model, cmd, true

	case *components.LevelUpSelector:
		model, cmd := m.handleLevelUpSelectorKeys(msg)
		return model, cmd, true

	case *components.DeLevelSelector:
		model, cmd := m.handleDeLevelSelectorKeys(msg)
		return model, cmd, true

	case *components.FightingStyleSelector:
		model, cmd := m.handleFightingStyleSelectorKeys(msg)
		return model, cmd, true

	case *components.CantripSelector:
		model, cmd := m.handleCantripSelectorKeys(msg)
		return model, cmd, true

	case *components.LeveledSpellSelector:
		model, cmd := m.handleLeveledSpellSelectorKeys(msg)
		return model, cmd, true

	case *components.SchoolSpellSelector:
		model, cmd := m.handleSchoolSpellSelectorKeys(msg)
		return model, cmd, true

	case *components.SpellbookEditor:
		model, cmd := m.handleSpellbookEditorKeys(msg)
		return model, cmd, true

	case *components.SpellPrepSelector:
		model, cmd := m.handleSpellPrepSelectorKeys(msg)
		return model, cmd, true

	case *components.SlotRestorer:
		model, cmd := m.handleSlotRestorerKeys(msg)
		return model, cmd, true

	case *components.ClassSkillSelector:
		debug.Log("ComponentRouter: routing to ClassSkillSelector handler, key=%s", msg.String())
		model, cmd := m.handleClassSkillSelectorKeys(msg)
		return model, cmd, true

	case *components.SubclassSelector:
		model, cmd := m.handleSubclassSelectorKeys(msg)
		return model, cmd, true

	case *components.BeastSelector:
		model, cmd, handled := m.handleBeastSelectorKeys(msg)
		return model, cmd, handled

	case *components.CompanionSelector:
		// Only handle if selector is actually visible
		if m.companionSelector != nil && m.companionSelector.IsVisible() {
			model, cmd := m.handleCompanionSelectorKeys(msg)
			return model, cmd, true
		}
		// If not visible, don't handle - let keys fall through to panel
		return m, nil, false

	case *components.ClassSelector:
		model, cmd := m.handleClassSelectorKeys(msg)
		return model, cmd, true

	case *components.SpeciesSelector:
		model, cmd := m.handleSpeciesSelectorKeys(msg)
		return model, cmd, true
	}

	return m, nil, false
}

// handleSpecialComponents handles components that need special processing before normal routing
// Returns true if handled, false if should continue to normal routing
func (m *Model) handleSpecialComponents(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	// Trait selector needs to be checked BEFORE global keys to allow 'p' in custom input
	if m.traitSelector.IsVisible() {
		model, cmd := m.handleTraitSelectorKeys(msg)
		return model, cmd, true
	}

	// Backstory editor needs to be checked BEFORE global keys to allow 'q' in editor
	if m.backstoryEditor.IsVisible() {
		model, cmd := m.handleBackstoryEditorKeys(msg)
		return model, cmd, true
	}

	// Divine Order selector has special tab passthrough logic
	if m.IsDivineOrderSelectorVisible() {
		// Check if Divine Order has already been applied
		divineOrderApplied := false
		if m.character.BenefitTracker != nil {
			for _, benefit := range m.character.BenefitTracker.Benefits {
				if benefit.Source.Type == "class_feature" && benefit.Source.Name == "Divine Order" {
					divineOrderApplied = true
					break
				}
			}
		}

		// If already applied, hide selector and allow navigation
		if divineOrderApplied {
			m.SetDivineOrderSelectorVisible(false)
			m.stateMachine.ClearContext("pendingDivineOrder")
			m.pendingDivineOrder = ""
			// Continue to normal routing (don't return handled)
			return m, nil, false
		}

		// Allow tab navigation to pass through even when selector is visible
		if msg.String() == "tab" || msg.String() == "shift+tab" {
			// Tab navigation - let it pass through to panel navigation check below
			return m, nil, false
		}

		// Handle other keys in selector
		model, cmd := m.handleDivineOrderSelectorKeys(msg)
		return model, cmd, true
	}

	// Primal Order selector has special tab passthrough logic
	if m.IsPrimalOrderSelectorVisible() {
		// Check if Primal Order has already been applied
		primalOrderApplied := false
		if m.character.BenefitTracker != nil {
			for _, benefit := range m.character.BenefitTracker.Benefits {
				if benefit.Source.Type == "class_feature" && benefit.Source.Name == "Primal Order" {
					primalOrderApplied = true
					break
				}
			}
		}

		// If already applied, hide selector and allow navigation
		if primalOrderApplied {
			m.SetPrimalOrderSelectorVisible(false)
			m.stateMachine.ClearContext("pendingPrimalOrder")
			m.pendingPrimalOrder = ""
			// Continue to normal routing (don't return handled)
			return m, nil, false
		}

		// Allow tab navigation to pass through even when selector is visible
		if msg.String() == "tab" || msg.String() == "shift+tab" {
			// Tab navigation - let it pass through to panel navigation check below
			return m, nil, false
		}

		// Handle other keys in selector
		model, cmd := m.handlePrimalOrderSelectorKeys(msg)
		return model, cmd, true
	}

	return m, nil, false
}
