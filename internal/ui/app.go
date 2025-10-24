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
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
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
	itemDetailPopup       *components.ItemDetailPopup
	spellDetailPopup      *components.SpellDetailPopup
	originSelector        *components.OriginSelector
	toolSelector          *components.ToolSelector
	itemSelector           *components.ItemSelector
	classSelector          *components.ClassSelector
	classSkillSelector     *components.ClassSkillSelector
	subclassSelector       *components.SubclassSelector
	fightingStyleSelector  *components.FightingStyleSelector
	cantripSelector        *components.CantripSelector
	spellPrepSelector      *components.SpellPrepSelector
	slotRestorer           *components.SlotRestorer
	statGenerator          *components.StatGenerator
	abilityRoller         *components.AbilityRoller
	abilityChoiceSelector *components.AbilityChoiceSelector
	attackRoller          *components.AttackRoller
	attackMenu            *components.AttackMenu
	weaponMasterySelector *components.WeaponMasterySelector
	levelUpSelector       *components.LevelUpSelector
	deLevelSelector       *components.DeLevelSelector

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

	// State
	currentPanel       PanelType
	focusArea          FocusArea
	width              int
	height             int
	ready              bool
	message            string
	quitting           bool
	pendingFeat        *models.Feat   // Temporarily store feat while choosing ability
	pendingOrigin      *models.Origin // Temporarily store origin while choosing ability
	pendingChanges     *models.PendingChanges // Transaction system for rollback support
}

// NewModel creates a new application model
func NewModel(char *models.Character, store *storage.Storage) *Model {
	return &Model{
		character:           char,
		storage:             store,
		tabs:                components.NewTabs(),
		help:                components.NewHelp(),
		speciesSelector:     components.NewSpeciesSelector(),
		subtypeSelector:     components.NewSubtypeSelector(),
		languageSelector:    components.NewLanguageSelector(),
		skillSelector:       components.NewSkillSelector(),
		spellSelector:       components.NewSpellSelector(),
		featSelector:          components.NewFeatSelector(),
		featDetailPopup:       components.NewFeatDetailPopup(),
		itemDetailPopup:       components.NewItemDetailPopup(),
		spellDetailPopup:      components.NewSpellDetailPopup(),
		originSelector:        components.NewOriginSelector(),
		toolSelector:          components.NewToolSelector(),
		itemSelector:           components.NewItemSelector(),
		classSelector:          components.NewClassSelector(char),
		classSkillSelector:     components.NewClassSkillSelector(),
		subclassSelector:       components.NewSubclassSelector(char),
		fightingStyleSelector:  components.NewFightingStyleSelector(),
		cantripSelector:        components.NewCantripSelector(char),
		spellPrepSelector:      components.NewSpellPrepSelector(char),
		slotRestorer:           components.NewSlotRestorer(char),
		statGenerator:          components.NewStatGenerator(),
		abilityRoller:         components.NewAbilityRoller(),
		abilityChoiceSelector: components.NewAbilityChoiceSelector(),
		attackRoller:          components.NewAttackRoller(),
		attackMenu:            components.NewAttackMenu(),
		weaponMasterySelector: components.NewWeaponMasterySelector(char),
		levelUpSelector:       components.NewLevelUpSelector(char),
		deLevelSelector:       components.NewDeLevelSelector(char),
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
		currentPanel:        StatsPanel,
		focusArea:           FocusMain,
		pendingChanges:      models.NewPendingChanges(),
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.ClearScreen
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

		// Panel navigation (only when focused on main and no popups)
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
			return m.handleFeatDetailPopupKeys(msg)
		}

		// Check if item detail popup is active
		if m.itemDetailPopup.IsVisible() {
			return m.handleItemDetailPopupKeys(msg)
		}

		// Check if spell detail popup is active
		if m.spellDetailPopup.IsVisible() {
			return m.handleSpellDetailPopupKeys(msg)
		}

		// Check if origin selector is active
		if m.originSelector.IsVisible() {
			return m.handleOriginSelectorKeys(msg)
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

		// Check if fighting style selector is active (highest priority in class flow)
		if m.fightingStyleSelector.IsVisible() {
			return m.handleFightingStyleSelectorKeys(msg)
		}

		// Check if cantrip selector is active
		if m.cantripSelector.IsVisible() {
			return m.handleCantripSelectorKeys(msg)
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
	case "v":
		// Open spell preparation selector
		if m.character.SpellBook.IsPreparedCaster {
			m.spellPrepSelector.Show()
			m.message = "Select spells to prepare..."
		} else {
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
	case "enter":
		// View spell details
		spell := m.spellsPanel.GetSelectedSpell()
		if spell != nil {
			m.spellDetailPopup.Show(*spell)
			m.message = "Viewing spell details..."
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
	case "u":
		// Use feature (decrement uses and apply effect)
		debug.Log("Features panel: 'u' pressed")
		debug.Log("Total features: %d", len(m.character.Features.Features))

		feature := m.featuresPanel.GetSelectedFeature()
		debug.Log("Selected feature: %v", feature != nil)

		if feature != nil {
			debug.Log("Feature name: %s, CurrentUses: %d, MaxUses: %d", feature.Name, feature.CurrentUses, feature.MaxUses)

			// Check if feature has uses remaining
			if feature.CurrentUses > 0 {
				debug.Log("Using feature: %s", feature.Name)

				// Apply feature effect based on name
				switch feature.Name {
				case "Second Wind":
					debug.Log("Applying Second Wind effect")
					// Roll 1d10 + level
					result, err := dice.Roll("1d10", dice.Normal)
					if err == nil {
						healing := result.Total + m.character.Level
						oldHP := m.character.CurrentHP
						m.character.CurrentHP += healing
						if m.character.CurrentHP > m.character.MaxHP {
							m.character.CurrentHP = m.character.MaxHP
						}
						debug.Log("Second Wind: Healed from %d to %d HP (rolled %d + level %d)", oldHP, m.character.CurrentHP, result.Total, m.character.Level)
						m.message = fmt.Sprintf("Second Wind: Healed %d HP (1d10[%d] + %d level)", healing, result.Total, m.character.Level)
					} else {
						debug.Log("Error rolling dice for Second Wind: %v", err)
						m.message = fmt.Sprintf("Second Wind: Healed %d HP", m.character.Level)
						m.character.CurrentHP += m.character.Level
					}
				default:
					debug.Log("Using generic feature: %s", feature.Name)
					m.message = fmt.Sprintf("%s used", feature.Name)
				}

				// Decrement uses
		m.featuresPanel.UseFeature()
				debug.Log("Feature uses decremented. New uses: %d", feature.CurrentUses)
		m.storage.Save(m.character)
			} else {
				debug.Log("Feature %s has no uses remaining", feature.Name)
				m.message = fmt.Sprintf("%s has no uses remaining", feature.Name)
			}
		} else {
			debug.Log("No feature selected")
			m.message = "No feature selected"
		}
	case "+", "=":
		// Restore one use
		m.featuresPanel.RestoreFeature()
		m.message = "Feature restored"
		m.storage.Save(m.character)
	case "d", "delete":
		// Delete feature
		m.featuresPanel.RemoveFeature()
		m.message = "Feature removed"
		m.storage.Save(m.character)
	case "a":
		// Add feature (simplified - in real app would show a form)
		m.message = "Add feature (not yet implemented)"
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
	case "t":
		// Add tool proficiency
		m.toolSelector.SetExcludeTools(m.character.ToolProficiencies)
		m.toolSelector.Show()
		m.message = "Select tool proficiency to add..."
	case "T":
		// Remove tool proficiency
		m.toolSelector.ShowForDeletion(m.character.ToolProficiencies)
		m.message = "Select tool proficiency to remove..."
	}
	return m, nil
}

// handleTraitsPanel handles traits panel specific keys
func (m *Model) handleTraitsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleTraitsPanel: key=%s", msg.String())
	switch msg.String() {
	case "up", "k":
		// Scroll viewport up freely
		m.traitsPanel.ScrollUp()
	case "down", "j":
		// Scroll viewport down freely
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
	case "+", "=":
		// Add Focus Point (Monk only)
		if m.character.IsMonk() {
			monk := m.character.GetMonkMechanics()
			monk.RestoreFocusPoints(1)
			current, max := monk.GetFocusPoints()
			m.message = fmt.Sprintf("Focus Point restored. Current: %d/%d", current, max)
		} else {
			m.message = "Only Monks can use Focus Points"
		}
		return m, nil
	case "-", "_":
		// Spend Focus Point (Monk only)
		if m.character.IsMonk() {
			monk := m.character.GetMonkMechanics()
			if monk.SpendFocusPoint(1) {
				current, max := monk.GetFocusPoints()
				m.message = fmt.Sprintf("Focus Point spent. Current: %d/%d", current, max)
			} else {
				m.message = "Not enough Focus Points"
			}
		} else {
			m.message = "Only Monks can use Focus Points"
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
		m.speciesSelector.Show()
		m.message = "Select a species..."
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

// handleAbilityRollerKeys handles ability roller specific keys
func (m *Model) handleAbilityRollerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.abilityRoller.Prev()
	case "down", "j":
		m.abilityRoller.Next()
	case "tab":
		m.abilityRoller.SwitchFocus()
		m.message = "Switched focus"
	case "space":
		m.abilityRoller.ToggleType()
	case "enter":
		// Roll the dice!
		expr := m.abilityRoller.GetRollExpression(m.character)
		description := m.abilityRoller.GetRollDescription(m.character)
		m.dicePanel.Roll(expr)
		m.message = fmt.Sprintf("%s: %s", description, m.dicePanel.LastMessage)
		m.abilityRoller.Hide()
	case "esc":
		m.abilityRoller.Hide()
		m.message = "Roll cancelled"
	}
	return m, nil
}

// handleSpeciesSelectorKeys handles species selector specific keys
func (m *Model) handleSpeciesSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.speciesSelector.Prev()
	case "down", "j":
		m.speciesSelector.Next()
	case "enter":
		selectedSpecies := m.speciesSelector.GetSelectedSpecies()
		if selectedSpecies != nil {
			// Check if species has subtypes
			if selectedSpecies.HasSubtypes && len(selectedSpecies.Subtypes) > 0 {
				// Show subtype selector
				subtypes := make([]components.SpeciesSubtype, len(selectedSpecies.Subtypes))
				for i, st := range selectedSpecies.Subtypes {
					subtypes[i] = components.SpeciesSubtype{
						Name:        st.Name,
						Description: st.Description,
						Modifier:    st.Modifier,
					}
					// For Dragonborn, show damage type
					if st.DamageType != "" {
						subtypes[i].Modifier = fmt.Sprintf("%s damage", st.DamageType)
					}
				}
				m.subtypeSelector.Show(selectedSpecies.Name, subtypes)
				m.message = fmt.Sprintf("Select %s subtype...", selectedSpecies.Name)
				m.speciesSelector.Hide()
				return m, nil
			}

			// No subtypes, apply species directly
			oldSpecies := m.character.Race
			models.ApplySpeciesToCharacter(m.character, selectedSpecies.Name)

			// Check if we need to select additional languages
			needsLanguageSelection := false
			for _, lang := range m.character.Languages {
				if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
					needsLanguageSelection = true
					break
				}
			}

			// Check for various selections needed
			needsSkillSelection := models.HasSkillChoice(selectedSpecies)
			needsSpellSelection := models.HasSpellChoice(selectedSpecies)
			needsFeatSelection := models.HasFeatChoice(selectedSpecies)

			if needsLanguageSelection {
				// Filter out languages the character already knows
				m.languageSelector.SetExcludeLanguages(m.character.Languages)
				m.languageSelector.Show()
				m.message = "Select your additional language..."
			} else if needsSkillSelection {
				m.skillSelector.Show()
				m.message = "Select your skill proficiency..."
			} else if needsSpellSelection {
				// Show wizard cantrip selector
				cantrips := models.GetWizardCantrips()
				m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
				m.spellSelector.Show()
				m.message = "Select your wizard cantrip..."
			} else if needsFeatSelection {
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
			} else {
				m.message = fmt.Sprintf("Species changed from %s to %s. Speed updated to %d ft.", oldSpecies, selectedSpecies.Name, m.character.Speed)
				// Save character when species change is complete (no additional selections)
				m.storage.Save(m.character)
			}
		}
		m.speciesSelector.Hide()
	case "esc":
		m.speciesSelector.Hide()
		m.message = "Species selection cancelled"
	}
	return m, nil
}

// handleSubtypeSelectorKeys handles subtype selector specific keys
func (m *Model) handleSubtypeSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.subtypeSelector.Prev()
	case "down", "j":
		m.subtypeSelector.Next()
	case "enter":
		selectedSubtype := m.subtypeSelector.GetSelectedSubtype()
		if selectedSubtype != nil {
			// ALWAYS use the species name from the selector (just selected)
			// Don't use m.character.Race as it might still have the old species
			speciesName := m.subtypeSelector.SpeciesName

			// Apply species with the selected subtype
			models.ApplySpeciesWithSubtype(m.character, speciesName, selectedSubtype.Name)

			// Get species info for additional checks
			species := models.GetSpeciesByName(speciesName)

			// Check if we need to select additional languages
			needsLanguageSelection := false
			for _, lang := range m.character.Languages {
				if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
					needsLanguageSelection = true
					break
				}
			}

			// Check for various selections needed
			needsSkillSelection := species != nil && models.HasSkillChoice(species)
			needsSpellSelection := species != nil && models.HasSpellChoice(species)
			needsFeatSelection := species != nil && models.HasFeatChoice(species)

			if needsLanguageSelection {
				// Filter out languages the character already knows
				m.languageSelector.SetExcludeLanguages(m.character.Languages)
				m.languageSelector.Show()
				m.message = "Select your additional language..."
			} else if needsSkillSelection {
				m.skillSelector.Show()
				m.message = "Select your skill proficiency..."
			} else if needsSpellSelection {
				// Show wizard cantrip selector (High Elf)
				cantrips := models.GetWizardCantrips()
				m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
				m.spellSelector.Show()
				m.message = "Select your wizard cantrip..."
			} else if needsFeatSelection {
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
			} else {
				m.message = fmt.Sprintf("%s (%s) selected! Speed: %d ft, Darkvision: %d ft",
					speciesName, selectedSubtype.Name, m.character.Speed, m.character.Darkvision)
				// Save character when species change is complete (no additional selections)
				m.storage.Save(m.character)
			}
		}
		m.subtypeSelector.Hide()
	case "esc":
		m.subtypeSelector.Hide()
		// Show species selector again to let user pick different species
		m.speciesSelector.Show()
		m.message = "Subtype selection cancelled"
	}
	return m, nil
}

// handleLanguageSelectorKeys handles language selector specific keys
func (m *Model) handleLanguageSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.languageSelector.Prev()
	case "down", "j":
		m.languageSelector.Next()
	case "enter":
		selectedLanguage := m.languageSelector.GetSelectedLanguage()
		if selectedLanguage != "" {
			// Check if we're in delete mode
			if m.languageSelector.IsDeleteMode() {
				// Remove the language
				for i, lang := range m.character.Languages {
					if lang == selectedLanguage {
						m.character.Languages = append(m.character.Languages[:i], m.character.Languages[i+1:]...)
						break
					}
				}
				m.message = fmt.Sprintf("Language removed: %s", selectedLanguage)
				m.storage.Save(m.character)
				m.languageSelector.Hide()
			} else {
				// Add mode: Check if adding a new language (from Traits panel) or replacing placeholder (from species selection)
				foundPlaceholder := false
			for i, lang := range m.character.Languages {
				if strings.Contains(strings.ToLower(lang), "additional") || strings.Contains(strings.ToLower(lang), "choice") {
					m.character.Languages[i] = selectedLanguage
						foundPlaceholder = true
					break
				}
			}

				// If no placeholder found, just append the new language
				if !foundPlaceholder {
					m.character.Languages = append(m.character.Languages, selectedLanguage)
				}

				// After language selection, check for skill, spell, or feat selection (only during species selection)
			species := models.GetSpeciesByName(m.character.Race)
				if species != nil && foundPlaceholder {
				if models.HasSkillChoice(species) {
					m.skillSelector.Show()
					m.message = "Select your skill proficiency..."
				} else if models.HasSpellChoice(species) {
					// Show wizard cantrip selector for High Elf
					cantrips := models.GetWizardCantrips()
					m.spellSelector.SetSpells(cantrips, "SELECT WIZARD CANTRIP")
					m.spellSelector.Show()
					m.message = "Select your wizard cantrip..."
				} else if models.HasFeatChoice(species) {
					// Show feat selector for origin feat
					m.featSelector.Show(m.character, true)
					m.message = "Select your origin feat..."
				} else {
					m.message = fmt.Sprintf("Language selected: %s (Total languages: %d)", selectedLanguage, len(m.character.Languages))
					// Save when selection is complete (no more selections needed)
					m.storage.Save(m.character)
				}
			} else {
					m.message = fmt.Sprintf("Language learned: %s!", selectedLanguage)
					// Save when adding a new language (not replacing placeholder)
				m.storage.Save(m.character)
		}
		m.languageSelector.Hide()
			}
		}
	case "esc":
		m.languageSelector.Hide()
		m.message = "Language selection cancelled"
	}
	return m, nil
}

// handleToolSelectorKeys handles tool selector specific keys
func (m *Model) handleToolSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.toolSelector.Prev()
	case "down", "j":
		m.toolSelector.Next()
	case "enter":
		selectedTool := m.toolSelector.GetSelected()
		if selectedTool != "" {
			// Check if we're in delete mode
			if m.toolSelector.IsDeleteMode() {
				// Remove the tool proficiency using benefit remover
				// Find and remove from ToolProficiencies
				for i, tool := range m.character.ToolProficiencies {
					if tool == selectedTool {
						m.character.ToolProficiencies = append(m.character.ToolProficiencies[:i], m.character.ToolProficiencies[i+1:]...)
						break
					}
				}

				// Also remove from BenefitTracker (all sources that granted this tool)
				// This is more complex - we need to iterate through all benefits
				allBenefits := m.character.BenefitTracker.Benefits
				for _, benefit := range allBenefits {
					if benefit.Type == models.BenefitTool && benefit.Target == selectedTool {
						m.character.BenefitTracker.RemoveBenefitsBySource(benefit.Source.Type, benefit.Source.Name)
						break
					}
				}

				m.message = fmt.Sprintf("Tool proficiency removed: %s", selectedTool)
				m.storage.Save(m.character)
				m.toolSelector.Hide()
			} else {
				// Add mode: Add tool proficiency directly (not from origin)
				// We'll add it as a "manual" benefit
				source := models.BenefitSource{Type: "manual", Name: "Tool Proficiency"}
				applier := models.NewBenefitApplier(m.character)
				applier.AddToolProficiency(source, selectedTool)

				m.message = fmt.Sprintf("Tool proficiency learned: %s!", selectedTool)
				m.storage.Save(m.character)
				m.toolSelector.Hide()
			}
		}
	case "esc":
		m.toolSelector.Hide()
		m.message = "Tool selection cancelled"
	}
	return m, nil
}

// handleWeaponMasterySelectorKeys handles weapon mastery selector specific keys
func (m *Model) handleWeaponMasterySelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.weaponMasterySelector, cmd = m.weaponMasterySelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case " ":
		// Check if toggle was successful and provide feedback
		if len(m.weaponMasterySelector.GetSelectedWeapons()) >= m.getWeaponMasteryCount() {
			selectedCount := len(m.weaponMasterySelector.GetSelectedWeapons())
			maxCount := m.getWeaponMasteryCount()
			if selectedCount > maxCount {
				m.message = "Maximum weapons already selected"
			}
		}
	case "enter":
		// Confirm selection
		if m.weaponMasterySelector.CanConfirm() {
			m.character.MasteredWeapons = m.weaponMasterySelector.GetSelectedWeapons()
			m.weaponMasterySelector.Hide()

			// Clear pending changes and complete class setup
			m.pendingChanges.Clear()
			m.storage.Save(m.character)

			debug.Log("Weapon mastery selection complete: %v", m.character.MasteredWeapons)
			m.message = fmt.Sprintf("Weapon mastery complete! Mastered %d weapons. Class setup complete. (HP: %d/%d)",
				len(m.character.MasteredWeapons), m.character.CurrentHP, m.character.MaxHP)
		} else {
			m.message = "Please select at least one weapon"
		}
	case "esc":
		// Cancel
		m.weaponMasterySelector.Hide()
		m.message = "Cancelled"
	}

	return m, cmd
}

// handleLevelUpSelectorKeys handles level-up selector keys
func (m *Model) handleLevelUpSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updated, cmd := m.levelUpSelector.Update(msg)
	m.levelUpSelector = &updated

	// Save character if level-up is complete
	if !m.levelUpSelector.IsVisible() {
		m.storage.Save(m.character)
		m.character.UpdateDerivedStats()
		m.message = "Character updated!"
	}

	return m, cmd
}

// handleDeLevelSelectorKeys handles de-level selector keys
func (m *Model) handleDeLevelSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleDeLevelSelectorKeys: key=%s, showConfirm=%v", msg.String(), m.deLevelSelector.IsShowingConfirmation())

	switch msg.String() {
	case "up", "k":
		if !m.deLevelSelector.IsShowingConfirmation() {
			m.deLevelSelector.Prev()
		}
	case "down", "j":
		if !m.deLevelSelector.IsShowingConfirmation() {
			m.deLevelSelector.Next()
		}
	case "enter":
		if m.deLevelSelector.IsShowingConfirmation() {
			// Second enter - execute de-level (already done in ShowConfirmation)
			m.deLevelSelector.ConfirmDeLevel()
			m.storage.Save(m.character)
			m.character.UpdateDerivedStats()

			// Show result message
			previewResult := m.deLevelSelector.GetPreviewResult()
			if previewResult != nil {
				if previewResult.ClassRemoved {
					m.message = fmt.Sprintf("%s class removed entirely (was level %d). Total level: %d",
						previewResult.ClassName, previewResult.OldClassLevel, previewResult.NewTotalLevel)
				} else {
					m.message = fmt.Sprintf("%s level reduced from %d to %d. Total level: %d",
						previewResult.ClassName, previewResult.OldClassLevel, previewResult.NewClassLevel, previewResult.NewTotalLevel)
				}
			} else {
				m.message = "De-level complete!"
			}
		} else {
			// First enter - show confirmation
			err := m.deLevelSelector.ShowConfirmation()
			if err != nil {
				m.message = fmt.Sprintf("Error: %v", err)
				m.deLevelSelector.Hide()
			} else {
				m.message = "Confirm de-level (Enter) or cancel (Esc)"
			}
		}
	case "esc":
		if m.deLevelSelector.IsShowingConfirmation() {
			// Cancel confirmation, go back to class list
			m.deLevelSelector.CancelConfirmation()
			m.message = "De-level cancelled"
		} else {
			// Cancel entire de-level process
			m.deLevelSelector.Hide()
			m.message = "De-level cancelled"
		}
	}

	return m, nil
}

// handleItemSelectorKeys handles item selector specific keys
func (m *Model) handleItemSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle the key once and store the result
	cmd := m.itemSelector.HandleKey(msg)

	// Check if we're confirming quantity entry (only in quantity mode)
	if msg.String() == "enter" && m.itemSelector.IsInQuantityMode() {
		selectedDef, quantity := m.itemSelector.GetSelectedItem()
		if selectedDef != nil {
			// Convert to inventory item and add
			item := models.ConvertToInventoryItem(*selectedDef, quantity)
			m.character.Inventory.AddItem(item)
			m.message = fmt.Sprintf("Added %dx %s to inventory", quantity, item.Name)
			m.storage.Save(m.character)
			m.itemSelector.Hide()
		}
	}

	return m, cmd
}

// handleClassSelectorKeys handles class selector specific keys
func (m *Model) handleClassSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleClassSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.classSelector, cmd = m.classSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		selectedClassName := m.classSelector.GetSelectedClass()
		debug.Log("Class selector: enter pressed, selected=%s", selectedClassName)

		if selectedClassName != "" {
			// Get the full class data to check skill choices
			classData := models.GetClassByName(selectedClassName)
			debug.Log("Class data loaded: %v (nil=%v)", selectedClassName, classData == nil)

			if classData == nil {
				debug.Log("ERROR: Class %s not found!", selectedClassName)
				m.message = fmt.Sprintf("Error: Class %s not found", selectedClassName)
				m.classSelector.Hide()
				return m, nil
			}

			debug.Log("Class %s: SkillChoices=%v", selectedClassName, classData.SkillChoices)
			if classData.SkillChoices != nil {
				debug.Log("  Choose=%d, From=%v", classData.SkillChoices.Choose, classData.SkillChoices.From)
			}

			// Check if class has skill choices
			if classData.SkillChoices != nil && classData.SkillChoices.Choose > 0 {
				// Show skill selector
				debug.Log("Showing skill selector for %s", selectedClassName)
				m.classSelector.Hide()
				m.classSkillSelector.Show(selectedClassName, classData.SkillChoices.From, classData.SkillChoices.Choose, m.character)
				m.message = fmt.Sprintf("Select skills for %s class...", selectedClassName)
			} else {
				// No skill choices, apply class directly
				debug.Log("No skill choices, applying class directly")
				err := models.ApplyClassToCharacter(m.character, selectedClassName)
				if err != nil {
					debug.Log("ERROR applying class: %v", err)
					m.message = fmt.Sprintf("Error applying class: %v", err)
				} else {
					debug.Log("Class applied successfully: %s (HP: %d/%d)", selectedClassName, m.character.CurrentHP, m.character.MaxHP)
					m.message = fmt.Sprintf("Class changed to: %s (HP: %d/%d)", selectedClassName, m.character.CurrentHP, m.character.MaxHP)
				}
				m.storage.Save(m.character)
				m.classSelector.Hide()
			}
		}
	case "esc":
		debug.Log("Class selector: cancelled - restoring previous state")
		// Restore previous class state
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.classSelector.Hide()
		m.storage.Save(m.character) // Save restored state
		m.message = "Class selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSubclassSelectorKeys handles subclass selector specific keys
func (m *Model) handleSubclassSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleSubclassSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.subclassSelector, cmd = m.subclassSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		selectedSubclass := m.subclassSelector.GetSelectedSubclass()
		if selectedSubclass != nil {
			debug.Log("Subclass selected: %s", selectedSubclass.Name)

			// Apply the subclass to the character's current class
			if len(m.character.Classes) > 0 {
				// Update the most recent class (should be the only one at level 1)
				m.character.Classes[len(m.character.Classes)-1].Subclass = selectedSubclass.Name
				debug.Log("Set character subclass to: %s", selectedSubclass.Name)

				// Grant subclass features for the current level
				className := m.character.Classes[len(m.character.Classes)-1].ClassName
				classLevel := m.character.Classes[len(m.character.Classes)-1].Level
				subclassFeatures := models.GrantSubclassFeatures(m.character, className, selectedSubclass.Name, classLevel)
				debug.Log("Granted %d subclass features: %v", len(subclassFeatures), subclassFeatures)
			}

			m.subclassSelector.Hide()

			// Check if we need cantrip selection next
			classData := models.GetClassByName(m.character.Class)
			needsCantrips := classData != nil && classData.Spellcasting != nil && classData.Spellcasting.CantripsKnown > 0
			debug.Log("Class %s needs cantrips: %v", m.character.Class, needsCantrips)

			if needsCantrips {
				// Show cantrip selector
				debug.Log("Showing cantrip selector for %d cantrips", classData.Spellcasting.CantripsKnown)
				m.cantripSelector.Show(m.character.Class, classData.Spellcasting.CantripsKnown)
				m.message = fmt.Sprintf("Select %d cantrips for %s...", classData.Spellcasting.CantripsKnown, m.character.Class)
			} else {
				// Complete class selection
				debug.Log("Saving character and completing class selection")
				m.pendingChanges.Clear()
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Class setup complete! %s - %s (HP: %d/%d)", m.character.Class, selectedSubclass.Name, m.character.CurrentHP, m.character.MaxHP)
			}
		}
	case "esc":
		debug.Log("Subclass selection cancelled")
		// Cancel - rollback
		m.subclassSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.storage.Save(m.character)
		m.message = "Subclass selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleClassSkillSelectorKeys handles class skill selector specific keys
func (m *Model) handleClassSkillSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleClassSkillSelectorKeys: key=%s", msg.String())

	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.classSkillSelector, cmd = m.classSkillSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case " ": // Space to toggle - provide feedback
		debug.Log("Skill selector: selected=%d/%d", len(m.classSkillSelector.SelectedSkills), m.classSkillSelector.MaxChoices)
	case "enter":
		canConfirm := m.classSkillSelector.CanConfirm()
		debug.Log("Skill selector: enter pressed, canConfirm=%v", canConfirm)

		if canConfirm {
			selectedSkills := m.classSkillSelector.GetSelectedSkills()
			selectedClassName := m.classSkillSelector.ClassName
			debug.Log("Applying class %s with skills: %v", selectedClassName, selectedSkills)

			// Note: Old class skills were already removed when backup was created
			// No need to remove them again here

			// Apply the class first
			err := models.ApplyClassToCharacter(m.character, selectedClassName)
			if err != nil {
				debug.Log("ERROR applying class: %v", err)
				m.message = fmt.Sprintf("Error applying class: %v", err)
				m.classSkillSelector.Hide()
				return m, nil
			}
			debug.Log("Class %s applied successfully", selectedClassName)

			// Apply selected skills
			for _, skillName := range selectedSkills {
				skillType := models.SkillType(skillName)
				skill := m.character.Skills.GetSkill(skillType)
				if skill != nil && skill.Proficiency == 0 {
					skill.Proficiency = 1 // Grant proficiency
					debug.Log("Granted proficiency in %s", skillName)
				}
				// Track this skill as coming from class
				m.character.ClassSkills = append(m.character.ClassSkills, skillType)
				debug.Log("Tracked %s as class skill", skillName)
			}

			// Record choices for rollback
			debug.Log("Recording class choice: %s with skills %v", selectedClassName, selectedSkills)
			m.character.Choices.RecordClassChoice(selectedClassName, "")
			m.character.Choices.RecordLevelChoice(1, selectedSkills, "", []string{}, "", nil)

			m.classSkillSelector.Hide()

			// Check if this class needs a subclass at level 1
			classData := models.GetClassByName(selectedClassName)
			needsSubclass := false
			if classData != nil && len(classData.Subclasses) > 0 {
				// Check if any subclass is for level 1 (or has SubclassLevel == 0 which means level 1)
				for _, subclass := range classData.Subclasses {
					if subclass.SubclassLevel == 0 || subclass.SubclassLevel == 1 {
						needsSubclass = true
						break
					}
				}
			}
			debug.Log("Class %s needs subclass at level 1: %v (has %d subclasses)", selectedClassName, needsSubclass, len(classData.Subclasses))

			if needsSubclass {
				// Show subclass selector
				debug.Log("Showing subclass selector for %s", selectedClassName)
				m.subclassSelector.Show(selectedClassName, 1)
				m.message = fmt.Sprintf("Select subclass for %s...", selectedClassName)
			} else {
				// No subclass at level 1, check for fighting style or cantrips
				needsFightingStyle := selectedClassName == "Fighter" || selectedClassName == "Paladin" || selectedClassName == "Ranger"
				debug.Log("Class %s needs fighting style: %v", selectedClassName, needsFightingStyle)

				needsCantrips := classData != nil && classData.Spellcasting != nil && classData.Spellcasting.CantripsKnown > 0
				debug.Log("Class %s needs cantrips: %v", selectedClassName, needsCantrips)

				if needsFightingStyle {
					// Show fighting style selector
					debug.Log("Showing fighting style selector")
					m.fightingStyleSelector.Show(selectedClassName)
					m.message = fmt.Sprintf("Select fighting style for %s...", selectedClassName)
				} else if needsCantrips {
					// Show cantrip selector
					debug.Log("Showing cantrip selector for %d cantrips", classData.Spellcasting.CantripsKnown)
					m.cantripSelector.Show(selectedClassName, classData.Spellcasting.CantripsKnown)
					m.message = fmt.Sprintf("Select %d cantrips for %s...", classData.Spellcasting.CantripsKnown, selectedClassName)
				} else {
					// Check for weapon mastery
					masteryCount := m.getWeaponMasteryCount()
					debug.Log("Class %s needs weapon mastery: count=%d", selectedClassName, masteryCount)

					if masteryCount > 0 {
						// Show weapon mastery selector
						debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
						m.weaponMasterySelector.Show(masteryCount)
						m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
					} else {
						// No fighting style, cantrips, weapon mastery, or subclass needed, complete class selection
						debug.Log("Saving character and completing class selection")
						m.pendingChanges.Clear() // Clear backup on successful completion
						m.storage.Save(m.character)
						m.message = fmt.Sprintf("Class changed to: %s with %d skill proficiencies (HP: %d/%d)", selectedClassName, len(selectedSkills), m.character.CurrentHP, m.character.MaxHP)
					}
				}
			}
		} else {
			debug.Log("Cannot confirm, need more skills")
			m.message = fmt.Sprintf("Please select %d more skill(s)", m.classSkillSelector.MaxChoices-len(m.classSkillSelector.SelectedSkills))
		}
	case "esc":
		debug.Log("Skill selector: cancelled - restoring previous state")
		// Restore previous class state
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.classSkillSelector.Hide()
		m.storage.Save(m.character) // Save restored state
		m.message = "Skill selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleCantripSelectorKeys handles cantrip selector specific keys
func (m *Model) handleCantripSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleCantripSelectorKeys: key=%s", msg.String())

	// Delegate navigation and selection to the component's Update method
	var cmd tea.Cmd
	*m.cantripSelector, cmd = m.cantripSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		if m.cantripSelector.CanConfirm() {
			selectedCantrips := m.cantripSelector.GetSelectedCantrips()
			debug.Log("Selected cantrips: %v", selectedCantrips)

			// Apply cantrips to character's spellbook
			m.character.SpellBook.Cantrips = selectedCantrips
			m.character.SpellBook.CantripsKnown = len(selectedCantrips)

			m.cantripSelector.Hide()

			// Check for weapon mastery
			masteryCount := m.getWeaponMasteryCount()
			debug.Log("After cantrips, checking weapon mastery: count=%d", masteryCount)

			if masteryCount > 0 {
				// Show weapon mastery selector
				debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
				m.weaponMasterySelector.Show(masteryCount)
				m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
			} else {
				// Complete class selection
				debug.Log("Saving character and completing class selection")
				m.pendingChanges.Clear() // Clear backup on successful completion
				m.storage.Save(m.character)
				m.message = fmt.Sprintf("Class selection complete! Selected %d cantrips", len(selectedCantrips))
			}
		} else {
			needed := m.cantripSelector.GetMaxCantrips() - m.cantripSelector.GetSelectedCount()
			m.message = fmt.Sprintf("Please select %d more cantrip(s)", needed)
		}
	case "esc":
		// Cancel - rollback to previous state
		debug.Log("Cantrip selection cancelled, rolling back changes")
		m.cantripSelector.Hide()
		m.pendingChanges.RestoreClass(m.character)
		m.storage.Save(m.character)
		m.message = "Class selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSpellPrepSelectorKeys handles keyboard input for the spell prep selector
func (m *Model) handleSpellPrepSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	*m.spellPrepSelector, cmd = m.spellPrepSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		m.spellPrepSelector.Hide()
		m.storage.Save(m.character)
		m.message = "Spell preparation updated!"
	case "esc":
		m.spellPrepSelector.Hide()
		m.message = "Spell preparation cancelled"
	}

	return m, cmd
}

// handleSlotRestorerKeys handles keyboard input for the slot restorer
func (m *Model) handleSlotRestorerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	*m.slotRestorer, cmd = m.slotRestorer.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter", " ":
		// Restore one slot
		restoredMsg, success := m.slotRestorer.RestoreSlot()
		if success {
			m.storage.Save(m.character)
			m.message = restoredMsg
		} else {
			m.message = restoredMsg
		}
	case "esc":
		m.slotRestorer.Hide()
		m.message = "Slot restoration cancelled"
	}

	return m, cmd
}

// handleFightingStyleSelectorKeys handles fighting style selector specific keys
func (m *Model) handleFightingStyleSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debug.Log("handleFightingStyleSelectorKeys: key=%s", msg.String())

	// Delegate navigation to the component's Update method
	var cmd tea.Cmd
	*m.fightingStyleSelector, cmd = m.fightingStyleSelector.Update(tea.KeyMsg(msg))

	switch msg.String() {
	case "enter":
		selectedStyle := m.fightingStyleSelector.GetSelectedStyle()
		debug.Log("Fighting style selector: enter pressed, selected=%s", selectedStyle)

		if selectedStyle != "" {
			// Apply fighting style
			err := models.ApplyFightingStyle(m.character, selectedStyle)
			if err != nil {
				debug.Log("ERROR applying fighting style: %v", err)
				m.message = fmt.Sprintf("Error applying fighting style: %v", err)
			} else {
				debug.Log("Fighting style '%s' applied successfully", selectedStyle)
				// Update the choice record with fighting style
				m.character.Choices.Class.FightingStyle = selectedStyle

				// Check if character also needs weapon mastery selection
				masteryCount := m.getWeaponMasteryCount()
				debug.Log("After fighting style, checking weapon mastery: count=%d", masteryCount)

				if masteryCount > 0 {
					// Show weapon mastery selector
					debug.Log("Showing weapon mastery selector for %d weapons", masteryCount)
					m.weaponMasterySelector.Show(masteryCount)
					m.message = fmt.Sprintf("Select up to %d weapons to master...", masteryCount)
					m.fightingStyleSelector.Hide()
				} else {
					// No weapon mastery needed, class setup complete
					m.pendingChanges.Clear() // Clear backup on successful completion
					m.message = fmt.Sprintf("Fighting style '%s' selected! Class setup complete. (HP: %d/%d)", selectedStyle, m.character.CurrentHP, m.character.MaxHP)
					m.fightingStyleSelector.Hide()
				}
			}
			m.storage.Save(m.character)
		}
	case "esc":
		debug.Log("Fighting style selector: cancelled - restoring previous state")
		// Restore previous class state
		m.pendingChanges.RestoreClass(m.character)
		m.pendingChanges.Clear()
		m.fightingStyleSelector.Hide()
		m.storage.Save(m.character) // Save restored state
		m.message = "Fighting style selection cancelled - restored previous state"
	}

	return m, cmd
}

// handleSkillSelectorKeys handles skill selector specific keys
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
		m.skillSelector.Hide()
	case "esc":
		m.skillSelector.Hide()
		m.message = "Skill selection cancelled"
	}
	return m, nil
}

// handleSpellSelectorKeys handles spell selector specific keys
func (m *Model) handleSpellSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.spellSelector.Prev()
	case "down", "j":
		m.spellSelector.Next()
	case "enter":
		selectedSpell := m.spellSelector.GetSelectedSpell()
		if selectedSpell.Name != "" {
			// Check if character already has this spell
			hasSpell := false
			for _, existing := range m.character.SpellBook.Spells {
				if existing.Name == selectedSpell.Name {
					hasSpell = true
					break
				}
			}

			if !hasSpell {
				// Add spell to spellbook and track it as a species spell
				m.character.SpellBook.AddSpell(selectedSpell)
				m.character.SpeciesSpells = append(m.character.SpeciesSpells, selectedSpell.Name)
				m.message = fmt.Sprintf("Spell learned: %s", selectedSpell.Name)
			} else {
				m.message = fmt.Sprintf("You already know %s", selectedSpell.Name)
			}

			// After spell selection, check if we need feat selection
			species := models.GetSpeciesByName(m.character.Race)
			if species != nil && models.HasFeatChoice(species) {
				m.spellSelector.Hide()
				// Show feat selector for origin feat
				m.featSelector.Show(m.character, true)
				m.message = "Select your origin feat..."
				return m, nil
			}

			// Save character after spell selection (final step)
			m.storage.Save(m.character)
		}
		m.spellSelector.Hide()
	case "esc":
		m.spellSelector.Hide()
		m.message = "Spell selection cancelled"
	}
	return m, nil
}

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

// handleFeatDetailPopupKeys handles keyboard input for the feat detail popup
func (m *Model) handleFeatDetailPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.featDetailPopup.Hide()
		m.message = "Closed feat details"
	}
	return m, nil
}

func (m *Model) handleItemDetailPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.itemDetailPopup.Hide()
		m.message = "Closed item details"
	}
	return m, nil
}

// handleSpellDetailPopupKeys handles keyboard input for the spell detail popup
func (m *Model) handleSpellDetailPopupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.spellDetailPopup.Hide()
		m.message = "Closed spell details"
	}
	return m, nil
}

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
			return "Traits", components.GetTraitsBindings()
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
			contextHelp = "[o] Change Origin • [t] Add Tool • [T] Remove Tool"
		}
	case FocusCharStats:
		panelName = "Character Info"
		contextHelp = "[n] Name • [r] Species • [h] HP • [+/-] ±1 • [i] Init"
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

	// Calculate heights to fit exactly within screen
	// Distribution: Line 1 (45%), Line 2 (45%), Status bar (4%), Gaps (6%)
	statusBarHeight := int(float64(m.height) * 0.04)
	if statusBarHeight < 1 {
		statusBarHeight = 1
	}

	// First row: Main panel (55%) + Character stats (43%)
	mainPanelWidth := int(float64(m.width) * 0.55)
	charStatsWidth := int(float64(m.width) * 0.43)

	// Calculate panel heights: top row 48%, bottom row 42%
	topRowHeight := int(float64(m.height) * 0.48)
	bottomHeight := int(float64(m.height) * 0.42)

	// Ensure minimum heights
	if topRowHeight < 10 {
		topRowHeight = 10
	}
	if bottomHeight < 8 {
		bottomHeight = 8
	}

	// Tab navigation (width accounts for border and padding)
	tabBarWidth := mainPanelWidth - 8 // Account for border (2) + horizontal padding (4)
	tabBar := m.tabs.View(tabBarWidth)
	tabHeight := lipgloss.Height(tabBar)

	// Main content height accounts for tabs and spacing, to fill the full topRowHeight
	// topRowHeight includes border and padding in the final render
	mainContentHeight := topRowHeight - tabHeight - 5 // border (2) + padding vertical (2) + spacing line (1)

	// Main panel content
	var mainPanelView string
	mainWidth := mainPanelWidth - 8 // Account for border + padding (2 for border, 4 for padding)

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

	// Combine tabs and content vertically (tabs inside the panel)
	tabsAndContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		"", // Add a line of spacing
		mainPanelView,
	)

	// Add border to combined tabs + main panel (focused = pink, unfocused = gray)
	// Set explicit height and width to fill space properly
	mainPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2). // Match skills/inventory padding
		Width(mainPanelWidth).
		Height(topRowHeight)

	if m.focusArea == FocusMain {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("240"))
	}

	mainPanelWithTabs := mainPanelStyle.Render(tabsAndContent)

	// Character stats panel (always visible, 45% of width)
	// Height should match the main panel exactly
	charStatsInnerWidth := charStatsWidth - 8 // Account for border + padding (2 for border, 4 for padding)
	charStatsInnerHeight := topRowHeight - 6 // Account for border (2) + vertical padding (4)
	charStatsView := m.characterStatsPanel.View(charStatsInnerWidth, charStatsInnerHeight)
	charStatsPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2). // Match skills/inventory padding
		Width(charStatsWidth).
		Height(topRowHeight)

	if m.focusArea == FocusCharStats {
		charStatsPanelStyle = charStatsPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		charStatsPanelStyle = charStatsPanelStyle.BorderForeground(lipgloss.Color("86"))
	}
	charStatsWithBorder := charStatsPanelStyle.Render(charStatsView)

	// Join main panel (with tabs) and character stats horizontally
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainPanelWithTabs,
		charStatsWithBorder,
	)

	// Bottom panels: Actions (50%) + Dice Roller (48%)
	actionsWidthRatio := int(float64(m.width) * 0.50)
	diceWidthRatio := int(float64(m.width) * 0.48)

	// Calculate inner dimensions: account for border + padding (2 for border, 4 for padding)
	actionsWidth := actionsWidthRatio - 8
	diceWidth := diceWidthRatio - 8
	bottomInnerHeight := bottomHeight - 6 // Account for border (2) + vertical padding (4)

	// Actions panel with border (focused = pink, unfocused = gray)
	actionsView := m.actionsPanel.View(actionsWidth, bottomInnerHeight)
	actionsPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2). // Match skills/inventory padding
		Width(actionsWidthRatio). // Set explicit width to fill space
		Height(bottomHeight) // Enforce 45% height

	if m.focusArea == FocusActions {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("240"))
	}
	actionsView = actionsPanelStyle.Render(actionsView)

	// Dice panel with border (focused = pink, unfocused = gray)
	diceView := m.dicePanel.View(diceWidth, bottomInnerHeight)
	dicePanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2). // Match skills/inventory padding
		Width(diceWidthRatio). // Set explicit width to fill space
		Height(bottomHeight) // Enforce 45% height

	if m.focusArea == FocusDice {
		dicePanelStyle = dicePanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		dicePanelStyle = dicePanelStyle.BorderForeground(lipgloss.Color("240"))
	}
	diceView = dicePanelStyle.Render(diceView)

	// Bottom row (actions + dice)
	bottomRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		actionsView,
		diceView,
	)

	// Status bar
	statusBar := m.buildStatusBar()

	// Combine all parts vertically
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		bottomRow,
		statusBar,
	)

	// Render popups/overlays (in priority order)
	// Calculate popup dimensions based on size category

	// Small popup dimensions (50% width, 60% height)
	popupSmallWidth := max(int(float64(m.width)*PopupSmallWidthPercent), PopupSmallMinWidth)
	popupSmallHeight := max(int(float64(m.height)*PopupSmallHeightPercent), PopupSmallMinHeight)

	// Medium popup dimensions (75% width, 80% height)
	popupMediumWidth := max(int(float64(m.width)*PopupMediumWidthPercent), PopupMediumMinWidth)
	popupMediumHeight := max(int(float64(m.height)*PopupMediumHeightPercent), PopupMediumMinHeight)

	// Large popup dimensions (85% width, 85% height)
	popupLargeWidth := max(int(float64(m.width)*PopupLargeWidthPercent), PopupLargeMinWidth)
	popupLargeHeight := max(int(float64(m.height)*PopupLargeHeightPercent), PopupLargeMinHeight)

	// Stat generator takes highest priority (Medium)
	if m.statGenerator.IsVisible() {
		return m.statGenerator.View(popupMediumWidth, popupMediumHeight)
	}

	// Ability roller takes high priority (Small)
	if m.abilityRoller.IsVisible() {
		return m.abilityRoller.View(popupSmallWidth, popupSmallHeight, m.character)
	}

	// Attack roller takes high priority (Medium)
	if m.attackRoller.IsVisible() {
		return m.attackRoller.View(popupMediumWidth, popupMediumHeight)
	}

	// Attack menu should show AFTER checking all full-screen popups
	// but BEFORE returning mainView

	// Spell selector takes high priority (Large)
	if m.spellSelector.IsVisible() {
		return m.spellSelector.View(popupLargeWidth, popupLargeHeight)
	}

	// Feat selector takes second priority (Medium)
	if m.featSelector.IsVisible() {
		return m.featSelector.View(popupMediumWidth, popupMediumHeight)
	}

	// Feat detail popup (Medium)
	if m.featDetailPopup.IsVisible() {
		return m.featDetailPopup.View(popupMediumWidth, popupMediumHeight)
	}

	// Item detail popup (Medium)
	if m.itemDetailPopup.IsVisible() {
		return m.itemDetailPopup.View(m.width, m.height)
	}

	// Spell detail popup (Medium)
	if m.spellDetailPopup.IsVisible() {
		return m.spellDetailPopup.View()
	}

	// Origin selector (Medium)
	if m.originSelector.IsVisible() {
		return m.originSelector.View(popupMediumWidth, popupMediumHeight)
	}

	// Ability choice selector (for feat ability choices) (Small)
	if m.abilityChoiceSelector.IsVisible() {
		return m.abilityChoiceSelector.View(popupSmallWidth, popupSmallHeight)
	}

	// Subtype selector takes third priority (Small)
	if m.subtypeSelector.IsVisible() {
		return m.subtypeSelector.View(popupSmallWidth, popupSmallHeight)
	}

	// Skill selector takes fourth priority (Small)
	if m.skillSelector.IsVisible() {
		return m.skillSelector.View(popupSmallWidth, popupSmallHeight)
	}

	// Language selector takes third priority (Small)
	if m.languageSelector.IsVisible() {
		return m.languageSelector.View(popupSmallWidth, popupSmallHeight)
	}

	// Tool selector takes fourth priority (Small)
	if m.toolSelector.IsVisible() {
		return m.toolSelector.View(popupSmallWidth, popupSmallHeight)
	}

	// Weapon mastery selector takes fifth priority (Medium)
	if m.weaponMasterySelector.IsVisible() {
		return m.weaponMasterySelector.View()
	}

	// Level-up selector takes sixth priority (Medium/Large)
	if m.levelUpSelector.IsVisible() {
		return m.levelUpSelector.View()
	}

	// De-level selector takes priority after level-up (Medium)
	if m.deLevelSelector.IsVisible() {
		return m.deLevelSelector.View(m.width, m.height)
	}

	// Item selector takes seventh priority (Large)
	if m.itemSelector.IsVisible() {
		return m.itemSelector.View(popupLargeWidth, popupLargeHeight)
	}

	// Fighting style selector takes sixth priority (Medium)
	if m.fightingStyleSelector.IsVisible() {
		return m.fightingStyleSelector.View(m.width, m.height)
	}

	// Cantrip selector takes seventh priority (Medium)
	if m.cantripSelector.IsVisible() {
		return m.cantripSelector.View()
	}

	// Spell prep selector takes eighth priority (Medium)
	if m.spellPrepSelector.IsVisible() {
		return m.spellPrepSelector.View()
	}

	// Slot restorer takes ninth priority (Small)
	if m.slotRestorer.IsVisible() {
		return m.slotRestorer.View()
	}

	// Class skill selector takes tenth priority (Medium)
	if m.classSkillSelector.IsVisible() {
		return m.classSkillSelector.View(m.width, m.height)
	}

	// Subclass selector takes priority after skills (Medium)
	if m.subclassSelector.IsVisible() {
		return m.subclassSelector.View()
	}

	// Class selector takes seventh priority (Medium)
	if m.classSelector.IsVisible() {
		return m.classSelector.View(popupMediumWidth, popupMediumHeight)
	}

	// Species selector takes seventh priority (Medium)
	if m.speciesSelector.IsVisible() {
		return m.speciesSelector.View(popupMediumWidth, popupMediumHeight)
	}

	// HP popup overlay if active (Small)
	hpPopup := m.characterStatsPanel.RenderHPPopup(popupSmallWidth, popupSmallHeight)
	if hpPopup != "" {
		return hpPopup
	}

	// XP popup overlay if active (Small)
	xpPopup := m.characterStatsPanel.RenderXPPopup(popupSmallWidth, popupSmallHeight)
	if xpPopup != "" {
		return xpPopup
	}

	// Attack menu takes priority (shows as centered overlay)
	// Note: This will hide the TUI underneath for simplicity
	if m.attackMenu.IsVisible() {
		return m.attackMenu.View(m.width, m.height)
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
