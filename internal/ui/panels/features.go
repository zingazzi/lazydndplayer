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
		// Separate consumable and passive features
		consumableFeatures := []models.Feature{}
		passiveFeatures := []models.Feature{}

		for _, feature := range p.character.Features.Features {
			if feature.MaxUses > 0 {
				consumableFeatures = append(consumableFeatures, feature)
			} else {
				passiveFeatures = append(passiveFeatures, feature)
			}
		}

		currentIndex := 0

		// Render CONSUMABLE features first (with rest type grouping)
		if len(consumableFeatures) > 0 {
			content = append(content, titleStyle.Render("=== CONSUMABLE FEATURES ==="))
			content = append(content, "")

			// Show Focus Points for Monk
			if p.character.IsMonk() {
				monk := p.character.GetMonkMechanics()
				currentFP, maxFP := monk.GetFocusPoints()
				if maxFP > 0 {
					fpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
					fpValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
					fpLine := fpStyle.Render("✧ Focus Points: ") + fpValueStyle.Render(fmt.Sprintf("%d/%d", currentFP, maxFP)) +
						normalStyle.Render(" (Short Rest) ") +
						lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[+/- to adjust]")
					content = append(content, fpLine)
					content = append(content, "")
				}
			}

			// Group consumable by rest type
			shortRest := []models.Feature{}
			longRest := []models.Feature{}
			daily := []models.Feature{}

			for _, f := range consumableFeatures {
				switch f.RestType {
				case models.ShortRest:
					shortRest = append(shortRest, f)
				case models.LongRest:
					longRest = append(longRest, f)
				case models.Daily:
					daily = append(daily, f)
				default:
					longRest = append(longRest, f) // Default to long rest
				}
			}

			if len(shortRest) > 0 {
				content = append(content, titleStyle.Render("⚡ Short Rest"))
				content = append(content, "")
				content = p.renderFeatureGroup(content, shortRest, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
				content = append(content, "")
			}

			if len(longRest) > 0 {
				content = append(content, titleStyle.Render("🌙 Long Rest"))
				content = append(content, "")
				content = p.renderFeatureGroup(content, longRest, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
				content = append(content, "")
			}

			if len(daily) > 0 {
				content = append(content, titleStyle.Render("📅 Daily"))
				content = append(content, "")
				content = p.renderFeatureGroup(content, daily, &currentIndex, normalStyle, selectedStyle, usedStyle, descStyle, width)
				content = append(content, "")
			}
		}

		// Render PASSIVE features second
		if len(passiveFeatures) > 0 {
			content = append(content, titleStyle.Render("=== PASSIVE FEATURES ==="))
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

		// Add Focus Point cost for Monk abilities
		fpCostInfo := ""
		if feature.Name == "Flurry of Blows" || feature.Name == "Patient Defense" || feature.Name == "Step of the Wind" {
			fpCostInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(" [Costs 1 FP]")
		}

		var featureLine string
		if isSelected {
			featureLine = selectedStyle.Render(fmt.Sprintf("  → %s%s", feature.Name, usageInfo)) + fpCostInfo
		} else {
			if feature.MaxUses > 0 && feature.CurrentUses == 0 {
				featureLine = usedStyle.Render(fmt.Sprintf("    %s%s [USED]", feature.Name, usageInfo)) + fpCostInfo
			} else {
				featureLine = normalStyle.Render(fmt.Sprintf("    %s%s", feature.Name, usageInfo)) + fpCostInfo
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
	// Count all features (consumable and passive)
	totalCount := len(p.character.Features.Features)

	if p.selectedIndex < totalCount-1 {
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
	actualIndex := p.getActualFeatureIndex()
	if actualIndex >= 0 {
		p.character.Features.UseFeature(actualIndex)
	}
}

func (p *FeaturesPanel) RestoreFeature() {
	actualIndex := p.getActualFeatureIndex()
	if actualIndex >= 0 {
		p.character.Features.RestoreFeature(actualIndex)
	}
}

func (p *FeaturesPanel) RemoveFeature() {
	actualIndex := p.getActualFeatureIndex()
	if actualIndex >= 0 {
		p.character.Features.RemoveFeature(actualIndex)

		// Count remaining features
		totalCount := len(p.character.Features.Features)

		// Adjust selection if needed
		if p.selectedIndex >= totalCount && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	}
}

// getActualFeatureIndex maps the selected index to the actual index in the full features array
func (p *FeaturesPanel) getActualFeatureIndex() int {
	// Now we show all features, so selected index maps directly
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Features.Features) {
		return p.selectedIndex
	}
	return -1
}

// GetSelectedIndex returns the currently selected feature index
func (p *FeaturesPanel) GetSelectedIndex() int {
	return p.selectedIndex
}

// GetSelectedFeature returns the currently selected feature (if any)
// Returns all features (consumable and passive)
func (p *FeaturesPanel) GetSelectedFeature() *models.Feature {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Features.Features) {
		return &p.character.Features.Features[p.selectedIndex]
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
