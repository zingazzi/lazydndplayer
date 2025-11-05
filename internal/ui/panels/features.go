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

// SelectableItem represents any selectable item in the features panel
type SelectableItem struct {
	ItemType    string // "consumable", "passive", "rage_effect"
	Consumable  *ConsumableItem
	Feature     *models.Feature
	Name        string
	Description string
}

type FeaturesPanel struct {
	character       *models.Character
	viewport        viewport.Model
	ready           bool
	selectedIndex   int
	consumables     []ConsumableItem     // Unified list of all consumables
	selectableItems []SelectableItem     // All selectable items (consumables + passives + rage effects)
	rageEffects     []SelectableItem     // Rage-specific effects
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

	// Add Warrior Dice for Path of the Zealot
	if p.character.IsZealot() && p.character.WarriorDice.Max > 0 {
		longRest = append(longRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "warrior_dice",
			Name:         fmt.Sprintf("Warrior Dice 1%s", p.character.WarriorDice.Size),
			Current:      p.character.WarriorDice.Current,
			Max:          p.character.WarriorDice.Max,
			RestType:     models.LongRest,
			Description:  "Divine energy dice used by Path of the Zealot Barbarians. As a bonus action, expend one die to heal yourself, rolling the die and regaining HP equal to the result. Regain all on long rest.",
		})
	}

	// Add Psionic Dice for Soulknife
	if p.character.IsSoulknife() && p.character.SoulknifePsiDice.Max > 0 {
		shortRest = append(shortRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "soulknife_psi_dice",
			Name:         fmt.Sprintf("Psionic Energy 1%s", p.character.SoulknifePsiDice.Size),
			Current:      p.character.SoulknifePsiDice.Current,
			Max:          p.character.SoulknifePsiDice.Max,
			RestType:     models.ShortRest,
			Description:  "Psionic energy dice used by Soulknife Rogues. Fuel various psionic powers like Psi-Bolstered Knack and Psychic Whispers. Regain 1 on short rest, all on long rest.",
		})
	}

	// Add Channel Divinity for Cleric/Paladin
	if p.character.ChannelDivinity.Max > 0 {
		// Determine rest type from Channel Divinity feature
		restType := models.ShortRest // Default
		desc := "Channel Divinity allows you to channel divine energy to fuel magical effects. "

		// Check if character has Channel Divinity feature to get details
		for i := range p.character.Features.Features {
			if p.character.Features.Features[i].Name == "Channel Divinity" {
				feature := &p.character.Features.Features[i]
				desc = feature.Description
				if feature.Mechanics != nil {
					if regainOne, ok := feature.Mechanics["regain_one_on_short_rest"].(bool); ok && regainOne {
						desc += " Regain 1 use on short rest, all uses on long rest."
					}
				}
				break
			}
		}

		shortRest = append(shortRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "channel_divinity",
			Name:         "Channel Divinity",
			Current:      p.character.ChannelDivinity.Current,
			Max:          p.character.ChannelDivinity.Max,
			RestType:     restType,
			Description:  desc,
		})
	}

	// Add Lay on Hands pool for Paladin
	if p.character.LayOnHands.Max > 0 {
		longRest = append(longRest, ConsumableItem{
			ItemType:     "resource",
			ResourceType: "lay_on_hands",
			Name:         "Lay on Hands (HP Pool)",
			Current:      p.character.LayOnHands.Current,
			Max:          p.character.LayOnHands.Max,
			RestType:     models.LongRest,
			Description:  "A pool of healing power that replenishes when you take a long rest. As a bonus action, you can touch a creature and restore a number of hit points from the pool, up to the creature's maximum hit points. Pool size: 5 × paladin level.",
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
			case "Channel Divinity":
				// Managed via ChannelDivinity resource
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

// buildRageEffects builds the list of rage effects for Barbarians
func (p *FeaturesPanel) buildRageEffects() []SelectableItem {
	var effects []SelectableItem

	// Only show rage effects if character is Barbarian and has Rage feature
	if !p.character.HasClass("Barbarian") || !p.character.HasFeature("Rage") {
		return effects
	}

	// Get Barbarian level for damage calculation
	barbarianLevel := p.character.GetBarbarianLevel()
	if barbarianLevel == 0 {
		barbarianLevel = 1
	}

	// Basic Rage Effects (always present)

	// Rage Damage
	rageDamage := models.GetRageDamageBonus(barbarianLevel)
	effects = append(effects, SelectableItem{
		ItemType:    "rage_effect",
		Name:        fmt.Sprintf("Rage Damage: +%d", rageDamage),
		Description: fmt.Sprintf("While raging, you deal +%d extra damage on melee weapon attacks using Strength.", rageDamage),
	})

	// Damage Resistance
	effects = append(effects, SelectableItem{
		ItemType:    "rage_effect",
		Name:        "Damage Resistance: B/P/S",
		Description: "While raging, you have resistance to bludgeoning, piercing, and slashing damage.",
	})

	// Strength Advantage
	effects = append(effects, SelectableItem{
		ItemType:    "rage_effect",
		Name:        "Strength Advantage",
		Description: "While raging, you have advantage on Strength checks and Strength saving throws.",
	})

	// No Concentration
	effects = append(effects, SelectableItem{
		ItemType:    "rage_effect",
		Name:        "No Concentration",
		Description: "While raging, you cannot maintain concentration on spells.",
	})

	// Subclass-Specific Rage Effects

	// Path of the Berserker - Frenzy
	if p.character.IsBerserker() {
		frenzyDice := rageDamage
		effects = append(effects, SelectableItem{
			ItemType:    "rage_effect",
			Name:        fmt.Sprintf("Frenzy: %dd6", frenzyDice),
			Description: fmt.Sprintf("When you use Reckless Attack while raging, deal an extra %dd6 damage to the first target you hit.", frenzyDice),
		})
	}

	// Path of the Wild Heart - Rage of the Wilds
	if p.character.IsWildHeart() {
		effects = append(effects, SelectableItem{
			ItemType:    "rage_effect",
			Name:        "Rage of the Wilds",
			Description: "Choose one animal form when entering rage:\n• Bear: Resistance to all damage except Force, Necrotic, Psychic, Radiant\n• Eagle: Dash or Disengage as bonus action (can use both)\n• Wolf: Allies within 5ft have advantage on attacks",
		})
	}

	// Path of the Zealot - Divine Fury
	if p.character.IsZealot() {
		halfLevel := barbarianLevel / 2
		effects = append(effects, SelectableItem{
			ItemType:    "rage_effect",
			Name:        fmt.Sprintf("Divine Fury: 1d6+%d", halfLevel),
			Description: fmt.Sprintf("First creature you hit each turn takes extra 1d6+%d damage. Choose Necrotic or Radiant damage type each time.", halfLevel),
		})
	}

	// Path of the World Tree - Vitality of the Tree
	if p.character.IsWorldTree() {
		vitalitySurge := barbarianLevel
		lifeGivingDice := rageDamage
		effects = append(effects, SelectableItem{
			ItemType:    "rage_effect",
			Name:        "Vitality of the Tree",
			Description: fmt.Sprintf("When entering rage:\n• Vitality Surge: Gain %d temp HP\n• Life-Giving Force: One ally within 10ft gains %dd6 temp HP", vitalitySurge, lifeGivingDice),
		})
	}

	return effects
}

// buildSelectableItems builds the complete list of selectable items
func (p *FeaturesPanel) buildSelectableItems() {
	p.selectableItems = []SelectableItem{}

	// Add Rage Effects (if applicable)
	p.rageEffects = p.buildRageEffects()
	for i := range p.rageEffects {
		p.selectableItems = append(p.selectableItems, p.rageEffects[i])
	}

	// Add consumables
	for i := range p.consumables {
		item := &p.consumables[i]
		p.selectableItems = append(p.selectableItems, SelectableItem{
			ItemType:    "consumable",
			Consumable:  item,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	// Add passive features
	for i := range p.character.Features.Features {
		feature := &p.character.Features.Features[i]
		if feature.MaxUses == 0 {
			p.selectableItems = append(p.selectableItems, SelectableItem{
				ItemType:    "passive",
				Feature:     feature,
				Name:        feature.Name,
				Description: feature.Description,
			})
		}
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

	// Build consumables list and selectable items
	p.consumables = p.buildConsumablesList()
	p.buildSelectableItems()

	if len(p.selectableItems) == 0 {
		content = append(content, emptyStyle.Render("No features yet"))
		content = append(content, "")
		content = append(content, normalStyle.Render("Features are limited-use abilities that recharge on rest."))
	} else {
		currentIndex := 0

		// Render RAGE EFFECTS (if applicable)
		if len(p.rageEffects) > 0 {
			content = append(content, titleStyle.Render("=== RAGE EFFECTS ==="))
			content = append(content, "")

			for i := range p.rageEffects {
				isSelected := currentIndex == p.selectedIndex
				effect := &p.rageEffects[i]

				var itemLine string
				if isSelected {
					itemLine = selectedStyle.Render(fmt.Sprintf("  → %s", effect.Name))
				} else {
					itemLine = normalStyle.Render(fmt.Sprintf("    %s", effect.Name))
				}
				content = append(content, itemLine)

				// Show short description when selected
				if isSelected && len(effect.Description) < 100 {
					wrapped := wrapFeatureText(effect.Description, width-8)
					for _, line := range wrapped {
						content = append(content, descStyle.Render("      "+line))
					}
				}

				currentIndex++
				content = append(content, "")
			}
		}

		// Render CONSUMABLE features
		if len(p.consumables) > 0 {
			content = append(content, titleStyle.Render("=== CONSUMABLE FEATURES ==="))
			content = append(content, "")

			// Group by rest type for rendering
			currentRestType := models.RestType("")

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

				// Special handling for Portent - show rolled d20 values
				if item.Name == "Portent" && item.Feature != nil && item.Feature.Mechanics != nil {
					if portentRolls, ok := item.Feature.Mechanics["portent_rolls"].([]interface{}); ok && len(portentRolls) > 0 {
						rollsStr := ""
						for i, roll := range portentRolls {
							if i > 0 {
								rollsStr += ", "
							}
							// Convert interface{} to int
							if rollInt, ok := roll.(int); ok {
								rollsStr += fmt.Sprintf("%d", rollInt)
							} else if rollFloat, ok := roll.(float64); ok {
								rollsStr += fmt.Sprintf("%d", int(rollFloat))
							}
						}
						if rollsStr != "" {
							usageInfo = fmt.Sprintf(" (rolls: %s) [%d/%d uses]", rollsStr, item.Current, item.Max)
						}
					}
				}

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
		passiveFeatures := []models.Feature{}
		for _, feature := range p.character.Features.Features {
			if feature.MaxUses == 0 {
				passiveFeatures = append(passiveFeatures, feature)
			}
		}

		if len(passiveFeatures) > 0 {
			content = append(content, titleStyle.Render("=== PASSIVE FEATURES ==="))
			content = append(content, "")

			for i := range passiveFeatures {
				feature := &passiveFeatures[i]
				isSelected := currentIndex == p.selectedIndex
				currentIndex++

				var featureLine string
				if isSelected {
					featureLine = selectedStyle.Render(fmt.Sprintf("  → %s", feature.Name))
				} else {
					featureLine = normalStyle.Render(fmt.Sprintf("    %s", feature.Name))
				}
				content = append(content, featureLine)

				// Show short description when selected
				if isSelected && len(feature.Description) < 100 {
					wrapped := wrapFeatureText(feature.Description, width-8)
					for _, line := range wrapped {
						content = append(content, descStyle.Render("      "+line))
					}
				}

				if feature.Source != "" && !isSelected {
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
	// Navigate through all selectable items
	if p.selectedIndex < len(p.selectableItems)-1 {
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

// GetSelectedItem returns the currently selected item (any type)
func (p *FeaturesPanel) GetSelectedItem() *SelectableItem {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.selectableItems) {
		return &p.selectableItems[p.selectedIndex]
	}
	return nil
}

// GetSelectedConsumable returns the currently selected consumable item
func (p *FeaturesPanel) GetSelectedConsumable() *ConsumableItem {
	item := p.GetSelectedItem()
	if item != nil && item.ItemType == "consumable" && item.Consumable != nil {
		return item.Consumable
	}
	return nil
}

// UseFeature decrements the uses of a feature
func (p *FeaturesPanel) UseFeature() {
	item := p.GetSelectedConsumable()
	if item != nil && item.ItemType == "feature" && item.Feature != nil {
		// Find the feature index in the character's features list
		featureIndex := -1
		for i, f := range p.character.Features.Features {
			if &f == item.Feature {
				featureIndex = i
				break
			}
		}

		if featureIndex >= 0 {
			// Use the FeatureList's UseFeature method to trigger any special logic
			p.character.Features.UseFeature(featureIndex)
		} else {
			// Fallback: direct modification if not found
			if item.Feature.CurrentUses > 0 {
				item.Feature.CurrentUses--
			}
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
