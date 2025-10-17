// internal/ui/styles.go
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("205") // Pink
	secondaryColor = lipgloss.Color("170") // Purple
	accentColor    = lipgloss.Color("86")  // Cyan
	successColor   = lipgloss.Color("42")  // Green
	dangerColor    = lipgloss.Color("196") // Red
	warningColor   = lipgloss.Color("214") // Orange
	mutedColor     = lipgloss.Color("240") // Gray

	// Title bar style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(primaryColor).
			Padding(0, 1)

	// Sidebar styles
	SidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	SidebarItemSelectedStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(primaryColor).
					Background(lipgloss.Color("237"))

	SidebarKeyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Panel styles
	PanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2)

	PanelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	// Content styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	ModifierPositiveStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	ModifierNegativeStyle = lipgloss.NewStyle().
				Foreground(dangerColor).
				Bold(true)

	// Status styles
	HPStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)

	HPLowStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	HPCriticalStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	// List item styles
	ListItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 2)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("237")).
				Padding(0, 2)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1)

	// Help styles
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Dice roll styles
	DiceRollStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Margin(0, 1)

	DiceRollCritStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Background(lipgloss.Color("235")).
				Bold(true).
				Padding(0, 1).
				Margin(0, 1)

	DiceRollFailStyle = lipgloss.NewStyle().
				Foreground(dangerColor).
				Background(lipgloss.Color("235")).
				Bold(true).
				Padding(0, 1).
				Margin(0, 1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)
)

// FormatModifier formats an ability modifier with + or -
func FormatModifier(mod int) string {
	style := ModifierPositiveStyle
	if mod < 0 {
		style = ModifierNegativeStyle
		return style.Render("-" + lipgloss.NewStyle().Render(fmt.Sprintf("%d", -mod)))
	}
	return style.Render("+" + lipgloss.NewStyle().Render(fmt.Sprintf("%d", mod)))
}

// FormatHP formats HP with appropriate styling
func FormatHP(current, max int) string {
	percentage := float64(current) / float64(max)
	style := HPStyle

	if percentage <= 0.25 {
		style = HPCriticalStyle
	} else if percentage <= 0.5 {
		style = HPLowStyle
	}

	return style.Render(fmt.Sprintf("%d/%d", current, max))
}
