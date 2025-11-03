// internal/ui/app.go
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
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

// Panel handlers moved to handlers_panel.go

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
// handleFeatSelectorKeys and handleAbilityChoiceSelectorKeys moved to handlers_feat.go

// Detail popup handlers have been moved to handlers package
// See handlers/detail_popup_handler.go

// handleOriginSelectorKeys, handleAlignmentSelectorKeys, handleTraitSelectorKeys,
// handleBackstoryEditorKeys, and handleInputPopupKeys moved to handlers_origin.go

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

// Origin panel handlers moved to handlers_origin.go
// renderDivineOrderSelector moved to handlers_class.go (Cleric-specific)

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
