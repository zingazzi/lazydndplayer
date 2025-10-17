# Features System Guide

## Overview

The **Features** tab replaces the old **Actions** tab and provides a comprehensive system for managing limited-use abilities that recharge on rest.

## What are Features?

Features are special abilities from your class, species, feats, or magic items that have:
- **Limited uses** per rest (e.g., "Rage" - 2 uses)
- **Recharge mechanics** (short rest, long rest, or daily)
- **Sources** (class, species, feat, item)

## Types of Features by Recharge

### ⚡ Short Rest Features
Recharge after a **short rest (1 hour)** or long rest:
- **Barbarian**: Rage
- **Fighter**: Second Wind, Action Surge
- **Monk**: Ki Points
- **Warlock**: Spell Slots

### 🌙 Long Rest Features
Recharge only after a **long rest (8 hours)**:
- **Wizard**: Arcane Recovery
- **Druid**: Wild Shape
- **Paladin**: Lay on Hands pool
- **Cleric**: Channel Divinity

### 📅 Daily Features
Recharge at dawn (similar to long rest):
- Some magic items
- Certain feats

### ✨ Passive Features
Always active, no charges:
- Ability modifiers
- Proficiency bonuses
- Always-on traits

## Usage

### Navigating Features
- **↑/↓** or **j/k**: Navigate through features
- **Ctrl+D/U**: Page down/up
- **Ctrl+E/Y**: Scroll down/up one line

### Managing Charges
- **u**: Use feature (consume one charge)
- **+** or **=**: Restore one charge manually
- **d**: Delete feature

### Resting
- **r**: Take a **short rest** (recovers short rest & daily features)
- **Shift+R**: Take a **long rest** (recovers all features)

### Adding Features
- **a**: Add new feature (not yet fully implemented)

## Example Features

### Barbarian Level 3
```
Feature: Rage
Max Uses: 3
Recharge: Long Rest
Description: Enter a rage as a bonus action. While raging, you gain advantage on Strength checks and saving throws, +2 damage with melee weapons using Strength, and resistance to bludgeoning, piercing, and slashing damage.
Source: Class: Barbarian

Feature: Reckless Attack
Max Uses: 0 (passive)
Recharge: None
Description: When you make your first attack on your turn, you can decide to attack recklessly. Doing so gives you advantage on melee weapon attack rolls using Strength, but attack rolls against you have advantage until your next turn.
Source: Class: Barbarian
```

### Dragonborn Species
```
Feature: Breath Weapon
Max Uses: 1
Recharge: Short Rest
Description: Exhale destructive energy in a 15-foot cone. Each creature in the area must make a Dexterity saving throw (DC = 8 + Constitution modifier + proficiency bonus). A creature takes 2d6 damage on a failed save, or half as much on a successful one.
Source: Species: Dragonborn
```

### Fighter Level 2
```
Feature: Action Surge
Max Uses: 1
Recharge: Short Rest
Description: On your turn, you can take one additional action. At level 17, you can use this twice before a rest.
Source: Class: Fighter

Feature: Second Wind
Max Uses: 1
Recharge: Short Rest
Description: You can use a bonus action to regain hit points equal to 1d10 + your fighter level.
Source: Class: Fighter
```

## Integration with Character Sheet

Features automatically:
- **Track usage** across sessions
- **Save** when modified
- **Recover** on appropriate rest type
- **Display** source and description

## Visual Indicators

- **White**: Available feature
- **Red + [USED]**: Depleted feature
- **Pink/Bold**: Selected feature
- **Gray**: Description and source

## Tips

1. **Add features as you level up** to track all your limited abilities
2. **Use short rests** between encounters to recover short rest features
3. **Long rest** at the end of the day to recover everything
4. **Delete outdated features** when multiclassing or retraining

## Future Enhancements

- [ ] Interactive form to add features with validation
- [ ] Auto-populate features based on class and level
- [ ] Feature templates for common abilities
- [ ] Usage history/log
- [ ] Spell slot integration as features
