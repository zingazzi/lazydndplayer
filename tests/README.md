# Tests Directory

This directory contains all test files for the lazydndplayer project.

## Structure

```
tests/
└── models/
    ├── feats_test.go       # Feat benefits application/removal tests
    └── feats_load_test.go  # Feat data loading tests
```

## Running Tests

### Run All Tests
```bash
go test ./tests/models/... -v
```

### Run Specific Test
```bash
go test ./tests/models -v -run TestApplyFeatBenefits_SingleAbility
```

### Run With Coverage
```bash
go test ./tests/models -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

### Feat System Tests (`feats_test.go`)
- ✅ **TestApplyFeatBenefits_SingleAbility** - Tests feats with single fixed ability (Actor)
- ✅ **TestApplyFeatBenefits_MultipleChoice** - Tests feats with ability choices (Athlete)
- ✅ **TestRemoveFeatBenefits_SingleAbility** - Tests benefit removal for single ability
- ✅ **TestRemoveFeatBenefits_MultipleChoice** - Tests benefit removal for choices
- ✅ **TestApplyFeatBenefits_Tough** - Tests special HP bonus feat
- ✅ **TestRemoveFeatBenefits_Tough** - Tests HP bonus removal
- ✅ **TestApplyFeatBenefits_Mobile** - Tests speed bonus feat
- ✅ **TestMultipleFeats** - Tests adding/removing multiple feats
- ✅ **TestHasAbilityChoice** - Tests ability choice detection
- ✅ **TestGetAbilityChoices** - Tests retrieving ability choices
- ✅ **TestAbilityScoreMax** - Tests ability score cap at 20
- ✅ **TestAbilityScoreMin** - Tests ability score floor at 1

### Feat Loading Tests (`feats_load_test.go`)
- ✅ **TestLoadAthleteFeat** - Verifies Athlete feat loads with correct choices
- ✅ **TestLoadActorFeat** - Verifies Actor feat loads with fixed ability

## Test Package Structure

Tests use the `models_test` package (black-box testing) to ensure they test only the public API of the models package. This follows Go testing best practices.

## Adding New Tests

1. Create test file in appropriate subdirectory under `tests/`
2. Use package name `<package>_test` (e.g., `models_test`)
3. Import the package under test: `"github.com/marcozingoni/lazydndplayer/internal/models"`
4. Follow existing naming conventions: `Test<FunctionName>_<Scenario>`
5. Update this README with test description

## CI/CD Integration

These tests should be run as part of the CI/CD pipeline before any deployment.
