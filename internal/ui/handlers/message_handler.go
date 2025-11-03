// internal/ui/handlers/message_handler.go
package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// HandleMessagePopupKeys handles message popup keyboard input
func HandleMessagePopupKeys(
	msg tea.KeyMsg,
	popup *components.MessagePopup,
) bool {
	switch msg.String() {
	case "esc", "enter", " ":
		popup.Hide()
		return true
	}
	return false
}

// HandleRestPopupKeys handles rest popup keyboard input
// Returns true if handled, and whether rest was confirmed
func HandleRestPopupKeys(
	msg tea.KeyMsg,
	popup *components.RestPopup,
) (bool, bool) {
	// Rest popup has its own Update method that handles the message
	// We just need to check if it's confirmed
	if popup.IsConfirmed() {
		return true, true
	}
	return false, false
}
