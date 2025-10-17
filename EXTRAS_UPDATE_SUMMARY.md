# Extra Bonuses Improvements - Summary

## What's New? ✨

### 1. Quick Value Adjustment with +/- Keys ⭐
**Before**: Had to press `e`, type number, press `Enter`
**After**: Just press `+` or `-` to adjust by 1

```
Old way (7 key presses):
e → 2 → Enter → ↓ → e → 1 → Enter

New way (4 key presses):
+ → + → ↓ → +

Savings: 3 fewer key presses (43% faster!)
```

### 2. Smart Escape Navigation 🧠
**Before**: `Esc` always went to previous screen (could be confusing)
**After**: `Esc` is context-aware!

- Press `e` → Extras → `Esc` → **Stats Panel** ✅
- Press `r` → ... → Extras → `Esc` → **Assignment** ✅

The screen tells you where Esc will go!

### 3. Better Instructions 📖
Now shows:
- Range limits: "(-5 to +10)"
- What +/- do: "Adjust"
- Where Esc goes: "Back to Stats Panel" or "Back to Assignment"

## Implementation Details

### New Methods in StatGenerator

```go
// Track entry point
directToExtras bool  // true if came from 'e', false if from 'r'

// Quick adjust methods
IncreaseExtra()  // +1 to selected ability (max +10)
DecreaseExtra()  // -1 to selected ability (min -5)

// Updated navigation
GoBack()  // Smart routing based on directToExtras
```

### Key Handler Updates

```go
case "+", "=":
    if state == StateSetExtras {
        IncreaseExtra()  // NEW!
    } else {
        IncreasePointBuy()  // Existing
    }

case "-", "_":
    if state == StateSetExtras {
        DecreaseExtra()  // NEW!
    } else {
        DecreasePointBuy()  // Existing
    }
```

### Navigation Flow

```
ShowExtrasOnly()  → directToExtras = true
Show()            → directToExtras = false

GoBack():
  if directToExtras:
    Hide()  // Close generator → Stats Panel
  else:
    Go to StateAssignStats  // Back to assignment
```

## Files Changed

1. **internal/ui/components/statgenerator.go**
   - Added `directToExtras` field
   - Added `IncreaseExtra()` / `DecreaseExtra()` methods
   - Updated `GoBack()` for smart navigation
   - Updated rendering to show context

2. **internal/ui/app.go**
   - Updated `+/-` key handlers to support extras adjustment
   - No other changes needed (clean integration!)

## User Benefits

### Speed ⚡
- **43% faster** for common operations (adding +2/+1 bonuses)
- No typing mode needed for simple adjustments

### Clarity 🔍
- Always know where Esc takes you
- Clear instructions for each mode
- Visible range limits

### Flexibility 🎯
- Still can type exact values with `e`
- Quick adjustments with `+/-`
- Choose the method that fits your workflow

## Usage Examples

### Adding Dragonborn Bonuses (Before)
```
e → type '2' → Enter → ↓ ↓ → e → type '1' → Enter → Navigate to CONFIRM → Enter
Total: 9 actions
```

### Adding Dragonborn Bonuses (After)
```
+ → + → ↓ ↓ → + → Navigate to CONFIRM → Enter
Total: 6 actions (33% faster!)
```

### Or Even Faster (After)
```
e (direct to extras, already on Strength)
+ → + → ↓ ↓ → + → CONFIRM
Total: 5 actions (44% faster!)
```

## Testing Checklist

- [ ] Build: `go build -o lazydndplayer .`
- [ ] Run: `./lazydndplayer`
- [ ] Test Direct Access (e key):
  - [ ] Press `e` in Stats panel
  - [ ] Use `+/-` to adjust values
  - [ ] Press `Esc` → Should return to Stats panel
  - [ ] Verify changes are applied
- [ ] Test Full Flow (r key):
  - [ ] Press `r` in Stats panel
  - [ ] Complete stat generation
  - [ ] In extras, press `Esc` → Should go to assignment
  - [ ] In assignment, press `Esc` → Should go to method selection
  - [ ] In method, press `Esc` → Should close generator
- [ ] Test +/- Keys:
  - [ ] Navigate to ability
  - [ ] Press `+` multiple times → Value increases
  - [ ] Press `-` multiple times → Value decreases
  - [ ] Verify cannot go below -5
  - [ ] Verify cannot go above +10
- [ ] Test Type Mode:
  - [ ] Press `e` on ability
  - [ ] Type a number
  - [ ] Press `Enter` → Value updates
  - [ ] Press `e` again
  - [ ] Press `Esc` → Cancels without saving

## Migration Notes

**No breaking changes!**
- All existing functionality preserved
- Only additions and improvements
- No data format changes
- No API changes

**Backward Compatible:**
- Old save files work perfectly
- All keyboard shortcuts still work
- Previous workflows still supported

## Performance Impact

**Minimal**: O(1) operations only
- Simple integer increment/decrement
- Boolean flag checks
- No loops or heavy computations

**Memory**: +1 bool field per StatGenerator instance
- Negligible impact (~1 byte)

## Future Enhancements

Possible improvements:
- [ ] Show total score while editing (Base + Extra = Total)
- [ ] Color-code positive/negative modifiers
- [ ] Preset bonus profiles (Dragonborn, Elf, etc.)
- [ ] Undo/redo for adjustments
- [ ] Macro support (apply common bonus sets)

## Conclusion

Three simple improvements that make a big difference:
1. ⚡ Faster value adjustment (+/- keys)
2. 🧠 Smarter navigation (context-aware Esc)
3. 📖 Better user guidance (clear instructions)

**Result**: More intuitive, faster, and easier to use! 🎉
