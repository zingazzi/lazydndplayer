// internal/ui/view/panel_renderer.go
package view

import (
	"github.com/marcozingoni/lazydndplayer/internal/ui/panels"
)

// PanelType represents the current active main panel
type PanelType int

const (
	StatsPanel PanelType = iota
	SkillsPanel
	InventoryPanel
	SpellsPanel
	FeaturesPanel
	TraitsPanel
	OriginPanel
	CompanionPanel
)

// PanelRenderer renders panels based on panel type
type PanelRenderer struct {
	panels map[PanelType]PanelView
}

// PanelView represents a panel that can be rendered
type PanelView interface {
	View(width, height int) string
}

// NewPanelRenderer creates a new panel renderer
func NewPanelRenderer() *PanelRenderer {
	return &PanelRenderer{
		panels: make(map[PanelType]PanelView),
	}
}

// RegisterPanel registers a panel for rendering
func (pr *PanelRenderer) RegisterPanel(panelType PanelType, panel PanelView) {
	pr.panels[panelType] = panel
}

// RenderPanel renders the specified panel with the given dimensions
func (pr *PanelRenderer) RenderPanel(panelType PanelType, width, height int) string {
	panel, exists := pr.panels[panelType]
	if !exists {
		return ""
	}
	return panel.View(width, height)
}

// RegisterPanelsFromModel registers all panels from the model
// This is a helper function to register all panels at once
func RegisterPanelsFromModel(
	statsPanel *panels.StatsPanel,
	skillsPanel *panels.SkillsPanel,
	inventoryPanel *panels.InventoryPanel,
	spellsPanel *panels.SpellsPanel,
	featuresPanel *panels.FeaturesPanel,
	traitsPanel *panels.TraitsPanel,
	originPanel *panels.OriginPanel,
	companionPanel *panels.CompanionPanel,
) *PanelRenderer {
	renderer := NewPanelRenderer()
	renderer.RegisterPanel(StatsPanel, statsPanel)
	renderer.RegisterPanel(SkillsPanel, skillsPanel)
	renderer.RegisterPanel(InventoryPanel, inventoryPanel)
	renderer.RegisterPanel(SpellsPanel, spellsPanel)
	renderer.RegisterPanel(FeaturesPanel, featuresPanel)
	renderer.RegisterPanel(TraitsPanel, traitsPanel)
	renderer.RegisterPanel(OriginPanel, originPanel)
	renderer.RegisterPanel(CompanionPanel, companionPanel)
	return renderer
}
