# Quick Start - LazyDnDPlayer v2.0

## What's New in v2.0?

🎨 **Complete UI Redesign!**

- **Tab navigation** instead of sidebar
- **Dice roller always visible** at bottom right
- **Actions always visible** at bottom left
- **Quick dice shortcuts**: `d6` = roll d20, `d2` = roll d6, etc.
- **Better layout**: More space for your character data

## Running the App

### Option 1: Quick Start (Recommended)
```bash
./RUN_ME.sh
```

### Option 2: Default Character
```bash
./lazydndplayer
```
Creates a new character at `~/.lazydndplayer/character.json`

### Option 3: Sample Character
```bash
./lazydndplayer -file ./data/sample_character.json
```
Use the pre-made Level 5 Fighter (Thorin Oakenshield)

## First Time User?

### 1. Start the App
```bash
./RUN_ME.sh
```

### 2. Learn the Layout
- **Top**: Character name and title bar
- **Second row**: Tab bar (Overview | Stats | Skills | Inventory | Spells)
- **Middle**: Main panel (changes based on selected tab)
- **Bottom Left**: Actions panel (always visible)
- **Bottom Right**: Dice roller (always visible)

### 3. Try These Keys
1. Press `?` - See all keyboard shortcuts
2. Press `2` - Jump to Stats tab
3. Press `3` - Jump to Skills tab
4. Press `d6` - Roll a d20 (watch the dice panel!)
5. Press `Tab` - Cycle through tabs
6. Press `r` in Skills - Roll a skill check
7. Press `s` - Save your character
8. Press `q` - Quit

## Essential Shortcuts

| Key | What It Does |
|-----|--------------|
| `1-5` | Jump to tabs (Overview, Stats, Skills, Inventory, Spells) |
| `Tab` | Next tab |
| `d6` | Roll d20 (works from anywhere!) |
| `d2` | Roll d6 |
| `r` | Roll selected skill check (in Skills tab) |
| `e` | Edit/Toggle (context-dependent) |
| `s` | Save character |
| `Shift+R` | Long rest (restore everything) |
| `?` | Help overlay |
| `q` | Quit |

## Common Tasks

### Roll a Skill Check
1. Press `3` to go to Skills tab
2. Use `↑`/`↓` to select a skill
3. Press `r` to roll with modifiers

### Roll Some Dice
- Press `d` then a number:
  - `d1` = d4
  - `d2` = d6
  - `d3` = d8
  - `d4` = d10
  - `d5` = d12
  - `d6` = d20 ⭐
  - `d7` = d100

### Manage Inventory
1. Press `4` for Inventory tab
2. Press `a` to add items
3. Use `↑`/`↓` to select items
4. Press `e` to toggle equipped

### Take a Long Rest
- Press `Shift+R` from anywhere
- Restores: HP, Spell Slots, Action uses

### Save Your Character
- Press `s` at any time
- Auto-saves to your character file

## Visual Guide

```
┌─────────────────────────────────────────────────────────┐
│ LazyDnDPlayer - Thorin (Fighter) - Level 5              │ ← Title
├─────────────────────────────────────────────────────────┤
│ [Overview] [Stats] [Skills] [Inventory] [Spells]       │ ← Tabs
├─────────────────────────────────────────────────────────┤
│                                                          │
│  SKILLS                          ← Main Panel           │
│  ● Athletics (STR)      +6       (Changes with tabs)    │
│  ● Perception (WIS)     +4                              │
│    Stealth (DEX)        +1                              │
│                                                          │
├──────────────────────────┬───────────────────────────────┤
│ ACTIONS                  │ DICE ROLLER                   │
│ • Attack                 │ Roll type: Normal             │
│ • Second Wind [1/1]      │ [d + 1-7 for quick]          │
│ • Action Surge [1/1]     │ History:                     │
│ • Dash                   │ • d20: [18] = 18             │
└──────────────────────────┴───────────────────────────────┘
  Status: Saved                                ← Status bar
```

## Pro Tips

💡 **During Combat**
- Keep Actions panel visible to see what you can do
- Use `d6` repeatedly to roll attacks
- Use `r` in Skills for ability checks

💡 **Character Building**
- Press `3` for Skills, then `e` to toggle proficiency
- Edit JSON file directly for detailed character creation
- Use `l` when you have enough XP to level up

💡 **Workflow**
- Dice rolls show in history automatically
- Critical hits (nat 20) and fails (nat 1) are highlighted
- Status bar shows your last action

## Troubleshooting

**Q: Application won't start?**
- Make sure you're in a terminal that supports colors
- Try: `export TERM=xterm-256color`

**Q: Want to create a new character?**
```bash
./lazydndplayer -file ./my-new-character.json
```

**Q: Want to edit character details?**
- Edit the JSON file with your favorite editor
- Or use the panels in the app

**Q: How do I share my character?**
```bash
./lazydndplayer -export ./shared-character.json
# Send the JSON file to others
```

## Next Steps

1. ✅ Run the app: `./RUN_ME.sh`
2. ✅ Learn the shortcuts: Press `?`
3. ✅ Try rolling some dice: `d6`
4. ✅ Navigate tabs: `1-5`
5. ✅ Read the full manual: `README.md`

## Need Help?

- Press `?` in the app for shortcuts
- Read `README.md` for detailed features
- Check `UI_REDESIGN.md` for v2.0 changes
- See `USAGE.md` for advanced tips

---

**Have fun adventuring!** 🎲✨
