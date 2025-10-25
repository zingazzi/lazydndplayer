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

// ConsumableItem represents a unified consumable resource or feature
type ConsumableItem struct {
	ItemType     string // "resource" or "feature"
	ResourceType string // "focus_points", "psi_dice", "superiority_dice", ""
	Feature      *models.Feature
	Name         string
	Current      int
	Max          int
	RestType     models.RestType
	Description  string
}

type FeaturesPanel struct {
	character     *models.Character
	viewport      viewport.Model
	ready         bool
	selectedIndex int
	consumables   []ConsumableItem // Unified list of all consumables
}

func NewFeaturesPanel(char *models.Character) *FeaturesPanel {
	return &FeaturesPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// buildConsumablesList builds a unified list of all consumable resources and features
func (p *FeaturesPanel) buildConsumablesList() []ConsumableItem {
	var items []ConsumableItem

	// Group by rest type
	shortRest := []ConsumableItem{}
	longRest := []ConsumableItem{}
	daily := []ConsumableItem{}

	// Add Focus Points for Monk
	if p.character.IsMonk() {
		monk := p.character.GetMonkMechanics()
		currentFP, maxFP := monk.GetFocusPoints()
		if maxFP > 0 {
			shortRest = append(shortRest, ConsumableItem{
				ItemType:     "resource",
				ResourceType: "focus_points",
				Name:         "Focus Points",
				Current:      currentFP,
				Max:          maxFP,
				RestType:     models.ShortRest,
				Description:  "Mystical energy used to power Monk techniques like Flurry of Blows, Patient Defense, and Step of the Wind.",
			})
		}
	}

	// Add Psi Dice for Psi Warrior
	if p.character.IsPsiWarrior() && p.character.PsiDice.Max > 0 {
		longRest = append(longRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "psi_dice",
			Name:         fmt.Sprintf("Psi Dice 1%s", p.character.PsiDice.Size),
			Current:      p.character.PsiDice.Current,
			Max:          p.character.PsiDice.Max,
			RestType:     models.LongRest,
			Description:  "Psionic energy dice used for Psi Warrior abilities like Protective Field, Psionic Strike, and Telekinetic Movement. Regain all on long rest, 1 on short rest.",
		})
	}

	// Add Superiority Dice for Battle Master
	if p.character.IsBattleMaster() && p.character.SuperiorityDice.Max > 0 {
		shortRest = append(shortRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "superiority_dice",
			Name:         fmt.Sprintf("Superiority Dice 1%s", p.character.SuperiorityDice.Size),
			Current:      p.character.SuperiorityDice.Current,
			Max:          p.character.SuperiorityDice.Max,
			RestType:     models.ShortRest,
			Description:  "Combat superiority dice used to fuel Battle Master maneuvers. Add the die result to attack rolls, damage, ability checks, or saving throws depending on the maneuver. Regain all on short or long rest.",
		})
	}

	// Add all consumable features (MaxUses > 0)
	// Skip features that are already represented as resources
	for i := range p.character.Features.Features {
		feature := &p.character.Features.Features[i]
		if feature.MaxUses > 0 {
			// Skip features that are managed as resources
			skipFeature := false
			switch feature.Name {
			case "Ki Points", "Focus Points", "Focus Point Improvement", "Ki Improvement":
				// These are managed via Monk mechanics
				skipFeature = true
			case "Psionic Power":
				// Managed via PsiDice
				skipFeature = true
			case "Combat Superiority":
				// Managed via SuperiorityDice
				skipFeature = true
			}

			if skipFeature {
				continue
			}

			item := ConsumableItem{
				ItemType:    "feature",
				Feature:     feature,
				Name:        feature.Name,
				Current:     feature.CurrentUses,
				Max:         feature.MaxUses,
				RestType:    feature.RestType,
				Description: feature.Description,
			}

			// Group by rest type
			switch feature.RestType {
			case models.ShortRest:
				shortRest = append(shortRest, item)
			case models.LongRest:
				longRest = append(longRest, item)
			case models.Daily:
				daily = append(daily, item)
			default:
				longRest = append(longRest, item) // Default to long rest
			}
		}
	}

	// Combine in order: Short Rest, Long Rest, Daily
	items = append(items, shortRest...)
	items = append(items, longRest...)
	items = append(items, daily...)

	return items
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

	// Build consumables list
	p.consumables = p.buildConsumablesList()

	// Separate passive features
	passiveFeatures := []models.Feature{}
	for _, feature := range p.character.Features.Features {
		if feature.MaxUses == 0 {
			passiveFeatures = append(passiveFeatures, feature)
		}
	}

	if len(p.consumables) == 0 && len(passiveFeatures) == 0 {
		content = append(content, emptyStyle.Render("No features yet"))
		content = append(content, "")
		content = append(content, normalStyle.Render("Features are limited-use abilities that recharge on rest."))
	} else {
		// Render CONSUMABLE features
		if len(p.consumables) > 0 {
			content = append(content, titleStyle.Render("=== CONSUMABLE FEATURES ==="))
			content = append(content, "")

			// Group by rest type for rendering
			currentRestType := models.RestType("")
			currentIndex := 0

			for _, item := range p.consumables {
				// Add rest type header if changed
				if item.RestType != currentRestType {
					currentRestType = item.RestType
					switch currentRestType {
					case models.ShortRest:
						content = append(content, titleStyle.Render("⚡ Short Rest"))
					case models.LongRest:
						content = append(content, titleStyle.Render("🌙 Long Rest"))
					case models.Daily:
						content = append(content, titleStyle.Render("📅 Daily"))
					}
					content = append(content, "")
				}

				// Render item
				isSelected := currentIndex == p.selectedIndex
				currentIndex++

				usageInfo := fmt.Sprintf(" (%d/%d)", item.Current, item.Max)

				var itemLine string
				if isSelected {
					itemLine = selectedStyle.Render(fmt.Sprintf("  → %s%s", item.Name, usageInfo))
				} else {
					if item.Current == 0 {
						itemLine = usedStyle.Render(fmt.Sprintf("    %s%s [USED]", item.Name, usageInfo))
					} else {
						itemLine = normalStyle.Render(fmt.Sprintf("    %s%s", item.Name, usageInfo))
					}
				}
				content = append(content, itemLine)

				// Show short description when selected
				if isSelected && len(item.Description) < 100 {
					wrapped := wrapFeatureText(item.Description, width-8)
					for _, line := range wrapped {
						content = append(content, descStyle.Render("      "+line))
					}
				}

				content = append(content, "")
			}
		}

		// Render PASSIVE features
		if len(passiveFeatures) > 0 {
			content = append(content, titleStyle.Render("=== PASSIVE FEATURES ==="))
			content = append(content, "")
			for _, feature := range passiveFeatures {
				featureLine := normalStyle.Render(fmt.Sprintf("    %s", feature.Name))
				content = append(content, featureLine)

				// Show short description
				if len(feature.Description) < 100 {
					wrapped := wrapFeatureText(feature.Description, width-8)
					for _, line := range wrapped {
						content = append(content, descStyle.Render("      "+line))
					}
				}

				if feature.Source != "" {
					content = append(content, descStyle.Render(fmt.Sprintf("      Source: %s", feature.Source)))
				}

				content = append(content, "")
			}
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

func (p *FeaturesPanel) Update(msg tea.Msg) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	_ = cmd
}

func (p *FeaturesPanel) Next() {
	// Navigate through consumables list
	if p.selectedIndex < len(p.consumables)-1 {
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

// GetSelectedConsumable returns the currently selected consumable item
func (p *FeaturesPanel) GetSelectedConsumable() *ConsumableItem {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.consumables) {
		return &p.consumables[p.selectedIndex]
	}
	return nil
}

// UseFeature decrements the uses of a feature
func (p *FeaturesPanel) UseFeature() {
	item := p.GetSelectedConsumable()
	if item != nil && item.ItemType == "feature" && item.Feature != nil {
		if item.Feature.CurrentUses > 0 {
			item.Feature.CurrentUses--
		}
	}
}

// RestoreFeature increments the uses of a feature
func (p *FeaturesPanel) RestoreFeature() {
	item := p.GetSelectedConsumable()
	if item != nil && item.ItemType == "feature" && item.Feature != nil {
		if item.Feature.CurrentUses < item.Feature.MaxUses {
			item.Feature.CurrentUses++
		}
	}
}

// RemoveFeature removes the selected feature
func (p *FeaturesPanel) RemoveFeature() {
	item := p.GetSelectedConsumable()
	if item != nil && item.ItemType == "feature" && item.Feature != nil {
		// Find feature in character's features list
		for i, f := range p.character.Features.Features {
			if &f == item.Feature {
				p.character.Features.RemoveFeature(i)
				// Adjust selection if needed
				if p.selectedIndex >= len(p.consumables) && p.selectedIndex > 0 {
					p.selectedIndex--
				}
				break
			}
		}
	}
}

// GetSelectedIndex returns the currently selected index
func (p *FeaturesPanel) GetSelectedIndex() int {
	return p.selectedIndex
}

// GetSelectedFeature returns the currently selected feature (for backward compatibility)
func (p *FeaturesPanel) GetSelectedFeature() *models.Feature {
	item := p.GetSelectedConsumable()
	if item != nil && item.ItemType == "feature" {
		return item.Feature
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
