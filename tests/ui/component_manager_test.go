// tests/ui/component_manager_test.go
package ui_test

import (
	"testing"

	"github.com/marcozingoni/lazydndplayer/internal/ui"
)

// mockComponent is a test component for ComponentManager tests
type mockComponent struct {
	visible bool
	name    string
}

func (m *mockComponent) IsVisible() bool {
	return m.visible
}

func TestNewComponentManager(t *testing.T) {
	cm := ui.NewComponentManager()

	if cm == nil {
		t.Fatal("NewComponentManager() returned nil")
	}

	if cm.HasVisibleComponent() {
		t.Error("Expected no visible components initially")
	}

	if cm.GetVisibleComponent() != nil {
		t.Error("Expected no visible component initially")
	}
}

func TestRegister(t *testing.T) {
	cm := ui.NewComponentManager()
	comp1 := &mockComponent{visible: true, name: "comp1"}
	comp2 := &mockComponent{visible: false, name: "comp2"}

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 20, "comp2")

	if !cm.HasVisibleComponent() {
		t.Error("Expected visible component after registration")
	}

	visible := cm.GetVisibleComponent()
	if visible != comp1 {
		t.Errorf("Expected comp1 to be visible, got %v", visible)
	}
}

func TestGetVisibleComponentPriority(t *testing.T) {
	cm := ui.NewComponentManager()

	comp1 := &mockComponent{visible: true, name: "comp1"}  // Priority 10
	comp2 := &mockComponent{visible: true, name: "comp2"}  // Priority 20
	comp3 := &mockComponent{visible: true, name: "comp3"}  // Priority 15

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 20, "comp2")
	cm.Register(comp3, 15, "comp3")

	visible := cm.GetVisibleComponent()
	if visible != comp2 {
		t.Errorf("Expected comp2 (priority 20) to be visible, got %v", visible)
	}
}

func TestGetVisibleComponentNoneVisible(t *testing.T) {
	cm := ui.NewComponentManager()

	comp1 := &mockComponent{visible: false, name: "comp1"}
	comp2 := &mockComponent{visible: false, name: "comp2"}

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 20, "comp2")

	if cm.HasVisibleComponent() {
		t.Error("Expected no visible components")
	}

	if cm.GetVisibleComponent() != nil {
		t.Error("Expected nil when no components are visible")
	}
}

func TestGetVisibleComponents(t *testing.T) {
	cm := ui.NewComponentManager()

	comp1 := &mockComponent{visible: true, name: "comp1"}  // Priority 10
	comp2 := &mockComponent{visible: false, name: "comp2"} // Not visible
	comp3 := &mockComponent{visible: true, name: "comp3"}  // Priority 30
	comp4 := &mockComponent{visible: true, name: "comp4"}  // Priority 20

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 20, "comp2")
	cm.Register(comp3, 30, "comp3")
	cm.Register(comp4, 20, "comp4")

	visible := cm.GetVisibleComponents()

	if len(visible) != 3 {
		t.Errorf("Expected 3 visible components, got %d", len(visible))
	}

	// Check order: should be sorted by priority (highest first)
	if visible[0] != comp3 {
		t.Error("Expected comp3 (priority 30) to be first")
	}
	if visible[1] != comp4 {
		t.Error("Expected comp4 (priority 20) to be second")
	}
	if visible[2] != comp1 {
		t.Error("Expected comp1 (priority 10) to be third")
	}
}

func TestHasVisibleComponent(t *testing.T) {
	cm := ui.NewComponentManager()

	if cm.HasVisibleComponent() {
		t.Error("Expected no visible components initially")
	}

	comp1 := &mockComponent{visible: true, name: "comp1"}
	cm.Register(comp1, 10, "comp1")

	if !cm.HasVisibleComponent() {
		t.Error("Expected visible component after registration")
	}

	comp1.visible = false
	if cm.HasVisibleComponent() {
		t.Error("Expected no visible components after hiding")
	}
}

func TestMultipleComponentsSamePriority(t *testing.T) {
	cm := ui.NewComponentManager()

	comp1 := &mockComponent{visible: true, name: "comp1"}
	comp2 := &mockComponent{visible: true, name: "comp2"}

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 10, "comp2")

	visible := cm.GetVisibleComponent()
	// When priorities are equal, the first one registered should be returned
	if visible != comp1 {
		t.Errorf("Expected comp1 (registered first), got %v", visible)
	}
}

func TestComponentVisibilityChange(t *testing.T) {
	cm := ui.NewComponentManager()

	comp1 := &mockComponent{visible: false, name: "comp1"}
	comp2 := &mockComponent{visible: false, name: "comp2"}

	cm.Register(comp1, 10, "comp1")
	cm.Register(comp2, 20, "comp2")

	if cm.GetVisibleComponent() != nil {
		t.Error("Expected no visible component")
	}

	comp1.visible = true
	if cm.GetVisibleComponent() != comp1 {
		t.Error("Expected comp1 to be visible after changing visibility")
	}

	comp1.visible = false
	comp2.visible = true
	if cm.GetVisibleComponent() != comp2 {
		t.Error("Expected comp2 to be visible after changing visibility")
	}
}
