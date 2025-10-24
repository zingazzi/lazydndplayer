// internal/ui/panels/traits.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

type TraitsPanel struct {
	character      *models.Character
	viewport       viewport.Model
	ready          bool
	selectedIndex  int
	selectedType   string // "language", "feat", "resistance", "trait", or "mastery"
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

	// Weapon Mastery Section (MOVED TO TOP for better visibility)
	rightCol = append(rightCol, titleStyle.Render("⚔️  WEAPON MASTERY"))
	rightCol = append(rightCol, "")

	if p.hasWeaponMasteryFeature(p.character) {
		// Get mastery count
		masteryCount := p.getWeaponMasteryCount(p.character)
		masteryInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("  Can master %d weapons", masteryCount))
		rightCol = append(rightCol, masteryInfo)
		rightCol = append(rightCol, "")

		// Debug: Check MasteredWeapons array
		debug.Log("TraitsPanel.View: MasteredWeapons count=%d, weapons=%v", len(p.character.MasteredWeapons), p.character.MasteredWeapons)

		// Show currently mastered weapons with descriptions
		if len(p.character.MasteredWeapons) > 0 {
			masteredStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)
			rightCol = append(rightCol, masteredStyle.Render(fmt.Sprintf("  Mastered: (%d weapons)", len(p.character.MasteredWeapons))))
			rightCol = append(rightCol, "")

			masteryNameStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true)

			for i, weapon := range p.character.MasteredWeapons {
				debug.Log("TraitsPanel.View: Processing weapon[%d]='%s'", i, weapon)
				itemDef := models.GetItemDefinitionByName(weapon)
				debug.Log("TraitsPanel.View: itemDef=%v, mastery='%s'", itemDef != nil, func() string {
					if itemDef != nil {
						return itemDef.Mastery
					}
					return "N/A"
				}())

				isSelected := p.selectedType == "mastery" && i == p.selectedIndex

				if itemDef != nil && itemDef.Mastery != "" {
					// Show weapon name with mastery type (selectable)
					weaponLine := fmt.Sprintf("    ✓ %s", weapon)
					if isSelected {
						rightCol = append(rightCol, selectedStyle.Render("  → "+weaponLine[4:]))
					} else {
						rightCol = append(rightCol, normalStyle.Render(weaponLine))
					}
					rightCol = append(rightCol, masteryNameStyle.Render(fmt.Sprintf("      %s", itemDef.Mastery)))

					// Show mastery description
					masteryDesc := models.GetMasteryDescription(itemDef.Mastery)
					debug.Log("TraitsPanel.View: masteryDesc length=%d", len(masteryDesc))
					if masteryDesc != "" {
						// Wrap description to fit in the column (adjusted for indent)
						wrapped := wrapText(masteryDesc, width/2-10)
						for _, line := range wrapped {
							rightCol = append(rightCol, descStyle.Render(fmt.Sprintf("        %s", line)))
						}
					}
					rightCol = append(rightCol, "")
				} else {
					// Weapon without mastery property (or not found)
					debug.Log("TraitsPanel.View: Weapon '%s' not found or no mastery", weapon)
					weaponLine := fmt.Sprintf("    ✓ %s", weapon)
					if isSelected {
						rightCol = append(rightCol, selectedStyle.Render("  → "+weaponLine[4:]))
					} else {
						rightCol = append(rightCol, normalStyle.Render(weaponLine))
					}
					if itemDef == nil {
						rightCol = append(rightCol, lipgloss.NewStyle().
							Foreground(lipgloss.Color("196")).
							Render("      (weapon data not found)"))
					} else {
						rightCol = append(rightCol, lipgloss.NewStyle().
							Foreground(lipgloss.Color("240")).
							Render("      (no mastery property)"))
					}
					rightCol = append(rightCol, "")
				}
			}
		} else {
			// No mastered weapons yet
			rightCol = append(rightCol, lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true).
				Render("  No weapons mastered yet"))
			rightCol = append(rightCol, "")
		}

		rightCol = append(rightCol, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("  Press 'm' to manage"))
	} else {
		debug.Log("TraitsPanel.View: No weapon mastery feature found")
		rightCol = append(rightCol, emptyStyle.Render("  No weapon mastery feature"))
	}

	rightCol = append(rightCol, "")
	rightCol = append(rightCol, "")

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

	// Proficiencies Section
	rightCol = append(rightCol, titleStyle.Render("⚔  PROFICIENCIES"))
	rightCol = append(rightCol, "")

	// Armor proficiencies
	if len(p.character.ArmorProficiencies) > 0 {
		rightCol = append(rightCol, normalStyle.Render("  Armor: "+strings.Join(p.character.ArmorProficiencies, ", ")))
	} else {
		rightCol = append(rightCol, emptyStyle.Render("  No armor proficiencies"))
	}

	// Weapon proficiencies
	if len(p.character.WeaponProficiencies) > 0 {
		rightCol = append(rightCol, normalStyle.Render("  Weapons: "+strings.Join(p.character.WeaponProficiencies, ", ")))
	} else {
		rightCol = append(rightCol, emptyStyle.Render("  No weapon proficiencies"))
	}

	// Tool proficiencies
	if len(p.character.ToolProficiencies) > 0 {
		rightCol = append(rightCol, normalStyle.Render("  Tools: "+strings.Join(p.character.ToolProficiencies, ", ")))
	}

	// Saving throw proficiencies
	if len(p.character.SavingThrowProficiencies) > 0 {
		rightCol = append(rightCol, normalStyle.Render("  Saves: "+strings.Join(p.character.SavingThrowProficiencies, ", ")))
	} else {
		rightCol = append(rightCol, emptyStyle.Render("  No saving throw proficiencies"))
	}

	rightCol = append(rightCol, "")
	rightCol = append(rightCol, "")

	// Fighting Style Section (if character has one)
	if p.character.FightingStyle != "" {
		rightCol = append(rightCol, titleStyle.Render("⚔️  FIGHTING STYLE"))
		rightCol = append(rightCol, "")

		// Get fighting style details for description
		fightingStyleData := models.GetFightingStyleByName(p.character.FightingStyle)
		if fightingStyleData != nil {
			rightCol = append(rightCol, normalStyle.Render("  "+p.character.FightingStyle))
			// Add description with wrapping
			descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			wrapped := wrapText(fightingStyleData.Description, width/2-6)
			for _, line := range wrapped {
				rightCol = append(rightCol, descStyle.Render("    "+line))
			}
		} else {
			rightCol = append(rightCol, normalStyle.Render("  "+p.character.FightingStyle))
		}

		rightCol = append(rightCol, "")
		// Add hint to change fighting style
		rightCol = append(rightCol, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("  Press 's' to change"))
		rightCol = append(rightCol, "")
	}

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
		} else if len(p.character.MasteredWeapons) > 0 {
			// Move to weapon mastery section
			p.selectedType = "mastery"
			p.selectedIndex = 0
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
	} else if p.selectedType == "mastery" {
		if p.selectedIndex < len(p.character.MasteredWeapons)-1 {
			p.selectedIndex++
			p.viewport.LineDown(5) // More space because of descriptions
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

// ScrollDown scrolls the viewport down by one line
func (p *TraitsPanel) ScrollDown() {
	p.viewport.LineDown(1)
}

// ScrollUp scrolls the viewport up by one line
func (p *TraitsPanel) ScrollUp() {
	p.viewport.LineUp(1)
}

// PageDown scrolls down by half a page
func (p *TraitsPanel) PageDown() {
	p.viewport.HalfViewDown()
}

// PageUp scrolls up by half a page
func (p *TraitsPanel) PageUp() {
	p.viewport.HalfViewUp()
}

// GotoTop scrolls to the top of the content
func (p *TraitsPanel) GotoTop() {
	p.viewport.GotoTop()
}

// GotoBottom scrolls to the bottom of the content
func (p *TraitsPanel) GotoBottom() {
	p.viewport.GotoBottom()
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
		} else if len(p.character.MasteredWeapons) > 0 {
			// Move to weapon mastery section
			p.selectedType = "mastery"
			p.selectedIndex = len(p.character.MasteredWeapons) - 1
			p.viewport.LineUp(3)
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
		} else if len(p.character.MasteredWeapons) > 0 {
			// Move to weapon mastery section
			p.selectedType = "mastery"
			p.selectedIndex = len(p.character.MasteredWeapons) - 1
			p.viewport.LineUp(3)
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "resistance" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(3)
		} else if len(p.character.MasteredWeapons) > 0 {
			// Move to weapon mastery section
			p.selectedType = "mastery"
			p.selectedIndex = len(p.character.MasteredWeapons) - 1
			p.viewport.LineUp(3)
		} else if len(p.character.Languages) > 0 {
			// Move to languages section
			p.selectedType = "language"
			p.selectedIndex = len(p.character.Languages) - 1
		}
	} else if p.selectedType == "mastery" {
		if p.selectedIndex > 0 {
			p.selectedIndex--
			p.viewport.LineUp(5) // More space because of descriptions
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

// GetSelectedMastery returns the currently selected weapon mastery (if any)
func (p *TraitsPanel) GetSelectedMastery() string {
	if p.selectedType == "mastery" && p.selectedIndex >= 0 && p.selectedIndex < len(p.character.MasteredWeapons) {
		return p.character.MasteredWeapons[p.selectedIndex]
	}
	return ""
}

// IsOnMastery returns true if currently on a weapon mastery
func (p *TraitsPanel) IsOnMastery() bool {
	return p.selectedType == "mastery" && len(p.character.MasteredWeapons) > 0
}

// hasWeaponMasteryFeature checks if the character has a weapon mastery feature
func (p *TraitsPanel) hasWeaponMasteryFeature(char *models.Character) bool {
	for _, feature := range char.Features.Features {
		if feature.Name == "Weapon Mastery" {
			return true
		}
	}
	return false
}

// getWeaponMasteryCount returns the number of weapons the character can master
func (p *TraitsPanel) getWeaponMasteryCount(char *models.Character) int {
	// Read weapons_mastered from feature mechanics (generic approach)
	for _, feature := range char.Features.Features {
		if feature.Name == "Weapon Mastery" && feature.Mechanics != nil {
			if weaponsMastered, ok := feature.Mechanics["weapons_mastered"].(float64); ok {
				return int(weaponsMastered)
			}
		}
	}
	return 0
}

// getAvailableWeaponsToMaster returns weapons the character has proficiency with but hasn't mastered
func (p *TraitsPanel) getAvailableWeaponsToMaster(char *models.Character) []string {
	availableWeapons := []string{}

	// Get all weapon items
	allItems := models.GetAllItemDefinitions()

	// Build a map of already mastered weapons for quick lookup
	masteredMap := make(map[string]bool)
	for _, weapon := range char.MasteredWeapons {
		masteredMap[weapon] = true
	}

	// Check each weapon
	for _, itemDef := range allItems {
		// Only consider weapons
		if itemDef.Category != "weapon" {
			continue
		}

		// Skip if already mastered
		if masteredMap[itemDef.Name] {
			continue
		}

		// Check if character has proficiency with this weapon
		if p.hasProficiencyForWeapon(char, itemDef.Subcategory) {
			availableWeapons = append(availableWeapons, itemDef.Name)
		}
	}

	return availableWeapons
}

// hasProficiencyForWeapon checks if the character has proficiency with a weapon type
func (p *TraitsPanel) hasProficiencyForWeapon(char *models.Character, subcategory string) bool {
	subcategory = strings.ToLower(subcategory)

	// Check weapon proficiencies
	for _, prof := range char.WeaponProficiencies {
		profLower := strings.ToLower(prof)

		// Check for "Simple" or "Martial" proficiency
		if profLower == "simple" && strings.Contains(subcategory, "simple") {
			return true
		}
		if profLower == "martial" && strings.Contains(subcategory, "martial") {
			return true
		}
	}

	return false
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
