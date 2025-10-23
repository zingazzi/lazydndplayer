# Class Selector Multiclass Logic Bug Fix

## Issue
Pressing 'c' in Character Info panel showed a black screen with no classes to select for a NEW character.

## Root Cause
**File**: `internal/ui/components/classselector.go`, Line 38 (function `Show()`)

The multiclass detection logic was **incorrect**:

```go
// WRONG:
c.isMulticlass = c.character.TotalLevel > 0
```

This treated a **new character with TotalLevel > 0 but NO class** as if they were multiclassing!

### Why This Failed:
1. New character has `TotalLevel = 1` (or more) but `Class = ""` (empty)
2. Code incorrectly thought this was multiclassing
3. Called `models.GetAvailableClasses(c.character)`
4. This function checks multiclass **prerequisites** (ability score 13+ in specific abilities)
5. New character doesn't meet prerequisites → returns **0 classes**
6. Selector had no classes to show → empty array
7. `GetSelectedClass()` always returned `""` → black screen

### Debug Log Evidence:
```
[11:34:18.646] Backed up class state:                    ← Empty class
[11:34:18.697] Successfully loaded 12 classes            ← GetAllClasses() works
[11:34:18.723] ClassSelector.Show() - Multiclass mode, loaded 0 classes  ← BUG!
[11:34:18.743] WARNING: ClassSelector.Show() - No classes loaded! isMulticlass=true
[11:34:19.595] ClassSelector.Next() - Already at end (index: 0, total: 0)  ← Can't navigate
```

## Fix Applied

```go
// CORRECT - Line 39:
c.isMulticlass = c.character.TotalLevel > 0 && c.character.Class != ""
```

Now properly checks **BOTH conditions**:
1. ✅ Character has levels (TotalLevel > 0)
2. ✅ Character already has a class assigned (Class != "")

### Logic Table:

| TotalLevel | Class     | isMulticlass | Loads            |
|------------|-----------|--------------|------------------|
| 0          | ""        | false        | All 12 classes   |
| 1          | ""        | false ✅     | All 12 classes   |
| 1          | "Fighter" | true         | Prerequisite met |
| 5          | "Fighter" | true         | Prerequisite met |

## Files Modified

- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/classselector.go`
  - Line 39: Added `&& c.character.Class != ""` to multiclass check
  - Line 40: Added debug log showing TotalLevel, Class, and isMulticlass

## Testing

### Build:
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go build -o lazydndplayer .
```

### Test (New Character):
```bash
rm -f lazydnd_debug.log
./lazydndplayer --debug
```

1. Press 'p' to focus Character Info panel
2. Press 'c' to open class selector
3. **Expected**: Popup shows all 12 classes (Barbarian, Bard, Cleric, etc.)
4. **Expected**: Can navigate with ↑/↓
5. **Expected**: Can select with Enter

### Expected Log:
```
[time] ClassSelector.Show() - TotalLevel=1, Class='', isMulticlass=false
[time] ClassSelector.Show() - First class mode, loaded 12 classes
[time] ClassSelector - First class: Barbarian, Last class: Wizard
```

### Test (Existing Character with Class):
If you already have a Fighter and try to change class:
```
[time] ClassSelector.Show() - TotalLevel=1, Class='Fighter', isMulticlass=true
[time] ClassSelector.Show() - Multiclass mode, loaded X classes
```

## Related Fixes
This completes the class selector bug fixes:
- ✅ `CRITICAL_BUG_FIX.md` - Fixed View() return statements
- ✅ `MISSING_RETURN_FIX.md` - Fixed handler return statement
- ✅ `SELECTOR_UPDATE_FIX.md` - Added Update() methods to selectors
- ✅ **This fix** - Fixed multiclass detection logic

## Status
✅ **FIXED** - Multiclass detection now correctly checks for existing class

The black screen issue for new characters should now be completely resolved!
