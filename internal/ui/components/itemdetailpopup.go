// internal/ui/components/itemdetailpopup.go
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marcozingoni/lazydndplayer/internal/models"
)

// ItemDetailPopup displays detailed information about an item
type ItemDetailPopup struct {
	visible bool
	item    *models.Item
	def     *models.ItemDefinition
}

// NewItemDetailPopup creates a new item detail popup
func NewItemDetailPopup() *ItemDetailPopup {
	return &ItemDetailPopup{
		visible: false,
	}
}

// Show displays the popup with item details
func (p *ItemDetailPopup) Show(item *models.Item) {
	p.visible = true
	p.item = item

	// Load item definition for additional details
	if item != nil {
		p.def = models.GetItemDefinitionByName(item.Name)
	}
}

// Hide closes the popup
func (p *ItemDetailPopup) Hide() {
	p.visible = false
	p.item = nil
	p.def = nil
}

// IsVisible returns whether the popup is visible
func (p *ItemDetailPopup) IsVisible() bool {
	return p.visible
}

// View renders the popup
func (p *ItemDetailPopup) View(width, height int) string {
	if !p.visible || p.item == nil {
		return ""
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	equippedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Calculate popup dimensions
	popupWidth := width * 60 / 100
	if popupWidth < 50 {
		popupWidth = 50
	}
	if popupWidth > 80 {
		popupWidth = 80
	}

	contentWidth := popupWidth - 6

	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("ITEM DETAILS"))
	content.WriteString("\n\n")

	// Item name
	content.WriteString(labelStyle.Render(p.item.Name))
	if p.item.Equipped {
		content.WriteString(" ")
		content.WriteString(equippedStyle.Render("[EQUIPPED]"))
	}
	content.WriteString("\n\n")

	// Basic info
	content.WriteString(labelStyle.Render("Type: "))
	content.WriteString(valueStyle.Render(string(p.item.Type)))
	content.WriteString("\n")

	if p.def != nil && p.def.Subcategory != "" {
		content.WriteString(labelStyle.Render("Category: "))
		content.WriteString(valueStyle.Render(p.def.Subcategory))
		content.WriteString("\n")
	}

	content.WriteString(labelStyle.Render("Quantity: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.item.Quantity)))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Weight: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%.1f lbs (%.1f total)", p.item.Weight, p.item.TotalWeight())))
	content.WriteString("\n")

	if p.def != nil {
		content.WriteString(labelStyle.Render("Value: "))
		content.WriteString(valueStyle.Render(fmt.Sprintf("%.2f gp each", p.def.PriceGP)))
		content.WriteString("\n")
	}

	// Weapon/Armor specific details
	if p.def != nil {
		if p.def.Damage != "" {
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Damage: "))
			content.WriteString(valueStyle.Render(fmt.Sprintf("%s %s", p.def.Damage, p.def.DamageType)))
			content.WriteString("\n")
		}

		if p.def.AC != "" {
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Armor Class: "))
			content.WriteString(valueStyle.Render(p.def.AC))
			content.WriteString("\n")

			if p.def.StrengthReq > 0 {
				content.WriteString(labelStyle.Render("Strength Required: "))
				content.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.def.StrengthReq)))
				content.WriteString("\n")
			}

			if p.def.StealthDisadvantage {
				content.WriteString(valueStyle.Render("⚠ Disadvantage on Stealth checks"))
				content.WriteString("\n")
			}
		}

		if len(p.def.Properties) > 0 {
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Properties: "))
			content.WriteString(valueStyle.Render(strings.Join(p.def.Properties, ", ")))
			content.WriteString("\n")
		}

		// Description
		if p.def.Description != "" {
			content.WriteString("\n")
			content.WriteString(labelStyle.Render("Description:"))
			content.WriteString("\n")
			wrapped := wrapItemDetailText(p.def.Description, contentWidth)
			content.WriteString(valueStyle.Render(wrapped))
			content.WriteString("\n")
		}
	}

	// Help text
	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Press ESC to close"))

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1, 2).
		Width(popupWidth)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content.String()),
	)
}

// wrapItemDetailText wraps text to fit within width
func wrapItemDetailText(text string, width int) string {
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
