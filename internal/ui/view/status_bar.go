// internal/ui/view/status_bar.go
package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
)

// StatusBarContext provides the context needed to build a status bar
type StatusBarContext struct {
	FocusArea    int
	CurrentPanel PanelType
	DiceMode     panels.DicePanelMode
	Width        int
}

// BuildStatusBar creates the status bar with contextual information
func BuildStatusBar(ctx StatusBarContext) string {
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

	const FocusMain = 0
	const FocusCharStats = 1
	const FocusActions = 2
	const FocusDice = 3

	switch ctx.FocusArea {
	case FocusMain:
		switch ctx.CurrentPanel {
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
			contextHelp = "[o] Origin • [Enter] Details • [a] Alignment • [h/w] Height/Weight • [t/i/b/f] Traits • [s] Species • [Shift+S] Save"
		}
	case FocusCharStats:
		panelName = "Character Info"
		contextHelp = "[n] Name • [s] Species • [Shift+S] Save • [h] HP • [r] Short Rest • [R] Long Rest • [+/-] ±1 • [i] Init"
	case FocusActions:
		panelName = "Actions"
		contextHelp = "[↑/↓] Navigate • [Enter] Activate"
	case FocusDice:
		panelName = "Dice Roller"
		switch ctx.DiceMode {
		case panels.DiceModeIdle:
			contextHelp = "[Enter] Input • [h] History • [r] Reroll"
		case panels.DiceModeInput:
			contextHelp = "Type dice notation • [Enter] Roll • [Esc] Cancel"
		case panels.DiceModeHistory:
			contextHelp = "[↑/↓] Navigate • [Enter] Reroll • [Esc] Back"
		}
	}

	// Build left section: panel + help
	leftSection := panelNameStyle.Render(" " + panelName + " ")

	if contextHelp != "" {
		leftSection += helpStyle.Render(" " + contextHelp + " ")
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
	padding := ctx.Width - leftWidth - rightWidth
	if padding < 0 {
		padding = 0
	}

	paddingStr := strings.Repeat(" ", padding)

	statusBarStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Width(ctx.Width)

	return statusBarStyle.Render(leftSection + paddingStr + rightSection)
}

// GetContextualHelp returns the panel name and contextual help bindings based on current focus
func GetContextualHelp(
	focusArea int,
	currentPanel PanelType,
	diceMode panels.DicePanelMode,
	isBattleMaster bool,
	hasManeuvers bool,
) (string, []components.HelpBinding) {
	const FocusMain = 0
	const FocusCharStats = 1
	const FocusActions = 2
	const FocusDice = 3

	switch focusArea {
	case FocusMain:
		switch currentPanel {
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
			if isBattleMaster && hasManeuvers {
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
		switch diceMode {
		case panels.DiceModeInput:
			mode = "input"
		case panels.DiceModeHistory:
			mode = "history"
		}
		return "Dice Roller", components.GetDiceBindings(mode)
	}
	return "Stats", components.GetStatsBindings()
}
