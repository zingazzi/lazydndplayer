# Species Spell Grants

This document details which spells are automatically granted by each species and at what character level they become available.

## Species with Spell Grants

### Aasimar
**Trait:** Light Bearer
- **Level 1+:** Light (cantrip)

### Elf (Drow)
**Trait:** Drow Magic
- **Level 1+:** Dancing Lights (cantrip)
- **Level 3+:** Faerie Fire (1st level, 1/day)
- **Level 5+:** Darkness (2nd level, 1/day)

### Tiefling (Abyssal)
**Trait:** Abyssal Legacy
- **Level 1+:** Dancing Lights (cantrip)
- **Level 3+:** Burning Hands (1st level, 1/day)
- **Level 5+:** Alter Self (2nd level, 1/day)

### Tiefling (Chthonic)
**Trait:** Chthonic Legacy
- **Level 1+:** Chill Touch (cantrip)
- **Level 3+:** False Life (1st level, 1/day)
- **Level 5+:** Ray of Enfeeblement (2nd level, 1/day)

### Tiefling (Infernal)
**Trait:** Infernal Legacy
- **Level 1+:** Thaumaturgy (cantrip)
- **Level 3+:** Hellish Rebuke (1st level, 1/day)
- **Level 5+:** Darkness (2nd level, 1/day)

## How It Works

### Automatic Spell Granting
When you select a species that grants spells:
1. **Cantrips** are immediately added to your spellbook (Level 1+)
2. **Higher-level spells** are automatically added when you reach the required level
3. All species-granted spells are marked as **Known** in your spellbook
4. **Cantrips** are always prepared, higher-level spells can be prepared/unprepared

### Level-Based Progression
As your character levels up, new spells become available:
- **Level 1:** All cantrips from your species
- **Level 3:** Additional 1st-level spells (if any)
- **Level 5:** Additional 2nd-level spells (if any)

### Changing Species
When you change your character's species:
- ✅ Old species spells are **automatically removed**
- ✅ New species spells are **automatically added** (based on your current level)
- ✅ Only spells that came from species are affected
- ✅ Spells learned through class/other means are **not removed**

### Tracking
The system tracks which spells came from your species in the `SpeciesSpells` field. This ensures:
- Species spells are removed when changing species
- Class/learned spells are preserved
- No duplicate spells are added

## Spell Details

### Cantrips (Level 0)

**Light** (Evocation)
- Casting Time: 1 action
- Range: Touch
- Duration: 1 hour
- Effect: Touched object sheds bright light (20 ft) and dim light (20 ft additional)

**Dancing Lights** (Evocation)
- Casting Time: 1 action
- Range: 120 feet
- Duration: Concentration, up to 1 minute
- Effect: Create up to 4 torch-sized hovering lights

**Chill Touch** (Necromancy)
- Casting Time: 1 action
- Range: 120 feet
- Duration: 1 round
- Effect: Ranged spell attack dealing 1d8 necrotic damage

**Thaumaturgy** (Transmutation)
- Casting Time: 1 action
- Range: 30 feet
- Duration: Up to 1 minute
- Effect: Manifest minor supernatural wonders (booming voice, flickering flames, etc.)

### 1st Level Spells

**Faerie Fire** (Evocation)
- Casting Time: 1 action
- Range: 60 feet
- Duration: Concentration, up to 1 minute
- Effect: Outline creatures/objects in light, granting advantage on attacks against them

**Burning Hands** (Evocation)
- Casting Time: 1 action
- Range: Self (15-foot cone)
- Duration: Instantaneous
- Effect: 3d6 fire damage in a cone (Dex save for half)

**False Life** (Necromancy)
- Casting Time: 1 action
- Range: Self
- Duration: 1 hour
- Effect: Gain 1d4 + 4 temporary hit points

**Hellish Rebuke** (Evocation)
- Casting Time: 1 reaction
- Range: 60 feet
- Duration: Instantaneous
- Effect: Retaliate with 2d10 fire damage when damaged (Dex save for half)

### 2nd Level Spells

**Darkness** (Evocation)
- Casting Time: 1 action
- Range: 60 feet
- Duration: Concentration, up to 10 minutes
- Effect: Create a 15-foot radius sphere of magical darkness

**Alter Self** (Transmutation)
- Casting Time: 1 action
- Range: Self
- Duration: Concentration, up to 1 hour
- Effect: Change appearance, adapt to water, or grow natural weapons

**Ray of Enfeeblement** (Necromancy)
- Casting Time: 1 action
- Range: 60 feet
- Duration: Concentration, up to 1 minute
- Effect: Target deals half damage with Strength-based weapon attacks

## Examples

### Example 1: Tiefling (Infernal) Level 1
- **Immediately gains:** Thaumaturgy (cantrip)
- **At Level 3:** Gains Hellish Rebuke
- **At Level 5:** Gains Darkness

### Example 2: Changing Species
Starting as **Drow** (Level 4):
- Has: Dancing Lights, Faerie Fire

Changes to **Tiefling (Infernal)**:
- **Removed:** Dancing Lights, Faerie Fire
- **Added:** Thaumaturgy, Hellish Rebuke (level 4 ≥ 3)

### Example 3: Aasimar Level 1
- **Immediately gains:** Light (cantrip)
- Simple but effective illumination spell

## Implementation Notes

- Species spells use the character's primary spellcasting ability (INT, WIS, or CHA)
- Cantrips scale with character level (where applicable)
- Species-granted spells are considered "innate" and don't require spell slots in D&D 5e 2024
- For simplicity, the app adds them to the regular spellbook with "Known" status

