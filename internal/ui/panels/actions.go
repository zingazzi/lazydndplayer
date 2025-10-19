// internal/ui/panels/actions.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ActionsPanel displays character actions
type ActionsPanel struct {
	character       *models.Character
	selectedIndex   int
	viewport        viewport.Model
	ready           bool
	attacks         []models.Attack // Cached attacks list
	totalItemCount  int             // Total number of items (attacks + actions)
}

// NewActionsPanel creates a new actions panel
func NewActionsPanel(char *models.Character) *ActionsPanel {
	return &ActionsPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// View renders the actions panel
func (p *ActionsPanel) View(width, height int) string {
	char := p.character

	// Initialize viewport if not ready
	if !p.ready {
		p.viewport = viewport.New(width, height)
		p.ready = true
	}
	p.viewport.Width = width
	p.viewport.Height = height

	// Generate attacks dynamically
	attackList := models.GenerateAttacks(char)
	p.attacks = attackList.Attacks

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	var lines []string
	lines = append(lines, titleStyle.Render("ACTIONS"))
	lines = append(lines, "")

	idx := 0

	// === ATTACKS SECTION ===
	if len(p.attacks) > 0 {
		lines = append(lines, sectionStyle.Render("⚔️  Attacks"))
		lines = append(lines, "")

		for _, attack := range p.attacks {
			line := fmt.Sprintf("%-20s %s", attack.Name, attack.GetAttackSummary())

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, normalStyle.Render("  "+line))
			}
			idx++
		}
		lines = append(lines, "")
	}

	// === OTHER ACTIONS SECTION ===
	// Filter out "Attack" from actions since we show specific attacks above
	otherActions := []models.Action{}
	for _, action := range char.Actions.Actions {
		if action.Name != "Attack" {
			otherActions = append(otherActions, action)
		}
	}

	if len(otherActions) > 0 {
		lines = append(lines, sectionStyle.Render("📋 Other Actions"))
		lines = append(lines, "")

		for _, action := range otherActions {
			usesStr := ""
			if action.UsesPerRest != -1 {
				usesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
				if action.UsesRemaining == 0 {
					usesStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
				}
				usesStr = usesStyle.Render(fmt.Sprintf(" [%d/%d]", action.UsesRemaining, action.UsesPerRest))
			}

			line := fmt.Sprintf("%-25s%s", action.Name, usesStr)

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, normalStyle.Render("  "+line))
			}
			idx++
		}
	}

	p.totalItemCount = idx

	lines = append(lines, "")

	// Show different help based on selection
	helpText := "↑/↓: Navigate"
	if p.selectedIndex < len(p.attacks) {
		// Attack selected
		helpText = "↑/↓: Navigate • 'r': Attack • 'a': Advantage • 'x': Disadvantage • 'd': Damage"
	} else {
		// Other action selected
		helpText = "↑/↓: Navigate • Enter: Activate"
	}

	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(helpText))

	content := strings.Join(lines, "\n")
	p.viewport.SetContent(content)

	// Add scroll indicators at bottom of viewport
	scrollInfo := ""
	if p.viewport.ScrollPercent() < 1.0 {
		scrollInfo = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("↓ %d%%", int(p.viewport.ScrollPercent()*100)))
	}

	return p.viewport.View() + scrollInfo
}

// Update handles updates for the actions panel
func (p *ActionsPanel) Update(char *models.Character) {
	p.character = char
}

// Next moves to next action
func (p *ActionsPanel) Next() {
	if p.selectedIndex < p.totalItemCount-1 {
		p.selectedIndex++
		p.viewport.LineDown(1)
	}
}

// Prev moves to previous action
func (p *ActionsPanel) Prev() {
	if p.selectedIndex > 0 {
		p.selectedIndex--
		p.viewport.LineUp(1)
	}
}

// GetSelectedAttack returns the currently selected attack (if in attack section)
func (p *ActionsPanel) GetSelectedAttack() *models.Attack {
	if p.selectedIndex < len(p.attacks) {
		return &p.attacks[p.selectedIndex]
	}
	return nil
}

// IsAttackSelected returns true if the selected item is an attack
func (p *ActionsPanel) IsAttackSelected() bool {
	return p.selectedIndex < len(p.attacks)
}
