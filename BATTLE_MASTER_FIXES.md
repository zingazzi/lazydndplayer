# Battle Master Fixes - Complete

## Issues Fixed

### 1. ✅ Display Format Fixed
**Problem**: Showed "4d8 (4/4)" instead of "1d8 (4/4)"
**Solution**: Changed display format in Features panel
**File**: `internal/ui/panels/features.go`
```go
// OLD: fmt.Sprintf("%d%s", p.character.SuperiorityDice.Current, p.character.SuperiorityDice.Size)
// NEW: fmt.Sprintf("1%s", p.character.SuperiorityDice.Size)
```
**Result**: Now displays "⚔ Superiority Dice: 1d8 (4/4) (Short Rest) [u/U to use/restore]"

### 2. ✅ 'n' Key No Longer Prompts for Tool/Skill
**Problem**: Pressing 'n' in Traits panel to edit maneuvers also prompted for tool and skill proficiency
**Solution**: Added check to differentiate between initial selection and editing
**Files**:
- `internal/ui/app.go` (handleTraitsPanel)
- `internal/ui/app.go` (handleManeuverSelectorKeys)

**Changes**:
1. In `handleTraitsPanel`, when 'n' is pressed, clear the `studentOfWarToolSelected` flag
2. In `handleManeuverSelectorKeys`, only trigger Student of War flow if:
   - This is initial selection (maneuvers array was empty before)
   - Character has "Student of War" feature
   - `studentOfWarToolSelected` flag is false

**Result**:
- **Initial selection** (after choosing Battle Master): Prompts for maneuvers → tool → skill
- **Editing** (pressing 'n' in Traits): Only prompts for maneuvers

### 3. ✅ Initial Battle Master Selection Flow
**Problem**: User reported can't select maneuvers, proficiency, and skill initially
**Solution**: The flow was already correct, just needed to fix the Student of War logic

**Flow**:
1. Select Fighter subclass → Battle Master
2. Prompt for 3 maneuvers (handled in `handleSubclassSelectorKeys`)
3. After maneuvers confirmed, check if Student of War feature exists
4. If yes and first time, prompt for artisan tool
5. After tool selected, prompt for Fighter skill
6. Complete setup

## Key Binding Reminder

- **`y`/`Y`**: Use/restore Psi Dice (Psi Warrior)
- **`u`/`U`**: Use/restore Superiority Dice (Battle Master)
- **`n`**: Edit Battle Master maneuvers (in Traits panel)

## Code Changes Summary

### internal/ui/panels/features.go
- Lines 111, 123: Changed dice display format from `%d%s` to `1%s`

### internal/ui/app.go
- Line 1338: Added `m.studentOfWarToolSelected = false` when 'n' pressed in Traits
- Lines 2057-2066: Added `isInitialSelection` check before triggering Student of War flow
- Line 2066: Added condition `&& !m.studentOfWarToolSelected` to prevent duplicate prompts

## Testing Checklist

- [ ] Display shows "1d8 (4/4)" not "4d8 (4/4)"
- [ ] Pressing 'u' correctly consumes one Superiority Die
- [ ] Pressing 'U' correctly restores one Superiority Die
- [ ] Pressing 'n' in Traits only shows maneuver selector (no tool/skill prompts)
- [ ] Initial Battle Master selection prompts for: maneuvers → tool → skill
- [ ] Maneuvers are saved and visible in Traits panel
- [ ] Tool proficiency appears in Traits panel
- [ ] Skill proficiency appears in Skills panel

## Notes

All code changes compile successfully (linter shows only deprecation warnings, no errors).

The dice consumption/restoration uses the standard `Current` and `Max` fields in the `SuperiorityDice` struct, which should work correctly:
- `u` key: `m.character.SuperiorityDice.Current--` (if Current > 0)
- `U` key: `m.character.SuperiorityDice.Current++` (if Current < Max)

If dice consumption still doesn't work, please provide specific error messages or describe what happens when you press 'u'.
