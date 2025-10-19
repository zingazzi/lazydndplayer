// internal/ui/panels/features.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type FeaturesPanel struct {
	character      *models.Character
	viewport       viewport.Model
	ready          bool
	selectedIndex  int
}

func NewFeaturesPanel(char *models.Character) *FeaturesPanel {
	return &FeaturesPanel{
		character:     char,
		selectedIndex: 0,
	}
}

func (p *FeaturesPanel) View(width, height int) string {
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

	usedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red for depleted
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	var content []string

	if len(p.character.Features.Features) == 0 {
		content = append(content, emptyStyle.Render("No features yet"))
		content = append(content, "")
		content = append(content, normalStyle.Render("Features are limited-use abilities that recharge on rest."))
		content = append(content, normalStyle.Render("Press 'a' to add a new feature."))
	} else {
		// Group features by rest type
		shortRestFeatures := []models.Feature{}
		longRestFeatures := []models.Feature{}
		dailyFeatures := []models.Feature{}
		passiveFeatures := []models.Feature{}

		for _, feature := range p.character.Features.Features {
			// Skip passive features with 0 max uses (like Fighting Style, Weapon Mastery)
			if feature.MaxUses == 0 {
				continue
			}

			switch feature.RestType {
			case models.ShortRest:
				shortRestFeatures = append(shortRestFeatures, feature)
			case models.LongRest:
				longRestFeatures = append(longRestFeatures, feature)
			case models.Daily:
				dailyFeatures = append(dailyFeatures, feature)
			case models.None:
				passiveFeatures = append(passiveFeatures, feature)
			}
		}

		// Render features by category
		currentIndex := 0

		if len(shortRestFeatures) > 0 {
			content = append(content, titleStyle.Render("⚡ SHORT REST FEATURES"))
			content = append(content, "")
			content = p.renderFeatureGroup(content, shortRestFeatures, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
			content = append(content, "")
		}

		if len(longRestFeatures) > 0 {
			content = append(content, titleStyle.Render("🌙 LONG REST FEATURES"))
			content = append(content, "")
			content = p.renderFeatureGroup(content, longRestFeatures, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
			content = append(content, "")
		}

		if len(dailyFeatures) > 0 {
			content = append(content, titleStyle.Render("📅 DAILY FEATURES"))
			content = append(content, "")
			content = p.renderFeatureGroup(content, dailyFeatures, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
			content = append(content, "")
		}

		if len(passiveFeatures) > 0 {
			content = append(content, titleStyle.Render("✨ PASSIVE FEATURES"))
			content = append(content, "")
			content = p.renderFeatureGroup(content, passiveFeatures, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
			content = append(content, "")
		}
	}

	contentStr := strings.Join(content, "\n")
	p.viewport.SetContent(contentStr)

	// Render viewport
	viewportContent := p.viewport.View()

	// Overlay scroll indicator if content is scrollable
	if p.viewport.TotalLineCount() > p.viewport.Height {
		scrollPercentage := int(p.viewport.ScrollPercent() * 100)
		scrollInfo := fmt.Sprintf("[%d%%]", scrollPercentage)

		scrollStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Right)

		lines := strings.Split(viewportContent, "\n")
		if len(lines) > 0 {
			paddedLine := lipgloss.NewStyle().Width(width).Render(lines[len(lines)-1])
			lines[len(lines)-1] = lipgloss.JoinHorizontal(lipgloss.Top, paddedLine)
			lines = append(lines[:len(lines)-1],
				lipgloss.PlaceHorizontal(width, lipgloss.Right, scrollStyle.Render(scrollInfo)))
			viewportContent = strings.Join(lines, "\n")
		}
	}

	return viewportContent
}

func (p *FeaturesPanel) renderFeatureGroup(
	content []string,
	features []models.Feature,
	currentIndex *int,
	normalStyle, selectedStyle, usedStyle, descStyle lipgloss.Style,
	width int,
) []string {
	for _, feature := range features {
		isSelected := *currentIndex == p.selectedIndex
		*currentIndex++

		// Feature name and usage
		usageInfo := ""
		if feature.MaxUses > 0 {
			usageInfo = fmt.Sprintf(" (%d/%d)", feature.CurrentUses, feature.MaxUses)
		}

		var featureLine string
		if isSelected {
			featureLine = selectedStyle.Render(fmt.Sprintf("  → %s%s", feature.Name, usageInfo))
		} else {
			if feature.MaxUses > 0 && feature.CurrentUses == 0 {
				featureLine = usedStyle.Render(fmt.Sprintf("    %s%s [USED]", feature.Name, usageInfo))
			} else {
				featureLine = normalStyle.Render(fmt.Sprintf("    %s%s", feature.Name, usageInfo))
			}
		}
		content = append(content, featureLine)

		// Show description when selected or if it's short
		if isSelected || len(feature.Description) < 60 {
			wrapped := wrapFeatureText(feature.Description, width-8)
			for _, line := range wrapped {
				content = append(content, descStyle.Render("      "+line))
			}
		}

		// Show source
		if feature.Source != "" {
			content = append(content, descStyle.Render(fmt.Sprintf("      Source: %s", feature.Source)))
		}

		content = append(content, "")
	}
	return content
}

func (p *FeaturesPanel) Update(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

func (p *FeaturesPanel) Next() {
	// Count only usable features (MaxUses > 0)
	usableCount := 0
	for _, feature := range p.character.Features.Features {
		if feature.MaxUses > 0 {
			usableCount++
		}
	}

	if p.selectedIndex < usableCount-1 {
		p.selectedIndex++
		p.viewport.LineDown(3)
	}
}

func (p *FeaturesPanel) Prev() {
	if p.selectedIndex > 0 {
		p.selectedIndex--
		p.viewport.LineUp(3)
	}
}

func (p *FeaturesPanel) ScrollDown() {
	p.viewport.LineDown(3)
}

func (p *FeaturesPanel) ScrollUp() {
	p.viewport.LineUp(3)
}

func (p *FeaturesPanel) PageDown() {
	p.viewport.HalfViewDown()
}

func (p *FeaturesPanel) PageUp() {
	p.viewport.HalfViewUp()
}

func (p *FeaturesPanel) UseFeature() {
	if len(p.character.Features.Features) > 0 {
		p.character.Features.UseFeature(p.selectedIndex)
	}
}

func (p *FeaturesPanel) RestoreFeature() {
	if len(p.character.Features.Features) > 0 {
		p.character.Features.RestoreFeature(p.selectedIndex)
	}
}

func (p *FeaturesPanel) RemoveFeature() {
	if len(p.character.Features.Features) > 0 {
		p.character.Features.RemoveFeature(p.selectedIndex)
		if p.selectedIndex >= len(p.character.Features.Features) && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	}
}

// GetSelectedIndex returns the currently selected feature index
func (p *FeaturesPanel) GetSelectedIndex() int {
	return p.selectedIndex
}

// GetSelectedFeature returns the currently selected feature (if any)
// Only returns features with MaxUses > 0 (usable features)
func (p *FeaturesPanel) GetSelectedFeature() *models.Feature {
	// Build list of usable features (skip passive ones with MaxUses == 0)
	usableFeatures := []int{}
	for i, feature := range p.character.Features.Features {
		if feature.MaxUses > 0 {
			usableFeatures = append(usableFeatures, i)
		}
	}

	if p.selectedIndex >= 0 && p.selectedIndex < len(usableFeatures) {
		actualIndex := usableFeatures[p.selectedIndex]
		return &p.character.Features.Features[actualIndex]
	}
	return nil
}

// wrapFeatureText wraps text to a specified width
func wrapFeatureText(text string, width int) []string {
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
