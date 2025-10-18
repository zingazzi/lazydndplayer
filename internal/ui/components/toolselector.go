// internal/ui/components/toolselector.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Common D&D 5e tools
var dndTools = []string{
	"Alchemist's Supplies",
	"Bagpipes",
	"Brewer's Supplies",
	"Calligrapher's Supplies",
	"Carpenter's Tools",
	"Cartographer's Tools",
	"Cobbler's Tools",
	"Cook's Utensils",
	"Disguise Kit",
	"Drum",
	"Dulcimer",
	"Flute",
	"Forgery Kit",
	"Gaming Set (Dice)",
	"Gaming Set (Playing Cards)",
	"Gaming Set (Chess)",
	"Glassblower's Tools",
	"Herbalism Kit",
	"Horn",
	"Jeweler's Tools",
	"Land Vehicles",
	"Leatherworker's Tools",
	"Lute",
	"Lyre",
	"Mason's Tools",
	"Navigator's Tools",
	"Painter's Supplies",
	"Pan Flute",
	"Poisoner's Kit",
	"Potter's Tools",
	"Shawm",
	"Smith's Tools",
	"Thieves' Tools",
	"Tinker's Tools",
	"Viol",
	"Water Vehicles",
	"Weaver's Tools",
	"Woodcarver's Tools",
}

// ToolSelector handles tool selection UI
type ToolSelector struct {
	allTools      []string
	tools         []string // Filtered list (excluding already known OR showing only known)
	selectedIndex int
	viewport      viewport.Model
	visible       bool
	deleteMode    bool   // If true, shows known tools for deletion
	title         string // Custom title
}

// NewToolSelector creates a new tool selector
func NewToolSelector() *ToolSelector {
	return &ToolSelector{
		allTools:      dndTools,
		tools:         dndTools,
		selectedIndex: 0,
		visible:       false,
		deleteMode:    false,
		title:         "SELECT TOOL PROFICIENCY",
	}
}

// SetExcludeTools filters out tools the character already has
func (ts *ToolSelector) SetExcludeTools(knownTools []string) {
	// Create a set of known tools for fast lookup
	knownSet := make(map[string]bool)
	for _, tool := range knownTools {
		// Normalize to handle variations
		normalizedTool := strings.ToLower(strings.TrimSpace(tool))
		// Skip placeholder texts
		if !strings.Contains(normalizedTool, "choose") &&
		   !strings.Contains(normalizedTool, "choice") {
			knownSet[normalizedTool] = true
		}
	}

	// Filter out known tools
	ts.tools = []string{}
	for _, tool := range ts.allTools {
		if !knownSet[strings.ToLower(tool)] {
			ts.tools = append(ts.tools, tool)
		}
	}

	// Reset selected index
	ts.selectedIndex = 0
}

// Show displays the tool selector (for adding tools)
func (ts *ToolSelector) Show() {
	ts.visible = true
	ts.selectedIndex = 0
	ts.deleteMode = false
	ts.title = "SELECT TOOL PROFICIENCY"
}

// ShowForDeletion displays the tool selector with known tools (for deleting)
func (ts *ToolSelector) ShowForDeletion(knownTools []string) {
	ts.visible = true
	ts.selectedIndex = 0
	ts.deleteMode = true
	ts.title = "SELECT TOOL TO REMOVE"

	// Filter out placeholder texts and set tools to known only
	ts.tools = []string{}
	for _, tool := range knownTools {
		normalizedTool := strings.ToLower(strings.TrimSpace(tool))
		// Skip placeholder texts
		if !strings.Contains(normalizedTool, "choose") &&
		   !strings.Contains(normalizedTool, "choice") {
			ts.tools = append(ts.tools, tool)
		}
	}

	ts.selectedIndex = 0
}

// Hide hides the tool selector
func (ts *ToolSelector) Hide() {
	ts.visible = false
	ts.deleteMode = false
}

// IsVisible returns whether the selector is visible
func (ts *ToolSelector) IsVisible() bool {
	return ts.visible
}

// IsDeleteMode returns whether we're in delete mode
func (ts *ToolSelector) IsDeleteMode() bool {
	return ts.deleteMode
}

// Next moves to the next tool
func (ts *ToolSelector) Next() {
	if ts.selectedIndex < len(ts.tools)-1 {
		ts.selectedIndex++
	}
}

// Prev moves to the previous tool
func (ts *ToolSelector) Prev() {
	if ts.selectedIndex > 0 {
		ts.selectedIndex--
	}
}

// GetSelected returns the currently selected tool
func (ts *ToolSelector) GetSelected() string {
	if len(ts.tools) == 0 {
		return ""
	}
	if ts.selectedIndex >= len(ts.tools) {
		ts.selectedIndex = len(ts.tools) - 1
	}
	return ts.tools[ts.selectedIndex]
}

// View renders the tool selector
func (ts *ToolSelector) View(width, height int) string {
	if !ts.visible {
		return ""
	}

	// Initialize viewport if needed
	if ts.viewport.Width == 0 {
		ts.viewport = viewport.New(width, height-6) // Leave space for title and borders
		ts.viewport.Style = lipgloss.NewStyle()
	}

	// Update viewport dimensions
	ts.viewport.Width = width
	ts.viewport.Height = height - 6

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Align(lipgloss.Center).
		Width(width - 4)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Background(lipgloss.Color("237")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(width - 4)

	// Build content
	var content strings.Builder
	
	content.WriteString(titleStyle.Render(ts.title))
	content.WriteString("\n\n")

	if len(ts.tools) == 0 {
		if ts.deleteMode {
			content.WriteString(emptyStyle.Render("No tool proficiencies to remove"))
		} else {
			content.WriteString(emptyStyle.Render("All tools already known"))
		}
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Press ESC to close"))
	} else {
		// Calculate scroll position to keep selected item visible
		visibleLines := ts.viewport.Height
		startIdx := 0
		endIdx := len(ts.tools)

		if len(ts.tools) > visibleLines {
			// Center the selected item
			startIdx = ts.selectedIndex - visibleLines/2
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx = startIdx + visibleLines
			if endIdx > len(ts.tools) {
				endIdx = len(ts.tools)
				startIdx = endIdx - visibleLines
				if startIdx < 0 {
					startIdx = 0
				}
			}
		}

		var toolList strings.Builder
		for i := startIdx; i < endIdx; i++ {
			tool := ts.tools[i]
			if i == ts.selectedIndex {
				toolList.WriteString(selectedStyle.Render("▶ " + tool))
			} else {
				toolList.WriteString(normalStyle.Render("  " + tool))
			}
			toolList.WriteString("\n")
		}

		ts.viewport.SetContent(toolList.String())
		content.WriteString(ts.viewport.View())
		content.WriteString("\n\n")
		
		if ts.deleteMode {
			content.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Remove • ESC: Cancel"))
		} else {
			content.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • ESC: Cancel"))
		}
	}

	return boxStyle.Render(content.String())
}

// Update handles tea.Msg for the viewport
func (ts *ToolSelector) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	ts.viewport, cmd = ts.viewport.Update(msg)
	return cmd
}

