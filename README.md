# LazyDnDPlayer

A terminal user interface (TUI) application for managing D&D 5e 2024 characters, inspired by lazygit.

## Features

- **Full TUI Interface**: Beautiful, keyboard-driven terminal interface
- **Character Management**: Complete character sheet with all D&D 5e attributes
- **Ability Scores & Stats**: Track STR, DEX, CON, INT, WIS, CHA with automatic modifier calculation
- **Skills System**: All 18 D&D skills with proficiency and expertise tracking
- **Inventory Management**: Track items, equipment, and encumbrance
- **Spell Management**: Organize spells by level, track spell slots
- **Action Tracking**: Manage actions, bonus actions, and reactions
- **Dice Roller**: Built-in dice roller with advantage/disadvantage support
- **Level Up System**: Track XP and level progression
- **JSON Import/Export**: Easy character data portability

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
git clone <repository-url>
cd lazydndplayer
go build -o lazydndplayer .
```

## Usage

### Basic Usage

Run the application:

```bash
./lazydndplayer
```

This will create a default character file at `~/.lazydndplayer/character.json`.

### Command Line Options

```bash
# Use a specific character file
./lazydndplayer -file /path/to/character.json

# Import a character from another file
./lazydndplayer -import /path/to/import.json

# Export your character
./lazydndplayer -export /path/to/export.json
```

## Keyboard Shortcuts

### Global Navigation

- `↑/↓` or `j/k` - Navigate lists
- `Tab` / `Shift+Tab` - Switch tabs
- `1-5` - Quick jump to tab (1=Overview, 2=Stats, 3=Skills, 4=Inventory, 5=Spells)
- `s` - Save character
- `l` - Level up (when enough XP)
- `Shift+R` - Long rest (restore HP, spells, and abilities)
- `?` - Toggle help
- `q` or `Ctrl+C` - Quit

### Quick Dice Rolling (Always Available)

- `d` + `1-7` - Quick roll (d1=d4, d2=d6, d3=d8, d4=d10, d5=d12, d6=d20, d7=d100)
- Examples: Press `d6` for d20, `d2` for d6

### Panel-Specific Keys

#### Skills Panel
- `r` - Roll selected skill check (automatically adds modifiers)
- `e` - Toggle proficiency (Not Proficient → Proficient → Expertise)

#### Inventory Panel
- `a` - Add new item
- `e` - Toggle equipped status
- Note: Delete functionality available in JSON editing

#### Spells Panel
- `a` - Add new spell
- Note: Spell slots restored with Shift+R (long rest)

### Fixed Panels (Always Visible)

The **Actions** and **Dice Roller** panels are always visible at the bottom of the screen for quick access during gameplay.

## UI Layout

The interface uses a modern tab-based layout with fixed panels at the bottom:

```
┌────────────────────────────────────────────────────────┐
│ Title: Character Name, Class, Level                    │
├────────────────────────────────────────────────────────┤
│ [Overview] [Stats] [Skills] [Inventory] [Spells]      │
├────────────────────────────────────────────────────────┤
│                                                         │
│            Main Panel (Switchable via tabs)            │
│                                                         │
├───────────────────────────┬─────────────────────────────┤
│    Actions Panel          │     Dice Roller Panel       │
│    (Always Visible)       │     (Always Visible)        │
└───────────────────────────┴─────────────────────────────┘
```

### Main Panels (Switch with Tab or 1-5)

**1. Overview**
Displays character name, race, class, level, XP, HP, AC, and other core stats.

**2. Stats**
Shows all six ability scores with modifiers and saving throw proficiencies.

**3. Skills**
Lists all 18 D&D skills with calculated bonuses. Roll skill checks directly from here.

**4. Inventory**
Manage your character's items, track weight and encumbrance, mark items as equipped.

**5. Spells**
View and manage spells organized by level, track spell slots, mark spells as prepared.

### Fixed Panels (Always Visible at Bottom)

**Actions** (Bottom Left)
List available actions, bonus actions, and reactions. Track limited-use abilities.

**Dice Roller** (Bottom Right)
Roll dice using standard D&D notation (e.g., "2d6+3", "1d20"). Quick roll with `d` + number keys.

## Character Data Format

Character data is stored in JSON format for easy editing and portability:

```json
{
  "name": "Thorin Oakenshield",
  "race": "Dwarf",
  "class": "Fighter",
  "level": 5,
  "experience": 6500,
  "max_hp": 52,
  "current_hp": 52,
  "armor_class": 18,
  "ability_scores": {
    "strength": 16,
    "dexterity": 12,
    "constitution": 16,
    "intelligence": 10,
    "wisdom": 12,
    "charisma": 10
  }
}
```

## Development

### Project Structure

```
lazydndplayer/
├── main.go                    # Entry point
├── internal/
│   ├── models/               # D&D data models
│   ├── storage/              # JSON persistence
│   ├── dice/                 # Dice rolling engine
│   ├── leveling/             # Level-up logic
│   └── ui/                   # TUI components
│       ├── panels/           # UI panels
│       └── components/       # Reusable components
└── data/                     # Character data
```

### Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## Roadmap

- [ ] Full level-up wizard implementation
- [ ] Complete spell list database
- [ ] Class-specific features
- [ ] Character creation wizard
- [ ] Multiple character management
- [ ] Combat tracker
- [ ] Initiative roller
- [ ] Condition tracking
- [ ] Rest management (short/long rest)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Acknowledgments

- Inspired by [lazygit](https://github.com/jesseduffield/lazygit)
- Built with [Charm](https://charm.sh/) tools
- D&D 5e is a trademark of Wizards of the Coast
