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
	selectedType   string // "language", "feat", or "resistance"
}

func NewTraitsPanel(char *models.Character) *TraitsPanel {
	return &TraitsPanel{
		character:     char,
		selectedIndex: 0,
		selectedType:  "language",
	}
}

func (p *TraitsPanel) View(width, height int) string {
	if !p.ready {
		p.viewport = viewport.New(width, height)
		p.viewport.Style = lipgloss.NewStyle()
		p.ready = true
	}

	if p.viewport.Width != width || p.viewport.Height != height {
		p.viewport.Width = width
		p.viewport.Height = height
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

	// Combine columns
	colWidth := width / 2
	leftContent := strings.Join(leftCol, "\n")
	rightContent := strings.Join(rightCol, "\n")

	leftStyle := lipgloss.NewStyle().Width(colWidth)
	rightStyle := lipgloss.NewStyle().Width(colWidth)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftContent),
		rightStyle.Render(rightContent),
	)

	p.viewport.SetContent(content)

	// Add scroll indicators
	scrollInfo := ""
	if p.viewport.TotalLineCount() > p.viewport.Height {
		scrollPercentage := int(p.viewport.ScrollPercent() * 100)
		scrollInfo = fmt.Sprintf(" [%d%%]", scrollPercentage)
	}

	return p.viewport.View() + scrollInfo
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
			p.viewport.LineDown(1)
		} else if len(p.character.Resistances) > 0 {
			// Move to resistances section
			p.selectedType = "resistance"
			p.selectedIndex = 0
		} else if len(p.character.Feats) > 0 {
			// Move to feats section
			p.selectedType = "feat"
			p.selectedIndex = 0
		}
	} else if p.selectedType == "resistance" {
		if p.selectedIndex < len(p.character.Resistances)-1 {
			p.selectedIndex++
			p.viewport.LineDown(1)
		} else if len(p.character.Feats) > 0 {
			// Move to feats section
			p.selectedType = "feat"
			p.selectedIndex = 0
		}
	} else if p.selectedType == "feat" {
		if p.selectedIndex < len(p.character.Feats)-1 {
			p.selectedIndex++
			p.viewport.LineDown(1)
		}
	}
}

func (p *TraitsPanel) Prev() {
	if p.selectedType == "feat" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(1)
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
			p.viewport.LineUp(1)
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "language" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(1)
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
}
