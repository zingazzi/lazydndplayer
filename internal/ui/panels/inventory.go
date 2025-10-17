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
		p.viewport = viewport.New(width, height)
		p.ready = true
	}
	p.viewport.Width = width
	p.viewport.Height = height

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

	// Items in two columns
	if len(char.Inventory.Items) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No items in inventory"))
	} else {
		numItems := len(char.Inventory.Items)
		midpoint := (numItems + 1) / 2

		columnWidth := width / 2
		if columnWidth > 40 {
			columnWidth = 40
		}

		// Build each row with two columns
		for i := 0; i < midpoint; i++ {
			leftItem := char.Inventory.Items[i]
			leftMarker := " "
			if leftItem.Equipped {
				leftMarker = "E"
			}

			leftLine := fmt.Sprintf("[%s] %-18s x%-2d %.1flbs",
				leftMarker,
				leftItem.Name,
				leftItem.Quantity,
				leftItem.TotalWeight(),
			)

			// Style left column
			var leftStyled string
			if i == p.selectedIndex {
				leftStyled = selectedStyle.Width(columnWidth).Render(leftLine)
			} else {
				leftStyled = normalStyle.Width(columnWidth).Render(leftLine)
			}

			// Right column (if exists)
			rightStyled := ""
			rightIdx := i + midpoint
			if rightIdx < numItems {
				rightItem := char.Inventory.Items[rightIdx]
				rightMarker := " "
				if rightItem.Equipped {
					rightMarker = "E"
				}

				rightLine := fmt.Sprintf("[%s] %-18s x%-2d %.1flbs",
					rightMarker,
					rightItem.Name,
					rightItem.Quantity,
					rightItem.TotalWeight(),
				)

				if rightIdx == p.selectedIndex {
					rightStyled = selectedStyle.Width(columnWidth).Render(rightLine)
				} else {
					rightStyled = normalStyle.Width(columnWidth).Render(rightLine)
				}
			}

			// Join columns
			row := lipgloss.JoinHorizontal(lipgloss.Left, leftStyled, rightStyled)
			lines = append(lines, row)
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("[E] = Equipped  |  Press 'a' to add item  |  Press 'e' to toggle equipped"))

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
