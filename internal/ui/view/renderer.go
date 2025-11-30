// internal/ui/view/renderer.go
package view

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// ComponentRenderInfo contains information needed to render a component
type ComponentRenderInfo struct {
	Component interface{}
	Size      PopupSize
	Name      string
}

// PopupRenderer handles rendering of popups and overlays
type PopupRenderer struct {
	character *models.Character
	width     int
	height    int
	layout    *LayoutDimensions
}

// NewPopupRenderer creates a new popup renderer
func NewPopupRenderer(character *models.Character, width, height int, layout *LayoutDimensions) *PopupRenderer {
	return &PopupRenderer{
		character: character,
		width:     width,
		height:    height,
		layout:    layout,
	}
}

// RenderPopups renders all visible popups in priority order
// This handles the complex logic of determining which popup to show
// Components are checked in priority order and the first visible one is rendered
func (pr *PopupRenderer) RenderPopups(
	statGenerator *components.StatGenerator,
	abilityRoller *components.AbilityRoller,
	attackRoller *components.AttackRoller,
	spellSelector *components.SpellSelector,
	featSelector *components.FeatSelector,
	featDetailPopup *components.FeatDetailPopup,
	masteryDetailPopup *components.MasteryDetailPopup,
	maneuverDetailPopup *components.ManeuverDetailPopup,
	consumableDetailPopup *components.ConsumableDetailPopup,
	featureDetailPopup *components.FeatureDetailPopup,
	itemDetailPopup *components.ItemDetailPopup,
	spellDetailPopup *components.SpellDetailPopup,
	originSelector *components.OriginSelector,
	alignmentSelector *components.AlignmentSelector,
	traitSelector *components.TraitSelector,
	backstoryEditor *components.BackstoryEditor,
	originDetailPopup *components.OriginDetailPopup,
	inputPopup *components.InputPopup,
	abilityChoiceSelector *components.AbilityChoiceSelector,
	subtypeSelector *components.SubtypeSelector,
	skillSelector *components.SkillSelector,
	languageSelector *components.LanguageSelector,
	toolSelector *components.ToolSelector,
	weaponMasterySelector *components.WeaponMasterySelector,
	expertiseSelector *components.ExpertiseSelector,
	maneuverSelector *components.ManeuverSelector,
	levelUpSelector *components.LevelUpSelector,
	deLevelSelector *components.DeLevelSelector,
	itemSelector *components.ItemSelector,
	fightingStyleSelector *components.FightingStyleSelector,
	cantripSelector *components.CantripSelector,
	leveledSpellSelector *components.LeveledSpellSelector,
	schoolSpellSelector *components.SchoolSpellSelector,
	spellbookEditor *components.SpellbookEditor,
	spellPrepSelector *components.SpellPrepSelector,
	slotRestorer *components.SlotRestorer,
	classSkillSelector *components.ClassSkillSelector,
	subclassSelector *components.SubclassSelector,
	beastSelector *components.BeastSelector,
	companionSelector *components.CompanionSelector,
	wildShapeFormSelector *components.WildShapeFormSelector,
	optionSelector *components.OptionSelector,
	classSelector *components.ClassSelector,
	speciesSelector *components.SpeciesSelector,
	messagePopup *components.MessagePopup,
	restPopup *components.RestPopup,
	attackMenu *components.AttackMenu,
	divineOrderSelectorVisible bool,
	renderDivineOrderSelector func() string,
	primalOrderSelectorVisible bool,
	renderPrimalOrderSelector func() string,
) string {
	// Check popups in priority order (highest priority first)
	// Stat generator takes highest priority (Medium)
	if statGenerator.IsVisible() {
		return statGenerator.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Ability roller takes high priority (Small)
	if abilityRoller.IsVisible() {
		return abilityRoller.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight, pr.character)
	}

	// Attack roller takes high priority (Medium)
	if attackRoller.IsVisible() {
		return attackRoller.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Spell selector takes high priority (Large)
	if spellSelector.IsVisible() {
		return spellSelector.View(pr.layout.PopupLargeWidth, pr.layout.PopupLargeHeight)
	}

	// Feat selector takes second priority (Medium)
	if featSelector.IsVisible() {
		return featSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Feat detail popup (Medium)
	if featDetailPopup.IsVisible() {
		return featDetailPopup.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Mastery detail popup (Medium)
	if masteryDetailPopup.IsVisible() {
		return masteryDetailPopup.View(pr.width, pr.height)
	}

	// Maneuver detail popup (Medium)
	if maneuverDetailPopup.IsVisible() {
		return maneuverDetailPopup.View(pr.width, pr.height)
	}

	// Consumable detail popup (Medium)
	if consumableDetailPopup.IsVisible() {
		return consumableDetailPopup.View(pr.width, pr.height)
	}

	// Feature detail popup (Medium)
	if featureDetailPopup.IsVisible() {
		return featureDetailPopup.View(pr.width, pr.height)
	}

	// Item detail popup (Medium)
	if itemDetailPopup.IsVisible() {
		return itemDetailPopup.View(pr.width, pr.height)
	}

	// Spell detail popup (Medium)
	if spellDetailPopup.IsVisible() {
		return spellDetailPopup.View()
	}

	// Origin selector (Medium)
	if originSelector.IsVisible() {
		return originSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Alignment selector (takes full screen)
	if alignmentSelector.IsVisible() {
		return alignmentSelector.View(pr.width, pr.height)
	}

	// Trait selector (Medium)
	if traitSelector.IsVisible() {
		return traitSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Backstory editor (takes full screen)
	if backstoryEditor.IsVisible() {
		return backstoryEditor.View(pr.width, pr.height)
	}

	// Origin detail popup (Medium)
	if originDetailPopup.IsVisible() {
		return originDetailPopup.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Input popup (Small)
	if inputPopup.IsVisible() {
		return inputPopup.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Ability choice selector (for feat ability choices) (Small)
	if abilityChoiceSelector.IsVisible() {
		return abilityChoiceSelector.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Subtype selector takes third priority (Small)
	if subtypeSelector.IsVisible() {
		return subtypeSelector.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Skill selector takes fourth priority (Small)
	if skillSelector.IsVisible() {
		return skillSelector.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Language selector takes third priority (Small)
	if languageSelector.IsVisible() {
		return languageSelector.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Tool selector takes fourth priority (Small)
	if toolSelector.IsVisible() {
		return toolSelector.View(pr.layout.PopupSmallWidth, pr.layout.PopupSmallHeight)
	}

	// Weapon mastery selector takes fifth priority (Medium)
	if weaponMasterySelector.IsVisible() {
		return weaponMasterySelector.View()
	}

	// Expertise selector (Medium)
	if expertiseSelector.IsVisible() {
		return expertiseSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Maneuver selector (Medium)
	if maneuverSelector.IsVisible() {
		return maneuverSelector.View()
	}

	// Level-up selector takes sixth priority (Medium/Large)
	if levelUpSelector.IsVisible() {
		return levelUpSelector.View()
	}

	// De-level selector takes priority after level-up (Medium)
	if deLevelSelector.IsVisible() {
		return deLevelSelector.View(pr.width, pr.height)
	}

	// Item selector takes seventh priority (Large)
	if itemSelector.IsVisible() {
		return itemSelector.View(pr.layout.PopupLargeWidth, pr.layout.PopupLargeHeight)
	}

	// Fighting style selector takes sixth priority (Medium)
	if fightingStyleSelector.IsVisible() {
		return fightingStyleSelector.View(pr.width, pr.height)
	}

	// Cantrip selector takes seventh priority (Medium)
	if cantripSelector.IsVisible() {
		return cantripSelector.View()
	}

	// Message popup takes highest priority (shown after selections complete)
	if messagePopup.IsVisible() {
		return messagePopup.View(pr.width, pr.height)
	}

	// Leveled spell selector takes priority
	if leveledSpellSelector.IsVisible() {
		return leveledSpellSelector.View(pr.width, pr.height)
	}

	// School spell selector takes priority
	if schoolSpellSelector.IsVisible() {
		return schoolSpellSelector.View(pr.width, pr.height)
	}

	// Spellbook editor takes priority
	if spellbookEditor.IsVisible() {
		return spellbookEditor.View(pr.width, pr.height)
	}

	// Spell prep selector takes eighth priority (Medium)
	if spellPrepSelector.IsVisible() {
		return spellPrepSelector.View()
	}

	// Slot restorer takes ninth priority (Small)
	if slotRestorer.IsVisible() {
		return slotRestorer.View()
	}

	// Class skill selector takes tenth priority (Medium)
	if classSkillSelector.IsVisible() {
		return classSkillSelector.View(pr.width, pr.height)
	}

	// Subclass selector takes priority after skills (Medium)
	if subclassSelector.IsVisible() {
		return subclassSelector.View()
	}

	// Beast selector takes priority after subclass (Medium)
	if beastSelector.IsVisible() {
		return beastSelector.View()
	}

	// Option selector (for Fey Gift, Hunter's Prey, etc.) (Medium)
	if optionSelector != nil && optionSelector.IsVisible() {
		return optionSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Companion selector (Medium)
	if companionSelector != nil && companionSelector.IsVisible() {
		return companionSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Wild Shape form selector (Medium)
	if wildShapeFormSelector != nil && wildShapeFormSelector.IsVisible() {
		return wildShapeFormSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Divine Order selector (for Cleric level 1) - show as popup (Medium)
	if divineOrderSelectorVisible {
		return renderDivineOrderSelector()
	}

	// Primal Order selector (for Druid level 1) - show as popup (Medium)
	if primalOrderSelectorVisible {
		return renderPrimalOrderSelector()
	}

	// Class selector takes seventh priority (Medium)
	if classSelector.IsVisible() {
		return classSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Species selector takes seventh priority (Medium)
	if speciesSelector.IsVisible() {
		return speciesSelector.View(pr.layout.PopupMediumWidth, pr.layout.PopupMediumHeight)
	}

	// Rest popup overlay if active
	if restPopup.IsVisible() {
		return restPopup.View()
	}

	// Attack menu takes priority (shows as centered overlay)
	// Note: This will hide the TUI underneath for simplicity
	if attackMenu.IsVisible() {
		return attackMenu.View(pr.width, pr.height)
	}

	// No popup to show
	return ""
}

// RenderMainView renders the main view with panels
func RenderMainView(
	width, height int,
	focusArea int, // FocusArea type from app.go
	layout *LayoutDimensions,
	tabBar string,
	mainPanelView string,
	charStatsView string,
	actionsView string,
	diceView string,
	statusBar string,
) string {
	// Main panel styling
	mainPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(layout.MainPanelWidth).
		Height(layout.TopRowHeight)

	const FocusMain = 0
	const FocusCharStats = 1
	const FocusActions = 2
	const FocusDice = 3

	if focusArea == FocusMain {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		mainPanelStyle = mainPanelStyle.BorderForeground(lipgloss.Color("240"))
	}

	// Combine tabs and content vertically
	tabsAndContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		"", // Add a line of spacing
		mainPanelView,
	)

	mainPanelWithTabs := mainPanelStyle.Render(tabsAndContent)

	// Character stats panel styling
	charStatsPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(layout.CharStatsWidth).
		Height(layout.TopRowHeight)

	if focusArea == FocusCharStats {
		charStatsPanelStyle = charStatsPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		charStatsPanelStyle = charStatsPanelStyle.BorderForeground(lipgloss.Color("86"))
	}

	charStatsWithBorder := charStatsPanelStyle.Render(charStatsView)

	// Join main panel (with tabs) and character stats horizontally
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainPanelWithTabs,
		charStatsWithBorder,
	)

	// Actions panel styling
	actionsPanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(layout.ActionsWidthRatio).
		Height(layout.BottomHeight)

	if focusArea == FocusActions {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		actionsPanelStyle = actionsPanelStyle.BorderForeground(lipgloss.Color("240"))
	}

	actionsViewStyled := actionsPanelStyle.Render(actionsView)

	// Dice panel styling
	dicePanelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(layout.DiceWidthRatio).
		Height(layout.BottomHeight)

	if focusArea == FocusDice {
		dicePanelStyle = dicePanelStyle.BorderForeground(lipgloss.Color("205"))
	} else {
		dicePanelStyle = dicePanelStyle.BorderForeground(lipgloss.Color("240"))
	}

	diceViewStyled := dicePanelStyle.Render(diceView)

	// Bottom row (actions + dice)
	bottomRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		actionsViewStyled,
		diceViewStyled,
	)

	// Combine all parts vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		bottomRow,
		statusBar,
	)
}
