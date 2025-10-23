# CRITICAL BUG FIX: Missing Return Statement

## Issue
Black screen when pressing 'c' to open class selector (and potentially other popups)

## Root Cause
**File**: `internal/ui/app.go`, Line 2841

Missing `return` statement in the View() method!

```go
// WRONG:
if m.itemSelector.IsVisible() {
    m.itemSelector.View(popupLargeWidth, popupLargeHeight)  // ❌ Missing return!
}
```

This caused the View to be calculated but not returned, so the render loop continued and showed the main view instead of the popup, making the screen appear black.

## Fix Applied
```go
// CORRECT:
if m.itemSelector.IsVisible() {
    return m.itemSelector.View(popupLargeWidth, popupLargeHeight)  // ✅ Now returns!
}
```

## Impact
This bug affected ALL popups rendered after the item selector in the priority chain:
- Fighting style selector
- Cantrip selector
- Spell prep selector  
- Slot restorer
- Class skill selector
- **Class selector** ← This is why 'c' showed black screen
- Species selector

## Test
1. Build and run the application
2. Press 'c' in Character Info panel
3. Class selector should now appear correctly
4. All other affected popups should also work

## Status
✅ **FIXED** - Build successful


