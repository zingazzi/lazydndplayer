// internal/ui/app.go
package ui

import (
	"fmt"
	"strings"

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
				return m, nil
			}
			// If not in main focus, let the focused panel handle it
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


// buildStatusBar creates the status bar with contextual information
func (m *Model) buildStatusBar() string {
	appNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("235"))

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
		case OverviewPanel:
			panelName = "Overview"
			contextHelp = ""
		case StatsPanel:
			panelName = "Stats"
			contextHelp = "[e] Edit"
		case SkillsPanel:
			panelName = "Skills"
			contextHelp = "[↑/↓] Navigate • [r] Roll • [e] Toggle Prof"
		case InventoryPanel:
			panelName = "Inventory"
			contextHelp = "[↑/↓] Navigate • [a] Add • [e] Equip • [d] Delete"
		case SpellsPanel:
			panelName = "Spells"
			contextHelp = "[a] Add • [r] Rest"
		}
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

	// Build left section: app name + panel + help
	leftSection := appNameStyle.Render(" lazydndplayer ") +
		panelNameStyle.Render(" "+panelName+" ")

	if contextHelp != "" {
		leftSection += helpStyle.Render(" "+contextHelp+" ")
	}

	// Build right section: global shortcuts
	rightSection := keyStyle.Render("[Tab]") + helpStyle.Render(" Switch tabs • ") +
		keyStyle.Render("[f]") + helpStyle.Render(" Focus • ") +
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
		return m.help.View(m.width, m.height)
	}

	// Title bar
	titleText := fmt.Sprintf(" LazyDnDPlayer - %s (%s) - Level %d ",
		m.character.Name, m.character.Class, m.character.Level)
	title := TitleStyle.Width(m.width).Render(titleText)

	// Tab navigation
	tabBar := m.tabs.View(m.width)

	// Calculate heights (better proportions)
	titleHeight := 1
	tabHeight := 3
	statusBarHeight := 1

	// Use proportional heights: 55% for main panel, 45% for bottom panels
	availableHeight := m.height - titleHeight - tabHeight - statusBarHeight - 2
	mainPanelHeight := int(float64(availableHeight) * 0.55)
	bottomHeight := availableHeight - mainPanelHeight

	// Ensure minimum heights
	if mainPanelHeight < 15 {
		mainPanelHeight = 15
	}
	if bottomHeight < 12 {
		bottomHeight = 12
	}

	// Main panel (full width at top)
	var mainPanelView string

	// Always account for padding
	mainWidth := m.width - 8

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

	// Add border (focused = pink, unfocused = gray)
	mainPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1)

	if m.focusArea == FocusMain {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("240"))
	}

	mainPanelView = mainPanelStyle.Render(mainPanelView)

	// Bottom panels (split 50/50)
	bottomWidth := m.width / 2

	// Always account for border padding
	actionsWidth := bottomWidth - 6
	diceWidth := bottomWidth - 6

	// Actions panel with border (focused = pink, unfocused = gray)
	actionsView := m.actionsPanel.View(actionsWidth, bottomHeight)
	actionsPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1)

	if m.focusArea == FocusActions {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("240"))
	}
	actionsView = actionsPanelStyle.Render(actionsView)

	// Dice panel with border (focused = pink, unfocused = gray)
	diceView := m.dicePanel.View(diceWidth, bottomHeight)
	dicePanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1)

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
