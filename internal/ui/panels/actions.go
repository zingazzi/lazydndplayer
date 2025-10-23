// internal/ui/panels/actions.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// MonkBonusAction represents a Monk-specific bonus action
type MonkBonusAction struct {
	Name        string
	Description string
	FPCost      int
}

// ActionsPanel displays character actions
type ActionsPanel struct {
	character       *models.Character
	selectedIndex   int
	viewport        viewport.Model
	ready           bool
	attacks         []models.Attack  // Cached attacks list
	actionSpells    []models.Spell   // Spells that are actions
	bonusSpells     []models.Spell   // Spells that are bonus actions
	reactionSpells  []models.Spell   // Spells that are reactions
	monkBonusActions []MonkBonusAction // Monk bonus actions
	totalItemCount  int              // Total number of items (attacks + spells + actions)
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

	// Generate attacks dynamically
	attackList := models.GenerateAttacks(char)
	p.attacks = attackList.Attacks

	// Filter prepared spells by action type
	p.actionSpells = []models.Spell{}
	p.bonusSpells = []models.Spell{}
	p.reactionSpells = []models.Spell{}

	for _, spell := range char.SpellBook.Spells {
		if !spell.Prepared || spell.Level == 0 {
			continue // Skip unprepared spells and cantrips
		}

		actionType := strings.ToLower(spell.ActionType)
		if strings.Contains(actionType, "action") && !strings.Contains(actionType, "bonus") && !strings.Contains(actionType, "reaction") {
			p.actionSpells = append(p.actionSpells, spell)
		} else if strings.Contains(actionType, "bonus") {
			p.bonusSpells = append(p.bonusSpells, spell)
		} else if strings.Contains(actionType, "reaction") {
			p.reactionSpells = append(p.reactionSpells, spell)
		}
	}

	// Build Monk bonus actions
	p.monkBonusActions = []MonkBonusAction{}
	if char.IsMonk() {
		// Check for each Monk bonus action feature
		if char.HasFeature("Flurry of Blows") {
			p.monkBonusActions = append(p.monkBonusActions, MonkBonusAction{
				Name:        "Flurry of Blows",
				Description: "Make two unarmed strikes as bonus action",
				FPCost:      1,
			})
		}
		if char.HasFeature("Patient Defense") {
			p.monkBonusActions = append(p.monkBonusActions, MonkBonusAction{
				Name:        "Patient Defense",
				Description: "Disengage + Dodge as bonus action",
				FPCost:      1,
			})
		}
		if char.HasFeature("Step of the Wind") {
			p.monkBonusActions = append(p.monkBonusActions, MonkBonusAction{
				Name:        "Step of the Wind",
				Description: "Disengage or Dash, jump doubles",
				FPCost:      1,
			})
		}
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	spellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")) // Purple for spells

	var lines []string
	lines = append(lines, titleStyle.Render("ACTIONS"))
	lines = append(lines, "")

	idx := 0

	// === ATTACKS SECTION ===
	lines = append(lines, sectionStyle.Render("⚔️  Attacks"))
	lines = append(lines, "")

	if len(p.attacks) > 0 {
		for _, attack := range p.attacks {
			line := fmt.Sprintf("%-20s %s", attack.Name, attack.GetAttackSummary())

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, normalStyle.Render("  "+line))
			}
			idx++
		}
	} else {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Render("  No attacks available"))
	}

	// Show action spells
	if len(p.actionSpells) > 0 {
		for _, spell := range p.actionSpells {
			slot := char.SpellBook.GetSlotByLevel(spell.Level)
			slotInfo := ""
			if slot != nil {
				slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
			}
			line := fmt.Sprintf("%-20s Lv%d%s", spell.Name, spell.Level, slotInfo)

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, spellStyle.Render("  "+line))
			}
			idx++
		}
	}
	lines = append(lines, "")

	// === BONUS ACTIONS SECTION ===
	lines = append(lines, sectionStyle.Render("⚡ Bonus Actions"))
	lines = append(lines, "")

	// Monk bonus actions
	monkActionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")) // Purple for Monk actions, matching FP color

	if len(p.monkBonusActions) > 0 {
		for _, action := range p.monkBonusActions {
			// Get current FP
			monk := char.GetMonkMechanics()
			currentFP, _ := monk.GetFocusPoints()

			// Show FP cost
			fpCostStr := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(fmt.Sprintf(" [%d FP]", action.FPCost))

			// Gray out if not enough FP
			var line string
			if currentFP < action.FPCost {
				line = fmt.Sprintf("%-20s%s", action.Name, fpCostStr)
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+line+" (Not enough FP)"))
			} else {
				line = fmt.Sprintf("%-20s%s", action.Name, fpCostStr)
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+line))
				} else {
					lines = append(lines, monkActionStyle.Render("  "+line))
				}
			}
			idx++
		}
	}

	// Bonus action spells
	if len(p.bonusSpells) > 0 {
		for _, spell := range p.bonusSpells {
			slot := char.SpellBook.GetSlotByLevel(spell.Level)
			slotInfo := ""
			if slot != nil {
				slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
			}
			line := fmt.Sprintf("%-20s Lv%d%s", spell.Name, spell.Level, slotInfo)

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, spellStyle.Render("  "+line))
			}
			idx++
		}
	}

	// Show "no bonus actions" only if there are no monk actions and no spells
	if len(p.monkBonusActions) == 0 && len(p.bonusSpells) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Render("  No bonus actions available"))
	}
	lines = append(lines, "")

	// === REACTIONS SECTION ===
	lines = append(lines, sectionStyle.Render("🛡️  Reactions"))
	lines = append(lines, "")

	if len(p.reactionSpells) > 0 {
		for _, spell := range p.reactionSpells {
			slot := char.SpellBook.GetSlotByLevel(spell.Level)
			slotInfo := ""
			if slot != nil {
				slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
			}
			line := fmt.Sprintf("%-20s Lv%d%s", spell.Name, spell.Level, slotInfo)

			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, spellStyle.Render("  "+line))
			}
			idx++
		}
	} else {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Render("  No reactions available"))
	}

	p.totalItemCount = idx

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
	if p.selectedIndex < p.totalItemCount-1 {
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

// GetSelectedAttack returns the currently selected attack (if in attack section)
func (p *ActionsPanel) GetSelectedAttack() *models.Attack {
	if p.selectedIndex < len(p.attacks) {
		return &p.attacks[p.selectedIndex]
	}
	return nil
}

// IsAttackSelected returns true if the selected item is an attack
func (p *ActionsPanel) IsAttackSelected() bool {
	return p.selectedIndex < len(p.attacks)
}

// GetSelectedSpell returns the currently selected spell (if any)
func (p *ActionsPanel) GetSelectedSpell() *models.Spell {
	idx := p.selectedIndex - len(p.attacks)

	// Check if in action spells
	if idx >= 0 && idx < len(p.actionSpells) {
		return &p.actionSpells[idx]
	}
	idx -= len(p.actionSpells)

	// Check if in bonus action spells
	if idx >= 0 && idx < len(p.bonusSpells) {
		return &p.bonusSpells[idx]
	}
	idx -= len(p.bonusSpells)

	// Check if in reaction spells
	if idx >= 0 && idx < len(p.reactionSpells) {
		return &p.reactionSpells[idx]
	}

	return nil
}

// IsSpellSelected returns true if the selected item is a spell
func (p *ActionsPanel) IsSpellSelected() bool {
	return p.selectedIndex >= len(p.attacks) && p.GetSelectedSpell() != nil
}

// CastSelectedSpell casts the currently selected spell and consumes a slot
func (p *ActionsPanel) CastSelectedSpell() (string, bool) {
	spell := p.GetSelectedSpell()
	if spell == nil {
		return "", false
	}

	// Get the appropriate spell slot
	slot := p.character.SpellBook.GetSlotByLevel(spell.Level)
	if slot == nil || slot.Current <= 0 {
		return fmt.Sprintf("No spell slots available for level %d!", spell.Level), false
	}

	// Consume the spell slot
	slot.Current--

	// Format the result message
	msg := fmt.Sprintf("Cast %s! (%d/%d level %d slots remaining)",
		spell.Name, slot.Current, slot.Maximum, spell.Level)

	return msg, true
}
