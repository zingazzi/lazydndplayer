# LazyDnDPlayer - Quick Start Guide

## Installation & First Run

1. **Build the application:**
   ```bash
   go build -o lazydndplayer .
   ```

2. **Run with default settings:**
   ```bash
   ./lazydndplayer
   ```
   This creates a new character at `~/.lazydndplayer/character.json`

3. **Run with sample character:**
   ```bash
   ./lazydndplayer -file ./data/sample_character.json
   ```

## Basic Workflow

### Creating Your First Character

1. Launch the application
2. Press `1` to go to Overview panel
3. Press `e` to edit (basic implementation - manual JSON editing recommended for now)
4. Navigate panels with `Tab` or number keys `1-7`

### Managing Character Stats

**Stats Panel (Press `2`):**
- View all ability scores and modifiers
- See saving throw proficiencies
- Modifiers are calculated automatically

**Skills Panel (Press `3`):**
- Use `↑`/`↓` to navigate skills
- Press `e` to cycle proficiency: None → Proficient → Expertise
- Press `r` to roll a skill check (automatically adds modifiers)

### Combat & Actions

**Actions Panel (Press `6`):**
- View available actions, bonus actions, reactions
- Track limited-use abilities (Action Surge, Second Wind, etc.)
- Press `r` for long rest to restore uses

### Rolling Dice

**Dice Panel (Press `7`):**
- Press `Enter` to focus input field
- Type dice notation: `1d20`, `2d6+3`, `d20-1`
- Press `Enter` to roll
- Quick rolls: `1`=d4, `2`=d6, `3`=d8, `4`=d10, `5`=d12, `6`=d20, `7`=d100
- Change roll type:
  - `n` = Normal
  - `a` = Advantage
  - `d` = Disadvantage

### Inventory Management

**Inventory Panel (Press `4`):**
- Press `a` to add items
- Press `d` to delete selected item
- Press `e` to toggle equipped status
- Automatic weight calculation
- Visual warning when overloaded

### Spells

**Spells Panel (Press `5`):**
- View spells organized by level
- Track spell slots (● = available, ○ = used)
- Press `r` for long rest to restore all slots
- Press `a` to add spells

## Tips & Tricks

1. **Quick Navigation:** Use number keys (1-7) to jump directly to panels
2. **Help Anytime:** Press `?` to see all keyboard shortcuts
3. **Auto-Save:** Press `s` to manually save your character
4. **Level Up:** When you have enough XP, press `l` to level up
5. **Skill Checks:** In Skills panel, select a skill and press `r` for automatic roll with modifiers

## Editing Character Details

For detailed character editing, edit the JSON file directly:

```bash
# With your preferred editor
vim ~/.lazydndplayer/character.json
# or
code ~/.lazydndplayer/character.json
```

Key fields to edit:
- `name`, `race`, `class`, `background`, `alignment`
- `ability_scores` - change the numeric values
- `max_hp`, `armor_class`, `speed`
- Add/modify items in `inventory.items[]`
- Add/modify spells in `spellbook.spells[]`

## Import/Export Characters

**Export your character:**
```bash
./lazydndplayer -export ./my-backup.json
```

**Import from another file:**
```bash
./lazydndplayer -import ./downloaded-character.json
```

**Use a specific character file:**
```bash
./lazydndplayer -file ./characters/gandalf.json
```

## Common Tasks

### Adding a New Weapon
1. Edit the JSON file or use the inventory panel
2. Set `equipped: true` for equipped items
3. The weight is automatically tracked

### Tracking Spell Slots
1. Go to Spells panel (Press `5`)
2. Use spells during play
3. Press `r` to take a long rest and restore all slots

### Rolling Skill Checks
1. Go to Skills panel (Press `3`)
2. Navigate to the skill (e.g., Perception)
3. Press `r` to roll with all modifiers applied

### Taking Damage
1. Edit the JSON file to modify `current_hp`
2. Or calculate manually and update

### Leveling Up
1. Add experience points by editing `experience` in JSON
2. When you have enough XP, press `l` in the application
3. Follow the level-up wizard (basic implementation)

## Keyboard Reference

| Key | Action |
|-----|--------|
| `1-7` | Quick jump to panels |
| `Tab` | Next panel |
| `Shift+Tab` | Previous panel |
| `↑`/`↓` or `j`/`k` | Navigate lists |
| `s` | Save character |
| `l` | Level up |
| `r` | Roll/Rest (context-dependent) |
| `e` | Edit/Toggle (context-dependent) |
| `a` | Add item/spell/action |
| `d` | Delete item/action |
| `?` | Toggle help |
| `q` or `Ctrl+C` | Quit |

## Troubleshooting

**Character file not found:**
- The app creates a default character on first run
- Check `~/.lazydndplayer/character.json`

**Application won't start:**
- Ensure terminal supports colors and UTF-8
- Try resizing terminal window
- Make sure you built with `go build`

**Changes not saving:**
- Press `s` to manually save
- Check file permissions on character file
- Verify the file path is writable

## Advanced Usage

### Multiple Characters

Create separate JSON files for each character:
```bash
./lazydndplayer -file ./characters/fighter.json
./lazydndplayer -file ./characters/wizard.json
./lazydndplayer -file ./characters/rogue.json
```

### Custom Items/Spells

Edit the JSON file to add custom items or spells with any properties you need. The application will display them correctly.

### Sharing Characters

Simply share the JSON file - it's human-readable and portable:
```bash
# Email or share
./lazydndplayer -export ./share/my-character.json

# On another machine
./lazydndplayer -import ./downloads/my-character.json
```

## Getting Help

- Press `?` in the application for keyboard shortcuts
- Check README.md for detailed feature list
- Refer to the sample character in `data/sample_character.json`

## Next Steps

1. Customize the sample character or create your own
2. Try rolling some skill checks
3. Manage your inventory
4. Track your character through a session
5. Level up when you earn enough XP!

Happy adventuring! 🎲
