# Class Progression System

This directory contains individual JSON files for each D&D class, with complete level-by-level progression data from levels 1-20.

## Directory Structure

```
data/classes/
├── README.md (this file)
├── barbarian.json
├── bard.json
├── cleric.json
├── druid.json
├── fighter.json
├── monk.json
├── paladin.json
├── ranger.json
├── rogue.json
├── sorcerer.json
├── warlock.json
└── wizard.json
```

## JSON Schema

Each class file follows this structure:

```json
{
  "name": "Class Name",
  "description": "Brief class description",
  "hit_die": 8,
  "primary_ability": "Ability Name",
  "saving_throws": ["Ability1", "Ability2"],
  "armor_proficiencies": ["Light", "Medium", etc.],
  "weapon_proficiencies": ["Simple", "Martial"],
  "tool_proficiencies": ["Tool names"],
  "skill_choices": {
    "choose": 2,
    "from": ["Skill1", "Skill2", ...]
  },
  "starting_equipment": ["Item descriptions"],
  "spellcasting": {
    "ability": "Charisma",
    "ritual_casting": true
  },
  "level_progression": [
    {
      "level": 1,
      "proficiency_bonus": 2,
      "features": [
        {
          "name": "Feature Name",
          "description": "Feature description",
          "uses": {
            "max": 2,
            "recharge": "long_rest"
          },
          "mechanics": {
            "custom_key": "custom_value"
          }
        }
      ],
      "spellcasting_info": {
        "cantrips_known": 2,
        "spells_known": 4,
        "spell_slots": {
          "1": 2,
          "2": 0
        }
      },
      "ability_score_improvement": false
    }
  ]
}
```

## Level Progression Fields

### Per Level
- **level**: The character level (1-20)
- **proficiency_bonus**: Proficiency bonus at this level (2-6)
- **features**: Array of features gained at this level
- **spellcasting_info**: Spell slots and known spells (for spellcasters)
- **ability_score_improvement**: Boolean indicating if ASI is available

### Feature Structure
Each feature can have:
- **name**: Feature name
- **description**: Full feature description
- **uses**: Optional, for limited-use features
  - **max**: Number of uses (or "proficiency_bonus", "unlimited")
  - **recharge**: "short_rest" or "long_rest"
- **mechanics**: Optional, game-mechanical data
  - Can contain any relevant data (damage_bonus, inspiration_die, etc.)
- **subclass_choice**: Boolean, indicates subclass selection
- **subclass_feature**: Boolean, indicates a subclass-specific feature

## Mechanics Tracking

The `mechanics` object within features can track:
- **damage_bonus**: Extra damage (e.g., Rage damage)
- **rage_count**: Number of rages
- **inspiration_die**: Bardic inspiration die size
- **healing_die**: Song of Rest die size
- **expertise_count**: Number of expertise selections
- **magical_secrets**: Number of Magical Secrets learned
- **speed_bonus**: Speed increase
- **brutal_critical_dice**: Extra critical hit dice
- **ability_increases**: Ability score increases (e.g., Primal Champion)
- **ability_max_increase**: New ability score maximum

## Spellcasting Info

For spellcasting classes, each level includes:
- **cantrips_known**: Number of cantrips known
- **spells_known**: Number of spells known (for known-spell casters)
- **spell_slots**: Dictionary of spell level to number of slots

Example:
```json
"spellcasting_info": {
  "cantrips_known": 3,
  "spells_known": 8,
  "spell_slots": {
    "1": 4,
    "2": 3,
    "3": 2
  }
}
```

## Implementation Status

- ✅ **Barbarian** - Complete (Levels 1-20)
- ✅ **Bard** - Complete (Levels 1-20)
- ✅ **Cleric** - Complete (Levels 1-20)
- ✅ **Druid** - Complete (Levels 1-20)
- ✅ **Fighter** - Complete (Levels 1-20)
- ✅ **Monk** - Complete (Levels 1-20)
- ✅ **Paladin** - Complete (Levels 1-20, Half-caster)
- ✅ **Ranger** - Complete (Levels 1-20, Half-caster)
- ✅ **Rogue** - Complete (Levels 1-20)
- ✅ **Sorcerer** - Complete (Levels 1-20)
- ✅ **Warlock** - Complete (Levels 1-20, Pact Magic)
- ✅ **Wizard** - Complete (Levels 1-20)

## Usage in Code

The application should:
1. Load individual class files from `data/classes/`
2. Parse level_progression array
3. Apply features automatically when character levels up
4. Track mechanical bonuses in character model
5. Update UI to show level-appropriate features

## Future Enhancements

- Subclass definitions with their own level progression
- Automated level-up wizard that applies features
- Feature dependencies and prerequisites
- Multi-classing support
