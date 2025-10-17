// internal/ui/app.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
)

// PanelType represents the current active main panel
type PanelType int

const (
	OverviewPanel PanelType = iota
	StatsPanel
	SkillsPanel
	InventoryPanel
	SpellsPanel
)

// FocusArea represents which area of the UI has focus
type FocusArea int

const (
	FocusMain FocusArea = iota
	FocusActions
	FocusDice
)

// Model is the main application model
type Model struct {
	character    *models.Character
	storage      *storage.Storage

	// UI Components
	tabs         *components.Tabs
	help         *components.Help

	// Main Panels (switchable)
	overviewPanel  *panels.OverviewPanel
	statsPanel     *panels.StatsPanel
	skillsPanel    *panels.SkillsPanel
	inventoryPanel *panels.InventoryPanel
	spellsPanel    *panels.SpellsPanel

	// Fixed Panels (always visible)
	actionsPanel   *panels.ActionsPanel
	dicePanel      *panels.DicePanel

	// State
	currentPanel PanelType
	focusArea    FocusArea
	width        int
	height       int
	ready        bool
	message      string
	quitting     bool
}

// NewModel creates a new application model
func NewModel(char *models.Character, store *storage.Storage) *Model {
	return &Model{
		character:      char,
		storage:        store,
		tabs:           components.NewTabs(),
		help:           components.NewHelp(),
		overviewPanel:  panels.NewOverviewPanel(char),
		statsPanel:     panels.NewStatsPanel(char),
		skillsPanel:    panels.NewSkillsPanel(char),
		inventoryPanel: panels.NewInventoryPanel(char),
		spellsPanel:    panels.NewSpellsPanel(char),
		actionsPanel:   panels.NewActionsPanel(char),
		dicePanel:      panels.NewDicePanel(char),
		currentPanel:   OverviewPanel,
		focusArea:      FocusMain,
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

		// Handle dice panel input when focused
		if m.dicePanel.GetInput() != "" {
			switch msg.String() {
			case "enter":
				m.dicePanel.Roll(m.dicePanel.GetInput())
				return m, nil
			case "esc":
				m.dicePanel.Blur()
				return m, nil
			default:
				return m, m.dicePanel.Update(msg)
			}
		}

		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "?":
			m.help.Toggle()
			return m, nil

		case "s":
			err := m.storage.Save(m.character)
			if err != nil {
				m.message = fmt.Sprintf("Error saving: %s", err.Error())
			} else {
				m.message = "Character saved!"
			}
			return m, nil

		case "l":
			if m.character.CanLevelUp() {
				m.message = "Level up! (Feature not fully implemented in this version)"
				// TODO: Implement level up wizard
			} else {
				m.message = "Not enough XP to level up"
			}
			return m, nil

		// Focus cycling - f key cycles through Main, Actions, Dice
		case "f":
			m.focusArea = (m.focusArea + 1) % 3
			switch m.focusArea {
			case FocusMain:
				m.message = "Focus: Main Panel"
			case FocusActions:
				m.message = "Focus: Actions Panel"
			case FocusDice:
				m.message = "Focus: Dice Roller"
			}
			return m, nil

		// Panel navigation (only when focused on main)
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

		// Number keys for quick panel selection (only when focused on main)
		case "1", "2", "3", "4", "5":
			if m.focusArea == FocusMain {
				idx := int(msg.String()[0] - '1')
				m.tabs.SetIndex(idx)
				m.currentPanel = PanelType(idx)
			}
			return m, nil
		}

		// Handle input based on current focus
		switch m.focusArea {
		case FocusMain:
			return m.handleMainPanelKeys(msg)
		case FocusActions:
			return m.handleActionsPanelKeys(msg)
		case FocusDice:
			return m.handleDicePanelKeys(msg)
		}

		// Global key 'r' for rest (affects actions and spells)
		switch msg.String() {
		case "R": // Shift+R for rest
			m.character.LongRest()
			m.message = "Long rest completed! HP, spells, and abilities restored."
			return m, nil
		}

		// Handle dice shortcuts (d + number for quick rolls)
		if len(msg.String()) == 2 && msg.String()[0] == 'd' {
			switch msg.String()[1] {
			case '1':
				m.dicePanel.RollQuick("d4")
				m.message = "Rolled d4"
			case '2':
				m.dicePanel.RollQuick("d6")
				m.message = "Rolled d6"
			case '3':
				m.dicePanel.RollQuick("d8")
				m.message = "Rolled d8"
			case '4':
				m.dicePanel.RollQuick("d10")
				m.message = "Rolled d10"
			case '5':
				m.dicePanel.RollQuick("d12")
				m.message = "Rolled d12"
			case '6':
				m.dicePanel.RollQuick("d20")
				m.message = "Rolled d20"
			case '7':
				m.dicePanel.RollQuick("d100")
				m.message = "Rolled d100"
			}
			return m, nil
		}
	}

	return m, nil
}

// handleMainPanelKeys handles keys when main panel has focus
func (m *Model) handleMainPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.currentPanel {
	case SkillsPanel:
		return m.handleSkillsPanel(msg)
	case InventoryPanel:
		return m.handleInventoryPanel(msg)
	case SpellsPanel:
		return m.handleSpellsPanel(msg)
	}
	return m, nil
}

// handleActionsPanelKeys handles keys when actions panel has focus
func (m *Model) handleActionsPanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.actionsPanel.Prev()
	case "down", "j":
		m.actionsPanel.Next()
	case "enter":
		m.message = "Action activated (not fully implemented)"
	}
	return m, nil
}

// handleDicePanelKeys handles keys when dice panel has focus
func (m *Model) handleDicePanelKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.dicePanel.Focus()
		return m, m.dicePanel.Update(msg)
	case "1":
		m.dicePanel.RollQuick("d4")
		m.message = "Rolled d4"
	case "2":
		m.dicePanel.RollQuick("d6")
		m.message = "Rolled d6"
	case "3":
		m.dicePanel.RollQuick("d8")
		m.message = "Rolled d8"
	case "4":
		m.dicePanel.RollQuick("d10")
		m.message = "Rolled d10"
	case "5":
		m.dicePanel.RollQuick("d12")
		m.message = "Rolled d12"
	case "6":
		m.dicePanel.RollQuick("d20")
		m.message = "Rolled d20"
	case "7":
		m.dicePanel.RollQuick("d100")
		m.message = "Rolled d100"
	case "n":
		m.dicePanel.SetRollType("normal")
		m.message = "Roll type: Normal"
	case "a":
		m.dicePanel.SetRollType("advantage")
		m.message = "Roll type: Advantage"
	case "d":
		m.dicePanel.SetRollType("disadvantage")
		m.message = "Roll type: Disadvantage"
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
	case "e":
		m.inventoryPanel.ToggleEquipped()
		m.message = "Item equipped status toggled"
	case "d":
		m.inventoryPanel.DeleteSelected()
		m.message = "Item deleted"
	case "a":
		// Add a sample item for demonstration
		m.character.Inventory.AddItem(models.Item{
			Name:     "New Item",
			Type:     models.Gear,
			Quantity: 1,
			Weight:   1.0,
		})
		m.message = "Item added (edit with 'e' - not fully implemented)"
	}
	return m, nil
}


// handleSpellsPanel handles spells panel specific keys
func (m *Model) handleSpellsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		m.character.SpellBook.LongRest()
		m.message = "Spell slots restored!"
	case "a":
		// Add a sample spell for demonstration
		m.character.SpellBook.AddSpell(models.Spell{
			Name:        "New Spell",
			Level:       1,
			School:      models.Evocation,
			CastingTime: "1 action",
			Range:       "60 feet",
			Components:  "V, S",
			Duration:    "Instantaneous",
			Description: "A new spell",
			Prepared:    false,
			Known:       true,
		})
		m.message = "Spell added (edit - not fully implemented)"
	}
	return m, nil
}


// View renders the application
func (m *Model) View() string {
	if !m.ready || m.quitting {
		return ""
	}

	// Show help overlay if visible
	if m.help.Visible {
		return m.help.View(m.width, m.height)
	}

	// Title bar
	titleText := fmt.Sprintf(" LazyDnDPlayer - %s (%s) - Level %d ",
		m.character.Name, m.character.Class, m.character.Level)
	title := TitleStyle.Width(m.width).Render(titleText)

	// Tab navigation
	tabBar := m.tabs.View(m.width)

	// Calculate heights
	titleHeight := 1
	tabHeight := 3
	bottomHeight := 20 // Height for actions + dice panels
	statusHeight := 1
	mainPanelHeight := m.height - titleHeight - tabHeight - bottomHeight - statusHeight - 2

	// Main panel (full width at top)
	var mainPanelView string

	// Adjust width for border if focused
	mainWidth := m.width - 4
	if m.focusArea == FocusMain {
		mainWidth = m.width - 8 // Account for border padding
	}

	switch m.currentPanel {
	case OverviewPanel:
		mainPanelView = m.overviewPanel.View(mainWidth, mainPanelHeight)
	case StatsPanel:
		mainPanelView = m.statsPanel.View(mainWidth, mainPanelHeight)
	case SkillsPanel:
		mainPanelView = m.skillsPanel.View(mainWidth, mainPanelHeight)
	case InventoryPanel:
		mainPanelView = m.inventoryPanel.View(mainWidth, mainPanelHeight)
	case SpellsPanel:
		mainPanelView = m.spellsPanel.View(mainWidth, mainPanelHeight)
	}

	// Add focus border to main panel
	if m.focusArea == FocusMain {
		mainPanelStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)
		mainPanelView = mainPanelStyle.Render(mainPanelView)
	}

	// Bottom panels (split 50/50)
	bottomWidth := m.width / 2

	// Adjust widths for borders
	actionsWidth := bottomWidth - 2
	diceWidth := bottomWidth - 2
	if m.focusArea == FocusActions {
		actionsWidth = bottomWidth - 6
	}
	if m.focusArea == FocusDice {
		diceWidth = bottomWidth - 6
	}

	// Actions panel with focus indicator
	actionsView := m.actionsPanel.View(actionsWidth, bottomHeight)
	if m.focusArea == FocusActions {
		actionsPanelStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)
		actionsView = actionsPanelStyle.Render(actionsView)
	}

	// Dice panel with focus indicator
	diceView := m.dicePanel.View(diceWidth, bottomHeight)
	if m.focusArea == FocusDice {
		dicePanelStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)
		diceView = dicePanelStyle.Render(diceView)
	}

	// Bottom row (actions + dice)
	bottomRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		actionsView,
		diceView,
	)

	// Status bar
	statusBar := ""
	if m.message != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Background(lipgloss.Color("235")).
			Width(m.width).
			Padding(0, 1)
		statusBar = statusStyle.Render(m.message)
	}

	// Combine all parts vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		tabBar,
		mainPanelView,
		bottomRow,
		statusBar,
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
