# Extra Bonuses Navigation Guide

## Overview

The Extra Bonuses screen now has improved navigation and easier value adjustment using +/- keys.

## Two Ways to Access Extras

### Method 1: Direct Access with 'e' ✨
**From Stats Panel → Press `e`**

```
Stats Panel
    ↓ press 'e'
Extra Bonuses Screen
    ↓ press 'esc'
Stats Panel (back to where you started)
```

**Use this when:**
- You just want to add/change modifiers
- You already have base stats set
- Adding species/feat bonuses

### Method 2: Full Flow with 'r' 🎲
**From Stats Panel → Press `r`**

```
Stats Panel
    ↓ press 'r'
Method Selection (4d6, Standard Array, Point Buy, Custom)
    ↓ select method & press enter
Stats Assignment (roll/assign/adjust values)
    ↓ press enter
Extra Bonuses Screen
    ↓ press 'esc'
Stats Assignment (go back to adjust base values)
    ↓ press 'esc'
Method Selection
    ↓ press 'esc'
Stats Panel
```

**Use this when:**
- Rolling new stats
- Starting a new character
- Want to change base values AND modifiers

## Editing Extra Bonuses

### Quick Adjust with +/- Keys (NEW! ⭐)

Navigate to an ability and use:
- **`+`** or **`=`** → Increase by 1
- **`-`** or **`_`** → Decrease by 1

```
Strength    : +2  ← selected
Dexterity   : +0
Constitution: +1

Press '+' → Strength becomes +3
Press '-' → Strength becomes +1
```

**Range**: -5 to +10
- Min: -5 (for penalties/debuffs)
- Max: +10 (for powerful bonuses)

### Type Exact Value with 'e' Key

For setting specific values quickly:

1. Navigate to ability with `↑/↓`
2. Press `e` to enter edit mode
3. Type the number (e.g., `2` or `-1`)
4. Press `Enter` to save

```
Strength    : +2█  ← typing mode
Type: 3
Press Enter → Saves as +3
```

## Navigation Map

### In Extra Bonuses Screen

**Not Editing:**
| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate abilities |
| `+` or `=` | Increase selected ability by 1 |
| `-` or `_` | Decrease selected ability by 1 |
| `e` | Enter typing mode for exact value |
| `Enter` | Confirm (on CONFIRM button) |
| `Esc` | Go back (depends on how you got here) |

**While Typing (after pressing `e`):**
| Key | Action |
|-----|--------|
| `0-9` | Type number |
| `+/-` | Type sign |
| `Backspace` | Delete character |
| `Enter` | Save value |
| `Esc` | Cancel typing |

## Escape Key Behavior

The `Esc` key is **context-aware**:

### If you pressed 'e' to get here:
```
Stats Panel → (e) → Extras
                    ↓ esc
               Stats Panel ✅
```
Shows: "Esc: Back to Stats Panel"

### If you pressed 'r' to get here:
```
Stats Panel → (r) → Method → Assignment → Extras
                                          ↓ esc
                              Assignment ✅
```
Shows: "Esc: Back to Assignment"

## Examples

### Example 1: Adding Dragonborn Bonuses

**Using Quick Adjust (+/-):**

1. Press `e` in Stats panel
2. Navigate to Strength (already selected)
3. Press `+` twice → +2 for Dragonborn
4. Press `↓` twice to Constitution
5. Press `+` once → +1 for Dragonborn
6. Navigate to CONFIRM and press `Enter`
7. Done! ✅

**Using Type Mode (e):**

1. Press `e` in Stats panel
2. On Strength, press `e` to type
3. Type `2` and press `Enter`
4. Press `↓` twice to Constitution
5. Press `e` to type
6. Type `1` and press `Enter`
7. Navigate to CONFIRM and press `Enter`
8. Done! ✅

### Example 2: Rolling Stats with Bonuses

1. Press `r` in Stats panel
2. Choose "4d6 Drop Lowest"
3. Assign all 6 rolled values using 1-6 keys
4. Press `Enter` → Goes to Extras
5. Add species bonuses with `+/-` keys
6. Press `Esc` → Goes back to assignment (if you want to change base values)
7. Or navigate to CONFIRM and press `Enter` → Apply everything! ✅

### Example 3: Negative Modifiers (Curses/Debuffs)

1. Press `e` in Stats panel
2. Navigate to affected ability
3. Press `-` to decrease (can go to -5)
4. Or press `e` and type `-2` for exact value
5. Navigate to CONFIRM and press `Enter`

## Visual Indicators

### Extra Bonuses Screen

```
SET EXTRA BONUSES

Add bonuses from species, feats, or other sources (range: -5 to +10):

▶ Strength    : +2   ← Selected (use +/- to adjust)
  Dexterity   : +0
  Constitution: +1
  Intelligence: +0
  Wisdom      : +0
  Charisma    : +0

  [CONFIRM]

↑/↓: Navigate  +/-: Adjust  e: Type Value  Enter: Confirm  Esc: Back to Stats Panel
```

### While Typing

```
▶ Strength    : +2█  ← Cursor blinking
  Dexterity   : +0
  ...

Type number  Enter: Save  Esc: Cancel
```

## Tips & Tricks

✅ **DO:**
- Use `+/-` for quick adjustments (most common use case)
- Use `e` for exact values (when you know the number)
- Check the bottom instruction to see where `Esc` goes
- Use `+/-` even in typing mode by pressing `Esc` first

❌ **DON'T:**
- Type letters (only numbers and +/- work)
- Go below -5 or above +10 (enforced limits)
- Forget to CONFIRM (changes aren't applied until you confirm)

## Common Workflows

### Just Adding Species Bonuses
```
Stats panel → e → +/- to adjust → CONFIRM → Done!
```
**Fastest**: 5 key presses total

### Rolling Fresh Stats
```
Stats panel → r → Choose method → Assign → +/- for bonuses → CONFIRM
```

### Changing Base AND Bonuses
```
Stats panel → r → Choose Custom → +/- to set base → Enter → +/- for bonuses → CONFIRM
```

### Fixing a Mistake
```
In Extras → Esc → Back to assignment → Fix base values → Enter → Adjust bonuses → CONFIRM
```

## Keyboard Summary

| Context | Key | Action |
|---------|-----|--------|
| Stats Panel | `e` | Direct to extras |
| Stats Panel | `r` | Full generator |
| Extras (not typing) | `+/-` | Adjust value ⭐ |
| Extras (not typing) | `e` | Type exact value |
| Extras (typing) | `0-9` | Type digits |
| Extras (typing) | `+/-` | Type sign |
| Extras (typing) | `Enter` | Save |
| Extras (typing) | `Esc` | Cancel typing |
| Extras (not typing) | `Esc` | Smart back (context-aware) |
| Extras (not typing) | `Enter` | Confirm (on button) |

## Troubleshooting

**Q: +/- keys don't work!**
A: Make sure you're not in typing mode. Press `Esc` if you see a cursor (█).

**Q: Where does Esc take me?**
A: Check the bottom instruction text. It will say either "Back to Stats Panel" or "Back to Assignment".

**Q: Can I have negative bonuses?**
A: Yes! Press `-` multiple times or type a negative number like `-2`. Range is -5 to +10.

**Q: How do I quickly set multiple bonuses?**
A: Use `↑/↓` to navigate and `+/-` to adjust each one. Much faster than typing!

**Q: I pressed Enter but nothing happened!**
A: Make sure you're on the `[CONFIRM]` button (use `↓` to get there) before pressing Enter.
