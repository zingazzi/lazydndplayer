// internal/ui/interfaces_test.go
package ui

import (
	"errors"
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// mockStorage implements StorageInterface for testing
type mockStorage struct {
	saveError   error
	loadError   error
	savedChar   *models.Character
	loadedChar  *models.Character
	saveCalled  bool
	loadCalled  bool
	exportError error
	importError error
}

func (m *mockStorage) Save(character *models.Character) error {
	m.saveCalled = true
	m.savedChar = character
	return m.saveError
}

func (m *mockStorage) Load() (*models.Character, error) {
	m.loadCalled = true
	if m.loadError != nil {
		return nil, m.loadError
	}
	if m.loadedChar != nil {
		return m.loadedChar, nil
	}
	return models.NewCharacter(), nil
}

func (m *mockStorage) Export(character *models.Character, exportPath string) error {
	return m.exportError
}

func (m *mockStorage) Import(importPath string) (*models.Character, error) {
	if m.importError != nil {
		return nil, m.importError
	}
	return m.loadedChar, nil
}

func TestStorageInterface(t *testing.T) {
	// Verify that mockStorage implements StorageInterface
	var _ StorageInterface = (*mockStorage)(nil)

	mock := &mockStorage{}

	char := models.NewCharacter()
	char.Name = "Test Character"

	// Test Save
	err := mock.Save(char)
	if err != nil {
		t.Errorf("Save() returned error: %v", err)
	}

	if !mock.saveCalled {
		t.Error("Save() was not called")
	}

	if mock.savedChar != char {
		t.Error("Save() did not store the character")
	}

	// Test Load
	loaded, err := mock.Load()
	if err != nil {
		t.Errorf("Load() returned error: %v", err)
	}

	if !mock.loadCalled {
		t.Error("Load() was not called")
	}

	if loaded == nil {
		t.Error("Load() returned nil character")
	}
}

func TestStorageInterfaceErrors(t *testing.T) {
	mock := &mockStorage{
		saveError: errors.New("save error"),
		loadError: errors.New("load error"),
	}

	char := models.NewCharacter()

	// Test Save error
	err := mock.Save(char)
	if err == nil {
		t.Error("Expected Save() to return error")
	}

	// Test Load error
	_, err = mock.Load()
	if err == nil {
		t.Error("Expected Load() to return error")
	}
}

func TestComponentFactoryInterface(t *testing.T) {
	// Verify that ComponentFactory implements ComponentFactoryInterface
	var _ ComponentFactoryInterface = (*ComponentFactory)(nil)

	factory := NewComponentFactory()

	// Test that factory is not nil
	if factory == nil {
		t.Fatal("NewComponentFactory() returned nil")
	}

	// Test creating a simple component
	tabs := factory.CreateTabs()
	if tabs == nil {
		t.Error("CreateTabs() returned nil")
	}

	help := factory.CreateHelp()
	if help == nil {
		t.Error("CreateHelp() returned nil")
	}
}

func TestComponentFactoryCreatesAllComponents(t *testing.T) {
	factory := NewComponentFactory()
	char := models.NewCharacter()

	// Test all component creation methods
	components := []struct {
		name string
		test func() bool
	}{
		{"Tabs", func() bool { return factory.CreateTabs() != nil }},
		{"Help", func() bool { return factory.CreateHelp() != nil }},
		{"SpeciesSelector", func() bool { return factory.CreateSpeciesSelector() != nil }},
		{"SubtypeSelector", func() bool { return factory.CreateSubtypeSelector() != nil }},
		{"LanguageSelector", func() bool { return factory.CreateLanguageSelector() != nil }},
		{"SkillSelector", func() bool { return factory.CreateSkillSelector() != nil }},
		{"SpellSelector", func() bool { return factory.CreateSpellSelector() != nil }},
		{"FeatSelector", func() bool { return factory.CreateFeatSelector() != nil }},
		{"FeatDetailPopup", func() bool { return factory.CreateFeatDetailPopup() != nil }},
		{"FeatureDetailPopup", func() bool { return factory.CreateFeatureDetailPopup() != nil }},
		{"ItemDetailPopup", func() bool { return factory.CreateItemDetailPopup() != nil }},
		{"MasteryDetailPopup", func() bool { return factory.CreateMasteryDetailPopup() != nil }},
		{"ManeuverDetailPopup", func() bool { return factory.CreateManeuverDetailPopup() != nil }},
		{"ConsumableDetailPopup", func() bool { return factory.CreateConsumableDetailPopup() != nil }},
		{"SpellDetailPopup", func() bool { return factory.CreateSpellDetailPopup() != nil }},
		{"OriginSelector", func() bool { return factory.CreateOriginSelector() != nil }},
		{"AlignmentSelector", func() bool { return factory.CreateAlignmentSelector() != nil }},
		{"TraitSelector", func() bool { return factory.CreateTraitSelector() != nil }},
		{"BackstoryEditor", func() bool { return factory.CreateBackstoryEditor() != nil }},
		{"OriginDetailPopup", func() bool { return factory.CreateOriginDetailPopup() != nil }},
		{"InputPopup", func() bool { return factory.CreateInputPopup() != nil }},
		{"ToolSelector", func() bool { return factory.CreateToolSelector() != nil }},
		{"ItemSelector", func() bool { return factory.CreateItemSelector() != nil }},
		{"ClassSelector", func() bool { return factory.CreateClassSelector(char) != nil }},
		{"ClassSkillSelector", func() bool { return factory.CreateClassSkillSelector() != nil }},
		{"SubclassSelector", func() bool { return factory.CreateSubclassSelector(char) != nil }},
		{"FightingStyleSelector", func() bool { return factory.CreateFightingStyleSelector() != nil }},
		{"CantripSelector", func() bool { return factory.CreateCantripSelector(char) != nil }},
		{"LeveledSpellSelector", func() bool { return factory.CreateLeveledSpellSelector(char) != nil }},
		{"SchoolSpellSelector", func() bool { return factory.CreateSchoolSpellSelector(char) != nil }},
		{"SpellbookEditor", func() bool { return factory.CreateSpellbookEditor(char) != nil }},
		{"SpellPrepSelector", func() bool { return factory.CreateSpellPrepSelector(char) != nil }},
		{"SlotRestorer", func() bool { return factory.CreateSlotRestorer(char) != nil }},
		{"StatGenerator", func() bool { return factory.CreateStatGenerator() != nil }},
		{"AbilityRoller", func() bool { return factory.CreateAbilityRoller() != nil }},
		{"AbilityChoiceSelector", func() bool { return factory.CreateAbilityChoiceSelector() != nil }},
		{"AttackRoller", func() bool { return factory.CreateAttackRoller() != nil }},
		{"AttackMenu", func() bool { return factory.CreateAttackMenu() != nil }},
		{"WeaponMasterySelector", func() bool { return factory.CreateWeaponMasterySelector(char) != nil }},
		{"ExpertiseSelector", func() bool { return factory.CreateExpertiseSelector(char) != nil }},
		{"ManeuverSelector", func() bool { return factory.CreateManeuverSelector() != nil }},
		{"LevelUpSelector", func() bool { return factory.CreateLevelUpSelector(char) != nil }},
		{"DeLevelSelector", func() bool { return factory.CreateDeLevelSelector(char) != nil }},
		{"MessagePopup", func() bool { return factory.CreateMessagePopup() != nil }},
		{"StatsPanel", func() bool { return factory.CreateStatsPanel(char) != nil }},
		{"SkillsPanel", func() bool { return factory.CreateSkillsPanel(char) != nil }},
		{"InventoryPanel", func() bool { return factory.CreateInventoryPanel(char) != nil }},
		{"SpellsPanel", func() bool { return factory.CreateSpellsPanel(char) != nil }},
		{"FeaturesPanel", func() bool { return factory.CreateFeaturesPanel(char) != nil }},
		{"TraitsPanel", func() bool { return factory.CreateTraitsPanel(char) != nil }},
		{"OriginPanel", func() bool { return factory.CreateOriginPanel(char) != nil }},
		{"DicePanel", func() bool { return factory.CreateDicePanel(char) != nil }},
		{"CharacterStatsPanel", func() bool { return factory.CreateCharacterStatsPanel(char) != nil }},
		{"ActionsPanel", func() bool { return factory.CreateActionsPanel(char) != nil }},
	}

	for _, comp := range components {
		t.Run(comp.name, func(t *testing.T) {
			if !comp.test() {
				t.Errorf("Create%s() returned nil", comp.name)
			}
		})
	}
}
