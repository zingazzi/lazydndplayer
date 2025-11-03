// internal/ui/handlers/interfaces.go
package handlers

import tea "github.com/charmbracelet/bubbletea"

// ComponentHandler defines the interface that all UI components must implement
type ComponentHandler interface {
	IsVisible() bool
}

// KeyHandler defines the interface for handling keyboard input
type KeyHandler interface {
	HandleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd)
}

// SelectorComponent defines the interface for selector components
type SelectorComponent interface {
	ComponentHandler
	Next()
	Prev()
	Hide()
}

// PopupComponent defines the interface for popup components
type PopupComponent interface {
	ComponentHandler
	Hide()
}

// DetailPopupComponent defines the interface for detail popup components
type DetailPopupComponent interface {
	PopupComponent
	GetTitle() string
}
