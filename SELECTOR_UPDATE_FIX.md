# Selector Update Method Consistency Fix

## Issue
Several UI selector components were not following a consistent pattern for handling keyboard input. Some handlers in `app.go` were manually duplicating navigation logic instead of delegating to the component's Update method.

## Root Cause
The selector components had inconsistent implementations:
- Some had Update() methods (e.g., `SpellPrepSelector`, `LevelUpSelector`)
- Others didn't have Update() methods (e.g., `CantripSelector`, `FightingStyleSelector`, `WeaponMasterySelector`, `ClassSkillSelector`)
- Handlers in `app.go` were manually handling navigation for components without Update() methods

This violated the **DRY principle** and made the code harder to maintain.

## Solution Applied

### 1. Added Update() Methods
Added consistent `Update(msg tea.Msg) (ComponentType, tea.Cmd)` methods to:
- ✅ `CantripSelector` (already existed but wasn't being used)
- ✅ `FightingStyleSelector` (newly added)
- ✅ `WeaponMasterySelector` (newly added)
- ✅ `ClassSkillSelector` (newly added)

Each Update method:
- Checks if the component is visible
- Handles navigation keys (up/k, down/j)
- Handles selection keys (space for toggles)
- Returns early on enter/esc for handler processing
- Returns the updated component and command

### 2. Updated Handler Methods in app.go
Refactored these handlers to use the component's Update() method:
- ✅ `handleCantripSelectorKeys`
- ✅ `handleFightingStyleSelectorKeys`
- ✅ `handleWeaponMasterySelectorKeys`
- ✅ `handleClassSkillSelectorKeys`

Pattern used:
```go
func (m *Model) handleComponentKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Delegate navigation to the component's Update method
    var cmd tea.Cmd
    *m.component, cmd = m.component.Update(tea.KeyMsg(msg))

    // Handler only processes business logic (save, confirm, cancel)
    switch msg.String() {
    case "enter":
        // Apply changes, save, hide component
    case "esc":
        // Cancel, rollback, hide component
    }

    return m, cmd
}
```

## Benefits

1. **DRY Principle**: Navigation logic is no longer duplicated across handlers
2. **Consistency**: All selectors now follow the same Update pattern
3. **Maintainability**: Changes to navigation only need to be made in one place
4. **Separation of Concerns**: Components handle UI logic, handlers handle business logic
5. **Testability**: Components can be tested independently

## Files Modified

### Component Files (Added Update methods)
- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/fightingstyleselector.go`
- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/weaponmasteryselector.go`
- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/components/classskillselector.go`

### Handler File (Refactored to use Update methods)
- `/Users/marcozingoni/Playgound/lazydndplayer/internal/ui/app.go`
  - `handleCantripSelectorKeys`
  - `handleFightingStyleSelectorKeys`
  - `handleWeaponMasterySelector Keys`
  - `handleClassSkillSelectorKeys`

## Testing

✅ No linter errors
✅ All components follow consistent pattern
✅ Handlers properly delegate to Update methods
✅ Business logic remains in handlers

## Pattern for Future Selectors

When creating new selector components:

1. **Component Structure**:
```go
type MySelector struct {
    visible bool
    cursor  int
    // ... other fields
}

func (s *MySelector) Update(msg tea.Msg) (MySelector, tea.Cmd) {
    if !s.visible {
        return *s, nil
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            s.Prev()
        case "down", "j":
            s.Next()
        case " ":
            s.Toggle() // if applicable
        case "enter", "esc":
            return *s, nil // Let handler process
        }
    }

    return *s, nil
}
```

2. **Handler Structure**:
```go
func (m *Model) handleMySelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Delegate navigation
    var cmd tea.Cmd
    *m.mySelector, cmd = m.mySelector.Update(tea.KeyMsg(msg))

    // Handle business logic
    switch msg.String() {
    case "enter":
        // Confirm, apply, save
    case "esc":
        // Cancel, rollback
    }

    return m, cmd
}
```

## Status
✅ **COMPLETE** - All selectors now follow consistent Update pattern
