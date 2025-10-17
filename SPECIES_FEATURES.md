# Species Features System

## Overview

The species features system automatically adds limited-use abilities to your character when you select a species. These features are dynamically calculated based on your character's level, proficiency bonus, and ability scores.

## How It Works

When you select a species (or change species), the system:
1. **Removes** old species features
2. **Parses** the species JSON for traits marked as features
3. **Calculates** max uses based on formulas
4. **Creates** Feature objects with proper rest types
5. **Adds** them to your character's Features list

## Species JSON Structure

Each species trait can now include feature metadata:

```json
{
  "name": "Healing Hands",
  "description": "You can touch a creature and restore hit points equal to your level. Once per long rest.",
  "is_feature": true,
  "max_uses": 1,
  "rest_type": "Long Rest",
  "uses_formula": "",
  "effect_formula": "level"
}
```

### Fields

- **`is_feature`**: `true` if this trait should create a limited-use feature
- **`max_uses`**: Static number or formula for max uses (e.g., `"1"`, `"proficiency"`)
- **`rest_type`**: When it recharges - `"Short Rest"`, `"Long Rest"`, or `"Daily"`
- **`uses_formula`**: Dynamic formula for calculating uses (e.g., `"proficiency"`, `"level"`)
- **`effect_formula`**: Formula for the effect value (e.g., `"level"`, `"1d12+con"`, `"2d6"`)

### Supported Formulas

#### Uses Formula
- `"proficiency"` - Uses = Proficiency Bonus
- `"level"` - Uses = Character Level
- `"con"` / `"str"` / `"dex"` / `"int"` / `"wis"` / `"cha"` - Uses = Ability Score
- Any integer (e.g., `"1"`, `"3"`)

#### Effect Formula
- `"level"` - Displays as "X HP" (for healing)
- `"proficiency"` - Displays as "+X"
- `"1d12+con"` - Displays as "1d12+X" (with CON modifier)
- `"2d6"`, `"3d6"`, etc. - Displays as "Xd6 damage"

## Example Species Features

### Aasimar - Healing Hands
```json
{
  "name": "Healing Hands",
  "description": "You can touch a creature and restore hit points equal to your level. Once per long rest.",
  "is_feature": true,
  "max_uses": 1,
  "rest_type": "Long Rest",
  "uses_formula": "",
  "effect_formula": "level"
}
```

**Result at Level 5:**
- **Name**: Healing Hands
- **Max Uses**: 1
- **Recharge**: Long Rest
- **Effect**: Heals 5 HP
- **Description**: "You can touch a creature... (Effect: 5 HP)"

### Dragonborn - Breath Weapon
```json
{
  "name": "Breath Weapon",
  "description": "Exhale destructive energy in a 15-foot cone. Deals 2d6 damage...",
  "is_feature": true,
  "max_uses": "proficiency",
  "rest_type": "Long Rest",
  "uses_formula": "proficiency",
  "effect_formula": "2d6"
}
```

**Result at Level 5 (Proficiency +3):**
- **Name**: Breath Weapon
- **Max Uses**: 3
- **Recharge**: Long Rest
- **Effect**: 2d6 damage
- **Description**: "Exhale destructive energy... (Effect: 2d6 damage)"

### Goliath - Stone's Endurance
```json
{
  "name": "Stone's Endurance",
  "description": "When you take damage, use reaction to reduce by 1d12 + CON modifier...",
  "is_feature": true,
  "max_uses": 1,
  "rest_type": "Long Rest",
  "uses_formula": "",
  "effect_formula": "1d12+con"
}
```

**Result with CON 16 (+3):**
- **Name**: Stone's Endurance
- **Max Uses**: 1
- **Recharge**: Long Rest
- **Effect**: 1d12+3
- **Description**: "When you take damage... (Effect: 1d12+3)"

### Orc - Adrenaline Rush
```json
{
  "name": "Adrenaline Rush",
  "description": "Take Dash action as bonus action...",
  "is_feature": true,
  "max_uses": "proficiency",
  "rest_type": "Long Rest",
  "uses_formula": "proficiency",
  "effect_formula": ""
}
```

**Result at Level 5 (Proficiency +3):**
- **Name**: Adrenaline Rush
- **Max Uses**: 3
- **Recharge**: Long Rest

### Orc - Relentless Endurance
```json
{
  "name": "Relentless Endurance",
  "description": "When reduced to 0 HP, drop to 1 HP instead. Once per long rest.",
  "is_feature": true,
  "max_uses": 1,
  "rest_type": "Long Rest",
  "uses_formula": "",
  "effect_formula": ""
}
```

**Result:**
- **Name**: Relentless Endurance
- **Max Uses**: 1
- **Recharge**: Long Rest

## Current Species with Features

| Species | Feature | Uses | Rest Type | Scales With |
|---------|---------|------|-----------|-------------|
| **Aasimar** | Healing Hands | 1 | Long Rest | Level (healing) |
| **Dragonborn** | Breath Weapon | Proficiency | Long Rest | Level (damage) |
| **Goliath** | Stone's Endurance | 1 | Long Rest | CON (reduction) |
| **Orc** | Adrenaline Rush | Proficiency | Long Rest | - |
| **Orc** | Relentless Endurance | 1 | Long Rest | - |

## Dynamic Updates

Features automatically update when:
- **Leveling up** - Uses recalculated (proficiency, level-based)
- **Ability score changes** - Effect formulas recalculated (CON, etc.)
- **Species change** - Old features removed, new ones added
- **Rest** - Uses restored based on rest type

## Using Species Features

1. **View Features**: Press `5` or Tab to the Features tab
2. **Use Feature**: Press `u` to consume a charge
3. **Rest**: Press `r` (short rest) or `Shift+R` (long rest) to recover
4. **Track Usage**: Features show current/max uses (e.g., "Breath Weapon (2/3)")

## Adding New Species Features

To add a feature to a species:

1. Open `data/species.json`
2. Find the species
3. Add a trait with `"is_feature": true`
4. Set `max_uses`, `rest_type`, `uses_formula`, and `effect_formula`
5. Save and restart the app
6. Select the species to see the feature

Example:
```json
{
  "name": "New Feature",
  "description": "Feature description here",
  "is_feature": true,
  "max_uses": 1,
  "rest_type": "Short Rest",
  "uses_formula": "proficiency",
  "effect_formula": "level"
}
```

## Technical Details

### Code Structure
- **`internal/models/species.go`**: Feature parsing and application
- **`internal/models/features.go`**: Feature data model
- **`data/species.json`**: Species definitions with feature metadata

### Key Functions
- `CalculateFeatureUses()` - Calculates max uses from formula
- `CalculateFeatureEffect()` - Formats effect value for display
- `ConvertTraitToFeature()` - Creates Feature from SpeciesTrait
- `ApplySpeciesFeatures()` - Adds features to character
- `RemoveSpeciesFeatures()` - Removes old species features

### Future Enhancements
- [ ] More complex formulas (e.g., "level/2", "proficiency*2")
- [ ] Conditional features based on level thresholds
- [ ] Feature choices (like High Elf cantrip selection)
- [ ] Multi-tier features (more uses at higher levels)
- [ ] Resource pools (Ki points, Channel Divinity)

