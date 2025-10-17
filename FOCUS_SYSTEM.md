# Focus Management System

## Overview

The UI now has a **focus management system** that allows you to cycle between three different areas:
1. **Main Panel** (tabs + content area)
2. **Actions Panel** (bottom left)
3. **Dice Roller Panel** (bottom right)

## How It Works

### Switching Focus

Press **`f`** to cycle through focus areas:
- **Main Panel** → **Actions Panel** → **Dice Roller** → **Main Panel** (loops)

The currently focused panel is highlighted with a **pink border**.

### Visual Indicators

```
┌─────────────────────────────────────────────┐
│ [Overview] [Stats] [Skills] ...             │
├═════════════════════════════════════════════┤ ← Pink border
│                                              │    when focused
│         MAIN PANEL (Focused)                │
│                                              │
└═════════════════════════════════════════════┘

┌──────────────────┐  ┌──────────────────────┐
│ ACTIONS          │  │ DICE ROLLER          │
│ (Not focused)    │  │ (Not focused)        │
└──────────────────┘  └──────────────────────┘
```

## Behavior by Focus Area

### When Main Panel is Focused

**Available Keys:**
- `Tab` / `Shift+Tab` - Switch between tabs
- `1-5` - Quick jump to specific tabs
- `↑`/`↓` or `j`/`k` - Navigate lists in current tab
- `r` - Roll skill check (in Skills tab)
- `e` - Edit/Toggle items
- `a` - Add items
- `d + 1-7` - Quick dice rolls (global shortcut)

**Example:**
```
You're in the Skills tab, scrolling through skills.
Press 'r' to roll a skill check.
Press '3' to jump to Skills tab.
Press 'Tab' to move to Inventory tab.
```

### When Actions Panel is Focused

**Available Keys:**
- `↑`/`↓` or `j`/`k` - Navigate through actions
- `Enter` - Activate selected action
- `d + 1-7` - Still works for quick dice rolls

**Example:**
```
Press 'f' to focus Actions panel (pink border appears).
Use arrow keys to select "Second Wind".
Press Enter to use it.
```

### When Dice Roller is Focused

**Available Keys:**
- `1-7` - Quick roll dice (d4, d6, d8, d10, d12, d20, d100)
- `n` - Set roll type to Normal
- `a` - Set roll type to Advantage
- `d` - Set roll type to Disadvantage
- `Enter` - Focus input field for custom dice notation

**Example:**
```
Press 'f' twice to focus Dice Roller (pink border appears).
Press '6' to roll a d20 immediately.
Press 'a' to set advantage mode.
Press '6' again to roll d20 with advantage.
Press 'n' to return to normal mode.
Press 'Enter' to type custom dice like "2d6+3".
```

## Common Workflows

### Combat Turn
1. **Stay in Main Panel focus** (default)
2. Check your stats/skills as needed with tabs
3. Use `d6` for quick d20 rolls (attacks, saves)
4. Press `f` to focus Dice Roller for multiple rolls
5. Press `f` again to focus Actions to see abilities

### Skill Check
1. **Stay in Main Panel focus**
2. Press `3` to go to Skills tab
3. Navigate to desired skill with `↑`/`↓`
4. Press `r` to roll with modifiers

### Multiple Dice Rolls
1. Press `f` twice to **focus Dice Roller**
2. Rapidly press number keys for quick rolls
3. `6` `6` `2` `2` = d20, d20, d6, d6
4. Watch results in the history

### Check Available Actions
1. Press `f` to **focus Actions panel**
2. Navigate with `↑`/`↓` to see what's available
3. Check remaining uses (e.g., "Second Wind [1/1]")
4. Press `f` to return to Main Panel

## Keyboard Shortcuts Summary

| Key | Global | Main Focus | Actions Focus | Dice Focus |
|-----|--------|------------|---------------|------------|
| `f` | Cycle focus | ✓ | ✓ | ✓ |
| `Tab` | - | Switch tabs | - | - |
| `1-5` | - | Jump to tabs | Roll dice | Roll dice |
| `6-7` | - | - | Roll dice | Roll dice |
| `↑`/`↓` | - | Navigate lists | Navigate actions | - |
| `r` | - | Roll skill | - | - |
| `Enter` | - | Activate | Use action | Dice input |
| `n/a/d` | - | - | - | Roll type |
| `d+1-7` | Quick roll | ✓ | ✓ | ✓ |
| `s` | Save | ✓ | ✓ | ✓ |
| `?` | Help | ✓ | ✓ | ✓ |
| `q` | Quit | ✓ | ✓ | ✓ |

## Tips

### 💡 Stay in Main Most of the Time
- Main panel focus is where you'll spend most time
- Only switch focus when you need specific interaction with Actions or Dice

### 💡 Quick Dice Always Works
- `d + number` works from **any focus area**
- Example: `d6` rolls a d20 whether you're in Main, Actions, or Dice focus

### 💡 Visual Feedback
- Always check the **pink border** to know where you are
- Status bar shows "Focus: [Area Name]" when you switch

### 💡 Fast Combat Workflow
```
Main Focus:
  d6 → roll attack
  d2 → roll damage
  f → switch to Actions

Actions Focus:
  ↑/↓ → check abilities
  f → back to Main
```

### 💡 Rapid Dice Rolling
```
f → f → (now in Dice focus)
6 → d20
6 → d20
6 → d20
2 → d6
2 → d6
```

## Troubleshooting

**Q: Tab key isn't working!**
- You're probably in Actions or Dice focus
- Press `f` to cycle back to Main Panel focus
- Tab only works when Main Panel is focused

**Q: Number keys 1-5 aren't switching tabs!**
- Same issue - you're not in Main Panel focus
- Press `f` until you see the pink border on the main area

**Q: I can't see which panel is focused!**
- Look for the **pink/magenta border** around the panel
- Check the status bar message after pressing `f`

**Q: How do I get back to normal navigation?**
- Press `f` until "Focus: Main Panel" appears
- The top main area should have the pink border

## Design Rationale

### Why Three Focus Areas?

1. **Main Panel**: Primary navigation and data viewing
2. **Actions Panel**: Quick reference without switching tabs
3. **Dice Roller**: Rapid dice rolling without disrupting flow

### Why 'f' Key?

- **F for Focus** - easy to remember
- Single key press, no modifiers needed
- Not commonly used for other commands
- Easy to reach on all keyboard layouts

### Why Cycle Instead of Direct Selection?

- Simpler mental model (one key to press)
- Only 3 areas, so cycling is fast
- Reduces cognitive load during gameplay
- Prevents accidental focus changes

## Advanced Usage

### Power User Combo
```bash
# Start turn in Main focus
r              # Roll initiative (if in Skills)
f f            # Jump to Dice focus
6 6            # Roll attacks
d2 d2          # Quick damage rolls
f              # Check Actions
↑ ↓            # Review abilities
f              # Back to Main
```

### DM Mode (quick NPC actions)
```bash
f f            # Dice focus
6              # Attack roll
2              # Damage
6              # Another attack
3              # More damage
n              # Normal mode
a              # Advantage for next
6              # Roll with advantage
```

## Future Enhancements

Potential improvements:
- [ ] Remember last focused panel between sessions
- [ ] Configurable focus key
- [ ] Different border colors for each area
- [ ] Focus indicator in status bar always visible
- [ ] Direct focus keys (Alt+1, Alt+2, Alt+3)

---

The focus system makes the UI much more keyboard-friendly while maintaining simplicity!
