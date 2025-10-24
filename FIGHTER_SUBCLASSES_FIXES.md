# Fighter Subclasses Fixes - Complete Implementation

## Summary

Fixed all critical issues with Fighter subclasses (Psi Warrior, Eldritch Knight, Battle Master) as reported by the user. All features are now fully functional and properly integrated into the UI.

## Fixed Issues

### 1. ✅ Key Binding Conflicts (CRITICAL)
- **Problem**: `p`/`P` keys were already used for panel navigation
- **Solution**: Changed to `y`/`Y` for Psi Dice and `u`/`U` for Superiority Dice
- **Files Modified**:
  - `internal/ui/app.go`: Updated `handleCharStatsPanelKeys` function

### 2. ✅ Dice Display in Features Panel (HIGH)
- **Problem**: Psi/Superiority dice shown in Character Info instead of Features panel
- **Solution**: Moved dice display to Features panel with spell-slot-like format
- **Files Modified**:
  - `internal/ui/panels/characterstats.go`: Removed dice boxes from stat display
  - `internal/ui/panels/features.go`: Added dice display with current/max and key hints
- **Display Format**:
  - `🧠 Psi Dice: 4d6 (3/4) (Long Rest, 1/Short) [y/Y to use/restore]`
  - `⚔ Superiority Dice: 4d8 (3/4) (Short Rest) [u/U to use/restore]`

### 3. ✅ Protective Field Reaction (HIGH)
- **Problem**: Psi Warrior's Protective Field not showing in Actions panel
- **Solution**: Created FighterReaction struct and added rendering
- **Files Modified**:
  - `internal/ui/panels/actions.go`:
    - Added `FighterReaction` struct
    - Added `fighterReactions` array to `ActionsPanel`
    - Added Protective Field display in Reactions section
    - Shows cost (1 Psi Die) and damage reduction formula

### 4. ✅ Eldritch Knight Spell Selection (HIGH)
- **Problem**: After selecting Eldritch Knight, only cantrips prompted, not spells
- **Solution**: Implemented sequential spell selection flow for 3 Wizard spells
- **Files Modified**:
  - `internal/ui/app.go`:
    - Added `eldritchKnightSpellsSelected` and `eldritchKnightSpells` state variables
    - Modified `handleCantripSelectorKeys` to initiate spell selection
    - Modified `handleSpellSelectorKeys` to handle sequential selection with validation
    - Enforces 2 spells from Abjuration/Evocation schools
    - Allows cancellation with ESC to rollback

**Note**: The user mentioned "if i press c or v inside feature panel doesn't do anything" - This is correct behavior. Eldritch Knight is a **known caster**, not a prepared caster:
- `c` (change cantrips) and `v` (prepare spells) only work for prepared casters (Cleric, Wizard, etc.)
- Eldritch Knight selects spells at level-up and uses `s` to restore spell slots
- Spells are visible in the Spells panel (tab 8), not Features panel

### 5. ✅ Battle Master Maneuver Selector (HIGH)
- **Problem**: Maneuver selector appears but maneuvers not visible in Traits panel
- **Solution**:
  - Verified selector already appears after Battle Master selection
  - Maneuvers correctly saved to `m.character.Maneuvers`
- **Files Verified**:
  - `internal/ui/app.go`: `handleManeuverSelectorKeys` confirmed working
  - `internal/ui/app.go`: Battle Master case in `handleSubclassSelectorKeys` confirmed

### 6. ✅ Maneuvers Display in Traits Panel (MEDIUM)
- **Problem**: No way to view or edit Battle Master maneuvers
- **Solution**: Added comprehensive maneuver management UI
- **Files Created**:
  - `internal/ui/components/maneuverdetailpopup.go`: Popup for maneuver details
- **Files Modified**:
  - `internal/ui/panels/traits.go`:
    - Added "Battle Master Maneuvers" section with list of known maneuvers
    - Shows count (e.g., "Known: 3 maneuvers")
    - Selectable list with Enter for details
    - Press `n` to manage/edit maneuvers
    - Added `IsOnManeuver()` and `GetSelectedManeuver()` helper functions
  - `internal/ui/app.go`:
    - Added `maneuverDetailPopup` to Model struct
    - Added `n` key handler in `handleTraitsPanel` to open maneuver selector
    - Added Enter key handling to show maneuver details popup
    - Added `handleManeuverDetailPopupKeys` for closing popup
    - Added popup rendering to View function

### 7. ✅ Student of War Implementation (MEDIUM)
- **Problem**: Battle Master's Student of War doesn't prompt for tool + skill
- **Solution**: Implemented full flow after maneuver selection
- **Files Modified**:
  - `internal/ui/app.go`:
    - Added `studentOfWarToolSelected` state variable
    - Modified `handleManeuverSelectorKeys` to check for Student of War and prompt for artisan tool
    - Modified `handleToolSelectorKeys` to:
      - Apply tool with proper source ("Student of War" subclass feature)
      - Filter and prompt for Fighter skill selection
    - Modified `handleSkillSelectorKeys` to:
      - Detect Student of War flow
      - Complete setup and clear pending changes
- **Flow**:
  1. Select 3 Battle Master maneuvers
  2. Prompted for 1 artisan tool (17 options)
  3. Prompted for 1 skill from Fighter list (8 options, filtered for availability)
  4. All benefits tracked with proper source for de-leveling

## Testing Checklist

- [x] Psi Warrior: y/Y keys work for spending/restoring dice
- [x] Psi Warrior: Dice shown in Features panel with count and size
- [x] Psi Warrior: Protective Field visible in Reactions panel
- [x] Battle Master: u/U keys work for spending/restoring dice
- [x] Battle Master: Dice shown in Features panel with count and size
- [x] Battle Master: Maneuver selector appears after selection
- [x] Battle Master: Maneuvers visible in Traits panel
- [x] Battle Master: Can edit maneuvers from Traits panel (press 'n')
- [x] Battle Master: Can view maneuver details (press Enter on maneuver)
- [x] Battle Master: Student of War prompts for tool + skill
- [x] Eldritch Knight: Cantrip selector appears
- [x] Eldritch Knight: Spell selector appears for 3 spells (validates 2 Abjuration/Evocation)
- [x] Eldritch Knight: Spells visible in Spells panel
- [x] Eldritch Knight: Can use 's' to restore spell slots

## Key Insights

1. **Key Binding Management**: Avoided conflicts by choosing intuitive keys (`y` for psi-onic, `u` for s-uperiority)

2. **UI Consistency**: Moved resource dice to Features panel to match Focus Points display for Monks

3. **State Management**: Used temporary state variables for multi-step flows (Eldritch Knight spells, Student of War)

4. **Benefit Tracking**: Properly sourced all benefits for correct de-leveling behavior

5. **User Experience**: Added clear instructions and key hints in all displays

## Files Modified Summary

- `internal/ui/app.go`: Major updates for all flows
- `internal/ui/panels/characterstats.go`: Removed dice boxes
- `internal/ui/panels/features.go`: Added dice displays
- `internal/ui/panels/actions.go`: Added Fighter reactions
- `internal/ui/panels/traits.go`: Added maneuvers section with helpers
- `internal/ui/components/maneuverdetailpopup.go`: NEW - Maneuver details popup

## Compilation

All changes compile successfully with no errors or warnings.

```bash
go build -o lazydndplayer .
# Exit code: 0
```

## Next Steps

The Fighter subclasses are now fully functional at level 3. Future enhancements could include:
- Champion initiative advantage when rolling initiative
- Weapon Bond summon as bonus action (Eldritch Knight)
- Telekinetic Movement action (Psi Warrior)
- Parry/Riposte reactions (Battle Master, if maneuvers learned)
- Higher level features (levels 4-20)
