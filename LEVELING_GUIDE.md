# Leveling Guide System

## Overview

The leveling guide system provides a comprehensive JSON-based reference for what characters gain at each level, including features, spell slots, HP increases, and species-specific scaling.

## Files

- **`data/leveling_guide.json`** - Complete leveling information for all classes and species
- **`internal/leveling/guide.go`** - Loader and helper functions for the leveling guide
- **`internal/leveling/levelup.go`** - Level-up logic with automatic species bonus calculation

## Features

### 1. Dwarven Toughness (Automatic HP Scaling)

Dwarves gain **+1 HP per level** from their "Dwarven Toughness" trait.

**How it works:**
- When you select Dwarf as your species: **+1 HP** immediately
- When you level up as a Dwarf: **+1 more HP** automatically
- When you change species from Dwarf: **-Level HP** (removes all Dwarven Toughness HP)

**Example:**
```
Level 1 Dwarf Fighter:
- Base HP: 10
- Dwarven Toughness: +1
- Total MaxHP: 11

Level 2 Dwarf Fighter:
- Base HP: 10
- HP gained at level 2: +6
- Dwarven Toughness: +2 (one per level)
- Total MaxHP: 18

Level 5 Dwarf Fighter:
- Base HP: 10
- HP from leveling: +24
- Dwarven Toughness: +5
- Total MaxHP: 39
```

### 2. Leveling Guide JSON

The guide contains:
- **Class progression** (features, spell slots, proficiency bonus)
- **Species scaling** (HP bonuses, feature upgrades)
- **General rules** (ASI levels, proficiency bonus progression)

## Using the Leveling Guide

### In the App

When you level up:
1. The app automatically applies Dwarven Toughness (if Dwarf)
2. HP increases are added (roll or average)
3. Proficiency bonus updates automatically
4. Species features scale (Dragonborn breath weapon, Aasimar healing, etc.)

### Manual Reference

Check `data/leveling_guide.json` for:
- What features you gain at each level
- When you get spell slots
- Class-specific information
- Notes and reminders

## JSON Structure

### Class Levels

```json
{
  "classes": {
    "Fighter": {
      "hit_die": 12,
      "levels": {
        "1": {
          "proficiency_bonus": 2,
          "features": ["Fighting Style", "Second Wind"],
          "spell_slots": {},
          "notes": "Choose a Fighting Style..."
        }
      }
    }
  }
}
```

### Species Scaling

```json
{
  "species_scaling": {
    "Dwarf": {
      "hp_per_level": 1,
      "note": "Dwarven Toughness: +1 HP per level"
    },
    "Dragonborn": {
      "features": {
        "5": "Breath Weapon damage increases to 3d6",
        "11": "Breath Weapon damage increases to 4d6",
        "17": "Breath Weapon damage increases to 5d6"
      }
    }
  }
}
```

## Level-Up Checklist

When leveling up, remember to:

### 1. Hit Points
- **Roll** your class's hit die + CON modifier
- **OR** take the average: (hit die ÷ 2) + 1 + CON modifier
- Dwarves automatically get +1 HP (Dwarven Toughness)

### 2. Proficiency Bonus
Automatically updates at levels: 5, 9, 13, 17

### 3. Ability Score Improvements (ASI)
At levels **4, 8, 12, 16, 19**:
- Increase one ability by 2, OR
- Increase two abilities by 1 each, OR
- Take a feat

### 4. Class Features
Check the guide for new features at your level:
- Fighters: Extra Attack (5), Action Surge (2), etc.
- Wizards: Arcane Tradition (2), new spell slots
- Clerics: Channel Divinity (2), domain features
- Rogues: Sneak Attack damage increases

### 5. Spell Slots (Spellcasters)
- New spell slots appear automatically in the guide
- Wizards: Add 2 spells to spellbook each level
- Clerics/Druids: Can prepare more spells

### 6. Species Features
Some scale with level:
- **Dragonborn**: Breath Weapon damage increases (5, 11, 17)
- **Aasimar**: Healing Hands heals for current level
- **Dwarf**: Dwarven Toughness continues to add +1 HP

## Programmatic Access

### Get Level Information

```go
import "github.com/marcozingoni/lazydndplayer/internal/leveling"

// Get info for Fighter level 5
info, err := leveling.GetLevelInfo("Fighter", 5)
if err != nil {
    // handle error
}

fmt.Println("Features:", info.Features)
// Output: Features: [Extra Attack]

fmt.Println("Notes:", info.Notes)
// Output: Notes: You can attack twice...
```

### Get Species Scaling

```go
scaling := leveling.GetSpeciesScaling("Dwarf")
if scaling != nil {
    fmt.Println("HP per level:", scaling.HPPerLevel)
    // Output: HP per level: 1
}
```

### Format Level Summary

```go
summary := leveling.FormatLevelUpSummary("Wizard", 5)
fmt.Println(summary)
// Output:
// === LEVEL 5 ===
//
// Proficiency Bonus: +3
//
// Spell Slots:
//   Level 1: 4 slots
//   Level 2: 3 slots
//   Level 3: 2 slots
//
// Add 2 spells to your spellbook
//
// Notes:
// Gain 3rd-level spell slots. Add 2 spells to spellbook
```

## Supported Classes

Currently detailed in the guide:
- ✅ **Fighter** (Levels 1-12)
- ✅ **Wizard** (Levels 1-6)
- ✅ **Cleric** (Levels 1-2)
- ✅ **Rogue** (Levels 1-3)

More classes and higher levels can be added to `data/leveling_guide.json`.

## Species with Special HP Scaling

| Species | HP Bonus | Description |
|---------|----------|-------------|
| **Dwarf** | +1 per level | Dwarven Toughness trait |
| All others | None | Standard HP progression |

## Example Level-Up Flow

### Dwarf Fighter: Level 1 → Level 2

**Before Level Up:**
- Level: 1
- MaxHP: 11 (10 base + 1 Dwarven Toughness)
- Proficiency: +2

**Rolling for HP:**
- d10 roll: 6
- CON modifier: +2
- HP increase: 8

**After Level Up:**
- Level: 2
- MaxHP: 20 (11 + 8 + 1 from Dwarven Toughness)
- Proficiency: +2
- New Features: Action Surge
- Species HP Bonus: 2 (tracked internally)

### Wizard: Level 4 → Level 5

**Before Level Up:**
- Level: 4
- Spell Slots: 1st (4), 2nd (3)

**After Level Up:**
- Level: 5
- Spell Slots: 1st (4), 2nd (3), **3rd (2)** ← NEW!
- New Features: None at this level
- Remember to: Add 2 spells to spellbook, prepare spells for new slots
- ASI Available: YES (level 4 was ASI, you may have taken it already)

## Technical Implementation

### Automatic HP Bonus Application

The `ApplySpeciesHPBonus()` function:
1. Checks species traits for HP-granting abilities
2. Calculates bonus based on character level
3. Adjusts MaxHP automatically
4. Tracks the bonus in `char.SpeciesHPBonus`

```go
// Automatically called on species selection and level up
func ApplySpeciesHPBonus(char *Character, species *SpeciesInfo) {
    oldBonus := char.SpeciesHPBonus
    char.SpeciesHPBonus = 0

    for _, trait := range species.Traits {
        if strings.Contains(strings.ToLower(trait.Name), "dwarven toughness") {
            char.SpeciesHPBonus = char.Level // +1 per level
        }
    }

    hpChange := char.SpeciesHPBonus - oldBonus
    char.MaxHP += hpChange
}
```

### Level Up Process

```go
func PerformLevelUp(char *Character, options LevelUpOptions) {
    char.Level++
    char.MaxHP += options.HPIncrease
    // ... ability scores, spell slots ...
    char.UpdateDerivedStats()

    // Recalculate species HP bonus (Dwarven Toughness grows)
    species := models.GetSpeciesByName(char.Race)
    if species != nil {
        models.ApplySpeciesHPBonus(char, species)
    }
}
```

## Adding More Classes/Levels

To add more to the leveling guide:

1. Open `data/leveling_guide.json`
2. Add your class under `"classes"`:

```json
{
  "classes": {
    "YourClass": {
      "hit_die": 8,
      "levels": {
        "1": {
          "proficiency_bonus": 2,
          "features": ["Feature 1", "Feature 2"],
          "spell_slots": {},
          "notes": "Your notes here"
        }
      }
    }
  }
}
```

3. The guide is automatically loaded and cached
4. Access it via `leveling.GetLevelInfo("YourClass", level)`

## Future Enhancements

Potential additions to the leveling system:
- [ ] Automatic feature tracking in Features tab
- [ ] Level-up UI wizard in the app
- [ ] Multi-class support
- [ ] Automatic spell slot updates
- [ ] HP roll vs average choice in UI
- [ ] ASI UI for selecting ability increases
- [ ] Feat selection UI

## Maintenance

When updating species traits:
1. Update `data/species.json` with the trait
2. Add scaling info to `data/leveling_guide.json` under `species_scaling`
3. If it's an HP bonus, update `ApplySpeciesHPBonus()` in `internal/models/species.go`

## FAQ

**Q: Does Dwarven Toughness work retroactively?**
A: Yes! If you select Dwarf at level 5, you immediately get +5 HP (1 per level).

**Q: What happens if I change species from Dwarf?**
A: All Dwarven Toughness HP is removed. If you're level 5, you lose 5 HP.

**Q: Is the leveling guide used automatically?**
A: Partially. HP bonuses are automatic. Features, spells, and notes are reference material for you to track manually (for now).

**Q: Can I add homebrew classes?**
A: Yes! Edit `data/leveling_guide.json` to add your custom class progression.

**Q: Do other species get HP bonuses?**
A: Currently only Dwarves. Other species get features that scale (like Dragonborn breath weapon damage).

**Q: Where is my species HP bonus shown?**
A: It's tracked internally in `character.SpeciesHPBonus`. Your MaxHP already includes it.

## Summary

The leveling guide system provides:
- ✅ Automatic Dwarven Toughness (+1 HP per level)
- ✅ JSON-based reference for all level-up decisions
- ✅ Species scaling for features
- ✅ Programmatic access to leveling data
- ✅ Extensible structure for adding classes/levels

Use `data/leveling_guide.json` as your reference when leveling up, and let the app handle Dwarven Toughness automatically!
