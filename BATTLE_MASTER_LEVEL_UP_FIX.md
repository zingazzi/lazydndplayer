# Battle Master Level-Up Fix

## Issue Found

When selecting Battle Master during level up (at Fighter level 3), the maneuver selector was not appearing.

### Root Cause

In `handleSubclassSelectorKeys` (line 2255), the Fighter subclass prompt logic was **outside** the `if len(m.character.Classes) > 0` block. This meant:

1. The `className` variable was declared inside the if block (line 2276)
2. But it was being accessed outside the if block (line 2285)
3. This caused the Battle Master check to never execute

### The Fix

**File**: `internal/ui/app.go` lines 2269-2325

Moved all the Fighter subclass handling code **inside** the `if len(m.character.Classes) > 0` block:

```go
if len(m.character.Classes) > 0 {
    // Set subclass
    m.character.Classes[len(m.character.Classes)-1].Subclass = selectedSubclass.Name

    // Grant features
    className := m.character.Classes[len(m.character.Classes)-1].ClassName
    classLevel := m.character.Classes[len(m.character.Classes)-1].Level
    subclassFeatures := models.GrantSubclassFeatures(m.character, className, selectedSubclass.Name, classLevel)

    m.subclassSelector.Hide()

    // NOW THE FIGHTER CHECK IS INSIDE THE SAME BLOCK
    if className == "Fighter" {
        switch selectedSubclass.Name {
        case "Battle Master":
            // Show maneuver selector
            m.maneuverSelector.Show(3)
            m.message = "Select 3 maneuvers for Battle Master..."
            return m, cmd
        case "Eldritch Knight":
            // Show cantrip selector
            m.cantripSelector.Show("Wizard", 2)
            m.message = "Select 2 cantrips from the Wizard spell list..."
            return m, cmd
        }
    }

    // Rest of the code...
}
```

### Added Debug Logging

Added extensive debug logging to help diagnose issues:
- Line 2284: Logs the className and subclass name being checked
- Line 2287: Logs when Fighter is detected
- Line 2297: Logs when Battle Master is selected
- Line 2302: Logs when Fighter subclass doesn't need prompts
- Line 2305: Logs when it's not a Fighter class

## Testing Instructions

1. Create new character
2. Add stats
3. Add Fighter level 1 (select skills, fighting style, weapon masteries)
4. Add Fighter level 2
5. Add Fighter level 3
6. **Select Battle Master subclass**

### Expected Result

You should now see:
1. ✅ Maneuver selector appears → "Select 3 maneuvers for Battle Master..."
2. ✅ After selecting 3 maneuvers → Tool selector appears → "Select an artisan tool for Student of War..."
3. ✅ After selecting tool → Skill selector appears → "Select a skill from the Fighter skill list for Student of War..."
4. ✅ After selecting skill → Setup complete!

### Debug Log Should Show

```
Subclass selected: Battle Master
Set character subclass to: Battle Master
Granted 2 subclass features: [Combat Superiority Student of War]
Checking for Fighter subclass prompts: className='Fighter', subclass='Battle Master'
Fighter detected, checking subclass type
Battle Master selected - prompting for maneuvers
```

## Other Fixes Included

1. **Display format**: Changed from "4d8 (4/4)" to "1d8 (4/4)"
2. **'n' key in Traits**: Only shows maneuver selector, doesn't re-trigger Student of War
3. **Initial vs Edit**: Logic distinguishes between first-time selection and editing

## Status

✅ Code compiles successfully
✅ No linter errors (only 1 deprecation warning)
✅ Logic flow corrected
✅ Ready for testing

Please test and let me know if the maneuver selector now appears!
