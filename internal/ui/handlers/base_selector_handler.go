// internal/ui/handlers/base_selector_handler.go
package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
)

// BaseSelectorHandler provides common navigation logic for selector components
type BaseSelectorHandler struct {
	selector SelectorComponent
}

// NewBaseSelectorHandler creates a new base selector handler
func NewBaseSelectorHandler(selector SelectorComponent) *BaseSelectorHandler {
	return &BaseSelectorHandler{
		selector: selector,
	}
}

// HandleNavigation handles common navigation keys (up/down/esc)
// This can be used by selectors that follow the standard pattern
func (h *BaseSelectorHandler) HandleNavigation(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "up", "k":
		h.selector.Prev()
		return true
	case "down", "j":
		h.selector.Next()
		return true
	case "esc":
		h.selector.Hide()
		return true
	}
	return false
}

// HandlePageNavigation handles page navigation keys
func (h *BaseSelectorHandler) HandlePageNavigation(msg tea.KeyMsg, selector PageNavigationSelector) bool {
	switch msg.String() {
	case "pgup", "ctrl+u":
		selector.PageUp()
		return true
	case "pgdown", "ctrl+d":
		selector.PageDown()
		return true
	}
	return false
}

// HandleCategoryNavigation handles category navigation keys
func (h *BaseSelectorHandler) HandleCategoryNavigation(msg tea.KeyMsg, selector CategoryNavigationSelector) bool {
	switch msg.String() {
	case "left", "h":
		selector.PrevCategory()
		return true
	case "right", "l":
		selector.NextCategory()
		return true
	}
	return false
}

// PageNavigationSelector defines selectors with page navigation
type PageNavigationSelector interface {
	SelectorComponent
	PageUp()
	PageDown()
}

// CategoryNavigationSelector defines selectors with category navigation
type CategoryNavigationSelector interface {
	SelectorComponent
	PrevCategory()
	NextCategory()
}
