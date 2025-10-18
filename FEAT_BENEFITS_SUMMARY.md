# Feat Benefits - Add & Remove System ✅

## What's Implemented

### ✅ Reversible Feat Benefits
When you add or remove a feat, ALL benefits are automatically applied or reversed:

| Feat | Benefit | On Add | On Remove |
|------|---------|--------|-----------|
| **Actor** | +1 Charisma | CHA +1 (max 20) | CHA -1 (min 1) |
| **Athlete** | +1 STR or DEX* | STR/DEX +1 | STR/DEX -1 |
| **Durable** | +1 Constitution | CON +1 | CON -1 |
| **Tough** | +2 HP per level | MaxHP +10 (lv5) | MaxHP -10 |
| **Mobile** | +10 speed | Speed +10 | Speed -10 |
| **Alert** | +5 initiative | (planned) | (planned) |

\* For feats with choices (not yet in UI, but tracked in backend)

### ✅ Character Tracking
Your character now stores:
```json
{
  "feats": ["Tough", "Durable"],
  "feat_choices": {
    "Athlete": "Dexterity"  // If you had chosen this
  }
}
```

### ✅ Safe Limits
- **Ability Scores**: Capped at 20 (max) and 1 (min)
- **HP**: Minimum 1
- **Speed**: Minimum 0

## How to Test

### Test 1: Actor Feat (+1 Charisma)
```bash
./lazydndplayer

1. Note your current Charisma (e.g., 10)
2. Go to Traits panel (Tab until you reach it)
3. Press 'f' → Select "Actor" → Press Enter
   ✅ Charisma should be 11
   ✅ Message: "Feat gained: Actor!"

4. Press 'F' → Select "Actor" → Press Enter
   ✅ Charisma should be 10 again
   ✅ Message: "Feat removed: Actor (benefits reversed)"
```

### Test 2: Tough Feat (+2 HP per level)
```bash
./lazydndplayer

1. Note your level and Max HP (e.g., Level 5, 40 HP)
2. Go to Traits panel
3. Press 'f' → Select "Tough" → Press Enter
   ✅ Max HP should be 50 (40 + 10)
   ✅ Current HP set to 50
   ✅ Message: "Feat gained: Tough!"

4. Press 'F' → Select "Tough" → Press Enter
   ✅ Max HP should be 40 again
   ✅ Current HP capped to 40
   ✅ Message: "Feat removed: Tough (benefits reversed)"
```

### Test 3: Mobile Feat (+10 speed)
```bash
./lazydndplayer

1. Note your current speed (e.g., 30 ft)
2. Go to Traits panel
3. Press 'f' → Select "Mobile" → Press Enter
   ✅ Speed should be 40 ft
   ✅ Message: "Feat gained: Mobile!"

4. Press 'F' → Select "Tough" → Press Enter
   ✅ Speed should be 30 ft again
   ✅ Message: "Feat removed: Mobile (benefits reversed)"
```

### Test 4: Durable Feat (+1 Constitution)
```bash
./lazydndplayer

1. Note your current Constitution (e.g., 14)
2. Go to Traits panel
3. Press 'f' → Select "Durable" → Press Enter
   ✅ Constitution should be 15
   ✅ Message: "Feat gained: Durable!"

4. Press 'F' → Select "Durable" → Press Enter
   ✅ Constitution should be 14 again
   ✅ Message: "Feat removed: Durable (benefits reversed)"
```

## Current Limitations

### ⏳ Ability Choices Not in UI Yet
Feats like **Athlete** (Strength OR Dexterity) will:
- Currently: Not apply any ability increase (needs manual selection UI)
- Future: Show a popup to choose which ability to increase

**Affected Feats:**
- Athlete (STR or DEX)
- Weapon Master (STR or DEX)
- Tavern Brawler (STR or CON)
- Observant (INT or WIS)
- Fey Touched (INT, WIS, or CHA)
- Shadow Touched (INT, WIS, or CHA)
- Telekinetic (INT, WIS, or CHA)
- Telepathic (INT, WIS, or CHA)
- Resilient (Any ability)
- Skill Expert (Any ability)

### ✅ Working Now
All feats with **single ability increases** work perfectly:
- Actor (+1 CHA)
- Durable (+1 CON)
- Heavy Armor Master (+1 STR)
- Keen Mind (+1 INT)

And feats with **special benefits**:
- Tough (+2 HP per level)
- Mobile (+10 speed)

## How It Works Behind the Scenes

### Adding a Feat
```
1. User presses 'f' in Traits panel
2. Selects feat (e.g., "Durable")
3. Presses Enter

Backend:
4. Adds feat name to character.Feats []
5. Calls ApplyFeatBenefits(character, feat, "")
   - Increases Constitution by 1
   - Updates derived stats
6. Saves character
```

### Removing a Feat
```
1. User presses 'F' in Traits panel
2. Selects feat from their list (e.g., "Durable")
3. Presses Enter

Backend:
4. Removes feat name from character.Feats []
5. Calls RemoveFeatBenefits(character, feat)
   - Decreases Constitution by 1
   - Updates derived stats
6. Saves character
```

## Benefits of This System

✅ **Automatic** - No manual stat tracking needed
✅ **Reversible** - Remove feats without breaking your character
✅ **Consistent** - All feats use the same system
✅ **Safe** - Ability scores stay within valid ranges (1-20)
✅ **Trackable** - All changes are recorded

## Example Scenarios

### Scenario 1: Testing Different Builds
```
"I want to try Tough vs Durable to see which I prefer"

1. Add Tough → Max HP +10
2. Play for a bit
3. Remove Tough → Max HP -10 (back to normal)
4. Add Durable → CON +1 (also affects HP modifier)
5. Compare which you like better!
```

### Scenario 2: Fixing Mistakes
```
"Oops, I accidentally added the wrong feat!"

1. Press 'F' (Shift+F)
2. Select the wrong feat
3. Press Enter
   ✅ All benefits removed
   ✅ Character back to pre-feat state
```

### Scenario 3: Level Up
```
"I'm level 4, time to choose my Ability Score Improvement (ASI) or feat!"

Current: Feats work, but ASI needs to be manual
Future: Will have dedicated ASI/Feat choice at level up
```

## Status

| Component | Status |
|-----------|--------|
| Backend feat tracking | ✅ Complete |
| Apply feat benefits | ✅ Complete |
| Remove feat benefits | ✅ Complete |
| Single ability increases | ✅ Working |
| Special benefits (Tough, Mobile) | ✅ Working |
| Ability choice feats (Athlete, etc.) | ⏳ Backend ready, UI pending |
| Ability choice UI selector | ✅ Created, ⏳ Not integrated |
| Alert feat (+5 initiative) | ⏳ Planned |

## Next Steps (For Future)

1. **Integrate Ability Choice Selector** - Show popup for feats with choices
2. **Add Alert Feat Handling** - Track and apply +5 initiative
3. **Level Up Wizard** - Auto-offer feats at levels 4, 8, 12, 16, 19
4. **Feat Details in Traits Panel** - Show chosen abilities next to feat names
5. **Feat Search/Filter** - Search feats by name or category

## To Run

```bash
./lazydndplayer
```

Then test adding and removing feats as described above!

## Notes

- ⚠️ Removing a feat with choices (like Athlete) will reverse the ability increase, but you won't see which ability was chosen in the UI yet
- ✅ All changes are auto-saved to your character file
- ✅ The system prevents you from having ability scores above 20 or below 1
- ✅ Removing Tough properly reduces max HP and current HP if needed

Enjoy your new reversible feat system! 🎉
