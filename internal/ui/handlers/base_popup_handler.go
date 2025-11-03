// internal/ui/handlers/base_popup_handler.go
package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
)

// BasePopupHandler provides common behavior for popup components
type BasePopupHandler struct {
	popup PopupComponent
}

// NewBasePopupHandler creates a new base popup handler
func NewBasePopupHandler(popup PopupComponent) *BasePopupHandler {
	return &BasePopupHandler{
		popup: popup,
	}
}

// HandleClose handles common close keys (esc/enter)
// Returns true if the key was handled, false otherwise
func (h *BasePopupHandler) HandleClose(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "esc", "enter":
		h.popup.Hide()
		return true
	}
	return false
}
