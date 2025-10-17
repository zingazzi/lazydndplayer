# UI Redesign - Version 2.0

## Overview

The UI has been completely redesigned to provide a better gameplay experience with always-accessible dice rolling and actions.

## What Changed

### Before (v1.0)
- Sidebar navigation with 7 switchable panels
- Had to switch to Dice panel to roll dice
- Had to switch to Actions panel to see available actions
- More navigation required during gameplay

### After (v2.0)
- Tab-based navigation with 5 main panels
- **Dice roller always visible** at bottom right
- **Actions panel always visible** at bottom left
- Less navigation, more focus on gameplay

## New Layout

```
┌─────────────────────────────────────────────────────────────┐
│  LazyDnDPlayer - Thorin Oakenshield (Fighter) - Level 5     │
├─────────────────────────────────────────────────────────────┤
│ [Overview] [Stats] [Skills] [Inventory] [Spells]  [Hints]  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│                    MAIN PANEL AREA                          │
│              (Full width, switchable content)               │
│                                                              │
│                                                              │
├──────────────────────────────┬───────────────────────────────┤
│     ACTIONS PANEL            │      DICE ROLLER             │
│     (Fixed - Always visible) │   (Fixed - Always visible)   │
│                              │                              │
│ • Attack                     │  Roll type: Normal           │
│ • Second Wind [1/1]          │  [d + 1-7 for quick roll]    │
│ • Action Surge [1/1]         │                              │
│ • Extra Attack               │  History:                    │
│                              │  d20: 18 = 18                │
└──────────────────────────────┴───────────────────────────────┘
│ Status: Ready                                                │
└──────────────────────────────────────────────────────────────┘
```

## Key Improvements

### 1. Always-Visible Dice Roller
- **Why**: Dice rolling is core to D&D gameplay
- **Benefit**: Roll dice anytime with `d` + number shortcuts
- **Examples**:
  - Press `d6` to roll a d20
  - Press `d2` to roll a d6
  - Roll history always visible

### 2. Always-Visible Actions
- **Why**: Quick reference during combat
- **Benefit**: See available actions without switching panels
- **Features**:
  - Shows use counters (e.g., "Second Wind [1/1]")
  - Updated in real-time
  - Shows all action types

### 3. Tab Navigation
- **Why**: Modern, clean interface
- **Benefit**: Less visual clutter, clear current location
- **Features**:
  - Quick shortcuts (1-5) visible in tab bar
  - Active tab highlighted
  - Keyboard hints included

### 4. More Screen Space
- **Main panel uses full width**
- **No sidebar taking up space**
- **Better use of vertical space**

## Keyboard Shortcuts Comparison

### Navigation

| Action | v1.0 | v2.0 |
|--------|------|------|
| Switch panels | Tab / ←→ / h/l | Tab / Shift+Tab |
| Quick select | 1-7 | 1-5 |
| Dice roller | Switch to panel 7 | Always visible |
| Actions | Switch to panel 6 | Always visible |

### New Shortcuts in v2.0

| Shortcut | Action |
|----------|--------|
| `d1` | Roll d4 |
| `d2` | Roll d6 |
| `d3` | Roll d8 |
| `d4` | Roll d10 |
| `d5` | Roll d12 |
| `d6` | Roll d20 |
| `d7` | Roll d100 |
| `Shift+R` | Long rest (global) |

## Technical Changes

### Files Modified
- `internal/ui/app.go` - Complete layout rewrite
- `internal/ui/components/tabs.go` - New tab component
- `internal/ui/components/help.go` - Updated shortcuts
- `internal/ui/panels/dice.go` - Improved for fixed display

### Files Removed
- No longer using `internal/ui/components/sidebar.go` for navigation

### Key Code Changes

1. **Panel Types Reduced**
   ```go
   // Only 5 main panels now
   const (
       OverviewPanel
       StatsPanel
       SkillsPanel
       InventoryPanel
       SpellsPanel
   )
   ```

2. **Fixed vs Switchable Panels**
   ```go
   // Main Panels (switchable)
   overviewPanel  *panels.OverviewPanel
   statsPanel     *panels.StatsPanel
   // ... etc

   // Fixed Panels (always visible)
   actionsPanel   *panels.ActionsPanel
   dicePanel      *panels.DicePanel
   ```

3. **New Layout Rendering**
   ```go
   // Three-section vertical layout
   - Title bar
   - Tab bar
   - Main panel (full width)
   - Bottom row (Actions | Dice)
   - Status bar
   ```

## Benefits for Gameplay

### During Combat
✅ Roll initiative without leaving overview
✅ See all available actions
✅ Quick dice rolls for attacks and damage
✅ No panel switching needed

### During Exploration
✅ Roll skill checks from Skills panel
✅ Roll any dice quickly
✅ Manage inventory while seeing actions
✅ Less distraction

### Character Management
✅ More screen space for spells/inventory
✅ Clear tab indication of current location
✅ Streamlined navigation
✅ Better visual hierarchy

## Migration Notes

If you're upgrading from v1.0:

1. **Muscle Memory**:
   - `6` and `7` no longer switch panels
   - Use `d` + number for dice instead

2. **New Shortcuts**:
   - Learn `d6` for d20 rolls
   - Use `Shift+R` for long rest

3. **Visual Changes**:
   - Look at tabs instead of sidebar
   - Actions and dice always at bottom

## Future Enhancements

Potential improvements for v2.1:
- [ ] Resizable panel heights
- [ ] Customizable bottom panel split
- [ ] Minimize bottom panels for more space
- [ ] Dice roll advantage/disadvantage toggle in UI
- [ ] Action quick-use buttons

## Feedback

The new layout prioritizes:
1. **Speed** - Less navigation
2. **Clarity** - Always see what you need
3. **Immersion** - Focus on gameplay, not UI

This design is inspired by modern gaming interfaces where critical information stays visible while navigating menus.
