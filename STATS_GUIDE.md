# Ability Score Generation Guide

## Overview

The Lazy D&D Player application now supports three official D&D 5e methods for generating ability scores, plus the ability to add extra bonuses from species, feats, and other sources.

## How to Use

1. **Open the Stats Tab**: Navigate to the Stats panel (default tab when app opens)
2. **Press 'e'**: This opens the Ability Score Generator popup
3. **Select Method**: Choose one of three methods:
   - **4d6 Drop Lowest**
   - **Standard Array**
   - **Point Buy**
4. **Assign/Generate Stats**: Depending on method chosen
5. **Set Extras**: Add bonuses from species, feats, etc.
6. **Confirm**: Apply the stats to your character

## The Three Methods

### 1. 4d6 Drop Lowest

**What it is**: The classic random rolling method
- Roll four six-sided dice (4d6)
- Drop the single lowest die result
- Add the remaining three dice together
- Repeat this process until you have six scores
- Assign these six scores to your abilities as you see fit

**In the App**:
- The app automatically rolls 4d6 six times and shows you each roll
- Shows the dice results and which die was dropped: `[2, 4, 5, 6] drop 2 = 15`
- Use number keys 1-6 to assign each rolled value to an ability
- Navigate with ↑/↓ keys, select with Enter

**Example**:
```
Roll Results:
1: [3, 4, 5, 6] drop 3 = 15
2: [1, 2, 4, 5] drop 1 = 11
3: [2, 3, 4, 6] drop 2 = 13
4: [1, 3, 5, 6] drop 1 = 14
5: [2, 4, 5, 5] drop 2 = 14
6: [1, 2, 3, 3] drop 1 = 8

Assign to abilities:
▶ Strength    : 15  (press '1' to assign roll #1)
  Dexterity   : 14  (press '4' to assign roll #4)
  ...
```

### 2. Standard Array

**What it is**: A pre-determined, balanced set of scores
- No randomness involved
- Uses the fixed set: **[15, 14, 13, 12, 10, 8]**
- Assign these scores to your six abilities in any order
- Ensures a balanced character without the variance of dice rolls

**In the App**:
- Shows all six values from the standard array
- Use number keys 1-6 to assign each value to an ability
- Navigate with ↑/↓ keys, select with Enter

**Example**:
```
Available: [15, 14, 13, 12, 10, 8]

▶ Strength    : 15  (assigned value 1)
  Dexterity   : 14  (assigned value 2)
  Constitution: 13  (assigned value 3)
  Intelligence: 12  (assigned value 4)
  Wisdom      : 10  (assigned value 5)
  Charisma    : 8   (assigned value 6)
```

### 3. Point Buy

**What it is**: A customizable, point-spending system
- Start with a pool of **27 points**
- All abilities start at 8
- Spend points to increase abilities up to 15
- Different scores have different costs (non-linear)

**Cost Table**:
| Score | Cost | To Reach (from 8) |
|-------|------|-------------------|
| 8     | 0    | 0 points          |
| 9     | 1    | 1 point           |
| 10    | 2    | 2 points          |
| 11    | 3    | 3 points          |
| 12    | 4    | 4 points          |
| 13    | 5    | 5 points          |
| 14    | 7    | 7 points          |
| 15    | 9    | 9 points          |

**In the App**:
- Shows remaining points: `Points Remaining: 27 / 27`
- Navigate abilities with ↑/↓
- Use + to increase, - to decrease
- Can't go below 8 or above 15
- Can't spend more than 27 points total

**Example**:
```
Points Remaining: 0 / 27

▶ Strength    : 15 (cost: 9)   [press + to increase]
  Dexterity   : 14 (cost: 7)   [press - to decrease]
  Constitution: 13 (cost: 5)
  Intelligence: 12 (cost: 4)
  Wisdom      : 10 (cost: 2)
  Charisma    : 8  (cost: 0)

Total: 9+7+5+4+2+0 = 27 points
```

## Setting Extra Bonuses

After choosing your base stats, you can add extra bonuses from:
- **Species**: Most species give +2 to one ability and +1 to another (or +1 to three)
- **Feats**: Some feats give +1 to an ability
- **Magic Items**: If your DM allows starting magic items
- **Other sources**: Blessings, curses, or custom campaign rules

**In the App**:
- Navigate to each ability with ↑/↓
- Press 'e' to edit the extra bonus
- Type a number (can be negative for penalties)
- Press Enter to save

**Example**:
```
SET EXTRA BONUSES

▶ Strength    : +2  [from Dragonborn]
  Dexterity   : +0
  Constitution: +1  [from Dragonborn]
  Intelligence: +0
  Wisdom      : +0
  Charisma    : +1  [from Fey Touched feat]

  [CONFIRM]
```

## Final Display

In the Stats panel, you'll see:
- **Total Score**: Base + Extra
- **Modifier**: Calculated from total score
- **Saving Throw**: Modifier + proficiency (if proficient)

**Example**:
```
ABILITY SCORES

Strength      15    +2 (STR)  Save: +2
Dexterity     14    +2 (DEX)  Save: +4 ●
Constitution  13    +1 (CON)  Save: +1
Intelligence  12    +1 (INT)  Save: +1
Wisdom        10    +0 (WIS)  Save: +0
Charisma      8    -1 (CHA)  Save: -1

● = Proficient in saving throw
```

Where the score (15, 14, etc.) is your **total** (base + extra).

## Keyboard Shortcuts Summary

### In Stat Generator:

**Method Selection**:
- `↑/↓` or `k/j`: Navigate methods
- `Enter`: Select method
- `Esc`: Cancel

**4d6 / Standard Array Assignment**:
- `↑/↓` or `k/j`: Navigate abilities
- `1-6`: Assign corresponding rolled/array value
- `Enter`: Continue to extras
- `Esc`: Go back

**Point Buy**:
- `↑/↓` or `k/j`: Navigate abilities
- `+` or `=`: Increase selected ability (costs points)
- `-` or `_`: Decrease selected ability (refunds points)
- `Enter`: Continue to extras
- `Esc`: Go back
- **Note**: Number keys (1-6) don't work in Point Buy - use +/- instead!

**Setting Extras**:
- `↑/↓` or `k/j`: Navigate abilities
- `e`: Edit extra for selected ability
- `Enter`: Confirm (when on [CONFIRM] button) or save extra (when editing)
- `Esc`: Cancel extra edit or go back

### In Stats Panel:
- `e`: Open ability score generator

## Tips

1. **Species Bonuses**: Apply extras AFTER selecting your species in Character Info panel
2. **Modifier Calculation**: Modifier = (Score - 10) ÷ 2, rounded down
3. **Point Buy Strategy**: Most optimal to have 15, 15, 15, 8, 8, 8 or 15, 14, 14, 12, 10, 8
4. **Save Often**: Press 's' to save your character after making changes
5. **Reroll 4d6**: If you don't like your rolls, just press Esc to go back and start over

## Technical Details

The app stores ability scores in three parts:
- **Base**: Your rolled/assigned/bought score
- **Extra**: Bonuses from species, feats, etc.
- **Total**: Base + Extra (this is what's shown in the stats panel)

This separation allows you to:
- See where bonuses come from
- Easily recalculate if bonuses change
- Track character progression accurately
