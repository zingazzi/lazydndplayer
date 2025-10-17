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
	character     *models.Character
	selectedIndex int
	viewport      viewport.Model
	ready         bool
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	var lines []string
	lines = append(lines, titleStyle.Render("ACTIONS"))
	lines = append(lines, "")

	// Group actions by type
	actionsByType := make(map[models.ActionType][]models.Action)
	for _, action := range char.Actions.Actions {
		actionsByType[action.Type] = append(actionsByType[action.Type], action)
	}

	// Display by type
	types := []models.ActionType{
		models.StandardAction,
		models.BonusAction,
		models.Reaction,
		models.FreeAction,
	}

	idx := 0
	for _, actionType := range types {
		actions := actionsByType[actionType]
		if len(actions) == 0 {
			continue
		}

		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true).
			Render(string(actionType)+"s"))

		for _, action := range actions {
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
				lines = append(lines, selectedStyle.Render(line))
			} else {
				lines = append(lines, normalStyle.Render(line))
			}
			idx++
		}
		lines = append(lines, "")
	}

	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/↓ Navigate • Enter Activate"))

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
	if p.selectedIndex < len(p.character.Actions.Actions)-1 {
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
