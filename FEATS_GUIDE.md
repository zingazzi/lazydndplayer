# Feats Guide - D&D 5e 2024 Edition

## Overview

Feats are special abilities that characters can gain to customize their capabilities. Instead of taking an Ability Score Improvement (ASI) at certain levels, you can choose a feat to give your character unique powers and advantages.

## When Can You Take Feats?

Feats can be taken at the following levels (when you would normally get an ASI):
- **Level 4**
- **Level 8**
- **Level 12**
- **Level 16**
- **Level 19** (Fighter gets additional ASI at level 6 and 14)

## Feat System

### Files
- **`data/feats.json`** - Complete list of all D&D 5e 2024 feats with their benefits
- **`internal/models/feats.go`** - Feat loading, validation, and application logic

### Feat Properties

Each feat has the following properties:

| Property | Description |
|----------|-------------|
| **Name** | The feat's name |
| **Category** | Type of feat (General, Combat, etc.) |
| **Prerequisite** | Requirements to take this feat |
| **Repeatable** | Whether you can take this feat multiple times |
| **Benefits** | List of mechanical and narrative benefits |
| **Description** | Summary of what the feat does |
| **Ability Increases** | Which ability scores increase (if any) |
| **Grants Spells** | Any spells granted by this feat |
| **Note** | Additional information |

## Popular Feats

### Combat Feats

#### Great Weapon Master
- **Prerequisite:** None
- **Benefits:**
  - Critical hit or kill → bonus action attack
  - Take -5 to hit for +10 damage (with heavy weapons)
- **Best for:** Barbarians, Fighters, Paladins

#### Sharpshooter
- **Prerequisite:** None
- **Benefits:**
  - No disadvantage at long range
  - Ignore half and 3/4 cover
  - Take -5 to hit for +10 damage (with ranged weapons)
- **Best for:** Rangers, Rogues, Fighters

#### Polearm Master
- **Prerequisite:** None
- **Benefits:**
  - Bonus action attack with opposite end (1d4)
  - Opportunity attacks when enemies enter reach
- **Best for:** Fighters, Paladins (reach weapon builds)

#### Sentinel
- **Prerequisite:** None
- **Benefits:**
  - Opportunity attacks reduce speed to 0
  - Attack enemies who attack your allies
  - Enemies can't Disengage from you
- **Best for:** Tanks, Defenders

### Defensive Feats

#### Tough
- **Prerequisite:** None
- **Benefits:**
  - +2 HP per level (retroactive and ongoing)
- **Best for:** Anyone who needs more survivability

#### Heavy Armor Master
- **Prerequisite:** Heavy armor proficiency
- **Benefits:**
  - +1 Strength
  - Reduce non-magical physical damage by 3
- **Best for:** Fighters, Paladins, Clerics

#### Shield Master
- **Prerequisite:** None
- **Benefits:**
  - Bonus action shove with shield
  - Add shield AC to DEX saves
  - Take no damage on successful DEX saves
- **Best for:** Anyone using a shield

### Utility Feats

#### Lucky
- **Prerequisite:** None
- **Benefits:**
  - 3 luck points per long rest
  - Reroll attack rolls, ability checks, or saving throws
- **Best for:** Everyone (universally strong)

#### Alert
- **Prerequisite:** None
- **Benefits:**
  - +5 to initiative
  - Can't be surprised
  - No advantage for hidden enemies
- **Best for:** Anyone who wants to go first

#### Mobile
- **Prerequisite:** None
- **Benefits:**
  - +10 ft speed
  - Dash through difficult terrain
  - No opportunity attacks from creatures you attack
- **Best for:** Monks, Rogues, skirmishers

#### Observant
- **Prerequisite:** None
- **Benefits:**
  - +1 INT or WIS
  - Read lips
  - +5 to passive Perception and Investigation
- **Best for:** Investigators, scouts

### Spellcasting Feats

#### War Caster
- **Prerequisite:** Ability to cast at least one spell
- **Benefits:**
  - Advantage on concentration saves
  - Cast with weapon/shield in hand
  - Cast spells as opportunity attacks
- **Best for:** Gish builds, Clerics, Paladins

#### Magic Initiate
- **Prerequisite:** None
- **Repeatable:** Yes
- **Benefits:**
  - Learn 2 cantrips + 1 1st-level spell
  - Choose from any class's spell list
  - Cast the spell once per long rest
- **Best for:** Non-casters gaining utility spells

#### Spell Sniper
- **Prerequisite:** Ability to cast at least one spell
- **Benefits:**
  - Double spell range
  - Ignore half and 3/4 cover
  - Learn one attack roll cantrip
- **Best for:** Ranged spellcasters

#### Ritual Caster
- **Prerequisite:** INT or WIS 13+
- **Benefits:**
  - Get a ritual book with 2 1st-level ritual spells
  - Can add more ritual spells to the book
- **Best for:** Utility casters

### Half-Feat Feats (+1 Ability Score)

These feats give you a +1 ability score increase along with other benefits:

| Feat | Ability | Best Benefits |
|------|---------|---------------|
| **Actor** | CHA +1 | Advantage on impersonation, mimic voices |
| **Athlete** | STR/DEX +1 | Easier climbing/jumping, fast stand-up |
| **Durable** | CON +1 | Better HP recovery from hit dice |
| **Heavy Armor Master** | STR +1 | -3 physical damage reduction |
| **Keen Mind** | INT +1 | Perfect memory, always know direction |
| **Observant** | INT/WIS +1 | +5 passive Perception/Investigation |
| **Resilient** | Any +1 | Gain proficiency in that ability's saves |
| **Skill Expert** | Any +1 | Gain proficiency + expertise in a skill |
| **Tavern Brawler** | STR/CON +1 | Better unarmed strikes, grapple bonus |
| **Weapon Master** | STR/DEX +1 | Proficiency with 4 weapons |
| **Fey Touched** | INT/WIS/CHA +1 | Misty Step + 1 divination/enchantment spell |
| **Shadow Touched** | INT/WIS/CHA +1 | Invisibility + 1 illusion/necromancy spell |
| **Telekinetic** | INT/WIS/CHA +1 | Mage Hand + bonus action shove |
| **Telepathic** | INT/WIS/CHA +1 | Telepathy + Detect Thoughts spell |

## Repeatable Feats

Some feats can be taken multiple times (choosing different options each time):

- **Elemental Adept** - Choose a different damage type each time
- **Magic Initiate** - Choose a different class's spell list each time
- **Resilient** - Choose a different ability score each time
- **Skill Expert** - Gain proficiency/expertise in different skills

## Feat Prerequisites

### Ability Score Requirements
- **STR 13+:** Grappler
- **DEX 13+:** Defensive Duelist, Skulker
- **INT 13+:** Ritual Caster (or WIS 13+)
- **WIS 13+:** Ritual Caster (or INT 13+)
- **CHA 13+:** Inspiring Leader

### Spellcasting Requirements
- **Must cast at least one spell:** Elemental Adept, Mage Slayer, Spell Sniper, War Caster

### Proficiency Requirements
- **Heavy armor proficiency:** Heavy Armor Master
- **Medium armor proficiency:** Medium Armor Master

## How to Choose a Feat

### 1. Consider Your Build
- **Damage dealer?** → Great Weapon Master, Sharpshooter, Polearm Master
- **Tank?** → Tough, Shield Master, Sentinel, Heavy Armor Master
- **Spellcaster?** → War Caster, Spell Sniper, Elemental Adept
- **Skill monkey?** → Skill Expert, Observant

### 2. Consider Your Role
- **Front-line fighter?** → Sentinel, Polearm Master, Great Weapon Master
- **Ranged attacker?** → Sharpshooter, Crossbow Expert
- **Support?** → Inspiring Leader, Healer
- **Scout?** → Alert, Mobile, Skulker

### 3. Consider Synergies
- **Polearm Master + Sentinel** = Area control powerhouse
- **Great Weapon Master + Barbarian's Reckless Attack** = Massive damage
- **War Caster + Booming Blade** = Amazing opportunity attacks
- **Mobile + Rogue** = Hit and run tactics

### 4. Half-Feats for Odd Ability Scores
If you have an odd ability score (like 15 or 17), consider a half-feat to get both the +1 and extra benefits:
- **15 STR → Heavy Armor Master** (16 STR + damage reduction)
- **17 WIS → Observant** (18 WIS + passive bonuses)
- **13 CHA → Fey Touched** (14 CHA + Misty Step)

## Feat vs ASI Decision

### Take a Feat When:
- ✅ You have a clear build strategy that the feat supports
- ✅ Your primary ability scores are already high (16-18)
- ✅ The feat provides unique capabilities you can't get elsewhere
- ✅ You have synergies with class features or other feats

### Take ASI When:
- ✅ Your primary ability score is below 16
- ✅ You're a spellcaster with low spell save DC
- ✅ You need to hit more often (attack bonus)
- ✅ You want to boost multiple abilities (STR/CON for melee)

## Programmatic Usage

### Load All Feats
```go
import "github.com/marcozingoni/lazydndplayer/internal/models"

feats := models.GetAllFeats()
for _, feat := range feats {
    fmt.Println(feat.Name)
}
```

### Get Specific Feat
```go
lucky := models.GetFeatByName("Lucky")
if lucky != nil {
    fmt.Println(models.FormatFeatForDisplay(*lucky))
}
```

### Check Prerequisites
```go
import "github.com/marcozingoni/lazydndplayer/internal/models"

char := models.NewCharacter()
sharpshooter := models.GetFeatByName("Sharpshooter")

if models.CanTakeFeat(char, *sharpshooter) {
    fmt.Println("Can take Sharpshooter!")
}
```

### Get Available Feats for Character
```go
availableFeats := models.GetFeatsForCharacter(char)
fmt.Printf("Character can take %d feats\n", len(availableFeats))
```

### Add Feat to Character
```go
err := models.AddFeatToCharacter(char, "Lucky")
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Println("Feat added successfully!")
}
```

### Apply Feat Benefits
```go
tough := models.GetFeatByName("Tough")
if tough != nil {
    // Automatically applies +2 HP per level
    models.ApplyFeatBenefits(char, *tough)
    fmt.Printf("New MaxHP: %d\n", char.MaxHP)
}
```

## JSON Structure

```json
{
  "name": "Lucky",
  "category": "General",
  "prerequisite": "None",
  "repeatable": false,
  "benefits": [
    "You have 3 luck points",
    "Spend a luck point to roll an additional d20",
    "Choose which d20 to use",
    "Regain all luck points on long rest"
  ],
  "description": "You have inexplicable luck...",
  "ability_increases": {},
  "note": "You have 3 luck points. Regain all on long rest."
}
```

## All Available Feats (D&D 5e 2024)

### Combat Feats
- Charger
- Crossbow Expert
- Defensive Duelist
- Dual Wielder
- Great Weapon Master
- Grappler
- Mage Slayer
- Martial Adept
- Polearm Master
- Savage Attacker
- Sentinel
- Sharpshooter

### Defensive Feats
- Durable
- Heavy Armor Master
- Medium Armor Master
- Shield Master
- Tough

### Movement Feats
- Mobile
- Skulker

### Utility Feats
- Actor
- Alert
- Athlete
- Healer
- Inspiring Leader
- Keen Mind
- Lucky
- Observant
- Tavern Brawler
- Weapon Master

### Spellcasting Feats
- Elemental Adept
- Magic Initiate
- Ritual Caster
- Spell Sniper
- War Caster

### Psionic/Magic Feats
- Fey Touched
- Shadow Touched
- Telekinetic
- Telepathic

### Skill Feats
- Skill Expert

### Save Feats
- Resilient

## Feat Recommendations by Class

### Fighter
1. **Great Weapon Master** or **Sharpshooter** (depending on weapon)
2. **Polearm Master** (reach weapon build)
3. **Sentinel** (tank build)
4. **Tough** (survivability)

### Wizard
1. **War Caster** (concentration)
2. **Resilient (CON)** (concentration saves)
3. **Alert** (go first to control battlefield)
4. **Lucky** (insurance on important saves)

### Rogue
1. **Sharpshooter** (if ranged)
2. **Mobile** (melee hit-and-run)
3. **Alert** (go first for burst damage)
4. **Skulker** (stealth builds)

### Cleric
1. **War Caster** (casting with shield)
2. **Resilient (CON)** (concentration)
3. **Heavy Armor Master** (damage reduction)
4. **Inspiring Leader** (support)

### Barbarian
1. **Great Weapon Master** (synergy with Reckless Attack)
2. **Polearm Master** (reach + reaction attacks)
3. **Sentinel** (lock down enemies)
4. **Tough** (even more HP)

## Future Enhancements

Planned additions to the feat system:
- [ ] UI for selecting feats during level-up
- [ ] Automatic feat benefit application
- [ ] Feat recommendation engine
- [ ] Track feat-granted spells in spellbook
- [ ] Track feat-granted proficiencies
- [ ] Feat synergy suggestions
- [ ] Character sheet feat display

## Summary

The feat system provides:
- ✅ Complete D&D 5e 2024 feat list (45+ feats)
- ✅ Prerequisite checking
- ✅ Repeatable feat support
- ✅ Automatic benefit application (HP, speed, etc.)
- ✅ Programmatic access API
- ✅ Formatted display for UI

Use `data/feats.json` as your reference when choosing feats at ASI levels!

