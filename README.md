# LazyDnDPlayer

A terminal user interface (TUI) application for managing D&D 5e 2024 characters, inspired by lazygit.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

### Core Features
- 🎮 **Full TUI Interface** - Beautiful keyboard-driven terminal interface
- 📊 **Complete Character Sheet** - All D&D 5e 2024 attributes and stats
- 🎲 **Advanced Dice Roller** - Support for `1d20+3`, `2d6`, advantage/disadvantage, complex rolls
- 🧙 **Species System** - All 14 D&D 5e 2024 species with traits, resistances, and darkvision
- 📦 **JSON Storage** - Easy character import/export and data management

### Character Management
- **Ability Scores** - Track STR, DEX, CON, INT, WIS, CHA with automatic modifiers
- **Skills** - All 18 D&D skills with proficiency and expertise tracking
- **Inventory** - Track items, equipment, weight, and encumbrance
- **Spells** - Organize spells by level, track slots and prepared spells
- **Features** - Manage limited-use abilities (class features, species abilities) with rest recovery
- **Traits** - Languages, feats, resistances, darkvision, and species traits

### Species Support
Aasimar, Dragonborn, Dwarf, Elf (Drow/High/Wood), Gnome, Goliath, Halfling, Human, Orc, Tiefling (Abyssal/Chthonic/Infernal)

## Quick Start

### Installation

**Prerequisites:** Go 1.21 or higher

```bash
# Clone and build
git clone <repository-url>
cd lazydndplayer

# Using Make (recommended)
make build

# Or using Go directly
go build -o lazydndplayer .

# Run
./lazydndplayer
```

### First Time Setup

On first run, a default character file is created at `~/.lazydndplayer/character.json`.

### Command Line Options

```bash
# Use specific character file
./lazydndplayer -file /path/to/character.json

# Import character
./lazydndplayer -import /path/to/backup.json

# Export character
./lazydndplayer -export /path/to/backup.json
```

## User Interface

### Layout

```
┌─────────────────────────────────────────────────────────────────┐
│ Main Panel (Tabs: Stats/Skills/Inventory/Spells/Features/      │
│                   Traits)                                       │
│                                                                 │
│ [55% width]                          Character Stats [45%]     │
│                                      - Name, Species, Level     │
│                                      - HP, AC, Initiative       │
│                                      - Speed, Proficiency       │
├─────────────────────────────────────┬───────────────────────────┤
│ Actions Panel [50%]                 │ Dice Roller [50%]        │
│ - Quick action reference            │ - Input field            │
│ - Combat actions                    │ - Roll history           │
└─────────────────────────────────────┴───────────────────────────┘
│ Status Bar: App • Help • Navigation                            │
└─────────────────────────────────────────────────────────────────┘
```

### Keyboard Shortcuts

#### Global
- `q` - Quit application
- `?` - Toggle help
- `f` - Cycle focus between panels
- `Tab` - Switch tabs (in main panel)
- `1-6` - Quick jump to tab (when main panel focused)
- `ctrl+s` - Save character

#### Navigation
- `↑/↓` or `j/k` - Navigate lists
- `←/→` or `h/l` - Navigate horizontally (where applicable)
- `Esc` - Cancel/close popup

#### Main Panel (Stats/Skills/Inventory/Spells/Features/Traits)
- `e` - Edit selected item
- `a` - Add new item
- `d` - Delete selected item
- `space` - Toggle (proficiency, equipped, prepared, etc.)
- `u` - Use feature (Features tab only)
- `+/=` - Restore feature charge (Features tab only)
- `r` - Short rest (Features/Spells tabs)
- `Shift+R` - Long rest (Features/Spells tabs)

#### Character Stats Panel
- `n` - Edit name
- `r` - Change species
- `+/-` - Add/remove HP
- `i` - Roll initiative

#### Dice Roller
- `Enter` - Start typing dice expression (input mode)
- `h` - View history mode
- `r` - Reroll last dice
- `↑/↓` - Navigate history (in history mode)
- `Esc` - Exit input/history mode

##### Dice Notation
- Basic: `1d20`, `2d6+3`
- Advantage/Disadvantage: `1d20 adv`, `1d20 dis`
- Complex: `2d8+3d4+2`
- Multiple rolls: `1d20+3, 2d6, 1d4` (comma-separated)

## Species System

### Selecting a Species

1. Focus on Character Stats panel (`f` key)
2. Press `r` to open species selector
3. Navigate with `↑/↓`
4. Press `Enter` to select
5. If prompted, select additional language
6. If prompted, select skill proficiency

### Automatic Features

When you select a species:
- ✅ Speed updated
- ✅ Languages applied
- ✅ Resistances applied
- ✅ Darkvision set
- ✅ Species traits added
- ✅ Automatic skill proficiencies applied (e.g., Elf → Perception)
- ✅ Old species proficiencies removed

### Adding Custom Species

Edit `data/species.json`:

```json
{
  "name": "Your Species",
  "size": "Medium",
  "speed": 30,
  "description": "Species description",
  "traits": [
    {
      "name": "Trait Name",
      "description": "Trait description"
    }
  ],
  "languages": ["Common", "Other Language"],
  "resistances": ["Fire"],
  "darkvision": 60
}
```

See `data/README.md` for full documentation.

## Data Management

### Character Files

Default location: `~/.lazydndplayer/character.json`

#### Backup Character
```bash
cp ~/.lazydndplayer/character.json ~/backup.json
```

#### Load Backup
```bash
./lazydndplayer -file ~/backup.json
```

### Species Data

Species definitions: `data/species.json`

Edit this file to add or modify species. Changes take effect on next application start (no recompilation needed).

## Development

### Project Structure

```
lazydndplayer/
├── main.go              # Entry point
├── data/
│   └── species.json     # Species data
├── internal/
│   ├── dice/            # Dice rolling engine
│   ├── leveling/        # Level up system
│   ├── models/          # Data models
│   ├── storage/         # JSON persistence
│   └── ui/              # TUI components
│       ├── app.go       # Main application
│       ├── components/  # Reusable UI components
│       └── panels/      # Panel implementations
└── README.md
```

### Building

```bash
# Build
go build -o lazydndplayer .

# Run tests
go test ./...

# Clean build
rm lazydndplayer
go build -o lazydndplayer .
```

### Adding Features

1. **New Species**: Edit `data/species.json`
2. **New Panel**: Create file in `internal/ui/panels/`
3. **New Component**: Create file in `internal/ui/components/`
4. **Data Models**: Add to `internal/models/`

## Troubleshooting

### Character File Issues

**Reset to default:**
```bash
rm ~/.lazydndplayer/character.json
./lazydndplayer
```

### Display Issues

Ensure your terminal:
- Supports 256 colors
- Has minimum size of 120x30
- Uses a modern terminal emulator (iTerm2, Alacritty, Windows Terminal, etc.)

### Species Not Loading

Check `data/species.json` for valid JSON syntax:
```bash
cat data/species.json | jq .
```

## Changelog

See `CHANGELOG.md` for version history and updates.

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

Inspired by [lazygit](https://github.com/jesseduffield/lazygit).
