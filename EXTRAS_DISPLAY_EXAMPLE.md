# Extra Bonuses Display - Visual Examples

## New Format: Base + Extra = Total (Modifier)

The extras screen now shows the complete calculation for each ability score!

## Example 1: Dragonborn Fighter

### After Rolling 4d6 and Assigning

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

▶ Strength    : 15 +0 = 15 (mod: +2)  ← Currently selected
  Dexterity   : 14 +0 = 14 (mod: +2)
  Constitution: 13 +0 = 13 (mod: +1)
  Intelligence:  8 +0 =  8 (mod: -1)
  Wisdom      : 12 +0 = 12 (mod: +1)
  Charisma    : 10 +0 = 10 (mod:  0)

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Assignment
```

### After Adding Dragonborn Bonuses (+2 STR, +1 CON)

Press `+` twice on Strength, then `↓↓` to Constitution, then `+` once:

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

  Strength    : 15 +2 = 17 (mod: +3)  ← Wow! +3 modifier now!
  Dexterity   : 14 +0 = 14 (mod: +2)
▶ Constitution: 13 +1 = 14 (mod: +2)  ← Just increased!
  Intelligence:  8 +0 =  8 (mod: -1)
  Wisdom      : 12 +0 = 12 (mod: +1)
  Charisma    : 10 +0 = 10 (mod:  0)

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Assignment
```

## Example 2: Using Point Buy

### After Setting Base Scores

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

▶ Strength    : 15 +0 = 15 (mod: +2)
  Dexterity   : 14 +0 = 14 (mod: +2)
  Constitution: 13 +0 = 13 (mod: +1)
  Intelligence: 12 +0 = 12 (mod: +1)
  Wisdom      : 10 +0 = 10 (mod:  0)
  Charisma    :  8 +0 =  8 (mod: -1)

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Assignment
```

### With Half-Elf Bonuses (+2 CHA, +1 any two)

Adding +1 to STR, WIS, and +2 to CHA:

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

  Strength    : 15 +1 = 16 (mod: +3)  ← Bumped to 16!
  Dexterity   : 14 +0 = 14 (mod: +2)
  Constitution: 13 +0 = 13 (mod: +1)
  Intelligence: 12 +0 = 12 (mod: +1)
  Wisdom      : 10 +1 = 11 (mod:  0)  ← Still +0 mod
▶ Charisma    :  8 +2 = 10 (mod:  0)  ← Went from -1 to 0!

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Stats Panel
```

Notice: "Esc: Back to Stats Panel" (because we came from 'e' key)

## Example 3: Direct Edit from Stats Panel

Press `e` in Stats panel with existing character:

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

▶ Strength    : 10 +2 = 12 (mod: +1)  ← Existing character
  Dexterity   : 12 +0 = 12 (mod: +1)
  Constitution: 14 +1 = 15 (mod: +2)
  Intelligence: 13 +0 = 13 (mod: +1)
  Wisdom      : 10 +0 = 10 (mod:  0)
  Charisma    :  8 +1 =  9 (mod: -1)

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Stats Panel
```

## Example 4: While Typing (Editing Mode)

When you press `e` to type an exact value:

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

▶ Strength    : 15 + 3█ = ? (mod: ?)  ← Typing cursor
  Dexterity   : 14 +0 = 14 (mod: +2)
  Constitution: 13 +1 = 14 (mod: +2)
  Intelligence:  8 +0 =  8 (mod: -1)
  Wisdom      : 12 +0 = 12 (mod: +1)
  Charisma    : 10 +0 = 10 (mod:  0)

  [CONFIRM]

Type number  Enter: Save  Esc: Cancel
```

Press Enter to save:

```
▶ Strength    : 15 +3 = 18 (mod: +4)  ← Calculated!
```

## Example 5: Negative Modifiers (Cursed Character)

```
SET EXTRA BONUSES

Format: Base + Extra = Total (Modifier)

▶ Strength    : 12 -2 = 10 (mod:  0)  ← Cursed: -2 to STR
  Dexterity   : 14 +0 = 14 (mod: +2)
  Constitution: 13 -1 = 12 (mod: +1)  ← Disease: -1 to CON
  Intelligence: 15 +0 = 15 (mod: +2)
  Wisdom      : 10 +0 = 10 (mod:  0)
  Charisma    :  8 -1 =  7 (mod: -2)  ← Scarred: -1 to CHA

  [CONFIRM]
```

## What You See

### Each Line Shows:

```
Ability Name : Base +Extra = Total (mod: +Modifier)
               ^^^^  ^^^^^   ^^^^^       ^^^^^^^^
               │     │       │           └─ Final modifier for rolls
               │     │       └─ Total ability score
               │     └─ Bonus from species/feats/items
               └─ Base rolled/assigned value
```

### Real Example Breakdown:

```
Strength    : 15 +2 = 17 (mod: +3)
              ││ ││  ││       ││
              ││ ││  ││       └└─ Modifier is (17-10)/2 = +3
              ││ ││  └└─ Total: 15 + 2 = 17
              ││ └└─ Extra from Dragonborn
              └└─ Base rolled value
```

## Benefits of New Display

### 1. **See Everything at Once** 👀
- No mental math needed
- Base value clearly shown
- Extra bonus clearly shown
- Total calculated for you
- Modifier calculated for you

### 2. **Understand Impact** 💡
```
Before: Strength: +2  (what does this mean?)
After:  Strength: 15 +2 = 17 (mod: +3)  (crystal clear!)
```

### 3. **Spot Breakpoints** 🎯
```
Constitution: 13 +1 = 14 (mod: +2)
                    ↑
            Adding +1 bumps modifier from +1 to +2!
```

### 4. **Plan Optimization** 📊
```
Dexterity: 15 +0 = 15 (mod: +2)
                        ↑
           Need +1 more to reach +3 modifier (score 16)
```

## Comparison: Before vs After

### Before (Old Display):
```
Strength    : +2
Dexterity   : +0
Constitution: +1
```
Questions:
- What's the base value?
- What's the total?
- What's the final modifier?

### After (New Display):
```
Strength    : 15 +2 = 17 (mod: +3)
Dexterity   : 14 +0 = 14 (mod: +2)
Constitution: 13 +1 = 14 (mod: +2)
```
Answers:
- ✅ Base: 15, 14, 13
- ✅ Total: 17, 14, 14
- ✅ Modifier: +3, +2, +2

## Tips for Using New Display

### Quick Scanning
Look at the rightmost numbers (modifiers) to see your final bonuses.

### Planning Bonuses
Check if adding +1 will bump you to next modifier level:
- Score 13 → 14: Modifier +1 → +2 ✅ Good!
- Score 14 → 15: Modifier +2 → +2 ❌ No change
- Score 15 → 16: Modifier +2 → +3 ✅ Good!

### Verifying Changes
Watch the numbers update in real-time as you press +/-

### Even vs Odd Strategy
Even numbers maximize efficiency:
- 14, 16, 18, 20 = Best modifier per point
- 13, 15, 17, 19 = Same modifier as one less

## Edge Cases

### Unassigned Stats (during generation):
```
▶ Strength    : 10 +0 = 10 (mod:  0)  ← Default 10 if not assigned
  Dexterity   : ---    (not assigned yet)
```

### Maximum Bonus:
```
Strength    : 15 +10 = 25 (mod: +7)  ← Max extra is +10
```

### Minimum Penalty:
```
Strength    : 15 -5 = 10 (mod:  0)  ← Min extra is -5
```

## Summary

The new display format shows:
```
Base + Extra = Total (Modifier)
 ↑      ↑      ↑        ↑
 │      │      │        └─ What you add to d20 rolls
 │      │      └─ Final ability score (for checks)
 │      └─ Your adjustments (species, feats, items)
 └─ Rolled/assigned value (unchangeable here)
```

**Result**: Complete transparency and no guesswork! 🎉
