// internal/ui/interfaces.go
package ui

import (
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
)

// Ensure storage.Storage implements StorageInterface at compile time
var _ StorageInterface = (*storage.Storage)(nil)

// StorageInterface defines the interface for character storage operations
// This allows for easy mocking in tests and future storage implementations (e.g., database)
type StorageInterface interface {
	Save(character *models.Character) error
	Load() (*models.Character, error)
	Export(character *models.Character, exportPath string) error
	Import(importPath string) (*models.Character, error)
}

// ComponentFactoryInterface defines the interface for creating UI components
// This allows for dependency injection and easier testing
type ComponentFactoryInterface interface {
	CreateTabs() *components.Tabs
	CreateHelp() *components.Help
	CreateSpeciesSelector() *components.SpeciesSelector
	CreateSubtypeSelector() *components.SubtypeSelector
	CreateLanguageSelector() *components.LanguageSelector
	CreateSkillSelector() *components.SkillSelector
	CreateSpellSelector() *components.SpellSelector
	CreateFeatSelector() *components.FeatSelector
	CreateFeatDetailPopup() *components.FeatDetailPopup
	CreateFeatureDetailPopup() *components.FeatureDetailPopup
	CreateItemDetailPopup() *components.ItemDetailPopup
	CreateMasteryDetailPopup() *components.MasteryDetailPopup
	CreateManeuverDetailPopup() *components.ManeuverDetailPopup
	CreateConsumableDetailPopup() *components.ConsumableDetailPopup
	CreateSpellDetailPopup() *components.SpellDetailPopup
	CreateOriginSelector() *components.OriginSelector
	CreateAlignmentSelector() *components.AlignmentSelector
	CreateTraitSelector() *components.TraitSelector
	CreateBackstoryEditor() *components.BackstoryEditor
	CreateOriginDetailPopup() *components.OriginDetailPopup
	CreateInputPopup() *components.InputPopup
	CreateToolSelector() *components.ToolSelector
	CreateItemSelector() *components.ItemSelector
	CreateClassSelector(character *models.Character) *components.ClassSelector
	CreateClassSkillSelector() *components.ClassSkillSelector
	CreateSubclassSelector(character *models.Character) *components.SubclassSelector
	CreateFightingStyleSelector() *components.FightingStyleSelector
	CreateCantripSelector(character *models.Character) *components.CantripSelector
	CreateLeveledSpellSelector(character *models.Character) *components.LeveledSpellSelector
	CreateSchoolSpellSelector(character *models.Character) *components.SchoolSpellSelector
	CreateSpellbookEditor(character *models.Character) *components.SpellbookEditor
	CreateSpellPrepSelector(character *models.Character) *components.SpellPrepSelector
	CreateSlotRestorer(character *models.Character) *components.SlotRestorer
	CreateStatGenerator() *components.StatGenerator
	CreateAbilityRoller() *components.AbilityRoller
	CreateAbilityChoiceSelector() *components.AbilityChoiceSelector
	CreateAttackRoller() *components.AttackRoller
	CreateAttackMenu() *components.AttackMenu
	CreateWeaponMasterySelector(character *models.Character) *components.WeaponMasterySelector
	CreateExpertiseSelector(character *models.Character) *components.ExpertiseSelector
	CreateManeuverSelector() *components.ManeuverSelector
	CreateLevelUpSelector(character *models.Character) *components.LevelUpSelector
	CreateDeLevelSelector(character *models.Character) *components.DeLevelSelector
	CreateRestPopup(character *models.Character, diceRoller models.DiceRoller) *components.RestPopup
	CreateMessagePopup() *components.MessagePopup
	CreateStatsPanel(character *models.Character) *panels.StatsPanel
	CreateSkillsPanel(character *models.Character) *panels.SkillsPanel
	CreateInventoryPanel(character *models.Character) *panels.InventoryPanel
	CreateSpellsPanel(character *models.Character) *panels.SpellsPanel
	CreateFeaturesPanel(character *models.Character) *panels.FeaturesPanel
	CreateTraitsPanel(character *models.Character) *panels.TraitsPanel
	CreateOriginPanel(character *models.Character) *panels.OriginPanel
	CreateCompanionPanel(character *models.Character) *panels.CompanionPanel
	CreateDicePanel(character *models.Character) *panels.DicePanel
	CreateCharacterStatsPanel(character *models.Character) *panels.CharacterStatsPanel
	CreateActionsPanel(character *models.Character) *panels.ActionsPanel
}

// ComponentFactory is the default implementation of ComponentFactoryInterface
type ComponentFactory struct{}

// NewComponentFactory creates a new component factory
func NewComponentFactory() *ComponentFactory {
	return &ComponentFactory{}
}

// Component creation methods - delegate to actual component constructors
func (f *ComponentFactory) CreateTabs() *components.Tabs {
	return components.NewTabs()
}

func (f *ComponentFactory) CreateHelp() *components.Help {
	return components.NewHelp()
}

func (f *ComponentFactory) CreateSpeciesSelector() *components.SpeciesSelector {
	return components.NewSpeciesSelector()
}

func (f *ComponentFactory) CreateSubtypeSelector() *components.SubtypeSelector {
	return components.NewSubtypeSelector()
}

func (f *ComponentFactory) CreateLanguageSelector() *components.LanguageSelector {
	return components.NewLanguageSelector()
}

func (f *ComponentFactory) CreateSkillSelector() *components.SkillSelector {
	return components.NewSkillSelector()
}

func (f *ComponentFactory) CreateSpellSelector() *components.SpellSelector {
	return components.NewSpellSelector()
}

func (f *ComponentFactory) CreateFeatSelector() *components.FeatSelector {
	return components.NewFeatSelector()
}

func (f *ComponentFactory) CreateFeatDetailPopup() *components.FeatDetailPopup {
	return components.NewFeatDetailPopup()
}

func (f *ComponentFactory) CreateFeatureDetailPopup() *components.FeatureDetailPopup {
	return components.NewFeatureDetailPopup()
}

func (f *ComponentFactory) CreateItemDetailPopup() *components.ItemDetailPopup {
	return components.NewItemDetailPopup()
}

func (f *ComponentFactory) CreateMasteryDetailPopup() *components.MasteryDetailPopup {
	return components.NewMasteryDetailPopup()
}

func (f *ComponentFactory) CreateManeuverDetailPopup() *components.ManeuverDetailPopup {
	return components.NewManeuverDetailPopup()
}

func (f *ComponentFactory) CreateConsumableDetailPopup() *components.ConsumableDetailPopup {
	return components.NewConsumableDetailPopup()
}

func (f *ComponentFactory) CreateSpellDetailPopup() *components.SpellDetailPopup {
	return components.NewSpellDetailPopup()
}

func (f *ComponentFactory) CreateOriginSelector() *components.OriginSelector {
	return components.NewOriginSelector()
}

func (f *ComponentFactory) CreateAlignmentSelector() *components.AlignmentSelector {
	return components.NewAlignmentSelector()
}

func (f *ComponentFactory) CreateTraitSelector() *components.TraitSelector {
	return components.NewTraitSelector()
}

func (f *ComponentFactory) CreateBackstoryEditor() *components.BackstoryEditor {
	return components.NewBackstoryEditor()
}

func (f *ComponentFactory) CreateOriginDetailPopup() *components.OriginDetailPopup {
	return components.NewOriginDetailPopup()
}

func (f *ComponentFactory) CreateInputPopup() *components.InputPopup {
	return components.NewInputPopup()
}

func (f *ComponentFactory) CreateToolSelector() *components.ToolSelector {
	return components.NewToolSelector()
}

func (f *ComponentFactory) CreateItemSelector() *components.ItemSelector {
	return components.NewItemSelector()
}

func (f *ComponentFactory) CreateClassSelector(character *models.Character) *components.ClassSelector {
	return components.NewClassSelector(character)
}

func (f *ComponentFactory) CreateClassSkillSelector() *components.ClassSkillSelector {
	return components.NewClassSkillSelector()
}

func (f *ComponentFactory) CreateSubclassSelector(character *models.Character) *components.SubclassSelector {
	return components.NewSubclassSelector(character)
}

func (f *ComponentFactory) CreateFightingStyleSelector() *components.FightingStyleSelector {
	return components.NewFightingStyleSelector()
}

func (f *ComponentFactory) CreateCantripSelector(character *models.Character) *components.CantripSelector {
	return components.NewCantripSelector(character)
}

func (f *ComponentFactory) CreateLeveledSpellSelector(character *models.Character) *components.LeveledSpellSelector {
	return components.NewLeveledSpellSelector(character)
}

func (f *ComponentFactory) CreateSchoolSpellSelector(character *models.Character) *components.SchoolSpellSelector {
	return components.NewSchoolSpellSelector(character)
}

func (f *ComponentFactory) CreateSpellbookEditor(character *models.Character) *components.SpellbookEditor {
	return components.NewSpellbookEditor(character)
}

func (f *ComponentFactory) CreateSpellPrepSelector(character *models.Character) *components.SpellPrepSelector {
	return components.NewSpellPrepSelector(character)
}

func (f *ComponentFactory) CreateSlotRestorer(character *models.Character) *components.SlotRestorer {
	return components.NewSlotRestorer(character)
}

func (f *ComponentFactory) CreateStatGenerator() *components.StatGenerator {
	return components.NewStatGenerator()
}

func (f *ComponentFactory) CreateAbilityRoller() *components.AbilityRoller {
	return components.NewAbilityRoller()
}

func (f *ComponentFactory) CreateAbilityChoiceSelector() *components.AbilityChoiceSelector {
	return components.NewAbilityChoiceSelector()
}

func (f *ComponentFactory) CreateAttackRoller() *components.AttackRoller {
	return components.NewAttackRoller()
}

func (f *ComponentFactory) CreateAttackMenu() *components.AttackMenu {
	return components.NewAttackMenu()
}

func (f *ComponentFactory) CreateWeaponMasterySelector(character *models.Character) *components.WeaponMasterySelector {
	return components.NewWeaponMasterySelector(character)
}

func (f *ComponentFactory) CreateExpertiseSelector(character *models.Character) *components.ExpertiseSelector {
	return components.NewExpertiseSelector(character)
}

func (f *ComponentFactory) CreateManeuverSelector() *components.ManeuverSelector {
	return components.NewManeuverSelector()
}

func (f *ComponentFactory) CreateLevelUpSelector(character *models.Character) *components.LevelUpSelector {
	return components.NewLevelUpSelector(character)
}

func (f *ComponentFactory) CreateDeLevelSelector(character *models.Character) *components.DeLevelSelector {
	return components.NewDeLevelSelector(character)
}

func (f *ComponentFactory) CreateRestPopup(character *models.Character, diceRoller models.DiceRoller) *components.RestPopup {
	return components.NewRestPopup(character, diceRoller)
}

func (f *ComponentFactory) CreateMessagePopup() *components.MessagePopup {
	return components.NewMessagePopup()
}

func (f *ComponentFactory) CreateStatsPanel(character *models.Character) *panels.StatsPanel {
	return panels.NewStatsPanel(character)
}

func (f *ComponentFactory) CreateSkillsPanel(character *models.Character) *panels.SkillsPanel {
	return panels.NewSkillsPanel(character)
}

func (f *ComponentFactory) CreateInventoryPanel(character *models.Character) *panels.InventoryPanel {
	return panels.NewInventoryPanel(character)
}

func (f *ComponentFactory) CreateSpellsPanel(character *models.Character) *panels.SpellsPanel {
	return panels.NewSpellsPanel(character)
}

func (f *ComponentFactory) CreateFeaturesPanel(character *models.Character) *panels.FeaturesPanel {
	return panels.NewFeaturesPanel(character)
}

func (f *ComponentFactory) CreateTraitsPanel(character *models.Character) *panels.TraitsPanel {
	return panels.NewTraitsPanel(character)
}

func (f *ComponentFactory) CreateOriginPanel(character *models.Character) *panels.OriginPanel {
	return panels.NewOriginPanel(character)
}

func (f *ComponentFactory) CreateCompanionPanel(character *models.Character) *panels.CompanionPanel {
	return panels.NewCompanionPanel(character)
}

func (f *ComponentFactory) CreateDicePanel(character *models.Character) *panels.DicePanel {
	return panels.NewDicePanel(character)
}

func (f *ComponentFactory) CreateCharacterStatsPanel(character *models.Character) *panels.CharacterStatsPanel {
	return panels.NewCharacterStatsPanel(character)
}

func (f *ComponentFactory) CreateActionsPanel(character *models.Character) *panels.ActionsPanel {
	return panels.NewActionsPanel(character)
}
