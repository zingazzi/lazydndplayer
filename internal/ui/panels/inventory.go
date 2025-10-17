// internal/ui/panels/inventory.go
package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// InventoryPanel displays character inventory
type InventoryPanel struct {
	character     *models.Character
	selectedIndex int
	viewport      viewport.Model
	ready         bool
}

// NewInventoryPanel creates a new inventory panel
func NewInventoryPanel(char *models.Character) *InventoryPanel {
	return &InventoryPanel{
		character:     char,
		selectedIndex: 0,
	}
}

// View renders the inventory panel
func (p *InventoryPanel) View(width, height int) string {
	char := p.character

	// Initialize viewport if not ready
	if !p.ready {
		p.viewport = viewport.New(width-4, height-2)
		p.ready = true
	}
	p.viewport.Width = width - 4
	p.viewport.Height = height - 2

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	var lines []string
	lines = append(lines, titleStyle.Render("INVENTORY"))
	lines = append(lines, "")

	// Currency
	currencyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	lines = append(lines, currencyStyle.Render(fmt.Sprintf("Gold: %d  Silver: %d  Copper: %d",
		char.Inventory.Gold, char.Inventory.Silver, char.Inventory.Copper)))
	lines = append(lines, "")

	// Encumbrance
	totalWeight := char.Inventory.TotalWeight()
	weightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	if char.Inventory.IsOverloaded() {
		weightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	}

	lines = append(lines, weightStyle.Render(fmt.Sprintf("Carry Weight: %.1f / %.1f lbs",
		totalWeight, char.Inventory.CarryCapacity)))
	lines = append(lines, "")

	// Items
	if len(char.Inventory.Items) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No items in inventory"))
	} else {
		for i, item := range char.Inventory.Items {
			equippedMarker := " "
			if item.Equipped {
				equippedMarker = "E"
			}

			line := fmt.Sprintf("[%s] %-25s x%-3d  %.1f lbs",
				equippedMarker,
				item.Name,
				item.Quantity,
				item.TotalWeight(),
			)

			if i == p.selectedIndex {
				lines = append(lines, selectedStyle.Render(line))
			} else {
				lines = append(lines, normalStyle.Render(line))
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("[E] = Equipped"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'a' to add item"))
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press 'e' to toggle equipped"))

	content := strings.Join(lines, "\n")
	p.viewport.SetContent(content)

	// Add scroll indicators
	scrollInfo := ""
	if p.viewport.ScrollPercent() < 1.0 {
		scrollInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf(" ↓ %d%%", int(p.viewport.ScrollPercent()*100)))
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(p.viewport.View() + scrollInfo)
}

// Update handles updates for the inventory panel
func (p *InventoryPanel) Update(char *models.Character) {
	p.character = char
}

// Next moves to next item
func (p *InventoryPanel) Next() {
	if len(p.character.Inventory.Items) > 0 && p.selectedIndex < len(p.character.Inventory.Items)-1 {
		p.selectedIndex++
		p.viewport.LineDown(1)
	}
}

// Prev moves to previous item
func (p *InventoryPanel) Prev() {
	if p.selectedIndex > 0 {
		p.selectedIndex--
		p.viewport.LineUp(1)
	}
}

// GetSelectedItem returns the currently selected item
func (p *InventoryPanel) GetSelectedItem() *models.Item {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Inventory.Items) {
		return &p.character.Inventory.Items[p.selectedIndex]
	}
	return nil
}

// ToggleEquipped toggles equipped status for selected item
func (p *InventoryPanel) ToggleEquipped() {
	if item := p.GetSelectedItem(); item != nil {
		item.Equipped = !item.Equipped
	}
}

// DeleteSelected deletes the selected item
func (p *InventoryPanel) DeleteSelected() {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.character.Inventory.Items) {
		p.character.Inventory.Items = append(
			p.character.Inventory.Items[:p.selectedIndex],
			p.character.Inventory.Items[p.selectedIndex+1:]...,
		)
		if p.selectedIndex >= len(p.character.Inventory.Items) && p.selectedIndex > 0 {
			p.selectedIndex--
		}
	}
}
