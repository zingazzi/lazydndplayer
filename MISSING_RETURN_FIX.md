# Missing Return Statement Fix - Class Selector

## Issue
Pressing 'c' in the Character Info panel causes the screen to go black.

## Root Cause
**File**: `internal/ui/app.go`, Line 1295 (function `handleCharStatsPanelKeys`)

Missing `return m, nil` statement after showing the class selector!

```go
// WRONG - Line 1290-1295:
case "c":
    // Backup current class state before opening selector
    m.pendingChanges.BackupClass(m.character)
    debug.Log("Backed up class state: %s", m.character.Class)
    m.classSelector.Show()
    m.message = "Select a class..."
    // ❌ MISSING RETURN - code continues executing!
```

Without the return statement, the code execution continued through the rest of the `handleCharStatsPanelKeys` function, which prevented the class selector from being displayed properly and caused other unintended behavior.

## Fix Applied

```go
// CORRECT - Line 1290-1296:
case "c":
    // Backup current class state before opening selector
    m.pendingChanges.BackupClass(m.character)
    debug.Log("Backed up class state: %s", m.character.Class)
    m.classSelector.Show()
    m.message = "Select a class..."
    return m, nil  // ✅ Now returns properly!
```

## Why This Happened

This is the **same type of bug** as documented in `CRITICAL_BUG_FIX.md`, but in a different location:
- CRITICAL_BUG_FIX.md: Missing return in the `View()` method
- This fix: Missing return in the `handleCharStatsPanelKeys()` method

Both bugs had the same symptom (black screen) because the selector was being shown but the code flow was incorrect.

## Files Modified

- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/app.go` (Line 1296 - added return statement)

## Testing Steps

1. Build the application:
   ```bash
   cd /Users/marcozingoni/Playgound/lazydndplayer
   go build -o lazydndplayer .
   ```

2. Run the application:
   ```bash
   ./lazydndplayer
   ```

3. Test the fix:
   - Press 'p' to cycle focus to Character Info panel (right side)
   - Press 'c' to open class selector
   - **Expected**: Class selector popup should appear with list of classes
   - **Expected**: You can navigate with ↑/↓ keys
   - **Expected**: You can select a class with Enter or cancel with Esc

## Status
✅ **FIXED** - Return statement added

## Related Bugs
This fix completes the work started in:
- `CRITICAL_BUG_FIX.md` - Fixed View() method returns
- `SELECTOR_UPDATE_FIX.md` - Standardized selector Update() methods
