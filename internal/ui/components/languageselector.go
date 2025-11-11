// internal/ui/components/languageselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Common D&D 5e languages
var dndLanguages = []string{
	"Abyssal",
	"Celestial",
	"Common",
	"Deep Speech",
	"Draconic",
	"Dwarvish",
	"Elvish",
	"Giant",
	"Gnomish",
	"Goblin",
	"Halfling",
	"Infernal",
	"Orc",
	"Primordial",
	"Sylvan",
	"Undercommon",
}

// LanguageSelector handles language selection UI
type LanguageSelector struct {
	allLanguages    []string
	languages       []string        // Filtered list (excluding already known OR showing only known)
	selectedIndex   int
	viewport        viewport.Model
	visible         bool
	deleteMode      bool            // If true, shows known languages for deletion
	title           string          // Custom title
	SelectedLanguages map[string]bool // Languages selected in this session (for multi-selection)
	existingLanguages map[string]bool // Languages character already has
	MaxChoices      int             // How many languages can be chosen (0 = unlimited for single selection)
	multiSelectMode bool            // If true, supports multi-selection with checkboxes
	excludeCommon   bool            // If true, exclude "Common" from the list
}

// NewLanguageSelector creates a new language selector
func NewLanguageSelector() *LanguageSelector {
	return &LanguageSelector{
		allLanguages:      dndLanguages,
		languages:         dndLanguages,
		selectedIndex:     0,
		visible:           false,
		deleteMode:        false,
		title:             "SELECT ADDITIONAL LANGUAGE",
		SelectedLanguages: make(map[string]bool),
		existingLanguages: make(map[string]bool),
		MaxChoices:        0, // 0 = single selection mode
		multiSelectMode:   false,
		excludeCommon:     false,
	}
}

// SetExcludeLanguages filters out languages the character already knows
func (ls *LanguageSelector) SetExcludeLanguages(knownLanguages []string) {
	// Create a set of known languages for fast lookup
	knownSet := make(map[string]bool)
	ls.existingLanguages = make(map[string]bool)
	for _, lang := range knownLanguages {
		// Normalize to handle variations like "One additional language of your choice"
		normalizedLang := strings.ToLower(strings.TrimSpace(lang))
		// Skip placeholder texts
		if !strings.Contains(normalizedLang, "additional") &&
		   !strings.Contains(normalizedLang, "choice") &&
		   !strings.Contains(normalizedLang, "extra") {
			knownSet[normalizedLang] = true
			ls.existingLanguages[lang] = true
		}
	}

	// Build language list
	ls.languages = []string{}
	for _, lang := range ls.allLanguages {
		// Exclude Common if excludeCommon is true
		if ls.excludeCommon && strings.ToLower(lang) == "common" {
			continue
		}
		// In multi-select mode, show all languages (including known ones)
		// In single-select mode, filter out known languages
		if ls.multiSelectMode || !knownSet[strings.ToLower(lang)] {
			ls.languages = append(ls.languages, lang)
		}
	}

	// Reset selected index
	ls.selectedIndex = 0
}

// ShowForMultiSelect displays the language selector with multi-selection support
func (ls *LanguageSelector) ShowForMultiSelect(knownLanguages []string, maxChoices int, excludeCommon bool) {
	ls.visible = true
	ls.selectedIndex = 0
	ls.deleteMode = false
	ls.multiSelectMode = true
	ls.MaxChoices = maxChoices
	ls.excludeCommon = excludeCommon
	ls.SelectedLanguages = make(map[string]bool)
	ls.title = fmt.Sprintf("SELECT %d LANGUAGE(S)", maxChoices)

	// Set up existing languages and filter
	ls.SetExcludeLanguages(knownLanguages)
}

// Show displays the language selector (for adding languages)
func (ls *LanguageSelector) Show() {
	ls.visible = true
	ls.selectedIndex = 0
	ls.deleteMode = false
	ls.multiSelectMode = false
	ls.MaxChoices = 0
	ls.excludeCommon = false
	ls.SelectedLanguages = make(map[string]bool)
	ls.title = "SELECT ADDITIONAL LANGUAGE"
}

// ShowForDeletion displays the language selector with known languages (for deleting)
func (ls *LanguageSelector) ShowForDeletion(knownLanguages []string) {
	ls.visible = true
	ls.selectedIndex = 0
	ls.deleteMode = true
	ls.title = "SELECT LANGUAGE TO REMOVE"

	// Filter out placeholder texts and set languages to known only
	ls.languages = []string{}
	for _, lang := range knownLanguages {
		normalizedLang := strings.ToLower(strings.TrimSpace(lang))
		// Skip placeholder texts
		if !strings.Contains(normalizedLang, "additional") &&
		   !strings.Contains(normalizedLang, "choice") &&
		   !strings.Contains(normalizedLang, "extra") {
			ls.languages = append(ls.languages, lang)
		}
	}

	// If no valid languages, close selector
	if len(ls.languages) == 0 {
		ls.visible = false
	}
}

// Hide hides the language selector
func (ls *LanguageSelector) Hide() {
	ls.visible = false
	ls.deleteMode = false
	ls.multiSelectMode = false
	ls.SelectedLanguages = make(map[string]bool)
}

// ToggleLanguage toggles selection of the current language (for multi-select mode)
func (ls *LanguageSelector) ToggleLanguage() bool {
	if !ls.multiSelectMode || ls.selectedIndex < 0 || ls.selectedIndex >= len(ls.languages) {
		return false
	}

	langName := ls.languages[ls.selectedIndex]

	// Can't select if already known
	if ls.existingLanguages[langName] {
		return false
	}

	// Toggle selection
	if ls.SelectedLanguages[langName] {
		delete(ls.SelectedLanguages, langName)
		return true
	}

	// Check if we can select more
	if ls.MaxChoices == 0 || len(ls.SelectedLanguages) < ls.MaxChoices {
		ls.SelectedLanguages[langName] = true
		return true
	}

	return false
}

// GetSelectedLanguages returns the list of selected languages (for multi-select mode)
func (ls *LanguageSelector) GetSelectedLanguages() []string {
	languages := []string{}
	for lang := range ls.SelectedLanguages {
		languages = append(languages, lang)
	}
	return languages
}

// CanConfirm returns true if the correct number of languages have been selected (for multi-select mode)
func (ls *LanguageSelector) CanConfirm() bool {
	if !ls.multiSelectMode {
		return false
	}
	return ls.MaxChoices > 0 && len(ls.SelectedLanguages) == ls.MaxChoices
}

// IsDeleteMode returns whether the selector is in delete mode
func (ls *LanguageSelector) IsDeleteMode() bool {
	return ls.deleteMode
}

// IsMultiSelectMode returns whether the selector is in multi-select mode
func (ls *LanguageSelector) IsMultiSelectMode() bool {
	return ls.multiSelectMode
}

// IsVisible returns whether the selector is visible
func (ls *LanguageSelector) IsVisible() bool {
	return ls.visible
}

// Next moves to the next language
func (ls *LanguageSelector) Next() {
	if ls.selectedIndex < len(ls.languages)-1 {
		ls.selectedIndex++
		ls.viewport.LineDown(1)
	}
}

// Prev moves to the previous language
func (ls *LanguageSelector) Prev() {
	if ls.selectedIndex > 0 {
		ls.selectedIndex--
		ls.viewport.LineUp(1)
	}
}

// GetSelectedLanguage returns the currently selected language
func (ls *LanguageSelector) GetSelectedLanguage() string {
	if ls.selectedIndex >= 0 && ls.selectedIndex < len(ls.languages) {
		return ls.languages[ls.selectedIndex]
	}
	return ""
}

// View renders the language selector
func (ls *LanguageSelector) View(screenWidth, screenHeight int) string {
	if !ls.visible {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	knownStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Italic(true)

	chosenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	counterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render(ls.title))
	content = append(content, "")

	// Show counter for multi-select mode
	if ls.multiSelectMode && ls.MaxChoices > 0 {
		content = append(content, counterStyle.Render(fmt.Sprintf("Choose %d language(s) • Selected: %d/%d", ls.MaxChoices, len(ls.SelectedLanguages), ls.MaxChoices)))
		content = append(content, "")
	}

	// Language list
	var languageList []string
	for i, lang := range ls.languages {
		var langLine string
		alreadyKnown := ls.existingLanguages[lang]
		alreadySelected := ls.SelectedLanguages[lang]

		// Build the line
		if ls.multiSelectMode {
			if alreadyKnown {
				langLine = fmt.Sprintf("  [✓] %s (Already Known)", lang)
				if i == ls.selectedIndex {
					languageList = append(languageList, selectedStyle.Render(langLine))
				} else {
					languageList = append(languageList, knownStyle.Render(langLine))
				}
			} else if alreadySelected {
				langLine = fmt.Sprintf("  [✓] %s", lang)
				if i == ls.selectedIndex {
					languageList = append(languageList, selectedStyle.Render(langLine))
				} else {
					languageList = append(languageList, chosenStyle.Render(langLine))
				}
			} else {
				langLine = fmt.Sprintf("  [ ] %s", lang)
				if i == ls.selectedIndex {
					languageList = append(languageList, selectedStyle.Render(langLine))
				} else {
					languageList = append(languageList, normalStyle.Render(langLine))
				}
			}
		} else {
			// Single selection mode (original behavior)
			langLine = " " + lang
			if i == ls.selectedIndex {
				languageList = append(languageList, selectedStyle.Render(langLine))
			} else {
				languageList = append(languageList, normalStyle.Render(langLine))
			}
		}
	}

	// Create viewport if needed
	listHeight := 20
	if ls.viewport.Width == 0 {
		ls.viewport = viewport.New(40, listHeight)
		ls.viewport.Style = lipgloss.NewStyle()
	}

	ls.viewport.SetContent(strings.Join(languageList, "\n"))
	content = append(content, ls.viewport.View())
	content = append(content, "")

	// Show legend for multi-select mode
	if ls.multiSelectMode {
		content = append(content, helpStyle.Render("Legend:"))
		content = append(content, helpStyle.Render("  [✓] = Selected or Already Known"))
		content = append(content, helpStyle.Render("  [ ] = Available"))
		content = append(content, "")
	}

	// Show different help text based on mode
	if ls.deleteMode {
		content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Remove • [Esc] Cancel"))
	} else if ls.multiSelectMode {
		if ls.CanConfirm() {
			content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("[↑/↓] Navigate • [Space] Toggle • [Enter] Confirm • [Esc] Cancel"))
		} else {
			content = append(content, helpStyle.Render(fmt.Sprintf("[↑/↓] Navigate • [Space] Toggle • Select %d more language(s) • [Esc] Cancel", ls.MaxChoices-len(ls.SelectedLanguages))))
		}
	} else {
		content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		Background(lipgloss.Color("235"))

	popup := boxStyle.Render(strings.Join(content, "\n"))

	// Center on screen
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, popup)
}

// Update handles viewport updates
func (ls *LanguageSelector) Update(msg tea.Msg) {
	var cmd tea.Cmd
	ls.viewport, cmd = ls.viewport.Update(msg)
	_ = cmd
}
