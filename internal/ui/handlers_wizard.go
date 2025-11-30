// internal/ui/handlers_wizard.go
package ui

import (
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/state"
)

// Wizard step constants
const (
	WizardStepClass        = 1
	WizardStepSpecies      = 2
	WizardStepOrigin       = 3
	WizardStepAbilityScores = 4
	WizardStepAlignment    = 5
)

// startWizard initializes and starts the character creation wizard
func (m *Model) startWizard() {
	debug.Log("startWizard: Starting character creation wizard")
	// Check if character already has a class
	if m.character.Class != "" || len(m.character.Classes) > 0 {
		debug.Log("startWizard: Character already has a class (Class='%s', Classes=%v), aborting", m.character.Class, len(m.character.Classes))
		m.message = "Character already has a class. Wizard only works for new characters."
		return
	}

	// Set wizard state
	debug.Log("startWizard: Setting wizard state to StateCharacterCreationWizard, step=%d", WizardStepClass)
	m.stateMachine.Transition(state.StateCharacterCreationWizard, map[string]interface{}{
		"wizardStep": WizardStepClass,
	})

	// Start with class selection
	m.classSelector.Show()
	m.message = "Step 1/5: Class Selection - Choose your class"
	debug.Log("startWizard: Wizard started, showing class selector")
}

// cancelWizard cancels the wizard and resets the character
func (m *Model) cancelWizard() {
	debug.Log("cancelWizard: Cancelling wizard and resetting character")
	// Reset character to defaults
	*m.character = *models.NewCharacter()

	// Clear wizard state
	m.stateMachine.Clear()
	m.stateMachine.Transition(state.StateIdle, nil)
	debug.Log("cancelWizard: Wizard state cleared, transitioned to StateIdle")

	// Hide all selectors
	m.classSelector.Hide()
	m.speciesSelector.Hide()
	m.subtypeSelector.Hide()
	m.originSelector.Hide()
	m.statGenerator.Hide()
	m.alignmentSelector.Hide()

	m.message = "Character creation cancelled - character reset"
	m.storage.Save(m.character)
	debug.Log("cancelWizard: Character reset and saved")
}

// advanceWizardStep moves to the next step in the wizard
func (m *Model) advanceWizardStep() {
	currentStep := m.GetWizardStep()
	debug.Log("advanceWizardStep: Current step=%d, IsInWizard=%v", currentStep, m.IsInWizard())

	if !m.IsInWizard() {
		debug.Log("advanceWizardStep: Not in wizard mode, aborting")
		return
	}

	switch currentStep {
	case WizardStepClass:
		// Move to species selection
		debug.Log("advanceWizardStep: Moving from Class (step 1) to Species (step 2)")
		m.SetWizardStep(WizardStepSpecies)
		m.speciesSelector.Show()
		m.message = "Step 2/5: Species Selection - Choose your species"
		debug.Log("advanceWizardStep: Species selector shown, wizard step set to %d", WizardStepSpecies)

	case WizardStepSpecies:
		// Move to origin selection
		debug.Log("advanceWizardStep: Moving from Species (step 2) to Origin (step 3)")
		m.SetWizardStep(WizardStepOrigin)
		m.originSelector.Show(m.character)
		m.message = "Step 3/5: Origin Selection - Choose your origin"
		debug.Log("advanceWizardStep: Origin selector shown, wizard step set to %d", WizardStepOrigin)

	case WizardStepOrigin:
		// Move to ability scores
		debug.Log("advanceWizardStep: Moving from Origin (step 3) to Ability Scores (step 4)")
		m.SetWizardStep(WizardStepAbilityScores)
		m.statGenerator.Show(&m.character.AbilityScores)
		m.message = "Step 4/5: Ability Scores - Set your ability scores"
		debug.Log("advanceWizardStep: Stat generator shown, wizard step set to %d", WizardStepAbilityScores)

	case WizardStepAbilityScores:
		// Move to alignment
		debug.Log("advanceWizardStep: Moving from Ability Scores (step 4) to Alignment (step 5)")
		m.SetWizardStep(WizardStepAlignment)
		m.alignmentSelector.Show()
		m.message = "Step 5/5: Alignment - Choose your alignment"
		debug.Log("advanceWizardStep: Alignment selector shown, wizard step set to %d", WizardStepAlignment)

	case WizardStepAlignment:
		// Complete wizard
		debug.Log("advanceWizardStep: Moving from Alignment (step 5) to completion")
		m.completeWizard()

	default:
		debug.Log("advanceWizardStep: Unknown wizard step %d", currentStep)
		m.message = "Unknown wizard step"
	}
}

// completeWizard finishes the wizard and saves the character
func (m *Model) completeWizard() {
	debug.Log("completeWizard: Completing wizard")
	// Clear wizard state
	m.stateMachine.Clear()
	m.stateMachine.Transition(state.StateIdle, nil)
	debug.Log("completeWizard: Wizard state cleared, transitioned to StateIdle")

	// Hide alignment selector
	m.alignmentSelector.Hide()

	// Save character
	m.storage.Save(m.character)
	debug.Log("completeWizard: Character saved, wizard complete")

	m.message = "Character creation complete! Your character is ready."
}

// getWizardProgressMessage returns a formatted progress message
func (m *Model) getWizardProgressMessage(step int) string {
	switch step {
	case WizardStepClass:
		return "Step 1/5: Class Selection - Choose your class"
	case WizardStepSpecies:
		return "Step 2/5: Species Selection - Choose your species"
	case WizardStepOrigin:
		return "Step 3/5: Origin Selection - Choose your origin"
	case WizardStepAbilityScores:
		return "Step 4/5: Ability Scores - Set your ability scores"
	case WizardStepAlignment:
		return "Step 5/5: Alignment - Choose your alignment"
	default:
		return "Character Creation Wizard"
	}
}

// checkAndAdvanceWizardAfterClassSetup checks if class setup is complete and advances wizard if needed
// This should be called after any class-related selection (subclass, cantrips, fighting style, etc.)
func (m *Model) checkAndAdvanceWizardAfterClassSetup() {
	debug.Log("checkAndAdvanceWizardAfterClassSetup: Checking if class setup is complete")
	if !m.IsInWizard() {
		debug.Log("checkAndAdvanceWizardAfterClassSetup: Not in wizard mode, aborting")
		return
	}

	// Check if any class-related selectors are still visible
	classVisible := m.classSelector.IsVisible()
	skillVisible := m.classSkillSelector.IsVisible()
	subclassVisible := m.subclassSelector.IsVisible()
	cantripVisible := m.cantripSelector.IsVisible()
	fightingStyleVisible := m.fightingStyleSelector.IsVisible()
	weaponMasteryVisible := m.weaponMasterySelector.IsVisible()
	expertiseVisible := m.expertiseSelector.IsVisible()
	divineOrderVisible := m.IsDivineOrderSelectorVisible()
	primalOrderVisible := m.IsPrimalOrderSelectorVisible()

	debug.Log("checkAndAdvanceWizardAfterClassSetup: Selector visibility - class=%v, skill=%v, subclass=%v, cantrip=%v, fightingStyle=%v, weaponMastery=%v, expertise=%v, divineOrder=%v, primalOrder=%v",
		classVisible, skillVisible, subclassVisible, cantripVisible, fightingStyleVisible, weaponMasteryVisible, expertiseVisible, divineOrderVisible, primalOrderVisible)

	if classVisible ||
		skillVisible ||
		subclassVisible ||
		cantripVisible ||
		fightingStyleVisible ||
		weaponMasteryVisible ||
		expertiseVisible ||
		divineOrderVisible ||
		primalOrderVisible {
		// Still have class setup to do, don't advance yet
		debug.Log("checkAndAdvanceWizardAfterClassSetup: Class setup not complete, waiting for more selections")
		return
	}

	// All class setup is complete, advance wizard
	debug.Log("checkAndAdvanceWizardAfterClassSetup: All class setup complete, advancing wizard")
	m.advanceWizardStep()
}

// checkAndAdvanceWizardAfterSpeciesSetup checks if species setup is complete and advances wizard if needed
// This should be called after any species-related selection (subtype, languages, skills, spells, feats)
func (m *Model) checkAndAdvanceWizardAfterSpeciesSetup() {
	debug.Log("checkAndAdvanceWizardAfterSpeciesSetup: Checking if species setup is complete")
	if !m.IsInWizard() {
		debug.Log("checkAndAdvanceWizardAfterSpeciesSetup: Not in wizard mode, aborting")
		return
	}

	// Check if any species-related selectors are still visible
	speciesVisible := m.speciesSelector.IsVisible()
	subtypeVisible := m.subtypeSelector.IsVisible()
	languageVisible := m.languageSelector.IsVisible()
	skillVisible := m.skillSelector.IsVisible()
	spellVisible := m.spellSelector.IsVisible()
	featVisible := m.featSelector.IsVisible()

	debug.Log("checkAndAdvanceWizardAfterSpeciesSetup: Selector visibility - species=%v, subtype=%v, language=%v, skill=%v, spell=%v, feat=%v",
		speciesVisible, subtypeVisible, languageVisible, skillVisible, spellVisible, featVisible)

	if speciesVisible ||
		subtypeVisible ||
		languageVisible ||
		skillVisible ||
		spellVisible ||
		featVisible {
		// Still have species setup to do, don't advance yet
		debug.Log("checkAndAdvanceWizardAfterSpeciesSetup: Species setup not complete, waiting for more selections")
		return
	}

	// All species setup is complete, advance wizard
	debug.Log("checkAndAdvanceWizardAfterSpeciesSetup: All species setup complete, advancing wizard")
	m.advanceWizardStep()
}
