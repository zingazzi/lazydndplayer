# Debug Instructions for Class Selector Issue

## Changes Made
1. Added Update() method to ClassSelector component
2. Updated handleClassSelectorKeys to use Update() method
3. Added comprehensive debug logging to track the issue

## How to Test

### Step 1: Rebuild
```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go build -o lazydndplayer .
```

### Step 2: Clear old logs and run with debug
```bash
rm -f lazydnd_debug.log
./lazydndplayer --debug
```

### Step 3: Reproduce the issue
1. Press 'p' to focus Character Info panel
2. Press 'c' to open class selector
3. Press down arrow 2-3 times
4. Press Enter
5. Press Esc to exit
6. Quit the app

### Step 4: Share the output

The debug output will now show:
- How many classes were loaded: `DEBUG: ClassSelector.Show() - First class mode, loaded X classes`
- When you navigate: `DEBUG: ClassSelector.Next() - selectedIndex now: X (total: Y)`
- When you try to select: `DEBUG: GetSelectedClass() - selectedIndex=X, classes.len=Y`

**Please share:**
1. The console output (you should see DEBUG messages)
2. The updated `lazydnd_debug.log` file

This will tell us exactly why GetSelectedClass() is returning empty!

## What We're Looking For

The debug will tell us if:
- ✅ Classes are being loaded (should be 12)
- ✅ selectedIndex is updating when you press down
- ✅ GetSelectedClass() can access the classes array
- ❌ Or if something else is wrong

## Expected Output Example
```
DEBUG: ClassSelector.Show() - First class mode, loaded 12 classes
DEBUG: ClassSelector - First class: Barbarian, Last class: Wizard
DEBUG: ClassSelector.Next() - selectedIndex now: 1 (total: 12)
DEBUG: ClassSelector.Next() - selectedIndex now: 2 (total: 12)
DEBUG: GetSelectedClass() - selectedIndex=2, classes.len=12
DEBUG: GetSelectedClass() - Returning: Cleric
```

If you see "classes.len=0" anywhere, that's the problem!
