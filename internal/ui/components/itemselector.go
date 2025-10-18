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
func (is *ItemSelector) Show() {
	is.visible = true
	is.mode = ItemModeCategory
	is.selectedCat = 0
	is.selectedItem = 0
	is.searchInput.SetValue("")
	is.quantityInput.SetValue("")
	is.currentCategory = ""
	is.items = []models.ItemDefinition{}
	is.selectedDef = nil
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
		is.items = models.GetItemsByCategory(is.currentCategory)
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
		is.items = models.SearchItems(query, is.currentCategory)
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

// View renders the item selector
func (is *ItemSelector) View(width, height int) string {
	if !is.visible {
		return ""
	}

	// Initialize viewport if needed
	if is.viewport.Width == 0 {
		is.viewport = viewport.New(width-10, height-15)
		is.viewport.Style = lipgloss.NewStyle()
	}
	is.viewport.Width = width - 10
	is.viewport.Height = height - 15

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Align(lipgloss.Center).
		Width(width - 10)

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
		content.WriteString(titleStyle.Render(fmt.Sprintf("ADD ITEM - %s", strings.ToUpper(is.currentCategory))))
		content.WriteString("\n\n")

		// Search box
		content.WriteString(labelStyle.Render("Search: "))
		content.WriteString(is.searchInput.View())
		content.WriteString(fmt.Sprintf(" (%d items)", len(is.items)))
		content.WriteString("\n\n")

		// Item list
		if len(is.items) == 0 {
			content.WriteString(helpStyle.Render("No items found"))
		} else {
			var itemList strings.Builder
			for i, item := range is.items {
				line := fmt.Sprintf("%-30s %6.2f gp  %.1f lbs",
					item.Name,
					item.PriceGP,
					item.Weight,
				)
				if i == is.selectedItem {
					itemList.WriteString(selectedStyle.Render("▶ " + line))
				} else {
					itemList.WriteString(normalStyle.Render("  " + line))
				}
				itemList.WriteString("\n")
			}
			is.viewport.SetContent(itemList.String())
			content.WriteString(is.viewport.View())
		}

		content.WriteString("\n\n")
		if is.mode == ItemModeSearch {
			content.WriteString(helpStyle.Render("Type to search • ↓: Navigate • Enter: Confirm • ESC: Back"))
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

// Update handles viewport updates
func (is *ItemSelector) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	is.viewport, cmd = is.viewport.Update(msg)
	return cmd
}
