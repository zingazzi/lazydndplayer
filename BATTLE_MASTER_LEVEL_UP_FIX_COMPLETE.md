# Battle Master Level-Up Fix - COMPLETE

## Problem Summary

When leveling up to Fighter level 3 and selecting Battle Master, the maneuver selector was not appearing. This was happening even though the Battle Master subclass was correctly applied and features were granted.

## Root Cause Analysis

There were TWO separate issues:

### Issue #1: Wrong Code Path (handleSubclassSelectorKeys)
The original fix in `app.go` (`handleSubclassSelectorKeys`) was correct, but it only applies when using the **initial class selection flow** (pressing 'c' in Character Info panel).

### Issue #2: LevelUpSelector Doesn't Trigger Fighter Prompts
When using the **level-up flow** (pressing '+' to level up), a different component (`LevelUpSelector`) handles the subclass selection. This component was NOT checking for Fighter-specific needs (maneuvers/cantrips).

## The Fix

### Part 1: LevelUpSelector Component
**File**: `internal/ui/components/levelupselector.go`

Added two flags to the struct (lines 49-51):
```go
// Fighter subclass flags
NeedsManeuverSelection bool // Battle Master needs maneuvers
NeedsCantripSelection  bool // Eldritch Knight needs cantrips
```

Added Fighter detection after subclass is granted (lines 151-164):
```go
// Check if Fighter subclass needs special handling for Battle Master/Eldritch Knight
if ls.selectedClass == "Fighter" {
    if ls.selectedSubclass == "Battle Master" {
        // Battle Master needs maneuver selection - set flag for app.go to handle
        ls.NeedsManeuverSelection = true
        ls.state = LevelUpComplete
        return *ls, cmd
    } else if ls.selectedSubclass == "Eldritch Knight" {
        // Eldritch Knight needs cantrip/spell selection - set flag for app.go to handle
        ls.NeedsCantripSelection = true
        ls.state = LevelUpComplete
        return *ls, cmd
    }
}
```

### Part 2: App.go Integration
**File**: `internal/ui/app.go`

Modified `handleLevelUpSelectorKeys` (lines 2106-2125):
```go
// Save character if level-up is complete
if !m.levelUpSelector.IsVisible() {
    m.storage.Save(m.character)
    m.character.UpdateDerivedStats()

    // Check if Battle Master needs maneuver selection
    if m.levelUpSelector.NeedsManeuverSelection {
        debug.Log("Battle Master detected - prompting for maneuver selection")
        m.levelUpSelector.NeedsManeuverSelection = false // Clear flag
        m.maneuverSelector.Show(3) // Battle Master starts with 3 maneuvers
        m.message = "Select 3 maneuvers for Battle Master..."
        return m, cmd
    }

    // Check if Eldritch Knight needs cantrip selection
    if m.levelUpSelector.NeedsCantripSelection {
        debug.Log("Eldritch Knight detected - prompting for cantrip selection")
        m.levelUpSelector.NeedsCantripSelection = false // Clear flag
        m.cantripSelector.Show("Wizard", 2) // Eldritch Knight starts with 2 cantrips
        m.message = "Select 2 cantrips from the Wizard spell list..."
        return m, cmd
    }

    m.message = "Character updated!"
}
```

## How It Works Now

### Scenario 1: First Class Selection (Character Creation)
1. Press 'c' in Character Info
2. Select Fighter → Select skills → Select fighting style → Select weapon masteries
3. Press '+' twice to reach level 3
4. When level 3 is reached, subclass selector appears
5. Select Battle Master
6. **→ Maneuver selector appears** ✅
7. After selecting 3 maneuvers → Tool selector appears (Student of War)
8. After selecting tool → Skill selector appears (Student of War)
9. Done!

### Scenario 2: Level-Up Flow (Your Test Case)
1. Create character with Fighter level 1
2. Level up to Fighter 2
3. Level up to Fighter 3 → Subclass selector appears
4. Select Battle Master
5. **→ LevelUpSelector sets `NeedsManeuverSelection = true`**
6. **→ App.go detects the flag and shows maneuver selector** ✅
7. After selecting 3 maneuvers → Tool selector appears (Student of War)
8. After selecting tool → Skill selector appears (Student of War)
9. Done!

## Testing Instructions

Please rebuild and test:

```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go build -o lazydndplayer .
./lazydndplayer
```

### Test Case (Your Exact Scenario):
1. Create new character
2. Add stats
3. Add first level of Fighter (select skills, fighting style, weapon masteries)
4. Add second level of Fighter
5. Add 3rd level of Fighter
6. Select Battle Master subclass

### Expected Result:
✅ Maneuver selector should appear: "Select 3 maneuvers for Battle Master..."
✅ After selecting 3 maneuvers → Tool selector: "Select an artisan tool for Student of War..."
✅ After selecting tool → Skill selector: "Select a skill from the Fighter skill list for Student of War..."
✅ Setup complete!

### Debug Log Should Show:
```
GrantSubclassFeatures: className=Fighter, subclassName=Battle Master, level=3
GrantSubclassFeatures: Found subclass Battle Master
GrantSubclassFeatures: Found 2 features for level 3
  Adding subclass feature: Combat Superiority (uses: 1/1, rest: Short Rest)
  Applying Combat Superiority benefits - Initializing Superiority Dice
  Initialized Superiority Dice: 4d8 (4/4)
  Adding subclass feature: Student of War (uses: 0/0, rest: None)
GrantSubclassFeatures: Granted 2 subclass features
Battle Master detected - prompting for maneuver selection
```

## Files Modified

1. `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/levelupselector.go`
   - Added `NeedsManeuverSelection` and `NeedsCantripSelection` flags
   - Added Fighter subclass detection logic after subclass is granted
   - Fixed indentation issues

2. `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/app.go`
   - Added flag checking in `handleLevelUpSelectorKeys`
   - Triggers maneuver/cantrip selectors when flags are set

## Status

✅ Code changes complete
✅ Logic flow corrected
✅ Indentation fixed
⏳ Ready for build and testing

Please build and test to confirm the maneuver selector now appears!
