# Feat System Testing Guide

## Overview
This document describes the comprehensive testing strategy for the modular benefit system, focusing on feat add/remove functionality.

## Test Files
- `internal/models/feats_test.go` - Unit tests for feat benefits

## Unit Tests Coverage

### 1. Single Ability Feat Tests
**TestApplyFeatBenefits_SingleAbility**
- Tests applying a feat with a single fixed ability (Actor: +1 Charisma)
- Verifies ability score increases
- Verifies benefit is tracked in BenefitTracker
- Verifies correct source, type, target, and value

**TestRemoveFeatBenefits_SingleAbility**
- Tests removing a feat with single ability
- Verifies ability score is restored to original value
- Verifies benefit is removed from BenefitTracker

### 2. Multiple Choice Feat Tests
**TestApplyFeatBenefits_MultipleChoice**
- Tests applying a feat with ability choices (Athlete: Strength or Dexterity)
- Verifies chosen ability increases
- Verifies other abilities remain unchanged
- Verifies correct ability is tracked

**TestRemoveFeatBenefits_MultipleChoice**
- Tests removing a feat with ability choice
- Verifies chosen ability is restored correctly
- Verifies other abilities remain unchanged
- Verifies benefit removal from tracker

### 3. Special Feat Tests
**TestApplyFeatBenefits_Tough**
- Tests Tough feat (+2 HP per level)
- Verifies HP calculation is correct
- Verifies benefit tracking for HP bonus

**TestRemoveFeatBenefits_Tough**
- Tests removing Tough feat
- Verifies HP is restored to original value

**TestApplyFeatBenefits_Mobile**
- Tests Mobile feat (+10 speed)
- Verifies speed increase
- Verifies benefit tracking

### 4. Multiple Feats Tests
**TestMultipleFeats**
- Tests adding multiple feats (Actor + Athlete)
- Verifies both feats apply correctly
- Verifies independent tracking
- Tests removing one feat while keeping the other
- Verifies correct benefit isolation

### 5. Helper Function Tests
**TestHasAbilityChoice**
- Tests detecting feats with ability choices
- Tests feats without ability choices
- Tests feats with no ability increases

**TestGetAbilityChoices**
- Tests retrieving available choices
- Verifies correct choice list

### 6. Edge Cases
**TestAbilityScoreMax**
- Tests ability scores cap at 20
- Verifies applying feat to maxed ability doesn't exceed cap

**TestAbilityScoreMin**
- Tests ability scores don't go below 1
- Verifies removing feat doesn't create negative scores

## Manual Testing Checklist

### Add Feat Flow (press 'f')

#### Test Case 1: Add Feat Without Choices
1. Launch app: `./lazydndplayer`
2. Navigate to Traits panel (Tab key)
3. Note current Charisma score
4. Press `f` to add feat
5. Navigate to "Actor" and press Enter
6. ✅ Verify: Charisma increased by 1
7. ✅ Verify: Message shows "Feat gained: Actor!"
8. ✅ Verify: "Actor" appears in feat list
9. Press `s` to save
10. Restart app and load character
11. ✅ Verify: Actor feat persists
12. ✅ Verify: Charisma still increased

**Expected Behavior:**
- Charisma: +1
- No ability selector shown
- Feat added to list immediately
- Character saved automatically

#### Test Case 2: Add Feat With Ability Choice
1. Navigate to Traits panel
2. Note current Strength and Dexterity
3. Press `f` to add feat
4. Navigate to "Athlete" and press Enter
5. ✅ Verify: Ability choice selector appears
6. ✅ Verify: Shows "Choose which ability to increase"
7. ✅ Verify: Lists "Strength" and "Dexterity"
8. Select "Strength" and press Enter
9. ✅ Verify: Strength increased by 1
10. ✅ Verify: Dexterity unchanged
11. ✅ Verify: Message shows "Feat gained: Athlete (+1 Strength)!"
12. ✅ Verify: "Athlete" in feat list
13. Press `s` to save
14. Restart and verify persistence

**Expected Behavior:**
- Ability selector appears
- Only chosen ability increases
- Feat added with choice recorded
- Character saved automatically

#### Test Case 3: Add Feat With Special Effects (Tough)
1. Navigate to Traits panel
2. Note current MaxHP and Level
3. Press `f` to add feat
4. Navigate to "Tough" and press Enter
5. ✅ Verify: MaxHP increased by (Level × 2)
6. ✅ Verify: CurrentHP also updated
7. ✅ Verify: Feat added to list

**Expected Behavior:**
- HP: +2 per character level
- Current HP updated to max
- Benefit tracked

#### Test Case 4: Add Feat With Special Effects (Mobile)
1. Navigate to Traits panel
2. Note current Speed
3. Press `f` to add feat
4. Navigate to "Mobile" and press Enter
5. ✅ Verify: Speed increased by 10
6. ✅ Verify: Feat added to list

**Expected Behavior:**
- Speed: +10 feet
- Benefit tracked

#### Test Case 5: Try to Add Duplicate Non-Repeatable Feat
1. Add "Actor" feat (if not already added)
2. Press `f` to add feat again
3. Navigate to "Actor" and press Enter
4. ✅ Verify: Error message: "You already have Actor and it's not repeatable"
5. ✅ Verify: Feat not added again
6. ✅ Verify: No duplicate in feat list

**Expected Behavior:**
- Error message shown
- Feat not duplicated
- No changes to character

#### Test Case 6: Cancel Feat Selection
1. Press `f` to add feat
2. Navigate to any feat
3. Press `Esc`
4. ✅ Verify: Selector closes
5. ✅ Verify: Message shows "Feat selection cancelled"
6. ✅ Verify: No feat added

**Expected Behavior:**
- Selector closes
- No changes to character

#### Test Case 7: Cancel Ability Choice
1. Press `f` to add feat
2. Select "Athlete" and press Enter
3. Ability selector appears
4. Press `Esc`
5. ✅ Verify: Selector closes
6. ✅ Verify: Feat NOT in feat list
7. ✅ Verify: No ability increases
8. ✅ Verify: Message shows "Feat selection cancelled"

**Expected Behavior:**
- Feat removed from list
- No ability increases
- Character saved

### Remove Feat Flow (press 'F')

#### Test Case 8: Remove Feat Without Choices
1. Add "Actor" feat if not present
2. Note current Charisma
3. Press `F` (Shift+F) to remove feat
4. ✅ Verify: Only owned feats shown
5. Navigate to "Actor" and press Enter
6. ✅ Verify: Charisma decreased by 1 (back to original)
7. ✅ Verify: Message shows "Feat removed: Actor (benefits reversed)"
8. ✅ Verify: "Actor" removed from feat list
9. Press `s` to save
10. Restart and verify removal persists

**Expected Behavior:**
- Charisma: -1 (restored)
- Feat removed from list
- Character saved automatically

#### Test Case 9: Remove Feat With Ability Choice
1. Add "Athlete" feat with Strength choice
2. Note Strength value
3. Press `F` to remove feat
4. Navigate to "Athlete" and press Enter
5. ✅ Verify: Strength decreased by 1 (back to original)
6. ✅ Verify: Feat removed from list
7. ✅ Verify: Message shows benefits reversed

**Expected Behavior:**
- Chosen ability restored
- Other abilities unchanged
- Feat removed completely

#### Test Case 10: Remove Feat With Special Effects (Tough)
1. Add "Tough" feat
2. Note MaxHP (should be base + Level × 2)
3. Press `F` to remove feat
4. Navigate to "Tough" and press Enter
5. ✅ Verify: MaxHP decreased by (Level × 2)
6. ✅ Verify: If CurrentHP > new MaxHP, it's adjusted

**Expected Behavior:**
- HP reduced correctly
- Current HP capped at new max

#### Test Case 11: Remove Feat With Special Effects (Mobile)
1. Add "Mobile" feat
2. Note Speed (should be base + 10)
3. Press `F` to remove feat
4. Navigate to "Mobile" and press Enter
5. ✅ Verify: Speed decreased by 10

**Expected Behavior:**
- Speed restored to original
- Benefit removed

#### Test Case 12: Try to Remove When No Feats
1. Remove all feats
2. Press `F` to remove feat
3. ✅ Verify: Message shows "No feats to remove"
4. ✅ Verify: No selector appears

**Expected Behavior:**
- Error message shown
- No selector

#### Test Case 13: Cancel Feat Removal
1. Have at least one feat
2. Press `F` to remove feat
3. Press `Esc`
4. ✅ Verify: Selector closes
5. ✅ Verify: Message shows "Feat selection cancelled"
6. ✅ Verify: No feat removed

**Expected Behavior:**
- Selector closes
- No changes to character

### Multi-Feat Tests

#### Test Case 14: Multiple Feats Independence
1. Add "Actor" feat (+1 Charisma)
2. Add "Athlete" feat (+1 Strength)
3. ✅ Verify: Both bonuses applied
4. Remove "Actor"
5. ✅ Verify: Charisma restored
6. ✅ Verify: Strength still has Athlete bonus
7. ✅ Verify: Athlete still in feat list

**Expected Behavior:**
- Feats tracked independently
- Removing one doesn't affect others
- Benefits properly isolated

#### Test Case 15: Same Ability from Multiple Sources
1. Add "Actor" feat (+1 Charisma)
2. Note Charisma (original + 1)
3. Add another feat that increases Charisma
4. ✅ Verify: Charisma increases again
5. Remove first feat
6. ✅ Verify: Only that feat's bonus removed
7. ✅ Verify: Other bonuses remain

**Expected Behavior:**
- Multiple bonuses stack
- Removal is source-specific
- No interference between sources

### Save/Load Tests

#### Test Case 16: Save Character With Feats
1. Add multiple feats with various benefits
2. Press `s` to save
3. ✅ Verify: Message shows "Character saved!"
4. Check JSON file contains:
   - Feats list
   - benefit_tracker section
   - Correct benefit sources and values

**Expected Behavior:**
- All feats persisted
- BenefitTracker serialized
- Character state complete

#### Test Case 17: Load Character With Feats
1. Restart app
2. Load saved character
3. ✅ Verify: All feats in list
4. ✅ Verify: All ability bonuses applied
5. ✅ Verify: BenefitTracker loaded
6. Remove a feat
7. ✅ Verify: Removal works correctly

**Expected Behavior:**
- Character fully restored
- Benefits tracked correctly
- System functional after load

## Integration Test Scenarios

### Scenario 1: Full Character Creation Flow
1. Create new character
2. Add 3 different feats:
   - One with single ability (Actor)
   - One with choice (Athlete)
   - One with special effect (Tough)
3. Verify all bonuses applied
4. Save character
5. Restart and load
6. Verify everything persists
7. Remove all feats one by one
8. Verify complete restoration

### Scenario 2: Error Recovery
1. Add feat with ability choice
2. Press Esc during ability selection
3. Verify feat not added
4. Add feat again
5. Complete selection this time
6. Verify works correctly

### Scenario 3: Edge Cases
1. Max out an ability to 20
2. Try to add feat that increases it
3. Verify caps at 20
4. Remove feat
5. Verify restoration works

## Test Execution

### Run All Tests
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go test ./internal/models -v
```

### Run Specific Test
```bash
go test ./internal/models -v -run TestApplyFeatBenefits_SingleAbility
```

### Run With Coverage
```bash
go test ./internal/models -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Success Criteria

### Unit Tests
- ✅ All tests pass
- ✅ No panics or errors
- ✅ 100% coverage of benefit application/removal logic

### Manual Tests
- ✅ All 17 test cases pass
- ✅ No data corruption
- ✅ Benefits properly tracked and reversed
- ✅ Save/load works correctly
- ✅ UI responsive and shows correct messages

### Integration Tests
- ✅ All 3 scenarios complete successfully
- ✅ No edge case failures
- ✅ System remains stable

## Known Issues
None currently identified.

## Future Test Coverage
- Class feature benefits
- Species benefits with choices
- Magic item benefits
- Temporary buffs/debuffs
- Multi-user concurrent testing
