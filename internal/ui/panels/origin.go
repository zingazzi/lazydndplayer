// internal/ui/panels/origin.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// OriginPanel displays character origin information
type OriginPanel struct {
	character *models.Character
	viewport  viewport.Model
	ready     bool
}

// NewOriginPanel creates a new origin panel
func NewOriginPanel(char *models.Character) *OriginPanel {
	return &OriginPanel{
		character: char,
		ready:     false,
	}
}

// View renders the origin panel with two-column layout
func (p *OriginPanel) View(width, height int) string {
	// Use all available height for the viewport
	viewportHeight := height

	if !p.ready {
		p.viewport = viewport.New(width, viewportHeight)
		p.viewport.Style = lipgloss.NewStyle()
		p.ready = true
	}

	if p.viewport.Width != width || p.viewport.Height != viewportHeight {
		p.viewport.Width = width
		p.viewport.Height = viewportHeight
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Align(lipgloss.Center)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	abilityStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	featStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Bold(true)

	skillStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Build content
	var content []string

	content = append(content, titleStyle.Render("CHARACTER ORIGIN"))
	content = append(content, "")

	// Check if origin is set
	if p.character.Origin == "" {
		content = append(content, emptyStyle.Render("No origin selected"))
		content = append(content, "")
		content = append(content, helpStyle.Render("Press 'o' to select an origin"))

		contentStr := strings.Join(content, "\n")
		p.viewport.SetContent(contentStr)
		return p.viewport.View()
	}

	// Get origin details
	origin := models.GetOriginByName(p.character.Origin)
	if origin == nil {
		content = append(content, emptyStyle.Render("Origin not found"))
		contentStr := strings.Join(content, "\n")
		p.viewport.SetContent(contentStr)
		return p.viewport.View()
	}

	// Calculate column widths (60% left, 40% right)
	leftWidth := int(float64(width) * 0.6)
	rightWidth := width - leftWidth - 3 // -3 for spacing

	// LEFT COLUMN - Origin Information
	var leftContent []string

	// Origin name
	leftContent = append(leftContent, sectionTitleStyle.Render(origin.Name))
	leftContent = append(leftContent, "")

	// Description
	wrapped := wrapText(origin.Description, leftWidth-4)
	for _, line := range wrapped {
		leftContent = append(leftContent, descStyle.Render(line))
	}
	leftContent = append(leftContent, "")

	// Ability Increases
	if origin.AbilityIncreases != nil {
		leftContent = append(leftContent, abilityStyle.Render("ABILITY INCREASE:"))
		if len(origin.AbilityIncreases.Choices) > 0 {
			// Show what was chosen (tracked in BenefitTracker)
			benefits := p.character.BenefitTracker.GetBenefitsBySource("origin", origin.Name)
			for _, benefit := range benefits {
				if benefit.Type == models.BenefitAbilityScore {
					leftContent = append(leftContent, valueStyle.Render(
						fmt.Sprintf("  +%d %s", benefit.Value, benefit.Target)))
				}
			}
		} else if origin.AbilityIncreases.Ability != "" {
			leftContent = append(leftContent, valueStyle.Render(
				fmt.Sprintf("  +%d %s",
					origin.AbilityIncreases.Amount,
					origin.AbilityIncreases.Ability)))
		}
		leftContent = append(leftContent, "")
	}

	// Granted Feat
	if origin.Feat != "" {
		leftContent = append(leftContent, featStyle.Render("GRANTED FEAT:"))
		leftContent = append(leftContent, valueStyle.Render("  "+origin.Feat))

		// Show if feat was applied
		hasFeat := false
		for _, feat := range p.character.Feats {
			if feat == origin.Feat {
				hasFeat = true
				break
			}
		}
		if hasFeat {
			leftContent = append(leftContent, labelStyle.Render("  (Applied ✓)"))
		}
		leftContent = append(leftContent, "")
	}

	// Equipment
	if len(origin.Equipment) > 0 {
		leftContent = append(leftContent, sectionTitleStyle.Render("STARTING EQUIPMENT:"))
		for _, item := range origin.Equipment {
			wrappedItem := wrapText(item, leftWidth-6)
			for _, line := range wrappedItem {
				leftContent = append(leftContent, labelStyle.Render("  • "+line))
			}
		}
		leftContent = append(leftContent, "")
	}

	// Help text
	leftContent = append(leftContent, "")
	leftContent = append(leftContent, helpStyle.Render("Press 'o' to change origin"))

	// RIGHT COLUMN - Tool Proficiencies
	var rightContent []string

	rightContent = append(rightContent, sectionTitleStyle.Render("TOOL PROFICIENCIES"))
	rightContent = append(rightContent, "")

	// Get all tool proficiencies (from all sources, not just origin)
	allTools := p.character.ToolProficiencies
	
	if len(allTools) > 0 {
		// Show all tools character has
		for _, tool := range allTools {
			wrappedTool := wrapText(tool, rightWidth-4)
			for i, line := range wrappedTool {
				if i == 0 {
					rightContent = append(rightContent, valueStyle.Render("  • "+line))
				} else {
					rightContent = append(rightContent, labelStyle.Render("    "+line))
				}
			}
		}
		rightContent = append(rightContent, "")
	} else {
		rightContent = append(rightContent, labelStyle.Render("  No tool proficiencies"))
		rightContent = append(rightContent, "")
	}

	rightContent = append(rightContent, "")
	rightContent = append(rightContent, helpStyle.Render("Press 't' to add tool"))
	rightContent = append(rightContent, helpStyle.Render("Press 'T' to remove tool"))

	// Languages (if any from origin)
	originLanguages := []string{}
	benefits := p.character.BenefitTracker.GetBenefitsBySource("origin", origin.Name)
	for _, benefit := range benefits {
		if benefit.Type == models.BenefitLanguage {
			originLanguages = append(originLanguages, benefit.Target)
		}
	}

	if len(originLanguages) > 0 {
		rightContent = append(rightContent, skillStyle.Render("LANGUAGES:"))
		for _, lang := range originLanguages {
			rightContent = append(rightContent, valueStyle.Render("  • "+lang+" ✓"))
		}
		rightContent = append(rightContent, "")
	}

	// Create bordered columns
	leftColumnStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Padding(0, 1)

	rightColumnStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderLeft(true).
		BorderForeground(lipgloss.Color("240"))

	// Render columns
	leftRendered := leftColumnStyle.Render(strings.Join(leftContent, "\n"))
	rightRendered := rightColumnStyle.Render(strings.Join(rightContent, "\n"))

	// Combine columns side by side
	combinedContent := lipgloss.JoinHorizontal(lipgloss.Top, leftRendered, rightRendered)

	content = append(content, combinedContent)

	contentStr := strings.Join(content, "\n")
	p.viewport.SetContent(contentStr)

	return p.viewport.View()
}

// Update handles viewport updates
func (p *OriginPanel) Update(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

// ScrollDown scrolls the viewport down
func (p *OriginPanel) ScrollDown() {
	p.viewport.LineDown(3)
}

// ScrollUp scrolls the viewport up
func (p *OriginPanel) ScrollUp() {
	p.viewport.LineUp(3)
}

// PageDown scrolls down by half a page
func (p *OriginPanel) PageDown() {
	p.viewport.HalfViewDown()
}

// PageUp scrolls up by half a page
func (p *OriginPanel) PageUp() {
	p.viewport.HalfViewUp()
}
