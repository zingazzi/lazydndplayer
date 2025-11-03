// internal/ui/handlers/selector_handler.go
package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marcozingoni/lazydndplayer/internal/models"
	"github.com/marcozingoni/lazydndplayer/internal/ui/components"
)

// HandleLanguageSelectorKeys handles language selector keyboard input
func HandleLanguageSelectorKeys(
	msg tea.KeyMsg,
	selector *components.LanguageSelector,
	character *models.Character,
	storage interface {
		Save(*models.Character)
	},
	message *string,
	onLanguageSelected func(string),
) (bool, string) {
	_ = NewBaseSelectorHandler(selector) // Base handler for future use

	switch msg.String() {
	case "up", "k":
		selector.Prev()
		return true, ""
	case "down", "j":
		selector.Next()
		return true, ""
	case "enter":
		selectedLanguage := selector.GetSelectedLanguage()
		if selectedLanguage != "" {
			// Check if we're in delete mode
			if selector.IsDeleteMode() {
				// Remove the language
				for i, lang := range character.Languages {
					if lang == selectedLanguage {
						character.Languages = append(character.Languages[:i], character.Languages[i+1:]...)
						break
					}
				}
				msg := "Language removed: " + selectedLanguage
				storage.Save(character)
				selector.Hide()
				return true, msg
			} else {
				// Add mode logic would go here
				// This is complex and depends on context, so we'll handle it in the caller
				if onLanguageSelected != nil {
					onLanguageSelected(selectedLanguage)
				}
				return true, ""
			}
		}
	case "esc":
		selector.Hide()
		return true, "Language selection cancelled"
	}

	return false, ""
}

// HandleSkillSelectorKeys handles skill selector keyboard input
func HandleSkillSelectorKeys(
	msg tea.KeyMsg,
	selector *components.SkillSelector,
	character *models.Character,
	storage interface {
		Save(*models.Character)
	},
	message *string,
) (bool, string) {
	baseHandler := NewBaseSelectorHandler(selector)
	if baseHandler.HandleNavigation(msg) {
		return true, ""
	}

	switch msg.String() {
	case "enter":
		selectedSkill := selector.GetSelectedSkill()
		if selectedSkill != "" {
			// Add skill logic would go here
			// This is complex and depends on context
			selector.Hide()
			return true, "Skill selected: " + selectedSkill
		}
	}

	return false, ""
}

// HandleToolSelectorKeys handles tool selector keyboard input
func HandleToolSelectorKeys(
	msg tea.KeyMsg,
	selector *components.ToolSelector,
	character *models.Character,
	storage interface {
		Save(*models.Character)
	},
	message *string,
) (bool, string) {
	baseHandler := NewBaseSelectorHandler(selector)
	if baseHandler.HandleNavigation(msg) {
		return true, ""
	}

	switch msg.String() {
	case "enter":
		selectedTool := selector.GetSelected()
		if selectedTool != "" {
			// Add/remove tool logic would go here
			// This is complex and depends on context
			selector.Hide()
			return true, "Tool selected: " + selectedTool
		}
	}

	return false, ""
}
