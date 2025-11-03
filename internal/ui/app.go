// internal/ui/app.go
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/dice"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/handlers"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
	"github.com/marcozingoni/lazydndplayer/internal/ui/state"
	"github.com/marcozingoni/lazydndplayer/internal/ui/view"
)

// PanelType represents the current active main panel
type PanelType int

const (
	StatsPanel PanelType = iota
	SkillsPanel
	InventoryPanel
	SpellsPanel
	FeaturesPanel
	TraitsPanel
	OriginPanel
)

// Popup size constants for different popup types
const (
	// Small popups (language, tool, ability choice selectors)
	PopupSmallWidthPercent  = 0.50  // 50% of screen width
	PopupSmallHeightPercent = 0.60  // 60% of screen height
	PopupSmallMinWidth      = 60    // Minimum width in characters
	PopupSmallMinHeight     = 20    // Minimum height in lines

	// Medium popups (feat, origin, species selectors)
	PopupMediumWidthPercent  = 0.75  // 75% of screen width
	PopupMediumHeightPercent = 0.80  // 80% of screen height
	PopupMediumMinWidth      = 80    // Minimum width in characters
	PopupMediumMinHeight     = 25    // Minimum height in lines

	// Large popups (item selector, spell selector)
	PopupLargeWidthPercent  = 0.85  // 85% of screen width
	PopupLargeHeightPercent = 0.85  // 85% of screen height
	PopupLargeMinWidth      = 90    // Minimum width in characters
	PopupLargeMinHeight     = 30    // Minimum height in lines
)

// FocusArea represents which area of the UI has focus
type FocusArea int

const (
	FocusMain FocusArea = iota
	FocusCharStats
	FocusActions
	FocusDice
)

// Model is the main application model
type Model struct {
	character    *models.Character
	storage      *storage.Storage

	// UI Components
	tabs             *components.Tabs
	help             *components.Help
	speciesSelector  *components.SpeciesSelector
	subtypeSelector  *components.SubtypeSelector
	languageSelector *components.LanguageSelector
	skillSelector    *components.SkillSelector
	spellSelector    *components.SpellSelector
	featSelector          *components.FeatSelector
	featDetailPopup       *components.FeatDetailPopup
	featureDetailPopup    *components.FeatureDetailPopup
	itemDetailPopup       *components.ItemDetailPopup
	masteryDetailPopup    *components.MasteryDetailPopup
	maneuverDetailPopup   *components.ManeuverDetailPopup
	consumableDetailPopup *components.ConsumableDetailPopup
	spellDetailPopup      *components.SpellDetailPopup
	originSelector        *components.OriginSelector
	alignmentSelector     *components.AlignmentSelector
	traitSelector         *components.TraitSelector
	backstoryEditor       *components.BackstoryEditor
	originDetailPopup     *components.OriginDetailPopup
	inputPopup            *components.InputPopup
	toolSelector          *components.ToolSelector
	itemSelector           *components.ItemSelector
	classSelector          *components.ClassSelector
	classSkillSelector     *components.ClassSkillSelector
	subclassSelector       *components.SubclassSelector
	fightingStyleSelector  *components.FightingStyleSelector
	cantripSelector        *components.CantripSelector
	leveledSpellSelector   *components.LeveledSpellSelector
	schoolSpellSelector    *components.SchoolSpellSelector
	spellbookEditor        *components.SpellbookEditor
	spellPrepSelector      *components.SpellPrepSelector
	slotRestorer           *components.SlotRestorer
	statGenerator          *components.StatGenerator
	abilityRoller         *components.AbilityRoller
	abilityChoiceSelector *components.AbilityChoiceSelector
	attackRoller          *components.AttackRoller
	attackMenu            *components.AttackMenu
	weaponMasterySelector *components.WeaponMasterySelector
	expertiseSelector     *components.ExpertiseSelector
	maneuverSelector      *components.ManeuverSelector
	levelUpSelector       *components.LevelUpSelector
	deLevelSelector       *components.DeLevelSelector
	restPopup             *components.RestPopup
	messagePopup          *components.MessagePopup

	// Main Panels (switchable)
	statsPanel     *panels.StatsPanel
	skillsPanel    *panels.SkillsPanel
	inventoryPanel *panels.InventoryPanel
	spellsPanel    *panels.SpellsPanel
	featuresPanel  *panels.FeaturesPanel
	traitsPanel    *panels.TraitsPanel
	originPanel    *panels.OriginPanel

	// Fixed Panels (always visible)
	dicePanel           *panels.DicePanel
	characterStatsPanel *panels.CharacterStatsPanel
	actionsPanel        *panels.ActionsPanel // Bottom panel for quick actions

	// Component Management
	componentManager *ComponentManager
	stateMachine     *state.StateMachine

	// State
	currentPanel       PanelType
	focusArea          FocusArea
	width              int
	height             int
	ready              bool
	message            string
	quitting           bool
	// Legacy state fields (will be replaced by StateMachine in Phase 2)
	// Keeping for now to maintain compatibility
	pendingFeat        *models.Feat   // Temporarily store feat while choosing ability
	pendingOrigin      *models.Origin // Temporarily store origin while choosing ability
	pendingDivineOrder string         // Temporarily store divine order choice ("Protector" or "Thaumaturgic")
	pendingDivineOrderSkill string    // For Thaumaturgic: chosen skill (Arcana or Religion)
	divineOrderSelectorVisible bool   // Flag for showing divine order selection
	pendingChanges     *models.PendingChanges // Transaction system for rollback support
	eldritchKnightSpellsSelected int   // Counter for Eldritch Knight spell selection (0-3)
	eldritchKnightSpells []models.Spell // Temporarily store selected spells
	studentOfWarToolSelected bool // Flag for Student of War tool selection flow
	inputPopupContext string // Context for what is being edited in inputPopup
}

// NewModel creates a new application model
func NewModel(char *models.Character, store *storage.Storage) *Model {
	// Initialize hit dice for the character
	classData, err := models.LoadClassesFromJSON("data/classes")
	if err == nil && classData != nil {
		// Convert []Class to map[string]*Class
		classMap := make(map[string]*models.Class)
		for i := range classData.Classes {
			classMap[classData.Classes[i].Name] = &classData.Classes[i]
		}
		char.InitializeHitDice(classMap)
	}

	// Create component manager and state machine
	componentManager := NewComponentManager()
	stateMachine := state.NewStateMachine()

	// Create components
	tabs := components.NewTabs()
	help := components.NewHelp()
	speciesSelector := components.NewSpeciesSelector()
	subtypeSelector := components.NewSubtypeSelector()
	languageSelector := components.NewLanguageSelector()
	skillSelector := components.NewSkillSelector()
	spellSelector := components.NewSpellSelector()
	featSelector := components.NewFeatSelector()
	featDetailPopup := components.NewFeatDetailPopup()
	featureDetailPopup := components.NewFeatureDetailPopup()
	itemDetailPopup := components.NewItemDetailPopup()
	masteryDetailPopup := components.NewMasteryDetailPopup()
	maneuverDetailPopup := components.NewManeuverDetailPopup()
	consumableDetailPopup := components.NewConsumableDetailPopup()
	spellDetailPopup := components.NewSpellDetailPopup()
	originSelector := components.NewOriginSelector()
	alignmentSelector := components.NewAlignmentSelector()
	traitSelector := components.NewTraitSelector()
	backstoryEditor := components.NewBackstoryEditor()
	originDetailPopup := components.NewOriginDetailPopup()
	inputPopup := components.NewInputPopup()
	toolSelector := components.NewToolSelector()
	itemSelector := components.NewItemSelector()
	classSelector := components.NewClassSelector(char)
	classSkillSelector := components.NewClassSkillSelector()
	subclassSelector := components.NewSubclassSelector(char)
	fightingStyleSelector := components.NewFightingStyleSelector()
	cantripSelector := components.NewCantripSelector(char)
	leveledSpellSelector := components.NewLeveledSpellSelector(char)
	schoolSpellSelector := components.NewSchoolSpellSelector(char)
	spellbookEditor := components.NewSpellbookEditor(char)
	spellPrepSelector := components.NewSpellPrepSelector(char)
	slotRestorer := components.NewSlotRestorer(char)
	statGenerator := components.NewStatGenerator()
	abilityRoller := components.NewAbilityRoller()
	abilityChoiceSelector := components.NewAbilityChoiceSelector()
	attackRoller := components.NewAttackRoller()
	attackMenu := components.NewAttackMenu()
	weaponMasterySelector := components.NewWeaponMasterySelector(char)
	expertiseSelector := components.NewExpertiseSelector(char)
	maneuverSelector := components.NewManeuverSelector()
	levelUpSelector := components.NewLevelUpSelector(char)
	deLevelSelector := components.NewDeLevelSelector(char)
	restPopup := components.NewRestPopup(char, models.NewStandardDiceRoller())
	messagePopup := components.NewMessagePopup()

	// Register components with priorities (higher number = higher priority)
	// Highest priority components first
	componentManager.Register(statGenerator, 100, "statGenerator")
	componentManager.Register(abilityRoller, 99, "abilityRoller")
	componentManager.Register(attackRoller, 98, "attackRoller")
	componentManager.Register(messagePopup, 97, "messagePopup")
	componentManager.Register(spellSelector, 96, "spellSelector")
	componentManager.Register(featSelector, 95, "featSelector")
	componentManager.Register(featDetailPopup, 94, "featDetailPopup")
	componentManager.Register(masteryDetailPopup, 93, "masteryDetailPopup")
	componentManager.Register(maneuverDetailPopup, 92, "maneuverDetailPopup")
	componentManager.Register(consumableDetailPopup, 91, "consumableDetailPopup")
	componentManager.Register(featureDetailPopup, 90, "featureDetailPopup")
	componentManager.Register(itemDetailPopup, 89, "itemDetailPopup")
	componentManager.Register(spellDetailPopup, 88, "spellDetailPopup")
	componentManager.Register(originSelector, 87, "originSelector")
	componentManager.Register(alignmentSelector, 86, "alignmentSelector")
	componentManager.Register(traitSelector, 85, "traitSelector")
	componentManager.Register(backstoryEditor, 84, "backstoryEditor")
	componentManager.Register(originDetailPopup, 83, "originDetailPopup")
	componentManager.Register(inputPopup, 82, "inputPopup")
	componentManager.Register(abilityChoiceSelector, 81, "abilityChoiceSelector")
	componentManager.Register(subtypeSelector, 80, "subtypeSelector")
	componentManager.Register(skillSelector, 79, "skillSelector")
	componentManager.Register(languageSelector, 78, "languageSelector")
	componentManager.Register(toolSelector, 77, "toolSelector")
	componentManager.Register(weaponMasterySelector, 76, "weaponMasterySelector")
	componentManager.Register(expertiseSelector, 75, "expertiseSelector")
	componentManager.Register(maneuverSelector, 74, "maneuverSelector")
	componentManager.Register(levelUpSelector, 73, "levelUpSelector")
	componentManager.Register(deLevelSelector, 72, "deLevelSelector")
	componentManager.Register(itemSelector, 71, "itemSelector")
	componentManager.Register(fightingStyleSelector, 70, "fightingStyleSelector")
	componentManager.Register(cantripSelector, 69, "cantripSelector")
	componentManager.Register(leveledSpellSelector, 68, "leveledSpellSelector")
	componentManager.Register(schoolSpellSelector, 67, "schoolSpellSelector")
	componentManager.Register(spellbookEditor, 66, "spellbookEditor")
	componentManager.Register(spellPrepSelector, 65, "spellPrepSelector")
	componentManager.Register(slotRestorer, 64, "slotRestorer")
	componentManager.Register(classSkillSelector, 63, "classSkillSelector")
	componentManager.Register(subclassSelector, 62, "subclassSelector")
	componentManager.Register(classSelector, 61, "classSelector")
	componentManager.Register(speciesSelector, 60, "speciesSelector")
	componentManager.Register(restPopup, 59, "restPopup")
	componentManager.Register(attackMenu, 58, "attackMenu")

	return &Model{
		character:           char,
		storage:             store,
		tabs:                tabs,
		help:                help,
		speciesSelector:     speciesSelector,
		subtypeSelector:     subtypeSelector,
		languageSelector:    languageSelector,
		skillSelector:       skillSelector,
		spellSelector:       spellSelector,
		featSelector:          featSelector,
		featDetailPopup:       featDetailPopup,
		featureDetailPopup:    featureDetailPopup,
		itemDetailPopup:       itemDetailPopup,
		masteryDetailPopup:    masteryDetailPopup,
		maneuverDetailPopup:   maneuverDetailPopup,
		consumableDetailPopup: consumableDetailPopup,
		spellDetailPopup:      spellDetailPopup,
		originSelector:        originSelector,
		alignmentSelector:     alignmentSelector,
		traitSelector:         traitSelector,
		backstoryEditor:       backstoryEditor,
		originDetailPopup:     originDetailPopup,
		inputPopup:            inputPopup,
		toolSelector:          toolSelector,
		itemSelector:           itemSelector,
		classSelector:          classSelector,
		classSkillSelector:     classSkillSelector,
		subclassSelector:       subclassSelector,
		fightingStyleSelector:  fightingStyleSelector,
		cantripSelector:        cantripSelector,
		leveledSpellSelector:   leveledSpellSelector,
		schoolSpellSelector:    schoolSpellSelector,
		spellbookEditor:        spellbookEditor,
		spellPrepSelector:      spellPrepSelector,
		slotRestorer:           slotRestorer,
		statGenerator:          statGenerator,
		abilityRoller:         abilityRoller,
		abilityChoiceSelector: abilityChoiceSelector,
		attackRoller:          attackRoller,
		attackMenu:            attackMenu,
		weaponMasterySelector: weaponMasterySelector,
		expertiseSelector:     expertiseSelector,
		maneuverSelector:      maneuverSelector,
		levelUpSelector:       levelUpSelector,
		deLevelSelector:       deLevelSelector,
		restPopup:             restPopup,
		messagePopup:          messagePopup,
		statsPanel:            panels.NewStatsPanel(char),
		skillsPanel:           panels.NewSkillsPanel(char),
		inventoryPanel:        panels.NewInventoryPanel(char),
		spellsPanel:           panels.NewSpellsPanel(char),
		featuresPanel:         panels.NewFeaturesPanel(char),
		traitsPanel:           panels.NewTraitsPanel(char),
		originPanel:           panels.NewOriginPanel(char),
		dicePanel:           panels.NewDicePanel(char),
		characterStatsPanel: panels.NewCharacterStatsPanel(char),
		actionsPanel:        panels.NewActionsPanel(char),
		componentManager:    componentManager,
		stateMachine:        stateMachine,
		currentPanel:        StatsPanel,
		focusArea:           FocusMain,
		pendingChanges:      models.NewPendingChanges(),
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.ClearScreen
}

// GetComponentManager returns the component manager
func (m *Model) GetComponentManager() *ComponentManager {
	return m.componentManager
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		// Help overlay takes priority
		if m.help.Visible {
			switch msg.String() {
			case "?", "esc":
				m.help.Toggle()
			}
			return m, nil
		}

		// Check if trait selector is active (BEFORE global keys to allow 'p' in custom input)
		if m.traitSelector.IsVisible() {
			return m.handleTraitSelectorKeys(msg)
		}

		// Check if backstory editor is active (BEFORE global keys to allow 'q' in editor)
		if m.backstoryEditor.IsVisible() {
			return m.handleBackstoryEditorKeys(msg)
		}

		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "?":
			m.help.Toggle()
			return m, nil

		// Note: 's' key is now handled in panel-specific handlers (e.g., Traits panel for Fighting Style, Spells panel for slot restore)
		// Removed global 's' handler (was for manual save) to avoid conflicts - auto-save happens on most actions
		// Note: 'l' and 'L' keys are now handled in Traits panel for language management
		// Removed global 'l' handler (was for level up) to avoid conflicts
		// Note: 'f' and 'F' keys are now handled in Traits panel for feat management

		// Focus cycling - p key cycles through Main, CharStats, Actions, Dice
		case "p":
			m.focusArea = (m.focusArea + 1) % 4
			switch m.focusArea {
			case FocusMain:
				m.message = "Focus: Main Panel"
			case FocusCharStats:
				m.message = "Focus: Character Stats"
			case FocusActions:
				m.message = "Focus: Actions Panel"
			case FocusDice:
				m.message = "Focus: Dice Roller"
			}
			return m, nil

		// Focus cycling backwards - Shift+P
		case "P":
			m.focusArea = (m.focusArea - 1 + 4) % 4
			switch m.focusArea {
			case FocusMain:
				m.message = "Focus: Main Panel"
			case FocusCharStats:
				m.message = "Focus: Character Stats"
			case FocusActions:
				m.message = "Focus: Actions Panel"
			case FocusDice:
				m.message = "Focus: Dice Roller"
			}
			return m, nil
		}

		// Check if stat generator is active first (BEFORE tab handling)
		if m.statGenerator.IsVisible() {
			return m.handleStatGeneratorKeys(msg)
		}

		// Check if ability roller is active (BEFORE tab handling)
		if m.abilityRoller.IsVisible() {
			return m.handleAbilityRollerKeys(msg)
		}

		// Check if rest popup is active (BEFORE tab handling)
		if m.restPopup.IsVisible() {
			return m.handleRestPopupKeys(msg)
		}

		// Check if message popup is active
		if m.messagePopup.IsVisible() {
			return m.handleMessagePopupKeys(msg)
		}

		// Check all popups/selectors BEFORE panel navigation (so Tab works in popups)
		// Check if spell selector is active
		if m.spellSelector.IsVisible() {
			return m.handleSpellSelectorKeys(msg)
		}

		// Check if feat selector is active
		if m.featSelector.IsVisible() {
			return m.handleFeatSelectorKeys(msg)
		}

		// Check if feat detail popup is active
		if m.featDetailPopup.IsVisible() {
			if handlers.HandleFeatDetailPopupKeys(msg, m.featDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if mastery detail popup is active
		if m.masteryDetailPopup.IsVisible() {
			if handlers.HandleMasteryDetailPopupKeys(msg, m.masteryDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if maneuver detail popup is active
		if m.maneuverDetailPopup.IsVisible() {
			if handlers.HandleManeuverDetailPopupKeys(msg, m.maneuverDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if consumable detail popup is active
		if m.consumableDetailPopup.IsVisible() {
			if handlers.HandleConsumableDetailPopupKeys(msg, m.consumableDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if feature detail popup is active
		if m.featureDetailPopup.IsVisible() {
			if handlers.HandleFeatureDetailPopupKeys(msg, m.featureDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if item detail popup is active
		if m.itemDetailPopup.IsVisible() {
			if handlers.HandleItemDetailPopupKeys(msg, m.itemDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if spell detail popup is active
		if m.spellDetailPopup.IsVisible() {
			if handlers.HandleSpellDetailPopupKeys(msg, m.spellDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if origin selector is active
		if m.originSelector.IsVisible() {
			return m.handleOriginSelectorKeys(msg)
		}

		// Check if alignment selector is active
		if m.alignmentSelector.IsVisible() {
			return m.handleAlignmentSelectorKeys(msg)
		}

		// Check if origin detail popup is active
		if m.originDetailPopup.IsVisible() {
			if handlers.HandleOriginDetailPopupKeys(msg, m.originDetailPopup, &m.message) {
				return m, nil
			}
		}

		// Check if input popup is active
		if m.inputPopup.IsVisible() {
			return m.handleInputPopupKeys(msg)
		}

		// Check if ability choice selector is active (for feat ability choices)
		if m.abilityChoiceSelector.IsVisible() {
			return m.handleAbilityChoiceSelectorKeys(msg)
		}

		// Check if subtype selector is active
		if m.subtypeSelector.IsVisible() {
			return m.handleSubtypeSelectorKeys(msg)
		}

		// Check if skill selector is active
		if m.skillSelector.IsVisible() {
			return m.handleSkillSelectorKeys(msg)
		}

		// Check if language selector is active
		if m.languageSelector.IsVisible() {
			return m.handleLanguageSelectorKeys(msg)
		}

		// Check if tool selector is active
		if m.toolSelector.IsVisible() {
			return m.handleToolSelectorKeys(msg)
		}

		// Check if weapon mastery selector is active
		if m.weaponMasterySelector.IsVisible() {
			return m.handleWeaponMasterySelectorKeys(msg)
		}
		if m.expertiseSelector.IsVisible() {
			return m.handleExpertiseSelectorKeys(msg)
		}

		// Check if maneuver selector is active
		if m.maneuverSelector.IsVisible() {
			return m.handleManeuverSelectorKeys(msg)
		}

		// Check if level-up selector is active
		if m.levelUpSelector.IsVisible() {
			return m.handleLevelUpSelectorKeys(msg)
		}

		// Check if de-level selector is active
		if m.deLevelSelector.IsVisible() {
			return m.handleDeLevelSelectorKeys(msg)
		}

		// Check if item selector is active
		if m.itemSelector.IsVisible() {
			return m.handleItemSelectorKeys(msg)
		}

		// Check if divine order selector is active (for Cleric level 1, before cantrip selection)
		// But allow tab navigation to pass through
		if m.divineOrderSelectorVisible {
			// Check if Divine Order has already been applied (in case selector wasn't properly hidden)
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
				m.divineOrderSelectorVisible = false
				m.pendingDivineOrder = ""
				// Continue to tab navigation below
			} else {
				// Allow tab navigation to pass through even when selector is visible
				if msg.String() == "tab" || msg.String() == "shift+tab" {
					// Tab navigation - let it pass through to panel navigation check below
					// Don't intercept tab keys
				} else {
					// Handle other keys in selector
					return m.handleDivineOrderSelectorKeys(msg)
				}
			}
		}

		// Check if fighting style selector is active (highest priority in class flow)
		if m.fightingStyleSelector.IsVisible() {
			return m.handleFightingStyleSelectorKeys(msg)
		}

		// Check if cantrip selector is active
		if m.cantripSelector.IsVisible() {
			return m.handleCantripSelectorKeys(msg)
		}

		// Check if leveled spell selector is active
		if m.leveledSpellSelector.IsVisible() {
			return m.handleLeveledSpellSelectorKeys(msg)
		}

		// Check if school spell selector is active
		if m.schoolSpellSelector.IsVisible() {
			return m.handleSchoolSpellSelectorKeys(msg)
		}

		// Check if spellbook editor is active
		if m.spellbookEditor.IsVisible() {
			return m.handleSpellbookEditorKeys(msg)
		}

		// Check if spell prep selector is active
		if m.spellPrepSelector.IsVisible() {
			return m.handleSpellPrepSelectorKeys(msg)
		}

		// Check if slot restorer is active
		if m.slotRestorer.IsVisible() {
			return m.handleSlotRestorerKeys(msg)
		}

		// Check if class skill selector is active
		if m.classSkillSelector.IsVisible() {
			return m.handleClassSkillSelectorKeys(msg)
		}

		// Check if subclass selector is active
		if m.subclassSelector.IsVisible() {
			return m.handleSubclassSelectorKeys(msg)
		}

		// Check if class selector is active
		if m.classSelector.IsVisible() {
			return m.handleClassSelectorKeys(msg)
		}

		// Check if species selector is active
		if m.speciesSelector.IsVisible() {
			return m.handleSpeciesSelectorKeys(msg)
		}

		// Panel navigation (AFTER all popups, only when focused on main and no popups active)
		switch msg.String() {
		case "tab":
			if m.focusArea == FocusMain {
				m.tabs.Next()
				m.currentPanel = PanelType(m.tabs.SelectedIndex)
			}
			return m, nil

		case "shift+tab":
			if m.focusArea == FocusMain {
				m.tabs.Prev()
				m.currentPanel = PanelType(m.tabs.SelectedIndex)
			}
			return m, nil
		}

		// Handle input based on current focus
		debug.Log("Update: Handling key in focusArea=%d (0=Main,1=CharStats,2=Actions,3=Dice)", m.focusArea)
		switch m.focusArea {
		case FocusMain:
			return m.handleMainPanelKeys(msg)
		case FocusCharStats:
			return m.handleCharStatsPanelKeys(msg)
		case FocusActions:
			return m.handleActionsPanelKeys(msg)
		case FocusDice:
			return m.handleDicePanelKeys(msg)
		}

		// Global key 'R' for rest (affects actions and spells)
		switch msg.String() {
		case "R": // Shift+R for rest
			m.character.LongRest()
			m.message = "Long rest completed! HP, spells, and abilities restored."
			return m, nil
		}
	}

	return m, nil
}

// handleMainPanelKeys handles keys when main panel has focus
func (m *Model) handleMainPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleMainPanelKeys: key=%s, currentPanel=%d (0=Stats,1=Skills,2=Inv,3=Spells,4=Features,5=Traits,6=Origin)", msg.String(), m.currentPanel)
	switch m.currentPanel {
	case StatsPanel:
		return m.handleStatsPanel(msg)
	case SkillsPanel:
		return m.handleSkillsPanel(msg)
	case InventoryPanel:
		return m.handleInventoryPanel(msg)
	case SpellsPanel:
		return m.handleSpellsPanel(msg)
	case FeaturesPanel:
		return m.handleFeaturesPanel(msg)
	case TraitsPanel:
		return m.handleTraitsPanel(msg)
	case OriginPanel:
		return m.handleOriginPanel(msg)
	}
	return m, nil
}

// handleStatsPanel handles stats panel specific keys
func (m *Model) handleStatsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.statsPanel.Prev()
	case "down", "j":
		m.statsPanel.Next()
	case "e":
		// Go directly to extras/modifier editing
		m.statGenerator.ShowExtrasOnly(&m.character.AbilityScores)
		m.message = "Edit ability modifiers..."
	case "r":
		// Open full stat generator for rolling/assigning stats
		m.statGenerator.Show(&m.character.AbilityScores)
		m.message = "Generate ability scores..."
	case "t":
		// Roll saving throw for selected ability
		selectedAbility := m.statsPanel.GetSelectedAbility()
		m.rollSavingThrow(selectedAbility)
	case "a":
		// Roll ability check for selected ability
		selectedAbility := m.statsPanel.GetSelectedAbility()
		m.rollAbilityCheck(selectedAbility)
	}
	return m, nil
}

// rollSavingThrow rolls a saving throw for the given ability
func (m *Model) rollSavingThrow(ability models.AbilityType) {
	char := m.character
	modifier := char.AbilityScores.GetModifier(ability)

	// Check if proficient in this saving throw
	isProficient := false
	abilityFullName := ""
	switch ability {
	case models.Strength:
		abilityFullName = "Strength"
	case models.Dexterity:
		abilityFullName = "Dexterity"
	case models.Constitution:
		abilityFullName = "Constitution"
	case models.Intelligence:
		abilityFullName = "Intelligence"
	case models.Wisdom:
		abilityFullName = "Wisdom"
	case models.Charisma:
		abilityFullName = "Charisma"
	}

	for _, prof := range char.SavingThrowProficiencies {
		if strings.EqualFold(prof, abilityFullName) {
			isProficient = true
			break
		}
	}

	// Add proficiency bonus if proficient
	if isProficient {
		modifier += char.ProficiencyBonus
	}

	// Roll 1d20 + modifier
	expression := fmt.Sprintf("1d20%+d", modifier)
	m.dicePanel.Roll(expression)

	profStr := ""
	if isProficient {
		profStr = " (proficient)"
	}
	m.message = fmt.Sprintf("Rolled %s saving throw%s: %s", abilityFullName, profStr, expression)
}

// rollAbilityCheck rolls an ability check for the given ability
func (m *Model) rollAbilityCheck(ability models.AbilityType) {
	char := m.character
	modifier := char.AbilityScores.GetModifier(ability)

	abilityFullName := ""
	switch ability {
	case models.Strength:
		abilityFullName = "Strength"
	case models.Dexterity:
		abilityFullName = "Dexterity"
	case models.Constitution:
		abilityFullName = "Constitution"
	case models.Intelligence:
		abilityFullName = "Intelligence"
	case models.Wisdom:
		abilityFullName = "Wisdom"
	case models.Charisma:
		abilityFullName = "Charisma"
	}

	// Roll 1d20 + modifier (no proficiency for raw ability checks)
	expression := fmt.Sprintf("1d20%+d", modifier)
	m.dicePanel.Roll(expression)

	m.message = fmt.Sprintf("Rolled %s ability check: %s", abilityFullName, expression)
}

// handleActionsPanelKeys handles keys when actions panel has focus
func (m *Model) handleActionsPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if attack menu is active
	if m.attackMenu.IsVisible() {
		return m.handleAttackMenuKeys(msg)
	}

	switch msg.String() {
	case "up", "k":
		m.actionsPanel.Prev()
	case "down", "j":
		m.actionsPanel.Next()
	case "r":
		// Roll normal attack
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "normal")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "a":
		// Roll attack with advantage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "advantage")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "x":
		// Roll attack with disadvantage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollAttackDirect(attack, "disadvantage")
				m.dicePanel.LastMessage = result
				m.message = result
			}
		}
	case "d":
		// Roll damage
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				result := m.rollDamageDirect(attack)
				m.dicePanel.LastMessage = result
				m.message = result
			}
		} else if m.actionsPanel.IsSneakAttackSelected() {
			// Roll Sneak Attack damage
			sneakAttackDice := m.actionsPanel.GetSneakAttackDice()
			m.dicePanel.Roll(sneakAttackDice)
			m.message = fmt.Sprintf("Sneak Attack damage: %s", m.dicePanel.LastMessage)
		}
	case "enter":
		// Show attack menu if attack is selected
		if m.actionsPanel.IsAttackSelected() {
			attack := m.actionsPanel.GetSelectedAttack()
			if attack != nil {
				debug.Log("Opening attack menu for: %s", attack.Name)
				m.attackMenu.Show(attack)
				m.message = "Select attack option..."
			} else {
				debug.Log("No attack selected (attack is nil)")
			}
		} else if m.actionsPanel.IsSneakAttackSelected() {
			// Roll Sneak Attack damage directly
			sneakAttackDice := m.actionsPanel.GetSneakAttackDice()
			m.dicePanel.Roll(sneakAttackDice)
			m.message = fmt.Sprintf("Sneak Attack damage: %s", m.dicePanel.LastMessage)
		} else if m.actionsPanel.IsSpellSelected() {
			// Cast spell if spell is selected
			spell := m.actionsPanel.GetSelectedSpell()
			if spell != nil {
				debug.Log("Attempting to cast spell: %s", spell.Name)
				msg, success := m.actionsPanel.CastSelectedSpell()
				if success {
					m.storage.Save(m.character)
					m.dicePanel.LastMessage = msg
					m.message = msg
					debug.Log("Spell cast successfully: %s", spell.Name)
				} else {
					m.message = msg
					debug.Log("Failed to cast spell: %s", msg)
				}
			}
		} else {
			// For non-attack/spell actions
			m.message = "Other actions not fully implemented"
		}
	}
	return m, nil
}

// handleAttackMenuKeys handles keys when attack menu is visible
func (m *Model) handleAttackMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleAttackMenuKeys: key=%s", msg.String())

	switch msg.String() {
	case "up", "k":
		m.attackMenu.Prev()
		debug.Log("Attack menu: moved up")
	case "down", "j":
		m.attackMenu.Next()
		debug.Log("Attack menu: moved down")
	case "enter":
		option := m.attackMenu.GetSelectedOption()
		attack := m.attackMenu.GetAttack()
		debug.Log("Attack menu: enter pressed, option=%s, attack=%v", option, attack != nil)

		if attack != nil {
			var result string
			switch {
			case option == "Attack with Advantage":
				result = m.rollAttackDirect(attack, "advantage")
			case option == "Attack with Disadvantage":
				result = m.rollAttackDirect(attack, "disadvantage")
			case option == "Attack (Normal)":
				result = m.rollAttackDirect(attack, "normal")
			case strings.HasPrefix(option, "1-Hand Damage"):
				result = m.rollDamageDirect(attack)
			case strings.HasPrefix(option, "1-Hand Critical"):
				result = m.rollCriticalDamage(attack, attack.DamageDice)
			case strings.HasPrefix(option, "2-Hands Damage"):
				result = m.rollVersatileDamage(attack)
			case strings.HasPrefix(option, "2-Hands Critical"):
				result = m.rollCriticalDamage(attack, attack.VersatileDamage)
			case strings.HasPrefix(option, "Damage"):
				result = m.rollDamageDirect(attack)
			case strings.HasPrefix(option, "Critical Hit"):
				result = m.rollCriticalDamage(attack, attack.DamageDice)
			default:
				result = fmt.Sprintf("Unknown option: %s", option)
			}
			debug.Log("Attack result: %s", result)
			m.dicePanel.LastMessage = result
			m.message = result
		}
		m.attackMenu.Hide()
		debug.Log("Attack menu hidden")
	case "esc":
		debug.Log("Attack menu: cancelled")
		m.attackMenu.Hide()
		m.message = ""
	}
	return m, nil
}

// rollAttackDirect performs an attack roll directly without popup
func (m *Model) rollAttackDirect(attack *models.Attack, rollType string) string {
	var diceRollType dice.RollType
	var advantageStr string

	switch rollType {
	case "advantage":
		diceRollType = dice.Advantage
		advantageStr = "Advantage"
	case "disadvantage":
		diceRollType = dice.Disadvantage
		advantageStr = "Disadvantage"
	default:
		diceRollType = dice.Normal
		advantageStr = ""
	}

	result, err := dice.Roll("1d20", diceRollType)
	if err != nil {
		return fmt.Sprintf("Error rolling: %v", err)
	}

	roll := 0
	if len(result.Rolls) > 0 {
		roll = result.Rolls[0]
	}

	total := roll + attack.AttackBonus

	// Check for critical hit based on character's critical range
	critRange := m.character.GetCriticalRange()
	if roll >= critRange {
		critText := "CRITICAL HIT!"
		if critRange < 20 {
			critText = "CRITICAL HIT! (19-20 range)"
		}
		return fmt.Sprintf("%s: %s [%d] + %d = %d %s",
			attack.Name, critText, roll, attack.AttackBonus, total, advantageStr)
	}

	return attack.FormatAttackRoll(roll, total, advantageStr)
}

// rollDamageDirect performs a damage roll directly without popup
func (m *Model) rollDamageDirect(attack *models.Attack) string {
	result, err := dice.Roll(attack.DamageDice, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling damage: %v", err)
	}

	total := result.Total + attack.DamageBonus
	return attack.FormatDamageRoll(result.Rolls, total)
}

// rollVersatileDamage performs a damage roll with two-handed versatile damage
func (m *Model) rollVersatileDamage(attack *models.Attack) string {
	result, err := dice.Roll(attack.VersatileDamage, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling damage: %v", err)
	}

	// Use TwoHandDamageBonus for two-handed attacks (no Dueling bonus)
	total := result.Total + attack.TwoHandDamageBonus
	return fmt.Sprintf("%s (2-Hands): Damage = %v +%d = %d %s",
		attack.Name, result.Rolls, attack.TwoHandDamageBonus, total, attack.DamageType)
}

// rollCriticalDamage performs a critical hit damage roll (double dice)
func (m *Model) rollCriticalDamage(attack *models.Attack, damageDice string) string {
	// Double the damage dice for critical hits
	// Parse dice notation (e.g., "1d8" -> "2d8", "2d6" -> "4d6")
	critDice := damageDice

	// Simple parsing: if it starts with a number, double it
	parts := strings.Split(damageDice, "d")
	if len(parts) == 2 {
		numDice := 1
		if parts[0] != "" {
			if n, err := fmt.Sscanf(parts[0], "%d", &numDice); err == nil && n == 1 {
				critDice = fmt.Sprintf("%dd%s", numDice*2, parts[1])
			}
		} else {
			// "d8" format, assume 1d8
			critDice = fmt.Sprintf("2d%s", parts[1])
		}
	}

	result, err := dice.Roll(critDice, dice.Normal)
	if err != nil {
		return fmt.Sprintf("Error rolling critical damage: %v", err)
	}

	// Determine which damage bonus to use
	damageBonus := attack.DamageBonus
	label := "Critical Hit"

	// If this is a two-handed critical (versatile weapon), use TwoHandDamageBonus
	if attack.VersatileDamage != "" && damageDice == attack.VersatileDamage {
		label = "Critical Hit (2-Hands)"
		damageBonus = attack.TwoHandDamageBonus
	}

	total := result.Total + damageBonus

	return fmt.Sprintf("%s %s: Damage = %v +%d = %d %s",
		attack.Name, label, result.Rolls, damageBonus, total, attack.DamageType)
}

// handleAttackRollerKeys handles keys when attack roller is visible
func (m *Model) handleAttackRollerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	state := m.attackRoller.GetState()

	switch state {
	case "select_attack":
		switch msg.String() {
		case "up", "k":
			m.attackRoller.Prev()
		case "down", "j":
			m.attackRoller.Next()
		case "enter":
			m.attackRoller.SelectAttack()
			m.message = "Choose action: 'a'ttack, 'd'amage, ad'v'antage, disadvantage (x)"
		case "esc":
			m.attackRoller.Hide()
			m.message = ""
		}
	case "select_roll_type":
		switch msg.String() {
		case "a":
			// Roll normal attack
			m.attackRoller.SetRollType("normal")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "d":
			// Roll damage
			result := m.attackRoller.RollDamage()
			m.dicePanel.LastMessage = result
			m.message = result
		case "v":
			// Roll attack with advantage
			m.attackRoller.SetRollType("advantage")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "x":
			// Roll attack with disadvantage
			m.attackRoller.SetRollType("disadvantage")
			result := m.attackRoller.RollAttack()
			m.dicePanel.LastMessage = result
			m.message = result
		case "esc":
			m.attackRoller.Hide()
			m.message = ""
		}
	}

	return m, nil
}

// handleDicePanelKeys handles keys when dice panel has focus
func (m *Model) handleDicePanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	mode := m.dicePanel.GetMode()

	switch mode {
	case panels.DiceModeIdle:
		// Idle mode - waiting for user to choose action
		switch msg.String() {
		case "enter":
			m.dicePanel.SetMode(panels.DiceModeInput)
			m.message = "Enter dice notation and press Enter to roll"
		case "h":
			m.dicePanel.SetMode(panels.DiceModeHistory)
			m.message = "Navigate history with ↑/↓, press Enter to reroll"
		case "r":
			m.dicePanel.RerollLast()
			m.message = "Rerolled last dice"
		}
		return m, nil

	case panels.DiceModeInput:
		// Input mode - typing dice notation
		switch msg.String() {
		case "esc":
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = ""
			return m, nil
		case "enter":
			if m.dicePanel.GetInput() != "" {
				m.dicePanel.Roll(m.dicePanel.GetInput())
				m.dicePanel.SetMode(panels.DiceModeIdle)
				m.message = ""
			}
			return m, nil
		}
		// Pass all other keys to input
		return m, m.dicePanel.Update(msg)

	case panels.DiceModeHistory:
		// History mode - browsing previous rolls
		switch msg.String() {
		case "esc":
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = ""
		case "up", "k":
			m.dicePanel.HistoryPrev()
		case "down", "j":
			m.dicePanel.HistoryNext()
		case "enter":
			m.dicePanel.RerollSelected()
			m.dicePanel.SetMode(panels.DiceModeIdle)
			m.message = "Rerolled selected dice"
		}
		return m, nil
	}

	return m, nil
}

// handleSkillsPanel handles skills panel specific keys
func (m *Model) handleSkillsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.skillsPanel.Prev()
	case "down", "j":
		m.skillsPanel.Next()
	case "e":
		m.skillsPanel.ToggleProficiency()
		m.message = "Proficiency toggled"
	case "r":
		if skill := m.skillsPanel.GetSelectedSkill(); skill != nil {
			abilityMod := m.character.AbilityScores.GetModifier(skill.Ability)
			bonus := skill.CalculateBonus(abilityMod, m.character.ProficiencyBonus)
			expr := fmt.Sprintf("1d20%+d", bonus)
			m.dicePanel.Roll(expr)
			m.message = fmt.Sprintf("Rolling %s: %s", skill.Name, m.dicePanel.LastMessage)
		}
	}
	return m, nil
}

// handleInventoryPanel handles inventory panel specific keys
func (m *Model) handleInventoryPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.inventoryPanel.Prev()
	case "down", "j":
		m.inventoryPanel.Next()
	case "enter":
		// Show item details
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			m.itemDetailPopup.Show(item)
			m.message = "Viewing item details..."
		}
	case "e":
		// Toggle equipped status for selected item
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			// Check if item is equippable
			def := models.GetItemDefinitionByName(item.Name)
			if def != nil && models.IsEquippable(*def) {
				// Only check proficiency when EQUIPPING (not unequipping)
				if !item.Equipped {
					// Check armor proficiency
					if item.Type == models.Armor {
						// Get armor subcategory (Light, Medium, Heavy, Shield)
						armorType := def.Subcategory

						// Check proficiency
						if !models.HasArmorProficiency(m.character, armorType) {
							m.message = fmt.Sprintf("Cannot equip %s: Not proficient with %s armor!", item.Name, armorType)
							return m, nil
						}

						// Unequip other armor pieces first
						models.UnequipOtherArmor(m.character, item)
					}

					// Check weapon proficiency
					if item.Type == models.Weapon {
						// Get weapon subcategory (simple melee, martial melee, etc.)
						weaponType := def.Subcategory

						// Check proficiency
						if !models.HasWeaponProficiency(m.character, weaponType) {
							m.message = fmt.Sprintf("Cannot equip %s: Not proficient with %s weapons!", item.Name, weaponType)
							return m, nil
						}
					}
				}

		m.inventoryPanel.ToggleEquipped()

				// Recalculate AC after equipping/unequipping
				m.character.UpdateDerivedStats()

				if item.Equipped {
					m.message = fmt.Sprintf("%s equipped (AC: %d)", item.Name, m.character.AC)
				} else {
					m.message = fmt.Sprintf("%s unequipped (AC: %d)", item.Name, m.character.AC)
				}
				m.storage.Save(m.character)
			} else {
				m.message = "This item cannot be equipped"
			}
		}
	case "d":
		// Delete selected item (decrease quantity by 1 or remove if quantity is 1)
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			wasEquipped := item.Equipped
			itemType := item.Type
			if item.Quantity > 1 {
				item.Quantity--
				m.message = fmt.Sprintf("%s quantity decreased to %d", item.Name, item.Quantity)
			} else {
				itemName := item.Name
		m.inventoryPanel.DeleteSelected()
				m.message = fmt.Sprintf("%s removed from inventory", itemName)
			}

			// Recalculate AC if armor was equipped
			if wasEquipped && itemType == models.Armor {
				m.character.UpdateDerivedStats()
				m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
			}
			m.storage.Save(m.character)
		}
	case "D":
		// Delete all of selected item
		item := m.inventoryPanel.GetSelectedItem()
		if item != nil {
			itemName := item.Name
			wasEquipped := item.Equipped
			itemType := item.Type
			m.inventoryPanel.DeleteSelected()
			m.message = fmt.Sprintf("All %s removed from inventory", itemName)

			// Recalculate AC if armor was equipped
			if wasEquipped && itemType == models.Armor {
				m.character.UpdateDerivedStats()
				m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
			}
			m.storage.Save(m.character)
		}
	case "a":
		// Open item selector to add items
		m.itemSelector.Show(m.character)
		m.message = "Select item category..."
	}
	return m, nil
}


// handleSpellsPanel handles spells panel specific keys
func (m *Model) handleSpellsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.spellsPanel.HandleKey(msg)
	case "down", "j":
		m.spellsPanel.HandleKey(msg)
	case "pgup":
		m.spellsPanel.HandleKey(msg)
	case "pgdown":
		m.spellsPanel.HandleKey(msg)
	case "r":
		// Rest - restore all spell slots
		m.spellsPanel.Rest()
		m.storage.Save(m.character)
		m.message = "Spell slots restored!"
	case "c":
		// Change cantrips
		if m.character.SpellBook.IsPreparedCaster {
			m.cantripSelector.Show(m.character.Class, m.character.SpellBook.CantripsKnown)
			m.message = fmt.Sprintf("Select %d cantrips...", m.character.SpellBook.CantripsKnown)
		} else {
			m.message = "Only prepared casters can change cantrips this way"
		}
	case "v", "b":
		// Open spellbook editor (for Wizards) or spell prep selector (for other prepared casters)
		debug.Log("=== KEY 'v' or 'b' PRESSED IN SPELLS PANEL")
		debug.Log("=== IsSpellbookCaster: %v, HasWizard: %v, IsPreparedCaster: %v",
			m.character.SpellBook.IsSpellbookCaster,
			m.character.HasClass("Wizard"),
			m.character.SpellBook.IsPreparedCaster)
		debug.Log("=== Current spells in spellbook: %d", len(m.character.SpellBook.Spells))

		if m.character.SpellBook.IsSpellbookCaster && (m.character.HasClass("Wizard") || m.character.IsArcaneTrickster()) {
			debug.Log("=== OPENING SPELLBOOK EDITOR")
			m.spellbookEditor.Show()
			m.message = "Managing spellbook... (Space: Prepare | a: Add | d/x: Remove | c: Cantrips | 0-9: Filter)"
			debug.Log("=== Spellbook editor visible: %v", m.spellbookEditor.IsVisible())
		} else if m.character.SpellBook.IsPreparedCaster || m.character.HasClass("Paladin") {
			// Paladins are prepared casters (like Clerics)
			// Check if Paladin has spellcasting (level 1+)
			if m.character.HasClass("Paladin") && m.character.GetClassLevel("Paladin") >= 1 {
				// Ensure spellcasting is initialized
				if !m.character.SpellBook.IsPreparedCaster {
					// Re-initialize spellcasting if needed
					class := models.GetClassByName("Paladin")
					if class != nil {
						models.InitializeSpellcasting(m.character, class)
						m.character.UpdateDerivedStats() // Recalculate prepared spell limit
					}
				}
			}
			debug.Log("=== OPENING SPELL PREP SELECTOR")
			m.spellPrepSelector.Show()
			m.message = "Managing spellbook... (Tab: Switch tabs • Space: Prepare/Add • a/d: Add/Remove cantrips)"
		} else {
			debug.Log("=== Not a prepared caster - showing error message")
			m.message = "Only prepared casters can prepare spells"
		}
	case "s":
		// Open slot restorer
		if m.character.SpellBook.SpellcastingMod != "" {
			m.slotRestorer.Show()
			m.message = "Select spell slot to restore..."
		} else {
			m.message = "Not a spellcaster"
		}
	case "a":
		// Add new spell to spellbook (Wizard only)
		debug.Log("=== KEY 'a' PRESSED IN SPELLS PANEL")
		debug.Log("=== IsSpellbookCaster: %v, HasWizard: %v",
			m.character.SpellBook.IsSpellbookCaster,
			m.character.HasClass("Wizard"))

		if m.character.SpellBook.IsSpellbookCaster && (m.character.HasClass("Wizard") || m.character.IsArcaneTrickster()) {
			var maxSpellLevel int
			if m.character.HasClass("Wizard") {
				wizardLevel := m.character.GetClassLevel("Wizard")
				maxSpellLevel = (wizardLevel + 1) / 2
			} else if m.character.IsArcaneTrickster() {
				rogueLevel := m.character.GetRogueLevel()
				maxSpellLevel = (rogueLevel + 2) / 3 // Third caster progression
			}
			if maxSpellLevel > 9 {
				maxSpellLevel = 9
			}
			debug.Log("=== SHOWING SPELLBOOK EDITOR: MaxSpellLevel=%d", maxSpellLevel)
			// Open spellbook editor to add spells
			m.spellbookEditor.Show()
			m.message = fmt.Sprintf("Add new spells to your spellbook (up to level %d)...", maxSpellLevel)
			debug.Log("=== Spellbook editor visible: %v", m.spellbookEditor.IsVisible())
		} else {
			debug.Log("=== Not a spellbook caster")
			m.message = "Only spellbook casters can add spells to their spellbook"
		}
	case "enter":
		// View spell details
		spell := m.spellsPanel.GetSelectedSpell()
		if spell != nil {
			m.spellDetailPopup.Show(*spell)
			m.message = "Viewing spell details..."
		}
	case "u":
		// Consume spell slot
		if m.spellsPanel.ConsumeSpellSlot() {
			m.storage.Save(m.character)
			m.message = "Spell slot consumed"
		} else {
			m.message = "Cannot consume spell slot (already at 0 or no slot selected)"
		}
	case "U":
		// Restore spell slot
		if m.spellsPanel.RestoreSpellSlot() {
			m.storage.Save(m.character)
			m.message = "Spell slot restored"
		} else {
			m.message = "Cannot restore spell slot (already at max or no slot selected)"
		}
	}
	return m, nil
}

// handleFeaturesPanel handles features panel specific keys
func (m *Model) handleFeaturesPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.featuresPanel.Prev()
	case "down", "j":
		m.featuresPanel.Next()
	case "ctrl+u", "pgup":
		m.featuresPanel.PageUp()
	case "ctrl+d", "pgdown":
		m.featuresPanel.PageDown()
	case "ctrl+y":
		m.featuresPanel.ScrollUp()
	case "ctrl+e":
		m.featuresPanel.ScrollDown()
	case "enter":
		// Show popup for selected item (any type)
		selectedItem := m.featuresPanel.GetSelectedItem()
		if selectedItem != nil {
			switch selectedItem.ItemType {
			case "consumable":
				// Show consumable detail popup
				item := selectedItem.Consumable
				if item != nil {
					restTypeStr := "Unknown"
					switch item.RestType {
					case models.ShortRest:
						restTypeStr = "Short Rest"
					case models.LongRest:
						restTypeStr = "Long Rest"
					case models.Daily:
						restTypeStr = "Daily"
					}
					m.consumableDetailPopup.Show(item.Name, item.Current, item.Max, restTypeStr, item.Description)
					m.message = "Viewing details..."
				}
			case "passive", "rage_effect":
				// Show feature detail popup
				usesStr := ""
				restTypeStr := ""

				// For passive features, show uses info if applicable
				if selectedItem.Feature != nil {
					if selectedItem.Feature.MaxUses > 0 {
						usesStr = fmt.Sprintf("%d/%d", selectedItem.Feature.CurrentUses, selectedItem.Feature.MaxUses)
					}
					restTypeStr = string(selectedItem.Feature.RestType)
				}

				m.featureDetailPopup.Show(selectedItem.Name, selectedItem.Description, usesStr, restTypeStr)
				m.message = "Viewing feature details..."
			}
		} else {
			m.message = "No item selected"
		}
		return m, nil
	case "u":
		// Consume selected item
		item := m.featuresPanel.GetSelectedConsumable()
		if item == nil {
			m.message = "No item selected"
			return m, nil
		}

		if item.Current <= 0 {
			m.message = fmt.Sprintf("%s has no uses remaining", item.Name)
			return m, nil
		}

		// Handle by type
		switch item.ItemType {
		case "resource":
			switch item.ResourceType {
			case "focus_points":
				monk := m.character.GetMonkMechanics()
				if monk.SpendFocusPoint(1) {
					current, max := monk.GetFocusPoints()
					m.message = fmt.Sprintf("Focus Point spent. Current: %d/%d", current, max)
				} else {
					m.message = "Not enough Focus Points"
				}
			case "psi_dice":
				m.character.PsiDice.Current--
				m.message = fmt.Sprintf("Psi Die spent. Current: %d/%d", m.character.PsiDice.Current, m.character.PsiDice.Max)
			case "superiority_dice":
				m.character.SuperiorityDice.Current--
				m.message = fmt.Sprintf("Superiority Die spent. Current: %d/%d", m.character.SuperiorityDice.Current, m.character.SuperiorityDice.Max)
			case "warrior_dice":
				m.character.WarriorDice.Current--
				m.message = fmt.Sprintf("Warrior Die spent. Current: %d/%d", m.character.WarriorDice.Current, m.character.WarriorDice.Max)
			case "soulknife_psi_dice":
				m.character.SoulknifePsiDice.Current--
				m.message = fmt.Sprintf("Psionic Energy die spent. Current: %d/%d", m.character.SoulknifePsiDice.Current, m.character.SoulknifePsiDice.Max)
			}
			m.storage.Save(m.character)
		case "feature":
			// Handle feature consumption with special effects
			if item.Feature != nil {
				switch item.Feature.Name {
				case "Second Wind":
					// Roll 1d10 + fighter level
					result, err := dice.Roll("1d10", dice.Normal)
					if err == nil {
						fighterLevel := m.character.GetFighterLevel()
						healing := result.Total + fighterLevel
						m.character.CurrentHP += healing
						if m.character.CurrentHP > m.character.MaxHP {
							m.character.CurrentHP = m.character.MaxHP
						}
						m.message = fmt.Sprintf("Second Wind: Healed %d HP (1d10[%d] + %d level)", healing, result.Total, fighterLevel)
					} else {
						fighterLevel := m.character.GetFighterLevel()
						m.message = fmt.Sprintf("Second Wind: Healed %d HP", fighterLevel)
						m.character.CurrentHP += fighterLevel
					}
				default:
					m.message = fmt.Sprintf("%s used", item.Feature.Name)
				}
				// Decrement uses
		m.featuresPanel.UseFeature()
		m.storage.Save(m.character)
			}
		}
		return m, nil
	case "U":
		// Restore selected item
		item := m.featuresPanel.GetSelectedConsumable()
		if item == nil {
			m.message = "No item selected"
			return m, nil
		}

		if item.Current >= item.Max {
			m.message = fmt.Sprintf("%s already at maximum", item.Name)
			return m, nil
		}

		// Handle by type
		switch item.ItemType {
		case "resource":
			switch item.ResourceType {
			case "focus_points":
				monk := m.character.GetMonkMechanics()
				monk.RestoreFocusPoints(1)
				current, max := monk.GetFocusPoints()
				m.message = fmt.Sprintf("Focus Point restored. Current: %d/%d", current, max)
			case "psi_dice":
				m.character.PsiDice.Current++
				m.message = fmt.Sprintf("Psi Die restored. Current: %d/%d", m.character.PsiDice.Current, m.character.PsiDice.Max)
			case "superiority_dice":
				m.character.SuperiorityDice.Current++
				m.message = fmt.Sprintf("Superiority Die restored. Current: %d/%d", m.character.SuperiorityDice.Current, m.character.SuperiorityDice.Max)
			case "warrior_dice":
				m.character.WarriorDice.Current++
				m.message = fmt.Sprintf("Warrior Die restored. Current: %d/%d", m.character.WarriorDice.Current, m.character.WarriorDice.Max)
			case "soulknife_psi_dice":
				m.character.SoulknifePsiDice.Current++
				m.message = fmt.Sprintf("Psionic Energy die restored. Current: %d/%d", m.character.SoulknifePsiDice.Current, m.character.SoulknifePsiDice.Max)
			}
			m.storage.Save(m.character)
		case "feature":
		m.featuresPanel.RestoreFeature()
			m.message = fmt.Sprintf("%s restored", item.Name)
		m.storage.Save(m.character)
		}
		return m, nil
	case "r":
		// Short rest
		m.character.ShortRest()
		m.message = "Short rest completed - features recovered"
		m.storage.Save(m.character)
	case "R":
		// Long rest
		m.character.LongRest()
		m.message = "Long rest completed - all features recovered"
		m.storage.Save(m.character)
	}
	return m, nil
}

// handleOriginPanel handles origin panel specific keys
func (m *Model) handleOriginPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.originPanel.ScrollUp()
	case "down", "j":
		m.originPanel.ScrollDown()
	case "ctrl+u", "pgup":
		m.originPanel.PageUp()
	case "ctrl+d", "pgdown":
		m.originPanel.PageDown()
	case "o":
		// Open origin selector
		m.originSelector.Show(m.character)
		m.message = "Select an origin..."
	case "enter":
		// Show origin details
		if m.character.Origin != "" {
			origin := models.GetOriginByName(m.character.Origin)
			if origin != nil {
				m.originDetailPopup.Show(origin)
			}
		}
	case "a":
		// Open alignment selector
		m.alignmentSelector.Show()
		m.message = "Select an alignment..."
	case "h":
		// Edit height
		m.inputPopup.Show("Edit Height", m.character.Height, "Enter height (e.g., 6'2\")...")
		m.inputPopupContext = "height"
		m.message = "Editing height..."
	case "w":
		// Edit weight
		m.inputPopup.Show("Edit Weight", m.character.Weight, "Enter weight (e.g., 180 lbs)...")
		m.inputPopupContext = "weight"
		m.message = "Editing weight..."
	case "t":
		// Edit personality
		m.traitSelector.Show(models.TraitPersonality, m.character.Personality)
		m.message = "Select or create personality traits (Space to toggle, Enter to confirm)..."
	case "i":
		// Edit ideal
		m.traitSelector.Show(models.TraitIdeal, m.character.Ideal)
		m.message = "Select or create ideals (Space to toggle, Enter to confirm)..."
	case "b":
		// Edit bond
		m.traitSelector.Show(models.TraitBond, m.character.Bond)
		m.message = "Select or create bonds (Space to toggle, Enter to confirm)..."
	case "f":
		// Edit flaw
		m.traitSelector.Show(models.TraitFlaw, m.character.Flaw)
		m.message = "Select or create flaws (Space to toggle, Enter to confirm)..."
	case "s":
		// Edit backstory
		m.backstoryEditor.Show(m.character.Backstory)
		m.message = "Edit your character's backstory..."
	}
	return m, nil
}

// handleTraitsPanel handles traits panel specific keys
func (m *Model) handleTraitsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleTraitsPanel: key=%s", msg.String())
	switch msg.String() {
	case "up", "k":
		// Navigate to previous selectable item
		debug.Log("handleTraitsPanel: calling Prev()")
		m.traitsPanel.Prev()
	case "down", "j":
		// Navigate to next selectable item
		debug.Log("handleTraitsPanel: calling Next()")
		m.traitsPanel.Next()
	case "shift+up":
		// Scroll viewport up
		debug.Log("handleTraitsPanel: scrolling viewport up")
		m.traitsPanel.ScrollUp()
	case "shift+down":
		// Scroll viewport down
		debug.Log("handleTraitsPanel: scrolling viewport down")
		m.traitsPanel.ScrollDown()
	case "ctrl+u", "pgup":
		// Page up
		m.traitsPanel.PageUp()
	case "ctrl+d", "pgdown":
		// Page down
		m.traitsPanel.PageDown()
	case "g":
		// Go to top
		m.traitsPanel.GotoTop()
	case "G":
		// Go to bottom (Shift+G)
		m.traitsPanel.GotoBottom()
	case "enter":
		// Show feat detail popup if on a feat
		if m.traitsPanel.IsOnFeat() {
			featName := m.traitsPanel.GetSelectedFeat()
			if featName != "" {
				m.featDetailPopup.Show(featName, m.character)
				m.message = "Viewing feat details..."
			}
		}
		// Show mastery detail popup if on a weapon mastery
		if m.traitsPanel.IsOnMastery() {
			weaponName, masteryType := m.traitsPanel.GetSelectedMastery()
			if weaponName != "" && masteryType != "" {
				m.masteryDetailPopup.Show(weaponName, masteryType)
				m.message = "Viewing weapon mastery details..."
			}
		}
		// Show maneuver detail popup if on a maneuver
		if m.traitsPanel.IsOnManeuver() {
			maneuverName := m.traitsPanel.GetSelectedManeuver()
			if maneuverName != "" {
				m.maneuverDetailPopup.Show(maneuverName)
				m.message = "Viewing maneuver details..."
			}
		}
	case "l":
		// Add language
		m.languageSelector.SetExcludeLanguages(m.character.Languages)
		m.languageSelector.Show()
		m.message = "Select a language to learn..."
	case "L": // Shift+L
		// Remove language
		if len(m.character.Languages) == 0 {
			m.message = "No languages to remove"
		} else {
			m.languageSelector.ShowForDeletion(m.character.Languages)
			m.message = "Select a language to remove..."
		}
	case "f":
		// Add feat
		m.featSelector.Show(m.character, false) // false = not an origin feat
		m.message = "Select a feat to acquire..."
	case "F": // Shift+F
		// Remove feat
		if len(m.character.Feats) == 0 {
			m.message = "No feats to remove"
		} else {
			m.featSelector.ShowForDeletion(m.character)
			m.message = "Select a feat to remove..."
		}
	case "s":
		// Change fighting style
		if m.character.FightingStyle != "" {
			// Character already has a fighting style, allow changing it
			m.fightingStyleSelector.Show(m.character.Class)
			m.message = "Select a new fighting style..."
		} else {
			m.message = "You don't have a Fighting Style to change"
		}
	case "m":
		// Manage weapon mastery
		debug.Log("handleTraitsPanel: 'm' key pressed - checking weapon mastery")
		// Check if character has weapon mastery feature
		masteryCount := m.getWeaponMasteryCount()
		debug.Log("handleTraitsPanel: masteryCount=%d", masteryCount)
		if masteryCount > 0 {
			debug.Log("handleTraitsPanel: Showing weapon mastery selector")
			m.weaponMasterySelector.Show(masteryCount)
			m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
		} else {
			debug.Log("handleTraitsPanel: No weapon mastery feature found")
			m.message = "You don't have the Weapon Mastery feature"
		}
	case "e":
		// Manage expertise
		debug.Log("handleTraitsPanel: 'e' key pressed - checking expertise")
		expertiseCount := m.getExpertiseCount()
		debug.Log("handleTraitsPanel: expertiseCount=%d", expertiseCount)
		if expertiseCount > 0 {
			debug.Log("handleTraitsPanel: Showing expertise selector")
			m.expertiseSelector.Show(expertiseCount)
			m.message = fmt.Sprintf("Select %d skill(s) for expertise...", expertiseCount)
		} else {
			debug.Log("handleTraitsPanel: No expertise feature found")
			m.message = "You don't have the Expertise feature"
		}
	case "n":
		// Manage Battle Master maneuvers (edit only, no Student of War benefits)
		debug.Log("handleTraitsPanel: 'n' key pressed - checking Battle Master maneuvers")
		// Check if character is a Battle Master
		if m.character.IsBattleMaster() {
			// Get maneuver count from Combat Superiority feature
			maneuverCount := 3 // Default for level 3
			if feature := m.character.GetFeature("Combat Superiority"); feature != nil && feature.Mechanics != nil {
				if count, ok := feature.Mechanics["maneuvers_known"].(float64); ok {
					maneuverCount = int(count)
				}
			}
			debug.Log("handleTraitsPanel: Showing maneuver selector for %d maneuvers", maneuverCount)
			// Clear Student of War flag to prevent prompting for tool/skill
			m.studentOfWarToolSelected = false
			// Show selector first (which clears selections), THEN load current maneuvers
			m.maneuverSelector.Show(maneuverCount)
			m.maneuverSelector.SetSelectedManeuvers(m.character.Maneuvers)
			debug.Log("handleTraitsPanel: Loaded %d existing maneuvers into selector", len(m.character.Maneuvers))
			m.message = fmt.Sprintf("Select up to %d maneuvers...", maneuverCount)
		} else {
			debug.Log("handleTraitsPanel: Not a Battle Master")
			m.message = "Only Battle Masters can learn maneuvers"
		}
	case "d", "x":
		m.traitsPanel.RemoveSelected()
		m.message = "Item removed"
		m.storage.Save(m.character)
	}
	return m, nil
}

// getLevelXP returns the XP required to reach a given level
func getLevelXP(level int) int {
	xpTable := map[int]int{
		1:  0,
		2:  300,
		3:  900,
		4:  2700,
		5:  6500,
		6:  14000,
		7:  23000,
		8:  34000,
		9:  48000,
		10: 64000,
		11: 85000,
		12: 100000,
		13: 120000,
		14: 140000,
		15: 165000,
		16: 195000,
		17: 225000,
		18: 265000,
		19: 305000,
		20: 355000,
	}
	if xp, ok := xpTable[level]; ok {
		return xp
	}
	return 355000 // Max level XP
}

// handleCharStatsPanelKeys handles character stats panel specific keys
func (m *Model) handleCharStatsPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	editMode := m.characterStatsPanel.GetEditMode()

	// If in edit mode, handle save/cancel
	if editMode != panels.CharStatsNormal {
		switch msg.String() {
		case "enter":
			if editMode == panels.CharStatsEditName {
				m.characterStatsPanel.SaveName()
				m.message = "Name updated"
			} else if editMode == panels.CharStatsEditRace {
				m.characterStatsPanel.SaveRace()
				m.message = "Race updated"
			} else if editMode == panels.CharStatsEditHP {
				amount, err := m.characterStatsPanel.SaveHP()
				if err != nil {
					m.message = fmt.Sprintf("Invalid HP value: %v", err)
				} else {
					m.message = fmt.Sprintf("HP adjusted by %+d. Current: %d/%d", amount, m.character.CurrentHP, m.character.MaxHP)
				}
			} else if editMode == panels.CharStatsEditXP {
				amount, err := m.characterStatsPanel.SaveXP()
				if err != nil {
					m.message = fmt.Sprintf("Invalid XP value: %v", err)
				} else {
					nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
					if m.character.Experience >= nextLevelXP {
						m.message = fmt.Sprintf("XP adjusted by %+d. Current: %d (LEVEL UP AVAILABLE!)", amount, m.character.Experience)
					} else {
						m.message = fmt.Sprintf("XP adjusted by %+d. Current: %d (next: %d)", amount, m.character.Experience, nextLevelXP)
					}
				}
			}
			return m, nil
		case "esc":
			m.characterStatsPanel.CancelEdit()
			m.message = "Edit cancelled"
			return m, nil
		default:
			// Pass key to input field
			return m, m.characterStatsPanel.HandleInput(msg)
		}
	}

	// In normal mode, allow class change
	switch msg.String() {
	case "c":
		// Backup current class state before opening selector
		m.pendingChanges.BackupClass(m.character)
		debug.Log("Backed up class state: %s", m.character.Class)
		m.classSelector.Show()
		m.message = "Select a class..."
		return m, nil
	case "l":
		// Open level-up selector (lowercase l)
		m.levelUpSelector.Show()
		m.message = "Level up your character..."
		return m, nil
	case "L":
		// Open de-level selector (Shift+L / uppercase L)
		if m.character.TotalLevel > 1 && len(m.character.Classes) > 0 {
			m.deLevelSelector.Show()
			m.message = "Select class to remove a level from..."
		} else {
			m.message = "Cannot de-level: Character is already at minimum level"
		}
		return m, nil
	case "]", "}":
		// Add 100 XP
		m.character.Experience += 100
		nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
		if m.character.Experience >= nextLevelXP {
			m.message = fmt.Sprintf("XP: %d (LEVEL UP AVAILABLE!)", m.character.Experience)
		} else {
			m.message = fmt.Sprintf("XP: %d (next level: %d)", m.character.Experience, nextLevelXP)
		}
		return m, nil
	case "[", "{":
		// Remove 100 XP (minimum 0)
		if m.character.Experience >= 100 {
			m.character.Experience -= 100
		} else {
			m.character.Experience = 0
		}
		nextLevelXP := getLevelXP(m.character.TotalLevel + 1)
		m.message = fmt.Sprintf("XP: %d (next level: %d)", m.character.Experience, nextLevelXP)
		return m, nil
	}

	// Normal mode - handle actions
	switch msg.String() {
	case "n":
		m.characterStatsPanel.EditName()
		m.message = "Editing name..."
	case "r":
		// Short rest
		m.restPopup.Show(components.ShortRestType)
		m.message = "Taking a short rest..."
	case "R":
		// Long rest
		m.restPopup.Show(components.LongRestType)
		m.message = "Taking a long rest..."
	case "h":
		m.characterStatsPanel.EditHP()
		m.message = "Enter HP change (+/- amount)..."
	case "x":
		m.characterStatsPanel.EditXP()
		m.message = "Enter XP change (+/- amount)..."
	case "+", "=":
		m.characterStatsPanel.AddHP(1)
		m.message = fmt.Sprintf("HP: %d/%d", m.character.CurrentHP, m.character.MaxHP)
	case "-", "_":
		m.characterStatsPanel.RemoveHP(1)
		m.message = fmt.Sprintf("HP: %d/%d", m.character.CurrentHP, m.character.MaxHP)
	case "i":
		// Roll initiative
		initMod := m.characterStatsPanel.GetInitiativeModifier()
		expr := fmt.Sprintf("1d20%+d", initMod)
		m.dicePanel.Roll(expr)
		m.message = fmt.Sprintf("Initiative rolled: %s", m.dicePanel.LastMessage)
	case "I":
		// Toggle inspiration
		m.characterStatsPanel.ToggleInspiration()
		if m.character.Inspiration {
			m.message = "✨ Inspiration gained!"
		} else {
			m.message = "Inspiration used"
		}
		m.storage.Save(m.character)
	}
	return m, nil
}

// handleStatGeneratorKeys handles stat generator specific keys
func (m *Model) handleStatGeneratorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if we're in editing mode for extras
	if m.statGenerator.IsVisible() {
		// Special handling for extras editing mode
		editingExtra := m.statGenerator.IsEditingExtra()

		switch msg.String() {
		case "up", "k":
			if !editingExtra {
				m.statGenerator.Prev()
			}
		case "down", "j":
			if !editingExtra {
				m.statGenerator.Next()
			}
		case "esc":
			// Cancel extra input or go back
			if editingExtra {
				m.statGenerator.CancelExtra()
			} else {
				m.statGenerator.GoBack()
				if !m.statGenerator.IsVisible() {
					m.message = "Stat generation cancelled"
				}
			}
		case "enter":
			// Save extra or continue to next step
			if editingExtra {
				m.statGenerator.SaveExtra()
			} else if m.statGenerator.CanContinue() {
				// Check if we're at the final step and on confirm button
				m.statGenerator.Continue()
				if !m.statGenerator.IsVisible() {
					// Apply stats and close
					m.statGenerator.ApplyToCharacter(m.character)
					m.message = "Ability scores updated!"
					m.storage.Save(m.character)
				}
			} else {
				m.message = "Please assign all stats before continuing"
			}
		case "e":
			// Edit extra in extras state
			if !editingExtra {
				m.statGenerator.StartEditingExtra()
			}
		case "1", "2", "3", "4", "5", "6":
			if !editingExtra {
				// Only assign stats for 4d6 and Standard Array methods
				method := m.statGenerator.GetMethod()
				state := m.statGenerator.GetState()
				if state == components.StateAssignStats &&
					(method == components.Method4d6DropLowest || method == components.MethodStandardArray) {
					idx := int(msg.String()[0] - '1')
					m.statGenerator.ToggleAssignment(idx)
				}
			}
		case "+", "=":
			if !editingExtra {
				// Increase in point buy state or extras
				state := m.statGenerator.GetState()
				if state == components.StateSetExtras {
					m.statGenerator.IncreaseExtra()
				} else {
					m.statGenerator.IncreasePointBuy()
				}
			}
		case "-", "_":
			if !editingExtra {
				// Decrease in point buy state or extras
				state := m.statGenerator.GetState()
				if state == components.StateSetExtras {
					m.statGenerator.DecreaseExtra()
				} else {
					m.statGenerator.DecreasePointBuy()
				}
			}
		case "backspace", "delete":
			// Delete character in extra input
			if editingExtra {
				m.statGenerator.DeleteExtraInput()
			}
		default:
			// Handle typing for extra input
			if editingExtra && len(msg.String()) == 1 {
				char := []rune(msg.String())[0]
				if (char >= '0' && char <= '9') || char == '+' || char == '-' {
					m.statGenerator.HandleExtraInput(char)
				}
			}
		}
	}
	return m, nil
}

// handleAbilityRollerKeys moved to handlers_ability.go

// handleRestPopupKeys and handleMessagePopupKeys moved to handlers_message.go

// handleSpeciesSelectorKeys and handleSubtypeSelectorKeys moved to handlers_species.go
// handleLanguageSelectorKeys moved to handlers_selectors.go

// handleToolSelectorKeys, handleWeaponMasterySelectorKeys, handleExpertiseSelectorKeys,
// handleManeuverSelectorKeys, and handleItemSelectorKeys moved to handlers_selectors.go

// handleLevelUpSelectorKeys and handleDeLevelSelectorKeys moved to handlers_levelup.go

// handleClassSelectorKeys, handleSubclassSelectorKeys, handleClassSkillSelectorKeys,
// and handleFightingStyleSelectorKeys moved to handlers_class.go

// handleSkillSelectorKeys handles skill selector specific keys (for species selection)
// This is different from handleClassSkillSelectorKeys - this one is for species skill choices
func (m *Model) handleSkillSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.skillSelector.Prev()
	case "down", "j":
		m.skillSelector.Next()
	case "enter":
		selectedSkill := m.skillSelector.GetSelectedSkill()
		if selectedSkill != "" {
			// Apply the skill proficiency and track it as a species skill
			skillNameLower := strings.ToLower(selectedSkill)
			var skillType models.SkillType
			switch skillNameLower {
			case "acrobatics":
				skillType = models.Acrobatics
			case "animal handling":
				skillType = models.AnimalHandling
			case "arcana":
				skillType = models.Arcana
			case "athletics":
				skillType = models.Athletics
			case "deception":
				skillType = models.Deception
			case "history":
				skillType = models.History
			case "insight":
				skillType = models.Insight
			case "intimidation":
				skillType = models.Intimidation
			case "investigation":
				skillType = models.Investigation
			case "medicine":
				skillType = models.Medicine
			case "nature":
				skillType = models.Nature
			case "perception":
				skillType = models.Perception
			case "performance":
				skillType = models.Performance
			case "persuasion":
				skillType = models.Persuasion
			case "religion":
				skillType = models.Religion
			case "sleight of hand":
				skillType = models.SleightOfHand
			case "stealth":
				skillType = models.Stealth
			case "survival":
				skillType = models.Survival
			}
			// Use the helper function to add and track the species skill
			models.AddSpeciesSkillChoice(m.character, skillType)

			// Check if this is Student of War skill selection
			if m.studentOfWarToolSelected {
				m.message = "Battle Master Student of War setup complete!"
				m.studentOfWarToolSelected = false
				m.pendingChanges.Clear()
				m.storage.Save(m.character)
			} else {
			// After skill selection, check if we need spell or feat selection
			species := models.GetSpeciesByName(m.character.Race)
			if species != nil && models.HasSpellChoice(species) {
				// Show wizard cantrip selector for High Elf
				cantrips := models.GetWizardCantrips()
				m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
				m.spellSelector.Show()
				m.message = "Select your wizard cantrip..."
			} else if species != nil && models.HasFeatChoice(species) {
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
			} else {
				m.message = fmt.Sprintf("Skill proficiency gained: %s", selectedSkill)
				// Save when selection is complete (no more selections needed)
				m.storage.Save(m.character)
				}
			}
		}
		m.skillSelector.Hide()
	case "esc":
		m.skillSelector.Hide()
		m.message = "Skill selection cancelled"
	}
	return m, nil
}

// handleSpellSelectorKeys, handleCantripSelectorKeys, handleLeveledSpellSelectorKeys, handleSchoolSpellSelectorKeys,
// handleSpellPrepSelectorKeys, handleSpellbookEditorKeys, handleSlotRestorerKeys moved to handlers_spell.go
// handleDivineOrderSelectorKeys moved to handlers_class.go

// All spell and class handlers moved to handlers_spell.go and handlers_class.go

// handleFeatSelectorKeys handles feat selector specific keys
func (m *Model) handleFeatSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.featSelector.Prev()
	case "down", "j":
		m.featSelector.Next()
	case "pgup", "ctrl+u":
		m.featSelector.PageUp()
	case "pgdown", "ctrl+d":
		m.featSelector.PageDown()
	case "left", "h":
		m.featSelector.PrevCategory()
	case "right", "l":
		m.featSelector.NextCategory()
	case "enter":
		selectedFeat := m.featSelector.GetSelectedFeat()
		if selectedFeat != nil {
			// Check if we're in delete mode
			if m.featSelector.IsDeleteMode() {
				// Remove the feat
				for i, featName := range m.character.Feats {
					if featName == selectedFeat.Name {
						m.character.Feats = append(m.character.Feats[:i], m.character.Feats[i+1:]...)
						break
					}
				}
				// Remove feat benefits (ability increases, HP, speed, etc.)
				models.RemoveFeatBenefits(m.character, *selectedFeat)

				m.message = fmt.Sprintf("Feat removed: %s (benefits reversed)", selectedFeat.Name)
				m.storage.Save(m.character)
				m.featSelector.Hide()
			} else {
				// Check if the feat can be selected (prerequisites met)
				if !m.featSelector.CanSelectCurrentFeat() {
					m.message = fmt.Sprintf("Cannot select %s: Prerequisites not met!", selectedFeat.Name)
					return m, nil
				}

				// Add mode: Check if character already has this feat
			if models.HasFeat(m.character, selectedFeat.Name) && !selectedFeat.Repeatable {
				m.message = fmt.Sprintf("You already have %s and it's not repeatable", selectedFeat.Name)
					m.featSelector.Hide()
			} else {
				// Add feat to character
				err := models.AddFeatToCharacter(m.character, selectedFeat.Name)
				if err != nil {
					m.message = fmt.Sprintf("Error adding feat: %v", err)
						m.featSelector.Hide()
				} else {
						// Check if this feat has ability choices
						if models.HasAbilityChoice(*selectedFeat) {
							// Store the feat and show ability choice selector
							m.pendingFeat = selectedFeat
							m.featSelector.Hide()
							choices := models.GetAbilityChoices(*selectedFeat)
							m.abilityChoiceSelector.Show(selectedFeat.Name, choices, m.character)
							m.message = "Choose which ability to increase"
						} else {
							// Apply feat benefits automatically (no ability choice)
							models.ApplyFeatBenefits(m.character, *selectedFeat, "")
					m.message = fmt.Sprintf("Feat gained: %s!", selectedFeat.Name)
					// Save character after feat selection
					m.storage.Save(m.character)
							m.featSelector.Hide()
				}
			}
		}
			}
		}
	case "esc":
		m.featSelector.Hide()
		m.message = "Feat selection cancelled"
	}
	return m, nil
}

// Detail popup handlers have been moved to handlers package
// See handlers/detail_popup_handler.go

// handleOriginSelectorKeys handles keyboard input for the origin selector
func (m *Model) handleOriginSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.originSelector.Prev()
	case "down", "j":
		m.originSelector.Next()
	case "enter":
		selectedOrigin := m.originSelector.GetSelected()
		if selectedOrigin != nil {
			// Check if origin has ability choice
			if models.HasOriginAbilityChoice(*selectedOrigin) {
				// Store origin temporarily and show ability choice selector
				m.pendingOrigin = selectedOrigin
				m.originSelector.Hide()
				choices := models.GetOriginAbilityChoices(*selectedOrigin)
				m.abilityChoiceSelector.Show(selectedOrigin.Name, choices, m.character)
				m.message = "Choose an ability score to increase..."
			} else {
				// Apply origin directly (no choice needed)
				// Remove old origin first
				if m.character.Origin != "" {
					oldOrigin := models.GetOriginByName(m.character.Origin)
					if oldOrigin != nil {
						models.RemoveOriginBenefits(m.character, *oldOrigin)
					}
				}

				// Apply new origin
				m.character.Origin = selectedOrigin.Name
				models.ApplyOriginBenefits(m.character, *selectedOrigin, "")
				m.storage.Save(m.character)
				m.originSelector.Hide()
				m.message = fmt.Sprintf("Origin changed to: %s", selectedOrigin.Name)
			}
		}
	case "esc":
		m.originSelector.Hide()
		m.message = "Origin selection cancelled"
	}
	return m, nil
}

// handleAbilityChoiceSelectorKeys handles keyboard input for the ability choice selector
func (m *Model) handleAbilityChoiceSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.abilityChoiceSelector.Prev()
	case "down", "j":
		m.abilityChoiceSelector.Next()
	case "enter":
		chosenAbility := m.abilityChoiceSelector.GetSelectedAbility()
		if chosenAbility == "" {
			return m, nil
		}

		// Handle feat ability choice
		if m.pendingFeat != nil {
			// Apply feat benefits with the chosen ability
			models.ApplyFeatBenefits(m.character, *m.pendingFeat, chosenAbility)
			m.message = fmt.Sprintf("Feat gained: %s (+1 %s)!", m.pendingFeat.Name, chosenAbility)
			m.storage.Save(m.character)
			m.pendingFeat = nil
			m.abilityChoiceSelector.Hide()
		}

		// Handle origin ability choice
		if m.pendingOrigin != nil {
			// Remove old origin first
			if m.character.Origin != "" {
				oldOrigin := models.GetOriginByName(m.character.Origin)
				if oldOrigin != nil {
					models.RemoveOriginBenefits(m.character, *oldOrigin)
				}
			}

			// Apply new origin with chosen ability
			m.character.Origin = m.pendingOrigin.Name
			models.ApplyOriginBenefits(m.character, *m.pendingOrigin, chosenAbility)
			m.message = fmt.Sprintf("Origin changed to: %s (+1 %s)!", m.pendingOrigin.Name, chosenAbility)
			m.storage.Save(m.character)
			m.pendingOrigin = nil
			m.abilityChoiceSelector.Hide()
		}
	case "esc":
		// Cancel ability choice
		if m.pendingFeat != nil {
			// Remove the feat from character since we're cancelling
			for i, featName := range m.character.Feats {
				if featName == m.pendingFeat.Name {
					m.character.Feats = append(m.character.Feats[:i], m.character.Feats[i+1:]...)
					break
				}
			}
			m.storage.Save(m.character)
			m.message = "Feat selection cancelled"
			m.pendingFeat = nil
		}

		if m.pendingOrigin != nil {
			m.message = "Origin selection cancelled"
			m.pendingOrigin = nil
		}

		m.abilityChoiceSelector.Hide()
	}
	return m, nil
}

// getContextualHelp returns the panel name and contextual help bindings based on current focus
func (m *Model) getContextualHelp() (string, []components.HelpBinding) {
	switch m.focusArea {
	case FocusMain:
		switch m.currentPanel {
		case StatsPanel:
			return "Stats", components.GetStatsBindings()
		case SkillsPanel:
			return "Skills", components.GetSkillsBindings()
		case InventoryPanel:
			return "Inventory", components.GetInventoryBindings()
		case SpellsPanel:
			return "Spells", components.GetSpellsBindings()
		case FeaturesPanel:
			return "Features", components.GetFeaturesBindings()
		case TraitsPanel:
			// Dynamic bindings for Traits panel based on character state
			bindings := components.GetTraitsBindings()
			// Add maneuver management key if character has maneuvers
			if m.character.IsBattleMaster() && len(m.character.Maneuvers) > 0 {
				bindings = append(bindings, components.HelpBinding{
					Key:  "n",
					Desc: "Manage Battle Master maneuvers",
				})
			}
			return "Traits", bindings
		case OriginPanel:
			return "Origin", components.GetGeneralBindings()
		}
	case FocusCharStats:
		return "Character Info", components.GetCharacterStatsBindings()
	case FocusActions:
		return "Actions", components.GetActionsBindings()
	case FocusDice:
		mode := "idle"
		switch m.dicePanel.GetMode() {
		case panels.DiceModeInput:
			mode = "input"
		case panels.DiceModeHistory:
			mode = "history"
		}
		return "Dice Roller", components.GetDiceBindings(mode)
	}
	return "Stats", components.GetStatsBindings()
}

// buildStatusBar creates the status bar with contextual information
func (m *Model) buildStatusBar() string {
	panelNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("235"))

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Background(lipgloss.Color("235"))

	// Get active panel name and contextual help
	var panelName, contextHelp string

	switch m.focusArea {
	case FocusMain:
		switch m.currentPanel {
		case StatsPanel:
			panelName = "Stats"
			contextHelp = "[r] Roll Stats • [e] Edit Modifiers • [t] Test/Save"
		case SkillsPanel:
			panelName = "Skills"
			contextHelp = "[↑/↓] Navigate • [r] Roll • [e] Toggle Prof"
		case InventoryPanel:
			panelName = "Inventory"
			contextHelp = "[a] Add Item • [e] Equip • [d] Remove 1 • [D] Remove All"
		case SpellsPanel:
			panelName = "Spells"
			contextHelp = "[↑/↓] Navigate • [c] Change Cantrips • [r] Rest"
		case FeaturesPanel:
			panelName = "Features"
			contextHelp = "[↑/↓] Navigate • [u] Use • [+] Restore"
		case TraitsPanel:
			panelName = "Traits"
			contextHelp = "[↑/↓] Navigate • [l] Add Lang • [f] Add Feat • [m] Weapon Mastery"
		case OriginPanel:
			panelName = "Origin"
			contextHelp = "[o] Origin • [Enter] Details • [a] Alignment • [h/w] Height/Weight • [t/i/b/f] Traits • [s] Backstory"
		}
	case FocusCharStats:
		panelName = "Character Info"
		contextHelp = "[n] Name • [h] HP • [r] Short Rest • [R] Long Rest • [+/-] ±1 • [i] Init"
	case FocusActions:
		panelName = "Actions"
		contextHelp = "[↑/↓] Navigate • [Enter] Activate"
	case FocusDice:
		panelName = "Dice Roller"
		switch m.dicePanel.GetMode() {
		case panels.DiceModeIdle:
			contextHelp = "[Enter] Input • [h] History • [r] Reroll"
		case panels.DiceModeInput:
			contextHelp = "Type dice notation • [Enter] Roll • [Esc] Cancel"
		case panels.DiceModeHistory:
			contextHelp = "[↑/↓] Navigate • [Enter] Reroll • [Esc] Back"
		}
	}

	// Build left section: panel + help
	leftSection := panelNameStyle.Render(" "+panelName+" ")

	if contextHelp != "" {
		leftSection += helpStyle.Render(" "+contextHelp+" ")
	}

	// Build right section: global shortcuts
	rightSection := keyStyle.Render("[Tab]") + helpStyle.Render(" Switch tabs • ") +
		keyStyle.Render("[p/P]") + helpStyle.Render(" Focus • ") +
		keyStyle.Render("[s]") + helpStyle.Render(" Save • ") +
		keyStyle.Render("[?]") + helpStyle.Render(" Help • ") +
		keyStyle.Render("[q]") + helpStyle.Render(" Quit ")

	// Calculate padding
	leftWidth := lipgloss.Width(leftSection)
	rightWidth := lipgloss.Width(rightSection)
	padding := m.width - leftWidth - rightWidth
	if padding < 0 {
		padding = 0
	}

	paddingStr := strings.Repeat(" ", padding)

	statusBarStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Width(m.width)

	return statusBarStyle.Render(leftSection + paddingStr + rightSection)
}

// View renders the application
func (m *Model) View() string {
	if !m.ready || m.quitting {
		return ""
	}

	// Show help overlay if visible
	if m.help.Visible {
		panelName, contextBindings := m.getContextualHelp()
		return m.help.ViewWithContext(m.width, m.height, panelName, contextBindings)
	}

	// Calculate layout dimensions
	layoutCalc := view.NewLayoutCalculator()
	layout := layoutCalc.CalculateLayout(m.width, m.height)

	// Tab navigation
	tabBarWidth := layout.MainPanelWidth - 8 // Account for border (2) + horizontal padding (4)
	tabBar := m.tabs.View(tabBarWidth)
	tabHeight := lipgloss.Height(tabBar)

	// Main content height accounts for tabs and spacing
	mainContentHeight := layout.TopRowHeight - tabHeight - 5 // border (2) + padding vertical (2) + spacing line (1)
	if mainContentHeight < 5 {
		mainContentHeight = 5
	}

	// Main panel content (without tabs - tabs will be combined in RenderMainView)
	mainWidth := layout.MainPanelWidth - 8 // Account for border + padding
	var mainPanelView string
	switch m.currentPanel {
	case StatsPanel:
		mainPanelView = m.statsPanel.View(mainWidth, mainContentHeight)
	case SkillsPanel:
		mainPanelView = m.skillsPanel.View(mainWidth, mainContentHeight)
	case InventoryPanel:
		mainPanelView = m.inventoryPanel.View(mainWidth, mainContentHeight)
	case SpellsPanel:
		mainPanelView = m.spellsPanel.View(mainWidth, mainContentHeight)
	case FeaturesPanel:
		mainPanelView = m.featuresPanel.View(mainWidth, mainContentHeight)
	case TraitsPanel:
		mainPanelView = m.traitsPanel.View(mainWidth, mainContentHeight)
	case OriginPanel:
		mainPanelView = m.originPanel.View(mainWidth, mainContentHeight)
	}

	// Character stats view
	charStatsView := m.characterStatsPanel.View(layout.CharStatsInnerWidth, layout.CharStatsInnerHeight)

	// Actions and dice views
	actionsView := m.actionsPanel.View(layout.ActionsWidth, layout.BottomInnerHeight)
	diceView := m.dicePanel.View(layout.DiceWidth, layout.BottomInnerHeight)

	// Status bar
	statusBar := m.buildStatusBar()

	// Render main view using view package (it will combine tabs and content)
	mainView := view.RenderMainView(
		m.width, m.height,
		int(m.focusArea),
		layout,
		tabBar,
		mainPanelView,
		charStatsView,
		actionsView,
		diceView,
		statusBar,
	)

	// Render popups using PopupRenderer
	popupRenderer := view.NewPopupRenderer(m.character, m.width, m.height, layout)
	popupView := popupRenderer.RenderPopups(
		m.statGenerator,
		m.abilityRoller,
		m.attackRoller,
		m.spellSelector,
		m.featSelector,
		m.featDetailPopup,
		m.masteryDetailPopup,
		m.maneuverDetailPopup,
		m.consumableDetailPopup,
		m.featureDetailPopup,
		m.itemDetailPopup,
		m.spellDetailPopup,
		m.originSelector,
		m.alignmentSelector,
		m.traitSelector,
		m.backstoryEditor,
		m.originDetailPopup,
		m.inputPopup,
		m.abilityChoiceSelector,
		m.subtypeSelector,
		m.skillSelector,
		m.languageSelector,
		m.toolSelector,
		m.weaponMasterySelector,
		m.expertiseSelector,
		m.maneuverSelector,
		m.levelUpSelector,
		m.deLevelSelector,
		m.itemSelector,
		m.fightingStyleSelector,
		m.cantripSelector,
		m.leveledSpellSelector,
		m.schoolSpellSelector,
		m.spellbookEditor,
		m.spellPrepSelector,
		m.slotRestorer,
		m.classSkillSelector,
		m.subclassSelector,
		m.classSelector,
		m.speciesSelector,
		m.messagePopup,
		m.restPopup,
		m.attackMenu,
		m.divineOrderSelectorVisible,
		func() string { return m.renderDivineOrderSelector() },
	)

	// If popup is visible, return it; otherwise return main view
	if popupView != "" {
		return popupView
	}

	// Check for HP/XP popups from character stats panel
	hpPopup := m.characterStatsPanel.RenderHPPopup(layout.PopupSmallWidth, layout.PopupSmallHeight)
	if hpPopup != "" {
		return hpPopup
	}

	xpPopup := m.characterStatsPanel.RenderXPPopup(layout.PopupSmallWidth, layout.PopupSmallHeight)
	if xpPopup != "" {
		return xpPopup
	}

	return mainView
}

// getWeaponMasteryCount returns the number of weapons the character can master
func (m *Model) getWeaponMasteryCount() int {
	debug.Log("getWeaponMasteryCount: Checking for Weapon Mastery feature")
	debug.Log("getWeaponMasteryCount: Character class=%s, total features=%d", m.character.Class, len(m.character.Features.Features))

	for i, feature := range m.character.Features.Features {
		debug.Log("getWeaponMasteryCount: Feature[%d]='%s'", i, feature.Name)
		if feature.Name == "Weapon Mastery" {
			debug.Log("getWeaponMasteryCount: Found Weapon Mastery feature!")

			// Read weapons_mastered from feature mechanics
			if feature.Mechanics != nil {
				if weaponsMastered, ok := feature.Mechanics["weapons_mastered"].(float64); ok {
					count := int(weaponsMastered)
					debug.Log("getWeaponMasteryCount: Returning %d from feature mechanics", count)
					return count
				}
			}

			// Fallback: if no mechanics data, return 0
			debug.Log("getWeaponMasteryCount: No mechanics data found, returning 0")
			return 0
		}
	}
	debug.Log("getWeaponMasteryCount: Weapon Mastery feature not found, returning 0")
	return 0
}

// getExpertiseCount returns the number of skills the character can have expertise in
func (m *Model) getExpertiseCount() int {
	debug.Log("getExpertiseCount: Checking for Expertise feature")

	// Check for Expertise feature
	for _, feature := range m.character.Features.Features {
		if feature.Name == "Expertise" && feature.Mechanics != nil {
			if expertiseCount, ok := feature.Mechanics["expertise_count"].(float64); ok {
				debug.Log("getExpertiseCount: Returning %d from feature mechanics", int(expertiseCount))
				return int(expertiseCount)
			}
		}
	}

	// Default: Rogue gets 2 expertise at level 1, 4 at level 6
	if m.character.IsRogue() {
		rogueLevel := m.character.GetRogueLevel()
		if rogueLevel >= 6 {
			debug.Log("getExpertiseCount: Rogue level %d, returning 4", rogueLevel)
			return 4
		} else if rogueLevel >= 1 {
			debug.Log("getExpertiseCount: Rogue level %d, returning 2", rogueLevel)
			return 2
		}
	}

	debug.Log("getExpertiseCount: No expertise feature found, returning 0")
	return 0
}

// handleAlignmentSelectorKeys handles alignment selector input
func (m *Model) handleAlignmentSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.alignmentSelector.Hide()
		m.message = ""
	case "up", "k":
		m.alignmentSelector.Prev()
	case "down", "j":
		m.alignmentSelector.Next()
	case "enter":
		selected := m.alignmentSelector.GetSelectedAlignment()
		m.character.Alignment = selected
		m.alignmentSelector.Hide()
		m.message = fmt.Sprintf("Alignment set to %s", selected)
		m.storage.Save(m.character)
	}
	return m, nil
}

// handleTraitSelectorKeys handles trait selector input
func (m *Model) handleTraitSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If in custom mode, handle it separately to allow all keys including 'p' and space
	if m.traitSelector.IsCustomMode() {
		switch msg.String() {
		case "esc":
			m.traitSelector.Hide()
			m.message = ""
		case "enter":
			// Save custom value and add to the list
			customValue := m.traitSelector.GetSelectedTrait()
			if customValue != "" {
				// Add to appropriate array
				switch m.traitSelector.TraitType {
				case models.TraitPersonality:
					m.character.Personality = append(m.character.Personality, customValue)
					m.message = "Custom personality trait added!"
				case models.TraitIdeal:
					m.character.Ideal = append(m.character.Ideal, customValue)
					m.message = "Custom ideal added!"
				case models.TraitBond:
					m.character.Bond = append(m.character.Bond, customValue)
					m.message = "Custom bond added!"
				case models.TraitFlaw:
					m.character.Flaw = append(m.character.Flaw, customValue)
					m.message = "Custom flaw added!"
				}
				m.traitSelector.Hide()
				m.storage.Save(m.character)
			}
		default:
			// Pass all other keys to text input (including 'p', space, etc.)
			m.traitSelector.Update(msg)
		}
		return m, nil
	}

	// Normal multi-select mode
	switch msg.String() {
	case "esc":
		m.traitSelector.Hide()
		m.message = ""
	case "up", "k":
		m.traitSelector.Prev()
	case "down", "j":
		m.traitSelector.Next()
	case " ":
		// Space key toggles the current item
		m.traitSelector.ToggleItem()
	case "enter":
		// Check if [Custom] is selected to enter custom mode
		m.traitSelector.ToggleCustomMode()

		// If not entering custom mode, save all selected items
		if !m.traitSelector.IsCustomMode() {
			selectedItems := m.traitSelector.GetSelectedItems()

			// Save the trait based on current type
			switch m.traitSelector.TraitType {
			case models.TraitPersonality:
				m.character.Personality = selectedItems
				m.message = "Personality traits saved!"
			case models.TraitIdeal:
				m.character.Ideal = selectedItems
				m.message = "Ideals saved!"
			case models.TraitBond:
				m.character.Bond = selectedItems
				m.message = "Bonds saved!"
			case models.TraitFlaw:
				m.character.Flaw = selectedItems
				m.message = "Flaws saved!"
			}
			m.traitSelector.Hide()
			m.storage.Save(m.character)
		}
	}
	return m, nil
}

// handleBackstoryEditorKeys handles backstory editor input
func (m *Model) handleBackstoryEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+enter":
		// Save and close
		m.character.Backstory = m.backstoryEditor.GetValue()
		m.backstoryEditor.Hide()
		m.message = "Backstory saved!"
		m.storage.Save(m.character)
	default:
		// Update textarea
		m.backstoryEditor.Update(msg)
	}
	return m, nil
}

// handleOriginDetailPopupKeys handles origin detail popup input
// handleOriginDetailPopupKeys moved to handlers package

// handleInputPopupKeys handles input popup input
func (m *Model) handleInputPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.inputPopup.Hide()
		m.inputPopupContext = ""
		m.message = "Cancelled"
	case "enter":
		// Save based on context
		value := m.inputPopup.GetValue()
		switch m.inputPopupContext {
		case "height":
			m.character.Height = value
			m.message = "Height updated!"
		case "weight":
			m.character.Weight = value
			m.message = "Weight updated!"
		}
		m.inputPopup.Hide()
		m.inputPopupContext = ""
		m.storage.Save(m.character)
	default:
		// Update text input
		m.inputPopup.Update(msg)
	}
	return m, nil
}

// renderDivineOrderSelector renders the Divine Order selection popup
func (m *Model) renderDivineOrderSelector() string {
	popupMediumWidth := int(float64(m.width) * 0.60)
	popupMediumHeight := int(float64(m.height) * 0.50)
	if popupMediumWidth < 70 {
		popupMediumWidth = 70
	}
	if popupMediumHeight < 20 {
		popupMediumHeight = 20
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		MarginBottom(1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center).
		MarginTop(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(2, 4).
		Width(popupMediumWidth - 8).
		Align(lipgloss.Center)

	var optionsText string
	if m.pendingDivineOrder == "Thaumaturgic" {
		// Show skill selection for Thaumaturgic
		optionsText = optionStyle.Render("Choose skill for Thaumaturgic:") + "\n\n" +
			highlightStyle.Render("[a]") + " Arcana\n" +
			highlightStyle.Render("[r]") + " Religion\n\n" +
			hintStyle.Render("Press 'a' for Arcana or 'r' for Religion")
	} else {
		// Show Divine Order selection
		optionsText = titleStyle.Render("Select Divine Order") + "\n\n" +
			optionStyle.Render("Choose how you channel your divine faith:") + "\n\n" +
			highlightStyle.Render("[1]") + " Protector\n" +
			"   Proficiency with martial weapons and heavy armor\n" +
			"   3 cantrips\n\n" +
			highlightStyle.Render("[2]") + " Thaumaturgic\n" +
			"   Extra cantrip (+1, for 4 total)\n" +
			"   Expertise in Arcana or Religion (your choice)\n\n" +
			hintStyle.Render("Press 1 for Protector or 2 for Thaumaturgic • ESC to cancel")
	}

	box := boxStyle.Render(optionsText)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

// Run runs the application
func Run(char *models.Character, store *storage.Storage) error {
	p := tea.NewProgram(
		NewModel(char, store),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
