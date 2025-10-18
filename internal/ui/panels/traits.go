// internal/ui/panels/traits.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type TraitsPanel struct {
	character      *models.Character
	viewport       viewport.Model
	ready          bool
	selectedIndex  int
	selectedType   string // "language", "feat", "resistance", or "trait"
}

func NewTraitsPanel(char *models.Character) *TraitsPanel {
	return &TraitsPanel{
		character:     char,
		selectedIndex: 0,
		selectedType:  "language",
	}
}

func (p *TraitsPanel) View(width, height int) string {
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

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build left column
	var leftCol []string

	// Darkvision Section
	leftCol = append(leftCol, titleStyle.Render("👁  DARKVISION"))
	leftCol = append(leftCol, "")
	if p.character.Darkvision > 0 {
		leftCol = append(leftCol, normalStyle.Render(fmt.Sprintf("    %d feet", p.character.Darkvision)))
	} else {
		leftCol = append(leftCol, emptyStyle.Render("  None"))
	}
	leftCol = append(leftCol, "")
	leftCol = append(leftCol, "")

	// Languages Section
	leftCol = append(leftCol, titleStyle.Render("🗣  LANGUAGES"))
	leftCol = append(leftCol, "")

	if len(p.character.Languages) == 0 {
		leftCol = append(leftCol, emptyStyle.Render("  No languages known"))
	} else {
		for i, lang := range p.character.Languages {
			if p.selectedType == "language" && i == p.selectedIndex {
				leftCol = append(leftCol, selectedStyle.Render(fmt.Sprintf("  → %s", lang)))
			} else {
				leftCol = append(leftCol, normalStyle.Render(fmt.Sprintf("    %s", lang)))
			}
		}
	}

	// Build right column
	var rightCol []string

	// Resistances Section
	rightCol = append(rightCol, titleStyle.Render("🛡  RESISTANCES"))
	rightCol = append(rightCol, "")

	if len(p.character.Resistances) == 0 {
		rightCol = append(rightCol, emptyStyle.Render("  No damage resistances"))
	} else {
		for i, resistance := range p.character.Resistances {
			if p.selectedType == "resistance" && i == p.selectedIndex {
				rightCol = append(rightCol, selectedStyle.Render(fmt.Sprintf("  → %s", resistance)))
			} else {
				rightCol = append(rightCol, normalStyle.Render(fmt.Sprintf("    %s", resistance)))
			}
		}
	}

	rightCol = append(rightCol, "")
	rightCol = append(rightCol, "")

	// Feats Section
	rightCol = append(rightCol, titleStyle.Render("⭐ FEATS"))
	rightCol = append(rightCol, "")

	if len(p.character.Feats) == 0 {
		rightCol = append(rightCol, emptyStyle.Render("  No feats acquired"))
	} else {
		for i, feat := range p.character.Feats {
			if p.selectedType == "feat" && i == p.selectedIndex {
				rightCol = append(rightCol, selectedStyle.Render(fmt.Sprintf("  → %s", feat)))
			} else {
				rightCol = append(rightCol, normalStyle.Render(fmt.Sprintf("    %s", feat)))
			}
		}
	}

	// Build traits section (full width, below columns)
	var traitsSection []string
	traitsSection = append(traitsSection, "")
	traitsSection = append(traitsSection, titleStyle.Render("✨ SPECIES TRAITS"))
	traitsSection = append(traitsSection, "")

	if len(p.character.SpeciesTraits) == 0 {
		traitsSection = append(traitsSection, emptyStyle.Render("  No species traits"))
	} else {
		for i, trait := range p.character.SpeciesTraits {
			if p.selectedType == "trait" && i == p.selectedIndex {
				traitsSection = append(traitsSection, selectedStyle.Render(fmt.Sprintf("  → %s", trait.Name)))
			} else {
				traitsSection = append(traitsSection, normalStyle.Render(fmt.Sprintf("    %s", trait.Name)))
			}
			// Add description with wrapping
			descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			wrapped := wrapText(trait.Description, width-6)
			for _, line := range wrapped {
				traitsSection = append(traitsSection, descStyle.Render("      "+line))
			}
			traitsSection = append(traitsSection, "")
		}
	}

	// Combine columns
	colWidth := width / 2
	leftContent := strings.Join(leftCol, "\n")
	rightContent := strings.Join(rightCol, "\n")

	leftStyle := lipgloss.NewStyle().Width(colWidth)
	rightStyle := lipgloss.NewStyle().Width(colWidth)

	columnsContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftContent),
		rightStyle.Render(rightContent),
	)

	// Add traits section below columns
	traitsContent := strings.Join(traitsSection, "\n")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		columnsContent,
		traitsContent,
	)

	p.viewport.SetContent(content)

	// Render viewport
	viewportContent := p.viewport.View()

	// Overlay scroll indicator if content is scrollable
	if p.viewport.TotalLineCount() > p.viewport.Height {
		scrollPercentage := int(p.viewport.ScrollPercent() * 100)
		scrollInfo := fmt.Sprintf("[%d%%]", scrollPercentage)

		// Position the scroll indicator at the bottom-right
		scrollStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Right)

		// Get the lines of the viewport
		lines := strings.Split(viewportContent, "\n")
		if len(lines) > 0 {
			// Replace the last line with scroll info on the right
			lastLine := lines[len(lines)-1]
			// Pad to full width and add scroll info
			paddedLine := lipgloss.NewStyle().Width(width).Render(lastLine)
			lines[len(lines)-1] = lipgloss.JoinHorizontal(
				lipgloss.Top,
				paddedLine,
			)
			// Overlay scroll info at bottom right
			lines = append(lines[:len(lines)-1],
				lipgloss.PlaceHorizontal(width, lipgloss.Right, scrollStyle.Render(scrollInfo)))
			viewportContent = strings.Join(lines, "\n")
		}
	}

	return viewportContent
}

func (p *TraitsPanel) Update(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

func (p *TraitsPanel) Next() {
	if p.selectedType == "language" {
		if p.selectedIndex < len(p.character.Languages)-1 {
			p.selectedIndex++
			p.viewport.LineDown(3)
		} else if len(p.character.Resistances) > 0 {
			// Move to resistances section
			p.selectedType = "resistance"
			p.selectedIndex = 0
		} else if len(p.character.Feats) > 0 {
			// Move to feats section
			p.selectedType = "feat"
			p.selectedIndex = 0
		} else if len(p.character.SpeciesTraits) > 0 {
			// Move to traits section
			p.selectedType = "trait"
			p.selectedIndex = 0
		}
	} else if p.selectedType == "resistance" {
		if p.selectedIndex < len(p.character.Resistances)-1 {
			p.selectedIndex++
			p.viewport.LineDown(3)
		} else if len(p.character.Feats) > 0 {
			// Move to feats section
			p.selectedType = "feat"
			p.selectedIndex = 0
		} else if len(p.character.SpeciesTraits) > 0 {
			// Move to traits section
			p.selectedType = "trait"
			p.selectedIndex = 0
		}
	} else if p.selectedType == "feat" {
		if p.selectedIndex < len(p.character.Feats)-1 {
			p.selectedIndex++
			p.viewport.LineDown(3)
		} else if len(p.character.SpeciesTraits) > 0 {
			// Move to traits section
			p.selectedType = "trait"
			p.selectedIndex = 0
		}
	} else if p.selectedType == "trait" {
		if p.selectedIndex < len(p.character.SpeciesTraits)-1 {
			p.selectedIndex++
			// Scroll more for traits since they have wrapped descriptions
			p.viewport.LineDown(5)
		}
	}
}

// ScrollDown scrolls the viewport down without changing selection
func (p *TraitsPanel) ScrollDown() {
	p.viewport.LineDown(3)
}

// ScrollUp scrolls the viewport up without changing selection
func (p *TraitsPanel) ScrollUp() {
	p.viewport.LineUp(3)
}

// PageDown scrolls down by half a page
func (p *TraitsPanel) PageDown() {
	p.viewport.HalfViewDown()
}

// PageUp scrolls up by half a page
func (p *TraitsPanel) PageUp() {
	p.viewport.HalfViewUp()
}

func (p *TraitsPanel) Prev() {
	if p.selectedType == "trait" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			// Scroll more for traits since they have wrapped descriptions
			p.viewport.LineUp(5)
		} else if len(p.character.Feats) > 0 {
			// Move to feats section
			p.selectedType = "feat"
			p.selectedIndex = len(p.character.Feats) - 1
		} else if len(p.character.Resistances) > 0 {
			// Move to resistances section
			p.selectedType = "resistance"
			p.selectedIndex = len(p.character.Resistances) - 1
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "feat" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(3)
		} else if len(p.character.Resistances) > 0 {
			// Move to resistances section
			p.selectedType = "resistance"
			p.selectedIndex = len(p.character.Resistances) - 1
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "resistance" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(3)
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "language" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(3)
		}
	}
}

func (p *TraitsPanel) AddLanguage(language string) {
	p.character.Languages = append(p.character.Languages, language)
}

func (p *TraitsPanel) AddFeat(feat string) {
	p.character.Feats = append(p.character.Feats, feat)
}

func (p *TraitsPanel) AddResistance(resistance string) {
	p.character.Resistances = append(p.character.Resistances, resistance)
}

func (p *TraitsPanel) RemoveSelected() {
	if p.selectedType == "language" && len(p.character.Languages) > 0 && p.selectedIndex < len(p.character.Languages) {
		p.character.Languages = append(
			p.character.Languages[:p.selectedIndex],
			p.character.Languages[p.selectedIndex+1:]...,
		)
		if p.selectedIndex >= len(p.character.Languages) && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	} else if p.selectedType == "resistance" && len(p.character.Resistances) > 0 && p.selectedIndex < len(p.character.Resistances) {
		p.character.Resistances = append(
			p.character.Resistances[:p.selectedIndex],
			p.character.Resistances[p.selectedIndex+1:]...,
		)
		if p.selectedIndex >= len(p.character.Resistances) && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	} else if p.selectedType == "feat" && len(p.character.Feats) > 0 && p.selectedIndex < len(p.character.Feats) {
		p.character.Feats = append(
			p.character.Feats[:p.selectedIndex],
			p.character.Feats[p.selectedIndex+1:]...,
		)
		if p.selectedIndex >= len(p.character.Feats) && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	}
	// Note: Species traits cannot be removed manually, they come from the species
}

// GetSelectedFeat returns the currently selected feat name (if any)
func (p *TraitsPanel) GetSelectedFeat() string {
	if p.selectedType == "feat" && p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Feats) {
		return p.character.Feats[p.selectedIndex]
	}
	return ""
}

// IsOnFeat returns true if currently on a feat
func (p *TraitsPanel) IsOnFeat() bool {
	return p.selectedType == "feat" && len(p.character.Feats) > 0
}

// wrapText wraps text to a specified width
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}
