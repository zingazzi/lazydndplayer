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

// MonkReaction represents a Monk-specific reaction
type MonkReaction struct {
	Name        string
	Description string
	FPCost      int
	Trigger     string
}

// FighterBonusAction represents a Fighter-specific bonus action
type FighterBonusAction struct {
	Name         string
	Description  string
	UsesFeature  string // Name of the feature that tracks uses (e.g., "Second Wind")
	CurrentUses  int
	MaxUses      int
}

// FighterReaction represents a Fighter-specific reaction
type FighterReaction struct {
	Name        string
	Description string
	Trigger     string
	Cost        string // e.g., "1 Psi Die", "1 Superiority Die"
}

// BarbarianBonusAction represents a Barbarian-specific bonus action
type BarbarianBonusAction struct {
	Name         string
	Description  string
	ResourceName string // Name of the resource (e.g., "Warrior Dice")
	CurrentUses  int
	MaxUses      int
	DiceSize     string // e.g., "d12"
}

// RogueBonusAction represents a Rogue-specific bonus action
type RogueBonusAction struct {
	Name        string
	Description string
	ActionType  string // "Cunning Action" or "Steady Aim"
}

// ClericReaction represents a Cleric-specific reaction
type ClericReaction struct {
	Name         string
	Description  string
	Trigger      string
	Cost         string // e.g., "Channel Divinity", "1 use"
	CurrentUses  int
	MaxUses      int
	FeatureName  string // Name of the feature this reaction comes from
}

// ActionsPanel displays character actions
type ActionsPanel struct {
	character          *models.Character
	selectedIndex      int
	viewport           viewport.Model
	ready              bool
	attacks            []models.Attack       // Cached attacks list
	actionSpells       []models.Spell        // Spells that are actions
	bonusSpells        []models.Spell        // Spells that are bonus actions
	reactionSpells     []models.Spell        // Spells that are reactions
	monkBonusActions      []MonkBonusAction      // Monk bonus actions
	monkReactions         []MonkReaction         // Monk reactions
	fighterBonusActions   []FighterBonusAction   // Fighter bonus actions
	fighterReactions      []FighterReaction      // Fighter reactions
	barbarianBonusActions []BarbarianBonusAction // Barbarian bonus actions
	rogueBonusActions     []RogueBonusAction     // Rogue bonus actions
	clericReactions       []ClericReaction       // Cleric reactions
	totalItemCount        int                    // Total number of items (attacks + spells + actions)
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

	// Build Fighter bonus actions
	p.fighterBonusActions = []FighterBonusAction{}
	if char.HasClass("Fighter") {
		// Second Wind
		if char.HasFeature("Second Wind") {
			// Find the Second Wind feature to get current/max uses
			for _, feature := range char.Features.Features {
				if feature.Name == "Second Wind" {
					fighterLevel := char.GetClassLevel("Fighter")
					p.fighterBonusActions = append(p.fighterBonusActions, FighterBonusAction{
						Name:         "Second Wind",
						Description:  fmt.Sprintf("Regain 1d10+%d HP", fighterLevel),
						UsesFeature:  "Second Wind",
						CurrentUses:  feature.CurrentUses,
						MaxUses:      feature.MaxUses,
					})
					break
				}
			}
		}
	}

	// Build Barbarian bonus actions
	p.barbarianBonusActions = []BarbarianBonusAction{}
	if char.IsZealot() && char.HasFeature("Warrior of the Gods") {
		if char.WarriorDice.Max > 0 {
			p.barbarianBonusActions = append(p.barbarianBonusActions, BarbarianBonusAction{
				Name:         "Warrior of the Gods",
				Description:  fmt.Sprintf("Heal yourself by rolling 1%s", char.WarriorDice.Size),
				ResourceName: "Warrior Dice",
				CurrentUses:  char.WarriorDice.Current,
				MaxUses:      char.WarriorDice.Max,
				DiceSize:     char.WarriorDice.Size,
			})
		}
	}

	// Build Rogue bonus actions
	p.rogueBonusActions = []RogueBonusAction{}
	if char.IsRogue() {
		// Cunning Action (level 2+)
		if char.HasFeature("Cunning Action") {
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Dash (Cunning Action)",
				Description: "Dash as a bonus action",
				ActionType:  "Cunning Action",
			})
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Disengage (Cunning Action)",
				Description: "Disengage as a bonus action",
				ActionType:  "Cunning Action",
			})
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Hide (Cunning Action)",
				Description: "Hide as a bonus action",
				ActionType:  "Cunning Action",
			})
		}

		// Steady Aim (level 3+)
		if char.HasFeature("Steady Aim") {
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Steady Aim",
				Description: "Gain advantage on next attack (can't move this turn)",
				ActionType:  "Steady Aim",
			})
		}

		// Arcane Trickster: Mage Hand Legerdemain
		if char.IsArcaneTrickster() && char.HasFeature("Mage Hand Legerdemain") {
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Mage Hand",
				Description: "Control Mage Hand as bonus action (invisible)",
				ActionType:  "Mage Hand Legerdemain",
			})
		}

		// Soulknife: Psychic Blade Off-Hand
		if char.IsSoulknife() && char.HasFeature("Psychic Blades") {
			p.rogueBonusActions = append(p.rogueBonusActions, RogueBonusAction{
				Name:        "Psychic Blade (Off-Hand)",
				Description: "Attack with second psychic blade (1d4 psychic)",
				ActionType:  "Psychic Blades",
			})
		}
	}

	// Build Monk reactions
	p.monkReactions = []MonkReaction{}
	if char.IsMonk() {
		if char.HasFeature("Deflect Missiles") {
			p.monkReactions = append(p.monkReactions, MonkReaction{
				Name:        "Deflect Missiles",
				Description: "Reduce damage by 1d10 + Dex mod. If reduced to 0, spend 1 FP to deflect back",
				FPCost:      0, // Base cost is 0, 1 FP to deflect back
				Trigger:     "When hit by ranged weapon attack",
			})
		}
		if char.HasFeature("Slow Fall") {
			monkLevel := char.GetMonkLevel()
			damageReduction := monkLevel * 5
			p.monkReactions = append(p.monkReactions, MonkReaction{
				Name:        "Slow Fall",
				Description: fmt.Sprintf("Reduce fall damage by %d", damageReduction),
				FPCost:      0,
				Trigger:     "When you fall",
			})
		}
	}

	// Build Fighter reactions
	p.fighterReactions = []FighterReaction{}
	if char.IsPsiWarrior() {
		if char.HasFeature("Psionic Power") {
			intMod := char.AbilityScores.GetModifier(models.Intelligence)
			p.fighterReactions = append(p.fighterReactions, FighterReaction{
				Name:        "Protective Field",
				Description: fmt.Sprintf("Reduce damage by %s + %d to you or an ally within 30 ft", char.PsiDice.Size, intMod),
				Trigger:     "When you or ally within 30 ft takes damage",
				Cost:        "1 Psi Die",
			})
		}
	}

	// Build Cleric reactions
	p.clericReactions = []ClericReaction{}
	if char.GetClassLevel("Cleric") > 0 {
		// Warding Flare (Light Domain)
		for _, feature := range char.Features.Features {
			if feature.Name == "Warding Flare" {
				if feature.Mechanics != nil {
					if actionType, ok := feature.Mechanics["action_type"].(string); ok && actionType == "reaction" {
						p.clericReactions = append(p.clericReactions, ClericReaction{
							Name:        "Warding Flare",
							Description: "Impose disadvantage on attack roll against you",
							Trigger:     "When attacked within 30 feet",
							Cost:        fmt.Sprintf("%d/%d uses", feature.CurrentUses, feature.MaxUses),
							CurrentUses: feature.CurrentUses,
							MaxUses:     feature.MaxUses,
							FeatureName: "Warding Flare",
						})
					}
				}
			}
		}
		// Guided Strike (War Domain)
		for _, feature := range char.Features.Features {
			if feature.Name == "Guided Strike" {
				if feature.Mechanics != nil {
					if actionType, ok := feature.Mechanics["action_type"].(string); ok && actionType == "reaction" {
						costStr := "Channel Divinity"
						if char.ChannelDivinity.Current > 0 {
							costStr = fmt.Sprintf("Channel Divinity (%d/%d)", char.ChannelDivinity.Current, char.ChannelDivinity.Max)
						}
						p.clericReactions = append(p.clericReactions, ClericReaction{
							Name:        "Guided Strike",
							Description: "Give +10 bonus to attack roll",
							Trigger:     "When you or ally within 30 ft makes attack roll",
							Cost:        costStr,
							CurrentUses: char.ChannelDivinity.Current,
							MaxUses:     char.ChannelDivinity.Max,
							FeatureName: "Guided Strike",
						})
					}
				}
			}
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

	// Show Sneak Attack for Rogues as a proper action
	if char.IsRogue() && char.HasFeature("Sneak Attack") {
		rogue := char.GetRogueMechanics()
		sneakAttackDice := rogue.GetSneakAttackDamage()
		sneakAttackStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")) // Green for Rogue features

		line := fmt.Sprintf("%-20s %s (When you have advantage)", "Sneak Attack", sneakAttackDice)
		if idx == p.selectedIndex {
			lines = append(lines, selectedStyle.Render("▶ "+line))
		} else {
			lines = append(lines, sneakAttackStyle.Render("  "+line))
		}
		idx++
	}

	// Show action spells
	if len(p.actionSpells) > 0 {
	for _, spell := range p.actionSpells {
		slot := char.SpellBook.GetSlotByLevel(spell.Level)
		slotInfo := ""
		if slot != nil {
			slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
		}
		markers := ""
		if spell.Concentration {
			markers += " (C)"
		}
		if spell.Ritual {
			markers += " (R)"
		}
		line := fmt.Sprintf("%-20s Lv%d%s%s", spell.Name, spell.Level, slotInfo, markers)

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

	// Fighter bonus actions
	fighterActionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")) // Orange for Fighter actions

	if len(p.fighterBonusActions) > 0 {
		for _, action := range p.fighterBonusActions {
			// Show uses
			usesStr := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(fmt.Sprintf(" [%d/%d uses]", action.CurrentUses, action.MaxUses))

			// Gray out if no uses left
			var line string
			if action.CurrentUses == 0 {
				line = fmt.Sprintf("%-20s%s", action.Name, usesStr)
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+line+" (No uses left)"))
			} else {
				line = fmt.Sprintf("%-20s%s", action.Name, usesStr)
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+line))
				} else {
					lines = append(lines, fighterActionStyle.Render("  "+line))
				}
			}
			idx++
		}
	}

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

	// Barbarian bonus actions
	barbarianActionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")) // Red for Barbarian actions

	if len(p.barbarianBonusActions) > 0 {
		for _, action := range p.barbarianBonusActions {
			// Show dice uses
			usesStr := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf(" [%d/%d %s]", action.CurrentUses, action.MaxUses, action.DiceSize))

			// Gray out if no dice left
			var line string
			if action.CurrentUses == 0 {
				line = fmt.Sprintf("%-20s%s", action.Name, usesStr)
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+line+" (No dice left)"))
			} else {
				line = fmt.Sprintf("%-20s%s", action.Name, usesStr)
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+line))
				} else {
					lines = append(lines, barbarianActionStyle.Render("  "+line))
				}
			}
			idx++
		}
	}

	// Rogue bonus actions
	rogueActionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")) // Green for Rogue actions

	if len(p.rogueBonusActions) > 0 {
		for _, action := range p.rogueBonusActions {
			line := fmt.Sprintf("%-20s %s", action.Name, action.Description)
			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, rogueActionStyle.Render("  "+line))
			}
			idx++
		}
	}

	// Soulknife Psionic Dice resource
	if char.IsSoulknife() && char.SoulknifePsiDice.Max > 0 {
		psiDiceStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")) // Green for Rogue resources

		usesStr := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render(fmt.Sprintf(" [%d/%d %s]", char.SoulknifePsiDice.Current, char.SoulknifePsiDice.Max, char.SoulknifePsiDice.Size))

		var line string
		if char.SoulknifePsiDice.Current == 0 {
			line = fmt.Sprintf("%-20s%s", "Psionic Energy", usesStr)
			lines = append(lines, lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("  "+line+" (No dice left)"))
		} else {
			line = fmt.Sprintf("%-20s%s", "Psionic Energy", usesStr)
			if idx == p.selectedIndex {
				lines = append(lines, selectedStyle.Render("▶ "+line))
			} else {
				lines = append(lines, psiDiceStyle.Render("  "+line))
			}
		}
		idx++
	}

	// Bonus action spells
	if len(p.bonusSpells) > 0 {
	for _, spell := range p.bonusSpells {
		slot := char.SpellBook.GetSlotByLevel(spell.Level)
		slotInfo := ""
		if slot != nil {
			slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
		}
		markers := ""
		if spell.Concentration {
			markers += " (C)"
		}
		if spell.Ritual {
			markers += " (R)"
		}
		line := fmt.Sprintf("%-20s Lv%d%s%s", spell.Name, spell.Level, slotInfo, markers)

		if idx == p.selectedIndex {
			lines = append(lines, selectedStyle.Render("▶ "+line))
		} else {
				lines = append(lines, spellStyle.Render("  "+line))
			}
			idx++
		}
	}

	// Show "no bonus actions" only if there are no class actions and no spells
	if len(p.fighterBonusActions) == 0 && len(p.monkBonusActions) == 0 && len(p.barbarianBonusActions) == 0 && len(p.rogueBonusActions) == 0 && len(p.bonusSpells) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Render("  No bonus actions available"))
	}
	lines = append(lines, "")

	// === REACTIONS SECTION ===
	lines = append(lines, sectionStyle.Render("🛡️  Reactions"))
	lines = append(lines, "")

	// Monk reactions
	if len(p.monkReactions) > 0 {
		for _, reaction := range p.monkReactions {
			// Get current FP
			monk := char.GetMonkMechanics()
			currentFP, _ := monk.GetFocusPoints()

			// Format the reaction line
			var line string
			if reaction.FPCost > 0 {
				fpCostStr := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(fmt.Sprintf(" [%d FP]", reaction.FPCost))
				line = fmt.Sprintf("%-20s%s", reaction.Name, fpCostStr)
			} else {
				line = fmt.Sprintf("%-20s", reaction.Name)
			}

			// Add trigger info
			triggerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			fullLine := line + " " + triggerStyle.Render("("+reaction.Trigger+")")

			// Check if enough FP (for deflecting back)
			if reaction.FPCost > 0 && currentFP < reaction.FPCost {
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+fullLine+" (Not enough FP)"))
			} else {
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+fullLine))
				} else {
					lines = append(lines, monkActionStyle.Render("  "+fullLine))
				}
			}
			idx++
		}
	}

	// Fighter reactions
	if len(p.fighterReactions) > 0 {
		fighterReactionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
		for _, reaction := range p.fighterReactions {
			// Format cost
			costStr := lipgloss.NewStyle().Foreground(lipgloss.Color("135")).Render(fmt.Sprintf(" [%s]", reaction.Cost))
			line := fmt.Sprintf("%-20s%s", reaction.Name, costStr)

			// Add trigger info
			triggerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			fullLine := line + " " + triggerStyle.Render("("+reaction.Trigger+")")

			// Check if enough resource (Psi Die or Superiority Die)
			hasResource := false
			if char.IsPsiWarrior() && char.PsiDice.Current > 0 {
				hasResource = true
			} else if char.IsBattleMaster() && char.SuperiorityDice.Current > 0 {
				hasResource = true
			}

			if !hasResource {
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+fullLine+" (Not enough dice)"))
			} else {
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+fullLine))
				} else {
					lines = append(lines, fighterReactionStyle.Render("  "+fullLine))
				}
			}
			idx++
		}
	}

	// Cleric reactions
	if len(p.clericReactions) > 0 {
		clericReactionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		for _, reaction := range p.clericReactions {
			// Format cost
			costStr := lipgloss.NewStyle().Foreground(lipgloss.Color("135")).Render(fmt.Sprintf(" [%s]", reaction.Cost))
			line := fmt.Sprintf("%-20s%s", reaction.Name, costStr)

			// Add trigger info
			triggerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			fullLine := line + " " + triggerStyle.Render("("+reaction.Trigger+")")

			// Check if enough uses
			hasUses := false
			if reaction.FeatureName == "Warding Flare" {
				hasUses = reaction.CurrentUses > 0
			} else if reaction.FeatureName == "Guided Strike" {
				hasUses = char.ChannelDivinity.Current > 0
			}

			if !hasUses {
				lines = append(lines, lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render("  "+fullLine+" (No uses available)"))
			} else {
				if idx == p.selectedIndex {
					lines = append(lines, selectedStyle.Render("▶ "+fullLine))
				} else {
					lines = append(lines, clericReactionStyle.Render("  "+fullLine))
				}
			}
			idx++
		}
	}

	// Reaction spells
	if len(p.reactionSpells) > 0 {
	for _, spell := range p.reactionSpells {
		slot := char.SpellBook.GetSlotByLevel(spell.Level)
		slotInfo := ""
		if slot != nil {
			slotInfo = fmt.Sprintf(" [%d/%d slots]", slot.Current, slot.Maximum)
		}
		markers := ""
		if spell.Concentration {
			markers += " (C)"
		}
		if spell.Ritual {
			markers += " (R)"
		}
		line := fmt.Sprintf("%-20s Lv%d%s%s", spell.Name, spell.Level, slotInfo, markers)

		if idx == p.selectedIndex {
			lines = append(lines, selectedStyle.Render("▶ "+line))
		} else {
				lines = append(lines, spellStyle.Render("  "+line))
			}
			idx++
		}
	}

	// Show "no reactions" only if there are none
	if len(p.monkReactions) == 0 && len(p.fighterReactions) == 0 && len(p.clericReactions) == 0 && len(p.reactionSpells) == 0 {
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

// IsSneakAttackSelected returns true if the selected item is Sneak Attack
func (p *ActionsPanel) IsSneakAttackSelected() bool {
	// Sneak Attack comes right after attacks
	sneakAttackIndex := len(p.attacks)
	return p.character.IsRogue() && p.character.HasFeature("Sneak Attack") && p.selectedIndex == sneakAttackIndex
}

// GetSneakAttackDice returns the Sneak Attack dice for rolling
func (p *ActionsPanel) GetSneakAttackDice() string {
	if p.character.IsRogue() && p.character.HasFeature("Sneak Attack") {
		rogue := p.character.GetRogueMechanics()
		return rogue.GetSneakAttackDamage()
	}
	return "0d6"
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
