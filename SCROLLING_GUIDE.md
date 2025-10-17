# Scrolling Guide

## Overview

Panels with content larger than the visible area now support **automatic scrolling**. This ensures you never lose access to your data, even on smaller terminal windows.

## Which Panels Support Scrolling?

✅ **Skills Panel** - Scrolls through all 18 D&D skills
✅ **Inventory Panel** - Scrolls through your items
✅ **Actions Panel** - Scrolls through available actions
✅ **Spells Panel** - (Fixed height, may need scrolling with many spells)

## How Scrolling Works

### Automatic Scrolling
When you navigate with `↑`/`↓` keys, the viewport automatically scrolls to keep the selected item visible:

```
SKILLS
━━━━━━━━━━━━━━━━━━
  Acrobatics (DEX)    +1
● Athletics (STR)     +6    ← Selected
  Arcana (INT)        +0
  Deception (CHA)     +0
... (more items below)
                    ↓ 25%   ← Scroll indicator
```

### Scroll Indicators

When content extends beyond the visible area, you'll see a scroll indicator:
- `↓ 25%` - Shows you're 25% through the content
- `↓ 75%` - Near the bottom
- No indicator - All content is visible

### Navigation Keys

| Panel | Keys | Behavior |
|-------|------|----------|
| **Skills** | `↑`/`↓` or `j`/`k` | Navigate and auto-scroll |
| **Inventory** | `↑`/`↓` or `j`/`k` | Navigate and auto-scroll |
| **Actions** | `↑`/`↓` or `j`/`k` | Navigate and auto-scroll |

## Examples

### Skills Panel with Many Items

```
Focus: Main Panel
Tab: Skills

SKILLS                    ← Title always visible
━━━━━━━━━━━━━━━━━━━━━━━━
  Acrobatics (DEX)    +1
● Animal Handling (WIS) +4  ← Currently selected
  Arcana (INT)        +0
  Athletics (STR)     +6
  Deception (CHA)     +0
  History (INT)       +0
  Insight (WIS)       +1
  (more below...)
                    ↓ 15%  ← Scroll indicator
```

Press `↓` to move down and scroll automatically.

### Inventory Panel with Many Items

```
INVENTORY
━━━━━━━━━━━━━━━━━━━━━━━━
Gold: 245  Silver: 18  Copper: 42

Carry Weight: 89.5 / 240.0 lbs

[E] Longsword           x1    3.0 lbs
[ ] Shield              x1    6.0 lbs
[E] Plate Armor         x1    65.0 lbs
[ ] Healing Potion      x3    1.5 lbs
[ ] Rope (50 ft)        x1    10.0 lbs
[ ] Torch               x5    5.0 lbs
... (more items)
                      ↓ 45%
```

### Actions Panel

```
ACTIONS
━━━━━━━━━━━━━━━━━━━━━━━━
Actions
  Attack
  Cast a Spell
  Dash
  Disengage
  Dodge

Bonus Actions
  Second Wind [1/1]

Reactions
  Opportunity Attack

↑/↓ Navigate • Enter Activate
```

If you have many custom actions, scroll indicators will appear.

## Visual Feedback

### Borders and Focus
- **Pink border** = Panel has focus
- **No border** = Panel visible but not focused
- Tabs are always visible at the top

### Layout
```
┌─────────────────────────────────────────┐
│ Title Bar                                │
├─────────────────────────────────────────┤
│ [Tab] [Tab] [Tab] [Tab] [Tab]          │ ← Always visible
├═════════════════════════════════════════┤
│                                          │ ← Pink border when
│         MAIN PANEL                       │   focused (FocusMain)
│         (Scrolls if needed)              │
│                      ↓ 30%               │ ← Scroll indicator
└═════════════════════════════════════════┘

┌──────────────────┐  ┌────────────────────┐
│ ACTIONS          │  │ DICE ROLLER        │
│ (Scrollable)     │  │ (Fixed height)     │
│            ↓ 10% │  │                    │
└──────────────────┘  └────────────────────┘
```

## Tips

### 💡 Use Focus to Your Advantage
1. Press `f` to cycle focus
2. When panel is focused, pink border appears
3. Use `↑`/`↓` to navigate and scroll

### 💡 Quick Navigation
```bash
# In Skills panel
f              # Focus Main (if not already)
3              # Go to Skills tab
↓ ↓ ↓          # Scroll down through skills
r              # Roll selected skill
```

### 💡 Check Scroll Position
- Look for the `↓ XX%` indicator
- No indicator = all content fits in view
- Indicator present = more content below

### 💡 Smooth Scrolling
- Each `↑`/`↓` moves one line
- Selection stays in view automatically
- No manual scrolling required

## Responsive Design

The scrolling system automatically adapts to:
- **Terminal size** - Smaller terminals show less content but still work
- **Content amount** - More items = automatic scrolling
- **Focus state** - Border appears without breaking layout

### Small Terminal Example
```
Terminal: 80x24

Main Panel gets:
- Title: 1 line
- Tab bar: 3 lines
- Status: 1 line
- Content: ~10 lines (viewport)
- Bottom panels: 20 lines

If you have 18 skills, you'll see ~10 at a time
and can scroll to see the rest.
```

### Large Terminal Example
```
Terminal: 120x40

Main Panel gets:
- More vertical space
- May fit all content without scrolling
- Scroll indicators won't appear if not needed
```

## Troubleshooting

**Q: I don't see the scroll indicator?**
- All content fits in the visible area
- This is normal and expected!

**Q: Content cuts off mid-line?**
- Viewport automatically handles line breaks
- Use `↑`/`↓` to see more

**Q: Tabs disappeared?**
- Check if focus border is taking too much space
- Try pressing `f` to cycle focus areas
- Tabs should always be visible

**Q: Can I scroll with mouse?**
- Currently keyboard-only navigation
- Use `↑`/`↓` or `j`/`k` keys

**Q: Content doesn't scroll?**
- Make sure the panel has focus (press `f`)
- Try resizing terminal window
- Check if you're in the correct tab

## Technical Details

### Implementation
- Uses `viewport` component from Bubbles
- Automatic content sizing
- Scroll position tracked per panel
- Smooth line-by-line scrolling

### Performance
- Efficient rendering
- Only visible content is drawn
- No performance impact from large lists

### Limitations
- Currently line-by-line scrolling only
- No page-up/page-down (yet)
- No mouse wheel support (yet)

## Keyboard Summary

| Key | In Main Focus | In Actions Focus | In Dice Focus |
|-----|---------------|------------------|---------------|
| `↑`/`k` | Scroll up | Scroll up | - |
| `↓`/`j` | Scroll down | Scroll down | - |
| `f` | Cycle focus | Cycle focus | Cycle focus |

---

**Scrolling makes the interface work on any terminal size!** 📜✨
