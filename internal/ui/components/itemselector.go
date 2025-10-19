// internal/ui/components/itemselector.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ItemSelectorMode represents the current interaction mode
type ItemSelectorMode int

const (
	ItemModeCategory ItemSelectorMode = iota // Selecting category
	ItemModeSearch                            // Typing search query
	ItemModeList                              // Browsing filtered items
	ItemModeQuantity                          // Entering quantity
)

// ItemSelector handles item selection with fuzzy search and category filtering
type ItemSelector struct {
	visible         bool
	mode            ItemSelectorMode
	categories      []string
	selectedCat     int
	items           []models.ItemDefinition
	selectedItem    int
	searchInput     textinput.Model
	quantityInput   textinput.Model
	viewport        viewport.Model
	currentCategory string
	selectedDef     *models.ItemDefinition
	character       *models.Character // For proficiency checks
}

// NewItemSelector creates a new item selector
func NewItemSelector() *ItemSelector {
	categories := []string{
		"All Items",
		"Weapons",
		"Armor",
		"Adventuring Gear",
		"Potions",
		"Magic Items",
		"Ammunition",
	}

	searchInput := textinput.New()
	searchInput.Placeholder = "Type to search..."
	searchInput.Focus()
	searchInput.CharLimit = 50
	searchInput.Width = 40

	quantityInput := textinput.New()
	quantityInput.Placeholder = "1"
	quantityInput.CharLimit = 4
	quantityInput.Width = 10

	return &ItemSelector{
		visible:       false,
		mode:          ItemModeCategory,
		categories:    categories,
		selectedCat:   0,
		items:         []models.ItemDefinition{},
		selectedItem:  0,
		searchInput:   searchInput,
		quantityInput: quantityInput,
	}
}

// Show displays the item selector
func (is *ItemSelector) Show(char *models.Character) {
	is.visible = true
	is.mode = ItemModeCategory
	is.selectedCat = 0
	is.selectedItem = 0
	is.searchInput.SetValue("")
	is.quantityInput.SetValue("")
	is.currentCategory = ""
	is.items = []models.ItemDefinition{}
	is.selectedDef = nil
	is.character = char
}

// Hide hides the item selector
func (is *ItemSelector) Hide() {
	is.visible = false
	is.mode = ItemModeCategory
	is.selectedCat = 0
	is.selectedItem = 0
	is.searchInput.SetValue("")
	is.quantityInput.SetValue("")
	is.searchInput.Blur()
	is.quantityInput.Blur()
	is.selectedDef = nil
}

// IsVisible returns whether the selector is visible
func (is *ItemSelector) IsVisible() bool {
	return is.visible
}

// IsInQuantityMode returns whether the selector is in quantity mode
func (is *ItemSelector) IsInQuantityMode() bool {
	return is.mode == ItemModeQuantity
}

// GetSelectedItem returns the selected item definition and quantity
func (is *ItemSelector) GetSelectedItem() (*models.ItemDefinition, int) {
	if is.selectedDef == nil {
		return nil, 0
	}

	quantity := 1
	if is.quantityInput.Value() != "" {
		fmt.Sscanf(is.quantityInput.Value(), "%d", &quantity)
	}
	if quantity < 1 {
		quantity = 1
	}

	return is.selectedDef, quantity
}

// HandleKey processes keyboard input
func (is *ItemSelector) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch is.mode {
	case ItemModeCategory:
		return is.handleCategoryKeys(msg)
	case ItemModeSearch:
		return is.handleSearchKeys(msg)
	case ItemModeList:
		return is.handleListKeys(msg)
	case ItemModeQuantity:
		return is.handleQuantityKeys(msg)
	}
	return nil
}

func (is *ItemSelector) handleCategoryKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if is.selectedCat > 0 {
			is.selectedCat--
		}
	case "down", "j":
		if is.selectedCat < len(is.categories)-1 {
			is.selectedCat++
		}
	case "enter":
		// Select category and move to search mode
		is.currentCategory = is.getCategoryFilter()
		is.items = is.getFilteredItems(models.GetItemsByCategory(is.currentCategory))
		is.selectedItem = 0
		is.mode = ItemModeSearch
		is.searchInput.Focus()
		return textinput.Blink
	case "esc":
		is.Hide()
	}
	return nil
}

func (is *ItemSelector) handleSearchKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// Back to category selection
		is.mode = ItemModeCategory
		is.searchInput.Blur()
		is.searchInput.SetValue("")
		return nil
	case "enter":
		// Move to item list
		is.mode = ItemModeList
		is.searchInput.Blur()
		return nil
	case "down", "ctrl+n":
		// Quick navigation to list
		is.mode = ItemModeList
		is.searchInput.Blur()
		return nil
	default:
		// Update search input
		var cmd tea.Cmd
		is.searchInput, cmd = is.searchInput.Update(msg)

		// Update filtered items
		query := is.searchInput.Value()
		is.items = is.getFilteredItems(models.SearchItems(query, is.currentCategory))
		is.selectedItem = 0

		return cmd
	}
}

func (is *ItemSelector) handleListKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if is.selectedItem > 0 {
			is.selectedItem--
			is.viewport.LineUp(1)
		}
	case "down", "j":
		if is.selectedItem < len(is.items)-1 {
			is.selectedItem++
			is.viewport.LineDown(1)
		}
	case "enter":
		// Select item and ask for quantity
		if is.selectedItem >= 0 && is.selectedItem < len(is.items) {
			is.selectedDef = &is.items[is.selectedItem]
			is.mode = ItemModeQuantity
			is.quantityInput.SetValue("1")
			is.quantityInput.Focus()
			return textinput.Blink
		}
	case "esc":
		// Back to search
		is.mode = ItemModeSearch
		is.searchInput.Focus()
		return textinput.Blink
	case "/":
		// Quick back to search
		is.mode = ItemModeSearch
		is.searchInput.Focus()
		return textinput.Blink
	}
	return nil
}

func (is *ItemSelector) handleQuantityKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// Back to list
		is.mode = ItemModeList
		is.quantityInput.Blur()
		is.selectedDef = nil
		return nil
	case "enter":
		// Confirm selection (handled by parent)
		return nil
	default:
		var cmd tea.Cmd
		is.quantityInput, cmd = is.quantityInput.Update(msg)
		return cmd
	}
}

func (is *ItemSelector) getCategoryFilter() string {
	switch is.selectedCat {
	case 0:
		return "all"
	case 1:
		return "weapon"
	case 2:
		return "armor"
	case 3:
		return "gear"
	case 4:
		return "potion"
	case 5:
		return "magic"
	case 6:
		return "ammunition"
	default:
		return "all"
	}
}

// getFilteredItems - No filtering needed, players can buy any item
// Proficiency checks happen only when EQUIPPING items
func (is *ItemSelector) getFilteredItems(items []models.ItemDefinition) []models.ItemDefinition {
	// Return all items - let players buy anything they want
	return items
}

// View renders the item selector
func (is *ItemSelector) View(width, height int) string {
	if !is.visible {
		return ""
	}

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

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(width - 4)

	var content strings.Builder

	switch is.mode {
	case ItemModeCategory:
		content.WriteString(titleStyle.Render("SELECT ITEM CATEGORY"))
		content.WriteString("\n\n")

		for i, cat := range is.categories {
			line := "  " + cat
			if i == is.selectedCat {
				content.WriteString(selectedStyle.Render("▶ " + cat))
			} else {
				content.WriteString(normalStyle.Render(line))
			}
			content.WriteString("\n")
		}

		content.WriteString("\n")
		content.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • ESC: Cancel"))

	case ItemModeSearch, ItemModeList:
		// Two-column layout: item list on left, details on right
		content.WriteString(titleStyle.Render(fmt.Sprintf("ADD ITEM - %s", strings.ToUpper(is.currentCategory))))
		content.WriteString("\n\n")

		// Calculate column widths
		leftWidth := (width - 8) / 2
		rightWidth := (width - 8) - leftWidth - 2

		// Left column: Item list
		var leftCol strings.Builder
		leftCol.WriteString(labelStyle.Render("Search: "))
		leftCol.WriteString(is.searchInput.View())
		leftCol.WriteString(fmt.Sprintf(" (%d)", len(is.items)))
		leftCol.WriteString("\n\n")

		if len(is.items) == 0 {
			leftCol.WriteString(helpStyle.Render("No items found"))
		} else {
			maxItems := (height - 12) // Limit displayed items
			startIdx := 0
			if is.selectedItem > maxItems/2 && len(is.items) > maxItems {
				startIdx = is.selectedItem - maxItems/2
				if startIdx+maxItems > len(is.items) {
					startIdx = len(is.items) - maxItems
				}
			}
			endIdx := startIdx + maxItems
			if endIdx > len(is.items) {
				endIdx = len(is.items)
			}

			for i := startIdx; i < endIdx; i++ {
				item := is.items[i]
				line := fmt.Sprintf("%-20s %5.0f gp",
					truncateItemName(item.Name, 20),
					item.PriceGP,
				)
				if i == is.selectedItem {
					leftCol.WriteString(selectedStyle.Render("▶ " + line))
				} else {
					leftCol.WriteString(normalStyle.Render("  " + line))
				}
				leftCol.WriteString("\n")
			}
		}

		// Right column: Item details
		var rightCol strings.Builder
		if is.selectedItem >= 0 && is.selectedItem < len(is.items) {
			selectedItem := is.items[is.selectedItem]

			rightCol.WriteString(labelStyle.Render("═══ ITEM DETAILS ═══"))
			rightCol.WriteString("\n\n")

			rightCol.WriteString(labelStyle.Render(selectedItem.Name))
			rightCol.WriteString("\n\n")

			rightCol.WriteString(labelStyle.Render("Category: "))
			rightCol.WriteString(normalStyle.Render(selectedItem.Category))
			rightCol.WriteString("\n")

			if selectedItem.Subcategory != "" {
				rightCol.WriteString(labelStyle.Render("Type: "))
				rightCol.WriteString(normalStyle.Render(selectedItem.Subcategory))
				rightCol.WriteString("\n")
			}

			rightCol.WriteString(labelStyle.Render("Price: "))
			rightCol.WriteString(normalStyle.Render(fmt.Sprintf("%.2f gp", selectedItem.PriceGP)))
			rightCol.WriteString("\n")

			rightCol.WriteString(labelStyle.Render("Weight: "))
			rightCol.WriteString(normalStyle.Render(fmt.Sprintf("%.1f lbs", selectedItem.Weight)))
			rightCol.WriteString("\n")

			// Weapon/Armor specific info
			if selectedItem.Damage != "" {
				rightCol.WriteString(labelStyle.Render("Damage: "))
				rightCol.WriteString(normalStyle.Render(fmt.Sprintf("%s %s", selectedItem.Damage, selectedItem.DamageType)))
				rightCol.WriteString("\n")
			}

			if selectedItem.AC != "" {
				rightCol.WriteString(labelStyle.Render("AC: "))
				rightCol.WriteString(normalStyle.Render(selectedItem.AC))
				rightCol.WriteString("\n")
			}

			if len(selectedItem.Properties) > 0 {
				rightCol.WriteString(labelStyle.Render("Properties: "))
				rightCol.WriteString(normalStyle.Render(strings.Join(selectedItem.Properties, ", ")))
				rightCol.WriteString("\n")
			}

			// Description
			if selectedItem.Description != "" {
				rightCol.WriteString("\n")
				rightCol.WriteString(labelStyle.Render("Description:"))
				rightCol.WriteString("\n")
				wrapped := wrapItemText(selectedItem.Description, rightWidth-2)
				rightCol.WriteString(normalStyle.Render(wrapped))
			}
		}

		// Combine columns
		leftColStr := lipgloss.NewStyle().
			Width(leftWidth).
			Render(leftCol.String())

		rightColStr := lipgloss.NewStyle().
			Width(rightWidth).
			Render(rightCol.String())

		content.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftColStr, "  ", rightColStr))
		content.WriteString("\n\n")

		if is.mode == ItemModeSearch {
			content.WriteString(helpStyle.Render("Type to search • ↓: Navigate • Enter: Select • ESC: Back"))
		} else {
			content.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • /: Search • ESC: Back"))
		}

	case ItemModeQuantity:
		content.WriteString(titleStyle.Render("SELECT QUANTITY"))
		content.WriteString("\n\n")

		if is.selectedDef != nil {
			content.WriteString(labelStyle.Render("Item: "))
			content.WriteString(normalStyle.Render(is.selectedDef.Name))
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Price: "))
			content.WriteString(normalStyle.Render(fmt.Sprintf("%.2f gp each", is.selectedDef.PriceGP)))
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Weight: "))
			content.WriteString(normalStyle.Render(fmt.Sprintf("%.1f lbs each", is.selectedDef.Weight)))
			content.WriteString("\n\n")

			content.WriteString(labelStyle.Render("Quantity: "))
			content.WriteString(is.quantityInput.View())
			content.WriteString("\n\n")

			// Calculate total
			quantity := 1
			if is.quantityInput.Value() != "" {
				fmt.Sscanf(is.quantityInput.Value(), "%d", &quantity)
			}
			if quantity < 1 {
				quantity = 1
			}
			totalPrice := float64(quantity) * is.selectedDef.PriceGP
			totalWeight := float64(quantity) * is.selectedDef.Weight
			content.WriteString(normalStyle.Render(fmt.Sprintf("Total: %.2f gp, %.1f lbs", totalPrice, totalWeight)))
		}

		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("Enter: Add to Inventory • ESC: Cancel"))
	}

	return boxStyle.Render(content.String())
}

// truncateItemName truncates an item name to maxLen
func truncateItemName(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// wrapItemText wraps text to fit within width
func wrapItemText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// Update handles viewport updates
func (is *ItemSelector) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	is.viewport, cmd = is.viewport.Update(msg)
	return cmd
}
