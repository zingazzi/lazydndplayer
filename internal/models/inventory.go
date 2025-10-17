// internal/models/inventory.go
package models

// ItemType represents the type of item
type ItemType string

const (
	Weapon ItemType = "Weapon"
	Armor  ItemType = "Armor"
	Tool   ItemType = "Tool"
	Gear   ItemType = "Gear"
	Magic  ItemType = "Magic"
	Potion ItemType = "Potion"
	Other  ItemType = "Other"
)

// Item represents an inventory item
type Item struct {
	Name        string   `json:"name"`
	Type        ItemType `json:"type"`
	Quantity    int      `json:"quantity"`
	Weight      float64  `json:"weight"` // Weight per item in lbs
	Description string   `json:"description,omitempty"`
	Equipped    bool     `json:"equipped"`
	Value       int      `json:"value"` // Value in gold pieces
}

// TotalWeight calculates the total weight of the item stack
func (i *Item) TotalWeight() float64 {
	return i.Weight * float64(i.Quantity)
}

// Inventory represents the character's inventory
type Inventory struct {
	Items         []Item  `json:"items"`
	Gold          int     `json:"gold"`
	Silver        int     `json:"silver"`
	Copper        int     `json:"copper"`
	CarryCapacity float64 `json:"carry_capacity"` // Based on STR score
}

// AddItem adds an item to inventory or increases quantity if it exists
func (inv *Inventory) AddItem(item Item) {
	for i := range inv.Items {
		if inv.Items[i].Name == item.Name && inv.Items[i].Type == item.Type {
			inv.Items[i].Quantity += item.Quantity
			return
		}
	}
	inv.Items = append(inv.Items, item)
}

// RemoveItem removes an item from inventory or decreases quantity
func (inv *Inventory) RemoveItem(name string, quantity int) bool {
	for i := range inv.Items {
		if inv.Items[i].Name == name {
			if inv.Items[i].Quantity <= quantity {
				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			} else {
				inv.Items[i].Quantity -= quantity
			}
			return true
		}
	}
	return false
}

// TotalWeight calculates the total weight of all items
func (inv *Inventory) TotalWeight() float64 {
	total := 0.0
	for _, item := range inv.Items {
		total += item.TotalWeight()
	}
	return total
}

// IsOverloaded checks if the character is carrying too much
func (inv *Inventory) IsOverloaded() bool {
	return inv.TotalWeight() > inv.CarryCapacity
}

// CalculateCarryCapacity calculates carry capacity from strength score (STR × 15)
func CalculateCarryCapacity(strength int) float64 {
	return float64(strength * 15)
}
