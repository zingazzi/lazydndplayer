// internal/ui/component_manager.go
package ui

import "github.com/marcozingoni/lazydndplayer/internal/debug"

// ComponentHandler defines the interface that all UI components must implement
type ComponentHandler interface {
	IsVisible() bool
}

// ComponentEntry represents a registered component with its priority
type ComponentEntry struct {
	Component ComponentHandler
	Priority  int
	Name      string
}

// ComponentManager manages component visibility and routing
type ComponentManager struct {
	components []ComponentEntry
}

// NewComponentManager creates a new component manager
func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make([]ComponentEntry, 0),
	}
}

// Register registers a component with the manager
func (cm *ComponentManager) Register(component ComponentHandler, priority int, name string) {
	cm.components = append(cm.components, ComponentEntry{
		Component: component,
		Priority:  priority,
		Name:      name,
	})
}

// GetVisibleComponent returns the highest priority visible component
func (cm *ComponentManager) GetVisibleComponent() ComponentHandler {
	highestPriority := -1
	var result ComponentHandler
	var resultName string

	for _, entry := range cm.components {
		if entry.Component.IsVisible() && entry.Priority > highestPriority {
			highestPriority = entry.Priority
			result = entry.Component
			resultName = entry.Name
		}
	}

	// Only log when we have a result to reduce log spam
	if result != nil {
		debug.Log("ComponentManager: Returning visible component '%s' (priority %d)", resultName, highestPriority)
	}

	return result
}

// GetVisibleComponents returns all visible components sorted by priority (highest first)
func (cm *ComponentManager) GetVisibleComponents() []ComponentHandler {
	visible := make([]ComponentEntry, 0)

	// Collect visible components
	for _, entry := range cm.components {
		if entry.Component.IsVisible() {
			visible = append(visible, entry)
		}
	}

	// Sort by priority (highest first)
	for i := 0; i < len(visible)-1; i++ {
		for j := i + 1; j < len(visible); j++ {
			if visible[i].Priority < visible[j].Priority {
				visible[i], visible[j] = visible[j], visible[i]
			}
		}
	}

	// Extract components
	result := make([]ComponentHandler, len(visible))
	for i, entry := range visible {
		result[i] = entry.Component
	}

	return result
}

// HasVisibleComponent checks if any component is visible
func (cm *ComponentManager) HasVisibleComponent() bool {
	return cm.GetVisibleComponent() != nil
}
