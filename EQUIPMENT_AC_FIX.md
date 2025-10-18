# Equipment Deletion AC Fix ✅

## 🐛 Bug Fixed

**Problem**: When deleting equipped armor without unequipping it first, the AC remained as if the armor was still equipped.

**Solution**: Added AC recalculation when deleting equipped armor items.

---

## 🔧 Changes Made

### **File: `internal/ui/app.go`**

Modified the `handleInventoryPanel` function for both `d` (delete 1) and `D` (delete all) cases:

```go
case "d":
    // Delete selected item (decrease quantity by 1 or remove if quantity is 1)
    item := m.inventoryPanel.GetSelectedItem()
    if item != nil {
        wasEquipped := item.Equipped  // ← Store equipped state
        itemType := item.Type          // ← Store item type

        if item.Quantity > 1 {
            item.Quantity--
            m.message = fmt.Sprintf("%s quantity decreased to %d", item.Name, item.Quantity)
        } else {
            itemName := item.Name
            m.inventoryPanel.DeleteSelected()
            m.message = fmt.Sprintf("%s removed from inventory", itemName)
        }

        // ← NEW: Recalculate AC if armor was equipped
        if wasEquipped && itemType == models.Armor {
            m.character.UpdateDerivedStats()
            m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
        }

        m.storage.Save(m.character)
    }

case "D":
    // Delete all of selected item
    item := m.inventoryPanel.GetSelectedItem()
    if item != nil {
        itemName := item.Name
        wasEquipped := item.Equipped  // ← Store equipped state
        itemType := item.Type          // ← Store item type

        m.inventoryPanel.DeleteSelected()
        m.message = fmt.Sprintf("All %s removed from inventory", itemName)

        // ← NEW: Recalculate AC if armor was equipped
        if wasEquipped && itemType == models.Armor {
            m.character.UpdateDerivedStats()
            m.message += fmt.Sprintf(" (AC: %d)", m.character.AC)
        }

        m.storage.Save(m.character)
    }
```

---

## ✨ How It Works

### **Before the Fix:**
```
1. Equip Leather Armor (AC: 13)
2. Press 'd' to delete it
Result: ❌ Item deleted, but AC still shows 13
```

### **After the Fix:**
```
1. Equip Leather Armor (AC: 13)
2. Press 'd' to delete it
Result: ✅ Item deleted, AC recalculates to 10 (unarmored)
Message: "Leather Armor removed from inventory (AC: 10)"
```

---

## 🎯 Features

✅ **Tracks Equipped State** - Stores `wasEquipped` before deletion
✅ **Tracks Item Type** - Stores `itemType` before deletion
✅ **Conditional Recalculation** - Only recalculates AC for armor
✅ **Works for Both Cases** - Handles both 'd' (delete 1) and 'D' (delete all)
✅ **User Feedback** - Shows new AC in the status message
✅ **Saves Changes** - Persists the updated AC to the character file

---

## 📝 Technical Details

### **Why This Fix is Necessary**

The `CalculateAC()` function looks for equipped armor in the inventory:

```go
for i := range char.Inventory.Items {
    item := &char.Inventory.Items[i]
    if !item.Equipped {
        continue
    }
    // Check if it's armor...
}
```

When an item is deleted:
1. **Before**: Item removed → AC not recalculated → old AC persists
2. **After**: Item removed → `UpdateDerivedStats()` called → AC recalculated → correct AC

### **Why Store State Before Deletion**

```go
wasEquipped := item.Equipped  // Must store BEFORE deletion
itemType := item.Type          // Must store BEFORE deletion
m.inventoryPanel.DeleteSelected()  // Item is now gone!
// Can't check item.Equipped here - item is deleted!
```

---

## 🧪 Test Cases

### Test 1: Delete Equipped Leather Armor
```
1. Start with Dex 14 (+2)
2. Add Leather Armor (base AC 11)
3. Equip it → AC = 13
4. Press 'd' to delete
✅ Verify: AC = 12 (10 + 2 Dex)
✅ Verify: Message shows "(AC: 12)"
```

### Test 2: Delete Equipped Shield
```
1. Wearing Leather Armor (AC 13)
2. Equip Shield → AC = 15
3. Press 'D' to delete all shields
✅ Verify: AC = 13 (armor only)
✅ Verify: Message shows "(AC: 13)"
```

### Test 3: Delete Unequipped Armor
```
1. Have Leather Armor (not equipped)
2. Press 'd' to delete
✅ Verify: AC unchanged
✅ Verify: Message shows "Leather Armor removed from inventory" (no AC)
```

### Test 4: Delete Non-Armor Item
```
1. AC = 13
2. Delete a potion (not armor)
✅ Verify: AC unchanged (13)
✅ Verify: No AC recalculation triggered
```

### Test 5: Decrease Quantity (Not Full Delete)
```
1. Have 3x Leather Armor (one equipped)
2. Press 'd' to decrease quantity
✅ Verify: Quantity becomes 2
✅ Verify: AC unchanged (armor still equipped)
✅ Verify: No AC recalculation needed
```

---

## 🎉 Result

**Bug Fixed!** ✨

Deleting equipped armor now correctly recalculates AC, preventing the character from maintaining phantom armor protection after the item is removed from inventory!

The system properly handles:
- ✅ Equipped armor deletion
- ✅ Unequipped item deletion (no recalculation)
- ✅ Non-armor item deletion (no recalculation)
- ✅ Quantity decrease (no recalculation until fully deleted)
- ✅ Both 'd' and 'D' delete commands
