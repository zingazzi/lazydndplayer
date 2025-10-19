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
	character         *models.Character
	selectedIndex     int
	visualToActualMap []int // Maps visual index to actual inventory index
	viewport          viewport.Model
	ready             bool
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
		p.viewport = viewport.New(width, height)
		p.ready = true
	}
	p.viewport.Width = width
	p.viewport.Height = height

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 1, 0)

	categoryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Underline(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("237"))

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	var lines []string
	lines = append(lines, titleStyle.Render("INVENTORY"))
	lines = append(lines, "")

	// Currency
	currencyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	lines = append(lines, currencyStyle.Render(fmt.Sprintf("💰 Gold: %d  Silver: %d  Copper: %d",
		char.Inventory.Gold, char.Inventory.Silver, char.Inventory.Copper)))
	lines = append(lines, "")

	// Encumbrance
	totalWeight := char.Inventory.TotalWeight()
	weightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	if char.Inventory.IsOverloaded() {
		weightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	}

	lines = append(lines, weightStyle.Render(fmt.Sprintf("⚖  Carry Weight: %.1f / %.1f lbs",
		totalWeight, char.Inventory.CarryCapacity)))
	lines = append(lines, "")

	// Group items by type
	if len(char.Inventory.Items) == 0 {
		lines = append(lines, emptyStyle.Render("No items in inventory"))
		lines = append(lines, "")
		lines = append(lines, emptyStyle.Render("Press 'a' to add items"))
		p.visualToActualMap = []int{} // Clear mapping
	} else {
		// Organize items by category
		itemsByType := make(map[models.ItemType][]int)
		categoryOrder := []models.ItemType{
			models.Weapon,
			models.Armor,
			"Ammunition",
			models.Potion,
			models.Magic,
			models.Tool,
			models.Gear,
			models.Other,
		}

		for i, item := range char.Inventory.Items {
			itemsByType[item.Type] = append(itemsByType[item.Type], i)
		}

		// Build mapping from visual index to actual inventory index
		p.visualToActualMap = []int{}
		visualIndex := 0

		// Display items by category
		for _, itemType := range categoryOrder {
			indices := itemsByType[itemType]
			if len(indices) == 0 {
				continue
			}

			// Category header
			lines = append(lines, categoryStyle.Render(fmt.Sprintf("═══ %s ═══", strings.ToUpper(string(itemType)))))

			// Display items in this category
			for _, actualIdx := range indices {
				item := char.Inventory.Items[actualIdx]

				// Check if item is equippable
				def := models.GetItemDefinitionByName(item.Name)
				isEquippable := def != nil && models.IsEquippable(*def)

				var line string
				if isEquippable {
					// Show checkbox for equippable items
					marker := " "
					if item.Equipped {
						marker = "E"
					}
					line = fmt.Sprintf("  [%s] %-25s x%-3d  %5.1f lbs",
						marker,
						truncateString(item.Name, 25),
						item.Quantity,
						item.TotalWeight(),
					)
				} else {
					// No checkbox for non-equippable items
					line = fmt.Sprintf("      %-25s x%-3d  %5.1f lbs",
						truncateString(item.Name, 25),
						item.Quantity,
						item.TotalWeight(),
					)
				}

				// Map visual index to actual inventory index
				p.visualToActualMap = append(p.visualToActualMap, actualIdx)

				if visualIndex == p.selectedIndex {
					lines = append(lines, selectedStyle.Render(line))
				} else {
					lines = append(lines, normalStyle.Render(line))
				}

				visualIndex++
			}

			lines = append(lines, "") // Space between categories
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("[E] = Equipped  |  'Enter' Details  |  'a' Add  |  'e' Equip  |  'd' Remove 1  |  'D' Remove All"))

	content := strings.Join(lines, "\n")
	p.viewport.SetContent(content)

	// Add scroll indicators at bottom of viewport
	scrollInfo := ""
	if p.viewport.ScrollPercent() < 1.0 {
		scrollInfo = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("↓ Scroll: %d%%", int(p.viewport.ScrollPercent()*100)))
	}

	return p.viewport.View() + scrollInfo
}

// truncateString truncates a string to maxLen, adding "..." if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// Update handles updates for the inventory panel
func (p *InventoryPanel) Update(char *models.Character) {
	p.character = char
}

// Next moves to next item
func (p *InventoryPanel) Next() {
	maxIndex := len(p.visualToActualMap) - 1
	if maxIndex < 0 {
		maxIndex = 0
	}
	if p.selectedIndex < maxIndex {
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
	// Use mapping to get actual inventory index
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.visualToActualMap) {
		actualIdx := p.visualToActualMap[p.selectedIndex]
		if actualIdx >= 0 && actualIdx < len(p.character.Inventory.Items) {
			return &p.character.Inventory.Items[actualIdx]
		}
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
	// Use mapping to get actual inventory index
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.visualToActualMap) {
		actualIdx := p.visualToActualMap[p.selectedIndex]
		if actualIdx >= 0 && actualIdx < len(p.character.Inventory.Items) {
			p.character.Inventory.Items = append(
				p.character.Inventory.Items[:actualIdx],
				p.character.Inventory.Items[actualIdx+1:]...,
			)
			// Adjust selection if needed
			if p.selectedIndex >= len(p.character.Inventory.Items) && p.selectedIndex > 0 {
				p.selectedIndex--
			}
		}
	}
}
