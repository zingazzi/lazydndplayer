// internal/ui/components/fightingstyleselector.go
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// FightingStyleSelector handles fighting style selection
type FightingStyleSelector struct {
	styles        []models.FightingStyle
	selectedIndex int
	visible       bool
	className     string
}

// NewFightingStyleSelector creates a new fighting style selector
func NewFightingStyleSelector() *FightingStyleSelector {
	return &FightingStyleSelector{
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the fighting style selector
func (fss *FightingStyleSelector) Show(className string) {
	fss.className = className
	fss.styles = models.GetAllFightingStyles()
	fss.selectedIndex = 0
	fss.visible = true
}

// Hide closes the fighting style selector
func (fss *FightingStyleSelector) Hide() {
	fss.visible = false
}

// IsVisible returns whether the selector is visible
func (fss *FightingStyleSelector) IsVisible() bool {
	return fss.visible
}

// Next moves to the next fighting style
func (fss *FightingStyleSelector) Next() {
	if fss.selectedIndex < len(fss.styles)-1 {
		fss.selectedIndex++
	}
}

// Prev moves to the previous fighting style
func (fss *FightingStyleSelector) Prev() {
	if fss.selectedIndex > 0 {
		fss.selectedIndex--
	}
}

// GetSelectedStyle returns the currently selected fighting style
func (fss *FightingStyleSelector) GetSelectedStyle() string {
	if fss.selectedIndex >= 0 && fss.selectedIndex < len(fss.styles) {
		return fss.styles[fss.selectedIndex].Name
	}
	return ""
}

// Update handles keyboard input for the fighting style selector
func (fss *FightingStyleSelector) Update(msg tea.Msg) (FightingStyleSelector, tea.Cmd) {
	if !fss.visible {
		return *fss, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			fss.Prev()
		case "down", "j":
			fss.Next()
		case "enter", "esc":
			return *fss, nil
		}
	}

	return *fss, nil
}

// View renders the fighting style selector
func (fss *FightingStyleSelector) View(width, height int) string {
	if !fss.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246")).
		Width(70).
		MarginTop(1)

	reqStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Italic(true)

	benefitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content []string
	content = append(content, titleStyle.Render(fmt.Sprintf("SELECT FIGHTING STYLE - %s", strings.ToUpper(fss.className))))
	content = append(content, "")

	// List all fighting styles (left column)
	for i, style := range fss.styles {
		styleName := fmt.Sprintf(" %s", style.Name)
		if i == fss.selectedIndex {
			content = append(content, selectedStyle.Render(styleName))
		} else {
			content = append(content, normalStyle.Render(styleName))
		}
	}

	content = append(content, "")
	content = append(content, strings.Repeat("─", 74))

	// Show details of selected style (right side)
	if fss.selectedIndex >= 0 && fss.selectedIndex < len(fss.styles) {
		selectedStyle := fss.styles[fss.selectedIndex]

		content = append(content, "")
		content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render(selectedStyle.Name))
		content = append(content, "")
		content = append(content, descStyle.Render(selectedStyle.Description))

		// Show requirements if any
		if len(selectedStyle.Requirements) > 0 {
			content = append(content, "")
			content = append(content, reqStyle.Render("Requirements:"))
			for _, req := range selectedStyle.Requirements {
				content = append(content, reqStyle.Render("  • "+req))
			}
		}

		// Show benefits
		if len(selectedStyle.Benefits) > 0 {
			content = append(content, "")
			content = append(content, benefitStyle.Render("Benefits:"))

			// Display benefits in a readable format
			for key, value := range selectedStyle.Benefits {
				switch key {
				case "ac_bonus":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • AC: +%v", value)))
				case "attack_bonus_ranged":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • Ranged Attack Rolls: +%v", value)))
				case "damage_bonus":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • Damage: +%v", value)))
				case "damage_bonus_thrown":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • Thrown Weapon Damage: +%v", value)))
				case "unarmed_damage":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • Unarmed Damage: %v", value)))
				case "special":
					content = append(content, benefitStyle.Render(fmt.Sprintf("  • %v", value)))
				}
			}
		}
	}

	content = append(content, "")
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

	// Wrap in a box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1, 2).
		Width(78)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(strings.Join(content, "\n")),
	)
}
