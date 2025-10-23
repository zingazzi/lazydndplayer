# Bug Fix: Class Selector Black Screen

## Issue
User reported that pressing 'c' in the Character Info panel caused the screen to become black.

## Root Cause Investigation
The class selector was showing a black screen when no classes were loaded. This could happen if:
1. Classes failed to load from JSON files
2. The `cachedClasses` variable was nil
3. The loading error was silently swallowed

## Changes Made

### 1. Added Debug Logging
**File**: `internal/models/classes.go`
- Added debug logs to `GetAllClasses()` to trace class loading
- Logs when classes are being loaded
- Logs the count of successfully loaded classes
- Logs errors during loading

### 2. Added Warning in Class Selector
**File**: `internal/ui/components/classselector.go`
- Added warning message when no classes are loaded
- Helps identify if the issue is with loading or displaying

### 3. Verification
- Confirmed 12 class JSON files exist in `data/classes/`
- Classes: Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard

## How to Test

### 1. Run with Debug Mode
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
./lazydndplayer --debug
```

### 2. Test the Class Selector
1. Start the application
2. Navigate to Character Info panel (use 'p' to cycle focus)
3. Press 'c' to open class selector
4. Check if classes appear
5. If screen is black, check the `lazydnd_debug.log` file

### 3. Check Debug Log
```bash
cat ~/lazydnd_debug.log | grep -A 5 "GetAllClasses"
```

Look for:
- "GetAllClasses: cachedClasses is nil, loading..."
- "Successfully loaded X classes"
- Any error messages

## Expected Behavior

When you press 'c' in Character Info:
- A popup should appear with a list of classes on the left
- Class details should appear on the right
- You should see all 12 classes
- You can navigate with up/down arrows
- Press Enter to select, Esc to cancel

## If Still Black Screen

Check these things:

### 1. Working Directory
The application must be run from the project root where `data/classes/` exists.

```bash
pwd  # Should be: /Users/marcozingoni/Playgound/lazydndplayer
ls data/classes/*.json | wc -l  # Should show: 12
```

### 2. File Permissions
```bash
ls -la data/classes/*.json
# All files should be readable
```

### 3. JSON Syntax
Test if JSON files are valid:
```bash
for file in data/classes/*.json; do
  echo "Checking $file..."
  python3 -m json.tool "$file" > /dev/null && echo "✓ Valid" || echo "✗ Invalid"
done
```

### 4. Check Console Output
When running with `--debug`, any errors should appear in the console AND in the log file.

## Additional Improvements

The debug logging will help identify:
- If classes are being loaded at all
- How many classes are loaded
- Any parsing errors from JSON files
- Timing of when classes are requested

## Next Steps If Issue Persists

1. Share the content of `lazydnd_debug.log`
2. Share the console output when running with `--debug`
3. Verify the working directory is correct
4. Check if any JSON files have syntax errors

## Files Modified

1. `/Users/marcozingoni/Playgound/lazydndplayer/internal/models/classes.go`
   - Added debug logging to `GetAllClasses()`

2. `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/classselector.go`
   - Added warning when no classes loaded in `Show()`

