// internal/ui/handlers/handler_context.go
package handlers

import (
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/storage"
	"github.com/marcozingoni/lazydndplayer/internal/ui/state"
)

// HandlerContext provides handlers with access to the model's components and state
// This allows handlers to work with the model without tight coupling
type HandlerContext struct {
	Character   *models.Character
	Storage     *storage.Storage
	StateMachine *state.StateMachine
	Message     *string // Pointer to message field so handlers can update it
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(
	character *models.Character,
	storage *storage.Storage,
	stateMachine *state.StateMachine,
	message *string,
) *HandlerContext {
	return &HandlerContext{
		Character:    character,
		Storage:      storage,
		StateMachine: stateMachine,
		Message:      message,
	}
}

// SetMessage sets the message
func (ctx *HandlerContext) SetMessage(msg string) {
	if ctx.Message != nil {
		*ctx.Message = msg
	}
}
