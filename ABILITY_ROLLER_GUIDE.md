# Ability Roller Guide - Roll Tests & Saving Throws

## Quick Access

From **Stats Panel**, press **`t`** to open the Ability Roller!

## What It Does

The Ability Roller lets you quickly roll:
- **Ability Checks** (d20 + modifier)
- **Saving Throws** (d20 + modifier + proficiency if proficient)

## Interface

```
ROLL ABILITY

► SELECT ABILITY:

  ▶ Strength     15 (+2) ●
    Dexterity    14 (+2) ●
    Constitution 13 (+1)  
    Intelligence  8 (-1)  
    Wisdom       12 (+1)  
    Charisma     10 (+0)  

  SELECT TYPE:

    Ability Check     - 1d20+2 (modifier only)
    Saving Throw      - 1d20+4 (modifier + 2 prof)

● = Proficient in saving throw

↑/↓: Navigate  Tab: Switch Section  Enter: Roll  Esc: Cancel
```

## Two Selection Modes

### 1. SELECT ABILITY (Top Section)
Navigate through the six abilities to choose which one to roll.

Shows:
- Ability name
- Current score
- Modifier
- ● if proficient in saving throw

### 2. SELECT TYPE (Bottom Section)
Choose between:
- **Ability Check**: Just the modifier
- **Saving Throw**: Modifier + proficiency (if proficient)

## Keyboard Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate within current section |
| `Tab` | Switch between Ability and Type sections |
| `Space` | Toggle type (Check ↔ Save) |
| `Enter` | Roll the dice! |
| `Esc` | Cancel |

## How It Works

### Step 1: Select Ability
```
► SELECT ABILITY:      ← You're here

  ▶ Dexterity    14 (+2) ●  ← Selected
```

Use `↑/↓` to choose the ability.

### Step 2: Choose Type
```
  SELECT ABILITY:

    Dexterity    14 (+2) ●

► SELECT TYPE:         ← Press Tab to get here

  ▶ Ability Check     - 1d20+2 (modifier only)
    Saving Throw      - 1d20+4 (modifier + 2 prof)
```

Use `↑/↓` or `Space` to choose type.

### Step 3: Roll!
Press `Enter` and see the result in the Dice panel!

## Examples

### Example 1: Dexterity Check (Sneak)
```
1. Press 't' in Stats panel
2. Navigate to Dexterity (↓)
3. Make sure "Ability Check" is selected
4. Press Enter

Result: Rolls 1d20+2 (your DEX modifier)
Message: "Dexterity Check: [Roll: 15] + 2 = 17"
```

### Example 2: Dexterity Saving Throw (Dodge Fireball)
```
1. Press 't' in Stats panel
2. Navigate to Dexterity (↓)
3. Press Tab to switch to type selection
4. Navigate to "Saving Throw" (↓)
5. Press Enter

Result: Rolls 1d20+4 (DEX modifier +2 + proficiency +2)
Message: "Dexterity Saving Throw: [Roll: 12] + 4 = 16"
```

### Example 3: Strength Check (Break Door)
```
1. Press 't' in Stats panel
2. Strength is already selected (▶)
3. "Ability Check" is already selected
4. Press Enter

Result: Rolls 1d20+2
Message: "Strength Check: [Roll: 18] + 2 = 20"
```

### Example 4: Wisdom Saving Throw (Resist Charm)
```
1. Press 't' in Stats panel
2. Navigate to Wisdom (↓↓↓↓)
3. Press Tab
4. Navigate to "Saving Throw" (↓)
5. Press Enter

Result: Rolls 1d20+1 (only modifier, not proficient)
Message: "Wisdom Saving Throw: [Roll: 14] + 1 = 15"
```

## Visual Breakdown

### What You See:

```
Strength     15 (+2) ●
             ││  ││  │
             ││  ││  └─ Proficient in STR saves (gets +prof)
             ││  └└─ Modifier: +2 (used in all rolls)
             └└─ Score: 15
```

### The Two Roll Types:

**Ability Check:**
```
Ability Check - 1d20+2 (modifier only)
                     ││
                     └└─ Just your DEX modifier
```

**Saving Throw (Not Proficient):**
```
Saving Throw - 1d20+1 (modifier)
                    ││
                    └└─ Just your WIS modifier (not proficient)
```

**Saving Throw (Proficient):**
```
Saving Throw - 1d20+4 (modifier + 2 prof)
                    ││           │
                    ││           └─ Proficiency bonus added!
                    └└─ Total bonus: +2 mod + +2 prof = +4
```

## When to Use Each Type

### Ability Checks (Use Modifier Only)
- **Strength**: Athletics, breaking things, lifting
- **Dexterity**: Acrobatics, Stealth, Sleight of Hand
- **Constitution**: Holding breath, resisting exhaustion
- **Intelligence**: Investigation, Arcana, History
- **Wisdom**: Perception, Insight, Survival, Animal Handling
- **Charisma**: Persuasion, Deception, Intimidation, Performance

### Saving Throws (Use Modifier + Proficiency if Proficient)
- **Strength**: Resisting being moved, grappled
- **Dexterity**: Dodging fireballs, traps, area effects
- **Constitution**: Resisting poison, disease, death
- **Intelligence**: Resisting illusions, mind effects
- **Wisdom**: Resisting charm, fear, enchantments
- **Charisma**: Resisting banishment, possession

## Quick Reference

### Common Rolls:

| Situation | Ability | Type | Roll |
|-----------|---------|------|------|
| Sneak past guard | DEX | Check | d20 + DEX mod |
| Break down door | STR | Check | d20 + STR mod |
| Spot hidden enemy | WIS | Check | d20 + WIS mod |
| Resist poison | CON | Save | d20 + CON mod (+prof if proficient) |
| Dodge fireball | DEX | Save | d20 + DEX mod (+prof if proficient) |
| Resist charm | WIS | Save | d20 + WIS mod (+prof if proficient) |
| Remember lore | INT | Check | d20 + INT mod |
| Lie convincingly | CHA | Check | d20 + CHA mod |

## Tips & Tricks

### 1. Check Proficiency Markers (●)
Before rolling a save, look for the ● symbol to see if you get proficiency bonus!

### 2. Use Space for Quick Toggle
Instead of `Tab` + `↑/↓`, just press `Space` to toggle between Check and Save.

### 3. Most Common: Ability Checks
Most rolls are ability checks (no proficiency). Saves are for resisting effects.

### 4. Watch the Formula
The roller shows you exactly what will be rolled:
```
Ability Check     - 1d20+2  ← This is what you'll roll
Saving Throw      - 1d20+4  ← Or this
```

### 5. Quick Navigation
- Start typing `t` + `Enter` for STR check (fastest!)
- `t` + `↓` + `Enter` for DEX check
- `t` + `Tab` + `↓` + `Enter` for STR save

## Comparison: Before vs After

### Before (Manual Process):
1. Look at Stats panel for modifier
2. Check if proficient in save
3. Do mental math (+2 mod + +2 prof = +4)
4. Press `f` to focus Dice panel
5. Press `Enter` to input mode
6. Type "1d20+4"
7. Press `Enter` to roll

**Total: 7+ steps and mental math**

### After (Ability Roller):
1. Press `t`
2. Select ability
3. Select type
4. Press `Enter`

**Total: 4 steps, no math needed!**

## Integration with Dice Panel

Results appear in the Dice Roller panel at the bottom:

```
┌─ DICE ROLLER ────────────────────────┐
│ Dexterity Saving Throw:              │
│ [Roll: 15] + 4 = 19                  │
│                                       │
│ History:                              │
│ 1. Dexterity Saving Throw: 19        │
│ 2. Strength Check: 20                │
│ 3. ...                                │
└───────────────────────────────────────┘
```

You can press `f` to focus the Dice panel and press `h` to see history, `r` to reroll!

## Workflow Examples

### Combat Round:
```
DM: "A goblin shoots an arrow at you. Roll Dexterity!"
You: Press 't' → Already on DEX → Enter
Result: Dexterity Check: 18 (you dodge!)
```

### Saving Against Spell:
```
DM: "The wizard casts Hold Person. Wisdom save!"
You: Press 't' → ↓↓↓↓ to WIS → Tab → ↓ to Save → Enter
Result: Wisdom Saving Throw: 15 (you resist!)
```

### Skill Check:
```
DM: "Roll Strength to break the chains."
You: Press 't' → Already on STR → Enter
Result: Strength Check: 12 (chains hold...)
```

## Summary

The Ability Roller is your quick access to:
- ✅ All six ability scores
- ✅ Ability checks (modifier only)
- ✅ Saving throws (with proficiency!)
- ✅ Automatic calculation
- ✅ Instant dice roll
- ✅ Clear visual feedback

**From Stats panel, press `t` and roll!** 🎲

