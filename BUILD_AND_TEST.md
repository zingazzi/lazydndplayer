# Build and Test Instructions

## Issues Fixed

### 1. Nil Pointer Dereference Crash (✅ FIXED)
**Problem**: When loading old character files that don't have `BenefitTracker`, the app crashed when trying to remove feats.

**Solution**: Added nil checks in three locations:

1. **`internal/storage/json.go`** - Initialize `BenefitTracker` when loading old character files
2. **`internal/models/benefit_applier.go`** - Ensure `BenefitTracker` exists in `NewBenefitApplier()`
3. **`internal/models/benefit_remover.go`** - Ensure `BenefitTracker` exists in `NewBenefitRemover()`

### 2. Ability Choice Selector Not Appearing
**Diagnosis Needed**: The feat selection flow should work like this:
- Press `f` in Traits panel
- Select "Athlete" feat
- Press Enter
- Should show ability choice selector (Strength or Dexterity)

## To Build and Test

### Step 1: Clean and Rebuild
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
rm -f lazydndplayer
go clean -cache
go build -o lazydndplayer .
```

### Step 2: Remove Old Character File (Important!)
```bash
# This will force creation of a new character with BenefitTracker
rm -f ~/.local/share/lazydndplayer/character.json
# Or wherever your character.json is stored
```

### Step 3: Run the App
```bash
./lazydndplayer
```

### Step 4: Test Feat Addition
1. Press Tab to navigate to Traits panel
2. Press `f` to add a feat
3. Navigate to "Athlete" (should be near the top)
4. Press Enter
5. **Expected**: Ability choice selector appears with "Strength" and "Dexterity"
6. **If nothing happens**: There's an issue with the feat loading or selector

### Step 5: Test Feat Removal
1. Ensure you have at least one feat
2. Press `F` (Shift+F) to remove a feat
3. Select a feat and press Enter
4. **Expected**: Feat removed, benefits reversed, no crash
5. **If crash**: Check if BenefitTracker is nil (should be fixed now)

## Debug: Check if Athlete Feat Loads Correctly

Create a simple test:
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go test ./internal/models -v -run TestLoadAthleteFeat
```

This will verify:
- Athlete feat has `AbilityIncreases` field populated
- `AbilityIncreases.Choices` contains ["Strength", "Dexterity"]
- `HasAbilityChoice()` returns true

## Possible Issues if Selector Still Doesn't Appear

### Issue 1: Feat Data Not Loading
Check if `data/feats.json` is being read correctly:
```go
// Add debug line in handleFeatSelectorKeys around line 1143
if models.HasAbilityChoice(*selectedFeat) {
    fmt.Printf("DEBUG: Feat has choice! Choices: %v\n", models.GetAbilityChoices(*selectedFeat))
    // ... rest of code
}
```

### Issue 2: Selector Component Not Visible
Check `abilityChoiceSelector` is rendering:
```go
// In app.go View() method, add debug around line where selector is shown
if m.abilityChoiceSelector.IsVisible() {
    fmt.Println("DEBUG: AbilityChoiceSelector is visible")
    return lipgloss.Place(...)
}
```

### Issue 3: Message Not Showing
The message "Choose which ability to increase" should appear at the bottom of the screen.

## Changes Made

### File: `internal/storage/json.go`
Added BenefitTracker initialization on load:
```go
// Initialize BenefitTracker if nil (for backwards compatibility with old saves)
if character.BenefitTracker == nil {
    character.BenefitTracker = models.NewBenefitTracker()
}
```

### File: `internal/models/benefit_applier.go`
Added safety check in constructor:
```go
func NewBenefitApplier(char *Character) *BenefitApplier {
    // Ensure BenefitTracker is initialized
    if char.BenefitTracker == nil {
        char.BenefitTracker = NewBenefitTracker()
    }
    return &BenefitApplier{char: char}
}
```

### File: `internal/models/benefit_remover.go`
Added safety check in constructor:
```go
func NewBenefitRemover(char *Character) *BenefitRemover {
    // Ensure BenefitTracker is initialized
    if char.BenefitTracker == nil {
        char.BenefitTracker = NewBenefitTracker()
    }
    return &BenefitRemover{char: char}
}
```

## Next Steps

1. **Build the app** - Use the commands above
2. **Delete old character file** - Very important to avoid loading old data
3. **Test feat addition** - Especially Athlete feat
4. **Test feat removal** - Should not crash anymore
5. **Report back** - Let me know if the ability selector appears or if there are other issues
