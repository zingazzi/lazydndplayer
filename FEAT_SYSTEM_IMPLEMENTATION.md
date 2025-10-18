# Feat System - Implementation Summary

## What's Been Implemented

✅ **Feat Choices Tracking** - Character now stores which ability was chosen for feats with options
✅ **Reversible Benefits** - Removing a feat now properly reverses ability score increases
✅ **Ability Choice Detection** - System automatically detects when a feat requires choosing an ability
✅ **Benefit Application** - Automatic application of ability increases, HP bonuses (Tough), speed bonuses (Mobile)
✅ **Benefit Removal** - Automatic removal of all feat benefits when feat is removed

## New Character Field

```go
type Character struct {
    ...
    Feats       []string              // List of feat names
    FeatChoices map[string]string     // feat name -> chosen ability
    ...
}
```

**Example:**
```json
{
  "feats": ["Athlete", "Tough"],
  "feat_choices": {
    "Athlete": "Dexterity"  // Player chose Dexterity over Strength
  }
}
```

## Feat Ability Increase Formats

The system handles these formats in `feats.json`:

### 1. Single Ability (No Choice)
```json
{
  "name": "Actor",
  "ability_increases": {"Charisma": 1}
}
```
→ Automatically increases Charisma by 1

### 2. Two Options ("or")
```json
{
  "name": "Athlete",
  "ability_increases": {"Strength or Dexterity": 1}
}
```
→ Player chooses Strength OR Dexterity

### 3. Multiple Options (comma-separated)
```json
{
  "name": "Fey Touched",
  "ability_increases": {"Intelligence, Wisdom, or Charisma": 1}
}
```
→ Player chooses one of the three

### 4. Any Ability
```json
{
  "name": "Resilient",
  "ability_increases": {"Any ability": 1}
}
```
→ Player chooses any of the six abilities

## Key Functions

### `HasAbilityChoice(feat Feat) bool`
Returns true if the feat requires the player to choose an ability.

```go
// Examples:
HasAbilityChoice(actorFeat)    // false - only Charisma
HasAbilityChoice(athleteFeat)  // true  - Strength or Dexterity
HasAbilityChoice(resilientFeat) // true  - Any ability
```

### `GetAbilityChoices(feat Feat) []string`
Returns the list of abilities that can be chosen.

```go
// Examples:
GetAbilityChoices(athleteFeat)
// → ["Strength", "Dexterity"]

GetAbilityChoices(resilientFeat)
// → ["Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"]
```

### `ApplyFeatBenefits(char *Character, feat Feat, chosenAbility string)`
Applies all feat benefits to the character.

```go
// Examples:
ApplyFeatBenefits(char, actorFeat, "")           // No choice needed
ApplyFeatBenefits(char, athleteFeat, "Dexterity") // With choice
ApplyFeatBenefits(char, toughFeat, "")           // +2 HP per level
ApplyFeatBenefits(char, mobileFeat, "")          // +10 speed
```

**What it does:**
- ✅ Increases ability scores (max 20)
- ✅ Stores the choice in `char.FeatChoices`
- ✅ Applies special benefits (Tough HP, Mobile speed)
- ✅ Updates derived stats

### `RemoveFeatBenefits(char *Character, feat Feat)`
Removes all feat benefits from the character.

```go
// Example:
RemoveFeatBenefits(char, athleteFeat)
// → Removes the +1 from whichever ability was chosen
// → Removes the choice from char.FeatChoices
```

**What it does:**
- ✅ Looks up the stored choice (if any)
- ✅ Removes ability score increases (min 1)
- ✅ Removes special benefits (Tough HP, Mobile speed)
- ✅ Cleans up `char.FeatChoices`
- ✅ Updates derived stats

## Future Integration (TODO)

The following components still need to be integrated into the UI:

### 1. Ability Choice Selector Component
File: `internal/ui/components/abilitychoiceselector.go`

**Status:** ✅ Created, ⏳ Needs integration

**Purpose:** Shows a popup when a feat with ability choices is selected.

**Example UI:**
```
╭─────────────────────────────────────╮
│   SELECT ABILITY FOR ATHLETE        │
│                                     │
│ Choose which ability to increase:   │
│                                     │
│ ▶ Strength (currently 14)           │
│   Dexterity (currently 16)          │
│                                     │
│ [↑/↓] Navigate • [Enter] Select     │
╰─────────────────────────────────────╯
```

### 2. App Integration Points

#### In `app.go`, add field:
```go
type Model struct {
    ...
    abilityChoiceSelector *components.AbilityChoiceSelector
    ...
}
```

#### In `NewModel()`:
```go
abilityChoiceSelector: components.NewAbilityChoiceSelector(),
```

#### In `handleFeatSelectorKeys()`, update the "enter" case:
```go
case "enter":
    selectedFeat := m.featSelector.GetSelectedFeat()
    if selectedFeat != nil {
        if m.featSelector.IsDeleteMode() {
            // REMOVE FEAT
            for i, featName := range m.character.Feats {
                if featName == selectedFeat.Name {
                    m.character.Feats = append(m.character.Feats[:i], m.character.Feats[i+1:]...)
                    break
                }
            }
            // REMOVE BENEFITS ← NEW!
            models.RemoveFeatBenefits(m.character, *selectedFeat)

            m.message = fmt.Sprintf("Feat removed: %s (benefits reversed)", selectedFeat.Name)
            m.storage.Save(m.character)
            m.featSelector.Hide()
        } else {
            // ADD FEAT
            if models.HasFeat(m.character, selectedFeat.Name) && !selectedFeat.Repeatable {
                m.message = fmt.Sprintf("You already have %s", selectedFeat.Name)
            } else {
                // Check if ability choice is needed ← NEW!
                if models.HasAbilityChoice(*selectedFeat) {
                    // Show ability choice selector
                    choices := models.GetAbilityChoices(*selectedFeat)
                    m.abilityChoiceSelector.Show(selectedFeat.Name, choices, m.character)
                    m.featSelector.Hide()
                    m.message = "Select which ability to increase..."
                } else {
                    // No choice needed, apply directly
                    err := models.AddFeatToCharacter(m.character, selectedFeat.Name)
                    if err != nil {
                        m.message = fmt.Sprintf("Error: %v", err)
                    } else {
                        models.ApplyFeatBenefits(m.character, *selectedFeat, "")
                        m.message = fmt.Sprintf("Feat gained: %s!", selectedFeat.Name)
                        m.storage.Save(m.character)
                    }
                    m.featSelector.Hide()
                }
            }
        }
    }
```

#### Add new handler function:
```go
func (m *Model) handleAbilityChoiceSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "up", "k":
        m.abilityChoiceSelector.Prev()
    case "down", "j":
        m.abilityChoiceSelector.Next()
    case "enter":
        selectedAbility := m.abilityChoiceSelector.GetSelectedAbility()
        featName := m.abilityChoiceSelector.featName  // Need to expose this

        if selectedAbility != "" {
            feat := models.GetFeatByName(featName)
            if feat != nil {
                err := models.AddFeatToCharacter(m.character, feat.Name)
                if err != nil {
                    m.message = fmt.Sprintf("Error: %v", err)
                } else {
                    // Apply with chosen ability
                    models.ApplyFeatBenefits(m.character, *feat, selectedAbility)
                    m.message = fmt.Sprintf("Feat gained: %s (+1 %s)!", feat.Name, selectedAbility)
                    m.storage.Save(m.character)
                }
            }
        }
        m.abilityChoiceSelector.Hide()
    case "esc":
        m.abilityChoiceSelector.Hide()
        m.message = "Feat selection cancelled"
    }
    return m, nil
}
```

#### In `Update()` method, add priority check:
```go
// Check if ability choice selector is active (high priority)
if m.abilityChoiceSelector.IsVisible() {
    return m.handleAbilityChoiceSelectorKeys(msg)
}
```

#### In `View()` method, add rendering:
```go
// Ability choice selector takes priority
if m.abilityChoiceSelector.IsVisible() {
    return m.abilityChoiceSelector.View(m.width, m.height)
}
```

## Special Feat Handling

### Tough Feat
- **On Add:** `char.MaxHP += char.Level * 2` (retroactive!)
- **On Remove:** `char.MaxHP -= char.Level * 2` (minimum 1)

### Mobile Feat
- **On Add:** `char.Speed += 10`
- **On Remove:** `char.Speed -= 10` (minimum 0)

### Alert Feat
- **Benefit:** +5 initiative
- **Note:** Currently not automatically applied (needs initiative tracking)

## Examples

### Example 1: Actor Feat (No Choice)
```go
// Actor: +1 Charisma (no choice)
feat := models.GetFeatByName("Actor")
models.ApplyFeatBenefits(char, *feat, "")

// Result:
// - char.AbilityScores.Charisma += 1 (max 20)
// - char.Feats contains "Actor"
// - char.FeatChoices has no entry (no choice needed)
```

### Example 2: Athlete Feat (With Choice)
```go
// Athlete: +1 Strength OR Dexterity
feat := models.GetFeatByName("Athlete")

// Player chooses Dexterity
models.ApplyFeatBenefits(char, *feat, "Dexterity")

// Result:
// - char.AbilityScores.Dexterity += 1 (max 20)
// - char.Feats contains "Athlete"
// - char.FeatChoices["Athlete"] = "Dexterity"

// Later, remove the feat
models.RemoveFeatBenefits(char, *feat)

// Result:
// - char.AbilityScores.Dexterity -= 1 (min 1)
// - char.FeatChoices["Athlete"] is deleted
// - char.Feats still contains "Athlete" (must be removed separately)
```

### Example 3: Tough Feat (Special Benefit)
```go
// Tough: +2 HP per level
char.Level = 5
char.MaxHP = 40

feat := models.GetFeatByName("Tough")
models.ApplyFeatBenefits(char, *feat, "")

// Result:
// - char.MaxHP = 50 (40 + 10)
// - char.CurrentHP = 50 (set to max)

// Remove the feat
models.RemoveFeatBenefits(char, *feat)

// Result:
// - char.MaxHP = 40 (50 - 10)
// - char.CurrentHP = 40 (capped to new max)
```

## Testing the System

### Test 1: Feat with No Choice
```bash
./lazydndplayer
# Go to Traits panel
# Press 'f'
# Select "Actor"
# Press Enter
# → Charisma should increase by 1
# Press 'F'
# Select "Actor"
# Press Enter
# → Charisma should decrease by 1
```

### Test 2: Feat with Choice (When Integrated)
```bash
./lazydndplayer
# Go to Traits panel
# Press 'f'
# Select "Athlete"
# Press Enter
# → Ability choice selector appears
# Select "Dexterity"
# Press Enter
# → Dexterity should increase by 1
# → FeatChoices should store "Athlete": "Dexterity"
# Press 'F'
# Select "Athlete"
# Press Enter
# → Dexterity should decrease by 1
# → FeatChoices should remove "Athlete" entry
```

### Test 3: Tough Feat
```bash
./lazydndplayer
# Note current Max HP (e.g., 40 at level 5)
# Go to Traits panel
# Press 'f'
# Select "Tough"
# Press Enter
# → Max HP should be 50 (40 + 10)
# Press 'F'
# Select "Tough"
# Press Enter
# → Max HP should be 40 again
```

## Benefits of This System

✅ **Consistent** - All feat modifications go through the same functions
✅ **Reversible** - Removing feats properly undoes all changes
✅ **Trackable** - All choices are stored in the character
✅ **Extensible** - Easy to add new special feat handling
✅ **Safe** - Ability scores capped at 20 (max) and 1 (min)

## Next Steps

To complete the implementation:
1. ✅ Update `ApplyFeatBenefits` calls in `app.go` to use new signature (add "")
2. ⏳ Integrate ability choice selector into feat selection flow
3. ⏳ Test with various feats (Actor, Athlete, Tough, Mobile, Resilient)
4. ⏳ Add Alert feat special handling for initiative
5. ⏳ Update UI to show chosen abilities in Traits panel (optional)

## Current Status

**Backend:** ✅ Complete
- Feat choice tracking
- Benefit application
- Benefit removal
- Ability choice detection

**UI:** ⏳ Partial
- Ability choice selector component created
- Needs integration into app flow
- Remove feat now calls RemoveFeatBenefits ✅

**Testing:** ⏳ Pending full integration
