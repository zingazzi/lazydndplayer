# Ability Roller - Navigation Fix

## Issue Fixed ✅

**Before**: Pressing Tab to switch sections didn't allow navigation in the Type section.
**After**: Tab works! Up/Down keys now work in BOTH sections.

## How It Works Now

### Step-by-Step Navigation:

```
1. Press 't' in Stats panel

ROLL ABILITY

► SELECT ABILITY:              ← Focus here (arrow shows it)

  ▶ Strength     15 (+2) ●     ← Selected
    Dexterity    14 (+2) ●
    ...

  SELECT TYPE:

    Ability Check     - 1d20+2
    Saving Throw      - 1d20+4
```

### Use Up/Down to Navigate Abilities:
```
Press ↓:

► SELECT ABILITY:

    Strength     15 (+2) ●
  ▶ Dexterity    14 (+2) ●     ← Moved down
    ...
```

### Press Tab to Switch to Type Section:
```
ROLL ABILITY

  SELECT ABILITY:              ← No arrow (not focused)

    Dexterity    14 (+2) ●

► SELECT TYPE:                 ← Arrow shows focus here!

  ▶ Ability Check     - 1d20+2  ← Selected
    Saving Throw      - 1d20+4
```

### Now Up/Down Works in Type Section:
```
Press ↓:

► SELECT TYPE:

    Ability Check     - 1d20+2
  ▶ Saving Throw      - 1d20+4  ← Moved down!
```

### Press Tab Again to Go Back:
```
► SELECT ABILITY:              ← Back to abilities

  ▶ Dexterity    14 (+2) ●

  SELECT TYPE:

    Saving Throw      - 1d20+4
```

## Complete Flow Example

**Goal**: Roll a Dexterity Saving Throw

```
1. Press 't'
   → Opens roller, Strength selected

2. Press '↓'
   → Moves to Dexterity

3. Press 'Tab'
   → Switches focus to Type section
   → Message: "Switched focus"

4. Press '↓'
   → Changes from "Ability Check" to "Saving Throw"

5. Press 'Enter'
   → Rolls 1d20+4 (DEX +2 + Prof +2)
   → Result: "Dexterity Saving Throw: 19"
```

## Visual Indicators

### Focus Indicator (►):
Shows which section is active:
```
► SELECT ABILITY:    ← You're here, ↑/↓ navigates abilities
  SELECT TYPE:       ← Not focused

Or:

  SELECT ABILITY:    ← Not focused
► SELECT TYPE:       ← You're here, ↑/↓ navigates types
```

### Selection Indicator (▶):
Shows which item is selected in current section:
```
► SELECT ABILITY:

    Strength     15 (+2) ●
  ▶ Dexterity    14 (+2) ●  ← Selected ability
    Constitution 13 (+1)
```

## Keyboard Summary

| Key | Action | Works In |
|-----|--------|----------|
| `↑/↓` | Navigate items | BOTH sections ✅ |
| `Tab` | Switch sections | Always |
| `Space` | Toggle type | Both (quick alternative to ↑/↓ in type) |
| `Enter` | Roll! | Always |
| `Esc` | Cancel | Always |

## What Changed in Code

### Before:
```go
func NextAbility() {
    if focusOnAbility {
        // Only worked in ability section
    }
}
```

### After:
```go
func Next() {
    if focusOnAbility {
        // Navigate abilities
    } else {
        // Navigate types ← NEW!
    }
}
```

## Quick Test

Try this to verify it works:

```
1. Open app: ./lazydndplayer
2. Press 't'
3. Press 'Tab'
4. Check message bar: "Switched focus" ✅
5. Press '↓'
6. Type section should change ✅
7. Press 'Tab' again
8. Press '↓'
9. Ability should change ✅
```

## Alternative: Space Key

Don't want to use Tab + ↑/↓? Use Space!

```
1. Press 't'
2. Navigate to ability with ↑/↓
3. Press 'Space' to toggle Check ↔ Save
4. Press 'Enter' to roll

No need for Tab! 🎉
```

## Common Workflows

### Quick Ability Check (No Tab Needed):
```
't' → ↑/↓ to select ability → Enter
```

### Quick Save (Use Space):
```
't' → ↑/↓ to select ability → Space → Enter
```

### Precise Selection (Use Tab):
```
't' → ↑/↓ for ability → Tab → ↑/↓ for type → Enter
```

## To Apply Fix

```bash
go build -o lazydndplayer .
./lazydndplayer
```

Now Tab + Up/Down works perfectly! 🎉
