// internal/ui/components/delevelselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/debug"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ClassLevelInfo represents a class that can be de-leveled
type ClassLevelInfo struct {
	ClassName string
	Level     int
	Subclass  string
	CanRemove bool
}

// DeLevelSelector handles the de-leveling process
type DeLevelSelector struct {
	visible       bool
	character     *models.Character
	classes       []ClassLevelInfo
	selectedIndex int
	showConfirm   bool
	confirmTarget string
	previewResult *models.DeLevelResult
}

// NewDeLevelSelector creates a new de-level selector
func NewDeLevelSelector(char *models.Character) *DeLevelSelector {
	return &DeLevelSelector{
		character: char,
	}
}

// Show displays the de-level selector
func (ds *DeLevelSelector) Show() {
	ds.visible = true
	ds.selectedIndex = 0
	ds.showConfirm = false
	ds.confirmTarget = ""
	ds.previewResult = nil
	ds.loadClasses()
	debug.Log("DeLevelSelector.Show(): Loaded %d classes", len(ds.classes))
}

// Hide hides the de-level selector
func (ds *DeLevelSelector) Hide() {
	ds.visible = false
	ds.showConfirm = false
	ds.confirmTarget = ""
	ds.previewResult = nil
}

// IsVisible returns whether the selector is visible
func (ds *DeLevelSelector) IsVisible() bool {
	return ds.visible
}

// loadClasses loads all classes with their current levels
func (ds *DeLevelSelector) loadClasses() {
	ds.classes = []ClassLevelInfo{}

	for _, cl := range ds.character.Classes {
		info := ClassLevelInfo{
			ClassName: cl.ClassName,
			Level:     cl.Level,
			Subclass:  cl.Subclass,
			CanRemove: cl.Level > 0,
		}
		ds.classes = append(ds.classes, info)
	}

	debug.Log("DeLevelSelector.loadClasses(): Loaded %d classes", len(ds.classes))
}

// Next moves to the next class
func (ds *DeLevelSelector) Next() {
	if ds.selectedIndex < len(ds.classes)-1 {
		ds.selectedIndex++
	}
}

// Prev moves to the previous class
func (ds *DeLevelSelector) Prev() {
	if ds.selectedIndex > 0 {
		ds.selectedIndex--
	}
}

// GetSelectedClass returns the currently selected class name
func (ds *DeLevelSelector) GetSelectedClass() string {
	if ds.selectedIndex >= 0 && ds.selectedIndex < len(ds.classes) {
		return ds.classes[ds.selectedIndex].ClassName
	}
	return ""
}

// ShowConfirmation generates a preview of what will be removed
func (ds *DeLevelSelector) ShowConfirmation() error {
	className := ds.GetSelectedClass()
	if className == "" {
		return fmt.Errorf("no class selected")
	}

	debug.Log("DeLevelSelector.ShowConfirmation(): Previewing de-level for %s", className)

	// Generate preview (don't actually de-level yet)
	preview, err := models.DeLevel(ds.character, className)
	if err != nil {
		return err
	}

	ds.showConfirm = true
	ds.confirmTarget = className
	ds.previewResult = preview

	return nil
}

// ConfirmDeLevel executes the de-level (already done in ShowConfirmation, so just hide)
func (ds *DeLevelSelector) ConfirmDeLevel() {
	debug.Log("DeLevelSelector.ConfirmDeLevel(): De-level confirmed")
	ds.Hide()
}

// CancelConfirmation cancels the confirmation dialog
func (ds *DeLevelSelector) CancelConfirmation() {
	ds.showConfirm = false
	ds.confirmTarget = ""
	ds.previewResult = nil
}

// IsShowingConfirmation returns whether the confirmation dialog is showing
func (ds *DeLevelSelector) IsShowingConfirmation() bool {
	return ds.showConfirm
}

// GetPreviewResult returns the preview result
func (ds *DeLevelSelector) GetPreviewResult() *models.DeLevelResult {
	return ds.previewResult
}

// View renders the de-level selector
func (ds *DeLevelSelector) View(width, height int) string {
	if !ds.visible {
		return ""
	}

	if ds.showConfirm {
		return ds.renderConfirmation(width, height)
	}

	return ds.renderClassList(width, height)
}

// renderClassList renders the class selection list
func (ds *DeLevelSelector) renderClassList(width, height int) string {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	var content strings.Builder
	content.WriteString(titleStyle.Render("REMOVE CLASS LEVEL") + "\n")
	content.WriteString(warningStyle.Render("⚠ WARNING: This will remove features and HP!") + "\n\n")

	if len(ds.classes) == 0 {
		content.WriteString(dimStyle.Render("No classes to remove") + "\n")
	} else {
		content.WriteString(normalStyle.Render("Select class to remove one level:") + "\n\n")

		for i, classInfo := range ds.classes {
			cursor := "  "
			style := normalStyle
			if i == ds.selectedIndex {
				cursor = "❯ "
				style = selectedStyle
			}

			subclassInfo := ""
			if classInfo.Subclass != "" {
				subclassInfo = fmt.Sprintf(" (%s)", classInfo.Subclass)
			}

			line := fmt.Sprintf("%s%s %d%s", cursor, classInfo.ClassName, classInfo.Level, subclassInfo)
			content.WriteString(style.Render(line) + "\n")
		}
	}

	content.WriteString("\n" + dimStyle.Render("↑/↓: Navigate • Enter: Confirm • Esc: Cancel"))

	// Popup style
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")). // Red for warning
		Padding(1, 2).
		Width(60)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}

// renderConfirmation renders the confirmation dialog
func (ds *DeLevelSelector) renderConfirmation(width, height int) string {
	if ds.previewResult == nil {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true)

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var content strings.Builder
	content.WriteString(titleStyle.Render("CONFIRM DE-LEVEL") + "\n\n")

	// Class change
	if ds.previewResult.ClassRemoved {
		content.WriteString(labelStyle.Render("Class: ") +
			valueStyle.Render(fmt.Sprintf("%s %d → REMOVED",
				ds.previewResult.ClassName, ds.previewResult.OldClassLevel)) + "\n")
	} else {
		content.WriteString(labelStyle.Render("Class: ") +
			valueStyle.Render(fmt.Sprintf("%s %d → %d",
				ds.previewResult.ClassName, ds.previewResult.OldClassLevel, ds.previewResult.NewClassLevel)) + "\n")
	}

	// Total level change
	content.WriteString(labelStyle.Render("Total Level: ") +
		valueStyle.Render(fmt.Sprintf("%d", ds.previewResult.NewTotalLevel)) + "\n\n")

	// HP loss
	content.WriteString(warningStyle.Render(fmt.Sprintf("HP Lost: -%d", ds.previewResult.HPLost)) + "\n\n")

	// Features removed
	if len(ds.previewResult.FeaturesRemoved) > 0 {
		content.WriteString(titleStyle.Render("FEATURES REMOVED:") + "\n")
		for _, feature := range ds.previewResult.FeaturesRemoved {
			content.WriteString(warningStyle.Render("  ✗ " + feature) + "\n")
		}
	} else {
		content.WriteString(dimStyle.Render("No features to remove") + "\n")
	}

	content.WriteString("\n" + dimStyle.Render("Enter: Confirm • Esc: Cancel"))

	// Popup style
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")). // Red for warning
		Padding(1, 2).
		Width(60)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popupStyle.Render(content.String()),
	)
}
