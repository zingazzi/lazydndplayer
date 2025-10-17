// internal/ui/components/languageselector.go
package components

import (
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
	languages     []string
	selectedIndex int
	viewport      viewport.Model
	visible       bool
}

// NewLanguageSelector creates a new language selector
func NewLanguageSelector() *LanguageSelector {
	return &LanguageSelector{
		languages:     dndLanguages,
		selectedIndex: 0,
		visible:       false,
	}
}

// Show displays the language selector
func (ls *LanguageSelector) Show() {
	ls.visible = true
	ls.selectedIndex = 0
}

// Hide hides the language selector
func (ls *LanguageSelector) Hide() {
	ls.visible = false
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
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("237"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Build content
	var content []string
	content = append(content, titleStyle.Render("SELECT ADDITIONAL LANGUAGE"))
	content = append(content, "")

	// Language list
	var languageList []string
	for i, lang := range ls.languages {
		langLine := " " + lang
		if i == ls.selectedIndex {
			languageList = append(languageList, selectedStyle.Render(langLine))
		} else {
			languageList = append(languageList, normalStyle.Render(langLine))
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
	content = append(content, helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [Esc] Cancel"))

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
