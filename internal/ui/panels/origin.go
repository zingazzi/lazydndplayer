// internal/ui/panels/origin.go
package panels

import (
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

// View renders the origin panel with simplified character information
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

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content []string

	content = append(content, titleStyle.Render("CHARACTER ORIGIN"))
	content = append(content, "")

	// Origin Section
	if p.character.Origin != "" {
		origin := models.GetOriginByName(p.character.Origin)
		if origin != nil {
			content = append(content, sectionTitleStyle.Render("Origin: ")+valueStyle.Render(origin.Name))
			content = append(content, hintStyle.Render("  [Press Enter for details]"))
			content = append(content, "")
			
			// Short description (first 100 chars)
			shortDesc := p.getOriginShortDescription(origin.Description)
			content = append(content, labelStyle.Render("  "+shortDesc))
		} else {
			content = append(content, sectionTitleStyle.Render("Origin: ")+emptyStyle.Render("Unknown"))
		}
	} else {
		content = append(content, sectionTitleStyle.Render("Origin: ")+emptyStyle.Render("Not set"))
		content = append(content, hintStyle.Render("  [Press 'o' to select]"))
	}
	content = append(content, "")
	
	// Alignment Section
	alignment := p.character.Alignment
	if alignment == "" {
		alignment = "Not set"
	}
	content = append(content, sectionTitleStyle.Render("Alignment: ")+valueStyle.Render(alignment))
	content = append(content, hintStyle.Render("  [Press 'a' to change]"))
	content = append(content, "")
	
	// Appearance Section
	content = append(content, sectionTitleStyle.Render("Appearance:"))
	
	heightValue := p.character.Height
	if heightValue == "" {
		heightValue = "Not set"
	}
	content = append(content, labelStyle.Render("  Height: ")+valueStyle.Render(heightValue)+hintStyle.Render("  [Press 'h' to edit]"))
	
	weightValue := p.character.Weight
	if weightValue == "" {
		weightValue = "Not set"
	}
	content = append(content, labelStyle.Render("  Weight: ")+valueStyle.Render(weightValue)+hintStyle.Render("  [Press 'w' to edit]"))
	content = append(content, "")
	
	// Personality Section
	personality := p.character.Personality
	if personality == "" {
		personality = "Not set"
	}
	displayPersonality := p.truncateText(personality, 60)
	content = append(content, sectionTitleStyle.Render("Personality: ")+valueStyle.Render(displayPersonality))
	content = append(content, hintStyle.Render("  [Press 'p' to change]"))
	content = append(content, "")
	
	// Ideal Section
	ideal := p.character.Ideal
	if ideal == "" {
		ideal = "Not set"
	}
	displayIdeal := p.truncateText(ideal, 60)
	content = append(content, sectionTitleStyle.Render("Ideal: ")+valueStyle.Render(displayIdeal))
	content = append(content, hintStyle.Render("  [Press 'i' to change]"))
	content = append(content, "")
	
	// Bond Section
	bond := p.character.Bond
	if bond == "" {
		bond = "Not set"
	}
	displayBond := p.truncateText(bond, 60)
	content = append(content, sectionTitleStyle.Render("Bond: ")+valueStyle.Render(displayBond))
	content = append(content, hintStyle.Render("  [Press 'b' to change]"))
	content = append(content, "")
	
	// Flaw Section
	flaw := p.character.Flaw
	if flaw == "" {
		flaw = "Not set"
	}
	displayFlaw := p.truncateText(flaw, 60)
	content = append(content, sectionTitleStyle.Render("Flaw: ")+valueStyle.Render(displayFlaw))
	content = append(content, hintStyle.Render("  [Press 'f' to change]"))
	content = append(content, "")
	
	// Backstory Section
	content = append(content, sectionTitleStyle.Render("Backstory:"))
	backstory := p.character.Backstory
	if backstory == "" {
		content = append(content, labelStyle.Render("  ")+emptyStyle.Render("No backstory written"))
	} else {
		backstoryPreview := p.getBackstoryPreview(backstory)
		content = append(content, labelStyle.Render("  "+backstoryPreview))
	}
	content = append(content, hintStyle.Render("  [Press 's' to edit]"))

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

// getOriginShortDescription truncates origin description to ~100 chars
func (p *OriginPanel) getOriginShortDescription(description string) string {
	if len(description) <= 100 {
		return description
	}
	
	// Find last space before 100 chars
	truncated := description[:100]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = description[:lastSpace]
	}
	
	return truncated + "..."
}

// getBackstoryPreview truncates backstory to ~100 chars
func (p *OriginPanel) getBackstoryPreview(backstory string) string {
	if len(backstory) <= 100 {
		return backstory
	}
	
	// Find last space before 100 chars
	truncated := backstory[:100]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = backstory[:lastSpace]
	}
	
	return truncated + "..."
}

// truncateText truncates text to specified length
func (p *OriginPanel) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	
	truncated := text[:maxLen-3]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = text[:lastSpace]
	}
	
	return truncated + "..."
}
