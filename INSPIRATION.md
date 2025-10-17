# Inspiration System

## Overview

Inspiration is a D&D 5e mechanic that rewards good roleplay. It's a simple binary state (you either have it or you don't) that can be used to gain advantage on any d20 roll.

## How It Works

### Display
In the **Character Stats Panel**, you'll see the Inspiration status:

```
☐ Inspiration          (when you don't have it)
☑ Inspiration          (when you have it - shown in green)
```

### For Humans
Humans have a special species trait called **"Resourceful"**:
```
☑ Inspiration (auto-restored on rest)
```

- **Automatically restored** on every long rest
- You don't need a DM to grant it
- This is part of the 2024 D&D rules for Humans

### For Other Species
For all other species, Inspiration is:
- **Manually toggled** by pressing `Shift+I`
- **Given by the DM** during gameplay for good roleplay
- **Used by the player** when they want advantage on a roll
- **Not automatically restored** (except through specific class features)

## Usage

### Toggle Inspiration
When focused on the **Character Stats Panel**:
1. Press `Shift+I` to toggle Inspiration on/off
2. Green checkmark (☑) means you have it
3. Gray checkbox (☐) means you don't have it

### Using Inspiration
In actual gameplay:
1. **When to use**: Before making any d20 roll (attack, save, check)
2. **Effect**: Roll 2d20 and take the higher result (advantage)
3. **After use**: Toggle it off (`Shift+I`) since you consumed it
4. **Limitation**: You can only have one inspiration at a time

### Gaining Inspiration (Non-Humans)
The DM awards Inspiration for:
- Excellent roleplay
- Creative problem solving
- Advancing character development
- Acting in accordance with your character's traits/bonds/flaws

When the DM says "You gain Inspiration":
1. Focus the Character Stats Panel (`f` key)
2. Press `Shift+I` to turn it on
3. The app will show: **"✨ Inspiration gained!"**

### Consuming Inspiration
When you use Inspiration during a roll:
1. Tell the DM you're using Inspiration
2. Roll with advantage (2d20, take higher)
3. Toggle it off: `Shift+I`
4. The app will show: **"Inspiration used"**

## Long Rest Behavior

### Humans
```
Before Long Rest: ☐ Inspiration
Press: Shift+R (Long Rest)
After Long Rest:  ☑ Inspiration (auto-restored)
```

### Other Species
```
Before Long Rest: ☐ Inspiration
Press: Shift+R (Long Rest)
After Long Rest:  ☐ Inspiration (unchanged)
```

## Keyboard Shortcuts

### Character Stats Panel
- `Shift+I` - Toggle Inspiration on/off
- `Shift+R` - Long rest (restores Inspiration for Humans only)

## Visual Examples

### Dragonborn (No Inspiration)
```
⚔ Thorin Dragonborn
Class: Fighter, Level 5

❤ HP        🛡 AC       ⚡ INIT
  45/45       18          +2

👣 SPD      ⭐ PRF
  30ft        +3

☐ Inspiration

[n] Name • [r] Species • [h] HP • [+/-] ±1 • [i] Init • [I] Inspiration
```

### Human (With Inspiration - Auto-restored)
```
⚔ Sarah Human
Class: Cleric, Level 3

❤ HP        🛡 AC       ⚡ INIT
  24/24       16          +1

👣 SPD      ⭐ PRF
  30ft        +2

☑ Inspiration (auto-restored on rest)

[n] Name • [r] Species • [h] HP • [+/-] ±1 • [i] Init • [I] Inspiration
```

## Integration with Character Data

Inspiration is:
- ✅ **Saved** when toggled
- ✅ **Persisted** across app restarts
- ✅ **Automatically restored** for Humans on long rest
- ✅ **Manually managed** for all other species

## Technical Details

### Character Model
```go
type Character struct {
    // ... other fields ...
    Inspiration bool `json:"inspiration"`
}
```

### Long Rest Logic
```go
func (c *Character) LongRest() {
    // ... HP, spells, features recovery ...
    
    // Humans regain Inspiration (Resourceful trait)
    if c.Race == "Human" {
        c.Inspiration = true
    }
}
```

### Species-Specific Behavior
Only **Humans** have the "Resourceful" trait from D&D 5e 2024:
- Trait: "You gain Inspiration whenever you finish a Long Rest"
- Implementation: Automatic `Inspiration = true` on long rest
- Display: Shows "(auto-restored on rest)" label

## D&D 5e 2024 Rules Reference

From the D&D 5e 2024 Player's Handbook:

### Inspiration (General)
- Binary state: you either have it or don't
- Used to gain advantage on any d20 roll
- You can't have multiple inspirations
- Awarded by DM for good roleplay

### Human "Resourceful" Trait
- "You gain Inspiration whenever you finish a Long Rest"
- Unique to Humans in the 2024 edition
- Makes Humans more versatile and reliable

## FAQ

**Q: Can I have multiple Inspirations?**
A: No, it's binary. You either have it or you don't.

**Q: Do other species get automatic Inspiration?**
A: No, only Humans have the "Resourceful" trait.

**Q: What if I'm a Human multiclass?**
A: As long as your species is "Human", you get the benefit.

**Q: Can I toggle Inspiration without a DM?**
A: Yes, for tracking purposes. But in actual gameplay, only use it when the DM grants it (or you're Human after a rest).

**Q: Does a short rest restore Inspiration?**
A: No, only long rest restores it for Humans.

**Q: What if I change species from Human to something else?**
A: You'll lose the auto-restore feature. Inspiration won't be automatically restored on long rests anymore.

## Tips

1. **Humans**: Take long rests regularly to maintain your Inspiration pool
2. **Other Species**: Save Inspiration for critical moments (death saves, important skill checks)
3. **Roleplay**: Engage actively to earn Inspiration from your DM
4. **Track Usage**: Toggle it off immediately after using it so you don't forget
5. **Communication**: Always tell your DM when you're using Inspiration before rolling

