# Weapon Mastery Implementation Summary

## Overview
Implemented a generic, data-driven weapon mastery system that works for any class with the "Weapon Mastery" feature.

## Changes Made

### 1. Data Structure Updates

#### **internal/models/feature_definitions.go**
- Added `Mechanics map[string]interface{}` field to `FeatureDefinition` struct
- Updated `ToFeature()` to copy mechanics data to Feature instances

#### **internal/models/features.go**
- Added `Mechanics map[string]interface{}` field to `Feature` struct

#### **data/classes/fighter.json**
- Added `mechanics.weapons_mastered: 3` to Level 1 "Weapon Mastery" feature
- This makes the mastery count data-driven instead of hardcoded

### 2. New Weapon Mastery Descriptions Module

#### **internal/models/weapon_masteries.go** (NEW FILE)
Created comprehensive weapon mastery descriptions with:
- `WeaponMasteryDescription` struct
- `GetWeaponMasteryDescriptions()` - returns all mastery descriptions
- `GetMasteryDescription(name)` - returns specific mastery description

Includes all 9 weapon masteries:
- **Cleave**: Hit second creature within 5 ft
- **Graze**: Deal ability modifier damage on miss
- **Nick**: Extra Light attack as part of Attack action
- **Push**: Push creature 10 ft away
- **Sap**: Target has disadvantage on next attack
- **Slow**: Reduce target's speed by 10 ft
- **Topple**: Force Constitution save or prone
- **Vex**: Advantage on next attack against target
- **Ensnare**: Restrain target (Net only)

### 3. Updated Weapon Mastery Selector

#### **internal/ui/components/weaponmasteryselector.go**
- Completely redesigned UI to show descriptions on the right side
- Two-column layout (similar to class/species selector):
  - **Left**: Weapon list with checkboxes and mastery types
  - **Right**: Detailed mastery description for selected weapon
- Added `wrapMasteryText()` helper function for text wrapping
- Adjusted dimensions: 50 chars width per box, 25 lines height

### 4. Updated Traits Panel

#### **internal/ui/panels/traits.go**
- Enhanced "Weapon Mastery" section to show mastery descriptions
- For each mastered weapon, now displays:
  - ✓ Weapon name
  - Mastery type (in bold purple)
  - Full mastery description (wrapped, italicized)
- Updated `getWeaponMasteryCount()` to read from feature mechanics (generic approach)
- Removed hardcoded class checks - now reads from `feature.Mechanics["weapons_mastered"]`

### 5. Updated App Controller

#### **internal/ui/app.go**
- Fixed `getWeaponMasteryCount()` to use generic approach:
  - Reads `weapons_mastered` from feature.Mechanics
  - Works for any class with the feature
  - No hardcoded class names
- Added weapon mastery selection to class creation flow:
  - After fighting style selection → check for weapon mastery
  - After cantrip selection → check for weapon mastery
  - After skill selection (no fighting style/cantrips) → check for weapon mastery
- Updated `handleWeaponMasterySelectorKeys` to complete class setup properly

## How It Works

### For Developers

1. **Add Weapon Mastery to Any Class**:
   ```json
   {
     "name": "Weapon Mastery",
     "description": "You can master N weapons...",
     "uses_formula": "0",
     "rest_type": "None",
     "mechanics": {
       "weapons_mastered": 3
     }
   }
   ```

2. **The System Automatically**:
   - Reads `weapons_mastered` count from feature mechanics
   - Prompts for weapon selection during class creation
   - Displays mastery descriptions in selector
   - Shows active masteries with descriptions in Traits panel

3. **Scaling**: For classes where weapon mastery count increases:
   ```json
   {
     "level": 4,
     "features": [
       {
         "name": "Weapon Mastery",
         "description": "You can master 4 weapons.",
         "mechanics": {
           "weapons_mastered": 4
         }
       }
     ]
   }
   ```

### For Users

1. **During Class Selection**:
   - After selecting Fighter (or any class with Weapon Mastery)
   - System automatically shows weapon mastery selector
   - Navigate with ↑/↓, select with Space
   - See mastery description on the right panel
   - Confirm with Enter

2. **In Traits Panel** (press 5 in main menu):
   - View all mastered weapons with descriptions
   - Press 'm' to change weapon masteries at any time
   - See available weapons to master

3. **Mastery Descriptions**:
   - Each weapon shows its mastery property
   - Full description explains the mechanical benefit
   - Wrapped text for readability

## Testing

To test the implementation:
1. Run `./lazydndplayer`
2. Press 'c' in Character Info panel to select Fighter class
3. Select 2 skills
4. Choose a fighting style
5. **NEW**: Weapon mastery selector appears automatically
6. Navigate and select 3 weapons
7. View them in Traits panel (press 5)

## Benefits

✅ **Generic**: Works for any class with Weapon Mastery feature
✅ **Data-Driven**: Mastery count stored in JSON, not hardcoded
✅ **Extensible**: Easy to add to other classes (Barbarian, Paladin, etc.)
✅ **User-Friendly**: Clear descriptions help players understand mechanics
✅ **Maintainable**: Single source of truth for mastery descriptions

## Future Enhancements

- Add weapon mastery count scaling for level-ups
- Implement mastery property effects in combat calculations
- Add visual indicators for mastery effects in Actions panel
- Support for changing individual masteries (not just all at once)
