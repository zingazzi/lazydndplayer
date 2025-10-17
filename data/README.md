# Data Directory

This directory contains JSON data files for D&D 5e game content.

## Files

### `species.json`
Contains all available species (races) for character creation.

#### Structure:
```json
{
  "species": [
    {
      "name": "Species Name",
      "size": "Medium|Small",
      "speed": 30,
      "description": "Flavor text about the species",
      "traits": [
        {
          "name": "Trait Name",
          "description": "Trait description and mechanics"
        }
      ],
      "languages": ["Language1", "Language2"],
      "resistances": ["Damage Type"],
      "darkvision": 60
    }
  ]
}
```

#### Special Values:
- **speed**: Walking speed in feet (typically 25-35)
- **darkvision**: Range in feet, use `0` for no darkvision
- **resistances**: Damage types like "Fire", "Necrotic", "Poison", etc.
- **languages**: Can include "One additional language" or "One additional language of your choice" to trigger language selection

#### Adding a New Species:
1. Add a new object to the `species` array
2. Fill in all required fields
3. Restart the application (no recompilation needed!)

#### Skill Proficiencies:
- Traits with "Keen Senses" or "Proficiency in Perception" automatically grant Perception proficiency
- Traits with "Skillful" or "proficiency in one skill of your choice" trigger the skill selector

## Notes
- The application caches species data after first load
- If this file is missing, the app uses hardcoded fallback data
- Always validate JSON syntax before running the app
