// internal/ui/handlers/detail_popup_handler.go
package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// HandleDetailPopupKeys handles keyboard input for detail popups
// This consolidates the common esc/enter behavior for detail popups
func HandleDetailPopupKeys(
	msg tea.KeyMsg,
	popup PopupComponent,
	message *string,
	messageText string,
) bool {
	baseHandler := NewBasePopupHandler(popup)
	if baseHandler.HandleClose(msg) {
		if message != nil {
			*message = messageText
		}
		return true
	}
	return false
}

// HandleFeatDetailPopupKeys handles feat detail popup keys
func HandleFeatDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.FeatDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed feat details")
}

// HandleMasteryDetailPopupKeys handles mastery detail popup keys
func HandleMasteryDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.MasteryDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed weapon mastery details")
}

// HandleManeuverDetailPopupKeys handles maneuver detail popup keys
func HandleManeuverDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.ManeuverDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed maneuver details")
}

// HandleConsumableDetailPopupKeys handles consumable detail popup keys
func HandleConsumableDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.ConsumableDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "")
}

// HandleFeatureDetailPopupKeys handles feature detail popup keys
func HandleFeatureDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.FeatureDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "")
}

// HandleItemDetailPopupKeys handles item detail popup keys
func HandleItemDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.ItemDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed item details")
}

// HandleSpellDetailPopupKeys handles spell detail popup keys
func HandleSpellDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.SpellDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed spell details")
}

// HandleOriginDetailPopupKeys handles origin detail popup keys
func HandleOriginDetailPopupKeys(
	msg tea.KeyMsg,
	popup *components.OriginDetailPopup,
	message *string,
) bool {
	return HandleDetailPopupKeys(msg, popup, message, "Closed origin details")
}
