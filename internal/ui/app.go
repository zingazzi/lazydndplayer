// internal/ui/app.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
	"github.com/marcozingoni/lazydndplayer/internal/ui/services"
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
	CompanionPanel
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
	storage      StorageInterface // Use interface instead of concrete type

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
	beastSelector         *components.BeastSelector

	// Main Panels (switchable)
	statsPanel     *panels.StatsPanel
	skillsPanel    *panels.SkillsPanel
	inventoryPanel *panels.InventoryPanel
	spellsPanel    *panels.SpellsPanel
	featuresPanel  *panels.FeaturesPanel
	traitsPanel    *panels.TraitsPanel
	originPanel    *panels.OriginPanel
	companionPanel *panels.CompanionPanel

	// Fixed Panels (always visible)
	dicePanel           *panels.DicePanel
	characterStatsPanel *panels.CharacterStatsPanel
	actionsPanel        *panels.ActionsPanel // Bottom panel for quick actions

	// Component Management
	componentManager *ComponentManager
	stateMachine     *state.StateMachine

	// Services
	rollService  *services.RollService
	featService  *services.FeatService
	originService *services.OriginService
	classService  *services.ClassService

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

// NewModel creates a new application model with dependency injection
// If factory is nil, it uses the default ComponentFactory
func NewModel(char *models.Character, store StorageInterface, factory ComponentFactoryInterface) *Model {
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

	// Use default factory if none provided
	if factory == nil {
		factory = NewComponentFactory()
	}

	// Create component manager and state machine
	componentManager := NewComponentManager()
	stateMachine := state.NewStateMachine()

	// Create services
	rollService := services.NewRollService()
	featService := services.NewFeatService()
	originService := services.NewOriginService()
	classService := services.NewClassService()

	// Create components using factory
	tabs := factory.CreateTabs()
	help := factory.CreateHelp()
	speciesSelector := factory.CreateSpeciesSelector()
	subtypeSelector := factory.CreateSubtypeSelector()
	languageSelector := factory.CreateLanguageSelector()
	skillSelector := factory.CreateSkillSelector()
	spellSelector := factory.CreateSpellSelector()
	featSelector := factory.CreateFeatSelector()
	featDetailPopup := factory.CreateFeatDetailPopup()
	featureDetailPopup := factory.CreateFeatureDetailPopup()
	itemDetailPopup := factory.CreateItemDetailPopup()
	masteryDetailPopup := factory.CreateMasteryDetailPopup()
	maneuverDetailPopup := factory.CreateManeuverDetailPopup()
	consumableDetailPopup := factory.CreateConsumableDetailPopup()
	spellDetailPopup := factory.CreateSpellDetailPopup()
	originSelector := factory.CreateOriginSelector()
	alignmentSelector := factory.CreateAlignmentSelector()
	traitSelector := factory.CreateTraitSelector()
	backstoryEditor := factory.CreateBackstoryEditor()
	originDetailPopup := factory.CreateOriginDetailPopup()
	inputPopup := factory.CreateInputPopup()
	toolSelector := factory.CreateToolSelector()
	itemSelector := factory.CreateItemSelector()
	classSelector := factory.CreateClassSelector(char)
	classSkillSelector := factory.CreateClassSkillSelector()
	subclassSelector := factory.CreateSubclassSelector(char)
	fightingStyleSelector := factory.CreateFightingStyleSelector()
	cantripSelector := factory.CreateCantripSelector(char)
	leveledSpellSelector := factory.CreateLeveledSpellSelector(char)
	schoolSpellSelector := factory.CreateSchoolSpellSelector(char)
	spellbookEditor := factory.CreateSpellbookEditor(char)
	spellPrepSelector := factory.CreateSpellPrepSelector(char)
	slotRestorer := factory.CreateSlotRestorer(char)
	statGenerator := factory.CreateStatGenerator()
	abilityRoller := factory.CreateAbilityRoller()
	abilityChoiceSelector := factory.CreateAbilityChoiceSelector()
	attackRoller := factory.CreateAttackRoller()
	attackMenu := factory.CreateAttackMenu()
	weaponMasterySelector := factory.CreateWeaponMasterySelector(char)
	expertiseSelector := factory.CreateExpertiseSelector(char)
	maneuverSelector := factory.CreateManeuverSelector()
	levelUpSelector := factory.CreateLevelUpSelector(char)
	deLevelSelector := factory.CreateDeLevelSelector(char)
	restPopup := factory.CreateRestPopup(char, models.NewStandardDiceRoller())
	messagePopup := factory.CreateMessagePopup()
	beastSelector := components.NewBeastSelector(char)

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
	componentManager.Register(beastSelector, 61, "beastSelector")
	componentManager.Register(classSelector, 60, "classSelector")
	componentManager.Register(speciesSelector, 59, "speciesSelector")
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
		statsPanel:            factory.CreateStatsPanel(char),
		skillsPanel:           factory.CreateSkillsPanel(char),
		inventoryPanel:        factory.CreateInventoryPanel(char),
		spellsPanel:           factory.CreateSpellsPanel(char),
		featuresPanel:         factory.CreateFeaturesPanel(char),
		traitsPanel:           factory.CreateTraitsPanel(char),
		originPanel:           factory.CreateOriginPanel(char),
		companionPanel:        factory.CreateCompanionPanel(char),
		beastSelector:         components.NewBeastSelector(char),
		dicePanel:           factory.CreateDicePanel(char),
		characterStatsPanel: factory.CreateCharacterStatsPanel(char),
		actionsPanel:        factory.CreateActionsPanel(char),
		componentManager:    componentManager,
		stateMachine:        stateMachine,
		rollService:         rollService,
		featService:         featService,
		originService:       originService,
		classService:        classService,
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

		// Character Creation Wizard - Shift+W
		case "W":
			// Only allow wizard if character has no class
			if len(m.character.Classes) == 0 && m.character.Class == "" {
				m.startWizard()
			} else {
				m.message = "Wizard only works for new characters (no class selected)"
			}
			return m, nil
		}

		// Handle special components that need processing before normal routing
		if model, cmd, handled := m.handleSpecialComponents(msg); handled {
			return model, cmd
		}

		// Check if levelUpSelector is visible and in class selection mode - it should get priority
		// when selecting a class to level up
		if m.levelUpSelector.IsVisible() && m.levelUpSelector.GetState() == components.LevelUpSelectClass {
			debug.Log("Update: levelUpSelector is visible in class selection mode, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.levelUpSelector, msg); handled {
				debug.Log("Update: levelUpSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Check if classSkillSelector is visible - it should take priority over levelUpSelector
		// when Deft Explorer/Scholar selection is needed
		if m.classSkillSelector.IsVisible() {
			debug.Log("Update: classSkillSelector is visible, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.classSkillSelector, msg); handled {
				debug.Log("Update: classSkillSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Check if languageSelector is visible in multi-select mode - it should take priority over levelUpSelector
		// when Deft Explorer language selection is needed
		if m.languageSelector.IsVisible() && m.languageSelector.IsMultiSelectMode() {
			debug.Log("Update: languageSelector (multi-select) is visible, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.languageSelector, msg); handled {
				debug.Log("Update: languageSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Check if fightingStyleSelector is visible - it should take priority over levelUpSelector
		// when fighting style selection is needed (e.g., Ranger level 2)
		if m.fightingStyleSelector.IsVisible() {
			debug.Log("Update: fightingStyleSelector is visible, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.fightingStyleSelector, msg); handled {
				debug.Log("Update: fightingStyleSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Check if cantripSelector is visible - it should take priority over levelUpSelector
		// when cantrip selection is needed (e.g., Druidic Warrior, Blessed Warrior)
		if m.cantripSelector.IsVisible() {
			debug.Log("Update: cantripSelector is visible, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.cantripSelector, msg); handled {
				debug.Log("Update: cantripSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Check if beastSelector is visible - it should take priority over levelUpSelector
		// when beast selection is needed (e.g., Beast Master level 3)
		if m.beastSelector.IsVisible() {
			debug.Log("Update: beastSelector is visible, routing key=%s directly to it", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(m.beastSelector, msg); handled {
				debug.Log("Update: beastSelector handler returned handled=true for key=%s", msg.String())
				return model, cmd
			}
		}

		// Panel navigation (BEFORE component routing, but only when focused on main)
		// This ensures tab navigation works even if components are visible
		if m.focusArea == FocusMain {
			switch msg.String() {
			case "tab":
				// Check if any component should handle tab for its own navigation
				// Components that use tab internally: SpellPrepSelector, SpellbookEditor
				shouldBlockTab := false
				if visibleComponent := m.componentManager.GetVisibleComponent(); visibleComponent != nil {
					switch visibleComponent.(type) {
					case *components.SpellPrepSelector, *components.SpellbookEditor:
						// These components use tab for internal navigation - let them handle it
						shouldBlockTab = true
					case *components.LevelUpSelector, *components.DeLevelSelector:
						// Block tab navigation during level up/de-level (only if actually visible)
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					case *components.SpeciesSelector, *components.SubtypeSelector:
						// Block tab navigation during character creation (only if actually visible)
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					case *components.ClassSelector:
						// Block tab navigation during class selection (only if actually visible)
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					// All other components should allow tab to pass through for panel navigation
					}
				}

				if !shouldBlockTab {
					m.tabs.Next()
					m.currentPanel = PanelType(m.tabs.SelectedIndex)
					debug.Log("Tab navigation: moved to panel %d", m.currentPanel)
					return m, nil
				}

			case "shift+tab":
				// Same logic for shift+tab
				shouldBlockTab := false
				if visibleComponent := m.componentManager.GetVisibleComponent(); visibleComponent != nil {
					switch visibleComponent.(type) {
					case *components.SpellPrepSelector, *components.SpellbookEditor:
						// These components use shift+tab for internal navigation
						shouldBlockTab = true
					case *components.LevelUpSelector, *components.DeLevelSelector:
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					case *components.SpeciesSelector, *components.SubtypeSelector:
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					case *components.ClassSelector:
						if visibleComponent.(interface{ IsVisible() bool }).IsVisible() {
							shouldBlockTab = true
						}
					}
				}

				if !shouldBlockTab {
					m.tabs.Prev()
					m.currentPanel = PanelType(m.tabs.SelectedIndex)
					debug.Log("Shift+Tab navigation: moved to panel %d", m.currentPanel)
					return m, nil
				}
			}
		}

		// Use ComponentManager to route to the highest priority visible component
		// This replaces 30+ individual if statements with a single priority-based routing
		// If a component doesn't handle the key (returns handled=false), it falls through to panel handlers
		if visibleComponent := m.componentManager.GetVisibleComponent(); visibleComponent != nil {
			debug.Log("Update: Found visible component, routing key=%s", msg.String())
			if model, cmd, handled := m.routeComponentToHandler(visibleComponent, msg); handled {
				debug.Log("Update: Component handler returned handled=true for key=%s", msg.String())
				return model, cmd
			} else {
				debug.Log("Update: Component handler returned handled=false for key=%s, falling through to panel handlers", msg.String())
				// Key not handled by component - continue to panel handlers below
			}
		} else {
			debug.Log("Update: No visible component found for key=%s", msg.String())
		}

		// Handle input based on current focus
		debug.Log("Update: Handling key in focusArea=%d (0=Main,1=CharStats,2=Actions,3=Dice)", m.focusArea)
		switch m.focusArea {
		case FocusMain:
			return m.handleMainPanelKeys(msg)
		case FocusCharStats:
			// When in CharStats panel, tab should switch focus back to main and navigate tabs
			if msg.String() == "tab" {
				m.focusArea = FocusMain
				m.tabs.Next()
				m.currentPanel = PanelType(m.tabs.SelectedIndex)
				debug.Log("Tab navigation: switched from CharStats to Main, moved to panel %d", m.currentPanel)
				return m, nil
			} else if msg.String() == "shift+tab" {
				m.focusArea = FocusMain
				m.tabs.Prev()
				m.currentPanel = PanelType(m.tabs.SelectedIndex)
				debug.Log("Shift+Tab navigation: switched from CharStats to Main, moved to panel %d", m.currentPanel)
				return m, nil
			}
			return m.handleCharStatsPanelKeys(msg)
		case FocusActions:
			return m.handleActionsPanelKeys(msg)
		case FocusDice:
			return m.handleDicePanelKeys(msg)
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

// handleSkillSelectorKeys moved to handlers_selectors.go
// handleSpellSelectorKeys, handleCantripSelectorKeys, handleLeveledSpellSelectorKeys, handleSchoolSpellSelectorKeys,
// handleSpellPrepSelectorKeys, handleSpellbookEditorKeys, handleSlotRestorerKeys moved to handlers_spell.go
// handleDivineOrderSelectorKeys moved to handlers_class.go

// All spell and class handlers moved to handlers_spell.go and handlers_class.go
// handleFeatSelectorKeys and handleAbilityChoiceSelectorKeys moved to handlers_feat.go

// Detail popup handlers have been moved to handlers package
// See handlers/detail_popup_handler.go

// handleOriginSelectorKeys, handleAlignmentSelectorKeys, handleTraitSelectorKeys,
// handleBackstoryEditorKeys, and handleInputPopupKeys moved to handlers_origin.go

// Status bar and contextual help methods moved to view/status_bar.go

// View renders the application
func (m *Model) View() string {
	if !m.ready || m.quitting {
		return ""
	}

	// Show help overlay if visible
	if m.help.Visible {
		panelName, contextBindings := view.GetContextualHelp(
			int(m.focusArea),
			view.PanelType(m.currentPanel),
			m.dicePanel.GetMode(),
			m.character.IsBattleMaster(),
			len(m.character.Maneuvers) > 0,
		)
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
	// Use PanelRenderer instead of switch statement
	panelRenderer := view.RegisterPanelsFromModel(
		m.statsPanel,
		m.skillsPanel,
		m.inventoryPanel,
		m.spellsPanel,
		m.featuresPanel,
		m.traitsPanel,
		m.originPanel,
		m.companionPanel,
	)
	mainPanelView := panelRenderer.RenderPanel(view.PanelType(m.currentPanel), mainWidth, mainContentHeight)

	// Character stats view
	charStatsView := m.characterStatsPanel.View(layout.CharStatsInnerWidth, layout.CharStatsInnerHeight)

	// Actions and dice views
	actionsView := m.actionsPanel.View(layout.ActionsWidth, layout.BottomInnerHeight)
	diceView := m.dicePanel.View(layout.DiceWidth, layout.BottomInnerHeight)

	// Status bar
	statusBar := view.BuildStatusBar(view.StatusBarContext{
		FocusArea:    int(m.focusArea),
		CurrentPanel: view.PanelType(m.currentPanel),
		DiceMode:     m.dicePanel.GetMode(),
		Width:        m.width,
	})

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
		m.beastSelector,
		m.classSelector,
		m.speciesSelector,
		m.messagePopup,
		m.restPopup,
		m.attackMenu,
		m.IsDivineOrderSelectorVisible(),
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

// Utility methods moved to utils.go
// getWeaponMasteryCount and getExpertiseCount are now package-level functions

// Origin panel handlers moved to handlers_origin.go
// renderDivineOrderSelector moved to handlers_class.go (Cleric-specific)

// Run runs the application
func Run(char *models.Character, store *storage.Storage) error {
	// Convert concrete storage to interface (storage.Storage implements StorageInterface)
	var storageInterface StorageInterface = store
	p := tea.NewProgram(
		NewModel(char, storageInterface, nil), // nil = use default factory
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
