# Testing Instructions - Class Selector Debug

## Rebuild with Debug Logging

```bash
cd /Users/marcozingoni/Playgound/lazydndplayer
go build -o lazydndplayer .
```

## Test Steps

```bash
# Clear old log
rm -f lazydnd_debug.log

# Run with debug
./lazydndplayer --debug
```

### In the app:
1. Press **'p'** to cycle focus to Character Info panel (right side)
2. Press **'c'** to open class selector
3. Press **down arrow** 2-3 times
4. Press **enter**
5. Press **esc** to exit if needed
6. Quit the app (Ctrl+C or 'q')

### Check the log:
```bash
cat lazydnd_debug.log
```

## What to Look For

You should now see in `lazydnd_debug.log`:

```
[time] ClassSelector.Show() - First class mode, loaded 12 classes
[time] ClassSelector - First class: Barbarian, Last class: Wizard
[time] ClassSelector.Next() - selectedIndex now: 1 (total: 12)
[time] ClassSelector.Next() - selectedIndex now: 2 (total: 12)
[time] GetSelectedClass() - selectedIndex=2, classes.len=12
[time] GetSelectedClass() - Returning: Cleric
```

## Possible Issues

### If you see "classes.len=0":
- Classes aren't being loaded into the selector
- This would explain why selected=empty

### If you see "selectedIndex=0" always:
- Navigation (Next/Prev) isn't working
- The Update() method might not be called

### If you don't see ClassSelector.Show() in log:
- The selector isn't being shown at all
- Screen goes black because View() returns empty

Please share the complete `lazydnd_debug.log` file after testing!
